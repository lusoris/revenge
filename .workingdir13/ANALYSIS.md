# Workingdir 13 — Codebase Audit Report

**Date:** 2026-02-12
**Scope:** Code Duplication · L1/L2 Cache Coverage · Flow Integration · River Jobs

---

## 1. Code Duplication Analysis

### Estimated total duplicated/near-duplicated code: ~3,400+ lines

### 1.1 Radarr / Sonarr Full-Package Clone — **HIGH SEVERITY**

The single largest duplication cluster. `internal/integration/radarr/` and `internal/integration/sonarr/` are structural clones across **12 files (6+6), ~1,500+ duplicated lines**.

| Component | Radarr File | Sonarr File | Duplicated |
|-----------|------------|------------|------------|
| Client struct, constructor, cache helpers | `radarr/client.go` | `sonarr/client.go` | Byte-for-byte identical except "radarr"→"sonarr" |
| Every API method (15-20 each) | `radarr/client.go` (658 lines) | `sonarr/client.go` (809 lines) | Same template: cache check → rate limit → HTTP → cache set |
| Shared types (MediaInfo, Quality, QualityProfile, RootFolder, SystemStatus, etc.) | `radarr/types.go` | `sonarr/types.go` | Byte-for-byte identical structs |
| SyncService delegation | `radarr/service.go` | `sonarr/service.go` | Identical pattern (SyncStatus, SyncResult, mutex guard) |
| Job workers | `radarr/jobs.go` (200 lines) | `sonarr/jobs.go` (203 lines) | Structural clone (SyncOp, SyncJobArgs, SyncWorker, WebhookJobArgs, WebhookWorker) |
| Webhook handler | `radarr/webhook_handler.go` | `sonarr/webhook_handler.go` | Same event-type switch |

**Recommendation:** Create `internal/integration/arrbase` package with:
- Shared types (MediaInfo, Quality, QualityProfile, RootFolder, SystemStatus, Image, Language)
- Generic `ArrClient[T any]` base client with cache+rate-limiter+HTTP pattern
- Generic `SyncWorker[TArgs river.JobArgs]` base worker
- Shared `SyncStatus`/`SyncResult` with adapter for Movies vs Series naming

### 1.2 Admin Authorization Boilerplate — **HIGH SEVERITY**

Every admin-only endpoint repeats the same 7-line authorization block (~150 duplicated lines across 21+ sites):

```go
if _, err := h.requireAdmin(ctx); err != nil {
    if errors.Is(err, errNotAuthenticated) {
        return &ogen.XxxUnauthorized{Code: 401, Message: "Authentication required"}, nil
    }
    if errors.Is(err, errNotAdmin) {
        return &ogen.XxxForbidden{Code: 403, Message: "Admin access required"}, nil
    }
    return nil, err
}
```

**Files:** `handler_radarr.go` (4x), `handler_sonarr.go` (4x), `handler_rbac.go` (11x), `handler_admin_users.go` (3x+), `handler_library.go` (multi), `handler_oidc.go` (2x)

**Recommendation:** Generic `adminGuard[T401, T403 any]()` helper using Go generics, or middleware-level admin enforcement.

### 1.3 API Handler Radarr/Sonarr Status/Sync — **HIGH SEVERITY**

`handler_radarr.go` (338 lines) vs `handler_sonarr.go` (394 lines) — every handler method pair is a structural clone (~400 duplicated lines). Both share methods like `IsHealthy`, `GetStatus`, `GetSystemStatus`, `GetQualityProfiles`, `GetRootFolders`.

**Recommendation:** Unify into a single `*arrService` interface and write the handler logic once.

### 1.4 Cached Service Decorator Pattern — **MEDIUM SEVERITY**

~500 duplicated lines across 4 files. Every cached method repeats:

```go
func (s *CachedService) GetXxx(ctx, ...) (*T, error) {
    if s.cache == nil { return s.Service.GetXxx(ctx, ...) }
    cacheKey := cache.XxxKey(...)
    var result T
    if err := s.cache.GetJSON(ctx, cacheKey, &result); err == nil { return &result, nil }
    result, err := s.Service.GetXxx(ctx, ...)
    if err != nil { return nil, err }
    go func() { s.cache.SetJSON(cacheCtx, cacheKey, result, ttl) }()
    return result, nil
}
```

**Files:** `movie/cached_service.go` (471 lines, ~12 methods), `tvshow/cached_service.go` (576 lines, ~15 methods), `library/cached_service.go` (223 lines), `session/cached_service.go` (188 lines)

**Recommendation:** Generic `CacheAside[T any]` helper function that collapses each method to ~3 lines.

### 1.5 Localize-then-Convert Pattern — **MEDIUM SEVERITY**

~100 duplicated lines across 10+ sites in `movie_handlers.go` and `tvshow_handlers.go`.

**Recommendation:** Generic `localizeAndConvert[TDomain, TOgen any]` helper.

### 1.6 Repository Error Wrapping — **MEDIUM SEVERITY**

~300 duplicated lines across 50+ sites in `movie/repository_postgres.go` (1030 lines) and `tvshow/repository_postgres.go` (1606 lines).

**Recommendation:** Generic `wrapQuery[T, R any]` helper.

### 1.7 Ogen Optional Field Mapping — **LOW-MEDIUM SEVERITY**

~400 duplicated lines in `movie_converters.go` (438 lines) and `tvshow_converters.go` (384 lines).

**Recommendation:** Helper functions like `setOptString(opt *ogen.OptString, val *string)`.

### 1.8 Service Thin Delegation — **LOW SEVERITY**

~60 duplicated lines in `radarr/service.go` and `sonarr/service.go` — pure delegation.

**Recommendation:** Embed Client in SyncService or use interface composition.

---

## 2. L1/L2 Cache Coverage Analysis

### 2.1 Architecture Overview

| Layer | Technology | Scope |
|-------|-----------|-------|
| **L1** | otter (W-TinyLFU eviction) | In-process, per-instance `cache.L1Cache[K,V]` |
| **L2** | rueidis → Dragonfly/Redis | Distributed, cross-instance `cache.Cache` |

- **Read path:** L1 → L2 (DoCache server-assisted) → populate L1 on L2 hit
- **Write path:** Write L1 + L2 simultaneously with TTL
- **Invalidation:** L1 prefix-based `DeleteByPrefix`; L2 `SCAN` + batch `DEL`

### 2.2 What IS Properly Cached (L1+L2)

| Service | Cached Methods |
|---------|---------------|
| **Session** | ValidateSession, CreateSession, RevokeSession, RevokeAllUserSessions |
| **RBAC** | Enforce, EnforceWithContext, GetUserRoles, HasRole, AssignRole, RemoveRole, policies |
| **Settings** | Server/User settings CRUD (except category-filtered variants) |
| **User** | GetUser, GetUserByUsername, UpdateUser, DeleteUser |
| **Library** | Get, List, Count, Create, Update, Delete, CompleteScan |
| **Movie** | GetMovie, ListMovies, RecentlyAdded, TopRated, Cast, Crew, Genres, Collection, ContinueWatching, Update, Delete, WatchProgress |
| **TV Show** | GetSeries, ListSeries, RecentlyAdded, Seasons, Episodes, Cast, Crew, Genres, Networks, ContinueWatching, WatchProgress, etc. |
| **Search** | Movie/TVShow/Season/Episode/Person/User search + Autocomplete + Facets |

### 2.3 L1-Only Caching (per-process)

| Component | What's Cached |
|-----------|---------------|
| **12 metadata providers** (TMDb, TVDb, Trakt, OMDb, TVMaze, MAL, Kitsu, Simkl, AniDB, AniList, FanartTV, Letterboxd) | HTTP API responses |
| **Radarr/Sonarr clients** | API responses |
| **Playback sessions** | Active sessions |
| **Playback probe** | ffprobe media info |
| **Transcode** | Active FFmpeg processes (with OnDeletion cleanup) |
| **HLS** | Master & media playlists |
| **Rate limiter** | Per-IP rate limiters |
| **MFA/WebAuthn** | Challenge sessions (actually uses L1+L2) |
| **Image service** | Filesystem cache, 7-day TTL |

### 2.4 NOT Cached But SHOULD Be — **HIGH Priority**

| Component | Method | Impact |
|-----------|--------|--------|
| **API Keys: `ValidateKey`** | Every API-key-authenticated request hits DB | **DB round-trip per request** |
| **Movie: `GetMovieByTMDbID`/`GetMovieByIMDbID`** | Called heavily during library scan/import | DB hit every scan |
| **Movie: `GetMovieFiles`** | Called on every playback initiation | DB round-trip per play |
| **Movie: `SearchMovies`** (content DB search) | UI search | Repeated DB queries |
| **Movie: `GetMoviesByGenre`** | Browse-by-genre | DB hit on every genre browse |
| **Movie/TVShow: `ListDistinctGenres`** | Filter sidebar, rarely changes | DB for near-static data |
| **TV Show: `ListByGenre`/`ListByNetwork`/`ListByStatus`** | Browse filters | DB hits on every filter |
| **TV Show: `ListRecentEpisodes`/`ListUpcomingEpisodes`** | Homepage widgets | Expensive queries on hot path |
| **TV Show: `GetEpisodeGuestStars`/`GetEpisodeCrew`** | Episode detail | DB per view |

### 2.5 NOT Cached But SHOULD Be — **MEDIUM Priority**

| Component | Method | Impact |
|-----------|--------|--------|
| **User: `GetUserByEmail`** | Auth lookup | DB per login |
| **User: `GetUserPreferences`** | Per-session load | DB per page load |
| **Library: `ListEnabled`/`ListByType`/`ListAccessible`** | Filter variants of List | DB for filtered views |
| **Library: `CheckPermission`/`CanAccess`/`CanDownload`/`CanManage`** | Called on every library access | DB per permission check |
| **Settings: `ListServerSettingsByCategory`/`ListUserSettingsByCategory`** | Not cached while other list methods are | Bypasses cache |
| **Settings: `ListUserSettings`** | UserSetting list | DB per full preference load |
| **Movie: `GetWatchHistory`/`GetUserStats`** | Profile pages | Per-user DB queries |
| **TV Show: `GetNextEpisode`/`GetSeriesWatchStats`** | Dashboard | Per-user DB queries |
| **Movie/TVShow: `CountMovies`/`CountSeries`** | Dashboard widget | Aggregate queries |

### 2.6 Cache Invalidation Issues

| Issue | Location | Severity |
|-------|----------|----------|
| **OIDC uses `sync.Map` instead of `L1Cache`** | `oidc/service.go` — unbounded memory, no TTL, violates project rules | Medium |
| **Session `RefreshSession` not overridden** in CachedService | Old token hash remains cached up to 30s after refresh | Medium |
| **Session `RevokeAllUserSessionsExcept` not overridden** | Revoked sessions served from cache until TTL | Medium |
| **Library permission changes don't invalidate** `ListAccessible`/`CanAccess` | `GrantPermission`/`RevokePermission` not overridden in CachedService | Medium |
| **Movie collection queries not invalidated** | `GetMoviesByCollection`/`GetCollectionForMovie` missed | Low |

---

## 3. Unintegrated / Partially Integrated Flows

### 3.1 Fully Integrated ✅

| Flow | Evidence |
|------|----------|
| Library Scanning (Movie) | API → River Job → Worker → complete chain |
| Metadata Refresh (Movie + TV Show) | API → River Job → Worker → Provider chain |
| Activity Logging | AsyncLogger → River → multiple services consume |
| Auth Token Cleanup | Periodic daily job |
| Radarr/Sonarr Integrations | Webhooks + API sync endpoints complete |
| Playback Lifecycle | Start/Get/Stop endpoints, session mgmt, cleanup job |
| SSE Real-Time Events | Agent registered, route active |
| Search Indexing (manual) | API trigger → River → Workers for both movie+tvshow |
| Stats Aggregation | Hourly periodic job with RunOnStart |
| RBAC, Session, API Keys, MFA, OIDC, Email, Settings, Storage | All fully wired |

### 3.2 Gaps & Partially Integrated ⚠️

#### **Notification Dispatch — Wired but ZERO callers** (HIGH)
- Full notification stack exists: `Dispatcher` → `NotificationWorker` → Discord/Gotify/Ntfy/Webhook/Email agents
- SSE agent registered at startup
- **BUT: No production code anywhere enqueues `NotificationArgs` jobs**
- Event types defined (`movie.added`, `library.scan_done`, `playback.started`) but never emitted
- **Fix:** Add dispatch calls at: library scan completion, movie/show added, playback start/stop, auth events

#### **TV Show Library Scan — No API trigger** (HIGH)
- `TriggerLibraryScan` in `handler_library.go` only enqueues `MovieLibraryScanArgs`
- TV show scan workers exist and are registered, but no API path enqueues TV show jobs
- **Fix:** Check library type and enqueue appropriate job type (`MovieLibraryScanArgs` vs `tvshowjobs.LibraryScanArgs`)

#### **Notification Agent Configuration — Only SSE auto-registered** (HIGH)
- Discord, Gotify, Ntfy, Webhook, Email agents fully implemented
- Only SSE agent is registered at startup
- No code reads agent config from config file/DB
- Users cannot configure notification targets
- **Fix:** Add config-based agent registration at startup

#### **Periodic Library Scan — Missing** (MEDIUM)
- Periodic jobs include auth/activity/library cleanup, playback health, stats
- No periodic library scanning job exists
- Libraries must always be scanned manually
- **Fix:** Add configurable periodic scan job (e.g., daily schedule)

#### **Search Index Updates on Content Change — Partial** (MEDIUM)
- Manual reindex via API works
- No automatic index update when content changes (after scan adds movie, after metadata refresh updates show)
- Search index and content DB can drift apart
- **Fix:** Enqueue `MovieSearchIndexArgs`/`SearchIndexArgs` at end of scan/refresh workers

#### **Playback Heartbeat — Missing** (MEDIUM)
- API only has Start/Get/Stop — no heartbeat/progress-reporting endpoint
- No way for clients to report playback position or extend session timeout during viewing
- **Fix:** Add `POST /api/v1/playback/sessions/{sessionId}/heartbeat` endpoint

#### **RefreshPersonWorker — Stub/No-Op** (MEDIUM)
- Worker registered and accepts jobs but does nothing
- `RefreshPersonArgs` are enqueued from metadata pipeline
- No `Person` content service exists
- Jobs succeed silently without doing work
- **Fix:** Return `river.JobCancel` or don't enqueue until person service exists

#### **DownloadImageWorker — Stale Comment** (LOW)
- Worker code works correctly, delegates to `imageService.FetchImage()`
- Comment says "stub — will use image service when available" but implementation is complete
- **Fix:** Update stale comment

---

## 4. River Background Job System Audit

### 4.1 Configuration Summary

| Setting | Value |
|---------|-------|
| MaxAttempts (global default) | 25 |
| FetchCooldown | 200ms |
| FetchPollInterval | 2s |
| RescueStuckJobsAfter | 30min |
| MaxWorkers (total cap) | 100 |

**5-level priority queue system:**

| Queue | Default MaxWorkers | Purpose |
|-------|-------------------|---------|
| `critical` | 20 | Security events, auth failures |
| `high` | 15 | Notifications, webhooks |
| `default` | 10 | Metadata fetch, sync |
| `low` | 5 | Cleanup, maintenance |
| `bulk` | 3 | Library scans, batch ops |

**Periodic Jobs:**

| ID | Schedule | Kind | RunOnStart |
|----|----------|------|------------|
| `auth_cleanup_daily` | 24h | `cleanup` | Yes |
| `activity_cleanup_daily` | 24h | `activity_cleanup` | No |
| `library_scan_cleanup_daily` | 24h | `library_scan_cleanup` | No |
| `playback_health_check` | 5min | `playback_cleanup` | No |
| `stats_aggregation_hourly` | 1h | `stats_aggregation` | Yes |

### 4.2 All 27 Job Types — Registration Status

**All 26 workers are properly registered. No orphaned workers found.** ✅

Complete list:

| Kind | Queue | MaxAttempts | Timeout | Unique | Registered |
|------|-------|-------------|---------|--------|------------|
| `cleanup` | `low` | 5 | 2min | None | ✅ |
| `notification` | `high` | 5 | 2min | ByArgs+ByPeriod(1h) | ✅ |
| `activity_log` | `low`/`critical` | 3 | 10s | None | ✅ |
| `activity_cleanup` | `low` | **25 (default)** | 2min | None | ✅ |
| `stats_aggregation` | `low` | **25 (default)** | 2min | None | ✅ |
| `library_scan_cleanup` | `low` | **25 (default)** | 2min | None | ✅ |
| `movie_library_scan` | `bulk` | **25 (default)** | 30min | None | ✅ |
| `metadata_refresh_movie` | `default` | **25 (default)** | 5min | None | ✅ |
| `movie_file_match` | `default` | **25 (default)** | 5min | None | ✅ |
| `movie_search_index` | `default` | **25 (default)** | 15min | None | ✅ |
| `tvshow_library_scan` | `bulk` | **25 (default)** | 30min | None | ✅ |
| `tvshow_metadata_refresh` | `default` | **25 (default)** | 15min | None | ✅ |
| `tvshow_file_match` | `default` | **25 (default)** | 5min | None | ✅ |
| `tvshow_search_index` | `bulk` | **25 (default)** | default | None | ✅ |
| `tvshow_series_refresh` | `default` | **25 (default)** | 10min | None | ✅ |
| `metadata_refresh_tvshow` | `default` | **25 (default)** | 10min | None | ✅ |
| `metadata_refresh_season` | `default` | **25 (default)** | 5min | None | ✅ |
| `metadata_refresh_episode` | `default` | **25 (default)** | 5min | None | ✅ |
| `metadata_refresh_person` | `default` | **25 (default)** | 2min | None (stub!) | ✅ |
| `metadata_enrich_content` | `default` | **25 (default)** | 5min | None | ✅ |
| `metadata_download_image` | `default` | **25 (default)** | 2min | None | ✅ |
| `radarr_sync` | `high` | **25 (default)** | 10min | None | ✅ |
| `radarr_webhook` | `high` | **25 (default)** | 1min | None | ✅ |
| `sonarr_sync` | `high` | **25 (default)** | 10min | None | ✅ |
| `sonarr_webhook` | `high` | **25 (default)** | 1min | None | ✅ |
| `playback_cleanup` | `low` (hardcoded!) | **25 (default)** | 1min | ByPeriod(5min) | ✅ |

### 4.3 Issues Found

#### HIGH

**H-1: 23 of 27 job types use global MaxAttempts=25 — WAY too many retries**

External API jobs (radarr/sonarr sync, all metadata refreshes, image download) will pound upstream services for hours with 25 retries + exponential backoff.

**Fix:** Set explicit MaxAttempts per job type:
- External API jobs: 3–5 attempts
- Library scans: 3 attempts
- Search indexing: 3–5 attempts
- Cleanup/maintenance: 3 attempts

**H-2: 13 metadata/content args types have NO `InsertOpts()` method**

`RefreshMovieArgs`, `RefreshTVShowArgs`, `RefreshSeasonArgs`, `RefreshEpisodeArgs`, `RefreshPersonArgs`, `EnrichContentArgs`, `DownloadImageArgs`, `MovieFileMatchArgs`, `MovieSearchIndexArgs`, `MetadataRefreshArgs` (tvshow), `FileMatchArgs` (tvshow), `SeriesRefreshArgs` — all inherit River defaults.

**Fix:** Add `InsertOpts()` to every job args type with appropriate queue, max attempts, unique constraints.

**H-3: No unique constraints on most job types — duplicate job accumulation**

Only `notification` and `playback_cleanup` use UniqueOpts. Vulnerable:
- `radarr_sync` / `sonarr_sync` — hammering sync button creates unbounded duplicates
- `movie_library_scan` — duplicate scans for same library
- All metadata refreshes — no dedup at all

**Fix:** Add UniqueOpts:
- `radarr_sync`: `ByArgs: true, ByPeriod: 5min`
- `sonarr_sync`: `ByArgs: true, ByPeriod: 5min`
- `movie_library_scan`: `ByArgs: true, ByPeriod: 10min`
- All metadata refreshes: `ByArgs: true, ByPeriod: 30min`

#### MEDIUM

**M-1: `movie_file_match` returns error on "unmatched", triggering 25 retries**

When a file can't be matched, the worker returns an error that triggers 25 retry attempts for a deterministic failure.

**Fix:** Return `nil` and log, or return `river.JobCancel` error.

**M-2: `playback_cleanup` hardcodes `"low"` instead of `QueueLow` constant**

Inconsistent and fragile.

**M-3: No dead-letter / discard policy**

No workers implement `NextRetry()` for custom backoff. A metadata refresh that gets 404 from TMDb will retry 25 times.

**Fix:** Detect HTTP 404/410 and return non-retryable error.

**M-4: `stats_aggregation` panics if queries is nil**

`collectStats()` calls `panic(...)` — will crash worker goroutine.

**Fix:** Replace with `return nil, fmt.Errorf(...)`.

**M-5: No `InsertTx` usage anywhere**

All job insertions happen outside DB transactions. If a service creates a resource then enqueues a job, failure between the two leaves inconsistent state.

#### LOW

**L-1: No `CompletedJobRetention` configured**

Completed/discarded jobs accumulate in `river_job` table indefinitely.

**Fix:** Set `CompletedJobRetention: 24h–72h`, `DiscardedJobRetention: 7d`.

**L-2: Inconsistent error return patterns across workers**

Some aggregate errors, some return on first error, some return nil on partial success.

**L-3: `playback_cleanup` ignores context**

Uses `_ context.Context` — never checks cancellation.

---

## 5. Prioritized Action Items

### Immediate (Before Next Release)

1. **[River] Add `MaxAttempts` + `UniqueOpts` to all 23 job types missing them** — prevents upstream API hammering and duplicate job storms
2. **[River] Fix `movie_file_match` non-retryable error** — stop 25x retry of deterministic failures
3. **[River] Fix `stats_aggregation` panic** — replace with error return
4. **[Flow] Fix TV show library scan API trigger** — enqueue correct job type based on library type
5. **[Cache] Cache `ValidateKey` in API Keys service** — eliminate DB round-trip per API-key request

### Short-Term (Next Sprint)

6. **[Flow] Wire notification dispatch** — add Dispatch() calls at lifecycle points (scan done, content added, playback events)
7. **[Flow] Auto-enqueue search index jobs** after library scan and metadata refresh
8. **[Cache] Cache `GetMovieFiles`** — called on every playback
9. **[Cache] Cache `ListDistinctGenres`** (movie + tvshow) — near-static data
10. **[Cache] Replace OIDC `sync.Map` with `L1Cache`** — violates project rules
11. **[Cache] Override `RefreshSession` + `RevokeAllUserSessionsExcept`** in session CachedService
12. **[River] Add `CompletedJobRetention` + `DiscardedJobRetention`** — prevent table growth

### Medium-Term (This Quarter)

13. **[Duplication] Create `internal/integration/arrbase` package** — deduplicate ~1,500 lines of radarr/sonarr clones
14. **[Duplication] Create generic `CacheAside[T any]` helper** — reduce ~500 lines across 4 cached services
15. **[Flow] Add periodic library scan job** — configurable schedule
16. **[Flow] Wire notification agent configuration** from config/DB
17. **[Flow] Add playback heartbeat endpoint**
18. **[Cache] Cache browse-path queries** (genre, network, status, recent/upcoming episodes)
19. **[Cache] Cache library permission checks** (CheckPermission, CanAccess, etc.)
20. **[River] Add `NextRetry()` custom backoff** for external API workers
21. **[River] Add error-type-based discard logic** (HTTP 404/410 → cancel)

### Long-Term (Technical Debt)

22. **[Duplication] Generic admin guard helper** to reduce auth boilerplate
23. **[Duplication] Generic `localizeAndConvert` + `wrapQuery` helpers**
24. **[Duplication] Helper functions for ogen optional field mapping**
25. **[Flow] Implement Person metadata refresh** (currently stub no-op)
