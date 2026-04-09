package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"project-example/internal/app"
	"project-example/internal/platform/config"
	"project-example/internal/platform/logger"
	"syscall"

	"go.uber.org/zap"
)

func main() {
	cfg := config.Load()
	log, logCloser, err := logger.New(cfg.Logging)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to initialize logger: %v\n", err)
		os.Exit(1)
	}

	defer logger.Close(log, logCloser)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	application, err := app.New(cfg, log)
	if err != nil {
		log.Error("failed to bootstrap application", zap.Error(err))
		os.Exit(1)
	}

	if err := application.Run(ctx); err != nil {
		log.Error("application stopped with error", zap.Error(err))
		os.Exit(1)
	}
}
