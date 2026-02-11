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

// contextKey is a private type for context keys in this package.
type contextKey int

const (
	queryStartTimeKey contextKey = iota
	querySQLKey
)

// QueryTracer implements pgx.QueryTracer to record metrics for every query
// while only logging slow queries and errors to reduce log spam.
type QueryTracer struct {
	logger             *slog.Logger
	slowQueryThreshold time.Duration
}

// NewQueryTracer creates a new QueryTracer that records metrics for all queries
// and logs slow queries (above threshold) and errors.
func NewQueryTracer(logger *slog.Logger, slowQueryThreshold time.Duration) *QueryTracer {
	return &QueryTracer{
		logger:             logger,
		slowQueryThreshold: slowQueryThreshold,
	}
}

// TraceQueryStart implements pgx.QueryTracer. Stores the start time and SQL in context.
func (t *QueryTracer) TraceQueryStart(ctx context.Context, _ *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	ctx = context.WithValue(ctx, queryStartTimeKey, time.Now())
	ctx = context.WithValue(ctx, querySQLKey, data.SQL)
	return ctx
}

// TraceQueryEnd implements pgx.QueryTracer. Records metrics and logs slow queries.
func (t *QueryTracer) TraceQueryEnd(ctx context.Context, _ *pgx.Conn, data pgx.TraceQueryEndData) {
	startTime, ok := ctx.Value(queryStartTimeKey).(time.Time)
	if !ok {
		return
	}
	duration := time.Since(startTime)
	sql, _ := ctx.Value(querySQLKey).(string)

	operation := extractOperation(sql)

	// Always record metrics for every query
	observability.DBQueryDuration.WithLabelValues(operation).Observe(duration.Seconds())

	if data.Err != nil {
		observability.DBQueryErrorsTotal.WithLabelValues(operation).Inc()
		t.logger.Error("query error",
			slog.String("operation", operation),
			slog.String("duration", FormatDuration(duration)),
			slog.Any("error", data.Err),
		)
		return
	}

	// Only log slow queries
	if t.slowQueryThreshold > 0 && duration >= t.slowQueryThreshold {
		t.logger.Warn("slow query",
			slog.String("operation", operation),
			slog.String("duration", FormatDuration(duration)),
			slog.String("sql", sql),
			slog.String("command_tag", data.CommandTag.String()),
		)
	}
}

// extractOperation returns the SQL operation type (select, insert, etc.) from a SQL string.
func extractOperation(sql string) string {
	if sql == "" {
		return "query"
	}

	// Skip sqlc comment prefix (e.g., "-- name: GetMovie :one\n")
	sqlToCheck := sql
	if strings.HasPrefix(sql, "-- ") {
		if idx := strings.Index(sql, "\n"); idx != -1 {
			sqlToCheck = strings.TrimSpace(sql[idx+1:])
		}
	}

	for _, prefix := range []string{"SELECT", "INSERT", "UPDATE", "DELETE", "BEGIN", "COMMIT", "ROLLBACK"} {
		if len(sqlToCheck) >= len(prefix) && strings.EqualFold(sqlToCheck[:len(prefix)], prefix) {
			return strings.ToLower(prefix)
		}
	}

	return "query"
}

// QueryLogger wraps slog.Logger for pgx query logging (used by tracelog.TraceLog).
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

// TracerConfig creates a pgx.QueryTracer that always records DB metrics
// and logs slow queries above the threshold.
func TracerConfig(logger *slog.Logger, _ tracelog.LogLevel, slowQueryThreshold time.Duration) (pgx.QueryTracer, *QueryTracer) {
	tracer := NewQueryTracer(logger, slowQueryThreshold)
	return tracer, tracer
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
