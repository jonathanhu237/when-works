package main

import (
	"log/slog"

	"github.com/jonathanhu237/when-works/backend/internal/config"
	"github.com/jonathanhu237/when-works/backend/internal/helpers"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	helpers.InitLogger(cfg)
	slog.Info("logger initialized successfully")
}
