package handlers

import (
	"fmt"
	"log"
	"net/http"

	database "github.com/Piccadilly98/goProjects/intelectHome2.0/src/dataBase"
)

func ProcessingURLParam(w http.ResponseWriter, r *http.Request, param string, db *database.DataBase) bool {
	if param == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"status":"error", "text":"invalid board_id"}`))
		return false
	}
	exist, code, err := db.GetExistWithBoardId(r.Context(), param)
	if err != nil {
		log.Printf("error: %s in %s\n", err.Error(), r.URL)
		if code == 0 {
			return false
		}
		w.WriteHeader(code)
		errResponse := `{"status":"error", "text":"server error"}`
		w.Write([]byte(errResponse))
		return false
	}
	if !exist {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 page not found"))
		return false
	}
	return true
}

func PorcessingURLParamControllerID(w http.ResponseWriter, r *http.Request, param string, db *database.DataBase) bool {
	if param == "" {
		w.WriteHeader(http.StatusNotFound)
		return false
	}

	exist, code, err := db.GetExistWithDeviceId(r.Context(), param)
	if err != nil {
		if code == 0 {
			return false
		}
		w.WriteHeader(code)
		errResponse := fmt.Sprintf(`{"status":"error", "text":"%s"}`, err.Error())
		w.Write([]byte(errResponse))
		return false
	}
	if !exist {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 page not found"))
		return false
	}
	return true
}
