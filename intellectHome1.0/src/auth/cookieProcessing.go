package auth

import (
	"net/http"
	"os"
)

func ProcessingCookie(r *http.Request) (bool, string) {
	cookie, err := r.Cookie(os.Getenv("COOKIE_NAME"))
	if err != nil {
		return false, ""
	}
	token := cookie.Value
	return true, token
}
