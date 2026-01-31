# ListenBrainz Integration

> Open-source music listening history tracking (MusicBrainz project)


<!-- TOC-START -->

## Table of Contents

- [Status](#status)
- [Overview](#overview)
- [Developer Resources](#developer-resources)
  - [API Documentation](#api-documentation)
  - [User Token Setup](#user-token-setup)
  - [Required Headers](#required-headers)
- [API Endpoints](#api-endpoints)
  - [Submit Listens (Scrobble)](#submit-listens-scrobble)
  - [Get Listens (Recent Tracks)](#get-listens-recent-tracks)
  - [Get Playing Now](#get-playing-now)
  - [Get Statistics](#get-statistics)
  - [Delete Listen](#delete-listen)
- [Implementation Checklist](#implementation-checklist)
  - [Phase 1: Client Setup](#phase-1-client-setup)
  - [Phase 2: API Implementation](#phase-2-api-implementation)
  - [Phase 3: Service Integration](#phase-3-service-integration)
  - [Phase 4: Testing](#phase-4-testing)
- [Integration Pattern](#integration-pattern)
  - [Real-time Scrobbling Flow](#real-time-scrobbling-flow)
  - [Listening History Sync Flow](#listening-history-sync-flow)
- [Sources & Cross-References](#sources-cross-references)
  - [Cross-Reference Indexes](#cross-reference-indexes)
  - [Referenced Sources](#referenced-sources)
- [Related Design Docs](#related-design-docs)
  - [In This Section](#in-this-section)
  - [Related Topics](#related-topics)
  - [Indexes](#indexes)
- [Related Documentation](#related-documentation)
- [Notes](#notes)
  - [Authentication (Simple Token)](#authentication-simple-token)
  - [Rate Limits](#rate-limits)
  - [Listen Types](#listen-types)
  - [MusicBrainz ID Integration](#musicbrainz-id-integration)
  - [Scrobbling Rules (Same as Last.fm)](#scrobbling-rules-same-as-lastfm)
  - [Batch Import](#batch-import)
  - [Pagination (Fetch History)](#pagination-fetch-history)
  - [Deduplication Strategy](#deduplication-strategy)
  - [Statistics (Top Artists/Releases/Recordings)](#statistics-top-artistsreleasesrecordings)
  - [Delete Listens](#delete-listens)
  - [Privacy & Open Data](#privacy-open-data)
  - [ListenBrainz vs Last.fm](#listenbrainz-vs-lastfm)
  - [Error Handling](#error-handling)
  - [Fallback Strategy (Music Scrobbling)](#fallback-strategy-music-scrobbling)

<!-- TOC-END -->

**Service**: ListenBrainz (https://listenbrainz.org)
**API**: REST API with user token authentication
**Category**: Scrobbling / Music
**Priority**: 🟡 MEDIUM (Open-source Last.fm alternative)

## Status

| Dimension | Status |
|-----------|--------|
| Design | ✅ |
| Sources | ✅ |
| Instructions | 🟡 |
| Code | 🔴 |
| Linting | 🔴 |
| Unit Testing | 🔴 |
| Integration Testing | 🔴 |
---

## Overview

**ListenBrainz** is an open-source music tracking service by MetaBrainz (creators of MusicBrainz). It tracks listening history and provides statistics, recommendations, and social features—all while respecting user privacy.

**Key Features**:
- **Music scrobbling**: Track what you listen to
- **Listening history**: Complete history of all tracks played
- **Statistics**: Top artists, releases, recordings
- **Recommendations**: Personalized music recommendations
- **MusicBrainz integration**: Native MusicBrainz ID support
- **Open data**: Export/import listening history (JSON format)
- **Privacy-focused**: Open-source, no ads, user data ownership

**Use Cases**:
- Automatic music scrobbling (open-source alternative to Last.fm)
- Sync listening history to ListenBrainz
- Import ListenBrainz listening history to Revenge
- Music statistics (top artists, releases, recordings)
- Privacy-focused music tracking

---

## Developer Resources

### API Documentation
- **Base URL**: https://api.listenbrainz.org/1/
- **Documentation**: https://listenbrainz.readthedocs.io/en/latest/dev/api/
- **Authentication**: User token (simple token-based auth)
- **Rate Limits**: 10 requests per second (per user token)

### User Token Setup
```
1. User creates account: https://listenbrainz.org/register
2. Generate user token: https://listenbrainz.org/profile/ → "User token"
3. Store token in Revenge (encrypted in database)
4. Use token in Authorization header: "Token {user_token}"
```

### Required Headers
```
Content-Type: application/json
Authorization: Token {user_token}
```

---

## API Endpoints

### Submit Listens (Scrobble)
```
POST /1/submit-listens
Headers:
  Authorization: Token {user_token}
Body:
{
  "listen_type": "single", // or "playing_now", "import"
  "payload": [
    {
      "listened_at": 1234567890, // Unix timestamp (omit for "playing_now")
      "track_metadata": {
        "artist_name": "Artist Name",
        "track_name": "Track Title",
        "release_name": "Album Name",
        "additional_info": {
          "recording_mbid": "musicbrainz_recording_id",
          "artist_mbids": ["musicbrainz_artist_id"],
          "release_mbid": "musicbrainz_release_id",
          "duration_ms": 240000,
          "tracknumber": 3
        }
      }
    }
  ]
}

Response:
{
  "status": "ok"
}
```

### Get Listens (Recent Tracks)
```
GET /1/user/{username}/listens?count=25&max_ts={max_timestamp}&min_ts={min_timestamp}

Response:
{
  "payload": {
    "count": 25,
    "latest_listen_ts": 1234567890,
    "listens": [
      {
        "listened_at": 1234567890,
        "track_metadata": {
          "artist_name": "Artist Name",
          "track_name": "Track Title",
          "release_name": "Album Name",
          "additional_info": {
            "recording_mbid": "musicbrainz_recording_id",
            "duration_ms": 240000
          }
        }
      }
    ]
  }
}
```

### Get Playing Now
```
GET /1/user/{username}/playing-now

Response:
{
  "payload": {
    "count": 1,
    "listens": [
      {
        "track_metadata": {
          "artist_name": "Artist Name",
          "track_name": "Track Title",
          "release_name": "Album Name"
        }
      }
    ]
  }
}
```

### Get Statistics
```
GET /1/stats/user/{username}/artists?range=all_time
GET /1/stats/user/{username}/releases?range=all_time
GET /1/stats/user/{username}/recordings?range=all_time

Ranges: week, month, year, all_time

Response:
{
  "payload": {
    "artists": [
      {
        "artist_name": "Artist Name",
        "artist_mbid": "musicbrainz_id",
        "listen_count": 1234
      }
    ]
  }
}
```

### Delete Listen
```
POST /1/delete-listen
Headers:
  Authorization: Token {user_token}
Body:
{
  "listened_at": 1234567890,
  "recording_mbid": "musicbrainz_recording_id"
}
```

---

## Implementation Checklist

### Phase 1: Client Setup
- [ ] Create client package structure
- [ ] Implement HTTP client
- [ ] Add token authentication (user token input, storage, validation)
- [ ] Implement rate limiting (10 req/sec)

### Phase 2: API Implementation
- [ ] Implement scrobble submission (playing now + single)
- [ ] Add history sync (import/export with pagination)
- [ ] Implement error handling (401, 429, 500 responses)

### Phase 3: Service Integration
- [ ] Create ListenBrainz service wrapper
- [ ] Add user preference storage (enable/disable)
- [ ] Implement playback event hooks (submit on playback)

### Phase 4: Testing
- [ ] Add unit tests (token validation, API mapping)
- [ ] Add integration tests (full scrobbling flow)

---

## Integration Pattern

### Real-time Scrobbling Flow
```
User starts playing music track
        ↓
Playback session starts (internal/service/playback/session.go)
        ↓
Check if user has ListenBrainz enabled (user.integrations.listenbrainz.enabled)
        ↓
        YES
        ↓
Lookup MusicBrainz IDs (recording_mbid, artist_mbids, release_mbid)
        ↓
Send "playing now" to ListenBrainz:
  POST /1/submit-listens
  {
    "listen_type": "playing_now",
    "payload": [{
      "track_metadata": {
        "artist_name": "Artist Name",
        "track_name": "Track Title",
        "release_name": "Album Name",
        "additional_info": {
          "recording_mbid": "abc123",
          "artist_mbids": ["def456"],
          "release_mbid": "ghi789",
          "duration_ms": 240000
        }
      }
    }]
  }
        ↓
Track continues playing
        ↓
Track finishes OR >= 50% played:
  POST /1/submit-listens
  {
    "listen_type": "single",
    "payload": [{
      "listened_at": 1234567890, // Unix timestamp when track started
      "track_metadata": {...}
    }]
  }
        ↓
ListenBrainz records listen (added to listening history)
```

### Listening History Sync Flow
```
User enables ListenBrainz integration (Settings → Integrations → ListenBrainz → Enter Token)
        ↓
Validate token (GET /1/validate-token)
        ↓
Store token (encrypted) in users.integrations.listenbrainz.user_token
        ↓
Initial sync:
  1. Fetch ListenBrainz listen history (GET /1/user/{username}/listens with pagination)
  2. Import to Revenge (create music_listen_history entries)
  3. Fetch Revenge listen history
  4. Export to ListenBrainz (POST /1/submit-listens with listen_type="import" batch)
  5. Merge & deduplicate
        ↓
Ongoing sync:
  - Real-time scrobbling (POST /1/submit-listens with listen_type="playing_now", "single")
  - Periodic sync (River job every 1 hour) → fetch ListenBrainz updates → merge
```

---


<!-- SOURCE-BREADCRUMBS-START -->

## Sources & Cross-References

> Auto-generated section linking to external documentation sources

### Cross-Reference Indexes

- [All Sources Index](../../../sources/SOURCES_INDEX.md) - Complete list of external documentation
- [Design ↔ Sources Map](../../../sources/DESIGN_CROSSREF.md) - Which docs reference which sources

### Referenced Sources

| Source | Documentation |
|--------|---------------|
| [Last.fm API](https://www.last.fm/api/intro) | [Local](../../../sources/apis/lastfm.md) |
| [River Job Queue](https://pkg.go.dev/github.com/riverqueue/river) | [Local](../../../sources/tooling/river.md) |

<!-- SOURCE-BREADCRUMBS-END -->

<!-- DESIGN-BREADCRUMBS-START -->

## Related Design Docs

> Auto-generated cross-references to related design documentation

**Category**: [Scrobbling](INDEX.md)

### In This Section

- [Last.fm Scrobbling Integration](LASTFM_SCROBBLE.md)
- [Letterboxd Integration](LETTERBOXD.md)
- [Simkl Integration](SIMKL.md)
- [Trakt Integration](TRAKT.md)

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

## Related Documentation

- [LASTFM_SCROBBLE.md](./LASTFM_SCROBBLE.md) - Last.fm scrobbling (commercial alternative)
- [MUSICBRAINZ.md](../metadata/music/MUSICBRAINZ.md) - MusicBrainz metadata (native integration)
- [TRAKT.md](./TRAKT.md) - Trakt scrobbling (movies/TV)

---

## Notes

### Authentication (Simple Token)
- **User token**: Long-lived personal token (generated by user)
- **No OAuth**: Simpler than Last.fm (no OAuth flow)
- **Token management**: Users generate/revoke tokens in ListenBrainz settings
- **Encryption**: Encrypt user tokens in database

### Rate Limits
- **Limit**: 10 requests per second (per user token)
- **Generous**: More lenient than Last.fm (5 req/sec)
- **Throttling**: Implement token bucket rate limiter (10 req/sec)
- **Retry**: Retry with exponential backoff on 429 (rate limit exceeded)

### Listen Types
- **"playing_now"**: Currently playing track (no timestamp, updates "now playing")
- **"single"**: Single listen (includes timestamp, added to history)
- **"import"**: Batch import (up to 100 listens per request, initial sync)

### MusicBrainz ID Integration
- **ListenBrainz loves MBIDs**: Native MusicBrainz ID support (recording_mbid, artist_mbids, release_mbid)
- **Enhanced matching**: MBIDs provide accurate matching (no artist/track name ambiguity)
- **Revenge → ListenBrainz**: Map Revenge tracks to MusicBrainz IDs (music.musicbrainz_recording_id)
- **Fallback**: If no MBID, send artist_name + track_name + release_name

### Scrobbling Rules (Same as Last.fm)
- **Minimum duration**: Track must be >= 30 seconds
- **Scrobble threshold**: Track played for >= 50% OR >= 4 minutes (whichever comes first)
- **Timestamp**: Unix timestamp when track STARTED playing (NOT when scrobble sent)
- **Playing now**: Update immediately when track starts (no threshold)

### Batch Import
- **Batch size**: Up to 100 listens per request
- **Use case**: Initial sync (import Revenge listen history to ListenBrainz)
- **Listen type**: "import" (not "single")

### Pagination (Fetch History)
- **max_ts**: Maximum timestamp (fetch listens BEFORE this timestamp)
- **min_ts**: Minimum timestamp (fetch listens AFTER this timestamp)
- **count**: Number of listens per page (default 25, max 100)
- **Strategy**: Fetch all listens by iterating with max_ts (start with current timestamp, paginate backwards)

### Deduplication Strategy
- **Duplicate detection**: Match by recording_mbid + listened_at (exact timestamp)
- **Fallback**: If no MBID, match by artist_name + track_name + listened_at (within 5 minutes)
- **Conflict resolution**: Keep ListenBrainz timestamp if earlier, Revenge if later
- **Merge strategy**: Union (both ListenBrainz and Revenge entries)

### Statistics (Top Artists/Releases/Recordings)
- **ListenBrainz statistics**: Top artists, releases, recordings by listen count
- **Time ranges**: week, month, year, all_time
- **Display**: Show in Revenge UI (statistics page, user profile)
- **Refresh**: Periodic sync (River job every 6 hours)

### Delete Listens
- **User control**: Users can delete individual listens
- **Endpoint**: POST /1/delete-listen
- **Use case**: Remove incorrect scrobbles, privacy (delete embarrassing listens)

### Privacy & Open Data
- **Open-source**: ListenBrainz is open-source (user privacy, no ads)
- **Data export**: Users can export all listening history (JSON format)
- **Data import**: Users can import history from Last.fm, Spotify, etc.
- **Public/private**: Listens are public by default (users can make profile private)

### ListenBrainz vs Last.fm
- **Open-source vs Commercial**: ListenBrainz is open-source, Last.fm is commercial
- **MusicBrainz IDs**: ListenBrainz has native MBID support, Last.fm has limited support
- **Privacy**: ListenBrainz is privacy-focused (no ads, open data), Last.fm has ads (free tier)
- **User base**: Last.fm has larger user base (social features), ListenBrainz is growing
- **Recommendations**: Both have recommendations (ListenBrainz uses collaborative filtering)

### Error Handling
- **400 Bad Request**: Invalid parameters (check required fields: artist_name, track_name)
- **401 Unauthorized**: Invalid user token → prompt user to re-enter token
- **429 Too Many Requests**: Rate limit exceeded → throttle requests (10 req/sec)
- **500 Server Error**: ListenBrainz down → retry with exponential backoff

### Fallback Strategy (Music Scrobbling)
- **Order**: Last.fm (primary, largest user base) → ListenBrainz (alternative, open-source, MusicBrainz integration)
- **Both**: Support both integrations (users can enable both simultaneously)
