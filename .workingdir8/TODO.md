# Backend Polish & Frontend Prep (2026-02-09)

Full analysis of stubs, SvelteKit readiness, River jobs, caching, Typesense, logger,
deduplication, and API structure.

---

## 1. Stubs & Unfinished Code

### 1A — Production Stubs (Critical)

| # | File | Issue | Fix |
|---|------|-------|-----|
| 1 | `internal/service/auth/mfa_integration.go:112` | WebAuthn `VerifyMFA` returns `"webauthn verification not yet implemented"` | Implement WebAuthn login assertion verification or gate the UI so users can't attempt WebAuthn MFA login |
| 2 | `internal/infra/database/db/apikeys.sql.go:182` | `GetAPIKeyUsageCount` is a placeholder — returns `last_used_at` instead of counting | Create `api_key_usage` tracking table or accept the placeholder and remove misleading name |

### 1B — Orphaned Job Definitions (High)

| # | File | Issue |
|---|------|-------|
| 3 | `internal/service/metadata/jobs/refresh.go` | 5 Args structs (`metadata_refresh_tvshow`, `metadata_refresh_season`, `metadata_refresh_episode`, `metadata_refresh_person`, `metadata_download_image`) have `Kind()` but **no registered Worker** — jobs are inserted but silently lost |
| 4 | `internal/service/metadata/jobs/refresh.go:86` | `metadata_enrich_content` (`EnrichContentArgs`) — same: no worker registered |

### 1C — TODO Comments (Medium)

| # | File:Line | Note |
|---|-----------|------|
| 5 | `providers/trakt/mapping.go:113` | `// TODO: store in a dedicated Certification field when available` (duplicated at L229) |
| 6 | `providers/simkl/mapping.go:58` | Same certification TODO |
| 7 | `service/session/service_test.go:294` | Stale comment `"Count not implemented yet"` — count IS now implemented, delete comment |

### 1D — Placeholder Modules (Low)

| # | Module | Status |
|---|--------|--------|
| 8 | `internal/content/qar/` | Entire module is `SELECT 1 AS placeholder` — v0.3.0+ scope |
| 9 | `internal/content/movie/db/placeholder.sql.go` | Vestigial placeholder alongside 50+ real queries — delete |
| 10 | `internal/content/tvshow/db/placeholder.sql.go` | Same — delete |

### 1E — Dead Dependencies

| # | Dependency | Issue |
|---|-----------|-------|
| 11 | `gobreaker` | In `go.mod` (indirect) but removed from code with comment in `search/module.go:76` |

---

## 2. SvelteKit + Tailwind Backend Prep

### 2A — Already Ready

| Capability | Status | Details |
|-----------|--------|---------|
| CORS for SPA | **READY** | Full impl in `middleware/cors.go` — origin reflection, credentials, Vary, 12h preflight |
| Cookie sessions for SSR | **READY** | `revenge_access_token` + `revenge_refresh_token` HttpOnly cookies, `CookieAuthMiddleware` auto-injects Bearer |
| CSRF protection | **READY** | Double-submit cookie `revenge_csrf` + `X-CSRF-Token` header, only on state-changing methods |
| SSE live updates | **READY** | `GET /api/v1/events` with category filtering, keepalive, auth |
| OpenAPI → TS client | **READY** | Spec at `/api/openapi.yaml` with caching — use `openapi-typescript` or `orval` |
| Image proxy/resize | **READY** | Multiple sizes (w185–original), ETag, `Cache-Control: immutable` |
| SSR auth validation | **READY** | Cookie→Bearer middleware, SvelteKit server hooks can forward cookies |

### 2B — Needs Work

| # | Gap | Severity | Fix |
|---|-----|----------|-----|
| 1 | **No Content Security Policy** headers | Medium | Add CSP middleware: `default-src 'self'; img-src 'self' image.tmdb.org; media-src 'self' blob:; script-src 'self'` |
| 2 | **No Cache-Control on JSON API responses** | Medium | Add `private, no-cache` for user-specific, `public, max-age=60` for library listings |
| 3 | **Query param naming inconsistency** | Medium | Movies use `orderBy` (camelCase), TV shows use `order_by` (snake_case) — standardize to snake_case |
| 4 | **Some list endpoints return bare arrays** | Low | `getRecentlyAdded`, `getTopRated`, cast/crew return arrays without `{items, total}` wrapping |
| 5 | **No admin user management endpoint** | Medium | Need `GET /api/v1/admin/users` (list), `DELETE /api/v1/users/{id}` |
| 6 | **No global genre endpoint** | Low | `GET /api/v1/genres` for filter UIs (currently only via search facets) |
| 7 | **No bulk episode watched endpoint** | Low | `POST /api/v1/shows/{id}/seasons/{num}/watched` to mark entire season |
| 8 | **TV show reindex not exposed** | Low | No `POST /api/v1/search/tvshows/reindex` endpoint (only movies) |

---

## 3. River Jobs — Gaps & Fixes

### 3A — Current Queue Layout (Good)

| Queue | Workers | Jobs Using It |
|-------|---------|---------------|
| `critical` | 20 | **NONE** — 20 workers idle |
| `high` | 15 | radarr_sync, radarr_webhook, sonarr_sync, sonarr_webhook, notification |
| `default` | 10 | movie_file_match, movie_search_index, metadata_refresh_movie, tvshow_file_match, tvshow_metadata_refresh, tvshow_series_refresh |
| `low` | 5 | cleanup, library_scan_cleanup, activity_cleanup, playback_cleanup |
| `bulk` | 3 | movie_library_scan, tvshow_library_scan, tvshow_search_index |

### 3B — Critical Issues

| # | Issue | Severity | Fix |
|---|-------|----------|-----|
| 1 | **`critical` queue has 0 jobs** — 20 workers allocated for nothing | High | Either use it for security/auth events or remove it and redistribute workers |
| 2 | **5 metadata job kinds have no Workers registered** (`metadata_refresh_tvshow`, `_season`, `_episode`, `_person`, `_download_image`) | **Critical** | Register workers or remove dead code. The `EnqueueRefresh*` and `EnqueueDownloadImage` functions in `metadata/jobs/queue.go` silently insert jobs that are never processed |
| 3 | **`metadata_enrich_content` has no Worker** | High | Create `EnrichContentWorker` or remove the Args definition |
| 4 | **Periodic cleanup never actually runs** — `Schedule*` functions take `*river.Client[any]` but app uses `*river.Client[pgx.Tx]` (type mismatch), and they're never called from DI lifecycle | **Critical** | Fix the type signatures and wire `Schedule*` into the fx `OnStart` hook |
| 5 | **Movie metadata refresh is synchronous in API handler** | High | `handler.go:201` calls service method directly — should `Insert` a `RefreshMovieArgs` job and return `202 Accepted` |
| 6 | **Activity logging is synchronous** — `service/activity/service.go:33` blocks on DB insert | Medium | Create `ActivityLogArgs` job, process on `low` queue |
| 7 | **Bare goroutine fallback** in `handler_radarr.go:145` / `handler_sonarr.go:155` when `riverClient == nil` | Low | Remove fallback or add basic error handling |

### 3C — Missing Jobs (Nice to Have)

| # | Operation | Current | Proposed |
|---|-----------|---------|----------|
| 8 | Image cache warming | Synchronous on miss | `DownloadImageWorker` pre-downloads in metadata refresh pipeline |
| 9 | Bulk search reindex | Only movies exposed via API | Add `tvshow_search_reindex` endpoint |
| 10 | Stats/analytics aggregation | Not implemented | Periodic job to compute library stats (counts, durations, etc.) |

---

## 4. Caching — Analysis & Gaps

### 4A — Architecture (Good)

- **L1**: Otter (W-TinyLFU) — 10k entries, 5 min default TTL
- **L2**: Rueidis (Dragonfly/Redis) — client-side cache, 16 MiB per conn
- **Pattern**: L1 → L2 → miss, write-through, invalidation via key prefix
- **7 CachedService wrappers**: Movie, User, Session, Settings, RBAC, Library, Search

### 4B — Gaps

| # | Gap | Severity | Fix |
|---|-----|----------|-----|
| 1 | **No CachedService for TV shows** — all TV show reads hit DB every time | **High** | Create `CachedTVShowService` mirroring the movie pattern |
| 2 | **No CachedTVShowSearchService** — TV show Typesense searches always hit Typesense | Medium | Create cached wrapper like `CachedMovieSearchService` |
| 3 | **L1 fully cleared on any pattern invalidation** — Otter doesn't support prefix delete | Medium | Switch to Otter's `DeleteByFunc` or use a separate cache instance per domain |
| 4 | **Redis `KEYS` command for pattern invalidation** — O(N), blocks | Low | Use `SCAN` instead, or restructure keys to avoid pattern matching |
| 5 | **Metadata provider caches are L1-only** — lost on restart, not shared across instances | Low | For now acceptable (API rate limits make full provider caching undesirable) |
| 6 | **No API key validation cache** — every API key request hits DB | Low | Add 30s cache like session validation |
| 7 | **Movie relational data (cast/crew/genres) not cached individually** — key prefixes exist but `CachedService` only wraps top-level movie | Low | Either cache in `CachedService` or accept DB reads for detail views |

---

## 5. Typesense Search — Analysis & Gaps

### 5A — Current State (Good)

| Collection | Fields | Searchable | Facetable |
|-----------|--------|-----------|-----------|
| `movies` | 31 fields | title, original_title (infix), overview, cast, directors | genres, year, status, directors, language, has_file, resolution, quality_profile |
| `tvshows` | 29 fields | title, original_title (infix), overview, cast, networks | genres, year, status, type, language, has_file, networks |

- Both synced via River background jobs (`movie_search_index`, `tvshow_search_index`)
- Movie search results cached 30s (via `CachedMovieSearchService`)
- All authenticated users can search (no RBAC differentiation)

### 5B — Gaps

| # | Gap | Severity | Fix |
|---|-----|----------|-----|
| 1 | **Episodes not indexed** — can't search by episode title ("The Rains of Castamere") | Medium | Create `episodes` collection with series context, or embed episodes in TV show docs |
| 2 | **People not a standalone collection** — cast embedded as string arrays, no cross-content person search | Medium | Create `people` collection (name, bio, photo, known_for, linked movies/shows) |
| 3 | **No unified search** — must hit `/search/movies` and `/search/tvshows` separately | Medium | Add `/search/multi` endpoint that queries both collections and merges results |
| 4 | **Settings not searchable** — but should they be? Settings are admin-only key/value pairs | Low | For an admin UI, a simple frontend filter is sufficient — Typesense overkill |
| 5 | **Users not searchable** — admin user management would benefit from user search | Low | Only needed if admin panel requires user search; do client-side in paginated list for now |
| 6 | **TV show reindex endpoint missing** | Low | Expose `POST /api/v1/search/tvshows/reindex` (admin-only) |
| 7 | **TV show search not cached** | Low | Create `CachedTVShowSearchService` (30s TTL like movies) |

---

## 6. Logger — Verdict: Clean

- **Unified on `log/slog`** (Go stdlib structured logging) everywhere
- Sub-loggers with component context (`slog.With("service", "...")`)
- No bare `fmt.Println` in production code (only in CLI `cmd/` — acceptable)
- No action needed

---

## 7. Deduplication Opportunities

| # | Pattern | Occurrences | Severity | Fix |
|---|---------|-------------|----------|-----|
| 1 | **Admin auth guard** — `getUserID` + `HasRole("admin")` + early return | ~20+ handlers | **High** | Extract `requireAdmin(ctx) (uuid.UUID, error)` helper in `handler.go` |
| 2 | **ExternalRatings→Ogen conversion** — identical `[]ogen.ExternalRating` mapping | 4 copies (movie_converters, tvshow_converters) | Medium | Extract `externalRatingsToOgen([]domain.ExternalRating) []ogen.ExternalRating` |
| 3 | **Movie field copy in continue/watched converters** — ~100 lines of identical optional field mapping | 2 functions in `movie_converters.go` | Medium | Extract `copyMovieFieldsToOgen(movie, ogenStruct)` or generate with ogen hooks |
| 4 | **Radarr/Sonarr webhook handler scaffolding** | 2 packages | Low | Acceptable — domains differ enough that a shared base isn't worth abstracting |
| 5 | **Optional field `nil → SetTo` pattern** | ~200 lines across converters | Low | Unavoidable with ogen types, no fix without codegen |

---

## Priority Order

### Tier 1 — Must Fix (Broken / Dead Code)

- [ ] **3B.2**: Register workers for 5 metadata job kinds or remove dead Args + `Enqueue*` functions
- [ ] **3B.4**: Fix periodic cleanup type mismatch and wire into fx `OnStart`
- [ ] **4B.1**: Create `CachedTVShowService`
- [ ] **1A.1**: Fix or gate WebAuthn `VerifyMFA` stub

### Tier 2 — High Value

- [ ] **3B.1**: Use `critical` queue or remove it
- [ ] **3B.5**: Make movie metadata refresh async (return 202)
- [ ] **7.1**: Extract `requireAdmin()` helper (dedup ~20 handlers)
- [ ] **2B.3**: Standardize query param naming to snake_case
- [ ] **2B.1**: Add CSP middleware
- [ ] **7.2**: Extract `externalRatingsToOgen()` helper

### Tier 3 — Medium Value

- [ ] **3B.6**: Make activity logging async
- [ ] **4B.2**: Create `CachedTVShowSearchService`
- [ ] **5B.3**: Add unified `/search/multi` endpoint
- [ ] **2B.2**: Add Cache-Control headers on JSON API responses
- [ ] **2B.5**: Add admin user management endpoints
- [ ] **5B.1**: Index episodes in Typesense
- [ ] **5B.2**: Create standalone people collection in Typesense
- [ ] **1A.2**: Fix or rename `GetAPIKeyUsageCount` placeholder

### Tier 4 — Polish

- [ ] **2B.4**: Wrap bare array list responses in `{items, total}`
- [ ] **2B.6**: Add `GET /api/v1/genres` endpoint
- [ ] **2B.7**: Add bulk episode watched endpoint
- [ ] **3B.3**: Create `EnrichContentWorker`
- [ ] **3C.8**: Create `DownloadImageWorker` for cache warming
- [ ] **4B.3**: Fix L1 pattern invalidation (avoid full cache clear)
- [ ] **4B.4**: Replace Redis `KEYS` with `SCAN`
- [ ] **1C.7**: Delete stale test comment
- [ ] **1D.9-10**: Delete vestigial sqlc placeholders
- [ ] **1E.11**: Remove gobreaker from go.mod

---

## Notes

- **Logger**: Already clean and unified on `slog` — no work needed
- **API structure**: Solid overall. OpenAPI 3.1.0 spec, ogen codegen, fx DI, proper error handling, SSE events
- **Caching**: Excellent L1/L2 architecture, just missing TV show coverage
- **River jobs**: Good queue design, but several orphaned job definitions and unscheduled periodic tasks
- **Search**: Strong movie/tvshow coverage, needs episodes + people + unified endpoint
- **Frontend readiness**: 7/10 — auth, CORS, cookies, SSE all ready; needs CSP, cache headers, API consistency
