// Package config provides utilities for reading and parsing application configuration.
package config

import (
	"errors"
	"os"

	"gopkg.in/yaml.v2"
)

// Custom errors for configuration handling
var (
	ErrPathNotProvided = errors.New("config path not provided")
	ErrUnableToRead    = errors.New("unable to read config file")
	ErrUnableToParse   = errors.New("unable to parse config file yaml")
)

// Config models the application configuration.
type Config struct {
	Server         *ServerConfig   `yaml:"server"`
	Database       *DatabaseConfig `yaml:"database"`
	Jwt            *Jwt            `yaml:"jwt"`
	PassportClient PassportClient  `yaml:"passport-client"`
}

// ServerConfig models server-specific configuration.
type ServerConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

// DatabaseConfig models database-specific configuration.
type DatabaseConfig struct {
	Host         string `yaml:"host"`
	Port         string `yaml:"port"`
	User         string `yaml:"user"`
	Password     string `yaml:"password"`
	DatabaseName string `yaml:"database-name"`
}

// Jwt models JWT-specific configuration.
type Jwt struct {
	Secret    string `yaml:"secret"`
	ValidMins int    `yaml:"valid-mins"`
}

// PassportClient models passport client-specific configuration.
type PassportClient struct {
	Host string `yaml:"host"`
}

// Read reads a YAML configuration file and unmarshals it into a Config object.
// Parameters:
// - path: The file path of the configuration file.
func Read(path string) (*Config, error) {
	// Validate the provided path
	if path == "" {
		return nil, ErrPathNotProvided
	}

	// Read the file content
	configFile, err := os.ReadFile(path)
	if err != nil {
		return nil, ErrUnableToRead
	}

	// Unmarshal into Config struct
	config := Config{}
	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		return nil, ErrUnableToParse
	}

	return &config, nil
}
