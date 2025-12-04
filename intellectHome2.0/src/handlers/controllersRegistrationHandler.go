package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	dto "github.com/Piccadilly98/goProjects/intellectHome2.0/src/DTO"
	database "github.com/Piccadilly98/goProjects/intellectHome2.0/src/dataBase"
	"github.com/go-chi/chi/v5"
)

const (
	TypeBinary = "binary"
	TypeSensor = "sensor"
)

type controllersRegistrationHandler struct {
	db *database.DataBase
}

func MakeControllersRegistrationHandler(db *database.DataBase) *controllersRegistrationHandler {
	return &controllersRegistrationHandler{db: db}
}

func (cr *controllersRegistrationHandler) Handler(w http.ResponseWriter, r *http.Request) {
	boardID := chi.URLParam(r, "board_id")
	w.Header().Set("Content-Type", "application/json")
	if !ProcessingURLParam(w, r, boardID, cr.db) {
		return
	}
	dto := &dto.RegistrationController{}
	if !cr.readBodyAndGetDTO(w, r, dto) {
		return
	}
	ok, isBinary, isSensor, json := cr.processingDTOandGetJson(w, r, dto)
	if !ok {
		return
	}
	code, err := cr.db.RegistrationController(r.Context(), json, boardID, isBinary, isSensor)
	if err != nil {
		if code == 0 {
			return
		}
		w.WriteHeader(code)
		errResponse := fmt.Sprintf(`{"status":"error", "text":"%s"}`, err.Error())
		w.Write([]byte(errResponse))
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(json)
}

func (cr *controllersRegistrationHandler) readBodyAndGetDTO(w http.ResponseWriter, r *http.Request, dto *dto.RegistrationController) bool {
	err := json.NewDecoder(r.Body).Decode(dto)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"status":"error", "text":"invalid body"}`))
		return false
	}
	if !dto.Validate() {
		log.Println("non valid")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"status":"error", "text":"invalid body"}`))
		return false
	}
	return true
}

func (cr *controllersRegistrationHandler) processingDTOandGetJson(w http.ResponseWriter, r *http.Request, dto *dto.RegistrationController) (bool, bool, bool, []byte) {
	exist, code, err := cr.db.GetExistWithControllerId(r.Context(), dto.ControllerID)
	var sensor, binary bool
	if err != nil {
		if code == 0 {
			return false, false, false, nil
		}
		w.WriteHeader(code)
		errResponse := fmt.Sprintf(`{"status":"error", "text":"%s"}`, err.Error())
		w.Write([]byte(errResponse))
		return false, false, false, nil
	}
	if exist {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(`{"status":"error", "text":"controller_id %s contains in board"}`, dto.ControllerID)))
		return false, false, false, nil
	}
	var b []byte
	if dto.ControllerType == TypeSensor {
		sensor = true
		res := dto.ToSensorController()
		if !res.Validate() {
			log.Println("non valid")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"status":"error", "text":"invalid body"}`))
			return false, false, false, nil
		}
		b, err = json.Marshal(res)
	} else if dto.ControllerType == TypeBinary {
		binary = true
		res := dto.ToBinaryController()
		if !res.Validate() {
			log.Println("non valid")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"status":"error", "text":"invalid body"}`))
			return false, false, false, nil
		}
		b, err = json.Marshal(res)
	}

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`{"status":"error", "text":"%s"}`, err.Error())))
		return false, false, false, nil
	}

	return true, binary, sensor, b
}
