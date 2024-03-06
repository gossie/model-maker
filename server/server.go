package server

import (
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/gossie/modelling-service/middleware"
)

type server struct {
	db         *sql.DB
	jwtSecrect string
}

func NewServer(db *sql.DB, jwtSecrect string) *server {
	s := server{
		db,
		jwtSecrect,
	}
	s.routes()
	return &s
}

func (s *server) routes() {
	http.HandleFunc("OPTIONS /", middleware.Any(func(w http.ResponseWriter, r *http.Request) {
		slog.InfoContext(r.Context(), "handling options request")
		w.WriteHeader(http.StatusNoContent)
	}))
	http.HandleFunc("POST /login", middleware.Any(s.login(s.jwtSecrect)))
	http.HandleFunc("POST /models", middleware.AuthenticatedRequest(s.jwtSecrect, s.createModel))
	http.HandleFunc("GET /models", middleware.AuthenticatedRequest(s.jwtSecrect, s.getModels))
	http.HandleFunc("GET /models/{modelId}", middleware.AuthenticatedRequest(s.jwtSecrect, s.getModel))
	http.HandleFunc("POST /models/{modelId}/parameters", middleware.AuthenticatedRequest(s.jwtSecrect, s.createParameter))
	http.HandleFunc("GET /models/{modelId}/parameters", middleware.AuthenticatedRequest(s.jwtSecrect, s.getParameters))
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	http.DefaultServeMux.ServeHTTP(w, r)
}
