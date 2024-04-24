package rest

import (
	"database/sql"
	"net/http"

	"github.com/gossie/modelling-service/domain"
	"github.com/gossie/modelling-service/middleware"
	"github.com/gossie/modelling-service/persistence"
	"github.com/gossie/modelling-service/views"
)

type Server struct {
	db                   *sql.DB
	modelRepository      domain.ModelRepository
	constraintRepository domain.ConstraintRepository
	parameterRepository  domain.ParameterRepository
	jwtSecrect           string
}

func NewServer(db *sql.DB, jwtSecrect string) *Server {
	modelRepo := persistence.NewPsqlModelRepository(db)
	paramRepo := persistence.NewPsqlParameterRepository(db)
	constRepo := persistence.NewPsqlConstraintRepository(db)

	s := Server{
		db,
		&modelRepo,
		&constRepo,
		&paramRepo,
		jwtSecrect,
	}
	s.routes()
	return &s
}

func (s *Server) routes() {
	http.HandleFunc("GET /", middleware.Any(s.GetIndex(views.NewView("index.html"))))
	http.HandleFunc("POST /login", middleware.Any(s.Login(s.jwtSecrect, views.NewView("model-catalog.html"))))
	http.HandleFunc("POST /models", middleware.Any(middleware.AuthenticatedRequest(s.jwtSecrect, s.PostModel(views.NewView("model-list")))))
	http.HandleFunc("GET /models", middleware.Any(middleware.AuthenticatedRequest(s.jwtSecrect, s.GetModels(views.NewView("model-catalog.html")))))
	http.HandleFunc("GET /models/{modelId}", middleware.Any(middleware.AuthenticatedRequest(s.jwtSecrect, middleware.Authorized(s.db, s.GetModel(views.NewView("model.html"))))))
	http.HandleFunc("POST /models/{modelId}/constraints", middleware.Any(middleware.AuthenticatedRequest(s.jwtSecrect, middleware.Authorized(s.db, s.PostConstraint))))
	http.HandleFunc("DELETE /models/{modelId}/constraints/{constraintId}", middleware.Any(middleware.AuthenticatedRequest(s.jwtSecrect, middleware.Authorized(s.db, s.DeleteConstraint))))
	http.HandleFunc("POST /models/{modelId}/parameters", middleware.Any(middleware.AuthenticatedRequest(s.jwtSecrect, middleware.Authorized(s.db, s.PostParameter(views.NewView("parameter-list"))))))
	http.HandleFunc("GET /models/{modelId}/parameters", middleware.Any(middleware.AuthenticatedRequest(s.jwtSecrect, middleware.Authorized(s.db, s.GetParameters(views.NewView("parameter-list"))))))
	http.HandleFunc("DELETE /models/{modelId}/parameters/{parameterId}", middleware.Any(middleware.AuthenticatedRequest(s.jwtSecrect, middleware.Authorized(s.db, s.DeleteParameter(views.NewView("parameter-list"))))))
	http.HandleFunc("GET /models/{modelId}/parameters/{parameterId}/translations", middleware.Any(middleware.AuthenticatedRequest(s.jwtSecrect, middleware.Authorized(s.db, s.GetParameterTranslations))))
	http.HandleFunc("PATCH /models/{modelId}/parameters/{parameterId}/translations", middleware.Any(middleware.AuthenticatedRequest(s.jwtSecrect, middleware.Authorized(s.db, s.PatchParameterTranslations))))
	http.HandleFunc("PATCH /models/{modelId}/parameters/{parameterId}/values", middleware.Any(middleware.AuthenticatedRequest(s.jwtSecrect, middleware.Authorized(s.db, s.PatchParameterValues))))

	// http.HandleFunc("GET /configuration-models/{modelId}", func(w http.ResponseWriter, r *http.Request) {
	// 	confModel := configurationmodel.Model{}

	// 	_, _ = s.modelRepository.FindById(r.Context(), r.PathValue("modelId"))
	// 	parameters, _ := s.parameterRepository.FindAllByModelId(r.Context(), r.PathValue("modelId"))

	// 	for _, p := range parameters {
	// 		var value configurationmodel.ValueModel
	// 		confModel.AddParameter(p.Name, value)
	// 	}

	// 	// encoder := json.NewEncoder(w)
	// 	// err = encoder.Encode(confModel)
	// })
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	http.DefaultServeMux.ServeHTTP(w, r)
}
