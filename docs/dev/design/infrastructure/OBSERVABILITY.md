# Observability Infrastructure

<!-- DESIGN: infrastructure -->

**Package**: `internal/infra/observability`
**fx Module**: `observability.Module`

> Prometheus metrics, pprof profiling, and HTTP instrumentation middleware

---

## Service Structure

```
internal/infra/observability/
├── metrics.go             # 40+ Prometheus metric definitions + recording functions
├── middleware.go           # ogen + standard HTTP metrics middleware
├── server.go              # Observability HTTP server (metrics + pprof)
├── pprof.go               # pprof handler registration
└── module.go              # fx module
```

## Observability Server

Runs on **port + 1000** (e.g., if app is on 8080, observability is on 9080).

**Endpoints**:
- `GET /metrics` - Prometheus scrape endpoint (always enabled)
- `GET /health/live` - Simple liveness check
- `GET /health/ready` - Simple readiness check
- `GET /debug/pprof/*` - pprof profiles (development mode only)

## Prometheus Metrics

All metrics use namespace `revenge`.

### HTTP

| Metric | Type | Labels | Description |
|--------|------|--------|-------------|
| `http_requests_total` | Counter | method, path, status | Total HTTP requests |
| `http_request_duration_seconds` | Histogram | method, path | Request latency (buckets: 1ms-10s) |
| `http_requests_in_flight` | Gauge | - | Currently processing requests |

### Cache

| Metric | Type | Labels | Description |
|--------|------|--------|-------------|
| `cache_hits_total` | Counter | cache, layer | Cache hits (L1/L2) |
| `cache_misses_total` | Counter | cache, layer | Cache misses (L1/L2) |
| `cache_operation_duration_seconds` | Histogram | cache, operation | Operation latency (buckets: 0.1ms-100ms) |
| `cache_size` | Gauge | cache | Current L1 item count |

### Database

| Metric | Type | Labels | Description |
|--------|------|--------|-------------|
| `db_query_duration_seconds` | Histogram | operation | Query latency (buckets: 1ms-5s) |
| `db_query_errors_total` | Counter | operation | Query errors |

Pool metrics defined separately in `database/metrics.go` (12 additional metrics).

### Job Queue

| Metric | Type | Labels | Description |
|--------|------|--------|-------------|
| `jobs_enqueued_total` | Counter | job_type | Jobs enqueued |
| `jobs_completed_total` | Counter | job_type, status | Jobs completed (success/failure) |
| `jobs_duration_seconds` | Histogram | job_type | Job processing time (buckets: 100ms-600s) |
| `jobs_queue_size` | Gauge | state | Queue size by state |

### Library Scanner

| Metric | Type | Labels | Description |
|--------|------|--------|-------------|
| `library_scan_duration_seconds` | Histogram | library_id | Scan duration (buckets: 1s-3600s) |
| `library_files_scanned_total` | Counter | library_id | Files scanned |
| `library_scan_errors_total` | Counter | library_id, error_type | Scan errors |

### Search

| Metric | Type | Labels | Description |
|--------|------|--------|-------------|
| `search_queries_total` | Counter | type | Search queries |
| `search_query_duration_seconds` | Histogram | type | Search latency (buckets: 1ms-1s) |

### Auth

| Metric | Type | Labels | Description |
|--------|------|--------|-------------|
| `auth_attempts_total` | Counter | method, status | Auth attempts |
| `ratelimit_hits_total` | Counter | limiter, action | Rate limit triggers |

### Session

| Metric | Type | Labels | Description |
|--------|------|--------|-------------|
| `sessions_active_total` | Gauge | - | Active sessions |

## HTTP Middleware

Two middleware variants:

**ogen middleware**: Integrates with ogen-generated API handlers. Extracts HTTP status from response type names (e.g., `GetMovieOK` → 200).

**Standard HTTP middleware**: Wraps `http.ResponseWriter` to capture status codes for non-ogen routes.

**Path normalization**: UUIDs, numeric IDs, and ID-like patterns replaced with `{id}` to prevent cardinality explosion.

## pprof Profiles

Available in development mode (`logging.development: true`):

| Endpoint | Profile |
|----------|---------|
| `/debug/pprof/heap` | Memory allocations |
| `/debug/pprof/goroutine` | Goroutine dump |
| `/debug/pprof/block` | Block contention |
| `/debug/pprof/mutex` | Mutex contention |
| `/debug/pprof/allocs` | All allocations |
| `/debug/pprof/threadcreate` | Thread creation |
| `/debug/pprof/profile` | CPU profile (30s) |
| `/debug/pprof/trace` | Execution trace |

## Dependencies

- `github.com/prometheus/client_golang` - Metrics registration + HTTP handler
- `net/http/pprof` - Standard library profiling

**Note**: No OpenTelemetry or distributed tracing. Project uses Prometheus + pprof.

## Related Documentation

- [LOGGING.md](LOGGING.md) - Structured logging (development flag enables pprof)
- [HEALTH.md](HEALTH.md) - Health check endpoints
- [DATABASE.md](DATABASE.md) - Database pool metrics
