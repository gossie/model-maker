package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type loginInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token string `json:"token"`
}

type modelCreationRequest struct {
	Name string `json:"name"`
}

type modelCreationResponse struct {
	ModelId int `json:"modelId"`
}

type model struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func (s *server) login(secret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		var info loginInfo
		_ = decoder.Decode(&info)
		// check username & password

		token, err := createToken(secret, info.Username)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		resp := loginResponse{Token: token}
		err = json.NewEncoder(w).Encode(resp)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		slog.InfoContext(r.Context(), "user was successfully loged in")
	}
}

func createToken(secret, username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": username,
			"exp":      time.Now().Add(time.Hour * 24).Unix(),
		})

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *server) createModel(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var cmr modelCreationRequest
	err := decoder.Decode(&cmr)
	if err != nil {
		slog.WarnContext(r.Context(), fmt.Sprintf("could not decode json: %v", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	sqlStatement := "INSERT INTO models (name) VALUES ($1) RETURNING id"
	modelId := 0
	err = s.db.QueryRowContext(r.Context(), sqlStatement, cmr.Name).Scan(&modelId)
	if err != nil {
		slog.WarnContext(r.Context(), fmt.Sprintf("error creating new model: %v", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	slog.InfoContext(r.Context(), fmt.Sprintf("created model with ID %v", modelId))

	resp := modelCreationResponse{ModelId: modelId}
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		slog.WarnContext(r.Context(), fmt.Sprintf("could not encode json: %v", err.Error()))
		http.Error(w, err.Error(), 500)
		return
	}
}

func (s *server) getModel(w http.ResponseWriter, r *http.Request) {
	modelId := r.PathValue("modelId")
	slog.InfoContext(r.Context(), fmt.Sprintf("retrieving model with id %v", modelId))

	sqlStatement := "SELECT * FROM models WHERE id = $1"
	var id int
	var name string
	row := s.db.QueryRowContext(r.Context(), sqlStatement, modelId)
	switch err := row.Scan(&id, &name); err {
	case sql.ErrNoRows:
		slog.InfoContext(r.Context(), fmt.Sprintf("could not find model with id %v", modelId))
		w.WriteHeader(http.StatusNotFound)
	case nil:
		resp := model{Id: id, Name: name}
		err = json.NewEncoder(w).Encode(resp)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	default:
		slog.WarnContext(r.Context(), fmt.Sprintf("unexpected error when retrieving model with id %v: %v", modelId, err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
	}
}
