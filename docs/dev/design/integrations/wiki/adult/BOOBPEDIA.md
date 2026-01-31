# Boobpedia Integration

> Adult performer encyclopedia with detailed profiles


<!-- TOC-START -->

## Table of Contents

- [Status](#status)
- [Overview](#overview)
- [Developer Resources](#developer-resources)
  - [API Documentation](#api-documentation)
  - [MediaWiki API Endpoints](#mediawiki-api-endpoints)
    - [Search for Performer](#search-for-performer)
    - [Get Page Extract](#get-page-extract)
    - [Get Full Page Content](#get-full-page-content)
    - [Get Page Images](#get-page-images)
  - [Authentication](#authentication)
  - [Data Coverage](#data-coverage)
  - [Go HTTP Client](#go-http-client)
- [API Details](#api-details)
  - [Search Endpoint](#search-endpoint)
  - [Page Extract Endpoint](#page-extract-endpoint)
  - [Full Page Content Endpoint](#full-page-content-endpoint)
  - [Images Endpoint](#images-endpoint)
- [Implementation Checklist](#implementation-checklist)
  - [Phase 1: REST Client (Adult Content - c schema)](#phase-1-rest-client-adult-content---c-schema)
  - [Phase 2: Content Enhancement](#phase-2-content-enhancement)
  - [Phase 3: Background Jobs (River)](#phase-3-background-jobs-river)
- [Integration Pattern](#integration-pattern)
  - [Performer Enrichment Flow](#performer-enrichment-flow)
  - [Rate Limiting Strategy](#rate-limiting-strategy)
- [Sources & Cross-References](#sources-cross-references)
  - [Cross-Reference Indexes](#cross-reference-indexes)
  - [Referenced Sources](#referenced-sources)
- [Related Design Docs](#related-design-docs)
  - [In This Section](#in-this-section)
  - [Related Topics](#related-topics)
  - [Indexes](#indexes)
- [Related Documentation](#related-documentation)
- [Notes](#notes)
  - [MediaWiki API (Same as Wikipedia)](#mediawiki-api-same-as-wikipedia)
  - [Adult Content Isolation (CRITICAL)](#adult-content-isolation-critical)
  - [User-Agent Requirement](#user-agent-requirement)
  - [Rate Limits (Conservative)](#rate-limits-conservative)
  - [Content Licensing](#content-licensing)
  - [Boobpedia vs Wikipedia](#boobpedia-vs-wikipedia)
  - [Boobpedia vs Babepedia](#boobpedia-vs-babepedia)
  - [Extract vs Full Content](#extract-vs-full-content)
  - [Infobox Parsing](#infobox-parsing)
  - [Search Accuracy](#search-accuracy)
  - [Image Quality](#image-quality)
  - [JSONB Storage (c schema)](#jsonb-storage-c-schema)
  - [Caching Strategy](#caching-strategy)
  - [Use Case: Performer Enrichment](#use-case-performer-enrichment)
  - [Content Quality](#content-quality)
  - [Priority: LOW (Alternative Source)](#priority-low-alternative-source)
  - [Fallback Strategy (Adult Performer Metadata)](#fallback-strategy-adult-performer-metadata)

<!-- TOC-END -->

**Service**: Boobpedia (https://www.boobpedia.com)
**API**: MediaWiki Action API (same as Wikipedia)
**Category**: Wiki / Encyclopedia (Adult Content)
**Priority**: 🟡 LOW (Alternative source, niche)

## Status

| Dimension | Status | Notes |
|-----------|--------|-------|
| Design | ✅ | Comprehensive API endpoints, data mapping, JSONB storage |
| Sources | ✅ | MediaWiki API documentation with examples |
| Instructions | ✅ | Phased implementation checklist |
| Code | 🔴 |  |
| Linting | 🔴 |  |
| Unit Testing | 🔴 |  |
| Integration Testing | 🔴 |  |
---

## Overview

**Boobpedia** is a MediaWiki-based encyclopedia dedicated to adult film performers, models, and related content. It uses the same MediaWiki software as Wikipedia, making it easier to integrate (official API available).

**Key Features**:
- **Performer profiles**: Biographies, career info, physical attributes
- **Model profiles**: Non-adult models, glamour models
- **Photos**: Performer images, body shots
- **MediaWiki API**: Same as Wikipedia (official API)
- **Community-edited**: Wiki format, community contributions
- **Adult focus**: Primarily adult performers, some mainstream models

**Use Cases**:
- Performer biography enrichment (alternative to Babepedia)
- Physical attribute metadata (measurements, tattoos, piercings)
- Career timeline tracking
- Photo gallery enrichment
- Alternative performer info source

**⚠️ CRITICAL: Adult Content Isolation**:
- **Database schema**: `c` schema ONLY (`c.performers`)
- **API namespace**: `/api/v1/legacy/wiki/boobpedia/*` (NOT `/api/v1/wiki/boobpedia/*`)
- **Module location**: `internal/content/c/wiki/boobpedia/` (NOT `internal/service/wiki/`)
- **Access control**: Mods/admins can see all data for monitoring, regular users see only their own library

---

## Developer Resources

### API Documentation
- **API**: MediaWiki Action API (same as Wikipedia)
- **Base URL**: https://www.boobpedia.com/wiki/api.php
- **Authentication**: None (public API, no API key required)
- **Rate Limits**: Undefined (use conservative 10 req/sec same as Wikipedia)

### MediaWiki API Endpoints

#### Search for Performer
```
GET https://www.boobpedia.com/wiki/api.php?action=query&list=search&srsearch=Performer+Name&format=json
```

#### Get Page Extract
```
GET https://www.boobpedia.com/wiki/api.php?action=query&prop=extracts&exintro=1&titles=Performer_Name&format=json
```

#### Get Full Page Content
```
GET https://www.boobpedia.com/wiki/api.php?action=parse&page=Performer_Name&format=json
```

#### Get Page Images
```
GET https://www.boobpedia.com/wiki/api.php?action=query&prop=pageimages&titles=Performer_Name&format=json
```

### Authentication
- **Method**: None (public API)
- **User-Agent**: REQUIRED (`User-Agent: Revenge/1.0 (contact@example.com)`)
- **Rate Limits**: Conservative 10 req/sec (same as Wikipedia)

### Data Coverage
- **Performers**: 20K+ performer profiles
- **Coverage**: Western adult performers (US/EU focus)
- **Languages**: Primarily English
- **Updates**: Community-edited (variable frequency)

### Go HTTP Client
- **HTTP client**: `net/http` (stdlib)
- **JSON parsing**: `encoding/json` (stdlib)
- **Reuse Wikipedia client**: Same MediaWiki API (different base URL)

---

## API Details

### Search Endpoint
```
GET /wiki/api.php?action=query&list=search&srsearch={performer_name}&format=json

Response:
{
  "query": {
    "search": [
      {
        "title": "Performer Name",
        "pageid": 12345,
        "snippet": "Biography snippet..."
      }
    ]
  }
}
```

### Page Extract Endpoint
```
GET /wiki/api.php?action=query&prop=extracts&exintro=1&titles={performer_name}&format=json

Response:
{
  "query": {
    "pages": {
      "12345": {
        "pageid": 12345,
        "title": "Performer Name",
        "extract": "Full biography text..."
      }
    }
  }
}
```

### Full Page Content Endpoint
```
GET /wiki/api.php?action=parse&page={performer_name}&format=json

Response:
{
  "parse": {
    "title": "Performer Name",
    "text": {
      "*": "<div>HTML content with infoboxes...</div>"
    }
  }
}
```

### Images Endpoint
```
GET /wiki/api.php?action=query&prop=pageimages&titles={performer_name}&format=json

Response:
{
  "query": {
    "pages": {
      "12345": {
        "title": "Performer Name",
        "thumbnail": {
          "source": "https://www.boobpedia.com/wiki/images/thumb/...",
          "width": 220,
          "height": 330
        }
      }
    }
  }
}
```

---

## Implementation Checklist

### Phase 1: REST Client (Adult Content - c schema)
- [ ] Reuse Wikipedia MediaWiki client (different base URL `https://www.boobpedia.com/wiki/api.php`)
- [ ] User-Agent configuration (REQUIRED)
- [ ] Performer page search (action=query&list=search)
- [ ] Page extract retrieval (prop=extracts&exintro=1)
- [ ] Full page content parsing (action=parse)
- [ ] Page images retrieval (prop=pageimages)
- [ ] **c schema storage**: `c.performers.metadata_json.boobpedia_data` (JSONB)

### Phase 2: Content Enhancement
- [ ] Performer bio extraction (biography, career info)
- [ ] Physical attributes parsing (measurements, height, weight)
- [ ] Tattoo/piercing tracking (parse from infobox)
- [ ] Photo extraction (performer images)
- [ ] Multi-language support (Boobpedia has limited multi-language content)

### Phase 3: Background Jobs (River)
- [ ] **Job**: `c.wiki.boobpedia.fetch_performer` (fetch performer profile)
- [ ] **Job**: `c.wiki.boobpedia.refresh` (periodic refresh)
- [ ] Rate limiting (conservative 10 req/sec)
- [ ] Retry logic (exponential backoff)

---

## Integration Pattern

### Performer Enrichment Flow
```
User views adult performer profile (c.performers)
        ↓
Check if Boobpedia data exists in cache
        ↓
        NO
        ↓
Search Boobpedia API (action=query&list=search&srsearch={performer_name})
        ↓
        MATCH FOUND
        ↓
Get page extract (prop=extracts&exintro=1&titles={performer_name})
        ↓
Parse performer data:
  - Biography
  - Physical attributes (measurements, height, weight)
  - Tattoos/piercings (parse from infobox)
  - Photos (performer images)
        ↓
Store in c.performers.metadata_json.boobpedia_data (c schema JSONB)
        ↓
Display in UI (performer profile page, c schema isolated)
        - Collapsible "Boobpedia" section
        - Attribution: "From Boobpedia"
        - Link to Boobpedia page
```

### Rate Limiting Strategy
```
Boobpedia API: 10 req/sec (conservative, same as Wikipedia)
- Cache API responses for 90 days (performer data changes infrequently)
- Background jobs: Use River queue (low priority)
- Prioritize user-initiated requests over background jobs
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
| [Dragonfly Documentation](https://www.dragonflydb.io/docs) | [Local](../../../../sources/infrastructure/dragonfly.md) |
| [River Job Queue](https://pkg.go.dev/github.com/riverqueue/river) | [Local](../../../../sources/tooling/river.md) |

<!-- SOURCE-BREADCRUMBS-END -->

<!-- DESIGN-BREADCRUMBS-START -->

## Related Design Docs

> Auto-generated cross-references to related design documentation

**Category**: [Adult](INDEX.md)

### In This Section

- [Babepedia Integration](BABEPEDIA.md)
- [IAFD Integration](IAFD.md)

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

- [BABEPEDIA.md](./BABEPEDIA.md) - Adult performer wiki (alternative source)
- [IAFD.md](./IAFD.md) - Internet Adult Film Database (filmography)
- [WIKIPEDIA.md](../WIKIPEDIA.md) - Wikipedia integration (same MediaWiki API)
- [STASHDB.md](../metadata/adult/STASHDB.md) - Primary adult metadata
- [THEPORNDB.md](../metadata/adult/THEPORNDB.md) - Adult metadata provider
- [ADULT_METADATA.md](../../ADULT_METADATA.md) - Adult metadata system architecture

---

## Notes

### MediaWiki API (Same as Wikipedia)
- **Boobpedia uses MediaWiki**: Same software as Wikipedia
- **API compatibility**: Reuse Wikipedia MediaWiki client (different base URL)
- **Endpoints**: Search, page extract, full content, images (same as Wikipedia)

### Adult Content Isolation (CRITICAL)
- **Database schema**: `c` schema ONLY
  - `c.performers.metadata_json.boobpedia_data` (JSONB)
  - NO data in public schema
- **API namespace**: `/api/v1/legacy/wiki/boobpedia/*` (isolated)
  - `/api/v1/legacy/wiki/boobpedia/search/{performer_name}`
  - `/api/v1/legacy/wiki/boobpedia/performers/{performer_id}`
- **Module location**: `internal/content/c/wiki/boobpedia/` (isolated)
- **Access control**: Mods/admins see all, regular users see only their library

### User-Agent Requirement
- **MUST set User-Agent**: MediaWiki API requires User-Agent (HTTP 403 if missing)
- **Format**: `Revenge/1.0 (https://github.com/lusoris/revenge; contact@example.com)`

### Rate Limits (Conservative)
- **No official limit**: Boobpedia doesn't publish rate limits
- **Conservative approach**: 10 req/sec (same as Wikipedia)
- **Caching**: Cache API responses for 90 days (performer data stable)

### Content Licensing
- **License**: Likely CC BY-SA (same as Wikipedia, MediaWiki default)
- **Attribution**: REQUIRED ("From Boobpedia")
- **Link**: Include Boobpedia page URL

### Boobpedia vs Wikipedia
- **Adult focus**: Boobpedia is adult performer-focused (Wikipedia has general bios)
- **Detail level**: Boobpedia has more detailed adult career info
- **Coverage**: Boobpedia covers adult performers (Wikipedia often lacks adult career details)
- **Use case**: Boobpedia for adult-specific info, Wikipedia for general biography

### Boobpedia vs Babepedia
- **API availability**: Boobpedia has official MediaWiki API (Babepedia requires web scraping)
- **Maintenance**: Boobpedia easier to maintain (API stable vs web scraping)
- **Coverage**: Babepedia has more performers (50K+ vs Boobpedia 20K+)
- **Data quality**: Variable (both community-edited)

### Extract vs Full Content
- **Extract** (`prop=extracts&exintro=1`): Plain text introduction (no infoboxes)
- **Full content** (`action=parse`): Complete HTML with infoboxes and images
- **Use case**: Extract for quick biography, full content for detailed parsing

### Infobox Parsing
- **Infoboxes**: HTML tables with performer data (measurements, birthdate, etc.)
- **Parsing**: Extract from HTML `<table class="infobox">`
- **Fields**: Height, weight, measurements, birthdate, birthplace, tattoos, piercings
- **Supplement**: Use infobox data to supplement StashDB/ThePornDB metadata

### Search Accuracy
- **Exact matches**: MediaWiki API returns exact title matches first
- **Disambiguation**: Check if title contains "(disambiguation)"
- **Multiple results**: Show user to select correct performer
- **Year matching**: Prefer results matching performer birth year/career years

### Image Quality
- **Thumbnails**: 220px width (default)
- **Full images**: Available via Wikimedia Commons links
- **Fallback**: If Boobpedia lacks images, use StashDB/ThePornDB images

### JSONB Storage (c schema)
- Store Boobpedia data in `c.performers.metadata_json.boobpedia_data`
- Fields:
  - `page_id`: Boobpedia page ID
  - `title`: Page title
  - `url`: Boobpedia page URL
  - `extract`: Biography extract (plain text)
  - `biography`: Full biography (HTML)
  - `measurements`: Body measurements (e.g., "34D-24-36")
  - `height_cm`: Height in centimeters
  - `weight_kg`: Weight in kilograms
  - `hair_color`: Hair color
  - `eye_color`: Eye color
  - `tattoos`: Array of tattoo descriptions
  - `piercings`: Array of piercing descriptions
  - `photos`: Array of image URLs
  - `last_fetched`: Timestamp

### Caching Strategy
- **Cache duration**: 90 days (performer data changes infrequently)
- **Invalidation**: Manual refresh OR periodic background job (quarterly)
- **Storage**: Store in `c.performers.metadata_json.boobpedia_data` (JSONB) + Dragonfly cache (c namespace)

### Use Case: Performer Enrichment
- **Primary source**: StashDB/ThePornDB for performer metadata
- **Supplement**: Boobpedia for detailed bio, measurements, career timeline
- **Display**: "Performer Bio (Boobpedia)" section in performer profile

### Content Quality
- **Community-edited**: Variable quality (active performers have better coverage)
- **Active moderation**: Popular performers have active moderation
- **Vandalism**: Rare (popular pages protected)
- **Supplementary**: Use as supplementary info (NOT primary metadata)

### Priority: LOW (Alternative Source)
- **Boobpedia is alternative**: StashDB/ThePornDB are primary sources
- **Use case**: Fallback when StashDB/ThePornDB lack data
- **Implementation**: LOW priority (after core features)
- **Niche**: Most users won't need Boobpedia data

### Fallback Strategy (Adult Performer Metadata)
- **Order**: StashDB (primary) → ThePornDB (supplementary) → IAFD (filmography) → Babepedia (bio/measurements) → Boobpedia (alternative/fallback)
