# Bug Report - Pre-existing Issues Found

## BUG-001: Force parameter is dead throughout metadata stack [FIXED]
**Severity**: HIGH
**Status**: FIXED
**Files affected**:
- `internal/content/movie/service.go` - `RefreshMovieMetadata` ignores Force
- `internal/content/tvshow/service.go` - `RefreshSeriesMetadata` ignores Force
- `internal/service/metadata/adapters/movie/adapter.go` - `ClearCache()` was no-op
- `internal/service/metadata/adapters/tvshow/adapter.go` - `ClearCache()` was no-op
- All metadata refresh workers log Force but never pass it

**Fix applied**:
- Added `MetadataRefreshOptions{Force, Languages}` to both content modules
- Service methods now accept variadic opts and pass through to adapter
- Adapters check `opts.Force` and call `service.ClearCache()` before fetching
- `ClearCache()` added to `metadata.Provider` interface and `metadata.Service` interface
- Both adapters delegate `ClearCache()` to the shared service which calls all providers
- Workers now construct `MetadataRefreshOptions` from job args and pass to service

---

## BUG-002: Languages parameter is dead throughout metadata stack [FIXED]
**Severity**: MEDIUM
**Status**: FIXED
**Files affected**:
- Same as BUG-001 plus all adapter `Enrich*` methods
- `internal/service/metadata/adapters/movie/adapter.go` - used `a.languages` (hardcoded at init)
- `internal/service/metadata/adapters/tvshow/adapter.go` - same

**Fix applied**:
- Adapters check `opts.Languages` and use them when provided, falling back to `a.languages`
- Workers pass `Languages` from job args to service as `MetadataRefreshOptions`
- Server-configured languages remain the default; per-request override is now possible

---

## BUG-003: CI govulncheck missing CGO dependencies [FIXED]
**Severity**: MEDIUM
**Status**: FIXED
**File**: `.github/workflows/ci.yml` vuln job

**Fix applied**: Added CGO dependency install step (libvips-dev, libav*-dev, pkg-config) matching security.yml's pattern.

---

## BUG-004: Tests reference wrong migrations path [FIXED]
**Severity**: HIGH
**Status**: FIXED
**Files affected**:
- `internal/testutil/testdb_migrate.go` - was `{projectRoot}/migrations/`
- `internal/testutil/database.go` - relative paths all wrong
- `internal/infra/database/migrations_test.go` - all `file://../../../migrations`

**Actual path**: `internal/infra/database/migrations/shared/`

**Fix applied**:
- `testdb_migrate.go`: Changed to `filepath.Join(projectRoot, "internal", "infra", "database", "migrations", "shared")`
- `database.go`: Updated `findMigrationsPath()` to search for correct relative paths
- `migrations_test.go`: Changed all references to `file://migrations/shared`

---

## BUG-005: metadata_refresh_test.go references deleted code [FIXED]
**Severity**: LOW (test-only)
**Status**: FIXED
**File**: `internal/content/movie/moviejobs/metadata_refresh_test.go`

**Fix applied**: Complete rewrite - tests now use `metadatajobs.RefreshMovieArgs` and check `worker.service`/`worker.jobClient` fields.

---

## BUG-006: All service tests use wrong NewService signature [FIXED]
**Severity**: HIGH (test-only, discovered during this session)
**Status**: FIXED
**Files affected**:
- `internal/content/movie/service_test.go` - ~38 calls to `NewService(repo)` with 1 arg
- `internal/content/tvshow/service_test.go` - ~15 calls to `NewService(repo)` with 1 arg

**Root cause**: `NewService` was changed to `NewService(repo, metadataProvider)` but tests were never updated.

**Fix applied**: All calls changed to `NewService(repo, nil)`. Assertions changed from "not implemented" to "metadata provider not configured".

---

## BUG-007: library_service_test.go mock missing ClearCache [FIXED]
**Severity**: LOW (test-only, discovered during this session)
**Status**: FIXED
**File**: `internal/content/movie/library_service_test.go`

**Fix applied**: Added `ClearCache()` to MockMetadataProvider and updated `EnrichMovie` signature.
