package middleware

import "net/http"

func addContentType(contentType string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", contentType)
		next(w, r)
	}
}
