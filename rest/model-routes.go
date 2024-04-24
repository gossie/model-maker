package rest

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gossie/modelling-service/domain"
	"github.com/gossie/modelling-service/middleware"
	"github.com/gossie/modelling-service/views"
)

func (s *Server) PostModel(v *views.View) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		modelName := r.FormValue("modelName")

		slog.InfoContext(r.Context(), fmt.Sprintf("creating new model with name %v", modelName))

		email := r.Context().Value(middleware.UserIdentifierKey).(string)

		cmr := domain.ModelCreationRequest{Name: modelName}

		_, err := s.modelRepository.SaveModel(r.Context(), email, cmr)
		if err != nil {
			slog.InfoContext(r.Context(), fmt.Sprintf("error creating new model: %v", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
		}

		renderModelCatalog(v, w, r, s.modelRepository, email)
	}
}

func (s *Server) GetModels(v *views.View) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.InfoContext(r.Context(), "retrieving models")
		email := r.Context().Value(middleware.UserIdentifierKey).(string)
		renderModelCatalog(v, w, r, s.modelRepository, email)
	}
}

func (s *Server) GetModel(v *views.View) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		modelId, _ := strconv.Atoi(r.PathValue("modelId"))
		slog.InfoContext(r.Context(), fmt.Sprintf("retrieving model with id %v", modelId))

		renderModel(v, w, r, s.modelRepository, s.parameterRepository, modelId)
	}
}

func renderModelCatalog(v *views.View, w http.ResponseWriter, r *http.Request, repo domain.ModelRepository, email string) {
	models, err := repo.FindAllByUser(r.Context(), email)
	if err != nil {
		slog.WarnContext(r.Context(), fmt.Sprintf("error retrieving models from database: %v", err.Error()))
		http.Error(w, err.Error(), 500)
		return
	}

	renderModels := make([]RenderModel, len(models))
	for i := range len(models) {
		renderModels[i] = RenderModel{Id: models[i].Id, Name: valueOrDefault(models[i].Name, models[i].Translation)}
	}

	v.Render(r.Context(), w, ModelCatalogRenderContext{Models: renderModels})
}

func renderModel(v *views.View, w http.ResponseWriter, r *http.Request, modelRepo domain.ModelRepository, paramRepo domain.ParameterRepository, modelId int) {
	model, err := retrieveData(nil, func() (domain.Model, error) {
		return modelRepo.FindById(r.Context(), modelId)
	})

	parameters, err := retrieveData(err, func() ([]domain.Parameter, error) {
		return paramRepo.FindAllByModelId(r.Context(), modelId, "")
	})

	if err != nil {
		slog.WarnContext(r.Context(), fmt.Sprintf("could not find model with id %v: %v", modelId, err.Error()))
		w.WriteHeader(http.StatusNotFound)
		return
	}

	parametersToRender := make([]RenderParameter, len(parameters))
	for i := range parameters {
		values := make([]string, len(parameters[i].Value.Values))
		for j := range parameters[i].Value.Values {
			values[j] = parameters[i].Value.Values[j].Translation
		}

		parametersToRender[i] = RenderParameter{
			Id:      parameters[i].Id,
			ModelId: modelId,
			Name:    valueOrDefault(parameters[i].Translation, parameters[i].Name),
			Values:  values,
		}
	}

	constraintsToRender := make([]RenderConstraint, len(model.Constraints))
	for i := range model.Constraints {
		constraintsToRender[i] = RenderConstraint{}
	}

	v.Render(r.Context(), w, ModelRenderContext{
		Model:       RenderModel{Id: model.Id, Name: valueOrDefault(model.Translation, model.Name)},
		Parameters:  parametersToRender,
		Constraints: constraintsToRender,
	})
}

func retrieveData[T any](err error, retriever func() (T, error)) (T, error) {
	if err == nil {
		return retriever()
	}
	var empty T
	return empty, err
}
