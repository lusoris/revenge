# Infrastructure Code Report

**Generated**: 2026-02-06
**Purpose**: Reference for writing Step 7 infra docs. Based on code exploration.

---

## Package Summary

| Package | Location | Files | Lines | Key Abstraction |
|---------|----------|-------|-------|-----------------|
| database | `internal/infra/database/` | 6 core + 8 tests | 587 core | Pool, migrations, query logging, metrics |
| cache | `internal/infra/cache/` | 4 core + 8 tests | ~1100 core | L1 (otter) + L2 (rueidis) unified cache |
| jobs | `internal/infra/jobs/` | 4 core | ~530 | River client, 5-level queues, progress tracking |
| health | `internal/infra/health/` | 4 core + 3 tests | 418 core | K8s probes, dependency checks |
| image | `internal/infra/image/` | 2 core + 1 test | 432 core | TMDb image download/proxy, dual-layer cache |
| logging | `internal/infra/logging/` | 2 core + 1 test | 194 core | slog + zap factories, tint dev mode |
| observability | `internal/infra/observability/` | 5 core + 1 test | 548 core | Prometheus metrics, pprof, HTTP middleware |
| search | `internal/infra/search/` | 1 core + 1 test | 240 core | Typesense client wrapper |

---

## Database (`internal/infra/database/`)

### Files
- `pool.go` (144) - pgxpool config, health check, stats
- `module.go` (39) - fx module: pool, queries, lifecycle hooks
- `migrate.go` (168) - golang-migrate with embedded FS (iofs)
- `logger.go` (81) - QueryLogger (slog), slow query detection
- `metrics.go` (95) - 12 Prometheus pool metrics
- `testing.go` (60) - embedded-postgres for tests

### Key Details
- **32 migrations** (64 files) in `migrations/shared/`
- **3-schema model**: `shared` (auth/config), `public` (content), `qar` (adult)
- **sqlc v1.30.0**: 425 generated methods in `db/` package
- **Pool**: CPU-aware (`CPU*2+1`), configurable max/min conns, health check period
- **Metrics**: acquire count/duration, active/idle/total conns, errors (Prometheus)
- **Config**: `database.url`, `database.max_conns`, `database.min_conns`, `database.max_conn_lifetime`, etc.

---

## Cache (`internal/infra/cache/`)

### Files
- `cache.go` (235) - Unified Cache struct (L1+L2), Get/Set/Delete, GetJSON/SetJSON, CacheAside, Invalidate
- `module.go` (155) - rueidis client init, fx lifecycle
- `otter.go` (86) - L1Cache[K,V] generic wrapper (otter v2, W-TinyLFU)
- `keys.go` (394) - 20+ key prefixes, 13 TTL constants, invalidation helpers

### Key Details
- **L1**: otter v2 (W-TinyLFU), 10k entries, 5min TTL
- **L2**: rueidis (Dragonfly/Redis), 16MiB client-side cache per conn
- **Smart TTL**: Short TTLs skip L1 to prevent stale divergence
- **Graceful degradation**: Works without L2 (L1-only mode)
- **Async population**: Cache sets in background goroutines
- **6 cached service wrappers**: session, user, rbac, settings, library, search
- **TTL ranges**: 30s (session validation) to 24h (image metadata)
- **Config**: `cache.url`, `cache.enabled`

---

## Jobs (`internal/infra/jobs/`)

### Files
- `river.go` (146) - Client wrapper (NewClient, Insert, Start/Stop)
- `module.go` (58) - fx module, lifecycle hooks
- `queues.go` (133) - 5-level priority queues, backoff strategies
- `progress.go` (84) - JobProgress tracking for frontend polling
- `notification_job.go` (213) - NotificationWorker
- `cleanup_job.go` (129) - Generic CleanupWorker

### Queue Config
| Queue | Workers | Purpose |
|-------|---------|---------|
| CRITICAL | 20 | Security events, urgent |
| HIGH | 15 | Notifications, webhooks |
| DEFAULT | 10 | Metadata, sync |
| LOW | 5 | Cleanup, maintenance |
| BULK | 3 | Library scans, batch |

### All 17 Workers (across codebase)
| Worker | Kind | Location | Queue | Timeout |
|--------|------|----------|-------|---------|
| Notification | `notification` | infra/jobs | HIGH | 2m |
| Cleanup (generic) | `cleanup` | infra/jobs | LOW | 2m |
| Library Scan Cleanup | `library_scan_cleanup` | service/library | LOW | 2m |
| Activity Cleanup | `activity_cleanup` | service/activity | LOW | 2m |
| Movie Library Scan | `movie_library_scan` | content/movie/moviejobs | BULK | 30m |
| Movie Metadata Refresh | `metadata_refresh_movie` | content/movie/moviejobs | DEFAULT | 5m |
| Movie File Match | `movie_file_match` | content/movie/moviejobs | DEFAULT | 5m |
| Movie Search Index | `movie_search_index` | content/movie/moviejobs | DEFAULT | 15m |
| TVShow Library Scan | `tvshow_library_scan` | content/tvshow/jobs | BULK | 30m |
| TVShow Metadata Refresh | `tvshow_metadata_refresh` | content/tvshow/jobs | DEFAULT | 15m |
| TVShow File Match | `tvshow_file_match` | content/tvshow/jobs | DEFAULT | 5m |
| TVShow Search Index | `tvshow_search_index` | content/tvshow/jobs | BULK | 10m |
| TVShow Series Refresh | `tvshow_series_refresh` | content/tvshow/jobs | DEFAULT | 10m |
| Radarr Sync | `radarr_sync` | integration/radarr | HIGH | 10m |
| Radarr Webhook | `radarr_webhook` | integration/radarr | HIGH | 1m |
| Sonarr Sync | `sonarr_sync` | integration/sonarr | HIGH | 10m |
| Sonarr Webhook | `sonarr_webhook` | integration/sonarr | HIGH | 1m |

- **Backoff**: Exponential (1s * 2^attempt, capped at 1h)
- **Default max attempts**: 25
- **Leader election**: Cleanup jobs skip if not leader
- **Config**: `jobs.max_workers`, `jobs.fetch_cooldown`

---

## Health (`internal/infra/health/`)

### Files
- `checks.go` (107) - CheckDatabase, CheckCache, CheckJobs, CheckAll
- `service.go` (165) - K8s probes: Liveness, Readiness, Startup
- `handler.go` (115) - HTTP handlers, RegisterRoutes
- `module.go` (31) - fx module

### Key Details
- **Liveness**: Always healthy (process alive check)
- **Readiness**: Database required, cache/jobs optional (degraded OK)
- **Startup**: Tracks initialization, marks complete via `MarkStartupComplete()`
- **Status enum**: healthy, unhealthy, degraded
- **CheckResult**: Name, Status, Message, Details (map)
- **Thread-safe**: RWMutex for startup state

---

## Image (`internal/infra/image/`)

### Files
- `service.go` (406) - FetchImage, StreamImage, ServeHTTP, ClearCache
- `module.go` (26) - fx module

### Key Details
- **NO govips** - HTTP download + filesystem cache, no image manipulation
- **TMDb integration**: Builds URLs with size variants (w185-original)
- **Dual cache**: sync.Map (memory) + filesystem (configurable dir, 7-day TTL)
- **HTTP proxy**: Implements http.Handler for `/images/{type}/{size}/{path}`
- **Security**: Validates MIME types, size limits (10MB), prevents directory traversal
- **Headers**: ETag, Cache-Control (immutable), CORS, If-None-Match
- **HTTP client**: imroc/req with 30s timeout, 3 retries

---

## Logging (`internal/infra/logging/`)

### Files
- `logging.go` (140) - NewLogger (slog), NewZapLogger (zap), NewTestLogger
- `module.go` (54) - fx module, ProvideSlogLogger, ProvideZapLogger

### Key Details
- **slog**: Default logger, context-aware, `slog.SetDefault()`
- **zap**: Structured logging, `zap.ReplaceGlobals()`
- **Dev mode**: tint handler (colorized, human-readable, source location)
- **Prod mode**: JSON handler (ISO8601 timestamps)
- **Config**: `logging.level`, `logging.format` (text/json), `logging.development`
- **Both loggers** provided via fx to all services

---

## Observability (`internal/infra/observability/`)

### Files
- `metrics.go` (235) - 40+ Prometheus metric definitions, recording functions
- `middleware.go` (188) - ogen + standard HTTP metrics middleware
- `server.go` (93) - Observability HTTP server (metrics + pprof)
- `pprof.go` (24) - pprof handler registration
- `module.go` (8) - fx module

### Metrics Categories
- **HTTP**: requests_total, duration, in_flight
- **Session**: active_total
- **Cache**: hits/misses (by layer), operation_duration, size
- **Database**: query_duration, query_errors (+ pool metrics in database pkg)
- **Jobs**: enqueued_total, completed_total, duration, queue_size
- **Library**: scan_duration, files_scanned, scan_errors
- **Search**: queries_total, query_duration
- **Auth**: attempts_total, ratelimit_hits

### Key Details
- **Dual-port**: Observability on `server.port + 1000`
- **pprof**: Dev mode only (heap, goroutine, block, mutex, allocs, trace)
- **Path normalization**: UUIDs/IDs → `{id}` (prevents cardinality explosion)
- **NO OpenTelemetry** - Prometheus + pprof only
- **Config**: Enabled via `logging.development` for pprof

---

## Search (`internal/infra/search/`)

### Files
- `module.go` (240) - Client wrapper, collection/document ops, health, fx hooks

### Key Details
- **Typesense client wrapper**: typesense-go SDK
- **Nil-safe**: All ops check `IsEnabled()`, graceful when disabled
- **Collection ops**: Create, delete, retrieve, list
- **Document ops**: Index, update, delete, bulk import
- **Search ops**: Single + multi-collection search
- **Startup**: 5 retries with exponential backoff, non-fatal failure
- **Config**: `search.url`, `search.api_key`, `search.enabled`
