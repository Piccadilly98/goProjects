package handlers

import (
	"encoding/json"
	"net/http"

	database "github.com/Piccadilly98/goProjects/intellectHome2.0/src/storage/dataBase"
	"github.com/go-chi/chi/v5"
)

type boardIDGet struct {
	db *database.DataBase
}

func MakeBoardIDGet(db *database.DataBase) *boardIDGet {
	return &boardIDGet{db: db}
}

func (bID *boardIDGet) Handler(w http.ResponseWriter, r *http.Request) {
	param := chi.URLParam(r, "board_id")
	w.Header().Set("Content-Type", "application/json")
	if !ProcessingURLParam(w, r, param, bID.db) {
		return
	}

	dto, code, err := bID.db.GetDtoWithId(r.Context(), param)
	if err != nil {
		if code == 0 {
			return
		}
		w.WriteHeader(code)
		w.Write([]byte(err.Error()))
		return
	}

	b, _ := json.MarshalIndent(dto, "", "	")
	w.Write(b)
}
