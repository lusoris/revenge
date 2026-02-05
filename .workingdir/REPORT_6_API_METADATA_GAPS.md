# API Metadata Endpoint Gaps Report (2026-02-06)

**Generated during Phase A12 API review**

---

## Summary

The shared metadata service has full support for both Movie and TV metadata providers (TMDb, TVDb). **Core TV metadata endpoints are now complete** (search, show details, season details, episode details). Extended endpoints (credits, images, content ratings, external IDs) remain as low priority future work.

---

## Current State

### Movie Metadata Endpoints (COMPLETE)

| Endpoint | Handler | Status |
|----------|---------|--------|
| `GET /api/v1/metadata/search/movie` | `SearchMoviesMetadata` | ✅ |
| `GET /api/v1/metadata/movie/{tmdbId}` | `GetMovieMetadata` | ✅ |
| `GET /api/v1/metadata/collection/{tmdbId}` | `GetCollectionMetadata` | ✅ |
| `GET /api/v1/images/proxy/{type}/{size}/{path}` | `GetProxiedImage` | ✅ |

### TV Metadata Endpoints (CORE COMPLETE)

| Service Method | Endpoint | Status |
|---------------|----------|--------|
| `SearchTVShow` | `GET /api/v1/metadata/search/tv` | ✅ |
| `GetTVShowMetadata` | `GET /api/v1/metadata/tv/{tmdbId}` | ✅ |
| `GetSeasonMetadata` | `GET /api/v1/metadata/tv/{tmdbId}/season/{seasonNumber}` | ✅ |
| `GetEpisodeMetadata` | `GET /api/v1/metadata/tv/{tmdbId}/season/{seasonNumber}/episode/{episodeNumber}` | ✅ |
| `GetTVShowCredits` | `GET /api/v1/metadata/tv/{tmdbId}/credits` | LOW |
| `GetTVShowImages` | `GET /api/v1/metadata/tv/{tmdbId}/images` | LOW |
| `GetTVShowContentRatings` | `GET /api/v1/metadata/tv/{tmdbId}/content-ratings` | LOW |
| `GetTVShowExternalIDs` | `GET /api/v1/metadata/tv/{tmdbId}/external-ids` | LOW |

### Person Metadata Endpoints (MISSING)

The TMDb provider has person lookup, but no API endpoints:

| Suggested Endpoint | Priority |
|-------------------|----------|
| `GET /api/v1/metadata/search/person` | LOW |
| `GET /api/v1/metadata/person/{tmdbId}` | LOW |

---

## Radarr/Sonarr Integration API Status

### Radarr Endpoints (COMPLETE)

| Endpoint | Handler | Status |
|----------|---------|--------|
| `GET /api/v1/admin/integrations/radarr/status` | `AdminGetRadarrStatus` | ✅ |
| `POST /api/v1/admin/integrations/radarr/sync` | `AdminTriggerRadarrSync` | ✅ |
| `GET /api/v1/admin/integrations/radarr/quality-profiles` | `AdminGetRadarrQualityProfiles` | ✅ |
| `GET /api/v1/admin/integrations/radarr/root-folders` | `AdminGetRadarrRootFolders` | ✅ |
| `POST /api/v1/webhooks/radarr` | `HandleRadarrWebhook` | ✅ |

### Sonarr Endpoints (COMPLETE)

| Endpoint | Handler | Status |
|----------|---------|--------|
| `GET /api/v1/admin/integrations/sonarr/status` | `AdminGetSonarrStatus` | ✅ |
| `POST /api/v1/admin/integrations/sonarr/sync` | `AdminTriggerSonarrSync` | ✅ |
| `GET /api/v1/admin/integrations/sonarr/quality-profiles` | `AdminGetSonarrQualityProfiles` | ✅ |
| `GET /api/v1/admin/integrations/sonarr/root-folders` | `AdminGetSonarrRootFolders` | ✅ |
| `POST /api/v1/webhooks/sonarr` | `HandleSonarrWebhook` | ✅ |

---

## Backend Service Readiness

### TMDb Provider (`internal/service/metadata/providers/tmdb/`)

- ✅ Movie search
- ✅ Movie details
- ✅ Movie credits
- ✅ Movie images
- ✅ Collection details
- ✅ **TV show search**
- ✅ **TV show details**
- ✅ **TV show credits**
- ✅ **TV show images**
- ✅ **TV show content ratings**
- ✅ **TV show translations**
- ✅ **TV show external IDs**
- ✅ **Season details**
- ✅ **Season credits**
- ✅ **Season images**
- ✅ **Episode details**
- ✅ **Episode credits**
- ✅ **Episode images**

### TVDb Provider (`internal/service/metadata/providers/tvdb/`)

- ✅ TV show search
- ✅ TV show details
- ✅ TV show credits
- ✅ TV show images
- ✅ TV show content ratings
- ✅ TV show translations
- ✅ TV show external IDs
- ✅ Season details
- ✅ Episode details

---

## Action Items

### ~~Priority 1 - HIGH (TV Metadata Search/Details)~~ COMPLETED

1. [x] Add `GET /api/v1/metadata/search/tv` to OpenAPI spec
2. [x] Add `GET /api/v1/metadata/tv/{tmdbId}` to OpenAPI spec
3. [x] Implement `SearchTVShowsMetadata` handler
4. [x] Implement `GetTVShowMetadata` handler
5. [x] Regenerate ogen code

### ~~Priority 2 - MEDIUM (Season/Episode Details)~~ COMPLETED

6. [x] Add `GET /api/v1/metadata/tv/{tmdbId}/season/{seasonNumber}` to OpenAPI spec
7. [x] Add `GET /api/v1/metadata/tv/{tmdbId}/season/{seasonNumber}/episode/{episodeNumber}` to OpenAPI spec
8. [x] Implement `GetSeasonMetadata` handler
9. [x] Implement `GetEpisodeMetadata` handler

### Priority 3 - LOW (Extended TV Metadata)

10. [ ] Add TV credits, images, content ratings, external IDs endpoints
11. [ ] Add person search/details endpoints

---

## Files to Modify

| File | Action |
|------|--------|
| `api/openapi/openapi.yaml` | Add TV metadata endpoints and schemas |
| `internal/api/handler_metadata.go` | Add TV metadata handler methods |
| `internal/api/ogen/*` | Regenerate with `make ogen` |

---

## Notes

- The backend metadata service is fully implemented for both Movie and TV content
- Only the API layer needs updating to expose TV metadata endpoints
- TMDb is the primary provider for both Movie and TV metadata
- TVDb provides additional TV-specific data and can be used as fallback
- Radarr/Sonarr integration APIs are complete and at feature parity
