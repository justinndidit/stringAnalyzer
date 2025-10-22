package routes

import (
	"net/http"

	chi "github.com/go-chi/chi/v5"
	"github.com/justinndidit/stringAnalyzer/internal/application"
)

func SetupAuthRoutes(app *application.Application) *chi.Mux {
	r := chi.NewRouter()
	r.Post("/strings", app.Handler.UploadString)
	r.Get("/strings", app.Handler.GetFilteredStrings)
	r.Get("/strings/{string_value}", app.Handler.GetString)
	r.Get("/strings/filter-by-natural-language", app.Handler.FilterByNaturalLanguage)
	r.Delete("/strings/{string_value}", app.Handler.DeleteString)
	r.Get("/kaithheathcheck", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("String Analyzer API"))
	})

	return r
}
