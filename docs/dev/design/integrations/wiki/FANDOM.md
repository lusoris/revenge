# FANDOM Integration

> Fan-curated wikis for movies, TV shows, games, and more


<!-- TOC-START -->

## Table of Contents

- [Status](#status)
- [Overview](#overview)
- [Developer Resources](#developer-resources)
  - [API Documentation](#api-documentation)
  - [Authentication](#authentication)
  - [Data Coverage](#data-coverage)
  - [Go Client Library](#go-client-library)
- [API Details](#api-details)
  - [REST Endpoints (Same as Wikipedia)](#rest-endpoints-same-as-wikipedia)
    - [Search Pages](#search-pages)
    - [Get Page Extract](#get-page-extract)
    - [Get Page Images](#get-page-images)
- [Implementation Checklist](#implementation-checklist)
  - [Phase 1: Core Integration](#phase-1-core-integration)
  - [Phase 2: Wiki Mapping](#phase-2-wiki-mapping)
  - [Phase 3: Background Jobs (River)](#phase-3-background-jobs-river)
- [Integration Pattern](#integration-pattern)
  - [Wiki Selection Flow](#wiki-selection-flow)
  - [Wiki Mapping Table](#wiki-mapping-table)
- [Sources & Cross-References](#sources-cross-references)
  - [Cross-Reference Indexes](#cross-reference-indexes)
  - [Referenced Sources](#referenced-sources)
- [Related Design Docs](#related-design-docs)
  - [In This Section](#in-this-section)
  - [Related Topics](#related-topics)
  - [Indexes](#indexes)
- [Related Documentation](#related-documentation)
- [Notes](#notes)
  - [Wiki Selection Strategy](#wiki-selection-strategy)
  - [Popular Wiki Mappings](#popular-wiki-mappings)
  - [FANDOM vs Wikipedia](#fandom-vs-wikipedia)
  - [Content Licensing](#content-licensing)
  - [User-Agent Requirement](#user-agent-requirement)
  - [Rate Limits](#rate-limits)
  - [JSONB Storage](#jsonb-storage)
  - [Caching Strategy](#caching-strategy)
  - [Content Quality](#content-quality)
  - [Episode Guides](#episode-guides)
  - [Character Bios](#character-bios)
  - [Spoilers Warning](#spoilers-warning)
  - [Multi-Language Support](#multi-language-support)
  - [Fallback Strategy](#fallback-strategy)

<!-- TOC-END -->

**Service**: FANDOM (https://www.fandom.com)
**API**: MediaWiki Action API (same as Wikipedia)
**Category**: Wiki / Knowledge Base (Fan Communities)
**Priority**: 🟢 MEDIUM (Niche content)

## Status

| Dimension | Status | Notes |
|-----------|--------|-------|
| Design | ✅ | Comprehensive API endpoints, wiki mapping, JSONB storage |
| Sources | ✅ | MediaWiki API documentation with examples |
| Instructions | ✅ | Phased implementation checklist |
| Code | 🔴 |  |
| Linting | 🔴 |  |
| Unit Testing | 🔴 |  |
| Integration Testing | 🔴 |  |
---

## Overview

**FANDOM** (formerly Wikia) is a wiki hosting platform with 250,000+ fan-curated wikis covering movies, TV shows, games, books, and more. Each wiki is community-managed and highly specialized.

**Key Features**:
- **Fan wikis**: Community-driven wikis for specific franchises
- **Detailed content**: Episode guides, character bios, plot details, behind-the-scenes
- **Rich media**: Images, videos, fan art
- **MediaWiki API**: Same API as Wikipedia (MediaWiki Action API)
- **Free access**: No API key required (rate-limited)

**Popular FANDOM Wikis**:
- **MCU**: Marvel Cinematic Universe Wiki (https://marvelcinematicuniverse.fandom.com)
- **Star Trek**: Memory Alpha (https://memory-alpha.fandom.com)
- **Star Wars**: Wookieepedia (https://starwars.fandom.com)
- **Harry Potter**: Harry Potter Wiki (https://harrypotter.fandom.com)
- **Game of Thrones**: Game of Thrones Wiki (https://gameofthrones.fandom.com)
- **The Witcher**: Witcher Wiki (https://witcher.fandom.com)

**Use Cases**:
- Franchise-specific information (MCU, Star Trek, Star Wars)
- Episode-by-episode guides
- Character biographies and relationships
- Behind-the-scenes trivia
- Fan theories and analysis

---

## Developer Resources

### API Documentation
- **Base URL**: `https://{wiki}.fandom.com/api.php` (e.g., `marvelcinematicuniverse.fandom.com`)
- **API Docs**: https://community.fandom.com/wiki/Help:Fandom_API (same as MediaWiki API)
- **Rate Limits**: Not officially documented (use conservative limits ~10 req/sec)

### Authentication
- **Method**: None (public API)
- **User-Agent**: REQUIRED (`User-Agent: Revenge/1.0 (contact@example.com)`)
- **Rate Limits**: Conservative 10 req/sec (no official limit)

### Data Coverage
- **Wikis**: 250,000+ wikis
- **Articles**: Millions of fan-curated pages
- **Languages**: Primarily English, some multi-language wikis
- **Updates**: Real-time (community-edited)

### Go Client Library
- **Official**: None
- **Recommended**: Use `net/http` with JSON parsing (same as Wikipedia)

---

## API Details

### REST Endpoints (Same as Wikipedia)

#### Search Pages
```
GET https://{wiki}.fandom.com/api.php?action=query&list=search&srsearch={query}&format=json

Example (MCU Wiki):
https://marvelcinematicuniverse.fandom.com/api.php?action=query&list=search&srsearch=Iron%20Man&format=json

Response:
{
  "query": {
    "search": [
      {
        "pageid": 1234,
        "title": "Iron Man",
        "snippet": "Tony Stark aka <b>Iron Man</b> is a genius billionaire..."
      }
    ]
  }
}
```

#### Get Page Extract
```
GET https://{wiki}.fandom.com/api.php?action=query&prop=extracts&exintro=1&titles={title}&format=json

Example (Memory Alpha - Star Trek):
https://memory-alpha.fandom.com/api.php?action=query&prop=extracts&exintro=1&titles=Jean-Luc%20Picard&format=json
```

#### Get Page Images
```
GET https://{wiki}.fandom.com/api.php?action=query&prop=pageimages&titles={title}&format=json
```

---

## Implementation Checklist

### Phase 1: Core Integration
- [ ] REST API client setup (reuse Wikipedia client, different base URL)
- [ ] Wiki selection mapping (movie → MCU Wiki, Star Trek → Memory Alpha, etc.)
- [ ] User-Agent configuration (REQUIRED)
- [ ] Page search (same as Wikipedia)
- [ ] Page extract fetch
- [ ] Image fetch
- [ ] JSONB storage (`metadata_json.fandom_data`)

### Phase 2: Wiki Mapping
- [ ] **Movie franchises**: MCU Wiki, DC Extended Universe Wiki, Star Wars Wiki, etc.
- [ ] **TV shows**: Memory Alpha (Star Trek), Wookieepedia (Star Wars), Breaking Bad Wiki, etc.
- [ ] **Genre-based**: Fallback to genre-specific wikis (Sci-Fi, Fantasy, Horror)
- [ ] User preference: Allow users to select preferred wiki

### Phase 3: Background Jobs (River)
- [ ] **Job**: `wiki.fandom.fetch_content` (fetch wiki page)
- [ ] **Job**: `wiki.fandom.refresh` (periodic refresh)
- [ ] Rate limiting (conservative 10 req/sec)
- [ ] Retry logic

---

## Integration Pattern

### Wiki Selection Flow
```
User views movie/TV show page (e.g., "Iron Man")
        ↓
Determine franchise/genre
  - MCU movies → MCU Wiki
  - Star Trek → Memory Alpha
  - Star Wars → Wookieepedia
  - Harry Potter → Harry Potter Wiki
  - Generic → Skip (no specific wiki)
        ↓
Search FANDOM wiki (action=query&list=search)
        ↓
Match found? → Get page extract
              ↓
              Store in metadata_json.fandom_data
              ↓
              Display in UI ("FANDOM Wiki" section)
        ↓
        NO MATCH
        ↓
Skip (no FANDOM data)
```

### Wiki Mapping Table
```
movies table:
  - tmdb_id: 1726 (Iron Man)
  - fandom_wiki: "marvelcinematicuniverse"
  - fandom_page: "Iron Man"

tvshows table:
  - thetvdb_id: 253463 (Star Trek: Discovery)
  - fandom_wiki: "memory-alpha"
  - fandom_page: "Star Trek: Discovery"
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
| [River Job Queue](https://pkg.go.dev/github.com/riverqueue/river) | [Local](../../../sources/tooling/river.md) |

<!-- SOURCE-BREADCRUMBS-END -->

<!-- DESIGN-BREADCRUMBS-START -->

## Related Design Docs

> Auto-generated cross-references to related design documentation

**Category**: [Wiki](INDEX.md)

### In This Section

- [TVTropes Integration](TVTROPES.md)
- [Wikipedia Integration](WIKIPEDIA.md)

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

- [WIKIPEDIA.md](./WIKIPEDIA.md) - General encyclopedia
- [TVTROPES.md](./TVTROPES.md) - Trope analysis
- [Wiki System](../../features/shared/WIKI_SYSTEM.md) - Built-in wiki system

---

## Notes

### Wiki Selection Strategy
- **Automatic detection**: Map franchise/genre to specific FANDOM wiki
- **User preference**: Allow users to override wiki selection
- **Fallback**: Wikipedia if no specific FANDOM wiki available

### Popular Wiki Mappings
- **MCU movies**: `marvelcinematicuniverse.fandom.com`
- **DC movies**: `dc.fandom.com`
- **Star Wars**: `starwars.fandom.com` (Wookieepedia)
- **Star Trek**: `memory-alpha.fandom.com`
- **Harry Potter**: `harrypotter.fandom.com`
- **Game of Thrones**: `gameofthrones.fandom.com`
- **The Witcher**: `witcher.fandom.com`
- **Breaking Bad**: `breakingbad.fandom.com`
- **Stranger Things**: `strangerthings.fandom.com`

### FANDOM vs Wikipedia
- **FANDOM**: Franchise-specific, detailed, fan-driven
- **Wikipedia**: General encyclopedia, neutral POV, broader coverage
- **Use case**: FANDOM for niche content (character bios, episode guides), Wikipedia for general info

### Content Licensing
- **License**: Creative Commons Attribution-Share Alike License (CC BY-SA)
- **Attribution**: MUST attribute FANDOM wiki in UI ("From {Wiki Name} on FANDOM")
- **Link**: Include FANDOM wiki URL

### User-Agent Requirement
- **MUST set User-Agent**: FANDOM uses MediaWiki (same requirement as Wikipedia)
- **Format**: `Revenge/1.0 (https://github.com/lusoris/revenge; contact@example.com)`

### Rate Limits
- **No official limit**: FANDOM doesn't publish rate limits
- **Conservative approach**: 10 req/sec (same as Wikipedia conservative limit)
- **Abuse prevention**: Excessive usage may result in IP ban

### JSONB Storage
- Store FANDOM data in `metadata_json.fandom_data`
- Fields:
  - `wiki`: Wiki subdomain (e.g., "marvelcinematicuniverse")
  - `page_id`: FANDOM page ID
  - `title`: Page title
  - `extract`: Summary text
  - `url`: FANDOM wiki URL
  - `thumbnail`: Image URL

### Caching Strategy
- **Cache duration**: 30 days (FANDOM content changes infrequently)
- **Invalidation**: Manual refresh OR automatic on content update

### Content Quality
- **Fan-driven**: Content quality varies by wiki
- **Popular wikis**: High quality (MCU, Star Wars, Star Trek - active moderation)
- **Obscure wikis**: Variable quality (less active communities)
- **Use case**: Supplementary information for fans

### Episode Guides
- **TV shows**: FANDOM wikis often have detailed episode guides
- **Use case**: Display episode summaries, trivia, behind-the-scenes
- **Integration**: Link FANDOM episode pages from Revenge episode pages

### Character Bios
- **Detailed bios**: FANDOM wikis excel at character biographies
- **Relationships**: Character relationships, family trees
- **Use case**: Display character bios in performer/character pages

### Spoilers Warning
- **Spoilers**: FANDOM wikis often contain spoilers
- **UI warning**: Display spoiler warning before showing FANDOM content
- **User preference**: Allow users to hide FANDOM content (spoiler-averse users)

### Multi-Language Support
- **Primarily English**: Most FANDOM wikis are English-language
- **Some multi-language**: Popular franchises have multi-language wikis
- **Detection**: Check if wiki has localized version (e.g., `de.memory-alpha.fandom.com`)

### Fallback Strategy
- **Order**: TMDb/TheTVDB (primary) → FANDOM (franchise-specific) → Wikipedia (general) → TVTropes (analysis)
