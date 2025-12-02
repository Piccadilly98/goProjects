package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Piccadilly98/goProjects/intellectHome1.0/src/models"
	"github.com/Piccadilly98/goProjects/intellectHome1.0/src/storage"
)

const (
	key          = "format"
	formatBoards = "text"
)

type boadsHandler struct {
	storage *storage.Storage
}

func MakeBoarsHandler(stor *storage.Storage) *boadsHandler {
	return &boadsHandler{storage: stor}
}

func (b *boadsHandler) BoardsHandler(w http.ResponseWriter, r *http.Request) {
	httpCode := http.StatusOK
	errors := ""
	attentions := make([]string, 0)
	jwtClaims, ok := r.Context().Value("jwtClaims").(*models.ClaimsJSON)
	if !ok {
		errors = "server error"
		w.WriteHeader(http.StatusInternalServerError)
		b.storage.NewLog(r, nil, httpCode, errors)
		w.Write([]byte(errors))
		return
	}
	defer func() {
		b.storage.NewLog(r, jwtClaims, httpCode, errors, attentions...)
	}()

	if r.Header.Get(key) == formatBoards {
		boards := b.storage.GetAllBoardsInfo()
		httpCode = http.StatusOK
		str := ""
		for _, v := range boards {
			str += v.String()
		}
		w.Write([]byte(str))
		return
	}
	res, err := json.MarshalIndent(b.storage.GetAllBoardsInfo(), "", "	")
	if err != nil {
		httpCode = http.StatusInternalServerError
		w.WriteHeader(httpCode)
		errors = err.Error()
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(httpCode)
	w.Write(res)
}
