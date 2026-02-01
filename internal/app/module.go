// Package app provides the main application module that wires all dependencies together.
package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/lusoris/revenge/internal/config"
	"github.com/lusoris/revenge/internal/infra/database"
	"github.com/lusoris/revenge/internal/infra/health"
	"github.com/lusoris/revenge/internal/infra/logging"
	"go.uber.org/fx"
)

// Module is the main application module that includes all sub-modules.
var Module = fx.Module("app",
	// Configuration
	config.Module,

	// Infrastructure
	logging.Module,
	database.Module,
	health.Module,

	// HTTP Server
	fx.Provide(NewHTTPServer),
	fx.Invoke(RegisterHTTPServer),
)

// NewHTTPServer creates a new HTTP server with health endpoints.
func NewHTTPServer(
	cfg *config.Config,
	logger *slog.Logger,
	healthService *health.Service,
) *http.Server {
	mux := http.NewServeMux()

	// Health endpoints
	mux.HandleFunc("/health/live", func(w http.ResponseWriter, r *http.Request) {
		result := healthService.Liveness(r.Context())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		// Simple JSON response (will be replaced by ogen later)
		fmt.Fprintf(w, `{"name":"%s","status":"%s","message":"%s"}`,
			result.Name, result.Status, result.Message)
	})

	mux.HandleFunc("/health/ready", func(w http.ResponseWriter, r *http.Request) {
		result := healthService.Readiness(r.Context())
		w.Header().Set("Content-Type", "application/json")
		if result.Status != health.StatusHealthy {
			w.WriteHeader(http.StatusServiceUnavailable)
		} else {
			w.WriteHeader(http.StatusOK)
		}
		fmt.Fprintf(w, `{"name":"%s","status":"%s","message":"%s"}`,
			result.Name, result.Status, result.Message)
	})

	mux.HandleFunc("/health/startup", func(w http.ResponseWriter, r *http.Request) {
		result := healthService.Startup(r.Context())
		w.Header().Set("Content-Type", "application/json")
		if result.Status != health.StatusHealthy {
			w.WriteHeader(http.StatusServiceUnavailable)
		} else {
			w.WriteHeader(http.StatusOK)
		}
		fmt.Fprintf(w, `{"name":"%s","status":"%s","message":"%s"}`,
			result.Name, result.Status, result.Message)
	})

	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler:      mux,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	return server
}

// RegisterHTTPServer registers lifecycle hooks for the HTTP server.
func RegisterHTTPServer(lc fx.Lifecycle, server *http.Server, cfg *config.Config, logger *slog.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("starting HTTP server",
				slog.String("addr", server.Addr),
			)

			// Start server in background
			go func() {
				if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					logger.Error("HTTP server error", slog.String("error", err.Error()))
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("stopping HTTP server")

			// Create shutdown context with timeout
			shutdownCtx, cancel := context.WithTimeout(ctx, cfg.Server.ShutdownTimeout)
			defer cancel()

			if err := server.Shutdown(shutdownCtx); err != nil {
				logger.Error("HTTP server shutdown error", slog.String("error", err.Error()))
				return err
			}

			logger.Info("HTTP server stopped")
			return nil
		},
	})
}
