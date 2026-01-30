# Metadata Service

> External metadata providers for media enrichment

**Location**: `internal/service/metadata/`

---

## Overview

The Metadata service provides access to external APIs for fetching media metadata:

- **TMDb** - Movies and TV shows (primary external source)
- **Radarr** - Movies via Servarr (primary when connected)

---

## Provider Architecture

```
┌─────────────────────────────────────────────┐
│            Metadata Service                  │
└──────────────┬──────────────────────────────┘
               │
    ┌──────────┴──────────┐
    ▼                     ▼
┌────────┐          ┌──────────┐
│ Radarr │          │   TMDb   │
│Provider│          │ Provider │
│Priority│          │ Priority │
│   1    │          │    2     │
└────────┘          └──────────┘
```

---

## TMDb Provider

**Location**: `internal/service/metadata/tmdb/`

### Interface

```go
type Provider struct {
    client *Client
    logger *slog.Logger
}

func (p *Provider) Name() string        // Returns "tmdb"
func (p *Provider) Priority() int       // Returns 2 (secondary to Servarr)
func (p *Provider) IsAvailable() bool   // Checks if API key configured
```

### Movie Search

```go
type MovieSearchResult struct {
    TMDbID    int
    Title     string
    Year      int
    Overview  string
    PosterURL string
    Score     float64
}

func (p *Provider) SearchMovies(ctx context.Context, query string, year int) ([]MovieSearchResult, error)
```

### Movie Metadata

```go
type MovieMetadata struct {
    TMDbID           int
    IMDbID           string
    Title            string
    OriginalTitle    string
    OriginalLanguage string
    Overview         string
    Tagline          string
    RuntimeMinutes   int
    ReleaseDate      time.Time
    ReleaseYear      int
    Budget           int64
    Revenue          int64
    Rating           float64
    VoteCount        int
    Popularity       float64
    Adult            bool
    Status           string
    Homepage         string
    PosterURL        string
    BackdropURL      string
    Genres           []string
    Studios          []StudioInfo
    Collection       *CollectionInfo
    Cast             []CastInfo
    Crew             []CrewInfo
    Images           []ImageInfo
    Videos           []VideoInfo
}

func (p *Provider) GetMovieMetadata(ctx context.Context, tmdbID int) (*MovieMetadata, error)
```

### Movie Matching

```go
// Match by IMDb ID first, then title/year search
func (p *Provider) MatchMovie(ctx context.Context, title string, year int, imdbID string) (*MovieMetadata, error)

// Find by IMDb ID
func (p *Provider) FindByIMDbID(ctx context.Context, imdbID string) (*MovieSearchResult, error)
```

---

## Radarr Provider

**Location**: `internal/service/metadata/radarr/`

Primary metadata source when Radarr is connected. Uses cached metadata from Radarr which itself aggregates from TMDb and other sources.

### Priority

Radarr takes priority over TMDb when connected:
- **Priority 1**: Radarr (Servarr-first principle)
- **Priority 2**: TMDb (fallback/enrichment)

---

## Data Types

### Studio Info

```go
type StudioInfo struct {
    TMDbID        int
    Name          string
    LogoURL       string
    OriginCountry string
}
```

### Cast/Crew Info

```go
type CastInfo struct {
    TMDbID     int
    Name       string
    Character  string
    Order      int
    ProfileURL string
}

type CrewInfo struct {
    TMDbID     int
    Name       string
    Department string
    Job        string
    ProfileURL string
}
```

### Images

```go
type ImageInfo struct {
    Type        string  // poster, backdrop, logo
    URL         string
    Width       int
    Height      int
    AspectRatio float64
    Language    string
    VoteAverage float64
    VoteCount   int
}
```

### Videos

```go
type VideoInfo struct {
    Key         string  // YouTube/Vimeo key
    Name        string
    Site        string  // YouTube, Vimeo
    Type        string  // Trailer, Teaser, Clip
    Size        int     // Resolution
    Official    bool
    PublishedAt string
    Language    string
}
```

---

## Configuration

```yaml
metadata:
  tmdb:
    api_key: "${TMDB_API_KEY}"
    enabled: true
  radarr:
    url: "http://localhost:7878"
    api_key: "${RADARR_API_KEY}"
    enabled: true
```

---

## Related

- [Metadata System](../architecture/03_METADATA_SYSTEM.md) - Architecture
- [TMDb Integration](../integrations/metadata/video/TMDB.md) - API details
- [Radarr Integration](../integrations/servarr/RADARR.md) - Servarr integration
