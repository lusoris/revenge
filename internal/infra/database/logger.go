package database

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/tracelog"

	"github.com/lusoris/revenge/internal/infra/observability"
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

	// Record DB query metrics
	if duration, ok := data["time"].(time.Duration); ok {
		operation := "query"
		if sql, ok := data["sql"].(string); ok && len(sql) > 0 {
			// Skip sqlc comment prefix (e.g., "-- name: GetMovie :one\n")
			sqlToCheck := sql
			if strings.HasPrefix(sql, "-- ") {
				if idx := strings.Index(sql, "\n"); idx != -1 {
					sqlToCheck = strings.TrimSpace(sql[idx+1:])
				}
			}
			// Extract operation type from SQL (SELECT, INSERT, UPDATE, DELETE)
			for _, prefix := range []string{"SELECT", "INSERT", "UPDATE", "DELETE", "BEGIN", "COMMIT", "ROLLBACK"} {
				if len(sqlToCheck) >= len(prefix) && strings.EqualFold(sqlToCheck[:len(prefix)], prefix) {
					operation = strings.ToLower(prefix)
					break
				}
			}
		}
		observability.DBQueryDuration.WithLabelValues(operation).Observe(duration.Seconds())
	}
	if level == tracelog.LogLevelError {
		operation := "query"
		if sql, ok := data["sql"].(string); ok && len(sql) > 0 {
			// Skip sqlc comment prefix
			sqlToCheck := sql
			if strings.HasPrefix(sql, "-- ") {
				if idx := strings.Index(sql, "\n"); idx != -1 {
					sqlToCheck = strings.TrimSpace(sql[idx+1:])
				}
			}
			for _, prefix := range []string{"SELECT", "INSERT", "UPDATE", "DELETE", "BEGIN", "COMMIT", "ROLLBACK"} {
				if len(sqlToCheck) >= len(prefix) && strings.EqualFold(sqlToCheck[:len(prefix)], prefix) {
					operation = strings.ToLower(prefix)
					break
				}
			}
		}
		observability.DBQueryErrorsTotal.WithLabelValues(operation).Inc()
	}

	l.logger.LogAttrs(ctx, slogLevel, msg, attrs...)
}

// TracerConfig creates a pgx.Tracer configuration for query logging.
// Returns the tracer and the underlying QueryLogger for lifecycle control.
func TracerConfig(logger *slog.Logger, logLevel tracelog.LogLevel, slowQueryThreshold time.Duration) (pgx.QueryTracer, *QueryLogger) {
	queryLogger := NewQueryLogger(logger, slowQueryThreshold)
	tracer := &tracelog.TraceLog{
		Logger:   queryLogger,
		LogLevel: logLevel,
	}
	return tracer, queryLogger
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
