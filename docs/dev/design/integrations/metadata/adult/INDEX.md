# Adult Metadata Providers

> Adult content metadata (isolated in `c` schema)

**⚠️ Adult Content**: All data stored in isolated PostgreSQL schema `c`.
API namespace: `/api/v1/c/`

---

## Overview

Adult metadata providers supply information for:
- Scene metadata
- Performer profiles
- Studio information
- Tags and categories

---

## Providers

### Scene Metadata

| Provider | Type | API | Status |
|----------|------|-----|--------|
| Whisparr v3 (eros) | Scenes | REST | 🟢 **Primary** (Servarr) |
| Stash App | Scenes | GraphQL | 🟡 Fallback (if connected) |
| [StashDB.org](STASHDB.md) | Scenes | GraphQL | 🟡 External Primary |
| [TPDB](THEPORNDB.md) | Scenes | REST | 🟡 External Fallback |

### Performer Metadata

| Provider | Type | API | Status |
|----------|------|-----|--------|
| Whisparr v3 (eros) | Performers | REST | 🟢 **Primary** (cached) |
| Stash App | Performers | GraphQL | 🟡 Fallback (if connected) |
| [StashDB.org](STASHDB.md) | Performers | GraphQL | 🟡 External Primary |
| [FreeOnes](../../external/adult/FREEONES.md) | Performers | REST | 🟡 Enrichment |
| [Babepedia](../../wiki/adult/BABEPEDIA.md) | Performers | Scraping | 🟡 Enrichment |
| [Boobpedia](../../wiki/adult/BOOBPEDIA.md) | Performers | MediaWiki | 🟡 Enrichment |
| [IAFD](../../wiki/adult/IAFD.md) | Performers | Scraping | 🟡 Enrichment |
| [TheNude](../../external/adult/THENUDE.md) | Performers | REST | 🟡 Enrichment |

---

## Provider Details

### Whisparr v3 (eros branch)
**Primary source - Servarr with cached metadata**

- ✅ Scene and performer metadata (cached from StashDB/TPDB)
- ✅ Studio information
- ✅ Automatic monitoring and downloads
- ✅ REST API (Servarr standard)
- ✅ Already curated and deduplicated

> **Servarr-First Principle**: Whisparr caches metadata from StashDB/TPDB. Use Whisparr as primary source to avoid redundant API calls. If user has Stash app connected, use as fallback/enrichment. External sources (StashDB.org first, then others) are for enrichment via background River jobs.

### Stash App
**Local fallback - user's personal library**

- ✅ Scene metadata from user's Stash instance
- ✅ Performer data (user curated)
- ✅ Fingerprint matching (phash)
- ✅ GraphQL API
- ⚠️ Requires user to have Stash app running
- ⚠️ Must be explicitly connected via API

### StashDB.org
**External Primary - community maintained**

- ✅ Scene fingerprinting (phash)
- ✅ Performer profiles
- ✅ Studio metadata
- ✅ Free, community-driven
- ✅ GraphQL API
- ✅ First priority among external sources

### TPDB (The Porn Database)
**Secondary scene database**

- ✅ Scene metadata
- ✅ DVD/series info (mapped to scene releases)
- ✅ REST API
- ⚠️ API key required

### FreeOnes
**Primary performer database**

- ✅ Performer biographies
- ✅ Physical attributes
- ✅ Career info
- ✅ Social links
- ⚠️ API key required

---

## Data Isolation

All adult content is isolated:

```sql
-- Separate PostgreSQL schema
CREATE SCHEMA IF NOT EXISTS c;

-- All tables prefixed
c.scenes
c.performers
c.studios
c.scene_performers
```

API namespace separation:
```
/api/v1/c/scenes
/api/v1/c/performers
/api/v1/c/studios
```

---

## Data Flow

```
Scan Library
    ↓
Check Whisparr cache (PRIMARY)
    ↓
Fallback to Stash App (if connected)
    ↓
Generate scene fingerprint (phash)
    ↓
Match via StashDB.org (external primary)
    ↓
Fallback to TPDB
    ↓
Enrich performer data (FreeOnes, Babepedia, etc.)
    ↓
Store in 'c' schema
```

> **Note**: Stash app integration requires user configuration. StashDB.org is always the first external source when Servarr/Stash don't have the data.

---

## Configuration

```yaml
# Adult content must be explicitly enabled
modules:
  adult:
    enabled: false  # Default disabled

metadata:
  adult:
    # Priority chain: Whisparr → Stash App → External Sources
    scene:
      primary: whisparr        # Servarr (cached metadata)
      fallback: stash_app      # Local Stash instance (if connected)
      external_primary: stashdb  # First among external sources
      external_fallback: tpdb
    performer:
      primary: whisparr
      fallback: stash_app
      external_primary: stashdb
      enrichment:
        - freeones
        - babepedia
        - boobpedia
        - iafd

# Stash App connection (optional)
integrations:
  stash:
    enabled: false
    url: "http://localhost:9999"
    api_key: "${STASH_API_KEY}"
```

---

## Privacy Considerations

- All data isolated in `c` schema
- Separate API namespace `/api/v1/c/`
- Can be completely disabled
- Separate user permissions
- No cross-referencing with regular content

---

## Related Documentation

- [Metadata Overview](../INDEX.md)
- [Adult Content System](../../../features/ADULT_CONTENT_SYSTEM.md)
- [Whisparr Integration](../../servarr/WHISPARR.md)
- [Social Links](../../external/adult/INDEX.md)
