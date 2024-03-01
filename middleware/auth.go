package middleware

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func AuthenticatedRequest(secret string, handler http.HandlerFunc) http.HandlerFunc {
	return Trace(
		ContentType(
			"application/json",
			func(w http.ResponseWriter, r *http.Request) {
				token, _ := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
				if !verifyToken(token, secret) {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
				handler(w, r)
			}))
}

func verifyToken(tokenStr, secret string) bool {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return false
	}

	return token.Valid
}
