package app

import (
	"net/http"

	"github.com/akxcix/nomadcore/pkg/handlers"
	authHandlers "github.com/akxcix/nomadcore/pkg/handlers/auth"
	calHandlers "github.com/akxcix/nomadcore/pkg/handlers/calendar"
	"github.com/akxcix/nomadcore/pkg/services/auth"
	"github.com/akxcix/nomadcore/pkg/services/calendar"
	"github.com/rs/zerolog/log"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	limiter "github.com/ulule/limiter/v3"
	mhttp "github.com/ulule/limiter/v3/drivers/middleware/stdlib"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

func createRoutes(authService *auth.Service, calService *calendar.Service) *chi.Mux {
	rateLimiter := NewRateLimiter()
	limiterMiddleware := rateLimiter.rateLimitMiddleware()
	calHandlers := calHandlers.New(calService)
	authHandlers := authHandlers.New(authService)
	corsHandler := getCorsHandler()
	r := chi.NewRouter()

	// global middlewares
	r.Use(corsHandler)
	r.Use(limiterMiddleware)
	r.Use(handlers.LogRequest)

	// general routes
	r.Get("/health", handlers.HealthCheck)

	// service routes
	r.Post("/calendars/new", authHandlers.AuthMiddleware(calHandlers.CreateCalendar))
	r.Post("/calendars/dates/new", authHandlers.AuthMiddleware(calHandlers.AddDatesToCalendar))
	r.Get("/calendars/public", authHandlers.AuthMiddleware(calHandlers.GetPublicCalendars))

	return r
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

func NewRateLimiter() *RateLimiter {
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
