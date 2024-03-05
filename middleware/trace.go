package middleware

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
)

type requestId string

const RequestIdKey = requestId("requestId")

func traceRequest(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		newRequestId := uuid.NewString()
		next(w, r.WithContext(context.WithValue(r.Context(), RequestIdKey, newRequestId)))
	}
}

func logIncomingRequests(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.InfoContext(r.Context(), "handle incoming request", "method", r.Method, "url", r.URL.Path)
		next(w, r)
	}
}
