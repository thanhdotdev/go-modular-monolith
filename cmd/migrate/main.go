package main

import (
	"fmt"
	"os"
	"project-example/internal/platform/config"
	"project-example/internal/platform/database"
	"project-example/internal/platform/logger"

	"go.uber.org/zap"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run ./cmd/migrate [up|down]")
		os.Exit(1)
	}

	cfg := config.Load()
	log, logCloser, err := logger.New(cfg.Logging)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Close(log, logCloser)

	if !cfg.Database.Enabled() {
		log.Error("DATABASE_DSN is required")
		os.Exit(1)
	}

	err = nil

	switch os.Args[1] {
	case "up":
		err = database.MigrateUp(cfg.Database)
	case "down":
		err = database.MigrateDown(cfg.Database)
	default:
		fmt.Fprintf(os.Stderr, "unknown migration command: %s\n", os.Args[1])
		os.Exit(1)
	}

	if err != nil {
		log.Error("migration command failed", zap.String("command", os.Args[1]), zap.Error(err))
		os.Exit(1)
	}

	log.Info("migration command completed", zap.String("command", os.Args[1]))
}
