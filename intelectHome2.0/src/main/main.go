package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	database "github.com/Piccadilly98/goProjects/intelectHome2.0/src/dataBase"
	"github.com/Piccadilly98/goProjects/intelectHome2.0/src/dataBase/status_worker"
	"github.com/Piccadilly98/goProjects/intelectHome2.0/src/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal(err)
	}
	db, err := database.MakeDataBase(os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USERNAME"), os.Getenv("DB_NAME_DB"), os.Getenv("DB_PASSWORD"))
	if err != nil {
		log.Fatal(err.Error())
	}
	defer func() {
		db.Close()
		fmt.Println("Сворачиваемся")
	}()
	sw := status_worker.MakeStatusWorker(db, 30*time.Second, 150*time.Second)
	errWork := database.MakeErrorWorker(db)
	errWork.Start()
	sw.Start()
	upateChan := sw.UpdateChan()
	r := chi.NewRouter()
	registration := handlers.MakeRegistrationHandler(db)
	update := handlers.MakeBoardUpdateHandler(db)
	get := handlers.MakeBoardIDGet(db)
	updateInfo := handlers.MakeUpdateBoardInfoHandler(db, upateChan)
	getInfo := handlers.MakeBoardInfoGetHandler(db)
	boardsGet := handlers.MakeBoardsGetHandler(db)
	controllersGet := handlers.MakeBoardControllersHandler(db)
	controllersReg := handlers.MakeControllersRegistrationHandler(db)
	controllerUpdate := handlers.MakeControllerUpdateHandler(db)
	r.Patch("/boards/{board_id}/controllers/{controller_id}", controllerUpdate.Handler)
	r.Get("/boards/{board_id}/controllers", controllersGet.Handler)
	r.Post("/boards/{board_id}/controllers", controllersReg.Handler)
	r.Get("/boards", boardsGet.Handler)
	r.Get("/boards/{board_id}/info", getInfo.Handler)
	r.Patch("/boards/{board_id}/info", updateInfo.Handler)
	r.Get("/boards/{board_id}", get.Handler)
	r.Post("/boards", registration.RegistrationHandler)
	r.Patch("/boards/{board_id}", update.Handler)
	err = http.ListenAndServe("localhost:8080", r)
	if err != nil {
		panic(err)
	}
}
