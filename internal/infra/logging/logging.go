// Package logging provides structured logging for the Revenge server.
// It uses tint for development (colorized, human-readable) and slog JSON handler for production.
package logging

import (
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/lmittmann/tint"
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

// NewTestLogger creates a logger suitable for testing.
// It writes to a discarding writer and uses debug level.
func NewTestLogger() *slog.Logger {
	return NewLogger(Config{
		Level:       "debug",
		Format:      "text",
		Development: true,
		Output:      io.Discard,
	})
}
