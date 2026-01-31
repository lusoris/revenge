# Simkl Integration

> TV tracker and movie scrobbler (alternative to Trakt)


<!-- TOC-START -->

## Table of Contents

- [Status](#status)
- [Overview](#overview)
- [Developer Resources](#developer-resources)
  - [API Documentation](#api-documentation)
  - [OAuth 2.0 Flow](#oauth-20-flow)
  - [Required Headers](#required-headers)
- [API Endpoints](#api-endpoints)
  - [Check-in (Mark as Watching)](#check-in-mark-as-watching)
  - [Scrobble (Mark as Watched)](#scrobble-mark-as-watched)
  - [Get Watch History](#get-watch-history)
  - [Add Rating](#add-rating)
  - [Get Ratings](#get-ratings)
  - [Add to Watchlist](#add-to-watchlist)
  - [Get Watchlist](#get-watchlist)
  - [Get User Statistics](#get-user-statistics)
- [Implementation Checklist](#implementation-checklist)
  - [Phase 1: Client Setup](#phase-1-client-setup)
  - [Phase 2: API Implementation](#phase-2-api-implementation)
  - [Phase 3: Service Integration](#phase-3-service-integration)
  - [Phase 4: Testing](#phase-4-testing)
- [Integration Pattern](#integration-pattern)
  - [Real-time Scrobbling Flow](#real-time-scrobbling-flow)
  - [Watch History Sync Flow](#watch-history-sync-flow)
  - [Rating Sync Flow](#rating-sync-flow)
- [Sources & Cross-References](#sources-cross-references)
  - [Cross-Reference Indexes](#cross-reference-indexes)
  - [Referenced Sources](#referenced-sources)
- [Related Design Docs](#related-design-docs)
  - [In This Section](#in-this-section)
  - [Related Topics](#related-topics)
  - [Indexes](#indexes)
- [Related Documentation](#related-documentation)
- [Notes](#notes)
  - [OAuth 2.0 Token Management](#oauth-20-token-management)
  - [Rate Limits](#rate-limits)
  - [Scrobbling Threshold](#scrobbling-threshold)
  - [Rating Normalization](#rating-normalization)
  - [Watch History Deduplication](#watch-history-deduplication)
  - [ID Mapping](#id-mapping)
  - [Bi-directional Sync Strategy](#bi-directional-sync-strategy)
  - [TV Shows Scrobbling](#tv-shows-scrobbling)
  - [Anime Support](#anime-support)
  - [Error Handling](#error-handling)
  - [Privacy Considerations](#privacy-considerations)
  - [Simkl vs Trakt](#simkl-vs-trakt)
  - [Fallback Strategy (Social Tracking)](#fallback-strategy-social-tracking)

<!-- TOC-END -->

**Service**: Simkl (https://simkl.com)
**API**: REST API with OAuth 2.0
**Category**: Scrobbling / Social
**Priority**: 🟡 MEDIUM (Alternative to Trakt)

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

**Simkl** is a TV tracker and movie scrobbler similar to Trakt, offering automatic scrobbling, watch history, ratings, and recommendations. It's a popular alternative to Trakt with a cleaner UI and additional features.

**Key Features**:
- **Watch history tracking**: Automatic scrobbling of movies and TV shows
- **Ratings & reviews**: User ratings and reviews
- **Watchlist**: Track what to watch next
- **Social features**: Follow users, comments, recommendations
- **Statistics**: Watch time, genres, trends
- **Cross-platform sync**: Sync across devices/apps
- **Anime support**: Strong anime tracking features (similar to AniList)

**Use Cases**:
- Automatic scrobbling (mark as watched)
- Sync watch history to Simkl
- Import Simkl watch history to Revenge
- Sync ratings/reviews
- Alternative to Trakt (user preference)

---

## Developer Resources

### API Documentation
- **Base URL**: https://api.simkl.com
- **Documentation**: https://simkl.docs.apiary.io/
- **Authentication**: OAuth 2.0
- **Rate Limits**: 1000 requests per hour (per user)

### OAuth 2.0 Flow
```
1. Authorization URL: https://simkl.com/oauth/authorize
   - client_id: Your app client ID
   - redirect_uri: https://revenge.example.com/api/v1/scrobble/simkl/callback
   - response_type: code

2. Exchange code for access token:
   POST https://api.simkl.com/oauth/token
   {
     "code": "authorization_code",
     "client_id": "your_client_id",
     "client_secret": "your_client_secret",
     "redirect_uri": "https://revenge.example.com/api/v1/scrobble/simkl/callback",
     "grant_type": "authorization_code"
   }

3. Access token (long-lived):
   - No refresh token (access token doesn't expire)
   - Store encrypted in database
```

### Required Headers
```
Content-Type: application/json
simkl-api-key: {client_id}
Authorization: Bearer {access_token}
```

---

## API Endpoints

### Check-in (Mark as Watching)
```
POST /checkin
{
  "movie": {
    "title": "Inception",
    "year": 2010,
    "ids": {
      "tmdb": 27205,
      "imdb": "tt1375666"
    }
  }
}
```

### Scrobble (Mark as Watched)
```
POST /sync/history
{
  "movies": [
    {
      "title": "Inception",
      "year": 2010,
      "ids": {
        "tmdb": 27205,
        "imdb": "tt1375666"
      },
      "watched_at": "2023-01-15T12:00:00Z"
    }
  ]
}
```

### Get Watch History
```
GET /sync/all-items/movies?date_from=2023-01-01
GET /sync/all-items/shows?date_from=2023-01-01

Response:
{
  "movies": [
    {
      "last_watched_at": "2023-01-15T12:00:00Z",
      "movie": {
        "title": "Inception",
        "year": 2010,
        "ids": {
          "simkl": 12345,
          "tmdb": 27205,
          "imdb": "tt1375666"
        }
      }
    }
  ]
}
```

### Add Rating
```
POST /sync/ratings
{
  "movies": [
    {
      "rating": 8,
      "title": "Inception",
      "year": 2010,
      "ids": {
        "tmdb": 27205,
        "imdb": "tt1375666"
      }
    }
  ]
}
```

### Get Ratings
```
GET /sync/ratings/movies
GET /sync/ratings/shows

Response:
{
  "movies": [
    {
      "rating": 8,
      "movie": {
        "title": "Inception",
        "year": 2010,
        "ids": {
          "simkl": 12345,
          "tmdb": 27205,
          "imdb": "tt1375666"
        }
      }
    }
  ]
}
```

### Add to Watchlist
```
POST /sync/watchlist
{
  "movies": [
    {
      "title": "Inception",
      "year": 2010,
      "ids": {
        "tmdb": 27205,
        "imdb": "tt1375666"
      }
    }
  ]
}
```

### Get Watchlist
```
GET /sync/watchlist/movies
GET /sync/watchlist/shows
```

### Get User Statistics
```
GET /users/settings

Response:
{
  "user": {
    "name": "username",
    "stats": {
      "movies_watched": 234,
      "shows_watched": 45,
      "episodes_watched": 1234
    }
  }
}
```

---

## Implementation Checklist

### Phase 1: Client Setup
- [ ] Create client package structure
- [ ] Implement HTTP client
- [ ] Add OAuth 2.0 authentication (authorization code flow)
- [ ] Implement rate limiting (1000 req/hour)

### Phase 2: API Implementation
- [ ] Implement scrobble submission (check-in + history)
- [ ] Add history sync (import/export with deduplication)
- [ ] Implement error handling (401, 404, 429, 500 responses)

### Phase 3: Service Integration
- [ ] Create Simkl service wrapper
- [ ] Add user preference storage (enable/disable)
- [ ] Implement playback event hooks (check-in on start, scrobble on finish)

### Phase 4: Testing
- [ ] Add unit tests (OAuth flow, rating normalization)
- [ ] Add integration tests (full scrobbling and sync flow)

---

## Integration Pattern

### Real-time Scrobbling Flow
```
User starts playback (movie/TV episode)
        ↓
Playback session starts (internal/service/playback/session.go)
        ↓
Check if user has Simkl enabled (user.integrations.simkl.enabled)
        ↓
        YES
        ↓
Send check-in to Simkl:
  POST /checkin
  {
    "movie": {"title": "Inception", "year": 2010, "ids": {"tmdb": 27205}}
  }
        ↓
Playback continues → track progress
        ↓
Playback finishes (progress >= 80%) → Send scrobble:
  POST /sync/history
  {
    "movies": [{
      "title": "Inception",
      "year": 2010,
      "ids": {"tmdb": 27205},
      "watched_at": "2023-01-15T12:00:00Z"
    }]
  }
        ↓
Simkl marks as watched (synced to Simkl watch history)
```

### Watch History Sync Flow
```
User enables Simkl integration (Settings → Integrations → Simkl → Connect)
        ↓
OAuth 2.0 flow (redirect to Simkl → authorize → callback)
        ↓
Store access token (encrypted) in users.integrations.simkl.access_token
        ↓
Initial sync:
  1. Fetch Simkl watch history (GET /sync/all-items/movies, /sync/all-items/shows)
  2. Import to Revenge (create watch_history entries)
  3. Fetch Revenge watch history
  4. Export to Simkl (POST /sync/history)
  5. Merge & deduplicate
        ↓
Ongoing sync:
  - Real-time scrobbling (POST /checkin, POST /sync/history)
  - Periodic sync (River job every 1 hour) → fetch Simkl updates → merge
```

### Rating Sync Flow
```
User rates movie/TV show in Revenge
        ↓
Store rating in Revenge database (ratings table)
        ↓
Check if user has Simkl enabled
        ↓
        YES
        ↓
Export rating to Simkl:
  POST /sync/ratings
  {
    "movies": [{
      "rating": 8, // Convert Revenge 0-5 to Simkl 1-10 (multiply by 2)
      "ids": {"tmdb": 27205}
    }]
  }
        ↓
Simkl rating synced
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
| [Simkl API](https://simkl.docs.apiary.io/) | [Local](../../../sources/apis/simkl.md) |

<!-- SOURCE-BREADCRUMBS-END -->

<!-- DESIGN-BREADCRUMBS-START -->

## Related Design Docs

> Auto-generated cross-references to related design documentation

**Category**: [Scrobbling](INDEX.md)

### In This Section

- [Last.fm Scrobbling Integration](LASTFM_SCROBBLE.md)
- [Letterboxd Integration](LETTERBOXD.md)
- [ListenBrainz Integration](LISTENBRAINZ.md)
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

- [TRAKT.md](./TRAKT.md) - Trakt scrobbling (primary alternative)
- [LETTERBOXD.md](./LETTERBOXD.md) - Letterboxd integration (movies only, no API)
- [ANILIST.md](../anime/ANILIST.md) - AniList anime tracking

---

## Notes

### OAuth 2.0 Token Management
- **Access token**: Long-lived (doesn't expire)
- **No refresh token**: Simpler than Trakt (access token is permanent until revoked)
- **Encryption**: Encrypt tokens in database (sensitive data)

### Rate Limits
- **Limit**: 1000 requests per hour (per user)
- **Headers**: Check `X-Ratelimit-Limit`, `X-Ratelimit-Remaining`, `X-Ratelimit-Reset`
- **Throttling**: Implement token bucket rate limiter
- **Retry**: Retry with exponential backoff on 429 (rate limit exceeded)

### Scrobbling Threshold
- **Mark as watched**: Progress >= 80% (same as Trakt)
- **Check-in**: Send when playback starts (progress >= 0%)
- **Scrobble**: Send when playback finishes (progress >= 80%)

### Rating Normalization
- **Simkl ratings**: 1-10 scale
- **Revenge ratings**: 0-5 stars (0.5 increments)
- **Conversion**: Revenge → Simkl (multiply by 2: 4.5 stars → 9/10)
- **Conversion**: Simkl → Revenge (divide by 2: 8/10 → 4 stars)

### Watch History Deduplication
- **Duplicate detection**: Match by TMDb ID + watch date (same day)
- **Conflict resolution**: Keep Simkl timestamp if earlier, Revenge if later
- **Merge strategy**: Union (both Simkl and Revenge entries)

### ID Mapping
- **Simkl IDs**: Simkl has own IDs (simkl_id)
- **TMDb IDs**: Simkl supports TMDb IDs (preferred for matching)
- **IMDb IDs**: Simkl supports IMDb IDs (fallback)
- **Revenge → Simkl**: Use TMDb ID from movies.tmdb_id field

### Bi-directional Sync Strategy
- **Export to Simkl**: POST /sync/history (send Revenge watch history to Simkl)
- **Import from Simkl**: GET /sync/all-items/movies, /sync/all-items/shows
- **Merge**: Union of both watch histories (deduplicate by TMDb ID + date)
- **Frequency**: Real-time scrobbling + periodic sync (every 1 hour)

### TV Shows Scrobbling
- **Episodes**: Scrobble individual episodes (POST /sync/history with episode info)
- **Season completion**: Simkl auto-detects season completion
- **Show completion**: Simkl auto-detects show completion

### Anime Support
- **Simkl anime**: Strong anime tracking features (similar to AniList)
- **Anime IDs**: Supports AniList, MyAnimeList, Kitsu IDs
- **Use case**: Alternative to AniList/MyAnimeList for anime tracking

### Error Handling
- **401 Unauthorized**: Access token invalid/revoked → re-authenticate
- **404 Not Found**: Content not found on Simkl → skip (log warning)
- **429 Rate Limit**: Retry with exponential backoff
- **500 Server Error**: Retry with exponential backoff

### Privacy Considerations
- **User opt-in**: Users must explicitly enable Simkl integration
- **OAuth consent**: Users authorize Revenge to access Simkl account
- **Data visibility**: Simkl watch history visibility controlled by user settings
- **Token storage**: Encrypt access tokens in database

### Simkl vs Trakt
- **Simkl**: Cleaner UI, anime support, long-lived tokens (no refresh)
- **Trakt**: Larger user base, more established, short-lived tokens (refresh required)
- **API similarity**: Both use OAuth 2.0, similar endpoints (easy to support both)
- **User preference**: Support both integrations (user choice)

### Fallback Strategy (Social Tracking)
- **Order**: Trakt (primary, most popular) → Simkl (alternative, similar features)
- **Both**: Support both integrations (users can enable both simultaneously)
