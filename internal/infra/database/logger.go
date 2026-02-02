package database

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/tracelog"
)

// QueryLogger wraps slog.Logger for pgx query logging.
type QueryLogger struct {
	logger             *slog.Logger
	slowQueryThreshold time.Duration
}

// NewQueryLogger creates a new QueryLogger.
func NewQueryLogger(logger *slog.Logger, slowQueryThreshold time.Duration) *QueryLogger {
	return &QueryLogger{
		logger:             logger,
		slowQueryThreshold: slowQueryThreshold,
	}
}

// Log implements pgx tracelog.Logger interface.
func (l *QueryLogger) Log(ctx context.Context, level tracelog.LogLevel, msg string, data map[string]interface{}) {
	var slogLevel slog.Level
	switch level {
	case tracelog.LogLevelTrace:
		slogLevel = slog.LevelDebug - 1
	case tracelog.LogLevelDebug:
		slogLevel = slog.LevelDebug
	case tracelog.LogLevelInfo:
		slogLevel = slog.LevelInfo
	case tracelog.LogLevelWarn:
		slogLevel = slog.LevelWarn
	case tracelog.LogLevelError:
		slogLevel = slog.LevelError
	default:
		slogLevel = slog.LevelInfo
	}

	attrs := make([]slog.Attr, 0, len(data))
	for k, v := range data {
		attrs = append(attrs, slog.Any(k, v))
	}

	if duration, ok := data["time"].(time.Duration); ok {
		if l.slowQueryThreshold > 0 && duration >= l.slowQueryThreshold {
			attrs = append(attrs, slog.Bool("slow_query", true))
			slogLevel = slog.LevelWarn
		}
	}

	l.logger.LogAttrs(ctx, slogLevel, msg, attrs...)
}

// TracerConfig creates a pgx.Tracer configuration for query logging.
func TracerConfig(logger *slog.Logger, logLevel tracelog.LogLevel, slowQueryThreshold time.Duration) pgx.QueryTracer {
	queryLogger := NewQueryLogger(logger, slowQueryThreshold)
	tracer := &tracelog.TraceLog{
		Logger:   queryLogger,
		LogLevel: logLevel,
	}
	return tracer
}

// FormatDuration formats a duration for logging.
func FormatDuration(d time.Duration) string {
	switch {
	case d >= time.Second:
		return fmt.Sprintf("%.2fs", d.Seconds())
	case d >= time.Millisecond:
		return fmt.Sprintf("%.2fms", float64(d.Microseconds())/1000.0)
	default:
		return fmt.Sprintf("%dus", d.Microseconds())
	}
}
