# MusicBrainz Integration

> Open music metadata database - primary metadata provider for music

**Service**: MusicBrainz
**Type**: Metadata Provider (Music)
**API Version**: v2 (REST JSON)
**Website**: https://musicbrainz.org
**API Docs**: https://musicbrainz.org/doc/MusicBrainz_API

---

## Overview

**MusicBrainz** is the open music encyclopedia that provides music metadata (artists, albums, tracks, relationships). It's the **primary metadata source** for the Music module.

**Why MusicBrainz**:
- Comprehensive music database (artists, releases, recordings, works)
- Free and open (CC0 license for data)
- Stable API (v2)
- MusicBrainz IDs used by Lidarr
- Cover Art Archive integration
- Relationships (artist collaborations, covers, remixes)

**Use Cases**:
- Artist metadata (biography, genres, formation year)
- Album metadata (release date, track listing, formats)
- Track metadata (duration, ISRC, recordings)
- Cover art (via Cover Art Archive)
- Artist relationships (member of, collaboration, tribute)

---

## Developer Resources

**API Documentation**: https://musicbrainz.org/doc/MusicBrainz_API
**Rate Limiting**: https://musicbrainz.org/doc/MusicBrainz_API/Rate_Limiting
**Search Syntax**: https://musicbrainz.org/doc/Indexed_Search_Syntax

**Authentication**: None (public API)
**Rate Limit**: 1 request/second (REQUIRED User-Agent)
**Free Tier**: Unlimited (respect rate limits)

---

## API Details

### Base URL
```
https://musicbrainz.org/ws/2/
```

### Authentication
None required, but **User-Agent header MANDATORY**:
```
User-Agent: Revenge/1.0.0 (https://github.com/lusoris/revenge)
```

### Rate Limiting
- **1 request/second** (strictly enforced)
- **User-Agent required** (requests without User-Agent are rejected)
- Respect rate limits or risk IP ban

### Key Endpoints

#### Get Artist
```bash
GET /ws/2/artist/{mbid}?inc=aliases+genres+ratings+url-rels
```

**Query Parameters**:
- `inc`: Include aliases, genres, ratings, url-rels (relationships), tags, releases, recordings

**Response**:
```json
{
  "id": "a74b1b7f-71a5-4011-9441-d0b5e4122711",
  "name": "Radiohead",
  "country": "GB",
  "life-span": {"begin": "1985"},
  "type": "Group",
  "genres": [{"name": "alternative rock", "count": 15}],
  "aliases": [{"name": "Radio Head", "locale": "en"}],
  "url-relations": [
    {"type": "official homepage", "url": {"resource": "https://radiohead.com"}}
  ]
}
```

#### Get Release Group (Album)
```bash
GET /ws/2/release-group/{mbid}?inc=artist-credits+releases+url-rels
```

**Response**:
```json
{
  "id": "b1392450-e666-3926-a536-22c65f834433",
  "title": "OK Computer",
  "first-release-date": "1997-05-21",
  "primary-type": "Album",
  "artist-credit": [{"artist": {"id": "a74b...", "name": "Radiohead"}}],
  "releases": [
    {"id": "...", "title": "OK Computer", "date": "1997-05-21", "country": "GB"}
  ]
}
```

#### Get Recording (Track)
```bash
GET /ws/2/recording/{mbid}?inc=artist-credits+isrcs+releases
```

**Response**:
```json
{
  "id": "6f9c8c32-3aae-4dad-b023-56389361cf6b",
  "title": "Karma Police",
  "length": 263000,
  "artist-credit": [{"artist": {"id": "a74b...", "name": "Radiohead"}}],
  "isrcs": ["GBAYE9700455"]
}
```

#### Search Artists
```bash
GET /ws/2/artist?query=radiohead&limit=10&offset=0
```

**Search Syntax**:
- `artist:radiohead` - Artist name
- `country:GB` - Country code
- `type:Group` - Artist type (Person, Group, Orchestra, etc.)

#### Cover Art Archive
```bash
GET https://coverartarchive.org/release-group/{mbid}/front
GET https://coverartarchive.org/release/{mbid}/front-500
```

**Cover Art Sizes**:
- `front` - Original resolution
- `front-250` - 250px thumbnail
- `front-500` - 500px thumbnail
- `front-1200` - 1200px (recommended)

---

## Implementation Checklist

### API Client (`internal/infra/metadata/provider_musicbrainz.go`)
- [ ] Base URL configuration
- [ ] HTTP client with **User-Agent header**
- [ ] Rate limiting (1 req/s with token bucket)
- [ ] Error handling (404, 503, rate limit exceeded)
- [ ] Response parsing (JSON unmarshalling)

### Artist Metadata
- [ ] Fetch artist by MusicBrainz ID
- [ ] Search artists by name
- [ ] Extract: name, country, formation year, type, genres
- [ ] Fetch artist relationships (members, collaborations)
- [ ] Store in `music_artists` table

### Album Metadata
- [ ] Fetch release group by MusicBrainz ID
- [ ] Extract: title, release date, primary type (Album, EP, Single)
- [ ] Fetch releases (different countries/formats)
- [ ] Store in `music_albums` table

### Track Metadata
- [ ] Fetch recording by MusicBrainz ID
- [ ] Extract: title, duration, ISRC
- [ ] Link to releases
- [ ] Store in `music_tracks` table

### Cover Art Handling
- [ ] Fetch cover art from Cover Art Archive
- [ ] Download high-res (front-1200)
- [ ] Generate Blurhash
- [ ] Convert to WebP
- [ ] Store locally (`data/music/covers/`)

### Error Handling
- [ ] Handle 404 (artist/release not found)
- [ ] Handle 503 (service unavailable)
- [ ] Handle rate limit exceeded (retry after delay)
- [ ] Log errors (obfuscated, no sensitive data)

---

## Integration Pattern

### Lidarr Webhook → MusicBrainz Metadata Sync
```go
// Webhook: Lidarr added new album
func (s *MusicService) HandleLidarrAlbumAdded(albumID string) error {
    // 1. Get album from Lidarr
    lidarrAlbum := s.lidarrClient.GetAlbum(albumID)
    mbid := lidarrAlbum.ForeignAlbumId // MusicBrainz ID

    // 2. Fetch metadata from MusicBrainz
    mbAlbum := s.musicbrainzClient.GetReleaseGroup(mbid)
    mbArtist := s.musicbrainzClient.GetArtist(mbAlbum.ArtistCredit[0].Artist.ID)

    // 3. Fetch cover art from Cover Art Archive
    coverURL := fmt.Sprintf("https://coverartarchive.org/release-group/%s/front-1200", mbid)
    coverPath := s.downloadCover(coverURL)

    // 4. Store in Revenge database
    s.db.InsertArtist(mbArtist)
    s.db.InsertAlbum(mbAlbum, coverPath)

    return nil
}
```

---

## Related Documentation

- **Music Module**: [docs/MODULE_IMPLEMENTATION_TODO.md](../../MODULE_IMPLEMENTATION_TODO.md) (Music section)
- **Lidarr Integration**: [../servarr/LIDARR.md](../servarr/LIDARR.md)
- **Last.fm Integration**: [LASTFM.md](LASTFM.md) (scrobbling + artist bio)
- **Spotify Integration**: [SPOTIFY.md](SPOTIFY.md) (cover art fallback)

---

## Notes

- **User-Agent REQUIRED**: All requests MUST include User-Agent header (or 403 Forbidden)
- **Rate limit strictly enforced**: 1 req/s (use token bucket, queue requests)
- **MusicBrainz IDs**: Primary identifier (used by Lidarr, stored in `foreign_artist_id`, `foreign_album_id`)
- **Cover Art Archive**: Free CDN (part of MusicBrainz, same rate limits)
- **Artist types**: Person, Group, Orchestra, Choir, Character, Other
- **Release types**: Album, EP, Single, Broadcast, Compilation, Soundtrack, Live, Remix, DJ-mix, Mixtape/Street, Interview, Audiobook, Audio drama, Demo, Other
- **Relationships**: member of, collaboration, cover, remix, tribute, etc. (fetch with `inc=artist-rels`)
- **Search syntax**: Lucene-based (use `artist:name AND country:GB`)
- **API v2 stable**: No breaking changes expected
- **Free and open**: CC0 license (data can be used freely)
- **Fallback**: Use Last.fm for artist bio/tags, Spotify for cover art if Cover Art Archive fails
