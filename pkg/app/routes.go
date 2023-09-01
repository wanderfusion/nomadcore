package app

import (
	"github.com/akxcix/nomadcore/pkg/handlers"
	calHandlers "github.com/akxcix/nomadcore/pkg/handlers/calendar"
	"github.com/akxcix/nomadcore/pkg/services/calendar"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func createRoutes(calService *calendar.Service) *chi.Mux {
	r := chi.NewRouter()

	// global middlewares
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
	r.Use(handlers.LogRequest)

	// general routes
	r.Get("/health", handlers.HealthCheck)

	authHandlers := calHandlers.New(calService)
	r.Post("/calendars/new", authHandlers.AuthMiddleware(authHandlers.CreateCalendar))
	r.Get("/users/calendars/public", authHandlers.AuthMiddleware(authHandlers.GetPublicCalendars))

	return r
}
