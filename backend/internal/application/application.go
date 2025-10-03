package application

import (
	"log/slog"

	"github.com/go-playground/validator/v10"
	"github.com/hibiken/asynq"
	"github.com/jonathanhu237/when-works/backend/internal/config"
	"github.com/jonathanhu237/when-works/backend/internal/models"
)

type Application struct {
	config      config.Config
	logger      *slog.Logger
	models      models.Models
	validator   *validator.Validate
	asynqClient *asynq.Client
}

func New(
	cfg config.Config,
	logger *slog.Logger,
	models models.Models,
	validator *validator.Validate,
	asynqClient *asynq.Client,
) *Application {
	return &Application{
		config:      cfg,
		logger:      logger,
		models:      models,
		validator:   validator,
		asynqClient: asynqClient,
	}
}
