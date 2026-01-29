# Music Metadata Providers

> Artists, albums, and tracks metadata

---

## Overview

Music metadata providers supply information for:
- Artist biographies and images
- Album artwork and release info
- Track listings and durations
- Genre classification
- Similar artists/albums

---

## Providers

| Provider | Type | API | Status |
|----------|------|-----|--------|
| [MusicBrainz](MUSICBRAINZ.md) | All Music | REST | 🟢 Primary |
| [Last.fm](LASTFM.md) | Tags, Similar | REST | 🟡 Secondary |
| [Spotify](SPOTIFY.md) | Popularity | REST/OAuth | 🟡 Supplementary |
| [Discogs](DISCOGS.md) | Physical | REST | 🟡 Supplementary |

---

## Provider Details

### MusicBrainz
**Primary provider - open music database**

- ✅ Comprehensive artist/album/track data
- ✅ Release groups and variants
- ✅ MBIDs for cross-referencing
- ✅ Free, no API key required
- ✅ Links to Cover Art Archive

### Last.fm
**Secondary for tags and recommendations**

- ✅ User-generated tags
- ✅ Similar artists/tracks
- ✅ Play statistics
- ✅ Artist images
- ⚠️ API key required

### Spotify
**Supplementary for popularity metrics**

- ✅ Popularity scores
- ✅ Audio features (tempo, energy)
- ✅ Genre classification
- ⚠️ OAuth required
- ⚠️ Rate limited

### Discogs
**Supplementary for physical releases**

- ✅ Vinyl pressings
- ✅ Label information
- ✅ Barcode/catalog numbers
- ✅ Marketplace pricing
- ⚠️ OAuth required

---

## Data Flow

```
Scan Library
    ↓
Identify via MusicBrainz (fingerprint/tags)
    ↓
Fetch metadata from MusicBrainz
    ↓
Enrich with Last.fm tags
    ↓
Add Spotify popularity (optional)
    ↓
Fetch artwork from Cover Art Archive
```

---

## Configuration

```yaml
metadata:
  music:
    primary: musicbrainz
    enrichment:
      - lastfm
      - spotify
    artwork:
      - coverartarchive
      - lastfm
```

---

## Related Documentation

- [Metadata Overview](../INDEX.md)
- [Last.fm Scrobbling](../../scrobbling/LASTFM_SCROBBLE.md)
- [ListenBrainz Scrobbling](../../scrobbling/LISTENBRAINZ.md)
