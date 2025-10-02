package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/jonathanhu237/when-works/backend/internal/application"
	"github.com/jonathanhu237/when-works/backend/internal/config"
	"github.com/jonathanhu237/when-works/backend/internal/logger"

	_ "github.com/jackc/pgx/v5/stdlib"
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
	// Open database
	// ------------------------------
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
	)
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		logger.Error("error opening database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	db.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	db.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	db.SetConnMaxIdleTime(time.Duration(cfg.Database.ConnMaxIdleTime) * time.Second)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.Database.PingTimeout)*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		logger.Error("error pinging database", "error", err)
		os.Exit(1)
	}
	logger.Info("database connection pool established")

	// ------------------------------
	// Start application
	// ------------------------------
	app := application.New(cfg, logger, db)
	if err := app.Serve(); err != nil {
		logger.Error("error starting server", "error", err)
		os.Exit(1)
	}
}
