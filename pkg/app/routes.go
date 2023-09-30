// Package app contains core application logic and routing.
package app

import (
	"net/http"

	// Import various packages and handlers
	"github.com/akxcix/nomadcore/pkg/handlers"
	authHandlers "github.com/akxcix/nomadcore/pkg/handlers/auth"
	groupHandlers "github.com/akxcix/nomadcore/pkg/handlers/group"
	usersHandlers "github.com/akxcix/nomadcore/pkg/handlers/users"
	"github.com/akxcix/nomadcore/pkg/services/auth"
	"github.com/akxcix/nomadcore/pkg/services/group"
	"github.com/akxcix/nomadcore/pkg/services/users"
	"github.com/rs/zerolog/log"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	limiter "github.com/ulule/limiter/v3"
	mhttp "github.com/ulule/limiter/v3/drivers/middleware/stdlib"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

// createRoutes initializes the router and applies routes and middleware.
func createRoutes(authService *auth.Service, groupService *group.Service, usersService *users.Service) *chi.Mux {
	r := chi.NewRouter()

	// Apply global middleware
	applyGlobalMiddlewares(r)

	// Initialize service handlers
	groupHandler := groupHandlers.New(groupService)
	authHandler := authHandlers.New(authService)
	usersHandler := usersHandlers.New(usersService)

	// Define API routes
	defineRoutes(r, authHandler, groupHandler, usersHandler)

	return r
}

// applyGlobalMiddlewares applies middleware that affect all routes.
func applyGlobalMiddlewares(r *chi.Mux) {
	r.Use(getCorsHandler())
	r.Use(newRateLimiter().rateLimitMiddleware())
	r.Use(handlers.LogRequest)
}

// defineRoutes specifies the API endpoints and their corresponding handlers.
func defineRoutes(r *chi.Mux, auth *authHandlers.Handlers, group *groupHandlers.Handlers, users *usersHandlers.Handlers) {
	// Health Check
	r.Get("/health", handlers.HealthCheck)

	// Group-related routes
	r.Route("/groups", func(r chi.Router) {
		r.Use(auth.AuthMiddleware)
		r.Get("/me", group.GetGroups)
		r.Get("/{groupID}", group.GetGroupDetails)
		r.Post("/new", group.CreateGroup)
		r.Post("/dates/new", group.AddDatesToGroup)
		r.Post("/users/new", group.AddUsersToGroup)
	})

	// User-related routes
	r.Route("/users", func(r chi.Router) {
		r.Get("/{username}", users.GetUserProfile)
	})
}

// CORS Middleware settings
func getCorsHandler() func(next http.Handler) http.Handler {
	return cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	})
}

// RateLimiter holds rate limiting settings.
type RateLimiter struct {
	store limiter.Store
	rate  limiter.Rate
}

// Initialize a new rate limiter.
func newRateLimiter() *RateLimiter {
	rate, err := limiter.NewRateFromFormatted("500-M")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize rate limiter")
	}

	store := memory.NewStore()

	return &RateLimiter{
		store: store,
		rate:  rate,
	}
}

// rateLimitMiddleware applies the rate limiting middleware.
func (l *RateLimiter) rateLimitMiddleware() func(h http.Handler) http.Handler {
	middleware := mhttp.NewMiddleware(limiter.New(
		l.store,
		l.rate,
		limiter.WithTrustForwardHeader(true),
	))

	return middleware.Handler
}
