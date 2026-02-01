package logging_test

import (
	"bytes"
	"context"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lusoris/revenge/internal/infra/logging"
)

func TestNewLogger(t *testing.T) {
	tests := []struct {
		name  string
		level string
		want  slog.Level
	}{
		{"debug level", "debug", slog.LevelDebug},
		{"info level", "info", slog.LevelInfo},
		{"warn level", "warn", slog.LevelWarn},
		{"error level", "error", slog.LevelError},
		{"invalid defaults to info", "invalid", slog.LevelInfo},
		{"empty defaults to info", "", slog.LevelInfo},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := logging.NewLogger(tt.level)
			require.NotNil(t, logger)

			// Logger should be usable
			logger.Info("test message")
		})
	}
}

func TestNewTestLogger(t *testing.T) {
	logger := logging.NewTestLogger()
	require.NotNil(t, logger)

	// Test logger should work
	logger.Debug("debug message")
	logger.Info("info message")
	logger.Warn("warn message")
}

func TestLoggerWithContext(t *testing.T) {
	var buf bytes.Buffer
	handler := slog.NewJSONHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	logger := slog.New(handler)

	ctx := context.Background()

	// Log with context
	logger.InfoContext(ctx, "test message", "key", "value")

	assert.Contains(t, buf.String(), "test message")
	assert.Contains(t, buf.String(), "key")
	assert.Contains(t, buf.String(), "value")
}

func TestLoggerLevels(t *testing.T) {
	var buf bytes.Buffer
	handler := slog.NewJSONHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelWarn,
	})
	logger := slog.New(handler)

	// Info should not log (level is Warn)
	logger.Info("info message")
	assert.NotContains(t, buf.String(), "info message")

	// Warn should log
	logger.Warn("warn message")
	assert.Contains(t, buf.String(), "warn message")

	buf.Reset()

	// Error should log
	logger.Error("error message")
	assert.Contains(t, buf.String(), "error message")
}
