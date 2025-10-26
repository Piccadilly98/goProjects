package handlers

import (
	"net/http"

	"github.com/Piccadilly98/goProjects/intelectHome/src/auth"
	"github.com/Piccadilly98/goProjects/intelectHome/src/models"
	"github.com/Piccadilly98/goProjects/intelectHome/src/storage"
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
	hash := r.URL.Query().Get("hash")
	if !q.sm.CheckActiveSession(hash) {
		httpCode = http.StatusNotFound
		w.WriteHeader(httpCode)
		w.Write([]byte("404 page not found"))
		errors = "not contains admin hash!"
		return
	}
	cookie := &http.Cookie{
		Name:     "jwt_token",
		Value:    q.token,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}

	http.SetCookie(w, cookie)
	w.Write([]byte("Authorization complete"))
}
