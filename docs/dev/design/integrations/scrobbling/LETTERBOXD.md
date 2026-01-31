# Letterboxd Integration

> Social network for movie lovers and film tracking


<!-- TOC-START -->

## Table of Contents

- [Status](#status)
- [Overview](#overview)
- [Developer Resources](#developer-resources)
  - [API Status](#api-status)
  - [Unofficial API Options](#unofficial-api-options)
  - [Authentication](#authentication)
- [Integration Approach](#integration-approach)
  - [Web Scraping Strategy](#web-scraping-strategy)
    - [User Profile Page](#user-profile-page)
    - [Diary (Watch History)](#diary-watch-history)
    - [Ratings](#ratings)
    - [Watchlist](#watchlist)
- [Implementation Checklist](#implementation-checklist)
  - [Phase 1: Client Setup](#phase-1-client-setup)
  - [Phase 2: API Implementation](#phase-2-api-implementation)
  - [Phase 3: Service Integration](#phase-3-service-integration)
  - [Phase 4: Testing](#phase-4-testing)
- [Integration Pattern](#integration-pattern)
  - [One-Time Import Flow](#one-time-import-flow)
  - [Rating Normalization](#rating-normalization)
- [Sources & Cross-References](#sources-cross-references)
  - [Cross-Reference Indexes](#cross-reference-indexes)
  - [Referenced Sources](#referenced-sources)
- [Related Design Docs](#related-design-docs)
  - [In This Section](#in-this-section)
  - [Related Topics](#related-topics)
  - [Indexes](#indexes)
- [Related Documentation](#related-documentation)
- [Notes](#notes)
  - [No Official API (Web Scraping Only)](#no-official-api-web-scraping-only)
  - [User-Agent Requirement](#user-agent-requirement)
  - [Rate Limits (Conservative)](#rate-limits-conservative)
  - [Read-Only Integration](#read-only-integration)
  - [Movie Matching Strategy](#movie-matching-strategy)
  - [Pagination](#pagination)
  - [HTML Parsing Challenges](#html-parsing-challenges)
  - [robots.txt Compliance](#robotstxt-compliance)
  - [Maintenance Burden](#maintenance-burden)
  - [Letterboxd vs Trakt](#letterboxd-vs-trakt)
  - [Privacy Considerations](#privacy-considerations)
  - [Fallback Strategy (Movie Tracking)](#fallback-strategy-movie-tracking)
  - [Future: Official API](#future-official-api)

<!-- TOC-END -->

**Service**: Letterboxd (https://letterboxd.com)
**API**: No official public API (web scraping OR unofficial API)
**Category**: Scrobbling / Social (Movies Only)
**Priority**: 🟡 LOW (No official API, niche audience)

## Status

| Dimension | Status |
|-----------|--------|
| Design | ✅ |
| Sources | ✅ |
| Instructions | 🟡 |
| Code | 🔴 |
| Linting | 🔴 |
| Unit Testing | 🔴 |
| Integration Testing | 🔴 |
---

## Overview

**Letterboxd** is a social network for film lovers, allowing users to rate, review, and track movies they've watched. It's popular among film enthusiasts for its clean UI and focus on movies (no TV shows).

**Key Features**:
- **Movie tracking**: Track what you watch (diary entries)
- **Ratings & reviews**: Rate movies (0.5-5 stars), write reviews
- **Lists**: Create custom movie lists (e.g., "Best of 2023")
- **Watchlist**: Track movies to watch
- **Social features**: Follow users, likes, comments
- **Statistics**: Watch count, top genres, decades
- **Movie metadata**: IMDb/TMDb integration

**Use Cases**:
- Import Letterboxd watch history to Revenge
- Export Revenge watch history to Letterboxd (if unofficial API available)
- Sync ratings/reviews
- Display Letterboxd profile in Revenge UI

**⚠️ NO OFFICIAL PUBLIC API**:
- Letterboxd has NO official public API (as of 2026)
- **Options**: Web scraping (HTML parsing) OR unofficial API (community-made, may break)
- **Priority**: LOW (high maintenance burden, niche audience)

---

## Developer Resources

### API Status
- **Official API**: NONE (Letterboxd has private API for iOS/Android apps only)
- **Unofficial API**: Community-made wrappers exist (may break anytime)
- **Web Scraping**: Alternative approach (parse HTML)

### Unofficial API Options
1. **letterboxd-api** (npm package): https://github.com/zaccolley/letterboxd
   - JavaScript library (can reverse-engineer for Go)
   - Scrapes Letterboxd website
   - No authentication (public data only)

2. **Custom web scraping**: Parse Letterboxd HTML directly
   - Libraries: `github.com/PuerkitoBio/goquery` (Go)
   - Rate limiting: Very conservative (1 req/sec)

### Authentication
- **No API authentication**: Letterboxd has no public API authentication
- **Web scraping**: Public data only (no user-specific actions like adding to diary)

---

## Integration Approach

### Web Scraping Strategy

**⚠️ CRITICAL: No Official API - Web Scraping Required**

#### User Profile Page
```
URL: https://letterboxd.com/{username}/

HTML Structure (example):
<section class="profile-stats">
  <h2>Profile Stats</h2>
  <p><strong>Films:</strong> 1,234</p>
  <p><strong>This year:</strong> 45</p>
  <p><strong>Lists:</strong> 12</p>
  <p><strong>Following:</strong> 67</p>
  <p><strong>Followers:</strong> 89</p>
</section>
```

#### Diary (Watch History)
```
URL: https://letterboxd.com/{username}/films/diary/

HTML Structure (example):
<table class="diary-table">
  <tbody>
    <tr class="diary-entry-row">
      <td class="td-day">
        <a href="/.../">15 Jan 2023</a>
      </td>
      <td class="td-film-details">
        <h3><a href="/film/inception-2010/">Inception</a></h3>
        <span>2010</span>
      </td>
      <td class="td-rating">
        <span class="rating">★★★★½</span>
      </td>
      <td class="td-review">
        <a href="/.../">Review text...</a>
      </td>
    </tr>
  </tbody>
</table>
```

#### Ratings
```
URL: https://letterboxd.com/{username}/films/ratings/

HTML Structure (example):
<ul class="poster-list">
  <li class="poster-container">
    <div class="film-poster">
      <img src="..." alt="Inception">
      <p class="poster-viewingdata">
        <span class="rating">★★★★½</span>
      </p>
    </div>
  </li>
</ul>
```

#### Watchlist
```
URL: https://letterboxd.com/{username}/watchlist/

HTML Structure (example):
<ul class="poster-list">
  <li class="poster-container">
    <div class="film-poster">
      <img src="..." alt="Movie Title">
      <h2>Movie Title</h2>
      <p>2023</p>
    </div>
  </li>
</ul>
```

---

## Implementation Checklist

### Phase 1: Client Setup
- [ ] Create client package structure
- [ ] Implement HTTP client with User-Agent
- [ ] Add web scraping with goquery
- [ ] Implement rate limiting (1 req/sec)

### Phase 2: API Implementation
- [ ] Implement profile scraping (stats)
- [ ] Implement watch history scraping (diary with pagination)
- [ ] Implement error handling (scraper failures)

### Phase 3: Service Integration
- [ ] Create Letterboxd scraper service wrapper
- [ ] Add username storage (user preference)
- [ ] Implement one-time import flow

### Phase 4: Testing
- [ ] Add unit tests (scraper parsing)
- [ ] Add integration tests (full import flow)

---

## Integration Pattern

### One-Time Import Flow
```
User enables Letterboxd import (Settings → Integrations → Letterboxd → Enter Username)
        ↓
Validate username (scrape profile page, check if exists)
        ↓
Store Letterboxd username in users.integrations.letterboxd.username
        ↓
Initial import (one-time):
  1. Scrape diary (https://letterboxd.com/{username}/films/diary/)
     - Parse watch dates, movie titles, years, ratings
     - Paginate (fetch all pages)
  2. Match movies to Revenge (by title + year → lookup TMDb/IMDb ID)
  3. Import to Revenge:
     - Create watch_history entries (movies.watch_history)
     - Import ratings (movies.user_ratings)
  4. Scrape watchlist (https://letterboxd.com/{username}/watchlist/)
     - Parse movie titles, years
     - Match to Revenge movies
     - Add to Revenge watchlist
        ↓
Import complete (display stats: "Imported 234 movies, 189 ratings")
```

### Rating Normalization
```
Letterboxd ratings: 0.5 - 5 stars (0.5 increments)
Revenge ratings: 0 - 5 stars (0.5 increments)

Conversion: Direct 1:1 mapping (no conversion needed)
  - Letterboxd 4.5 stars → Revenge 4.5 stars
  - Letterboxd 3 stars → Revenge 3 stars
```

---


<!-- SOURCE-BREADCRUMBS-START -->

## Sources & Cross-References

> Auto-generated section linking to external documentation sources

### Cross-Reference Indexes

- [All Sources Index](../../../sources/SOURCES_INDEX.md) - Complete list of external documentation
- [Design ↔ Sources Map](../../../sources/DESIGN_CROSSREF.md) - Which docs reference which sources

### Referenced Sources

| Source | Documentation |
|--------|---------------|
| [Go io](https://pkg.go.dev/io) | [Local](../../../sources/go/stdlib/io.md) |
| [Last.fm API](https://www.last.fm/api/intro) | [Local](../../../sources/apis/lastfm.md) |
| [Letterboxd API](https://api-docs.letterboxd.com/) | [Local](../../../sources/apis/letterboxd.md) |

<!-- SOURCE-BREADCRUMBS-END -->

<!-- DESIGN-BREADCRUMBS-START -->

## Related Design Docs

> Auto-generated cross-references to related design documentation

**Category**: [Scrobbling](INDEX.md)

### In This Section

- [Last.fm Scrobbling Integration](LASTFM_SCROBBLE.md)
- [ListenBrainz Integration](LISTENBRAINZ.md)
- [Simkl Integration](SIMKL.md)
- [Trakt Integration](TRAKT.md)

### Related Topics

- [Revenge - Architecture v2](../../architecture/01_ARCHITECTURE.md) _Architecture_
- [Revenge - Design Principles](../../architecture/02_DESIGN_PRINCIPLES.md) _Architecture_
- [Revenge - Metadata System](../../architecture/03_METADATA_SYSTEM.md) _Architecture_
- [Revenge - Player Architecture](../../architecture/04_PLAYER_ARCHITECTURE.md) _Architecture_
- [Plugin Architecture Decision](../../architecture/05_PLUGIN_ARCHITECTURE_DECISION.md) _Architecture_

### Indexes

- [Design Index](../../DESIGN_INDEX.md) - All design docs by category/topic
- [Source of Truth](../../00_SOURCE_OF_TRUTH.md) - Package versions and status

<!-- DESIGN-BREADCRUMBS-END -->

## Related Documentation

- [TRAKT.md](./TRAKT.md) - Trakt scrobbling (movies + TV shows, official API)
- [SIMKL.md](./SIMKL.md) - Simkl tracking (alternative to Trakt)
- [TMDB.md](../metadata/video/TMDB.md) - TMDb metadata (movie matching)

---

## Notes

### No Official API (Web Scraping Only)
- **Letterboxd has NO public API**: All data retrieval requires web scraping
- **Scraping risks**: HTML structure changes can break scraper
- **Maintenance**: Regular scraper updates required
- **Priority**: LOW (high maintenance burden, niche audience)

### User-Agent Requirement
- **MUST set User-Agent**: Web scraping requires proper User-Agent
- **Format**: `Revenge/1.0 (https://github.com/lusoris/revenge; contact@example.com)`

### Rate Limits (Conservative)
- **No official limit**: Letterboxd doesn't publish rate limits
- **Very conservative approach**: 1 req/sec (respect server load)
- **Caching**: Cache scraped data aggressively (avoid re-scraping)

### Read-Only Integration
- **Web scraping**: Public data only (no authentication)
- **Import only**: Can import FROM Letterboxd (watch history, ratings, watchlist)
- **No export**: Cannot export TO Letterboxd (no API for write operations)
- **Use case**: One-time import for users migrating from Letterboxd to Revenge

### Movie Matching Strategy
- **Letterboxd data**: Movie title + year
- **Revenge matching**: Lookup movie by title + year → get TMDb/IMDb ID
- **Ambiguity**: Multiple matches → prefer exact title match
- **Unmatched**: Log unmatched movies (manual review/import later)

### Pagination
- **Letterboxd pagination**: URL format `https://letterboxd.com/{username}/films/diary/page/{page_number}/`
- **Scraping**: Iterate pages until no more entries
- **Page size**: ~100 entries per page (Letterboxd default)

### HTML Parsing Challenges
- **Dynamic structure**: HTML structure can change without notice
- **Inconsistent formatting**: Different pages have different formats
- **JavaScript rendering**: Some content may require JavaScript (use headless browser if needed)
- **Fallback**: If scraping fails, skip Letterboxd import (non-critical)

### robots.txt Compliance
- **Check**: https://letterboxd.com/robots.txt
- **Respect directives**: Follow User-agent, Disallow, Crawl-delay rules
- **Legal**: Not legally binding but ethical to respect

### Maintenance Burden
- **Web scraping**: Ongoing maintenance required
- **HTML changes**: Letterboxd can change HTML structure anytime
- **Breaking changes**: Scraper breaks without notice
- **Monitoring**: Implement health checks to detect failures

### Letterboxd vs Trakt
- **Letterboxd**: Movies only, no TV shows, no official API, social features, clean UI
- **Trakt**: Movies + TV shows, official API, scrobbling, bi-directional sync
- **Use case**: Letterboxd for one-time import (migration), Trakt for ongoing scrobbling

### Privacy Considerations
- **Public data only**: Web scraping accesses public Letterboxd profiles
- **User opt-in**: Users must explicitly provide Letterboxd username
- **Data import**: Import data into Revenge (user controls)

### Fallback Strategy (Movie Tracking)
- **Order**: Trakt (primary, official API, bi-directional sync) → Simkl (alternative) → Letterboxd (one-time import, no API)
- **Letterboxd use case**: One-time import for migration (not ongoing sync)

### Future: Official API
- **Letterboxd API**: If Letterboxd releases official public API in future, update integration
- **Priority**: LOW (implement only if official API released)
