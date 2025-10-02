package main

import (
	"os"

	"github.com/jonathanhu237/when-works/backend/internal/application"
	"github.com/jonathanhu237/when-works/backend/internal/config"
	"github.com/jonathanhu237/when-works/backend/internal/helpers"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	logger := helpers.InitLogger(cfg)
	logger.Info("logger initialized successfully")

	app := application.New(cfg, logger)
	if err := app.Serve(); err != nil {
		logger.Error("error starting server", "error", err)
		os.Exit(1)
	}
}
