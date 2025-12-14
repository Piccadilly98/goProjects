package events

import (
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

const (
	DefaultBufferSize         = 50
	DefaultIntervarCheckQueue = 10
)

type EventBus struct {
	mu            sync.RWMutex
	subscribers   map[string]map[int]*topicSubscriberIn
	bufferSize    int
	subscribersID int
	dq            *DroppedQueueProcessor
	isStartedDq   atomic.Bool
}

func NewEventBus(bufferSize int, intervarCheckQueue time.Duration) *EventBus {
	if bufferSize <= 0 {
		bufferSize = DefaultBufferSize
		log.Printf("Invalid buffer size in NewEventBus, set default buffer size: %d\n", DefaultBufferSize)
	}
	eb := &EventBus{
		subscribers: make(map[string]map[int]*topicSubscriberIn),
		bufferSize:  bufferSize,
		dq:          newDroppedQueueProcessor(intervarCheckQueue),
	}
	return eb
}

func (eb *EventBus) Subscribe(topic, name string) *TopicSubscriberOut {
	eb.mu.Lock()
	defer eb.mu.Unlock()
	eb.subscribersID++
	ch := make(chan Event, eb.bufferSize)
	ts := &TopicSubscriberOut{
		Topic: topic,
		Chan:  ch,
		ID:    eb.subscribersID,
	}

	if eb.subscribers[topic] == nil {
		eb.subscribers[topic] = make(map[int]*topicSubscriberIn)
	}
	in := &topicSubscriberIn{
		ch:       ch,
		name:     name,
		isClosed: false,
	}

	eb.subscribers[topic][ts.ID] = in

	log.Printf("Зарегистрирован новый подписчик: %s\n", name)
	return ts
}

func (eb *EventBus) Publish(topic string, event Event, id int) error {
	if event.DatePublic.IsZero() {
		event.DatePublic = time.Now()
	}
	eb.mu.RLock()
	defer eb.mu.RUnlock()
	subscribers, ok := eb.subscribers[topic]
	if !ok {
		return fmt.Errorf("invalid topic, please subscribe")
	}
	_, ok = eb.subscribers[topic][id]
	if !ok {
		return fmt.Errorf("invalid id, please subscribe")
	}

	go func() {
		for idMap, sub := range subscribers {
			if id == idMap {
				continue
			}
			select {
			case sub.ch <- event:
			default:
				log.Printf("event in topic '%s', for id: %d not sent, dropped in queue\n", topic, id)
				eb.dq.AddToQueue(topic, event, sub, id)
			}
		}
	}()
	if !eb.isStartedDq.Load() {
		go eb.dq.StartProcessingQueue()
		eb.isStartedDq.Store(true)
	}
	return nil
}

func (eb *EventBus) UnSubscribe(sub *TopicSubscriberOut) {
	eb.mu.Lock()
	defer func() {
		eb.mu.Unlock()
	}()
	log.Printf("запрос на удаление: %d\n", sub.ID)
	eb.subscribers[sub.Topic][sub.ID].isClosed = true
	eb.dq.DeleteFromQueue(sub.Topic, sub.ID)
	log.Printf("удалили из очереди %d\n", sub.ID)
	close(eb.subscribers[sub.Topic][sub.ID].ch)
}
