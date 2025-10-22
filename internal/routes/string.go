package routes

import (
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

	return r
}
