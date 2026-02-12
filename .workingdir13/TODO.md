# Workingdir 13 — Master TODO

## Phase 1: Critical Bugs & River Jobs (Priority: IMMEDIATE)

- [x] 1.1 Create .workingdir13 directory
- [x] 1.2 Add `InsertOpts()` with `MaxAttempts` + `UniqueOpts` to ALL job args types that lack them
  - [x] `metadata/jobs/refresh.go`: RefreshMovieArgs, RefreshTVShowArgs, RefreshSeasonArgs, RefreshEpisodeArgs, RefreshPersonArgs, EnrichContentArgs, DownloadImageArgs
  - [x] `moviejobs/file_match.go`: MovieFileMatchArgs
  - [x] `moviejobs/search_index.go`: MovieSearchIndexArgs
  - [x] `tvshow/jobs/jobs.go`: MetadataRefreshArgs, FileMatchArgs, SeriesRefreshArgs
  - [x] `radarr/jobs.go`: RadarrSyncJobArgs, RadarrWebhookJobArgs (add MaxAttempts + UniqueOpts)
  - [x] `sonarr/jobs.go`: SonarrSyncJobArgs, SonarrWebhookJobArgs (add MaxAttempts + UniqueOpts)
  - [x] `activity/cleanup.go`: ActivityCleanupArgs (add MaxAttempts)
  - [x] `analytics/stats_worker.go`: StatsAggregationArgs (add MaxAttempts)
  - [x] `library/cleanup.go`: LibraryScanCleanupArgs (add MaxAttempts)
  - [x] `playback/jobs/cleanup.go`: CleanupArgs (use QueueLow constant instead of "low")
- [x] 1.3 Fix `movie_file_match` returning error on unmatched files (causes 25 retries)
- [x] 1.4 Fix `stats_aggregation` panic when queries is nil
- [x] 1.5 Add `CompletedJobRetentionPeriod` + `DiscardedJobRetentionPeriod` to River config
- [x] 1.6 Reduce global default MaxAttempts from 25 → 5
- [x] 1.7 Add `JobTimeout: -1` to River config (per-worker Timeout() handles this)

## Phase 2: Flow Integration Fixes (Priority: HIGH)

- [x] 2.1 Fix TV show library scan API trigger (currently only enqueues movie scan)
- [x] 2.2 Wire notification dispatch at lifecycle points
- [x] 2.3 Auto-enqueue search index jobs after library scan completion
- [x] 2.4 Fix RefreshPersonWorker to return `river.JobCancel` instead of silently succeeding

## Phase 3: Cache Coverage (Priority: HIGH)

- [x] 3.1 Fix OIDC `sync.Map` → `L1Cache` (violates project rules)
- [x] 3.2 Fix session CachedService: override `RefreshSession` + `RevokeAllUserSessionsExcept`
- [ ] 3.3 Cache `ValidateKey` in API Keys service (requires interface extraction — deferred)
- [x] 3.4 Cache `GetMovieFiles` (called on every playback)
- [x] 3.5 Cache `ListDistinctGenres` (movie + tvshow, near-static data)

## Phase 4: Code Deduplication (Priority: MEDIUM)

- [ ] 4.1 Create `internal/integration/arrbase` shared package (types, client base)
- [x] 4.2 Create generic `CacheAside[T]` helper function
- [ ] 4.3 Create `setOpt*` helpers for ogen optional field mapping

## Phase 5: Additional Features (Priority: LOW)

- [x] 5.1 Add periodic library scan job (configurable schedule)
- [ ] 5.2 Add notification agent configuration from config/DB
- [ ] 5.3 Add playback heartbeat endpoint
- [x] 5.4 Update stale DownloadImageWorker comment
