package middleware

import (
	"net/http"
	"time"

	"github.com/YEgorLu/kafka-test/internal/logger"
)

func Logger(log logger.Logger) Middleware {
	return func(next http.Handler) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			log.Info("request ", r.Method, r.URL)
			next.ServeHTTP(w, r)
			log.Info("response ", " ", time.Since(start).Milliseconds(), " ", r.Method, r.URL, r.ContentLength)
		}
	}
}
