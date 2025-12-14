package events

import "time"

type Event struct {
	Type         string
	BoardID      string
	ControllerID *string
	Payload      any
	Description  *string
	DatePublic   time.Time
	Publisher    string
}

type TopicSubscriberOut struct {
	Topic string
	Chan  <-chan Event
	ID    int
}

type topicSubscriberIn struct {
	isClosed bool
	ch       chan Event
	name     string
}

type topicQueueDropping struct {
	event       Event
	subscribers map[int]*topicSubscriberIn
}
