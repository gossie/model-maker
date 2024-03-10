package modellingservice

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gossie/modelling-service/domain"
)

func (s *server) getParameters(w http.ResponseWriter, r *http.Request) {
	slog.InfoContext(r.Context(), "trying to retrieve models")

	modelId := r.PathValue("modelId")

	models, err := s.parameterRepository.FindAllByModelId(r.Context(), modelId)
	if err != nil {
		slog.WarnContext(r.Context(), fmt.Sprintf("error retrieving parameters for model ID %v from database: %v", modelId, err.Error()))
		http.Error(w, err.Error(), 500)
		return
	}

	err = json.NewEncoder(w).Encode(models)
	if err != nil {
		slog.WarnContext(r.Context(), fmt.Sprintf("could not encode json: %v", err.Error()))
		http.Error(w, err.Error(), 500)
		return
	}
}

func (s *server) postParameter(w http.ResponseWriter, r *http.Request) {
	slog.InfoContext(r.Context(), "creating new parameter")

	// email := r.Context().Value(middleware.UserIdentifierKey).(string)

	decoder := json.NewDecoder(r.Body)
	var pmr domain.ParameterCreationRequest
	err := decoder.Decode(&pmr)
	if err != nil {
		slog.WarnContext(r.Context(), fmt.Sprintf("could not decode json: %v", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	parameterId, err := s.parameterRepository.SaveParameter(r.Context(), r.PathValue("modelId"), pmr)
	if err != nil {
		slog.WarnContext(r.Context(), fmt.Sprintf("error creating new parameter: %v", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := parameterCreationResponse{ParameterId: parameterId}
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		slog.WarnContext(r.Context(), fmt.Sprintf("could not encode json: %v", err.Error()))
		http.Error(w, err.Error(), 500)
		return
	}
}

func (s *server) deleteParameter(w http.ResponseWriter, r *http.Request) {
	modelId, parameterId := r.PathValue("modelId"), r.PathValue("parameterId")
	slog.InfoContext(r.Context(), fmt.Sprintf("deleting parameter - modelId: %v, parameterId: %v", modelId, parameterId))

	err := s.parameterRepository.DeleteParameter(r.Context(), modelId, parameterId)
	if err != nil {
		slog.WarnContext(r.Context(), fmt.Sprintf("error deleting parameter - modelId = %v, parameterId = %v: %v", modelId, parameterId, err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(200)
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

func (s *server) putParameterTranslations(w http.ResponseWriter, r *http.Request) {
	modelId, parameterId := r.PathValue("modelId"), r.PathValue("parameterId")
	slog.InfoContext(r.Context(), fmt.Sprintf("saving parameter translations - modelId: %v, parameterId: %v", modelId, parameterId))

	decoder := json.NewDecoder(r.Body)
	var tcr domain.TranslationModificationRequest
	err := decoder.Decode(&tcr)
	if err != nil {
		slog.WarnContext(r.Context(), fmt.Sprintf("could not decode json: %v", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = s.parameterRepository.SaveTranslations(r.Context(), parameterId, tcr)
	if err != nil {
		slog.WarnContext(r.Context(), fmt.Sprintf("could not save translations: %v", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(200)
}
