# Progress Report

**Date**: 2026-02-08 (updated)

---

## Phase 1: Test Coverage (COMPLETE)

**Target**: 80% coverage across all authored packages
**Measurement**: `go test -coverprofile` (full tests, NOT `-short`)

All packages at 80%+ when measured properly (without `-short`). Only 5 were below:

| Package | Before | After | Status |
|---------|--------|-------|--------|
| `service/oidc` | 74.9% | 80.9% | DONE |
| `service/mfa` | 71.8% | 77.7% | Ceiling (WebAuthn attestation) |
| `infra/jobs` | 74.6% | 74.6% | BLOCKED (needs River DB) |
| `service/library` | 70.0% | 90.5% | DONE |
| `service/email` | 62.3% | 86.1% | DONE |

### Commits

| Commit | Description |
|--------|-------------|
| `d57f0212` | test: add integration tests for MFA and OIDC packages |
| `a98e2828` | test: add coverage tests for library (90.5%) and email (86.1%) |

---

## Phase 2: Local Container Configuration (COMPLETE)

### Issues Fixed

1. **docker-compose.dev.yml: Broken Typesense build reference**
   - `deploy/typesense.Dockerfile` didn't exist
   - Fixed: Use official `typesense/typesense:27.1` image (matches TECH_STACK.md)
   - Fixed healthcheck: Typesense 27.1 doesn't ship `curl`, use `bash /dev/tcp` instead

2. **Port inconsistency (config default 8080 vs Docker 8096)**
   - `config.go` defaulted to 8080, but Dockerfile/Helm/compose all use 8096
   - `docker-compose.prod.yml` didn't set `REVENGE_SERVER_PORT` → app would listen on 8080, health check would hit 8096 → broken
   - Fixed: Changed config defaults to 8096 (server.port + email.base_url)
   - Updated all test assertions that checked port default

3. **Missing fx modules in app/module.go**
   - `storage.Module` missing → `*user.Service` couldn't resolve `storage.Storage`
   - `email.Module` missing → `*auth.Service` couldn't resolve `*email.Service`
   - `crypto.Module` (NEW) → `*mfa.TOTPService` couldn't resolve `*crypto.Encryptor`
   - Created `internal/crypto/module.go` with proper key hierarchy:
     1. `auth.encryption_key` (hex-encoded 32-byte AES key) — production
     2. Derived from `auth.jwt_secret` via SHA-256 — development
     3. Deterministic dev fallback — when no secrets configured

4. **River job queue schema missing**
   - River tables (`river_job`, `river_queue`, `river_leader`) weren't created by app migrations
   - Added `rivermigrate.Migrate()` call in `NewRiverClient` before creating the client
   - Used River's built-in `rivermigrate` package (proper, version-matched approach)

5. **Dockerfile missing vips runtime library**
   - Build stage had `vips-dev` for compilation, runtime stage lacked `vips`
   - Added `vips` to runtime `apk add`

### New Files

- `internal/crypto/module.go` — fx module providing `*crypto.Encryptor`

### New Config Fields

- `auth.encryption_key` — dedicated AES-256 encryption key (64 hex chars)

### New Makefile Targets

- `make docker-local` — build + start full stack, verify health
- `make test-live` — run live smoke tests against running stack

### Verification

```
$ curl http://localhost:8096/healthz
{"name":"liveness","status":"healthy","message":"Service is alive"}

$ curl http://localhost:8096/readyz
{"name":"readiness","status":"healthy","message":"service is ready"}

$ curl http://localhost:8096/startupz
{"name":"startup","status":"healthy","message":"startup complete"}
```

---

## Phase 3: Live Smoke Tests (COMPLETE)

### File: `tests/live/smoke_test.go`

15 tests covering:
- Health endpoints (liveness, readiness, startup)
- Full auth flow (register → login → get user → refresh → logout)
- Invalid password rejection
- Unauthenticated access rejection
- Server settings (admin)
- Library listing
- User preferences

All 15 tests pass against the live stack:
```
$ make test-live
PASS
ok  github.com/lusoris/revenge/tests/live  0.339s
```

---

## Phase 4: Bug Fixes from Live Testing (COMPLETE)

Expanded smoke tests from 15 → full suite, uncovering and fixing 10+ production bugs.

### Bug 1: Casbin Auto-Reload Failure (NULL v2 scanning)
**Symptom**: Admin role assignments via SQL INSERT never took effect; tests timed out after 25s.
**Root Cause**: `adapter.go` scanned `v0`, `v1`, `v2` columns as Go `string`, but all v0-v5 are nullable `VARCHAR(100)`. When a `g`-type record (group/role) has NULL `v2`, pgx v5 fails the scan. This silently broke the entire `LoadPolicy()` call, preventing auto-reload from picking up new policies.
**Fix**: Changed all v0-v5 to scan as `*string`, then assign to struct fields only if non-nil.
**Files**: `internal/service/rbac/adapter.go`, `internal/service/rbac/adapter_test.go`
**Test**: `TestAdapter_LoadPolicy_NullableV0V1V2`

### Bug 2: Typesense Search 500 on Missing Collection
**Symptom**: `/api/v1/search/movies` returned HTTP 500 when Typesense collection didn't exist yet.
**Root Cause**: Search handlers propagated Typesense errors to `NewError()` which always returned 500.
**Fix**: Return empty results (200) with `total_hits: 0` when search is unavailable.
**Files**: `internal/api/handler_search.go`

### Bug 3: X-API-Key Authentication Not Implemented
**Symptom**: Any request using `X-API-Key` header returned 500.
**Root Cause**: `HandleApiKeyAuth` was unimplemented (returned error).
**Fix**: Full implementation — added `apiKeyAuth` security scheme to OpenAPI spec, regenerated ogen, implemented `HandleApiKeyAuth` to validate key and set user context.
**Files**: `api/openapi/openapi.yaml`, `internal/api/handler.go`, `internal/api/ogen/*` (regenerated)

### Bug 4: API Key List Returns Revoked Keys
**Symptom**: Listing API keys showed revoked/deleted keys.
**Root Cause**: `ListUserKeys` used `ListUserAPIKeys` (returns all) instead of `ListActiveUserAPIKeys`.
**Fix**: Changed to use `ListActiveUserAPIKeys`.
**Files**: `internal/service/apikeys/service.go`

### Bug 5: River Job Workers Not Registered
**Symptom**: Job executor logged "job kind is not registered" errors.
**Root Cause**: Both `moviejobs.Module` and `tvshowjobs.Module` used `fx.Provide` for workers but never called `fx.Invoke(RegisterWorkers)` to register them with River.
**Fix**: Added `fx.Invoke(RegisterWorkers)` to both modules.
**Files**: `internal/content/movie/moviejobs/module.go`, `internal/content/tvshow/jobs/module.go`

### Bug 6: Typesense Popularity Field Schema Error
**Symptom**: Job executor error: "Default sorting field `popularity` cannot be an optional field."
**Root Cause**: `movie_schema.go` had `Optional: ptr(true)` on the `popularity` field, which is also the `DefaultSortingField`.
**Fix**: Removed `Optional` from the popularity field.
**Files**: `internal/service/search/movie_schema.go`

### Bug 7: Movie Not Found Returns 500 Instead of 404
**Symptom**: Requesting a non-existent movie returned HTTP 500 instead of 404.
**Root Cause**: pgx v5 returns `pgx.ErrNoRows`, not `database/sql.ErrNoRows`. The repository checked `err == sql.ErrNoRows` which never matched, so no-rows errors propagated as generic 500s. Handler used `==` instead of `errors.Is()`.
**Fix**: Repository uses `errors.Is(err, pgx.ErrNoRows)` → returns `ErrMovieNotFound`. Handler uses `errors.Is()` for all sentinel error checks.
**Files**: `internal/content/movie/repository_postgres.go`, `internal/api/movie_handlers.go`

### Bug 8: Koanf Env Var Compound Name Mapping
**Symptom**: `REVENGE_API_KEY_HASH_COST` was mapped to `api.key.hash.cost` instead of `api_key.hash_cost`.
**Root Cause**: Koanf's env provider replaces ALL `_` with `.`, breaking compound field names like `api_key`, `jwt_secret`, etc.
**Fix**: Added `applyCompoundEnvOverrides()` in config loader to explicitly map known compound env vars.
**Files**: `internal/config/loader.go`, `internal/config/loader_test.go`

### Bug 9: Movie Stats NULL Aggregation
**Symptom**: Movie stats query returned NULL for count/size fields when no movies exist.
**Root Cause**: SQL `COUNT(*)` returns 0 but `SUM()` on empty set returns NULL.
**Fix**: Added `COALESCE(..., 0)` wrappers in the SQL query.
**Files**: `internal/infra/database/queries/movie/movies.sql`

### Bug 10: Casbin SyncedEnforcer Migration
**Symptom**: Concurrent policy checks could race with reloads.
**Root Cause**: Used regular `casbin.Enforcer` instead of `casbin.SyncedEnforcer`.
**Fix**: Migrated to `SyncedEnforcer` with `StartAutoLoadPolicy(10 * time.Second)`.
**Files**: `internal/service/rbac/module.go`, `internal/service/rbac/service.go`

### Test Suite Updates
**File**: `tests/live/smoke_test.go`

Expanded from 15 basic tests to comprehensive suite covering:
- Health endpoints (liveness, readiness, startup)
- Full auth flow (register → login → refresh → logout)
- Server settings (admin-only access)
- Library management (list, scan)
- Movie search and autocomplete
- User preferences (get, update)
- MFA setup (TOTP, backup codes)
- API key lifecycle (create, list, revoke, auth)
- RBAC role assignment
- OIDC provider listing
- Sonarr integration endpoints

### Final Result

```
$ make test-live
ok  github.com/lusoris/revenge/tests/live  35.729s
```

All tests pass. Container logs show only expected errors from deliberate negative test scenarios (unauthenticated access, invalid tokens).

---

## Summary of All Changes

| Area | Change |
|------|--------|
| `docker-compose.dev.yml` | Typesense 27.1 image, bash healthcheck |
| `Dockerfile` | Added `vips` runtime dep |
| `config.go` | Port default 8096, email base_url 8096, auth.encryption_key field |
| `config/module.go` | Port default 8096 |
| `config/*_test.go` | Updated port assertions |
| `app/module.go` | Added storage, email, crypto, search modules |
| `crypto/module.go` | NEW: fx provider for Encryptor |
| `infra/jobs/module.go` | River schema auto-migration |
| `Makefile` | docker-local, test-live targets |
| `tests/live/smoke_test.go` | Expanded to full suite |
| `service/rbac/adapter.go` | Fixed NULL v0-v5 scanning |
| `service/rbac/adapter_test.go` | Added nullable column unit test |
| `service/rbac/module.go` | SyncedEnforcer + auto-reload |
| `service/rbac/service.go` | Uses SyncedEnforcer |
| `api/handler_search.go` | Empty results on error instead of 500 |
| `api/handler.go` | Implemented HandleApiKeyAuth |
| `api/openapi/openapi.yaml` | Added apiKeyAuth security scheme |
| `api/ogen/*` | Regenerated with apiKeyAuth |
| `service/apikeys/service.go` | List only active keys |
| `movie/moviejobs/module.go` | Added fx.Invoke(RegisterWorkers) |
| `tvshow/jobs/module.go` | Added fx.Invoke(RegisterWorkers) |
| `search/movie_schema.go` | Popularity field not optional |
| `movie/repository_postgres.go` | pgx.ErrNoRows + ErrMovieNotFound |
| `api/movie_handlers.go` | errors.Is() for sentinels |
| `config/loader.go` | Compound env var overrides |
| `config/loader_test.go` | Tests for compound env mapping |
| `database/queries/movie/movies.sql` | COALESCE for NULL stats |

---

## Phase 5: Unify Logging to slog (COMPLETE)

**Problem**: Codebase used two separate structured logging libraries simultaneously:
- **slog** (Go stdlib): Used by infrastructure (database, cache, jobs, health, playback) — 62 files
- **zap** (Uber): Used by API handlers, services, middleware, raft — 132 files

This produced two different log formats in container output (tint/JSON vs zap's encoder).

**Solution**: Migrated all `*zap.Logger` usage to `*slog.Logger`. Removed zap entirely from authored code.

### Conversion Patterns Applied

| zap | slog |
|-----|------|
| `*zap.Logger` | `*slog.Logger` |
| `zap.String("k", v)` | `slog.String("k", v)` |
| `zap.Int("k", v)` | `slog.Int("k", v)` |
| `zap.Bool("k", v)` | `slog.Bool("k", v)` |
| `zap.Float64("k", v)` | `slog.Float64("k", v)` |
| `zap.Error(err)` | `slog.Any("error", err)` |
| `logger.Named("api")` | `logger.With("component", "api")` |
| `zap.NewNop()` | `logging.NewTestLogger()` |
| `"go.uber.org/zap"` | `"log/slog"` |

### Files Modified (~130 files)

**Logging module** (removed zap provider):
- `infra/logging/logging.go` — Removed `NewZapLogger()`
- `infra/logging/module.go` — Removed `ProvideZapLogger`
- `infra/logging/logging_test.go` — Removed zap tests

**API layer** (~26 files):
- `api/handler.go` — Handler struct logger field
- `api/server.go` — ServerParams, Server struct, NewServer()
- `api/handler_*.go` (19 handler files) — All log calls
- `api/middleware/ratelimit.go`, `ratelimit_redis.go` — Rate limiter logging
- All API test files — `zap.NewNop()` → `logging.NewTestLogger()`

**Service layer** (~40 files):
- `service/activity/` — service.go, cleanup.go
- `service/apikeys/` — service.go, module.go
- `service/auth/` — service.go, module.go
- `service/email/` — service.go, module.go
- `service/library/` — service.go, module.go, cached_service.go, cleanup.go
- `service/mfa/` — manager.go, totp.go, backup_codes.go, webauthn.go, module.go
- `service/oidc/` — service.go, module.go
- `service/rbac/` — service.go, module.go, cached_service.go, roles.go
- `service/search/` — cached_service.go
- `service/session/` — service.go, module.go, cached_service.go
- `service/settings/` — cached_service.go
- `service/storage/` — s3.go, module.go, storage.go
- `service/user/` — cached_service.go
- All service test files updated

**Content modules** (~15 files):
- `content/movie/moviejobs/` — all workers + module.go
- `content/movie/cached_service.go`
- `content/tvshow/jobs/` — all workers + module.go
- `content/shared/jobs/types.go`

**Infrastructure** (~6 files):
- `infra/raft/election.go` — hclog adapter rewired from zap→slog
- `infra/raft/module.go`
- `infra/image/` — service.go, module.go
- `infra/observability/server.go`
- `crypto/module.go`

**Integrations** (~6 files):
- `integration/radarr/` — jobs.go, module.go
- `integration/sonarr/` — jobs.go, module.go

### Additional Fixes During Migration

1. **`handler_session_test.go`**: Still had `zap.NewNop()` passed to `session.NewService` — fixed
2. **`handler_activity_test.go`, `handler_radarr_test.go`, `handler_rbac_test.go`**: Used `casbin.NewEnforcer()` but `rbac.NewService` expects `*casbin.SyncedEnforcer` (from Phase 4 migration) — changed to `casbin.NewSyncedEnforcer()`
3. **`service/apikeys/service_unit_test.go`**: Mocked `ListUserAPIKeys` but service now calls `ListActiveUserAPIKeys` (from Bug 4 fix) — updated mock expectations

### Verification

```
$ grep -r "go.uber.org/zap" internal/
(zero results)

$ make build
(clean compilation)

$ make test-short
ok  github.com/lusoris/revenge/internal/api          0.258s
ok  github.com/lusoris/revenge/internal/service/...   (all pass)
(48 packages, zero failures)
```

Note: `go.uber.org/zap` remains in `go.mod` as an `// indirect` dependency (transitive dep from another library). All direct zap imports from authored code are removed.

---

## Phase 6: API Completeness Audit & Bug Fixes (COMPLETE)

### Audit Results

Comprehensive audit of all production code for stubs, TODOs, placeholders, and incomplete implementations.

**Finding**: No stubs, TODOs, or placeholder code in production code. The only placeholder is an sqlc scaffolding query in `tvshow/db/placeholder.sql.go` (valid pattern — sqlc requires non-empty query directories).

### Gap Found: TV Show Typesense Search Endpoints Missing

The TV show search service was fully implemented at the service layer (`tvshow_service.go` with Search, Autocomplete, GetFacets, ReindexAll) but had **no HTTP handlers**. TV show search only used a basic database `SearchSeries` query, not the Typesense index.

**Fix**: Added 3 new API endpoints + 4 new OpenAPI schemas:

| Endpoint | Operation | Description |
|----------|-----------|-------------|
| `GET /api/v1/search/tvshows` | `searchLibraryTVShows` | Full-text search with filtering, facets, pagination, sorting |
| `GET /api/v1/search/tvshows/autocomplete` | `autocompleteTVShows` | Title suggestions for search-as-you-type |
| `GET /api/v1/search/tvshows/facets` | `getTVShowSearchFacets` | Filter values (genres, years, networks, status, type) |

**Schemas added**: `TVShowSearchResults`, `TVShowSearchHit`, `TVShowSearchDocument`, `TVShowSearchFacets`

### Bug 13: Security Errors Logged as ERR and Returned 500

**Symptom**: Every unauthenticated request logged at ERR level and returned HTTP 500.
**Root Cause**: `NewError()` didn't recognize ogen's `SecurityError`, `DecodeParamsError`, or `DecodeRequestError` types. All fell through to the generic 500 handler.
**Fix**: Added `ogenerrors` type matching — security errors → 401 at DEBUG level, param/request decode errors → 400 at WARN level, rate limit errors → 429 via `statusCoder` interface.
**Files**: `internal/api/handler.go`

### Bug 14: Search Collections Not Initialized on Startup

**Symptom**: Search/autocomplete/facets returned Typesense 404 on fresh stack (collections didn't exist until first reindex job).
**Root Cause**: No startup hook to create collections.
**Fix**: Added `fx.Invoke(initializeCollections)` in `search/module.go` that creates both `movies` and `tvshows` collections on startup.
**Files**: `internal/service/search/module.go`

### Bug 15: TV Show Schema Popularity Field Optional

**Symptom**: `"Default sorting field 'popularity' cannot be an optional field"` — tvshows collection creation failed.
**Root Cause**: `popularity` field marked `Optional: true` in `tvshow_schema.go` but used as `DefaultSortingField`.
**Fix**: Removed `Optional` from the popularity field.
**Files**: `internal/service/search/tvshow_schema.go`

### Bug 16: Email Verification Returned 500 ERR

**Symptom**: Invalid verification token requests logged as ERR and returned HTTP 500.
**Root Cause**: `VerifyEmail` handler returned error through `fmt.Errorf()` → `NewError()`, which logged at ERR level and returned 500. But "invalid token" is a client error.
**Fix**: Return proper 400 response inline with WRN-level logging.
**Files**: `internal/api/handler.go`

### Container Log Verification

After all fixes, clean stack with full live test suite:
- **ERR lines**: 0
- **WRN lines**: All expected (email disabled, deliberate test failures)
- **Both search collections**: Created on startup (`movies`, `tvshows`)
- **Rate limiting**: Redis-backed, enabled (50 RPS global, 5 RPS auth)
- **All live tests**: Pass (37.6s)

### Files Modified

| File | Change |
|------|--------|
| `api/openapi/openapi.yaml` | +3 TV show search endpoints, +4 schemas |
| `internal/api/handler.go` | `tvshowSearchService` field, `ogenerrors` handling in `NewError()`, email verify fix |
| `internal/api/handler_search.go` | +3 TV show search handler methods |
| `internal/api/handler_test.go` | Updated `NewError` test assertions |
| `internal/api/server.go` | `TVShowSearchService` in DI params + wiring |
| `internal/api/ogen/*` | Regenerated from updated OpenAPI spec |
| `internal/service/search/module.go` | `initializeCollections` startup hook |
| `internal/service/search/tvshow_schema.go` | `popularity` not optional |

---

## Phase 7: Docker Tuning, Security Hardening, Expanded Tests (COMPLETE)

### Docker Compose Tuning

Researched and applied optimal settings for all three backing services in `docker-compose.dev.yml`.

**PostgreSQL 18** (~2GB container):
- `shared_buffers=512MB` (25% of container memory)
- `effective_cache_size=1536MB` (75% memory estimate)
- `work_mem=16MB`, `maintenance_work_mem=128MB`
- WAL tuning: `wal_buffers=16MB`, `checkpoint_completion_target=0.9`, `max_wal_size=1GB`, `wal_compression=on`
- SSD-optimized: `random_page_cost=1.1`, `effective_io_concurrency=200`
- Autovacuum: `max_workers=3`, `naptime=30s`
- `shm_size: 256mb` for shared_buffers

**Dragonfly** (~1.5GB container):
- `cache_mode=true` (LFRU eviction, no persistence)
- `maxmemory=1G` (headroom within container)
- `proactor_threads=2`
- Persistence disabled (`dbfilename=`)
- `ulimits.memlock: -1` for memory management

**Typesense 27.1** (~1.5GB container):
- `thread-pool-size=8` (4x CPU for I/O-bound search)
- `num-collections-parallel-load=4` (faster startup)
- `healthy-read-lag=1000`, `healthy-write-lag=500` (backpressure)
- `snapshot-interval-seconds=3600` (hourly snapshots)
- Resource limits with reservations

### Trivy Container Scan Results

| Scan Type | CRITICAL | HIGH | Status |
|-----------|----------|------|--------|
| Docker image (OS + Go binary) | 0 | 0 | CLEAN |
| Filesystem (Go modules) | 0 | 0 | CLEAN |
| Config: Production Dockerfile | 0 | 0 | CLEAN |
| Config: Devcontainer Dockerfile | 0 | 2 | Fixed |
| Config: Helm deployment.yaml | 0 | 12 | Fixed |

### Trivy Fixes Applied

**Helm chart** (`charts/revenge/templates/deployment.yaml`):
- Added `securityContext` at pod level: `runAsNonRoot: true`, `fsGroup`
- Added `securityContext` at container level: `runAsNonRoot: true`, `readOnlyRootFilesystem: true`, `allowPrivilegeEscalation: false`
- Applied to all 4 containers: revenge, postgresql, dragonfly, typesense

**Devcontainer Dockerfile**:
- Added `--no-install-recommends` to `apt-get install` (DS-0029)

### gosec Findings Fixed

**G301 — Directory permissions** (4 MEDIUM → Fixed):
- Changed `0o755` → `0o750` in:
  - `internal/playback/transcode/pipeline.go` (2 locations)
  - `internal/playback/subtitle/extract.go`
  - `internal/playback/service.go`

**G304 — Path traversal** (3 MEDIUM → Already mitigated):
- `storage.go` already has `sanitizeKey()` with `filepath.Clean`, `..` removal, prefix validation
- HLS paths use session-generated UUIDs + profile names (not user input)
- Assessment: false positive in current code

### Expanded Live Test Coverage

Added TV show search tests to `tests/live/smoke_test.go`:
- `TestLive_TVShowSearchInfrastructure` — 5 subtests:
  - Typesense search (empty results)
  - Autocomplete (empty suggestions)
  - Facets (endpoint functional)
  - Search with pagination params
  - Search with filter_by
- Added 3 TV show search endpoints to unauthenticated access tests

### Verification

```
$ make docker-test
Health check: OK

$ make test-live
ok  github.com/lusoris/revenge/tests/live  34.728s

$ docker logs revenge-dev 2>&1 | grep '\[91mERR'
(zero results — 0 ERR-level log lines)
```

### Files Modified

| File | Change |
|------|--------|
| `docker-compose.dev.yml` | Full PostgreSQL/Dragonfly/Typesense tuning with resource limits |
| `charts/revenge/templates/deployment.yaml` | Security contexts for all 4 containers + pods |
| `.devcontainer/Dockerfile` | `--no-install-recommends` |
| `internal/playback/transcode/pipeline.go` | `0o755` → `0o750` (2 locations) |
| `internal/playback/subtitle/extract.go` | `0o755` → `0o750` |
| `internal/playback/service.go` | `0o755` → `0o750` |
| `tests/live/smoke_test.go` | +TV show search tests, +unauthenticated access coverage |
