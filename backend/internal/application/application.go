package application

import (
	"database/sql"
	"log/slog"

	"github.com/jonathanhu237/when-works/backend/internal/config"
)

type Application struct {
	config config.Config
	logger *slog.Logger
	db     *sql.DB
}

func New(cfg config.Config, logger *slog.Logger, db *sql.DB) *Application {
	return &Application{
		config: cfg,
		logger: logger,
		db:     db,
	}
}
