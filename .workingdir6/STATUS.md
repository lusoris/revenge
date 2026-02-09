# Frontend Readiness Status (2026-02-08)

## Completed

### C1: CORS Middleware (CRITICAL) - DONE
- `internal/api/middleware/cors.go` + tests
- Origin reflection, credentials, preflight 204
- Wired as outermost HTTP wrapper in server.go
- Config: `server.cors.allowed_origins`, `server.cors.allow_credentials`, `server.cors.max_age`

### H1: List Endpoint Pagination - DONE
- Added `total` (int64) to 10 bounded list response schemas (sessions, policies, roles, etc.)
- Added `MovieListResponse` and `TVShowListResponse` wrapper schemas with `items` + `total`
- Movie/TVShow handlers now call `CountMovies()`/`CountSeries()` for accurate totals
- All handlers updated to populate Total field

### H3: 201 Response Bodies - DONE
- `addPolicy` 201 now returns `Policy` object
- `assignRole` 201 now returns `RoleListResponse`

### M2: OpenAPI Spec + Scalar API Docs - DONE
- `GET /api/openapi.yaml` serves embedded spec
- `GET /api/docs` serves Scalar API reference UI
- `api/openapi/embed.go` embeds the spec via go:embed

### M3: Webhook Security - DONE
- Added `webhookAuth` security scheme (X-Webhook-Secret header) to OpenAPI
- Applied to both `/api/v1/webhooks/radarr` and `/api/v1/webhooks/sonarr`
- Config: `integrations.radarr.webhook_secret`, `integrations.sonarr.webhook_secret`
- Backwards-compatible: if no secret configured, allows unauthenticated webhooks

### Phase 1: Expose Existing Metadata Capabilities - DONE
- 16 new metadata API endpoints (movie credits/images/similar/recommendations/external-ids, TV credits/images/content-ratings/external-ids/season+episode images, person search/details/credits/images)
- OpenAPI spec additions + ogen regeneration
- Handler implementations in `handler_metadata.go`
- Commits: `c86586a5`, `78084fed`

### Phase 2: Multi-Provider Search - DONE
- Added `provider` (enum: tmdb, tvdb) and `language` (ISO 639-1) query params to all 3 search endpoints
- Service layer routes to specific provider when `ProviderID` is set
- TVDb search now accessible via `?provider=tvdb`
- Commit: `555fa7f6`

### Phase 3: Metadata Providers - DONE
- **Fanart.tv** (`b83a1313`): HD artwork provider (logos, clearart, disc art, banners, character art) for movies + TV shows. Priority 60. Config: `metadata.fanarttv.api_key`, `metadata.fanarttv.client_key`
- **OMDb** (`141a9b87`): External ratings provider (IMDb, Rotten Tomatoes, Metacritic). Added `ExternalRatings` field to MovieMetadata/TVShowMetadata. Priority 40. Config: `metadata.omdb.api_key`
- **TVmaze** (`4454051a`): Free TV metadata provider (episodes, cast/crew, seasons, images, external ID cross-refs). No API key needed. Priority 50. Config: `metadata.tvmaze.enabled`

### Observability - DONE
- Wired all 31 Prometheus metrics (auth, rate limit, jobs, search, library scan, sessions, cache, DB)
- Job completion metrics via River event subscription
- Periodic stats collector goroutine
- pprof always-on (removed dev-mode gate)
- Makefile bench/pprof targets
- Benchmark tests for cache, search, rate limiter, metrics
- Commit: `12240c41`

### Flaky Test Fixes - DONE
- `TestCache_Eviction`: Replaced fixed 100ms sleep with `require.Eventually` polling for otter async eviction
- `TestStreamHandler_ServeMediaPlaylist_Cached`: Added 50ms sleep for otter async Set before asserting cache hit

## Provider Priority Chain

| Provider | Priority | Capabilities | Auth |
|----------|----------|-------------|------|
| TMDb | 100 | Movies, TV, People, Images, Collections | API key |
| TVDb | 80 | TV shows, Images | API key |
| Fanart.tv | 60 | HD artwork (movies + TV) | API key |
| TVmaze | 50 | TV shows, Episodes, Cast | Free |
| OMDb | 40 | Ratings (IMDb/RT/Metacritic) | API key |

## Not Done / Deferred

### H2: Standardize Search Pagination
- Search uses `page`/`per_page` (maps to Typesense), lists use `limit`/`offset`
- Decided: keep as-is. Both patterns are standard.

### M1: SSE Real-Time Events
- Deferred. Not blocking frontend MVP.

### L1/L2: Cookie Auth + CSRF
- Deferred to SSR phase. Token-in-header is fine for SPA.

### Provider List Endpoint
- `GET /api/v1/metadata/providers` — returns available providers with capabilities
- Not yet implemented (P2B from plan)

## Commits (this session chain)
- `ccf9b96c` feat: add CORS middleware, Scalar API docs, and pagination totals
- `60518cfe` feat: add webhook security, movie/tvshow pagination wrappers, fix lint
- `c86586a5` feat(api): expose 16 metadata endpoints
- `78084fed` feat(api): implement all 16 metadata endpoint handlers
- `555fa7f6` feat(api): add multi-provider search with provider+language params
- `12240c41` feat(observability): wire job completion metrics, add missing benchmarks
- `b83a1313` feat: implement Fanart.tv provider, fix flaky otter cache tests
- `141a9b87` feat: implement OMDb provider with external ratings support
- `4454051a` feat: implement TVmaze provider for free TV metadata
