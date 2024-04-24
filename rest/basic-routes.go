package rest

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/gossie/modelling-service/views"
)

func (s *Server) GetIndex(v *views.View) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.InfoContext(r.Context(), "retrieving index page")
		v.Render(r.Context(), w, nil)
	}
}

func (s *Server) Login(secret string, v *views.View) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email := r.FormValue("email")

		slog.InfoContext(r.Context(), fmt.Sprintf("loging in %v", email))
		// check username & password

		var userId int
		row := s.db.QueryRowContext(r.Context(), "SELECT id FROM users WHERE email = $1", email)
		switch err := row.Scan(&userId); err {
		case sql.ErrNoRows:
			slog.InfoContext(r.Context(), fmt.Sprintf("could not find user with email %v", email))
			w.WriteHeader(http.StatusNotFound)
		case nil:
			slog.InfoContext(r.Context(), fmt.Sprintf("found user with email %v", email))
			token, err := createToken(secret, email)
			if err != nil {
				slog.WarnContext(r.Context(), fmt.Sprintf("could not create token: %v", err.Error()))
				http.Error(w, err.Error(), 500)
				return
			}

			expiration := time.Now().Add(24 * time.Hour)
			cookie := http.Cookie{Name: "accessToken", Value: token, Expires: expiration}
			http.SetCookie(w, &cookie)

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

			v.Render(r.Context(), w, ModelCatalogRenderContext{Models: renderModels})
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
