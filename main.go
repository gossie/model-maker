package main

import (
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/gossie/modelling-service/middleware"
	_ "github.com/lib/pq"
)

const (
	defaultDBHost = "localhost"
	defaultDBPort = "5432"
	defaultDBUser = "postgres"
	defaultDBName = "modelling"
)

func connectToDB() *sql.DB {
	dbHost := getOrDefault("DB_HOST", defaultDBHost)
	dbPort := getOrDefault("DB_PORT", defaultDBPort)
	dbUser := getOrDefault("DB_USER", defaultDBUser)
	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		panic("no DB_PASSWORD was passed")
	}
	dbName := getOrDefault("DB_NAME", defaultDBName)
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	return db
}

func getOrDefault(env, defaultValue string) string {
	value := os.Getenv(env)
	if value == "" {
		slog.Info("no " + env + " was set, using default: " + defaultValue)
		value = defaultValue
	}
	return value
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	jwtSecrect := os.Getenv("JWT_SECRET")
	if jwtSecrect == "" {
		panic("no JWT_SECRET was passed")
	}

	db := connectToDB()
	defer db.Close()

	http.HandleFunc("POST /login", middleware.Trace(middleware.ContentType("application/json", login(jwtSecrect))))
	http.HandleFunc("POST /models", middleware.AuthenticatedRequest(jwtSecrect, createModel))
	http.HandleFunc("GET /models/{modelId}", middleware.AuthenticatedRequest(jwtSecrect, getModel))

	slog.Info("starting server on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
