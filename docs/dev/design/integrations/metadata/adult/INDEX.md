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

| Provider | Type | API | Status |
|----------|------|-----|--------|
| [StashDB](STASHDB.md) | Scenes | GraphQL | 🟢 Primary |
| [TPDB](TPDB.md) | Scenes | REST | 🟡 Secondary |
| [FreeOnes](FREEONES.md) | Performers | REST | 🟢 Primary |

---

## Provider Details

### StashDB
**Primary scene database - community maintained**

- ✅ Scene fingerprinting (phash)
- ✅ Performer profiles
- ✅ Studio metadata
- ✅ Free, community-driven
- ✅ GraphQL API

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
Generate scene fingerprint (phash)
    ↓
Match via StashDB
    ↓
Fallback to TPDB
    ↓
Fetch performer data from FreeOnes
    ↓
Store in 'c' schema
```

---

## Configuration

```yaml
# Adult content must be explicitly enabled
modules:
  adult:
    enabled: false  # Default disabled

metadata:
  adult:
    scene:
      primary: stashdb
      fallback: tpdb
    performer:
      primary: freeones
      fallback: stashdb
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
