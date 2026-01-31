# Last.fm Scrobbling Integration

<!-- SOURCES: lastfm-api, river -->

<!-- DESIGN: integrations/scrobbling, 01_ARCHITECTURE, 02_DESIGN_PRINCIPLES, 03_METADATA_SYSTEM -->


> Music scrobbling and listening history tracking


<!-- TOC-START -->

## Table of Contents

- [Status](#status)
- [Overview](#overview)
- [Developer Resources](#developer-resources)
  - [API Documentation](#api-documentation)
  - [API Key Setup](#api-key-setup)
  - [Required Parameters](#required-parameters)
  - [API Signature](#api-signature)
- [API Endpoints](#api-endpoints)
  - [Scrobble Track](#scrobble-track)
  - [Update Now Playing](#update-now-playing)
  - [Get Recent Tracks](#get-recent-tracks)
  - [Get Loved Tracks](#get-loved-tracks)
  - [Love/Unlove Track](#loveunlove-track)
- [Implementation Checklist](#implementation-checklist)
  - [Phase 1: Authentication Setup](#phase-1-authentication-setup)
  - [Phase 2: Scrobbling (Real-time)](#phase-2-scrobbling-real-time)
  - [Phase 3: Listening History Sync](#phase-3-listening-history-sync)
  - [Phase 4: Loved Tracks Sync](#phase-4-loved-tracks-sync)
  - [Phase 5: Background Jobs (River)](#phase-5-background-jobs-river)
- [Integration Pattern](#integration-pattern)
  - [Real-time Scrobbling Flow](#real-time-scrobbling-flow)
  - [Listening History Sync Flow](#listening-history-sync-flow)
  - [Loved Tracks Sync Flow](#loved-tracks-sync-flow)
- [Sources & Cross-References](#sources-cross-references)
  - [Cross-Reference Indexes](#cross-reference-indexes)
  - [Referenced Sources](#referenced-sources)
- [Related Design Docs](#related-design-docs)
  - [In This Section](#in-this-section)
  - [Related Topics](#related-topics)
  - [Indexes](#indexes)
- [Related Documentation](#related-documentation)
- [Notes](#notes)
  - [Authentication (Desktop Auth Flow)](#authentication-desktop-auth-flow)
  - [Rate Limits](#rate-limits)
  - [Scrobbling Rules (Last.fm Standard)](#scrobbling-rules-lastfm-standard)
  - [API Signature Calculation](#api-signature-calculation)
  - [Batch Scrobbling](#batch-scrobbling)
  - [Deduplication Strategy](#deduplication-strategy)
  - [MusicBrainz IDs](#musicbrainz-ids)
  - [Loved Tracks vs Favorites](#loved-tracks-vs-favorites)
  - [Artist/Album Metadata Enrichment](#artistalbum-metadata-enrichment)
  - [Error Handling](#error-handling)
  - [Privacy Considerations](#privacy-considerations)
  - [Last.fm vs ListenBrainz](#lastfm-vs-listenbrainz)
  - [Fallback Strategy (Music Scrobbling)](#fallback-strategy-music-scrobbling)

<!-- TOC-END -->

**Service**: Last.fm (https://www.last.fm)
**API**: REST API with API key + session authentication
**Category**: Scrobbling / Music
**Priority**: 🟢 HIGH (Popular music tracking service)

## Status

| Dimension | Status |
|-----------|--------|
| Design | ✅ |
| Sources | ✅ |
| Instructions | ✅ |
| Code | 🔴 |
| Linting | 🔴 |
| Unit Testing | 🔴 |
| Integration Testing | 🔴 |---

## Overview

**Last.fm** is the world's largest music tracking service, allowing users to scrobble (track) their listening history and discover new music through personalized recommendations.

**Key Features**:
- **Music scrobbling**: Automatic tracking of what you listen to
- **Listening history**: Complete history of all tracks played
- **Statistics**: Top artists, albums, tracks, genres
- **Recommendations**: Personalized music recommendations
- **Social features**: Friends, recent tracks, loved tracks
- **Artist/album metadata**: Metadata enrichment

**Use Cases**:
- Automatic music scrobbling (mark tracks as played)
- Sync listening history to Last.fm
- Import Last.fm listening history to Revenge
- Music statistics (top artists, albums, tracks)
- Recommendations based on listening habits

---

## Developer Resources

### API Documentation
- **Base URL**: https://ws.audioscrobbler.com/2.0/
- **Documentation**: https://www.last.fm/api
- **Authentication**: API key + session key (user authorization)
- **Rate Limits**: 5 requests per second (per API key)

### API Key Setup
```
1. Register app: https://www.last.fm/api/account/create
   - Get API key and shared secret

2. User authentication (desktop auth flow):
   - Get auth token: http://www.last.fm/api/auth/?api_key={API_KEY}
   - User authorizes on Last.fm website
   - Get session key: auth.getSession (with token)

3. Session key (long-lived):
   - Store encrypted in database
   - Use for all authenticated API calls
```

### Required Parameters
```
api_key: Your API key
method: API method (e.g., track.scrobble)
format: json (response format)
api_sig: MD5 signature (for authenticated calls)
sk: Session key (for authenticated calls)
```

### API Signature
```
MD5 signature format (authenticated calls only):
1. Sort params alphabetically (exclude format and callback)
2. Concatenate: key1value1key2value2...
3. Append shared secret
4. Calculate MD5 hash

Example:
  Params: {api_key: "abc", method: "track.scrobble", sk: "xyz", timestamp: "123"}
  String: "api_keyabcmethodtrack.scrobbleskxyztimestamp123" + shared_secret
  Signature: MD5(string)
```

---

## API Endpoints

### Scrobble Track
```
POST /2.0/
Parameters:
  method: track.scrobble
  artist: Artist name
  track: Track title
  timestamp: Unix timestamp (when track started playing)
  album: Album name (optional)
  albumArtist: Album artist (optional)
  duration: Track duration in seconds (optional)
  api_key: Your API key
  sk: Session key
  api_sig: MD5 signature

Response:
{
  "scrobbles": {
    "@attr": {
      "accepted": 1,
      "ignored": 0
    },
    "scrobble": {
      "artist": {"#text": "Artist Name"},
      "track": {"#text": "Track Title"},
      "album": {"#text": "Album Name"},
      "timestamp": "1234567890"
    }
  }
}
```

### Update Now Playing
```
POST /2.0/
Parameters:
  method: track.updateNowPlaying
  artist: Artist name
  track: Track title
  album: Album name (optional)
  albumArtist: Album artist (optional)
  duration: Track duration in seconds (optional)
  api_key: Your API key
  sk: Session key
  api_sig: MD5 signature
```

### Get Recent Tracks
```
GET /2.0/?method=user.getRecentTracks&user={username}&api_key={api_key}&format=json

Response:
{
  "recenttracks": {
    "track": [
      {
        "artist": {"#text": "Artist Name"},
        "name": "Track Title",
        "album": {"#text": "Album Name"},
        "date": {"uts": "1234567890"},
        "mbid": "musicbrainz_id"
      }
    ]
  }
}
```

### Get Loved Tracks
```
GET /2.0/?method=user.getLovedTracks&user={username}&api_key={api_key}&format=json
```

### Love/Unlove Track
```
POST /2.0/
Parameters:
  method: track.love (or track.unlove)
  artist: Artist name
  track: Track title
  api_key: Your API key
  sk: Session key
  api_sig: MD5 signature
```

---

## Implementation Checklist

### Phase 1: Authentication Setup
- [ ] API key configuration (store in config)
- [ ] Desktop authentication flow (auth token → session key)
- [ ] Session key storage (encrypt in database)
- [ ] API signature generation (MD5 hash)

### Phase 2: Scrobbling (Real-time)
- [ ] **Update now playing** (POST track.updateNowPlaying when track starts)
- [ ] **Scrobble track** (POST track.scrobble when track finishes or >= 50% played)
- [ ] Timestamp tracking (Unix timestamp when track started)
- [ ] Duration tracking (track length in seconds)
- [ ] Real-time sync (sync on playback events)

### Phase 3: Listening History Sync
- [ ] **Export to Last.fm** (POST track.scrobble batch - send Revenge listen history to Last.fm)
- [ ] **Import from Last.fm** (GET user.getRecentTracks - fetch Last.fm listen history)
- [ ] Bi-directional sync (merge listen histories)
- [ ] Deduplication (handle duplicate scrobbles)

### Phase 4: Loved Tracks Sync
- [ ] **Export loved tracks to Last.fm** (POST track.love)
- [ ] **Import loved tracks from Last.fm** (GET user.getLovedTracks)
- [ ] Bi-directional sync (merge loved tracks)

### Phase 5: Background Jobs (River)
- [ ] **Job**: `scrobble.lastfm.sync_history` (periodic history sync)
- [ ] **Job**: `scrobble.lastfm.sync_loved_tracks` (periodic loved tracks sync)
- [ ] Rate limiting (5 req/sec)
- [ ] Retry logic (exponential backoff)

---

## Integration Pattern

### Real-time Scrobbling Flow
```
User starts playing music track
        ↓
Playback session starts (internal/service/playback/session.go)
        ↓
Check if user has Last.fm enabled (user.integrations.lastfm.enabled)
        ↓
        YES
        ↓
Send "now playing" to Last.fm:
  POST /2.0/
  {
    method: "track.updateNowPlaying",
    artist: "Artist Name",
    track: "Track Title",
    album: "Album Name",
    duration: 240
  }
        ↓
Track continues playing
        ↓
Track finishes OR >= 50% played (Last.fm scrobbling rule):
  POST /2.0/
  {
    method: "track.scrobble",
    artist: "Artist Name",
    track: "Track Title",
    album: "Album Name",
    timestamp: 1234567890, // Unix timestamp when track started
    duration: 240
  }
        ↓
Last.fm marks track as scrobbled (added to listening history)
```

### Listening History Sync Flow
```
User enables Last.fm integration (Settings → Integrations → Last.fm → Connect)
        ↓
Desktop authentication flow:
  1. Redirect to Last.fm authorization page (http://www.last.fm/api/auth/?api_key={API_KEY})
  2. User authorizes Revenge app
  3. Get session key (POST auth.getSession with token)
  4. Store session key (encrypted) in users.integrations.lastfm.session_key
        ↓
Initial sync:
  1. Fetch Last.fm listen history (GET user.getRecentTracks)
  2. Import to Revenge (create music_listen_history entries)
  3. Fetch Revenge listen history
  4. Export to Last.fm (POST track.scrobble batch)
  5. Merge & deduplicate
        ↓
Ongoing sync:
  - Real-time scrobbling (POST track.updateNowPlaying, track.scrobble)
  - Periodic sync (River job every 1 hour) → fetch Last.fm updates → merge
```

### Loved Tracks Sync Flow
```
User loves track in Revenge (click heart icon)
        ↓
Store in Revenge database (music_user_favorites table)
        ↓
Check if user has Last.fm enabled
        ↓
        YES
        ↓
Export to Last.fm:
  POST /2.0/
  {
    method: "track.love",
    artist: "Artist Name",
    track: "Track Title"
  }
        ↓
Last.fm marks track as loved
```

---


## Related Documentation

- [LISTENBRAINZ.md](./LISTENBRAINZ.md) - ListenBrainz scrobbling (open-source alternative)
- [TRAKT.md](./TRAKT.md) - Trakt scrobbling (movies/TV)
- [MUSICBRAINZ.md](../metadata/music/MUSICBRAINZ.md) - MusicBrainz metadata
- [SPOTIFY.md](../metadata/music/SPOTIFY.md) - Spotify integration

---

## Notes

### Authentication (Desktop Auth Flow)
- **API key**: Public (identify Revenge app)
- **Shared secret**: Private (sign API requests)
- **Session key**: User-specific (long-lived token)
- **Token expiry**: Session keys don't expire (unless user revokes)
- **Encryption**: Encrypt session keys in database

### Rate Limits
- **Limit**: 5 requests per second (per API key)
- **Headers**: No rate limit headers (monitor 429 responses)
- **Throttling**: Implement token bucket rate limiter (5 req/sec)
- **Retry**: Retry with exponential backoff on 429 (rate limit exceeded)

### Scrobbling Rules (Last.fm Standard)
- **Minimum duration**: Track must be >= 30 seconds
- **Scrobble threshold**: Track played for >= 50% OR >= 4 minutes (whichever comes first)
- **Timestamp**: Unix timestamp when track STARTED playing (NOT when scrobble sent)
- **Now playing**: Update immediately when track starts (no threshold)

### API Signature Calculation
```go
func calculateAPISignature(params map[string]string, secret string) string {
    // 1. Remove 'format' and 'callback' params
    delete(params, "format")
    delete(params, "callback")

    // 2. Sort keys alphabetically
    keys := make([]string, 0, len(params))
    for k := range params {
        keys = append(keys, k)
    }
    sort.Strings(keys)

    // 3. Concatenate key-value pairs
    var builder strings.Builder
    for _, k := range keys {
        builder.WriteString(k)
        builder.WriteString(params[k])
    }

    // 4. Append shared secret
    builder.WriteString(secret)

    // 5. Calculate MD5 hash
    hash := md5.Sum([]byte(builder.String()))
    return hex.EncodeToString(hash[:])
}
```

### Batch Scrobbling
- **Batch size**: Up to 50 tracks per request
- **Parameters**: Use indexed params (artist[0], track[0], timestamp[0], artist[1], etc.)
- **Use case**: Initial sync (import Revenge listen history to Last.fm)

### Deduplication Strategy
- **Duplicate detection**: Match by artist + track + timestamp (within 5 minutes)
- **Conflict resolution**: Keep Last.fm timestamp if earlier, Revenge if later
- **Merge strategy**: Union (both Last.fm and Revenge entries)

### MusicBrainz IDs
- **Last.fm supports**: MusicBrainz IDs (artist, album, track)
- **Enhanced matching**: Send MusicBrainz IDs for better accuracy
- **Parameter**: `mbid` (MusicBrainz ID)

### Loved Tracks vs Favorites
- **Last.fm "loved"**: Equivalent to Revenge "favorites"
- **Bi-directional sync**: Loved tracks ↔ favorites
- **Icon**: Heart icon (same in both systems)

### Artist/Album Metadata Enrichment
- **Last.fm metadata**: Artist bio, similar artists, top tracks, top albums
- **Endpoints**: artist.getInfo, album.getInfo, artist.getSimilar
- **Use case**: Supplement MusicBrainz metadata with Last.fm data

### Error Handling
- **6 - Invalid parameters**: Check required params (artist, track, timestamp)
- **9 - Invalid session key**: Session key expired/revoked → re-authenticate
- **11 - Service offline**: Last.fm down → retry later
- **16 - Service temporarily unavailable**: Retry with exponential backoff
- **29 - Rate limit exceeded**: Throttle requests (5 req/sec)

### Privacy Considerations
- **User opt-in**: Users must explicitly enable Last.fm integration
- **Desktop auth**: Users authorize Revenge to scrobble on their behalf
- **Data visibility**: Last.fm scrobbles are public by default (users can change in Last.fm settings)
- **Session key storage**: Encrypt session keys in database

### Last.fm vs ListenBrainz
- **Last.fm**: Commercial, largest user base, social features, recommendations
- **ListenBrainz**: Open-source, MusicBrainz-backed, privacy-focused
- **Both**: Support both integrations (user choice)

### Fallback Strategy (Music Scrobbling)
- **Order**: Last.fm (primary) → ListenBrainz (alternative/supplement)
