package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	dto "github.com/Piccadilly98/goProjects/intellectHome2.0/src/DTO"
	database "github.com/Piccadilly98/goProjects/intellectHome2.0/src/dataBase"
	"github.com/go-chi/chi/v5"
)

type controllerUpdateHandler struct {
	db *database.DataBase
}

func MakeControllerUpdateHandler(db *database.DataBase) *controllerUpdateHandler {
	return &controllerUpdateHandler{db: db}
}

func (cu *controllerUpdateHandler) Handler(w http.ResponseWriter, r *http.Request) {
	boardID := chi.URLParam(r, "board_id")
	controllerID := chi.URLParam(r, "controller_id")
	w.Header().Set("Content-Type", "application/json")
	if !ProcessingURLParam(w, r, boardID, cu.db) {
		return
	}
	controllerType, ok := ProccesingControllerIDGetType(w, r, controllerID, boardID, cu.db)
	if !ok {
		return
	}

	dto := &dto.ControllerUpdateDTO{}

	if err := json.NewDecoder(r.Body).Decode(dto); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	if !dto.ValidateWithType(controllerType) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"status":"error", "text":"invalid body"}`))
		return
	}

	b, code, err := cu.db.UpdateControllerData(r.Context(), boardID, dto, controllerType, controllerID)
	if err != nil {
		if code == 0 {
			return
		}
		w.WriteHeader(code)
		errResponse := fmt.Sprintf(`{"status":"error", "text":"%s"}`, err.Error())
		w.Write([]byte(errResponse))
		return
	}
	w.Write(b)

}
