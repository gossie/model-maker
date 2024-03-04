package middleware

import (
	"log/slog"
	"net/http"
)

func Trace(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("handling", "method", r.Method, "url", r.URL)
		handler(w, r)
	}
}
