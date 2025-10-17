package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/Piccadilly98/goProjects/intelectHome/src/models"
	"github.com/Piccadilly98/goProjects/intelectHome/src/storage"
	"github.com/go-chi/chi/v5"
)

type boarsIDHandl struct {
	storage *storage.Storage
}

func MakeBoarsIDHandler(st *storage.Storage) *boarsIDHandl {
	return &boarsIDHandl{storage: st}
}

func (b *boarsIDHandl) BoardsIDHandler(w http.ResponseWriter, r *http.Request) {
	boardID := chi.URLParam(r, "boardID")

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

	switch r.Method {
	case http.MethodGet:
		db, err := b.storage.GetBoardInfo(boardID)
		if err != nil {
			errors = err.Error()
			httpCode = http.StatusBadRequest
			w.WriteHeader(httpCode)
			w.Write([]byte(errors))
			return
		}
		res, err := json.Marshal(db)
		if err != nil {
			httpCode = http.StatusInternalServerError
			errors = err.Error()
			w.WriteHeader(httpCode)
			w.Write([]byte(err.Error()))
			return
		}
		w.Write(res)
		return

	case http.MethodPost:
		db := models.DataBoard{}
		res, err := io.ReadAll(r.Body)
		if err != nil {
			httpCode = http.StatusInternalServerError
			errors = err.Error()
			w.WriteHeader(httpCode)
			w.Write([]byte(err.Error()))
			return
		}
		json.Unmarshal(res, &db)
		if db.BoardId != boardID {
			errors = "Discrepancy boardID"
			httpCode = http.StatusBadRequest
			w.WriteHeader(httpCode)
			w.Write([]byte(errors))
			return
		}
		if !b.storage.AddNewBoardInfo(&db) {
			errors = "Error!Invalid board id"
			httpCode = http.StatusBadRequest
			w.WriteHeader(httpCode)
			w.Write([]byte(errors))
			return
		}
		return
	}
}
