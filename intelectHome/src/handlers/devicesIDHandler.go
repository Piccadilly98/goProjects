package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Piccadilly98/goProjects/intelectHome/src/models"
	"github.com/Piccadilly98/goProjects/intelectHome/src/storage"
	"github.com/go-chi/chi/v5"
)

type devicesIDHandler struct {
	storage *storage.Storage
}

func MakeDevicesIDHandler(stor *storage.Storage) *devicesIDHandler {
	return &devicesIDHandler{storage: stor}
}

func (d *devicesIDHandler) DevicesIDHandler(w http.ResponseWriter, r *http.Request) {

	httpCode := http.StatusOK
	errors := ""
	attentions := make([]string, 0)
	jwtClaims, ok := r.Context().Value("jwtClaims").(*models.ClaimsJSON)
	if !ok {
		errors = "server error"
		w.WriteHeader(http.StatusInternalServerError)
		d.storage.NewLog(r, nil, httpCode, errors)
		w.Write([]byte(errors))
		return
	}
	defer func() {
		d.storage.NewLog(r, jwtClaims, httpCode, errors, attentions...)
	}()

	deviceID := chi.URLParam(r, "deviceID")
	if !d.storage.CheckIdDevice(deviceID) {
		httpCode = http.StatusBadRequest
		errors = "ivalid device id"
		w.WriteHeader(httpCode)
		w.Write([]byte(errors))
		return
	}
	device, err := d.storage.GetDeviceInfo(deviceID)
	if err != nil {
		httpCode = http.StatusBadRequest
		errors = err.Error()
		w.WriteHeader(httpCode)
		w.Write([]byte(err.Error()))
		return
	}
	b, err := json.Marshal(device)
	if err != nil {
		httpCode = http.StatusInternalServerError
		errors = err.Error()
		w.WriteHeader(httpCode)
		return
	}
	w.Write(b)
}
