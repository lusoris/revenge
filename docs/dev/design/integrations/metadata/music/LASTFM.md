# Last.fm Integration

> Music scrobbling and metadata provider - artist bio, tags, similar artists

**Service**: Last.fm
**Type**: Scrobbling + Metadata Provider (Music)
**API Version**: 2.0
**Website**: https://www.last.fm
**API Docs**: https://www.last.fm/api

---

## Overview

**Last.fm** provides music scrobbling (listening history tracking) and music metadata (artist bio, tags, similar artists, album info).

**Why Last.fm**:
- Scrobbling (track listening history)
- Artist biography and tags
- Similar artists recommendations
- User listening statistics
- Free API (generous limits)
- OAuth authentication

**Use Cases**:
- **Scrobbling**: Track user's music listening (sync to Last.fm profile)
- **Metadata enrichment**: Artist bio, tags, similar artists (fallback to MusicBrainz)
- **Recommendations**: Similar artists, top tracks, user stats
- **Social features**: User profiles, friends, charts

---

## Developer Resources

**API Documentation**: https://www.last.fm/api/intro
**Authentication**: https://www.last.fm/api/authentication
**Scrobbling**: https://www.last.fm/api/scrobbling

**Authentication**: API Key + OAuth (for scrobbling)
**Rate Limit**: Reasonable (no strict limit, respect fair use)
**Free Tier**: Unlimited (API key required)

---

## API Details

### Base URL
```
https://ws.audioscrobbler.com/2.0/
```

### Authentication
- **API Key**: Required for all requests (query param `api_key`)
- **OAuth**: Required for scrobbling (write operations)

**API Key Registration**: https://www.last.fm/api/account/create

### OAuth Flow (for Scrobbling)
```
1. Get API Key + Shared Secret (from Last.fm)
2. Generate auth token: http://www.last.fm/api/auth/?api_key={API_KEY}
3. User authorizes → redirected back with token
4. Get session key: api_key={KEY}&method=auth.getSession&token={TOKEN}&api_sig={SIGNATURE}
5. Use session key for scrobbling
```

### Rate Limiting
- No strict rate limit (fair use)
- Scrobbling: Max 50 tracks per request
- Cache responses where possible

### Key Endpoints

#### Get Artist Info
```bash
GET /2.0/?method=artist.getinfo&artist=Radiohead&api_key={API_KEY}&format=json
```

**Response**:
```json
{
  "artist": {
    "name": "Radiohead",
    "mbid": "a74b1b7f-71a5-4011-9441-d0b5e4122711",
    "url": "https://www.last.fm/music/Radiohead",
    "image": [
      {"#text": "https://...", "size": "large"}
    ],
    "bio": {
      "summary": "Radiohead are an English rock band...",
      "content": "Full biography..."
    },
    "tags": {
      "tag": [
        {"name": "alternative rock", "url": "..."},
        {"name": "experimental", "url": "..."}
      ]
    },
    "similar": {
      "artist": [
        {"name": "Thom Yorke", "url": "..."}
      ]
    }
  }
}
```

#### Get Album Info
```bash
GET /2.0/?method=album.getinfo&artist=Radiohead&album=OK+Computer&api_key={API_KEY}&format=json
```

#### Search Artists
```bash
GET /2.0/?method=artist.search&artist=radiohead&api_key={API_KEY}&format=json
```

#### Scrobble Track (requires OAuth session)
```bash
POST /2.0/
artist=Radiohead&track=Karma+Police&timestamp=1234567890&api_key={API_KEY}&api_sig={SIGNATURE}&sk={SESSION_KEY}&method=track.scrobble
```

**Required Parameters**:
- `artist`: Artist name
- `track`: Track name
- `timestamp`: Unix timestamp when track started playing
- `api_key`: API key
- `sk`: Session key (from OAuth)
- `api_sig`: MD5 signature

**API Signature Calculation**:
```
MD5(api_key{API_KEY}artist{ARTIST}method{METHOD}sk{SESSION_KEY}timestamp{TIMESTAMP}track{TRACK}{SHARED_SECRET})
```

#### Update Now Playing
```bash
POST /2.0/
artist=Radiohead&track=Karma+Police&api_key={API_KEY}&api_sig={SIGNATURE}&sk={SESSION_KEY}&method=track.updateNowPlaying
```

---

## Implementation Checklist

### API Client (`internal/infra/scrobble/provider_lastfm.go`)
- [ ] Base URL configuration
- [ ] API key configuration
- [ ] OAuth flow (session key generation)
- [ ] API signature generation (MD5)
- [ ] Error handling (6: Invalid parameters, 9: Invalid session key, 11: Service offline)
- [ ] Response parsing (JSON unmarshalling)

### Artist Metadata
- [ ] Fetch artist info (bio, tags, similar artists)
- [ ] Extract: bio, tags, similar artists, images
- [ ] Store in `music_artists` table (enrich MusicBrainz data)
- [ ] Cache responses (1 day TTL)

### Scrobbling
- [ ] OAuth flow (user authorization)
- [ ] Store session key per user
- [ ] Scrobble track (when user plays music)
- [ ] Update now playing (real-time)
- [ ] Batch scrobbling (max 50 tracks per request)
- [ ] Retry on failure (network errors)

### User Stats
- [ ] Fetch user's top artists
- [ ] Fetch user's top tracks
- [ ] Fetch recent tracks
- [ ] Display in user profile (optional)

### Error Handling
- [ ] Handle 6 (Invalid parameters)
- [ ] Handle 9 (Invalid session key - re-authenticate)
- [ ] Handle 11 (Service offline - retry)
- [ ] Log errors (no sensitive data)

---

## Integration Pattern

### Scrobbling Workflow
```go
// User plays track → Scrobble to Last.fm
func (s *ScrobbleService) ScrobbleTrack(ctx context.Context, userID uuid.UUID, trackID uuid.UUID) error {
    // 1. Get user's Last.fm session key
    sessionKey := s.db.GetUserLastFmSession(userID)
    if sessionKey == "" {
        return errors.New("user not connected to Last.fm")
    }

    // 2. Get track metadata
    track := s.db.GetTrack(trackID)

    // 3. Update now playing
    s.lastfmClient.UpdateNowPlaying(sessionKey, track.Artist, track.Title)

    // 4. Scrobble after 30 seconds OR 50% playback
    time.AfterFunc(30*time.Second, func() {
        timestamp := time.Now().Unix()
        s.lastfmClient.Scrobble(sessionKey, track.Artist, track.Title, timestamp)
    })

    return nil
}
```

### Metadata Enrichment
```go
// Enrich artist with Last.fm bio/tags
func (s *MusicService) EnrichArtistMetadata(artistID uuid.UUID) error {
    artist := s.db.GetArtist(artistID)

    // Fetch from Last.fm
    lastfmArtist := s.lastfmClient.GetArtistInfo(artist.Name)

    // Update artist with bio/tags
    s.db.UpdateArtist(artistID, map[string]interface{}{
        "biography":      lastfmArtist.Bio.Content,
        "tags":           lastfmArtist.Tags,
        "similar_artists": lastfmArtist.Similar,
    })

    return nil
}
```

---

## Related Documentation

- **Music Module**: [MODULE_IMPLEMENTATION_TODO.md](../../../planning/MODULE_IMPLEMENTATION_TODO.md) (Music section)
- **MusicBrainz Integration**: [MUSICBRAINZ.md](MUSICBRAINZ.md) (primary metadata)
- **Scrobbling Overview**: [docs/SCROBBLING.md](../../../SCROBBLING.md)
- **ListenBrainz Integration**: [../../scrobbling/LISTENBRAINZ.md](../../scrobbling/LISTENBRAINZ.md) (open alternative)

---

## Notes

- **API Key required**: Register at https://www.last.fm/api/account/create
- **OAuth for scrobbling**: Write operations require user authorization
- **Scrobble rules**: Min 30 seconds playback OR 50% of track (whichever comes first)
- **Batch scrobbling**: Max 50 tracks per request (useful for offline sync)
- **API signature**: MD5 hash of sorted params + shared secret
- **MusicBrainz IDs**: Last.fm returns MusicBrainz IDs (cross-reference)
- **Artist bio source**: Wikipedia-based (good for enrichment)
- **Tags**: User-generated (folksonomy, useful for genre discovery)
- **Similar artists**: Based on listening habits (recommendations)
- **Rate limiting**: No strict limit, but respect fair use (cache responses)
- **Error codes**: 6 (Invalid parameters), 9 (Invalid session key), 11 (Service offline), 26 (Suspended API key)
- **Fallback**: Use MusicBrainz for core metadata, Last.fm for bio/tags/similar artists
- **User privacy**: Scrobbling opt-in (user must authorize)
- **Session key storage**: Store per user in `user_integrations` table (encrypted)
- **Scrobbling alternatives**: ListenBrainz (open source, Last.fm compatible API)
