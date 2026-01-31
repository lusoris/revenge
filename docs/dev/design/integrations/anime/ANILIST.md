# AniList Integration

<!-- SOURCES: anilist, anilist-graphql, go-blurhash, river -->

<!-- DESIGN: integrations/anime, 01_ARCHITECTURE, 02_DESIGN_PRINCIPLES, 03_METADATA_SYSTEM -->


> Primary metadata and tracking provider for anime and manga


<!-- TOC-START -->

## Table of Contents

- [Status](#status)
- [Overview](#overview)
- [Developer Resources](#developer-resources)
- [API Details](#api-details)
  - [Key Queries](#key-queries)
  - [Key Mutations](#key-mutations)
- [Data Mapping](#data-mapping)
  - [AniList → Revenge Mapping](#anilist-revenge-mapping)
  - [Status Mapping](#status-mapping)
- [Implementation Checklist](#implementation-checklist)
- [Configuration](#configuration)
- [Database Schema](#database-schema)
- [Rate Limiting Strategy](#rate-limiting-strategy)
- [Error Handling](#error-handling)
- [Sources & Cross-References](#sources-cross-references)
  - [Cross-Reference Indexes](#cross-reference-indexes)
  - [Referenced Sources](#referenced-sources)
- [Related Design Docs](#related-design-docs)
  - [In This Section](#in-this-section)
  - [Related Topics](#related-topics)
  - [Indexes](#indexes)
- [Related Documentation](#related-documentation)

<!-- TOC-END -->

**Priority**: 🟡 MEDIUM (Phase 6 - Anime Module)
**Type**: GraphQL API client with OAuth

## Status

| Dimension | Status | Notes |
|-----------|--------|-------|
| Design | ✅ | Comprehensive GraphQL API spec, data mapping, database schema |
| Sources | ✅ | API docs, GraphQL endpoint, OAuth portal linked |
| Instructions | ✅ | Detailed implementation checklist |
| Code | 🔴 |  |
| Linting | 🔴 |  |
| Unit Testing | 🔴 |  |
| Integration Testing | 🔴 |  |---

## Overview

AniList is a modern anime/manga tracking and discovery platform with a comprehensive GraphQL API. Revenge uses AniList as the primary metadata source for anime content:
- Anime metadata (title, synopsis, staff, studios, genres)
- Manga metadata (if manga module added)
- User list tracking (watching, completed, plan to watch)
- Scoring and progress synchronization
- Airing schedule information

**Integration Points**:
- **GraphQL API**: Query anime, manga, characters, staff
- **OAuth2**: User authentication for list sync
- **Webhooks**: N/A (polling required)
- **Rate limiting**: 90 requests per minute

---

## Developer Resources

- 📚 **API Docs**: https://anilist.gitbook.io/anilist-apiv2-docs/
- 🔗 **GraphQL Endpoint**: https://graphql.anilist.co
- 🔗 **OAuth Portal**: https://anilist.co/settings/developer
- 🔗 **GraphQL Playground**: https://anilist.co/graphiql
- 🔗 **GitHub Examples**: https://github.com/AniList/ApiV2-GraphQL-Docs

---

## API Details

**Endpoint**: `https://graphql.anilist.co`
**Authentication**:
- Public queries: No auth required
- User data: OAuth2 Bearer token
**Rate Limits**: 90 requests per minute
**Free Tier**: Available (OAuth app registration required)
**i18n Support**: `titleLanguage` preference (ROMAJI, ENGLISH, NATIVE)

### Key Queries

```graphql
# Search anime by title
query {
  Media(search: "Attack on Titan", type: ANIME) {
    id
    title { romaji english native }
    description
    episodes
    status
    genres
    averageScore
    coverImage { large medium }
    bannerImage
    studios { nodes { name } }
  }
}

# Get anime by AniList ID
query {
  Media(id: 16498, type: ANIME) {
    id
    idMal
    title { romaji english native }
    # ... fields
  }
}

# Get user's anime list
query {
  MediaListCollection(userId: $userId, type: ANIME) {
    lists {
      name
      entries {
        media { id title { romaji } }
        status
        progress
        score
      }
    }
  }
}
```

### Key Mutations

```graphql
# Update anime progress
mutation {
  SaveMediaListEntry(mediaId: $mediaId, progress: $progress, status: $status) {
    id
    progress
    status
  }
}

# Delete entry
mutation {
  DeleteMediaListEntry(id: $entryId) {
    deleted
  }
}
```

---

## Data Mapping

### AniList → Revenge Mapping

| AniList Field | Revenge Field | Notes |
|---------------|---------------|-------|
| `id` | `anilist_id` | Primary AniList identifier |
| `idMal` | `mal_id` | MyAnimeList cross-reference |
| `title.romaji` | `title` | Default display title |
| `title.english` | `title_en` | English title |
| `title.native` | `title_native` | Japanese title |
| `description` | `overview` | HTML-encoded, needs sanitization |
| `episodes` | `episode_count` | Total episode count |
| `duration` | `episode_duration` | Minutes per episode |
| `status` | `airing_status` | FINISHED, RELEASING, NOT_YET_RELEASED, CANCELLED |
| `genres` | `genres[]` | Array of genre strings |
| `averageScore` | `anilist_score` | 0-100 scale |
| `coverImage.large` | `poster_url` | Poster image |
| `bannerImage` | `backdrop_url` | Banner/backdrop image |
| `studios.nodes[]` | `studios[]` | Production studios |
| `staff.edges[]` | `staff[]` | Directors, writers, etc. |

### Status Mapping

| AniList Status | Revenge Status |
|----------------|----------------|
| `CURRENT` | `watching` |
| `COMPLETED` | `completed` |
| `PAUSED` | `on_hold` |
| `DROPPED` | `dropped` |
| `PLANNING` | `plan_to_watch` |
| `REPEATING` | `rewatching` |

---

## Implementation Checklist

- [ ] **GraphQL Client** (`internal/service/metadata/provider_anilist.go`)
  - [ ] GraphQL client setup (gqlgen or manual)
  - [ ] Anime metadata fetching
  - [ ] Anime search by title
  - [ ] Character/staff data
  - [ ] Rate limiting (90 req/min)
  - [ ] Error handling & retries

- [ ] **OAuth2 Integration** (`internal/service/oidc/anilist.go`)
  - [ ] OAuth2 authorization flow
  - [ ] Token storage (per user)
  - [ ] Token refresh handling
  - [ ] Account linking

- [ ] **List Sync** (`internal/service/sync/anilist_sync.go`)
  - [ ] Import user's anime list
  - [ ] Sync watch progress → AniList
  - [ ] Sync ratings → AniList
  - [ ] Conflict resolution (local vs remote)
  - [ ] Periodic sync job (River)

- [ ] **Image Handling**
  - [ ] Download cover images
  - [ ] Download banner images
  - [ ] Generate Blurhash placeholders
  - [ ] Image caching

---

## Configuration

```yaml
# configs/config.yaml
integrations:
  anilist:
    enabled: true
    client_id: "${REVENGE_ANILIST_CLIENT_ID}"
    client_secret: "${REVENGE_ANILIST_CLIENT_SECRET}"
    rate_limit:
      requests_per_minute: 90
    sync:
      enabled: true
      interval: "6h"
      direction: "bidirectional"  # import_only, export_only, bidirectional
    default_title_language: "romaji"  # romaji, english, native
```

---

## Database Schema

```sql
-- Anime external IDs
CREATE TABLE anime_external_ids (
    anime_id UUID PRIMARY KEY REFERENCES anime(id),
    anilist_id INTEGER UNIQUE,
    mal_id INTEGER,
    kitsu_id INTEGER,
    thetvdb_id INTEGER,
    imdb_id VARCHAR(20),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- User AniList tokens
CREATE TABLE user_anilist_tokens (
    user_id UUID PRIMARY KEY REFERENCES users(id),
    anilist_user_id INTEGER NOT NULL,
    access_token TEXT NOT NULL,
    refresh_token TEXT,
    token_expires_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Sync state tracking
CREATE TABLE anilist_sync_state (
    user_id UUID PRIMARY KEY REFERENCES users(id),
    last_sync_at TIMESTAMPTZ,
    last_sync_status VARCHAR(20),
    entries_synced INTEGER DEFAULT 0,
    errors JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

---

## Rate Limiting Strategy

AniList allows 90 requests per minute. Strategy:

1. **Token bucket**: 90 tokens, refills 1.5 tokens/second
2. **Batch queries**: Use GraphQL to fetch multiple items in single request
3. **Caching**: Cache metadata for 24 hours (configurable)
4. **Retry with backoff**: On 429 errors, exponential backoff (1s, 2s, 4s...)

---

## Error Handling

| Error Code | Meaning | Action |
|------------|---------|--------|
| 400 | Invalid query | Log error, check query syntax |
| 401 | Unauthorized | Refresh OAuth token |
| 404 | Media not found | Mark as unavailable |
| 429 | Rate limited | Backoff and retry |
| 500 | Server error | Retry with exponential backoff |

---


## Related Documentation

- [MyAnimeList Integration](MYANIMELIST.md)
- [Kitsu Integration](KITSU.md)
- [Scrobbling Overview](../scrobbling/INDEX.md)
