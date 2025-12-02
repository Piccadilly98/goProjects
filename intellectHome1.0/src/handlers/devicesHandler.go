package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Piccadilly98/goProjects/intelectHome/src/models"
	"github.com/Piccadilly98/goProjects/intelectHome/src/storage"
)

type devicesHandler struct {
	storage *storage.Storage
}

func MakeDevicesHandler(stor *storage.Storage) *devicesHandler {
	return &devicesHandler{storage: stor}
}

func (d *devicesHandler) DevicesHandler(w http.ResponseWriter, r *http.Request) {

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

	switch r.Method {
	case http.MethodGet:
		b, err := json.Marshal(d.storage.GetAllDevicesStatusJson())
		if err != nil {
			httpCode = http.StatusInternalServerError
			errors = err.Error()
			w.WriteHeader(httpCode)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(httpCode)
		w.Write(b)
		return

	case http.MethodPost:
		var deviceStat []models.Device_data
		res, err := io.ReadAll(r.Body)
		if err != nil {
			httpCode = http.StatusInternalServerError
			errors = err.Error()
			w.WriteHeader(httpCode)
			w.Write([]byte(err.Error()))
			return
		}
		err = json.Unmarshal(res, &deviceStat)
		if err != nil {
			var device models.Device_data
			err = json.Unmarshal(res, &device)
			if err != nil {
				httpCode = http.StatusBadRequest
				errors = "empty body!"
				w.WriteHeader(httpCode)
				w.Write([]byte(errors))
				return
			}
			if deviceInServer, _ := d.storage.GetDeviceInfo(device.ID); deviceInServer.BoadrId != device.BoadrId {
				httpCode = http.StatusBadRequest
				errors = "discrepancy boardID"
				w.WriteHeader(httpCode)
				w.Write([]byte(errors))
				return
			}
			if !d.storage.CheckIdDevice(device.ID) {
				httpCode = http.StatusBadRequest
				errors = fmt.Sprintf("invalid device: %s", device.ID)
				w.WriteHeader(httpCode)
				w.Write([]byte(errors))
				return
			}
			statusServer, _ := d.storage.GetDeviceInfo(device.ID)
			if statusServer.Status != device.Status {
				attentions = append(attentions, fmt.Sprintf("Attention! Device %s in esp = %s, server = %s", device.ID, device.Status, statusServer.Status))
			} else if statusServer.BoadrId != device.BoadrId {
				httpCode = http.StatusBadRequest
				errors = "discrepancy boardID"
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(errors))
				return
			}
			if !d.storage.UpdateStatusDevice(device.ID, device.Status) {
				httpCode = http.StatusBadRequest
				errors = fmt.Sprintf("error in update Status device: %s", device.ID)
				w.WriteHeader(httpCode)
				w.Write([]byte(errors))
				return
			}
			return
		}
		for _, v := range deviceStat {
			if !d.storage.CheckIdDevice(v.ID) {
				httpCode = http.StatusBadRequest
				errors = fmt.Sprintf("invalid device: %s", v.ID)
				w.WriteHeader(httpCode)
				w.Write([]byte(errors))
				return
			}
			statusServer, _ := d.storage.GetDeviceInfo(v.ID)
			if statusServer.Status != v.Status {
				attentions = append(attentions, fmt.Sprintf("Attention! Device %s in esp = %s, server = %s", v.ID, v.Status, statusServer.Status))
			} else if statusServer.BoadrId != v.BoadrId {
				httpCode = http.StatusBadRequest
				errors = "discrepancy boardID"
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(errors))
				return
			}
			if !d.storage.UpdateStatusDevice(v.ID, v.Status) {
				httpCode = http.StatusBadRequest
				errors = fmt.Sprintf("error in update Status device: %s", v.ID)
				w.WriteHeader(httpCode)
				w.Write([]byte(errors))
				return
			}
		}
		return
	}
}
