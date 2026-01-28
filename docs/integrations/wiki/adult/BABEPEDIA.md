# Babepedia Integration

> Adult performer wiki with biographies and filmographies

**Service**: Babepedia (https://www.babepedia.com)
**API**: None (web scraping required)
**Category**: Wiki / Knowledge Base (Adult Content)
**Priority**: 🟢 MEDIUM (Adult performer info)
**Status**: 🔴 DESIGN PHASE

---

## Overview

**Babepedia** is a community-driven wiki dedicated to adult film performers with comprehensive biographies, filmographies, photos, and measurements. It provides detailed information about performers' careers, physical attributes, and personal information.

**Key Features**:
- **Performer profiles**: Detailed bios, career info, aliases
- **Measurements**: Height, weight, body measurements, tattoos, piercings
- **Filmography**: Complete list of scenes/movies
- **Photos**: Performer headshots and body shots
- **Career timeline**: Debut date, retirement date, active years
- **No official API**: Requires web scraping

**Use Cases**:
- Performer biography enrichment
- Physical attribute metadata
- Career timeline tracking
- Filmography completion
- Alternative name/alias resolution

**⚠️ CRITICAL: Adult Content Isolation**:
- **Database schema**: `c` schema ONLY (`c.performers`, `c.adult_movies`)
- **API namespace**: `/api/v1/c/wiki/babepedia/*` (NOT `/api/v1/wiki/babepedia/*`)
- **Module location**: `internal/content/c/wiki/babepedia/` (NOT `internal/service/wiki/`)
- **Access control**: Mods/admins can see all data for monitoring, regular users see only their own library

---

## Developer Resources

### API Documentation
- **API**: NONE (no official API)
- **Web Scraping**: Required (parse HTML)
- **Base URL**: https://www.babepedia.com/babe/{PerformerName}
- **Rate Limits**: Undefined (use conservative scraping ~1 req/sec)

### Authentication
- **Method**: None (public website)
- **User-Agent**: REQUIRED (`User-Agent: Revenge/1.0 (contact@example.com)`)
- **Rate Limits**: Conservative 1 req/sec (avoid server overload)

### Data Coverage
- **Performers**: 50K+ performer profiles
- **Coverage**: Western adult performers (US/EU focus)
- **Languages**: Primarily English
- **Updates**: Community-edited (variable frequency)

### Go Scraping Library
- **Recommended**: `github.com/PuerkitoBio/goquery` (jQuery-like HTML parsing)
- **Alternative**: `github.com/gocolly/colly` (web scraping framework)

---

## Integration Approach

### Web Scraping Strategy

**⚠️ CRITICAL: No Official API - Web Scraping Required**

#### Performer Page Structure
```
URL: https://www.babepedia.com/babe/Performer_Name

HTML Structure (example):
<div class="performer-profile">
  <h1>Performer Name</h1>

  <div class="bio-section">
    <p><strong>Real Name:</strong> Jane Doe</p>
    <p><strong>Aliases:</strong> Alias 1, Alias 2</p>
    <p><strong>Born:</strong> May 15, 1990</p>
    <p><strong>Birthplace:</strong> Los Angeles, CA, USA</p>
    <p><strong>Ethnicity:</strong> Caucasian</p>
    <p><strong>Career Start:</strong> 2015</p>
    <p><strong>Career End:</strong> Active</p>
  </div>

  <div class="measurements">
    <p><strong>Height:</strong> 5'6" (168 cm)</p>
    <p><strong>Weight:</strong> 125 lbs (57 kg)</p>
    <p><strong>Measurements:</strong> 34D-24-36</p>
    <p><strong>Hair Color:</strong> Blonde</p>
    <p><strong>Eye Color:</strong> Blue</p>
    <p><strong>Tattoos:</strong> Lower back tribal design</p>
    <p><strong>Piercings:</strong> Navel</p>
  </div>

  <div class="filmography">
    <h2>Filmography</h2>
    <ul>
      <li>Scene Title 1 (2023)</li>
      <li>Scene Title 2 (2022)</li>
    </ul>
  </div>

  <div class="photos">
    <img src="/images/performer-headshot.jpg" alt="Performer Name">
  </div>
</div>
```

---

## Implementation Checklist

### Phase 1: Web Scraping (Adult Content - c schema)
- [ ] HTML scraping setup (`goquery` OR `colly`)
- [ ] User-Agent configuration (REQUIRED)
- [ ] URL normalization (performer name → Babepedia URL)
- [ ] Bio extraction (birthdate, ethnicity, career dates)
- [ ] Measurements extraction (height, weight, body measurements)
- [ ] Filmography extraction (scene list)
- [ ] Photo extraction (headshot, body shots)
- [ ] **c schema storage**: `c.performers.metadata_json.babepedia_data` (JSONB)

### Phase 2: Data Enrichment
- [ ] Alias resolution (match alternative performer names)
- [ ] Career timeline (debut → retirement)
- [ ] Tattoo/piercing tracking (detailed descriptions)
- [ ] Physical attribute updates (measurements change over time)

### Phase 3: Background Jobs (River)
- [ ] **Job**: `c.wiki.babepedia.scrape_performer` (scrape performer profile)
- [ ] **Job**: `c.wiki.babepedia.refresh` (periodic refresh)
- [ ] Rate limiting (conservative 1 req/sec)
- [ ] Retry logic (exponential backoff)

---

## Integration Pattern

### Performer Enrichment Flow
```
User views adult performer profile (e.g., c.performers)
        ↓
Check if Babepedia data exists in cache
        ↓
        NO
        ↓
Construct Babepedia URL (https://www.babepedia.com/babe/{name})
        ↓
Scrape HTML page (goquery)
        ↓
Parse performer data:
  - Bio (birthdate, ethnicity, career dates)
  - Measurements (height, weight, body measurements)
  - Tattoos/piercings (detailed descriptions)
  - Filmography (scene list)
  - Photos (headshot, body shots)
        ↓
Store in c.performers.metadata_json.babepedia_data (c schema JSONB)
        ↓
Display in UI (performer profile page, c schema isolated)
```

### Rate Limiting Strategy
```
Babepedia scraping: 1 req/sec (very conservative)
- Cache scraped data for 90 days (performer data changes infrequently)
- Background jobs: Use River queue (low priority)
- Respect server load
```

---

## Related Documentation

- [IAFD.md](./IAFD.md) - Internet Adult Film Database (alternative source)
- [BOOBPEDIA.md](./BOOBPEDIA.md) - Adult performer encyclopedia (alternative source)
- [STASHDB.md](../metadata/adult/STASHDB.md) - Primary adult metadata (performers, studios, scenes)
- [THEPORNDB.md](../metadata/adult/THEPORNDB.md) - Adult metadata provider
- [ADULT_METADATA.md](../../ADULT_METADATA.md) - Adult metadata system architecture

---

## Notes

### No Official API (Web Scraping)
- **Babepedia has NO API**: All data retrieval requires web scraping
- **Scraping risks**: HTML structure changes can break scraper
- **Maintenance**: Regular scraper updates required
- **Priority**: MEDIUM (useful for performer enrichment)

### Adult Content Isolation (CRITICAL)
- **Database schema**: `c` schema ONLY
  - `c.performers.metadata_json.babepedia_data` (JSONB)
  - NO data in public schema
- **API namespace**: `/api/v1/c/wiki/babepedia/*` (isolated)
  - `/api/v1/c/wiki/babepedia/search/{performer_name}`
  - `/api/v1/c/wiki/babepedia/performers/{performer_id}`
- **Module location**: `internal/content/c/wiki/babepedia/` (isolated)
- **Access control**: Mods/admins see all, regular users see only their library

### User-Agent Requirement
- **MUST set User-Agent**: Web scraping requires proper User-Agent
- **Format**: `Revenge/1.0 (https://github.com/lusoris/revenge; contact@example.com)`

### Rate Limits (Conservative)
- **No official limit**: Babepedia doesn't publish rate limits
- **Conservative approach**: 1 req/sec (very conservative)
- **Caching**: Cache scraped data for 90 days (performer data stable)

### Content Licensing
- **License**: Unclear (community-contributed content)
- **Attribution**: Recommended to attribute Babepedia in UI ("Data from Babepedia")
- **Link**: Include Babepedia profile URL

### URL Normalization
- **Name formatting**: Spaces → underscores (e.g., "Performer Name" → "Performer_Name")
- **Special characters**: Handle apostrophes, hyphens (e.g., "O'Connor" → "O'Connor")
- **Aliases**: Try performer name first, fallback to aliases

### HTML Parsing Challenges
- **Dynamic structure**: HTML structure can change
- **Inconsistent formatting**: Different performer pages have different formats
- **Missing data**: Not all performers have complete data
- **Fallback**: If scraping fails, skip Babepedia data (non-critical)

### JSONB Storage (c schema)
- Store Babepedia data in `c.performers.metadata_json.babepedia_data`
- Fields:
  - `url`: Babepedia profile URL
  - `real_name`: Real name (if available)
  - `aliases`: Array of aliases
  - `birthdate`: Date of birth
  - `birthplace`: Place of birth
  - `ethnicity`: Ethnicity
  - `career_start`: Career start year
  - `career_end`: Career end year (null if active)
  - `height_cm`: Height in centimeters
  - `weight_kg`: Weight in kilograms
  - `measurements`: Body measurements (e.g., "34D-24-36")
  - `hair_color`: Hair color
  - `eye_color`: Eye color
  - `tattoos`: Array of tattoo descriptions
  - `piercings`: Array of piercing descriptions
  - `filmography`: Array of scene titles
  - `photos`: Array of image URLs
  - `last_scraped`: Timestamp

### Caching Strategy
- **Cache duration**: 90 days (performer data changes infrequently)
- **Invalidation**: Manual refresh OR periodic background job (quarterly)
- **Storage**: Store in `c.performers.metadata_json.babepedia_data` (JSONB) + Dragonfly cache (c namespace)

### Use Case: Performer Enrichment
- **Primary source**: StashDB/ThePornDB for performer metadata
- **Supplement**: Babepedia for detailed bio, measurements, career timeline
- **Display**: "Performer Bio (Babepedia)" section in performer profile

### Measurements Tracking
- **Measurements change**: Body measurements can change over time
- **Historical tracking**: Store measurement history (timestamp + measurements)
- **Display**: Show current measurements + historical changes (if available)

### Career Timeline
- **Debut date**: First scene/movie release
- **Retirement**: Career end date (if retired)
- **Active status**: Detect if still active (recent scenes)

### Alias Resolution
- **Alternative names**: Performers often use multiple aliases
- **Matching**: Use aliases for searching/matching (fuzzy matching)
- **Normalization**: Normalize performer names (lowercase, remove special chars)

### Scraping Ethics
- **Respect server**: Use conservative rate limiting (1 req/sec)
- **Cache aggressively**: 90 day cache reduces scraping frequency
- **User-Agent**: Identify as scraper
- **Fallback**: If scraping fails, degrade gracefully

### Fallback Strategy (Adult Performer Metadata)
- **Order**: StashDB (primary) → ThePornDB (supplementary) → Babepedia (bio/measurements) → IAFD (filmography) → Boobpedia (alternative)
