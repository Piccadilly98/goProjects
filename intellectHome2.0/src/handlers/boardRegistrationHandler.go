package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	dto "github.com/Piccadilly98/goProjects/intelectHome2.0/src/DTO"
	database "github.com/Piccadilly98/goProjects/intelectHome2.0/src/dataBase"
)

type boardRegistration struct {
	db *database.DataBase
}

func MakeRegistrationHandler(db *database.DataBase) *boardRegistration {
	return &boardRegistration{db: db}
}

func (br *boardRegistration) RegistrationHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	upload := &dto.UploadBoardDataDto{}
	w.Header().Set("Content-Type", "application/json")
	if !br.readBodyWriteHeader(w, r, upload) {
		return
	}

	code, err := br.db.RegistrationBoard(ctx, upload.BoardId, upload.Name, upload.BoardType, upload.BoardState)
	if err != nil {
		log.Println(err.Error())
		if code == 0 {
			return
		}
		w.WriteHeader(code)
		errResponse := fmt.Sprintf(`{"status":"error", "text":"%s"}`, err.Error())
		w.Write([]byte(errResponse))
		return
	}
	w.WriteHeader(code)
	w.Write([]byte(`{"status":"ok"}`))
}

func (br *boardRegistration) readBodyWriteHeader(w http.ResponseWriter, r *http.Request, uploadInfo *dto.UploadBoardDataDto) bool {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("empty body"))
		return false
	}
	err = json.Unmarshal(body, uploadInfo)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid body"))
		log.Println(err)
		return false
	}
	if !uploadInfo.ValidateAndDefault() {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"errror":"board id incorrect"}`))
		return false
	}
	return true
}
