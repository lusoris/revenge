# Pornhub Integration

> Adult content platform with performer pages and scene metadata


<!-- TOC-START -->

## Table of Contents

- [Status](#status)
- [Overview](#overview)
- [Developer Resources](#developer-resources)
  - [API Status](#api-status)
  - [Authentication (Web Scraping)](#authentication-web-scraping)
  - [Data Coverage](#data-coverage)
  - [Go Scraping Library](#go-scraping-library)
- [Integration Approach](#integration-approach)
  - [Web Scraping Strategy](#web-scraping-strategy)
    - [Performer Page Structure](#performer-page-structure)
- [Implementation Checklist](#implementation-checklist)
  - [Phase 1: Web Scraping (Adult Content - c schema)](#phase-1-web-scraping-adult-content---c-schema)
  - [Phase 2: Link Enrichment](#phase-2-link-enrichment)
  - [Phase 3: Background Jobs (River)](#phase-3-background-jobs-river)
- [Integration Pattern](#integration-pattern)
  - [Performer Profile URL Enrichment Flow](#performer-profile-url-enrichment-flow)
- [Sources & Cross-References](#sources-cross-references)
  - [Cross-Reference Indexes](#cross-reference-indexes)
  - [Referenced Sources](#referenced-sources)
- [Related Design Docs](#related-design-docs)
  - [In This Section](#in-this-section)
  - [Related Topics](#related-topics)
  - [Indexes](#indexes)
- [Related Documentation](#related-documentation)
- [Notes](#notes)
  - [Official API (Restricted)](#official-api-restricted)
  - [Adult Content Isolation (CRITICAL)](#adult-content-isolation-critical)
  - [Age Verification Cookie](#age-verification-cookie)
  - [Verified Performers](#verified-performers)
  - [Social Links Extraction](#social-links-extraction)
  - [JSONB Storage (c schema)](#jsonb-storage-c-schema)
  - [Priority: LOW](#priority-low)
  - [Rate Limiting & Blocking](#rate-limiting-blocking)
  - [Caching Strategy](#caching-strategy)
  - [Fallback Strategy (Adult Performer External Links)](#fallback-strategy-adult-performer-external-links)

<!-- TOC-END -->

## Status

| Dimension | Status | Notes |
|-----------|--------|-------|
| Design | ✅ | Web scraping spec, HTML structure, JSONB schema |
| Sources | ✅ | Partner API noted, goquery library documented |
| Instructions | ✅ | Phased implementation checklist with c schema isolation |
| Code | 🔴 |  |
| Linting | 🔴 |  |
| Unit Testing | 🔴 |  |
| Integration Testing | 🔴 |  |
---

## Overview

**Pornhub** is the world's largest adult content platform with performer profile pages, scene metadata, and community features. While primarily a content platform, performer pages provide useful metadata and external links.

**Key Features**:
- **Performer pages**: Profile info, bio, social media links
- **Scene pages**: Video metadata, tags, performers
- **View counts**: Scene popularity metrics
- **Social links**: Twitter, Instagram links on performer pages
- **Verified performers**: Official performer verification

**Use Cases**:
- Performer profile URL enrichment
- Social media link collection
- Scene view count/popularity metrics
- Verified performer status

**⚠️ CRITICAL: Adult Content Isolation**:
- **Database schema**: `c` schema ONLY (`c.performers`, `c.movies`)
- **API namespace**: `/api/v1/legacy/external/pornhub/*` (NOT `/api/v1/external/pornhub/*`)
- **Module location**: `internal/content/c/external/pornhub/` (NOT `internal/service/external/`)
- **Access control**: Mods/admins can see all data for monitoring, regular users see only their own library

---

## Developer Resources

### API Status
- **Official API**: EXISTS (requires partnership agreement)
  - **URL**: https://www.pornhub.com/partners (content partner program)
  - **Access**: Restricted (requires approval)
  - **Use case**: Content partners only (not for media servers)
- **Alternative**: Web scraping (parse HTML)

### Authentication (Web Scraping)
- **Method**: None (public performer pages)
- **User-Agent**: REQUIRED (`User-Agent: Revenge/1.0 (contact@example.com)`)
- **Rate Limits**: Very conservative (1 req/sec, avoid being blocked)
- **Cookies**: May require cookie handling (age verification)

### Data Coverage
- **Performers**: 100K+ verified performers
- **Scenes**: Millions of scenes
- **Focus**: View counts, popularity, social links

### Go Scraping Library
- **Recommended**: `github.com/PuerkitoBio/goquery`

---

## Integration Approach

### Web Scraping Strategy

**⚠️ CRITICAL: Web Scraping Only (Official API restricted)**

#### Performer Page Structure
```
URL: https://www.pornhub.com/pornstar/{performer_name}

HTML Structure (example):
<div class="performer-page">
  <h1 class="performer-name">Performer Name</h1>
  <span class="verified-badge">✓ Verified</span>

  <div class="performer-bio">
    <p>Biography text...</p>
  </div>

  <div class="performer-stats">
    <span class="video-views">123.4M views</span>
    <span class="video-count">234 videos</span>
    <span class="rank">Rank: #42</span>
  </div>

  <div class="social-links">
    <a href="https://twitter.com/performer" class="twitter">Twitter</a>
    <a href="https://instagram.com/performer" class="instagram">Instagram</a>
  </div>

  <div class="performer-info">
    <p><strong>Born:</strong> May 15, 1990</p>
    <p><strong>Career Start:</strong> 2015</p>
  </div>
</div>
```

---

## Implementation Checklist

### Phase 1: Web Scraping (Adult Content - c schema)
- [ ] HTML scraping setup (`goquery`)
- [ ] User-Agent configuration (REQUIRED)
- [ ] Cookie handling (age verification)
- [ ] URL construction (performer name → Pornhub URL)
- [ ] Verified status extraction (verified badge)
- [ ] Social links extraction (Twitter, Instagram)
- [ ] Stats extraction (view counts, video count, rank)
- [ ] Bio extraction (career start, birthdate if available)
- [ ] **c schema storage**: `c.performers.metadata_json.pornhub_data` (JSONB)
- [ ] **c schema storage**: `c.performers.external_urls` (Pornhub profile link)

### Phase 2: Link Enrichment
- [ ] Pornhub profile URL (add to performer external URLs)
- [ ] Social media links (Twitter, Instagram from Pornhub page)
- [ ] Verified status flag (c.performers.pornhub_verified)

### Phase 3: Background Jobs (River)
- [ ] **Job**: `c.external.pornhub.scrape_performer` (scrape performer profile)
- [ ] **Job**: `c.external.pornhub.refresh` (periodic refresh)
- [ ] Rate limiting (very conservative 1 req/sec)
- [ ] Retry logic (exponential backoff)
- [ ] Cookie management (refresh age verification cookies)

---

## Integration Pattern

### Performer Profile URL Enrichment Flow
```
User views adult performer profile (c.performers)
        ↓
Check if Pornhub profile URL exists
        ↓
        NO
        ↓
Construct Pornhub URL (https://www.pornhub.com/pornstar/{name})
        ↓
Scrape performer page (with age verification cookie)
        ↓
Parse data:
  - Verified status
  - Social links (Twitter, Instagram)
  - Stats (view counts, video count, rank)
  - Bio (career start, birthdate)
        ↓
Store:
  - c.performers.metadata_json.pornhub_data (stats, bio)
  - c.performers.external_urls (Pornhub profile link)
  - c.performers.pornhub_verified (verified status flag)
        ↓
Display in UI:
  - External links section (Pornhub icon → profile link)
  - Verified badge (if verified performer)
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
| [Go io](https://pkg.go.dev/io) | [Local](../../../../sources/go/stdlib/io.md) |
| [River Job Queue](https://pkg.go.dev/github.com/riverqueue/river) | [Local](../../../../sources/tooling/river.md) |

<!-- SOURCE-BREADCRUMBS-END -->

<!-- DESIGN-BREADCRUMBS-START -->

## Related Design Docs

> Auto-generated cross-references to related design documentation

**Category**: [Adult](INDEX.md)

### In This Section

- [FreeOnes Integration](FREEONES.md)
- [Instagram Integration](INSTAGRAM.md)
- [OnlyFans Integration](ONLYFANS.md)
- [TheNude Integration](THENUDE.md)
- [Twitter/X Integration](TWITTER_X.md)

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

- [FREEONES.md](./FREEONES.md) - FreeOnes performer database
- [ONLYFANS.md](./ONLYFANS.md) - OnlyFans content platform
- [STASHDB.md](../metadata/adult/STASHDB.md) - Primary adult metadata

---

## Notes

### Official API (Restricted)
- **Pornhub API**: Exists but restricted to content partners
- **Access**: Requires partnership agreement (not for media servers)
- **Alternative**: Web scraping (public performer pages)

### Adult Content Isolation (CRITICAL)
- **Database schema**: `c` schema ONLY
  - `c.performers.metadata_json.pornhub_data` (JSONB)
  - `c.performers.external_urls` (Pornhub profile link)
  - `c.performers.pornhub_verified` (BOOLEAN verified flag)
  - NO data in public schema
- **API namespace**: `/api/v1/legacy/external/pornhub/*` (isolated)
- **Module location**: `internal/content/c/external/pornhub/` (isolated)
- **Access control**: Mods/admins see all, regular users see only their library

### Age Verification Cookie
- **Pornhub**: Requires age verification (18+ cookie)
- **Cookie handling**: Set cookie `age_verified=1` OR handle age gate redirect
- **Scraping**: Must handle cookie to access performer pages

### Verified Performers
- **Verified badge**: Performers with verified accounts (official performers)
- **Trust**: Verified performers provide more accurate info (social links, bio)
- **Storage**: Store verified status in `c.performers.pornhub_verified` (BOOLEAN)

### Social Links Extraction
- **Twitter**: Extract from `<a class="twitter">` OR `href` contains "twitter.com"
- **Instagram**: Extract from `<a class="instagram">` OR `href` contains "instagram.com"
- **Validation**: Cross-check with FreeOnes/StashDB social links

### JSONB Storage (c schema)
```json
{
  "pornhub_url": "https://www.pornhub.com/pornstar/performer-name",
  "verified": true,
  "stats": {
    "video_views": 123400000,
    "video_count": 234,
    "rank": 42
  },
  "social_links": {
    "twitter": "https://twitter.com/performer",
    "instagram": "https://instagram.com/performer"
  },
  "bio": "Biography text...",
  "career_start": 2015,
  "last_scraped": "2023-01-15T10:00:00Z"
}
```

### Priority: LOW
- **Pornhub**: Not primary metadata source (FreeOnes/StashDB are primary)
- **Use case**: Primarily for Pornhub profile link + verified status
- **Implementation**: LOW priority (after core features)

### Rate Limiting & Blocking
- **Very conservative**: 1 req/sec (Pornhub actively blocks scrapers)
- **User-Agent**: Must set proper User-Agent
- **Respect robots.txt**: Check https://www.pornhub.com/robots.txt
- **Fallback**: If blocked, skip Pornhub integration (non-critical)

### Caching Strategy
- **Cache duration**: 90 days (performer pages change infrequently)
- **Stats refresh**: Periodic refresh (quarterly) for view counts/rank

### Fallback Strategy (Adult Performer External Links)
- **Order**: StashDB (primary urls field) → FreeOnes (comprehensive external links) → Pornhub (profile link + social links) → TheNude (cross-references)
