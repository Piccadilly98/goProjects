package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Piccadilly98/goProjects/intelectHome/src/storage"
)

type boadsHandler struct {
	storage *storage.Storage
}

func MakeBoarsHandler(stor *storage.Storage) *boadsHandler {
	return &boadsHandler{storage: stor}
}

func (b *boadsHandler) BoardsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		err := []byte("InvalidMethod")
		b.storage.NewLog(r, nil, http.StatusMethodNotAllowed, string(err))
		w.Write(err)
		return
	} else {
		if r.Header.Get("format") == "text" {
			a := b.storage.GetAllBoardsInfo()
			str := ""
			for _, v := range a {
				str += v.String()
			}
			b.storage.NewLog(r, []byte(str), http.StatusOK, "")
			w.Write([]byte(str))
			return
		}
		res, err := json.MarshalIndent(b.storage.GetAllBoardsInfo(), "", "	")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			b.storage.NewLog(r, res, http.StatusInternalServerError, err.Error())
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusOK)
		b.storage.NewLog(r, res, http.StatusOK, "")
		w.Write(res)
		return
	}
}
