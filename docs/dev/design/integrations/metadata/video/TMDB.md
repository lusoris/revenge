# TMDb (The Movie Database) Integration

> Primary metadata provider for movies and TV shows


<!-- TOC-START -->

## Table of Contents

- [Status](#status)
- [Overview](#overview)
- [Developer Resources](#developer-resources)
- [API Details](#api-details)
  - [Key Endpoints](#key-endpoints)
- [Implementation Checklist](#implementation-checklist)
- [Revenge Integration Pattern](#revenge-integration-pattern)
  - [Go Client Example](#go-client-example)
- [Sources & Cross-References](#sources-cross-references)
  - [Cross-Reference Indexes](#cross-reference-indexes)
  - [Referenced Sources](#referenced-sources)
- [Related Design Docs](#related-design-docs)
  - [In This Section](#in-this-section)
  - [Related Topics](#related-topics)
  - [Indexes](#indexes)
- [Related Documentation](#related-documentation)
- [Image Sizes](#image-sizes)
  - [Posters](#posters)
  - [Backdrops](#backdrops)
  - [Logos](#logos)
  - [Profile Pictures (Actors/Directors)](#profile-pictures-actorsdirectors)
- [External ID Mapping](#external-id-mapping)
- [Notes](#notes)

<!-- TOC-END -->

## Status

| Dimension | Status | Notes |
|-----------|--------|-------|
| Design | ✅ | Comprehensive REST API v3 spec, image CDN, external ID mapping |
| Sources | ✅ | API docs v3/v4, image CDN, status page linked |
| Instructions | ✅ | Detailed implementation checklist with i18n support |
| Code | 🔴 |  |
| Linting | 🔴 |  |
| Unit Testing | 🔴 |  |
| Integration Testing | 🔴 |  |
---

## Overview

TMDb is the industry-standard metadata provider for movies and TV shows. Revenge uses TMDb as the primary metadata source for:
- Movie metadata (title, overview, cast, crew, release dates)
- TV show metadata (series info, seasons, episodes)
- Images (posters, backdrops, logos, profile pictures)
- External ID mapping (IMDb ID, TheTVDB ID, etc.)

**Integration Points**:
- **API client**: Query movies, TV shows, people, search
- **Image CDN**: Download posters, backdrops, logos
- **Blurhash generation**: Generate blur placeholders for images
- **Rate limiting**: 40 requests per 10 seconds

---

## Developer Resources

- 📚 **API Docs**: https://developers.themoviedb.org/3
- 🔗 **API v3**: https://api.themoviedb.org/3/
- 🔗 **API v4**: https://developers.themoviedb.org/4 (experimental, OAuth-based)
- 🔗 **Image CDN**: https://image.tmdb.org/t/p/
- 🔗 **Status Page**: https://status.themoviedb.org/

---

## API Details

**Base URL**: `https://api.themoviedb.org/3/`
**Authentication**: API Key query parameter `?api_key={key}`
**Rate Limits**: 40 requests per 10 seconds
**Free Tier**: Available (requires API key registration)
**i18n Support**: `language` parameter (e.g., `language=en-US`, `language=fr-FR`)

### Key Endpoints

| Endpoint | Purpose |
|----------|---------|
| `/movie/{movie_id}` | Get movie details |
| `/movie/{movie_id}/credits` | Get cast & crew |
| `/movie/{movie_id}/images` | Get posters, backdrops |
| `/movie/{movie_id}/external_ids` | Get IMDb/TheTVDB IDs |
| `/tv/{tv_id}` | Get TV show details |
| `/tv/{tv_id}/season/{season_number}` | Get season details |
| `/tv/{tv_id}/season/{season_number}/episode/{episode_number}` | Get episode details |
| `/person/{person_id}` | Get person details (actor/director) |
| `/search/movie` | Search movies by title |
| `/search/tv` | Search TV shows by title |
| `/configuration` | Get API configuration (image base URLs, sizes) |

---

## Implementation Checklist

- [ ] **API Client** (`internal/service/metadata/provider_tmdb.go`)
  - [ ] Movie metadata fetching
  - [ ] TV show metadata fetching
  - [ ] Season/episode metadata fetching
  - [ ] Cast & crew metadata
  - [ ] Search functionality
  - [ ] Rate limiting (40 req/10s)
  - [ ] Error handling & retries (circuit breaker)

- [ ] **Image Handling**
  - [ ] Download posters (w500, original sizes)
  - [ ] Download backdrops (w1280, original sizes)
  - [ ] Download logos (w500 size)
  - [ ] Generate Blurhash for placeholders
  - [ ] Store images locally (configurable path)
  - [ ] Image optimization (WebP conversion)

- [ ] **Metadata Enrichment**
  - [ ] Map TMDb movie → Revenge `movies` table
  - [ ] Map TMDb TV show → Revenge `tvshows` table
  - [ ] Map TMDb cast → Revenge `movie_people`, `tvshow_people`
  - [ ] Map TMDb genres → Revenge `movie_genres`, `tvshow_genres`
  - [ ] Extract external IDs (IMDb, TheTVDB) for cross-referencing

- [ ] **i18n Support**
  - [ ] Fetch metadata in multiple languages
  - [ ] Store translations in `movie_metadata` JSONB field
  - [ ] Fallback to English (`en-US`) if translation unavailable

---

## Revenge Integration Pattern

```
User requests movie metadata (TMDb ID: 603)
           ↓
Revenge queries TMDb API: /movie/603
           ↓
TMDb returns movie details + credits + images
           ↓
Revenge stores in PostgreSQL (movies, movie_people, movie_genres)
           ↓
Download posters/backdrops from TMDb CDN
           ↓
Generate Blurhash for placeholders
           ↓
Update Typesense search index
           ↓
Movie metadata available
```

### Go Client Example

```go
type TMDbClient struct {
    baseURL string
    apiKey  string
    client  *http.Client
    limiter *rate.Limiter  // 40 req/10s
}

func (c *TMDbClient) GetMovie(ctx context.Context, movieID int) (*Movie, error) {
    c.limiter.Wait(ctx)  // Rate limiting

    url := fmt.Sprintf("%s/movie/%d?api_key=%s", c.baseURL, movieID, c.apiKey)
    req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)

    resp, err := c.client.Do(req)
    if err != nil {
        return nil, fmt.Errorf("failed to get movie: %w", err)
    }
    defer resp.Body.Close()

    var movie Movie
    json.NewDecoder(resp.Body).Decode(&movie)
    return &movie, nil
}

func (c *TMDbClient) GetMovieCredits(ctx context.Context, movieID int) (*Credits, error) {
    c.limiter.Wait(ctx)

    url := fmt.Sprintf("%s/movie/%d/credits?api_key=%s", c.baseURL, movieID, c.apiKey)
    req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)

    resp, err := c.client.Do(req)
    if err != nil {
        return nil, fmt.Errorf("failed to get credits: %w", err)
    }
    defer resp.Body.Close()

    var credits Credits
    json.NewDecoder(resp.Body).Decode(&credits)
    return &credits, nil
}

func (c *TMDbClient) DownloadImage(ctx context.Context, imagePath string, size string) ([]byte, error) {
    // size = "w500", "w780", "original"
    url := fmt.Sprintf("https://image.tmdb.org/t/p/%s%s", size, imagePath)
    req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)

    resp, err := c.client.Do(req)
    if err != nil {
        return nil, fmt.Errorf("failed to download image: %w", err)
    }
    defer resp.Body.Close()

    return io.ReadAll(resp.Body)
}
```

---


<!-- SOURCE-BREADCRUMBS-START -->

## Sources & Cross-References

> Auto-generated section linking to external documentation sources

### Cross-Reference Indexes

- [All Sources Index](../../../../sources/SOURCES_INDEX.md) - Complete list of external documentation
- [Design ↔ Sources Map](../../../../sources/DESIGN_CROSSREF.md) - Which docs reference which sources

### Referenced Sources

| Source | Documentation |
|--------|---------------|
| [PostgreSQL Arrays](https://www.postgresql.org/docs/current/arrays.html) | [Local](../../../../sources/database/postgresql-arrays.md) |
| [PostgreSQL JSON Functions](https://www.postgresql.org/docs/current/functions-json.html) | [Local](../../../../sources/database/postgresql-json.md) |
| [Typesense API](https://typesense.org/docs/latest/api/) | [Local](../../../../sources/infrastructure/typesense.md) |
| [Typesense Go Client](https://github.com/typesense/typesense-go) | [Local](../../../../sources/infrastructure/typesense-go.md) |
| [go-blurhash](https://pkg.go.dev/github.com/bbrks/go-blurhash) | [Local](../../../../sources/media/go-blurhash.md) |
| [golang.org/x/time](https://pkg.go.dev/golang.org/x/time) | [Local](../../../../sources/go/x/time.md) |
| [pgx PostgreSQL Driver](https://pkg.go.dev/github.com/jackc/pgx/v5) | [Local](../../../../sources/database/pgx.md) |

<!-- SOURCE-BREADCRUMBS-END -->

<!-- DESIGN-BREADCRUMBS-START -->

## Related Design Docs

> Auto-generated cross-references to related design documentation

**Category**: [Video](INDEX.md)

### In This Section

- [OMDb (Open Movie Database) Integration](OMDB.md)
- [ThePosterDB Integration](THEPOSTERDB.md)
- [TheTVDB Integration](THETVDB.md)

### Related Topics

- [Revenge - Architecture v2](../../../architecture/01_ARCHITECTURE.md) _Architecture_
- [Revenge - Design Principles](../../../architecture/02_DESIGN_PRINCIPLES.md) _Architecture_
- [Revenge - Metadata System](../../../architecture/03_METADATA_SYSTEM.md) _Architecture_
- [Revenge - Player Architecture](../../../architecture/04_PLAYER_ARCHITECTURE.md) _Architecture_
- [Plugin Architecture Decision](../../../architecture/05_PLUGIN_ARCHITECTURE_DECISION.md) _Architecture_

### Indexes

- [Design Index](../../../DESIGN_INDEX.md) - All design docs by category/topic
- [Source of Truth](../../../00_SOURCE_OF_TRUTH.md) - Package versions and status

<!-- DESIGN-BREADCRUMBS-END -->

## Related Documentation

- [Movie Module](../../features/video/MOVIE_MODULE.md)
- [TV Show Module](../../features/video/TVSHOW_MODULE.md)
- [TheTVDB Integration](THETVDB.md) - Primary TV metadata source
- [OMDb Integration](OMDB.md) - Fallback metadata + IMDb ratings
- [Metadata Enrichment Pattern](../../patterns/METADATA_ENRICHMENT.md)

---

## Image Sizes

### Posters
- `w92` (92px width)
- `w154`
- `w185`
- `w342`
- `w500` ⭐ **Recommended for UI**
- `w780`
- `original` (full resolution)

### Backdrops
- `w300`
- `w780`
- `w1280` ⭐ **Recommended for UI**
- `original`

### Logos
- `w45`
- `w92`
- `w154`
- `w185`
- `w300`
- `w500` ⭐ **Recommended for UI**
- `original`

### Profile Pictures (Actors/Directors)
- `w45`
- `w185` ⭐ **Recommended for UI**
- `h632`
- `original`

---

## External ID Mapping

TMDb provides external IDs for cross-referencing:

```json
{
  "id": 603,
  "imdb_id": "tt0133093",
  "facebook_id": "thematrixmovie",
  "instagram_id": "thematrixmovie",
  "twitter_id": "thematrixmovie"
}
```

**For TV Shows**:
```json
{
  "id": 1399,
  "imdb_id": "tt0944947",
  "freebase_mid": "/m/0524b41",
  "freebase_id": "/en/game_of_thrones",
  "tvdb_id": 121361,
  "tvrage_id": 24493
}
```

Use these IDs to link with TheTVDB, IMDb, etc.

---

## Notes

- **TMDb API v3 is stable** (widely used, production-ready)
- **API v4 is experimental** (OAuth-based, not required for Revenge)
- **Rate limit**: 40 requests per 10 seconds (use `golang.org/x/time/rate` limiter)
- **Free tier available** (register for API key at https://www.themoviedb.org/settings/api)
- **i18n support**: Use `language` parameter (e.g., `?language=fr-FR`) for localized metadata
- **Image CDN**: TMDb provides free CDN for images (no additional API calls)
- **Blurhash**: Generate locally from downloaded images (not provided by TMDb)
- **External IDs**: Use IMDb/TheTVDB IDs to cross-reference with other services
- **Cast & crew**: TMDb provides detailed cast/crew info (actor names, roles, profile pictures)
- **Genres**: TMDb uses genre IDs (map to Revenge genre names in config)
- **Release dates**: TMDb provides multiple release dates by region (use primary release date)
- **Certification**: TMDb provides content ratings (PG, PG-13, R, etc.) by region
