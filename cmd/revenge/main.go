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

	gen "github.com/lusoris/revenge/api/generated"
	"github.com/lusoris/revenge/internal/api"
	"github.com/lusoris/revenge/internal/infra/cache"
	"github.com/lusoris/revenge/internal/infra/database"
	"github.com/lusoris/revenge/internal/infra/jobs"
	"github.com/lusoris/revenge/internal/infra/search"
	"github.com/lusoris/revenge/internal/service/activity"
	"github.com/lusoris/revenge/internal/service/apikeys"
	"github.com/lusoris/revenge/internal/service/auth"
	"github.com/lusoris/revenge/internal/service/library"
	"github.com/lusoris/revenge/internal/service/oidc"
	"github.com/lusoris/revenge/internal/service/rbac"
	"github.com/lusoris/revenge/internal/service/session"
	"github.com/lusoris/revenge/internal/service/settings"
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
			NewBuildInfo,
		),

		// Infrastructure modules
		database.Module,
		cache.Module,
		search.Module,
		jobs.Module,

		// Service modules
		auth.Module,
		user.Module,
		session.Module,
		library.Module,
		oidc.Module,
		rbac.Module,
		activity.Module,
		settings.Module,
		apikeys.Module,

		// API module (ogen-generated handlers)
		api.Module,

		// HTTP modules
		fx.Provide(
			NewMux,
			NewServer,
		),
		fx.Invoke(RegisterRoutes),
		fx.Invoke(RegisterHealthChecks),
		fx.Invoke(MountAPIServer),
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
		stats := pool.Stat()
		resp := map[string]int32{
			"total_conns":    stats.TotalConns(),
			"idle_conns":     stats.IdleConns(),
			"acquired_conns": stats.AcquiredConns(),
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp) //nolint:errcheck // best-effort encode
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

	// TODO: Add API endpoints when handlers are implemented

	logger.Info("Routes registered",
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

// NewBuildInfo creates build info for the API handlers.
func NewBuildInfo() api.BuildInfo {
	return api.BuildInfo{
		Version:   Version,
		BuildTime: BuildTime,
		GitCommit: GitCommit,
	}
}

// MountAPIServer mounts the ogen-generated API server to the mux.
func MountAPIServer(
	mux *http.ServeMux,
	handler *api.Handler,
	security *api.SecurityHandler,
	logger *slog.Logger,
) error {
	srv, err := gen.NewServer(handler, security)
	if err != nil {
		return fmt.Errorf("create API server: %w", err)
	}

	// Mount API under /api/v1
	mux.Handle("/api/v1/", http.StripPrefix("/api/v1", srv))

	logger.Info("API server mounted", slog.String("prefix", "/api/v1"))
	return nil
}
