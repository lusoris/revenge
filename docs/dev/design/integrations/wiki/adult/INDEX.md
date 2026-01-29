# Adult Wiki Providers

> Performer information from wiki sources (isolated in `c` schema)

**⚠️ Adult Content**: All data stored in isolated PostgreSQL schema `c`.
API namespace: `/api/v1/c/`

---

## Overview

Adult wiki providers supply supplementary performer information:
- Extended biographies
- Career timelines
- Filmographies
- Awards and nominations

---

## Providers

| Provider | Type | Status |
|----------|------|--------|
| [IAFD](IAFD.md) | Performer Wiki | 🟡 Planned |
| [Babepedia](BABEPEDIA.md) | Performer Wiki | 🟡 Planned |
| [Boobpedia](BOOBPEDIA.md) | Performer Wiki | 🟡 Planned |
| AFDB | Adult Film DB | 🟡 Planned |

---

## Provider Details

### IAFD (Internet Adult Film Database)
**Comprehensive performer database**

- ✅ Performer filmographies
- ✅ Career dates
- ✅ Awards
- ⚠️ No official API (scraping required)

### AFDB (Adult Film Database)
**Film-focused database**

- ✅ Scene listings
- ✅ Studio information
- ✅ Release dates
- ⚠️ Limited API

### Boobpedia
**Wiki-style performer info**

- ✅ Detailed biographies
- ✅ Physical attributes
- ✅ Career history
- ⚠️ MediaWiki API

---

## Data Isolation

All adult wiki data is isolated:

```sql
-- Stored in 'c' schema
c.performer_wiki_data
c.performer_awards
c.performer_filmography
```

---

## Configuration

```yaml
metadata:
  adult:
    wiki:
      iafd:
        enabled: false
      afdb:
        enabled: false
      boobpedia:
        enabled: false
```

---

## Related Documentation

- [Wiki Overview](../INDEX.md)
- [Adult Metadata](../../metadata/adult/INDEX.md)
- [FreeOnes](../../metadata/adult/FREEONES.md)
