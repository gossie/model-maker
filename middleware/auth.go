package middleware

import (
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

func Auth(secret string, handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if !verifyToken(token, secret) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		handler(w, r)
	}
}

func verifyToken(tokenStr, secret string) bool {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if err != nil {
		return false
	}

	return token.Valid
}
