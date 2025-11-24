package handlers

import (
	"encoding/json"
	"net/http"

	database "github.com/Piccadilly98/goProjects/intelectHome2.0/src/dataBase"
	"github.com/go-chi/chi/v5"
)

type boardInfoGetHanlder struct {
	db *database.DataBase
}

func MakeBoardInfoGetHandler(db *database.DataBase) *boardInfoGetHanlder {
	return &boardInfoGetHanlder{db: db}
}

func (bi *boardInfoGetHanlder) Handler(w http.ResponseWriter, r *http.Request) {
	param := chi.URLParam(r, "board_id")

	if !ProcessingURLParam(w, r, param, bi.db) {
		return
	}
	dto, code, err := bi.db.GetInfoDtoWithId(r.Context(), param)
	if err != nil {
		if code == 0 {
			return
		}
		w.WriteHeader(code)
		w.Write([]byte(err.Error()))
		return
	}
	b, err := json.MarshalIndent(dto, "", "		")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Write(b)
}
