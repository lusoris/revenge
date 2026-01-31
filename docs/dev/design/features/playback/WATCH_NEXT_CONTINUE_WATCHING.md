# Watch Next & Continue Watching System

> Intelligent playback continuation and recommendation system

## Status

| Dimension           | Status | Notes |
| ------------------- | ------ | ----- |
| Design              | ✅     |       |
| Sources             | ✅     |       |
| Instructions        | ✅     |       |
| Code                | 🔴     |       |
| Linting             | 🔴     |       |
| Unit Testing        | 🔴     |       |
| Integration Testing | 🔴     |       |

**Priority**: 🔴 HIGH (Core UX Feature)
**Related**: [METADATA_SYSTEM.md](../architecture/03_METADATA_SYSTEM.md), [SCROBBLING.md](SCROBBLING.md)
**Location**: `internal/feature/watchnext/`

---

## Developer Resources

| Source       | URL                                                                   | Purpose                    |
| ------------ | --------------------------------------------------------------------- | -------------------------- |
| Jellyfin API | [api.jellyfin.org](https://api.jellyfin.org/)                         | Continue watching patterns |
| TMDb API     | [developers.themoviedb.org/3](https://developers.themoviedb.org/3)    | Next episode metadata      |
| Trakt API    | [trakt.docs.apiary.io](https://trakt.docs.apiary.io/)                 | Watch history sync         |

---

## Overview

The Watch Next / Continue Watching system provides users with intelligent suggestions for what to watch based on their viewing history. This is a core UX feature that significantly impacts user engagement.

### Goals

1. **Continue Watching** - Resume content in progress
2. **Watch Next** - Suggest the next episode/movie to watch
3. **Up Next** - Auto-play queue for binge watching
4. **Cross-Device Sync** - Resume anywhere

---

## Feature Components

### 1. Continue Watching

Shows content the user has started but not finished.

**Inclusion Criteria:**
- Position > 5% of total duration (skip accidental plays)
- Position < 90% of total duration (not finished)
- Last watched within 30 days (configurable)
- Not marked as "watched" by user

**Sorting:**
- Primary: Last played timestamp (most recent first)
- Secondary: Content type priority (episodes > movies)

**Database Schema:**

```sql
-- Movies: watch_history table
-- Episodes: episode_watch_history table

-- Query for Continue Watching (movies)
SELECT m.*, wh.position_ticks, wh.duration_ticks, wh.last_played_at
FROM movies m
JOIN watch_history wh ON m.id = wh.movie_id
WHERE wh.user_id = $1
  AND wh.completed = false
  AND wh.position_ticks > (wh.duration_ticks * 0.05)  -- > 5%
  AND wh.position_ticks < (wh.duration_ticks * 0.90)  -- < 90%
  AND wh.last_played_at > NOW() - INTERVAL '30 days'
ORDER BY wh.last_played_at DESC
LIMIT 20;

-- Query for Continue Watching (episodes)
SELECT e.*, s.title as series_title, ewh.position_ticks, ewh.duration_ticks, ewh.last_played_at
FROM episodes e
JOIN series s ON e.series_id = s.id
JOIN episode_watch_history ewh ON e.id = ewh.episode_id
WHERE ewh.user_id = $1
  AND ewh.completed = false
  AND ewh.position_ticks > (ewh.duration_ticks * 0.05)
  AND ewh.position_ticks < (ewh.duration_ticks * 0.90)
  AND ewh.last_played_at > NOW() - INTERVAL '30 days'
ORDER BY ewh.last_played_at DESC
LIMIT 20;
```

### 2. Watch Next (Series)

For TV series, determines the next episode to watch after the current one.

**Algorithm:**

```go
func GetNextEpisode(ctx context.Context, userID, seriesID uuid.UUID) (*Episode, error) {
    // 1. Get series watch progress
    progress, err := repo.GetSeriesWatchProgress(ctx, userID, seriesID)
    if err != nil {
        return nil, err
    }

    // 2. If user has progress, get next after last watched
    if progress != nil && progress.LastEpisodeID != uuid.Nil {
        lastEp, _ := repo.GetEpisodeByID(ctx, progress.LastEpisodeID)
        if lastEp != nil {
            // Try next episode in same season
            next, err := repo.GetNextEpisode(ctx, seriesID, lastEp.SeasonNumber, lastEp.EpisodeNumber)
            if err == nil && next != nil {
                return next, nil
            }
            // If no next in season, try first episode of next season
            nextSeason, err := repo.GetSeasonByNumber(ctx, seriesID, lastEp.SeasonNumber + 1)
            if err == nil && nextSeason != nil {
                return repo.GetEpisodeByNumber(ctx, seriesID, nextSeason.SeasonNumber, 1)
            }
        }
    }

    // 3. No progress - return first unwatched episode
    return repo.GetFirstUnwatchedEpisode(ctx, userID, seriesID)
}

func GetFirstUnwatchedEpisode(ctx context.Context, userID, seriesID uuid.UUID) (*Episode, error) {
    // Find first episode user hasn't watched
    return db.Query(`
        SELECT e.* FROM episodes e
        LEFT JOIN episode_watch_history ewh ON e.id = ewh.episode_id AND ewh.user_id = $1
        WHERE e.series_id = $2
          AND (ewh.completed IS NULL OR ewh.completed = false)
          AND e.season_number > 0  -- Skip specials
        ORDER BY e.season_number ASC, e.episode_number ASC
        LIMIT 1
    `, userID, seriesID)
}
```

### 3. Series Watch Progress

Tracks overall progress through a series.

**Database Schema:**

```sql
CREATE TABLE series_watch_progress (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    series_id UUID NOT NULL REFERENCES series(id) ON DELETE CASCADE,

    -- Progress tracking
    last_episode_id UUID REFERENCES episodes(id),
    last_season_number INT,
    last_episode_number INT,

    -- Statistics
    watched_episodes INT DEFAULT 0,
    total_episodes INT DEFAULT 0,
    watched_runtime_ticks BIGINT DEFAULT 0,
    total_runtime_ticks BIGINT DEFAULT 0,

    -- Timestamps
    started_at TIMESTAMPTZ DEFAULT NOW(),
    last_watched_at TIMESTAMPTZ DEFAULT NOW(),
    completed_at TIMESTAMPTZ,

    -- Status
    is_completed BOOLEAN DEFAULT false,

    UNIQUE(user_id, series_id)
);

CREATE INDEX idx_series_watch_progress_user ON series_watch_progress(user_id);
CREATE INDEX idx_series_watch_progress_last_watched ON series_watch_progress(last_watched_at DESC);
```

**Progress Update Logic:**

```go
func UpdateSeriesProgress(ctx context.Context, userID uuid.UUID, episode *Episode) error {
    // Get or create progress
    progress, _ := repo.GetSeriesWatchProgress(ctx, userID, episode.SeriesID)
    if progress == nil {
        progress = &SeriesWatchProgress{
            UserID:   userID,
            SeriesID: episode.SeriesID,
        }
    }

    // Update last watched
    progress.LastEpisodeID = episode.ID
    progress.LastSeasonNumber = episode.SeasonNumber
    progress.LastEpisodeNumber = episode.EpisodeNumber
    progress.LastWatchedAt = time.Now()

    // Recalculate statistics
    watchedCount, _ := repo.CountWatchedEpisodesBySeries(ctx, userID, episode.SeriesID)
    totalCount, _ := repo.CountEpisodesBySeries(ctx, episode.SeriesID)

    progress.WatchedEpisodes = watchedCount
    progress.TotalEpisodes = totalCount

    // Check if completed
    if watchedCount >= totalCount {
        progress.IsCompleted = true
        progress.CompletedAt = time.Now()
    }

    return repo.UpsertSeriesWatchProgress(ctx, progress)
}
```

### 4. Up Next / Auto-Play Queue

For binge watching, automatically plays the next episode.

**Implementation:**

```go
type UpNextQueue struct {
    CurrentItem  *PlaybackItem
    NextItem     *PlaybackItem
    AutoPlayIn   time.Duration  // Countdown timer (e.g., 15 seconds)
    QueuedItems  []*PlaybackItem // Future items
}

type PlaybackItem struct {
    ID           uuid.UUID
    Type         string  // "movie", "episode"
    Title        string
    SeriesTitle  string  // For episodes
    SeasonNum    int
    EpisodeNum   int
    ThumbnailURL string
    RuntimeTicks int64
}

func BuildUpNextQueue(ctx context.Context, userID uuid.UUID, currentEpisode *Episode) (*UpNextQueue, error) {
    queue := &UpNextQueue{
        CurrentItem: toPlaybackItem(currentEpisode),
        AutoPlayIn:  15 * time.Second,  // Configurable
    }

    // Get next 5 episodes
    episodes, err := repo.GetNextEpisodes(ctx, currentEpisode.SeriesID,
        currentEpisode.SeasonNumber, currentEpisode.EpisodeNumber, 5)
    if err != nil {
        return queue, nil
    }

    if len(episodes) > 0 {
        queue.NextItem = toPlaybackItem(episodes[0])
        for _, ep := range episodes[1:] {
            queue.QueuedItems = append(queue.QueuedItems, toPlaybackItem(ep))
        }
    }

    return queue, nil
}
```

### 5. Cross-Device Sync

Uses WebSocket for real-time sync and polling for offline devices.

**Real-time Sync (WebSocket):**

```go
// When playback position updates
func (h *PlaybackHandler) OnPositionUpdate(ctx context.Context, userID uuid.UUID, update PositionUpdate) {
    // 1. Save to database
    h.repo.UpdateWatchProgress(ctx, userID, update.ContentID, update.PositionTicks)

    // 2. Broadcast to user's other sessions
    h.wsHub.BroadcastToUser(userID, WebSocketMessage{
        Type: "playback_sync",
        Data: PlaybackSyncData{
            ContentID:     update.ContentID,
            PositionTicks: update.PositionTicks,
            Timestamp:     time.Now(),
        },
    })
}
```

**Polling Fallback:**

```go
// For clients that don't support WebSocket
// GET /api/sync/playback?since={timestamp}
func (h *SyncHandler) GetPlaybackChanges(ctx context.Context, userID uuid.UUID, since time.Time) []PlaybackChange {
    return h.repo.GetPlaybackChangesSince(ctx, userID, since)
}
```

---

## API Endpoints

### Continue Watching

```
GET /api/users/{userId}/continue-watching
Query Parameters:
  - limit (int, default: 20)
  - include_types (string[], default: ["movie", "episode"])

Response:
{
  "items": [
    {
      "id": "uuid",
      "type": "episode",
      "title": "Episode Title",
      "series_id": "uuid",
      "series_title": "Series Name",
      "season_number": 1,
      "episode_number": 3,
      "position_ticks": 18000000000,
      "duration_ticks": 36000000000,
      "percent_complete": 50.0,
      "last_played_at": "2026-01-29T10:00:00Z",
      "thumbnail_url": "/items/{id}/images/still"
    }
  ]
}
```

### Watch Next (Series)

```
GET /api/series/{seriesId}/watch-next
Query Parameters:
  - user_id (uuid, required)

Response:
{
  "episode": {
    "id": "uuid",
    "title": "Next Episode",
    "season_number": 2,
    "episode_number": 1,
    "overview": "...",
    "runtime_ticks": 36000000000,
    "still_url": "/episodes/{id}/images/still"
  },
  "series_progress": {
    "watched_episodes": 10,
    "total_episodes": 24,
    "percent_complete": 41.67
  }
}
```

### Up Next Queue

```
GET /api/playback/up-next
Query Parameters:
  - current_id (uuid, required)
  - current_type (string, required: "movie" | "episode")
  - queue_size (int, default: 5)

Response:
{
  "next": { ... },
  "auto_play_seconds": 15,
  "queue": [ ... ]
}
```

---

## User Settings

```sql
ALTER TABLE users ADD COLUMN auto_play_enabled BOOLEAN DEFAULT true;
ALTER TABLE users ADD COLUMN auto_play_delay_seconds INT DEFAULT 15;
ALTER TABLE users ADD COLUMN continue_watching_days INT DEFAULT 30;
ALTER TABLE users ADD COLUMN mark_watched_percent INT DEFAULT 90;  -- Mark as watched when 90% complete
```

---

## Event Triggers

### On Playback Start

```go
func OnPlaybackStart(ctx context.Context, userID uuid.UUID, contentID uuid.UUID, contentType string) {
    // 1. Create or update watch history
    // 2. Update "last played" timestamp
    // 3. Emit event for scrobbling services
}
```

### On Playback Progress

```go
func OnPlaybackProgress(ctx context.Context, userID uuid.UUID, contentID uuid.UUID, positionTicks int64) {
    // 1. Update position in watch history
    // 2. Check if passed "mark as watched" threshold
    // 3. Sync to other devices via WebSocket
}
```

### On Playback Stop

```go
func OnPlaybackStop(ctx context.Context, userID uuid.UUID, contentID uuid.UUID, positionTicks int64) {
    // 1. Final position update
    // 2. Calculate if completed
    // 3. Update series progress (if episode)
    // 4. Emit scrobble event if completed
}
```

---

## Implementation Checklist

- [ ] **Database**
  - [ ] `series_watch_progress` table and indexes
  - [ ] Update `watch_history` with position tracking
  - [ ] Update `episode_watch_history` with position tracking

- [ ] **Repository**
  - [ ] `GetContinueWatchingMovies(userID, limit)`
  - [ ] `GetContinueWatchingEpisodes(userID, limit)`
  - [ ] `GetSeriesWatchProgress(userID, seriesID)`
  - [ ] `GetNextEpisode(seriesID, seasonNum, episodeNum)`
  - [ ] `GetFirstUnwatchedEpisode(userID, seriesID)`
  - [ ] `CountWatchedEpisodesBySeries(userID, seriesID)`

- [ ] **Service**
  - [ ] `ContinueWatchingService` with merged results
  - [ ] `WatchNextService` for series navigation
  - [ ] `UpNextQueueService` for auto-play

- [ ] **API Handlers**
  - [ ] `GET /api/users/{userId}/continue-watching`
  - [ ] `GET /api/series/{seriesId}/watch-next`
  - [ ] `GET /api/playback/up-next`

- [ ] **WebSocket**
  - [ ] Playback sync messages
  - [ ] Position update broadcasts

- [ ] **Client Integration**
  - [ ] Continue Watching carousel on home
  - [ ] "Resume" button on content pages
  - [ ] Up Next overlay at end of playback
  - [ ] Auto-play countdown

---


<!-- SOURCE-BREADCRUMBS-START -->

## Sources & Cross-References

> Auto-generated section linking to external documentation sources

### Cross-Reference Indexes

- [All Sources Index](../../../sources/SOURCES_INDEX.md) - Complete list of external documentation
- [Design ↔ Sources Map](../../../sources/DESIGN_CROSSREF.md) - Which docs reference which sources

### Referenced Sources

| Source | Documentation |
|--------|---------------|
| [Jellyfin API](https://api.jellyfin.org/) | [Local](../../../sources/apis/jellyfin.md) |

<!-- SOURCE-BREADCRUMBS-END -->

<!-- DESIGN-BREADCRUMBS-START -->

## Related Design Docs

> Auto-generated cross-references to related design documentation

**Category**: [Playback](INDEX.md)

### In This Section

- [Revenge - Media Enhancement Features](MEDIA_ENHANCEMENTS.md)
- [Release Calendar System](RELEASE_CALENDAR.md)
- [Skip Intro / Credits Detection](SKIP_INTRO.md)
- [SyncPlay (Watch Together)](SYNCPLAY.md)
- [Trickplay (Timeline Thumbnails)](TRICKPLAY.md)

### Related Topics

- [Revenge - Architecture v2](../../architecture/01_ARCHITECTURE.md) _Architecture_
- [Revenge - Design Principles](../../architecture/02_DESIGN_PRINCIPLES.md) _Architecture_
- [Revenge - Metadata System](../../architecture/03_METADATA_SYSTEM.md) _Architecture_
- [Revenge - Player Architecture](../../architecture/04_PLAYER_ARCHITECTURE.md) _Architecture_
- [Plugin Architecture Decision](../../architecture/05_PLUGIN_ARCHITECTURE_DECISION.md) _Architecture_

### Indexes

- [Design Index](../../DESIGN_INDEX.md) - All design docs by category/topic
- [Source of Truth](../../00_SOURCE_OF_TRUTH.md) - Package versions and status

<!-- DESIGN-BREADCRUMBS-END -->

## Related Documentation

- [SCROBBLING.md](SCROBBLING.md) - External service sync (Trakt, etc.)
- [USER_EXPERIENCE_FEATURES.md](USER_EXPERIENCE_FEATURES.md) - Other UX features
- [PLAYER_ARCHITECTURE.md](../architecture/04_PLAYER_ARCHITECTURE.md) - Video player design
