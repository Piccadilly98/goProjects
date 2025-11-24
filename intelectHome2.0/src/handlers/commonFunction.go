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
	exist, code, err := db.GetExistWithId(r.Context(), param)
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
		w.WriteHeader(http.StatusNotFound)
		return false
	}
	return true
}
