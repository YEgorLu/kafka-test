package middleware

import "net/http"

type Middleware func(next http.Handler) http.HandlerFunc

func Combine(m ...Middleware) Middleware {
	return func(next http.Handler) http.HandlerFunc {
		last := next
		for _, v := range m {
			last = v(last)
		}
		return func(w http.ResponseWriter, r *http.Request) {
			last.ServeHTTP(w, r)
		}
	}
}
