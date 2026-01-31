# Discogs Integration

> Music marketplace and database - vinyl releases, marketplace data, detailed credits


<!-- TOC-START -->

## Table of Contents

- [Status](#status)
- [Overview](#overview)
- [Developer Resources](#developer-resources)
- [API Details](#api-details)
  - [Base URL](#base-url)
  - [Authentication (Personal Access Token - Recommended)](#authentication-personal-access-token---recommended)
  - [Rate Limiting](#rate-limiting)
  - [Key Endpoints](#key-endpoints)
    - [Search Releases](#search-releases)
    - [Get Release](#get-release)
    - [Get Artist](#get-artist)
    - [Get Master Release](#get-master-release)
- [Implementation Checklist](#implementation-checklist)
  - [API Client (`internal/infra/metadata/provider_discogs.go`)](#api-client-internalinframetadataprovider-discogsgo)
  - [Release Metadata (Physical Media)](#release-metadata-physical-media)
  - [Detailed Credits](#detailed-credits)
  - [Marketplace Data (Future Feature)](#marketplace-data-future-feature)
  - [Label Information](#label-information)
  - [Error Handling](#error-handling)
- [Integration Pattern](#integration-pattern)
  - [Enrich Album with Discogs Metadata](#enrich-album-with-discogs-metadata)
- [Sources & Cross-References](#sources-cross-references)
  - [Cross-Reference Indexes](#cross-reference-indexes)
  - [Referenced Sources](#referenced-sources)
- [Related Design Docs](#related-design-docs)
  - [In This Section](#in-this-section)
  - [Related Topics](#related-topics)
  - [Indexes](#indexes)
- [Related Documentation](#related-documentation)
- [Notes](#notes)

<!-- TOC-END -->

**Service**: Discogs
**Type**: Metadata Provider (Music)
**API Version**: 2.0
**Website**: https://www.discogs.com
**API Docs**: https://www.discogs.com/developers

## Status

| Dimension | Status | Notes |
|-----------|--------|-------|
| Design | ✅ | Comprehensive REST API endpoints, marketplace data, credits |
| Sources | ✅ | API docs, authentication, database linked |
| Instructions | ✅ | Detailed implementation checklist |
| Code | 🔴 |  |
| Linting | 🔴 |  |
| Unit Testing | 🔴 |  |
| Integration Testing | 🔴 |  |
---

## Overview

**Discogs** is a music database and marketplace with comprehensive release information, especially for vinyl, rare editions, and detailed credits.

**Why Discogs**:
- Detailed release information (vinyl, cassette, CD, digital)
- Marketplace data (prices, listings, sales history)
- Comprehensive credits (musicians, producers, engineers)
- Label information
- Release notes and versions
- Free API

**Use Cases**:
- **Vinyl/physical media metadata**: Release formats, pressing details, catalog numbers
- **Marketplace integration**: Price estimates, availability, buy links (future feature)
- **Detailed credits**: Musicians, producers, engineers, studios
- **Label information**: Record labels, catalog numbers
- **Release versions**: Different pressings, reissues, remasters

**NOT Primary Source**:
- MusicBrainz is primary for digital music
- Discogs is fallback/enrichment for physical releases

---

## Developer Resources

**API Documentation**: https://www.discogs.com/developers
**Authentication**: https://www.discogs.com/developers#page:authentication
**Database**: https://www.discogs.com/developers#page:database

**Authentication**: OAuth 1.0a OR Personal Access Token
**Rate Limit**: 60 requests per minute (authenticated), 25 requests per minute (unauthenticated)
**Free Tier**: Unlimited (registration required)

---

## API Details

### Base URL
```
https://api.discogs.com/
```

### Authentication (Personal Access Token - Recommended)
```bash
# 1. Generate token: https://www.discogs.com/settings/developers
# 2. Use in requests:
Authorization: Discogs token={PERSONAL_ACCESS_TOKEN}
```

### Rate Limiting
- **Authenticated**: 60 requests per minute
- **Unauthenticated**: 25 requests per minute
- **Retry-After header**: Use for 429 errors

### Key Endpoints

#### Search Releases
```bash
GET /database/search?q=OK+Computer&type=release&artist=Radiohead&format=vinyl
Authorization: Discogs token={TOKEN}
```

**Response**:
```json
{
  "results": [
    {
      "id": 67853,
      "type": "release",
      "title": "Radiohead - OK Computer",
      "year": "1997",
      "format": ["Vinyl", "LP", "Album"],
      "label": ["Parlophone"],
      "catno": "7243 8 55229 1 5",
      "thumb": "https://i.discogs.com/...",
      "cover_image": "https://i.discogs.com/..."
    }
  ]
}
```

#### Get Release
```bash
GET /releases/{release_id}
```

**Response**:
```json
{
  "id": 67853,
  "title": "OK Computer",
  "artists": [{"name": "Radiohead", "id": 3840}],
  "year": 1997,
  "formats": [
    {
      "name": "Vinyl",
      "qty": "2",
      "descriptions": ["LP", "Album", "Gatefold"]
    }
  ],
  "labels": [
    {"name": "Parlophone", "catno": "7243 8 55229 1 5", "id": 3840}
  ],
  "tracklist": [
    {
      "position": "A1",
      "title": "Airbag",
      "duration": "4:44",
      "artists": [{"name": "Radiohead"}]
    }
  ],
  "extraartists": [
    {"name": "Nigel Godrich", "role": "Producer"},
    {"name": "Stanley Donwood", "role": "Artwork"}
  ],
  "community": {
    "rating": {"average": 4.67, "count": 1234}
  },
  "lowest_price": 25.00,
  "num_for_sale": 15
}
```

#### Get Artist
```bash
GET /artists/{artist_id}
```

#### Get Master Release
```bash
GET /masters/{master_id}
```

**Master Release**: Canonical version (groups all versions/pressings of same album)

---

## Implementation Checklist

### API Client (`internal/infra/metadata/provider_discogs.go`)
- [ ] Base URL configuration
- [ ] Personal Access Token configuration
- [ ] Rate limiting (60 req/min with token bucket)
- [ ] Error handling (401: Invalid token, 404: Not found, 429: Rate limit exceeded)
- [ ] Response parsing (JSON unmarshalling)

### Release Metadata (Physical Media)
- [ ] Search releases by artist + album
- [ ] Fetch release details (formats, labels, catalog numbers)
- [ ] Extract: format (Vinyl, CD, Cassette, Digital), pressing details, catalog number
- [ ] Store in `music_albums.discogs_id`, `music_albums.format`, `music_albums.catalog_number`

### Detailed Credits
- [ ] Fetch extraartists (producers, engineers, mixing, mastering, artwork)
- [ ] Store in `music_album_credits` table (role, name)
- [ ] Display in album detail page

### Marketplace Data (Future Feature)
- [ ] Fetch marketplace stats (lowest_price, num_for_sale, sales_history)
- [ ] Store in `music_albums.marketplace_data` (JSONB)
- [ ] Display "Buy on Discogs" link (affiliate program)

### Label Information
- [ ] Fetch label details (name, catalog number)
- [ ] Store in `music_labels` table
- [ ] Link albums to labels

### Error Handling
- [ ] Handle 401 (Invalid token - check configuration)
- [ ] Handle 404 (Release not found)
- [ ] Handle 429 (Rate limit exceeded - retry after `Retry-After` seconds)
- [ ] Log errors (no sensitive data)

---

## Integration Pattern

### Enrich Album with Discogs Metadata
```go
// Background job: Enrich albums with Discogs metadata (vinyl releases)
func (s *MusicService) EnrichAlbumWithDiscogs(albumID uuid.UUID) error {
    album := s.db.GetAlbum(albumID)

    // 1. Search Discogs
    results := s.discogsClient.SearchReleases(album.ArtistName, album.Title, "vinyl")
    if len(results) == 0 {
        return errors.New("release not found on Discogs")
    }

    // 2. Get release details
    release := s.discogsClient.GetRelease(results[0].ID)

    // 3. Extract metadata
    format := release.Formats[0].Name // "Vinyl"
    catalogNumber := release.Labels[0].Catno
    credits := make([]Credit, len(release.ExtraArtists))
    for i, artist := range release.ExtraArtists {
        credits[i] = Credit{Name: artist.Name, Role: artist.Role}
    }

    // 4. Update album
    s.db.UpdateAlbum(albumID, map[string]interface{}{
        "discogs_id":      release.ID,
        "format":          format,
        "catalog_number":  catalogNumber,
        "credits":         credits,
        "marketplace_data": map[string]interface{}{
            "lowest_price":  release.LowestPrice,
            "num_for_sale":  release.NumForSale,
        },
    })

    return nil
}
```

---


<!-- SOURCE-BREADCRUMBS-START -->

## Sources & Cross-References

> Auto-generated section linking to external documentation sources

### Cross-Reference Indexes

- [All Sources Index](../../../../sources/SOURCES_INDEX.md) - Complete list of external documentation
- [Design ↔ Sources Map](../../../../sources/DESIGN_CROSSREF.md) - Which docs reference which sources

### Referenced Sources

| Source | Documentation |
|--------|---------------|
| [Discogs API](https://www.discogs.com/developers) | [Local](../../../../sources/apis/discogs.md) |
| [Last.fm API](https://www.last.fm/api/intro) | [Local](../../../../sources/apis/lastfm.md) |

<!-- SOURCE-BREADCRUMBS-END -->

<!-- DESIGN-BREADCRUMBS-START -->

## Related Design Docs

> Auto-generated cross-references to related design documentation

**Category**: [Music](INDEX.md)

### In This Section

- [Last.fm Integration](LASTFM.md)
- [MusicBrainz Integration](MUSICBRAINZ.md)
- [Spotify Integration](SPOTIFY.md)

### Related Topics

- [Revenge - Architecture v2](../../../architecture/01_ARCHITECTURE.md) _Architecture_
- [Revenge - Design Principles](../../../architecture/02_DESIGN_PRINCIPLES.md) _Architecture_
- [Revenge - Metadata System](../../../architecture/03_METADATA_SYSTEM.md) _Architecture_
- [Revenge - Player Architecture](../../../architecture/04_PLAYER_ARCHITECTURE.md) _Architecture_
- [Plugin Architecture Decision](../../../architecture/05_PLUGIN_ARCHITECTURE_DECISION.md) _Architecture_

### Indexes

- [Design Index](../../../DESIGN_INDEX.md) - All design docs by category/topic
- [Source of Truth](../../../00_SOURCE_OF_TRUTH.md) - Package versions and status

<!-- DESIGN-BREADCRUMBS-END -->

## Related Documentation

- **MusicBrainz Integration**: [MUSICBRAINZ.md](MUSICBRAINZ.md) (primary metadata)
- **Lidarr Integration**: [../servarr/LIDARR.md](../servarr/LIDARR.md)

---

## Notes

- **Personal Access Token recommended**: Easier than OAuth 1.0a (no callback URL)
- **Rate limit**: 60 req/min (authenticated), 25 req/min (unauthenticated)
- **Master vs Release**: Master = canonical version, Release = specific pressing/version
- **Format types**: Vinyl, CD, Cassette, Digital, Minidisc, 8-Track, etc.
- **Vinyl details**: LP (12"), EP (7"), Single (7"), Gatefold, Picture Disc, Colored Vinyl
- **Catalog numbers**: Unique identifier from label (e.g., "7243 8 55229 1 5")
- **Extraartists**: Producers, engineers, mixing, mastering, artwork, photography
- **Marketplace data**: Lowest price, number for sale, sales history (requires OAuth for full access)
- **Community ratings**: Average rating, number of ratings (user-generated)
- **Search syntax**: `q` (query), `type` (release, master, artist, label), `artist`, `format`, `year`
- **API stable**: v2.0 stable, no breaking changes expected
- **Free tier**: Unlimited requests (respect rate limits)
- **Use case**: Fallback/enrichment for physical media, MusicBrainz is primary for digital
- **Cover images**: Available (use as fallback if Cover Art Archive + Spotify fail)
- **Affiliate program**: Discogs marketplace (future feature - "Buy on Discogs" links)
- **Token generation**: https://www.discogs.com/settings/developers
- **Error codes**: 401 (Invalid token), 404 (Not found), 429 (Rate limit exceeded)
- **Retry-After header**: Use for 429 errors (seconds to wait)
