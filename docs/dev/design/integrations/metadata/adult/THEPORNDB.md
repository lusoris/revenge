# ThePornDB Integration

> Alternative adult metadata provider with scene/performer/studio data

**Service**: ThePornDB (https://theporndb.net)
**API**: REST API v1 (https://api.theporndb.net)
**Category**: Metadata Provider (Adult Content)
**Priority**: 🟢 MEDIUM (Fallback to StashDB)

## Status

| Dimension | Status | Notes |
| --------- | ------ | ----- |
| Design | ✅ | Comprehensive REST API endpoints, data mapping, fallback logic |
| Sources | ✅ | Base URL, API docs, Swagger UI linked |
| Instructions | ✅ | Phased implementation checklist with fallback strategy |
| Code | 🔴 | |
| Linting | 🔴 | |
| Unit Testing | 🔴 | |
| Integration Testing | 🔴 | |

---

## Overview

**ThePornDB** (TPDb) is a community-driven adult content metadata database similar to TMDb but for adult content. It provides comprehensive metadata for scenes, performers, and studios with a REST API.

**Key Features**:
- **Scenes**: Scene titles, release dates, performers, studios, tags, cover images, descriptions
- **Performers**: Names, aliases, birthdate, measurements, tattoos, piercings, career info, images
- **Studios**: Studio names, logos, parent companies, networks
- **Sites**: Site-specific content (e.g., Brazzers.com, Reality Kings, etc.)
- **Free API**: REST API with API Key authentication (free registration)
- **Multi-language**: Supports multiple languages (English primary)
- **Community-driven**: User submissions with moderation

**Use Cases**:
- Fallback adult metadata source (when StashDB lacks data)
- Scene/performer/studio metadata enrichment
- Alternative to StashDB for specific content gaps

**⚠️ CRITICAL: Adult Content Isolation**:
- **Database schema**: `qar` schema ONLY (`qar.expeditions`, `qar.voyages`, `qar.crew`, `qar.ports`)
- **API namespace**: `/api/v1/qar/metadata/theporndb/*` (NOT `/api/v1/metadata/theporndb/*`)
- **Module location**: `internal/content/qar/metadata/theporndb/` (NOT `internal/service/metadata/`)
- **Access control**: Mods/admins can see all data for monitoring, regular users see only their own library

---

## Developer Resources

### API Documentation
- **Base URL**: https://api.theporndb.net
- **API Docs**: https://api.theporndb.net/docs
- **Swagger UI**: https://api.theporndb.net/swagger
- **GitHub**: None (closed-source API)

### Authentication
- **Method**: API Key (header-based)
- **Header**: `Authorization: Bearer YOUR_API_KEY`
- **Registration**: https://theporndb.net/register (free)
- **Rate Limits**: 120 requests/minute (2 requests/second)

### Data Coverage
- **Scenes**: 500K+ scenes (English-language focus)
- **Performers**: 100K+ performers (active + retired)
- **Studios**: 5K+ studios
- **Sites**: 1K+ adult sites (Brazzers, Reality Kings, etc.)
- **Updates**: Daily (community submissions + moderation)

### Go Client Library
- **Official**: None
- **Recommended**: Use `net/http` with JSON parsing
- **Alternative**: Create custom client wrapper

---

## API Details

### REST Endpoints

#### Search Scenes
```
GET /scenes/search
Query Parameters:
  - q: Search query (scene title)
  - page: Page number (default 1)
  - per_page: Results per page (max 25, default 25)

Response:
{
  "data": [
    {
      "id": "uuid",
      "title": "Scene Title",
      "description": "Scene description...",
      "release_date": "2024-01-15",
      "site": {
        "id": "site-uuid",
        "name": "Brazzers",
        "url": "https://brazzers.com"
      },
      "studio": {
        "id": "studio-uuid",
        "name": "Brazzers"
      },
      "performers": [
        {
          "id": "performer-uuid",
          "name": "Performer Name"
        }
      ],
      "tags": [
        {"id": "tag-uuid", "name": "POV"}
      ],
      "posters": [
        {"url": "https://cdn.theporndb.net/...", "size": "large"}
      ]
    }
  ],
  "meta": {
    "current_page": 1,
    "total_pages": 10,
    "total_results": 250
  }
}
```

#### Get Scene Details
```
GET /scenes/{id}

Response:
{
  "data": {
    "id": "uuid",
    "title": "Scene Title",
    "description": "Full scene description...",
    "release_date": "2024-01-15",
    "duration": 2400,  // seconds
    "site": {
      "id": "site-uuid",
      "name": "Brazzers",
      "url": "https://brazzers.com",
      "network": {
        "id": "network-uuid",
        "name": "MindGeek"
      }
    },
    "studio": {
      "id": "studio-uuid",
      "name": "Brazzers",
      "parent": {
        "id": "parent-uuid",
        "name": "MindGeek"
      }
    },
    "performers": [
      {
        "id": "performer-uuid",
        "name": "Performer Name",
        "role": "Female"
      }
    ],
    "tags": [
      {"id": "tag-uuid", "name": "POV", "category": "Position"}
    ],
    "posters": [
      {"url": "https://cdn.theporndb.net/...", "size": "large"},
      {"url": "https://cdn.theporndb.net/...", "size": "medium"}
    ],
    "backgrounds": [
      {"url": "https://cdn.theporndb.net/...", "size": "full"}
    ]
  }
}
```

#### Search Performers
```
GET /performers/search
Query Parameters:
  - q: Search query (performer name)
  - page: Page number
  - per_page: Results per page (max 25)

Response:
{
  "data": [
    {
      "id": "uuid",
      "name": "Performer Name",
      "disambiguation": null,
      "aliases": ["Alias1", "Alias2"],
      "birthdate": "1990-05-15",
      "gender": "Female",
      "ethnicity": "Caucasian",
      "country": "USA",
      "measurements": "34D-24-36",
      "height": 165,  // cm
      "weight": 55,   // kg
      "posters": [
        {"url": "https://cdn.theporndb.net/...", "size": "large"}
      ]
    }
  ]
}
```

#### Get Performer Details
```
GET /performers/{id}

Response:
{
  "data": {
    "id": "uuid",
    "name": "Performer Name",
    "disambiguation": null,
    "aliases": ["Alias1", "Alias2"],
    "birthdate": "1990-05-15",
    "death_date": null,
    "gender": "Female",
    "ethnicity": "Caucasian",
    "country": "USA",
    "eye_color": "Blue",
    "hair_color": "Blonde",
    "measurements": "34D-24-36",
    "height": 165,
    "weight": 55,
    "tattoos": [
      {"location": "Lower back", "description": "Tribal design"}
    ],
    "piercings": [
      {"location": "Navel", "description": "Standard"}
    ],
    "career_start_year": 2015,
    "career_end_year": null,
    "posters": [
      {"url": "https://cdn.theporndb.net/...", "size": "large"}
    ],
    "biography": "Performer biography..."
  }
}
```

#### Search Studios
```
GET /studios/search
Query Parameters:
  - q: Search query (studio name)
  - page: Page number

Response:
{
  "data": [
    {
      "id": "uuid",
      "name": "Brazzers",
      "parent": {
        "id": "parent-uuid",
        "name": "MindGeek"
      },
      "logo": "https://cdn.theporndb.net/..."
    }
  ]
}
```

#### Get Studio Details
```
GET /studios/{id}

Response:
{
  "data": {
    "id": "uuid",
    "name": "Brazzers",
    "parent": {
      "id": "parent-uuid",
      "name": "MindGeek"
    },
    "subsidiaries": [
      {"id": "sub-uuid", "name": "Brazzers Network"}
    ],
    "logo": "https://cdn.theporndb.net/...",
    "description": "Studio description..."
  }
}
```

---

## Implementation Checklist

### Phase 1: Core Integration
- [ ] REST API client setup (Go `net/http`)
- [ ] API Key configuration (`configs/config.yaml` - `theporndb.api_key`)
- [ ] **Adult schema**: Use existing `qar.crew`, `qar.ports`, `qar.voyages` tables
- [ ] **API namespace**: `/api/v1/qar/metadata/theporndb/*` endpoints
- [ ] **Module location**: `internal/content/qar/metadata/theporndb/` (isolated)
- [ ] Basic scene search (REST `/scenes/search`)
- [ ] Scene details fetch (REST `/scenes/{id}`)
- [ ] Performer search (REST `/performers/search`)
- [ ] Performer details fetch (REST `/performers/{id}`)
- [ ] Studio search (REST `/studios/search`)
- [ ] Image downloads (posters, performer images, studio logos)
- [ ] JSONB storage (`qar.expeditions.metadata_json.theporndb_data`)

### Phase 2: Fallback Logic
- [ ] Fallback to ThePornDB when StashDB lacks data
- [ ] Merge StashDB + ThePornDB metadata (prefer StashDB, supplement with TPDb)
- [ ] Conflict resolution (prefer StashDB data, use TPDb for missing fields)
- [ ] Data quality scoring (prefer source with more complete data)

### Phase 3: Background Jobs (River)
- [ ] **Job**: `qar.metadata.theporndb.fetch_voyage` (fetch voyage/scene metadata)
- [ ] **Job**: `qar.metadata.theporndb.fetch_crew` (fetch crew/performer metadata)
- [ ] **Job**: `qar.metadata.theporndb.refresh_metadata` (monthly refresh fallback data)
- [ ] Rate limiting (120 req/min = 2 req/sec token bucket)
- [ ] Retry logic (exponential backoff for failures)

---

## Integration Pattern

### Fallback Metadata Fetch Flow
```
Scene missing metadata
        ↓
Check StashDB first (primary source)
        ↓
StashDB data found? → Use StashDB metadata
        ↓
        NO
        ↓
Search ThePornDB (fallback)
        ↓
ThePornDB data found? → Fetch scene details (REST /scenes/{id})
              ↓
              Extract: title, release_date, performers, studio, tags, images
              ↓
              Store in qar.expeditions.metadata_json.theporndb_data
              ↓
              Download voyage cover image
              ↓
              Fetch crew details (REST /performers/{id})
              ↓
              Store in qar.crew table
              ↓
              Download crew images
              ↓
              Fetch port details (REST /studios/{id})
              ↓
              Store in qar.ports table
              ↓
              Update Typesense search index
              ↓
              Notify user: "Metadata updated from ThePornDB"
        ↓
        NO
        ↓
Manual entry (user adds metadata manually)
```

### Rate Limiting Strategy
```
ThePornDB rate limit: 120 req/min (2 req/sec)
- Token bucket: 120 tokens, refill 2 tokens/sec
- Batch requests: NOT supported (REST API, no batch endpoint)
- Caching: Cache performer/studio data for 30 days
- Background jobs: Use River queue (prioritize user-initiated requests)
- Exponential backoff: Retry with delay if 429 errors
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
| [ThePornDB API](https://api.theporndb.net/docs) | [Local](../../../../sources/apis/theporndb.md) |

<!-- SOURCE-BREADCRUMBS-END -->

<!-- DESIGN-BREADCRUMBS-START -->

## Related Design Docs

> Auto-generated cross-references to related design documentation

**Category**: [Adult](INDEX.md)

### In This Section

- [Stash Integration](STASH.md)
- [StashDB Integration](STASHDB.md)
- [Whisparr v3 (eros) - Adult Content Structure Analysis](WHISPARR_V3_ANALYSIS.md)

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

- [STASHDB.md](./STASHDB.md) - Primary adult metadata provider (use first)
- [STASH.md](./STASH.md) - Stash self-hosted organizer
- [WHISPARR.md](../../servarr/WHISPARR.md) - Adult content management
- [ADULT_METADATA.md](../../../ADULT_METADATA.md) - Adult metadata system architecture

---

## Notes

### API Key (Free)
- Register at https://theporndb.net/register
- Free API access (no paid tiers)
- Rate limit: 120 req/min (reasonable for most use cases)

### Rate Limits (120 req/min)
- **Hard limit**: 120 requests per minute (2 requests per second)
- **Enforcement**: HTTP 429 status if exceeded
- **Token bucket**: Implement client-side rate limiter (120 tokens, refill 2/sec)
- **Retry-After header**: Respect `Retry-After` header in 429 responses

### Data Quality
- **Community-driven**: User submissions with moderation
- **Moderate quality**: Less comprehensive than StashDB for some content
- **Site-specific**: Excellent for scene-site relationships (Brazzers, Reality Kings, etc.)
- **Use case**: Fallback when StashDB lacks data

### Adult Content Isolation (CRITICAL)
- **Database schema**: `qar` schema ONLY
  - `qar.expeditions.metadata_json.theporndb_data` (JSONB)
  - `qar.voyages.metadata_json.theporndb_data` (JSONB)
  - `qar.crew` (shared with StashDB)
  - `qar.ports` (shared with StashDB)
- **API namespace**: `/api/v1/qar/metadata/theporndb/*` (isolated)
  - `/api/v1/qar/metadata/theporndb/search/voyages`
  - `/api/v1/qar/metadata/theporndb/voyages/{tpdb_id}`
  - `/api/v1/qar/metadata/theporndb/crew/{tpdb_id}`
- **Module location**: `internal/content/qar/metadata/theporndb/` (isolated)
- **Access control**: Mods/admins see all, regular users see only their library

### Fallback Strategy
- **StashDB primary**: Always check StashDB first (better fingerprinting, better performer data)
- **ThePornDB fallback**: Use ThePornDB when:
  - StashDB lacks scene data
  - StashDB performer data incomplete
  - Site-specific metadata needed (e.g., Brazzers scene ID)
- **Merge metadata**: Combine StashDB + ThePornDB data (prefer StashDB, fill gaps with TPDb)

### JSONB Storage
- Store ThePornDB response in `qar.expeditions.metadata_json.theporndb_data`
- Separate from StashDB data (`stashdb_data` field)
- Allows querying both sources independently

### Site vs Studio
- **Site**: Specific website (e.g., Brazzers.com, RealityKings.com)
- **Studio**: Production company (e.g., Brazzers studio, Reality Kings studio)
- **Network**: Parent company (e.g., MindGeek owns Brazzers, Reality Kings, etc.)
- **Hierarchy**: Network → Studio → Site → Scene

### Image URLs
- **Direct URLs**: ThePornDB provides direct CDN URLs (no auth required)
- **Download**: Download and store locally (avoid hotlinking)
- **Sizes**: Multiple sizes (large, medium, small)

### Performer Aliases
- **Aliases array**: List of stage names/alternate spellings
- **Disambiguation**: Optional field (e.g., "Performer Name (Studio)")
- **Matching**: Search by aliases when matching performers

### Use Case: Supplementary Metadata
- **Primary source**: StashDB (better fingerprinting, better performer data)
- **Supplement**: ThePornDB for site-specific metadata (Brazzers scene ID, site URL)
- **Example**: StashDB provides performer measurements → ThePornDB provides Brazzers.com URL
