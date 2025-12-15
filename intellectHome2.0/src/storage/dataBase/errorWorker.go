package database

import (
	"fmt"
	"log"
	"time"

	"github.com/Piccadilly98/goProjects/intellectHome2.0/src/core/events"
)

type ErrorWorker struct {
	db             *DataBase
	subErrrorDb    *events.TopicSubscriberOut
	dataBaseStatus *events.TopicSubscriberOut
}

// всю логику сюда

func MakeErrorWorker(db *DataBase) *ErrorWorker {
	return &ErrorWorker{
		db:             db,
		subErrrorDb:    db.eventBus.Subscribe(events.TopicErrorsDB, TopicPublisherNameErrorWorker),
		dataBaseStatus: db.eventBus.Subscribe(events.TopicDataBaseStatus, TopicPublisherNameErrorWorker),
	}
}

func (ew *ErrorWorker) Start() {
	go func() {
		for event := range ew.subErrrorDb.Chan {
			log.Printf("DB error detected: %v by: %s → starting recovery", event.Payload, event.Publisher)
			err := ew.db.eventBus.Publish(ew.dataBaseStatus.Topic, events.Event{
				Type:       ew.dataBaseStatus.Topic,
				Payload:    fmt.Sprintf("message received by: %s, message: %v,  DataBase fail, start Recover\n", event.Publisher, event.Payload),
				Publisher:  TopicPublisherNameErrorWorker,
				DatePublic: time.Now(),
			}, ew.dataBaseStatus.ID)
			if err != nil {
				log.Println(err)
			}
			if !ew.db.Recover() {
				log.Fatalf("FATAL: cannot recover database: %v", event.Payload)
				err := ew.db.eventBus.Publish(ew.dataBaseStatus.Topic, events.Event{
					Type:       ew.dataBaseStatus.Topic,
					Payload:    "Data base not recover, server off\n",
					Publisher:  TopicPublisherNameErrorWorker,
					DatePublic: time.Now(),
				}, ew.dataBaseStatus.ID)
				if err != nil {
					log.Println(err)
				}
			} else {
				log.Println("DB recovered successfully")
				err := ew.db.eventBus.Publish(ew.dataBaseStatus.Topic, events.Event{
					Type:       ew.dataBaseStatus.Topic,
					Payload:    "Data base recovered successfully\n",
					Publisher:  TopicPublisherNameErrorWorker,
					DatePublic: time.Now(),
				}, ew.dataBaseStatus.ID)
				if err != nil {
					log.Println(err)
				}
			}
		}
	}()
}
