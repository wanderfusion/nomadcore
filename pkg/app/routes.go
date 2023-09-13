package app

import (
	"net/http"

	"github.com/akxcix/nomadcore/pkg/handlers"
	authHandlers "github.com/akxcix/nomadcore/pkg/handlers/auth"
	groupHandlers "github.com/akxcix/nomadcore/pkg/handlers/group"
	"github.com/akxcix/nomadcore/pkg/services/auth"
	"github.com/akxcix/nomadcore/pkg/services/group"
	"github.com/rs/zerolog/log"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	limiter "github.com/ulule/limiter/v3"
	mhttp "github.com/ulule/limiter/v3/drivers/middleware/stdlib"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

func createRoutes(authService *auth.Service, groupService *group.Service) *chi.Mux {
	r := chi.NewRouter()

	// Apply global middlewares
	applyGlobalMiddlewares(r)

	// Initialize handlers
	groupHandler := groupHandlers.New(groupService)
	authHandler := authHandlers.New(authService)

	// Define routes
	defineRoutes(r, authHandler, groupHandler)

	return r
}

func applyGlobalMiddlewares(r *chi.Mux) {
	r.Use(getCorsHandler())
	r.Use(newRateLimiter().rateLimitMiddleware())
	r.Use(handlers.LogRequest)
}

func defineRoutes(r *chi.Mux, auth *authHandlers.Handlers, group *groupHandlers.Handlers) {
	// Health Check Route
	r.Get("/health", handlers.HealthCheck)

	// Group Routes
	r.Route("/groups", func(r chi.Router) {
		r.Use(auth.AuthMiddleware)
		r.Get("/me", group.GetGroups)
		r.Get("/{groupID}", group.GetGroupDetails) // Get group by ID
		r.Post("/new", group.CreateGroup)
		r.Post("/dates/new", group.AddDatesToGroup)
		r.Post("/users/new", group.AddUsersToGroup) // Add users to group
	})

	// // User Routes
	// r.Route("/users", func(r chi.Router) {
	// 	r.Use(auth.AuthMiddleware)
	// 	r.Get("/{username}/profile", auth.GetUserProfileByID) // Get user profile by ID
	// })
}

// cors -------------------------------------------------------------------------------------------
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

// rate limiter -----------------------------------------------------------------------------------
type RateLimiter struct {
	store limiter.Store
	rate  limiter.Rate
}

func newRateLimiter() *RateLimiter {
	rate, err := limiter.NewRateFromFormatted("500-M")
	if err != nil {
		log.Fatal().Err(err).Msg("unable to initialise limiter")
	}

	store := memory.NewStore()

	limiter := RateLimiter{
		store: store,
		rate:  rate,
	}
	return &limiter
}

func (l *RateLimiter) rateLimitMiddleware() func(h http.Handler) http.Handler {
	middleware := mhttp.NewMiddleware(limiter.New(
		l.store,
		l.rate,
		limiter.WithTrustForwardHeader(true),
	))

	return middleware.Handler
}
