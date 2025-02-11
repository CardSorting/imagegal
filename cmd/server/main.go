package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"image/internal/app"
	"image/internal/domain/ports"
	"image/internal/handlers/health"
	"image/internal/handlers/text2img"
	"image/internal/infrastructure/config"
	"image/internal/infrastructure/http"
	registry "image/internal/infrastructure/registry"
	"image/internal/infrastructure/validation"
	"image/internal/services/modelslab"
	"image/pkg/logger"
)

func main() {
	// Initialize logger
	appLogger := logger.New()

	// Load configuration
	cfg, err := config.New()
	if err != nil {
		appLogger.Error("Failed to load configuration", err)
		os.Exit(1)
	}

	// Initialize validator
	validator := validation.New()

	// Initialize HTTP client
	httpClient := http.NewClient(
		cfg.ModelsLab.BaseURL,
		cfg.ModelsLab.APIKey,
		appLogger,
		http.WithMaxRetries(cfg.ModelsLab.MaxRetries),
		http.WithTimeout(30*time.Second),
	)

	// Initialize model registry
	modelRegistry := registry.NewModelRegistry()

	// Initialize services
	modelsLabService := modelslab.NewService(httpClient, validator, appLogger, modelRegistry)

	// Initialize handlers
	handlers := make(map[string]ports.Handler)
	handlers["text2img"] = text2img.NewHandler(modelsLabService, appLogger)
	handlers["health"] = health.NewHandler(appLogger)

	// Create and configure server
	server := app.NewServer(cfg, appLogger, handlers)

	// Handle graceful shutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := server.Start(); err != nil {
			appLogger.Error("Server failed", err)
			os.Exit(1)
		}
	}()

	appLogger.Info("Server is running", "port", cfg.Server.Port)

	// Wait for interrupt signal
	<-done
	appLogger.Info("Shutting down server...")

	// Create shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		appLogger.Error("Server shutdown failed", err)
		os.Exit(1)
	}

	appLogger.Info("Server stopped gracefully")
}
