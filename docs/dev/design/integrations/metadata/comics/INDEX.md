# Comics Metadata Providers

> Comics, manga, and graphic novels metadata

---

## Overview

Comics metadata providers supply information for:
- Comic series and issues
- Manga volumes and chapters
- Publisher information
- Creator credits
- Cover artwork
- Character appearances

---

## Providers

| Provider | Type | API | Status |
|----------|------|-----|--------|
| [ComicVine](COMICVINE.md) | Comics | REST | 🟢 Primary |
| [Marvel API](MARVEL_API.md) | Marvel | REST | 🟡 Supplementary |
| [Grand Comics DB](GRAND_COMICS_DATABASE.md) | Archive | REST | 🟡 Supplementary |

---

## Provider Details

### ComicVine
**Primary provider - comprehensive comics database**

- ✅ Western comics
- ✅ Manga
- ✅ Characters and teams
- ✅ Story arcs
- ✅ Publisher info
- ⚠️ API key required

### Marvel API
**Official Marvel comics data**

- ✅ Marvel comics only
- ✅ Character appearances
- ✅ Events and crossovers
- ✅ High quality artwork
- ⚠️ Marvel content only
- ⚠️ API key required

### Grand Comics Database (GCD)
**Archival comics data**

- ✅ Historical comics
- ✅ Detailed credits
- ✅ Print runs
- ✅ Variant covers
- ✅ Free access

---

## Content Types

| Type | Primary Provider | Notes |
|------|-----------------|-------|
| Western Comics | ComicVine | DC, Image, Dark Horse, etc. |
| Marvel Comics | Marvel API + ComicVine | Official + community data |
| Manga | ComicVine | Limited, AniList better for anime |
| Graphic Novels | ComicVine | Collected editions |
| Indie Comics | ComicVine + GCD | Independent publishers |

---

## Data Flow

```
Scan Library
    ↓
Parse ComicInfo.xml (if exists)
    ↓
Identify via ComicVine
    ↓
Enrich Marvel content with Marvel API
    ↓
Cross-reference with GCD
    ↓
Download cover artwork
```

---

## ComicInfo.xml Support

Comics often include ComicInfo.xml with metadata:

```xml
<ComicInfo>
  <Series>Batman</Series>
  <Number>1</Number>
  <Volume>2016</Volume>
  <Title>I Am Gotham, Part One</Title>
  <Publisher>DC Comics</Publisher>
  <Writer>Tom King</Writer>
  <Penciller>David Finch</Penciller>
</ComicInfo>
```

---

## Configuration

```yaml
metadata:
  comics:
    primary: comicvine
    supplementary:
      - marvel
      - gcd
    parse_comicinfo: true
```

---

## Related Documentation

- [Metadata Overview](../INDEX.md)
- [Books](../books/INDEX.md)
- [Anime](../../anime/INDEX.md)
