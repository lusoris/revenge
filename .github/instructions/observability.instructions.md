---
applyTo: "**/pkg/metrics/**/*.go,**/internal/api/**/*.go"
---

# Observability & Metrics Instructions

## Overview

`pkg/metrics` provides lightweight metrics collection:

- Counters for totals
- Gauges for current values
- Histograms for distributions
- Timers for latencies
- HTTP middleware

## Metrics Types

### Counter

Monotonically increasing value (resets on restart):

```go
registry := metrics.NewRegistry()

// Get or create counter
requestCounter := registry.Counter("http_requests_total")
errorCounter := registry.Counter("errors_total")

// Increment
requestCounter.Inc()       // +1
errorCounter.Add(5)        // +5
```

### Gauge

Value that can go up and down:

```go
activeStreams := registry.Gauge("active_streams")
cacheSize := registry.Gauge("cache_size_bytes")

// Set/modify
activeStreams.Inc()           // +1
activeStreams.Dec()           // -1
cacheSize.Set(1024 * 1024)    // Set absolute value
```

### Timer

Tracks durations with percentiles:

```go
timer := registry.Timer("request_duration")

// Option 1: Defer
func handleRequest() {
    done := timer.Time()
    defer done()

    // ... handle request
}

// Option 2: Manual
start := time.Now()
// ... do work
timer.ObserveDuration(time.Since(start))

// Get percentiles
p50 := timer.P50()  // 50th percentile
p95 := timer.P95()  // 95th percentile
p99 := timer.P99()  // 99th percentile
```

### Histogram

Track value distributions:

```go
responseSizes := metrics.NewHistogram([]int64{
    100, 1000, 10000, 100000, 1000000, // Bucket boundaries
})

responseSizes.Observe(4500)  // Goes into 10000 bucket

p90 := responseSizes.Percentile(90)
```

## HTTP Middleware

Automatic HTTP metrics:

```go
httpMetrics := metrics.NewHTTPMetrics()

// Use as middleware
mux := http.NewServeMux()
handler := httpMetrics.Middleware(mux)
http.ListenAndServe(":8080", handler)

// Get stats
stats := httpMetrics.Stats()
// {
//   "requests_total": 1234,
//   "requests_active": 5,
//   "response_time_p50": 12,
//   "response_time_p95": 45,
//   "response_time_p99": 120,
//   "status_200": 1100,
//   "status_404": 50,
//   "status_500": 10,
// }
```

## Metrics Endpoint

```go
mux.HandleFunc("GET /metrics", func(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(registry.Snapshot())
})
```

## Key Metrics to Track

### Application Metrics

| Metric                     | Type    | Description              |
| -------------------------- | ------- | ------------------------ |
| `http_requests_total`      | Counter | Total HTTP requests      |
| `http_request_duration_ms` | Timer   | Request latency          |
| `active_sessions`          | Gauge   | Current active sessions  |
| `active_streams`           | Gauge   | Current playback streams |

### Playback Metrics

| Metric                       | Type    | Description                 |
| ---------------------------- | ------- | --------------------------- |
| `playback_started`           | Counter | Playback sessions started   |
| `playback_completed`         | Counter | Playback completed          |
| `playback_errors`            | Counter | Playback errors             |
| `transcodes_active`          | Gauge   | Active transcoding sessions |
| `transcode_cache_hits`       | Counter | Cache hit count             |
| `transcode_cache_size_bytes` | Gauge   | Cache size                  |

### Database Metrics

| Metric                  | Type    | Description        |
| ----------------------- | ------- | ------------------ |
| `db_connections_active` | Gauge   | Active connections |
| `db_connections_idle`   | Gauge   | Idle connections   |
| `db_query_duration_ms`  | Timer   | Query latency      |
| `db_errors`             | Counter | Database errors    |

### External Service Metrics

| Metric                | Type    | Description         |
| --------------------- | ------- | ------------------- |
| `blackbeard_requests` | Counter | Transcoder requests |
| `blackbeard_errors`   | Counter | Transcoder errors   |
| `tmdb_requests`       | Counter | TMDb API calls      |
| `tmdb_errors`         | Counter | TMDb errors         |

## Structured Logging

Use `slog` for structured logs:

```go
logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
    Level:     slog.LevelInfo,
    AddSource: true,
}))

// Context logging
logger.Info("request handled",
    slog.String("method", r.Method),
    slog.String("path", r.URL.Path),
    slog.Int("status", status),
    slog.Duration("duration", duration),
)

// Error logging
logger.Error("database query failed",
    slog.String("query", query),
    slog.Any("error", err),
)
```

### Request-Scoped Logging

```go
func requestLogger(r *http.Request) *slog.Logger {
    return logger.With(
        slog.String("request_id", r.Header.Get("X-Request-ID")),
        slog.String("user_id", getUserID(r)),
        slog.String("remote_addr", r.RemoteAddr),
    )
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
    log := requestLogger(r)

    log.Info("request started")
    // ... handle
    log.Info("request completed", slog.Int("status", 200))
}
```

## Integration with fx

```go
func NewMetricsRegistry() *metrics.Registry {
    return metrics.NewRegistry()
}

func NewHTTPMetrics(registry *metrics.Registry) *metrics.HTTPMetrics {
    return metrics.NewHTTPMetrics()
}

var MetricsModule = fx.Options(
    fx.Provide(NewMetricsRegistry),
    fx.Provide(NewHTTPMetrics),
)
```

## DO's

- ✅ Track request counts and latencies
- ✅ Track error rates
- ✅ Track active connections/sessions
- ✅ Use consistent naming (snake_case)
- ✅ Add context to logs (request_id, user_id)

## DON'Ts

- ❌ Track sensitive data (passwords, tokens)
- ❌ Create unbounded label cardinality
- ❌ Log request/response bodies
- ❌ Use counters for values that can decrease
- ❌ Forget to track errors

## Log Levels

| Level | Use Case                           |
| ----- | ---------------------------------- |
| Debug | Development, verbose tracing       |
| Info  | Normal operations, key events      |
| Warn  | Recoverable issues, degraded state |
| Error | Failures requiring attention       |

```go
logger.Debug("cache lookup", slog.String("key", key))
logger.Info("user logged in", slog.String("user_id", userID))
logger.Warn("rate limit approaching", slog.Int("current", current))
logger.Error("database connection failed", slog.Any("error", err))
```
