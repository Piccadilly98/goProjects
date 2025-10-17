package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/Piccadilly98/goProjects/intelectHome/src/models"
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
	httpCode := http.StatusOK
	errors := ""
	attentions := make([]string, 0)
	jwtClaims, ok := r.Context().Value("jwtClaims").(*models.ClaimsJSON)
	if !ok {
		errors = "server error"
		w.WriteHeader(http.StatusInternalServerError)
		h.storage.NewLog(r, nil, httpCode, errors)
		w.Write([]byte(errors))
		return
	}
	defer func() {
		h.storage.NewLog(r, jwtClaims, httpCode, errors, attentions...)
	}()
	reqInfo := &RequestInfo{}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		httpCode = http.StatusInternalServerError
		errors = err.Error()
		w.WriteHeader(httpCode)
		w.Write([]byte(err.Error()))
		return
	}
	err = json.Unmarshal(body, reqInfo)
	if err != nil {
		httpCode = http.StatusInternalServerError
		errors = err.Error()
		w.WriteHeader(httpCode)
		w.Write([]byte(errors))
		return
	}
	if !h.storage.CheckIdDevice(reqInfo.ID) {
		httpCode = http.StatusBadRequest
		errors = "error: device not found"
		w.WriteHeader(httpCode)
		w.Write([]byte(errors))
		return
	}
	if !reqInfo.Validate() {
		httpCode = http.StatusBadRequest
		errors = "invalid request"
		w.WriteHeader(httpCode)
		w.Write([]byte(errors))
		return
	}
	h.storage.UpdateStatusDevice(reqInfo.ID, reqInfo.Status)
}
