# Video Metadata Providers

> Movies and TV shows metadata

---

## Overview

Video metadata providers supply information for movies and TV shows:
- Titles and translations
- Plot summaries
- Cast & crew
- Artwork (posters, backdrops, logos)
- Ratings and reviews
- Release information

---

## Providers

| Provider | Type | API | Status |
|----------|------|-----|--------|
| [TMDB](TMDB.md) | Movies, TV | REST | 🟢 Primary |
| [TVDB](TVDB.md) | TV Shows | REST v4 | 🟡 Secondary |
| [OMDB](OMDB.md) | Movies | REST | 🟡 Fallback |
| [Fanart.tv](FANART_TV.md) | Artwork | REST | 🟡 Supplementary |

---

## Provider Details

### TMDB (The Movie Database)
**Primary provider for all video content**

- ✅ Movies - Comprehensive coverage
- ✅ TV Shows - Good coverage
- ✅ People - Cast & crew profiles
- ✅ Images - High quality artwork
- ✅ Free API with generous limits

### TVDB
**Secondary provider for TV-specific data**

- ✅ TV Shows - Episode-level detail
- ✅ Airdate tracking
- ✅ Series status
- ⚠️ Paid API subscription required

### OMDB
**Fallback for IMDb ratings**

- ✅ IMDb ratings
- ✅ Rotten Tomatoes scores
- ✅ Basic movie info
- ⚠️ Limited free tier

### Fanart.tv
**Supplementary artwork source**

- ✅ HD Clearlogos
- ✅ Character art
- ✅ Season/disc art
- ✅ Unique artwork not on TMDB

---

## Priority Chain

```
Movie: TMDB → OMDB → Fanart.tv
TV Show: TMDB → TVDB → Fanart.tv
```

---

## Configuration

```yaml
metadata:
  video:
    movie:
      primary: tmdb
      fallback: [omdb]
      artwork: [tmdb, fanart]
    tvshow:
      primary: tmdb
      fallback: [tvdb]
      artwork: [tmdb, tvdb, fanart]
```

---

## Related Documentation

- [Metadata Overview](../INDEX.md)
- [Anime Providers](../../anime/INDEX.md)
- [Servarr Integration](../../servarr/INDEX.md)
