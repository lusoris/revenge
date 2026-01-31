# Sonarr Integration

<!-- SOURCES: fx, pgx, postgresql-arrays, postgresql-json, river, servarr-wiki, sonarr-docs, typesense, typesense-go -->

<!-- DESIGN: integrations/servarr, 01_ARCHITECTURE, 02_DESIGN_PRINCIPLES, 03_METADATA_SYSTEM -->


> TV show management automation


<!-- TOC-START -->

## Table of Contents

- [Status](#status)
- [Overview](#overview)
- [Developer Resources](#developer-resources)
- [API Details](#api-details)
  - [Key Endpoints](#key-endpoints)
- [Webhook Events](#webhook-events)
  - [On Import (Episode Downloaded & Imported)](#on-import-episode-downloaded-imported)
  - [On Episode Added (New Episode Tracked)](#on-episode-added-new-episode-tracked)
  - [On Episode File Delete](#on-episode-file-delete)
  - [On Series Delete](#on-series-delete)
  - [On Rename](#on-rename)
  - [On Health Issue](#on-health-issue)
- [Implementation Checklist](#implementation-checklist)
  - [Phase 1: Client Setup](#phase-1-client-setup)
  - [Phase 2: API Implementation](#phase-2-api-implementation)
  - [Phase 3: Service Integration](#phase-3-service-integration)
  - [Phase 4: Testing](#phase-4-testing)
- [Revenge Integration Pattern](#revenge-integration-pattern)
  - [Go Client Example](#go-client-example)
- [Sources & Cross-References](#sources-cross-references)
  - [Cross-Reference Indexes](#cross-reference-indexes)
  - [Referenced Sources](#referenced-sources)
- [Related Design Docs](#related-design-docs)
  - [In This Section](#in-this-section)
  - [Related Topics](#related-topics)
  - [Indexes](#indexes)
- [Related Documentation](#related-documentation)
- [Quality Profile Mapping](#quality-profile-mapping)
- [Notes](#notes)

<!-- TOC-END -->

## Status

| Dimension | Status |
|-----------|--------|
| Design | ✅ |
| Sources | ✅ |
| Instructions | 🟡 |
| Code | 🔴 |
| Linting | 🔴 |
| Unit Testing | 🔴 |
| Integration Testing | 🔴 |**Priority**: 🔴 CRITICAL (Phase 3 - TV Show Module)
**Type**: Webhook listener + API client for metadata sync

---

## Overview

Sonarr is the industry-standard TV show management automation tool. Revenge integrates with Sonarr to:
- Receive webhook notifications when TV episodes are imported
- Sync TV show & episode metadata
- Monitor Sonarr download/import status
- Map Sonarr quality profiles to Revenge quality tiers

**Integration Points**:
- **Webhook listener**: Process Sonarr events (On Import, On Episode Added, etc.)
- **API client**: Query TV shows, episodes, series metadata
- **Metadata sync**: Enrich Revenge metadata with Sonarr data
- **Quality mapping**: Sonarr quality profiles → Revenge quality tiers

---

## Developer Resources

- 📚 **API Docs**: https://sonarr.tv/docs/api/
- 🔗 **OpenAPI Spec**: https://github.com/Sonarr/Sonarr/blob/develop/src/Sonarr.Api.V3/openapi.json
- 🔗 **GitHub**: https://github.com/Sonarr/Sonarr
- 🔗 **Wiki**: https://wiki.servarr.com/sonarr

---

## API Details

**Base Path**: `/api/v3/`
**Authentication**: `X-Api-Key` header (API key from Sonarr settings)
**Rate Limits**: None (self-hosted)

### Key Endpoints

| Endpoint | Method | Purpose |
|----------|--------|---------|
| `/series` | GET | List all TV shows |
| `/series/{id}` | GET | Get specific TV show details |
| `/episode` | GET | List episodes (filterable by series) |
| `/episode/{id}` | GET | Get specific episode details |
| `/importlist` | GET | List configured import lists |
| `/metadata` | GET | Get metadata settings |
| `/qualityprofile` | GET | List quality profiles |
| `/system/status` | GET | Get Sonarr version & status |
| `/health` | GET | Check Sonarr health |

---

## Webhook Events

Sonarr can send webhooks for the following events:

### On Import (Episode Downloaded & Imported)
```json
{
  "eventType": "Download",
  "series": {
    "id": 1,
    "title": "Breaking Bad",
    "year": 2008,
    "tvdbId": 81189,
    "imdbId": "tt0903747",
    "path": "/media/TV Shows/Breaking Bad"
  },
  "episodes": [
    {
      "id": 123,
      "episodeNumber": 1,
      "seasonNumber": 1,
      "title": "Pilot",
      "airDate": "2008-01-20",
      "overview": "High school chemistry teacher...",
      "episodeFile": {
        "id": 456,
        "relativePath": "Season 01/Breaking Bad - S01E01 - Pilot.mkv",
        "quality": "Bluray-1080p",
        "size": 2147483648
      }
    }
  ]
}
```

### On Episode Added (New Episode Tracked)
Triggered when Sonarr starts monitoring a new episode.

### On Episode File Delete
Triggered when episode file is deleted from Sonarr.

### On Series Delete
Triggered when TV show is removed from Sonarr.

### On Rename
Triggered when episode files are renamed.

### On Health Issue
Triggered when Sonarr detects health issues.

---

## Implementation Checklist

### Phase 1: Client Setup
- [ ] Create client package structure
- [ ] Implement HTTP client with resty
- [ ] Add API key authentication
- [ ] Implement rate limiting

### Phase 2: API Implementation
- [ ] Implement core API methods
- [ ] Add response type definitions
- [ ] Implement error handling

### Phase 3: Service Integration
- [ ] Create service wrapper
- [ ] Add caching layer
- [ ] Implement fx module wiring

### Phase 4: Testing
- [ ] Add unit tests with mocks
- [ ] Add integration tests

---

## Revenge Integration Pattern

```
Sonarr imports episode (S01E01)
           ↓
Sends webhook to Revenge
           ↓
Revenge processes webhook
           ↓
Stores series/episode in PostgreSQL (tvshows, tvshow_episodes)
           ↓
Enriches metadata from TheTVDB (posters, fanart, ratings)
           ↓
Updates Typesense search index
           ↓
Episode available for playback
```

### Go Client Example

```go
type SonarrClient struct {
    baseURL string
    apiKey  string
    client  *http.Client
}

func (c *SonarrClient) GetSeries(ctx context.Context, seriesID int) (*Series, error) {
    url := fmt.Sprintf("%s/api/v3/series/%d", c.baseURL, seriesID)
    req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
    req.Header.Set("X-Api-Key", c.apiKey)

    resp, err := c.client.Do(req)
    if err != nil {
        return nil, fmt.Errorf("failed to get series: %w", err)
    }
    defer resp.Body.Close()

    var series Series
    json.NewDecoder(resp.Body).Decode(&series)
    return &series, nil
}

func (c *SonarrClient) GetEpisodes(ctx context.Context, seriesID int) ([]Episode, error) {
    url := fmt.Sprintf("%s/api/v3/episode?seriesId=%d", c.baseURL, seriesID)
    req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
    req.Header.Set("X-Api-Key", c.apiKey)

    resp, err := c.client.Do(req)
    if err != nil {
        return nil, fmt.Errorf("failed to get episodes: %w", err)
    }
    defer resp.Body.Close()

    var episodes []Episode
    json.NewDecoder(resp.Body).Decode(&episodes)
    return episodes, nil
}
```

---


## Related Documentation

- [Radarr Integration](RADARR.md) - Similar workflow for movies
- [Lidarr Integration](LIDARR.md) - Similar workflow for music
- [TV Show Module](../../features/video/TVSHOW_MODULE.md)
- [Arr Integration Pattern](../../patterns/ARR_INTEGRATION.md)
- [Webhook Handling](../../patterns/WEBHOOK_PATTERNS.md)

---

## Quality Profile Mapping

| Sonarr Quality | Revenge Quality | Max Bitrate | Resolution |
|----------------|-----------------|-------------|------------|
| WEB-2160p | `4K` | 80 Mbps | 3840x2160 |
| Bluray-2160p | `4K` | 80 Mbps | 3840x2160 |
| WEB-1080p | `1080p` | 20 Mbps | 1920x1080 |
| Bluray-1080p | `1080p` | 20 Mbps | 1920x1080 |
| HDTV-1080p | `1080p` | 15 Mbps | 1920x1080 |
| WEB-720p | `720p` | 8 Mbps | 1280x720 |
| Bluray-720p | `720p` | 8 Mbps | 1280x720 |
| HDTV-720p | `720p` | 6 Mbps | 1280x720 |
| SDTV | `480p` | 3 Mbps | 720x480 |
| Any | `auto` | Varies | Varies |

---

## Notes

- **TheTVDB is primary metadata source** (consistency with Sonarr)
- Sonarr API v3 is stable (widely adopted)
- Self-hosted = no rate limits (unlike cloud APIs)
- Quality profiles are customizable in Sonarr (respect user settings)
- Sonarr handles season packs automatically (individual episode tracking)
- Episode air dates: Sonarr uses TheTVDB air dates (UTC timezone)
- Missing episodes: Sonarr tracks "monitored" status (not yet aired = no file)
