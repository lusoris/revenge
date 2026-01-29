# Wiki Providers

> Supplementary information from wiki sources

---

## Overview

Wiki providers supply supplementary information for:
- Detailed biographies
- Plot summaries
- Trivia and production notes
- External links
- Character information

---

## Providers

### Standard Wiki

| Provider | Type | Status |
|----------|------|--------|
| [Wikipedia](WIKIPEDIA.md) | General | 🟢 Active |
| Wikidata | Structured | 🟡 Planned |
| [Fandom](FANDOM.md) | Fan Wikis | 🟡 Planned |

### Adult Wiki (Isolated)

| Provider | Type | Status |
|----------|------|--------|
| [IAFD](adult/IAFD.md) | Performer Wiki | 🟡 Planned |
| AFDB | Adult Film DB | 🟡 Planned |
| [Boobpedia](adult/BOOBPEDIA.md) | Performer Wiki | 🟡 Planned |

---

## Provider Details

### Wikipedia
**General encyclopedia**

- ✅ Movie/TV plot summaries
- ✅ Actor biographies
- ✅ Production information
- ✅ Free, no API key
- ✅ Multi-language support

### Wikidata
**Structured data**

- ✅ Cross-references (IMDb, TMDB, etc.)
- ✅ Relationships
- ✅ Identifiers
- ✅ Free, no API key

### Fandom (Wikia)
**Fan-maintained wikis**

- ✅ Franchise-specific details
- ✅ Character information
- ✅ Episode guides
- ⚠️ Quality varies by wiki

---

## Use Cases

Wiki data enriches primary metadata:

| Content | Wiki Usage |
|---------|------------|
| Movies | Production trivia, box office |
| TV Shows | Episode summaries, cast changes |
| Music | Band history, album reception |
| People | Extended biography |

---

## Configuration

```yaml
metadata:
  wiki:
    wikipedia:
      enabled: true
      languages: ["en", "de", "fr"]
    wikidata:
      enabled: true
    fandom:
      enabled: false  # Quality varies
```

---

## Related Documentation

- [Metadata Providers](../metadata/INDEX.md)
- [Adult Wiki](adult/INDEX.md)
