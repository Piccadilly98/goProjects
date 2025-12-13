package server

import (
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	database "github.com/Piccadilly98/goProjects/intellectHome2.0/src/dataBase"
	"github.com/Piccadilly98/goProjects/intellectHome2.0/src/dataBase/status_worker"
	"github.com/Piccadilly98/goProjects/intellectHome2.0/src/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

type Server struct {
	Db              *database.DataBase
	StatusWorker    *status_worker.StatusWorker
	ErrorWorker     *database.ErrorWorker
	R               *chi.Mux
	UpdateChan      chan string
	Wg              sync.WaitGroup
	ErrorServerChan chan error
}

func NewServer(testing bool, intervalUpdateStatus time.Duration, timeForStatusOffline time.Duration, workers bool) (*Server, error) {
	err := loadConfig("/Users/flowerma/Desktop/goProjects/intellectHome2.0/src/main/.env")
	if err != nil {
		return nil, err
	}
	serv := &Server{}
	serv.ErrorServerChan = make(chan error)
	serv.R = chi.NewMux()
	var db *database.DataBase
	if !testing {
		db, err = database.MakeDataBase(os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USERNAME"), os.Getenv("DB_NAME_DB"), os.Getenv("DB_PASSWORD"))
		if err != nil {
			return nil, err
		}
	} else {
		db, err = database.MakeDataBase(os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USERNAME"), os.Getenv("DB_NAME_TEST"), os.Getenv("DB_PASSWORD"))
		if err != nil {
			return nil, err
		}
	}
	serv.Db = db
	serv.ErrorWorker = database.MakeErrorWorker(serv.Db)
	serv.StatusWorker = status_worker.MakeStatusWorker(serv.Db, intervalUpdateStatus, timeForStatusOffline)
	if workers {
		serv.ErrorWorker.Start()
		serv.StatusWorker.Start()
	}
	serv.UpdateChan = serv.StatusWorker.UpdateChan()
	registration := handlers.MakeRegistrationHandler(db)
	update := handlers.MakeBoardUpdateHandler(db)
	get := handlers.MakeBoardIDGet(db)
	updateInfo := handlers.MakeUpdateBoardInfoHandler(db, serv.UpdateChan)
	getInfo := handlers.MakeBoardInfoGetHandler(db)
	boardsGet := handlers.MakeBoardsGetHandler(db)
	controllersGet := handlers.MakeBoardControllersHandler(db)
	controllersReg := handlers.MakeControllersRegistrationHandler(db)
	controllerUpdate := handlers.MakeControllerUpdateHandler(db)
	serv.R.Post("/boards", registration.RegistrationHandler)
	serv.R.Get("/boards", boardsGet.Handler)
	serv.R.Get("/boards/{board_id}", get.Handler)
	serv.R.Patch("/boards/{board_id}", update.Handler)
	serv.R.Get("/boards/{board_id}/info", getInfo.Handler)
	serv.R.Patch("/boards/{board_id}/info", updateInfo.Handler)
	serv.R.Get("/boards/{board_id}/controllers", controllersGet.Handler)
	serv.R.Post("/boards/{board_id}/controllers", controllersReg.Handler)
	serv.R.Patch("/boards/{board_id}/controllers/{controller_id}", controllerUpdate.Handler)
	return serv, nil
}

func loadConfig(pathToConfig string) error {
	err := godotenv.Load(pathToConfig)
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) Start(addr string) {
	isErr := false
	go func() {
		err := http.ListenAndServe(addr, s.R)
		if err != nil {
			isErr = true
			s.ErrorServerChan <- err
		}
	}()
	time.Sleep(3 * time.Second)
	if !isErr {
		log.Printf("Server start in: %s\n", addr)
	}
}
