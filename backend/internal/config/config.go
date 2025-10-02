package config

import "github.com/caarlos0/env/v11"

type Config struct {
	Database DatabaseConfig `envPrefix:"DATABASE_"`
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
