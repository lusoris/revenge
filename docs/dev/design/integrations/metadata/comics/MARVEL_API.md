# Marvel API Integration

> Official Marvel Comics metadata (Marvel Universe content only)

**Priority**: 🟢 LOW (Phase 7 - Comics Module, optional enrichment)
**Provider**: Marvel Entertainment

## Status

| Dimension | Status | Notes |
| --------- | ------ | ----- |
| Design | ✅ | Comprehensive REST API spec, authentication, rate limiting |
| Sources | ✅ | Base URL, documentation, authentication examples linked |
| Instructions | ✅ | Phased implementation checklist with enrichment flow |
| Code | 🔴 | |
| Linting | 🔴 | |
| Unit Testing | 🔴 | |
| Integration Testing | 🔴 | |

---

## Overview

**Marvel API** is the official API from Marvel Comics, providing metadata for Marvel Universe content (characters, comics, creators, events, series). It's useful as a **supplementary source** for Marvel-specific comics, but ComicVine remains the primary metadata provider due to publisher-agnostic coverage.

**Why Marvel API** (as supplement):
- Official Marvel data (authoritative for Marvel content)
- Character bios and comic appearances
- Event/crossover data (Civil War, Secret Wars, etc.)
- Creator information (Stan Lee, Jack Kirby, etc.)
- Free API with reasonable rate limits (3000 requests/day)

**Limitations**:
- **Marvel content only** (no DC, Image, Dark Horse, etc.)
- **Rate limits**: 3000 requests per day (lower than ComicVine)
- **Less comprehensive**: ComicVine often has more detailed metadata
- **No cover images**: Marvel API does not provide high-quality cover images
- **Requires attribution**: "Data provided by Marvel. © 2024 Marvel" required in UI

**Use Case**: Optional enrichment for Marvel comics (character bios, event data), NOT primary metadata source.

---

## Developer Resources

### API Documentation
- **Base URL**: `https://gateway.marvel.com/v1/public/`
- **Documentation**: https://developer.marvel.com/documentation/getting_started
- **API Version**: v1 (stable)
- **Authentication**: API Key + Timestamp + MD5 Hash
- **Rate Limits**: 3000 requests per day
- **Response Format**: JSON

### Key Features
- **Comics**: Comic issue metadata (title, description, dates, prices)
- **Characters**: Character bios, appearances in comics
- **Creators**: Creator information (writers, artists)
- **Events**: Crossover events (Civil War, Infinity War, etc.)
- **Series**: Comic series metadata (Amazing Spider-Man, X-Men, etc.)

### Authentication
**Method**: API Key + Timestamp + Hash
**Required Parameters**:
- `apikey`: Your public API key
- `ts`: Timestamp (e.g., `1`)
- `hash`: MD5 hash of `ts + private_key + public_key`

**Example Hash Generation** (Go):
```go
import (
    "crypto/md5"
    "fmt"
    "time"
)

func generateMarvelHash(publicKey, privateKey string) (ts string, hash string) {
    ts = fmt.Sprintf("%d", time.Now().Unix())
    data := ts + privateKey + publicKey
    hashBytes := md5.Sum([]byte(data))
    hash = fmt.Sprintf("%x", hashBytes)
    return ts, hash
}

// Usage:
// ts, hash := generateMarvelHash(publicKey, privateKey)
// url := fmt.Sprintf("https://gateway.marvel.com/v1/public/comics?apikey=%s&ts=%s&hash=%s", publicKey, ts, hash)
```

---

## API Details

### Authentication
**Method**: API Key + Timestamp + MD5 Hash
**Rate Limit**: 3000 requests per day (resets at midnight EST)

**Example Request**:
```bash
# Generate hash: md5(ts + privateKey + publicKey)
# Example: ts=1, privateKey=abcdef, publicKey=123456
# hash = md5("1abcdef123456") = "ffd275c5130566a2916217b101f26150"

curl "https://gateway.marvel.com/v1/public/comics/82967?apikey=123456&ts=1&hash=ffd275c5130566a2916217b101f26150"
```

### Core Endpoints

#### 1. Search Comics
```http
GET /v1/public/comics?titleStartsWith={title}&apikey={key}&ts={ts}&hash={hash}
```

**Parameters**:
- `titleStartsWith`: Comic title (e.g., "Amazing Spider-Man")
- `issueNumber`: Issue number (e.g., "1")
- `limit`: Results per page (max 100, default 20)
- `offset`: Pagination offset
- `orderBy`: Sort order (`focDate`, `onsaleDate`, `title`, `issueNumber`)

**Response**:
```json
{
  "code": 200,
  "status": "Ok",
  "data": {
    "offset": 0,
    "limit": 20,
    "total": 500,
    "count": 20,
    "results": [
      {
        "id": 82967,
        "digitalId": 0,
        "title": "Amazing Spider-Man (2018) #1",
        "issueNumber": 1,
        "description": "...",
        "isbn": "",
        "pageCount": 80,
        "dates": [
          {"type": "onsaleDate", "date": "2018-07-11T00:00:00-0400"},
          {"type": "focDate", "date": "2018-06-18T00:00:00-0400"}
        ],
        "prices": [
          {"type": "printPrice", "price": 5.99}
        ],
        "thumbnail": {
          "path": "http://i.annihil.us/u/prod/marvel/i/mg/c/e0/5b2c4ff46c4c8",
          "extension": "jpg"
        },
        "creators": {
          "items": [
            {"name": "Nick Spencer", "role": "writer"},
            {"name": "Ryan Ottley", "role": "penciler"}
          ]
        },
        "characters": {
          "items": [
            {"name": "Spider-Man (Peter Parker)"}
          ]
        },
        "events": {
          "items": []
        }
      }
    ]
  }
}
```

#### 2. Get Comic Details
```http
GET /v1/public/comics/{comic_id}?apikey={key}&ts={ts}&hash={hash}
```

#### 3. Get Character Details
```http
GET /v1/public/characters/{character_id}?apikey={key}&ts={ts}&hash={hash}
```

**Response**:
```json
{
  "code": 200,
  "data": {
    "results": [
      {
        "id": 1009610,
        "name": "Spider-Man (Peter Parker)",
        "description": "Bitten by a radioactive spider...",
        "thumbnail": {
          "path": "http://i.annihil.us/u/prod/marvel/i/mg/3/50/526548a343e4b",
          "extension": "jpg"
        },
        "comics": {
          "available": 4500,
          "items": [
            {"name": "Amazing Spider-Man #1"}
          ]
        }
      }
    ]
  }
}
```

#### 4. Get Event Details
```http
GET /v1/public/events/{event_id}?apikey={key}&ts={ts}&hash={hash}
```

**Example**: Civil War, Secret Wars, Infinity War, etc.

---

## Implementation Checklist

### Phase 1: Core Integration
- [ ] Register for Marvel API keys (https://developer.marvel.com/account)
- [ ] Implement authentication (timestamp + MD5 hash generation)
- [ ] Implement rate limiting (3000 requests/day, track usage)
- [ ] Implement comic search (by title, issue number)
- [ ] Implement comic metadata fetching (supplementary to ComicVine)
- [ ] Store Marvel Comic IDs in `metadata_json` JSONB field

### Phase 2: Character & Event Enrichment
- [ ] Fetch character information (bios, appearances)
- [ ] Store character bios in `metadata_json` (link to characters table)
- [ ] Fetch event/crossover data (Civil War, Secret Wars, etc.)
- [ ] Link comics to story arcs/events
- [ ] Implement creator metadata fetching (Marvel-specific credits)

### Phase 3: Optimization
- [ ] Cache Marvel API responses (reduce daily request count)
- [ ] Use ComicVine as primary, Marvel API as enrichment only
- [ ] Implement daily request counter (reset at midnight EST)
- [ ] Add Marvel attribution in UI ("Data provided by Marvel. © 2024 Marvel")
- [ ] Throttle requests (3000/day = ~2 requests/minute average)

### Phase 4: Background Jobs (River)
- [ ] Job: `EnrichMarvelComicArgs` (enrich Marvel comics with official data)
- [ ] Job: `FetchMarvelCharacterArgs` (fetch character bios)
- [ ] Job: `FetchMarvelEventArgs` (fetch crossover event data)

---

## Integration Pattern

### Metadata Enrichment Flow (Marvel Comics Only)
```
Comic scanned from library (CBZ file)
                ↓
Fetch ComicVine metadata (primary)
                ↓
Check if Marvel comic (publisher_id = Marvel)
                ↓
Query Marvel API for supplementary data (character bios, events)
                ↓
Merge Marvel API data into metadata_json
                ↓
Store enriched metadata in PostgreSQL
```

### Character Enrichment Flow
```
User views Marvel comic (e.g., Amazing Spider-Man #1)
                ↓
Display characters from ComicVine (character appearances)
                ↓
Background: Fetch character bios from Marvel API
                ↓
Cache character bios (reduce API calls)
                ↓
Display enriched character info in UI
```

### Rate Limiting Strategy
```go
// Marvel API: 3000 requests per day
// Track daily usage, reset at midnight EST

type MarvelAPIClient struct {
    requestCounter int
    lastReset      time.Time
    maxRequests    int // 3000
}

func (c *MarvelAPIClient) checkRateLimit() error {
    now := time.Now()
    if now.Day() != c.lastReset.Day() {
        c.requestCounter = 0
        c.lastReset = now
    }
    if c.requestCounter >= c.maxRequests {
        return fmt.Errorf("Marvel API rate limit exceeded (3000/day)")
    }
    c.requestCounter++
    return nil
}
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
- [Grand Comics Database (GCD) Integration](GRAND_COMICS_DATABASE.md)

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
- **Grand Comics Database**: `docs/integrations/metadata/comics/GRAND_COMICS_DATABASE.md` (historical fallback)
- **Comics Module**: `docs/features/COMICS_MODULE.md` (overall comics architecture)

---

## Notes

- **NOT Primary Source**: ComicVine remains primary metadata provider (publisher-agnostic)
- **Marvel Content Only**: API only covers Marvel Comics (no DC, Image, etc.)
- **Rate Limits**: 3000 requests/day (strict enforcement, no hourly bucket)
  - Track daily usage carefully
  - Use as enrichment only (NOT for every comic)
- **Attribution Required**: "Data provided by Marvel. © 2024 Marvel" in UI
- **No Cover Images**: Marvel API does not provide high-quality cover images (use ComicVine)
- **Hash Authentication**: MD5 hash of `ts + privateKey + publicKey` required for every request
- **Timestamp**: Can use Unix timestamp OR any arbitrary number (commonly use `1` for simplicity)
- **Thumbnail Images**: Low resolution (~100x150), use ComicVine for covers
- **Character IDs**: Marvel character IDs useful for linking to Marvel Universe wiki/database
- **Event Data**: Excellent source for crossover event information (Civil War, Secret Wars, etc.)
- **Use Case**: Best for Marvel-specific enrichment (character bios, event context), NOT comprehensive metadata
- **Fallback Strategy**: ComicVine (primary) → Marvel API (Marvel enrichment) → GCD (historical)
- **JSONB Storage**: Store Marvel API responses in `metadata_json.marvel_data` for future use
- **Daily Reset**: Rate limit resets at midnight EST (account for timezone)
