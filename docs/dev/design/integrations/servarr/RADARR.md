# Radarr Integration

> Movie management automation and metadata synchronization

**Status**: ✅ Researched | API v3 | Self-hosted
**Priority**: 🔴 CRITICAL (Phase 2 - Movie Module)
**Authentication**: API Key (`X-Api-Key` header)

---

## Overview

**Purpose**: Automatic movie downloading, metadata management, library organization

**Integration Points**:
- Webhook listener for import/download events
- API client for metadata synchronization
- Quality profile mapping
- Root folder management

---

## Developer Resources

- **API Documentation**: https://radarr.video/docs/api/
- **OpenAPI Spec**: https://github.com/Radarr/Radarr/blob/develop/src/Radarr.Api.V3/openapi.json
- **GitHub Repository**: https://github.com/Radarr/Radarr
- **Wiki**: https://wiki.servarr.com/radarr

---

## API Details

### Base Configuration
- **Base Path**: `/api/v3/`
- **Authentication**: API Key header `X-Api-Key: {your_api_key}`
- **Rate Limit**: None (self-hosted)
- **Response Format**: JSON

### Key Endpoints

#### Movies
```bash
GET  /api/v3/movie              # List all movies
GET  /api/v3/movie/{id}         # Get movie by ID
POST /api/v3/movie              # Add new movie
PUT  /api/v3/movie/{id}         # Update movie
DELETE /api/v3/movie/{id}       # Delete movie
```

#### Import Lists
```bash
GET  /api/v3/importlist         # List import lists
POST /api/v3/importlist         # Add import list
```

#### Metadata
```bash
GET  /api/v3/metadata           # List metadata providers
```

#### Media Management
```bash
GET  /api/v3/mediamanagement    # Get media management config
PUT  /api/v3/mediamanagement    # Update media management config
```

#### System
```bash
GET  /api/v3/system/status      # Get system status
GET  /api/v3/health             # Get health check
```

---

## Webhook Events

Radarr sends webhooks for:
- **On Import**: Movie file imported
- **On Upgrade**: Movie file upgraded
- **On Rename**: Movie file renamed
- **On Movie Added**: Movie added to library
- **On Movie Delete**: Movie deleted from library
- **On Movie File Delete**: Movie file deleted
- **On Health Issue**: Health check issue

### Webhook Payload Example
```json
{
  "eventType": "MovieAdded",
  "movie": {
    "id": 1,
    "title": "The Matrix",
    "year": 1999,
    "tmdbId": 603,
    "imdbId": "tt0133093",
    "overview": "...",
    "images": [
      {
        "coverType": "poster",
        "url": "https://..."
      }
    ],
    "folderPath": "/movies/The Matrix (1999)"
  }
}
```

---

## Implementation Checklist

- [ ] Create API client: `internal/service/metadata/provider_radarr.go`
- [ ] Implement webhook handler: `internal/api/handlers/webhook_radarr.go`
- [ ] Metadata sync (title, overview, posters, cast, crew)
- [ ] Quality profile mapping
- [ ] Root folder management
- [ ] Event processing (import, upgrade, rename, delete)
- [ ] Error handling and retries
- [ ] Rate limiting (if needed for external instances)

---

## Revenge Integration Pattern

### Metadata Flow
```
Radarr (Add Movie) → Webhook → Revenge (Process Event)
                                       ↓
                               Store in PostgreSQL (movies table)
                                       ↓
                               Enrich with TMDb metadata
                                       ↓
                               Update search index (Typesense)
```

### Client Example
```go
type RadarrClient struct {
    baseURL string
    apiKey  string
    client  *http.Client
}

func (c *RadarrClient) GetMovie(ctx context.Context, id int) (*Movie, error) {
    req, _ := http.NewRequestWithContext(ctx, "GET",
        fmt.Sprintf("%s/api/v3/movie/%d", c.baseURL, id), nil)
    req.Header.Set("X-Api-Key", c.apiKey)

    resp, err := c.client.Do(req)
    if err != nil {
        return nil, fmt.Errorf("radarr request failed: %w", err)
    }
    defer resp.Body.Close()

    var movie Movie
    if err := json.NewDecoder(resp.Body).Decode(&movie); err != nil {
        return nil, fmt.Errorf("decode failed: %w", err)
    }

    return &movie, nil
}
```

---

## Related Documentation

- [Sonarr Integration](SONARR.md) - TV show management
- [Lidarr Integration](LIDARR.md) - Music management
- [Movie Module](../../architecture/ARCHITECTURE_V2.md#movie-module) - Revenge movie module design
- [Arr Integration Pattern](../../patterns/arr_integration.md) - Common patterns for all *arr services
- [Webhook Handling](../../patterns/webhook_patterns.md) - Webhook processing patterns

---

## Quality Profile Mapping

| Radarr Profile | Revenge Quality | Max Bitrate | Resolution |
|----------------|-----------------|-------------|------------|
| Ultra HD | ultra_hd | 80 Mbps | 2160p (4K) |
| HD-1080p | hd_1080p | 15 Mbps | 1080p |
| HD-720p | hd_720p | 8 Mbps | 720p |
| SD | sd | 3 Mbps | 480p |
| Any | any | Variable | Any |

---

## Notes

- Radarr uses TMDb as primary metadata source - Revenge should too for consistency
- Radarr v3 API is stable, breaking changes unlikely
- Self-hosted means no rate limits, but external instances may implement them
- Quality profiles are customizable - mapping should be configurable
- Root folders define library organization - respect Radarr's structure
