# FreeOnes Integration

> Adult performer database with comprehensive profiles and links

## Status

| Dimension | Status | Notes |
| --------- | ------ | ----- |
| Design | ✅ | Comprehensive web scraping spec, HTML structure, JSONB schema |
| Sources | ✅ | Website URL, goquery library documented |
| Instructions | ✅ | Detailed phased implementation checklist with c schema isolation |
| Code | 🔴 | |
| Linting | 🔴 | |
| Unit Testing | 🔴 | |
| Integration Testing | 🔴 | |

---

## Overview

**FreeOnes** is one of the largest and most comprehensive adult performer databases on the internet, featuring detailed profiles, biographies, photos, external links, and social media connections for performers.

**Key Features**:
- **Performer profiles**: 100K+ performers with detailed bios
- **External links**: Links to performer social media, fan sites, official sites
- **Photos**: Large photo galleries (headshots, body shots)
- **Measurements**: Physical attributes, tattoos, piercings
- **Career info**: Career start/end, aliases, awards
- **Social media**: Twitter, Instagram, OnlyFans, Fansly links
- **Awards**: AVN, XBIZ, and other industry awards

**Use Cases**:
- Performer profile enrichment (bio, career info)
- External link collection (social media, fan sites, official sites)
- Photo gallery enrichment
- Award information tracking
- Alternative name/alias resolution

**⚠️ CRITICAL: Adult Content Isolation**:
- **Database schema**: `c` schema ONLY (`c.performers`)
- **API namespace**: `/api/v1/legacy/external/freeones/*` (NOT `/api/v1/external/freeones/*`)
- **Module location**: `internal/content/c/external/freeones/` (NOT `internal/service/external/`)
- **Access control**: Mods/admins can see all data for monitoring, regular users see only their own library

---

## Developer Resources

### API Status
- **Official API**: NONE (no official API)
- **Web Scraping**: Required (parse HTML)
- **Base URL**: https://www.freeones.com/performers/{performer_name}

### Authentication
- **Method**: None (public website)
- **User-Agent**: REQUIRED (`User-Agent: Revenge/1.0 (contact@example.com)`)
- **Rate Limits**: Very conservative (1 req/sec)

### Data Coverage
- **Performers**: 100K+ performer profiles
- **Coverage**: Western performers (US/EU focus), some international
- **Updates**: Community-driven + admin curation
- **Quality**: High (major database, well-maintained)

### Go Scraping Library
- **Recommended**: `github.com/PuerkitoBio/goquery` (jQuery-like HTML parsing)

---

## Integration Approach

### Web Scraping Strategy

**⚠️ CRITICAL: No Official API - Web Scraping Required**

#### Performer Page Structure
```
URL: https://www.freeones.com/performers/{performer_name}

HTML Structure (example):
<div class="performer-page">
  <h1 class="performer-name">Performer Name</h1>

  <div class="performer-bio">
    <p><strong>Born:</strong> May 15, 1990 (33 years old)</p>
    <p><strong>Birthplace:</strong> Los Angeles, CA, USA</p>
    <p><strong>Aliases:</strong> Alias1, Alias2</p>
    <p><strong>Career:</strong> 2015 - Active</p>
    <p><strong>Ethnicity:</strong> Caucasian</p>
  </div>

  <div class="measurements">
    <p><strong>Height:</strong> 5'6" (168 cm)</p>
    <p><strong>Weight:</strong> 125 lbs (57 kg)</p>
    <p><strong>Measurements:</strong> 34D-24-36</p>
    <p><strong>Hair:</strong> Blonde</p>
    <p><strong>Eyes:</strong> Blue</p>
    <p><strong>Tattoos:</strong> Lower back, left shoulder</p>
    <p><strong>Piercings:</strong> Navel, ears</p>
  </div>

  <div class="external-links">
    <h2>External Links</h2>
    <ul>
      <li><a href="https://twitter.com/performer" class="twitter-link">Twitter</a></li>
      <li><a href="https://instagram.com/performer" class="instagram-link">Instagram</a></li>
      <li><a href="https://onlyfans.com/performer" class="onlyfans-link">OnlyFans</a></li>
      <li><a href="https://performer-official.com" class="official-site">Official Site</a></li>
    </ul>
  </div>

  <div class="awards">
    <h2>Awards</h2>
    <ul>
      <li>AVN Award - Best New Starlet (2016)</li>
      <li>XBIZ Award - Female Performer of the Year (2018)</li>
    </ul>
  </div>

  <div class="photo-gallery">
    <img src="/photos/performer-1.jpg">
    <img src="/photos/performer-2.jpg">
  </div>
</div>
```

---

## Implementation Checklist

### Phase 1: Web Scraping (Adult Content - c schema)
- [ ] HTML scraping setup (`goquery`)
- [ ] User-Agent configuration (REQUIRED)
- [ ] URL normalization (performer name → FreeOnes URL)
- [ ] Bio extraction (birthdate, ethnicity, career dates, aliases)
- [ ] Measurements extraction (height, weight, body measurements)
- [ ] External links extraction (Twitter, Instagram, OnlyFans, Fansly, official sites)
- [ ] Awards extraction (AVN, XBIZ, etc.)
- [ ] Photo extraction (profile pics, gallery images)
- [ ] **c schema storage**: `c.performers.metadata_json.freeones_data` (JSONB)

### Phase 2: Link Enrichment
- [ ] Social media links (Twitter, Instagram)
- [ ] Content platform links (OnlyFans, Fansly, ManyVids)
- [ ] Official site links (performer websites)
- [ ] Fan site links (performer fan clubs)
- [ ] Store in `c.performers.external_urls` (separate table for links)

### Phase 3: Background Jobs (River)
- [ ] **Job**: `c.external.freeones.scrape_performer` (scrape performer profile)
- [ ] **Job**: `c.external.freeones.refresh` (periodic refresh)
- [ ] Rate limiting (very conservative 1 req/sec)
- [ ] Retry logic (exponential backoff)

---

## Integration Pattern

### Performer Enrichment Flow
```
User views adult performer profile (c.performers)
        ↓
Check if FreeOnes data exists in cache
        ↓
        NO
        ↓
Construct FreeOnes URL (https://www.freeones.com/performers/{name})
        ↓
Scrape HTML page (goquery)
        ↓
Parse performer data:
  - Bio (birthdate, ethnicity, career dates, aliases)
  - Measurements (height, weight, body measurements)
  - External links (Twitter, Instagram, OnlyFans, official sites)
  - Awards (AVN, XBIZ, etc.)
  - Photos (profile pics, gallery images)
        ↓
Store in c.performers.metadata_json.freeones_data (c schema JSONB)
Store external links in c.performers.external_urls (c schema)
        ↓
Display in UI (performer profile page, c schema isolated)
        - External links section (Twitter, Instagram, OnlyFans icons)
        - Awards section
        - Bio section
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

<!-- SOURCE-BREADCRUMBS-END -->

<!-- DESIGN-BREADCRUMBS-START -->

## Related Design Docs

> Auto-generated cross-references to related design documentation

**Category**: [Adult](INDEX.md)

### In This Section

- [Instagram Integration](INSTAGRAM.md)
- [OnlyFans Integration](ONLYFANS.md)
- [Pornhub Integration](PORNHUB.md)
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

- [THENUDE.md](./THENUDE.md) - TheNude performer database (alternative)
- [PORNHUB.md](./PORNHUB.md) - Pornhub performer platform pages
- [STASHDB.md](../metadata/adult/STASHDB.md) - Primary adult metadata (urls field tracks FreeOnes links)
- [BABEPEDIA.md](../wiki/adult/BABEPEDIA.md) - Adult performer wiki

---

## Notes

### No Official API (Web Scraping Only)
- **FreeOnes has NO API**: All data retrieval requires web scraping
- **Scraping risks**: HTML structure changes can break scraper
- **Maintenance**: Regular scraper updates required
- **Priority**: MEDIUM (useful for external links, bio enrichment)

### Adult Content Isolation (CRITICAL)
- **Database schema**: `c` schema ONLY
  - `c.performers.metadata_json.freeones_data` (JSONB)
  - `c.performers.external_urls` (table for external links)
  - NO data in public schema
- **API namespace**: `/api/v1/legacy/external/freeones/*` (isolated)
  - `/api/v1/legacy/external/freeones/search/{performer_name}`
  - `/api/v1/legacy/external/freeones/performers/{performer_id}`
- **Module location**: `internal/content/c/external/freeones/` (isolated)
- **Access control**: Mods/admins see all, regular users see only their library

### External Links Table Schema
```sql
CREATE TABLE c.performer_external_urls (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  performer_id UUID NOT NULL REFERENCES c.performers(id) ON DELETE CASCADE,
  source VARCHAR(50) NOT NULL, -- 'freeones', 'stashdb', etc.
  platform VARCHAR(50) NOT NULL, -- 'twitter', 'instagram', 'onlyfans', etc.
  url TEXT NOT NULL,
  verified BOOLEAN DEFAULT FALSE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE(performer_id, platform, url)
);

CREATE INDEX idx_performer_external_urls_performer_id ON c.performer_external_urls(performer_id);
CREATE INDEX idx_performer_external_urls_platform ON c.performer_external_urls(platform);
```

### URL Extraction
- **Twitter**: Extract from `<a class="twitter-link">` OR `href` contains "twitter.com"
- **Instagram**: Extract from `<a class="instagram-link">` OR `href` contains "instagram.com"
- **OnlyFans**: Extract from `<a class="onlyfans-link">` OR `href` contains "onlyfans.com"
- **Official sites**: Extract from `<a class="official-site">` OR rel="external"
- **Normalization**: Clean URLs (remove tracking params, normalize format)

### Photo Gallery
- **FreeOnes photos**: Large photo galleries (headshots, body shots, glamour shots)
- **Download**: Download photos to Revenge storage (with attribution)
- **Storage**: Store in `c.performers.photos` table (photo URLs + local paths)

### Awards Tracking
- **AVN Awards**: Adult Video News Awards (industry standard)
- **XBIZ Awards**: XBIZ Industry Awards (second major award)
- **Other awards**: XRCO, NightMoves, etc.
- **Storage**: Store in `c.performers.metadata_json.freeones_data.awards` (JSONB array)

### Alias Resolution
- **Alternative names**: FreeOnes tracks comprehensive alias lists
- **Matching**: Use aliases for performer matching/searching
- **Normalization**: Normalize performer names (lowercase, remove special chars)

### JSONB Storage (c schema)
- Store FreeOnes data in `c.performers.metadata_json.freeones_data`
- Fields:
  - `freeones_url`: FreeOnes profile URL
  - `birthdate`: Date of birth
  - `birthplace`: Place of birth
  - `ethnicity`: Ethnicity
  - `career_start`: Career start year
  - `career_end`: Career end year (null if active)
  - `aliases`: Array of aliases
  - `height_cm`: Height in centimeters
  - `weight_kg`: Weight in kilograms
  - `measurements`: Body measurements (e.g., "34D-24-36")
  - `hair_color`: Hair color
  - `eye_color`: Eye color
  - `tattoos`: Array of tattoo descriptions
  - `piercings`: Array of piercing descriptions
  - `awards`: Array of award objects (year, award name, category)
  - `last_scraped`: Timestamp

### Caching Strategy
- **Cache duration**: 90 days (performer data changes infrequently)
- **Invalidation**: Manual refresh OR periodic background job (quarterly)
- **Storage**: Store in `c.performers.metadata_json.freeones_data` (JSONB) + Dragonfly cache (c namespace)

### Robots.txt Compliance
- **Check**: https://www.freeones.com/robots.txt
- **Respect directives**: Follow User-agent, Disallow, Crawl-delay rules
- **Ethical**: Not legally binding but ethical to respect

### Fallback Strategy (Adult Performer Metadata)
- **Order**: StashDB (primary) → ThePornDB (supplementary) → FreeOnes (external links, bio) → IAFD (filmography) → Babepedia (alternative)
- **Use case**: FreeOnes excels at external links (Twitter, Instagram, OnlyFans) and awards tracking
