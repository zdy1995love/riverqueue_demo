package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
	"zdy/worker/addone"
	"zdy/worker/addthree"
	"zdy/worker/multiplytwo"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
)

// Config represents the application configuration
type Config struct {
	RiverDatabaseURL string  `json:"river_database_url"`
	RiverMaxWorkers  float64 `json:"river_max_workers"`
	RiverTestOnly    bool    `json:"river_test_only"`
}

// loadConfig loads configuration from config file
func loadConfig() (*Config, error) {
	data, err := os.ReadFile("setting/config_DEV.jsonc")
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &config, nil
}

func main() {
	ctx := context.Background()

	// Create logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	logger.Info("🚀 Starting River Queue Demo Project")

	// Load configuration
	config, err := loadConfig()
	if err != nil {
		logger.Error("Failed to load config", "error", err)
		panic(err)
	}

	logger.Info("Configuration loaded successfully",
		"database_url", config.RiverDatabaseURL,
		"max_workers", config.RiverMaxWorkers,
		"test_only", config.RiverTestOnly,
	)

	// Connect to database
	dbPool, err := pgxpool.New(ctx, config.RiverDatabaseURL)
	if err != nil {
		logger.Error("Failed to connect to database", "error", err)
		panic(err)
	}
	defer dbPool.Close()

	logger.Info("✅ Database connected successfully")

	// Create and register workers (MultiplyTwoWorker without client for now)
	workers := river.NewWorkers()

	// Register AddOneWorker
	addOneWorker := addone.NewWorker(logger.With("worker", "add_one"))
	river.AddWorker(workers, addOneWorker)

	// Register MultiplyTwoWorker (will set client later)
	multiplyTwoWorker := multiplytwo.NewWorker(logger.With("worker", "multiply_two"), nil)
	river.AddWorker(workers, multiplyTwoWorker)

	// Register AddThreeWorker
	addThreeWorker := addthree.NewWorker(logger.With("worker", "add_three"))
	river.AddWorker(workers, addThreeWorker)

	logger.Info("✅ Workers registered successfully")

	// Create River client
	riverClient, err := river.NewClient(riverpgxv5.New(dbPool), &river.Config{
		Logger: logger,
		Queues: map[string]river.QueueConfig{
			river.QueueDefault: {MaxWorkers: int(config.RiverMaxWorkers)},
		},
		TestOnly:    config.RiverTestOnly,
		Workers:     workers,
		JobTimeout:  time.Minute * 5,
		MaxAttempts: 3,
	})
	if err != nil {
		logger.Error("Failed to create River client", "error", err)
		panic(err)
	}

	// Set MultiplyTwoWorker's riverClient to support chained tasks
	multiplyTwoWorker.SetClient(riverClient)

	logger.Info("✅ River client created successfully")

	// Start River
	if err := riverClient.Start(ctx); err != nil {
		logger.Error("Failed to start River", "error", err)
		panic(err)
	}

	logger.Info("✅ River started successfully")
	logger.Info("⏳ River Queue is running, waiting for tasks...")
	logger.Info("💡 Tip: Insert tasks using other programs")
	logger.Info("   - AddOne: kind='add_one', args={'number': N}")
	logger.Info("   - MultiplyTwo: kind='multiply_two', args={'number': N}")
	logger.Info("   - AddThree: kind='add_three', args={'number': N}")
	logger.Info("🛑 Press Ctrl+C to stop")

	// Graceful shutdown
	sigintOrTerm := make(chan os.Signal, 1)
	signal.Notify(sigintOrTerm, syscall.SIGINT, syscall.SIGTERM)

	<-sigintOrTerm
	logger.Info("⏹️  Received shutdown signal, stopping...")

	stopCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := riverClient.Stop(stopCtx); err != nil {
		logger.Error("Failed to stop River", "error", err)
	} else {
		logger.Info("✅ River stopped successfully")
	}
}
