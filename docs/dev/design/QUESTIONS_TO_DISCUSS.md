# Questions & Gaps to Discuss

> Temporary file for tracking discrepancies, gaps, and questions that need resolution

**Created**: 2026-01-30
**Status**: 🔄 Active collection

---

## Critical Discrepancies Found

### Package Version Mismatches

| Package | go.mod (actual) | SOURCES.yaml | SOURCE_OF_TRUTH.md | Resolution |
|---------|-----------------|--------------|---------------------|------------|
| otter | v1.2.4 | v2.3.0 | v1.2.4 ✅ | Update SOURCES.yaml |
| typesense-go | v4.0.0-alpha2 ❌ | v3.2.0 | v3.2.0 ✅ | **CODE FIX NEEDED**: go.mod must use v3.2.0 stable |
| resty | go-resty/resty/v2 v2.17.1 ✅ | resty v3.0.0-b6 | v2.17.1 ✅ | Update SOURCES.yaml |

### Resolved (2026-01-30)

1. **Typesense Go Client**: ✅ DECISION MADE
   - **Answer**: Use stable v3.2.0, NOT alpha v4
   - **Principle**: "Bleeding edge stable" = newest STABLE version, never alphas/RCs
   - **Action needed**: Fix go.mod when doing code work (change v4 alpha → v3 stable)
   - **Verified**: [typesense-go releases](https://github.com/typesense/typesense-go/releases) - v3.2.0 is latest stable (March 2025)

2. **Resty HTTP Client**: ✅ CORRECT
   - go.mod uses v2 stable which is correct
   - SOURCES.yaml mentioned v3 beta incorrectly - needs update

---

## Missing from SOURCE_OF_TRUTH

### Not Yet Documented
- [ ] Full API endpoint inventory
- [ ] Complete database table list
- [ ] All environment variables with defaults
- [ ] OpenAPI spec file locations
- [ ] Migration file naming conventions
- [ ] Test coverage requirements
- [ ] CI/CD pipeline stages

### Need Design Docs
- [ ] Fingerprint service design doc
- [ ] Grants service design doc
- [ ] Search service design doc
- [ ] Analytics service design doc
- [ ] Notification service design doc

---

## Outdated References Found

### Old Namespace 'c' → 'qar' ✅ COMPLETE

**All Fixed (2026-01-30):**
- [x] INDEX.md, docs/dev/INDEX.md, MODULE_IMPLEMENTATION_TODO.md
- [x] ARCHITECTURE.md - updated to `qar` schema + SOT reference
- [x] WHISPARR.md - all `/api/v1/c/` → `/api/v1/legacy/`
- [x] POSTGRESQL.md - updated to `qar` schema
- [x] features/adult/INDEX.md - SOT reference added
- [x] REQUEST_SYSTEM.md - QAR terminology (expedition, voyage)
- [x] RBAC_CASBIN.md - `qar` schema reference
- [x] All external/adult integrations (6 files) - `/api/v1/legacy/`
- [x] All wiki/adult integrations (3 files) - `/api/v1/legacy/`
- [x] NEWS_SYSTEM.md, WIKI_SYSTEM.md
- [x] WHISPARR_STASHDB_SCHEMA.md, WHISPARR_V3_ANALYSIS.md
- [x] PLUGIN_ARCHITECTURE_DECISION.md

### Broken Internal Links (RESOLVED)

All broken links have been updated to point to SOURCE_OF_TRUTH.md sections:

| File | Old Link | New Link |
|------|----------|----------|
| ARCHITECTURE.md | `PROJECT_STRUCTURE.md` | `SOURCE_OF_TRUTH.md#project-structure` |
| DATA_RECONCILIATION.md | `RIVER_JOBS.md` | `SOURCE_OF_TRUTH.md#river-job-queue-patterns` |
| SKIP_INTRO.md | `RIVER_JOBS.md` | `SOURCE_OF_TRUTH.md#river-job-queue-patterns` |
| TRICKPLAY.md | `RIVER_JOBS.md` | `SOURCE_OF_TRUTH.md#river-job-queue-patterns` |
| NEWS_SYSTEM.md | `RIVER_JOBS.md` | `SOURCE_OF_TRUTH.md#river-job-queue-patterns` |

**Status:** RESOLVED - Links now point to SOT master sections

### Unreferenced Documentation

- [ ] GALLERY_MODULE.md - not in features/adult/INDEX.md
- [ ] STASH.md - not in metadata/adult/INDEX.md

---

## Questions for Project Owner

1. **Database Strategy**:
   - Is SQLite truly needed for single-user deployments?
   - What's the priority for dual DB support vs other features?

2. **Typesense Version**:
   - v4 alpha in go.mod - intentional bleeding edge?
   - Should we pin to v3.2.0 stable?

3. **Package Update Policy**:
   - How aggressive on updates? (bleeding edge vs 1 version behind)
   - Who monitors changelogs/breaking changes?

4. **Documentation Priority**:
   - Should we complete SOURCE_OF_TRUTH first or start scaffolding?
   - User clarified: docs first, no scaffolding until consistent ✅

---

## Live Docs Verification Needed

These need to be verified against actual package/API documentation:

| Source | Last Verified | Action |
|--------|---------------|--------|
| Go 1.25.6 release notes | Never | Check new features |
| pgx v5.8.0 changelog | Never | Breaking changes? |
| River v0.30.2 changelog | Never | Breaking changes? |
| ogen v1.18.0 changelog | Never | Breaking changes? |
| rueidis v1.0.71 | Never | Breaking changes? |
| Typesense API v30 | Never | Verify latest |
| StashDB GraphQL schema | Never | Schema changes? |

---

## Resolution Log

| Date | Item | Resolution |
|------|------|------------|
| 2026-01-30 | otter version | SOURCES.yaml was wrong, go.mod is correct |
| 2026-01-30 | resty version | SOURCES.yaml references v3 beta, we use stable v2 |
| | | |

---

## Document Naming Convention Proposal

**Problem:** Current docs lack numbering and categorization in filenames, making order and relationships unclear.

**Proposed Convention:**
```
[NN]-[CATEGORY]-[name].md

Examples:
01-arch-overview.md          # Architecture overview (read first)
02-arch-principles.md        # Design principles
03-arch-data-flow.md         # Data flow patterns
10-svc-auth.md               # Auth service
11-svc-user.md               # User service
20-mod-movie.md              # Movie module
21-mod-tvshow.md             # TV Show module
30-int-tmdb.md               # TMDb integration
31-int-radarr.md             # Radarr integration
90-ref-api-endpoints.md      # Reference: API endpoints
91-ref-db-tables.md          # Reference: Database tables
```

**Categories:**
- `arch` - Architecture & design
- `svc` - Services
- `mod` - Content modules
- `int` - Integrations
- `ops` - Operations
- `ref` - Reference tables

**Benefits:**
- Clear reading order (numbers)
- Category visible in filename
- Alphabetical sorting = logical order
- Easy to identify doc purpose

**Decision needed:** Adopt this convention? Rename existing files?

### Detailed Renaming Plan

**Phase 1: Core Architecture (00-09)**
```
architecture/
├── 00-arch-source-of-truth.md      # SOURCE_OF_TRUTH.md (master)
├── 01-arch-overview.md             # ARCHITECTURE.md
├── 02-arch-principles.md           # DESIGN_PRINCIPLES.md
├── 03-arch-player.md               # PLAYER_ARCHITECTURE.md
```

**Phase 2: Technical Infrastructure (10-19)**
```
technical/
├── 10-tech-stack.md                # TECH_STACK.md
├── 11-tech-configuration.md        # CONFIGURATION.md
├── 12-tech-audio-streaming.md      # AUDIO_STREAMING.md
├── 13-tech-offloading.md           # OFFLOADING.md
```

**Phase 3: Content Modules (20-39)**
```
features/video/
├── 20-mod-movie.md                 # MOVIE_MODULE.md
├── 21-mod-tvshow.md                # TVSHOW_MODULE.md

features/music/
├── 22-mod-music.md                 # MUSIC_MODULE.md

features/audiobooks/
├── 23-mod-audiobook.md             # (new)

features/podcasts/
├── 24-mod-podcast.md               # PODCASTS.md

features/photos/
├── 25-mod-photos.md                # PHOTOS_LIBRARY.md

features/comics/
├── 26-mod-comics.md                # COMICS_MODULE.md

features/livetv/
├── 27-mod-livetv.md                # LIVE_TV_DVR.md

features/adult/ (qar namespace)
├── 28-qar-performers.md            # PERFORMERS.md
├── 29-qar-scenes.md                # SCENES.md
├── 2A-qar-studios.md               # STUDIOS.md
├── 2B-qar-galleries.md             # GALLERY_MODULE.md
├── 2C-qar-reconciliation.md        # DATA_RECONCILIATION.md
```

**Phase 4: Shared Features (40-49)**
```
features/shared/
├── 40-feat-access-controls.md      # ACCESS_CONTROLS.md
├── 41-feat-analytics.md            # ANALYTICS_SERVICE.md
├── 42-feat-client-support.md       # CLIENT_SUPPORT.md
├── 43-feat-content-rating.md       # CONTENT_RATING.md
├── 44-feat-i18n.md                 # I18N.md
├── 45-feat-news.md                 # NEWS_SYSTEM.md
├── 46-feat-notifications.md        # NOTIFICATIONS.md
├── 47-feat-scrobbling.md           # SCROBBLING.md
├── 48-feat-ticketing.md            # TICKETING_SYSTEM.md
├── 49-feat-voice-control.md        # VOICE_CONTROL.md
```

**Phase 5: Playback Features (50-59)**
```
features/playback/
├── 50-play-media-enhancements.md   # MEDIA_ENHANCEMENTS.md
├── 51-play-skip-intro.md           # SKIP_INTRO.md
├── 52-play-syncplay.md             # SYNCPLAY.md
├── 53-play-trickplay.md            # TRICKPLAY.md
├── 54-play-watch-next.md           # WATCH_NEXT_CONTINUE_WATCHING.md
├── 55-play-release-calendar.md     # RELEASE_CALENDAR.md
```

**Phase 6: Integrations (60-79)**
```
integrations/metadata/video/
├── 60-int-tmdb.md                  # TMDB.md
├── 61-int-tvdb.md                  # THETVDB.md
├── 62-int-omdb.md                  # OMDB.md

integrations/servarr/
├── 63-int-radarr.md                # (exists in metadata)
├── 64-int-sonarr.md                # SONARR.md
├── 65-int-lidarr.md                # LIDARR.md

integrations/metadata/adult/
├── 70-qar-stashdb.md               # STASHDB.md
├── 71-qar-tpdb.md                  # THEPORNDB.md
├── 72-qar-whisparr.md              # WHISPARR.md
```

**Phase 7: Operations (80-89)**
```
operations/
├── 80-ops-setup.md                 # SETUP.md
├── 81-ops-gitflow.md               # GITFLOW.md
├── 82-ops-db-healing.md            # DATABASE_AUTO_HEALING.md
├── 83-ops-best-practices.md        # BEST_PRACTICES.md
├── 84-ops-reverse-proxy.md         # REVERSE_PROXY.md
```

**Phase 8: References (90-99)**
```
├── 90-ref-api-endpoints.md         # API reference (generate from OpenAPI)
├── 91-ref-db-schemas.md            # Database schemas
├── 92-ref-env-vars.md              # Environment variables
├── 93-ref-error-codes.md           # Error codes
├── 99-ref-glossary.md              # Terms and definitions
```

**INDEX.md files:** Keep as-is (no numbering), they serve as folder navigation.

**Execution:** This is a significant refactor. Recommend doing after initial MVP to avoid constant link updates during active development.

---

## Design Strategy Requirements (MANDATORY)

**Principle:** All advanced design patterns and coding strategies MUST be documented in SOURCE_OF_TRUTH.md BEFORE implementation.

**Required in SOT:**
- [ ] All performance patterns (caching, pooling, batching)
- [ ] All resilience patterns (circuit breaker, retry, fallback)
- [ ] All security patterns (auth, RBAC, isolation)
- [ ] All async patterns (jobs, queues, workers)
- [ ] All data patterns (transactions, consistency, partitioning)
- [ ] All API patterns (versioning, errors, pagination)

**Why:**
- No "write first, fix later" approach
- Prevents wasted time on reiteration
- Every implementation inherits patterns from SOT
- Consistency across all modules

**Current gaps to add to SOT:**
- [ ] Error handling patterns (Go errors, API errors)
- [ ] Testing patterns (unit, integration, mocks)
- [ ] Logging patterns (slog, structured, levels)
- [ ] Metrics patterns (Prometheus, OTel)
- [ ] Validation patterns (input, business rules)
- [ ] Pagination patterns (cursor, offset)

---

## Notes

- This file is temporary - items should be resolved and moved to appropriate docs
- When resolved, update SOURCE_OF_TRUTH.md and remove from here
- Questions for owner should be asked and answers documented in design docs
