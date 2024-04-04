package rest

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gossie/modelling-service/domain"
)

type RenderModel struct {
	Id   int
	Name string
}

type RenderParameter struct {
	Id      int
	ModelId int
	Name    string
	Values  []string
}

type RenderConstraint struct {
}

type ModelCatalogRenderContext struct {
	Models []RenderModel
}

type ModelRenderContext struct {
	Model       RenderModel
	Parameters  []RenderParameter
	Constraints []RenderConstraint
}

func valueOrDefault(value, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}

func (s *server) renderModelCatalog(w http.ResponseWriter, r *http.Request, email string) {

	models, err := s.modelRepository.FindAllByUser(r.Context(), email)
	if err != nil {
		slog.WarnContext(r.Context(), fmt.Sprintf("error retrieving models from database: %v", err.Error()))
		http.Error(w, err.Error(), 500)
		return
	}

	renderModels := make([]RenderModel, len(models))
	for i := range len(models) {
		renderModels[i] = RenderModel{Id: models[i].Id, Name: valueOrDefault(models[i].Name, models[i].Translation)}
	}

	s.tmpl.ExecuteTemplate(w, "model-catalog.html", ModelCatalogRenderContext{Models: renderModels})
}

func (s *server) renderModelList(w http.ResponseWriter, r *http.Request, email string) {
	models, err := s.modelRepository.FindAllByUser(r.Context(), email)
	if err != nil {
		slog.WarnContext(r.Context(), fmt.Sprintf("error retrieving models from database: %v", err.Error()))
		http.Error(w, err.Error(), 500)
		return
	}

	renderModels := make([]RenderModel, len(models))
	for i := range len(models) {
		renderModels[i] = RenderModel{Id: models[i].Id, Name: valueOrDefault(models[i].Name, models[i].Translation)}
	}

	s.tmpl.ExecuteTemplate(w, "model-list", ModelCatalogRenderContext{Models: renderModels})
}

func (s *server) renderModel(w http.ResponseWriter, r *http.Request, modelId int) {
	model, err := retrieveData(nil, func() (domain.Model, error) {
		return s.modelRepository.FindById(r.Context(), modelId)
	})

	parameters, err := retrieveData(err, func() ([]domain.Parameter, error) {
		return s.parameterRepository.FindAllByModelId(r.Context(), modelId)
	})

	if err != nil {
		slog.WarnContext(r.Context(), fmt.Sprintf("could not find model with id %v", modelId))
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

	s.tmpl.ExecuteTemplate(w, "model.html", ModelRenderContext{
		Model:       RenderModel{Id: model.Id, Name: valueOrDefault(model.Translation, model.Name)},
		Parameters:  parametersToRender,
		Constraints: constraintsToRender,
	})
}

func (s *server) renderParameters(w http.ResponseWriter, r *http.Request, modelId int) {
	parameters, err := retrieveData(nil, func() ([]domain.Parameter, error) {
		return s.parameterRepository.FindAllByModelId(r.Context(), modelId)
	})

	if err != nil {
		slog.WarnContext(r.Context(), fmt.Sprintf("could not find model with id %v", modelId))
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

	s.tmpl.ExecuteTemplate(w, "parameter-list", parametersToRender)
}

func retrieveData[T any](err error, retriever func() (T, error)) (T, error) {
	if err == nil {
		return retriever()
	}
	var empty T
	return empty, err
}
