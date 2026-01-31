# Trakt Integration

> Social platform for tracking movies and TV shows

**Service**: Trakt (https://trakt.tv)
**API**: REST API with OAuth 2.0
**Category**: Scrobbling / Social
**Priority**: 🟢 HIGH (Popular tracking service)

## Status

| Dimension | Status | Notes |
| --------- | ------ | ----- |
| Design | ✅ | |
| Sources | ✅ | |
| Instructions | ✅ | |
| Code | 🔴 | |
| Linting | 🔴 | |
| Unit Testing | 🔴 | |
| Integration Testing | 🔴 | |

---

## Overview

**Trakt** is a platform that tracks what you watch across streaming services, generates personalized recommendations, and connects you with friends to see what they're watching.

**Key Features**:
- **Watch history tracking**: Automatic scrobbling of watched content
- **Ratings & reviews**: User ratings and reviews
- **Watchlists**: Track what to watch next
- **Social features**: Friends, comments, recommendations
- **Statistics**: Watch time, genres, trends
- **Cross-platform sync**: Sync across devices/apps

**Use Cases**:
- Automatic scrobbling (mark as watched)
- Sync watch history to Trakt
- Import Trakt watch history to Revenge
- Sync ratings/reviews
- Social sharing (what you're watching)

---

## Developer Resources

### API Documentation
- **Base URL**: https://api.trakt.tv
- **Documentation**: https://trakt.docs.apiary.io/
- **API Version**: v2
- **Authentication**: OAuth 2.0
- **Rate Limits**: 1000 requests per 5 minutes (per user)

### OAuth 2.0 Flow
```
1. Authorization URL: https://trakt.tv/oauth/authorize
   - client_id: Your app client ID
   - redirect_uri: https://revenge.example.com/api/v1/scrobble/trakt/callback
   - response_type: code

2. Exchange code for access token:
   POST https://api.trakt.tv/oauth/token
   {
     "code": "authorization_code",
     "client_id": "your_client_id",
     "client_secret": "your_client_secret",
     "redirect_uri": "https://revenge.example.com/api/v1/scrobble/trakt/callback",
     "grant_type": "authorization_code"
   }

3. Refresh token:
   POST https://api.trakt.tv/oauth/token
   {
     "refresh_token": "your_refresh_token",
     "client_id": "your_client_id",
     "client_secret": "your_client_secret",
     "redirect_uri": "https://revenge.example.com/api/v1/scrobble/trakt/callback",
     "grant_type": "refresh_token"
   }
```

### Required Headers
```
Content-Type: application/json
trakt-api-version: 2
trakt-api-key: {client_id}
Authorization: Bearer {access_token}
```

---

## API Endpoints

### Scrobble (Mark as Watching)
```
POST /scrobble/start
{
  "movie": {
    "title": "Inception",
    "year": 2010,
    "ids": {
      "tmdb": 27205,
      "imdb": "tt1375666"
    }
  },
  "progress": 10.0
}
```

### Scrobble (Mark as Watched)
```
POST /scrobble/stop
{
  "movie": {
    "title": "Inception",
    "year": 2010,
    "ids": {
      "tmdb": 27205,
      "imdb": "tt1375666"
    }
  },
  "progress": 100.0
}
```

### Add to Watch History
```
POST /sync/history
{
  "movies": [
    {
      "watched_at": "2023-01-15T12:00:00Z",
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

### Get Watch History
```
GET /sync/history/movies
GET /sync/history/shows

Response:
[
  {
    "id": 12345,
    "watched_at": "2023-01-15T12:00:00Z",
    "action": "watch",
    "type": "movie",
    "movie": {
      "title": "Inception",
      "year": 2010,
      "ids": {
        "trakt": 12345,
        "tmdb": 27205,
        "imdb": "tt1375666"
      }
    }
  }
]
```

### Add Rating
```
POST /sync/ratings
{
  "movies": [
    {
      "rating": 8,
      "rated_at": "2023-01-15T12:00:00Z",
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
[
  {
    "rated_at": "2023-01-15T12:00:00Z",
    "rating": 8,
    "type": "movie",
    "movie": {
      "title": "Inception",
      "year": 2010,
      "ids": {
        "trakt": 12345,
        "tmdb": 27205,
        "imdb": "tt1375666"
      }
    }
  }
]
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

---

## Implementation Checklist

### Phase 1: OAuth 2.0 Setup
- [ ] OAuth 2.0 client (authorization code flow)
- [ ] Token storage (encrypt access tokens)
- [ ] Token refresh (automatic refresh before expiry)
- [ ] User authorization flow (redirect to Trakt → callback)

### Phase 2: Scrobbling (Real-time)
- [ ] **Scrobble start** (POST /scrobble/start when playback starts)
- [ ] **Scrobble pause** (POST /scrobble/pause when playback pauses)
- [ ] **Scrobble stop** (POST /scrobble/stop when playback finishes)
- [ ] Progress tracking (send progress % to Trakt)
- [ ] Real-time sync (sync on playback events)

### Phase 3: Watch History Sync
- [ ] **Export to Trakt** (POST /sync/history - send Revenge watch history to Trakt)
- [ ] **Import from Trakt** (GET /sync/history - fetch Trakt watch history)
- [ ] Bi-directional sync (merge watch histories)
- [ ] Conflict resolution (handle duplicate entries)

### Phase 4: Ratings Sync
- [ ] **Export ratings to Trakt** (POST /sync/ratings)
- [ ] **Import ratings from Trakt** (GET /sync/ratings)
- [ ] Bi-directional sync (merge ratings)
- [ ] Rating normalization (Trakt 1-10, Revenge 0-5)

### Phase 5: Watchlist Sync
- [ ] **Export watchlist to Trakt** (POST /sync/watchlist)
- [ ] **Import watchlist from Trakt** (GET /sync/watchlist)
- [ ] Bi-directional sync (merge watchlists)

### Phase 6: Background Jobs (River)
- [ ] **Job**: `scrobble.trakt.sync_history` (periodic history sync)
- [ ] **Job**: `scrobble.trakt.sync_ratings` (periodic ratings sync)
- [ ] **Job**: `scrobble.trakt.sync_watchlist` (periodic watchlist sync)
- [ ] Rate limiting (1000 req / 5 min)
- [ ] Retry logic (exponential backoff)

---

## Integration Pattern

### Real-time Scrobbling Flow
```
User starts playback (movie/TV episode)
        ↓
Playback session starts (internal/service/playback/session.go)
        ↓
Check if user has Trakt enabled (user.integrations.trakt.enabled)
        ↓
        YES
        ↓
Send scrobble start to Trakt:
  POST /scrobble/start
  {
    "movie": {"title": "Inception", "year": 2010, "ids": {"tmdb": 27205}},
    "progress": 10.0
  }
        ↓
Playback continues → periodically send progress updates
        ↓
Playback finishes (progress >= 80%) → Send scrobble stop:
  POST /scrobble/stop
  {
    "movie": {"title": "Inception", "year": 2010, "ids": {"tmdb": 27205}},
    "progress": 100.0
  }
        ↓
Trakt marks as watched (synced to Trakt watch history)
```

### Watch History Sync Flow
```
User enables Trakt integration (Settings → Integrations → Trakt → Connect)
        ↓
OAuth 2.0 flow (redirect to Trakt → authorize → callback)
        ↓
Store access token (encrypted) in users.integrations.trakt.access_token
        ↓
Initial sync:
  1. Fetch Trakt watch history (GET /sync/history/movies, /sync/history/shows)
  2. Import to Revenge (create watch_history entries)
  3. Fetch Revenge watch history
  4. Export to Trakt (POST /sync/history)
  5. Merge & deduplicate
        ↓
Ongoing sync:
  - Real-time scrobbling (POST /scrobble/start, /scrobble/stop)
  - Periodic sync (River job every 1 hour) → fetch Trakt updates → merge
```

### Rating Sync Flow
```
User rates movie/TV show in Revenge
        ↓
Store rating in Revenge database (ratings table)
        ↓
Check if user has Trakt enabled
        ↓
        YES
        ↓
Export rating to Trakt:
  POST /sync/ratings
  {
    "movies": [
      {
        "rating": 8, // Convert Revenge 0-5 to Trakt 1-10 (multiply by 2)
        "rated_at": "2023-01-15T12:00:00Z",
        "ids": {"tmdb": 27205}
      }
    ]
  }
        ↓
Trakt rating synced
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
| [Trakt API](https://trakt.tv/b/api-docs) | [Local](../../../sources/apis/trakt.md) |

<!-- SOURCE-BREADCRUMBS-END -->

<!-- DESIGN-BREADCRUMBS-START -->

## Related Design Docs

> Auto-generated cross-references to related design documentation

**Category**: [Scrobbling](INDEX.md)

### In This Section

- [Last.fm Scrobbling Integration](LASTFM_SCROBBLE.md)
- [Letterboxd Integration](LETTERBOXD.md)
- [ListenBrainz Integration](LISTENBRAINZ.md)
- [Simkl Integration](SIMKL.md)

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

- [LASTFM_SCROBBLE.md](./LASTFM_SCROBBLE.md) - Last.fm scrobbling (music)
- [LISTENBRAINZ.md](./LISTENBRAINZ.md) - ListenBrainz scrobbling (music)
- [LETTERBOXD.md](./LETTERBOXD.md) - Letterboxd integration (movies)
- [SIMKL.md](./SIMKL.md) - Simkl tracking (alternative to Trakt)

---

## Notes

### OAuth 2.0 Token Management
- **Access token**: Short-lived (3 months)
- **Refresh token**: Long-lived (use to get new access token)
- **Automatic refresh**: Refresh access token before expiry (background job)
- **Encryption**: Encrypt tokens in database (sensitive data)

### Rate Limits
- **Limit**: 1000 requests per 5 minutes (per user)
- **Headers**: Check `X-Ratelimit-Limit`, `X-Ratelimit-Remaining`, `X-Ratelimit-Reset`
- **Throttling**: Implement token bucket rate limiter
- **Retry**: Retry with exponential backoff on 429 (rate limit exceeded)

### Scrobbling Threshold
- **Mark as watched**: Progress >= 80% (Trakt standard)
- **Scrobble start**: Send when playback starts (progress >= 0%)
- **Scrobble pause**: Optional (send when user pauses)
- **Scrobble stop**: Send when playback finishes (progress >= 80%)

### Rating Normalization
- **Trakt ratings**: 1-10 scale
- **Revenge ratings**: 0-5 stars (0.5 increments)
- **Conversion**: Revenge → Trakt (multiply by 2: 4.5 stars → 9/10)
- **Conversion**: Trakt → Revenge (divide by 2: 8/10 → 4 stars)

### Watch History Deduplication
- **Duplicate detection**: Match by TMDb ID + watch date (same day)
- **Conflict resolution**: Keep Trakt timestamp if earlier, Revenge if later
- **Merge strategy**: Union (both Trakt and Revenge entries)

### ID Mapping
- **Trakt IDs**: Trakt has own IDs (trakt.movie, trakt.show, trakt.episode)
- **TMDb IDs**: Trakt supports TMDb IDs (preferred for matching)
- **IMDb IDs**: Trakt supports IMDb IDs (fallback)
- **Revenge → Trakt**: Use TMDb ID from movies.tmdb_id field

### Bi-directional Sync Strategy
- **Export to Trakt**: POST /sync/history (send Revenge watch history to Trakt)
- **Import from Trakt**: GET /sync/history (fetch Trakt watch history)
- **Merge**: Union of both watch histories (deduplicate by TMDb ID + date)
- **Frequency**: Real-time scrobbling + periodic sync (every 1 hour)

### TV Shows Scrobbling
- **Episodes**: Scrobble individual episodes (POST /scrobble/start with episode info)
- **Season completion**: Trakt auto-detects season completion
- **Show completion**: Trakt auto-detects show completion

### Trakt VIP Features
- **Free tier**: Basic scrobbling, watch history, ratings
- **VIP tier**: Advanced stats, no ads, priority support
- **API access**: Same for both free and VIP users

### Error Handling
- **401 Unauthorized**: Access token expired → refresh token
- **404 Not Found**: Content not found on Trakt → skip (log warning)
- **429 Rate Limit**: Retry with exponential backoff
- **500 Server Error**: Retry with exponential backoff

### Privacy Considerations
- **User opt-in**: Users must explicitly enable Trakt integration
- **OAuth consent**: Users authorize Revenge to access Trakt account
- **Data visibility**: Trakt watch history is public by default (users can change in Trakt settings)
- **Token storage**: Encrypt access tokens in database

### Fallback Strategy (Social Tracking)
- **Order**: Trakt (primary) → Simkl (alternative) → Letterboxd (movies only)
