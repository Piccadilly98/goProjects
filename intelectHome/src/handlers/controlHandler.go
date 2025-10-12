package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/Piccadilly98/goProjects/intelectHome/src/storage"
)

type RequestInfo struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

func (r *RequestInfo) Validate() bool {
	if r.ID == "" {
		return false
	} else if r.Status != storage.StatusOFF && r.Status != storage.StatusON {
		return false
	}
	return true
}

type HandlerControl struct {
	storage *storage.Storage
}

func MakeHandlerControl(storage *storage.Storage) *HandlerControl {
	return &HandlerControl{storage: storage}
}

func (h *HandlerControl) Contorol(w http.ResponseWriter, r *http.Request) {
	reqInfo := &RequestInfo{}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		h.storage.NewLog(r, body, http.StatusInternalServerError, err.Error())
		return
	}
	json.Unmarshal(body, reqInfo)

	if !h.storage.CheckIdDevice(reqInfo.ID) {
		code := http.StatusBadRequest
		w.WriteHeader(code)
		h.storage.NewLog(r, body, code, "errors: device not found")
		return
	}
	if !reqInfo.Validate() {
		code := http.StatusBadRequest
		w.WriteHeader(code)
		h.storage.NewLog(r, body, code, "invalid request")
		return
	}
	h.storage.NewLog(r, body, http.StatusOK, "")
	h.storage.UpdateStatusDevice(reqInfo.ID, reqInfo.Status)

}
