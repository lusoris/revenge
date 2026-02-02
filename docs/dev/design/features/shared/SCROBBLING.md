## Table of Contents

- [Revenge - External Scrobbling & Sync](#revenge-external-scrobbling-sync)
  - [Status](#status)
  - [Architecture](#architecture)
    - [Database Schema](#database-schema)
    - [Module Structure](#module-structure)
    - [Component Interaction](#component-interaction)
  - [Implementation](#implementation)
    - [Key Interfaces](#key-interfaces)
    - [Dependencies](#dependencies)
  - [Configuration](#configuration)
    - [Environment Variables](#environment-variables)
    - [Config Keys](#config-keys)
  - [API Endpoints](#api-endpoints)
    - [Content Management](#content-management)
  - [Related Documentation](#related-documentation)
    - [Design Documents](#design-documents)
    - [External Sources](#external-sources)

# Revenge - External Scrobbling & Sync


**Created**: 2026-01-31
**Status**: ✅ Complete
**Category**: feature


> Content module for 

> Sync playback data to external services like Trakt, Last.fm, ListenBrainz, etc.

---


## Status

| Dimension | Status | Notes |
|-----------|--------|-------|
| Design | ✅ | - |
| Sources | ✅ | - |
| Instructions | ✅ | - |
| Code | 🔴 | - |
| Linting | 🔴 | - |
| Unit Testing | 🔴 | - |
| Integration Testing | 🔴 | - |

**Overall**: ✅ Complete



---


## Architecture

```mermaid
flowchart TD
    node1([Client<br/>[Web/App]])
    node2[[API Handler<br/>[ogen]]]
    node3[[Service<br/>[Logic]]]
    node4["Repository<br/>[sqlc]"]
    node5[[Metadata<br/>Service]]
    node6[(Cache<br/>[otter])]
    node7[(PostgreSQL<br/>[pgx])]
    node8([External<br/>APIs])
    node1 --> node2
    node2 --> node3
    node4 --> node5
    node5 --> node6
    node7 --> node8
    node3 --> node4
    node6 --> node7
```

### Database Schema

**Schema**: `public`

<!-- Schema diagram -->

### Module Structure

```
internal/content/revenge___external_scrobbling_&_sync/
├── module.go              # fx module definition
├── repository.go          # Database operations
├── service.go             # Business logic
├── handler.go             # HTTP handlers (ogen)
├── types.go               # Domain types
└── revenge___external_scrobbling_&_sync_test.go
```

### Component Interaction

<!-- Component interaction diagram -->
## Implementation

### Key Interfaces

```go
type ScrobbleService interface {
  ConnectService(ctx context.Context, userID uuid.UUID, service string, oauthCode string) (*ScrobbleConnection, error)
  DisconnectService(ctx context.Context, connectionID uuid.UUID) error
  ListConnections(ctx context.Context, userID uuid.UUID) ([]ScrobbleConnection, error)
  UpdateConnection(ctx context.Context, connectionID uuid.UUID, settings ConnectionSettings) error

  QueueScrobble(ctx context.Context, connectionID uuid.UUID, scrobble ScrobbleRequest) error
  ProcessQueue(ctx context.Context) (int, error)

  ImportHistory(ctx context.Context, connectionID uuid.UUID) (int, error)
  GetHistory(ctx context.Context, connectionID uuid.UUID, limit int) ([]ScrobbleHistoryItem, error)
}

type ScrobbleConnection struct {
  ID                      uuid.UUID  `db:"id" json:"id"`
  UserID                  uuid.UUID  `db:"user_id" json:"user_id"`
  Service                 string     `db:"service" json:"service"`
  Enabled                 bool       `db:"enabled" json:"enabled"`
  ScrobbleThresholdPercent int       `db:"scrobble_threshold_percent" json:"scrobble_threshold_percent"`
  SyncWatchStatus         bool       `db:"sync_watch_status" json:"sync_watch_status"`
  SyncRatings             bool       `db:"sync_ratings" json:"sync_ratings"`
  LastSyncAt              *time.Time `db:"last_sync_at" json:"last_sync_at,omitempty"`
  LastSyncStatus          string     `db:"last_sync_status" json:"last_sync_status"`
}

type ScrobbleRequest struct {
  ContentType    string                 `json:"content_type"`
  ContentID      uuid.UUID              `json:"content_id"`
  Action         string                 `json:"action"`
  WatchedAt      time.Time              `json:"watched_at"`
  ProgressPercent int                   `json:"progress_percent,omitempty"`
  Rating         *float64               `json:"rating,omitempty"`
  ExternalIDs    map[string]interface{} `json:"external_ids"`
}

type ScrobbleClient interface {
  Scrobble(ctx context.Context, req ScrobbleRequest) error
  ImportHistory(ctx context.Context, since time.Time) ([]ExternalHistoryItem, error)
  GetExternalIDs(ctx context.Context, contentType string, contentID uuid.UUID) (map[string]string, error)
}

type TraktClient interface {
  ScrobbleClient
  OAuth(code string) (*OAuthTokens, error)
  RefreshToken(refreshToken string) (*OAuthTokens, error)
}
```


### Dependencies
**Go Packages**:
- `github.com/google/uuid`
- `github.com/jackc/pgx/v5`
- `github.com/riverqueue/river` - Background job queue for scrobble processing
- `golang.org/x/oauth2` - OAuth2 client
- `go.uber.org/fx`

**External APIs**:
- Trakt.tv API v2
- Last.fm API
- ListenBrainz API
- Simkl API
- AniList GraphQL API

## Configuration

### Environment Variables

```bash
TRAKT_CLIENT_ID=your_client_id
TRAKT_CLIENT_SECRET=your_client_secret
TRAKT_REDIRECT_URI=http://localhost:8080/api/v1/scrobble/trakt/callback

LASTFM_API_KEY=your_api_key
LASTFM_API_SECRET=your_secret

LISTENBRAINZ_TOKEN=your_token

SIMKL_CLIENT_ID=your_client_id
SIMKL_CLIENT_SECRET=your_client_secret

ANILIST_CLIENT_ID=your_client_id
ANILIST_CLIENT_SECRET=your_client_secret
```


### Config Keys
```yaml
scrobble:
  enabled: true
  queue_process_interval: 5m
  max_retries: 3
  retry_backoff: 5m

  services:
    trakt:
      enabled: true
      client_id: ${TRAKT_CLIENT_ID}
      client_secret: ${TRAKT_CLIENT_SECRET}
    lastfm:
      enabled: true
      api_key: ${LASTFM_API_KEY}
      api_secret: ${LASTFM_API_SECRET}
    listenbrainz:
      enabled: true
    simkl:
      enabled: true
      client_id: ${SIMKL_CLIENT_ID}
    anilist:
      enabled: true
      client_id: ${ANILIST_CLIENT_ID}
```

## API Endpoints

### Content Management
```
GET  /api/v1/scrobble/connections              # List user's connections
POST /api/v1/scrobble/connect/{service}        # Initiate OAuth flow
GET  /api/v1/scrobble/{service}/callback       # OAuth callback
DELETE /api/v1/scrobble/connections/{id}       # Disconnect service
PUT  /api/v1/scrobble/connections/{id}         # Update connection settings

POST /api/v1/scrobble/queue                    # Manually queue scrobble
GET  /api/v1/scrobble/history                  # Get scrobble history

POST /api/v1/scrobble/import/{service}         # Import history from service
GET  /api/v1/scrobble/imported                 # List imported items
```

## Related Documentation
### Design Documents
- [01_ARCHITECTURE](../../architecture/01_ARCHITECTURE.md)
- [02_DESIGN_PRINCIPLES](../../architecture/02_DESIGN_PRINCIPLES.md)
- [03_METADATA_SYSTEM](../../architecture/03_METADATA_SYSTEM.md)

### External Sources
- [Uber fx](../../../sources/tooling/fx.md) - Auto-resolved from fx
- [Last.fm API](../../../sources/apis/lastfm.md) - Auto-resolved from lastfm-api
- [ogen OpenAPI Generator](../../../sources/tooling/ogen.md) - Auto-resolved from ogen
- [River Job Queue](../../../sources/tooling/river.md) - Auto-resolved from river
- [sqlc](../../../sources/database/sqlc.md) - Auto-resolved from sqlc
- [sqlc Configuration](../../../sources/database/sqlc-config.md) - Auto-resolved from sqlc-config

