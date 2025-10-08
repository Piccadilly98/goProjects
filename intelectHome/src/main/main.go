package main

import (
	"log"
	"net/http"

	"github.com/Piccadilly98/goProjects/intelectHome/src/handlers"
	"github.com/Piccadilly98/goProjects/intelectHome/src/storage"
	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()

	st := storage.MakeStorage()
	control := handlers.MakeHandlerControl(st)
	boardsID := handlers.MakeBoarsIDHandler(st)
	boards := handlers.MakeBoarsHandler(st)
	devices := handlers.MakeDevicesHandler(st)
	devicesID := handlers.MakeDevicesIDHandler(st)
	r.HandleFunc("/control", control.Contorol)
	r.HandleFunc("/boards/{boardID}", boardsID.BoardsIDHandler)
	r.HandleFunc("/boards", boards.BoardsHandler)
	r.HandleFunc("/devices", devices.DevicesHandler)
	r.HandleFunc("/devices/{deviceID}", devicesID.DevicesIDHandler)
	err := http.ListenAndServe("localhost:8080", r)
	if err != nil {
		log.Fatal(err)
	}
}
