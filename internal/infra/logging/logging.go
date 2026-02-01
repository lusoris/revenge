// Package logging provides structured logging for the Revenge server.
// It uses tint for development (colorized, human-readable) and zap for production (JSON).
package logging

import (
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/lmittmann/tint"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Config holds logging configuration.
type Config struct {
	// Level is the minimum log level (debug, info, warn, error).
	Level string

	// Format is the log format (text, json).
	Format string

	// Development enables development mode.
	Development bool

	// Output is where logs are written (defaults to os.Stdout).
	Output io.Writer
}

// NewLogger creates a new slog.Logger based on the configuration.
func NewLogger(cfg Config) *slog.Logger {
	// Default to stdout if no output specified
	if cfg.Output == nil {
		cfg.Output = os.Stdout
	}

	// Parse log level
	level := parseLevel(cfg.Level)

	// Create handler based on format
	var handler slog.Handler
	if cfg.Development || cfg.Format == "text" {
		// Development mode: use tint for colorized output
		handler = tint.NewHandler(cfg.Output, &tint.Options{
			Level:      level,
			TimeFormat: "15:04:05",
			AddSource:  cfg.Development,
		})
	} else {
		// Production mode: use JSON handler
		handler = slog.NewJSONHandler(cfg.Output, &slog.HandlerOptions{
			Level:     level,
			AddSource: false,
		})
	}

	return slog.New(handler)
}

// NewZapLogger creates a zap.Logger for production use.
// This is provided for components that require a zap logger directly.
func NewZapLogger(cfg Config) (*zap.Logger, error) {
	// Parse log level
	var level zapcore.Level
	switch strings.ToLower(cfg.Level) {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn", "warning":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	default:
		level = zapcore.InfoLevel
	}

	// Create encoder config
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// Create core
	var encoder zapcore.Encoder
	if cfg.Development {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	}

	// Default to stdout if no output specified
	output := cfg.Output
	if output == nil {
		output = os.Stdout
	}

	core := zapcore.NewCore(
		encoder,
		zapcore.AddSync(output),
		level,
	)

	// Create logger
	logger := zap.New(core)

	// Add caller and stack trace in development
	if cfg.Development {
		logger = logger.WithOptions(zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	}

	return logger, nil
}

// parseLevel converts a string log level to slog.Level.
func parseLevel(levelStr string) slog.Level {
	switch strings.ToLower(levelStr) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
