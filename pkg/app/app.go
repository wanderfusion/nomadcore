// Package app handles the main application logic.
package app

import (
	"fmt"
	"net/http"

	"github.com/akxcix/nomadcore/pkg/config"
	"github.com/akxcix/nomadcore/pkg/services/auth"
	"github.com/akxcix/nomadcore/pkg/services/group"
	"github.com/akxcix/nomadcore/pkg/services/users"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// application struct holds the app configuration.
type application struct {
	host   string
	port   string
	routes *chi.Mux
}

// readConfigs reads the configuration file and returns a Config struct.
func readConfigs() *config.Config {
	config, err := config.Read("./config.yml")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to read config")
	}

	return config
}

// createServices initializes and returns instances of service layers.
func createServices(conf *config.Config) (*auth.Service, *group.Service, *users.Service) {
	if conf == nil {
		log.Fatal().Msg("Configuration is nil")
	}

	authService := auth.New(conf.Jwt)
	groupService := group.New(conf.Database, &conf.PassportClient)
	usersService := users.New(conf.Database, &conf.PassportClient)

	return authService, groupService, usersService
}

// new initializes the application and returns its instance.
func new() *application {
	config := readConfigs()

	authService, groupService, usersService := createServices(config)
	routes := createRoutes(authService, groupService, usersService)

	return &application{
		host:   config.Server.Host,
		port:   config.Server.Port,
		routes: routes,
	}
}

// Run starts the application.
func Run() {
	// Set logging time format
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	// Initialize the application
	app := new()

	// Construct the address and start the HTTP server
	addr := fmt.Sprintf("%s:%s", app.host, app.port)
	log.Info().Msg(fmt.Sprintf("Running application at %s", addr))
	if err := http.ListenAndServe(addr, app.routes); err != nil {
		log.Fatal().Err(err).Msg("Application crashed")
	}
}
