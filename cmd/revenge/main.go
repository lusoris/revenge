package main

import (
	"context"
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
	"github.com/lusoris/revenge/internal/config"
	"github.com/lusoris/revenge/internal/infra/cache"
	"github.com/lusoris/revenge/internal/infra/database"
	"github.com/lusoris/revenge/internal/infra/database/db"
	infrahealth "github.com/lusoris/revenge/internal/infra/health"
	"github.com/lusoris/revenge/internal/infra/jobs"
	"github.com/lusoris/revenge/internal/infra/search"
	"github.com/lusoris/revenge/internal/service/auth"
	"github.com/lusoris/revenge/internal/service/library"
	"github.com/lusoris/revenge/internal/service/oidc"
	"github.com/lusoris/revenge/internal/service/session"
	"github.com/lusoris/revenge/internal/service/user"
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
	app := fx.New(
		// Core infrastructure
		config.Module,
		fx.Provide(NewLogger),

		// Infrastructure modules
		database.Module,
		fx.Provide(NewQueries), // sqlc queries
		cache.Module,
		jobs.Module,
		search.Module,
		infrahealth.Module,

		// Shared services
		user.Module,
		session.Module,
		auth.Module,
		library.Module,
		oidc.Module,

		// API handlers
		api.Module,
		fx.Provide(NewAPIServer),

		// Content modules will be added here as they're implemented
		// movie.Module,
		// tvshow.Module,
		// music.Module,

		// HTTP server
		fx.Provide(
			NewMux,
			NewServer,
		),
		fx.Invoke(RegisterRoutes),
		fx.Invoke(RunServer),
	)

	app.Run()
}

// NewLogger creates a new structured logger using slog with tint for console output.
func NewLogger(cfg *config.Config) *slog.Logger {
	var handler slog.Handler

	level := parseLogLevel(cfg.Logging.Level)

	if cfg.Logging.Format == "json" {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: level,
		})
	} else {
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

// NewQueries creates sqlc queries instance.
func NewQueries(pool *pgxpool.Pool) *db.Queries {
	return db.New(pool)
}

// NewAPIServer creates the ogen-generated API server.
func NewAPIServer(handler *api.Handler, securityHandler *api.SecurityHandler) (*gen.Server, error) {
	return gen.NewServer(handler, securityHandler)
}

// NewMux creates a new HTTP router using Go 1.22+ enhanced ServeMux.
func NewMux(logger *slog.Logger) *http.ServeMux {
	mux := http.NewServeMux()
	logger.Info("HTTP router initialized")
	return mux
}

// RegisterRoutes registers all HTTP routes.
func RegisterRoutes(
	mux *http.ServeMux,
	logger *slog.Logger,
	apiServer *gen.Server,
) {
	// Mount the ogen-generated API server on /api/v1
	mux.Handle("/api/v1/", http.StripPrefix("/api/v1", apiServer))

	logger.Info("Routes registered",
		slog.String("api_prefix", "/api/v1"),
	)
}

// NewServer creates a new HTTP server.
func NewServer(mux *http.ServeMux, cfg *config.Config, logger *slog.Logger) *http.Server {
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)

	srv := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}

	logger.Info("Server configured", slog.String("address", addr))
	return srv
}

// RunServer starts the HTTP server with graceful shutdown.
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
}
