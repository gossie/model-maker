package middleware

import (
	"log"
	"net/http"
)

func Trace(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Default().Println("handling", r.Method, r.URL)
		defer log.Default().Println("handled", r.Method, r.URL)

		handler(w, r)
	}
}
