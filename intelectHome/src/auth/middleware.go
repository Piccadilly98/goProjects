package auth

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/Piccadilly98/goProjects/intelectHome/src/models"
	"github.com/Piccadilly98/goProjects/intelectHome/src/storage"
)

const (
	key    = "Authorization"
	method = "Bearer"
	ctxKey = "jwtClaims"
)

func MiddlewareAuth(stor *storage.Storage, sm *SessionManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			httpCode := http.StatusOK
			errors := ""
			attentions := make([]string, 0)
			var jwtClaims *models.ClaimsJSON = nil
			deferNeed := true
			defer func() {
				if deferNeed {
					stor.NewLog(r, jwtClaims, httpCode, errors, attentions...)
				}
			}()
			if r.URL.Path == "/login" {
				deferNeed = false
				next.ServeHTTP(w, r)
				return
			}
			token, err := validateHeaderGetToken(r.Header)
			if err != nil {
				errors = err.Error()
				httpCode = http.StatusBadRequest
				w.WriteHeader(httpCode)
				w.Write([]byte(errors))
				return
			}
			ok, claims := ValidateToken(token, stor)
			if !ok {
				httpCode = http.StatusUnauthorized
				errors = "Invalid token!"
				w.WriteHeader(httpCode)
				w.Write([]byte(errors))
				return
			}
			jwtClaims = claims
			_, err = sm.CheckTokenValid(token, claims)
			if err != nil {
				httpCode = http.StatusUnauthorized
				w.WriteHeader(httpCode)
				errors = err.Error()
				return
			}
			if claims.Role == "ADMIN" {
				ctx := context.WithValue(r.Context(), ctxKey, claims)
				next.ServeHTTP(w, r.WithContext(ctx))
				deferNeed = false
				return
			}
			if strings.HasPrefix(claims.Role, "ESP32") {
				if strings.HasPrefix(r.URL.Path, "/boards/") || strings.HasPrefix(r.URL.Path, "/devices/") {
					ctx := context.WithValue(r.Context(), ctxKey, claims)
					next.ServeHTTP(w, r.WithContext(ctx))
					deferNeed = false
					return
				}
				w.WriteHeader(http.StatusForbidden)
				w.Write([]byte("no access rights"))
				return
			}
		})
	}
}

func validateHeaderGetToken(r http.Header) (string, error) {
	token := ""
	header := r.Get(key)
	str := strings.Split(header, " ")
	if len(str) < 2 {
		return token, fmt.Errorf("invalid headers")
	}
	if str[0] == method {
		if str[1] != "" {
			token = str[1]
			return token, nil
		}
	}
	return token, fmt.Errorf("no token or method")
}
