package database

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/tracelog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewQueryLogger(t *testing.T) {
	logger := slog.Default()
	slowQueryThreshold := 100 * time.Millisecond

	queryLogger := NewQueryLogger(logger, slowQueryThreshold)

	require.NotNil(t, queryLogger)
	assert.Equal(t, logger, queryLogger.logger)
	assert.Equal(t, slowQueryThreshold, queryLogger.slowQueryThreshold)
}

func TestQueryLoggerLog(t *testing.T) {
	tests := []struct {
		name  string
		level tracelog.LogLevel
		msg   string
		data  map[string]interface{}
	}{
		{
			name:  "trace level",
			level: tracelog.LogLevelTrace,
			msg:   "trace message",
			data:  map[string]interface{}{"key": "value"},
		},
		{
			name:  "debug level",
			level: tracelog.LogLevelDebug,
			msg:   "debug message",
			data:  map[string]interface{}{"key": "value"},
		},
		{
			name:  "info level",
			level: tracelog.LogLevelInfo,
			msg:   "info message",
			data:  map[string]interface{}{"key": "value"},
		},
		{
			name:  "slow query detection",
			level: tracelog.LogLevelInfo,
			msg:   "query completed",
			data: map[string]interface{}{
				"sql":  "SELECT * FROM users",
				"time": 200 * time.Millisecond,
			},
		},
		{
			name:  "fast query",
			level: tracelog.LogLevelInfo,
			msg:   "query completed",
			data: map[string]interface{}{
				"sql":  "SELECT 1",
				"time": 50 * time.Millisecond,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := slog.Default()
			queryLogger := NewQueryLogger(logger, 100*time.Millisecond)

			assert.NotPanics(t, func() {
				queryLogger.Log(context.Background(), tt.level, tt.msg, tt.data)
			})
		})
	}
}

func TestTracerConfig(t *testing.T) {
	logger := slog.Default()
	tracer, queryTracer := TracerConfig(logger, tracelog.LogLevelInfo, 100*time.Millisecond)

	require.NotNil(t, tracer)
	require.NotNil(t, queryTracer)
}

func TestExtractOperation(t *testing.T) {
	tests := []struct {
		sql  string
		want string
	}{
		{"SELECT * FROM users", "select"},
		{"INSERT INTO users VALUES (1)", "insert"},
		{"UPDATE users SET name='x'", "update"},
		{"DELETE FROM users WHERE id=1", "delete"},
		{"BEGIN", "begin"},
		{"COMMIT", "commit"},
		{"ROLLBACK", "rollback"},
		{"-- name: GetMovie :one\nSELECT * FROM movies", "select"},
		{"-- name: CreateUser :exec\nINSERT INTO users", "insert"},
		{"", "query"},
		{"UNKNOWN STATEMENT", "query"},
	}

	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			got := extractOperation(tt.sql)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestNewQueryTracer(t *testing.T) {
	logger := slog.Default()
	tracer := NewQueryTracer(logger, 100*time.Millisecond)

	require.NotNil(t, tracer)
	assert.Equal(t, logger, tracer.logger)
	assert.Equal(t, 100*time.Millisecond, tracer.slowQueryThreshold)
}

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
		want     string
	}{
		{
			name:     "seconds",
			duration: 2 * time.Second,
			want:     "2.00s",
		},
		{
			name:     "milliseconds",
			duration: 250 * time.Millisecond,
			want:     "250.00ms",
		},
		{
			name:     "microseconds",
			duration: 500 * time.Microsecond,
			want:     "500us",
		},
		{
			name:     "mixed seconds",
			duration: 1500 * time.Millisecond,
			want:     "1.50s",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatDuration(tt.duration)
			assert.Equal(t, tt.want, got)
		})
	}
}
