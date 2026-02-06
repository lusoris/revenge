# Observability Pattern

> Metrics, logging, profiling, and the observability server. Covers what's implemented and what's planned. Written from code as of 2026-02-06.

---

## Architecture

Observability runs on a **separate HTTP server** from the main API, on port `API_PORT + 1000` (e.g., API on 8096 -> observability on 9096). No authentication required.

```
internal/infra/observability/
  metrics.go      # Prometheus metric definitions + recording functions
  middleware.go   # HTTP metrics middleware (ogen-compatible)
  pprof.go        # pprof profiling endpoints (dev mode)
  server.go       # Observability HTTP server
  module.go       # fx wiring
```

Wired via `observability.Module` in `internal/app/module.go`.

---

## Metrics (Prometheus)

All metrics use namespace `revenge` with subsystem-specific prefixes.

### Implemented and Active

| Metric | Type | Labels | Recorded In |
|--------|------|--------|-------------|
| `revenge_http_requests_total` | Counter | method, path, status | `middleware.go` |
| `revenge_http_request_duration_seconds` | Histogram | method, path | `middleware.go` |
| `revenge_http_requests_in_flight` | Gauge | — | `middleware.go` |
| `revenge_cache_hits_total` | Counter | cache, layer | `cache/cache.go` |
| `revenge_cache_misses_total` | Counter | cache, layer | `cache/cache.go` |
| `revenge_cache_operation_duration_seconds` | Histogram | cache, operation | `cache/cache.go` |
| `revenge_cache_size` | Gauge | cache | `cache/cache.go` |

### Defined but Not Yet Instrumented

| Metric | Type | Labels | Status |
|--------|------|--------|--------|
| `revenge_db_query_duration_seconds` | Histogram | operation | Needs pgx hooks |
| `revenge_db_query_errors_total` | Counter | operation | Needs pgx hooks |
| `revenge_jobs_enqueued_total` | Counter | job_type | Needs River hooks |
| `revenge_jobs_completed_total` | Counter | job_type, status | Needs River hooks |
| `revenge_jobs_duration_seconds` | Histogram | job_type | Needs River hooks |
| `revenge_jobs_queue_size` | Gauge | state | Needs River hooks |
| `revenge_library_scan_duration_seconds` | Histogram | library_id | Needs scan instrumentation |
| `revenge_library_files_scanned_total` | Counter | library_id | Needs scan instrumentation |
| `revenge_library_scan_errors_total` | Counter | library_id, error_type | Needs scan instrumentation |
| `revenge_search_queries_total` | Counter | type | Needs Typesense instrumentation |
| `revenge_search_query_duration_seconds` | Histogram | type | Needs Typesense instrumentation |
| `revenge_auth_attempts_total` | Counter | method, status | Needs auth service calls |
| `revenge_ratelimit_hits_total` | Counter | limiter, action | Needs middleware calls |
| `revenge_sessions_active_total` | Gauge | — | Needs session service calls |

Recording functions exist (`RecordAuthAttempt`, `RecordJobEnqueued`, etc.) but aren't called from the relevant services yet.

### Histogram Buckets

Buckets are tuned per subsystem:

| Subsystem | Buckets (seconds) |
|-----------|-------------------|
| HTTP | 0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10 |
| Cache | 0.0001, 0.0005, 0.001, 0.005, 0.01, 0.025, 0.05, 0.1 |
| Database | 0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5 |
| Jobs | 0.1, 0.5, 1, 5, 10, 30, 60, 120, 300, 600 |
| Library scan | 1, 5, 10, 30, 60, 120, 300, 600, 1800, 3600 |
| Search | 0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1 |

### Adding a New Metric

```go
// In metrics.go
var MyMetric = promauto.NewCounterVec(prometheus.CounterOpts{
    Namespace: "revenge",
    Subsystem: "myservice",
    Name:      "operations_total",
    Help:      "Total number of operations.",
}, []string{"operation", "status"})

func RecordMyOperation(operation, status string) {
    MyMetric.WithLabelValues(operation, status).Inc()
}
```

---

## HTTP Middleware

### Ogen-Compatible Metrics Middleware

The `HTTPMetricsMiddleware` in `middleware.go` wraps ogen-generated handlers:

1. Increments `HTTPRequestsInFlight` on request start, decrements on completion
2. Records `HTTPRequestsTotal` with method, normalized path, and status code
3. Observes `HTTPRequestDuration` latency

**Path normalization** prevents metric cardinality explosion:
- Numeric IDs: `/movies/123` -> `/movies/{id}`
- UUIDs: `/users/550e8400-...` -> `/users/{id}`
- Short IDs: `/items/abc12345` -> `/items/{id}`

**Status extraction**: Extracts HTTP status codes from ogen type-safe response objects (e.g., `GetMovieOK` -> 200, `NotFound` -> 404).

### Middleware Chain Order

```
1. RequestIDMiddleware()        # Generate/extract X-Request-ID
2. RequestMetadataMiddleware()  # Extract IP, User-Agent, Accept-Language
3. HTTPMetricsMiddleware()      # Record request metrics
4. RateLimiter.Middleware()     # In-memory or Redis rate limiting
5. Handler execution
```

---

## Logging

### Framework

| Mode | Handler | Format | Features |
|------|---------|--------|----------|
| Development | `tint` | Colorized text | Source location, short timestamps |
| Production | `slog.JSONHandler` | JSON | ISO8601 timestamps, machine-parseable |

Both modes use `log/slog` as the primary logger. A `zap.Logger` is also provided for components that need it.

### Configuration

```go
type LoggingConfig struct {
    Level       string // "debug", "info", "warn", "error"
    Format      string // "text" or "json"
    Development bool   // Enables colored output, source info
}
```

### Pattern

```go
logger.Info("Operation completed",
    slog.String("movie_id", id.String()),
    slog.Int("count", count),
    slog.Duration("elapsed", elapsed),
)
```

---

## Profiling (pprof)

Available in development mode only (`config.Logging.Development == true`).

Served on the observability server:

| Endpoint | Purpose |
|----------|---------|
| `/debug/pprof/` | Index page |
| `/debug/pprof/heap` | Heap profile |
| `/debug/pprof/goroutine` | Goroutine stacks |
| `/debug/pprof/block` | Block contention |
| `/debug/pprof/mutex` | Mutex contention |
| `/debug/pprof/allocs` | Memory allocations |
| `/debug/pprof/profile` | CPU profile (30s default) |
| `/debug/pprof/trace` | Execution trace |

---

## Observability Server

| Endpoint | Purpose |
|----------|---------|
| `GET /metrics` | Prometheus scrape target |
| `GET /health/live` | Kubernetes liveness probe |
| `GET /health/ready` | Kubernetes readiness probe |
| `GET /debug/pprof/*` | Profiling (dev mode only) |

Started as a background goroutine during fx app startup. Graceful shutdown with context timeout.

---

## Rate Limiting

Two implementations, automatic fallback:

| Implementation | Backend | Scope |
|---------------|---------|-------|
| In-memory | `golang.org/x/time/rate` | Per-instance, per-IP |
| Redis-backed | Sliding window (Lua script via rueidis) | Shared across instances |

Falls back to in-memory if Redis is unavailable. Stricter limits on auth endpoints.

---

## Planned Work

### Phase 1: Fill Instrumentation Gaps
- Wire `RecordAuthAttempt` calls into auth service
- Wire `RecordJobEnqueued`/`RecordJobCompleted` into River hooks
- Wire library scan, search, and database metrics
- Wire rate limit hit recording in middleware

### Phase 2: Per-Handler RED Metrics
- Add per-endpoint Rate, Errors, Duration metrics via ogen middleware
- Create example Prometheus scrape config
- Build Grafana dashboard JSONs

### Phase 3: Distributed Tracing
- OpenTelemetry SDK integration
- W3C Trace Context (evolve X-Request-ID)
- OTLP exporter to Jaeger/Tempo

---

## Key Dashboard Metrics

**Application health**:
- `rate(revenge_http_requests_total{status=~"5.."}[5m])` — error rate
- `histogram_quantile(0.95, revenge_http_request_duration_seconds)` — p95 latency

**Cache performance**:
- `revenge_cache_hits_total / (hits_total + misses_total)` — hit ratio by layer

**Job queue health**:
- `revenge_jobs_queue_size{state="available"}` — pending work
- `rate(revenge_jobs_completed_total{status="failed"}[5m])` — failure rate

---

## Related Documentation

- [Cache Strategy](CACHE_STRATEGY.md) — L1/L2 caching with metrics integration
- [River Workers](RIVER_WORKERS.md) — Job workers that should record metrics
- [Error Handling](ERROR_HANDLING.md) — Error flow that feeds into error rate metrics
