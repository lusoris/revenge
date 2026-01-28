# Wikipedia Integration

> General encyclopedia information via Wikipedia API

**Service**: Wikipedia (https://www.wikipedia.org)
**API**: MediaWiki Action API (https://www.mediawiki.org/wiki/API:Main_page)
**Category**: Wiki / Knowledge Base
**Priority**: 🟢 MEDIUM (Supplementary info)
**Status**: 🔴 DESIGN PHASE

---

## Overview

**Wikipedia** is the world's largest free encyclopedia with comprehensive information about movies, TV shows, music, books, people, and more. The MediaWiki Action API provides programmatic access to Wikipedia content.

**Key Features**:
- **General information**: Plot summaries, production info, cast/crew, release history
- **Biography data**: Person biographies, career info, personal life
- **Historical context**: Background information, cultural impact, trivia
- **Multi-language**: 300+ language editions
- **Free access**: No API key required (rate-limited)
- **Structured data**: Infoboxes, categories, links

**Use Cases**:
- Plot summaries for movies/TV shows
- Artist/performer biographies
- Historical context for media
- Trivia and background information
- Links to related topics

---

## Developer Resources

### API Documentation
- **Base URL**: `https://en.wikipedia.org/w/api.php` (English Wikipedia)
- **API Docs**: https://www.mediawiki.org/wiki/API:Main_page
- **Query Examples**: https://www.mediawiki.org/wiki/API:Query
- **Rate Limits**: 200 requests/second (user-agent required)

### Authentication
- **Method**: None (public API)
- **User-Agent**: REQUIRED (`User-Agent: Revenge/1.0 (contact@example.com)`)
- **Rate Limits**: 200 req/sec (respect rate limits)

### Data Coverage
- **Articles**: 60M+ articles (all languages)
- **English**: 6.8M+ articles
- **Languages**: 300+ language editions
- **Updates**: Real-time (community-edited)

### Go Client Library
- **Official**: None
- **Recommended**: Use `net/http` with JSON parsing
- **Alternative**: `github.com/trietmn/go-wiki` (unofficial)

---

## API Details

### REST Endpoints

#### Search Pages
```
GET /w/api.php?action=query&list=search&srsearch={query}&format=json

Example:
https://en.wikipedia.org/w/api.php?action=query&list=search&srsearch=The%20Matrix&format=json

Response:
{
  "query": {
    "search": [
      {
        "pageid": 45678,
        "title": "The Matrix",
        "snippet": "The <b>Matrix</b> is a 1999 science fiction action film..."
      }
    ]
  }
}
```

#### Get Page Extract (Summary)
```
GET /w/api.php?action=query&prop=extracts&exintro=1&titles={title}&format=json

Example:
https://en.wikipedia.org/w/api.php?action=query&prop=extracts&exintro=1&titles=The%20Matrix&format=json

Response:
{
  "query": {
    "pages": {
      "45678": {
        "pageid": 45678,
        "title": "The Matrix",
        "extract": "The Matrix is a 1999 science fiction action film written and directed by the Wachowskis..."
      }
    }
  }
}
```

#### Get Full Page Content
```
GET /w/api.php?action=parse&page={title}&format=json

Response includes:
- Full HTML content
- Infobox data
- Categories
- Links
- Images
```

#### Get Page Images
```
GET /w/api.php?action=query&prop=pageimages&titles={title}&format=json

Response:
{
  "query": {
    "pages": {
      "45678": {
        "thumbnail": {
          "source": "https://upload.wikimedia.org/wikipedia/en/thumb/c/c1/The_Matrix_Poster.jpg/220px-The_Matrix_Poster.jpg",
          "width": 220,
          "height": 326
        }
      }
    }
  }
}
```

---

## Implementation Checklist

### Phase 1: Core Integration
- [ ] REST API client setup (Go `net/http`)
- [ ] User-Agent configuration (REQUIRED: `Revenge/1.0 (contact@example.com)`)
- [ ] Page search (MediaWiki `action=query&list=search`)
- [ ] Page extract fetch (summary text)
- [ ] Full page content fetch (HTML parsing)
- [ ] Image fetch (page thumbnails)
- [ ] JSONB storage (`metadata_json.wikipedia_data`)

### Phase 2: Content Enhancement
- [ ] Plot summary extraction (movies/TV shows)
- [ ] Biography extraction (actors/directors/musicians)
- [ ] Trivia extraction (parse sections)
- [ ] Production info extraction (parse infoboxes)
- [ ] Multi-language support (fallback to English)
- [ ] Link extraction (related topics)

### Phase 3: Background Jobs (River)
- [ ] **Job**: `wiki.wikipedia.fetch_summary` (fetch page summary)
- [ ] **Job**: `wiki.wikipedia.refresh` (periodic refresh for popular content)
- [ ] Rate limiting (200 req/sec max, use conservative limits)
- [ ] Retry logic (exponential backoff)

---

## Integration Pattern

### Content Enrichment Flow
```
User views movie/TV show/music page
        ↓
Check if Wikipedia summary exists in cache
        ↓
        NO
        ↓
Search Wikipedia API (action=query&list=search)
        ↓
Match found? → Get page extract (action=query&prop=extracts)
              ↓
              Extract summary text (intro paragraph)
              ↓
              Store in metadata_json.wikipedia_data
              ↓
              Display in UI (collapsible "Wikipedia" section)
        ↓
        NO MATCH
        ↓
Skip (no Wikipedia data available)
```

### Rate Limiting Strategy
```
Wikipedia rate limit: 200 req/sec (official limit, use conservative approach)
- Conservative limit: 10 req/sec (token bucket)
- Caching: Cache Wikipedia data for 30 days (reduce API calls)
- Background jobs: Use River queue (prioritize user-initiated requests)
- Exponential backoff: Retry with delay if 429 errors
```

---

## Related Documentation

- [FANDOM.md](./FANDOM.md) - Fan wikis (MCU, Memory Alpha, Wookieepedia)
- [TVTROPES.md](./TVTROPES.md) - Trope analysis
- [INTERNAL_WIKI.md](../../features/INTERNAL_WIKI.md) - Built-in wiki system

---

## Notes

### User-Agent Requirement (CRITICAL)
- **MUST set User-Agent**: Wikipedia requires User-Agent header
- **Format**: `Revenge/1.0 (https://github.com/lusoris/revenge; contact@example.com)`
- **Failure**: Requests without User-Agent will be blocked (HTTP 403)

### Rate Limits (200 req/sec)
- **Official limit**: 200 requests/second
- **Conservative approach**: Use 10 req/sec to avoid issues
- **Token bucket**: Implement client-side rate limiter
- **Respect limits**: Aggressive usage can result in IP ban

### Content Licensing
- **License**: Creative Commons Attribution-ShareAlike 4.0 (CC BY-SA 4.0)
- **Attribution**: MUST attribute Wikipedia in UI ("From Wikipedia, the free encyclopedia")
- **Share-alike**: Derived content must use same license
- **Citation**: Include Wikipedia page URL

### Extract vs Full Content
- **Extract** (`prop=extracts`): Plain text summary (intro paragraph)
  - Use case: Quick info display in UI
  - No images, no infoboxes
- **Full content** (`action=parse`): Complete HTML
  - Use case: Detailed wiki page
  - Includes infoboxes, images, tables
  - Requires HTML parsing

### Infobox Parsing
- **Infoboxes**: Structured data tables (e.g., movie release date, cast, budget)
- **Parsing**: Extract from HTML (`<table class="infobox">`)
- **Use case**: Enrich metadata (supplement TMDb/TheTVDB data)

### Multi-Language Support
- **Language editions**: 300+ languages (en, de, fr, es, ja, etc.)
- **URL format**: `https://{lang}.wikipedia.org/w/api.php` (e.g., `de.wikipedia.org`)
- **Fallback**: Try user's language → fallback to English
- **Detection**: Use `Accept-Language` header OR explicit language selection

### Disambiguation Pages
- **Disambiguation pages**: Multiple meanings for same term
- **Detection**: Check if page title contains "(disambiguation)"
- **Handling**: Parse disambiguation page → list options → user selects correct page

### Redirect Handling
- **Redirects**: Wikipedia uses redirects (e.g., "The Matrix 1" → "The Matrix")
- **Automatic**: MediaWiki API follows redirects by default
- **Detection**: Check `redirects` field in response

### JSONB Storage
- Store Wikipedia data in `metadata_json.wikipedia_data`
- Fields:
  - `page_id`: Wikipedia page ID
  - `title`: Page title
  - `extract`: Summary text (intro paragraph)
  - `url`: Wikipedia page URL
  - `thumbnail`: Image URL
  - `categories`: List of categories
  - `last_fetched`: Timestamp

### Caching Strategy
- **Cache duration**: 30 days (Wikipedia content changes infrequently)
- **Invalidation**: Manual refresh OR automatic on content update
- **Storage**: Store in `metadata_json` (JSONB) + Dragonfly cache (fast access)

### Use Case: Plot Summaries
- **Movies/TV shows**: Extract plot summary from Wikipedia
- **Supplement TMDb**: TMDb has short overviews, Wikipedia has detailed plots
- **Display**: Collapsible section "Wikipedia Plot Summary" in UI

### Use Case: Biographies
- **Actors/directors/musicians**: Fetch biography from Wikipedia
- **Supplement**: Enrich performer/artist pages with background info
- **Display**: "Biography (Wikipedia)" section in performer profile

### Search Accuracy
- **Exact matches**: Search for exact title (e.g., "The Matrix")
- **Disambiguation**: Handle multiple results (e.g., "The Matrix" vs "The Matrix Reloaded")
- **Year matching**: Prefer results matching release year (filter by year)

### Image Quality
- **Thumbnails**: Wikipedia provides low-res thumbnails (220px width)
- **Full images**: Fetch original images from Wikimedia Commons (higher resolution)
- **Use case**: Fallback images if TMDb/TheTVDB lacks posters

### API Response Caching
- **Cache responses**: Cache API responses in Dragonfly (reduce redundant calls)
- **TTL**: 30 days for page extracts, 7 days for search results
- **Invalidation**: On user-initiated refresh

### Content Quality
- **Reliability**: Wikipedia content is community-edited (variable quality)
- **Moderation**: Active moderation on popular pages (movies, TV shows)
- **Vandalism**: Rare on popular pages, more common on obscure pages
- **Use case**: Supplementary information (NOT primary metadata source)

### Fallback Strategy
- **Wikipedia supplementary**: Use Wikipedia for plot summaries, trivia, biographies
- **Primary sources**: TMDb/TheTVDB for core metadata (titles, release dates, cast)
- **Order**: TMDb/TheTVDB (primary) → Wikipedia (supplementary) → FANDOM/TVTropes (niche)
