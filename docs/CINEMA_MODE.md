# Revenge - Cinema Mode

> Theatrical experience for home viewing with preroll, trailers, and intermissions.

## Overview

Cinema Mode transforms home viewing into a theatrical experience with customizable sequences of content before, during, and after the main feature.

---

## Feature Summary

| Feature | Description |
|---------|-------------|
| **Preroll** | Content before main feature (studio logos, custom clips) |
| **Trailers** | Automatic or curated trailer selection |
| **Intermission** | Mid-movie breaks for long films |
| **Postroll** | Content after credits (bonus, up-next) |
| **Ambient Mode** | Dimming, theater sounds between content |

---

## Preroll System

### Preroll Types

| Type | Description | Example |
|------|-------------|---------|
| `studio_intro` | Studio/network fanfare | THX, Dolby Atmos, custom |
| `custom_clip` | User-uploaded video | Personal intro, family clip |
| `seasonal` | Date-based content | Halloween, Christmas themes |
| `random` | Random from pool | Variety on each play |
| `library_specific` | Per-library intros | Different for movies vs shows |

### Preroll Configuration

```yaml
cinema:
  preroll:
    enabled: true

    # Global prerolls (all libraries)
    global:
      - path: /prerolls/thx.mp4
        weight: 1
      - path: /prerolls/dolby-atmos.mp4
        weight: 1

    # Seasonal prerolls (date-triggered)
    seasonal:
      - name: halloween
        start: "10-01"
        end: "10-31"
        paths:
          - /prerolls/halloween-1.mp4
          - /prerolls/halloween-2.mp4
      - name: christmas
        start: "12-01"
        end: "12-31"
        paths:
          - /prerolls/christmas.mp4

    # Per-library prerolls
    libraries:
      "Kids Movies":
        paths:
          - /prerolls/disney-intro.mp4
      "Horror":
        paths:
          - /prerolls/horror-warning.mp4

    # Selection mode
    mode: random           # random, sequential, weighted
    max_prerolls: 2        # Maximum prerolls per session
    skip_on_resume: true   # Skip prerolls when resuming playback
```

### Preroll Database Schema

```sql
-- Preroll definitions
CREATE TABLE prerolls (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    path VARCHAR(1024) NOT NULL,
    type VARCHAR(50) NOT NULL,  -- studio_intro, custom_clip, seasonal
    weight INT DEFAULT 1,
    duration_ms INT,

    -- Scheduling
    seasonal_start VARCHAR(5),  -- MM-DD
    seasonal_end VARCHAR(5),

    -- Targeting
    library_ids UUID[],         -- NULL = all libraries
    content_ratings VARCHAR(50)[],  -- G, PG, PG-13, R

    enabled BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Preroll playback log
CREATE TABLE preroll_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    preroll_id UUID NOT NULL REFERENCES prerolls(id),
    media_item_id UUID NOT NULL,
    played_at TIMESTAMPTZ DEFAULT NOW(),
    skipped BOOLEAN DEFAULT false,
    skip_time_ms INT
);
```

---

## Trailer System

### Trailer Selection Modes

| Mode | Description |
|------|-------------|
| `genre_match` | Trailers matching main feature genre |
| `rating_match` | Same content rating or lower |
| `upcoming` | Unreleased movies in library |
| `unwatched` | Movies user hasn't seen |
| `curated` | Admin-selected trailer playlist |
| `random` | Random from trailer pool |

### Trailer Sources

| Source | Priority | Description |
|--------|----------|-------------|
| Local files | 1 | Trailers in library |
| TMDb | 2 | Fetched from TMDb API |
| YouTube | 3 | YouTube trailer links (proxy) |

### Trailer Configuration

```yaml
cinema:
  trailers:
    enabled: true
    count: 2                    # Number of trailers before feature
    max_duration_sec: 180       # Skip trailers longer than 3 min

    selection:
      mode: genre_match         # Primary selection mode
      fallback: random          # Fallback if not enough matches

      # Preferences
      prefer_unwatched: true    # Prioritize unwatched movies
      prefer_upcoming: true     # Prioritize unreleased
      same_rating_only: false   # Allow lower ratings

    # Content filters
    exclude_watched: true       # Don't show trailers for watched movies
    exclude_library_ids: []     # Exclude specific libraries

    # Sources
    sources:
      local: true
      tmdb: true
      youtube: false            # Requires proxy setup
```

### Trailer Database Schema

```sql
-- Local trailer files
CREATE TABLE trailers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    movie_id UUID REFERENCES movies(id),
    path VARCHAR(1024) NOT NULL,
    source VARCHAR(50) NOT NULL,  -- local, tmdb, youtube
    external_id VARCHAR(255),      -- TMDb/YouTube ID
    duration_ms INT,
    resolution VARCHAR(20),
    language VARCHAR(10),
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Trailer playlists (curated)
CREATE TABLE trailer_playlists (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_by UUID REFERENCES users(id),
    is_default BOOLEAN DEFAULT false,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE trailer_playlist_items (
    playlist_id UUID NOT NULL REFERENCES trailer_playlists(id) ON DELETE CASCADE,
    trailer_id UUID NOT NULL REFERENCES trailers(id) ON DELETE CASCADE,
    position INT NOT NULL,
    PRIMARY KEY (playlist_id, trailer_id)
);
```

---

## Intermission System

### Automatic Intermission

Intermissions are automatically suggested for long movies:

| Film Duration | Intermission |
|---------------|--------------|
| < 2h 30m | None |
| 2h 30m - 3h | Optional (user preference) |
| > 3h | Suggested at midpoint |
| > 4h | Suggested (2 intermissions) |

### Intermission Content

| Type | Description |
|------|-------------|
| `countdown` | Animated countdown timer |
| `trivia` | Movie trivia slides |
| `behind_scenes` | BTS clips if available |
| `concession` | "Get snacks" reminder with timer |
| `custom` | User-uploaded intermission video |

### Intermission Configuration

```yaml
cinema:
  intermission:
    enabled: true

    # Automatic triggers
    auto_suggest_above_minutes: 150  # 2.5 hours
    auto_insert_above_minutes: 210   # 3.5 hours

    # Timing
    duration_seconds: 300            # 5 minute intermission
    position: midpoint               # midpoint, chapter, custom

    # Content
    content_type: countdown          # countdown, trivia, custom
    countdown_style: classic         # classic, modern, minimal

    # Custom intermission video
    custom_video_path: /intermissions/default.mp4

    # Audio
    play_ambient_music: true
    ambient_music_path: /intermissions/music/
```

### Intermission Markers

```sql
-- Manual intermission markers (for specific movies)
CREATE TABLE intermission_markers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    movie_id UUID NOT NULL REFERENCES movies(id) ON DELETE CASCADE,
    position_ms BIGINT NOT NULL,     -- Position in movie
    duration_ms INT DEFAULT 300000,  -- Intermission duration
    reason VARCHAR(255),             -- "Original theatrical intermission"
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMPTZ DEFAULT NOW(),

    UNIQUE(movie_id, position_ms)
);
```

### Classic Intermission Movies

Pre-populate intermission markers for classic films:

| Movie | Original Intermission Point |
|-------|----------------------------|
| Lawrence of Arabia | After desert crossing |
| Ben-Hur | Before chariot race |
| Gone with the Wind | After Atlanta burns |
| 2001: A Space Odyssey | After "Jupiter Mission" title |
| The Godfather Part II | Midpoint |

---

## Postroll System

### Postroll Types

| Type | Description |
|------|-------------|
| `credits` | Full credits with chapter skip option |
| `up_next` | Recommendation/next episode |
| `bonus` | Bonus content if available |
| `collection` | Other movies in collection |
| `custom` | Custom end card |

### Postroll Sequence

```
[Movie Ends]
    ↓
[Credits Begin] ← Option to skip
    ↓
[Post-credits Scene] ← Auto-detect & alert
    ↓
[Up Next Card] ← 15 second countdown
    ↓
[Auto-play or Return to Menu]
```

### Post-Credits Scene Detection

```go
// PostCreditsInfo detected from metadata or user contributions
type PostCreditsInfo struct {
    HasPostCredits  bool
    HasMidCredits   bool
    MidCreditsAt    time.Duration  // Position of mid-credits scene
    PostCreditsAt   time.Duration  // Position of post-credits scene
    Description     string         // "Thanos teaser" etc.
}
```

### Postroll Configuration

```yaml
cinema:
  postroll:
    enabled: true

    # Credits handling
    credits:
      auto_skip: false           # Never auto-skip
      show_skip_button: true     # Show skip option
      alert_post_credits: true   # Alert if post-credits scene exists

    # Up Next
    up_next:
      enabled: true
      countdown_seconds: 15
      auto_play: true

      # Selection for movies
      movie_selection: collection  # collection, similar, watchlist

      # Selection for TV
      show_selection: next_episode

    # Bonus content
    bonus:
      show_if_available: true
      types:
        - deleted_scenes
        - bloopers
        - behind_the_scenes
```

---

## Ambient Mode

### Theater Ambiance

Create a theatrical atmosphere between content:

| Feature | Description |
|---------|-------------|
| **Screen Dimming** | Fade to black between content |
| **Curtain Animation** | Virtual curtains open/close |
| **Ambient Sounds** | Theater murmur, popcorn sounds |
| **Lighting Sync** | Hue/LIFX dimming integration |

### Ambient Configuration

```yaml
cinema:
  ambient:
    enabled: true

    # Transitions
    fade_duration_ms: 1500
    curtain_animation: true
    curtain_style: velvet        # velvet, modern, minimal

    # Audio
    ambient_audio: true
    ambient_volume: 0.3
    audio_file: /ambient/theater.mp3

    # Smart home
    lighting:
      enabled: false
      provider: hue              # hue, lifx, home_assistant
      dim_to_percent: 10
      color_temperature: 2700    # Warm white
```

---

## Cinema Session Flow

### Full Session Example

```
┌─────────────────────────────────────────────────────────┐
│                  CINEMA SESSION                         │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  1. [Lights Dim + Curtain Opens]     (5 sec)           │
│                                                         │
│  2. [Preroll: THX Logo]              (30 sec)          │
│                                                         │
│  3. [Preroll: Dolby Atmos]           (20 sec)          │
│                                                         │
│  4. [Trailer 1: Genre Match]         (2 min)           │
│                                                         │
│  5. [Trailer 2: Upcoming]            (2 min)           │
│                                                         │
│  6. [Feature Presentation Card]      (5 sec)           │
│                                                         │
│  7. [MAIN FEATURE]                   (2h 30m)          │
│                                                         │
│  8. [Intermission - Optional]        (5 min)           │
│     └─ Countdown + Trivia                              │
│                                                         │
│  9. [MAIN FEATURE Continues]                           │
│                                                         │
│  10. [Credits Roll]                                    │
│      └─ "Post-credits scene detected!"                 │
│                                                         │
│  11. [Post-Credits Scene]                              │
│                                                         │
│  12. [Up Next Card]                  (15 sec countdown)│
│                                                         │
│  13. [Curtain Closes + Lights Up]                      │
│                                                         │
└─────────────────────────────────────────────────────────┘
```

---

## Implementation

### Cinema Service

```go
// internal/service/cinema/service.go
package cinema

import (
    "context"
    "time"
)

type Service struct {
    db           *pgxpool.Pool
    prerolls     *PrerollService
    trailers     *TrailerService
    intermission *IntermissionService
    config       *CinemaConfig
    logger       *slog.Logger
}

// CinemaSession represents a full cinema experience
type CinemaSession struct {
    ID            uuid.UUID
    UserID        uuid.UUID
    MediaItemID   uuid.UUID

    // Sequence
    Prerolls      []PrerollItem
    Trailers      []TrailerItem
    Intermissions []IntermissionPoint
    PostCredits   *PostCreditsInfo
    UpNext        *UpNextItem

    // State
    CurrentPhase  SessionPhase
    StartedAt     time.Time
}

type SessionPhase string

const (
    PhasePreroll      SessionPhase = "preroll"
    PhaseTrailers     SessionPhase = "trailers"
    PhaseFeature      SessionPhase = "feature"
    PhaseIntermission SessionPhase = "intermission"
    PhaseCredits      SessionPhase = "credits"
    PhasePostCredits  SessionPhase = "post_credits"
    PhaseUpNext       SessionPhase = "up_next"
)

// BuildSession creates a cinema session for a media item
func (s *Service) BuildSession(ctx context.Context, userID, mediaID uuid.UUID) (*CinemaSession, error) {
    session := &CinemaSession{
        ID:          uuid.New(),
        UserID:      userID,
        MediaItemID: mediaID,
        StartedAt:   time.Now(),
    }

    // Get media info
    media, err := s.getMediaInfo(ctx, mediaID)
    if err != nil {
        return nil, fmt.Errorf("get media info: %w", err)
    }

    // Select prerolls
    if s.config.Preroll.Enabled {
        session.Prerolls, err = s.prerolls.Select(ctx, media, s.config.Preroll.MaxPrerolls)
        if err != nil {
            s.logger.Warn("failed to select prerolls", "error", err)
        }
    }

    // Select trailers
    if s.config.Trailers.Enabled {
        session.Trailers, err = s.trailers.Select(ctx, userID, media, s.config.Trailers.Count)
        if err != nil {
            s.logger.Warn("failed to select trailers", "error", err)
        }
    }

    // Check for intermission points
    if s.config.Intermission.Enabled && media.Duration > s.config.Intermission.AutoSuggestAbove {
        session.Intermissions = s.intermission.GetPoints(ctx, media)
    }

    // Check for post-credits
    session.PostCredits = s.getPostCreditsInfo(ctx, mediaID)

    // Get up-next recommendation
    if s.config.Postroll.UpNext.Enabled {
        session.UpNext, _ = s.getUpNext(ctx, userID, media)
    }

    return session, nil
}

// GetSessionPlaylist returns ordered playlist for playback
func (s *Service) GetSessionPlaylist(session *CinemaSession) []PlaylistItem {
    var items []PlaylistItem

    // Add prerolls
    for _, p := range session.Prerolls {
        items = append(items, PlaylistItem{
            Type:     ItemTypePreroll,
            Path:     p.Path,
            Duration: p.Duration,
            Skippable: true,
        })
    }

    // Add trailers
    for _, t := range session.Trailers {
        items = append(items, PlaylistItem{
            Type:      ItemTypeTrailer,
            Path:      t.Path,
            Duration:  t.Duration,
            Skippable: true,
            Metadata:  t.MovieMetadata,
        })
    }

    // Add feature presentation card
    items = append(items, PlaylistItem{
        Type:     ItemTypeCard,
        CardType: "feature_presentation",
        Duration: 5 * time.Second,
    })

    // Main feature is handled separately (not in playlist)

    return items
}
```

### Preroll Selection

```go
// internal/service/cinema/preroll.go
package cinema

type PrerollService struct {
    db     *pgxpool.Pool
    config *PrerollConfig
}

func (s *PrerollService) Select(ctx context.Context, media *MediaInfo, maxCount int) ([]PrerollItem, error) {
    var candidates []PrerollItem

    // Get seasonal prerolls (highest priority)
    seasonal, _ := s.getSeasonalPrerolls(ctx)
    candidates = append(candidates, seasonal...)

    // Get library-specific prerolls
    if media.LibraryID != uuid.Nil {
        libraryPrerolls, _ := s.getLibraryPrerolls(ctx, media.LibraryID)
        candidates = append(candidates, libraryPrerolls...)
    }

    // Get global prerolls
    global, _ := s.getGlobalPrerolls(ctx)
    candidates = append(candidates, global...)

    // Select based on mode
    switch s.config.Mode {
    case "random":
        return s.selectRandom(candidates, maxCount), nil
    case "weighted":
        return s.selectWeighted(candidates, maxCount), nil
    case "sequential":
        return s.selectSequential(candidates, maxCount), nil
    default:
        return s.selectRandom(candidates, maxCount), nil
    }
}

func (s *PrerollService) getSeasonalPrerolls(ctx context.Context) ([]PrerollItem, error) {
    now := time.Now()
    monthDay := fmt.Sprintf("%02d-%02d", now.Month(), now.Day())

    rows, err := s.db.Query(ctx, `
        SELECT id, name, path, duration_ms, weight
        FROM prerolls
        WHERE type = 'seasonal'
          AND enabled = true
          AND seasonal_start <= $1
          AND seasonal_end >= $1
    `, monthDay)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var items []PrerollItem
    for rows.Next() {
        var item PrerollItem
        if err := rows.Scan(&item.ID, &item.Name, &item.Path, &item.Duration, &item.Weight); err != nil {
            continue
        }
        items = append(items, item)
    }

    return items, nil
}
```

### Trailer Selection

```go
// internal/service/cinema/trailers.go
package cinema

type TrailerService struct {
    db     *pgxpool.Pool
    tmdb   *tmdb.Client
    config *TrailerConfig
}

func (s *TrailerService) Select(ctx context.Context, userID uuid.UUID, media *MediaInfo, count int) ([]TrailerItem, error) {
    var trailers []TrailerItem

    switch s.config.Selection.Mode {
    case "genre_match":
        trailers = s.selectByGenre(ctx, userID, media.Genres, count)
    case "upcoming":
        trailers = s.selectUpcoming(ctx, userID, count)
    case "unwatched":
        trailers = s.selectUnwatched(ctx, userID, count)
    case "curated":
        trailers = s.selectFromPlaylist(ctx, count)
    default:
        trailers = s.selectRandom(ctx, userID, count)
    }

    // Apply fallback if not enough trailers
    if len(trailers) < count && s.config.Selection.Fallback != "" {
        remaining := count - len(trailers)
        fallback := s.selectByMode(ctx, userID, s.config.Selection.Fallback, remaining)
        trailers = append(trailers, fallback...)
    }

    return trailers, nil
}

func (s *TrailerService) selectByGenre(ctx context.Context, userID uuid.UUID, genres []string, count int) []TrailerItem {
    rows, err := s.db.Query(ctx, `
        SELECT t.id, t.path, t.duration_ms, m.title, m.year, m.poster_path
        FROM trailers t
        JOIN movies m ON t.movie_id = m.id
        WHERE m.id NOT IN (
            SELECT movie_id FROM movie_watch_history WHERE user_id = $1
        )
        AND m.genres && $2
        AND t.duration_ms <= $3
        ORDER BY RANDOM()
        LIMIT $4
    `, userID, genres, s.config.MaxDurationSec*1000, count)

    if err != nil {
        return nil
    }
    defer rows.Close()

    var items []TrailerItem
    for rows.Next() {
        var item TrailerItem
        rows.Scan(&item.ID, &item.Path, &item.Duration,
            &item.MovieMetadata.Title, &item.MovieMetadata.Year, &item.MovieMetadata.Poster)
        items = append(items, item)
    }

    return items
}
```

---

## API Endpoints

### Cinema Session

```yaml
/api/v1/cinema/session:
  post:
    summary: Create cinema session
    requestBody:
      content:
        application/json:
          schema:
            type: object
            required: [media_id]
            properties:
              media_id:
                type: string
                format: uuid
              include_prerolls:
                type: boolean
                default: true
              include_trailers:
                type: boolean
                default: true
              trailer_count:
                type: integer
                default: 2
    responses:
      200:
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/CinemaSession'

/api/v1/cinema/session/{id}/playlist:
  get:
    summary: Get session playlist
    responses:
      200:
        content:
          application/json:
            schema:
              type: array
              items:
                $ref: '#/components/schemas/PlaylistItem'

/api/v1/cinema/session/{id}/skip:
  post:
    summary: Skip current item
    requestBody:
      content:
        application/json:
          schema:
            type: object
            properties:
              skip_to:
                type: string
                enum: [next, feature, credits]
```

### Preroll Management

```yaml
/api/v1/cinema/prerolls:
  get:
    summary: List all prerolls
  post:
    summary: Upload preroll

/api/v1/cinema/prerolls/{id}:
  put:
    summary: Update preroll
  delete:
    summary: Delete preroll
```

### Intermission

```yaml
/api/v1/cinema/intermission/{movie_id}:
  get:
    summary: Get intermission points for movie
  post:
    summary: Add custom intermission marker
```

---

## User Preferences

### Per-User Cinema Settings

```sql
CREATE TABLE user_cinema_preferences (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,

    -- Prerolls
    prerolls_enabled BOOLEAN DEFAULT true,
    max_prerolls INT DEFAULT 2,

    -- Trailers
    trailers_enabled BOOLEAN DEFAULT true,
    trailer_count INT DEFAULT 2,
    trailer_mode VARCHAR(50) DEFAULT 'genre_match',

    -- Intermission
    auto_intermission BOOLEAN DEFAULT true,
    intermission_threshold_minutes INT DEFAULT 180,

    -- Postroll
    auto_play_next BOOLEAN DEFAULT true,
    up_next_countdown_seconds INT DEFAULT 15,

    -- Ambient
    curtain_animation BOOLEAN DEFAULT true,
    ambient_sounds BOOLEAN DEFAULT false,

    updated_at TIMESTAMPTZ DEFAULT NOW()
);
```

---

## Client Integration

### Client Requirements

| Feature | Web | Mobile | TV | Desktop |
|---------|-----|--------|-----|---------|
| Prerolls | ✅ | ✅ | ✅ | ✅ |
| Trailers | ✅ | ✅ | ✅ | ✅ |
| Intermission | ✅ | ✅ | ✅ | ✅ |
| Curtain Animation | ✅ | ❌ | ✅ | ✅ |
| Lighting Sync | ❌ | ❌ | ✅ | ✅ |
| Ambient Audio | ✅ | ❌ | ✅ | ✅ |

### Client Events

```typescript
// Cinema session events
interface CinemaEvents {
  'cinema:session_start': { sessionId: string };
  'cinema:preroll_start': { prerollId: string };
  'cinema:preroll_skip': { prerollId: string; skipTime: number };
  'cinema:trailer_start': { trailerId: string; movieTitle: string };
  'cinema:trailer_skip': { trailerId: string };
  'cinema:feature_start': { mediaId: string };
  'cinema:intermission_start': { duration: number };
  'cinema:intermission_skip': {};
  'cinema:credits_start': { hasPostCredits: boolean };
  'cinema:post_credits_alert': { position: number };
  'cinema:up_next_countdown': { mediaId: string; seconds: number };
  'cinema:session_end': { sessionId: string };
}
```

---

## Inspiration Sources

### Plex Cinema Trailers
- Pre-roll video support
- Trailer selection from library
- Customizable trailer count

### Emby Cinema Mode
- Theater intros
- Genre-based trailers
- Up-next functionality

### Kodi Experience
- Curtain animations
- Ambient audio
- Lighting integration via add-ons

### Real Theater Experience
- THX/Dolby intros
- Intermissions for epics
- Post-credits scenes

---

## Summary

| Component | Status |
|-----------|--------|
| Preroll system | Documented |
| Trailer selection | Documented |
| Intermission support | Documented |
| Postroll/Up-next | Documented |
| Ambient mode | Documented |
| API endpoints | Documented |
| Database schema | Documented |

Cinema Mode brings the theatrical experience home with customizable prerolls, intelligent trailer selection, intermission support for epics, and seamless up-next transitions.
