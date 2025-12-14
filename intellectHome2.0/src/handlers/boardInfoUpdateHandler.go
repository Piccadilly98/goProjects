package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	dto "github.com/Piccadilly98/goProjects/intellectHome2.0/src/DTO"
	database "github.com/Piccadilly98/goProjects/intellectHome2.0/src/storage/dataBase"
	"github.com/go-chi/chi/v5"
)

type updateBoardInfoHandler struct {
	db         *database.DataBase
	chanUpdate chan string
}

func MakeUpdateBoardInfoHandler(db *database.DataBase, chanUpdate chan string) *updateBoardInfoHandler {
	if chanUpdate == nil {
		return nil
	}
	return &updateBoardInfoHandler{db: db, chanUpdate: chanUpdate}
}

func (ub *updateBoardInfoHandler) Handler(w http.ResponseWriter, r *http.Request) {
	param := chi.URLParam(r, "board_id")
	data := &dto.UpdateBoardInfo{}
	w.Header().Set("Content-Type", "application/json")
	if !ub.processingURLAndBody(w, r, param, data) {
		return
	}
	code, err := ub.db.UpdateBoardInfo(r.Context(), param, data)
	if err != nil {
		w.WriteHeader(code)
		strErr := fmt.Sprintf(`{"status":"error", "text":"%s"}`, err.Error())
		w.Write([]byte(strErr))
		return
	}
	ub.chanUpdate <- param
	w.Write([]byte(`{"status":"ok"}`))
}

func (ub *updateBoardInfoHandler) processingURLAndBody(w http.ResponseWriter, r *http.Request, param string, data *dto.UpdateBoardInfo) bool {
	exist, code, err := ub.db.GetExistWithBoardId(r.Context(), param)
	if err != nil {
		log.Println(err.Error())
		if code == 0 {
			return false
		}
		w.WriteHeader(code)
		errResponse := fmt.Sprintf(`{"status":"error", "text":"%s"}`, err.Error())
		w.Write([]byte(errResponse))
		return false
	}
	if !exist {
		log.Println("!exist")
		w.WriteHeader(http.StatusNotFound)
		return false
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"status":"error", "text":"invalid body"}`))
		return false
	}
	err = json.Unmarshal(body, data)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"status":"error", "text":"invalid body"}`))
		return false
	}
	if !data.Validate() {
		log.Println("non valid")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"status":"error", "text":"invalid body"}`))
		return false
	}
	return true
}
