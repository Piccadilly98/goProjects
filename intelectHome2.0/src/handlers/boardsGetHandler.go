package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	database "github.com/Piccadilly98/goProjects/intelectHome2.0/src/dataBase"
)

type boardsGetHandler struct {
	db *database.DataBase
}

func MakeBoardsGetHandler(db *database.DataBase) *boardsGetHandler {
	return &boardsGetHandler{db: db}
}

func (bg *boardsGetHandler) Handler(w http.ResponseWriter, r *http.Request) {
	boardId := r.URL.Query().Get("id")
	boardType := r.URL.Query().Get("type")
	boardState := r.URL.Query().Get("state")
	boardName := r.URL.Query().Get("name")
	w.Header().Set("Content-Type", "application/json")

	res, code, err := bg.db.GetAllBoardsWithConditions(r.Context(), boardState, boardId, boardType, boardName)
	if err != nil {
		if code == 0 {
			return
		}
		w.WriteHeader(code)
		errResponse := fmt.Sprintf(`{"status":"error", "text":"%s"}`, err.Error())
		w.Write([]byte(errResponse))
		return
	}
	b, err := json.MarshalIndent(res, "", "	")
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`{"status":"error", "text":"%s"}`, err.Error())))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}
