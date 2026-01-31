# TheNude Integration

> Adult performer database with aliases and measurements

## Status

| Dimension | Status | Notes |
| --------- | ------ | ----- |
| Design | ✅ | Web scraping spec, alias resolution flow, JSONB schema |
| Sources | ✅ | Website URL, goquery library documented |
| Instructions | ✅ | Phased implementation checklist with alias matching |
| Code | 🔴 | |
| Linting | 🔴 | |
| Unit Testing | 🔴 | |
| Integration Testing | 🔴 | |

---

## Overview

**TheNude** is an adult performer database focusing on aliases, measurements, and cross-referencing performers across different names and studios. It's particularly useful for alias resolution and finding performers under multiple stage names.

**Key Features**:
- **Alias tracking**: Comprehensive alias/stage name tracking
- **Measurements**: Physical attributes, tattoos, piercings
- **Cross-referencing**: Links to other databases (FreeOnes, IAFD, Babepedia)
- **Bio info**: Career dates, birthdate, ethnicity
- **Photos**: Performer photos (limited compared to FreeOnes)

**Use Cases**:
- Alias resolution (match performer across different names)
- Measurements enrichment
- Cross-reference with other databases
- Performer bio supplementation

**⚠️ CRITICAL: Adult Content Isolation**:
- **Database schema**: `c` schema ONLY (`c.performers`)
- **API namespace**: `/api/v1/legacy/external/thenude/*` (NOT `/api/v1/external/thenude/*`)
- **Module location**: `internal/content/c/external/thenude/` (NOT `internal/service/external/`)
- **Access control**: Mods/admins can see all data for monitoring, regular users see only their own library

---

## Developer Resources

### API Status
- **Official API**: NONE (no official API)
- **Web Scraping**: Required (parse HTML)
- **Base URL**: https://www.thenude.com/{performer_id}.htm

### Authentication
- **Method**: None (public website)
- **User-Agent**: REQUIRED (`User-Agent: Revenge/1.0 (contact@example.com)`)
- **Rate Limits**: Very conservative (1 req/sec)

### Data Coverage
- **Performers**: 50K+ performer profiles
- **Specialty**: Alias tracking, cross-referencing
- **Quality**: Good (focused database)

### Go Scraping Library
- **Recommended**: `github.com/PuerkitoBio/goquery`

---

## Integration Approach

### Web Scraping Strategy

**⚠️ CRITICAL: No Official API - Web Scraping Required**

#### Performer Page Structure
```
URL: https://www.thenude.com/Performer_Name_12345.htm

HTML Structure (example):
<div class="performer-profile">
  <h1>Performer Name</h1>

  <div class="aliases">
    <h2>Also Known As</h2>
    <ul>
      <li>Alias 1</li>
      <li>Alias 2</li>
      <li>Alias 3</li>
    </ul>
  </div>

  <div class="bio">
    <p><strong>Born:</strong> May 15, 1990</p>
    <p><strong>Birthplace:</strong> Los Angeles, CA</p>
    <p><strong>Ethnicity:</strong> Caucasian</p>
    <p><strong>Career:</strong> 2015-2023</p>
  </div>

  <div class="measurements">
    <p><strong>Height:</strong> 168 cm</p>
    <p><strong>Weight:</strong> 57 kg</p>
    <p><strong>Measurements:</strong> 34D-24-36</p>
  </div>

  <div class="external-links">
    <h2>External Links</h2>
    <ul>
      <li><a href="https://www.freeones.com/...">FreeOnes</a></li>
      <li><a href="https://www.iafd.com/...">IAFD</a></li>
      <li><a href="https://www.babepedia.com/...">Babepedia</a></li>
    </ul>
  </div>
</div>
```

---

## Implementation Checklist

### Phase 1: Web Scraping (Adult Content - c schema)
- [ ] HTML scraping setup (`goquery`)
- [ ] User-Agent configuration (REQUIRED)
- [ ] URL construction (performer ID → TheNude URL)
- [ ] Alias extraction (comprehensive alias list)
- [ ] Bio extraction (birthdate, ethnicity, career dates)
- [ ] Measurements extraction (height, weight, body measurements)
- [ ] External links extraction (FreeOnes, IAFD, Babepedia cross-references)
- [ ] **c schema storage**: `c.performers.metadata_json.thenude_data` (JSONB)

### Phase 2: Alias Resolution
- [ ] Alias matching (fuzzy matching across aliases)
- [ ] Cross-database linking (match TheNude → FreeOnes/IAFD)
- [ ] Performer deduplication (identify same performer under different names)

### Phase 3: Background Jobs (River)
- [ ] **Job**: `c.external.thenude.scrape_performer` (scrape performer profile)
- [ ] **Job**: `c.external.thenude.refresh` (periodic refresh)
- [ ] Rate limiting (very conservative 1 req/sec)
- [ ] Retry logic (exponential backoff)

---

## Integration Pattern

### Alias Resolution Flow
```
Performer name mismatch detected (e.g., "Jane Doe" in Revenge vs "J. Doe" in StashDB)
        ↓
Search TheNude for performer aliases
        ↓
Scrape TheNude profile (https://www.thenude.com/{performer_id}.htm)
        ↓
Parse aliases:
  - "Jane Doe"
  - "J. Doe"
  - "Janie D"
        ↓
Store aliases in c.performers.metadata_json.thenude_data.aliases
        ↓
Use aliases for performer matching across databases
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

- [FreeOnes Integration](FREEONES.md)
- [Instagram Integration](INSTAGRAM.md)
- [OnlyFans Integration](ONLYFANS.md)
- [Pornhub Integration](PORNHUB.md)
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

- [FREEONES.md](./FREEONES.md) - FreeOnes performer database (primary)
- [STASHDB.md](../metadata/adult/STASHDB.md) - Primary adult metadata
- [BABEPEDIA.md](../wiki/adult/BABEPEDIA.md) - Adult performer wiki
- [IAFD.md](../wiki/adult/IAFD.md) - Internet Adult Film Database

---

## Notes

### No Official API (Web Scraping Only)
- **TheNude has NO API**: All data retrieval requires web scraping
- **Scraping risks**: HTML structure changes can break scraper
- **Priority**: LOW (FreeOnes is primary, TheNude is supplementary for aliases)

### Adult Content Isolation (CRITICAL)
- **Database schema**: `c` schema ONLY
  - `c.performers.metadata_json.thenude_data` (JSONB)
  - NO data in public schema
- **API namespace**: `/api/v1/legacy/external/thenude/*` (isolated)
- **Module location**: `internal/content/c/external/thenude/` (isolated)
- **Access control**: Mods/admins see all, regular users see only their library

### Specialty: Alias Tracking
- **TheNude excels**: Comprehensive alias tracking (performers under 5-10+ different names)
- **Use case**: Resolve performer aliases across different studios/scenes
- **Example**: "Jane Doe" = "J. Doe" = "Janie D" = "Jane D." = all same performer

### Cross-Database Linking
- **TheNude links**: FreeOnes, IAFD, Babepedia, Boobpedia profiles
- **Use case**: Cross-reference performer across databases
- **Validation**: Use TheNude links to validate performer matches

### JSONB Storage (c schema)
```json
{
  "thenude_url": "https://www.thenude.com/Performer_Name_12345.htm",
  "thenude_id": "12345",
  "aliases": ["Alias 1", "Alias 2", "Alias 3"],
  "birthdate": "1990-05-15",
  "birthplace": "Los Angeles, CA",
  "ethnicity": "Caucasian",
  "career_start": 2015,
  "career_end": 2023,
  "height_cm": 168,
  "weight_kg": 57,
  "measurements": "34D-24-36",
  "external_links": {
    "freeones": "https://www.freeones.com/...",
    "iafd": "https://www.iafd.com/...",
    "babepedia": "https://www.babepedia.com/..."
  },
  "last_scraped": "2023-01-15T10:00:00Z"
}
```

### Caching Strategy
- **Cache duration**: 90 days (alias data stable)
- **Use case**: One-time lookup for alias resolution

### Fallback Strategy (Adult Performer Aliases)
- **Order**: StashDB (primary aliases) → TheNude (comprehensive alias tracking) → IAFD (alternative) → Babepedia (alternative)
