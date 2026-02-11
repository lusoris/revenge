# Bugs Found During Live Test Development

## Bug 1: Empty 404 Response Bodies (FIXED)
**File**: `internal/api/movie_handlers.go`, `internal/api/tvshow_handlers.go`
**Issue**: Handlers returned empty ogen `NotFound` structs like `&ogen.GetMovieNotFound{}`
**Result**: API returned `{"code":0,"message":""}` instead of `{"code":404,"message":"Movie not found"}`
**Fix**: Added `OgenNotFound()` helper in `internal/api/errors.go` and updated 32+ handlers to return properly populated error responses
**Status**: FIXED (this session)

## Bug 2: CreateMovieFile SQL Missing file_name (FIXED - commit d3885a7f)
**File**: `internal/infra/database/queries/movie/movies.sql`
**Issue**: INSERT statement was missing the `file_name` column which has NOT NULL constraint
**Status**: FIXED

## Bug 3: CreateEpisodeFile nil arrays (FIXED - commit bc4f0e5c)
**File**: `internal/content/tvshow/repository_postgres.go`
**Issue**: Passing nil slices for audio_languages/subtitle_languages caused postgres errors
**Status**: FIXED (defaulting to []string{})

## Bug 4: GetSeriesWatchStats NULL SUM (FIXED - commit bc4f0e5c)
**File**: `internal/infra/database/queries/tvshow/watch_progress.sql`
**Issue**: SUM(watch_count) returns NULL when no rows, causing scan errors
**Status**: FIXED (COALESCE(..., 0)::bigint)

---

# Provider Coupling Audit — TMDb Hardcoding Issues

These are systemic design issues where the codebase is tightly coupled to TMDb
as a metadata/image provider instead of being provider-agnostic. The Provider
interfaces (provider.go) are well-designed, but the Service layer, database
schema, and API surface re-introduce TMDb-specific coupling.

## Issue 1: Image System Hardcoded to TMDb (HIGH) — FIXED (commit 603c026c)
**Problem**: The image proxy endpoint (`/api/v1/images/{type}/{size}/{path}`)
only works with TMDb. Fanart.tv returns full URLs — piping them through
`GetImageURL()` produced broken double-prefixed URLs.

**Fix applied**: `GetURL()` and all size-specific methods in `images.go` now
detect full URLs (http/https prefix) and pass through unchanged. `GetImageURL()`
in `image/service.go` does the same via `isFullURL()` helper. Tests added for
fanart.tv, Radarr, and http URLs.
**Status**: FIXED

## Issue 2: Genre System Uses `tmdb_genre_id` as Universal Key (HIGH)
**Problem**: Genre IDs from TMDb (28=Action, 35=Comedy) are used as the sole
deduplication key everywhere. Non-TMDb providers must fake TMDb IDs:
- AniList/Kitsu genres like "Ecchi", "Mecha", "Slice of Life" get invented
  IDs in the 90000+ range
- Multiple genres collapse (AniList "Action" + "Adventure" → TMDb 10759)
- TMDb movie vs TV genre IDs overlap inconsistently (movie Action=28, TV=10759)

**Files affected (~30 locations)**:
- DB migrations: `000025_create_movie_genres_table.up.sql`, `000032_create_tvshow_schema.up.sql`
- SQL queries: `queries/movie/movies.sql` (5 refs), `queries/tvshow/genres.sql` (5 refs)
- Generated sqlc: `content/movie/db/movies.sql.go`, `content/tvshow/db/genres.sql.go`
- Domain types: `internal/content/types.go` (`GenreSummary.TMDbGenreID`)
- OpenAPI: `api/openapi/openapi.yaml` (Genre schema requires `tmdb_genre_id`)
- Provider mappings: anilist, kitsu, tvmaze, mal, letterboxd all have hardcoded genre→TMDb ID maps

**Fix**: Create a first-class `genres` table with canonical slugs + a
`genre_external_ids` junction table for provider-specific ID mappings. Requires
DB migration + touching ~30 files.

## Issue 3: Metadata Service Interface Takes `tmdbID int32` (CRITICAL) — FIXED
**Problem**: The `Service` interface accepted `tmdbID int32` for all ~20 methods,
forcing TMDb-specific coupling. Provider interfaces already used `id string`.

**Fix applied**: Changed all Service interface methods from `tmdbID int32` to
`id string`. Removed `id := fmt.Sprintf("%d", tmdbID)` conversions from all
implementations. Updated all callers:
- `internal/service/metadata/service.go` — interface + 20 implementations
- `internal/api/handler_metadata.go` — 19 handler calls now use `strconv.Itoa()`
- `internal/content/movie/metadata_provider.go` — `tmdbID int` → `providerID string`
- `internal/content/tvshow/metadata_provider.go` — `seriesTMDbID int32` → `seriesProviderID string`
- `internal/service/metadata/adapters/movie/adapter.go` — updated all method sigs
- `internal/service/metadata/adapters/tvshow/adapter.go` — updated all method sigs
- `internal/content/movie/service.go`, `library_service.go`, `library_matcher.go`
- `internal/content/tvshow/service.go`, `jobs/jobs.go`
- `internal/service/metadata/jobs/queue.go`, `refresh.go`, `workers.go`
- All test mocks updated across 5 test files
**Status**: FIXED

## Issue 4: API Endpoints Use `{tmdbId}` Path Parameter (HIGH)
**Problem**: All 19 `/api/v1/metadata/*` endpoints use `{tmdbId}` in the path,
coupling the public API surface to TMDb. Changing this is a breaking change for
all API clients.

**Endpoints**: `/api/v1/metadata/movie/{tmdbId}`, `.../credits`, `.../images`,
`.../similar`, `.../recommendations`, `.../external-ids`, and the equivalent
TV/person/collection endpoints (19 total).

**Fix**: Rename `{tmdbId}` to `{id}` + add optional `?provider=tmdb` query
parameter (defaulting to `tmdb` for backwards compat). Or use
`/api/v1/metadata/movie/{provider}:{id}` URI format.

## Issue 5: Image Paths Stored as TMDb-Relative Fragments (MEDIUM) — FIXED (with #1)
**Problem**: TMDb-relative paths vs full URLs from Radarr/fanart.tv in the
same DB columns.

**Fix applied**: Covered by Issue #1 fix — `GetURL()` and `GetImageURL()` now
detect full URLs and pass through unchanged. Both storage formats work.
**Status**: FIXED

## Issue 6: External IDs Not Normalized (MEDIUM)
**Problem**: `tmdb_id` is a top-level field on Movie, Series, Episode, Network
structs — not in a generic external ID map. Each entity type has ad-hoc
external ID columns (`imdb_id`, `tvdb_id`, `sonarr_id`) with inconsistent
indexing. The search schemas also index on `tmdb_id` specifically.

**Notable**: The `networks` table has `tmdb_id INTEGER UNIQUE NOT NULL` — the
only table where TMDb ID is required, making it impossible to create a network
from a non-TMDb provider.

**Fix**: Add an `external_ids` junction table or JSONB column for arbitrary
provider→ID mappings. Make `tmdb_id` nullable on networks.

## Issue 7: Refresh Jobs Tied to TMDb ID (LOW) — FIXED (with #3)
**Problem**: `RefreshPersonArgs` had `TMDbID int32` in job payloads.

**Fix applied**: Changed to `ProviderID string` in `refresh.go`, `queue.go`,
and `workers.go`.
**Status**: FIXED

---

## Summary Priority Matrix

| Issue | Severity | Effort | Order |
|-------|----------|--------|-------|
| #3 Service interface `tmdbID int32` | CRITICAL | Large | 1st — unblocks everything |
| #2 Genre `tmdb_genre_id` | HIGH | Large | 2nd — DB migration |
| #1 Image system TMDb-only | HIGH | Medium | 3rd — consolidate constants |
| #4 API `{tmdbId}` paths | HIGH | Large | 4th — breaking API change, do with v2 |
| #5 Image paths format | MEDIUM | Medium | 5th — with #1 |
| #6 External IDs | MEDIUM | Medium | 6th — with #3 |
| #7 Refresh jobs | LOW | Small | 7th — trivial after #3 |

NOTE: The `Provider` interface design and the certification/rating system are
already provider-agnostic and well-designed. The problem is the Service layer
and database schema re-introducing coupling that the Provider layer avoids.

---
Last updated: 2026-02-11
