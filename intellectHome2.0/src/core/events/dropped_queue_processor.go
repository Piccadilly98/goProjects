package events

import (
	"log"
	"sync"
	"time"
)

type DroppedQueueProcessor struct {
	queue              map[string]*topicQueueDropping
	mu                 sync.RWMutex
	intervalCheckQueue time.Duration
}

func newDroppedQueueProcessor(intervalCheckQueue time.Duration) *DroppedQueueProcessor {
	if intervalCheckQueue <= 0 {
		intervalCheckQueue = DefaultIntervarCheckQueue * time.Second
	}
	return &DroppedQueueProcessor{
		queue:              make(map[string]*topicQueueDropping),
		intervalCheckQueue: intervalCheckQueue,
	}
}

func (dq *DroppedQueueProcessor) AddToQueue(topic string, event Event, sub *topicSubscriberIn, id int) {
	dq.mu.Lock()
	defer dq.mu.Unlock()
	if dq.queue[topic] == nil {
		dq.queue[topic] = &topicQueueDropping{
			event:       event,
			subscribers: make(map[int]*topicSubscriberIn),
		}
	}
	dq.queue[topic].subscribers[id] = sub
}

func (dq *DroppedQueueProcessor) DeleteFromQueue(topic string, id int) {
	dq.mu.Lock()
	defer dq.mu.Unlock()
	if dq.queue[topic] == nil {
		return
	}
	delete(dq.queue[topic].subscribers, id)
	if len(dq.queue[topic].subscribers) == 0 {
		log.Printf("удалили топик: %s - %d", topic, id)
		delete(dq.queue, topic)
	}
}

func (dq *DroppedQueueProcessor) StartProcessingQueue() {
	time.Sleep(dq.intervalCheckQueue)
	ticker := time.NewTicker(dq.intervalCheckQueue)
	for range ticker.C {
		dq.mu.Lock()
		for topic, queue := range dq.queue {
			for id, sub := range queue.subscribers {
				if sub.isClosed {
					log.Printf("Канал закрыт, удаляем\n")
					delete(queue.subscribers, id)
					continue
				} else {
					select {
					case sub.ch <- queue.event:
						log.Printf("event in topic: '%s', for id: %d sent\n", topic, id)
						delete(queue.subscribers, id)
						continue
					default:
						log.Printf("event in topic: '%s', for id: %d not sent\n", topic, id)
						delete(queue.subscribers, id)
						continue
					}
				}
			}
			if len(queue.subscribers) == 0 {
				delete(dq.queue, topic)
			}
		}
		dq.mu.Unlock()
	}
}
