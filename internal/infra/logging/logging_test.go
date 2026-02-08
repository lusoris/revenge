package logging

import (
	"bytes"
	"io"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewLogger_DefaultConfig(t *testing.T) {
	cfg := Config{}
	logger := NewLogger(cfg)

	assert.NotNil(t, logger)
}

func TestNewLogger_DevelopmentMode(t *testing.T) {
	var buf bytes.Buffer
	cfg := Config{
		Level:       "debug",
		Format:      "text",
		Development: true,
		Output:      &buf,
	}

	logger := NewLogger(cfg)
	require.NotNil(t, logger)

	logger.Info("test message", slog.String("key", "value"))
	output := buf.String()

	// Tint handler produces colorized output
	assert.Contains(t, output, "test message")
}

func TestNewLogger_ProductionMode(t *testing.T) {
	var buf bytes.Buffer
	cfg := Config{
		Level:       "info",
		Format:      "json",
		Development: false,
		Output:      &buf,
	}

	logger := NewLogger(cfg)
	require.NotNil(t, logger)

	logger.Info("test message", slog.String("key", "value"))
	output := buf.String()

	// JSON format should contain JSON-like structure
	assert.Contains(t, output, "test message")
	assert.Contains(t, output, "key")
	assert.Contains(t, output, "value")
}

func TestNewLogger_LogLevels(t *testing.T) {
	tests := []struct {
		level       string
		logMessage  string
		shouldShow  bool
		logFunc     func(*slog.Logger)
	}{
		{"debug", "debug msg", true, func(l *slog.Logger) { l.Debug("debug msg") }},
		{"info", "info msg", true, func(l *slog.Logger) { l.Info("info msg") }},
		{"warn", "warn msg", true, func(l *slog.Logger) { l.Warn("warn msg") }},
		{"error", "error msg", true, func(l *slog.Logger) { l.Error("error msg") }},
		{"info", "debug filtered", false, func(l *slog.Logger) { l.Debug("debug filtered") }},
		{"warn", "info filtered", false, func(l *slog.Logger) { l.Info("info filtered") }},
		{"error", "warn filtered", false, func(l *slog.Logger) { l.Warn("warn filtered") }},
	}

	for _, tt := range tests {
		t.Run(tt.level+"_"+tt.logMessage, func(t *testing.T) {
			var buf bytes.Buffer
			cfg := Config{
				Level:  tt.level,
				Format: "text",
				Output: &buf,
			}

			logger := NewLogger(cfg)
			tt.logFunc(logger)
			output := buf.String()

			if tt.shouldShow {
				assert.Contains(t, output, tt.logMessage)
			} else {
				assert.NotContains(t, output, tt.logMessage)
			}
		})
	}
}

func TestParseLevel(t *testing.T) {
	tests := []struct {
		input    string
		expected slog.Level
	}{
		{"debug", slog.LevelDebug},
		{"DEBUG", slog.LevelDebug},
		{"info", slog.LevelInfo},
		{"INFO", slog.LevelInfo},
		{"warn", slog.LevelWarn},
		{"WARN", slog.LevelWarn},
		{"warning", slog.LevelWarn},
		{"error", slog.LevelError},
		{"ERROR", slog.LevelError},
		{"invalid", slog.LevelInfo},
		{"", slog.LevelInfo},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := parseLevel(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNewTestLogger(t *testing.T) {
	logger := NewTestLogger()

	assert.NotNil(t, logger)

	// Should not panic when logging
	logger.Debug("debug message")
	logger.Info("info message")
	logger.Warn("warn message")
	logger.Error("error message")
}

func TestConfig_Defaults(t *testing.T) {
	cfg := Config{}

	// Verify zero values
	assert.Equal(t, "", cfg.Level)
	assert.Equal(t, "", cfg.Format)
	assert.False(t, cfg.Development)
	assert.Nil(t, cfg.Output)
}

func TestNewLogger_OutputToDiscard(t *testing.T) {
	cfg := Config{
		Level:  "debug",
		Output: io.Discard,
	}

	logger := NewLogger(cfg)
	require.NotNil(t, logger)

	// Should not panic
	logger.Info("discarded message")
}
