## Table of Contents

- [Revenge - Media Enhancement Features](#revenge-media-enhancement-features)
  - [Status](#status)
  - [Architecture](#architecture)
    - [Database Schema](#database-schema)
    - [Module Structure](#module-structure)
    - [Component Interaction](#component-interaction)
  - [Implementation](#implementation)
    - [File Structure](#file-structure)
    - [Key Interfaces](#key-interfaces)
    - [Dependencies](#dependencies)
  - [Configuration](#configuration)
    - [Environment Variables](#environment-variables)
- [Cinema mode](#cinema-mode)
- [Trailers](#trailers)
- [Themes](#themes)
- [YouTube](#youtube)
- [Chapter detection](#chapter-detection)
    - [Config Keys](#config-keys)
  - [API Endpoints](#api-endpoints)
    - [Content Management](#content-management)
- [Cinema Mode](#cinema-mode)
- [Trailers](#trailers)
- [Themes](#themes)
- [Chapters](#chapters)
- [Picture-in-Picture](#picture-in-picture)
  - [Related Documentation](#related-documentation)
    - [Design Documents](#design-documents)
    - [External Sources](#external-sources)

# Revenge - Media Enhancement Features


**Created**: 2026-01-31
**Status**: ✅ Complete
**Category**: feature


> Content module for 

> Advanced playback features: trailers, themes, intros, trickplay, cinema mode, and live TV.

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

```
  ┌─────────────┐     ┌──────────────┐     ┌─────────────┐
  │   Client    │────▶│  API Handler │────▶│   Service   │
  │  (Web/App)  │◀────│   (ogen)     │◀────│   (Logic)   │
  └─────────────┘     └──────────────┘     └──────┬──────┘
                                                   │
                            ┌──────────────────────┼────────────┐
                            ▼                      ▼            ▼
                      ┌──────────┐          ┌───────────┐  ┌────────┐
                      │Repository│          │ Metadata  │  │  Cache │
                      │  (sqlc)  │          │  Service  │  │(otter) │
                      └────┬─────┘          └─────┬─────┘  └────────┘
                           │                      │
                           ▼                      ▼
                    ┌─────────────┐        ┌──────────┐
                    │ PostgreSQL  │        │ External │
                    │   (pgx)     │        │   APIs   │
                    └─────────────┘        └──────────┘
  ```

### Database Schema

**Schema**: `public`

<!-- Schema diagram -->

### Module Structure

```
internal/content/revenge___media_enhancement_features/
├── module.go              # fx module definition
├── repository.go          # Database operations
├── service.go             # Business logic
├── handler.go             # HTTP handlers (ogen)
├── types.go               # Domain types
└── revenge___media_enhancement_features_test.go
```

### Component Interaction

<!-- Component interaction diagram -->


## Implementation

### File Structure

**Key Files**:
- `internal/playback/enhancements/service.go` - Core enhancement logic
- `internal/playback/enhancements/trailers/*.go` - Trailer management
- `internal/playback/enhancements/chapters/detector.go` - Chapter detection
- `web/src/lib/components/player/CinemaMode.svelte` - Cinema mode UI
- `migrations/017_media_enhancements.sql` - Database schema


### Key Interfaces

```go
// EnhancementService manages playback enhancements
type EnhancementService interface {
  // Cinema Mode
  GetCinemaPlaylist(ctx context.Context, movieID uuid.UUID, userID uuid.UUID) (*CinemaPlaylist, error)
  GetCinemaModePreferences(ctx context.Context, userID uuid.UUID) (*CinemaModePreferences, error)
  UpdateCinemaModePreferences(ctx context.Context, userID uuid.UUID, prefs CinemaModePreferences) error

  // Trailers
  GetTrailers(ctx context.Context, movieID uuid.UUID) ([]Trailer, error)
  DownloadTrailer(ctx context.Context, movieID uuid.UUID, source string, sourceID string) (*Trailer, error)
  ScanLocalTrailers(ctx context.Context) (int, error)

  // Themes
  GetTheme(ctx context.Context, seriesID uuid.UUID) (*Theme, error)
  ExtractTheme(ctx context.Context, seriesID uuid.UUID) (*Theme, error)
  FetchFromTVThemes(ctx context.Context, seriesID uuid.UUID) (*Theme, error)

  // Chapters
  GetChapters(ctx context.Context, contentType string, contentID uuid.UUID) ([]Chapter, error)
  DetectChapters(ctx context.Context, contentType string, contentID uuid.UUID) ([]Chapter, error)
  CreateChapter(ctx context.Context, chapter Chapter) error
  UpdateChapter(ctx context.Context, chapterID uuid.UUID, updates ChapterUpdate) error
  DeleteChapter(ctx context.Context, chapterID uuid.UUID) error

  // Picture-in-Picture
  StartPiP(ctx context.Context, session PiPSession) (*PiPSession, error)
  GetActivePiP(ctx context.Context, userID uuid.UUID) (*PiPSession, error)
  EndPiP(ctx context.Context, sessionID uuid.UUID) error
}

// EnhancementRepository handles database operations
type EnhancementRepository interface {
  // Trailers
  CreateTrailer(ctx context.Context, trailer Trailer) error
  GetTrailers(ctx context.Context, movieID uuid.UUID) ([]Trailer, error)
  DeleteTrailer(ctx context.Context, trailerID uuid.UUID) error

  // Themes
  CreateTheme(ctx context.Context, theme Theme) error
  GetTheme(ctx context.Context, seriesID uuid.UUID) (*Theme, error)
  UpdateTheme(ctx context.Context, theme Theme) error

  // Chapters
  CreateChapter(ctx context.Context, chapter Chapter) error
  GetChapters(ctx context.Context, contentType string, contentID uuid.UUID) ([]Chapter, error)
  UpdateChapter(ctx context.Context, chapter Chapter) error
  DeleteChapter(ctx context.Context, chapterID uuid.UUID) error

  // Cinema Mode
  GetCinemaModePreferences(ctx context.Context, userID uuid.UUID) (*CinemaModePreferences, error)
  UpsertCinemaModePreferences(ctx context.Context, prefs CinemaModePreferences) error

  // PiP
  CreatePiPSession(ctx context.Context, session PiPSession) error
  GetActivePiPSession(ctx context.Context, userID uuid.UUID) (*PiPSession, error)
  EndPiPSession(ctx context.Context, sessionID uuid.UUID) error
}

// TrailerService downloads and manages trailers
type TrailerService interface {
  // DownloadFromYouTube downloads trailer from YouTube
  DownloadFromYouTube(ctx context.Context, youtubeID string, destPath string) error

  // FetchFromTMDb fetches trailer metadata from TMDb
  FetchFromTMDb(ctx context.Context, tmdbID int) ([]TrailerMetadata, error)

  // ScanLocal scans for local trailer files
  ScanLocal(ctx context.Context, moviePath string) ([]string, error)
}

// ChapterDetector auto-detects chapter markers
type ChapterDetector interface {
  // DetectBySilence detects chapters by silence gaps
  DetectBySilence(ctx context.Context, filePath string, threshold float64) ([]Chapter, error)

  // DetectByBlackFrames detects chapters by black frames
  DetectByBlackFrames(ctx context.Context, filePath string) ([]Chapter, error)

  // DetectIntroOutro detects intro/outro segments
  DetectIntroOutro(ctx context.Context, contentType string, contentID uuid.UUID) (*Chapter, *Chapter, error)
}

// Types
type Trailer struct {
  ID              uuid.UUID       `db:"id" json:"id"`
  MovieID         uuid.UUID       `db:"movie_id" json:"movie_id"`
  Title           string          `db:"title" json:"title"`
  FilePath        string          `db:"file_path" json:"file_path"`
  DurationSeconds *int            `db:"duration_seconds" json:"duration_seconds,omitempty"`
  Resolution      *string         `db:"resolution" json:"resolution,omitempty"`
  Source          string          `db:"source" json:"source"`
  SourceID        *string         `db:"source_id" json:"source_id,omitempty"`
  TrailerType     string          `db:"trailer_type" json:"trailer_type"`
  Language        string          `db:"language" json:"language"`
  ReleaseDate     *time.Time      `db:"release_date" json:"release_date,omitempty"`
  Priority        int             `db:"priority" json:"priority"`
  CreatedAt       time.Time       `db:"created_at" json:"created_at"`
  UpdatedAt       time.Time       `db:"updated_at" json:"updated_at"`
}

type Theme struct {
  ID                uuid.UUID       `db:"id" json:"id"`
  SeriesID          uuid.UUID       `db:"series_id" json:"series_id"`
  FilePath          string          `db:"file_path" json:"file_path"`
  IntroStartSeconds *float64        `db:"intro_start_seconds" json:"intro_start_seconds,omitempty"`
  IntroEndSeconds   *float64        `db:"intro_end_seconds" json:"intro_end_seconds,omitempty"`
  DurationSeconds   *int            `db:"duration_seconds" json:"duration_seconds,omitempty"`
  Source            string          `db:"source" json:"source"`
  CreatedAt         time.Time       `db:"created_at" json:"created_at"`
}

type Chapter struct {
  ID                uuid.UUID       `db:"id" json:"id"`
  ContentType       string          `db:"content_type" json:"content_type"`
  ContentID         uuid.UUID       `db:"content_id" json:"content_id"`
  StartTimeSeconds  float64         `db:"start_time_seconds" json:"start_time_seconds"`
  EndTimeSeconds    float64         `db:"end_time_seconds" json:"end_time_seconds"`
  Title             string          `db:"title" json:"title"`
  ChapterType       *string         `db:"chapter_type" json:"chapter_type,omitempty"`
  ConfidenceScore   *float64        `db:"confidence_score" json:"confidence_score,omitempty"`
  DetectedBy        *string         `db:"detected_by" json:"detected_by,omitempty"`
  CreatedAt         time.Time       `db:"created_at" json:"created_at"`
  UpdatedAt         time.Time       `db:"updated_at" json:"updated_at"`
}

type CinemaModePreferences struct {
  UserID        uuid.UUID       `db:"user_id" json:"user_id"`
  Enabled       bool            `db:"enabled" json:"enabled"`
  PlayTrailers  bool            `db:"play_trailers" json:"play_trailers"`
  MaxTrailers   int             `db:"max_trailers" json:"max_trailers"`
  TrailerTypes  []string        `db:"trailer_types" json:"trailer_types"`
  PlayTheme     bool            `db:"play_theme" json:"play_theme"`
  DimLights     bool            `db:"dim_lights" json:"dim_lights"`
  CreatedAt     time.Time       `db:"created_at" json:"created_at"`
  UpdatedAt     time.Time       `db:"updated_at" json:"updated_at"`
}

type PiPSession struct {
  ID                     uuid.UUID       `db:"id" json:"id"`
  UserID                 uuid.UUID       `db:"user_id" json:"user_id"`
  PrimaryContentType     string          `db:"primary_content_type" json:"primary_content_type"`
  PrimaryContentID       uuid.UUID       `db:"primary_content_id" json:"primary_content_id"`
  PiPContentType         string          `db:"pip_content_type" json:"pip_content_type"`
  PiPContentID           uuid.UUID       `db:"pip_content_id" json:"pip_content_id"`
  PrimaryPositionSeconds *float64        `db:"primary_position_seconds" json:"primary_position_seconds,omitempty"`
  PiPPositionSeconds     *float64        `db:"pip_position_seconds" json:"pip_position_seconds,omitempty"`
  PiPSize                string          `db:"pip_size" json:"pip_size"`
  PiPPosition            string          `db:"pip_position" json:"pip_position"`
  StartedAt              time.Time       `db:"started_at" json:"started_at"`
  EndedAt                *time.Time      `db:"ended_at" json:"ended_at,omitempty"`
}

type CinemaPlaylist struct {
  Trailers []Trailer       `json:"trailers"`
  Movie    *Movie          `json:"movie"`
  Theme    *Theme          `json:"theme,omitempty"`
}
```


### Dependencies
**Go Packages**:
- `github.com/google/uuid` - UUID handling
- `github.com/jackc/pgx/v5` - PostgreSQL driver
- `github.com/maypok86/otter` - L1 in-memory cache
- `github.com/riverqueue/river` - Background job queue
- `go.uber.org/fx` - Dependency injection
- `go.uber.org/zap` - Structured logging
- `github.com/asticode/go-astiav` - FFmpeg bindings for detection
- `github.com/kkdai/youtube/v2` - YouTube video downloader

**Frontend Packages**:
- `@sveltejs/kit` - SvelteKit framework
- `svelte` - Svelte 5 with runes
- `hls.js` - HLS playback
- `video.js` - Video player

**External APIs**:
- TMDb API - Trailer metadata
- YouTube API - Trailer downloads
- TV Themes database (optional)







## Configuration

### Environment Variables

```bash
# Cinema mode
CINEMA_MODE_ENABLED=true                     # Enable cinema mode globally
CINEMA_MODE_DEFAULT_TRAILERS=2               # Default max trailers

# Trailers
TRAILERS_PATH=/path/to/trailers              # Storage path for trailers
TRAILERS_AUTO_DOWNLOAD=false                 # Auto-download trailers for new movies

# Themes
THEMES_PATH=/path/to/themes                  # Storage path for theme songs
THEMES_AUTO_EXTRACT=true                     # Auto-extract themes from episodes

# YouTube
YOUTUBE_ENABLED=true                         # Enable YouTube trailer downloads
YOUTUBE_QUALITY=1080p                        # Preferred quality

# Chapter detection
CHAPTERS_AUTO_DETECT=true                    # Auto-detect chapters
CHAPTERS_SILENCE_THRESHOLD=-30dB             # Silence detection threshold
```


### Config Keys
```yaml
enhancements:
  cinema_mode:
    enabled: true
    default_max_trailers: 2
    default_trailer_types:
      - trailer
      - teaser
    auto_download_trailers: false

  trailers:
    storage_path: /path/to/trailers
    youtube:
      enabled: true
      quality: 1080p
      format: mp4
    tmdb:
      enabled: true
      preferred_language: en

  themes:
    storage_path: /path/to/themes
    auto_extract: true
    extraction:
      min_episodes: 3                # Need at least 3 episodes to detect pattern
      confidence_threshold: 0.85     # Minimum confidence for auto-detection

  chapters:
    auto_detect: true
    detection:
      silence_threshold: -30         # dB
      silence_duration: 1.0          # seconds
      black_frame_threshold: 32      # pixel value (0-255)
      black_frame_duration: 2.0      # seconds

  pip:
    enabled: true
    default_size: small
    default_position: bottom-right
    allowed_combinations:
      - primary: live_tv
        pip: [movie, episode, live_tv]
      - primary: movie
        pip: [live_tv]
```



## API Endpoints

### Content Management
**Endpoints**:
```
# Cinema Mode
GET    /api/v1/users/:id/cinema-mode                  # Get user cinema preferences
PUT    /api/v1/users/:id/cinema-mode                  # Update cinema preferences
GET    /api/v1/movies/:id/cinema-playlist             # Get cinema mode playlist

# Trailers
GET    /api/v1/movies/:id/trailers                    # Get trailers for movie
POST   /api/v1/movies/:id/trailers                    # Download/add trailer
DELETE /api/v1/trailers/:id                           # Delete trailer

# Themes
GET    /api/v1/series/:id/theme                       # Get theme for series
POST   /api/v1/series/:id/theme/extract               # Extract theme from episodes
DELETE /api/v1/series/:id/theme                       # Delete theme

# Chapters
GET    /api/v1/:type/:id/chapters                     # Get chapters
POST   /api/v1/:type/:id/chapters                     # Create chapter
POST   /api/v1/:type/:id/chapters/detect              # Auto-detect chapters
PUT    /api/v1/chapters/:id                           # Update chapter
DELETE /api/v1/chapters/:id                           # Delete chapter

# Picture-in-Picture
POST   /api/v1/playback/pip                           # Start PiP session
GET    /api/v1/playback/pip                           # Get active PiP
DELETE /api/v1/playback/pip/:id                       # End PiP session
```

**Request/Response Examples**:

**Get Cinema Playlist**:
```http
GET /api/v1/movies/550e8400-e29b-41d4-a716-446655440000/cinema-playlist

Response 200 OK:
{
  "trailers": [
    {
      "id": "660e8400-e29b-41d4-a716-446655440001",
      "title": "Official Trailer",
      "file_path": "/trailers/movie-trailer-1.mp4",
      "duration_seconds": 150,
      "resolution": "1080p",
      "source": "youtube",
      "trailer_type": "trailer"
    }
  ],
  "movie": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "title": "Example Movie",
    "file_path": "/movies/example.mp4"
  }
}
```

**Extract Theme**:
```http
POST /api/v1/series/770e8400-e29b-41d4-a716-446655440002/theme/extract

Response 202 Accepted:
{
  "message": "Theme extraction job queued",
  "job_id": "880e8400-e29b-41d4-a716-446655440003"
}
```

**Get Chapters**:
```http
GET /api/v1/episode/990e8400-e29b-41d4-a716-446655440004/chapters

Response 200 OK:
{
  "chapters": [
    {
      "id": "aa0e8400-e29b-41d4-a716-446655440005",
      "content_type": "episode",
      "content_id": "990e8400-e29b-41d4-a716-446655440004",
      "start_time_seconds": 0,
      "end_time_seconds": 90,
      "title": "Recap",
      "chapter_type": "recap",
      "confidence_score": 0.92,
      "detected_by": "chromaprint"
    },
    {
      "id": "bb0e8400-e29b-41d4-a716-446655440006",
      "start_time_seconds": 90,
      "end_time_seconds": 150,
      "title": "Intro",
      "chapter_type": "intro",
      "confidence_score": 0.95,
      "detected_by": "chromaprint"
    }
  ]
}
```

**Start PiP**:
```http
POST /api/v1/playback/pip
{
  "user_id": "cc0e8400-e29b-41d4-a716-446655440007",
  "primary_content_type": "live_tv",
  "primary_content_id": "dd0e8400-e29b-41d4-a716-446655440008",
  "pip_content_type": "movie",
  "pip_content_id": "550e8400-e29b-41d4-a716-446655440000",
  "pip_size": "small",
  "pip_position": "bottom-right"
}

Response 201 Created:
{
  "id": "ee0e8400-e29b-41d4-a716-446655440009",
  "user_id": "cc0e8400-e29b-41d4-a716-446655440007",
  "primary_content_type": "live_tv",
  "pip_content_type": "movie",
  "pip_size": "small",
  "pip_position": "bottom-right",
  "started_at": "2026-01-31T15:30:00Z"
}
```








## Related Documentation
### Design Documents
- [01_ARCHITECTURE](../../architecture/01_ARCHITECTURE.md)
- [02_DESIGN_PRINCIPLES](../../architecture/02_DESIGN_PRINCIPLES.md)
- [03_METADATA_SYSTEM](../../architecture/03_METADATA_SYSTEM.md)

### External Sources
- [FFmpeg Documentation](../../../sources/media/ffmpeg.md) - Auto-resolved from ffmpeg
- [FFmpeg Codecs](../../../sources/media/ffmpeg-codecs.md) - Auto-resolved from ffmpeg-codecs
- [FFmpeg Formats](../../../sources/media/ffmpeg-formats.md) - Auto-resolved from ffmpeg-formats
- [go-astiav (FFmpeg bindings)](../../../sources/media/go-astiav.md) - Auto-resolved from go-astiav
- [go-astiav GitHub README](../../../sources/media/go-astiav-guide.md) - Auto-resolved from go-astiav-docs
- [M3U8 Extended Format](../../../sources/protocols/m3u8.md) - Auto-resolved from m3u8
- [River Job Queue](../../../sources/tooling/river.md) - Auto-resolved from river
- [Svelte 5 Runes](../../../sources/frontend/svelte-runes.md) - Auto-resolved from svelte-runes
- [Svelte 5 Documentation](../../../sources/frontend/svelte5.md) - Auto-resolved from svelte5
- [SvelteKit Documentation](../../../sources/frontend/sveltekit.md) - Auto-resolved from sveltekit

