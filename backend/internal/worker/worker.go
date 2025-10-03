package worker

import (
	"log/slog"

	"github.com/jonathanhu237/when-works/backend/internal/config"
	"github.com/jonathanhu237/when-works/backend/internal/mailer"
)

type Worker struct {
	config config.Config
	logger *slog.Logger
	mailer *mailer.Mailer
}

func New(cfg config.Config, logger *slog.Logger, mailer *mailer.Mailer) *Worker {
	return &Worker{
		config: cfg,
		logger: logger,
		mailer: mailer,
	}
}
