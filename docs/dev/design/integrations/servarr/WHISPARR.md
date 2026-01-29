# Whisparr v3 Integration

> Adult content management automation (eros branch)

**Status**: 🟡 PLANNED
**Priority**: 🟡 MEDIUM (Phase 7 - Adult Modules)
**Type**: Webhook listener + API client for metadata sync
**Schema Isolation**: PostgreSQL schema `c` (see [Adult Content System](../../features/ADULT_CONTENT_SYSTEM.md))
**Branch**: `eros` (Whisparr v3)

---

## Overview

Whisparr v3 (eros branch) is the adult content management automation tool. Revenge integrates with Whisparr to:
- Receive webhook notifications when adult scenes are imported
- Sync performer, studio, and scene metadata
- Monitor Whisparr download/import status
- Respect privacy isolation (schema `c`, `/c/` API namespace)

**Integration Points**:
- **Webhook listener**: Process Whisparr events (On Import, On Scene Added, etc.)
- **API client**: Query scenes, performers, studios
- **Metadata sync**: Enrich Revenge metadata with Whisparr data + StashDB
- **Privacy isolation**: All adult data stored in PostgreSQL schema `c`, API namespace `/c/`

**⚠️ Important**: Whisparr uses Radarr codebase (fork) but treats "scenes" as individual videos (NOT series/episodes). Folder structure differs from TV shows. See [WHISPARR_STASHDB_SCHEMA.md](../../features/WHISPARR_STASHDB_SCHEMA.md) for details.

---

## Developer Resources

- 📚 **API Docs**: https://whisparr.com/docs/api/ (similar to Radarr v3)
- 🔗 **GitHub**: https://github.com/Whisparr/Whisparr (branch: `eros`)
- 🔗 **Based on**: Radarr v3 codebase (API structure similar)
- 🔗 **Metadata**: StashDB integration (performer/studio data)
- ℹ️ **Note**: Use `eros` branch for v3 features

---

## API Details

**Base Path**: `/api/v3/` (assumed, Radarr-based)
**Authentication**: `X-Api-Key` header (API key from Whisparr settings)
**Rate Limits**: None (self-hosted)

### Key Endpoints (Assumed - Radarr-like)

| Endpoint | Method | Purpose |
|----------|--------|---------|
| `/movie` | GET | List all scenes (called "movies" in Whisparr) |
| `/movie/{id}` | GET | Get specific scene details |
| `/importlist` | GET | List configured import lists |
| `/metadata` | GET | Get metadata settings |
| `/qualityprofile` | GET | List quality profiles |
| `/system/status` | GET | Get Whisparr version & status |
| `/health` | GET | Check Whisparr health |

**Note**: Whisparr API is not fully documented. Assume Radarr v3 API structure with adult-specific metadata fields (performers, studios).

---

## Webhook Events (Assumed - Radarr-like)

### On Import (Scene Downloaded & Imported)
```json
{
  "eventType": "Download",
  "movie": {
    "id": 1,
    "title": "Scene Title",
    "year": 2024,
    "stashdbId": "abc123",  // StashDB scene ID
    "overview": "Scene description...",
    "images": [
      {
        "coverType": "poster",
        "url": "https://example.com/scene-poster.jpg"
      }
    ],
    "path": "/media/Adult/Scene Title (2024)",
    "performers": [
      {"name": "Performer A", "stashdbId": "perf-123"},
      {"name": "Performer B", "stashdbId": "perf-456"}
    ],
    "studio": {
      "name": "Studio Name",
      "stashdbId": "studio-789"
    },
    "tags": ["tag1", "tag2"],
    "movieFile": {
      "id": 456,
      "relativePath": "Scene Title (2024).mkv",
      "quality": "Bluray-1080p",
      "size": 3221225472
    }
  }
}
```

### On Movie Added (New Scene Tracked)
Triggered when Whisparr starts monitoring a new scene.

### On Movie Delete
Triggered when scene is removed from Whisparr.

### On Movie File Delete
Triggered when scene file is deleted from Whisparr.

### On Rename
Triggered when scene files are renamed.

### On Health Issue
Triggered when Whisparr detects health issues.

---

## Implementation Checklist

- [ ] **API Client** (`internal/content/c/movie/provider_whisparr.go`)
  - [ ] Scene listing & detail fetching
  - [ ] Performer metadata extraction
  - [ ] Studio metadata extraction
  - [ ] Tag/category handling
  - [ ] Quality profile mapping
  - [ ] Health check integration

- [ ] **Webhook Handler** (`internal/api/handlers/c/webhook_whisparr.go`)
  - [ ] Parse webhook payload (On Download event)
  - [ ] Extract scene + performer + studio metadata
  - [ ] Trigger metadata enrichment (StashDB)
  - [ ] Store in PostgreSQL schema `c` (`c.scenes`, `c.performers`, `c.studios`)
  - [ ] Update Typesense search index (schema `c`)

- [ ] **Metadata Sync**
  - [ ] Map Whisparr scenes → Revenge `c.scenes` table
  - [ ] Map performers → Revenge `c.performers` table
  - [ ] Map studios → Revenge `c.studios` table
  - [ ] Map Whisparr quality profiles → Revenge quality tiers
  - [ ] Handle multi-performer scenes

- [ ] **Privacy Controls**
  - [ ] **Schema Isolation**: All queries use `SET search_path TO c, public;`
  - [ ] **API Namespace**: All endpoints use `/api/v1/c/movies`, `/api/v1/c/scenes`, `/api/v1/c/performers`
  - [ ] **NSFW Toggle**: Frontend toggle to show/hide adult content
  - [ ] **Secure Access**: Require NSFW permission for schema `c` access

- [ ] **Quality Profile Mapping**
  - [ ] Ultra HD (4K) → `quality='4K'`, `max_bitrate=80000`
  - [ ] HD-1080p → `quality='1080p'`, `max_bitrate=20000`
  - [ ] HD-720p → `quality='720p'`, `max_bitrate=8000`
  - [ ] SD → `quality='480p'`, `max_bitrate=3000`

- [ ] **Error Handling**
  - [ ] Retry failed API calls (circuit breaker)
  - [ ] Log webhook failures (obfuscate adult content in logs)
  - [ ] Handle missing scenes (not yet released)

---

## Revenge Integration Pattern

```
Whisparr imports scene
           ↓
Sends webhook to Revenge
           ↓
Revenge processes webhook
           ↓
Stores scene/performers/studio in PostgreSQL schema `c`
           ↓
Enriches metadata from StashDB (performer bios, studio info)
           ↓
Updates Typesense search index (schema `c`)
           ↓
Scene available for playback (requires NSFW permission)
```

### Go Client Example

```go
type WhisparrClient struct {
    baseURL string
    apiKey  string
    client  *http.Client
}

func (c *WhisparrClient) GetScene(ctx context.Context, sceneID int) (*Scene, error) {
    url := fmt.Sprintf("%s/api/v3/movie/%d", c.baseURL, sceneID)
    req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
    req.Header.Set("X-Api-Key", c.apiKey)

    resp, err := c.client.Do(req)
    if err != nil {
        return nil, fmt.Errorf("failed to get scene: %w", err)
    }
    defer resp.Body.Close()

    var scene Scene
    json.NewDecoder(resp.Body).Decode(&scene)
    return &scene, nil
}
```

---

## Related Documentation

- [Adult Content System](../../features/ADULT_CONTENT_SYSTEM.md) - Schema `c` isolation
- [Whisparr/StashDB Schema](../../features/WHISPARR_STASHDB_SCHEMA.md) - Database design
- [StashDB Integration](../metadata/adult/STASHDB.md) - Performer/studio metadata
- [Radarr Integration](RADARR.md) - Similar API structure
- [Arr Integration Pattern](../../patterns/arr_integration.md)
- [Webhook Handling](../../patterns/webhook_patterns.md)

---

## Quality Profile Mapping

| Whisparr Quality | Revenge Quality | Max Bitrate | Resolution |
|------------------|-----------------|-------------|------------|
| WEB-2160p | `4K` | 80 Mbps | 3840x2160 |
| Bluray-2160p | `4K` | 80 Mbps | 3840x2160 |
| WEB-1080p | `1080p` | 20 Mbps | 1920x1080 |
| Bluray-1080p | `1080p` | 20 Mbps | 1920x1080 |
| WEB-720p | `720p` | 8 Mbps | 1280x720 |
| HDTV-720p | `720p` | 8 Mbps | 1280x720 |
| SDTV | `480p` | 3 Mbps | 720x480 |
| Any | `auto` | Varies | Varies |

---

## Notes

- **Whisparr uses Radarr codebase** but treats "scenes" as individual videos (NOT series/episodes)
- **Folder structure differs from TV shows** - Sonarr codebase NOT applicable
- **StashDB is primary metadata source** for performer/studio data
- Whisparr API v3 assumed (Radarr-based), NOT fully documented
- Self-hosted = no rate limits (unlike cloud APIs)
- Quality profiles are customizable in Whisparr (respect user settings)
- **Privacy isolation CRITICAL**: PostgreSQL schema `c`, API namespace `/c/`, NSFW toggle required
- **Logs must be obfuscated** - never log adult content titles/performers in plain text
- **NSFW permission required** - users must explicitly enable adult content access
- Scene metadata: Whisparr stores performer names, studio, tags, release date
- Multi-performer scenes: Whisparr tracks performer order (primary performer first)
