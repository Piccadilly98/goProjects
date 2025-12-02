package handlers

import (
	"fmt"
	"net/http"

	database "github.com/Piccadilly98/goProjects/intelectHome2.0/src/dataBase"
	"github.com/go-chi/chi/v5"
)

type boardControllersHandler struct {
	db *database.DataBase
}

func MakeBoardControllersHandler(db *database.DataBase) *boardControllersHandler {
	return &boardControllersHandler{db: db}
}

func (bc *boardControllersHandler) Handler(w http.ResponseWriter, r *http.Request) {
	param := chi.URLParam(r, "board_id")
	w.Header().Set("Content-Type", "application/json")
	if !ProcessingURLParam(w, r, param, bc.db) {
		return
	}
	res, code, err := bc.db.GetControllersByte(r.Context(), param)
	if err != nil {
		if code == 0 {
			return
		}
		w.WriteHeader(code)
		errResponse := fmt.Sprintf(`{"status":"error", "text":"%s"}`, err.Error())
		w.Write([]byte(errResponse))
		return
	}
	w.Write(res)
}
