package application

import (
	"log/slog"
	"sync"

	"github.com/go-playground/validator/v10"
	"github.com/jonathanhu237/when-works/backend/internal/config"
	"github.com/jonathanhu237/when-works/backend/internal/mailer"
	"github.com/jonathanhu237/when-works/backend/internal/models"
)

type Application struct {
	config    config.Config
	logger    *slog.Logger
	models    models.Models
	validator *validator.Validate
	mailer    *mailer.Mailer
	wg        sync.WaitGroup
}

func New(
	cfg config.Config,
	logger *slog.Logger,
	models models.Models,
	validator *validator.Validate,
	mailer *mailer.Mailer,
) *Application {
	return &Application{
		config:    cfg,
		logger:    logger,
		models:    models,
		validator: validator,
		mailer:    mailer,
	}
}
