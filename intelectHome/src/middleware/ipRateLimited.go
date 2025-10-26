package middleware

import (
	"net"
	"net/http"

	"github.com/Piccadilly98/goProjects/intelectHome/src/rate_limit"
	"github.com/Piccadilly98/goProjects/intelectHome/src/storage"
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
				r.URL.Path != "/control" {
				next.ServeHTTP(w, r)
				return
			}
			host, _, err := net.SplitHostPort(r.RemoteAddr)
			// fmt.Println(host)
			if err != nil {
				httpCode = http.StatusInternalServerError
				errors = err.Error()
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
