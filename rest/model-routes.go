package rest

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gossie/modelling-service/domain"
	"github.com/gossie/modelling-service/middleware"
)

func (s *server) login(secret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		var info loginInfo
		err := decoder.Decode(&info)
		if err != nil {
			slog.WarnContext(r.Context(), fmt.Sprintf("could not decode json: %v", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
		}

		slog.InfoContext(r.Context(), fmt.Sprintf("loging in %v", info.Email))
		// check username & password

		var userId int
		row := s.db.QueryRowContext(r.Context(), "SELECT id FROM users WHERE email = $1", info.Email)
		switch err := row.Scan(&userId); err {
		case sql.ErrNoRows:
			slog.InfoContext(r.Context(), fmt.Sprintf("could not find user with email %v", info.Email))
			w.WriteHeader(http.StatusNotFound)
		case nil:
			slog.InfoContext(r.Context(), fmt.Sprintf("found user with email %v", info.Email))
			token, err := createToken(secret, info.Email)
			if err != nil {
				slog.WarnContext(r.Context(), fmt.Sprintf("could not create token: %v", err.Error()))
				http.Error(w, err.Error(), 500)
				return
			}

			resp := loginResponse{Token: token}
			err = json.NewEncoder(w).Encode(resp)
			if err != nil {
				slog.WarnContext(r.Context(), fmt.Sprintf("could not encode json: %v", err.Error()))
				http.Error(w, err.Error(), 500)
				return
			}
			slog.InfoContext(r.Context(), "user was successfully loged in")
		}
	}
}

func createToken(secret string, email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"sub": email,
			"exp": time.Now().Add(time.Hour * 24).Unix(),
		})

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *server) postModel(w http.ResponseWriter, r *http.Request) {
	slog.InfoContext(r.Context(), "creating new model")

	email := r.Context().Value(middleware.UserIdentifierKey).(string)

	decoder := json.NewDecoder(r.Body)
	var cmr domain.ModelCreationRequest
	err := decoder.Decode(&cmr)
	if err != nil {
		slog.WarnContext(r.Context(), fmt.Sprintf("could not decode json: %v", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	modelId, err := s.modelRepository.SaveModel(r.Context(), email, cmr)
	if err != nil {
		slog.InfoContext(r.Context(), fmt.Sprintf("error creating new model: %v", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
	}

	resp := modelCreationResponse{ModelId: modelId}
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		slog.WarnContext(r.Context(), fmt.Sprintf("could not encode json: %v", err.Error()))
		http.Error(w, err.Error(), 500)
		return
	}
}

func (s *server) getModels(w http.ResponseWriter, r *http.Request) {
	slog.InfoContext(r.Context(), "retrieving models")

	email := r.Context().Value(middleware.UserIdentifierKey).(string)

	models, err := s.modelRepository.FindAllByUser(r.Context(), email)
	if err != nil {
		slog.WarnContext(r.Context(), fmt.Sprintf("error retrieving models from database: %v", err.Error()))
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

func (s *server) getModel(w http.ResponseWriter, r *http.Request) {
	modelId := r.PathValue("modelId")
	slog.InfoContext(r.Context(), fmt.Sprintf("retrieving model with id %v", modelId))

	response, err := s.modelRepository.FindById(r.Context(), modelId)
	if err != nil {
		slog.WarnContext(r.Context(), fmt.Sprintf("could not find model with id %v", modelId))
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		slog.WarnContext(r.Context(), fmt.Sprintf("could not encode json: %v", err.Error()))
		http.Error(w, err.Error(), 500)
		return
	}
}