package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	dbUrl := os.Getenv("DATABASE_URL")
	if dbUrl == "" {
		fmt.Fprintf(os.Stderr, "No DATABASE_URL environment variable!\n")
		os.Exit(1)
	}
	app, err := getApp(dbUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to db: %v\n", err)
		os.Exit(1)
	}
	app.mountHandlers()
	defer app.closeConnection()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.AllowContentType("application/json"))
	// Health check
	r.Get("/", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	r.Mount("/quotes", app.Router)
	http.ListenAndServe(":8080", r)
}
