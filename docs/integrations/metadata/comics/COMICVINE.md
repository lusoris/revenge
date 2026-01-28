# ComicVine API Integration

> Primary comics metadata provider (GameSpot's comprehensive comics database)

**Status**: 🟡 PLANNED
**Priority**: 🟡 MEDIUM (Phase 7 - Comics Module)
**Provider**: ComicVine (GameSpot)

---

## Overview

**ComicVine** is the primary metadata source for the Comics module, providing comprehensive data for Western comics (Marvel, DC, Image, etc.), graphic novels, and manga. Owned by GameSpot/Fandom, it offers one of the most complete comics databases available.

**Why ComicVine**:
- Largest Western comics database (1M+ issues, 100K+ series)
- Publisher-agnostic (Marvel, DC, Image, Dark Horse, etc.)
- Comprehensive metadata (writers, artists, colorists, inkers, letterers)
- Cover images (high-resolution available)
- Character/team/location/concept relationships
- Story arcs and crossover events
- Free API with generous rate limits
- Used by Mylar3, ComicTagger, and other comics tools

**Alternatives**:
- **Marvel API**: Official Marvel Comics API (Marvel content only, limited to 3000 req/day)
- **Grand Comics Database (GCD)**: Open-source database (historical focus, no official API)
- **AniList**: Manga metadata (Japanese comics, better suited for manga module)

**Use Case**: Primary metadata source for all Western comics and graphic novels.

---

## Developer Resources

### API Documentation
- **Base URL**: `https://comicvine.gamespot.com/api/`
- **Documentation**: https://comicvine.gamespot.com/api/documentation
- **API Version**: v1 (stable)
- **Authentication**: API Key (request via GameSpot account)
- **Rate Limits**: 200 requests per resource per hour (resets every hour)
- **Response Format**: JSON, XML (JSON recommended)

### Key Features
- **Search**: Search for issues, volumes, characters, creators, publishers
- **Detailed metadata**: Issue details (writers, artists, cover date, page count, description)
- **Relationships**: Characters appearing in issues, story arcs, crossover events
- **Cover images**: Multiple sizes (small, medium, super, original)
- **Publisher data**: Publisher information, imprints
- **Creator credits**: Writers, pencillers, inkers, colorists, letterers, cover artists

### Authentication
```http
GET /api/issues/?api_key=YOUR_API_KEY&format=json
```

**API Key**: Required for all requests (obtain from https://comicvine.gamespot.com/api/)

---

## API Details

### Authentication
**Method**: API Key in query parameter
**Header**: None required (User-Agent recommended for identification)
**Rate Limit**: 200 requests per resource per hour

**Example Request**:
```bash
curl "https://comicvine.gamespot.com/api/issue/4000-140529/?api_key=YOUR_API_KEY&format=json&field_list=id,name,volume,issue_number,cover_date,image"
```

### Core Endpoints

#### 1. Search Issues
```http
GET /api/search/?api_key={key}&format=json&query={query}&resources=issue
```

**Parameters**:
- `query`: Search term (e.g., "Amazing Spider-Man")
- `resources`: `issue` (search issues only)
- `field_list`: Comma-separated list of fields to return (optional, reduces response size)
- `limit`: Results per page (max 100, default 10)
- `offset`: Pagination offset

**Response**:
```json
{
  "status_code": 1,
  "number_of_page_results": 10,
  "number_of_total_results": 800,
  "results": [
    {
      "id": 140529,
      "name": "The Night Gwen Stacy Died",
      "issue_number": "121",
      "volume": {
        "id": 2127,
        "name": "The Amazing Spider-Man"
      },
      "cover_date": "1973-06-01",
      "image": {
        "super_url": "https://comicvine.gamespot.com/...",
        "medium_url": "https://comicvine.gamespot.com/..."
      }
    }
  ]
}
```

#### 2. Get Issue Details
```http
GET /api/issue/4000-{issue_id}/?api_key={key}&format=json
```

**Response**:
```json
{
  "status_code": 1,
  "results": {
    "id": 140529,
    "name": "The Night Gwen Stacy Died",
    "issue_number": "121",
    "volume": {
      "id": 2127,
      "name": "The Amazing Spider-Man",
      "publisher": {
        "id": 31,
        "name": "Marvel Comics"
      }
    },
    "cover_date": "1973-06-01",
    "store_date": null,
    "description": "...",
    "person_credits": [
      {"id": 40439, "name": "Stan Lee", "role": "writer"},
      {"id": 40440, "name": "Gil Kane", "role": "penciller"},
      {"id": 40441, "name": "John Romita Sr.", "role": "inker"}
    ],
    "character_credits": [
      {"id": 1443, "name": "Spider-Man"},
      {"id": 41233, "name": "Green Goblin"}
    ],
    "story_arc_credits": [
      {"id": 55654, "name": "The Death of Gwen Stacy"}
    ],
    "image": {
      "super_url": "https://comicvine.gamespot.com/a/uploads/scale_large/12/124259/8235305-00.jpg",
      "original_url": "https://comicvine.gamespot.com/a/uploads/original/12/124259/8235305-00.jpg"
    }
  }
}
```

#### 3. Get Volume (Series) Details
```http
GET /api/volume/4050-{volume_id}/?api_key={key}&format=json
```

**Response**:
```json
{
  "status_code": 1,
  "results": {
    "id": 2127,
    "name": "The Amazing Spider-Man",
    "start_year": "1963",
    "publisher": {
      "id": 31,
      "name": "Marvel Comics"
    },
    "count_of_issues": 801,
    "description": "...",
    "image": {
      "super_url": "https://comicvine.gamespot.com/..."
    }
  }
}
```

#### 4. Get Publisher Details
```http
GET /api/publisher/4010-{publisher_id}/?api_key={key}&format=json
```

---

## Implementation Checklist

### Phase 1: Core Integration
- [ ] Register for ComicVine API key (https://comicvine.gamespot.com/api/)
- [ ] Implement API client with rate limiting (200 req/hour per resource)
- [ ] Implement search functionality (issues, volumes, publishers)
- [ ] Implement issue metadata fetching (title, number, cover date, description)
- [ ] Store ComicVine IDs as `external_id` in `comics` table
- [ ] Implement cover image download and caching
- [ ] Handle 429 (rate limit exceeded) responses gracefully

### Phase 2: Metadata Enrichment
- [ ] Fetch creator credits (writers, artists, colorists, inkers, letterers)
- [ ] Store creator data in `comic_creators` and `comic_creator_roles` tables
- [ ] Fetch character appearances (link to characters database)
- [ ] Fetch story arc information (crossover events)
- [ ] Implement publisher metadata caching (reduce API calls)
- [ ] Implement volume (series) metadata caching

### Phase 3: Advanced Features
- [ ] Implement batch metadata fetching (queue system for new issues)
- [ ] Implement automatic metadata refresh (weekly/monthly)
- [ ] Use `field_list` parameter to reduce API payload size
- [ ] Implement fallback to Marvel API for Marvel-specific content (optional)
- [ ] Implement fallback to GCD for historical/public domain comics (optional)
- [ ] Add ComicVine attribution (required by ToS)

### Phase 4: Background Jobs (River)
- [ ] Job: `FetchComicMetadataArgs` (fetch metadata for single issue)
- [ ] Job: `RefreshComicsMetadataArgs` (batch refresh for series)
- [ ] Job: `DownloadComicCoverArgs` (download and cache cover images)
- [ ] Job: `SyncComicVinePublishersArgs` (sync publisher data)

---

## Integration Pattern

### Metadata Fetch Flow
```
User adds comic to library (CBZ file scanned)
                ↓
Extract ComicInfo.xml from CBZ (if present)
                ↓
Parse series name, volume, issue number
                ↓
Search ComicVine API for matching issue
                ↓
Fetch detailed metadata (creators, characters, story arcs)
                ↓
Download cover image (if not present in CBZ)
                ↓
Store metadata in PostgreSQL (comics table, metadata_json JSONB)
                ↓
Link creators (comic_creators, comic_creator_roles)
```

### Search Flow
```
User searches for "Amazing Spider-Man #121"
                ↓
Query ComicVine Search API (resources=issue)
                ↓
Parse results (match volume + issue number)
                ↓
Display results with covers (user selects)
                ↓
Fetch full metadata for selected issue
                ↓
Optionally: Add to download queue (Mylar3 integration)
```

### Rate Limiting Strategy
```go
// Rate limiter: 200 requests per hour per resource
// Use separate buckets for different resources (issues, volumes, publishers)

type ComicVineClient struct {
    issueRateLimiter    *rate.Limiter  // 200 per hour
    volumeRateLimiter   *rate.Limiter  // 200 per hour
    publisherRateLimiter *rate.Limiter // 200 per hour
}

// Example: 200 requests per hour = ~3.33 requests per minute
// Use token bucket: refill 1 token every 18 seconds
limiter := rate.NewLimiter(rate.Every(18*time.Second), 200)
```

---

## Related Documentation

- **Comics Module**: `docs/features/COMICS_MODULE.md` (overall comics architecture)
- **Marvel API Integration**: `docs/integrations/metadata/comics/MARVEL_API.md` (Marvel-specific metadata)
- **Grand Comics Database**: `docs/integrations/metadata/comics/GRAND_COMICS_DATABASE.md` (historical fallback)
- **Mylar3 Integration**: `docs/integrations/servarr/MYLAR3.md` (comics download automation - future)

---

## Notes

- **ComicVine API Key**: Free but requires GameSpot account registration
- **Rate Limits**: 200 requests per resource per hour (strict enforcement)
  - Use separate rate limiters for `issues`, `volumes`, `publishers` resources
  - Cache metadata aggressively to minimize API calls
- **Attribution Required**: ComicVine ToS requires attribution ("Powered by ComicVine API")
- **Cover Images**: Use `super_url` (large) or `original_url` (full resolution)
- **ComicInfo.xml**: Parse from CBZ files first (reduces API calls, instant metadata)
- **Fallback Strategy**: ComicVine primary → Marvel API (Marvel only) → GCD (historical)
- **Volume vs Issue**: Volume = Series (e.g., "The Amazing Spider-Man Vol. 1"), Issue = Single comic (#121)
- **Store Date vs Cover Date**: `store_date` (actual release) often null, use `cover_date` instead
- **Person Credits**: ComicVine uses `person_credits` array with roles (writer, penciller, inker, colorist, letterer, cover artist, editor)
- **JSONB Storage**: Store full ComicVine response in `metadata_json` for flexibility (future-proofing)
- **No Manga**: ComicVine has limited manga coverage, use AniList for Japanese comics
- **Historical Comics**: GCD better for Golden/Silver Age comics (pre-1980), ComicVine better for modern comics
