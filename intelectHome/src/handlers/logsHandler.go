package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Piccadilly98/goProjects/intelectHome/src/models"
	"github.com/Piccadilly98/goProjects/intelectHome/src/storage"
)

const (
	headersKey = "format"
	format     = "JSON"
)

type logsHandler struct {
	storage *storage.Storage
}

func MakeLogsHandler(stor *storage.Storage) *logsHandler {
	return &logsHandler{storage: stor}
}

func (l *logsHandler) LogsHandler(w http.ResponseWriter, r *http.Request) {
	formatKey := r.Header.Get(headersKey)
	logsIDstr := r.URL.Query().Get("logsID")
	jwtIDstr := r.URL.Query().Get("jwtID")

	httpCode := http.StatusOK
	errors := ""
	attentions := make([]string, 0)
	jwtClaims, ok := r.Context().Value("jwtClaims").(*models.ClaimsJSON)
	if !ok {
		errors = "server error"
		w.WriteHeader(http.StatusInternalServerError)
		l.storage.NewLog(r, nil, httpCode, errors)
		w.Write([]byte(errors))
		return
	}
	defer func() {
		l.storage.NewLog(r, jwtClaims, httpCode, errors, attentions...)
	}()
	if jwtIDstr != "" {
		logs := l.storage.GetLogsJWTIDJSON(jwtIDstr)
		if len(logs) == 0 {
			errors = "invalid jwtID"
			httpCode = http.StatusBadRequest
			w.WriteHeader(httpCode)
			w.Write([]byte(errors))
			return
		}
		if formatKey == format {
			b, err := json.Marshal(logs)
			if err != nil {
				errors = err.Error()
				httpCode = http.StatusInternalServerError
				w.WriteHeader(httpCode)
				w.Write([]byte(errors))
				return
			}
			w.Write(b)
			return
		}
	}
	if logsIDstr != "" {
		logsID, err := strconv.Atoi(logsIDstr)
		if err != nil {
			httpCode = http.StatusBadRequest
			errors = "invalid ID"
			w.WriteHeader(httpCode)
			w.Write([]byte(errors))
			return
		}
		log, err := l.storage.GetLog(logsID)
		if err != nil {
			httpCode = http.StatusBadRequest
			errors = err.Error()
			w.WriteHeader(httpCode)
			w.Write([]byte(err.Error()))
			return
		}
		if formatKey != "" {
			if formatKey != format {
				errors = "invalid format key"
				httpCode = http.StatusBadRequest
				w.WriteHeader(httpCode)
				w.Write([]byte(errors))
				return
			}
			log := l.storage.GetLogJson(logsID)
			b, err := json.Marshal(log)
			if err != nil {
				httpCode = http.StatusInternalServerError
				errors = err.Error()
				w.WriteHeader(httpCode)
				return
			}
			w.Write(b)
			return
		}
		w.Write([]byte(log))
		return
	}
	if formatKey != "" {
		if formatKey == format {
			logs := l.storage.GetAllLogsJSON()
			b, err := json.Marshal(logs)
			if err != nil {
				httpCode = http.StatusInternalServerError
				errors = err.Error()
				w.WriteHeader(httpCode)
				w.Write([]byte(err.Error()))
				return
			}
			w.Write(b)
			return
		}
		errors = "invalid formatKey"
		httpCode = http.StatusBadRequest
		w.WriteHeader(httpCode)
		w.Write([]byte(errors))
		return
	}
	logs := l.storage.GetAllLogs()
	b, err := json.Marshal(logs)
	if err != nil {
		httpCode = http.StatusInternalServerError
		errors = err.Error()
		w.WriteHeader(httpCode)
		return
	}
	w.WriteHeader(httpCode)
	w.Write(b)
}
