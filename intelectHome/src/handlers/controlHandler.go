package handlers

import (
	"encoding/json"
	"fmt"
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
	defer func() {
		h.storage.PrintDataBoards()
		h.storage.PrintDataDevice()
		h.storage.PrintLogs()
	}()
	reqInfo := &RequestInfo{}
	if r.Method == http.MethodPost {
		json.NewDecoder(r.Body).Decode(reqInfo)

		if !h.storage.CheckIdDevice(reqInfo.ID) {
			code := http.StatusBadRequest
			w.WriteHeader(code)
			h.storage.NewLogPost(r, reqInfo.ID, reqInfo.Status, code)
			return
		}
		if !reqInfo.Validate() {
			code := http.StatusBadRequest
			w.WriteHeader(code)
			h.storage.NewLogPost(r, reqInfo.ID, reqInfo.Status, code)
			return
		}
		h.storage.NewLogPost(r, reqInfo.ID, reqInfo.Status, http.StatusOK)
		h.storage.UpdateStatusDevice(reqInfo.ID, reqInfo.Status)
	} else if r.Method == http.MethodGet {
		js := h.storage.GetAllDevicesStatusJson()
		b, err := json.Marshal(js)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			h.storage.NewLogGet(r, b, http.StatusInternalServerError)
			return
		}
		h.storage.NewLogGet(r, b, http.StatusOK)
		w.Write(b)
		return
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid method request"))
		return
	}
}
