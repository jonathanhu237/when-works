package logger

import (
	"log/slog"
	"os"

	"github.com/jonathanhu237/when-works/backend/internal/config"
)

func Init(cfg config.Config) {
	var handler slog.Handler

	if cfg.Environment == config.Development {
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	} else {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)
}
