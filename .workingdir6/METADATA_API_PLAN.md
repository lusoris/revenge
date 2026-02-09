# Metadata & Provider API Plan (2026-02-08)

## Current State

**Service layer**: 41 methods across movie, tvshow, person, collection, images, refresh
**API endpoints exposed**: ~8 of 41 (19%)
**Providers implemented**: TMDb (full), TVDb (TV-only, partial images)
**Providers defined but not implemented**: Fanart.tv, OMDb

### What IS Exposed

| Endpoint | Method |
|----------|--------|
| `GET /api/v1/metadata/search/movie` | `SearchMovie()` (TMDb only) |
| `GET /api/v1/metadata/search/tv` | `SearchTVShow()` (TMDb only) |
| `GET /api/v1/metadata/movie/{tmdbId}` | `GetMovieMetadata()` |
| `GET /api/v1/metadata/tv/{tmdbId}` | `GetTVShowMetadata()` |
| `GET /api/v1/metadata/tv/{tmdbId}/season/{sn}` | `GetSeasonMetadata()` |
| `GET /api/v1/metadata/tv/{tmdbId}/season/{sn}/episode/{en}` | `GetEpisodeMetadata()` |
| `GET /api/v1/metadata/collection/{tmdbId}` | `GetCollectionMetadata()` |
| `GET /api/v1/images/proxy/{type}/{size}/{path}` | Image proxy (TMDb CDN) |

### What is NOT Exposed (33 capabilities)

**Movie** (6 hidden): Credits, Images list, Release dates, External IDs, Similar, Recommendations
**TV Show** (5 hidden): Credits, Images list, Content ratings, External IDs, Season/Episode images+credits
**Person** (4 hidden): Search, Metadata, Credits (filmography), Images
**Operations** (3 hidden): RefreshMovie, RefreshTVShow, ClearCache
**TVDb search**: Exists in service layer but not exposed via separate endpoint

---

## Phase 1: Expose Existing Movie/TV Metadata (HIGH priority)

These are all implemented in the service layer. Just need OpenAPI spec + handler code.

### 1A: Movie Metadata Endpoints

| Endpoint | Service Method | Handler |
|----------|---------------|---------|
| `GET /api/v1/metadata/movie/{tmdbId}/credits` | `GetMovieCredits()` | NEW |
| `GET /api/v1/metadata/movie/{tmdbId}/images` | `GetMovieImages()` | NEW |
| `GET /api/v1/metadata/movie/{tmdbId}/similar` | `GetSimilarMovies()` | NEW |
| `GET /api/v1/metadata/movie/{tmdbId}/recommendations` | `GetMovieRecommendations()` | NEW |
| `GET /api/v1/metadata/movie/{tmdbId}/external-ids` | `GetMovieExternalIDs()` | NEW |

### 1B: TV Show Metadata Endpoints

| Endpoint | Service Method | Handler |
|----------|---------------|---------|
| `GET /api/v1/metadata/tv/{tmdbId}/credits` | `GetTVShowCredits()` | NEW |
| `GET /api/v1/metadata/tv/{tmdbId}/images` | `GetTVShowImages()` | NEW |
| `GET /api/v1/metadata/tv/{tmdbId}/content-ratings` | `GetTVShowContentRatings()` | NEW |
| `GET /api/v1/metadata/tv/{tmdbId}/external-ids` | `GetTVShowExternalIDs()` | NEW |
| `GET /api/v1/metadata/tv/{tmdbId}/season/{sn}/images` | Season images | NEW |
| `GET /api/v1/metadata/tv/{tmdbId}/season/{sn}/episode/{en}/images` | Episode images | NEW |

### 1C: Person Metadata Endpoints

| Endpoint | Service Method | Handler |
|----------|---------------|---------|
| `GET /api/v1/metadata/search/person` | `SearchPerson()` | NEW |
| `GET /api/v1/metadata/person/{tmdbId}` | `GetPersonMetadata()` | NEW |
| `GET /api/v1/metadata/person/{tmdbId}/credits` | `GetPersonCredits()` | NEW |
| `GET /api/v1/metadata/person/{tmdbId}/images` | `GetPersonImages()` | NEW |

**Files to modify**:
- `api/openapi/openapi.yaml` — add ~16 endpoint definitions + response schemas
- `internal/api/handler_metadata.go` — add ~16 handler methods
- Regenerate ogen

**Estimated scope**: ~16 new endpoints, all thin wrappers around existing service methods.

---

## Phase 2: Multi-Provider Search (HIGH priority)

### 2A: Add `provider` Query Parameter to Existing Search

Instead of separate TVDb endpoints, add an optional `provider` query param:

```
GET /api/v1/metadata/search/movie?q=inception&provider=tmdb    (default)
GET /api/v1/metadata/search/tv?q=breaking+bad&provider=tmdb    (default)
GET /api/v1/metadata/search/tv?q=breaking+bad&provider=tvdb    (TVDb search)
GET /api/v1/metadata/search/person?q=brad+pitt&provider=tmdb   (default)
```

The handler routes to the correct provider based on param. This keeps the API clean while supporting multiple providers.

### 2B: Provider List Endpoint

```
GET /api/v1/metadata/providers
```

Returns list of available providers with capabilities:
```json
{
  "providers": [
    {"id": "tmdb", "name": "The Movie Database", "movies": true, "tvshows": true, "people": true},
    {"id": "tvdb", "name": "TheTVDB", "movies": false, "tvshows": true, "people": true}
  ]
}
```

**Files to modify**:
- `api/openapi/openapi.yaml` — add `provider` param to search endpoints, add provider list endpoint
- `internal/api/handler_metadata.go` — route search by provider
- `internal/service/metadata/service.go` — add `SearchMovieByProvider()` or refactor `SearchMovie()` to accept provider param

---

## Phase 3: Fanart.tv Provider (MEDIUM priority)

Fanart.tv provides high-quality artwork not available from TMDb:
- HD clearlogos, HD clearart, movie disc art
- Character art, TV thumbs, season posters
- Banner art, background art in ultra-high resolution

### 3A: Implement Provider

New files:
- `internal/service/metadata/providers/fanarttv/provider.go`
- `internal/service/metadata/providers/fanarttv/client.go`
- `internal/service/metadata/providers/fanarttv/types.go`

Implements `ImageProvider` interface:
- `GetImageURL(imageType, size, path) string`
- `DownloadImage(ctx, url) ([]byte, error)`
- Maps Fanart.tv image types to generic image categories

### 3B: Config

```go
type FanartTVConfig struct {
    Enabled bool   `koanf:"enabled"`
    APIKey  string `koanf:"api_key"`
}
```

### 3C: Wire Into Image Endpoints

After Fanart.tv is registered as a provider, the image listing endpoints from Phase 1 automatically include Fanart.tv images alongside TMDb images (provider chain with fallback).

**Files**:
- `internal/service/metadata/providers/fanarttv/` (new package)
- `internal/config/config.go` — add FanartTVConfig
- `internal/service/metadata/metadatafx/module.go` — register provider

---

## Phase 4: Frontend Polish (LOW priority)

### 4A: SSE Real-Time Events (M1)
- `GET /api/v1/events/stream` with `text/event-stream`
- Events: library scan progress, new content, job completion
- Auth via query param token (SSE doesn't support headers)

### 4B: Cookie Auth + CSRF (L1/L2)
- Deferred to SvelteKit SSR phase
- Token-in-header is fine for SPA MVP

---

## Merged TODO List

### Already Done (this session)
- [x] C1: CORS middleware + config
- [x] H1: List endpoint pagination (`total` field)
- [x] H3: 201 response bodies (addPolicy, assignRole)
- [x] M2: OpenAPI spec + Scalar API docs
- [x] M3: Webhook security (X-Webhook-Secret)
- [x] H1-ext: MovieListResponse + TVShowListResponse wrappers
- [x] Fix 10 CI lint issues

### Phase 1: Expose Existing Capabilities - DONE
- [x] 1A: Movie metadata endpoints (credits, images, similar, recommendations, external-ids)
- [x] 1B: TV show metadata endpoints (credits, images, content-ratings, external-ids, season/episode images)
- [x] 1C: Person metadata endpoints (search, details, credits, images)
- [x] OpenAPI spec additions + ogen regeneration
- [x] 16 handler implementations

### Phase 2: Multi-Provider Search - DONE
- [x] 2A: Add `provider` + `language` query params to search endpoints (tmdb/tvdb routing)
- [x] 2B: Provider list endpoint (`GET /api/v1/metadata/providers`)
- [x] TVDb search exposed via provider param

### Phase 3: Metadata Providers - DONE
- [x] 3A: Fanart.tv provider (HD artwork for movies + TV shows)
- [x] 3B: OMDb provider (IMDb/RT/Metacritic ratings, ExternalRatings type)
- [x] 3C: TVmaze provider (free TV metadata, episodes, cast/crew)
- [x] Config integration (api_key, enabled toggles)
- [x] Wire all into metadatafx module

### Phase 4: Frontend Polish (deferred)
- [ ] 4A: SSE events endpoint
- [ ] 4B: Cookie auth + CSRF (only for SSR)

---

## Files Impact Summary

| Phase | New Files | Modified Files |
|-------|-----------|----------------|
| 1 | 0 | openapi.yaml, handler_metadata.go, ogen/* |
| 2 | 0 | openapi.yaml, handler_metadata.go, service.go, ogen/* |
| 3 | 3+ | providers/fanarttv/*, config.go, metadatafx/module.go |
| 4 | 1-2 | handler_events.go (new), server.go, openapi.yaml |

## Verification

After Phase 1:
1. `make build` compiles
2. `make test` passes
3. `curl localhost:8096/api/v1/metadata/movie/550/credits` returns Fight Club credits
4. `curl localhost:8096/api/v1/metadata/search/person?q=brad+pitt` returns results
5. Scalar docs at `/api/docs` shows all new endpoints

After Phase 2:
1. `curl localhost:8096/api/v1/metadata/search/tv?q=breaking+bad&provider=tvdb` returns TVDb results
2. `curl localhost:8096/api/v1/metadata/providers` lists tmdb + tvdb

After Phase 3:
1. `curl localhost:8096/api/v1/metadata/movie/550/images` includes fanart.tv artwork
2. Config `integrations.fanarttv.api_key` works
