package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/Piccadilly98/goProjects/intelectHome/src/models"
	"github.com/Piccadilly98/goProjects/intelectHome/src/storage"
)

func MiddlewareAuth(stor *storage.Storage, sm *sessionManager) func(http.Handler) http.Handler {
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

			header := r.Header.Get("Authorization")
			if header == "" {
				httpCode = http.StatusForbidden
				errors = "No authorizations method and token!"
				w.WriteHeader(httpCode)
				w.Write([]byte(errors))
				return
			}
			if !strings.Contains(header, "Bearer") {
				httpCode = http.StatusForbidden
				errors = "No authorizations method!"
				w.WriteHeader(httpCode)
				w.Write([]byte(errors))
				return
			}
			header = strings.ReplaceAll(header, "Bearer", "")
			header = strings.TrimSpace(header)
			if header == "" {
				httpCode = http.StatusForbidden
				errors = "NO JWT TOKEN!"
				w.WriteHeader(httpCode)
				w.Write([]byte(errors))
				return
			}
			ok, claims := ValidateToken(header, stor)
			if !ok {
				httpCode = http.StatusForbidden
				errors = "Invalid token!"
				w.WriteHeader(httpCode)
				w.Write([]byte(errors))
				return
			}
			// if sm.CheckBlackListJWT(claims.TokenID) {
			// 	httpCode = http.StatusBadRequest
			// 	errors = "jwt in BL"
			// 	w.WriteHeader(httpCode)
			// 	return
			// }
			jwtClaims = claims
			if claims.Role == "ADMIN" {
				ctx := context.WithValue(r.Context(), "jwtClaims", claims)
				next.ServeHTTP(w, r.WithContext(ctx))
				deferNeed = false
				return
			}
			if strings.HasPrefix(claims.Role, "ESP32") {
				if strings.HasPrefix(r.URL.Path, "/boards") || strings.HasPrefix(r.URL.Path, "/devices") {
					ctx := context.WithValue(r.Context(), "jwtClaims", claims)
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
