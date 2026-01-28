# TVTropes Integration

> Trope analysis and storytelling patterns

**Service**: TVTropes (https://tvtropes.org)
**API**: None (web scraping required)
**Category**: Wiki / Trope Analysis
**Priority**: 🟡 LOW (Niche content for enthusiasts)
**Status**: 🔴 DESIGN PHASE

---

## Overview

**TVTropes** is a wiki dedicated to cataloging storytelling tropes, patterns, and conventions across all media (movies, TV, books, games, etc.). It analyzes narrative devices, character archetypes, plot structures, and cultural references.

**Key Features**:
- **Trope catalog**: 30,000+ documented tropes (narrative patterns)
- **Analysis**: Detailed analysis of storytelling conventions
- **Cross-media**: Covers movies, TV, books, games, anime, comics
- **Community-driven**: Fan-curated content
- **No official API**: Requires web scraping

**Use Cases**:
- Trope analysis for movies/TV shows
- Storytelling pattern recognition
- Character archetype identification
- Cultural references and allusions
- Fan engagement (niche audience)

**Example Tropes**:
- "The Hero's Journey" - Classic narrative structure
- "Chekhov's Gun" - Foreshadowing device
- "MacGuffin" - Plot device object
- "Red Herring" - Misleading clue
- "Deus Ex Machina" - Contrived plot resolution

---

## Developer Resources

### API Documentation
- **API**: NONE (no official API)
- **Web Scraping**: Required (parse HTML)
- **Base URL**: https://tvtropes.org/pmwiki/pmwiki.php/Main/{TropeName}
- **Rate Limits**: Undefined (use conservative scraping ~1 req/sec)

### Authentication
- **Method**: None (public website)
- **User-Agent**: REQUIRED (`User-Agent: Revenge/1.0 (contact@example.com)`)
- **Rate Limits**: Conservative 1 req/sec (avoid overwhelming server)

### Data Coverage
- **Tropes**: 30,000+ documented tropes
- **Works**: Hundreds of thousands of analyzed works
- **Languages**: Primarily English
- **Updates**: Real-time (community-edited)

### Go Scraping Library
- **Recommended**: `github.com/PuerkitoBio/goquery` (jQuery-like HTML parsing)
- **Alternative**: `github.com/gocolly/colly` (web scraping framework)

---

## Integration Approach

### Web Scraping Strategy

**⚠️ CRITICAL: No Official API - Web Scraping Required**

#### Movie/TV Show Page Structure
```
URL: https://tvtropes.org/pmwiki/pmwiki.php/Film/TheMatrix

HTML Structure:
<div id="main-article">
  <h1>The Matrix</h1>
  <p>The Matrix is a 1999 science fiction action film...</p>

  <h2>Tropes</h2>
  <ul>
    <li>
      <a href="/pmwiki/pmwiki.php/Main/TheHerosJourney">The Hero's Journey</a>
      - Neo follows the classic monomyth structure...
    </li>
    <li>
      <a href="/pmwiki/pmwiki.php/Main/ChekhovsGun">Chekhov's Gun</a>
      - The red pill introduced early, becomes crucial later...
    </li>
  </ul>
</div>
```

#### Trope Page Structure
```
URL: https://tvtropes.org/pmwiki/pmwiki.php/Main/TheHerosJourney

HTML Structure:
<div id="main-article">
  <h1>The Hero's Journey</h1>
  <p>Also known as the Monomyth, this narrative pattern...</p>

  <h2>Examples</h2>
  <h3>Film</h3>
  <ul>
    <li><i>The Matrix</i>: Neo's journey from ordinary hacker to The One...</li>
    <li><i>Star Wars</i>: Luke Skywalker's transformation...</li>
  </ul>
</div>
```

---

## Implementation Checklist

### Phase 1: Web Scraping (Low Priority)
- [ ] HTML scraping setup (`goquery` OR `colly`)
- [ ] User-Agent configuration (REQUIRED)
- [ ] URL pattern mapping (Film/TheMatrix, Series/BreakingBad, etc.)
- [ ] Trope list extraction (parse HTML `<ul>` lists)
- [ ] Trope description extraction
- [ ] JSONB storage (`metadata_json.tvtropes_data`)

### Phase 2: Content Parsing
- [ ] Parse trope names (extract from links)
- [ ] Parse trope descriptions (extract text)
- [ ] Category detection (Film, Series, Western Animation, Anime, etc.)
- [ ] Related tropes (cross-references)

### Phase 3: Background Jobs (River) - Low Priority
- [ ] **Job**: `wiki.tvtropes.scrape_tropes` (scrape trope list)
- [ ] **Job**: `wiki.tvtropes.refresh` (periodic refresh)
- [ ] Rate limiting (conservative 1 req/sec, avoid server overload)
- [ ] Retry logic (exponential backoff, respect 429/503 errors)

---

## Integration Pattern

### Trope Scraping Flow
```
User views movie/TV show page (e.g., "The Matrix")
        ↓
Check if TVTropes data exists in cache
        ↓
        NO
        ↓
Construct TVTropes URL (https://tvtropes.org/pmwiki/pmwiki.php/Film/TheMatrix)
        ↓
Scrape HTML page (goquery)
        ↓
Parse trope list:
  - Extract trope names (from <a> links)
  - Extract trope descriptions (from <li> text)
        ↓
Store in metadata_json.tvtropes_data
        ↓
Display in UI (collapsible "TVTropes" section)
        ↓
User clicks trope? → Navigate to trope detail page (future feature)
```

### URL Mapping
```
Movies: /pmwiki/pmwiki.php/Film/{Title}
TV Shows: /pmwiki/pmwiki.php/Series/{Title}
Anime: /pmwiki/pmwiki.php/Anime/{Title}
Western Animation: /pmwiki/pmwiki.php/WesternAnimation/{Title}
```

### Rate Limiting Strategy
```
TVTropes scraping: 1 req/sec (very conservative)
- Avoid overwhelming server
- Cache scraped data for 30 days (reduce scraping frequency)
- Background jobs: Use River queue (low priority)
- Respect robots.txt (check if scraping allowed)
```

---

## Related Documentation

- [WIKIPEDIA.md](./WIKIPEDIA.md) - General encyclopedia
- [FANDOM.md](./FANDOM.md) - Fan wikis
- [INTERNAL_WIKI.md](../../features/INTERNAL_WIKI.md) - Built-in wiki system

---

## Notes

### No Official API (Web Scraping)
- **TVTropes has NO API**: All data retrieval requires web scraping
- **Scraping risks**: HTML structure changes can break scraper
- **Maintenance**: Regular scraper updates required
- **Priority**: LOW (implement only if time permits)

### User-Agent Requirement
- **MUST set User-Agent**: Web scraping requires proper User-Agent
- **Format**: `Revenge/1.0 (https://github.com/lusoris/revenge; contact@example.com)`
- **Politeness**: Identify as scraper to avoid being blocked

### Rate Limits (Conservative)
- **No official limit**: TVTropes doesn't publish rate limits
- **Conservative approach**: 1 req/sec (very conservative, respect server)
- **Caching**: Cache scraped data for 30+ days (reduce scraping frequency)
- **robots.txt**: Check `https://tvtropes.org/robots.txt` (respect directives)

### Content Licensing
- **License**: Creative Commons Attribution-NonCommercial-ShareAlike 3.0 Unported (CC BY-NC-SA 3.0)
- **Attribution**: MUST attribute TVTropes in UI ("Tropes from TVTropes.org")
- **Non-commercial**: Revenge is FOSS (non-commercial OK)
- **Link**: Include TVTropes page URL

### URL Normalization
- **Title normalization**: TVTropes uses CamelCase (e.g., "TheMatrix", not "The Matrix")
- **Spaces removed**: Remove spaces from titles (e.g., "Breaking Bad" → "BreakingBad")
- **Special characters**: Handle special characters (e.g., ":" → empty)

### HTML Parsing Challenges
- **Dynamic structure**: HTML structure can change without notice
- **Inconsistent formatting**: Different pages have different structures
- **Maintenance burden**: Scraper requires regular maintenance
- **Fallback**: If scraping fails, skip TVTropes data (non-critical)

### JSONB Storage
- Store TVTropes data in `metadata_json.tvtropes_data`
- Fields:
  - `url`: TVTropes page URL
  - `tropes`: Array of trope objects [{name, description, url}]
  - `last_scraped`: Timestamp
  - `scraper_version`: Scraper version (for maintenance tracking)

### Caching Strategy
- **Cache duration**: 30+ days (TVTropes content changes infrequently)
- **Invalidation**: Manual refresh OR periodic background job
- **Storage**: Store in `metadata_json` (JSONB) + Dragonfly cache

### Use Case: Niche Audience
- **Target audience**: Enthusiasts, film students, narrative analysis fans
- **Priority**: LOW (most users won't use TVTropes data)
- **Implementation**: Low priority (implement after core features complete)

### Trope Display UI
- **Collapsible section**: "Tropes (TVTropes)" in movie/TV show page
- **List format**: Display trope names with descriptions
- **Links**: Link to TVTropes pages (external links)
- **Spoiler warning**: TVTropes often contains spoilers (display warning)

### Scraping Ethics
- **Respect server**: Use conservative rate limiting (1 req/sec)
- **Cache aggressively**: Reduce scraping frequency (30+ day cache)
- **User-Agent**: Identify as scraper (proper User-Agent)
- **robots.txt**: Respect robots.txt directives
- **Fallback**: If scraping fails, degrade gracefully (skip TVTropes data)

### robots.txt Compliance
- **Check robots.txt**: `https://tvtropes.org/robots.txt`
- **Respect directives**: Follow robots.txt rules (User-agent, Disallow, Crawl-delay)
- **Legal**: robots.txt is not legally binding, but ethical to respect

### Maintenance Burden
- **High maintenance**: Web scraping requires ongoing maintenance
- **HTML changes**: TVTropes can change HTML structure anytime
- **Breaking changes**: Scraper can break without notice
- **Monitoring**: Implement scraper health checks (detect failures)

### Fallback Strategy
- **Order**: TMDb/TheTVDB (primary) → FANDOM (franchise-specific) → Wikipedia (general) → TVTropes (analysis, LOW priority)
- **Optional**: TVTropes is entirely optional (non-critical feature)
