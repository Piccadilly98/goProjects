package main

import (
	"log"
	"net/http"
	"os"
	"time"

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
	logs := handlers.MakeLogsHandler(st)
	r.HandleFunc("/control", control.Contorol)
	r.HandleFunc("/boards/{boardID}", boardsID.BoardsIDHandler)
	r.HandleFunc("/boards", boards.BoardsHandler)
	r.HandleFunc("/devices", devices.DevicesHandler)
	r.HandleFunc("/devices/{deviceID}", devicesID.DevicesIDHandler)
	r.HandleFunc("/logs", logs.LogsHandler)
	go func() {
		time.Sleep(5 * time.Minute)
		os.Exit(1)
	}()
	err := http.ListenAndServe("localhost:8080", r)
	if err != nil {
		log.Fatal(err)
	}
}
