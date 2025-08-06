package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

func main() {

	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file")
		return
	}

	chi := chi.NewRouter()
	chi.Use(middleware.Logger)
	chi.Use(middleware.Recoverer)

	chi.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Welcome to the Chi Web Service!")
	})

	http.ListenAndServe(":8080", chi)
}
