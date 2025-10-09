package handlers

import (
	"encoding/json"
	"fmt"
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
	if r.Method == http.MethodPost {
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
	} else if r.Method == http.MethodGet {
		js := h.storage.GetAllDevicesStatusJson()
		b, err := json.Marshal(js)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			h.storage.NewLog(r, b, http.StatusInternalServerError, err.Error())
			return
		}
		h.storage.NewLog(r, b, http.StatusOK, "")
		w.Write(b)
		return
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid method request"))
		return
	}
}
