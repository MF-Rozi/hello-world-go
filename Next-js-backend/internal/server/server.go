package server

//TODO: - make Chi Server and Routes
import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func routes() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Welcome to the Next.js Backend!")
	})
	return r
}
func StartServer() {
	r := routes()

	fmt.Println("Starting server on :8080")
	http.ListenAndServe(":8080", r)
}
