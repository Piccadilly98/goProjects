package handlers

import (
	"encoding/json"
	"net/http"

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
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		d.storage.NewLog(r, nil, http.StatusMethodNotAllowed, "invalid method")
		w.Write([]byte("invalid method"))
		return
	}

	deviceID := chi.URLParam(r, "deviceID")
	if !d.storage.CheckIdDevice(deviceID) {
		w.WriteHeader(http.StatusBadRequest)
		d.storage.NewLog(r, nil, http.StatusBadRequest, "ivalid device id")
		w.Write([]byte("invalid device id"))
		return
	}
	b, err := json.Marshal(d.storage.GetAllDevicesStatusJson())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		d.storage.NewLog(r, b, http.StatusInternalServerError, err.Error())
		return
	}
	d.storage.NewLog(r, b, http.StatusOK, "")
	w.Write(b)
}
