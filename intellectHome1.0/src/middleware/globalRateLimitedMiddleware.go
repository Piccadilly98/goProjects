package middleware

import (
	"net"
	"net/http"
	"strings"

	"github.com/Piccadilly98/goProjects/intelectHome/src/rate_limit"
	"github.com/Piccadilly98/goProjects/intelectHome/src/storage"
)

func GlobalRateLimiterToMiddleware(rl *rate_limit.GlobalRateLimiter, stor *storage.Storage) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			httpCode := http.StatusOK
			errors := ""
			attentions := make([]string, 0)
			deferNeed := true
			defer func() {
				if deferNeed {
					stor.NewLog(r, nil, httpCode, errors, attentions...)
				}
			}()
			host, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				httpCode = http.StatusInternalServerError
				errors = err.Error()
				return
			}
			if strings.Contains(host, "::1") {
				deferNeed = false
				next.ServeHTTP(w, r)
				return
			}
			if !rl.Allow() {
				httpCode = http.StatusTooManyRequests
				errors = "global rate limited to many requests, request rejected!"
				w.WriteHeader(httpCode)
				return
			}
			deferNeed = false
			next.ServeHTTP(w, r)
		})
	}
}
