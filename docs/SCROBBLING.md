# Revenge - External Scrobbling & Sync

> Sync playback data to external services like Trakt, Last.fm, ListenBrainz, etc.

## Overview

Revenge supports scrobbling (reporting playback) to various external services:

| Service | Content Types | Features |
|---------|---------------|----------|
| **Trakt** | Movies, TV Shows | Watch history, ratings, watchlist, collections |
| **Last.fm** | Music | Scrobbles, now playing, loved tracks |
| **ListenBrainz** | Music | Scrobbles, listening history |
| **Goodreads** | Books | Read status, ratings, reviews |
| **Simkl** | Movies, TV, Anime | Watch history, ratings |
| **AniList** | Anime | Watch status, ratings |
| **MyAnimeList** | Anime | Watch status, ratings |
| **Letterboxd** | Movies | Watch history, ratings (via Trakt) |

---

## Architecture

```
┌──────────────────────────────────────────────────────────────────────────┐
│                         Scrobbling Flow                                   │
└──────────────────────────────────────────────────────────────────────────┘

┌─────────┐    ┌───────────┐    ┌────────────────┐    ┌──────────────────┐
│ Client  │ ──→│  Revenge  │ ──→│ Scrobble Queue │ ──→│ External Service │
│ Playback│    │  Server   │    │   (River)      │    │ (Trakt, Last.fm) │
└─────────┘    └───────────┘    └────────────────┘    └──────────────────┘
     │              │                   │                      │
     │  Progress    │                   │                      │
     │  Update      │                   │                      │
     │─────────────→│                   │                      │
     │              │                   │                      │
     │              │  Check threshold  │                      │
     │              │  (50% music,      │                      │
     │              │   80% video)      │                      │
     │              │──────────────────→│                      │
     │              │                   │                      │
     │              │                   │  Async scrobble      │
     │              │                   │─────────────────────→│
     │              │                   │                      │
     │              │                   │  Rate limit aware    │
     │              │                   │  Retry on failure    │
```

---

## Trakt Integration

### OAuth Setup

```go
type TraktConfig struct {
    ClientID     string `yaml:"client_id"`
    ClientSecret string `yaml:"client_secret"`
    RedirectURI  string `yaml:"redirect_uri"`
}

type TraktClient struct {
    config     TraktConfig
    httpClient *http.Client
    baseURL    string
}

// OAuth flow
func (c *TraktClient) GetAuthURL(state string) string {
    return fmt.Sprintf(
        "https://trakt.tv/oauth/authorize?response_type=code&client_id=%s&redirect_uri=%s&state=%s",
        c.config.ClientID,
        url.QueryEscape(c.config.RedirectURI),
        state,
    )
}

func (c *TraktClient) ExchangeCode(ctx context.Context, code string) (*TraktToken, error) {
    payload := map[string]string{
        "code":          code,
        "client_id":     c.config.ClientID,
        "client_secret": c.config.ClientSecret,
        "redirect_uri":  c.config.RedirectURI,
        "grant_type":    "authorization_code",
    }

    resp, err := c.post(ctx, "/oauth/token", payload, nil)
    if err != nil {
        return nil, err
    }

    var token TraktToken
    if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
        return nil, err
    }

    return &token, nil
}
```

### Scrobble Actions

```go
type TraktScrobbleAction string

const (
    TraktScrobbleStart TraktScrobbleAction = "start"
    TraktScrobblePause TraktScrobbleAction = "pause"
    TraktScrobbleStop  TraktScrobbleAction = "stop"
)

type TraktScrobbleRequest struct {
    Movie   *TraktMovie   `json:"movie,omitempty"`
    Episode *TraktEpisode `json:"episode,omitempty"`
    Show    *TraktShow    `json:"show,omitempty"`
    Progress float64      `json:"progress"`
}

type TraktMovie struct {
    IDs TraktIDs `json:"ids"`
}

type TraktEpisode struct {
    Season int      `json:"season"`
    Number int      `json:"number"`
    IDs    TraktIDs `json:"ids,omitempty"`
}

type TraktShow struct {
    IDs TraktIDs `json:"ids"`
}

type TraktIDs struct {
    Trakt  int    `json:"trakt,omitempty"`
    IMDB   string `json:"imdb,omitempty"`
    TMDB   int    `json:"tmdb,omitempty"`
    TVDB   int    `json:"tvdb,omitempty"`
    Slug   string `json:"slug,omitempty"`
}

func (c *TraktClient) ScrobbleMovie(ctx context.Context, token string, action TraktScrobbleAction, tmdbID int, progress float64) error {
    req := TraktScrobbleRequest{
        Movie:    &TraktMovie{IDs: TraktIDs{TMDB: tmdbID}},
        Progress: progress,
    }

    endpoint := fmt.Sprintf("/scrobble/%s", action)
    _, err := c.post(ctx, endpoint, req, &token)
    return err
}

func (c *TraktClient) ScrobbleEpisode(ctx context.Context, token string, action TraktScrobbleAction, tvdbID, season, episode int, progress float64) error {
    req := TraktScrobbleRequest{
        Episode:  &TraktEpisode{Season: season, Number: episode},
        Show:     &TraktShow{IDs: TraktIDs{TVDB: tvdbID}},
        Progress: progress,
    }

    endpoint := fmt.Sprintf("/scrobble/%s", action)
    _, err := c.post(ctx, endpoint, req, &token)
    return err
}
```

### Sync History

```go
// Sync watch history bidirectionally
func (c *TraktClient) SyncHistory(ctx context.Context, token string, since time.Time) (*TraktSyncResult, error) {
    // Get history from Trakt
    params := url.Values{}
    if !since.IsZero() {
        params.Set("start_at", since.Format(time.RFC3339))
    }

    resp, err := c.get(ctx, "/sync/history", params, &token)
    if err != nil {
        return nil, err
    }

    var history []TraktHistoryItem
    if err := json.NewDecoder(resp.Body).Decode(&history); err != nil {
        return nil, err
    }

    return &TraktSyncResult{Items: history}, nil
}

// Push local history to Trakt
func (c *TraktClient) AddToHistory(ctx context.Context, token string, items []TraktHistoryItem) error {
    payload := map[string][]TraktHistoryItem{
        "movies":   filterMovies(items),
        "episodes": filterEpisodes(items),
    }

    _, err := c.post(ctx, "/sync/history", payload, &token)
    return err
}
```

---

## Last.fm Integration

### Scrobbling

```go
type LastFMClient struct {
    apiKey    string
    apiSecret string
    baseURL   string
}

type LastFMScrobble struct {
    Artist    string    `json:"artist"`
    Track     string    `json:"track"`
    Album     string    `json:"album,omitempty"`
    Timestamp time.Time `json:"timestamp"`
    Duration  int       `json:"duration,omitempty"` // seconds
    MBID      string    `json:"mbid,omitempty"`     // MusicBrainz ID
}

func (c *LastFMClient) Scrobble(ctx context.Context, sessionKey string, scrobble LastFMScrobble) error {
    params := url.Values{
        "method":    {"track.scrobble"},
        "artist":    {scrobble.Artist},
        "track":     {scrobble.Track},
        "timestamp": {fmt.Sprintf("%d", scrobble.Timestamp.Unix())},
        "sk":        {sessionKey},
        "api_key":   {c.apiKey},
    }

    if scrobble.Album != "" {
        params.Set("album", scrobble.Album)
    }
    if scrobble.MBID != "" {
        params.Set("mbid", scrobble.MBID)
    }

    // Sign request
    params.Set("api_sig", c.signRequest(params))
    params.Set("format", "json")

    resp, err := http.PostForm(c.baseURL, params)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("last.fm scrobble failed: %d", resp.StatusCode)
    }

    return nil
}

// Update "Now Playing" (doesn't count as scrobble)
func (c *LastFMClient) UpdateNowPlaying(ctx context.Context, sessionKey string, track LastFMScrobble) error {
    params := url.Values{
        "method":  {"track.updateNowPlaying"},
        "artist":  {track.Artist},
        "track":   {track.Track},
        "sk":      {sessionKey},
        "api_key": {c.apiKey},
    }

    if track.Album != "" {
        params.Set("album", track.Album)
    }
    if track.Duration > 0 {
        params.Set("duration", fmt.Sprintf("%d", track.Duration))
    }

    params.Set("api_sig", c.signRequest(params))
    params.Set("format", "json")

    _, err := http.PostForm(c.baseURL, params)
    return err
}

func (c *LastFMClient) signRequest(params url.Values) string {
    // Sort parameters alphabetically
    keys := make([]string, 0, len(params))
    for k := range params {
        if k != "format" && k != "callback" {
            keys = append(keys, k)
        }
    }
    sort.Strings(keys)

    // Build signature string
    var sig strings.Builder
    for _, k := range keys {
        sig.WriteString(k)
        sig.WriteString(params.Get(k))
    }
    sig.WriteString(c.apiSecret)

    // MD5 hash
    hash := md5.Sum([]byte(sig.String()))
    return hex.EncodeToString(hash[:])
}
```

### Last.fm Rules

- Track must be longer than 30 seconds
- Must listen to at least 50% OR 4 minutes
- Scrobble timestamp should be when playback started
- Rate limit: max 5 scrobbles per second per user

---

## ListenBrainz Integration

```go
type ListenBrainzClient struct {
    baseURL string
}

type ListenBrainzSubmission struct {
    ListenType string              `json:"listen_type"` // "single", "playing_now", "import"
    Payload    []ListenBrainzTrack `json:"payload"`
}

type ListenBrainzTrack struct {
    ListenedAt    int64                  `json:"listened_at,omitempty"`
    TrackMetadata ListenBrainzMetadata   `json:"track_metadata"`
}

type ListenBrainzMetadata struct {
    ArtistName  string                 `json:"artist_name"`
    TrackName   string                 `json:"track_name"`
    ReleaseName string                 `json:"release_name,omitempty"`
    AdditionalInfo map[string]any     `json:"additional_info,omitempty"`
}

func (c *ListenBrainzClient) SubmitListen(ctx context.Context, token string, track ListenBrainzTrack) error {
    submission := ListenBrainzSubmission{
        ListenType: "single",
        Payload:    []ListenBrainzTrack{track},
    }

    body, _ := json.Marshal(submission)
    req, _ := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/1/submit-listens", bytes.NewReader(body))
    req.Header.Set("Authorization", "Token "+token)
    req.Header.Set("Content-Type", "application/json")

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("listenbrainz submit failed: %d", resp.StatusCode)
    }

    return nil
}

func (c *ListenBrainzClient) NowPlaying(ctx context.Context, token string, track ListenBrainzTrack) error {
    submission := ListenBrainzSubmission{
        ListenType: "playing_now",
        Payload:    []ListenBrainzTrack{track},
    }

    body, _ := json.Marshal(submission)
    req, _ := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/1/submit-listens", bytes.NewReader(body))
    req.Header.Set("Authorization", "Token "+token)
    req.Header.Set("Content-Type", "application/json")

    _, err := http.DefaultClient.Do(req)
    return err
}
```

---

## Scrobble Service

### Job Definitions

```go
// Video scrobble job (Trakt, Simkl)
type VideoScrobbleArgs struct {
    UserID     uuid.UUID           `json:"user_id"`
    ItemID     uuid.UUID           `json:"item_id"`
    ItemType   string              `json:"item_type"` // "movie", "episode"
    Action     TraktScrobbleAction `json:"action"`
    Progress   float64             `json:"progress"`
    WatchedAt  time.Time           `json:"watched_at"`

    // IDs for external services
    IMDBID     string              `json:"imdb_id,omitempty"`
    TMDBID     int                 `json:"tmdb_id,omitempty"`
    TVDBID     int                 `json:"tvdb_id,omitempty"`
    Season     int                 `json:"season,omitempty"`
    Episode    int                 `json:"episode,omitempty"`
}

func (VideoScrobbleArgs) Kind() string { return "scrobble.video" }

// Music scrobble job (Last.fm, ListenBrainz)
type MusicScrobbleArgs struct {
    UserID      uuid.UUID `json:"user_id"`
    TrackID     uuid.UUID `json:"track_id"`
    Artist      string    `json:"artist"`
    Track       string    `json:"track"`
    Album       string    `json:"album,omitempty"`
    Duration    int       `json:"duration"` // seconds
    MBID        string    `json:"mbid,omitempty"`
    ListenedAt  time.Time `json:"listened_at"`
    IsNowPlaying bool     `json:"is_now_playing"`
}

func (MusicScrobbleArgs) Kind() string { return "scrobble.music" }
```

### Scrobble Workers

```go
type VideoScrobbleWorker struct {
    river.WorkerDefaults[VideoScrobbleArgs]
    trakt   *TraktClient
    simkl   *SimklClient
    users   *UserRepository
    logger  *slog.Logger
}

func (w *VideoScrobbleWorker) Work(ctx context.Context, job *river.Job[VideoScrobbleArgs]) error {
    args := job.Args

    // Get user's connected services
    user, err := w.users.GetByID(ctx, args.UserID)
    if err != nil {
        return fmt.Errorf("get user: %w", err)
    }

    var errs []error

    // Trakt
    if token := user.GetServiceToken("trakt"); token != "" {
        var scrobbleErr error
        if args.ItemType == "movie" {
            scrobbleErr = w.trakt.ScrobbleMovie(ctx, token, args.Action, args.TMDBID, args.Progress)
        } else {
            scrobbleErr = w.trakt.ScrobbleEpisode(ctx, token, args.Action, args.TVDBID, args.Season, args.Episode, args.Progress)
        }
        if scrobbleErr != nil {
            errs = append(errs, fmt.Errorf("trakt: %w", scrobbleErr))
        }
    }

    // Simkl
    if token := user.GetServiceToken("simkl"); token != "" {
        // Similar implementation
    }

    if len(errs) > 0 {
        return errors.Join(errs...)
    }

    return nil
}

type MusicScrobbleWorker struct {
    river.WorkerDefaults[MusicScrobbleArgs]
    lastfm      *LastFMClient
    listenbrainz *ListenBrainzClient
    users       *UserRepository
    logger      *slog.Logger
}

func (w *MusicScrobbleWorker) Work(ctx context.Context, job *river.Job[MusicScrobbleArgs]) error {
    args := job.Args

    user, err := w.users.GetByID(ctx, args.UserID)
    if err != nil {
        return fmt.Errorf("get user: %w", err)
    }

    var errs []error

    // Last.fm
    if sessionKey := user.GetServiceToken("lastfm"); sessionKey != "" {
        scrobble := LastFMScrobble{
            Artist:    args.Artist,
            Track:     args.Track,
            Album:     args.Album,
            Timestamp: args.ListenedAt,
            Duration:  args.Duration,
            MBID:      args.MBID,
        }

        if args.IsNowPlaying {
            if err := w.lastfm.UpdateNowPlaying(ctx, sessionKey, scrobble); err != nil {
                errs = append(errs, fmt.Errorf("lastfm now playing: %w", err))
            }
        } else {
            if err := w.lastfm.Scrobble(ctx, sessionKey, scrobble); err != nil {
                errs = append(errs, fmt.Errorf("lastfm scrobble: %w", err))
            }
        }
    }

    // ListenBrainz
    if token := user.GetServiceToken("listenbrainz"); token != "" {
        track := ListenBrainzTrack{
            ListenedAt: args.ListenedAt.Unix(),
            TrackMetadata: ListenBrainzMetadata{
                ArtistName:  args.Artist,
                TrackName:   args.Track,
                ReleaseName: args.Album,
            },
        }

        if args.IsNowPlaying {
            if err := w.listenbrainz.NowPlaying(ctx, token, track); err != nil {
                errs = append(errs, fmt.Errorf("listenbrainz now playing: %w", err))
            }
        } else {
            if err := w.listenbrainz.SubmitListen(ctx, token, track); err != nil {
                errs = append(errs, fmt.Errorf("listenbrainz scrobble: %w", err))
            }
        }
    }

    if len(errs) > 0 {
        return errors.Join(errs...)
    }

    return nil
}
```

### Scrobble Trigger Service

```go
type ScrobbleService struct {
    jobs   *river.Client[pgx.Tx]
    logger *slog.Logger
}

// Called when progress is updated
func (s *ScrobbleService) OnProgressUpdate(ctx context.Context, userID uuid.UUID, update ProgressUpdate) {
    // Check if this is a scrobbleable event
    switch update.ItemType {
    case "movie", "episode":
        s.handleVideoProgress(ctx, userID, update)
    case "track":
        s.handleMusicProgress(ctx, userID, update)
    }
}

func (s *ScrobbleService) handleVideoProgress(ctx context.Context, userID uuid.UUID, update ProgressUpdate) {
    progress := float64(update.PositionMs) / float64(update.DurationMs) * 100

    // Trakt scrobble thresholds:
    // - start: when playback begins
    // - pause: when paused
    // - stop: at 80% or when stopped

    var action TraktScrobbleAction

    switch {
    case update.Event == "start":
        action = TraktScrobbleStart
    case update.Event == "pause":
        action = TraktScrobblePause
    case update.Event == "stop" || progress >= 80:
        action = TraktScrobbleStop
    default:
        return // No scrobble needed
    }

    // Queue scrobble job
    s.jobs.Insert(ctx, &VideoScrobbleArgs{
        UserID:   userID,
        ItemID:   update.ItemID,
        ItemType: update.ItemType,
        Action:   action,
        Progress: progress,
        // ... external IDs filled from metadata
    })
}

func (s *ScrobbleService) handleMusicProgress(ctx context.Context, userID uuid.UUID, update ProgressUpdate) {
    progress := float64(update.PositionMs) / float64(update.DurationMs) * 100

    // Last.fm rules:
    // - Track > 30 seconds
    // - Listened to 50% OR 4 minutes

    if update.DurationMs < 30000 {
        return // Too short
    }

    // Now Playing: at start
    if update.Event == "start" {
        s.jobs.Insert(ctx, &MusicScrobbleArgs{
            UserID:       userID,
            TrackID:      update.ItemID,
            IsNowPlaying: true,
            // ... track metadata
        })
        return
    }

    // Scrobble: at 50% or 4 minutes
    minProgress := min(50.0, float64(4*60*1000)/float64(update.DurationMs)*100)

    if progress >= minProgress && update.Event == "scrobble_point" {
        s.jobs.Insert(ctx, &MusicScrobbleArgs{
            UserID:       userID,
            TrackID:      update.ItemID,
            IsNowPlaying: false,
            ListenedAt:   time.Now().Add(-time.Duration(update.PositionMs) * time.Millisecond),
            // ... track metadata
        })
    }
}
```

---

## User Service Connections

### Database Schema

```sql
CREATE TABLE user_external_services (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    service         VARCHAR(50) NOT NULL,  -- 'trakt', 'lastfm', 'listenbrainz', etc.

    -- OAuth tokens
    access_token    TEXT,
    refresh_token   TEXT,
    token_expires   TIMESTAMPTZ,

    -- Service-specific
    username        VARCHAR(255),          -- Display username on service
    profile_url     TEXT,

    -- Settings
    enabled         BOOLEAN DEFAULT true,
    sync_history    BOOLEAN DEFAULT true,  -- Import history from service
    sync_ratings    BOOLEAN DEFAULT true,  -- Sync ratings
    sync_watchlist  BOOLEAN DEFAULT true,  -- Sync watchlist

    -- Timestamps
    connected_at    TIMESTAMPTZ DEFAULT NOW(),
    last_sync_at    TIMESTAMPTZ,

    UNIQUE (user_id, service)
);

CREATE INDEX idx_user_services_user ON user_external_services(user_id);
```

### Connection Management

```go
type ExternalServiceConnection struct {
    ID           uuid.UUID  `json:"id"`
    Service      string     `json:"service"`
    Username     string     `json:"username,omitempty"`
    ProfileURL   string     `json:"profile_url,omitempty"`
    Enabled      bool       `json:"enabled"`
    SyncHistory  bool       `json:"sync_history"`
    SyncRatings  bool       `json:"sync_ratings"`
    SyncWatchlist bool      `json:"sync_watchlist"`
    ConnectedAt  time.Time  `json:"connected_at"`
    LastSyncAt   *time.Time `json:"last_sync_at,omitempty"`
}

type ExternalServiceService struct {
    db     *pgxpool.Pool
    trakt  *TraktClient
    lastfm *LastFMClient
}

func (s *ExternalServiceService) Connect(ctx context.Context, userID uuid.UUID, service, code string) (*ExternalServiceConnection, error) {
    var token *TokenResponse
    var username string
    var err error

    switch service {
    case "trakt":
        token, err = s.trakt.ExchangeCode(ctx, code)
        if err != nil {
            return nil, err
        }
        // Get username
        user, _ := s.trakt.GetUser(ctx, token.AccessToken)
        username = user.Username

    case "lastfm":
        // Last.fm uses session key, not OAuth
        session, err := s.lastfm.GetSession(ctx, code)
        if err != nil {
            return nil, err
        }
        token = &TokenResponse{AccessToken: session.Key}
        username = session.Username

    case "listenbrainz":
        // ListenBrainz uses user token (no OAuth)
        token = &TokenResponse{AccessToken: code}
        user, _ := s.listenbrainz.ValidateToken(ctx, code)
        username = user.Username

    default:
        return nil, fmt.Errorf("unknown service: %s", service)
    }

    // Store connection
    conn := &ExternalServiceConnection{
        ID:        uuid.New(),
        Service:   service,
        Username:  username,
        Enabled:   true,
    }

    _, err = s.db.Exec(ctx, `
        INSERT INTO user_external_services (id, user_id, service, access_token, refresh_token, token_expires, username, enabled)
        VALUES ($1, $2, $3, $4, $5, $6, $7, true)
        ON CONFLICT (user_id, service) DO UPDATE SET
            access_token = EXCLUDED.access_token,
            refresh_token = EXCLUDED.refresh_token,
            token_expires = EXCLUDED.token_expires,
            username = EXCLUDED.username,
            connected_at = NOW()
    `, conn.ID, userID, service, token.AccessToken, token.RefreshToken, token.ExpiresAt, username)

    return conn, err
}

func (s *ExternalServiceService) Disconnect(ctx context.Context, userID uuid.UUID, service string) error {
    _, err := s.db.Exec(ctx, `
        DELETE FROM user_external_services
        WHERE user_id = $1 AND service = $2
    `, userID, service)
    return err
}
```

---

## History Sync

### Import from External Services

```go
type HistorySyncArgs struct {
    UserID  uuid.UUID `json:"user_id"`
    Service string    `json:"service"`
    Since   time.Time `json:"since"`
}

func (HistorySyncArgs) Kind() string { return "sync.history" }

type HistorySyncWorker struct {
    river.WorkerDefaults[HistorySyncArgs]
    trakt    *TraktClient
    movies   *MovieRepository
    tvshows  *TVShowRepository
    users    *UserRepository
}

func (w *HistorySyncWorker) Work(ctx context.Context, job *river.Job[HistorySyncArgs]) error {
    args := job.Args

    user, err := w.users.GetByID(ctx, args.UserID)
    if err != nil {
        return err
    }

    switch args.Service {
    case "trakt":
        return w.syncFromTrakt(ctx, user, args.Since)
    }

    return nil
}

func (w *HistorySyncWorker) syncFromTrakt(ctx context.Context, user *User, since time.Time) error {
    token := user.GetServiceToken("trakt")
    if token == "" {
        return errors.New("trakt not connected")
    }

    // Get watched history from Trakt
    history, err := w.trakt.SyncHistory(ctx, token, since)
    if err != nil {
        return err
    }

    for _, item := range history.Items {
        switch item.Type {
        case "movie":
            // Find movie by TMDB ID
            movie, err := w.movies.GetByTMDBID(ctx, item.Movie.IDs.TMDB)
            if err != nil {
                continue // Movie not in library
            }

            // Mark as watched
            w.movies.AddToWatchHistory(ctx, user.ID, movie.ID, item.WatchedAt)

        case "episode":
            // Find episode
            episode, err := w.tvshows.GetEpisodeByTVDBID(ctx, item.Episode.IDs.TVDB)
            if err != nil {
                continue
            }

            w.tvshows.AddToWatchHistory(ctx, user.ID, episode.ID, item.WatchedAt)
        }
    }

    // Update last sync time
    w.users.UpdateServiceLastSync(ctx, user.ID, "trakt", time.Now())

    return nil
}
```

### Export to External Services

```go
type ExportHistoryArgs struct {
    UserID  uuid.UUID `json:"user_id"`
    Service string    `json:"service"`
}

func (ExportHistoryArgs) Kind() string { return "sync.export_history" }

func (w *HistorySyncWorker) exportToTrakt(ctx context.Context, user *User) error {
    token := user.GetServiceToken("trakt")
    if token == "" {
        return errors.New("trakt not connected")
    }

    // Get local watch history not yet synced
    lastSync := user.GetServiceLastSync("trakt")

    movies, err := w.movies.GetWatchHistorySince(ctx, user.ID, lastSync)
    if err != nil {
        return err
    }

    episodes, err := w.tvshows.GetWatchHistorySince(ctx, user.ID, lastSync)
    if err != nil {
        return err
    }

    // Build Trakt history items
    items := make([]TraktHistoryItem, 0, len(movies)+len(episodes))

    for _, m := range movies {
        items = append(items, TraktHistoryItem{
            Type:      "movie",
            WatchedAt: m.WatchedAt,
            Movie:     &TraktMovie{IDs: TraktIDs{TMDB: m.TMDBID}},
        })
    }

    for _, e := range episodes {
        items = append(items, TraktHistoryItem{
            Type:      "episode",
            WatchedAt: e.WatchedAt,
            Episode:   &TraktEpisode{Season: e.Season, Number: e.Number},
            Show:      &TraktShow{IDs: TraktIDs{TVDB: e.ShowTVDBID}},
        })
    }

    // Push to Trakt
    return w.trakt.AddToHistory(ctx, token, items)
}
```

---

## Configuration

```yaml
# configs/config.yaml
scrobbling:
  enabled: true

  # Scrobble thresholds
  video_threshold: 80       # % watched to scrobble
  music_threshold: 50       # % listened OR 4 minutes

  services:
    trakt:
      enabled: true
      client_id: "${TRAKT_CLIENT_ID}"
      client_secret: "${TRAKT_CLIENT_SECRET}"
      redirect_uri: "https://revenge.example.com/api/v1/auth/callback/trakt"

    lastfm:
      enabled: true
      api_key: "${LASTFM_API_KEY}"
      api_secret: "${LASTFM_API_SECRET}"

    listenbrainz:
      enabled: true
      # Uses user-provided tokens

    simkl:
      enabled: false
      client_id: "${SIMKL_CLIENT_ID}"
      client_secret: "${SIMKL_CLIENT_SECRET}"

  # Sync settings
  sync:
    import_on_connect: true    # Import history when connecting
    auto_sync_interval: 24h    # Periodic sync
    bidirectional: true        # Sync both ways
```

---

## API Endpoints

```yaml
# api/openapi/scrobbling.yaml
paths:
  /api/v1/user/services:
    get:
      summary: List connected external services

  /api/v1/user/services/{service}/connect:
    get:
      summary: Get OAuth URL for service
    post:
      summary: Complete OAuth connection

  /api/v1/user/services/{service}/disconnect:
    post:
      summary: Disconnect service

  /api/v1/user/services/{service}/sync:
    post:
      summary: Trigger manual sync

  /api/v1/user/services/{service}/settings:
    patch:
      summary: Update service settings
```

---

## Summary

| Service | Content | Auth | Features |
|---------|---------|------|----------|
| Trakt | Video | OAuth2 | Scrobble, history, ratings, watchlist |
| Last.fm | Music | Session key | Scrobble, now playing, loved |
| ListenBrainz | Music | User token | Scrobble, now playing |
| Simkl | Video | OAuth2 | Scrobble, history |

| Aspect | Implementation |
|--------|----------------|
| Queue | River jobs with retry |
| Rate Limiting | Per-service limits respected |
| Sync | Bidirectional, incremental |
| Storage | user_external_services table |
