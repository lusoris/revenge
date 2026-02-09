# Ratings & Sync Overhaul (2026-02-09)

## Problem

Three distinct systems are broken or missing:

1. **External provider ratings** (IMDb, RT, Metacritic, TMDb, Trakt, etc.) — partially fetched but **never persisted or exposed to the API**
2. **User ratings** (personal scores synced to Trakt etc.) — **does not exist at all**
3. **Scrobbling & sync** (Trakt watched states, library sync, Simkl sync) — **not implemented** (design docs exist at `docs/dev/design/planned/integrations/scrobbling/`)

---

## Current State

### What works
- `ExternalRating` type exists in `metadata/types.go` with `Source`, `Value`, `Score` (0-100)
- OMDb provider fetches IMDb + Rotten Tomatoes (Tomatometer only) + Metacritic into `ExternalRatings`
- AniList, AniDB, MAL, Kitsu add their own scores as `ExternalRatings`
- `VoteAverage` / `VoteCount` fields exist on movie/show metadata types
- `movie_watched` table tracks per-user watch progress (progress_seconds, is_completed, watch_count)
- Extensive design docs already exist for Trakt (`TRAKT.md`), Simkl (`SIMKL.md`), Letterboxd (`LETTERBOXD.md`)

### What's broken
- **ExternalRatings are never persisted** — no DB column, adapter silently drops them
- **ExternalRatings are never exposed** — no OpenAPI schema, invisible to frontend
- **No cross-provider enrichment** — each provider used independently, TMDb metadata never enriched with OMDb ratings
- **Most providers don't create ExternalRatings for their own score** (TMDb, Trakt, TVDb, TVmaze, Letterboxd only set VoteAverage)
- **Trakt abuses ExternalRatings** for content certifications instead of ratings
- **VoteCount missing** on many providers even when the API has it (TVDb, TVmaze, MAL, Letterboxd)
- **Simkl movie mapping** drops MAL rating (show mapping captures it)
- **Rotten Tomatoes**: OMDb only provides Tomatometer (critics) — RT Audience Score is NOT available from OMDb
- **No user rating system** — no table, no service, no endpoints
- **No scrobbling** — Trakt/Simkl are metadata-only, no user-level OAuth, no watch state sync
- **Trakt OAuth not per-user** — currently app-level API key only; needs per-user Bearer + refresh token flow
- **No library sync** — Trakt/Simkl can sync what users have in their library/collection, not implemented
- **Playback service doesn't persist** — in-memory only, no scrobble hooks

### Existing infrastructure
- `shared.oidc_user_links` stores encrypted OAuth tokens for SSO — similar pattern needed for Trakt/Simkl
- `shared.users` + `shared.sessions` provide user/auth foundation
- `movie_watched` table exists but playback service doesn't write to it
- v0.8.0 planning already specifies `scrobble_connections` + `scrobble_queue` tables + River jobs

---

## Fix Plan

### P1: Persist & Expose ExternalRatings ✅ DONE

- [x] Add `external_ratings JSONB` column to movies + series tables (migration 000035)
- [x] Update sqlc queries (CreateMovie, UpdateMovie, CreateSeries, UpdateSeries)
- [x] Regenerate sqlc models
- [x] Add `ExternalRating` domain type to movie + tvshow packages
- [x] Update repository params (Create + Update) for both movie + tvshow
- [x] Update repository postgres impl (marshal/unmarshal, CRUD wiring)
- [x] Update service param converters (`movieToUpdateParams`, `seriesToUpdateParams`)
- [x] Update metadata adapters (`mapMetadataToMovie`, `mapMetadataToSeries`)
- [x] Add `ExternalRating` schema to OpenAPI spec
- [x] Add `external_ratings` array to Movie, TVSeries, ContinueWatchingItem, WatchedMovieItem
- [x] Regenerate ogen
- [x] Update API converters (`movieToOgen`, `seriesToOgen`, continue watching, watched items)
- [x] Build + lint pass cleanly

### P2: All Providers Populate ExternalRatings for Their Own Score ✅ DONE

Each provider should add itself as an ExternalRating entry alongside setting VoteAverage:

| Provider | Status | Fix |
|----------|--------|-----|
| TMDb | ✅ Done | Add `{Source: "TMDb", Value: "7.5/10", Score: 75}` |
| TVDb | ✅ Done | Add TVDb ExternalRating |
| TVmaze | ✅ Done | Add TVmaze ExternalRating |
| Trakt | ✅ Done | Replace certification abuse with real Trakt rating |
| Letterboxd | ✅ Done | Add `{Source: "Letterboxd", Value: "4.2/5", Score: 84}` |
| Simkl | ✅ Done | Add Simkl ExternalRating |
| OMDb | Done | Already creates IMDb, RT (Tomatometer), Metacritic entries |
| AniList | Done | Already creates AniList entry |
| AniDB | Done | Already creates Permanent + Temporary entries |
| MAL | Done | Already creates MyAnimeList entry |
| Kitsu | Done | Already creates Kitsu entry |

**RT Audience Score note**: OMDb only has Tomatometer (critics). RT Audience Score is not available from any free API. If we ever want it, we'd need to scrape or use a paid source. For now, Tomatometer from OMDb is what we expose as "Rotten Tomatoes" — the Value already says e.g. "96%" which is the critics score. We should name it clearly: `Source: "Rotten Tomatoes (Tomatometer)"`.

### P3: Fix Missing VoteCount ✅ DONE

| Provider | Data Available | Status |
|----------|---------------|--------|
| TVDb | Score only (no count) | N/A |
| TVmaze | Weight (not vote count) | N/A |
| MAL | NumScoringUsers | Search only |
| Letterboxd | FilmStatistics.Counts.Ratings | ✅ Done (P1) |
| Kitsu | UserCount | Already used |
| Simkl | Ratings.Simkl.Votes | ✅ Done |

### P4: Cross-Provider Rating Enrichment ✅ DONE

- [x] Implement enrichment in `metadata/service.go` (config: `EnableEnrichment: true`)
- [x] After fetching from primary provider (TMDb), concurrently fetch from OMDb (IMDb ID lookup)
- [x] Merge ExternalRatings into single slice with dedup by Source
- [x] Optionally fetch Trakt/Letterboxd community scores if configured (all secondary MovieProvider/TVShowProvider used)

### P5: Fix Provider-Specific Issues ✅ DONE

- [x] Trakt: Stop abusing ExternalRatings for certifications — use proper field or drop them (done in P2)
- [x] Simkl movies: Add MAL rating (already done for shows, missing for movies) (done in P2)
- [x] Simkl types: Add `Tmdb`, `Trakt`, `Letterboxd` to Ratings struct + map as ExternalRatings
- [x] Letterboxd: Fetch FilmStatistics for vote count (already done in P1)

### P6: User Rating System (new feature)

- [ ] New migration: `user_ratings` table (`user_id`, `media_type`, `media_id`, `rating DECIMAL(3,1)`, `review TEXT`, `rated_at`)
- [ ] Repository: CRUD for user ratings
- [ ] Service: `internal/service/rating/` — rate, unrate, get, list, average
- [ ] Handler: `POST /api/v1/movies/{id}/rate`, `DELETE`, `GET`
- [ ] Handler: `POST /api/v1/shows/{id}/rate`, `DELETE`, `GET`
- [ ] OpenAPI spec: Add rating endpoints + schemas
- [ ] Trakt sync: Bidirectional rating sync when Trakt OAuth configured

### P7: Scrobbling & Sync Service (v0.8.0 scope — design already exists)

Design docs: `docs/dev/design/planned/integrations/scrobbling/`

#### Infrastructure
- [ ] New migration: `shared.scrobble_connections` table (user_id, provider, access_token_encrypted, refresh_token_encrypted, token_expires_at, enabled, sync_ratings, sync_watch_status, sync_library, scrobble_threshold_percent)
- [ ] New migration: `shared.scrobble_queue` table (user_id, provider, media_type, media_id, action, payload JSONB, status, retry_count, next_retry_at)
- [ ] Scrobble service: `internal/service/scrobbling/service.go`
- [ ] River background jobs for queue processing + periodic history sync

#### Trakt Integration (per-user OAuth2 Bearer + refresh token)
- [ ] OAuth2 flow: `/api/v1/scrobbling/trakt/connect` → redirect to Trakt → callback with code → exchange for Bearer token
- [ ] Token storage in `scrobble_connections` (encrypted, with refresh token rotation)
- [ ] Scrobble: POST `/scrobble/start`, `/scrobble/pause`, `/scrobble/stop` (playback events)
- [ ] Watch history sync: bidirectional — import Trakt history, export local watches to Trakt
- [ ] Library/collection sync: sync what items exist in user's Trakt collection
- [ ] Rating sync: bidirectional user rating sync
- [ ] Watchlist sync: bidirectional
- [ ] Playback hook: wire `playback.SessionManager` to emit scrobble events

#### Simkl Integration (PIN-based OAuth)
- [ ] PIN auth flow: generate PIN → user enters on simkl.com → poll for token
- [ ] Checkin/scrobble support
- [ ] Watch history + watchlist sync
- [ ] Rating sync

#### Letterboxd (CSV import/export only — no real-time API scrobbling)
- [ ] Import diary, watchlist, ratings from CSV
- [ ] Export diary, watchlist to CSV

---

## Priority Order

1. **P1** — Persist & expose (without this, all provider ratings are invisible)
2. **P2** — All providers create ExternalRatings (data completeness)
3. **P3** — Fix VoteCount (low effort, high value)
4. **P4** — Cross-provider enrichment (TMDb + OMDb ratings together)
5. **P5** — Provider-specific fixes (cleanup)
6. **P6** — User rating system (new feature, bigger scope)
7. **P7** — Scrobbling & sync (v0.8.0 scope, depends on P6 for rating sync)
