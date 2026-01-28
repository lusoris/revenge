---
applyTo: "**/internal/service/scrobble/**/*.go,**/internal/service/sync/**/*.go"
---

# External Services Integration Instructions

> Guidelines for integrating with external services like Trakt, Last.fm, ListenBrainz.

## General Patterns

### Service Connection Model

```go
type UserServiceConnection struct {
    ID          uuid.UUID
    UserID      uuid.UUID
    Service     string    // "trakt", "lastfm", "listenbrainz"
    Enabled     bool
    AccessToken string    // Encrypted
    TokenExpiry time.Time // OAuth2 only
    Settings    map[string]any
}
```

### OAuth2 Services (Trakt, Simkl)

```go
// Good: Use oauth2 package
oauth2Conf := &oauth2.Config{
    ClientID:     cfg.ClientID,
    ClientSecret: cfg.ClientSecret,
    Scopes:       []string{"sync"},
    Endpoint: oauth2.Endpoint{
        AuthURL:  "https://trakt.tv/oauth/authorize",
        TokenURL: "https://trakt.tv/oauth/token",
    },
}

// Good: Auto-refresh tokens
token, err := oauth2Conf.TokenSource(ctx, storedToken).Token()
if err != nil {
    // Token refresh failed, user needs to re-authenticate
}
```

### API Key Services (Last.fm, ListenBrainz)

```go
// Good: Session key validation
func (c *LastFMClient) ValidateSession(ctx context.Context, sessionKey string) (bool, error) {
    // Call user.getInfo to verify session is still valid
}
```

## Scrobbling Guidelines

### When to Scrobble

| Content Type | Trigger Point             | Minimum Duration |
| ------------ | ------------------------- | ---------------- |
| Music track  | 50% complete OR 4 minutes | 30 seconds       |
| Podcast      | 80% complete              | 5 minutes        |
| Video        | 80% complete              | 10 minutes       |

```go
// Good: Check multiple thresholds
func shouldScrobble(progress PlaybackProgress) bool {
    // Music: 50% OR 4 minutes
    if progress.Type == "track" {
        return progress.Percentage >= 50 ||
               progress.WatchedMs >= 4*60*1000
    }
    // Video/Podcast: 80%
    return progress.Percentage >= 80
}
```

### Scrobble Job Pattern

```go
// Good: Use River jobs for reliability
type ScrobbleArgs struct {
    UserID      uuid.UUID `json:"user_id"`
    Service     string    `json:"service"`
    ItemID      uuid.UUID `json:"item_id"`
    ItemType    string    `json:"item_type"`
    WatchedAt   time.Time `json:"watched_at"`
}

func (ScrobbleArgs) Kind() string { return "scrobble" }

// Good: Handle rate limits with retry
func (w *ScrobbleWorker) Work(ctx context.Context, job *river.Job[ScrobbleArgs]) error {
    err := w.sendScrobble(ctx, job.Args)
    if isRateLimited(err) {
        return river.JobSnooze(5 * time.Minute)
    }
    return err
}
```

### Deduplication

```go
// Good: Track scrobbled items to prevent duplicates
CREATE TABLE scrobble_history (
    id          UUID PRIMARY KEY,
    user_id     UUID NOT NULL,
    service     VARCHAR(50) NOT NULL,
    item_id     UUID NOT NULL,
    scrobbled_at TIMESTAMPTZ NOT NULL,
    UNIQUE (user_id, service, item_id, scrobbled_at)
);
```

## Sync Guidelines

### Import History

```go
// Good: Batch import with progress tracking
type ImportProgress struct {
    Total     int
    Processed int
    Matched   int
    Failed    int
}

// Good: Match by external IDs, not title
func (s *TraktSync) matchItem(traktItem TraktItem) (uuid.UUID, error) {
    // Try TMDb ID first
    if traktItem.Movie.IDs.Tmdb > 0 {
        return s.movies.GetByTmdbID(ctx, traktItem.Movie.IDs.Tmdb)
    }
    // Fall back to IMDb
    if traktItem.Movie.IDs.Imdb != "" {
        return s.movies.GetByImdbID(ctx, traktItem.Movie.IDs.Imdb)
    }
    return uuid.Nil, ErrNoMatch
}
```

### Export Ratings

```go
// Good: Queue exports as background jobs
type ExportRatingArgs struct {
    UserID   uuid.UUID `json:"user_id"`
    Service  string    `json:"service"`
    ItemID   uuid.UUID `json:"item_id"`
    ItemType string    `json:"item_type"`
    Rating   float64   `json:"rating"`
}
```

## Rate Limiting

| Service      | Rate Limit | Strategy      |
| ------------ | ---------- | ------------- |
| Trakt        | 1000/5 min | Token bucket  |
| Last.fm      | 5/second   | Token bucket  |
| ListenBrainz | 1/second   | Delay between |
| TMDb         | 40/10 sec  | Token bucket  |

```go
// Good: Per-service rate limiters
rateLimiters := map[string]*resilience.TokenBucketLimiter{
    "trakt":       resilience.NewTokenBucket(200, time.Minute),
    "lastfm":      resilience.NewTokenBucket(5, time.Second),
    "listenbrainz": resilience.NewTokenBucket(1, time.Second),
}
```

## Error Handling

```go
// Good: Service-specific error types
var (
    ErrServiceUnavailable = errors.New("service temporarily unavailable")
    ErrInvalidToken       = errors.New("token expired or revoked")
    ErrRateLimited        = errors.New("rate limit exceeded")
    ErrItemNotFound       = errors.New("item not found on service")
)

// Good: Handle auth errors by notifying user
if errors.Is(err, ErrInvalidToken) {
    // Disable connection and notify user to re-authenticate
    s.disableConnection(ctx, userID, service)
    s.notifyUser(ctx, userID, "reconnect_required", service)
}
```

## Testing

```go
// Good: Mock external services
type MockTraktClient struct {
    ScrobbleFn func(ctx context.Context, item TraktItem) error
}

// Good: Use VCR pattern for integration tests
func TestTraktSync(t *testing.T) {
    // Record/playback HTTP interactions
    r := recorder.NewCassette("trakt_sync")
    defer r.Stop()

    client := &http.Client{Transport: r}
    // ...
}
```
