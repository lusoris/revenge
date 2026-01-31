# StashDB Integration

> Adult metadata database for performers, studios, and scenes

**Service**: StashDB (https://stashdb.org)
**API**: GraphQL API (https://stashdb.org/graphql)
**Category**: Metadata Provider (Adult Content)
**Priority**: 🟡 HIGH (Adult module core metadata)

## Status

| Dimension | Status | Notes |
| --------- | ------ | ----- |
| Design | ✅ | Comprehensive GraphQL API spec, fingerprinting, data mapping |
| Sources | ✅ | GraphQL schema, playground, API docs, GitHub linked |
| Instructions | ✅ | Phased implementation checklist with fingerprinting |
| Code | 🔴 | |
| Linting | 🔴 | |
| Unit Testing | 🔴 | |
| Integration Testing | 🔴 | |

---

## Overview

**StashDB** is a community-driven adult content metadata database with comprehensive information about performers, studios, and scenes. It is the primary metadata provider for Stash (self-hosted adult media organizer) and the de facto standard for adult content metadata.

**Key Features**:
- **Performers**: Names, aliases, measurements, tattoos, piercings, career start/end dates, images
- **Studios**: Studio names, parent companies, logos, URLs
- **Scenes**: Scene titles, release dates, performers, studios, tags, cover images, scene markers
- **Community-driven**: User submissions with moderation workflow
- **Free API**: GraphQL API with API Key authentication (free registration)
- **Rich relationships**: Performer-scene, studio-scene, parent/subsidiary studios
- **Fingerprinting**: Perceptual hashing (phash) for scene identification
- **Tags**: Comprehensive tag system (positions, acts, settings, fetishes)

**Use Cases**:
- Primary adult content metadata source
- Automatic scene identification via fingerprinting
- Performer/studio metadata enrichment
- Tag-based content organization
- Duplicate scene detection

**⚠️ CRITICAL: Adult Content Isolation**:
- **Database schema**: `qar` schema ONLY (`qar.expeditions`, `qar.voyages`, `qar.crew`, `qar.ports`)
- **API namespace**: `/api/v1/qar/metadata/stashdb/*` (NOT `/api/v1/metadata/stashdb/*`)
- **Module location**: `internal/content/qar/metadata/stashdb/` (NOT `internal/service/metadata/`)
- **Access control**: Mods/admins can see all data for monitoring, regular users see only their own library

---

## Developer Resources

### API Documentation
- **GraphQL Schema**: https://stashdb.org/graphql
- **GraphQL Playground**: https://stashdb.org/graphql (interactive query editor)
- **API Docs**: https://github.com/stashapp/stashdb/blob/develop/docs/API.md
- **GitHub**: https://github.com/stashapp/stashdb

### Authentication
- **Method**: API Key (header-based)
- **Header**: `ApiKey: YOUR_API_KEY`
- **Registration**: https://stashdb.org/register (free)
- **Rate Limits**: Reasonable usage (no hard limit, but avoid abuse)

### Data Coverage
- **Scenes**: 1M+ scenes (English-language focus, Western studios)
- **Performers**: 200K+ performers (active + retired)
- **Studios**: 10K+ studios (parent/subsidiary relationships)
- **Tags**: 500+ tags (positions, acts, settings, fetishes)
- **Updates**: Real-time (community submissions + moderation)

### Go Client Library
- **Official**: None (use standard GraphQL client)
- **Recommended**: `github.com/machinebox/graphql` OR `github.com/Khan/genqlient`
- **Alternative**: `net/http` + manual GraphQL JSON requests

---

## API Details

### GraphQL Queries

#### Search Scenes
```graphql
query SearchScenes($query: String!, $page: Int) {
  searchScene(term: $query, page: $page) {
    count
    scenes {
      id
      title
      release_date
      studio {
        id
        name
        parent_studio {
          id
          name
        }
      }
      performers {
        performer {
          id
          name
          disambiguation
          gender
        }
        as
      }
      tags {
        id
        name
        description
      }
      images {
        url
        width
        height
      }
      fingerprints {
        hash
        algorithm
        duration
      }
      duration
      details
      url
    }
  }
}
```

#### Get Scene Details
```graphql
query GetScene($id: ID!) {
  findScene(id: $id) {
    id
    title
    release_date
    studio {
      id
      name
      url
      parent_studio {
        id
        name
      }
    }
    performers {
      performer {
        id
        name
        disambiguation
        gender
        birthdate
        ethnicity
        country
        eye_color
        hair_color
        height
        measurements
        tattoos {
          location
          description
        }
        piercings {
          location
          description
        }
        images {
          url
          width
          height
        }
      }
      as
    }
    tags {
      id
      name
      description
      category
    }
    images {
      url
      width
      height
    }
    fingerprints {
      hash
      algorithm
      duration
    }
    duration
    director
    details
    url
    code
  }
}
```

#### Search Performers
```graphql
query SearchPerformers($query: String!, $page: Int) {
  searchPerformer(term: $query, page: $page) {
    count
    performers {
      id
      name
      disambiguation
      gender
      birthdate
      ethnicity
      country
      eye_color
      hair_color
      height
      measurements
      career_start_year
      career_end_year
      tattoos {
        location
        description
      }
      piercings {
        location
        description
      }
      images {
        url
        width
        height
      }
      aliases
    }
  }
}
```

#### Get Performer Details
```graphql
query GetPerformer($id: ID!) {
  findPerformer(id: $id) {
    id
    name
    disambiguation
    gender
    birthdate
    death_date
    ethnicity
    country
    eye_color
    hair_color
    height
    measurements
    breast_type
    career_start_year
    career_end_year
    tattoos {
      location
      description
    }
    piercings {
      location
      description
    }
    images {
      url
      width
      height
    }
    aliases
    urls {
      url
      site {
        name
        icon
      }
    }
  }
}
```

#### Search Studios
```graphql
query SearchStudios($query: String!, $page: Int) {
  searchStudio(term: $query, page: $page) {
    count
    studios {
      id
      name
      url
      parent_studio {
        id
        name
      }
      child_studios {
        id
        name
      }
      images {
        url
        width
        height
      }
    }
  }
}
```

#### Fingerprint Lookup (Scene Identification)
```graphql
query FingerprintLookup($fingerprints: [FingerprintQueryInput!]!) {
  queryFingerprints(fingerprints: $fingerprints) {
    hash
    duration
    submissions {
      id
      title
      release_date
      studio {
        name
      }
      performers {
        performer {
          name
        }
      }
    }
  }
}
```

---

## Implementation Checklist

### Phase 1: Core Integration
- [ ] GraphQL client setup (`machinebox/graphql` OR `genqlient`)
- [ ] API Key configuration (`configs/config.yaml` - `stashdb.api_key`)
- [ ] **Adult schema**: `qar.crew`, `qar.ports`, `qar.voyages` tables (NOT public schema)
- [ ] **API namespace**: `/api/v1/qar/metadata/stashdb/*` endpoints
- [ ] **Module location**: `internal/content/qar/metadata/stashdb/` (isolated from public metadata)
- [ ] Basic scene search (GraphQL `searchScene` query)
- [ ] Scene details fetch (GraphQL `findScene` query)
- [ ] Performer search (GraphQL `searchPerformer` query)
- [ ] Performer details fetch (GraphQL `findPerformer` query)
- [ ] Studio search (GraphQL `searchStudio` query)
- [ ] Image downloads (performer images, scene covers, studio logos)
- [ ] JSONB storage (`qar.expeditions.metadata_json.stashdb_data`)

### Phase 2: Fingerprinting & Auto-Identification
- [ ] Perceptual hashing (phash) generation for video files
- [ ] Fingerprint lookup API (GraphQL `queryFingerprints`)
- [ ] Automatic scene identification (match video → StashDB scene)
- [ ] Duplicate detection (same fingerprint = duplicate scene)
- [ ] Match confidence scoring (duration match, fingerprint similarity)
- [ ] Manual match override (user can correct mismatches)

### Phase 3: Performer & Studio Management
- [ ] Performer profiles (bio, measurements, tattoos, piercings, career dates)
- [ ] Performer image galleries (headshots, body shots)
- [ ] Performer aliases (stage names, alternate spellings)
- [ ] Studio hierarchy (parent/subsidiary relationships)
- [ ] Studio logos & branding
- [ ] Tag management (positions, acts, settings, fetishes)
- [ ] Performer scene filmography (all scenes with performer)

### Phase 4: Background Jobs (River)
- [ ] **Job**: `qar.metadata.stashdb.fetch_voyage` (fetch voyage/scene metadata)
- [ ] **Job**: `qar.metadata.stashdb.fetch_crew` (fetch crew/performer metadata)
- [ ] **Job**: `qar.metadata.stashdb.identify_voyages` (fingerprint-based identification)
- [ ] **Job**: `qar.metadata.stashdb.refresh_metadata` (weekly refresh for active content)
- [ ] Rate limiting (reasonable usage, avoid API abuse)
- [ ] Retry logic (exponential backoff for failures)

---

## Integration Pattern

### Scene Metadata Fetch Flow
```
User adds video file
        ↓
Generate perceptual hash (phash)
        ↓
StashDB fingerprint lookup (GraphQL queryFingerprints)
        ↓
Match found? → Fetch scene details (findScene)
              ↓
              Extract: title, release_date, performers, studio, tags, images
              ↓
              Store in qar.expeditions OR qar.voyages (qar schema ONLY)
              ↓
              metadata_json.stashdb_data = full GraphQL response
              ↓
              Download voyage cover image
              ↓
              Fetch crew details (findPerformer for each crew member)
              ↓
              Store in qar.crew table
              ↓
              Download crew images
              ↓
              Fetch port details (findStudio)
              ↓
              Store in qar.ports table
              ↓
              Link crew → voyage (qar.voyage_crew junction table)
              ↓
              Update Typesense search index (qar_voyages collection)
              ↓
              Notify user: "Metadata updated from StashDB"

Match not found? → Manual search UI (user searches StashDB by title)
                   ↓
                   User selects correct scene
                   ↓
                   Fetch scene details (same flow as above)
```

### Crew Search Flow
```
User searches crew name
        ↓
GraphQL searchPerformer query
        ↓
Display results: name, image, cargo (measurements), career dates
        ↓
User selects crew member
        ↓
Fetch crew details (findPerformer)
        ↓
Store in qar.crew table
        ↓
Download crew images
        ↓
Display crew profile page:
  - Bio (birthdate, ethnicity, country, measurements)
  - Images (gallery)
  - Tattoos & piercings (detailed)
  - Career info (start/end year, active status)
  - Filmography (all voyages with this crew member)
```

### Rate Limiting Strategy
```
StashDB has no hard rate limit, but avoid abuse:
- Batch requests: Fetch multiple scenes in one query (GraphQL batch)
- Caching: Cache performer/studio data for 30 days (reduce redundant requests)
- Background jobs: Use River queue (avoid overwhelming API)
- Exponential backoff: Retry with delay if 429/500 errors
- User-initiated: Prioritize user-initiated searches over background jobs
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
| [Khan/genqlient](https://pkg.go.dev/github.com/Khan/genqlient) | [Local](../../../../sources/tooling/genqlient.md) |
| [StashDB GraphQL API](https://stashdb.org/graphql) | [Local](../../../../sources/apis/stashdb-schema.graphql) |
| [genqlient GitHub README](https://github.com/Khan/genqlient) | [Local](../../../../sources/tooling/genqlient-guide.md) |

<!-- SOURCE-BREADCRUMBS-END -->

<!-- DESIGN-BREADCRUMBS-START -->

## Related Design Docs

> Auto-generated cross-references to related design documentation

**Category**: [Adult](INDEX.md)

### In This Section

- [Stash Integration](STASH.md)
- [ThePornDB Integration](THEPORNDB.md)
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

- [STASH.md](./STASH.md) - Stash self-hosted organizer (uses StashDB API)
- [THEPORNDB.md](./THEPORNDB.md) - Alternative adult metadata provider
- [WHISPARR.md](../../servarr/WHISPARR.md) - Adult content management (uses StashDB)
- [ADULT_METADATA.md](../../../ADULT_METADATA.md) - Adult metadata system architecture
- [ADULT_CONTENT_SYSTEM.md](../../../ADULT_CONTENT_SYSTEM.md) - Adult content isolation design

---

## Notes

### API Key (Free)
- Register at https://stashdb.org/register
- Free API access (no paid tiers)
- Reasonable usage policy (no hard rate limit)

### GraphQL vs REST
- StashDB uses GraphQL (not REST)
- Flexible queries (request only needed fields)
- Batch requests (fetch multiple scenes in one query)
- Strong typing (GraphQL schema validation)

### Data Quality
- **Community-driven**: User submissions with moderation
- **High quality**: Active moderation team reviews submissions
- **Performer focus**: Best performer metadata (measurements, tattoos, piercings)
- **Scene fingerprinting**: Accurate scene identification via phash

### Fingerprinting (Perceptual Hashing)
- **Algorithm**: Perceptual hash (phash) - detects duplicate/similar videos
- **Generation**: Use `ffmpeg` + perceptual hashing library (e.g., `goimagehash`)
- **Matching**: Compare phash values (Hamming distance < threshold = match)
- **Duration**: Include video duration for additional validation

### Adult Content Isolation (CRITICAL)
- **Database schema**: `qar` schema ONLY (`qar.expeditions`, `qar.voyages`, `qar.crew`, `qar.ports`)
  - `qar.expeditions.metadata_json.stashdb_data` (JSONB)
  - `qar.voyages.metadata_json.stashdb_data` (JSONB)
  - `qar.crew` (dedicated table)
  - `qar.ports` (dedicated table)
  - `qar.expedition_crew` (junction table)
  - `qar.voyage_crew` (junction table)
- **API namespace**: `/api/v1/qar/metadata/stashdb/*` (isolated)
  - `/api/v1/qar/metadata/stashdb/search/voyages`
  - `/api/v1/qar/metadata/stashdb/search/crew`
  - `/api/v1/qar/metadata/stashdb/voyages/{stashdb_id}`
  - `/api/v1/qar/metadata/stashdb/crew/{stashdb_id}`
  - `/api/v1/qar/metadata/stashdb/identify` (fingerprint lookup)
- **Module location**: `internal/content/qar/metadata/stashdb/` (NOT `internal/service/metadata/`)
- **Access control**: Mods/admins see all data for monitoring, regular users see only their library

### Tags & Categories
- **Categories**: Position, Act, Setting, Fetish, Hair Color, Ethnicity, etc.
- **Hierarchical**: Tags have categories (e.g., "Doggy Style" → category "Position")
- **Comprehensive**: 500+ tags covering all aspects
- **User-defined**: Custom tags supported (in addition to StashDB tags)

### Parent/Subsidiary Studios
- **Hierarchy**: Studios can have parent companies (e.g., "Brazzers" → parent "MindGeek")
- **Branding**: Display parent studio branding in UI
- **Filtering**: Filter scenes by parent studio (e.g., all MindGeek content)

### JSONB Storage
- Store full StashDB GraphQL response in `qar.expeditions.metadata_json.stashdb_data`
- Future-proofing: If StashDB adds new fields, they're automatically stored
- Querying: Use PostgreSQL JSONB operators for advanced queries

### Image URLs
- **Direct URLs**: StashDB provides direct image URLs (no authentication required)
- **CDN**: Images served via CDN (fast delivery)
- **Sizes**: Multiple sizes available (thumbnail, medium, full-resolution)
- **Download**: Download and store locally (avoid hotlinking, ensure availability)

### Disambiguation
- **Performer names**: Disambiguation field (e.g., "Performer Name (Studio)" OR "Performer Name (Birthdate)")
- **Use case**: Multiple performers with same name (rare but possible)

### Community Contributions
- **User submissions**: Users can submit new scenes/performers/studios
- **Moderation**: Active moderation team reviews submissions
- **Edits**: Users can propose edits to existing data
- **Trust system**: Established contributors have higher trust levels

### Fallback Strategy
- **StashDB primary**: Use StashDB as primary adult metadata source
- **ThePornDB fallback**: Use ThePornDB if StashDB lacks data (rare)
- **Manual entry**: Allow users to manually add metadata if both sources lack data

### Use Case: Whisparr Integration
- Whisparr uses StashDB for scene identification
- When Whisparr downloads scene → fetch StashDB metadata → update Revenge database
- Unified metadata across Whisparr + Revenge
