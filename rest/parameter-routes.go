package rest

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	configurationmodel "github.com/gossie/configuration-model"
	"github.com/gossie/modelling-service/domain"
)

func (s *server) getParameters(w http.ResponseWriter, r *http.Request) {
	slog.InfoContext(r.Context(), "retrieving parameters")

	modelId, _ := strconv.Atoi(r.PathValue("modelId"))

	searchValue := r.URL.Query().Get("parameterName")
	if searchValue == "" {
		searchValue = "*"
	}
	s.renderParameters(w, r, modelId, searchValue)
}

func (s *server) postParameter(w http.ResponseWriter, r *http.Request) {
	slog.InfoContext(r.Context(), "creating new parameter")

	// email := r.Context().Value(middleware.UserIdentifierKey).(string)

	modelId, _ := strconv.Atoi(r.PathValue("modelId"))
	name := r.FormValue("parameterName")
	valueType, _ := strconv.Atoi(r.FormValue("valueType"))

	pmr := domain.ParameterCreationRequest{Name: name, ValueType: configurationmodel.ValueType(valueType)}

	_, err := s.parameterRepository.SaveParameter(r.Context(), modelId, pmr)
	if err != nil {
		slog.WarnContext(r.Context(), fmt.Sprintf("error creating new parameter: %v", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	s.renderParameters(w, r, modelId, "")
}

func (s *server) deleteParameter(w http.ResponseWriter, r *http.Request) {
	modelId, _ := strconv.Atoi(r.PathValue("modelId"))
	parameterId, _ := strconv.Atoi(r.PathValue("parameterId"))

	slog.InfoContext(r.Context(), fmt.Sprintf("deleting parameter - modelId: %v, parameterId: %v", modelId, parameterId))

	err := s.parameterRepository.DeleteParameter(r.Context(), modelId, parameterId)
	if err != nil {
		slog.WarnContext(r.Context(), fmt.Sprintf("error deleting parameter - modelId = %v, parameterId = %v: %v", modelId, parameterId, err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	s.renderParameters(w, r, modelId, "")
}

func (s *server) getParameterTranslations(w http.ResponseWriter, r *http.Request) {
	modelId, parameterId := r.PathValue("modelId"), r.PathValue("parameterId")
	slog.InfoContext(r.Context(), fmt.Sprintf("retrieving parameter translations - modelId: %v, parameterId: %v", modelId, parameterId))

	translations, err := s.parameterRepository.FindAllTranslations(r.Context(), parameterId)
	if err != nil {
		slog.WarnContext(r.Context(), fmt.Sprintf("could not retrieve translations: %v", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(translations)
	if err != nil {
		slog.WarnContext(r.Context(), fmt.Sprintf("could not encode json: %v", err.Error()))
		http.Error(w, err.Error(), 500)
		return
	}
}

func (s *server) patchParameterTranslations(w http.ResponseWriter, r *http.Request) {
	modelId, parameterId := r.PathValue("modelId"), r.PathValue("parameterId")
	slog.InfoContext(r.Context(), fmt.Sprintf("saving parameter translations - modelId: %v, parameterId: %v", modelId, parameterId))

	decoder := json.NewDecoder(r.Body)
	var tmr domain.TranslationModificationRequest
	err := decoder.Decode(&tmr)
	if err != nil {
		slog.WarnContext(r.Context(), fmt.Sprintf("could not decode json: %v", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = s.parameterRepository.SaveTranslations(r.Context(), parameterId, tmr)
	if err != nil {
		slog.WarnContext(r.Context(), fmt.Sprintf("could not save translations: %v", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(200)
}

func (s *server) patchParameterValues(w http.ResponseWriter, r *http.Request) {
	modelId, parameterId := r.PathValue("modelId"), r.PathValue("parameterId")
	slog.InfoContext(r.Context(), fmt.Sprintf("saving parameter values - modelId: %v, parameterId: %v", modelId, parameterId))

	decoder := json.NewDecoder(r.Body)
	var vmr domain.ValueModificationRequest
	err := decoder.Decode(&vmr)
	if err != nil {
		slog.WarnContext(r.Context(), fmt.Sprintf("could not decode json: %v", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = s.parameterRepository.SaveValues(r.Context(), parameterId, vmr)
	if err != nil {
		slog.WarnContext(r.Context(), fmt.Sprintf("could not save translations: %v", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(200)
}
