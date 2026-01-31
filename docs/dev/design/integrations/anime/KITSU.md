# Kitsu Integration

> Modern anime tracking platform with social features

**Priority**: 🟢 LOW (Phase 6 - Anime Module)
**Type**: JSON:API client with OAuth2

## Status

| Dimension | Status | Notes |
| --------- | ------ | ----- |
| Design | ✅ | Comprehensive JSON:API spec, data mapping, OAuth flow |
| Sources | ✅ | API docs, endpoint, OAuth docs, GitHub linked |
| Instructions | ✅ | Detailed implementation checklist |
| Code | 🔴 | |
| Linting | 🔴 | |
| Unit Testing | 🔴 | |
| Integration Testing | 🔴 | |

---

## Overview

Kitsu is a modern, open-source anime/manga tracking platform with strong social features. Revenge uses Kitsu as a tertiary metadata source and sync target for users who prefer Kitsu:
- Anime/manga metadata
- User list synchronization
- Social features (activities, reactions)
- Modern UI/UX focused community

**Integration Points**:
- **JSON:API**: Query anime, manga, users, libraries
- **OAuth2**: User authentication for list sync
- **ID mapping**: Cross-reference with AniList, MAL
- **Rate limiting**: Generous (no documented limits)

---

## Developer Resources

- 📚 **API Docs**: https://kitsu.docs.apiary.io/
- 🔗 **API Endpoint**: https://kitsu.io/api/edge/
- 🔗 **OAuth Docs**: https://kitsu.docs.apiary.io/#introduction/authentication
- 🔗 **GitHub**: https://github.com/hummingbird-me/kitsu-tools

---

## API Details

**Base URL**: `https://kitsu.io/api/edge/`
**Authentication**:
- Public queries: No auth required
- User data: OAuth2 Bearer token
**Rate Limits**: No documented limits (be respectful)
**Free Tier**: Fully free, open-source
**Spec**: JSON:API (https://jsonapi.org/)

### Key Endpoints

| Endpoint | Purpose |
|----------|---------|
| `GET /anime/{id}` | Get anime details |
| `GET /anime?filter[text]={query}` | Search anime by title |
| `GET /anime?filter[season]={season}&filter[seasonYear]={year}` | Seasonal anime |
| `GET /trending/anime` | Get trending anime |
| `GET /users/{id}/library-entries` | Get user's library |
| `POST /library-entries` | Add to user's library |
| `PATCH /library-entries/{id}` | Update library entry |
| `DELETE /library-entries/{id}` | Remove from library |

### JSON:API Features

```bash
# Include related resources
GET /anime/1?include=genres,categories,staff

# Sparse fieldsets
GET /anime/1?fields[anime]=canonicalTitle,synopsis,posterImage

# Filtering
GET /anime?filter[text]=attack%20on%20titan

# Pagination
GET /anime?page[limit]=20&page[offset]=0

# Sorting
GET /anime?sort=-averageRating
```

---

## Data Mapping

### Kitsu → Revenge Mapping

| Kitsu Field | Revenge Field | Notes |
|-------------|---------------|-------|
| `id` | `kitsu_id` | Kitsu identifier |
| `attributes.canonicalTitle` | `title` | Default display title |
| `attributes.titles.en` | `title_en` | English title |
| `attributes.titles.ja_jp` | `title_native` | Japanese title |
| `attributes.synopsis` | `overview` | Description |
| `attributes.episodeCount` | `episode_count` | Total episodes |
| `attributes.episodeLength` | `episode_duration` | Minutes |
| `attributes.status` | `airing_status` | See status mapping |
| `attributes.averageRating` | `kitsu_score` | 0-100 scale |
| `attributes.ratingRank` | `kitsu_rank` | Rating ranking |
| `attributes.popularityRank` | `kitsu_popularity` | Popularity rank |
| `attributes.posterImage.large` | `poster_url` | Poster image |
| `attributes.coverImage.large` | `backdrop_url` | Cover/backdrop |

### Status Mapping

| Kitsu Status | Revenge Status |
|--------------|----------------|
| `current` | `watching` |
| `completed` | `completed` |
| `on_hold` | `on_hold` |
| `dropped` | `dropped` |
| `planned` | `plan_to_watch` |

### Score Mapping

Kitsu uses 0-100 internally (displayed as 0-10 stars with half-stars):
```
revenge_score = kitsu_score  # Already 0-100
```

---

## OAuth2 Flow

Kitsu uses standard OAuth2 (no PKCE required):

```go
// Password grant (for trusted apps)
POST https://kitsu.io/api/oauth/token
Content-Type: application/x-www-form-urlencoded

grant_type=password
username={email}
password={password}

// Authorization code grant (recommended)
// 1. Redirect to authorize
GET https://kitsu.io/api/oauth/authorize?response_type=code&client_id={id}&redirect_uri={uri}

// 2. Exchange code
POST https://kitsu.io/api/oauth/token
grant_type=authorization_code
code={code}
client_id={id}
client_secret={secret}
redirect_uri={uri}

// 3. Refresh token
POST https://kitsu.io/api/oauth/token
grant_type=refresh_token
refresh_token={token}
```

---

## Implementation Checklist

- [ ] **JSON:API Client** (`internal/service/metadata/provider_kitsu.go`)
  - [ ] JSON:API parsing (relationships, includes)
  - [ ] Anime metadata fetching
  - [ ] Anime search by title
  - [ ] Trending/seasonal queries
  - [ ] Error handling & retries

- [ ] **OAuth2 Integration** (`internal/service/oidc/kitsu.go`)
  - [ ] Authorization code flow
  - [ ] Token storage (per user)
  - [ ] Token refresh handling
  - [ ] Account linking

- [ ] **List Sync** (`internal/service/sync/kitsu_sync.go`)
  - [ ] Import user's library
  - [ ] Sync watch progress → Kitsu
  - [ ] Sync ratings → Kitsu
  - [ ] Conflict resolution
  - [ ] Periodic sync job (River)

- [ ] **ID Mapping**
  - [ ] Kitsu ID ↔ AniList ID mapping
  - [ ] Kitsu ID ↔ MAL ID mapping
  - [ ] Use Kitsu's mappings endpoint

---

## Configuration

```yaml
# configs/config.yaml
integrations:
  kitsu:
    enabled: true
    client_id: "${REVENGE_KITSU_CLIENT_ID}"
    client_secret: "${REVENGE_KITSU_CLIENT_SECRET}"
    rate_limit:
      requests_per_second: 10  # Generous
    sync:
      enabled: true
      interval: "24h"
      direction: "import_only"  # Kitsu as tertiary
    use_as_primary: false
```

---

## Database Schema

```sql
-- Kitsu tokens
CREATE TABLE user_kitsu_tokens (
    user_id UUID PRIMARY KEY REFERENCES users(id),
    kitsu_user_id INTEGER NOT NULL,
    kitsu_username VARCHAR(255),
    access_token TEXT NOT NULL,
    refresh_token TEXT NOT NULL,
    token_expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Kitsu sync state
CREATE TABLE kitsu_sync_state (
    user_id UUID PRIMARY KEY REFERENCES users(id),
    last_sync_at TIMESTAMPTZ,
    last_sync_status VARCHAR(20),
    entries_synced INTEGER DEFAULT 0,
    errors JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

---

## ID Mapping via Kitsu

Kitsu provides mappings to other services:

```bash
GET /anime/{id}/mappings
```

Response includes:
- MyAnimeList ID
- AniDB ID
- TheTVDB ID
- AniList ID (sometimes)

Use this to cross-reference anime across services.

---

## Error Handling

| Error Code | Meaning | Action |
|------------|---------|--------|
| 400 | Bad request | Check JSON:API formatting |
| 401 | Unauthorized | Refresh OAuth token |
| 403 | Forbidden | Check permissions |
| 404 | Not found | Mark as unavailable |
| 500+ | Server error | Retry with backoff |

---

## Kitsu Advantages

1. **Open Source**: Backend is open-source, community-driven
2. **No API Key Required**: Public endpoints work without auth
3. **JSON:API Standard**: Well-documented, consistent format
4. **ID Mappings**: Built-in cross-references to MAL, AniDB, etc.
5. **Social Features**: Activities, reactions, groups (if needed)

---

## Kitsu Limitations

1. **Smaller Database**: Fewer entries than AniList/MAL
2. **Slower Updates**: New anime added slower than AniList
3. **Less Community**: Smaller user base
4. **OAuth Required**: For any user-specific operations

---

## Priority in Anime Stack

| Priority | Service | Role |
|----------|---------|------|
| 1st | AniList | Primary metadata + user sync |
| 2nd | MyAnimeList | Fallback metadata, MAL scores |
| 3rd | Kitsu | User sync (if preferred), ID mapping |

---


<!-- SOURCE-BREADCRUMBS-START -->

## Sources & Cross-References

> Auto-generated section linking to external documentation sources

### Cross-Reference Indexes

- [All Sources Index](../../../sources/SOURCES_INDEX.md) - Complete list of external documentation
- [Design ↔ Sources Map](../../../sources/DESIGN_CROSSREF.md) - Which docs reference which sources

<!-- SOURCE-BREADCRUMBS-END -->

<!-- DESIGN-BREADCRUMBS-START -->

## Related Design Docs

> Auto-generated cross-references to related design documentation

**Category**: [Anime](INDEX.md)

### In This Section

- [AniList Integration](ANILIST.md)
- [MyAnimeList (MAL) Integration](MYANIMELIST.md)

### Related Topics

- [Revenge - Architecture v2](../../architecture/01_ARCHITECTURE.md) _Architecture_
- [Revenge - Design Principles](../../architecture/02_DESIGN_PRINCIPLES.md) _Architecture_
- [Revenge - Metadata System](../../architecture/03_METADATA_SYSTEM.md) _Architecture_
- [Revenge - Player Architecture](../../architecture/04_PLAYER_ARCHITECTURE.md) _Architecture_
- [Plugin Architecture Decision](../../architecture/05_PLUGIN_ARCHITECTURE_DECISION.md) _Architecture_

### Indexes

- [Design Index](../../DESIGN_INDEX.md) - All design docs by category/topic
- [Source of Truth](../../00_SOURCE_OF_TRUTH.md) - Package versions and status

<!-- DESIGN-BREADCRUMBS-END -->

## Related Documentation

- [AniList Integration](ANILIST.md)
- [MyAnimeList Integration](MYANIMELIST.md)
- [Scrobbling Overview](../scrobbling/INDEX.md)
