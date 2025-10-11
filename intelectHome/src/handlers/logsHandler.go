package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Piccadilly98/goProjects/intelectHome/src/storage"
)

const (
	key        = "key123"
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
	if r.Method != http.MethodGet {
		err := "invalid method"
		w.WriteHeader(http.StatusMethodNotAllowed)
		l.storage.NewLog(r, nil, http.StatusMethodNotAllowed, err)
		w.Write([]byte(err))
		return
	}
	keyParam := r.URL.Query().Get("key")
	formatKey := r.Header.Get(headersKey)
	if keyParam != key {
		err := "error key"
		w.WriteHeader(http.StatusBadRequest)
		l.storage.NewLog(r, []byte(fmt.Sprintf("%s:%s", err, keyParam)), http.StatusBadRequest, err)
		w.Write([]byte(err))
		return
	}
	logsIDstr := r.URL.Query().Get("logsID")
	if logsIDstr != "" {
		logsID, err := strconv.Atoi(logsIDstr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			err := "invalid ID"
			l.storage.NewLog(r, []byte(fmt.Sprintf("%s:%s", err, logsIDstr)), http.StatusBadRequest, err)
			w.Write([]byte(err))
			return
		}
		log, err := l.storage.GetLog(logsID)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			l.storage.NewLog(r, []byte(fmt.Sprintf("%s:%s", err.Error(), log)), http.StatusBadRequest, err.Error())
			w.Write([]byte(err.Error()))
			return
		}
		if formatKey != "" {
			if formatKey != format {
				w.WriteHeader(http.StatusBadRequest)
				err := "invalid format key"
				l.storage.NewLog(r, []byte(formatKey), http.StatusBadRequest, err)
				w.Write([]byte(err))
				return
			}
			log := l.storage.GetLogJson(logsID)
			b, err := json.Marshal(log)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
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
				w.WriteHeader(http.StatusInternalServerError)
				l.storage.NewLog(r, []byte(formatKey), http.StatusInternalServerError, err.Error())
				w.Write([]byte(err.Error()))
				return
			}
			w.Write(b)
			return
		}
		err := "invalid formatkey"
		w.WriteHeader(http.StatusBadRequest)
		l.storage.NewLog(r, []byte(formatKey), http.StatusBadRequest, err)
		w.Write([]byte(err))
		return
	}
	logs := l.storage.GetAllLogs()
	// logs = strings.ReplaceAll(logs, "\n", "\\n")
	// logs = strings.ReplaceAll(logs, "\t", "\\t")
	// logs = strings.ReplaceAll(logs, "\r", "\\r")
	b, err := json.Marshal(logs)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		l.storage.NewLog(r, b, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}
