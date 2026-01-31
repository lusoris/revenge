# Revenge - Media Enhancement Features

> Advanced playback features: trailers, themes, intros, trickplay, cinema mode, and live TV.


<!-- TOC-START -->

## Table of Contents

- [Status](#status)
- [Developer Resources](#developer-resources)
- [Feature Overview](#feature-overview)
- [Cinema Mode (Theatrical Experience)](#cinema-mode-theatrical-experience)
  - [Overview](#overview)
  - [Preroll System](#preroll-system)
  - [Intermission System](#intermission-system)
  - [Postroll / Up Next](#postroll-up-next)
  - [Cinema Mode Configuration](#cinema-mode-configuration)
- [Sources & Cross-References](#sources-cross-references)
  - [Cross-Reference Indexes](#cross-reference-indexes)
  - [Referenced Sources](#referenced-sources)
- [Related Design Docs](#related-design-docs)
  - [In This Section](#in-this-section)
  - [Related Topics](#related-topics)
  - [Indexes](#indexes)
- [Summary](#summary)

<!-- TOC-END -->

## Status

| Dimension | Status |
|-----------|--------|
| Design | ✅ |
| Sources | ✅ |
| Instructions | ✅ |
| Code | 🔴 |
| Linting | 🔴 |
| Unit Testing | 🔴 |
| Integration Testing | 🔴 |
**Location**: `internal/feature/enhancements/`

---

## Developer Resources

| Source       | URL                                                                            | Purpose                 |
| ------------ | ------------------------------------------------------------------------------ | ----------------------- |
| FFmpeg       | [ffmpeg.org/documentation.html](https://ffmpeg.org/documentation.html)         | Video/audio processing  |
| HDR10+       | [hdr10plus.org](https://hdr10plus.org/)                                        | HDR metadata            |
| Dolby Vision | [developer.dolby.com](https://developer.dolby.com/)                            | Dolby Vision support    |
| LPCM Wiki    | [wiki.multimedia.cx/index.php/PCM](https://wiki.multimedia.cx/index.php/PCM)   | Audio formats reference |

---

## Feature Overview

| Feature | Description | Status |
|---------|-------------|--------|
| **Trailers** | Movie/show trailers before playback | 🔶 Planned |
| **Audio Themes** | Theme music on item hover/selection | 🔶 Planned |
| **Intro/Outro Detection** | Netflix-style skip buttons | 🔶 Planned |
| **Trickplay** | Video scrubbing thumbnails | 🔶 Planned |
| **Chapters** | Chapter markers with images | 🔶 Planned |
| **Cinema Mode** | Prerolls, intermissions, postrolls | 🔶 Planned |
| **Live TV** | EPG, DVR, timeshift | 🔶 Planned |
| **PVR Integration** | External PVR server support | 🔶 Planned |

---

## Cinema Mode (Theatrical Experience)

### Overview

Cinema Mode transforms home viewing into a theatrical experience with customizable sequences.

```
┌─────────────────────────────────────────────────────────────────┐
│                    CINEMA SESSION FLOW                          │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  [Preroll] → [Trailers] → [Feature] → [Intermission?] →        │
│  [Credits] → [Post-Credits Alert] → [Up Next]                  │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### Preroll System

| Type | Description | Example |
|------|-------------|---------|
| `studio_intro` | Studio/network fanfare | THX, Dolby Atmos |
| `custom_clip` | User-uploaded video | Personal intro |
| `seasonal` | Date-based content | Halloween, Christmas |
| `library_specific` | Per-library intros | Different for Kids vs Horror |

```sql
CREATE TABLE prerolls (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    path VARCHAR(1024) NOT NULL,
    type VARCHAR(50) NOT NULL,
    weight INT DEFAULT 1,
    duration_ms INT,
    seasonal_start VARCHAR(5),  -- MM-DD
    seasonal_end VARCHAR(5),
    library_ids UUID[],
    enabled BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

### Intermission System

Auto-suggest for long movies:

| Duration | Intermission |
|----------|--------------|
| < 2h 30m | None |
| 2h 30m - 3h | Optional |
| > 3h | Suggested at midpoint |

```sql
CREATE TABLE intermission_markers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    movie_id UUID NOT NULL REFERENCES movies(id) ON DELETE CASCADE,
    position_ms BIGINT NOT NULL,
    duration_ms INT DEFAULT 300000,  -- 5 min default
    reason VARCHAR(255),
    UNIQUE(movie_id, position_ms)
);
```

### Postroll / Up Next

| Type | Description |
|------|-------------|
| `credits` | Full credits with skip option |
| `post_credits_alert` | Alert for post-credits scenes |
| `up_next` | Next episode/similar movie countdown |
| `collection` | Other movies in collection |

### Cinema Mode Configuration

```yaml
cinema:
  preroll:
    enabled: true
    max_prerolls: 2
    mode: random  # random, sequential, weighted

  intermission:
    enabled: true
    auto_suggest_above_minutes: 150
    duration_seconds: 300

  postroll:
    up_next:
      enabled: true
      countdown_seconds: 15
      auto_play: true

---

## Trailer System

### Architecture

```
┌─────────────────────────────────────────────────────────────────────────┐
│                         Trailer Flow                                     │
└─────────────────────────────────────────────────────────────────────────┘

┌──────────┐     ┌───────────────┐     ┌───────────────┐     ┌──────────┐
│  Play    │ ──→ │ Trailer Queue │ ──→ │ Play Trailer  │ ──→ │  Main    │
│  Movie   │     │  (Optional)   │     │ (Skippable)   │     │ Content  │
└──────────┘     └───────────────┘     └───────────────┘     └──────────┘
                       │
          ┌────────────┼────────────┐
          ▼            ▼            ▼
    ┌──────────┐ ┌──────────┐ ┌──────────┐
    │  Local   │ │  Servarr │ │  YouTube │
    │ Trailers │ │ Trailers │ │   API    │
    └──────────┘ └──────────┘ └──────────┘
```

### Data Model

```sql
-- Trailer storage
CREATE TABLE movie_trailers (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    movie_id        UUID REFERENCES movies(id) ON DELETE CASCADE,
    name            VARCHAR(255) NOT NULL,
    type            VARCHAR(50) NOT NULL,  -- 'trailer', 'teaser', 'featurette', 'clip'
    language        VARCHAR(10) NOT NULL,
    source          VARCHAR(50) NOT NULL,  -- 'local', 'youtube', 'tmdb'
    source_id       VARCHAR(100),          -- YouTube ID, TMDb video ID
    local_path      TEXT,                  -- Path if locally stored
    thumbnail_url   TEXT,
    duration_seconds INT,
    is_primary      BOOLEAN DEFAULT false,
    created_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_movie_trailers_movie ON movie_trailers(movie_id);

-- Same for TV shows
CREATE TABLE series_trailers (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    series_id       UUID REFERENCES series(id) ON DELETE CASCADE,
    season_number   INT,                   -- NULL = series trailer
    episode_number  INT,                   -- NULL = season trailer
    name            VARCHAR(255) NOT NULL,
    type            VARCHAR(50) NOT NULL,
    language        VARCHAR(10) NOT NULL,
    source          VARCHAR(50) NOT NULL,
    source_id       VARCHAR(100),
    local_path      TEXT,
    thumbnail_url   TEXT,
    duration_seconds INT,
    is_primary      BOOLEAN DEFAULT false,
    created_at      TIMESTAMPTZ DEFAULT NOW()
);
```

### Trailer Sources

```go
// TrailerService manages trailer discovery and playback
type TrailerService struct {
    movieRepo    MovieRepository
    trailerRepo  TrailerRepository
    tmdbClient   *TMDbClient
    youtubeAPI   *YouTubeClient
    radarrClient *RadarrClient
    logger       *slog.Logger
}

// TrailerSource priorities
// 1. Local files (in movie folder: *-trailer.*, trailer.*)
// 2. Radarr metadata (has YouTube IDs)
// 3. TMDb API (official trailers)
// 4. YouTube search (fallback)

func (s *TrailerService) GetTrailers(ctx context.Context, movieID uuid.UUID) ([]Trailer, error) {
    // Check cached
    trailers, err := s.trailerRepo.GetByMovie(ctx, movieID)
    if err == nil && len(trailers) > 0 {
        return trailers, nil
    }

    movie, err := s.movieRepo.GetByID(ctx, movieID)
    if err != nil {
        return nil, err
    }

    // Try sources in order
    trailers = s.findLocalTrailers(movie.Path)
    if len(trailers) == 0 {
        trailers = s.getFromRadarr(ctx, movie.RadarrID)
    }
    if len(trailers) == 0 {
        trailers = s.getFromTMDb(ctx, movie.TmdbID)
    }

    // Cache results
    for _, t := range trailers {
        s.trailerRepo.Create(ctx, movieID, &t)
    }

    return trailers, nil
}

// Local trailer file patterns
var trailerPatterns = []string{
    "*-trailer.*",
    "*-trailers/*",
    "trailer.*",
    "trailers/*",
    "*-teaser.*",
}
```

### User Preferences

```go
type TrailerPreferences struct {
    EnableTrailers     bool `json:"enable_trailers"`
    MaxTrailerCount    int  `json:"max_trailer_count"`    // 0-5
    PreferredLanguage  string `json:"preferred_language"` // ISO 639-1
    ShowRelatedTrailers bool `json:"show_related_trailers"` // Upcoming movies
    SkipAfterSeconds   int  `json:"skip_after_seconds"`   // Auto-show skip button
}
```

---

## Audio Theme System

### Overview

Theme music plays when hovering over or selecting a movie/show, similar to Netflix/Disney+.

```
┌─────────────────────────────────────────────────────────────────────────┐
│                      Audio Theme Flow                                    │
└─────────────────────────────────────────────────────────────────────────┘

   User hovers/selects item
            │
            ▼
   ┌─────────────────┐
   │ Check for theme │
   └─────────────────┘
            │
    ┌───────┴───────┐
    ▼               ▼
┌────────┐    ┌────────────┐
│ Local  │    │ Extract    │
│ Theme  │    │ From Media │
└────────┘    └────────────┘
    │               │
    └───────┬───────┘
            ▼
   ┌─────────────────┐
   │ Fade In Audio   │
   │ (2-3 sec)       │
   └─────────────────┘
            │
            ▼
   ┌─────────────────┐
   │ Play 15-30 sec  │
   │ (configurable)  │
   └─────────────────┘
            │
            ▼
   ┌─────────────────┐
   │ Fade Out        │
   │ (on mouse leave)│
   └─────────────────┘
```

### Data Model

```sql
CREATE TABLE media_themes (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    item_type       VARCHAR(50) NOT NULL,  -- 'movie', 'series'
    item_id         UUID NOT NULL,
    source          VARCHAR(50) NOT NULL,  -- 'local', 'extracted', 'spotify'
    local_path      TEXT,
    start_time      FLOAT DEFAULT 0,       -- Start position in seconds
    duration        FLOAT DEFAULT 30,      -- Theme duration
    volume          FLOAT DEFAULT 0.3,     -- Relative volume (0-1)
    created_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_media_themes_item ON media_themes(item_type, item_id);
```

### Implementation

```go
type ThemeService struct {
    themeRepo  ThemeRepository
    extractor  *AudioExtractor
    cache      *cache.Client
    logger     *slog.Logger
}

// GetTheme returns theme audio for an item
func (s *ThemeService) GetTheme(ctx context.Context, itemType string, itemID uuid.UUID) (*Theme, error) {
    // Check cache
    cacheKey := fmt.Sprintf("theme:%s:%s", itemType, itemID)
    if cached, ok := s.cache.Get(ctx, cacheKey); ok {
        return cached.(*Theme), nil
    }

    // Check database
    theme, err := s.themeRepo.Get(ctx, itemType, itemID)
    if err == nil {
        s.cache.Set(ctx, cacheKey, theme, 1*time.Hour)
        return theme, nil
    }

    return nil, ErrNoTheme
}

// ExtractTheme extracts theme from media file
func (s *ThemeService) ExtractTheme(ctx context.Context, mediaPath string) (*ThemeExtract, error) {
    // Use FFmpeg to extract audio segment
    // Analyze for "interesting" segment (not silence, has energy)
    // Return extracted theme metadata
}

// Theme file patterns
var themePatterns = []string{
    "theme.*",
    "theme-music.*",
    "*-theme.*",
    "themes/*",
}
```

### Frontend Integration

```typescript
// Svelte component for theme audio
class ThemeAudioManager {
    private audio: HTMLAudioElement | null = null;
    private fadeTimeout: number | null = null;
    private currentItemId: string | null = null;

    async playTheme(itemId: string, itemType: 'movie' | 'series') {
        if (this.currentItemId === itemId) return;

        await this.stopTheme(); // Fade out current

        const theme = await fetchTheme(itemId, itemType);
        if (!theme) return;

        this.audio = new Audio(theme.url);
        this.audio.volume = 0;
        this.audio.currentTime = theme.startTime;
        this.audio.play();

        // Fade in
        this.fadeVolume(0, theme.volume, 2000);
        this.currentItemId = itemId;

        // Auto-stop after duration
        this.fadeTimeout = setTimeout(() => {
            this.fadeVolume(this.audio!.volume, 0, 2000);
        }, theme.duration * 1000);
    }

    async stopTheme() {
        if (!this.audio) return;

        if (this.fadeTimeout) clearTimeout(this.fadeTimeout);
        await this.fadeVolume(this.audio.volume, 0, 500);

        this.audio.pause();
        this.audio = null;
        this.currentItemId = null;
    }

    private fadeVolume(from: number, to: number, duration: number): Promise<void> {
        return new Promise(resolve => {
            const steps = 20;
            const stepTime = duration / steps;
            const volumeStep = (to - from) / steps;
            let currentStep = 0;

            const interval = setInterval(() => {
                currentStep++;
                this.audio!.volume = Math.max(0, Math.min(1, from + volumeStep * currentStep));

                if (currentStep >= steps) {
                    clearInterval(interval);
                    resolve();
                }
            }, stepTime);
        });
    }
}
```

---

## Intro/Outro Detection (Netflix-Style)

### Overview

Automatically detect and mark intros, outros, and credits for skipping.

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    Episode Structure                                     │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│  ┌───────┐ ┌─────────────────────────────────┐ ┌───────┐ ┌───────────┐ │
│  │ Recap │ │            Main Content         │ │Credits│ │Next Ep    │ │
│  │ 0-90s │ │                                 │ │ ~60s  │ │Preview    │ │
│  └───────┘ └─────────────────────────────────┘ └───────┘ └───────────┘ │
│      ▲                    ▲                        ▲           ▲        │
│   Skip Recap          Skip Intro              Skip Credits  Auto-Play   │
│                                                                          │
└─────────────────────────────────────────────────────────────────────────┘
```

### Data Model

```sql
CREATE TABLE episode_segments (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    episode_id      UUID REFERENCES episodes(id) ON DELETE CASCADE,
    segment_type    VARCHAR(50) NOT NULL,  -- 'intro', 'outro', 'recap', 'credits', 'preview'
    start_time      FLOAT NOT NULL,        -- Seconds from start
    end_time        FLOAT NOT NULL,        -- Seconds from start
    confidence      FLOAT DEFAULT 1.0,     -- Detection confidence (0-1)
    source          VARCHAR(50) NOT NULL,  -- 'manual', 'detected', 'chromaprint', 'imported'
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_episode_segments_episode ON episode_segments(episode_id);
CREATE INDEX idx_episode_segments_type ON episode_segments(segment_type);
```

### Detection Methods

```go
// SegmentDetector detects intros, outros, and other segments
type SegmentDetector struct {
    chromaprint *ChromaprintClient  // Audio fingerprinting
    blackbeard  *TranscoderClient   // For video analysis
    chapterRepo ChapterRepository
    logger      *slog.Logger
}

// DetectionMethod priorities:
// 1. Imported from external tools (e.g., intro-skipper plugin data)
// 2. Chromaprint audio fingerprinting (compare episodes)
// 3. Video analysis (scene detection, black frames)
// 4. Manual user input

func (d *SegmentDetector) DetectIntro(ctx context.Context, seriesID uuid.UUID, seasonNum int) error {
    episodes, err := d.getSeasonEpisodes(ctx, seriesID, seasonNum)
    if err != nil || len(episodes) < 2 {
        return err
    }

    // Extract audio fingerprints from first 5 minutes of each episode
    fingerprints := make([][]byte, len(episodes))
    for i, ep := range episodes {
        fp, err := d.chromaprint.Fingerprint(ctx, ep.Path, 0, 300) // 5 min
        if err != nil {
            d.logger.Warn("fingerprint failed", "episode", ep.ID, "error", err)
            continue
        }
        fingerprints[i] = fp
    }

    // Find common audio sequence (intro)
    intro := d.findCommonSequence(fingerprints)
    if intro != nil {
        for _, ep := range episodes {
            d.saveSegment(ctx, ep.ID, "intro", intro.Start, intro.End, intro.Confidence, "chromaprint")
        }
    }

    return nil
}

// Chromaprint integration for audio fingerprinting
type ChromaprintClient struct {
    binaryPath string  // Path to fpcalc binary
}

func (c *ChromaprintClient) Fingerprint(ctx context.Context, path string, start, duration float64) ([]byte, error) {
    args := []string{
        "-raw", "-json",
        "-offset", fmt.Sprintf("%.1f", start),
        "-length", fmt.Sprintf("%.1f", duration),
        path,
    }

    cmd := exec.CommandContext(ctx, c.binaryPath, args...)
    output, err := cmd.Output()
    if err != nil {
        return nil, fmt.Errorf("chromaprint: %w", err)
    }

    var result struct {
        Fingerprint []int `json:"fingerprint"`
    }
    json.Unmarshal(output, &result)

    // Convert to bytes for storage/comparison
    return encodeFingerprint(result.Fingerprint), nil
}
```

### External Tool Import

```go
// Import from Jellyfin intro-skipper plugin format
type IntroSkipperImport struct {
    EpisodeID string  `json:"EpisodeId"`
    Valid     bool    `json:"Valid"`
    IntroStart float64 `json:"IntroStart"`
    IntroEnd   float64 `json:"IntroEnd"`
}

func (d *SegmentDetector) ImportIntroSkipperData(ctx context.Context, dataPath string) error {
    // Read intro-skipper JSON files
    files, _ := filepath.Glob(filepath.Join(dataPath, "*.json"))

    for _, file := range files {
        data, _ := os.ReadFile(file)
        var intro IntroSkipperImport
        json.Unmarshal(data, &intro)

        if intro.Valid {
            // Map EpisodeID to our episode
            episode, err := d.findEpisodeByPath(ctx, intro.EpisodeID)
            if err == nil {
                d.saveSegment(ctx, episode.ID, "intro", intro.IntroStart, intro.IntroEnd, 1.0, "imported")
            }
        }
    }
    return nil
}
```

### Playback Integration

```go
// SegmentHandler provides segment info during playback
type SegmentHandler struct {
    segmentRepo SegmentRepository
}

// GetSegments returns segments for an episode
func (h *SegmentHandler) GetSegments(ctx context.Context, episodeID uuid.UUID) ([]Segment, error) {
    return h.segmentRepo.GetByEpisode(ctx, episodeID)
}

// API Response
type PlaybackSegments struct {
    EpisodeID  uuid.UUID `json:"episode_id"`
    Intro      *Segment  `json:"intro,omitempty"`
    Outro      *Segment  `json:"outro,omitempty"`
    Recap      *Segment  `json:"recap,omitempty"`
    Credits    *Segment  `json:"credits,omitempty"`
    NextEpisode *Episode `json:"next_episode,omitempty"`
}
```

### Frontend Skip UI

```typescript
// Skip button component
interface SkipButtonProps {
    segment: Segment;
    currentTime: number;
    onSkip: () => void;
}

function SkipButton({ segment, currentTime, onSkip }: SkipButtonProps) {
    const isVisible = currentTime >= segment.startTime &&
                      currentTime < segment.endTime - 5; // Hide 5s before end

    if (!isVisible) return null;

    const label = {
        'intro': 'Skip Intro',
        'recap': 'Skip Recap',
        'credits': 'Skip Credits',
        'outro': 'Next Episode'
    }[segment.type];

    return (
        <button
            class="skip-button"
            onClick={() => {
                if (segment.type === 'credits' && nextEpisode) {
                    playNextEpisode();
                } else {
                    onSkip();
                }
            }}
        >
            {label}
        </button>
    );
}
```

---

## Trickplay (Video Scrubbing Thumbnails)

### Overview

Generate thumbnail strips for video scrubbing, showing preview images during seek.

```
┌─────────────────────────────────────────────────────────────────────────┐
│                         Trickplay Preview                                │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│     ┌───────────────────┐                                               │
│     │   [Thumbnail]     │  ← Shows on hover/drag                        │
│     │   00:42:15        │                                               │
│     └───────────────────┘                                               │
│                 ▼                                                        │
│  ═══════════════●════════════════════════════════════════════════       │
│  0:00                                                          1:45:00  │
│                                                                          │
└─────────────────────────────────────────────────────────────────────────┘
```

### Data Model

```sql
CREATE TABLE trickplay_manifests (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    media_id        UUID NOT NULL,
    media_type      VARCHAR(50) NOT NULL,  -- 'movie', 'episode'
    width           INT NOT NULL,          -- Thumbnail width
    height          INT NOT NULL,          -- Thumbnail height
    tile_width      INT NOT NULL,          -- Images per row in sprite
    tile_height     INT NOT NULL,          -- Rows per sprite
    interval_ms     INT NOT NULL,          -- Milliseconds between thumbnails
    sprite_count    INT NOT NULL,          -- Number of sprite images
    base_path       TEXT NOT NULL,         -- Path to sprite images
    created_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_trickplay_media ON trickplay_manifests(media_type, media_id, width);
```

### Generation

```go
type TrickplayGenerator struct {
    blackbeard *TranscoderClient  // For thumbnail extraction
    config     TrickplayConfig
    logger     *slog.Logger
}

type TrickplayConfig struct {
    Widths        []int   // [320, 480] - Generate multiple resolutions
    IntervalSec   float64 // 10 seconds between thumbnails
    TileWidth     int     // 10 images per row
    TileHeight    int     // 10 rows per sprite
    Quality       int     // JPEG quality (75-90)
    BasePath      string  // Storage path
}

func DefaultTrickplayConfig() TrickplayConfig {
    return TrickplayConfig{
        Widths:      []int{320},
        IntervalSec: 10,
        TileWidth:   10,
        TileHeight:  10,
        Quality:     80,
        BasePath:    "/var/cache/revenge/trickplay",
    }
}

// GenerateTrickplay creates thumbnail sprites for a video
func (g *TrickplayGenerator) Generate(ctx context.Context, mediaID uuid.UUID, mediaPath string, duration float64) error {
    for _, width := range g.config.Widths {
        // Calculate dimensions
        height := width * 9 / 16 // Assume 16:9
        totalThumbs := int(duration / g.config.IntervalSec)
        thumbsPerSprite := g.config.TileWidth * g.config.TileHeight
        spriteCount := (totalThumbs + thumbsPerSprite - 1) / thumbsPerSprite

        // Generate sprites via Blackbeard
        for i := 0; i < spriteCount; i++ {
            startTime := float64(i*thumbsPerSprite) * g.config.IntervalSec
            spritePath := filepath.Join(g.config.BasePath, mediaID.String(),
                fmt.Sprintf("%d_%d.jpg", width, i))

            err := g.generateSprite(ctx, mediaPath, spritePath, width, height, startTime)
            if err != nil {
                return fmt.Errorf("generate sprite %d: %w", i, err)
            }
        }

        // Save manifest
        manifest := &TrickplayManifest{
            MediaID:     mediaID,
            MediaType:   "movie", // Or episode
            Width:       width,
            Height:      height,
            TileWidth:   g.config.TileWidth,
            TileHeight:  g.config.TileHeight,
            IntervalMS:  int(g.config.IntervalSec * 1000),
            SpriteCount: spriteCount,
            BasePath:    filepath.Join(g.config.BasePath, mediaID.String()),
        }
        g.saveManifest(ctx, manifest)
    }

    return nil
}

// River job for background generation
type GenerateTrickplayArgs struct {
    MediaID   uuid.UUID `json:"media_id"`
    MediaType string    `json:"media_type"`
    MediaPath string    `json:"media_path"`
    Duration  float64   `json:"duration"`
}

func (GenerateTrickplayArgs) Kind() string { return "trickplay.generate" }
```

### Client Support

```go
// Client capability check for trickplay
func (c *ClientCapabilities) SupportsTrickplay() bool {
    // Most modern clients support trickplay
    switch c.DeviceType {
    case "browser":
        return true // CSS sprites work everywhere
    case "mobile":
        return true // iOS/Android support
    case "tv":
        return true // Modern smart TVs
    default:
        return true
    }
}

// API response with trickplay info
type PlaybackInfo struct {
    MediaID      uuid.UUID          `json:"media_id"`
    StreamURL    string             `json:"stream_url"`
    Trickplay    *TrickplayInfo     `json:"trickplay,omitempty"`
    Segments     *PlaybackSegments  `json:"segments,omitempty"`
}

type TrickplayInfo struct {
    Width      int    `json:"width"`
    Height     int    `json:"height"`
    TileWidth  int    `json:"tile_width"`
    TileHeight int    `json:"tile_height"`
    IntervalMS int    `json:"interval_ms"`
    BaseURL    string `json:"base_url"`  // URL pattern for sprites
}
```

### Frontend Implementation

```typescript
// Trickplay thumbnail component
class TrickplayManager {
    private manifest: TrickplayInfo;
    private spriteCache: Map<number, HTMLImageElement> = new Map();

    constructor(manifest: TrickplayInfo) {
        this.manifest = manifest;
    }

    // Get thumbnail position for a given time
    getThumbnail(timeMs: number): ThumbnailPosition {
        const thumbIndex = Math.floor(timeMs / this.manifest.intervalMs);
        const thumbsPerSprite = this.manifest.tileWidth * this.manifest.tileHeight;

        const spriteIndex = Math.floor(thumbIndex / thumbsPerSprite);
        const indexInSprite = thumbIndex % thumbsPerSprite;

        const row = Math.floor(indexInSprite / this.manifest.tileWidth);
        const col = indexInSprite % this.manifest.tileWidth;

        return {
            spriteUrl: `${this.manifest.baseUrl}/${this.manifest.width}_${spriteIndex}.jpg`,
            x: col * this.manifest.width,
            y: row * this.manifest.height,
            width: this.manifest.width,
            height: this.manifest.height
        };
    }

    // Preload sprites around current position
    preload(currentTimeMs: number) {
        const thumbIndex = Math.floor(currentTimeMs / this.manifest.intervalMs);
        const thumbsPerSprite = this.manifest.tileWidth * this.manifest.tileHeight;
        const currentSprite = Math.floor(thumbIndex / thumbsPerSprite);

        // Preload current and adjacent sprites
        for (let i = Math.max(0, currentSprite - 1); i <= currentSprite + 1; i++) {
            if (!this.spriteCache.has(i)) {
                const img = new Image();
                img.src = `${this.manifest.baseUrl}/${this.manifest.width}_${i}.jpg`;
                this.spriteCache.set(i, img);
            }
        }
    }
}
```

---

## Chapter System

### Data Model

```sql
CREATE TABLE media_chapters (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    media_id        UUID NOT NULL,
    media_type      VARCHAR(50) NOT NULL,  -- 'movie', 'episode'
    chapter_index   INT NOT NULL,
    title           VARCHAR(255),
    start_time      FLOAT NOT NULL,        -- Seconds
    end_time        FLOAT,                 -- Seconds (NULL = until next/end)
    thumbnail_path  TEXT,                  -- Chapter thumbnail
    source          VARCHAR(50) NOT NULL,  -- 'embedded', 'manual', 'detected'
    created_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_media_chapters_media ON media_chapters(media_type, media_id);
```

### Chapter Extraction

```go
type ChapterService struct {
    chapterRepo ChapterRepository
    blackbeard  *TranscoderClient
    logger      *slog.Logger
}

// ExtractChapters extracts chapters from media file
func (s *ChapterService) ExtractChapters(ctx context.Context, mediaID uuid.UUID, mediaPath string) ([]Chapter, error) {
    // Use FFprobe to extract chapter metadata
    chapters, err := s.probeChapters(ctx, mediaPath)
    if err != nil {
        return nil, err
    }

    // Generate thumbnails for each chapter
    for i := range chapters {
        thumbPath := s.generateChapterThumb(ctx, mediaPath, chapters[i].StartTime)
        chapters[i].ThumbnailPath = thumbPath
    }

    // Save to database
    for _, ch := range chapters {
        ch.MediaID = mediaID
        ch.Source = "embedded"
        s.chapterRepo.Create(ctx, &ch)
    }

    return chapters, nil
}

func (s *ChapterService) probeChapters(ctx context.Context, path string) ([]Chapter, error) {
    // FFprobe command to get chapters
    cmd := exec.CommandContext(ctx, "ffprobe",
        "-v", "quiet",
        "-print_format", "json",
        "-show_chapters",
        path)

    output, err := cmd.Output()
    if err != nil {
        return nil, err
    }

    var probe struct {
        Chapters []struct {
            ID        int    `json:"id"`
            TimeBase  string `json:"time_base"`
            Start     int64  `json:"start"`
            End       int64  `json:"end"`
            Tags      struct {
                Title string `json:"title"`
            } `json:"tags"`
        } `json:"chapters"`
    }
    json.Unmarshal(output, &probe)

    chapters := make([]Chapter, len(probe.Chapters))
    for i, ch := range probe.Chapters {
        // Convert to seconds
        chapters[i] = Chapter{
            ChapterIndex: ch.ID,
            Title:        ch.Tags.Title,
            StartTime:    float64(ch.Start) / 1000000, // Microseconds to seconds
            EndTime:      float64(ch.End) / 1000000,
        }
    }

    return chapters, nil
}
```

---

## Live TV System

### Architecture

```
┌─────────────────────────────────────────────────────────────────────────┐
│                         Live TV Architecture                             │
└─────────────────────────────────────────────────────────────────────────┘

┌──────────────────┐     ┌─────────────────┐     ┌──────────────────────┐
│  External PVR    │ ──→ │  Revenge Live   │ ──→ │      Client          │
│   Servers        │     │   TV Service    │     │    (EPG + Stream)    │
└──────────────────┘     └─────────────────┘     └──────────────────────┘
        │                        │
        │               ┌────────┴────────┐
        │               ▼                 ▼
        │        ┌──────────┐      ┌──────────────┐
        │        │   EPG    │      │   DVR        │
        │        │  Cache   │      │  Scheduling  │
        │        └──────────┘      └──────────────┘
        │
┌───────┴───────────────────────────────────────────────────┐
│                    PVR Backends                            │
├──────────────┬──────────────┬──────────────┬──────────────┤
│ TVHeadend    │ Tvhproxy     │ NextPVR      │ IPTV Simple  │
└──────────────┴──────────────┴──────────────┴──────────────┘
```

### Data Model

```sql
-- Live TV channels
CREATE TABLE livetv_channels (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    pvr_id          VARCHAR(100) NOT NULL,  -- External PVR channel ID
    pvr_source      VARCHAR(50) NOT NULL,   -- 'tvheadend', 'nextpvr', etc.
    name            VARCHAR(255) NOT NULL,
    number          INT,
    logo_url        TEXT,
    stream_url      TEXT,                   -- Direct stream or PVR URL
    group_name      VARCHAR(255),           -- Channel group
    is_hd           BOOLEAN DEFAULT false,
    is_favorite     BOOLEAN DEFAULT false,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(pvr_source, pvr_id)
);

-- EPG (Electronic Program Guide)
CREATE TABLE livetv_programs (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    channel_id      UUID REFERENCES livetv_channels(id) ON DELETE CASCADE,
    pvr_program_id  VARCHAR(100),
    title           VARCHAR(255) NOT NULL,
    subtitle        VARCHAR(255),
    description     TEXT,
    start_time      TIMESTAMPTZ NOT NULL,
    end_time        TIMESTAMPTZ NOT NULL,
    category        VARCHAR(100),
    episode_info    VARCHAR(100),          -- "S01E05" format
    image_url       TEXT,
    is_movie        BOOLEAN DEFAULT false,
    is_sports       BOOLEAN DEFAULT false,
    is_news         BOOLEAN DEFAULT false,
    year            INT,
    rating          VARCHAR(20),           -- Content rating
    created_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_livetv_programs_channel ON livetv_programs(channel_id);
CREATE INDEX idx_livetv_programs_time ON livetv_programs(start_time, end_time);

-- DVR recordings
CREATE TABLE livetv_recordings (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID REFERENCES users(id),
    channel_id      UUID REFERENCES livetv_channels(id),
    program_id      UUID REFERENCES livetv_programs(id),
    pvr_recording_id VARCHAR(100),
    status          VARCHAR(50) NOT NULL,  -- 'scheduled', 'recording', 'completed', 'failed'
    start_time      TIMESTAMPTZ NOT NULL,
    end_time        TIMESTAMPTZ NOT NULL,
    pre_padding     INT DEFAULT 0,         -- Minutes before
    post_padding    INT DEFAULT 5,         -- Minutes after
    file_path       TEXT,                  -- Recorded file path
    file_size       BIGINT,
    created_at      TIMESTAMPTZ DEFAULT NOW()
);

-- Series recordings (record all episodes)
CREATE TABLE livetv_series_timers (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID REFERENCES users(id),
    series_name     VARCHAR(255) NOT NULL,
    channel_id      UUID REFERENCES livetv_channels(id),
    time_pattern    VARCHAR(50),           -- Time of day preference
    days_of_week    INT[],                 -- Days to record
    pre_padding     INT DEFAULT 0,
    post_padding    INT DEFAULT 5,
    keep_recordings INT DEFAULT 0,         -- 0 = keep all
    is_active       BOOLEAN DEFAULT true,
    created_at      TIMESTAMPTZ DEFAULT NOW()
);
```

### PVR Backend Interface

```go
// PVRBackend is the interface for external PVR servers
type PVRBackend interface {
    // Connection
    Connect(ctx context.Context) error
    Disconnect() error

    // Channels
    GetChannels(ctx context.Context) ([]Channel, error)
    GetChannelStream(ctx context.Context, channelID string) (string, error)

    // EPG
    GetEPG(ctx context.Context, channelID string, start, end time.Time) ([]Program, error)
    RefreshEPG(ctx context.Context) error

    // DVR
    GetRecordings(ctx context.Context) ([]Recording, error)
    ScheduleRecording(ctx context.Context, req RecordingRequest) error
    CancelRecording(ctx context.Context, recordingID string) error
    DeleteRecording(ctx context.Context, recordingID string) error

    // Timers
    GetTimers(ctx context.Context) ([]Timer, error)
    CreateSeriesTimer(ctx context.Context, req SeriesTimerRequest) error
    DeleteTimer(ctx context.Context, timerID string) error
}
```

### TVHeadend Backend

```go
type TVHeadendBackend struct {
    baseURL  string
    username string
    password string
    client   *http.Client
}

func (t *TVHeadendBackend) GetChannels(ctx context.Context) ([]Channel, error) {
    resp, err := t.get(ctx, "/api/channel/grid")
    if err != nil {
        return nil, err
    }

    var result struct {
        Entries []struct {
            UUID   string `json:"uuid"`
            Name   string `json:"name"`
            Number int    `json:"number"`
            Icon   string `json:"icon_public_url"`
        } `json:"entries"`
    }
    json.Unmarshal(resp, &result)

    channels := make([]Channel, len(result.Entries))
    for i, e := range result.Entries {
        channels[i] = Channel{
            PVRID:    e.UUID,
            Name:     e.Name,
            Number:   e.Number,
            LogoURL:  e.Icon,
        }
    }
    return channels, nil
}

func (t *TVHeadendBackend) GetChannelStream(ctx context.Context, channelID string) (string, error) {
    // Return TVHeadend stream URL
    return fmt.Sprintf("%s/stream/channel/%s", t.baseURL, channelID), nil
}
```

### Live TV Service

```go
type LiveTVService struct {
    backends    map[string]PVRBackend
    channelRepo ChannelRepository
    programRepo ProgramRepository
    recordRepo  RecordingRepository
    cache       *cache.Client
    logger      *slog.Logger
}

// GetGuide returns EPG for specified time range
func (s *LiveTVService) GetGuide(ctx context.Context, start, end time.Time, channelIDs []uuid.UUID) ([]GuideEntry, error) {
    // Check cache first
    cacheKey := fmt.Sprintf("epg:%s:%s", start.Format(time.RFC3339), end.Format(time.RFC3339))
    if cached, ok := s.cache.Get(ctx, cacheKey); ok {
        return cached.([]GuideEntry), nil
    }

    // Fetch from database
    programs, err := s.programRepo.GetByTimeRange(ctx, start, end, channelIDs)
    if err != nil {
        return nil, err
    }

    // Group by channel
    guide := s.groupByChannel(programs)
    s.cache.Set(ctx, cacheKey, guide, 5*time.Minute)

    return guide, nil
}

// RefreshEPG updates EPG data from all backends
func (s *LiveTVService) RefreshEPG(ctx context.Context) error {
    for name, backend := range s.backends {
        channels, err := backend.GetChannels(ctx)
        if err != nil {
            s.logger.Error("failed to get channels", "backend", name, "error", err)
            continue
        }

        for _, ch := range channels {
            // Fetch EPG for next 7 days
            programs, err := backend.GetEPG(ctx, ch.PVRID, time.Now(), time.Now().AddDate(0, 0, 7))
            if err != nil {
                continue
            }

            for _, p := range programs {
                s.programRepo.Upsert(ctx, &p)
            }
        }
    }
    return nil
}
```

### Timeshift Support

```go
// TimeshiftService manages timeshift buffers for live TV
type TimeshiftService struct {
    bufferDuration time.Duration  // How far back to buffer (e.g., 2 hours)
    storagePath    string
    buffers        map[uuid.UUID]*TimeshiftBuffer
    mu             sync.RWMutex
}

type TimeshiftBuffer struct {
    ChannelID    uuid.UUID
    StartTime    time.Time
    Segments     []TimeshiftSegment
    CurrentIndex int
}

// StartTimeshift begins buffering a channel
func (s *TimeshiftService) StartTimeshift(ctx context.Context, channelID uuid.UUID, streamURL string) (*TimeshiftBuffer, error) {
    s.mu.Lock()
    defer s.mu.Unlock()

    if buf, ok := s.buffers[channelID]; ok {
        return buf, nil // Already buffering
    }

    buf := &TimeshiftBuffer{
        ChannelID: channelID,
        StartTime: time.Now(),
        Segments:  make([]TimeshiftSegment, 0),
    }

    // Start background recording
    go s.recordBuffer(ctx, buf, streamURL)

    s.buffers[channelID] = buf
    return buf, nil
}

// SeekTo seeks to a position in the timeshift buffer
func (s *TimeshiftService) SeekTo(channelID uuid.UUID, position time.Duration) (string, error) {
    s.mu.RLock()
    buf, ok := s.buffers[channelID]
    s.mu.RUnlock()

    if !ok {
        return "", errors.New("no timeshift buffer for channel")
    }

    // Find segment at position
    targetTime := buf.StartTime.Add(position)
    for _, seg := range buf.Segments {
        if targetTime.After(seg.StartTime) && targetTime.Before(seg.EndTime) {
            return seg.Path, nil
        }
    }

    return "", errors.New("position not in buffer")
}
```

---

## Configuration

```yaml
# config.yaml - Media Enhancement Features

media_enhancements:
  # Trailers
  trailers:
    enabled: true
    max_per_movie: 3
    preferred_sources:
      - local
      - radarr
      - tmdb
    youtube_api_key: "${YOUTUBE_API_KEY}"
    default_skip_delay: 5  # Show skip button after 5 seconds

  # Audio Themes
  themes:
    enabled: true
    default_duration: 30  # seconds
    fade_duration: 2      # seconds
    default_volume: 0.3
    extract_on_scan: true

  # Intro/Outro Detection
  segments:
    enabled: true
    detect_on_scan: true
    chromaprint_path: "/usr/bin/fpcalc"
    import_introskipper: true
    import_path: "/config/introskipper"

  # Trickplay
  trickplay:
    enabled: true
    generate_on_scan: true
    widths: [320]
    interval_seconds: 10
    tile_size: 10
    jpeg_quality: 80
    storage_path: "/var/cache/revenge/trickplay"

  # Chapters
  chapters:
    enabled: true
    extract_on_scan: true
    generate_thumbnails: true

  # Live TV
  livetv:
    enabled: true
    epg_refresh_hours: 6
    epg_days_ahead: 7
    timeshift:
      enabled: true
      buffer_hours: 2
      storage_path: "/var/cache/revenge/timeshift"
    backends:
      - type: tvheadend
        url: "http://tvheadend:9981"
        username: "revenge"
        password: "${TVH_PASSWORD}"
      - type: iptv
        m3u_url: "http://iptv.example.com/playlist.m3u"
        epg_url: "http://iptv.example.com/epg.xml"
```


<!-- SOURCE-BREADCRUMBS-START -->

## Sources & Cross-References

> Auto-generated section linking to external documentation sources

### Cross-Reference Indexes

- [All Sources Index](../../../sources/SOURCES_INDEX.md) - Complete list of external documentation
- [Design ↔ Sources Map](../../../sources/DESIGN_CROSSREF.md) - Which docs reference which sources

### Referenced Sources

| Source | Documentation |
|--------|---------------|
| [FFmpeg Codecs](https://ffmpeg.org/ffmpeg-codecs.html) | [Local](../../../sources/media/ffmpeg-codecs.md) |
| [FFmpeg Documentation](https://ffmpeg.org/ffmpeg.html) | [Local](../../../sources/media/ffmpeg.md) |
| [FFmpeg Formats](https://ffmpeg.org/ffmpeg-formats.html) | [Local](../../../sources/media/ffmpeg-formats.md) |
| [M3U8 Extended Format](https://datatracker.ietf.org/doc/html/rfc8216) | [Local](../../../sources/protocols/m3u8.md) |
| [River Job Queue](https://pkg.go.dev/github.com/riverqueue/river) | [Local](../../../sources/tooling/river.md) |
| [Svelte 5 Documentation](https://svelte.dev/docs/svelte/overview) | [Local](../../../sources/frontend/svelte5.md) |
| [Svelte 5 Runes](https://svelte.dev/docs/svelte/$state) | [Local](../../../sources/frontend/svelte-runes.md) |
| [SvelteKit Documentation](https://svelte.dev/docs/kit/introduction) | [Local](../../../sources/frontend/sveltekit.md) |
| [go-astiav (FFmpeg bindings)](https://pkg.go.dev/github.com/asticode/go-astiav) | [Local](../../../sources/media/go-astiav.md) |

<!-- SOURCE-BREADCRUMBS-END -->

<!-- DESIGN-BREADCRUMBS-START -->

## Related Design Docs

> Auto-generated cross-references to related design documentation

**Category**: [Playback](INDEX.md)

### In This Section

- [Release Calendar System](RELEASE_CALENDAR.md)
- [Skip Intro / Credits Detection](SKIP_INTRO.md)
- [SyncPlay (Watch Together)](SYNCPLAY.md)
- [Trickplay (Timeline Thumbnails)](TRICKPLAY.md)
- [Watch Next & Continue Watching System](WATCH_NEXT_CONTINUE_WATCHING.md)

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

---

## Summary

| Feature | Detection/Source | Storage | Client Support |
|---------|------------------|---------|----------------|
| Trailers | Local, Radarr, TMDb, YouTube | DB + Files | All |
| Themes | Local, Extracted | DB + Files | Web, Apps |
| Intros | Chromaprint, Import, Manual | DB | All |
| Trickplay | Generated (FFmpeg) | Sprites | Capability-aware |
| Chapters | Embedded, Manual | DB + Thumbnails | All |
| Live TV | External PVR | DB (EPG) | All |
| DVR | External PVR | DB + Files | All |
