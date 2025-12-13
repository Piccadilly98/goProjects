package database

import (
	"log"
)

type ErrorWorker struct {
	db *DataBase
}

// всю логику сюда

func MakeErrorWorker(db *DataBase) *ErrorWorker {
	return &ErrorWorker{db: db}
}

func (ew *ErrorWorker) Start() {
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
