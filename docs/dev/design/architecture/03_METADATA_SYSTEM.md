# Revenge - Metadata System

> Servarr-first metadata with intelligent fallback and multi-language support.

## Design Philosophy

1. **Local Sources First** - Servarr suite provides curated, cached metadata
2. **Native Audio Content** - Audiobooks/Podcasts managed natively with metadata providers
3. **Avoid External Calls** - Only fetch what local sources don't provide
4. **Language Awareness** - UI shows available data immediately, fetches translations async
5. **Progressive Enhancement** - Never block UI waiting for metadata

---

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────────────────┐
│                          Metadata Flow                                   │
└─────────────────────────────────────────────────────────────────────────┘

┌──────────┐     ┌───────────────────┐     ┌─────────────────┐     ┌──────────┐
│  Library │ ──→ │  Local Importers  │ ──→ │ Metadata Store  │ ←── │   UI     │
│   Scan   │     │    (Primary)      │     │  (PostgreSQL)   │     │ Request  │
└──────────┘     └───────────────────┘     └─────────────────┘     └──────────┘
                       │                          │
        ┌──────────────┼──────────────┐           │
        ▼              ▼              ▼           ▼
┌─────────────┐ ┌─────────────┐ ┌─────────────┐ ┌─────────────────┐
│   Servarr   │ │   Native    │ │ StashApp    │ │  Translation    │
│ (Arr Suite) │ │ Audio/Books │ │ (Adult)     │ │   Job Queue     │
└─────────────┘ └─────────────┘ └─────────────┘ └─────────────────┘
        │              │              │           │
        └──────────────┼──────────────┘           ▼
                       ▼                  ┌─────────────────┐
                ┌──────────────┐          │  Fetch Missing  │
                │   Missing?   │          │   Languages     │
                │   Fallback   │          └─────────────────┘
                └──────────────┘
                       │
                       ▼
                ┌──────────────┐
                │  TMDb/TVDB   │
                │  MusicBrainz │
                │  Audnexus    │
                └──────────────┘
```

---

## Data Sources Priority

### Movies

| Priority | Source | Data Provided |
|----------|--------|---------------|
| 1 | Radarr | Title, year, overview, genres, ratings, poster, fanart, trailer |
| 2 | TMDb | Translations, cast, crew, keywords, collection info |
| 3 | IMDb | Ratings (IMDb score), content rating |
| 4 | OMDb | Rotten Tomatoes, Metacritic scores |
| 5 | Fanart.tv | Additional artwork (logo, disc art, clearart) |

### TV Shows

| Priority | Source | Data Provided |
|----------|--------|---------------|
| 1 | Sonarr | Title, overview, genres, network, poster, fanart |
| 2 | TheTVDB | Episode details, translations, actor images |
| 3 | TMDb | Alternative titles, translations, keywords |
| 4 | Fanart.tv | Additional artwork |

### Music

| Priority | Source | Data Provided |
|----------|--------|---------------|
| 1 | Lidarr | Artist, album, tracks, cover art |
| 2 | MusicBrainz | Detailed credits, release info, genres |
| 3 | Last.fm | Similar artists, tags, play counts |
| 4 | Spotify | Popularity, audio features (if API available) |
| 5 | Fanart.tv | Artist images, album art variations |

### Audiobooks

| Priority | Source | Data Provided |
|----------|--------|---------------|
| 1 | **Native Scanner** | File metadata (ID3/M4B), duration, chapters, cover |
| 2 | Audnexus | Chapters, metadata from Audible |
| 3 | Chaptarr | Series info, Goodreads/Hardcover metadata (Readarr API spec) |
| 4 | Goodreads | Reviews, ratings, series info |
| 5 | OpenLibrary | ISBN, editions, subjects |

### Books (E-Books)

| Priority | Source | Data Provided |
|----------|--------|---------------|
| 1 | **Native Scanner** | EPUB/PDF metadata, cover extraction |
| 2 | Chaptarr | Series info, Goodreads/Hardcover metadata |
| 3 | Goodreads | Reviews, ratings, lists |
| 4 | Hardcover | Modern book database, lists |
| 5 | OpenLibrary | ISBN, editions, subjects |

### Podcasts

| Priority | Source | Data Provided |
|----------|--------|---------------|
| 1 | **RSS Feed** | Episodes, artwork, metadata (native parsing) |
| 2 | iTunes API | Categories, ratings, reviews |
| 3 | Podcast Index | Chapters (podcast namespace), value4value |
| 4 | Podchaser | Guest info, ratings |

### Adult Content

| Priority | Source | Data Provided |
|----------|--------|---------------|
| 1 | **Whisparr-v3** | Title, performers, studio, tags, cover |
| 2 | TPDB (ThePornDB) | Extended metadata, performer info |
| 3 | Stash-Box | Scene fingerprints, performer database |

---

## Multi-Language Support

### Language Priority Resolution

```go
// User's language preference chain
type LanguagePreference struct {
    Primary   string   // e.g., "de"
    Fallbacks []string // e.g., ["en", "original"]
}

// Resolution order for German user:
// 1. German metadata (if available)
// 2. English metadata (fallback)
// 3. Original language metadata
```

### Translation Storage

```sql
-- Per-module translation tables
CREATE TABLE movie_translations (
    movie_id        UUID REFERENCES movies(id),
    language        VARCHAR(5) NOT NULL,  -- ISO 639-1 + region
    title           TEXT,
    tagline         TEXT,
    overview        TEXT,
    source          VARCHAR(50),          -- 'arr', 'tmdb', 'manual'
    fetched_at      TIMESTAMPTZ,
    PRIMARY KEY (movie_id, language)
);

-- Index for fast lookups
CREATE INDEX idx_movie_translations_lang ON movie_translations(language);
```

### Translation Fetch Flow

```
User requests movie in German
         │
         ▼
┌─────────────────────────────────┐
│ Check movie_translations        │
│ WHERE language = 'de'           │
└─────────────────────────────────┘
         │
         ├── Found? ──→ Return German metadata
         │
         ▼ Not found
┌─────────────────────────────────┐
│ Return English metadata         │
│ (from Servarr cache)            │
│ + Queue translation job         │
└─────────────────────────────────┘
         │
         ▼
┌─────────────────────────────────┐
│ River Job: FetchTranslation     │
│ - Fetch from TMDb in German     │
│ - Store in movie_translations   │
│ - Notify UI via WebSocket       │
└─────────────────────────────────┘
```

### UI Behavior

1. **Immediate Response**: Show available metadata (even if wrong language)
2. **Loading Indicator**: Mark fields being translated
3. **Real-time Update**: WebSocket pushes translated content
4. **Graceful Fallback**: Never show empty fields if any translation exists

```typescript
// Frontend translation state
interface MovieMetadata {
  id: string;
  title: string;
  titleLanguage: string;         // Actual language of title
  overview: string;
  overviewLanguage: string;
  isTranslating: boolean;        // Translation job in progress
  availableLanguages: string[];  // Languages we have
}
```

---

## Local Data Sources

### Servarr Suite

The complete Servarr suite provides curated, cached metadata for most content types:

| Service | Content Type | API Version | Data Provided |
|---------|--------------|-------------|---------------|
| **Radarr** | Movies | v3 | Full metadata, ratings, artwork |
| **Sonarr** | TV Shows | v3 | Series, seasons, episodes |
| **Lidarr** | Music | v1 | Artists, albums, tracks |
| **Whisparr-v3** | Adult | v3 | Scenes, performers, studios |
| **Chaptarr** | Books | Readarr spec | Goodreads, Hardcover integration |

### Native Audio Content Management

Audiobooks and podcasts are managed natively without external dependencies:

```go
// Native audiobook metadata extraction
type AudiobookScanner struct {
    tagReader  TagReader      // ID3/M4B metadata
    audnexus   *AudnexusClient // Chapter lookup
    covers     CoverExtractor
}

// Extract metadata from audio files
func (s *AudiobookScanner) ScanFile(ctx context.Context, path string) (*Audiobook, error) {
    // 1. Read embedded metadata (ID3, M4B atoms)
    // 2. Extract embedded cover art
    // 3. Parse chapter markers
    // 4. Lookup additional data from Audnexus if ASIN available
}

// Native podcast RSS parsing
type PodcastService struct {
    parser    *gofeed.Parser
    downloader EpisodeDownloader
    scheduler  *river.Client
}

// Subscribe and parse RSS feed
func (s *PodcastService) Subscribe(ctx context.Context, feedURL string) (*Podcast, error) {
    feed, err := s.parser.ParseURL(feedURL)
    if err != nil {
        return nil, fmt.Errorf("parse feed: %w", err)
    }

    podcast := &Podcast{
        Title:       feed.Title,
        Author:      feed.Author.Name,
        Description: feed.Description,
        FeedURL:     feedURL,
        ImageURL:    feed.Image.URL,
    }

    // Convert episodes
    for _, item := range feed.Items {
        podcast.Episodes = append(podcast.Episodes, convertEpisode(item))
    }

    return podcast, nil
}

// Schedule episode downloads via River
func (s *PodcastService) DownloadEpisode(ctx context.Context, episodeID uuid.UUID) error {
    return s.scheduler.Insert(ctx, &DownloadEpisodeArgs{EpisodeID: episodeID})
}
```

---

## Servarr Integration

### Radarr API Integration

```go
type RadarrClient struct {
    baseURL string
    apiKey  string
    client  *http.Client
}

type RadarrMovie struct {
    ID              int       `json:"id"`
    Title           string    `json:"title"`
    OriginalTitle   string    `json:"originalTitle"`
    Year            int       `json:"year"`
    Overview        string    `json:"overview"`
    TmdbID          int       `json:"tmdbId"`
    ImdbID          string    `json:"imdbId"`
    Genres          []string  `json:"genres"`
    Runtime         int       `json:"runtime"`
    Ratings         Ratings   `json:"ratings"`
    Images          []Image   `json:"images"`
    Path            string    `json:"path"`
    HasFile         bool      `json:"hasFile"`
    MovieFile       MovieFile `json:"movieFile,omitempty"`
}

// Sync from Radarr
func (c *RadarrClient) SyncMovies(ctx context.Context) ([]RadarrMovie, error) {
    resp, err := c.get(ctx, "/api/v3/movie")
    if err != nil {
        return nil, fmt.Errorf("radarr sync: %w", err)
    }
    // Parse and return movies
}
```

### Sonarr API Integration

```go
type SonarrClient struct {
    baseURL string
    apiKey  string
    client  *http.Client
}

type SonarrSeries struct {
    ID              int       `json:"id"`
    Title           string    `json:"title"`
    OriginalTitle   string    `json:"originalTitle"`
    Year            int       `json:"year"`
    Overview        string    `json:"overview"`
    TvdbID          int       `json:"tvdbId"`
    ImdbID          string    `json:"imdbId"`
    Genres          []string  `json:"genres"`
    Network         string    `json:"network"`
    Status          string    `json:"status"`
    Seasons         []Season  `json:"seasons"`
    Images          []Image   `json:"images"`
    Path            string    `json:"path"`
}

// Sync from Sonarr
func (c *SonarrClient) SyncSeries(ctx context.Context) ([]SonarrSeries, error) {
    resp, err := c.get(ctx, "/api/v3/series")
    if err != nil {
        return nil, fmt.Errorf("sonarr sync: %w", err)
    }
    // Parse and return series
}
```

### Webhook Listeners

```go
// Receive webhooks from Servarr for real-time updates
func (h *WebhookHandler) HandleRadarrWebhook(w http.ResponseWriter, r *http.Request) {
    var event RadarrWebhookEvent
    if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
        http.Error(w, "invalid payload", http.StatusBadRequest)
        return
    }

    switch event.EventType {
    case "Download":
        // New movie downloaded - update metadata
        h.jobs.Insert(ctx, &UpdateMovieMetadataArgs{
            RadarrID: event.Movie.ID,
        })
    case "MovieDelete":
        // Movie deleted - mark as unavailable
        h.movies.MarkUnavailable(ctx, event.Movie.ID)
    case "Rename":
        // File renamed - update paths
        h.movies.UpdatePath(ctx, event.Movie.ID, event.Movie.Path)
    }

    w.WriteHeader(http.StatusOK)
}
```

---

## Fallback Providers

### TMDb Provider

```go
type TMDbProvider struct {
    apiKey string
    client *http.Client
    cache  *cache.Client
}

type TMDbMovie struct {
    ID               int                `json:"id"`
    Title            string             `json:"title"`
    OriginalTitle    string             `json:"original_title"`
    Overview         string             `json:"overview"`
    Tagline          string             `json:"tagline"`
    ReleaseDate      string             `json:"release_date"`
    Runtime          int                `json:"runtime"`
    Genres           []TMDbGenre        `json:"genres"`
    ProductionCompanies []TMDbCompany   `json:"production_companies"`
    Credits          *TMDbCredits       `json:"credits,omitempty"`
    Translations     *TMDbTranslations  `json:"translations,omitempty"`
    ExternalIDs      *TMDbExternalIDs   `json:"external_ids,omitempty"`
}

// Fetch with translations
func (p *TMDbProvider) GetMovie(ctx context.Context, tmdbID int, language string) (*TMDbMovie, error) {
    cacheKey := fmt.Sprintf("tmdb:movie:%d:%s", tmdbID, language)

    // Check cache first
    if cached, err := p.cache.Get(ctx, cacheKey); err == nil {
        var movie TMDbMovie
        json.Unmarshal(cached, &movie)
        return &movie, nil
    }

    // Fetch from TMDb
    url := fmt.Sprintf("%s/movie/%d?language=%s&append_to_response=credits,translations,external_ids",
        tmdbBaseURL, tmdbID, language)

    // ... fetch and cache
}
```

### MusicBrainz Provider

```go
type MusicBrainzProvider struct {
    client *http.Client
    cache  *cache.Client
}

type MBArtist struct {
    ID            string   `json:"id"`
    Name          string   `json:"name"`
    SortName      string   `json:"sort-name"`
    Type          string   `json:"type"`
    Country       string   `json:"country"`
    Disambiguation string  `json:"disambiguation"`
    LifeSpan      LifeSpan `json:"life-span"`
    Tags          []Tag    `json:"tags"`
}

// Fetch artist with releases
func (p *MusicBrainzProvider) GetArtist(ctx context.Context, mbid string) (*MBArtist, error) {
    // MusicBrainz rate limit: 1 req/sec
    // Use circuit breaker + rate limiter
}
```

---

## Metadata Jobs

### Job Definitions

```go
// Fetch missing metadata for a movie
type FetchMovieMetadataArgs struct {
    MovieID   uuid.UUID `json:"movie_id"`
    TmdbID    int       `json:"tmdb_id"`
    FetchFull bool      `json:"fetch_full"` // Include credits, keywords
}

func (FetchMovieMetadataArgs) Kind() string { return "metadata.fetch_movie" }

// Fetch translation for specific language
type FetchTranslationArgs struct {
    Module   string    `json:"module"`   // "movie", "tvshow", etc.
    ItemID   uuid.UUID `json:"item_id"`
    Language string    `json:"language"` // ISO 639-1
}

func (FetchTranslationArgs) Kind() string { return "metadata.fetch_translation" }

// Refresh all metadata for a library
type RefreshLibraryMetadataArgs struct {
    LibraryID uuid.UUID `json:"library_id"`
    Force     bool      `json:"force"` // Ignore cache
}

func (RefreshLibraryMetadataArgs) Kind() string { return "metadata.refresh_library" }

// Sync with Servarr
type SyncServarrArgs struct {
    ServarrType string    `json:"servarr_type"` // "radarr", "sonarr", "lidarr"
    ServarrID   uuid.UUID `json:"servarr_id"`
}

func (SyncServarrArgs) Kind() string { return "metadata.sync_servarr" }
```

### Worker Implementation

```go
type FetchTranslationWorker struct {
    river.WorkerDefaults[FetchTranslationArgs]
    tmdb    *TMDbProvider
    movies  *MovieRepository
    tvshows *TVShowRepository
}

func (w *FetchTranslationWorker) Work(ctx context.Context, job *river.Job[FetchTranslationArgs]) error {
    args := job.Args

    switch args.Module {
    case "movie":
        movie, err := w.movies.GetByID(ctx, args.ItemID)
        if err != nil {
            return fmt.Errorf("get movie: %w", err)
        }

        // Fetch translation from TMDb
        tmdbMovie, err := w.tmdb.GetMovie(ctx, movie.TmdbID, args.Language)
        if err != nil {
            return fmt.Errorf("fetch tmdb: %w", err)
        }

        // Store translation
        err = w.movies.SaveTranslation(ctx, args.ItemID, args.Language, Translation{
            Title:    tmdbMovie.Title,
            Tagline:  tmdbMovie.Tagline,
            Overview: tmdbMovie.Overview,
            Source:   "tmdb",
        })
        if err != nil {
            return fmt.Errorf("save translation: %w", err)
        }

        // Notify UI via WebSocket
        w.notifyTranslationReady(ctx, args.Module, args.ItemID, args.Language)

    case "tvshow":
        // Similar for TV shows
    }

    return nil
}
```

---

## Image Management

### Image Types

| Type | Description | Aspect Ratio |
|------|-------------|--------------|
| Poster | Main promotional art | 2:3 |
| Backdrop/Fanart | Background art | 16:9 |
| Logo | Transparent logo | varies |
| Thumb | Episode thumbnail | 16:9 |
| Banner | Series banner | ~5:1 |
| Disc | DVD/Blu-ray disc art | 1:1 |
| Clearart | Transparent art | varies |

### Image Storage

```sql
CREATE TABLE movie_images (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    movie_id    UUID REFERENCES movies(id) ON DELETE CASCADE,
    type        VARCHAR(20) NOT NULL,     -- 'poster', 'backdrop', 'logo'
    language    VARCHAR(5),               -- NULL for no text
    url         TEXT NOT NULL,            -- Original URL
    local_path  TEXT,                     -- Local cached path
    width       INT,
    height      INT,
    vote_count  INT DEFAULT 0,            -- From TMDb
    vote_avg    FLOAT DEFAULT 0,
    blurhash    VARCHAR(50),              -- For loading placeholder
    is_primary  BOOLEAN DEFAULT false,
    source      VARCHAR(50),              -- 'arr', 'tmdb', 'fanart', 'manual'
    created_at  TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_movie_images_type ON movie_images(movie_id, type, is_primary);
```

### Blurhash Generation

```go
// Generate blurhash for image placeholders
func GenerateBlurhash(imagePath string) (string, error) {
    file, err := os.Open(imagePath)
    if err != nil {
        return "", err
    }
    defer file.Close()

    img, _, err := image.Decode(file)
    if err != nil {
        return "", err
    }

    // Resize for faster processing
    resized := resize.Resize(32, 0, img, resize.Lanczos3)

    hash, err := blurhash.Encode(4, 3, resized)
    if err != nil {
        return "", err
    }

    return hash, nil
}
```

---

## Caching Strategy

### Cache Layers

```
┌─────────────────────────────────────────────────────────────────────────┐
│                        Metadata Cache Layers                             │
└─────────────────────────────────────────────────────────────────────────┘

┌──────────────┐    ┌──────────────┐    ┌──────────────┐    ┌──────────────┐
│  PostgreSQL  │ ←─ │  Dragonfly   │ ←─ │   In-Memory  │ ←─ │   Request    │
│  (Permanent) │    │  (1h TTL)    │    │  (5min TTL)  │    │   Context    │
└──────────────┘    └──────────────┘    └──────────────┘    └──────────────┘
```

### Cache Keys

```go
const (
    // Metadata cache
    CacheKeyMovie         = "meta:movie:%s"           // meta:movie:{id}
    CacheKeyMovieTrans    = "meta:movie:%s:trans:%s"  // meta:movie:{id}:trans:{lang}
    CacheKeyTVShow        = "meta:tvshow:%s"
    CacheKeyEpisode       = "meta:episode:%s"

    // Provider cache (prevent rate limiting)
    CacheKeyTMDb          = "provider:tmdb:%s:%s"     // provider:tmdb:{type}:{id}
    CacheKeyTVDB          = "provider:tvdb:%s:%s"
    CacheKeyMusicBrainz   = "provider:mb:%s:%s"

    // Servarr sync state
    CacheKeyServarrSync   = "servarr:%s:sync:%s"      // servarr:{type}:sync:{id}
)
```

### TTL Configuration

```yaml
# configs/config.yaml
metadata:
  cache:
    # In-memory (go-cache)
    memory_ttl: 5m
    memory_cleanup: 10m

    # Dragonfly
    redis_ttl: 1h

    # Provider cache (respect rate limits)
    provider_ttl: 24h

  # Refresh intervals
  refresh:
    automatic: true
    interval: 7d              # Full refresh every 7 days
    on_play: true             # Check for updates on playback
```

---

## Configuration

### Servarr Configuration

```yaml
# configs/config.yaml
metadata:
  servarr:
    radarr:
      - name: "Main Radarr"
        url: "http://radarr:7878"
        api_key: "${RADARR_API_KEY}"
        sync_interval: 15m
        webhook_enabled: true

    sonarr:
      - name: "Main Sonarr"
        url: "http://sonarr:8989"
        api_key: "${SONARR_API_KEY}"
        sync_interval: 15m
        webhook_enabled: true

    lidarr:
      - name: "Music"
        url: "http://lidarr:8686"
        api_key: "${LIDARR_API_KEY}"
        sync_interval: 30m

    readarr:
      - name: "Books"
        url: "http://readarr:8787"
        api_key: "${READARR_API_KEY}"
        sync_interval: 1h

  providers:
    tmdb:
      api_key: "${TMDB_API_KEY}"
      include_adult: false
      rate_limit: 40          # requests per 10 seconds

    tvdb:
      api_key: "${TVDB_API_KEY}"
      pin: "${TVDB_PIN}"

    musicbrainz:
      user_agent: "Revenge/1.0"
      rate_limit: 1           # request per second

    fanart:
      api_key: "${FANART_API_KEY}"
      personal_key: "${FANART_PERSONAL_KEY}"  # Optional, higher limits

    omdb:
      api_key: "${OMDB_API_KEY}"
```

### Language Configuration

```yaml
metadata:
  languages:
    default: "en"
    ui_default: "en"
    available:
      - "en"
      - "de"
      - "fr"
      - "es"
      - "it"
      - "pt"
      - "ja"
      - "zh"

    # Fetch translations for these languages on new content
    prefetch:
      - "en"
      - "de"                  # German translation auto-fetched
```

---

## API Endpoints

### Metadata Endpoints

```yaml
# api/openapi/metadata.yaml
paths:
  /api/v1/movies/{id}/metadata:
    get:
      summary: Get movie metadata
      parameters:
        - name: language
          in: query
          schema:
            type: string
            default: en
      responses:
        200:
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/MovieMetadata'

  /api/v1/movies/{id}/metadata/refresh:
    post:
      summary: Refresh movie metadata
      responses:
        202:
          description: Refresh job queued

  /api/v1/movies/{id}/translations:
    get:
      summary: Get available translations
      responses:
        200:
          content:
            application/json:
              schema:
                type: array
                items:
                  $ref: '#/components/schemas/Translation'

  /api/v1/movies/{id}/translations/{language}:
    post:
      summary: Request translation for language
      responses:
        202:
          description: Translation job queued
```

---

## Summary

| Aspect | Decision |
|--------|----------|
| Primary Source | Arrs (Radarr, Sonarr, Lidarr, Chaptarr) |
| Fallback | TMDb, TVDB, MusicBrainz, etc. |
| Translation Strategy | Show available immediately, fetch async |
| Cache Layers | Memory (5min) → Dragonfly (1h) → PostgreSQL |
| Image Processing | Blurhash placeholders, local caching |
| Sync Method | Webhooks (real-time) + Polling (backup) |
| Job Queue | River for all background fetching |
