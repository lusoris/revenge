// Package main is the entry point for the Revenge media server.
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/lusoris/revenge/internal/app"
	"github.com/lusoris/revenge/internal/version"
	"go.uber.org/fx"
)

func main() {
	// Check for version flag
	if len(os.Args) > 1 && os.Args[1] == "version" {
		fmt.Println("revenge", version.Info())
		return
	}

	// Check for migrate command
	if len(os.Args) > 1 && os.Args[1] == "migrate" {
		runMigrate()
		return
	}

	// Create fx application
	application := fx.New(
		app.Module,
		fx.NopLogger, // Suppress fx logs (we use slog)
	)

	// Handle signals for graceful shutdown
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// Start application
	if err := application.Start(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to start application: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Revenge server started. Press Ctrl+C to stop.")

	// Wait for interrupt signal
	<-ctx.Done()

	fmt.Println("\nShutting down gracefully...")

	// Stop application
	stopCtx, stopCancel := context.WithTimeout(context.Background(), application.StopTimeout())
	defer stopCancel()

	if err := application.Stop(stopCtx); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to stop application: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Server stopped")
}

// runMigrate is implemented in migrate.go
