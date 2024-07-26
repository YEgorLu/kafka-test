package middleware

import (
	"net/http"
	"strings"
)

type CORSConfig struct {
	AllowedOrigins   []string
	AllowCredentials bool
}

func CORS(conf CORSConfig) Middleware {
	return func(next http.Handler) http.HandlerFunc {
		originsString := strings.Join(conf.AllowedOrigins, ",")
		allowCredentials := "false"
		if conf.AllowCredentials {
			allowCredentials = "true"
		}
		return func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodOptions {
				next.ServeHTTP(w, r)
			} else {
				w.Header().Add("Access-Control-Allow-Origin", originsString)
				w.Header().Add("Access-Control-Allow-Credentials", allowCredentials)
			}
		}
	}
}
