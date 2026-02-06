# Arr Integration Pattern

> Prescriptive template for all Servarr integrations. Radarr (movies) and Sonarr (TV shows) are the reference implementations. Written from code as of 2026-02-06.

---

## Architecture: 8 Layers

Every arr integration follows the same layered architecture:

```
internal/integration/{arr}/
  types.go             # 1. API response types (match Servarr API v3)
  errors.go            # 2. Integration-specific sentinel errors
  client.go            # 3. HTTP client (req/v3, rate limiting, cache)
  mapper.go            # 4. Servarr types -> domain types
  service.go           # 5. SyncService (full + single sync)
  jobs.go              # 6. River workers (sync + webhook jobs)
  webhook_handler.go   # 7. Webhook event routing
  module.go            # 8. fx wiring (optional deps)
```

API handlers live in `internal/api/handler_{arr}.go` — they're part of the API layer, not the integration package.

---

## Layer Details

### 1. Types (`types.go`)

Mirror the Servarr API v3 response schemas exactly. JSON tags match API field names.

| Radarr | Sonarr | Purpose |
|--------|--------|---------|
| `Movie` | `Series` | Main content type |
| `MovieFile` | `EpisodeFile` | Media file metadata |
| `MediaInfo` | `MediaInfo` | Codec, bitrate, resolution |
| `Quality` | `Quality` | Quality definition + revision |
| `QualityProfile` | `QualityProfile` | Quality config |
| `RootFolder` | `RootFolder` | Monitored paths |
| `WebhookPayload` | `WebhookPayload` | Incoming webhook events |
| `Command` | `Command` | Background task tracking |
| `CalendarEntry` | — | Upcoming releases |
| — | `Episode` | Individual episode |
| — | `SeasonInfo` | Season stats + monitoring |

**Webhook events** (shared pattern, different content):

| Event | Radarr | Sonarr |
|-------|--------|--------|
| Grab | Release grabbed | Release grabbed |
| Download | File imported (**main sync trigger**) | Episode imported (**main sync trigger**) |
| Rename | Files renamed | Files renamed |
| Delete | `MovieDelete`, `MovieFileDelete` | `SeriesDelete`, `EpisodeFileDelete` |
| Health | System health issues | System health issues |
| Test | Webhook test | Webhook test |

### 2. Client (`client.go`)

All clients use `imroc/req/v3` with identical patterns:

```go
client := req.C().
    SetBaseURL(config.BaseURL + "/api/v3").
    SetTimeout(30 * time.Second).
    SetCommonHeader("X-Api-Key", config.APIKey).
    SetCommonHeader("Content-Type", "application/json").
    SetCommonRetryCount(3).
    SetCommonRetryBackoffInterval(1*time.Second, 10*time.Second)
```

**Rate limiting**: `golang.org/x/time/rate` — 10 req/s default, burst 20 (local service, higher than external APIs).

**Caching**: `sync.Map` with 5-minute TTL. Cache keys:
- `movies:all` / `series:all` — bulk lists
- `movie:{id}` / `series:{id}` — single items
- `movie:tmdb:{id}` / `series:tvdb:{id}` — external ID lookups
- `qualityprofiles`, `rootfolders` — config data

Cache is invalidated manually after write operations (Add/Delete).

**Error handling**: Specific sentinel errors per integration:

```go
// Radarr
ErrMovieNotFound = "movie not found in radarr"
ErrNotConfigured = "radarr integration not configured"
ErrUnauthorized  = "radarr api key is invalid"

// Sonarr
ErrSeriesNotFound      = "sonarr: series not found"
ErrEpisodeNotFound     = "sonarr: episode not found"
ErrEpisodeFileNotFound = "sonarr: episode file not found"
ErrConnectionFailed    = "sonarr: connection failed"
ErrRateLimited         = "sonarr: rate limited"
```

### 3. Mapper (`mapper.go`)

Converts Servarr API types to internal domain types. Stateless — no dependencies.

Key transformations:
- Generates new UUIDs for all domain objects
- Maps external IDs: Radarr uses `TMDbID`, Sonarr uses `TVDbID`
- `float64` ratings to `decimal.Decimal` (avoids float precision loss)
- Image URLs: prefers `RemoteURL`, falls back to `URL`
- Quality resolution: integer to string (`2160` -> `"4K"`, `1080` -> `"1080p"`)
- Audio channels: float to pattern (`5.1`, `7.1`)
- Release date priority (Radarr): `DigitalRelease` > `PhysicalRelease` > `InCinemas`
- Ratings priority: TMDb > IMDb
- Season naming (Sonarr): Season 0 = "Specials", others = "Season N"

### 4. SyncService (`service.go`)

Orchestrates sync operations with mutex-protected status tracking.

**Full sync** (`SyncLibrary`):
1. Acquire `sync.Mutex`, set `IsRunning = true`
2. Fetch all items from Servarr
3. Load existing items from local DB by external ID
4. For each item: skip if no files, update if exists, add if new
5. Sync related data (genres, files, collections/seasons)
6. Detect removed items (log only — no auto-delete)
7. Update status, release lock

**Single sync** (`SyncMovie` / `SyncSeries`):
1. Fetch single item by Servarr ID
2. Check local DB by external ID
3. Add or update accordingly
4. Sync dependent data

**Sync hierarchy**:

| Radarr (3 levels) | Sonarr (4 levels) |
|--------------------|-------------------|
| Movie | Series |
| MovieFile | Season |
| Collection | Episode |
| | EpisodeFile |

### 5. Jobs (`jobs.go`)

Two River job types per integration:

**Sync job** (`{Arr}SyncJobArgs`):
- Operations: `"full"` or `"single"`
- Queue: `QueueHigh`
- Timeout: 10 minutes

**Webhook job** (`{Arr}WebhookJobArgs`):
- Contains full `WebhookPayload`
- Queue: `QueueHigh`
- Timeout: 1 minute

Workers follow the standard River pattern (see [River Workers](RIVER_WORKERS.md)).

### 6. Webhook Handler (`webhook_handler.go`)

Routes webhook events to appropriate actions:

| Event | Action |
|-------|--------|
| Download/Import | `SyncMovie(id)` or `SyncSeries(id)` — **main trigger** |
| Rename | `SyncMovie(id)` / `SyncSeries(id)` — update file paths |
| Delete | Log only (no auto-delete from local DB) |
| File Delete | Sync to reflect removed files |
| Grab | No action (wait for download) |
| Health | Log warning only |
| Test | Acknowledge and log |

### 7. Module (`module.go`)

fx wiring with optional dependencies:

```go
Module = fx.Module("{arr}",
    fx.Provide(
        NewClientFromConfig,      // Config -> Client (nil if disabled)
        NewMapper,                // Stateless
        NewSyncServiceFromDeps,   // Returns nil if Client is nil
    ),
)
```

Client creation checks `config.{Arr}.Enabled` and validates `BaseURL` + `APIKey`. If disabled, provides `nil` — downstream services handle this gracefully.

### 8. API Handlers (`handler_{arr}.go`)

Three endpoint categories:

| Category | Endpoints | Auth |
|----------|-----------|------|
| Status | `GET /admin/integrations/{arr}/status` | Admin |
| Control | `POST /admin/integrations/{arr}/sync` | Admin |
| Config | `GET /admin/integrations/{arr}/quality-profiles`, `root-folders` | Admin |
| Webhook | `POST /webhooks/{arr}` | None (firewall-based) |

Sync endpoint: checks service availability (503 if nil), prevents duplicates (409 if running), queues River job if available, falls back to goroutine with timeout.

Webhook endpoint: converts ogen types to internal types, queues River job, returns 202 Accepted immediately.

---

## Configuration

```go
type {Arr}Config struct {
    Enabled      bool   // Master enable flag
    BaseURL      string // e.g., http://localhost:7878
    APIKey       string // From arr settings
    AutoSync     bool   // Enable scheduled syncs
    SyncInterval int    // Seconds between syncs
}
```

Loaded via koanf from `integrations.{arr}` config section.

---

## Key Differences: Radarr vs Sonarr

| Aspect | Radarr | Sonarr |
|--------|--------|--------|
| Primary external ID | TMDbID | TVDbID |
| Content unit | Movie (flat) | Series -> Season -> Episode (hierarchical) |
| Files | Per movie | Per episode |
| Related data | Collections | Seasons (with statistics) |
| Release dates | 3 sources (cinema, physical, digital) | `FirstAired`, `LastAired` |
| Sync complexity | 3 levels | 4 levels |
| Extra fields | — | `SeriesType` (standard/daily/anime), `FinaleType` |

---

## Adding a New Arr Integration

To add a new arr (e.g., Lidarr for music):

1. Create `internal/integration/lidarr/` with all 8 layers
2. Use Radarr as the template — copy structure, adapt types
3. Add `LidarrConfig` to `internal/config/config.go`
4. Add API handlers in `internal/api/handler_lidarr.go`
5. Register module in `internal/app/module.go`
6. Add OpenAPI endpoints in `api/openapi/`

Planned integrations (no code exists yet):
- **Lidarr** — music (for music module)
- **Whisparr** — adult content (for QAR module)
- **Readarr/Chaptarr** — books/comics (for book module)

---

## Related Documentation

- [Radarr Integration](../integrations/servarr/RADARR.md) — Radarr-specific details
- [Sonarr Integration](../integrations/servarr/SONARR.md) — Sonarr-specific details
- [Webhook Patterns](WEBHOOKS.md) — Webhook handling pattern
- [River Workers](RIVER_WORKERS.md) — Background job processing
- [Metadata Enrichment](METADATA.md) — How arr data feeds into metadata
- [HTTP Client](HTTP_CLIENT.md) — Client patterns (arr clients use same library)
