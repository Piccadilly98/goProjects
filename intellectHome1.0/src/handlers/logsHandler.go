package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Piccadilly98/goProjects/intellectHome1.0/src/models"
	"github.com/Piccadilly98/goProjects/intellectHome1.0/src/storage"
)

const (
	headersKey = "format"
	format     = "text"
)

type logsHandler struct {
	storage *storage.Storage
}

func MakeLogsHandler(stor *storage.Storage) *logsHandler {
	return &logsHandler{storage: stor}
}

func (l *logsHandler) LogsHandler(w http.ResponseWriter, r *http.Request) {
	formatValue := r.Header.Get(headersKey)
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

	if logsIDstr != "" {
		b, httpCode, err := l.GetLogIdFormat(logsIDstr, formatValue)
		if err != nil {
			errors = err.Error()
			w.WriteHeader(httpCode)
			w.Write([]byte(errors))
			return
		}
		w.Write(b)
		return
	}
	if jwtIDstr != "" {
		b, httpCode, err := l.GetLogsJwtIdFormat(jwtIDstr, formatValue)
		if err != nil {
			errors = err.Error()
			w.WriteHeader(httpCode)
			w.Write([]byte(err.Error()))
			return
		}
		w.Write(b)
		return
	}

	if formatValue != format {
		logsJson := l.storage.GetAllLogsJSON()
		b, err := json.Marshal(logsJson)
		if err != nil {
			httpCode = http.StatusInternalServerError
			errors = err.Error()
			w.WriteHeader(httpCode)
			w.Write([]byte(errors))
			return
		}
		w.Write(b)
		return
	}
	logs := l.storage.GetAllLogs()
	b, err := json.Marshal(logs)
	if err != nil {
		httpCode = http.StatusInternalServerError
		errors = err.Error()
		w.WriteHeader(httpCode)
		w.Write([]byte(errors))
		return
	}
	w.Write(b)

}

func (l *logsHandler) GetLogIdFormat(logsIDstr string, formatValue string) ([]byte, int, error) {
	httpCode := 200
	logsID, err := strconv.Atoi(logsIDstr)
	if err != nil {
		httpCode = http.StatusInternalServerError
		return nil, httpCode, err
	}

	if formatValue != format {
		logsIDjson := l.storage.GetLogJson(logsID)
		b, err := json.Marshal(logsIDjson)
		if err != nil {
			httpCode = http.StatusInternalServerError
			return nil, httpCode, err
		}
		return b, httpCode, nil
	} else if formatValue == format {
		logString, err := l.storage.GetLog(logsID)
		if err != nil {
			httpCode = http.StatusBadRequest
			return nil, httpCode, err
		}
		b, err := json.Marshal(logString)
		if err != nil {
			httpCode = http.StatusInternalServerError
			return nil, httpCode, err
		}
		return b, httpCode, err
	} else {
		httpCode = http.StatusBadRequest
		return nil, httpCode, fmt.Errorf("invalid format")
	}
}

func (l *logsHandler) GetLogsJwtIdFormat(jwtIDstr string, formatValue string) ([]byte, int, error) {
	httpCode := 200

	if formatValue != format {
		logsIDjson := l.storage.GetLogsJWTIDJSON(jwtIDstr)
		b, err := json.Marshal(logsIDjson)
		if err != nil {
			httpCode = http.StatusInternalServerError
			return nil, httpCode, err
		}
		return b, httpCode, nil
	} else if formatValue == format {
		logString := l.storage.GetLogsJWTIDString(jwtIDstr)
		b, err := json.Marshal(logString)
		if err != nil {
			httpCode = http.StatusInternalServerError
			return nil, httpCode, err
		}
		return b, httpCode, err
	} else {
		httpCode = http.StatusBadRequest
		return nil, httpCode, fmt.Errorf("invalid format")
	}
}
