# Release Calendar System

> Upcoming releases and recent additions calendar via Servarr integration

**Status**: 🔵 PLANNED
**Priority**: 🟡 MEDIUM
**Module**: `internal/calendar`
**Dependencies**: [Servarr Integration](../integrations/servarr/INDEX.md)

---

## Overview

The Release Calendar provides users with a unified view of:
- **Upcoming Releases** - Content coming soon (from Servarr calendars)
- **Recently Aired** - Episodes that just aired
- **Recently Added** - Content added to their library
- **Premieres** - New series/seasons starting

This feature integrates with Servarr calendar APIs to pull release dates and sync them locally.

## Goals

- Unified calendar across all content types (movies, episodes, music, books)
- Real-time sync with Servarr calendar APIs
- User-specific calendar based on monitored content
- Push notifications for upcoming releases (optional)

## Non-Goals

- Manual calendar event creation (we source from Servarr)
- Calendar editing (changes go through Servarr)
- Shared/family calendars

---

## Technical Design

### Database Schema

```sql
-- Calendar entries synced from Servarr
CREATE TABLE calendar_entries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Content reference (may not exist in library yet)
    content_type VARCHAR(20) NOT NULL,  -- 'movie', 'episode', 'album', 'book'
    local_content_id UUID,              -- FK if exists in library

    -- Servarr reference
    servarr_type VARCHAR(20) NOT NULL,  -- 'radarr', 'sonarr', 'lidarr', 'chaptarr'
    servarr_id INT NOT NULL,            -- ID in Servarr

    -- Event info
    title VARCHAR(500) NOT NULL,
    series_title VARCHAR(500),          -- For episodes
    season_number INT,
    episode_number INT,

    -- Release info
    release_date DATE NOT NULL,
    release_time TIME,                  -- For precise air times
    release_type VARCHAR(50),           -- 'theatrical', 'digital', 'physical', 'premiere', 'finale'

    -- Status
    is_downloaded BOOLEAN DEFAULT false,
    is_available BOOLEAN DEFAULT false, -- In library

    -- Metadata
    poster_url VARCHAR(512),
    overview TEXT,
    runtime_minutes INT,

    -- Sync tracking
    synced_at TIMESTAMPTZ DEFAULT NOW(),

    UNIQUE(servarr_type, servarr_id, release_date)
);

CREATE INDEX idx_calendar_date ON calendar_entries(release_date);
CREATE INDEX idx_calendar_type_date ON calendar_entries(content_type, release_date);
CREATE INDEX idx_calendar_servarr ON calendar_entries(servarr_type, servarr_id);

-- User calendar preferences
CREATE TABLE user_calendar_settings (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,

    -- Display settings
    show_movies BOOLEAN DEFAULT true,
    show_episodes BOOLEAN DEFAULT true,
    show_music BOOLEAN DEFAULT true,
    show_books BOOLEAN DEFAULT true,

    -- Filter settings
    show_monitored_only BOOLEAN DEFAULT true,  -- Only tracked content
    show_downloaded BOOLEAN DEFAULT true,      -- Show after download

    -- Notification settings
    notify_releases BOOLEAN DEFAULT false,
    notify_hours_before INT DEFAULT 24,

    updated_at TIMESTAMPTZ DEFAULT NOW()
);
```

### Servarr Calendar API Integration

#### Radarr Calendar

```go
// GET /api/v3/calendar?start={date}&end={date}
type RadarrCalendarEntry struct {
    ID              int       `json:"id"`
    Title           string    `json:"title"`
    SortTitle       string    `json:"sortTitle"`
    Status          string    `json:"status"`
    Overview        string    `json:"overview"`
    InCinemas       time.Time `json:"inCinemas"`
    PhysicalRelease time.Time `json:"physicalRelease"`
    DigitalRelease  time.Time `json:"digitalRelease"`
    HasFile         bool      `json:"hasFile"`
    Monitored       bool      `json:"monitored"`
    Images          []Image   `json:"images"`
}

func (c *RadarrClient) GetCalendar(ctx context.Context, start, end time.Time) ([]RadarrCalendarEntry, error) {
    url := fmt.Sprintf("%s/api/v3/calendar?start=%s&end=%s",
        c.baseURL,
        start.Format("2006-01-02"),
        end.Format("2006-01-02"),
    )
    // ...
}
```

#### Sonarr Calendar

```go
// GET /api/v3/calendar?start={date}&end={date}&includeSeries=true
type SonarrCalendarEntry struct {
    ID              int       `json:"id"`
    SeriesID        int       `json:"seriesId"`
    EpisodeFileID   int       `json:"episodeFileId"`
    SeasonNumber    int       `json:"seasonNumber"`
    EpisodeNumber   int       `json:"episodeNumber"`
    Title           string    `json:"title"`
    AirDate         string    `json:"airDate"`
    AirDateUtc      time.Time `json:"airDateUtc"`
    Overview        string    `json:"overview"`
    HasFile         bool      `json:"hasFile"`
    Monitored       bool      `json:"monitored"`
    Series          *Series   `json:"series"`  // When includeSeries=true
}

func (c *SonarrClient) GetCalendar(ctx context.Context, start, end time.Time) ([]SonarrCalendarEntry, error) {
    url := fmt.Sprintf("%s/api/v3/calendar?start=%s&end=%s&includeSeries=true",
        c.baseURL,
        start.Format("2006-01-02"),
        end.Format("2006-01-02"),
    )
    // ...
}
```

#### Lidarr Calendar

```go
// GET /api/v1/calendar?start={date}&end={date}
type LidarrCalendarEntry struct {
    ID          int       `json:"id"`
    ArtistID    int       `json:"artistId"`
    Title       string    `json:"title"`
    ReleaseDate time.Time `json:"releaseDate"`
    AlbumType   string    `json:"albumType"`
    Monitored   bool      `json:"monitored"`
    Artist      *Artist   `json:"artist"`
}
```

#### Chaptarr Calendar (Readarr API)

```go
// GET /api/v1/calendar?start={date}&end={date}
type ChaptarrCalendarEntry struct {
    ID          int       `json:"id"`
    AuthorID    int       `json:"authorId"`
    Title       string    `json:"title"`
    ReleaseDate time.Time `json:"releaseDate"`
    Monitored   bool      `json:"monitored"`
    Author      *Author   `json:"author"`
}
```

### Repository Interface

```go
type CalendarRepository interface {
    // Entries
    GetEntriesByDateRange(ctx context.Context, start, end time.Time, types []string) ([]CalendarEntry, error)
    GetUpcomingEntries(ctx context.Context, days int, types []string) ([]CalendarEntry, error)
    GetRecentEntries(ctx context.Context, days int, types []string) ([]CalendarEntry, error)
    UpsertEntry(ctx context.Context, entry *CalendarEntry) error
    DeleteOldEntries(ctx context.Context, before time.Time) error

    // User settings
    GetUserSettings(ctx context.Context, userID uuid.UUID) (*UserCalendarSettings, error)
    UpdateUserSettings(ctx context.Context, settings *UserCalendarSettings) error

    // Sync tracking
    GetLastSyncTime(ctx context.Context, servarrType string) (time.Time, error)
    UpdateLastSyncTime(ctx context.Context, servarrType string, syncedAt time.Time) error
}
```

### Service Layer

```go
type CalendarService struct {
    repo     CalendarRepository
    radarr   *RadarrClient
    sonarr   *SonarrClient
    lidarr   *LidarrClient
    chaptarr *ChaptarrClient
    logger   *slog.Logger
}

// GetCalendar returns unified calendar entries
func (s *CalendarService) GetCalendar(ctx context.Context, userID uuid.UUID, start, end time.Time) (*Calendar, error) {
    settings, _ := s.repo.GetUserSettings(ctx, userID)

    var types []string
    if settings.ShowMovies {
        types = append(types, "movie")
    }
    if settings.ShowEpisodes {
        types = append(types, "episode")
    }
    // ...

    entries, err := s.repo.GetEntriesByDateRange(ctx, start, end, types)
    if err != nil {
        return nil, err
    }

    return &Calendar{
        Entries:   entries,
        StartDate: start,
        EndDate:   end,
    }, nil
}

// SyncFromServarr pulls calendar data from all Servarr instances
func (s *CalendarService) SyncFromServarr(ctx context.Context) error {
    start := time.Now().AddDate(0, -1, 0)  // 1 month ago
    end := time.Now().AddDate(0, 3, 0)     // 3 months ahead

    // Sync in parallel
    var wg sync.WaitGroup
    errCh := make(chan error, 4)

    wg.Add(4)
    go func() { defer wg.Done(); errCh <- s.syncRadarr(ctx, start, end) }()
    go func() { defer wg.Done(); errCh <- s.syncSonarr(ctx, start, end) }()
    go func() { defer wg.Done(); errCh <- s.syncLidarr(ctx, start, end) }()
    go func() { defer wg.Done(); errCh <- s.syncChaptarr(ctx, start, end) }()

    wg.Wait()
    close(errCh)

    // Collect errors
    var errs []error
    for err := range errCh {
        if err != nil {
            errs = append(errs, err)
        }
    }

    if len(errs) > 0 {
        return fmt.Errorf("calendar sync errors: %v", errs)
    }
    return nil
}
```

### API Endpoints

```
GET /api/calendar
Query Parameters:
  - start (date, default: today - 7 days)
  - end (date, default: today + 30 days)
  - types (string[], default: all enabled)

Response:
{
  "entries": [
    {
      "id": "uuid",
      "content_type": "episode",
      "title": "Episode Title",
      "series_title": "Series Name",
      "season_number": 2,
      "episode_number": 5,
      "release_date": "2026-02-01",
      "release_time": "21:00:00",
      "release_type": "premiere",
      "is_downloaded": false,
      "is_available": false,
      "poster_url": "https://...",
      "overview": "..."
    }
  ],
  "start_date": "2026-01-22",
  "end_date": "2026-02-28"
}

GET /api/calendar/upcoming
Query Parameters:
  - days (int, default: 7)
  - limit (int, default: 20)

GET /api/calendar/recent
Query Parameters:
  - days (int, default: 7)
  - limit (int, default: 20)

GET /api/users/{userId}/calendar/settings
PUT /api/users/{userId}/calendar/settings
```

### River Jobs

```go
// Sync calendar from Servarr instances
type CalendarSyncArgs struct{}

func (CalendarSyncArgs) Kind() string { return "calendar.sync" }

func (CalendarSyncArgs) InsertOpts() river.InsertOpts {
    return river.InsertOpts{
        Queue:    "default",
        UniqueOpts: river.UniqueOpts{
            ByPeriod: 15 * time.Minute,  // Don't sync more than once per 15 min
        },
    }
}

type CalendarSyncWorker struct {
    river.WorkerDefaults[CalendarSyncArgs]
    service *CalendarService
    logger  *slog.Logger
}

func (w *CalendarSyncWorker) Work(ctx context.Context, job *river.Job[CalendarSyncArgs]) error {
    return w.service.SyncFromServarr(ctx)
}

// Schedule periodic sync
func ScheduleCalendarSync(client *river.Client) error {
    return client.InsertPeriodic(
        river.PeriodicJobOpts{
            Every: 1 * time.Hour,
        },
        CalendarSyncArgs{},
    )
}
```

---

## Configuration

```yaml
calendar:
  enabled: true
  sync_interval: 1h        # How often to sync from Servarr
  upcoming_days: 90        # How far ahead to fetch
  past_days: 30            # How far back to keep
  cleanup_days: 180        # Delete entries older than this
```

---

## Implementation Checklist

### Database
- [ ] Migration: `calendar_entries` table
- [ ] Migration: `user_calendar_settings` table
- [ ] Indexes for date range queries

### sqlc Queries
- [ ] `queries/calendar/entries.sql`
- [ ] `queries/calendar/settings.sql`

### Repository
- [ ] `internal/calendar/repository.go`
- [ ] `internal/calendar/repository_pg.go`

### Servarr Clients
- [ ] Radarr calendar endpoint
- [ ] Sonarr calendar endpoint
- [ ] Lidarr calendar endpoint
- [ ] Chaptarr calendar endpoint

### Service
- [ ] `internal/calendar/service.go`
- [ ] Unified calendar generation
- [ ] Servarr sync logic

### Jobs
- [ ] `internal/calendar/jobs.go`
- [ ] Periodic sync job
- [ ] Notification job (optional)

### API
- [ ] `GET /api/calendar`
- [ ] `GET /api/calendar/upcoming`
- [ ] `GET /api/calendar/recent`
- [ ] User settings endpoints

### Frontend
- [ ] Calendar view component
- [ ] Week/month view toggle
- [ ] Content type filters

---

## Related Documentation

- [Servarr Integration](../integrations/servarr/INDEX.md) - Servarr API details
- [Radarr](../integrations/servarr/RADARR.md) - Movie calendar
- [Sonarr](../integrations/servarr/SONARR.md) - Episode calendar
- [Lidarr](../integrations/servarr/LIDARR.md) - Album calendar
- [Chaptarr](../integrations/servarr/CHAPTARR.md) - Book calendar
