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
