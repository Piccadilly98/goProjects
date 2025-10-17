package auth

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Piccadilly98/goProjects/intelectHome/src/models"
	"github.com/Piccadilly98/goProjects/intelectHome/src/storage"
)

type loginHandlers struct {
	storage     *storage.Storage
	tokenWorker *TokenWorker
	sm          *sessionManager
}

func MakeLoginHandlers(stor *storage.Storage, sm *sessionManager, tw *TokenWorker) *loginHandlers {
	return &loginHandlers{storage: stor, tokenWorker: tw, sm: sm}
}

func (l *loginHandlers) LoginHandler(w http.ResponseWriter, r *http.Request) {
	httpCode := http.StatusOK
	errors := ""
	attentions := make([]string, 0)
	jwtClaims := &models.ClaimsJSON{}
	defer func() {
		l.storage.NewLog(r, jwtClaims, httpCode, errors, attentions...)
	}()

	ok, login, role := ValidateLoginData(r.Body, l.storage)
	if !ok {
		httpCode = http.StatusBadRequest
		errors = "INVALID LOGIN ON PASSWORD"
		w.WriteHeader(httpCode)
		w.Write([]byte("invalid login or password"))
		return
	}
	token, id, err := l.tokenWorker.CreateToken(login, role, 24*time.Hour)
	if err != nil {
		httpCode = http.StatusInternalServerError
		errors = err.Error() + "Errors in generation jwt token"
		w.WriteHeader(httpCode)
		log.Println(err.Error())
		return
	}
	ok = l.sm.NewSession(login, role, token, 24*time.Hour, id)
	if !ok {
		l.tokenWorker.tokenIDCount.Add(-1)
		httpCode = http.StatusBadRequest
		errors = "Repeat get jwt token!"
		w.WriteHeader(httpCode)
		w.Write([]byte(errors))
		return
	}
	attentions = append(attentions, fmt.Sprintf("CREATE NEW JWT KEY: ID: %d", id))
	jwtClaims.Subject = login
	jwtClaims.Role = role
	jwtClaims.TokenID = id
	// jwtClaims.SessionID = sessionId

	b, err := l.tokenWorker.TokenToJSON(token)
	if err != nil {
		httpCode = http.StatusInternalServerError
		errors = err.Error() + "Errors in token to json response"
		w.WriteHeader(httpCode)
		return
	}
	w.Write(b)
}
