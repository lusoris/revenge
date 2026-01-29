package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
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
	"github.com/lusoris/revenge/pkg/graceful"
	"github.com/lusoris/revenge/pkg/health"
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
			NewHealthChecker,
			NewShutdowner,
		),

		// Infrastructure modules
		database.Module,
		cache.Module,
		search.Module,
		jobs.Module,

		// Service modules
		auth.Module,
		user.Module,
		library.Module,
		rating.Module,
		oidc.Module,
		genre.Module,
		playback.Module,

		// API modules
		fx.Provide(
			middleware.NewAuth,
			handlers.NewAuthHandler,
			handlers.NewUserHandler,
			handlers.NewLibraryHandler,
			handlers.NewRatingHandler,
		),

		// HTTP modules
		fx.Provide(
			NewMux,
			NewServer,
		),
		fx.Invoke(RegisterRoutes),
		fx.Invoke(RegisterHealthChecks),
		fx.Invoke(StartShutdowner),
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

	logger.Info("Revenge Go starting",
		slog.String("version", Version),
		slog.String("build_time", BuildTime),
		slog.String("git_commit", GitCommit),
	)

	return logger
}

// NewHealthChecker creates a health checker with proper configuration
func NewHealthChecker(logger *slog.Logger) *health.Checker {
	return health.NewChecker(logger)
}

// NewShutdowner creates a graceful shutdown handler
func NewShutdowner(cfg *config.Config, logger *slog.Logger) *graceful.Shutdowner {
	shutdownCfg := graceful.DefaultShutdownConfig()
	if cfg.Server.ShutdownTimeout > 0 {
		shutdownCfg.Timeout = time.Duration(cfg.Server.ShutdownTimeout) * time.Second
	}
	return graceful.NewShutdowner(shutdownCfg, logger)
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
	checker *health.Checker,
	authMiddleware *middleware.Auth,
	authHandler *handlers.AuthHandler,
	userHandler *handlers.UserHandler,
	libraryHandler *handlers.LibraryHandler,
	ratingHandler *handlers.RatingHandler,
) {
	// Health check endpoints using pkg/health (Go 1.22+ pattern matching)
	mux.HandleFunc("GET /health/live", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK")) //nolint:errcheck // best-effort write
	})

	// Readiness check - comprehensive health check
	mux.HandleFunc("GET /health/ready", func(w http.ResponseWriter, r *http.Request) {
		status := checker.Check(r.Context())

		w.Header().Set("Content-Type", "application/json")

		if status.Status == health.StatusHealthy {
			w.WriteHeader(http.StatusOK)
		} else if status.Status == health.StatusDegraded {
			w.WriteHeader(http.StatusOK) // Still operational but degraded
		} else {
			w.WriteHeader(http.StatusServiceUnavailable)
		}

		_ = json.NewEncoder(w).Encode(status) //nolint:errcheck // best-effort encode
	})

	// Detailed health status endpoint
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		status := checker.Check(r.Context())
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(status) //nolint:errcheck // best-effort encode
	})

	// Database stats endpoint
	mux.HandleFunc("GET /health/db", func(w http.ResponseWriter, r *http.Request) {
		stats := database.Stats(pool)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(stats) //nolint:errcheck // best-effort encode
	})

	// Version endpoint with structured response
	mux.HandleFunc("GET /version", func(w http.ResponseWriter, r *http.Request) {
		version := map[string]string{
			"version":    Version,
			"build_time": BuildTime,
			"git_commit": GitCommit,
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(version) //nolint:errcheck // best-effort encode
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

	logger.Info("Routes registered",
		slog.Int("auth_routes", 4),
		slog.Int("user_routes", 7),
		slog.Int("library_routes", 6),
		slog.Int("rating_routes", 7),
		slog.Int("health_routes", 4),
	)
}

// RegisterHealthChecks registers all health checks with the checker
func RegisterHealthChecks(
	checker *health.Checker,
	pool *pgxpool.Pool,
	logger *slog.Logger,
) {
	// Database health check (critical)
	checker.RegisterFunc("database", health.CategoryCritical, func(ctx context.Context) error {
		return database.HealthCheck(ctx, pool)
	})

	// Cache health check (warm) - TODO: add once cache client is available
	// checker.RegisterFunc("cache", health.CategoryWarm, func(ctx context.Context) error {
	//     return cacheClient.Ping(ctx).Err()
	// })

	// Search health check (cold) - TODO: add once search client is available
	// checker.RegisterFunc("search", health.CategoryCold, func(ctx context.Context) error {
	//     return searchClient.Health(ctx)
	// })

	logger.Info("Health checks registered", slog.Int("count", 1))
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

// StartShutdowner initializes graceful shutdown handling
func StartShutdowner(lifecycle fx.Lifecycle, shutdowner *graceful.Shutdowner, logger *slog.Logger) {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			shutdowner.Start()
			logger.Info("Graceful shutdown handler started")
			return nil
		},
	})
}

// RunServer starts the HTTP server with graceful shutdown
func RunServer(lifecycle fx.Lifecycle, srv *http.Server, shutdowner *graceful.Shutdowner, logger *slog.Logger) {
	// Register server shutdown hook
	shutdowner.RegisterFunc("http_server", 100, func(ctx context.Context) error {
		logger.Info("Shutting down HTTP server")
		return srv.Shutdown(ctx)
	})

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
			// Fx handles this, but we trigger our graceful shutdown
			shutdowner.Trigger()
			return nil
		},
	})
}
