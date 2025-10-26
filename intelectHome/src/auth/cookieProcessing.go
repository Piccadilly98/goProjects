package auth

import (
	"net/http"
)

func ProcessingCookie(r *http.Request) (bool, string) {
	cookie, err := r.Cookie("jwt_token")
	if err != nil {
		return false, ""
	}
	token := cookie.Value
	return true, token
}
