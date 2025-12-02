package handlers

import (
	"net/http"
	"os"

	"github.com/Piccadilly98/goProjects/intellectHome1.0/src/auth"
	"github.com/Piccadilly98/goProjects/intellectHome1.0/src/models"
	"github.com/Piccadilly98/goProjects/intellectHome1.0/src/storage"
)

type QuickAdmin struct {
	stor  *storage.Storage
	sm    *auth.SessionManager
	token string
}

func NewQuickAdmin(stor *storage.Storage, sm *auth.SessionManager, token string) *QuickAdmin {
	if token == "" {
		return nil
	}
	return &QuickAdmin{stor: stor, sm: sm, token: token}
}

func (q *QuickAdmin) AddAdminCookie(w http.ResponseWriter, r *http.Request) {
	httpCode := http.StatusOK
	errors := ""
	attentions := make([]string, 0)
	jwtClaims := &models.ClaimsJSON{}
	defer func() {
		q.stor.NewLog(r, jwtClaims, httpCode, errors, attentions...)
	}()
	ok, jwtClaims := auth.ValidateToken(q.token, q.stor)
	if !ok {
		errors = "invalid admin token"
		httpCode = http.StatusNotFound
		w.WriteHeader(httpCode)
		w.Write([]byte("404 page not found"))
		return
	}
	hash := r.URL.Query().Get("hash")
	if !q.sm.CheckActiveSession(hash) {
		httpCode = http.StatusNotFound
		w.WriteHeader(httpCode)
		w.Write([]byte("404 page not found"))
		errors = "not contains admin hash!"
		return
	}
	cookie := &http.Cookie{
		Name:     os.Getenv("COOKIE_NAME"),
		Value:    q.token,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   60 * 60 * 100,
	}

	http.SetCookie(w, cookie)
	w.Write([]byte("Authorization complete"))
}
