package app

import (
	"fmt"
	"net/http"

	"github.com/akxcix/nomadcore/pkg/config"
	"github.com/akxcix/nomadcore/pkg/services/calendar"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type application struct {
	Config     *config.Config
	CalService *calendar.Service
	Routes     *chi.Mux
}

func readConfigs() *config.Config {
	config, err := config.Read("./config.yml")
	if err != nil {
		log.Fatal().Err(err)
	}

	return config
}

func createServices(conf *config.Config) *calendar.Service {
	if conf == nil {
		log.Fatal().Msg("Conf is nil")
	}

	authService := calendar.New(conf.Database, conf.Jwt)

	return authService
}

func new() *application {
	config := readConfigs()

	calService := createServices(config)
	routes := createRoutes(calService)

	app := application{
		Config:     config,
		CalService: calService,
		Routes:     routes,
	}

	return &app
}

func Run() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	app := new()

	addr := fmt.Sprintf("%s:%s", app.Config.Server.Host, app.Config.Server.Port)
	log.Info().Msg(fmt.Sprintf("Running application at %s", addr))
	err := http.ListenAndServe(addr, app.Routes)
	log.Fatal().Err(err).Msg("Crashed")
}
