package server

import (
	"database/sql"
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
	http.HandleFunc("POST /login", middleware.Trace(middleware.ContentType("application/json", s.login(s.jwtSecrect))))
	http.HandleFunc("POST /models", middleware.AuthenticatedRequest(s.jwtSecrect, s.createModel))
	http.HandleFunc("GET /models/{modelId}", middleware.AuthenticatedRequest(s.jwtSecrect, s.getModel))

}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	http.DefaultServeMux.ServeHTTP(w, r)
}
