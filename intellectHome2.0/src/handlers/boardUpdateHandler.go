package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	dto "github.com/Piccadilly98/goProjects/intellectHome2.0/src/DTO"
	database "github.com/Piccadilly98/goProjects/intellectHome2.0/src/dataBase"
	"github.com/go-chi/chi/v5"
)

const (
	urlPath = "board_id"
)

type boardUpdate struct {
	db *database.DataBase
}

func MakeBoardUpdateHandler(db *database.DataBase) *boardUpdate {
	return &boardUpdate{db: db}
}

func (bu *boardUpdate) Handler(w http.ResponseWriter, r *http.Request) {
	param := chi.URLParam(r, urlPath)
	w.Header().Set("Content-Type", "application/json")
	if !ProcessingURLParam(w, r, param, bu.db) {
		return
	}

	update := &dto.UpdateBoardDataDto{}
	if !bu.readBodyWriteHeader(w, r, update) {
		return
	}

	code, err := bu.db.UpdateBoard(r.Context(), param, update)
	if err != nil {
		log.Println(err.Error())
		if code == 0 {
			return
		}
		w.WriteHeader(code)
		errResponse := fmt.Sprintf(`{"status":"error", "text":"%s"}`, err.Error())
		w.Write([]byte(errResponse))
		return
	}
	w.Write([]byte(`{"status":"ok"}`))
}

func (bu *boardUpdate) readBodyWriteHeader(w http.ResponseWriter, r *http.Request, update *dto.UpdateBoardDataDto) bool {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"status":"error", "text":"empty body"}`))
		return false
	}
	err = json.Unmarshal(body, update)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"status":"error", "text":"invalid body"}`))
		return false
	}

	if !update.Validate() {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"status":"error", "text":"invalid body"}`))
		return false
	}
	return true
}
