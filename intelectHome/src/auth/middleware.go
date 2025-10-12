package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/Piccadilly98/goProjects/intelectHome/src/storage"
)

func MiddlewareAuth(stor *storage.Storage) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if header == "" {
				w.WriteHeader(http.StatusForbidden)
				w.Write([]byte("No authorizations!"))
				return
			}
			if !strings.Contains(header, "Bearer") {
				w.WriteHeader(http.StatusForbidden)
				w.Write([]byte("No authorizations method!"))
				return
			}
			header = strings.ReplaceAll(header, "Bearer", "")
			header = strings.TrimSpace(header)
			if header == "" {
				w.WriteHeader(http.StatusForbidden)
				w.Write([]byte("No JWT Token!"))
				return
			}
			ok, claims := ValidateToken(header, stor)
			if !ok {
				w.WriteHeader(http.StatusForbidden)
				w.Write([]byte("Invalid token!"))
				return
			}
			if claims.Role == "ADMIN" {
				id, err := claims.GetSubject()
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				ctx := context.WithValue(r.Context(), "userID", id)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
			if strings.HasPrefix(claims.Role, "ESP32") {
				if strings.HasPrefix(r.URL.Path, "/boards") || strings.HasPrefix(r.URL.Path, "/devices") {
					id, err := claims.GetSubject()
					if err != nil {
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
					ctx := context.WithValue(r.Context(), "userID", id)
					next.ServeHTTP(w, r.WithContext(ctx))
					return
				}
				w.WriteHeader(http.StatusForbidden)
				w.Write([]byte("no access rights"))
				return
			}
		})
	}
}
