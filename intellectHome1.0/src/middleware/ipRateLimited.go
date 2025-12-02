package middleware

import (
	"net"
	"net/http"
	"strings"

	"github.com/Piccadilly98/goProjects/intellectHome1.0/src/rate_limit"
	"github.com/Piccadilly98/goProjects/intellectHome1.0/src/storage"
)

func IpRateLimiter(ipRate *rate_limit.IpRateLimiter, stor *storage.Storage) func(http.Handler) http.Handler {
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
			if r.URL.Path != "/login" && r.URL.Path != "/boards" && r.URL.Path != "/boards/esp32_1" &&
				r.URL.Path != "/devices" && r.URL.Path != "/devices/led1" && r.URL.Path != "/logs" &&
				r.URL.Path != "/control" && !strings.Contains(r.URL.Path, "/quick-auth-admin") && !strings.Contains(r.URL.Path, "/admin/global") {
				httpCode = http.StatusNotFound
				errors = "404 page not found"
				w.WriteHeader(httpCode)
				w.Write([]byte(errors))
				return
			}
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
			if !ipRate.Allow(host) {
				httpCode = http.StatusTooManyRequests
				errors = "many requests in ip: " + host + " rejected in ipRateLimited"
				w.WriteHeader(httpCode)
				w.Write([]byte(errors))
				return
			}
			deferNeed = false
			next.ServeHTTP(w, r)
		})
	}
}
