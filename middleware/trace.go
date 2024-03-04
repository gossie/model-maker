package middleware

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
)

type requestId string

const RequestIdKey = requestId("requestId")

func Trace(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		newRequestId := uuid.NewString()
		logInComingRequests(handler)(w, r.WithContext(context.WithValue(r.Context(), RequestIdKey, newRequestId)))
	}
}

func logInComingRequests(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.InfoContext(r.Context(), "handle incoming request", "method", r.Method, "url", r.URL)
		handler(w, r)
	}
}
