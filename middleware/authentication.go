package middleware

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

type userIdentifier string

const UserIdentifierKey = userIdentifier("userIdentifier")

func AuthenticatedRequest(secret string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("accessToken")
		if err != nil {
			slog.InfoContext(r.Context(), "no access token")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		token := verifyToken(cookie.Value, secret)
		if token == nil || !token.Valid {
			slog.InfoContext(r.Context(), "token is not valid")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		subject, _ := token.Claims.GetSubject() // TODO: handle err
		next(w, r.WithContext(context.WithValue(r.Context(), UserIdentifierKey, subject)))
	}
}

func verifyToken(tokenStr, secret string) *jwt.Token {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return nil
	}

	return token
}
