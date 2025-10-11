package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Piccadilly98/goProjects/intelectHome/src/storage"
)

const (
	key = "key123"
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
	param := r.URL.Query().Get("key")

	if param != key {
		err := "error key"
		w.WriteHeader(http.StatusBadRequest)
		l.storage.NewLog(r, nil, http.StatusBadRequest, err)
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
