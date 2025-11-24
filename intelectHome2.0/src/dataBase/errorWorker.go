package database

import (
	"log"
)

type errorWorker struct {
	db *DataBase
}

func MakeErrorWorker(db *DataBase) *errorWorker {
	return &errorWorker{db: db}
}

func (ew *errorWorker) Start() {
	go func() {
		for err := range ew.db.ErrChan() {
			log.Printf("DB error detected: %v → starting recovery", err)
			if !ew.db.Recover() {
				log.Fatalf("FATAL: cannot recover database: %v", err)
			} else {
				log.Println("DB recovered successfully")
			}
		}
	}()
}
