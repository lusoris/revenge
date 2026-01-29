# Servarr Stack Integration

> Automated media management tools

---

## Overview

The *arr stack provides automated media acquisition and organization:
- Automated searching and downloading
- Quality management
- Library organization
- Metadata enrichment

---

## Applications

| Application | Content | Status |
|-------------|---------|--------|
| [Radarr](RADARR.md) | Movies | 🟢 Supported |
| [Sonarr](SONARR.md) | TV Shows | 🟢 Supported |
| [Lidarr](LIDARR.md) | Music | 🟢 Supported |
| [Chaptarr](CHAPTARR.md) | Books & Audiobooks | 🟢 Supported |
| [Whisparr v3](WHISPARR.md) | Adult (eros branch) | 🟢 Supported |

> **Note**: Chaptarr uses the Readarr API but handles both books and audiobooks. Readarr itself is currently inactive.

---

## Application Details

### Radarr
**Movie collection management**

- ✅ Movie searching
- ✅ Quality profiles
- ✅ Automatic upgrades
- ✅ List sync

### Sonarr
**TV series management**

- ✅ Series tracking
- ✅ Episode management
- ✅ Season packs
- ✅ Anime support

### Lidarr
**Music collection management**

- ✅ Artist tracking
- ✅ Album management
- ✅ Quality profiles
- ✅ Import lists

### Chaptarr
**Book & audiobook collection management** (uses Readarr API)

- ✅ Author tracking
- ✅ Book management
- ✅ Audiobook support
- ✅ Calibre integration
- ℹ️ Replaces Readarr (which is currently inactive)

### Whisparr v3 (eros)
**Adult content management**

- ✅ Scene tracking (NOT series/episodes)
- ✅ Performer monitoring
- ✅ Studio management
- ⚠️ Isolated in `c` schema
- ⚠️ Custom UI/UX required (differs from TV structure)
- ℹ️ Uses `eros` branch for v3 features

---

## Integration Modes

### Read-Only (Recommended)
Revenge reads library paths from Servarr:

```yaml
servarr:
  radarr:
    url: "http://radarr:7878"
    api_key: "${RADARR_API_KEY}"
    mode: "read"  # Only read library paths
```

### Full Sync
Bi-directional synchronization:

```yaml
servarr:
  radarr:
    url: "http://radarr:7878"
    api_key: "${RADARR_API_KEY}"
    mode: "sync"  # Full sync
    sync:
      watched_status: true
      ratings: true
```

---

## Common API Patterns

All *arr applications share similar API structure:

```
GET /api/v3/movie           # List all
GET /api/v3/movie/{id}      # Get by ID
POST /api/v3/movie          # Add new
PUT /api/v3/movie/{id}      # Update
DELETE /api/v3/movie/{id}   # Delete
```

Authentication via `X-Api-Key` header.

---

## Library Path Integration

```go
// Fetch root folders from Radarr
func (c *RadarrClient) GetRootFolders(ctx context.Context) ([]RootFolder, error) {
    resp, err := c.get(ctx, "/api/v3/rootfolder")
    // Returns paths like /movies, /movies-4k
}

// Use paths in Revenge library configuration
type Library struct {
    Name string
    Path string  // From Radarr root folder
    Type string  // "movie"
}
```

---

## Configuration

```yaml
servarr:
  enabled: true

  radarr:
    enabled: true
    url: "http://radarr:7878"
    api_key: "${RADARR_API_KEY}"

  sonarr:
    enabled: true
    url: "http://sonarr:8989"
    api_key: "${SONARR_API_KEY}"

  lidarr:
    enabled: true
    url: "http://lidarr:8686"
    api_key: "${LIDARR_API_KEY}"

  chaptarr:
    enabled: true
    url: "http://chaptarr:8787"
    api_key: "${CHAPTARR_API_KEY}"

  whisparr:  # v3 (eros branch)
    enabled: false  # Adult content disabled by default
    url: "http://whisparr:6969"
    api_key: "${WHISPARR_API_KEY}"
```

---

## Related Documentation

- [Metadata Providers](../metadata/INDEX.md)
- [Scrobbling Services](../scrobbling/INDEX.md)
- [External Services](../external/INDEX.md)
