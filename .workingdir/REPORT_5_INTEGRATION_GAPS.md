# Integration Gaps Report (2026-02-05)

**Generated after Phase A12/A12b completion**
**Updated: 2026-02-06 - TVShow jobs fully implemented (no more stubs)**
**Updated: 2026-02-06 - Sonarr API endpoints fully implemented**

---

## Summary Table

| Integration Component | Movie | TVShow | Radarr | Sonarr |
|---|---|---|---|---|
| Module exists | ✓ | ✓ | ✓ | ✓ |
| Module in app init | ✓ | ✓ | ✓ | ✓ |
| Jobs module exists | ✓ | ✓ | ✓ | ✓ |
| Jobs in app init | ✓ | ✓ | ✓ | ✓ |
| Service implementation | ✓ | ✓ | ✓ | ✓ |
| API Handlers | ✓ | ✓ | ✓ | ✓ |
| Handler integration | ✓ | ✓ | ✓ | ✓ |
| Metadata refresh impl | ✓ | ✓ | ✓ | ✓ |
| Feature parity | — | — | COMPLETE | COMPLETE |

**All Radarr/Sonarr integrations are now at feature parity**

---

## Completed Fixes (2026-02-05)

### 1. ✅ TVShow Module Added to App Initialization

**File:** `internal/app/module.go`
- Added `tvshow.Module` to Content Modules section

### 2. ✅ Sonarr Module Added to App Initialization

**File:** `internal/app/module.go`
- Added `sonarr.Module` to Integrations section

### 3. ✅ TVShow Jobs Module Created and Added

**Files:**
- Created `internal/content/tvshow/jobs/module.go`
- Added `tvshowjobs.Module` to Job Workers section in app/module.go

### 4. ✅ Sonarr API Handler Interface Created

**File:** `internal/api/handler_sonarr.go`
- Created `sonarrService` interface matching `radarrService` pattern
- Ready for API endpoint implementation when OpenAPI spec is updated

### 5. ✅ Sonarr Service Wired into Server

**Files:**
- `internal/api/server.go` - Added `SonarrService` to `ServerParams`
- `internal/api/handler.go` - Added `sonarrService` field to `Handler`
- `internal/api/server.go` - Wired `SonarrService` in `NewServer()`

### 6. ✅ Metadata Service Module Added

**File:** `internal/app/module.go`
- Added `metadatafx.Module` to Metadata Service section

### 7. ✅ Metadata Refresh Implemented (Movie)

**File:** `internal/content/movie/service.go`
- Added `metadataProvider MetadataProvider` to `movieService` struct
- Implemented `RefreshMovieMetadata`:
  - Gets movie from database
  - Enriches via `MetadataProvider.EnrichMovie`
  - Updates movie in database
  - Refreshes credits and genres from TMDb

**File:** `internal/content/movie/module.go`
- Updated to inject `MetadataProvider` into service

### 8. ✅ Metadata Refresh Implemented (TVShow)

**File:** `internal/content/tvshow/service.go`
- Added `metadataProvider MetadataProvider` to `tvService` struct
- Implemented `RefreshSeriesMetadata`:
  - Enriches series via `MetadataProvider.EnrichSeries`
  - Updates credits and genres from TMDb
- Implemented `RefreshSeasonMetadata`:
  - Enriches season via `MetadataProvider.EnrichSeason`
- Implemented `RefreshEpisodeMetadata`:
  - Enriches episode via `MetadataProvider.EnrichEpisode`

**File:** `internal/content/tvshow/module.go`
- Updated to inject `MetadataProvider` into service

---

## Remaining Issues (Low Priority)

### 1. ~~TVShow Jobs Incomplete~~ COMPLETED (2026-02-06)

**Location:** `internal/content/tvshow/jobs/jobs.go`

All TODO placeholders have been replaced with full implementations:
- ✅ **LibraryScanWorker** - Uses shared scanner with TVShowFileParser, walks directories, parses filenames (SxxExx patterns), auto-creates series/season/episode when enabled
- ✅ **MetadataRefreshWorker** - Full batch refresh with pagination (default 50 per batch), refreshes all series when no ID specified
- ✅ **FileMatchWorker** - Parses filenames, searches/creates series via TMDb, creates season/episode records, links files
- ✅ **SearchIndexWorker** - Gracefully handles missing search service, validates series exist, ready for Typesense integration
- ✅ **SeriesRefreshWorker** - Uses service.RefreshSeriesMetadata, cascade refresh for seasons and episodes when requested

### 2. ~~Sonarr API Endpoints Not in OpenAPI Spec~~ COMPLETED (2026-02-06)

All Sonarr API endpoints have been added to OpenAPI spec and handlers implemented:
- ✅ `GET /api/v1/admin/integrations/sonarr/status` - `AdminGetSonarrStatus`
- ✅ `POST /api/v1/admin/integrations/sonarr/sync` - `AdminTriggerSonarrSync`
- ✅ `GET /api/v1/admin/integrations/sonarr/quality-profiles` - `AdminGetSonarrQualityProfiles`
- ✅ `GET /api/v1/admin/integrations/sonarr/root-folders` - `AdminGetSonarrRootFolders`
- ✅ `POST /api/v1/webhooks/sonarr` - `HandleSonarrWebhook`

**Implementation Details:**
- OpenAPI spec updated with Sonarr schemas: `SonarrStatus`, `SonarrSyncStatus`, `SonarrSyncResponse`, `SonarrQualityProfile`, `SonarrRootFolder`, `SonarrWebhookPayload`, etc.
- ogen regenerated successfully
- Handlers in `internal/api/handler_sonarr.go` follow exact Radarr pattern
- Full webhook payload conversion for Series, Episodes, EpisodeFiles, Releases

---

## Action Items (Updated)

### ~~Priority 1 (Critical - Must Fix)~~ COMPLETED

1. [x] Add TVShow module to app initialization
2. [x] Add Sonarr module to app initialization
3. [x] Create TVShow jobs module
4. [x] Add Metadata service module to app initialization

### ~~Priority 2 (High - API Parity)~~ COMPLETED

5. [x] Create `handler_sonarr.go` with service interface
6. [x] Add `SonarrService` to `ServerParams`
7. [x] Add `sonarrService` to `Handler` struct
8. [x] Wire Sonarr service in `NewServer()`

### ~~Priority 3 (Medium - Functionality)~~ COMPLETED

9. [x] Implement metadata refresh in movie service
10. [x] Implement metadata refresh in tvshow service
11. [x] Complete TVShow jobs implementations ✅ DONE (2026-02-06)
12. [x] Add Sonarr endpoints to OpenAPI spec ✅ DONE (2026-02-06)

### ~~Priority 4 (Low - Code Quality)~~ COMPLETED

13. [x] Remove or implement all TODO comments in tvshow jobs ✅ DONE (2026-02-06)

---

## Files Modified

| File | Action | Status |
|------|--------|--------|
| `internal/app/module.go` | Add tvshow, sonarr, tvshowjobs, metadatafx modules | ✅ DONE |
| `internal/api/server.go` | Add SonarrService to ServerParams + wiring | ✅ DONE |
| `internal/api/handler.go` | Add sonarrService field | ✅ DONE |
| `internal/api/handler_sonarr.go` | CREATE - sonarrService interface + full handlers | ✅ DONE |
| `internal/content/tvshow/module.go` | UPDATE - inject MetadataProvider | ✅ DONE |
| `internal/content/tvshow/jobs/module.go` | CREATE - fx module | ✅ DONE |
| `internal/content/movie/service.go` | Implement RefreshMovieMetadata | ✅ DONE |
| `internal/content/movie/module.go` | UPDATE - inject MetadataProvider | ✅ DONE |
| `internal/content/tvshow/service.go` | Implement Refresh*Metadata methods | ✅ DONE |
| `internal/content/tvshow/jobs/jobs.go` | Implement all job workers (no stubs) | ✅ DONE |
| `internal/content/tvshow/jobs/module.go` | UPDATE - provider functions for optional MetadataProvider | ✅ DONE |
| `api/openapi/openapi.yaml` | Add Sonarr endpoints and schemas | ✅ DONE |
| `internal/api/ogen/*` | Regenerated with `make ogen` | ✅ DONE |

---

## Build Verification

All changes verified with successful `go build ./...` - no errors.
