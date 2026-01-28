# IAFD Integration

> Internet Adult Film Database - Adult performer and scene database

**Service**: IAFD (https://www.iafd.com)
**API**: None (web scraping required)
**Category**: Wiki / Database (Adult Content)
**Priority**: 🟢 MEDIUM (Filmography tracking)
**Status**: 🔴 DESIGN PHASE

---

## Overview

**IAFD** (Internet Adult Film Database) is the adult industry's equivalent of IMDb, providing comprehensive filmography data for performers, studios, and scenes. It's one of the oldest and most complete adult film databases on the internet.

**Key Features**:
- **Performer filmography**: Complete scene/movie lists
- **Studio database**: Production company information
- **Scene metadata**: Release dates, directors, studios
- **Cross-referencing**: Links between performers, scenes, studios
- **Historical data**: Covers adult films from 1960s to present
- **No official API**: Requires web scraping

**Use Cases**:
- Complete filmography tracking
- Scene release date verification
- Studio/production company info
- Performer career timeline
- Cross-referencing performers across scenes

**⚠️ CRITICAL: Adult Content Isolation**:
- **Database schema**: `c` schema ONLY (`c.performers`, `c.adult_movies`, `c.studios`)
- **API namespace**: `/api/v1/c/wiki/iafd/*` (NOT `/api/v1/wiki/iafd/*`)
- **Module location**: `internal/content/c/wiki/iafd/` (NOT `internal/service/wiki/`)
- **Access control**: Mods/admins can see all data for monitoring, regular users see only their own library

---

## Developer Resources

### API Documentation
- **API**: NONE (no official API)
- **Web Scraping**: Required (parse HTML)
- **Base URL**: https://www.iafd.com
  - Performers: `/person.rme/perfid={id}/gender={gender}/{name}`
  - Movies: `/title.rme/title={title}/year={year}/{title}-{year}.htm`
  - Studios: `/studio.rme/studioid={id}/{name}`
- **Rate Limits**: Undefined (use conservative scraping ~1 req/sec)

### Authentication
- **Method**: None (public website)
- **User-Agent**: REQUIRED (`User-Agent: Revenge/1.0 (contact@example.com)`)
- **Rate Limits**: Conservative 1 req/sec (respect server load)

### Data Coverage
- **Performers**: 200K+ performer profiles
- **Scenes/Movies**: 500K+ titles
- **Studios**: 10K+ production companies
- **Historical Coverage**: 1960s to present
- **Languages**: Primarily English
- **Updates**: Community-edited (frequent updates)

### Go Scraping Library
- **Recommended**: `github.com/PuerkitoBio/goquery` (jQuery-like HTML parsing)
- **Alternative**: `github.com/gocolly/colly` (web scraping framework)

---

## Integration Approach

### Web Scraping Strategy

**⚠️ CRITICAL: No Official API - Web Scraping Required**

#### Performer Page Structure
```
URL: https://www.iafd.com/person.rme/perfid=performer_id/gender=f/performer_name

HTML Structure (example):
<div id="perfbio">
  <h1>Performer Name</h1>

  <div class="biodata">
    <p><b>Born:</b> May 15, 1990</p>
    <p><b>Birthplace:</b> Los Angeles, CA, USA</p>
    <p><b>Years Active:</b> 2015-2023</p>
    <p><b>Aliases:</b> Alias1, Alias2</p>
  </div>

  <div class="perffilm">
    <table id="personal">
      <thead>
        <tr>
          <th>Title</th>
          <th>Year</th>
          <th>Studio</th>
          <th>Notes</th>
        </tr>
      </thead>
      <tbody>
        <tr>
          <td><a href="/title.rme/...">Scene Title 1</a></td>
          <td>2023</td>
          <td><a href="/studio.rme/...">Studio Name</a></td>
          <td>Facial</td>
        </tr>
      </tbody>
    </table>
  </div>
</div>
```

#### Movie/Scene Page Structure
```
URL: https://www.iafd.com/title.rme/title=movie_title/year=2023/movie-title-2023.htm

HTML Structure (example):
<div id="title-page">
  <h1>Movie Title (2023)</h1>

  <div class="movie-info">
    <p><b>Studio:</b> Studio Name</p>
    <p><b>Director:</b> Director Name</p>
    <p><b>Release Date:</b> January 15, 2023</p>
    <p><b>Runtime:</b> 45 minutes</p>
  </div>

  <div class="cast">
    <h2>Cast</h2>
    <ul>
      <li><a href="/person.rme/...">Performer 1</a></li>
      <li><a href="/person.rme/...">Performer 2</a></li>
    </ul>
  </div>
</div>
```

#### Studio Page Structure
```
URL: https://www.iafd.com/studio.rme/studioid=studio_id/studio-name

HTML Structure (example):
<div id="studio-page">
  <h1>Studio Name</h1>

  <div class="studio-info">
    <p><b>Parent Company:</b> Parent Studio Name</p>
    <p><b>Years Active:</b> 2010-Present</p>
    <p><b>Website:</b> https://studio-website.com</p>
  </div>

  <div class="studio-films">
    <table>
      <tr>
        <td><a href="/title.rme/...">Movie Title 1</a></td>
        <td>2023</td>
      </tr>
    </table>
  </div>
</div>
```

---

## Implementation Checklist

### Phase 1: Web Scraping (Adult Content - c schema)
- [ ] HTML scraping setup (`goquery` OR `colly`)
- [ ] User-Agent configuration (REQUIRED)
- [ ] URL construction (performer/movie/studio)
- [ ] Performer bio extraction (birthdate, years active, aliases)
- [ ] Filmography extraction (scene list, release dates)
- [ ] Scene metadata extraction (studio, director, cast)
- [ ] Studio info extraction (parent company, website)
- [ ] **c schema storage**: `c.performers.metadata_json.iafd_data`, `c.adult_movies.metadata_json.iafd_data`, `c.studios.metadata_json.iafd_data` (JSONB)

### Phase 2: Data Enrichment
- [ ] Filmography completion (fill gaps in performer filmographies)
- [ ] Scene cross-referencing (link performers ↔ scenes ↔ studios)
- [ ] Release date verification (validate scene release dates)
- [ ] Studio hierarchy (parent company relationships)

### Phase 3: Background Jobs (River)
- [ ] **Job**: `c.wiki.iafd.scrape_performer` (scrape performer filmography)
- [ ] **Job**: `c.wiki.iafd.scrape_scene` (scrape scene metadata)
- [ ] **Job**: `c.wiki.iafd.scrape_studio` (scrape studio info)
- [ ] **Job**: `c.wiki.iafd.refresh` (periodic refresh)
- [ ] Rate limiting (conservative 1 req/sec)
- [ ] Retry logic (exponential backoff)

---

## Integration Pattern

### Performer Filmography Enrichment Flow
```
User views adult performer profile (c.performers)
        ↓
Check if IAFD data exists in cache
        ↓
        NO
        ↓
Construct IAFD URL (https://www.iafd.com/person.rme/perfid={id}/gender={gender}/{name})
        ↓
Scrape HTML page (goquery)
        ↓
Parse performer filmography:
  - Scene titles
  - Release years
  - Studios
  - Notes (scene categories)
        ↓
Store in c.performers.metadata_json.iafd_data (c schema JSONB)
        ↓
Display in UI (performer filmography section, c schema isolated)
```

### Scene Metadata Enrichment Flow
```
User views adult scene (c.adult_movies)
        ↓
Check if IAFD data exists in cache
        ↓
        NO
        ↓
Construct IAFD URL (https://www.iafd.com/title.rme/title={title}/year={year}/{title}-{year}.htm)
        ↓
Scrape HTML page (goquery)
        ↓
Parse scene metadata:
  - Studio
  - Director
  - Release date
  - Cast (performers)
        ↓
Store in c.adult_movies.metadata_json.iafd_data (c schema JSONB)
        ↓
Display in UI (scene details page, c schema isolated)
```

### Rate Limiting Strategy
```
IAFD scraping: 1 req/sec (very conservative)
- Cache scraped data for 90 days (filmography changes infrequently)
- Background jobs: Use River queue (low priority)
- Respect server load
```

---

## Related Documentation

- [BABEPEDIA.md](./BABEPEDIA.md) - Adult performer wiki (alternative source)
- [BOOBPEDIA.md](./BOOBPEDIA.md) - Adult performer encyclopedia (alternative source)
- [STASHDB.md](../metadata/adult/STASHDB.md) - Primary adult metadata (performers, studios, scenes)
- [THEPORNDB.md](../metadata/adult/THEPORNDB.md) - Adult metadata provider
- [ADULT_METADATA.md](../../ADULT_METADATA.md) - Adult metadata system architecture

---

## Notes

### No Official API (Web Scraping)
- **IAFD has NO API**: All data retrieval requires web scraping
- **Scraping risks**: HTML structure changes can break scraper
- **Maintenance**: Regular scraper updates required
- **Priority**: MEDIUM (useful for filmography completion)

### Adult Content Isolation (CRITICAL)
- **Database schema**: `c` schema ONLY
  - `c.performers.metadata_json.iafd_data` (JSONB)
  - `c.adult_movies.metadata_json.iafd_data` (JSONB)
  - `c.studios.metadata_json.iafd_data` (JSONB)
  - NO data in public schema
- **API namespace**: `/api/v1/c/wiki/iafd/*` (isolated)
  - `/api/v1/c/wiki/iafd/search/{performer_name}`
  - `/api/v1/c/wiki/iafd/performers/{performer_id}`
  - `/api/v1/c/wiki/iafd/movies/{movie_id}`
  - `/api/v1/c/wiki/iafd/studios/{studio_id}`
- **Module location**: `internal/content/c/wiki/iafd/` (isolated)
- **Access control**: Mods/admins see all, regular users see only their library

### User-Agent Requirement
- **MUST set User-Agent**: Web scraping requires proper User-Agent
- **Format**: `Revenge/1.0 (https://github.com/lusoris/revenge; contact@example.com)`

### Rate Limits (Conservative)
- **No official limit**: IAFD doesn't publish rate limits
- **Conservative approach**: 1 req/sec (very conservative)
- **Caching**: Cache scraped data for 90 days (filmography stable)

### Content Licensing
- **License**: Unclear (community-contributed content)
- **Attribution**: Recommended to attribute IAFD in UI ("Data from IAFD")
- **Link**: Include IAFD profile URL

### URL Construction
- **Performer URL**: `/person.rme/perfid={id}/gender={m|f}/performer-name`
- **Movie URL**: `/title.rme/title={title}/year={year}/{title}-{year}.htm`
- **Studio URL**: `/studio.rme/studioid={id}/studio-name`
- **URL normalization**: Spaces → hyphens, lowercase, remove special chars

### HTML Parsing Challenges
- **Dynamic structure**: HTML structure can change
- **Inconsistent formatting**: Different pages have different formats
- **Missing data**: Not all pages have complete data
- **Fallback**: If scraping fails, skip IAFD data (non-critical)

### JSONB Storage (c schema)

#### Performer Data (`c.performers.metadata_json.iafd_data`)
```json
{
  "url": "https://www.iafd.com/person.rme/perfid=1234/gender=f/performer-name",
  "performer_id": "1234",
  "birthdate": "1990-05-15",
  "birthplace": "Los Angeles, CA, USA",
  "years_active": "2015-2023",
  "aliases": ["Alias1", "Alias2"],
  "filmography": [
    {
      "title": "Scene Title 1",
      "year": 2023,
      "studio": "Studio Name",
      "notes": "Facial"
    }
  ],
  "last_scraped": "2023-01-15T10:00:00Z"
}
```

#### Scene Data (`c.adult_movies.metadata_json.iafd_data`)
```json
{
  "url": "https://www.iafd.com/title.rme/title=movie-title/year=2023/movie-title-2023.htm",
  "title": "Movie Title",
  "year": 2023,
  "studio": "Studio Name",
  "director": "Director Name",
  "release_date": "2023-01-15",
  "runtime_minutes": 45,
  "cast": ["Performer 1", "Performer 2"],
  "last_scraped": "2023-01-15T10:00:00Z"
}
```

#### Studio Data (`c.studios.metadata_json.iafd_data`)
```json
{
  "url": "https://www.iafd.com/studio.rme/studioid=1234/studio-name",
  "studio_id": "1234",
  "parent_company": "Parent Studio Name",
  "years_active": "2010-Present",
  "website": "https://studio-website.com",
  "last_scraped": "2023-01-15T10:00:00Z"
}
```

### Caching Strategy
- **Cache duration**: 90 days (filmography changes infrequently)
- **Invalidation**: Manual refresh OR periodic background job (quarterly)
- **Storage**: Store in JSONB (`c.performers/adult_movies/studios.metadata_json.iafd_data`) + Dragonfly cache (c namespace)

### Use Case: Filmography Completion
- **Primary source**: StashDB/ThePornDB for performer metadata
- **Supplement**: IAFD for complete filmography (historical + recent)
- **Display**: "Filmography (IAFD)" section in performer profile

### Filmography Tracking
- **Historical data**: IAFD has extensive historical coverage (1960s+)
- **Recent scenes**: Community-edited (frequent updates)
- **Cross-referencing**: Link scenes ↔ performers ↔ studios

### Release Date Verification
- **IAFD release dates**: Often more accurate than other sources
- **Validation**: Use IAFD release dates to validate other sources
- **Display**: Show release date with source attribution

### Studio Hierarchy
- **Parent companies**: Many studios are owned by parent companies (e.g., Brazzers → MindGeek)
- **IAFD tracks**: Parent company relationships
- **Schema**: `c.studios.parent_studio_id` (foreign key to `c.studios.id`)

### Alias Resolution
- **Alternative names**: Performers often use multiple aliases
- **IAFD tracks**: Comprehensive alias lists
- **Matching**: Use aliases for searching/matching (fuzzy matching)

### Scraping Ethics
- **Respect server**: Use conservative rate limiting (1 req/sec)
- **Cache aggressively**: 90 day cache reduces scraping frequency
- **User-Agent**: Identify as scraper
- **Fallback**: If scraping fails, degrade gracefully

### Robots.txt Compliance
- **Check**: https://www.iafd.com/robots.txt
- **Respect directives**: Follow User-agent, Disallow, Crawl-delay rules
- **Legal**: Not legally binding but ethical to respect

### Maintenance Burden
- **Web scraping**: Ongoing maintenance required
- **HTML changes**: IAFD can change HTML structure anytime
- **Breaking changes**: Scraper breaks without notice
- **Monitoring**: Implement health checks to detect failures

### Fallback Strategy (Adult Filmography)
- **Order**: StashDB (primary) → ThePornDB (supplementary) → IAFD (complete filmography) → Babepedia (bio) → Boobpedia (alternative)
