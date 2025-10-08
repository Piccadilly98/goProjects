package handlers

import (
	"fmt"
	"net/http"

	"github.com/Piccadilly98/goProjects/intelectHome/src/storage"
)

type boarsHandl struct {
	storage *storage.Storage
}

func MakeBoarsHandler(st *storage.Storage) *boarsHandl {
	return &boarsHandl{storage: st}
}

func (b *boarsHandl) BoardsHendler(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query().Get("BoardId")
	fmt.Println(v)
	switch r.Method {
	case http.MethodGet:
		w.Write([]byte(b.storage.GetBoardInfo(v)))
		return
	}
}
