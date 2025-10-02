package application

import (
	"log/slog"

	"github.com/jonathanhu237/when-works/backend/internal/config"
)

type Application struct {
	config config.Config
	logger *slog.Logger
}

func New(cfg config.Config, logger *slog.Logger) *Application {
	return &Application{
		config: cfg,
		logger: logger,
	}
}
