package middleware

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

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

func profile(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		defer func() { slog.InfoContext(r.Context(), fmt.Sprintf("took %v ms", time.Since(start).Milliseconds())) }()

		next(w, r)
	}
}

func logIncomingRequests(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.InfoContext(r.Context(), "handle incoming request", "method", r.Method, "url", r.URL.Path)
		next(w, r)
	}
}
