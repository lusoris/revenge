# Consolidated TODO — Frontend Readiness + Metadata API (2026-02-08)

## DONE (this session)

- [x] C1: CORS middleware (origin reflection, credentials, preflight) — `ccf9b96c`
- [x] H1: `total` field on all list response schemas (10 bounded + movies/tvshows wrappers)
- [x] H3: 201 response bodies (addPolicy, assignRole)
- [x] M2: Scalar API docs at `/api/docs` + raw spec at `/api/openapi.yaml`
- [x] M3: Webhook security (`webhookAuth` / `X-Webhook-Secret`) — `60518cfe`
- [x] Fix 10 CI lint issues (govet unusedwrite + ineffassign)

---

## P1: Expose Existing Metadata Capabilities (16 endpoints) — DONE

- [x] 5 movie metadata endpoints (credits, images, similar, recommendations, external-ids) — `c86586a5`, `78084fed`
- [x] 6 TV show metadata endpoints (credits, images, content-ratings, external-ids, season/episode images)
- [x] 4 person metadata endpoints (search, details, credits, images)
- [x] Provider list endpoint (`GET /api/v1/metadata/providers`)
- [x] OpenAPI spec + ogen regeneration

---

## P2: Multi-Provider Search — DONE

- [x] `provider` + `language` query params on all 3 search endpoints — `555fa7f6`
- [x] TVDb search via `?provider=tvdb`
- [x] *arr metadata proxying (Radarr/Sonarr lookup API) — `a3f72838`

---

## P3: New Metadata Providers

### Tier 1 — Movies & TV — DONE

- [x] **Fanart.tv** — HD logos, clearart, disc art, banners (movies + TV) — `b83a1313`
- [x] **OMDb** — IMDb/RT/Metacritic ratings + `ExternalRatings` type — `141a9b87`
- [x] **TVmaze** — Free TV metadata (episodes, cast, seasons, images) — `4454051a`
- [ ] **ThePosterDB** — Curated posters (deferred)

### Tier 2 — Anime (TV show providers)

- [x] **AniList** — Anime metadata via GraphQL (search, show, credits, images, external IDs) — priority 45
- [x] **Kitsu** — Anime metadata via JSON:API (search, show, episodes, images, ratings, external IDs) — priority 42
- [x] **AniDB** — Anime metadata via XML HTTP API, title dump search, 0.5 req/s rate limit — priority 44
- [x] **MyAnimeList** — Anime metadata via REST API v2, X-MAL-CLIENT-ID auth — priority 43

### Tier 3 — Scrobbling & Tracking — DONE

- [x] **Trakt** (priority 38) — Movie + TV metadata, extended info, rate-limited — `f64f23d3`
- [x] **Simkl** (priority 36) — Movie + TV + anime metadata, cross-reference IDs — `f64f23d3`
- [x] **Letterboxd** (priority 34) — Movie-only metadata, OAuth2 client credentials — `f64f23d3`

### Tier 4 — Music (when music module lands)

| Provider | Type | Purpose | API | Auth |
|----------|------|---------|-----|------|
| **MusicBrainz** | Music metadata | Artist/album/track + CoverArt Archive | REST | None (1 req/s) |
| **Last.fm** | Music metadata + scrobble | Artist info, similar artists, scrobbling | REST | API key |
| **Spotify** | Music metadata | Audio features, popularity, previews | REST | OAuth 2.0 |
| **Discogs** | Music metadata | Vinyl/physical release data | REST | OAuth 1.0a |

### Tier 5 — Other media (future)

| Provider | Type | Purpose |
|----------|------|---------|
| **OpenLibrary** | Books | Book metadata, covers |
| **Goodreads** | Books | Reviews, ratings |
| **Audible** | Audiobooks | Audiobook metadata |
| **ComicVine** | Comics | Comic metadata |
| **Lidarr** | Music *arr | Music library management |
| **Whisparr** | Adult *arr | Adult content management |

### Provider Architecture

Each new provider follows the existing pattern:
```
internal/service/metadata/providers/{name}/
  ├── client.go     — HTTP client with rate limiting
  ├── provider.go   — Implements Provider/MovieProvider/TVShowProvider/ImageProvider interfaces
  ├── types.go      — API response types
  └── mapping.go    — Map external types → internal metadata types
```

Registered in `internal/service/metadata/metadatafx/module.go` via fx.

---

## P4: Frontend Polish — DONE

- [x] M1: SSE real-time events (`GET /api/v1/events`) — `acb290b3`
- [x] L1: Cookie-based auth (for SvelteKit SSR) — `92150a01`
- [x] L2: CSRF protection (only with cookies) — `92150a01`

---

## Priority Order

1. **P1** — Expose existing service methods (pure boilerplate, high value)
2. **P2** — Multi-provider search + user-facing *arr data (unlocks content discovery)
3. **P3 Tier 1** — Fanart.tv + OMDb + TVmaze (most-requested, straightforward APIs)
4. **P3 Tier 2** — Anime providers (needed for anime libraries)
5. **P3 Tier 3** — Scrobbling (nice-to-have, user engagement)
6. **P4** — SSE events + cookie auth
7. **P3 Tier 4-5** — Music/books/comics (when those content modules exist)

---

## Research Sources

- [Jellyfin Metadata Docs](https://jellyfin.org/docs/general/server/metadata/)
- [Jellyfin Plugin List](https://jellyfin.org/docs/general/server/plugins/)
- [Fanart.tv API](https://fanarttv.docs.apiary.io/)
- [TVmaze API](https://www.tvmaze.com/api)
- [Radarr Metadata System (Skyhook)](https://deepwiki.com/radarr/radarr/3.7-metadata-system)
- [ThePosterDB](https://theposterdb.com/)
- [OMDb API](https://www.omdbapi.com/)

### Providers used by other media servers (Jellyfin/Plex/Emby/Kodi)

Confirmed in use across the ecosystem:
- TMDb, TVDb, Fanart.tv, OMDb — universal across all servers
- AniDB, AniList, Kitsu, MAL — anime-specific, Jellyfin plugins
- MusicBrainz + CoverArt Archive — music metadata
- TVmaze — TV show alternative to TVDb (used by Emby, Kodi, Plex)
- ThePosterDB — curated poster art (community-driven)
- OpenSubtitles — subtitle downloads (not metadata, but related)

### Missing from our design docs (found via research)
- **TVmaze** — free TV metadata API, no auth needed, good TVDb alternative
- **AniDB** — anime episode grouping (Jellyfin has a dedicated plugin)
- **Anisearch** — anime metadata (Jellyfin plugin)
- **CoverArt Archive** — album art (companion to MusicBrainz)
- **OpenSubtitles** — subtitle provider (v2 REST API)
