package main

import (
	"os"

	"github.com/jonathanhu237/when-works/backend/internal/config"
	"github.com/jonathanhu237/when-works/backend/internal/logger"
	"github.com/jonathanhu237/when-works/backend/internal/mailer"
	"github.com/jonathanhu237/when-works/backend/internal/worker"
)

func main() {
	// ------------------------------
	// Load config
	// ------------------------------
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	// ------------------------------
	// Initialize logger
	// ------------------------------
	logger := logger.Init(cfg)
	logger.Info("logger initialized successfully")

	// ------------------------------
	// Initialize mailer
	// ------------------------------
	mailer, err := mailer.New(cfg)
	if err != nil {
		logger.Error("error initializing mailer", "error", err)
		os.Exit(1)
	}
	logger.Info("mailer initialized successfully")

	// ------------------------------
	// Initialize worker
	// ------------------------------
	worker := worker.New(cfg, logger, mailer)

	// ------------------------------
	// Run worker
	// ------------------------------
	if err := worker.Run(); err != nil {
		logger.Error("error running worker", "error", err)
		os.Exit(1)
	}
}
