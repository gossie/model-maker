package server

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

func (s *server) getParameters(w http.ResponseWriter, r *http.Request) {
	slog.InfoContext(r.Context(), "trying to retrieve models")

	modelId := r.PathValue("modelId")

	models, err := s.findAllParametersByModelId(r.Context(), modelId)
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

func (s *server) createParameter(w http.ResponseWriter, r *http.Request) {
	slog.InfoContext(r.Context(), "creating new parameter")

	// email := r.Context().Value(middleware.UserIdentifierKey).(string)

	decoder := json.NewDecoder(r.Body)
	var pmr parameterCreationRequest
	err := decoder.Decode(&pmr)
	if err != nil {
		slog.WarnContext(r.Context(), fmt.Sprintf("could not decode json: %v", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	parameterId, err := s.saveParameter(r.Context(), r.PathValue("modelId"), pmr)
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

	err := s.deleteParameterFromModel(r.Context(), modelId, parameterId)
	if err != nil {
		slog.WarnContext(r.Context(), fmt.Sprintf("error deleting parameter - modelId = %v, parameterId = %v: %v", modelId, parameterId, err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(200)
}
