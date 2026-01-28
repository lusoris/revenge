# Metadata Providers

> External services for media metadata enrichment

---

## Overview

Revenge integrates with multiple metadata providers to enrich media libraries with:
- Titles, descriptions, and taglines
- Artwork (posters, backdrops, logos)
- Cast & crew information
- Ratings and reviews
- Release dates and runtime
- External identifiers

---

## Categories

### 🎬 [Video](video/INDEX.md)
Movies, TV shows, and video content metadata.

| Provider | Type | Status |
|----------|------|--------|
| [TMDB](video/TMDB.md) | Movies, TV | 🟢 Primary |
| [TVDB](video/TVDB.md) | TV Shows | 🟡 Secondary |
| [OMDB](video/OMDB.md) | Movies | 🟡 Fallback |
| [Fanart.tv](video/FANART_TV.md) | Artwork | 🟡 Supplementary |

### 🎵 [Music](music/INDEX.md)
Artists, albums, and tracks metadata.

| Provider | Type | Status |
|----------|------|--------|
| [MusicBrainz](music/MUSICBRAINZ.md) | All Music | 🟢 Primary |
| [Last.fm](music/LASTFM.md) | Tags, Similar | 🟡 Secondary |
| [Spotify](music/SPOTIFY.md) | Popularity | 🟡 Supplementary |
| [Discogs](music/DISCOGS.md) | Vinyl/Physical | 🟡 Supplementary |

### 📚 [Books](books/INDEX.md)
Books and literature metadata.

| Provider | Type | Status |
|----------|------|--------|
| [Open Library](books/OPENLIBRARY.md) | Books | 🟢 Primary |
| [Google Books](books/GOOGLE_BOOKS.md) | Books | 🟡 Secondary |
| [Goodreads](books/GOODREADS.md) | Reviews | 🟡 Supplementary |
| [ISBN DB](books/ISBNDB.md) | ISBN Lookup | 🟡 Fallback |

### 📖 [Comics](comics/INDEX.md)
Comics, manga, and graphic novels.

| Provider | Type | Status |
|----------|------|--------|
| [ComicVine](comics/COMICVINE.md) | Comics | 🟢 Primary |
| [Marvel API](comics/MARVEL.md) | Marvel | 🟡 Supplementary |
| [Grand Comics DB](comics/GCD.md) | Archive | 🟡 Supplementary |

### 🔞 [Adult](adult/INDEX.md)
Adult content metadata (isolated in `c` schema).

| Provider | Type | Status |
|----------|------|--------|
| [StashDB](adult/STASHDB.md) | Scenes | 🟢 Primary |
| [TPDB](adult/TPDB.md) | Scenes | 🟡 Secondary |
| [FreeOnes](adult/FREEONES.md) | Performers | 🟢 Primary |

---

## Provider Priority

When multiple providers have data for the same item:

1. **Primary** - First source checked, most trusted
2. **Secondary** - Checked if primary fails or for additional data
3. **Supplementary** - Merged data for enrichment
4. **Fallback** - Last resort if others fail

---

## Common Patterns

### Metadata Service Interface

```go
type MetadataProvider interface {
    Search(ctx context.Context, query string) ([]SearchResult, error)
    GetByID(ctx context.Context, id string) (*Metadata, error)
    GetImages(ctx context.Context, id string) ([]Image, error)
}
```

### Provider Configuration

```yaml
metadata:
  providers:
    tmdb:
      enabled: true
      api_key: "${TMDB_API_KEY}"
      priority: 1
    tvdb:
      enabled: true
      api_key: "${TVDB_API_KEY}"
      priority: 2
```

---

## Related Documentation

- [Wiki Providers](../wiki/INDEX.md)
- [Scrobbling Services](../scrobbling/INDEX.md)
- [Servarr Stack](../servarr/INDEX.md)
