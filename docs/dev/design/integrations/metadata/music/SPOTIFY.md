# Spotify Integration

> Music metadata and cover art provider - popularity scores, high-quality images

**Service**: Spotify
**Type**: Metadata Provider (Music)
**API Version**: Web API v1
**Website**: https://www.spotify.com
**API Docs**: https://developer.spotify.com/documentation/web-api

## Status

| Dimension | Status | Notes |
| --------- | ------ | ----- |
| Design | ✅ | Comprehensive REST API endpoints, OAuth flow, token management |
| Sources | ✅ | API docs, authentication, console linked |
| Instructions | ✅ | Detailed implementation checklist |
| Code | 🔴 | |
| Linting | 🔴 | |
| Unit Testing | 🔴 | |
| Integration Testing | 🔴 | |

---

## Overview

**Spotify** provides music metadata (artists, albums, tracks) and high-quality cover art. Used as **fallback** for cover art and **popularity scores**.

**Why Spotify**:
- High-quality cover art (640x640, 300x300)
- Popularity scores (0-100)
- Rich metadata (genres, release precision)
- Free API (generous limits)
- OAuth authentication

**Use Cases**:
- **Cover art fallback**: When Cover Art Archive fails
- **Popularity scores**: Track/album/artist popularity (0-100)
- **Metadata enrichment**: Genres, release precision
- **Search**: Fallback search when MusicBrainz fails

**NOT Used For**:
- Streaming (requires Spotify Premium, separate integration)
- Scrobbling (use Last.fm/ListenBrainz)

---

## Developer Resources

**API Documentation**: https://developer.spotify.com/documentation/web-api
**Authentication**: https://developer.spotify.com/documentation/web-api/tutorials/getting-started
**Console**: https://developer.spotify.com/console/

**Authentication**: OAuth 2.0 (Client Credentials for metadata)
**Rate Limit**: 180 requests per minute (per app)
**Free Tier**: Unlimited (API app required)

---

## API Details

### Base URL
```
https://api.spotify.com/v1/
```

### Authentication (Client Credentials Flow)
```bash
# 1. Get access token
POST https://accounts.spotify.com/api/token
Content-Type: application/x-www-form-urlencoded
Authorization: Basic {BASE64(CLIENT_ID:CLIENT_SECRET)}

grant_type=client_credentials

# Response:
{
  "access_token": "NgCXRK...MzYjw",
  "token_type": "Bearer",
  "expires_in": 3600
}

# 2. Use token in requests
Authorization: Bearer {ACCESS_TOKEN}
```

### Rate Limiting
- **180 requests per minute** (per app)
- **429 Too Many Requests** (retry after `Retry-After` header)

### Key Endpoints

#### Search Artists
```bash
GET /v1/search?q=radiohead&type=artist&limit=10
Authorization: Bearer {ACCESS_TOKEN}
```

**Response**:
```json
{
  "artists": {
    "items": [
      {
        "id": "4Z8W4fKeB5YxbusRsdQVPb",
        "name": "Radiohead",
        "genres": ["alternative rock", "art rock", "melancholia", "permanent wave", "rock"],
        "popularity": 79,
        "images": [
          {"height": 640, "width": 640, "url": "https://i.scdn.co/image/..."}
        ],
        "external_urls": {"spotify": "https://open.spotify.com/artist/..."}
      }
    ]
  }
}
```

#### Get Artist
```bash
GET /v1/artists/{id}
```

#### Get Album
```bash
GET /v1/albums/{id}
```

**Response**:
```json
{
  "id": "6dVIqQ8qmQ5GBnJ9shOYGE",
  "name": "OK Computer",
  "release_date": "1997-05-21",
  "release_date_precision": "day",
  "total_tracks": 12,
  "images": [
    {"height": 640, "width": 640, "url": "https://i.scdn.co/image/..."},
    {"height": 300, "width": 300, "url": "https://i.scdn.co/image/..."}
  ],
  "artists": [{"id": "4Z8W4fKeB5YxbusRsdQVPb", "name": "Radiohead"}],
  "popularity": 76
}
```

#### Get Track
```bash
GET /v1/tracks/{id}
```

**Response**:
```json
{
  "id": "6ORfhpZbYM8PDJ5CBAUphD",
  "name": "Karma Police",
  "duration_ms": 263000,
  "popularity": 78,
  "artists": [{"id": "4Z8W4fKeB5YxbusRsdQVPb", "name": "Radiohead"}],
  "album": {
    "id": "6dVIqQ8qmQ5GBnJ9shOYGE",
    "name": "OK Computer",
    "images": [{"height": 640, "width": 640, "url": "..."}]
  }
}
```

---

## Implementation Checklist

### API Client (`internal/infra/metadata/provider_spotify.go`)
- [ ] Base URL configuration
- [ ] Client ID + Client Secret configuration
- [ ] OAuth token management (Client Credentials flow)
- [ ] Token refresh (expires after 3600s)
- [ ] Rate limiting (180 req/min with token bucket)
- [ ] Error handling (401: Token expired, 429: Rate limit exceeded)
- [ ] Response parsing (JSON unmarshalling)

### Token Management
- [ ] Request access token on startup
- [ ] Cache token in memory (3600s TTL)
- [ ] Auto-refresh token (background job every 3500s)
- [ ] Handle 401 (token expired - request new token)

### Cover Art Fallback
- [ ] Search album by artist + title (when Cover Art Archive fails)
- [ ] Extract highest resolution image (640x640)
- [ ] Download cover art
- [ ] Generate Blurhash
- [ ] Convert to WebP
- [ ] Store locally (`data/music/covers/`)

### Popularity Scores
- [ ] Fetch artist/album/track popularity (0-100)
- [ ] Store in `music_artists.popularity`, `music_albums.popularity`, `music_tracks.popularity`
- [ ] Update periodically (weekly background job)

### Metadata Enrichment
- [ ] Extract genres (artist-level)
- [ ] Extract release precision (day, month, year)
- [ ] Store in database

### Error Handling
- [ ] Handle 401 (Token expired - refresh token)
- [ ] Handle 404 (Artist/album/track not found)
- [ ] Handle 429 (Rate limit exceeded - retry after `Retry-After` seconds)
- [ ] Log errors (no sensitive data)

---

## Integration Pattern

### Cover Art Fallback Workflow
```go
// Fetch cover art: Cover Art Archive → Spotify fallback
func (s *MusicService) FetchAlbumCover(albumID uuid.UUID) (string, error) {
    album := s.db.GetAlbum(albumID)
    mbid := album.MusicBrainzID

    // 1. Try Cover Art Archive first
    coverURL := fmt.Sprintf("https://coverartarchive.org/release-group/%s/front-1200", mbid)
    coverPath, err := s.downloadCover(coverURL)
    if err == nil {
        return coverPath, nil
    }

    // 2. Fallback to Spotify
    spotifyAlbum := s.spotifyClient.SearchAlbum(album.ArtistName, album.Title)
    if spotifyAlbum == nil {
        return "", errors.New("cover art not found")
    }

    // 3. Download highest resolution image (640x640)
    coverPath, err = s.downloadCover(spotifyAlbum.Images[0].URL)
    if err != nil {
        return "", err
    }

    return coverPath, nil
}
```

### Popularity Scores Update
```go
// Background job: Update popularity scores (weekly)
func (s *MusicService) UpdatePopularityScores(ctx context.Context) error {
    albums := s.db.GetAlbumsWithoutPopularity()

    for _, album := range albums {
        // Search Spotify by artist + album
        spotifyAlbum := s.spotifyClient.SearchAlbum(album.ArtistName, album.Title)
        if spotifyAlbum == nil {
            continue
        }

        // Update popularity
        s.db.UpdateAlbum(album.ID, map[string]interface{}{
            "popularity":  spotifyAlbum.Popularity,
            "spotify_id":  spotifyAlbum.ID,
        })

        time.Sleep(350 * time.Millisecond) // Rate limit: 180/min = ~333ms per request
    }

    return nil
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
| [Spotify Web API](https://developer.spotify.com/documentation/web-api) | [Local](../../../../sources/apis/spotify.md) |

<!-- SOURCE-BREADCRUMBS-END -->

<!-- DESIGN-BREADCRUMBS-START -->

## Related Design Docs

> Auto-generated cross-references to related design documentation

**Category**: [Music](INDEX.md)

### In This Section

- [Discogs Integration](DISCOGS.md)
- [Last.fm Integration](LASTFM.md)
- [MusicBrainz Integration](MUSICBRAINZ.md)

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

- **Music Module**: [MODULE_IMPLEMENTATION_TODO.md](../../../planning/MODULE_IMPLEMENTATION_TODO.md) (Music section)
- **MusicBrainz Integration**: [MUSICBRAINZ.md](MUSICBRAINZ.md) (primary metadata)
- **Cover Art Archive**: [MUSICBRAINZ.md](MUSICBRAINZ.md) (primary cover art source)
- **Last.fm Integration**: [LASTFM.md](LASTFM.md) (scrobbling + bio)

---

## Notes

- **Client Credentials flow**: Metadata access only (no user-specific data)
- **Token expires**: 3600s (refresh in background job every 3500s)
- **Rate limit**: 180 req/min (use token bucket, ~333ms per request)
- **Cover art quality**: 640x640 (high quality), 300x300 (medium), 64x64 (low)
- **Popularity scores**: 0-100 (based on recent streams, updated frequently)
- **Release precision**: `day` (1997-05-21), `month` (1997-05), `year` (1997)
- **Genres**: Artist-level only (no album/track genres)
- **Search syntax**: Simple query string (artist + album title)
- **External IDs**: Spotify IDs (store in `music_albums.spotify_id`, `music_artists.spotify_id`)
- **API stable**: v1 stable, no breaking changes expected
- **NOT for streaming**: Metadata API only (streaming requires Spotify Premium + separate integration)
- **Free tier**: Unlimited requests (respect rate limits)
- **Error codes**: 401 (Token expired), 403 (Forbidden), 404 (Not found), 429 (Rate limit exceeded)
- **Retry-After header**: Use for 429 errors (seconds to wait)
- **Fallback strategy**: Cover Art Archive (free, primary) → Spotify (fallback)
- **Privacy**: No user data collected (Client Credentials flow)
- **Token storage**: In-memory cache (no persistent storage required)
