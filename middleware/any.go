package middleware

import (
	"context"
	"net/http"
)

type language string

const LanguageKey = language("language")

func Any(next http.HandlerFunc) http.HandlerFunc {
	return traceRequest(
		profile(
			logIncomingRequests(
				withLanguage(next))))
}

func withLanguage(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		lang := r.URL.Query().Get("lang")
		if lang == "" {
			lang = "de" // TODO: define default language somewhere
		}

		next(w, r.WithContext(context.WithValue(r.Context(), LanguageKey, lang)))
	}
}
