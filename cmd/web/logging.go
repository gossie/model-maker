package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/gossie/modelling-service/middleware"
)

type LogHandler struct {
	slog.Handler
}

func (lh *LogHandler) Handle(ctx context.Context, r slog.Record) error {
	if requestId, ok := ctx.Value(middleware.RequestIdKey).(string); ok {
		r.AddAttrs(slog.String("requestId", requestId))
	}

	if userIdentifier, ok := ctx.Value(middleware.UserIdentifierKey).(string); ok {
		r.AddAttrs(slog.String("userIdentifier", userIdentifier))
	}

	return lh.Handler.Handle(ctx, r)
}
func customizeLogging() {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	wrapper := LogHandler{handler}
	logger := slog.New(&wrapper)
	slog.SetDefault(logger)

	slog.Info("logging was customized, requestId was added")
}
