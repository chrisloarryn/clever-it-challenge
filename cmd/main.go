package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"beers-challenge/internal/infrastructure/dependencies"
)

func main() {
	// Create dependency injection container
	container, err := dependencies.NewContainer()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create container: %v\n", err)
		os.Exit(1)
	}
	defer func() {
		if err := container.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "Error closing container: %v\n", err)
		}
	}()

	logger := container.GetLogger()
	logger.Info(context.Background(), "Starting Beer API", map[string]interface{}{
		"version": "1.0.0",
	})

	// Get HTTP server
	server := container.GetHTTPServer()

	// Start server in a goroutine
	go func() {
		if err := server.Start(); err != nil {
			logger.Error(context.Background(), "Server failed to start", err, nil)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info(context.Background(), "Shutting down server", nil)

	// Create a context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := server.Stop(ctx); err != nil {
		logger.Error(ctx, "Server forced to shutdown", err, nil)
		os.Exit(1)
	}

	logger.Info(context.Background(), "Server shutdown completed", nil)
}
