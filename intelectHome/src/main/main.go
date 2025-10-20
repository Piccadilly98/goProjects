package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Piccadilly98/goProjects/intelectHome/src/auth"
	"github.com/Piccadilly98/goProjects/intelectHome/src/handlers"
	"github.com/Piccadilly98/goProjects/intelectHome/src/middleware"
	"github.com/Piccadilly98/goProjects/intelectHome/src/rate_limit"
	"github.com/Piccadilly98/goProjects/intelectHome/src/storage"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load("/Users/flowerma/Desktop/goProjects/intelectHome/.env")
	r := chi.NewRouter()
	st := storage.MakeStorage("ADMIN", "ESP32_1", "ESP32_2")
	sm := auth.MakeSessionManager()
	tw := &auth.TokenWorker{}
	middlewareAuth := auth.MiddlewareAuth(st, sm)
	control := handlers.MakeHandlerControl(st)
	boardsID := handlers.MakeBoarsIDHandler(st)
	boards := handlers.MakeBoarsHandler(st)
	devices := handlers.MakeDevicesHandler(st)
	devicesID := handlers.MakeDevicesIDHandler(st)
	logs := handlers.MakeLogsHandler(st)
	login := auth.MakeLoginHandlers(st, sm, tw)
	globalRateLimiter := rate_limit.MakeGlobalRateLimiter(50, 50)
	ipRl := rate_limit.MakeIpRateLimiter(2, 2)

	r.Use(middleware.GlobalRateLimiterToMiddleware(globalRateLimiter, st))
	r.Use(middleware.IpRateLimiter(ipRl, st))
	r.With(middlewareAuth).Route("/", func(r chi.Router) {
		r.Post("/control", control.Control)
		r.HandleFunc("/boards/{boardID}", boardsID.BoardsIDHandler)
		r.Get("/boards", boards.BoardsHandler)
		r.HandleFunc("/devices", devices.DevicesHandler)
		r.Get("/devices/{deviceID}", devicesID.DevicesIDHandler)
		r.Get("/logs", logs.LogsHandler)
		r.Post("/login", login.LoginHandler)
	})

	go func() {
		time.Sleep(15 * time.Minute)
		os.Exit(1)
	}()
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal(err)
	}
}
