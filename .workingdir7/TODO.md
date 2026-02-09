# Ratings System Overhaul (2026-02-09)

## Problem

Two distinct rating concepts exist but both are broken:

1. **External provider ratings** (IMDb, RT, Metacritic, TMDb, Trakt, etc.) — partially fetched but **never persisted or exposed to the API**
2. **User ratings** (personal scores synced to Trakt etc.) — **does not exist at all**

---

## Current State

### What works
- `ExternalRating` type exists in `metadata/types.go` with `Source`, `Value`, `Score` (0-100)
- OMDb provider fetches IMDb + Rotten Tomatoes + Metacritic into `ExternalRatings`
- AniList, AniDB, MAL, Kitsu add their own scores as `ExternalRatings`
- `VoteAverage` / `VoteCount` fields exist on movie/show metadata types

### What's broken
- **ExternalRatings are never persisted** — no DB column, adapter silently drops them
- **ExternalRatings are never exposed** — no OpenAPI schema, invisible to frontend
- **No cross-provider enrichment** — each provider used independently, TMDb metadata never enriched with OMDb ratings
- **Most providers don't create ExternalRatings for their own score** (TMDb, Trakt, TVDb, TVmaze, Letterboxd only set VoteAverage)
- **Trakt abuses ExternalRatings** for content certifications instead of ratings
- **VoteCount missing** on many providers even when the API has it (TVDb, TVmaze, MAL, Letterboxd)
- **Simkl movie mapping** drops MAL rating (show mapping captures it)
- **No user rating system** — no table, no service, no endpoints

---

## Fix Plan

### P1: Persist & Expose ExternalRatings

- [ ] Add `external_ratings JSONB` column to movies + series tables (new migration)
- [ ] Regenerate sqlc models
- [ ] Update movie/tvshow adapter to persist ExternalRatings (currently dropped in `mapMetadataToMovie`)
- [ ] Add `ExternalRating` schema to OpenAPI spec
- [ ] Add `external_ratings` array to Movie/TVShow response schemas
- [ ] Regenerate ogen

### P2: All Providers Populate ExternalRatings for Their Own Score

Each provider should add itself as an ExternalRating entry alongside setting VoteAverage:

| Provider | Status | Fix |
|----------|--------|-----|
| TMDb | Missing | Add `{Source: "TMDb", Value: "7.5/10", Score: 75}` |
| TVDb | Missing | Add TVDb ExternalRating |
| TVmaze | Missing | Add TVmaze ExternalRating |
| Trakt | Wrong (has certifications) | Replace certification abuse with real Trakt rating |
| Letterboxd | Missing | Add `{Source: "Letterboxd", Value: "4.2/5", Score: 84}` |
| Simkl | Missing (only VoteAverage) | Add Simkl ExternalRating |
| OMDb | Done | Already creates IMDb, RT, Metacritic entries |
| AniList | Done | Already creates AniList entry |
| AniDB | Done | Already creates Permanent + Temporary entries |
| MAL | Done | Already creates MyAnimeList entry |
| Kitsu | Done | Already creates Kitsu entry |

### P3: Fix Missing VoteCount

| Provider | Data Available | Currently Set |
|----------|---------------|---------------|
| TVDb | Yes (score) | No VoteCount |
| TVmaze | Weight field | No |
| MAL | NumScoringUsers | Search only, not metadata |
| Letterboxd | FilmStatistics.Counts.Ratings | Not fetched |
| Kitsu | UserCount | Used, but RatingFrequencies sum would be better |
| Simkl | No direct count | N/A |

### P4: Cross-Provider Rating Enrichment

- [ ] Implement enrichment in `metadata/service.go` (config: `EnableEnrichment: true`)
- [ ] After fetching from primary provider (TMDb), concurrently fetch from OMDb (IMDb ID lookup)
- [ ] Merge ExternalRatings into single slice with dedup by Source
- [ ] Optionally fetch Trakt/Letterboxd community scores if configured

### P5: Fix Provider-Specific Issues

- [ ] Trakt: Stop abusing ExternalRatings for certifications — use proper field or drop them
- [ ] Simkl movies: Add MAL rating (already done for shows, missing for movies)
- [ ] Simkl types: Add `Tmdb`, `Trakt`, `Letterboxd` to Ratings struct (API returns them)
- [ ] Letterboxd: Fetch FilmStatistics for vote count

### P6: User Rating System (new feature)

- [ ] New migration: `user_ratings` table (`user_id`, `media_type`, `media_id`, `rating DECIMAL(3,1)`, `review TEXT`, `rated_at`)
- [ ] Repository: CRUD for user ratings
- [ ] Service: `internal/service/rating/` — rate, unrate, get, list, average
- [ ] Handler: `POST /api/v1/movies/{id}/rate`, `DELETE`, `GET`
- [ ] Handler: `POST /api/v1/shows/{id}/rate`, `DELETE`, `GET`
- [ ] OpenAPI spec: Add rating endpoints + schemas
- [ ] Trakt sync: Bidirectional rating sync when Trakt OAuth configured

---

## Priority Order

1. **P1** — Persist & expose (without this, all provider ratings are invisible)
2. **P2** — All providers create ExternalRatings (data completeness)
3. **P3** — Fix VoteCount (low effort, high value)
4. **P4** — Cross-provider enrichment (TMDb + OMDb ratings together)
5. **P5** — Provider-specific fixes (cleanup)
6. **P6** — User rating system (new feature, bigger scope)
