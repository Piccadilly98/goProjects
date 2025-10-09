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
	switch r.Method {
	case http.MethodGet:
		b, err := json.Marshal(d.storage.GetAllDevicesStatusJson())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			d.storage.NewLog(r, b, http.StatusInternalServerError, err.Error())
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusOK)
		d.storage.NewLog(r, b, http.StatusOK, "")
		w.Write(b)
		return

	case http.MethodPost:
		var deviceStat []models.Device_data
		body := r.Body
		attentions := ""
		res, err := io.ReadAll(body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			d.storage.NewLog(r, res, http.StatusInternalServerError, err.Error())
			return
		}
		err = json.Unmarshal(res, &deviceStat)
		if err != nil {
			var device models.Device_data
			err = json.Unmarshal(res, &device)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				d.storage.NewLog(r, res, http.StatusInternalServerError, err.Error())
				return
			}
			if deviceInServer, _ := d.storage.GetDeviceInfo(device.ID); deviceInServer.BoadrId != device.BoadrId {
				w.WriteHeader(http.StatusBadRequest)
				err := "discrepancy boardID"
				d.storage.NewLog(r, res, http.StatusBadRequest, err)
				w.Write([]byte(err))
				return
			}
			if !d.storage.CheckIdDevice(device.ID) {
				w.WriteHeader(http.StatusBadRequest)
				err := fmt.Sprintf("invalid device: %s", device.ID)
				w.Write([]byte(err))
				d.storage.NewLog(r, res, http.StatusBadRequest, err)
				return
			}
			if !d.storage.UpdateStatusDevice(device.ID, device.Status) {
				w.WriteHeader(http.StatusBadRequest)
				err := fmt.Sprintf("error in update Status device: %s", device.ID)
				w.Write([]byte(err))
				d.storage.NewLog(r, res, http.StatusBadRequest, err)
				return
			}
			if statusServer := d.storage.GetDeviceStatus(device.ID); statusServer != device.ID {
				attentions += fmt.Sprintf("Attention! Device %s in esp = %s, server = %s\n", device.ID, statusServer, device.Status)
			}
			d.storage.NewLog(r, res, http.StatusOK, attentions)
			return
		}
		for _, v := range deviceStat {
			if !d.storage.CheckIdDevice(v.ID) {
				w.WriteHeader(http.StatusBadRequest)
				err := fmt.Sprintf("invalid device: %s", v.ID)
				w.Write([]byte(err))
				d.storage.NewLog(r, res, http.StatusBadRequest, err)
				return
			}
			if !d.storage.UpdateStatusDevice(v.ID, v.Status) {
				w.WriteHeader(http.StatusBadRequest)
				err := fmt.Sprintf("error in update Status device: %s", v.ID)
				w.Write([]byte(err))
				d.storage.NewLog(r, res, http.StatusBadRequest, err)
				return
			}
			statusServer, _ := d.storage.GetDeviceInfo(v.ID)
			if statusServer.Status != v.Status {
				attentions += fmt.Sprintf("Attention! Device %s in esp = %s, server = %s\n", v.ID, statusServer, v.Status)
			} else if statusServer.BoadrId != v.BoadrId {
				w.WriteHeader(http.StatusBadRequest)
				err := "discrepancy boardID"
				d.storage.NewLog(r, res, http.StatusBadRequest, err)
				w.Write([]byte(err))
			}
		}
		d.storage.NewLog(r, res, http.StatusOK, attentions)
		return
	}
}
