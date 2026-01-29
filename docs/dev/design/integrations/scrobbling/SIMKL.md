# Simkl Integration

> TV tracker and movie scrobbler (alternative to Trakt)

**Service**: Simkl (https://simkl.com)
**API**: REST API with OAuth 2.0
**Category**: Scrobbling / Social
**Priority**: 🟡 MEDIUM (Alternative to Trakt)
**Status**: 🔴 DESIGN PHASE

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

### Phase 1: OAuth 2.0 Setup
- [ ] OAuth 2.0 client (authorization code flow)
- [ ] Token storage (encrypt access tokens)
- [ ] User authorization flow (redirect to Simkl → callback)

### Phase 2: Scrobbling (Real-time)
- [ ] **Check-in** (POST /checkin when playback starts)
- [ ] **Scrobble** (POST /sync/history when playback finishes)
- [ ] Progress tracking (send progress % to Simkl)
- [ ] Real-time sync (sync on playback events)

### Phase 3: Watch History Sync
- [ ] **Export to Simkl** (POST /sync/history - send Revenge watch history to Simkl)
- [ ] **Import from Simkl** (GET /sync/all-items/movies, /sync/all-items/shows)
- [ ] Bi-directional sync (merge watch histories)
- [ ] Conflict resolution (handle duplicate entries)

### Phase 4: Ratings Sync
- [ ] **Export ratings to Simkl** (POST /sync/ratings)
- [ ] **Import ratings from Simkl** (GET /sync/ratings/movies, /sync/ratings/shows)
- [ ] Bi-directional sync (merge ratings)
- [ ] Rating normalization (Simkl 1-10, Revenge 0-5)

### Phase 5: Watchlist Sync
- [ ] **Export watchlist to Simkl** (POST /sync/watchlist)
- [ ] **Import watchlist from Simkl** (GET /sync/watchlist/movies, /sync/watchlist/shows)
- [ ] Bi-directional sync (merge watchlists)

### Phase 6: Background Jobs (River)
- [ ] **Job**: `scrobble.simkl.sync_history` (periodic history sync)
- [ ] **Job**: `scrobble.simkl.sync_ratings` (periodic ratings sync)
- [ ] **Job**: `scrobble.simkl.sync_watchlist` (periodic watchlist sync)
- [ ] Rate limiting (1000 req/hour)
- [ ] Retry logic (exponential backoff)

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
