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
	adultmovie "github.com/lusoris/revenge/internal/content/c/movie"
	adultscene "github.com/lusoris/revenge/internal/content/c/scene"
	"github.com/lusoris/revenge/internal/content/movie"
	"github.com/lusoris/revenge/internal/infra/cache"
	"github.com/lusoris/revenge/internal/infra/database"
	"github.com/lusoris/revenge/internal/infra/database/db"
	infrahealth "github.com/lusoris/revenge/internal/infra/health"
	"github.com/lusoris/revenge/internal/infra/jobs"
	"github.com/lusoris/revenge/internal/infra/search"
	"github.com/lusoris/revenge/internal/service/auth"
	"github.com/lusoris/revenge/internal/service/library"
	"github.com/lusoris/revenge/internal/service/metadata/radarr"
	"github.com/lusoris/revenge/internal/service/metadata/tmdb"
	"github.com/lusoris/revenge/internal/service/oidc"
	"github.com/lusoris/revenge/internal/service/session"
	"github.com/lusoris/revenge/internal/service/user"
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
	cfg, err := config.LoadDefault()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
		os.Exit(1)
	}

	options := []fx.Option{
		// Core infrastructure
		fx.Supply(cfg),
		fx.Provide(NewLogger),
		fx.Provide(func() api.BuildInfo {
			return api.BuildInfo{
				Version:   Version,
				BuildTime: BuildTime,
				GitCommit: GitCommit,
			}
		}),
		fx.Provide(func(cfg *config.Config) tmdb.Config {
			return tmdb.Config{
				APIKey:     cfg.Metadata.TMDb.APIKey,
				BaseURL:    cfg.Metadata.TMDb.BaseURL,
				ImageURL:   cfg.Metadata.TMDb.ImageURL,
				Timeout:    time.Duration(cfg.Metadata.TMDb.Timeout) * time.Second,
				CacheTTL:   time.Duration(cfg.Metadata.TMDb.CacheTTL) * time.Second,
				CacheSize:  cfg.Metadata.TMDb.CacheSize,
				RetryCount: cfg.Metadata.TMDb.RetryCount,
			}
		}),

		// Infrastructure modules
		database.Module,
		fx.Provide(NewQueries), // sqlc queries
		cache.Module,
		jobs.Module,
		search.Module,
		infrahealth.Module,
		radarr.Module,
		tmdb.Module,

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
	}

	if cfg.Modules.Movie {
		options = append(options, movie.ModuleWithRiver)
	}
	if cfg.Modules.Adult {
		options = append(options, adultmovie.Module, adultscene.Module)
	}

	options = append(options,
		// HTTP server
		fx.Provide(
			NewMux,
			NewServer,
		),
		fx.Invoke(RegisterRoutes),
		fx.Invoke(RunServer),
	)

	app := fx.New(options...)

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
	pool *pgxpool.Pool,
	healthChecker *health.Checker,
	cfg *config.Config,
) {
	// Mount the ogen-generated API server on /api/v1
	mux.Handle("/api/v1/", http.StripPrefix("/api/v1", apiServer))

	mux.HandleFunc("GET /health/live", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	mux.HandleFunc("GET /health/ready", func(w http.ResponseWriter, r *http.Request) {
		if healthChecker.IsReady(r.Context()) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("Ready"))
			return
		}
		w.WriteHeader(http.StatusServiceUnavailable)
		_, _ = w.Write([]byte("Not Ready"))
	})

	mux.HandleFunc("GET /health/db", func(w http.ResponseWriter, r *http.Request) {
		stats := pool.Stat()
		resp := map[string]int32{
			"total_connections":  stats.TotalConns(),
			"idle_connections":   stats.IdleConns(),
			"active_connections": stats.AcquiredConns(),
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	})

	mux.HandleFunc("GET /version", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(Version))
	})

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
		ReadTimeout:       cfg.Server.ReadTimeout,
		WriteTimeout:      cfg.Server.WriteTimeout,
		IdleTimeout:       cfg.Server.IdleTimeout,
		ReadHeaderTimeout: cfg.Server.ReadHeaderTimeout,
		MaxHeaderBytes:    cfg.Server.MaxHeaderBytes,
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
