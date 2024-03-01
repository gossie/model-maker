package main

import (
	"encoding/json"
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

func login(secret string) http.HandlerFunc {
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

func createModel(w http.ResponseWriter, r *http.Request) {

}

func getModel(w http.ResponseWriter, r *http.Request) {

}
