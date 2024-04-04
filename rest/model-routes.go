package rest

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gossie/modelling-service/domain"
	"github.com/gossie/modelling-service/middleware"
)

func (s *server) postModel(w http.ResponseWriter, r *http.Request) {
	modelName := r.FormValue("modelName")

	slog.InfoContext(r.Context(), fmt.Sprintf("creating new model with name %v", modelName))

	email := r.Context().Value(middleware.UserIdentifierKey).(string)

	cmr := domain.ModelCreationRequest{Name: modelName}

	_, err := s.modelRepository.SaveModel(r.Context(), email, cmr)
	if err != nil {
		slog.InfoContext(r.Context(), fmt.Sprintf("error creating new model: %v", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
	}

	s.renderModelList(w, r, email)
}

func (s *server) getModels(w http.ResponseWriter, r *http.Request) {
	slog.InfoContext(r.Context(), "retrieving models")

	email := r.Context().Value(middleware.UserIdentifierKey).(string)

	s.renderModelCatalog(w, r, email)
}

func (s *server) getModel(w http.ResponseWriter, r *http.Request) {
	modelId, _ := strconv.Atoi(r.PathValue("modelId"))
	slog.InfoContext(r.Context(), fmt.Sprintf("retrieving model with id %v", modelId))

	s.renderModel(w, r, modelId)
}
