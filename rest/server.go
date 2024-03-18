package rest

import (
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/gossie/modelling-service/domain"
	"github.com/gossie/modelling-service/middleware"
	"github.com/gossie/modelling-service/persistence"
)

type server struct {
	db                   *sql.DB
	modelRepository      domain.ModelRepository
	constraintRepository domain.ConstraintRepository
	parameterRepository  domain.ParameterRepository
	jwtSecrect           string
}

func NewServer(db *sql.DB, jwtSecrect string) *server {
	modelRepo := persistence.NewPsqlModelRepository(db)
	paramRepo := persistence.NewPsqlParameterRepository(db)
	constRepo := persistence.NewPsqlConstraintRepository(db)

	s := server{
		db,
		&modelRepo,
		&constRepo,
		&paramRepo,
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
	http.HandleFunc("POST /models", middleware.Any(middleware.AuthenticatedRequest(s.jwtSecrect, s.postModel)))
	http.HandleFunc("GET /models", middleware.Any(middleware.AuthenticatedRequest(s.jwtSecrect, s.getModels)))
	http.HandleFunc("GET /models/{modelId}", middleware.Any(middleware.AuthenticatedRequest(s.jwtSecrect, middleware.Authorized(s.db, s.getModel))))
	http.HandleFunc("POST /models/{modelId}/constraints", middleware.Any(middleware.AuthenticatedRequest(s.jwtSecrect, middleware.Authorized(s.db, s.postConstraint))))
	http.HandleFunc("DELETE /models/{modelId}/constraints/{constraintId}", middleware.Any(middleware.AuthenticatedRequest(s.jwtSecrect, middleware.Authorized(s.db, s.deleteConstraint))))
	http.HandleFunc("POST /models/{modelId}/parameters", middleware.Any(middleware.AuthenticatedRequest(s.jwtSecrect, middleware.Authorized(s.db, s.postParameter))))
	http.HandleFunc("GET /models/{modelId}/parameters", middleware.Any(middleware.AuthenticatedRequest(s.jwtSecrect, middleware.Authorized(s.db, s.getParameters))))
	http.HandleFunc("DELETE /models/{modelId}/parameters/{parameterId}", middleware.Any(middleware.AuthenticatedRequest(s.jwtSecrect, middleware.Authorized(s.db, s.deleteParameter))))
	http.HandleFunc("GET /models/{modelId}/parameters/{parameterId}/translations", middleware.Any(middleware.AuthenticatedRequest(s.jwtSecrect, middleware.Authorized(s.db, s.getParameterTranslations))))
	http.HandleFunc("PATCH /models/{modelId}/parameters/{parameterId}/translations", middleware.Any(middleware.AuthenticatedRequest(s.jwtSecrect, middleware.Authorized(s.db, s.patchParameterTranslations))))
	http.HandleFunc("PATCH /models/{modelId}/parameters/{parameterId}/values", middleware.Any(middleware.AuthenticatedRequest(s.jwtSecrect, middleware.Authorized(s.db, s.patchParameterValues))))
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	http.DefaultServeMux.ServeHTTP(w, r)
}
