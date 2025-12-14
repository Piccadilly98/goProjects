package handlers

import (
	"fmt"
	"log"
	"net/http"

	database "github.com/Piccadilly98/goProjects/intellectHome2.0/src/storage/dataBase"
)

func ValidateURLParam(param ...string) bool {
	for _, v := range param {
		if v == "" {
			return false
		}
	}
	return true
}

func ProcessingURLParam(w http.ResponseWriter, r *http.Request, param string, db *database.DataBase) bool {
	if !ValidateURLParam(param) {
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

func ProccesingControllerIDGetType(w http.ResponseWriter, r *http.Request, conrollerID, boardID string, db *database.DataBase) (string, bool) {
	if !ValidateURLParam(conrollerID) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"status":"error", "text":"invalid board_id"}`))
		return "", false
	}
	controllerType, boardIDRes, code, err := db.GetControllerTypeAndBoardID(r.Context(), conrollerID)
	if err != nil {
		if code == 0 {
			return "", false
		}
		if code == http.StatusBadRequest {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 page not found"))
			return "", false
		}
		w.WriteHeader(code)
		str := fmt.Sprintf(`{"status":"error", "%s"}`, err.Error())
		w.Write([]byte(str))
		return "", false
	}
	if boardID != boardIDRes {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 page not found"))
		return "", false
	}

	return controllerType, true
}
