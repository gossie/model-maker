package middleware

import "net/http"

func ContentType(contentType string, handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		handler(w, r)
	}
}
