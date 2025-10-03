package application

import (
	"log/slog"

	"github.com/jonathanhu237/when-works/backend/internal/config"
	"github.com/jonathanhu237/when-works/backend/internal/models"
)

type Application struct {
	config config.Config
	logger *slog.Logger
	models models.Models
}

func New(cfg config.Config, logger *slog.Logger, models models.Models) *Application {
	return &Application{
		config: cfg,
		logger: logger,
		models: models,
	}
}
