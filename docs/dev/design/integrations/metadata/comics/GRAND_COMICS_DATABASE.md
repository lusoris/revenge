# Grand Comics Database (GCD) Integration

> Open-source historical comics database (Golden/Silver Age focus)

**Priority**: 🟢 LOW (Phase 7 - Comics Module, historical fallback)
**Provider**: Grand Comics Database Project

## Status

| Dimension | Status | Notes |
| --------- | ------ | ----- |
| Design | ✅ | Comprehensive database import strategy, SQL examples |
| Sources | ✅ | Website, database dumps, data license linked |
| Instructions | ✅ | Phased implementation checklist with import strategy |
| Code | 🔴 | |
| Linting | 🔴 | |
| Unit Testing | 🔴 | |
| Integration Testing | 🔴 | |

---

## Overview

**Grand Comics Database (GCD)** is an open-source, volunteer-driven database focused on historical comics, particularly Golden Age (1938-1956) and Silver Age (1956-1970) comics. It's the most comprehensive source for public domain and historical comics metadata.

**Why GCD** (as fallback):
- **Historical focus**: Best source for pre-1980 comics
- **Public domain comics**: Extensive coverage of public domain works
- **Detailed credits**: Comprehensive creator credits (writers, artists, editors)
- **Cover scans**: High-quality scans of historical covers
- **Free and open**: Open-source project, no API key required
- **International comics**: Strong coverage of non-US comics (European, British, etc.)

**Limitations**:
- **No official API**: Data available via database dumps, web scraping, or community tools
- **Modern comics**: ComicVine better for post-1980 comics
- **Update frequency**: Slower updates compared to commercial databases
- **Manual integration**: Requires custom parsing or third-party tools

**Use Case**: Fallback for historical comics (pre-1980), public domain comics, when ComicVine/Marvel API lack data.

---

## Developer Resources

### Data Access
- **Website**: https://www.comics.org/
- **Database Dumps**: https://docs.comics.org/wiki/Downloading_the_Data
- **API**: No official API (use database dumps or community tools)
- **Data License**: CC BY-SA 4.0 (Creative Commons Attribution-ShareAlike)
- **Format**: MySQL dumps (import into PostgreSQL with conversion)

### Key Features
- **Comprehensive historical data**: 1M+ issues, 300K+ series (1930s-present)
- **Detailed credits**: Writers, pencillers, inkers, colorists, letterers, editors
- **Cover gallery**: High-quality scans of comic covers (historical focus)
- **Story details**: Individual story credits within issues (not just issue-level)
- **Publisher information**: Extensive publisher/imprint data
- **Reprints**: Tracks reprint editions and collections

### Data Structure (Simplified)
```sql
-- GCD uses normalized relational database
-- Key tables: gcd_issue, gcd_series, gcd_story, gcd_creator, gcd_publisher

-- Example: gcd_issue table
CREATE TABLE gcd_issue (
    id INT PRIMARY KEY,
    series_id INT,
    number VARCHAR(50),
    publication_date VARCHAR(20),
    key_date VARCHAR(10),
    page_count DECIMAL(10,2),
    editing TEXT,
    notes TEXT
);

-- Example: gcd_story table (individual stories within issues)
CREATE TABLE gcd_story (
    id INT PRIMARY KEY,
    issue_id INT,
    title VARCHAR(255),
    feature VARCHAR(255),
    type_id INT,  -- cover, story, advertisement, etc.
    page_count DECIMAL(10,2),
    script TEXT,  -- writer credits
    pencils TEXT, -- penciller credits
    inks TEXT,
    colors TEXT,
    letters TEXT,
    editing TEXT
);
```

---

## Integration Approaches

### Approach 1: Database Import (Recommended for Historical Comics)
**Method**: Import GCD MySQL dump into Revenge PostgreSQL database

**Pros**:
- Full offline access to GCD data
- No rate limits or API dependencies
- Fast queries (local database)
- Complete historical coverage

**Cons**:
- Large database (~500MB-1GB)
- Requires periodic updates (manual or automated sync)
- Schema conversion (MySQL → PostgreSQL)

**Implementation Steps**:
1. Download GCD MySQL dump (https://docs.comics.org/wiki/Downloading_the_Data)
2. Convert MySQL schema to PostgreSQL (use `pgloader` or manual conversion)
3. Import into separate PostgreSQL schema (e.g., `gcd` schema)
4. Create indexes for common queries (series_id, issue_number, publication_date)
5. Implement GCD metadata service (query local GCD tables)
6. Schedule monthly/quarterly updates (download new dumps, merge changes)

**Example Query** (after import):
```sql
-- Find issue in GCD database
SELECT
    i.id,
    i.number,
    i.publication_date,
    s.name AS series_name,
    p.name AS publisher_name
FROM gcd.gcd_issue i
JOIN gcd.gcd_series s ON i.series_id = s.id
JOIN gcd.gcd_publisher p ON s.publisher_id = p.id
WHERE s.name ILIKE '%Action Comics%'
  AND i.number = '1';
```

### Approach 2: Web Scraping (NOT Recommended)
**Method**: Scrape comics.org website for metadata

**Pros**:
- No database import required
- Always up-to-date data

**Cons**:
- Fragile (website changes break scraper)
- Slow (HTTP requests for each lookup)
- Violates GCD's ToS (use database dumps instead)
- Rate limiting concerns

**Verdict**: Use Approach 1 (database import) instead.

### Approach 3: Community Tools (Optional)
**Method**: Use third-party GCD API wrappers or tools

**Options**:
- **pyGCD**: Python library for GCD data (https://github.com/comictagger/pyGCD)
- **ComicTagger**: Desktop app with GCD integration (open-source)

**Verdict**: Useful for reference, but database import provides more control.

---

## Implementation Checklist

### Phase 1: Database Import
- [ ] Download GCD MySQL dump (https://docs.comics.org/wiki/Downloading_the_Data)
- [ ] Convert MySQL schema to PostgreSQL (`pgloader` or manual)
- [ ] Create `gcd` schema in PostgreSQL
- [ ] Import GCD tables (gcd_issue, gcd_series, gcd_story, gcd_creator, gcd_publisher)
- [ ] Create indexes (series_id, issue_number, publication_date)
- [ ] Verify data integrity (sample queries, count checks)

### Phase 2: Metadata Service
- [ ] Implement GCD metadata service (query local GCD tables)
- [ ] Search GCD for historical comics (pre-1980 fallback)
- [ ] Fetch issue metadata (title, publication date, page count)
- [ ] Fetch story credits (writers, artists, editors)
- [ ] Map GCD IDs to Revenge comics table (metadata_json.gcd_id)
- [ ] Store GCD metadata in metadata_json JSONB field

### Phase 3: Cover Images
- [ ] Download GCD cover scans (https://www.comics.org/issue/{id}/cover/)
- [ ] Store covers in Revenge media storage (cache locally)
- [ ] Implement cover fallback (ComicVine → GCD covers)

### Phase 4: Sync & Maintenance
- [ ] Schedule GCD database updates (monthly/quarterly)
- [ ] Implement incremental update strategy (detect changes, merge)
- [ ] Add GCD attribution ("Data from Grand Comics Database (CC BY-SA 4.0)")
- [ ] Document GCD schema for maintainability

---

## Integration Pattern

### Historical Comics Fallback Flow
```
Comic scanned from library (CBZ file)
                ↓
Extract metadata (ComicInfo.xml, filename parsing)
                ↓
Search ComicVine API (primary)
                ↓
NO RESULTS (pre-1980 comic or obscure title)
                ↓
Fallback: Search local GCD database
                ↓
Query gcd_series + gcd_issue tables
                ↓
Fetch metadata (publication date, credits, page count)
                ↓
Download GCD cover scan (if available)
                ↓
Store GCD metadata in metadata_json.gcd_data
```

### Golden Age Comics Example
```sql
-- User adds "Action Comics #1" (1938, Superman's first appearance)
-- ComicVine may have incomplete data for 1930s comics

-- Step 1: Search ComicVine (primary attempt)
-- Result: Partial data OR no results

-- Step 2: Fallback to GCD database
SELECT
    i.id AS gcd_issue_id,
    i.number,
    i.publication_date,
    i.page_count,
    s.name AS series_name,
    p.name AS publisher_name
FROM gcd.gcd_issue i
JOIN gcd.gcd_series s ON i.series_id = s.id
JOIN gcd.gcd_publisher p ON s.publisher_id = p.id
WHERE s.name = 'Action Comics'
  AND i.number = '1'
  AND i.publication_date LIKE '1938%';

-- Step 3: Fetch story credits (Superman story)
SELECT
    st.title,
    st.script AS writer,
    st.pencils AS artist,
    st.inks AS inker
FROM gcd.gcd_story st
WHERE st.issue_id = {gcd_issue_id}
  AND st.type_id = 19  -- story type (not advertisement)
ORDER BY st.sequence_number;

-- Result: "Superman, Champion of the Oppressed"
-- Writer: Jerry Siegel, Artist: Joe Shuster
```

---


<!-- SOURCE-BREADCRUMBS-START -->

## Sources & Cross-References

> Auto-generated section linking to external documentation sources

### Cross-Reference Indexes

- [All Sources Index](../../../../sources/SOURCES_INDEX.md) - Complete list of external documentation
- [Design ↔ Sources Map](../../../../sources/DESIGN_CROSSREF.md) - Which docs reference which sources

<!-- SOURCE-BREADCRUMBS-END -->

<!-- DESIGN-BREADCRUMBS-START -->

## Related Design Docs

> Auto-generated cross-references to related design documentation

**Category**: [Comics](INDEX.md)

### In This Section

- [ComicVine API Integration](COMICVINE.md)
- [Marvel API Integration](MARVEL_API.md)

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

- **ComicVine Integration**: `docs/integrations/metadata/comics/COMICVINE.md` (primary metadata source)
- **Marvel API Integration**: `docs/integrations/metadata/comics/MARVEL_API.md` (Marvel enrichment)
- **Comics Module**: `docs/features/COMICS_MODULE.md` (overall comics architecture)

---

## Notes

- **Historical Focus**: GCD excels at Golden Age (1938-1956) and Silver Age (1956-1970) comics
- **Modern Comics**: ComicVine better for post-1980 comics (GCD still covers modern, but slower updates)
- **Public Domain**: Best source for public domain comics (pre-1964 US comics)
- **Database Size**: ~500MB-1GB (compressed), ~2-3GB (uncompressed)
- **Update Frequency**: GCD releases database dumps monthly/quarterly (check https://docs.comics.org/)
- **License**: CC BY-SA 4.0 (attribution + share-alike required)
  - Attribution: "Data from Grand Comics Database (https://www.comics.org/) - CC BY-SA 4.0"
- **Schema Complexity**: GCD uses highly normalized schema (many JOIN queries required)
- **Story-Level Credits**: GCD tracks credits per story (not just per issue), very detailed
- **No Official API**: Use database dumps only (web scraping violates ToS)
- **Incremental Updates**: Track `modified` timestamps in GCD tables, sync changes only
- **Cover Scans**: Available at `https://www.comics.org/issue/{gcd_issue_id}/cover/` (4 sizes: thumbnail, medium, large, original)
- **Reprint Tracking**: GCD tracks reprints (useful for collections, trade paperbacks)
- **International Comics**: Strong coverage of non-US comics (Franco-Belgian, British, etc.)
- **Fallback Strategy**: ComicVine (primary) → Marvel API (Marvel) → GCD (historical/obscure)
- **JSONB Storage**: Store GCD data in `metadata_json.gcd_data` (preserve full dataset)
- **PostgreSQL Import**: Use `pgloader` for MySQL → PostgreSQL conversion
  ```bash
  # Install pgloader
  brew install pgloader  # macOS
  apt-get install pgloader  # Ubuntu

  # Convert MySQL dump to PostgreSQL
  pgloader mysql://user:pass@localhost/gcd postgresql://user:pass@localhost/revenge_gcd
  ```
- **Use Case Priority**: Low priority (most users have modern comics), implement only if historical comics support needed
