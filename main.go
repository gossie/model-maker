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
	if jwtSecrect == "" {
		panic("no JWT_SECRET was passed")
	}

	http.HandleFunc("POST /login", middleware.Trace(middleware.ContentType("application/json", login(jwtSecrect))))
	http.HandleFunc("POST /models", middleware.AuthenticatedRequest(jwtSecrect, createModel))
	http.HandleFunc("GET /models/{modelId}", middleware.AuthenticatedRequest(jwtSecrect, getModel))

	log.Default().Println("starting server on port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
