# Lidarr Integration

> Music management automation

**Status**: 🟡 PLANNED
**Priority**: 🟡 HIGH (Phase 4 - Music Module)
**Type**: Webhook listener + API client for metadata sync

---

## Overview

Lidarr is the industry-standard music management automation tool. Revenge integrates with Lidarr to:
- Receive webhook notifications when albums/tracks are imported
- Sync artist, album, and track metadata
- Monitor Lidarr download/import status
- Map Lidarr quality profiles to Revenge quality tiers

**Integration Points**:
- **Webhook listener**: Process Lidarr events (On Import, On Album Added, etc.)
- **API client**: Query artists, albums, tracks
- **Metadata sync**: Enrich Revenge metadata with Lidarr data
- **Quality mapping**: Lidarr quality profiles → Revenge audio quality tiers

---

## Developer Resources

- 📚 **API Docs**: https://lidarr.audio/docs/api/
- 🔗 **OpenAPI Spec**: https://github.com/Lidarr/Lidarr/blob/develop/src/Lidarr.Api.V1/openapi.json
- 🔗 **GitHub**: https://github.com/Lidarr/Lidarr
- 🔗 **Wiki**: https://wiki.servarr.com/lidarr

---

## API Details

**Base Path**: `/api/v1/`
**Authentication**: `X-Api-Key` header (API key from Lidarr settings)
**Rate Limits**: None (self-hosted)

### Key Endpoints

| Endpoint | Method | Purpose |
|----------|--------|---------|
| `/artist` | GET | List all artists |
| `/artist/{id}` | GET | Get specific artist details |
| `/album` | GET | List albums (filterable by artist) |
| `/album/{id}` | GET | Get specific album details |
| `/track` | GET | List tracks (filterable by album) |
| `/track/{id}` | GET | Get specific track details |
| `/importlist` | GET | List configured import lists |
| `/metadata` | GET | Get metadata settings |
| `/qualityprofile` | GET | List quality profiles |
| `/metadataprofile` | GET | List metadata profiles |
| `/system/status` | GET | Get Lidarr version & status |
| `/health` | GET | Check Lidarr health |

---

## Webhook Events

Lidarr can send webhooks for the following events:

### On Import (Album Downloaded & Imported)
```json
{
  "eventType": "Download",
  "artist": {
    "id": 1,
    "name": "Radiohead",
    "foreignArtistId": "a74b1b7f-71a5-4011-9441-d0b5e4122711",  // MusicBrainz ID
    "path": "/media/Music/Radiohead"
  },
  "album": {
    "id": 123,
    "title": "OK Computer",
    "releaseDate": "1997-05-21",
    "foreignAlbumId": "e7e0490e-07a8-3419-b527-2baa90616e83",  // MusicBrainz ID
    "overview": "Radiohead's third studio album...",
    "images": [
      {
        "coverType": "cover",
        "url": "https://coverartarchive.org/release/e7e0490e-07a8-3419-b527-2baa90616e83/front.jpg"
      }
    ]
  },
  "tracks": [
    {
      "id": 1,
      "trackNumber": "1",
      "title": "Airbag",
      "duration": "00:04:44",
      "trackFile": {
        "id": 456,
        "relativePath": "OK Computer/01 - Airbag.flac",
        "quality": "FLAC",
        "size": 32145728
      }
    }
  ]
}
```

### On Album Added (New Album Tracked)
Triggered when Lidarr starts monitoring a new album.

### On Track File Delete
Triggered when track file is deleted from Lidarr.

### On Artist Delete
Triggered when artist is removed from Lidarr.

### On Rename
Triggered when track files are renamed.

### On Health Issue
Triggered when Lidarr detects health issues.

---

## Implementation Checklist

- [ ] **API Client** (`internal/service/metadata/provider_lidarr.go`)
  - [ ] Artist listing & detail fetching
  - [ ] Album listing & detail fetching
  - [ ] Track listing & detail fetching
  - [ ] Quality profile mapping
  - [ ] Metadata profile mapping
  - [ ] Health check integration

- [ ] **Webhook Handler** (`internal/api/handlers/webhook_lidarr.go`)
  - [ ] Parse webhook payload (On Download event)
  - [ ] Extract artist + album + track metadata
  - [ ] Trigger metadata enrichment (MusicBrainz, Last.fm)
  - [ ] Store in PostgreSQL (`music_artists`, `music_albums`, `music_tracks`)
  - [ ] Update Typesense search index

- [ ] **Metadata Sync**
  - [ ] Map Lidarr artists → Revenge `music_artists` table
  - [ ] Map Lidarr albums → Revenge `music_albums` table
  - [ ] Map Lidarr tracks → Revenge `music_tracks` table
  - [ ] Map Lidarr quality profiles → Revenge audio quality tiers
  - [ ] Handle multi-disc albums

- [ ] **Quality Profile Mapping**
  - [ ] Lossless (FLAC) → `quality='lossless'`, `bitrate=variable`
  - [ ] High (320kbps MP3) → `quality='high'`, `bitrate=320`
  - [ ] Medium (192kbps MP3) → `quality='medium'`, `bitrate=192`
  - [ ] Low (128kbps MP3) → `quality='low'`, `bitrate=128`

- [ ] **Error Handling**
  - [ ] Retry failed API calls (circuit breaker)
  - [ ] Log webhook failures
  - [ ] Handle missing tracks (not yet released)

---

## Revenge Integration Pattern

```
Lidarr imports album (OK Computer)
           ↓
Sends webhook to Revenge
           ↓
Revenge processes webhook
           ↓
Stores artist/album/tracks in PostgreSQL (music_artists, music_albums, music_tracks)
           ↓
Enriches metadata from MusicBrainz (artist bio, genres)
           ↓
Enriches metadata from Last.fm (play counts, similar artists)
           ↓
Updates Typesense search index
           ↓
Album available for playback
```

### Go Client Example

```go
type LidarrClient struct {
    baseURL string
    apiKey  string
    client  *http.Client
}

func (c *LidarrClient) GetArtist(ctx context.Context, artistID int) (*Artist, error) {
    url := fmt.Sprintf("%s/api/v1/artist/%d", c.baseURL, artistID)
    req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
    req.Header.Set("X-Api-Key", c.apiKey)
    
    resp, err := c.client.Do(req)
    if err != nil {
        return nil, fmt.Errorf("failed to get artist: %w", err)
    }
    defer resp.Body.Close()
    
    var artist Artist
    json.NewDecoder(resp.Body).Decode(&artist)
    return &artist, nil
}

func (c *LidarrClient) GetAlbumsByArtist(ctx context.Context, artistID int) ([]Album, error) {
    url := fmt.Sprintf("%s/api/v1/album?artistId=%d", c.baseURL, artistID)
    req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
    req.Header.Set("X-Api-Key", c.apiKey)
    
    resp, err := c.client.Do(req)
    if err != nil {
        return nil, fmt.Errorf("failed to get albums: %w", err)
    }
    defer resp.Body.Close()
    
    var albums []Album
    json.NewDecoder(resp.Body).Decode(&albums)
    return albums, nil
}
```

---

## Related Documentation

- [Radarr Integration](RADARR.md) - Similar workflow for movies
- [Sonarr Integration](SONARR.md) - Similar workflow for TV shows
- [Music Module](../../architecture/modules/MUSIC.md)
- [MusicBrainz Integration](../metadata/music/MUSICBRAINZ.md)
- [Last.fm Integration](../metadata/music/LASTFM.md)
- [Arr Integration Pattern](../../patterns/arr_integration.md)
- [Webhook Handling](../../patterns/webhook_patterns.md)

---

## Quality Profile Mapping

| Lidarr Quality | Revenge Quality | Bitrate | Format |
|----------------|-----------------|---------|--------|
| FLAC | `lossless` | Variable (avg 1000+ kbps) | FLAC |
| ALAC | `lossless` | Variable (avg 1000+ kbps) | ALAC |
| MP3-320 | `high` | 320 kbps | MP3 (CBR) |
| MP3-V0 | `high` | 220-260 kbps | MP3 (VBR) |
| MP3-256 | `high` | 256 kbps | MP3 (CBR) |
| AAC-256 | `high` | 256 kbps | AAC |
| MP3-192 | `medium` | 192 kbps | MP3 (CBR) |
| AAC-192 | `medium` | 192 kbps | AAC |
| MP3-128 | `low` | 128 kbps | MP3 (CBR) |
| AAC-128 | `low` | 128 kbps | AAC |
| Any | `auto` | Varies | Varies |

---

## Notes

- **MusicBrainz is primary metadata source** (consistency with Lidarr)
- Lidarr API v1 is stable (widely adopted)
- Self-hosted = no rate limits (unlike cloud APIs)
- Quality profiles are customizable in Lidarr (respect user settings)
- Lidarr uses MusicBrainz IDs (`foreignArtistId`, `foreignAlbumId`)
- Multi-disc albums: Lidarr tracks disc number + track number
- Metadata profiles: Control which release types to monitor (studio albums, live, compilations, etc.)
- Release date: Lidarr uses earliest release date from MusicBrainz
- Wanted missing: Lidarr tracks "monitored" status (not yet released = no file)
