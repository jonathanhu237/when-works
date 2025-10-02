package config

import (
	"fmt"
	"strings"

	"github.com/caarlos0/env/v11"
)

type Environment string

const (
	Development Environment = "development"
	Production  Environment = "production"
)

func (e *Environment) UnmarshalText(text []byte) error {
	s := strings.ToLower(string(text))
	switch s {
	case "development":
		*e = Development
	case "production":
		*e = Production
	default:
		return fmt.Errorf("invalid environment value: %q, must be one of 'development' or 'production'", s)
	}
	return nil
}

type Config struct {
	Environment Environment    `env:"ENVIRONMENT"`
	Database    DatabaseConfig `envPrefix:"DATABASE_"`
}

type DatabaseConfig struct {
	Host     string `env:"HOST"`
	Port     string `env:"PORT"`
	User     string `env:"USER"`
	Password string `env:"PASSWORD"`
	Name     string `env:"NAME"`
}

func LoadConfig() (Config, error) {
	cfg := Config{}
	if err := env.ParseWithOptions(&cfg, env.Options{RequiredIfNoDef: true}); err != nil {
		return Config{}, err
	}
	return cfg, nil
}
