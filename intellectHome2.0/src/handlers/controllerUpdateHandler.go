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
	if !PorcessingURLParamControllerID(w, r, controllerID, cu.db) {
		return
	}

	dto := &dto.ControllerUpdateDTO{}

	if err := json.NewDecoder(r.Body).Decode(dto); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	if !dto.Validate() {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	sets, args, argnum := cu.db.GetJSONBuilderArgs(boardID, dto)
	queryBinary := cu.db.GetQueryToUpdateConroller(sets, args, argnum, true, false)
	querySensor := cu.db.GetQueryToUpdateConroller(sets, args, argnum, false, true)
	fmt.Println(queryBinary, args)
	w.Write([]byte(querySensor))
}
