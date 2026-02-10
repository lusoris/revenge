# Backend Polish & Frontend Prep (2026-02-09)

Full analysis of stubs, SvelteKit readiness, River jobs, caching, Typesense, logger,
deduplication, dependency audit, schema architecture, and API structure.

---

## 1. Stubs & Unfinished Code

### 1A — Production Stubs (Critical)

| # | File | Issue | Fix |
|---|------|-------|-----|
| 1 | `internal/service/auth/mfa_integration.go:112` | WebAuthn `VerifyMFA` returns `"webauthn verification not yet implemented"` | **Implement** WebAuthn login assertion verification (we want it testable) |
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
| 3 | **Query param naming inconsistency** | **High** | Movies use `orderBy` (camelCase), TV shows use `order_by` (snake_case) — **standardize to snake_case** (REST convention, matches Go/PostgreSQL, TS client can transform) |
| 4 | ~~**All list endpoints need pagination**~~ | ✅ DONE | ~~`getRecentlyAdded`, `getTopRated`, cast/crew return bare arrays~~ — `4d8f8919`: 7 endpoints now return `{items, total, page, page_size}` with proper caching |
| 5 | **No admin user management endpoint** | **High** | Need `GET /api/v1/admin/users` (list+search), `DELETE /api/v1/users/{id}` — backend must be ready before frontend |
| 6 | **No global genre endpoint** | Medium | `GET /api/v1/genres` for filter UIs (currently only via search facets) |
| 7 | **No bulk episode watched endpoint** | Medium | `POST /api/v1/shows/{id}/seasons/{num}/watched` to mark entire season |
| 8 | **TV show reindex not exposed** | Medium | No `POST /api/v1/search/tvshows/reindex` endpoint (only movies) |

---

## 3. River Jobs — Gaps & Fixes

### 3A — Current Queue Layout

| Queue | Workers | Jobs Using It |
|-------|---------|---------------|
| `critical` | 20 | Auth/login, frontend SSR, streaming/playback, SSE events (TO BE WIRED) |
| `high` | 15 | radarr_sync, radarr_webhook, sonarr_sync, sonarr_webhook, notification |
| `default` | 10 | movie_file_match, movie_search_index, metadata_refresh_movie, tvshow_file_match, tvshow_metadata_refresh, tvshow_series_refresh |
| `low` | 5 | cleanup, library_scan_cleanup, activity_cleanup, playback_cleanup |
| `bulk` | 3 | movie_library_scan, tvshow_library_scan, tvshow_search_index |

### 3B — Critical Issues

| # | Issue | Severity | Fix |
|---|-------|----------|-----|
| 1 | **`critical` queue has 0 jobs wired** — 20 workers idle | **High** | Wire auth/login events, frontend-critical tasks, streaming/playback session start, SSE delivery to `critical` queue |
| 2 | **5 metadata job kinds have no Workers registered** (`metadata_refresh_tvshow`, `_season`, `_episode`, `_person`, `_download_image`) | **Critical** | Register workers or remove dead code. The `EnqueueRefresh*` and `EnqueueDownloadImage` functions in `metadata/jobs/queue.go` silently insert jobs that are never processed |
| 3 | **`metadata_enrich_content` has no Worker** | High | Create `EnrichContentWorker` — enrichment must run as River job, NOT synchronous in request path |
| 4 | **Periodic cleanup never actually runs** — `Schedule*` functions take `*river.Client[any]` but app uses `*river.Client[pgx.Tx]` (type mismatch), and they're never called from DI lifecycle | **Critical** | Fix the type signatures and wire `Schedule*` into the fx `OnStart` hook |
| 5 | **Movie metadata refresh is synchronous in API handler** | **High** | `handler.go:201` calls service method directly — must `Insert` a `RefreshMovieArgs` job and return `202 Accepted` |
| 6 | ~~**Rating enrichment is synchronous in request path**~~ | **High** | ~~`enrichMovieRatings()` / `enrichTVShowRatings()` call external provider APIs synchronously — must be River jobs. User gets cached/partial data immediately, enrichment runs async~~ ✅ `d2281594` — moved to adapter Enrich methods (River workers only) |
| 7 | **Activity logging is synchronous** — `service/activity/service.go:33` blocks on DB insert | Medium | Create `ActivityLogArgs` job, process on `low` queue |
| 8 | **Bare goroutine fallback** in `handler_radarr.go:145` / `handler_sonarr.go:155` when `riverClient == nil` | Low | Remove fallback or add basic error handling |

### 3C — Missing Jobs (Needed)

| # | Operation | Current | Proposed |
|---|-----------|---------|----------|
| 9 | Image cache warming | Synchronous on miss | `DownloadImageWorker` pre-downloads in metadata refresh pipeline |
| 10 | Bulk search reindex | Only movies exposed via API | Add `tvshow_search_reindex` endpoint + job |
| 11 | Stats/analytics aggregation | Not implemented | Periodic job to compute library stats (counts, durations, etc.) |

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
| 1 | ~~**No CachedService for TV shows**~~ — all TV show reads hit DB every time | **Critical** | ~~Create `CachedTVShowService` mirroring the movie pattern~~ ✅ `962f8192` — created + wired in DI (movie too) |
| 2 | **No CachedTVShowSearchService** — TV show Typesense searches always hit Typesense | **High** | Create cached wrapper like `CachedMovieSearchService` |
| 3 | **L1 fully cleared on any pattern invalidation** — Otter doesn't support prefix delete | Medium | Verify Otter docs for `DeleteByFunc`, or use a separate cache instance per domain |
| 4 | **Redis `KEYS` command for pattern invalidation** — O(N), blocks | Medium | Verify current impl — use `SCAN` instead, or restructure keys to avoid pattern matching |
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
| 1 | **Episodes not indexed** — can't search by episode title | **High** | Create `episodes` collection with episode title, season name, series context |
| 2 | **Seasons not indexed** — many shows have named seasons (e.g. "American Horror Story: Coven") | **High** | Create `seasons` collection with season name, number, series context |
| 3 | **People not a standalone collection** — cast embedded as string arrays | **High** | Create `people` collection (name, bio, photo, known_for, linked movies/shows) — people get their own pages in frontend |
| 4 | **No unified search** — must hit `/search/movies` and `/search/tvshows` separately | **High** | Add `/search/multi` endpoint that queries all collections and merges results |
| 5 | **Users not searchable** — admin user management needs user search | Medium | Create `users` collection for admin panel search |
| 6 | **TV show reindex endpoint missing** | Medium | Expose `POST /api/v1/search/tvshows/reindex` (admin-only) |
| 7 | **TV show search not cached** | **High** | Create `CachedTVShowSearchService` (30s TTL like movies) |

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
| 4 | **Pointer helpers** — `stringPtr`, `ptrToString`, `deref*` duplicated across 10+ files | 10+ copies | **High** | Create `internal/util/ptr/` with generic `To[T](*T)` and `From[T](*T) T` |
| 5 | **Provider ErrNotFound stubs** — ~150 identical methods across 10 providers | ~150 methods | **High** | Create `BaseMovieProvider`/`BaseTVShowProvider` embedding structs with default `ErrNotFound` returns |
| 6 | **Test helpers** — `createTestUser` in 5+ files, `setupTestService` in 7+ files | 12+ copies | Medium | Move to `internal/testutil/helpers.go` |

---

## 8. Monitoring & Observability

### 8A — Fixes Applied

| # | Issue | Fix |
|---|-------|-----|
| 1 | Datasource UID mismatch → all Grafana panels NO DATA | Added `uid: PBFA97CFB590B2093` to `deploy/grafana/provisioning/datasources/prometheus.yml` |
| 2 | Metrics port 9096 not exposed | Added `EXPOSE 9096` to Dockerfile, added `9096:9096` to `docker-compose.dev.yml` |

### 8B — Outstanding

| # | Issue | Severity | Fix |
|---|-------|----------|-----|
| 3 | `database/metrics.go` — dead code (`RecordPoolMetrics` has 0 callers, 0% coverage) | Low | Delete file + test — `observability/collector.go` already handles pgxpool metrics under `revenge_pgxpool_*` namespace |
| 4 | Metrics wiring is correct (47+ metrics registered, all incremented) but untested end-to-end | Low | After UID fix, smoke test: `docker compose -f docker-compose.dev.yml up`, then check `http://localhost:9096/metrics` and Grafana panels |

---

## 9. Schema Architecture — Movie Module Needs Own Schema

### Problem

TV Shows have their own `tvshow` schema. Movies are dumped in `public` alongside
shared infrastructure tables (`libraries`, `activity_log`, `library_permissions`).
This is inconsistent and violates the "each content module owns its schema" principle.

**Current layout (broken):**

| Schema | Contains |
|--------|----------|
| `shared` | users, auth, sessions, settings, RBAC, MFA, OIDC (~17 tables) |
| `public` | libraries + activity **mixed with** movies, movie_files, movie_credits, movie_collections, movie_genres, movie_watched |
| `tvshow` | series, seasons, episodes, episode_files, series_credits, series_genres, networks (~10 tables) |
| `qar` | empty (later: own DB) |

**Target layout:**

| Schema | Contains |
|--------|----------|
| `shared` | users, auth, sessions, settings, RBAC, MFA, OIDC |
| `public` | libraries, library_scans, library_permissions, activity_log |
| `movie` | movies, movie_files, movie_credits, movie_collections, movie_collection_members, movie_genres, movie_watched |
| `tvshow` | (unchanged) |
| `qar` | later: separate DB + own pool |

### Steps

- [x] **9.1** New migration `000036_create_movie_schema.up.sql`:
  - `CREATE SCHEMA IF NOT EXISTS movie;`
  - `ALTER TABLE public.movies SET SCHEMA movie;`
  - `ALTER TABLE public.movie_files SET SCHEMA movie;`
  - `ALTER TABLE public.movie_credits SET SCHEMA movie;`
  - `ALTER TABLE public.movie_collections SET SCHEMA movie;`
  - `ALTER TABLE public.movie_collection_members SET SCHEMA movie;`
  - `ALTER TABLE public.movie_genres SET SCHEMA movie;`
  - `ALTER TABLE public.movie_watched SET SCHEMA movie;`
  - Down migration moves them back to `public`
- [x] **9.2** Update `sqlc.yaml`: Add rename mappings (singularized keys) to preserve Go type names
- [x] **9.3** Update SQL query files: `public.movie*` → `movie.movie*` in all 587 lines
- [x] **9.4** Update connection `search_path` to `public,shared,movie,tvshow` in:
  - `internal/testutil/containers.go` (tests)
  - `docker-compose.dev.yml`
  - `docker-compose.prod.yml`
  - `.devcontainer/docker-compose.yml`
- [x] **9.5** Regenerate sqlc (`sqlc generate`) — all types correct: Movie, MovieFile, etc.
- [x] **9.6** Build + vet + test — all 42 packages pass
- [x] **9.7** Update `docs/dev/design/infrastructure/DATABASE.md` to reflect schema-per-module convention — `d59f2b88`

### Architecture Decision

- **Movie/TVShow/shared**: Schemas in one DB (cross-schema JOINs for `user_id` FKs, shared buffer pool)
- **QAR**: Separate DB (`revenge_qar`) + own connection pool when implemented (full isolation, separate backups, deletable, compliance)

---

## 10. Dependency Audit

### 10A — Circuit Breaker

**Decision: gobreaker NOT needed.** All external provider API calls (TMDb, Trakt, Simkl, OMDb,
Letterboxd, Fanart) must run as **River jobs** — River provides retry with exponential backoff,
max attempts, and dead-letter. Adding a circuit breaker on top of River would be double
resilience logic and counterproductive.

Internal infrastructure (Typesense, Redis/Dragonfly, PostgreSQL) already has built-in
resilience: Typesense has own retry, Rueidis has auto-reconnect, pgxpool has health checks.

The synchronous external calls currently in the request path (metadata fetch, rating enrichment)
are **bugs** — they must be moved to River jobs (see 3B.5, 3B.6).

`gobreaker` can be removed from `go.mod` after `go mod tidy`.

### 10B — Other Issues

| # | Issue | Severity | Fix |
|---|-------|----------|-----|
| 1 | **`lib/pq` + `pgx/v5` — two PG drivers** | Medium | `lib/pq` only used in 1 smoke test (`tests/live/smoke_test.go:22`) as `_ "github.com/lib/pq"`. Switch to `_ "github.com/jackc/pgx/v5/stdlib"` and remove `lib/pq` |
| 2 | **`hashicorp/raft-boltdb` uses archived BoltDB** (`boltdb/bolt`) | Low | Consider `hashicorp/raft-boltdb/v2` which uses `go.etcd.io/bbolt` |
| 3 | **Empty `pkg/` directory** | Low | Contains only `.gitkeep`, no Go files import from it — delete or document |

### 10C — Needs Verification Against Package Docs

Before implementing cache-related fixes, verify against current package documentation:

| # | Claim | Package | What to verify |
|---|-------|---------|---------------|
| 1 | Otter supports `DeleteByFunc` for prefix invalidation | `maypok86/otter` | Check if method exists in current API |
| 2 | We use `KEYS` for Redis pattern invalidation | `redis/rueidis` | Check actual impl — might already use `SCAN` |
| 3 | `*river.Client[any]` vs `[pgx.Tx]` type mismatch | `riverqueue/river` | ✅ Verified and FIXED — `aac680d5` changed queue.go to use `*infrajobs.Client` |

---

## Priority Order

### Tier 1 — Must Fix (Broken / Architecture)

- [x] **9.1-9.7**: Move movie tables to own `movie` schema — `7036592b`
- [x] **3B.2**: Register workers for 5+1 metadata job kinds — `aac680d5`
- [x] **3B.4**: Wire periodic cleanup jobs into River via fx — `6157182e`
- [x] **3B.5**: Move movie metadata refresh to River job (return 202) — `37717e60`
- [x] **3B.6**: Move rating enrichment to background workers — `d2281594`
- [x] **4B.1**: Create `CachedTVShowService` + wire DI (movie too) — `962f8192`
- [x] **4B.2**: ~~Create `CachedTVShowSearchService`~~ ✅ `796971ec` — created + wired in DI
- [x] **1A.1**: Implement WebAuthn `VerifyMFA` ✅ `9bdbc1b9` — VerifyWebAuthn on MFAManager, wired in VerifyMFA
- [x] **3B.1**: Wire auth/login, streaming, SSE jobs to `critical` queue ✅ `13c480e9` — AsyncLogger routes security actions to critical queue

### Tier 2 — High Value (API Readiness + Dedup)

- [x] **2B.4**: Add pagination to ALL list endpoints (`{items, total, page, page_size}`) ✅ `4d8f8919` — 7 endpoints updated, proper caching with pagination in cache keys
- [x] **2B.3**: Standardize query param naming to snake_case ✅ `9a9babe2` — renamed `orderBy`→`order_by`, `minVotes`→`min_votes` in OpenAPI spec + regenerated ogen
- [x] **2B.5**: Add admin user management endpoints (list+search, delete)
- [ ] **5B.1**: Index episodes in Typesense (with series context) ✅ `72056505` — EpisodeDocument (20 fields), EpisodeSearchService, CachedEpisodeSearchService, worker integration, 19 tests
- [x] **5B.2**: Index seasons in Typesense (named seasons)
- [x] **5B.3**: Create standalone people collection in Typesense (own pages)
- [x] **5B.4**: Add `/search/multi` unified search endpoint
- [x] **7.1**: Extract `requireAdmin()` helper — migrated 42 handlers across 7 files (commit `6523726e`)
- [x] **7.4**: Create `internal/util/ptr/` generic pointer package — `ptr.To`, `ptr.Value`, `ptr.ValueOr`, etc (commit `3847c32b`)
- [x] **7.5**: Create `BaseMovieProvider`/`BaseTVShowProvider` embedding structs ✅ `c880e1eb` — base structs + tests, embedded in 10 providers, removed ~85 stub methods
- [x] **2B.1**: Add CSP middleware ✅ `87f458c6` — SecurityHeadersMiddleware with CSP, X-Frame-Options, nosniff, Referrer-Policy, Permissions-Policy
- [x] **10B.1**: Remove `lib/pq`, switch smoke test to `pgx/v5/stdlib` ✅ `8e4a3b3e` — switched driver to pgx, lib/pq remains indirect via embedded-postgres

### Tier 3 — Medium Value

- [x] **3B.7**: Make activity logging async (River job on `low` queue) ✅ done in 3B.1 (`13c480e9`)
- [x] **7.2**: Extract `externalRatingsToOgen()` helper ✅ `cc2d0d4c` — shared `content.ExternalRating` type + type aliases + one converter replacing 4 copies
- [x] **7.6**: Move test helpers to `internal/testutil/helpers.go` ✅ `56678c6a` — added `UniqueUser()`, deduplicated 3 `createTestUser` functions
- [x] **2B.2**: Add Cache-Control headers on JSON API responses ✅ `42ee5a3b` — ogen middleware categorizing ~50 cacheable catalog ops (`private, max-age=60`) vs user-specific (`private, no-store`) vs mutations (`no-store`)
- [x] **2B.6**: Add `GET /api/v1/genres` endpoint — `deb694c1` unified genre aggregation from movies + TV shows
- [x] **2B.7**: Add bulk episode watched endpoint — `9d075c40` batch SQL with CTE, POST /api/v1/tvshows/episodes/bulk-watched
- [x] **2B.8**: Add TV show reindex endpoint — `db0a2f21` POST /api/v1/search/tvshows/reindex, enqueues River job
- [x] **5B.5**: Create users collection in Typesense (admin search)
- [x] **1A.2**: Fix or rename `GetAPIKeyUsageCount` placeholder ✅ `ed2028b9` — renamed to `GetAPIKeyLastUsedAt`, removed misleading comments
- [x] **7.3**: Extract `copyMovieFieldsToOgen` helper — **SKIPPED**: ogen generates flat types with no shared interface; any extraction (19-field interface / reflection) adds more complexity than the 3 copies of schema-driven boilerplate
- [x] **4B.3**: Fix L1 pattern invalidation ✅ `bd2be17f` — `DeleteByPrefix()` on Otter Keys() iterator, `simpleGlobPrefix()` detects prefix* patterns
- [x] **4B.4**: Verify + fix Redis pattern invalidation ✅ `bd2be17f` — replaced `KEYS` with cursor-based `SCAN` + batch delete

### Tier 4 — Polish

- [x] **8B.3**: Delete dead `database/metrics.go` + test ✅ `033bc406`
- [x] **8B.4**: Smoke test Grafana dashboards after UID fix ✅ validated: all 32 metric names match Go source, datasource UID `PBFA97CFB590B2093` consistent, dashboard JSON valid (schemaVersion 39, 27 panels)
- [x] **3B.3**: Create `EnrichContentWorker` — done in 3B.2 (`aac680d5`)
- [x] **3C.9**: Create `DownloadImageWorker` for cache warming
- [x] **3C.11**: Periodic stats/analytics aggregation job ✅ `8896fdbb` — migration 000037, StatsAggregationWorker (hourly, low queue), 10 aggregate stat keys, 13 tests
- [x] **1C.5-7**: Clean up TODO comments and stale test comments ✅ `72f19cdc`
- [x] **1D.9-10**: Delete vestigial sqlc placeholders ✅ `033bc406` — removed source SQL + regenerated sqlc
- [x] **10B.2**: Upgrade `raft-boltdb` to v2 ✅ `c837af04` — switched to v2 (go.etcd.io/bbolt), API-compatible
- [x] **10B.3**: Delete empty `pkg/` directory ✅ `033bc406`
- [x] **3B.8**: Remove bare goroutine fallback in radarr/sonarr handlers ✅ `f6326dab` — return 503 when River unavailable

---

## Notes

- **Circuit Breaker**: NOT needed — all external API calls run as River jobs with built-in retry/backoff. Internal infra has own resilience. `gobreaker` removed.
- **Logger**: Already clean and unified on `slog` — no work needed
- **API structure**: Solid overall. OpenAPI 3.1.0 spec, ogen codegen, fx DI, proper error handling, SSE events
- **Caching**: Excellent L1/L2 architecture — TV show coverage is the critical gap
- **River jobs**: Good queue design. `critical` queue reserved for auth/login, streaming, SSE, frontend-critical tasks. Several orphaned job definitions. Synchronous external calls in request path must move to jobs.
- **Search**: Strong movie/tvshow coverage. Needs episodes, seasons, people, users, and unified `/search/multi` endpoint
- **Frontend readiness**: 7/10 — auth, CORS, cookies, SSE all ready. Needs: pagination on all lists, snake_case params, CSP, cache headers, admin endpoints
- **Schema**: Movie module inconsistently in `public` — needs own `movie` schema to match `tvshow` pattern
- **QAR**: Future separate DB (`revenge_qar`) — not current scope
- **Package docs**: Cache-related fixes (Otter, Rueidis, River) need verification against current package docs before implementation
