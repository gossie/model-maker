package main

import (
	"database/sql"
	"embed"
	"fmt"
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/gossie/modelling-service/rest"
	_ "github.com/lib/pq"
)

type InputField struct {
	Label       string
	Name        string
	Type        string
	Placeholder string
}

type Option struct {
	Key, Value string
}

type SelectBox struct {
	Label   string
	Name    string
	Options []Option
}

type PrimaryButton struct {
	Label string
}

//go:embed templates/*
var htmlTemplates embed.FS

var tmpl = template.Must(template.New("").Funcs(template.FuncMap{
	"inputField": func(label, name, fieldType, placeholder string) InputField {
		return InputField{Label: label, Name: name, Type: fieldType, Placeholder: placeholder}
	},
	"options": func(args ...string) []Option {
		options := make([]Option, 0, len(args)/2)
		for i := 0; i < len(args); i += 2 {
			options = append(options, Option{Key: args[i], Value: args[i+1]})
		}
		return options
	},
	"selectBox": func(label, name string, options []Option) SelectBox {
		return SelectBox{Label: label, Name: name, Options: options}
	},
	"primaryButton": func(label string) PrimaryButton {
		return PrimaryButton{Label: label}
	},
}).ParseFS(htmlTemplates, "templates/*.html"))

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
	customizeLogging()

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

	svr := rest.NewServer(tmpl, db, jwtSecrect)

	slog.Info("starting server on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, svr))
}
