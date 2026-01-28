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
	"github.com/lmittmann/tint"
	"go.uber.org/fx"

	"github.com/lusoris/revenge/internal/api/handlers"
	"github.com/lusoris/revenge/internal/api/middleware"
	"github.com/lusoris/revenge/internal/infra/cache"
	"github.com/lusoris/revenge/internal/infra/database"
	"github.com/lusoris/revenge/internal/infra/jobs"
	"github.com/lusoris/revenge/internal/infra/search"
	"github.com/lusoris/revenge/internal/service/auth"
	"github.com/lusoris/revenge/internal/service/genre"
	"github.com/lusoris/revenge/internal/service/library"
	"github.com/lusoris/revenge/internal/service/oidc"
	"github.com/lusoris/revenge/internal/service/playback"
	"github.com/lusoris/revenge/internal/service/rating"
	"github.com/lusoris/revenge/internal/service/user"
	"github.com/lusoris/revenge/pkg/config"
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
		cache.Module,
		jobs.Module,
		search.Module,

		// Service modules
		auth.Module,
		user.Module,
		library.Module,
		rating.Module,
		genre.Module,
		oidc.Module,
		playback.Module,

		// API modules
		fx.Provide(
			middleware.NewAuth,
			handlers.NewAuthHandler,
			handlers.NewUserHandler,
			handlers.NewLibraryHandler,
			handlers.NewRatingHandler,
			handlers.NewGenreHandler,
			handlers.NewOIDCHandler,
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

	logger.Info("Revenge starting",
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
	cacheClient *cache.Client,
	searchClient *search.Client,
	authMiddleware *middleware.Auth,
	authHandler *handlers.AuthHandler,
	userHandler *handlers.UserHandler,
	libraryHandler *handlers.LibraryHandler,
	ratingHandler *handlers.RatingHandler,
	genreHandler *handlers.GenreHandler,
	oidcHandler *handlers.OIDCHandler,
) {
	// Health check endpoints (Go 1.22+ pattern matching)
	mux.HandleFunc("GET /health/live", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	// Readiness check - verifies database and cache connectivity
	mux.HandleFunc("GET /health/ready", func(w http.ResponseWriter, r *http.Request) {
		// Check database
		if err := database.HealthCheck(r.Context(), pool); err != nil {
			logger.Error("Database readiness check failed", slog.Any("error", err))
			http.Error(w, "Database not ready", http.StatusServiceUnavailable)
			return
		}

		// Check cache
		if err := cacheClient.Ping(r.Context()); err != nil {
			logger.Warn("Cache readiness check failed", slog.Any("error", err))
			// Cache failure is not critical, just log warning
		}

		// Check search
		if err := searchClient.Health(r.Context()); err != nil {
			logger.Warn("Search readiness check failed", slog.Any("error", err))
			// Search failure is not critical, just log warning
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("Ready"))
	})

	// Full health check with detailed status
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		status := map[string]any{
			"status":  "healthy",
			"version": Version,
			"checks":  map[string]string{},
		}
		checks := status["checks"].(map[string]string)

		// Database
		if err := database.HealthCheck(r.Context(), pool); err != nil {
			checks["database"] = "unhealthy: " + err.Error()
			status["status"] = "degraded"
		} else {
			checks["database"] = "healthy"
		}

		// Cache
		if err := cacheClient.Ping(r.Context()); err != nil {
			checks["cache"] = "unhealthy: " + err.Error()
			status["status"] = "degraded"
		} else {
			checks["cache"] = "healthy"
		}

		// Search
		if err := searchClient.Health(r.Context()); err != nil {
			checks["search"] = "unhealthy: " + err.Error()
			status["status"] = "degraded"
		} else {
			checks["search"] = "healthy"
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(status)
	})

	// Database stats endpoint
	mux.HandleFunc("GET /health/db", func(w http.ResponseWriter, r *http.Request) {
		stats := database.Stats(pool)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(stats)
	})

	// Version endpoint with structured response
	mux.HandleFunc("GET /version", func(w http.ResponseWriter, r *http.Request) {
		version := map[string]string{
			"version":    Version,
			"build_time": BuildTime,
			"git_commit": GitCommit,
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(version)
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

	// Library endpoints
	libraryHandler.RegisterRoutes(mux, authMiddleware)

	// Rating endpoints
	ratingHandler.RegisterRoutes(mux, authMiddleware)

	// Genre endpoints
	genreHandler.RegisterRoutes(mux, authMiddleware)

	// OIDC endpoints
	oidcHandler.RegisterRoutes(mux, authMiddleware)

	logger.Info("Routes registered",
		slog.Int("auth_routes", 4),
		slog.Int("user_routes", 7),
		slog.Int("library_routes", 6),
		slog.Int("rating_routes", 7),
		slog.Int("genre_routes", 4),
		slog.Int("oidc_routes", 4),
		slog.Int("health_routes", 4),
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
