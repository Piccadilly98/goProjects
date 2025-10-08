package main

import (
	"log"
	"net/http"

	"github.com/Piccadilly98/goProjects/intelectHome/src/handlers"
	"github.com/Piccadilly98/goProjects/intelectHome/src/storage"
)

func main() {
	st := storage.MakeStorage()
	control := handlers.MakeHandlerControl(st)
	board := handlers.MakeBoarsHandler(st)
	http.HandleFunc("/control", control.Contorol)
	http.HandleFunc("/sensor/data", board.BoardsHendler)
	err := http.ListenAndServe("localhost:9091", nil)
	if err != nil {
		log.Fatal(err)
	}
}
