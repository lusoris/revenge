# TheTVDB Integration

> Primary metadata provider for TV shows

**Status**: 🟡 PLANNED
**Priority**: 🔴 CRITICAL (Phase 3 - TV Show Module)
**Type**: HTTP API client with JWT authentication

---

## Overview

TheTVDB is the industry-standard metadata provider for TV shows. Revenge uses TheTVDB as the primary metadata source for:
- TV series metadata (title, overview, network, status)
- Season metadata (season numbers, posters)
- Episode metadata (episode titles, air dates, overviews)
- Artwork (series posters, fanart, season posters, episode thumbnails)

**Integration Points**:
- **API client**: Query series, seasons, episodes, artwork
- **JWT authentication**: Login to obtain JWT token, refresh periodically
- **Rate limiting**: Varies by subscription tier
- **Artwork CDN**: Download posters, fanart, banners

---

## Developer Resources

- 📚 **API Docs**: https://thetvdb.github.io/v4-api/
- 🔗 **API v4**: https://api4.thetvdb.com/v4/
- 🔗 **GitHub**: https://github.com/thetvdb/v4-api
- 🔗 **Swagger UI**: https://thetvdb.github.io/v4-api/#/

---

## API Details

**Base URL**: `https://api4.thetvdb.com/v4/`
**Authentication**: JWT Bearer token (POST `/login` first)
**Rate Limits**: Varies by subscription tier (free tier has limited requests)
**Free Tier**: Available (requires API key registration)

### Authentication Flow

1. **Login** (POST `/login`):
```json
{
  "apikey": "your-api-key",
  "pin": ""
}
```

2. **Receive JWT Token**:
```json
{
  "status": "success",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

3. **Use Token** (Header: `Authorization: Bearer {token}`)

4. **Refresh Token** (periodically, token expires after 24 hours)

### Key Endpoints

| Endpoint | Purpose |
|----------|---------|
| `/login` | Authenticate & obtain JWT token |
| `/series/{id}` | Get series details |
| `/series/{id}/extended` | Get series + episodes + seasons |
| `/seasons/{id}` | Get season details |
| `/seasons/{id}/extended` | Get season + episodes |
| `/episodes/{id}` | Get episode details |
| `/episodes/{id}/extended` | Get episode + translations |
| `/artwork/{id}` | Get artwork details |
| `/artwork/types` | List artwork types |
| `/search` | Search series by name |

---

## Implementation Checklist

- [ ] **API Client** (`internal/service/metadata/provider_thetvdb.go`)
  - [ ] JWT authentication (login, token refresh)
  - [ ] Series metadata fetching
  - [ ] Season metadata fetching
  - [ ] Episode metadata fetching
  - [ ] Artwork fetching (posters, fanart, season posters, episode thumbnails)
  - [ ] Search functionality
  - [ ] Rate limiting (tier-dependent)
  - [ ] Error handling & retries (circuit breaker)

- [ ] **JWT Token Management**
  - [ ] Login on startup
  - [ ] Store token in memory (OR cache)
  - [ ] Refresh token every 23 hours
  - [ ] Handle token expiration (401 response → re-login)

- [ ] **Metadata Enrichment**
  - [ ] Map TheTVDB series → Revenge `tvshows` table
  - [ ] Map TheTVDB seasons → Revenge `tvshow_seasons` table
  - [ ] Map TheTVDB episodes → Revenge `tvshow_episodes` table
  - [ ] Map artwork → local storage (posters, fanart)
  - [ ] Extract external IDs (IMDb, TMDb) for cross-referencing

- [ ] **Artwork Handling**
  - [ ] Download series posters
  - [ ] Download series fanart (backdrops)
  - [ ] Download season posters
  - [ ] Download episode thumbnails
  - [ ] Generate Blurhash for placeholders
  - [ ] Store images locally (configurable path)

---

## Revenge Integration Pattern

```
User requests TV show metadata (TheTVDB ID: 121361)
           ↓
Revenge queries TheTVDB API: /series/121361/extended
           ↓
TheTVDB returns series + seasons + episodes
           ↓
Revenge stores in PostgreSQL (tvshows, tvshow_seasons, tvshow_episodes)
           ↓
Query /artwork for posters, fanart
           ↓
Download artwork from TheTVDB CDN
           ↓
Generate Blurhash for placeholders
           ↓
Update Typesense search index
           ↓
TV show metadata available
```

### Go Client Example

```go
type TheTVDBClient struct {
    baseURL string
    apiKey  string
    token   string
    client  *http.Client
    mu      sync.RWMutex
}

func (c *TheTVDBClient) Login(ctx context.Context) error {
    url := fmt.Sprintf("%s/login", c.baseURL)
    body := map[string]string{"apikey": c.apiKey, "pin": ""}
    jsonBody, _ := json.Marshal(body)

    req, _ := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(jsonBody))
    req.Header.Set("Content-Type", "application/json")

    resp, err := c.client.Do(req)
    if err != nil {
        return fmt.Errorf("login failed: %w", err)
    }
    defer resp.Body.Close()

    var result struct {
        Data struct {
            Token string `json:"token"`
        } `json:"data"`
    }
    json.NewDecoder(resp.Body).Decode(&result)

    c.mu.Lock()
    c.token = result.Data.Token
    c.mu.Unlock()

    return nil
}

func (c *TheTVDBClient) GetSeries(ctx context.Context, seriesID int) (*Series, error) {
    c.mu.RLock()
    token := c.token
    c.mu.RUnlock()

    url := fmt.Sprintf("%s/series/%d/extended", c.baseURL, seriesID)
    req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
    req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

    resp, err := c.client.Do(req)
    if err != nil {
        return nil, fmt.Errorf("failed to get series: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode == 401 {
        // Token expired, re-login
        c.Login(ctx)
        return c.GetSeries(ctx, seriesID)
    }

    var result struct {
        Data Series `json:"data"`
    }
    json.NewDecoder(resp.Body).Decode(&result)
    return &result.Data, nil
}
```

---

## Related Documentation

- [TV Show Module](../../architecture/modules/TVSHOW.md)
- [TMDb Integration](TMDB.md) - Alternative TV metadata source
- [Sonarr Integration](../../servarr/SONARR.md) - TV show management
- [Metadata Enrichment Pattern](../../patterns/metadata_enrichment.md)

---

## Artwork Types

TheTVDB provides multiple artwork types:

| Type | Description |
|------|-------------|
| `poster` | Series poster (vertical) |
| `fanart` | Backdrop/fanart (horizontal) |
| `banner` | Series banner (wide) |
| `season` | Season poster |
| `seasonwide` | Season banner |
| `clearlogo` | Logo with transparent background |
| `clearart` | Artwork with transparent background |
| `background` | Background image |

Query `/artwork/types` for full list.

---

## External ID Mapping

TheTVDB provides external IDs for cross-referencing:

```json
{
  "id": 121361,
  "remoteIds": [
    {
      "id": "tt0944947",
      "type": 2,
      "sourceName": "IMDB"
    },
    {
      "id": "1399",
      "type": 4,
      "sourceName": "TheMovieDB.com"
    }
  ]
}
```

Use these IDs to link with TMDb, IMDb, etc.

---

## Notes

- **TheTVDB API v4 is current** (v3 deprecated)
- **JWT authentication required** (token expires after 24 hours, refresh periodically)
- **Rate limits vary by tier** (free tier has limited requests, paid tiers have higher limits)
- **Free tier available** (register for API key at https://thetvdb.com/api-information)
- **Artwork CDN**: TheTVDB provides URLs to artwork (download separately)
- **Blurhash**: Generate locally from downloaded images (not provided by TheTVDB)
- **External IDs**: Use IMDb/TMDb IDs to cross-reference with other services
- **Episode air dates**: TheTVDB uses UTC timezone (convert to user's timezone)
- **Episode numbering**: Absolute episode numbers available (useful for anime)
- **Translations**: TheTVDB supports multiple languages (query `/episodes/{id}/translations`)
- **Network/Studio**: TheTVDB provides network info (e.g., "HBO", "Netflix")
- **Status**: Series status (Continuing, Ended, Upcoming)
- **Token refresh strategy**: Use background job (every 23 hours) to refresh token
