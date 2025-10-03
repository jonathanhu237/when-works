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
	Environment  Environment        `env:"ENVIRONMENT"`
	Server       ServerConfig       `envPrefix:"SERVER_"`
	Database     DatabaseConfig     `envPrefix:"DATABASE_"`
	InitialAdmin InitialAdminConfig `envPrefix:"INITIAL_ADMIN_"`
	JWT          JWTConfig          `envPrefix:"JWT_"`
	Redis        RedisConfig        `envPrefix:"REDIS_"`
	Asynq        AsynqConfig        `envPrefix:"ASYNQ_"`
	SMTP         SMTPConfig         `envPrefix:"SMTP_"`
}

type ServerConfig struct {
	IdleTimeout     int `env:"IDLE_TIMEOUT"`
	ReadTimeout     int `env:"READ_TIMEOUT"`
	WriteTimeout    int `env:"WRITE_TIMEOUT"`
	ShutdownTimeout int `env:"SHUTDOWN_TIMEOUT"`
}

type DatabaseConfig struct {
	Host            string `env:"HOST"`
	Port            string `env:"PORT"`
	User            string `env:"USER"`
	Password        string `env:"PASSWORD"`
	Name            string `env:"NAME"`
	MaxOpenConns    int    `env:"MAX_OPEN_CONNS"`
	MaxIdleConns    int    `env:"MAX_IDLE_CONNS"`
	ConnMaxIdleTime int    `env:"CONN_MAX_IDLE_TIME"`
	PingTimeout     int    `env:"PING_TIMEOUT"`
	QueryTimeout    int    `env:"QUERY_TIMEOUT"`
}

type InitialAdminConfig struct {
	Username string `env:"USERNAME"`
	Password string `env:"PASSWORD"`
	Email    string `env:"EMAIL"`
}

type JWTConfig struct {
	Secret     string `env:"SECRET"`
	Expiration int    `env:"EXPIRATION"`
}

type RedisConfig struct {
	Host     string `env:"HOST"`
	Port     string `env:"PORT"`
	Password string `env:"PASSWORD"`
	DB       int    `env:"DB"`
}

type AsynqConfig struct {
	MaxRetry    int `env:"MAX_RETRY"`
	Timeout     int `env:"TIMEOUT"`
	Concurrency int `env:"CONCURRENCY"`
}

type SMTPConfig struct {
	Host     string `env:"HOST"`
	Port     int    `env:"PORT"`
	Username string `env:"USERNAME"`
	Password string `env:"PASSWORD"`
	From     string `env:"FROM"`
}

func LoadConfig() (Config, error) {
	cfg := Config{}
	if err := env.ParseWithOptions(&cfg, env.Options{RequiredIfNoDef: true}); err != nil {
		return Config{}, err
	}
	return cfg, nil
}
