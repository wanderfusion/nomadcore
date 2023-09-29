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

type application struct {
	host   string
	port   string
	routes *chi.Mux
}

func readConfigs() *config.Config {
	config, err := config.Read("./config.yml")
	if err != nil {
		log.Fatal().Err(err)
	}

	return config
}

func createServices(conf *config.Config) (*auth.Service, *group.Service, *users.Service) {
	if conf == nil {
		log.Fatal().Msg("Conf is nil")
	}

	authService := auth.New(conf.Jwt)
	groupService := group.New(conf.Database, &conf.PassportClient)
	usersService := users.New(conf.Database, &conf.PassportClient)

	return authService, groupService, usersService
}

func new() *application {
	config := readConfigs()

	authService, groupService, usersService := createServices(config)
	routes := createRoutes(authService, groupService, usersService)

	app := application{
		host:   config.Server.Host,
		port:   config.Server.Port,
		routes: routes,
	}

	return &app
}

func Run() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	app := new()

	addr := fmt.Sprintf("%s:%s", app.host, app.port)
	log.Info().Msg(fmt.Sprintf("Running application at %s", addr))
	err := http.ListenAndServe(addr, app.routes)
	log.Fatal().Err(err).Msg("Crashed")
}
