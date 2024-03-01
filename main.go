package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gossie/modelling-service/middleware"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	jwtSecrect := os.Getenv("JWT_SECRET")

	http.HandleFunc("POST /models", middleware.Trace(middleware.Auth(jwtSecrect, createModel)))
	http.HandleFunc("GET /models/{modelId}", middleware.Trace(middleware.Auth(jwtSecrect, getModel)))

	log.Default().Println("starting server on port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
