package auth

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Piccadilly98/goProjects/intellectHome1.0/src/models"
	"github.com/Piccadilly98/goProjects/intellectHome1.0/src/storage"
)

const (
	headerKey   = "Content-Type"
	headerValue = "application/json"
)

type loginHandlers struct {
	storage     *storage.Storage
	tokenWorker *TokenWorker
	sm          *SessionManager
}

func MakeLoginHandlers(stor *storage.Storage, sm *SessionManager, tw *TokenWorker) *loginHandlers {
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

	value := r.Header.Get(headerKey)
	if value != headerValue {
		httpCode = http.StatusBadRequest
		errors = "invalid format body data"
		w.WriteHeader(httpCode)
		w.Write([]byte(errors))
		return
	}

	ok, login, role := ValidateLoginData(r.Body, l.storage)
	if !ok {
		httpCode = http.StatusUnauthorized
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
		httpCode = http.StatusConflict
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
