package config

import "github.com/caarlos0/env/v11"

type Config struct {
	Database DatabaseConfig `envPrefix:"DATABASE_"`
}

type DatabaseConfig struct {
	Host     string `env:"HOST,required"`
	Port     string `env:"PORT,required"`
	User     string `env:"USER,required"`
	Password string `env:"PASSWORD,required"`
	Name     string `env:"NAME,required"`
}

func LoadConfig() (Config, error) {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}
