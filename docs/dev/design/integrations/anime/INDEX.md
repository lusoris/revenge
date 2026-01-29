# Anime Providers

> Anime-specific metadata and tracking

---

## Overview

Anime providers supply specialized metadata for:
- Anime series and films
- Airing schedules
- Watch order (including OVAs, specials)
- Japanese titles and romanization
- Anime-specific ratings

---

## Providers

| Provider | Type | API | Status |
|----------|------|-----|--------|
| [AniList](ANILIST.md) | Metadata + Tracking | GraphQL | 🟢 Primary |
| [MyAnimeList](MYANIMELIST.md) | Metadata + Tracking | REST/OAuth | 🟡 Secondary |
| [Kitsu](KITSU.md) | Metadata + Tracking | JSON:API | 🟡 Alternative |

---

## Provider Details

### AniList
**Modern anime database**

- ✅ Comprehensive anime/manga data
- ✅ Watch list tracking
- ✅ GraphQL API (flexible queries)
- ✅ User ratings sync
- ✅ Free, generous rate limits

### MyAnimeList (MAL)
**Original anime database**

- ✅ Largest user base
- ✅ Historical data
- ✅ Community reviews
- ✅ OAuth authentication
- ⚠️ API requires approval

### Kitsu
**Social anime platform**

- ✅ Good API design (JSON:API)
- ✅ Streaming links
- ✅ Social features
- ✅ Free, no approval needed

---

## Why Separate from TMDB/TVDB?

Anime has unique requirements:

| Feature | TMDB/TVDB | Anime Providers |
|---------|-----------|-----------------|
| Episode numbering | Western style | Absolute + seasonal |
| Specials/OVAs | Limited | Comprehensive |
| Airing info | Basic | Detailed schedules |
| Watch order | None | Recommended order |
| Japanese titles | Sometimes | Always |

---

## Data Flow

```
Scan Library (anime folder)
    ↓
Identify via filename/AniDB hash
    ↓
Fetch metadata from AniList
    ↓
Fallback to MAL/Kitsu
    ↓
Map episodes (absolute ↔ seasonal)
    ↓
Sync watch progress (if enabled)
```

---

## Configuration

```yaml
metadata:
  anime:
    enabled: true
    primary: anilist
    fallback: [mal, kitsu]

    # Episode mapping
    episode_mapping:
      absolute_to_seasonal: true
      use_tvdb_mapping: false

    # Tracking sync
    tracking:
      anilist:
        enabled: true
      mal:
        enabled: false
```

---

## Episode Mapping

Anime often has different numbering systems:

```
AniList: Episode 25 (S1E25)
TVDB:    S1E25 or S2E01
File:    Episode 25 (absolute)
```

Revenge maintains mapping tables for consistent display.

---

## Related Documentation

- [Video Metadata](../metadata/video/INDEX.md)
- [Scrobbling Services](../scrobbling/INDEX.md)
- [Simkl (Anime tracking)](../scrobbling/SIMKL.md)
