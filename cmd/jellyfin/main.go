package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jellyfin/jellyfin-go/internal/api/handlers"
	"github.com/jellyfin/jellyfin-go/internal/api/middleware"
	"github.com/jellyfin/jellyfin-go/internal/infra/database"
	"github.com/jellyfin/jellyfin-go/internal/service/auth"
	"github.com/jellyfin/jellyfin-go/internal/service/user"
	"github.com/jellyfin/jellyfin-go/pkg/config"
	"github.com/lmittmann/tint"
	"go.uber.org/fx"
)

var (
	// Version is set at build time
	Version = "dev"
	// BuildTime is set at build time
	BuildTime = "unknown"
	// GitCommit is set at build time
	GitCommit = "unknown"
)

func main() {
	// Create Fx application with modern dependency injection
	app := fx.New(
		// Core modules
		fx.Provide(
			config.New,
			NewLogger,
		),

		// Infrastructure modules
		database.Module,

		// Service modules
		auth.Module,
		user.Module,

		// API modules
		fx.Provide(
			middleware.NewAuth,
			handlers.NewAuthHandler,
			handlers.NewUserHandler,
		),

		// HTTP modules
		fx.Provide(
			NewMux,
			NewServer,
		),
		fx.Invoke(RegisterRoutes),
		fx.Invoke(RunServer),
	)

	// Start the application
	app.Run()
}

// NewLogger creates a new structured logger using slog (Go 1.21+)
func NewLogger(cfg *config.Config) *slog.Logger {
	var handler slog.Handler

	level := parseLogLevel(cfg.Log.Level)

	if cfg.Log.Format == "json" {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: level,
		})
	} else {
		// Use tint for beautiful colored console output
		handler = tint.NewHandler(os.Stdout, &tint.Options{
			Level:      level,
			TimeFormat: time.DateTime,
		})
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)

	logger.Info("Jellyfin Go starting",
		slog.String("version", Version),
		slog.String("build_time", BuildTime),
		slog.String("git_commit", GitCommit),
	)

	return logger
}

// parseLogLevel converts string to slog.Level
func parseLogLevel(level string) slog.Level {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// NewMux creates a new HTTP router using Go 1.22+ enhanced ServeMux
func NewMux(logger *slog.Logger) *http.ServeMux {
	mux := http.NewServeMux()
	logger.Info("HTTP router initialized")
	return mux
}

// RegisterRoutes registers all HTTP routes
func RegisterRoutes(
	mux *http.ServeMux,
	logger *slog.Logger,
	pool *pgxpool.Pool,
	authMiddleware *middleware.Auth,
	authHandler *handlers.AuthHandler,
	userHandler *handlers.UserHandler,
) {
	// Health check endpoints (Go 1.22+ pattern matching)
	mux.HandleFunc("GET /health/live", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Readiness check - verifies database connectivity
	mux.HandleFunc("GET /health/ready", func(w http.ResponseWriter, r *http.Request) {
		if err := database.HealthCheck(r.Context(), pool); err != nil {
			logger.Error("Readiness check failed", slog.Any("error", err))
			http.Error(w, "Database not ready", http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Ready"))
	})

	// Database stats endpoint
	mux.HandleFunc("GET /health/db", func(w http.ResponseWriter, r *http.Request) {
		stats := database.Stats(pool)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(stats)
	})

	// Version endpoint with structured response
	mux.HandleFunc("GET /version", func(w http.ResponseWriter, r *http.Request) {
		version := map[string]string{
			"version":    Version,
			"build_time": BuildTime,
			"git_commit": GitCommit,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(version)
	})

	// Auth endpoints
	mux.HandleFunc("POST /Users/AuthenticateByName", authHandler.Login)
	mux.Handle("POST /Sessions/Logout", authMiddleware.Required(http.HandlerFunc(authHandler.Logout)))
	mux.Handle("POST /Auth/Refresh", authMiddleware.Required(http.HandlerFunc(authHandler.RefreshToken)))
	mux.Handle("POST /Users/{userId}/Password", authMiddleware.Required(http.HandlerFunc(authHandler.ChangePassword)))

	// User endpoints
	mux.Handle("GET /Users/Me", authMiddleware.Required(http.HandlerFunc(userHandler.GetCurrentUser)))
	mux.Handle("GET /Users", authMiddleware.Required(http.HandlerFunc(userHandler.ListUsers)))
	mux.Handle("GET /Users/{userId}", authMiddleware.Required(http.HandlerFunc(userHandler.GetUser)))
	mux.Handle("POST /Users/New", authMiddleware.AdminRequired(http.HandlerFunc(userHandler.CreateUser)))
	mux.Handle("POST /Users", authMiddleware.AdminRequired(http.HandlerFunc(userHandler.CreateUser)))
	mux.Handle("POST /Users/{userId}", authMiddleware.Required(http.HandlerFunc(userHandler.UpdateUser)))
	mux.Handle("DELETE /Users/{userId}", authMiddleware.AdminRequired(http.HandlerFunc(userHandler.DeleteUser)))

	logger.Info("Routes registered",
		slog.Int("auth_routes", 4),
		slog.Int("user_routes", 7),
		slog.Int("health_routes", 3),
	)
}

// NewServer creates a new HTTP server with modern settings
func NewServer(mux *http.ServeMux, cfg *config.Config, logger *slog.Logger) *http.Server {
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)

	srv := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		MaxHeaderBytes:    1 << 20, // 1 MB
	}

	logger.Info("Server configured", slog.String("address", addr))
	return srv
}

// RunServer starts the HTTP server with graceful shutdown
func RunServer(lifecycle fx.Lifecycle, srv *http.Server, logger *slog.Logger) {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				logger.Info("Starting HTTP server", slog.String("address", srv.Addr))
				if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
					logger.Error("Server error", slog.Any("error", err))
					os.Exit(1)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Shutting down HTTP server")

			shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			if err := srv.Shutdown(shutdownCtx); err != nil {
				logger.Error("Server shutdown error", slog.Any("error", err))
				return err
			}

			logger.Info("Server stopped gracefully")
			return nil
		},
	})

	// Setup signal handling
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		logger.Info("Received shutdown signal")
	}()
}
