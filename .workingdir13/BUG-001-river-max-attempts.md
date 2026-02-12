# BUG-001: 23 of 27 River job types use global MaxAttempts=25
**Status: RESOLVED**

**Severity:** HIGH  
**Category:** River Jobs  
**Status:** RESOLVED

## Description

Most job types either don't define `InsertOpts()` at all (inheriting global default of 25 retries) or define `InsertOpts()` without setting `MaxAttempts`. For jobs hitting external APIs (TMDb, Radarr, Sonarr), 25 retries with exponential backoff will pound upstream services for hours.

## Affected Files

### No `InsertOpts()` at all (inherit defaults for everything):
- `internal/service/metadata/jobs/refresh.go` — RefreshMovieArgs, RefreshTVShowArgs, RefreshSeasonArgs, RefreshEpisodeArgs, RefreshPersonArgs, EnrichContentArgs, DownloadImageArgs (7 types!)
- `internal/content/movie/moviejobs/file_match.go` — MovieFileMatchArgs
- `internal/content/movie/moviejobs/search_index.go` — MovieSearchIndexArgs
- `internal/content/tvshow/jobs/jobs.go` — MetadataRefreshArgs, FileMatchArgs, SeriesRefreshArgs

### Have `InsertOpts()` but missing `MaxAttempts`:
- `internal/integration/radarr/jobs.go` — RadarrSyncJobArgs, RadarrWebhookJobArgs
- `internal/integration/sonarr/jobs.go` — SonarrSyncJobArgs, SonarrWebhookJobArgs
- `internal/service/activity/cleanup.go` — ActivityCleanupArgs
- `internal/service/analytics/stats_worker.go` — StatsAggregationArgs
- `internal/service/library/cleanup.go` — LibraryScanCleanupArgs
- `internal/content/movie/moviejobs/library_scan.go` — MovieLibraryScanArgs
- `internal/content/tvshow/jobs/jobs.go` — LibraryScanArgs, SearchIndexArgs

### Properly configured (reference):
- `internal/infra/jobs/cleanup_job.go` — CleanupArgs (MaxAttempts: 5) ✓
- `internal/infra/jobs/notification_job.go` — NotificationArgs (MaxAttempts: 5, UniqueOpts) ✓
- `internal/service/activity/job.go` — ActivityLogArgs (MaxAttempts: 3) ✓

## Fix

Add appropriate `MaxAttempts` to every job type:
- External API jobs (radarr/sonarr sync, metadata refresh, image download): **5**
- Library scans: **3**
- Search indexing: **5**
- Cleanup/maintenance: **3**
- File matching: **3**
- Webhook processing: **5**

Also add `UniqueOpts` to prevent duplicate jobs where appropriate.
