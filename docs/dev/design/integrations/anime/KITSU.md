## Table of Contents

- [Kitsu](#kitsu)
  - [Status](#status)
  - [Architecture](#architecture)
    - [Integration Structure](#integration-structure)
    - [Data Flow](#data-flow)
    - [Provides](#provides)
  - [Implementation](#implementation)
    - [Key Interfaces](#key-interfaces)
    - [Dependencies](#dependencies)
  - [Configuration](#configuration)
    - [Environment Variables](#environment-variables)
- [Kitsu OAuth configuration (per-user)](#kitsu-oauth-configuration-per-user)
- [Users configure via web UI, not environment variables](#users-configure-via-web-ui-not-environment-variables)
- [Global settings](#global-settings)
- [Sync settings](#sync-settings)
    - [Config Keys](#config-keys)
  - [API Endpoints](#api-endpoints)
- [OAuth (password grant - user provides Kitsu credentials)](#oauth-password-grant-user-provides-kitsu-credentials)
- [Status](#status)
- [Operations](#operations)
  - [Related Documentation](#related-documentation)
    - [Design Documents](#design-documents)
    - [External Sources](#external-sources)

# Kitsu


**Created**: 2026-01-31
**Status**: ✅ Complete
**Category**: integration


> Integration with Kitsu

> Modern anime tracking platform with social features
**API Base URL**: `https://kitsu.io/api/edge`
**Authentication**: oauth

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
    node1([Revenge<br/>Web Player])
    node2["Anime Item<br/>(TV Show)<br/>River Queue"]
    node3[[Scrobbling<br/>Service]]
    node4["River Queue<br/>(Background)"]
    node5["Kitsu<br/>Scrobbler"]
    node6[(Database<br/>(history))]
    node7[[Kitsu API<br/>(REST)]]
    node8["Kitsu<br/>User Library"]
    node1 --> node2
    node3 --> node4
    node6 --> node7
    node2 --> node3
    node4 --> node5
    node5 --> node6
    node7 --> node8
```

### Integration Structure

```
internal/integration/kitsu/
├── client.go              # API client
├── types.go               # Response types
├── mapper.go              # Map external → internal types
├── cache.go               # Response caching
└── client_test.go         # Tests
```

### Data Flow

<!-- Data flow diagram -->

### Provides
<!-- Data provided by integration -->
## Implementation

### Key Interfaces

```go
// KitsuClient manages Kitsu JSON:API
type KitsuClient interface {
    // Search for anime
    SearchAnime(ctx context.Context, query string) ([]AnimeResult, error)

    // Get anime by ID
    GetAnime(ctx context.Context, kitsuID string) (*AnimeDetails, error)

    // Get user's library
    GetLibrary(ctx context.Context, userID string, status LibraryStatus) (*LibraryEntryCollection, error)

    // Update library entry progress
    UpdateProgress(ctx context.Context, entryID string, progress int, status LibraryStatus) error

    // Add anime to library
    AddToLibrary(ctx context.Context, animeID string, status LibraryStatus) (*LibraryEntry, error)

    // Remove from library
    RemoveFromLibrary(ctx context.Context, entryID string) error

    // Get current user
    GetCurrentUser(ctx context.Context) (*User, error)
}

// KitsuScrobbler implements scrobbling to Kitsu
type KitsuScrobbler struct {
    client *KitsuClient
    config *KitsuConfig
}

type KitsuConfig struct {
    Email        string        // User email for OAuth
    Password     string        // User password for OAuth
    ClientID     string        // OAuth client ID
    ClientSecret string        // OAuth client secret
    Enabled      bool
    SyncInterval time.Duration
}

type JSONAPIResource struct {
    Type          string                    `json:"type"`
    ID            string                    `json:"id"`
    Attributes    map[string]interface{}    `json:"attributes,omitempty"`
    Relationships map[string]JSONAPIRelationship `json:"relationships,omitempty"`
    Links         map[string]string         `json:"links,omitempty"`
}

type JSONAPIRelationship struct {
    Data  interface{}       `json:"data,omitempty"`   // Single or array of {type, id}
    Links map[string]string `json:"links,omitempty"`
}

type JSONAPIResponse struct {
    Data     interface{}              `json:"data"`     // Single resource or array
    Included []JSONAPIResource        `json:"included,omitempty"`
    Links    map[string]string        `json:"links,omitempty"`
    Meta     map[string]interface{}   `json:"meta,omitempty"`
}

type AnimeResult struct {
    ID            string
    Slug          string
    Titles        Titles
    Synopsis      string
    PosterImage   PosterImage
    CoverImage    CoverImage
    EpisodeCount  int
    EpisodeLength int  // Minutes
    Status        AnimeStatus
    AgeRating     string
    AverageRating string  // "75.32"
    UserCount     int
    StartDate     string
    EndDate       string
}

type AnimeDetails struct {
    ID            string
    Slug          string
    Titles        Titles
    Synopsis      string
    Description   string
    PosterImage   PosterImage
    CoverImage    CoverImage
    EpisodeCount  int
    EpisodeLength int
    Status        AnimeStatus
    Subtype       AnimeSubtype  // TV, movie, ONA, OVA, special, music
    AgeRating     string        // G, PG, R, R18
    AverageRating string
    UserCount     int
    FavoritesCount int
    StartDate     string
    EndDate       string
    NSFW          bool
    Categories    []Category    // From included resources
    MediaRelationships []MediaRelationship
}

type Titles struct {
    En      string `json:"en"`
    EnJp    string `json:"en_jp"`  // Romaji
    JaJp    string `json:"ja_jp"`  // Native
}

type PosterImage struct {
    Tiny     string
    Small    string
    Medium   string
    Large    string
    Original string
}

type CoverImage struct {
    Tiny     string
    Small    string
    Large    string
    Original string
}

type AnimeStatus string

const (
    StatusCurrent   AnimeStatus = "current"    // Airing
    StatusFinished  AnimeStatus = "finished"   // Completed
    StatusTBA       AnimeStatus = "tba"        // Announced
    StatusUnreleased AnimeStatus = "unreleased"
    StatusUpcoming  AnimeStatus = "upcoming"
)

type AnimeSubtype string

const (
    SubtypeTV      AnimeSubtype = "TV"
    SubtypeMovie   AnimeSubtype = "movie"
    SubtypeONA     AnimeSubtype = "ONA"
    SubtypeOVA     AnimeSubtype = "OVA"
    SubtypeSpecial AnimeSubtype = "special"
    SubtypeMusic   AnimeSubtype = "music"
)

type LibraryStatus string

const (
    LibStatusCurrent   LibraryStatus = "current"    // Watching
    LibStatusPlanned   LibraryStatus = "planned"    // Plan to watch
    LibStatusCompleted LibraryStatus = "completed"  // Finished
    LibStatusOnHold    LibraryStatus = "on_hold"    // Paused
    LibStatusDropped   LibraryStatus = "dropped"    // Abandoned
)

type LibraryEntry struct {
    ID           string
    AnimeID      string
    Status       LibraryStatus
    Progress     int           // Episodes watched
    Rating       float64       // 0-20 (Kitsu uses 20-point scale)
    Notes        string
    Private      bool
    ProgressedAt time.Time
    StartedAt    string        // YYYY-MM-DD
    FinishedAt   string        // YYYY-MM-DD
    UpdatedAt    time.Time
    Anime        *AnimeDetails // From included resources
}

type LibraryEntryCollection struct {
    Entries []LibraryEntry
    Links   map[string]string  // Pagination links
    Meta    PaginationMeta
}

type PaginationMeta struct {
    Count int `json:"count"`
}

type User struct {
    ID       string
    Name     string
    Slug     string
    Avatar   Image
    CoverImage Image
    About    string
    Location string
    Waifus   int
    LifeSpentOnAnime int  // Minutes
}

type Image struct {
    Tiny     string
    Small    string
    Medium   string
    Large    string
    Original string
}

type Category struct {
    ID          string
    Title       string
    Description string
    Slug        string
}

type OAuthTokenResponse struct {
    AccessToken  string `json:"access_token"`
    TokenType    string `json:"token_type"`
    ExpiresIn    int    `json:"expires_in"`
    RefreshToken string `json:"refresh_token"`
    Scope        string `json:"scope"`
    CreatedAt    int64  `json:"created_at"`
}
```


### Dependencies
**Go Packages**:
- `github.com/riverqueue/river` - Background job queue
- `golang.org/x/oauth2` - OAuth 2.0 client
- `encoding/json` - JSON encoding
- `net/http` - HTTP client

**External Services**:
- Kitsu API (https://kitsu.io)

## Configuration

### Environment Variables

```bash
# Kitsu OAuth configuration (per-user)
# Users configure via web UI, not environment variables

# Global settings
REVENGE_KITSU_ENABLED=true
REVENGE_KITSU_CLIENT_ID=dd031b32d2f56c990b1425efe6c42ad847e7fe3ab46bf1299f05ecd856bdb7dd
REVENGE_KITSU_CLIENT_SECRET=54d7307928f63414defd96399fc31ba847961ceaecef3a5fd93144e960c0e151

# Sync settings
REVENGE_KITSU_SYNC_INTERVAL=24h
```


### Config Keys
```yaml
scrobbling:
  kitsu:
    enabled: true
    client_id: dd031b32d2f56c990b1425efe6c42ad847e7fe3ab46bf1299f05ecd856bdb7dd
    client_secret: 54d7307928f63414defd96399fc31ba847961ceaecef3a5fd93144e960c0e151
    sync_interval: 24h
    auto_sync: true
    scrobble_threshold: 0.9
    metadata_enabled: true
    metadata_priority: 15     # Lower priority than AniList
```

## API Endpoints
**Revenge API Endpoints**:

```
# OAuth (password grant - user provides Kitsu credentials)
POST /api/v1/scrobbling/kitsu/connect
POST /api/v1/scrobbling/kitsu/disconnect

# Status
GET  /api/v1/scrobbling/kitsu/status

# Operations
POST /api/v1/scrobbling/kitsu/import
POST /api/v1/scrobbling/kitsu/sync
```

**Example - Connect**:
```json
POST /api/v1/scrobbling/kitsu/connect
{
  "email": "user@example.com",
  "password": "kitsu-password"
}

Response:
{
  "connected": true,
  "user": {
    "id": "123456",
    "name": "AnimeWatcher",
    "slug": "animewatcher"
  }
}
```

## Related Documentation
### Design Documents
- [01_ARCHITECTURE](../../architecture/01_ARCHITECTURE.md)
- [02_DESIGN_PRINCIPLES](../../architecture/02_DESIGN_PRINCIPLES.md)
- [03_METADATA_SYSTEM](../../architecture/03_METADATA_SYSTEM.md)

### External Sources
- [River Job Queue](../../../sources/tooling/river.md) - Auto-resolved from river

