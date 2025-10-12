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
	godotenv.Load("/Users/flowerma/Desktop/goProjects/intelectHome/.env")
	// worker := auth.TokenWorker{}
	// str, _ := worker.CreateToken("123", "ADMIN", 24*time.Hour)
	// fmt.Println(str)
	// os.Exit(1)
	r := chi.NewRouter()

	st := storage.MakeStorage("ADMIN", "ESP32_1")

	middleware := auth.MiddlewareAuth(st)
	control := handlers.MakeHandlerControl(st)
	boardsID := handlers.MakeBoarsIDHandler(st)
	boards := handlers.MakeBoarsHandler(st)
	devices := handlers.MakeDevicesHandler(st)
	devicesID := handlers.MakeDevicesIDHandler(st)
	logs := handlers.MakeLogsHandler(st)
	login := auth.MakeLoginHandlers(st)

	r.With(middleware).Route("/", func(r chi.Router) {
		r.Post("/control", control.Contorol)
		r.HandleFunc("/boards/{boardID}", boardsID.BoardsIDHandler)
		r.Get("/boards", boards.BoardsHandler)
		r.HandleFunc("/devices", devices.DevicesHandler)
		r.HandleFunc("/devices/{deviceID}", devicesID.DevicesIDHandler)
		r.Get("/logs", logs.LogsHandler)
	})
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
