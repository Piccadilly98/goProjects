package auth

import (
	"net/http"

	"github.com/Piccadilly98/goProjects/intelectHome/src/storage"
)

func MiddlewareAuth(next http.Handler, stor *storage.Storage) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	})

}
