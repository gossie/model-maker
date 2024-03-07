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
	http.HandleFunc("POST /models", middleware.Any(middleware.AuthenticatedRequest(s.jwtSecrect, s.createModel)))
	http.HandleFunc("GET /models", middleware.Any(middleware.AuthenticatedRequest(s.jwtSecrect, s.getModels)))
	http.HandleFunc("GET /models/{modelId}", middleware.Any(middleware.AuthenticatedRequest(s.jwtSecrect, middleware.Authorized(s.db, s.getModel))))
	http.HandleFunc("POST /models/{modelId}/parameters", middleware.Any(middleware.AuthenticatedRequest(s.jwtSecrect, middleware.Authorized(s.db, s.createParameter))))
	http.HandleFunc("GET /models/{modelId}/parameters", middleware.Any(middleware.AuthenticatedRequest(s.jwtSecrect, middleware.Authorized(s.db, s.getParameters))))
	http.HandleFunc("DELETE /models/{modelId}/parameters/{parameterId}", middleware.Any(middleware.AuthenticatedRequest(s.jwtSecrect, middleware.Authorized(s.db, s.deleteParameter))))
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	http.DefaultServeMux.ServeHTTP(w, r)
}
