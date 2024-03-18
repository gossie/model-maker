package rest

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gossie/modelling-service/domain"
)

func (s *server) postConstraint(w http.ResponseWriter, r *http.Request) {
	slog.InfoContext(r.Context(), "creating new parameter")

	// email := r.Context().Value(middleware.UserIdentifierKey).(string)

	decoder := json.NewDecoder(r.Body)
	var ccr domain.ConstraintCreationRequest
	err := decoder.Decode(&ccr)
	if err != nil {
		slog.WarnContext(r.Context(), fmt.Sprintf("could not decode json: %v", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	parameterId, err := s.constraintRepository.SaveConstraint(r.Context(), r.PathValue("modelId"), ccr)
	if err != nil {
		slog.WarnContext(r.Context(), fmt.Sprintf("error creating new constraint: %v", err.Error()))
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

func (s *server) deleteConstraint(w http.ResponseWriter, r *http.Request) {
	modelId, constraintId := r.PathValue("modelId"), r.PathValue("constraintId")
	slog.InfoContext(r.Context(), fmt.Sprintf("deleting constraint - modelId: %v, constraintId: %v", modelId, constraintId))

	err := s.constraintRepository.DeleteConstraint(r.Context(), modelId, constraintId)
	if err != nil {
		slog.WarnContext(r.Context(), fmt.Sprintf("error deleting constraint - modelId = %v, parameterId = %v: %v", modelId, constraintId, err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(200)
}
