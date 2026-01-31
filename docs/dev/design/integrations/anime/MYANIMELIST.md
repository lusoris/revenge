# MyAnimeList (MAL) Integration

<!-- SOURCES: myanimelist, river -->

<!-- DESIGN: integrations/anime, 01_ARCHITECTURE, 02_DESIGN_PRINCIPLES, 03_METADATA_SYSTEM -->


> Legacy anime tracking platform with extensive database


<!-- TOC-START -->

## Table of Contents

- [Status](#status)
- [Overview](#overview)
- [Developer Resources](#developer-resources)
- [API Details](#api-details)
  - [Key Endpoints](#key-endpoints)
  - [Request Headers](#request-headers)
  - [Response Fields](#response-fields)
- [Data Mapping](#data-mapping)
  - [MAL → Revenge Mapping](#mal-revenge-mapping)
  - [Status Mapping](#status-mapping)
  - [Score Mapping](#score-mapping)
- [OAuth2 Flow (PKCE)](#oauth2-flow-pkce)
- [Implementation Checklist](#implementation-checklist)
- [Configuration](#configuration)
- [Database Schema](#database-schema)
- [Rate Limiting Strategy](#rate-limiting-strategy)
- [Error Handling](#error-handling)
- [MAL vs AniList Priority](#mal-vs-anilist-priority)
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
**Type**: REST API client with OAuth2

## Status

| Dimension | Status | Notes |
|-----------|--------|-------|
| Design | ✅ | Comprehensive REST API spec, PKCE flow, data mapping |
| Sources | ✅ | API docs, endpoint, OAuth portal linked |
| Instructions | ✅ | Detailed implementation checklist |
| Code | 🔴 |  |
| Linting | 🔴 |  |
| Unit Testing | 🔴 |  |
| Integration Testing | 🔴 |  |---

## Overview

MyAnimeList is the oldest and most comprehensive anime/manga database with millions of users. While AniList is the primary provider, MAL support is essential for:
- User list import/migration from MAL
- Cross-referencing MAL IDs
- Additional metadata (MAL scores, popularity rankings)
- Users who prefer MAL ecosystem

**Integration Points**:
- **REST API v2**: Query anime, manga, user lists
- **OAuth2**: User authentication for list sync
- **MAL ID mapping**: Cross-reference with AniList, Kitsu
- **Rate limiting**: 5 requests per second

---

## Developer Resources

- 📚 **API Docs**: https://myanimelist.net/apiconfig/references/api/v2
- 🔗 **API Endpoint**: https://api.myanimelist.net/v2/
- 🔗 **OAuth Portal**: https://myanimelist.net/apiconfig
- 🔗 **Status**: https://status.myanimelist.net/

---

## API Details

**Base URL**: `https://api.myanimelist.net/v2/`
**Authentication**:
- Public queries: Client ID header (`X-MAL-CLIENT-ID`)
- User data: OAuth2 Bearer token
**Rate Limits**: ~5 requests per second (undocumented, be conservative)
**Free Tier**: Available (OAuth app registration required)
**PKCE Required**: Yes, for OAuth2 flow

### Key Endpoints

| Endpoint | Purpose |
|----------|---------|
| `GET /anime/{anime_id}` | Get anime details |
| `GET /anime?q={query}` | Search anime by title |
| `GET /anime/ranking` | Get top anime rankings |
| `GET /anime/season/{year}/{season}` | Get seasonal anime |
| `GET /users/@me/animelist` | Get user's anime list |
| `PATCH /anime/{anime_id}/my_list_status` | Update anime in user's list |
| `DELETE /anime/{anime_id}/my_list_status` | Remove anime from user's list |

### Request Headers

```http
X-MAL-CLIENT-ID: {client_id}
Authorization: Bearer {access_token}
```

### Response Fields

Request specific fields with `fields` parameter:
```
?fields=id,title,main_picture,synopsis,mean,rank,popularity,num_episodes,status,genres,studios
```

---

## Data Mapping

### MAL → Revenge Mapping

| MAL Field | Revenge Field | Notes |
|-----------|---------------|-------|
| `id` | `mal_id` | MyAnimeList identifier |
| `title` | `title` | Default title (usually romaji) |
| `alternative_titles.en` | `title_en` | English title |
| `alternative_titles.ja` | `title_native` | Japanese title |
| `synopsis` | `overview` | Description |
| `num_episodes` | `episode_count` | Total episodes |
| `average_episode_duration` | `episode_duration` | Seconds |
| `status` | `airing_status` | See status mapping |
| `mean` | `mal_score` | 0-10 scale |
| `rank` | `mal_rank` | Overall ranking |
| `popularity` | `mal_popularity` | Popularity rank |
| `genres[].name` | `genres[]` | Genre strings |
| `studios[].name` | `studios[]` | Studio names |
| `main_picture.large` | `poster_url` | Poster image |

### Status Mapping

| MAL Status | Revenge Status |
|------------|----------------|
| `watching` | `watching` |
| `completed` | `completed` |
| `on_hold` | `on_hold` |
| `dropped` | `dropped` |
| `plan_to_watch` | `plan_to_watch` |

### Score Mapping

MAL uses 0-10 scale, Revenge normalizes to 0-100:
```
revenge_score = mal_score * 10
```

---

## OAuth2 Flow (PKCE)

MAL requires PKCE (Proof Key for Code Exchange):

```go
// 1. Generate code verifier (43-128 chars)
codeVerifier := generateRandomString(128)

// 2. Generate code challenge (plain or S256)
codeChallenge := codeVerifier // plain method

// 3. Authorization URL
authURL := fmt.Sprintf(
    "https://myanimelist.net/v1/oauth2/authorize?response_type=code&client_id=%s&code_challenge=%s&code_challenge_method=plain&redirect_uri=%s",
    clientID, codeChallenge, redirectURI,
)

// 4. Exchange code for tokens
tokenURL := "https://myanimelist.net/v1/oauth2/token"
// POST with: grant_type=authorization_code, code, code_verifier, client_id, redirect_uri
```

---

## Implementation Checklist

- [ ] **REST Client** (`internal/service/metadata/provider_mal.go`)
  - [ ] HTTP client with rate limiting
  - [ ] Anime metadata fetching
  - [ ] Anime search by title
  - [ ] Seasonal anime queries
  - [ ] Ranking queries
  - [ ] Error handling & retries

- [ ] **OAuth2 Integration** (`internal/service/oidc/mal.go`)
  - [ ] PKCE flow implementation
  - [ ] Token storage (per user)
  - [ ] Token refresh handling
  - [ ] Account linking

- [ ] **List Sync** (`internal/service/sync/mal_sync.go`)
  - [ ] Import user's anime list
  - [ ] Sync watch progress → MAL
  - [ ] Sync ratings → MAL (convert 0-100 to 0-10)
  - [ ] Conflict resolution
  - [ ] Periodic sync job (River)

- [ ] **ID Mapping**
  - [ ] MAL ID ↔ AniList ID mapping
  - [ ] MAL ID ↔ Kitsu ID mapping
  - [ ] Fallback search by title

---

## Configuration

```yaml
# configs/config.yaml
integrations:
  myanimelist:
    enabled: true
    client_id: "${REVENGE_MAL_CLIENT_ID}"
    client_secret: "${REVENGE_MAL_CLIENT_SECRET}"
    rate_limit:
      requests_per_second: 3  # Conservative
    sync:
      enabled: true
      interval: "12h"
      direction: "import_only"  # MAL as secondary
    use_as_primary: false  # AniList is primary
```

---

## Database Schema

```sql
-- MAL tokens (separate from AniList)
CREATE TABLE user_mal_tokens (
    user_id UUID PRIMARY KEY REFERENCES users(id),
    mal_user_id INTEGER NOT NULL,
    mal_username VARCHAR(255),
    access_token TEXT NOT NULL,
    refresh_token TEXT NOT NULL,
    token_expires_at TIMESTAMPTZ NOT NULL,
    code_verifier TEXT,  -- For PKCE refresh
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- MAL sync state
CREATE TABLE mal_sync_state (
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

MAL is stricter than AniList. Strategy:

1. **Conservative limit**: 3 requests/second (documented is ~5)
2. **Request batching**: Use `fields` parameter to get all data in one request
3. **Caching**: Cache metadata for 48 hours (MAL data changes less frequently)
4. **Retry with backoff**: On errors, exponential backoff (2s, 4s, 8s...)

---

## Error Handling

| Error Code | Meaning | Action |
|------------|---------|--------|
| 400 | Bad request | Check parameters, log error |
| 401 | Unauthorized | Refresh OAuth token |
| 403 | Forbidden | Check app permissions |
| 404 | Anime not found | Mark as unavailable |
| 429 | Rate limited | Backoff 60 seconds |
| 500+ | Server error | Retry with backoff |

---

## MAL vs AniList Priority

When both sources available:

1. **Primary metadata**: AniList (better API, more frequent updates)
2. **Fallback metadata**: MAL (if AniList missing)
3. **Scores**: Display both separately (MAL score, AniList score)
4. **User sync**: User chooses preferred platform

---


## Related Documentation

- [AniList Integration](ANILIST.md)
- [Kitsu Integration](KITSU.md)
- [Scrobbling Overview](../scrobbling/INDEX.md)
