# Full Work Plan - Metadata System + CI Fixes

## Status Legend
- [ ] Not started
- [x] Completed
- [~] In progress

---

## 1. METADATA SYSTEM - Force + Languages Plumbing

### 1.1 Add MetadataRefreshOptions to content modules
- [x] `internal/content/movie/metadata_provider.go` - Add `MetadataRefreshOptions` struct
- [x] `internal/content/tvshow/metadata_provider.go` - Add `MetadataRefreshOptions` struct

### 1.2 Update Service interfaces + implementations
- [x] `internal/content/movie/service.go` - Change `RefreshMovieMetadata(ctx, id)` to `RefreshMovieMetadata(ctx, id, opts ...MetadataRefreshOptions)`
- [x] `internal/content/movie/service.go` - Passes opts through to adapter EnrichMovie
- [x] `internal/content/tvshow/service.go` - Change `RefreshSeriesMetadata/Season/Episode` to variadic opts
- [x] `internal/content/tvshow/service.go` - Passes opts through to adapter Enrich methods

### 1.3 Update MetadataProvider interfaces to support per-request options
- [x] `internal/content/movie/metadata_provider.go` - `EnrichMovie(ctx, mov, opts ...MetadataRefreshOptions)`
- [x] `internal/content/movie/metadata_provider.go` - `ClearCache()` method added
- [x] `internal/content/tvshow/metadata_provider.go` - `EnrichSeries/Season/Episode` with opts
- [x] `internal/content/tvshow/metadata_provider.go` - `ClearCache()` method added

### 1.4 Add ClearCache to shared metadata Provider + Service
- [x] `internal/service/metadata/provider.go` - `ClearCache()` added to `Provider` base interface
- [x] `internal/service/metadata/service.go` - `ClearCache()` added to `Service` interface
- [x] `internal/service/metadata/service.go` - Implementation delegates to all registered providers

### 1.5 Update Adapters to support per-request languages + ClearCache
- [x] `internal/service/metadata/adapters/movie/adapter.go` - `EnrichMovie` uses opts languages when provided, clears cache on Force
- [x] `internal/service/metadata/adapters/movie/adapter.go` - `ClearCache()` delegates to service
- [x] `internal/service/metadata/adapters/tvshow/adapter.go` - `EnrichSeries/Season/Episode` uses opts languages, clears cache on Force
- [x] `internal/service/metadata/adapters/tvshow/adapter.go` - `ClearCache()` delegates to service

### 1.6 Update Workers to pass Force + Languages through
- [x] `internal/content/movie/moviejobs/metadata_refresh.go` - Passes Force/Languages as MetadataRefreshOptions to service
- [x] `internal/content/tvshow/jobs/jobs.go` - MetadataRefreshWorker passes Force as opts
- [x] `internal/content/tvshow/jobs/jobs.go` - SeriesRefreshWorker passes Languages as opts

### 1.7 Handler changes
- [x] No handler changes needed - variadic `opts...` is backward compatible, handlers call with no opts

---

## 2. BROKEN TESTS - Fixed

### 2.1 Fix moviejobs test file
- [x] `internal/content/movie/moviejobs/metadata_refresh_test.go` - Complete rewrite:
  - Removed tests for deleted `formatTimePtr` and `formatDecimalPtr`
  - Removed tests for deleted `MovieMetadataRefreshArgs` type (now uses `metadatajobs.RefreshMovieArgs`)
  - Updated worker field assertions from `worker.movieRepo` to `worker.service` and `worker.jobClient`

### 2.2 Fix mock signatures
- [x] `internal/content/movie/handler_test.go` - MockService.RefreshMovieMetadata updated to variadic opts
- [x] `internal/content/movie/library_service_test.go` - MockMetadataProvider updated:
  - EnrichMovie signature updated to variadic opts
  - ClearCache() method added
- [x] `internal/content/movie/service_test.go` - All `NewService(repo)` → `NewService(repo, nil)`, assertion updated
- [x] `internal/content/tvshow/service_test.go` - All `NewService(repo)` → `NewService(repo, nil)`, assertion updated

---

## 3. CI WORKFLOW - CGO Dependencies

### 3.1 Fix govulncheck in ci.yml
- [x] `.github/workflows/ci.yml` - Added CGO dependencies (libvips-dev, libav*-dev, pkg-config) to `vuln` job

---

## 4. CI WORKFLOW - Migrations Path

### 4.1 Fix testutil migrations path
- [x] `internal/testutil/testdb_migrate.go` - Changed from `filepath.Join(projectRoot, "migrations")` to `filepath.Join(projectRoot, "internal", "infra", "database", "migrations", "shared")`

### 4.2 Fix database.go migrations path
- [x] `internal/testutil/database.go` - Fixed `findMigrationsPath()` to use correct relative paths

### 4.3 Fix migrations_test.go paths
- [x] `internal/infra/database/migrations_test.go` - Fixed all `file://../../../migrations` to `file://migrations/shared`

---

## 5. PROGRESS REPORTING - Workers

### 5.1 Progress already injected (from previous session)
- [x] TV show LibraryScanWorker
- [x] TV show MetadataRefreshWorker
- [x] TV show SeriesRefreshWorker
- [x] Movie MetadataRefreshWorker

### 5.2 Workers that DON'T need progress (short/single-item jobs)
- FileMatchWorker (both movie and tvshow) - single file, fast
- SearchIndexWorker (single movie index) - fast
- Cleanup workers - maintenance, not user-facing
- Webhook workers - event processing, not user-facing

---

## 6. COMPILE CHECK

- [x] `go build` passes for all changed packages
- [x] `go vet` passes for all changed packages (including tests)

---

## 7. COMMIT + PUSH

- [ ] Stage all changes
- [ ] Commit with proper conventional commit message
- [ ] Push to develop
