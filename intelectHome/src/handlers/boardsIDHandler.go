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
	switch r.Method {
	case http.MethodGet:
		db, err := b.storage.GetBoardInfo(boardID)
		if err != nil {
			w.Write([]byte(err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			b.storage.NewLog(r, nil, http.StatusBadRequest, err.Error())
			return
		}
		res, err := json.Marshal(db)
		if err != nil {
			w.Write([]byte(err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			b.storage.NewLog(r, nil, http.StatusInternalServerError, err.Error())
			return
		}
		w.Write(res)
		b.storage.NewLog(r, []byte(db.String()), http.StatusOK, "")
		return
	case http.MethodPost:
		db := models.DataBoard{}
		res, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			b.storage.NewLog(r, nil, http.StatusInternalServerError, err.Error())
			return
		}
		json.Unmarshal(res, &db)
		b.storage.AddNewBoardInfo(&db)
		b.storage.NewLog(r, res, http.StatusOK, "")
		return
	}
}
