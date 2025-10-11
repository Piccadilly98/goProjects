package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Piccadilly98/goProjects/intelectHome/src/auth"
	"github.com/Piccadilly98/goProjects/intelectHome/src/handlers"
	"github.com/Piccadilly98/goProjects/intelectHome/src/storage"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load("/Users/flowerma/Desktop/goProjects/intelectHome/src/.env")

	r := chi.NewRouter()
	st := storage.MakeStorage()
	control := handlers.MakeHandlerControl(st)
	boardsID := handlers.MakeBoarsIDHandler(st)
	boards := handlers.MakeBoarsHandler(st)
	devices := handlers.MakeDevicesHandler(st)
	devicesID := handlers.MakeDevicesIDHandler(st)
	logs := handlers.MakeLogsHandler(st)
	login := auth.MakeLoginHandlers(st)
	r.HandleFunc("/control", control.Contorol)
	r.HandleFunc("/boards/{boardID}", boardsID.BoardsIDHandler)
	r.HandleFunc("/boards", boards.BoardsHandler)
	r.HandleFunc("/devices", devices.DevicesHandler)
	r.HandleFunc("/devices/{deviceID}", devicesID.DevicesIDHandler)
	r.HandleFunc("/logs", logs.LogsHandler)
	r.HandleFunc("/login", login.LoginHandler)
	go func() {
		time.Sleep(15 * time.Minute)
		os.Exit(1)
	}()
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal(err)
	}
}
