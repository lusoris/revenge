

---
sources:
  - name: MyAnimeList API
    url: ../../../sources/apis/myanimelist.md
    note: Auto-resolved from myanimelist
  - name: River Job Queue
    url: ../../../sources/tooling/river.md
    note: Auto-resolved from river
design_refs:
  - title: 01_ARCHITECTURE
    path: ../../architecture/01_ARCHITECTURE.md
  - title: 02_DESIGN_PRINCIPLES
    path: ../../architecture/02_DESIGN_PRINCIPLES.md
  - title: 03_METADATA_SYSTEM
    path: ../../architecture/03_METADATA_SYSTEM.md
---

## Table of Contents

- [MyAnimeList (MAL)](#myanimelist-mal)
  - [Status](#status)
  - [Architecture](#architecture)
    - [Integration Structure](#integration-structure)
    - [Data Flow](#data-flow)
    - [Provides](#provides)
  - [Implementation](#implementation)
    - [File Structure](#file-structure)
    - [Key Interfaces](#key-interfaces)
    - [Dependencies](#dependencies)
  - [Configuration](#configuration)
    - [Environment Variables](#environment-variables)
- [MAL OAuth configuration](#mal-oauth-configuration)
- [Sync settings](#sync-settings)
    - [Config Keys](#config-keys)
  - [API Endpoints](#api-endpoints)
- [OAuth flow with PKCE](#oauth-flow-with-pkce)
- [Status](#status)
- [Operations](#operations)
  - [Testing Strategy](#testing-strategy)
    - [Unit Tests](#unit-tests)
    - [Integration Tests](#integration-tests)
    - [Test Coverage](#test-coverage)
  - [Related Documentation](#related-documentation)
    - [Design Documents](#design-documents)
    - [External Sources](#external-sources)


# MyAnimeList (MAL)


**Created**: 2026-01-31
**Status**: ✅ Complete
**Category**: integration


> Integration with MyAnimeList (MAL)

> Legacy anime tracking platform with extensive database
**API Base URL**: `https://api.myanimelist.net/v2`
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

### Integration Structure

```
internal/integration/myanimelist_mal/
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

### File Structure

<!-- File structure -->

### Key Interfaces

```go
// MALClient manages MyAnimeList API v2
type MALClient interface {
    // Search for anime
    SearchAnime(ctx context.Context, query string, limit int) ([]AnimeResult, error)

    // Get anime details
    GetAnime(ctx context.Context, malID int, fields []string) (*AnimeDetails, error)

    // Get user's anime list
    GetUserList(ctx context.Context, status ListStatus, sort string, limit int, offset int) (*UserAnimeList, error)

    // Update watch status
    UpdateStatus(ctx context.Context, malID int, status ListStatus, numWatchedEpisodes int) (*ListStatus, error)

    // Delete from list
    DeleteFromList(ctx context.Context, malID int) error

    // Get current user
    GetCurrentUser(ctx context.Context) (*User, error)
}

// MALScrobbler implements scrobbling to MAL
type MALScrobbler struct {
    client *MALClient
    config *MALConfig
}

type MALConfig struct {
    ClientID     string
    ClientSecret string
    RedirectURL  string
    Enabled      bool
    SyncInterval time.Duration
}

type AnimeResult struct {
    Node struct {
        ID           int           `json:"id"`
        Title        string        `json:"title"`
        MainPicture  Picture       `json:"main_picture"`
        AlternativeTitles AlternativeTitles `json:"alternative_titles,omitempty"`
        StartDate    string        `json:"start_date,omitempty"`
        EndDate      string        `json:"end_date,omitempty"`
        Synopsis     string        `json:"synopsis,omitempty"`
        Mean         float64       `json:"mean,omitempty"`         // Average score
        Rank         int           `json:"rank,omitempty"`
        Popularity   int           `json:"popularity,omitempty"`
        NumListUsers int           `json:"num_list_users,omitempty"`
        NumEpisodes  int           `json:"num_episodes,omitempty"`
        MediaType    string        `json:"media_type,omitempty"`   // tv, movie, ova, etc.
        Status       AnimeStatus   `json:"status,omitempty"`
        Genres       []Genre       `json:"genres,omitempty"`
        Studios      []Studio      `json:"studios,omitempty"`
    } `json:"node"`
}

type AnimeDetails struct {
    ID                int                  `json:"id"`
    Title             string               `json:"title"`
    MainPicture       Picture              `json:"main_picture"`
    AlternativeTitles AlternativeTitles    `json:"alternative_titles"`
    StartDate         string               `json:"start_date"`
    EndDate           string               `json:"end_date"`
    Synopsis          string               `json:"synopsis"`
    Mean              float64              `json:"mean"`
    Rank              int                  `json:"rank"`
    Popularity        int                  `json:"popularity"`
    NumListUsers      int                  `json:"num_list_users"`
    NumScoringUsers   int                  `json:"num_scoring_users"`
    NSFW              string               `json:"nsfw"`  // white, gray, black
    Genres            []Genre              `json:"genres"`
    CreatedAt         time.Time            `json:"created_at"`
    UpdatedAt         time.Time            `json:"updated_at"`
    MediaType         string               `json:"media_type"`
    Status            AnimeStatus          `json:"status"`
    NumEpisodes       int                  `json:"num_episodes"`
    StartSeason       Season               `json:"start_season,omitempty"`
    Broadcast         Broadcast            `json:"broadcast,omitempty"`
    Source            string               `json:"source,omitempty"`  // manga, light_novel, etc.
    AverageEpisodeDuration int             `json:"average_episode_duration,omitempty"`
    Rating            string               `json:"rating,omitempty"`  // g, pg, pg_13, r, r+, rx
    Studios           []Studio             `json:"studios,omitempty"`
    Pictures          []Picture            `json:"pictures,omitempty"`
    Background        string               `json:"background,omitempty"`
    RelatedAnime      []RelatedAnime       `json:"related_anime,omitempty"`
    Recommendations   []Recommendation     `json:"recommendations,omitempty"`
    Statistics        Statistics           `json:"statistics,omitempty"`
}

type Picture struct {
    Medium string `json:"medium"`
    Large  string `json:"large"`
}

type AlternativeTitles struct {
    Synonyms []string `json:"synonyms"`
    En       string   `json:"en"`
    Ja       string   `json:"ja"`
}

type Genre struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
}

type Studio struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
}

type Season struct {
    Year   int    `json:"year"`
    Season string `json:"season"`  // winter, spring, summer, fall
}

type Broadcast struct {
    DayOfTheWeek string `json:"day_of_the_week"`  // monday, tuesday, etc.
    StartTime    string `json:"start_time"`       // HH:MM
}

type RelatedAnime struct {
    Node                AnimeResult `json:"node"`
    RelationType        string      `json:"relation_type"`  // sequel, prequel, etc.
    RelationTypeFormatted string    `json:"relation_type_formatted"`
}

type Recommendation struct {
    Node             AnimeResult `json:"node"`
    NumRecommendations int      `json:"num_recommendations"`
}

type Statistics struct {
    NumListUsers int                `json:"num_list_users"`
    Status       StatusStatistics    `json:"status"`
}

type StatusStatistics struct {
    Watching    int `json:"watching"`
    Completed   int `json:"completed"`
    OnHold      int `json:"on_hold"`
    Dropped     int `json:"dropped"`
    PlanToWatch int `json:"plan_to_watch"`
}

type AnimeStatus string

const (
    StatusFinishedAiring  AnimeStatus = "finished_airing"
    StatusCurrentlyAiring AnimeStatus = "currently_airing"
    StatusNotYetAired     AnimeStatus = "not_yet_aired"
)

type ListStatus string

const (
    ListWatching    ListStatus = "watching"
    ListCompleted   ListStatus = "completed"
    ListOnHold      ListStatus = "on_hold"
    ListDropped     ListStatus = "dropped"
    ListPlanToWatch ListStatus = "plan_to_watch"
)

type UserAnimeListEntry struct {
    Node       AnimeResult     `json:"node"`
    ListStatus UserListStatus  `json:"list_status"`
}

type UserListStatus struct {
    Status             ListStatus `json:"status"`
    Score              int        `json:"score"`              // 0-10
    NumEpisodesWatched int        `json:"num_episodes_watched"`
    IsRewatching       bool       `json:"is_rewatching"`
    StartDate          string     `json:"start_date"`         // YYYY-MM-DD
    FinishDate         string     `json:"finish_date"`        // YYYY-MM-DD
    Priority           int        `json:"priority"`           // 0-2 (low, medium, high)
    NumTimesRewatched  int        `json:"num_times_rewatched"`
    RewatchValue       int        `json:"rewatch_value"`      // 0-5
    Tags               []string   `json:"tags"`
    Comments           string     `json:"comments"`
    UpdatedAt          time.Time  `json:"updated_at"`
}

type UserAnimeList struct {
    Data   []UserAnimeListEntry `json:"data"`
    Paging Paging               `json:"paging"`
}

type Paging struct {
    Previous string `json:"previous,omitempty"`
    Next     string `json:"next,omitempty"`
}

type User struct {
    ID       int       `json:"id"`
    Name     string    `json:"name"`
    Picture  string    `json:"picture"`
    Gender   string    `json:"gender"`
    Birthday string    `json:"birthday"`
    Location string    `json:"location"`
    JoinedAt time.Time `json:"joined_at"`
}

type PKCEChallenge struct {
    CodeVerifier  string
    CodeChallenge string
    State         string
}
```


### Dependencies

**Go Packages**:
- `github.com/riverqueue/river` - Background job queue
- `golang.org/x/oauth2` - OAuth 2.0 client
- `crypto/sha256` - PKCE code challenge
- `encoding/base64` - PKCE encoding
- `crypto/rand` - Random code generation

**External Services**:
- MyAnimeList API v2 (https://myanimelist.net)






## Configuration
### Environment Variables

```bash
# MAL OAuth configuration
REVENGE_MAL_ENABLED=true
REVENGE_MAL_CLIENT_ID=your-client-id
REVENGE_MAL_CLIENT_SECRET=your-client-secret
REVENGE_MAL_REDIRECT_URL=https://revenge.local/api/v1/scrobbling/mal/callback

# Sync settings
REVENGE_MAL_SYNC_INTERVAL=24h
```


### Config Keys

```yaml
scrobbling:
  myanimelist:
    enabled: true
    client_id: ${REVENGE_MAL_CLIENT_ID}
    client_secret: ${REVENGE_MAL_CLIENT_SECRET}
    redirect_url: https://revenge.local/api/v1/scrobbling/mal/callback
    sync_interval: 24h
    auto_sync: true
    scrobble_threshold: 0.9
    metadata_enabled: true
    metadata_priority: 20       # Lower priority than AniList/Kitsu
```



## API Endpoints
**Revenge API Endpoints**:

```
# OAuth flow with PKCE
GET  /api/v1/scrobbling/mal/authorize
GET  /api/v1/scrobbling/mal/callback

# Status
GET  /api/v1/scrobbling/mal/status
POST /api/v1/scrobbling/mal/disconnect

# Operations
POST /api/v1/scrobbling/mal/import
POST /api/v1/scrobbling/mal/sync
```



## Testing Strategy

### Unit Tests

<!-- Unit test strategy -->

### Integration Tests

<!-- Integration test strategy -->

### Test Coverage

Target: **80% minimum**







## Related Documentation
### Design Documents
- [01_ARCHITECTURE](../../architecture/01_ARCHITECTURE.md)
- [02_DESIGN_PRINCIPLES](../../architecture/02_DESIGN_PRINCIPLES.md)
- [03_METADATA_SYSTEM](../../architecture/03_METADATA_SYSTEM.md)

### External Sources
- [MyAnimeList API](../../../sources/apis/myanimelist.md) - Auto-resolved from myanimelist
- [River Job Queue](../../../sources/tooling/river.md) - Auto-resolved from river

