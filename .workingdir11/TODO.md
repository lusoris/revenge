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

## Issue 1: Image System Hardcoded to TMDb (HIGH)
**Problem**: The image proxy endpoint (`/api/v1/images/{type}/{size}/{path}`)
only works with TMDb. Three independent copies of TMDb image size constants
exist across the codebase. The `ImageProvider` interface exists and fanart.tv
implements it, but the infra image `Service` and content-level `ImageURLBuilder`
bypass the interface entirely and always construct TMDb URLs. Fanart.tv returns
full URLs — piping them through `GetImageURL()` produces broken double-prefixed
URLs like `image.tmdb.org/t/p/w500/https://...`.

**Files affected**:
- `internal/infra/image/service.go` — hardcoded BaseURL + TMDb size constants
- `internal/content/shared/metadata/images.go` — duplicate TMDb constants + `NewImageURLBuilder()` hardcodes TMDb URL
- `internal/service/metadata/provider.go` — third copy of TMDb `ImageSize` constants
- `internal/service/metadata/providers/tmdb/client.go` — fourth `ImageBaseURL` constant

**Fix**: Consolidate into one `ImageURLBuilder` that delegates to the registered
`ImageProvider`. For providers returning full URLs (fanart.tv), pass through
unchanged. Remove the three duplicate constant sets.

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

## Issue 3: Metadata Service Interface Takes `tmdbID int32` (CRITICAL)
**Problem**: The `Service` interface in `service/metadata/service.go` accepts
`tmdbID int32` as the lookup parameter for every method (~20 methods). This
forces every consumer to have a TMDb ID before fetching metadata. The
underlying `Provider` interface correctly uses `id string`, but the Service
layer re-narrows it to TMDb-specific `int32`.

Non-TMDb providers (TVmaze, AniList, Kitsu) can only serve as enrichment
sources, never as primary providers.

**Files affected**:
- `internal/service/metadata/service.go` — all 20+ interface methods + implementations
- `internal/content/movie/metadata_provider.go` — param named `tmdbID`
- `internal/content/tvshow/metadata_provider.go` — `EnrichSeason(ctx, season, seriesTMDbID int32)`
- All metadata adapters that call the service

**Fix**: Change Service interface to accept `(provider ProviderID, id string)`
or a generic `ExternalRef{Provider, ID}` struct. Keep tmdb_id columns in DB but
don't make them the only lookup path.

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

## Issue 5: Image Paths Stored as TMDb-Relative Fragments (MEDIUM)
**Problem**: `poster_path`, `backdrop_path`, `profile_path`, `still_path`
columns store TMDb-relative paths (e.g. `/abc123.jpg`), not full URLs. These
only work when prefixed with `image.tmdb.org/t/p/{size}`. But Radarr and
fanart.tv store full URLs in the same columns — inconsistent format in DB.

**Files affected**: All DB image columns across movies, tvshows, seasons,
episodes, credits tables + all adapters that persist paths.

**Fix**: Either always store full URLs at write time, or store paths + an
`image_source` field and resolve at read time.

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

## Issue 7: Refresh Jobs Tied to TMDb ID (LOW)
**Problem**: Metadata refresh job payloads include `TMDbID int32`, tying the
job queue to TMDb.

**Files**: `service/metadata/jobs/refresh.go`, `workers.go`

**Fix**: Use `ExternalRef{Provider, ID}` in job payloads instead.

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
