package auth

import (
	"log"
	"net/http"
	"time"

	"github.com/Piccadilly98/goProjects/intelectHome/src/storage"
)

type loginHandlers struct {
	storage     *storage.Storage
	tokenWorker *TokenWorker
}

func MakeLoginHandlers(stor *storage.Storage) *loginHandlers {
	return &loginHandlers{storage: stor, tokenWorker: &TokenWorker{}}
}

func (l *loginHandlers) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("invalid method"))
		return
	}
	ok, login, role := ValidateLoginData(r.Body, l.storage)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid login or password"))
		return
	}
	token, err := l.tokenWorker.CreateToken(login, role, 24*time.Hour)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}
	b, err := l.tokenWorker.TokenToJSON(token)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(b)
}
