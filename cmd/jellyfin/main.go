package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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
	// Initialize configuration
	initConfig()

	// Create Fx application
	app := fx.New(
		fx.Provide(
			NewLogger,
			NewRouter,
			NewServer,
		),
		fx.Invoke(RunServer),
	)

	// Start the application
	app.Run()
}

// initConfig initializes the application configuration
func initConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath(".")

	// Set defaults
	viper.SetDefault("server.port", 8096)
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("log.level", "info")
	viper.SetDefault("db.type", "sqlite")
	viper.SetDefault("db.path", "./data/jellyfin.db")

	// Environment variables
	viper.SetEnvPrefix("JELLYFIN")
	viper.AutomaticEnv()

	// Read config file (optional)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Fatalf("Error reading config file: %v", err)
		}
	}
}

// NewLogger creates a new zap logger
func NewLogger() (*zap.Logger, error) {
	logLevel := viper.GetString("log.level")

	config := zap.NewProductionConfig()
	config.Level = zap.NewAtomicLevelAt(parseLogLevel(logLevel))

	logger, err := config.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to create logger: %w", err)
	}

	logger.Info("Jellyfin Go starting",
		zap.String("version", Version),
		zap.String("build_time", BuildTime),
		zap.String("git_commit", GitCommit),
	)

	return logger, nil
}

// parseLogLevel converts string log level to zap level
func parseLogLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zap.DebugLevel
	case "info":
		return zap.InfoLevel
	case "warn":
		return zap.WarnLevel
	case "error":
		return zap.ErrorLevel
	default:
		return zap.InfoLevel
	}
}

// NewRouter creates a new HTTP router
func NewRouter(logger *zap.Logger) *mux.Router {
	r := mux.NewRouter()

	// Health check endpoints
	r.HandleFunc("/health/live", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	r.HandleFunc("/health/ready", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Ready"))
	}).Methods("GET")

	// Version endpoint
	r.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"version":"%s","build_time":"%s","git_commit":"%s"}`,
			Version, BuildTime, GitCommit)
	}).Methods("GET")

	logger.Info("Router initialized")
	return r
}

// NewServer creates a new HTTP server
func NewServer(router *mux.Router, logger *zap.Logger) *http.Server {
	host := viper.GetString("server.host")
	port := viper.GetInt("server.port")
	addr := fmt.Sprintf("%s:%d", host, port)

	srv := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	logger.Info("Server configured", zap.String("address", addr))
	return srv
}

// RunServer starts the HTTP server
func RunServer(lifecycle fx.Lifecycle, srv *http.Server, logger *zap.Logger) {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				logger.Info("Starting HTTP server", zap.String("address", srv.Addr))
				if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					logger.Fatal("Failed to start server", zap.Error(err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Shutting down HTTP server")

			shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			if err := srv.Shutdown(shutdownCtx); err != nil {
				logger.Error("Server shutdown error", zap.Error(err))
				return err
			}

			logger.Info("Server stopped gracefully")
			return nil
		},
	})

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		logger.Info("Received shutdown signal")
	}()
}
