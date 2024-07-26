package middleware

import (
	"net/http"

	"github.com/YEgorLu/kafka-test/internal/logger"
)

func Recover(log logger.Logger) Middleware {
	return func(next http.Handler) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					log.Error(err)
					w.WriteHeader(http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		}
	}
}
