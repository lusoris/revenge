# Documentation Restructuring - Step-by-Step Execution Plan

> Complete action plan for documentation restructuring (4-6 days estimated)

**Status**: 🔴 READY TO EXECUTE
**Created**: 2026-01-28
**Total Steps**: 120 discrete operations

---

## Pre-Flight Checklist

- [x] Git on develop branch
- [x] No uncommitted changes
- [x] Documentation analysis complete
- [ ] Create backup branch: `git checkout -b docs-restructure-backup`

---

## Phase 1: Split Massive Files (40 operations, 2-3 days)

### Step 1.1: Split EXTERNAL_INTEGRATIONS_TODO.md (1445 lines → 66 files)

**Directory Structure**:
```
docs/integrations/
├── INDEX.md
├── servarr/
│   ├── INDEX.md
│   ├── RADARR.md
│   ├── SONARR.md
│   ├── LIDARR.md
│   ├── WHISPARR.md
│   └── READARR.md
├── metadata/
│   ├── INDEX.md
│   ├── video/
│   │   ├── INDEX.md
│   │   ├── TMDB.md
│   │   ├── THETVDB.md
│   │   ├── OMDB.md
│   │   └── THEPOSTERDB.md
│   ├── music/
│   │   ├── INDEX.md
│   │   ├── MUSICBRAINZ.md
│   │   ├── LASTFM.md
│   │   ├── SPOTIFY.md
│   │   └── DISCOGS.md
│   ├── books/
│   │   ├── INDEX.md
│   │   ├── GOODREADS.md
│   │   ├── OPENLIBRARY.md
│   │   ├── AUDIBLE.md
│   │   └── HARDCOVER.md
│   ├── comics/
│   │   ├── INDEX.md
│   │   ├── COMICVINE.md
│   │   ├── MARVEL_API.md
│   │   └── GRAND_COMICS_DATABASE.md
│   └── adult/
│       ├── INDEX.md
│       ├── STASHDB.md
│       ├── THEPORNDB.md
│       └── STASH.md
├── scrobbling/
│   ├── INDEX.md
│   ├── TRAKT.md
│   ├── LASTFM_SCROBBLE.md
│   ├── LISTENBRAINZ.md
│   ├── LETTERBOXD.md
│   └── SIMKL.md
├── anime/
│   ├── INDEX.md
│   ├── ANILIST.md
│   ├── MYANIMELIST.md
│   └── KITSU.md
├── auth/
│   ├── INDEX.md
│   ├── AUTHELIA.md
│   ├── AUTHENTIK.md
│   ├── KEYCLOAK.md
│   └── GENERIC_OIDC.md
├── request/
│   ├── INDEX.md
│   ├── OVERSEERR.md
│   └── JELLYSEERR.md
├── audiobook/
│   ├── INDEX.md
│   └── AUDIOBOOKSHELF.md
├── transcoding/
│   ├── INDEX.md
│   └── BLACKBEARD.md
├── livetv/
│   ├── INDEX.md
│   ├── TVHEADEND.md
│   └── NEXTPVR.md
├── casting/
│   ├── INDEX.md
│   ├── CHROMECAST.md
│   └── DLNA.md
└── infrastructure/
    ├── INDEX.md
    ├── POSTGRESQL.md
    ├── DRAGONFLY.md
    ├── TYPESENSE.md
    └── RIVER.md
```

**Operations** (18 operations):
1. Create `docs/integrations/` directory
2. Create subdirectories: servarr/, metadata/, scrobbling/, anime/, auth/, request/, audiobook/, transcoding/, livetv/, casting/, infrastructure/
3. Create metadata subdirectories: video/, music/, books/, comics/, adult/
4. Read EXTERNAL_INTEGRATIONS_TODO.md sections
5. Extract Servarr services (lines 17-146) → 5 files
6. Extract Video metadata (lines 147-280) → 4 files
7. Extract Music metadata (lines 281-380) → 4 files
8. Extract Book metadata (lines 381-480) → 4 files
9. Extract Comics metadata (from COMICS_MODULE.md) → 3 files
10. Extract Adult metadata (lines 606-680) → 3 files
11. Extract Scrobbling services (lines 481-580) → 5 files
12. Extract Anime services (lines 700-850) → 3 files
13. Extract Auth services (lines 851-950) → 4 files
14. Extract Request services (lines 200-250) → 2 files
15. Extract Audiobook (from existing docs) → 1 file
16. Extract Transcoding (internal) → 1 file
17. Extract Live TV (lines 951-1050) → 2 files
18. Extract Casting (lines 1051-1150) → 2 files
19. Extract Infrastructure (internal) → 4 files
20. Create all INDEX.md files (11 category indexes + 5 subcategory indexes)
21. Create master integrations/INDEX.md
22. Delete EXTERNAL_INTEGRATIONS_TODO.md
23. Commit: "docs: split EXTERNAL_INTEGRATIONS_TODO.md into 66 service files"

### Step 1.2: Split PREPARATION_MASTER_PLAN.md (1292 lines → 8 files)

**Directory Structure**:
```
docs/planning/
├── PREPARATION_MASTER_PLAN.md (condensed to 200 lines)
└── phases/
    ├── INDEX.md
    ├── PHASE_01_CORE_INFRASTRUCTURE.md
    ├── PHASE_02_MOVIE_MODULE.md
    ├── PHASE_03_TV_MODULE.md
    ├── PHASE_04_MUSIC_MODULE.md
    ├── PHASE_05_FRONTEND_UX.md
    └── APPENDIX_A_UX_UI_RESOURCES.md
```

**Operations** (8 operations):
1. Create `docs/planning/phases/` directory
2. Read PREPARATION_MASTER_PLAN.md
3. Extract Phase 1 (lines 100-300) → PHASE_01_CORE_INFRASTRUCTURE.md
4. Extract Phase 2 (lines 301-500) → PHASE_02_MOVIE_MODULE.md
5. Extract Phase 3 (lines 501-700) → PHASE_03_TV_MODULE.md
6. Extract Phase 4 (lines 701-900) → PHASE_04_MUSIC_MODULE.md
7. Extract Phase 5 (lines 901-1100) → PHASE_05_FRONTEND_UX.md
8. Extract Appendix A (lines 1101-1292) → APPENDIX_A_UX_UI_RESOURCES.md
9. Create phases/INDEX.md
10. Condense PREPARATION_MASTER_PLAN.md to 200-line executive summary
11. Commit: "docs: split PREPARATION_MASTER_PLAN.md into phase files"

### Step 1.3: Split MEDIA_ENHANCEMENTS.md (1319 lines → 11 files)

**Directory Structure**:
```
docs/features/
├── MEDIA_ENHANCEMENTS.md (condensed to 200 lines)
└── enhancements/
    ├── INDEX.md
    ├── TRAILERS.md
    ├── EXTRAS.md
    ├── SUBTITLES.md
    ├── LYRICS.md
    ├── ALBUM_ART.md
    ├── POSTERS_FANART.md
    ├── PREVIEWS_TRICKPLAY.md
    ├── CREDITS_CHAPTERS.md
    ├── RECOMMENDATIONS.md
    └── COLLECTIONS.md
```

**Operations** (13 operations):
1. Create `docs/features/enhancements/` directory
2. Read MEDIA_ENHANCEMENTS.md
3. Extract Trailers (lines 1-150) → TRAILERS.md
4. Extract Extras (lines 151-250) → EXTRAS.md
5. Extract Subtitles (lines 251-400) → SUBTITLES.md
6. Extract Lyrics (lines 401-500) → LYRICS.md
7. Extract Album Art (lines 501-600) → ALBUM_ART.md
8. Extract Posters/Fanart (lines 601-750) → POSTERS_FANART.md
9. Extract Previews/Trickplay (lines 751-950) → PREVIEWS_TRICKPLAY.md
10. Extract Credits/Chapters (lines 951-1100) → CREDITS_CHAPTERS.md
11. Extract Recommendations (lines 1101-1220) → RECOMMENDATIONS.md
12. Extract Collections (lines 1221-1319) → COLLECTIONS.md
13. Create enhancements/INDEX.md
14. Condense MEDIA_ENHANCEMENTS.md to 200-line overview
15. Commit: "docs: split MEDIA_ENHANCEMENTS.md into enhancement files"

### Step 1.4: Modularize PLAYER_ARCHITECTURE.md (1293 lines → 7 files)

**Directory Structure**:
```
docs/architecture/
├── PLAYER_ARCHITECTURE.md (condensed to 400 lines)
└── player/
    ├── INDEX.md
    ├── FORMATS.md
    ├── PROTOCOLS.md
    ├── TRANSCODING.md
    ├── BANDWIDTH.md
    ├── CLIENT_DETECTION.md
    └── EXAMPLES.md
```

**Operations** (9 operations):
1. Create `docs/architecture/player/` directory
2. Read PLAYER_ARCHITECTURE.md
3. Extract Formats/Codecs (lines 200-400) → FORMATS.md
4. Extract Protocols (lines 401-600) → PROTOCOLS.md
5. Extract Transcoding (lines 601-800) → TRANSCODING.md
6. Extract Bandwidth/ABR (lines 801-950) → BANDWIDTH.md
7. Extract Client Detection (lines 951-1100) → CLIENT_DETECTION.md
8. Extract Examples (lines 1101-1293) → EXAMPLES.md
9. Create player/INDEX.md
10. Condense PLAYER_ARCHITECTURE.md to 400-line high-level design
11. Commit: "docs: modularize PLAYER_ARCHITECTURE.md into component files"

### Step 1.5: Modularize AUDIO_STREAMING.md (1195 lines → 6 files)

**Directory Structure**:
```
docs/technical/
├── AUDIO_STREAMING.md (condensed to 400 lines)
└── audio/
    ├── INDEX.md
    ├── CODECS.md
    ├── PROTOCOLS.md
    ├── TRANSCODING.md
    ├── BITRATE_PROFILES.md
    └── EXAMPLES.md
```

**Operations** (7 operations):
1. Create `docs/technical/audio/` directory
2. Read AUDIO_STREAMING.md
3. Extract Codecs (lines 150-450) → CODECS.md
4. Extract Protocols (lines 451-650) → PROTOCOLS.md
5. Extract Transcoding (lines 651-850) → TRANSCODING.md
6. Extract Bitrate Profiles (lines 851-1000) → BITRATE_PROFILES.md
7. Extract Examples (lines 1001-1195) → EXAMPLES.md
8. Create audio/INDEX.md
9. Condense AUDIO_STREAMING.md to 400-line high-level design
10. Commit: "docs: modularize AUDIO_STREAMING.md into codec/protocol files"

---

## Phase 2: Eliminate Duplication (15 operations, 1 day)

### Step 2.1: Extract Adult Content Schema

**Operations** (5 operations):
1. Create `docs/schemas/` directory
2. Read WHISPARR_STASHDB_SCHEMA.md lines 97-250 (full schema)
3. Create `schemas/adult_content.sql` with canonical schema
4. Update WHISPARR_STASHDB_SCHEMA.md: Replace schema with link to schemas/adult_content.sql
5. Update ADULT_CONTENT_SYSTEM.md: Replace schema with link
6. Update MODULE_IMPLEMENTATION_TODO.md: Replace schema with link
7. Update ARCHITECTURE_V2.md: Replace schema reference with link
8. Commit: "docs: extract adult content schema to canonical schemas/adult_content.sql"

### Step 2.2: Consolidate Whisparr Integration Info

**Operations** (4 operations):
1. Read WHISPARR_STASHDB_SCHEMA.md lines 25-70 (Whisparr API)
2. Create `integrations/servarr/WHISPARR.md` with full Whisparr API details
3. Update WHISPARR_STASHDB_SCHEMA.md: Remove API details, link to integrations/servarr/WHISPARR.md
4. Update ADULT_METADATA.md: Remove duplicate Whisparr details, link to integration doc
5. Commit: "docs: consolidate Whisparr integration to integrations/servarr/WHISPARR.md"

### Step 2.3: Consolidate StashDB Integration Info

**Operations** (3 operations):
1. Read WHISPARR_STASHDB_SCHEMA.md lines 269-290 (StashDB API)
2. Create `integrations/metadata/adult/STASHDB.md` with full StashDB details
3. Update WHISPARR_STASHDB_SCHEMA.md: Remove StashDB details, link to integrations/metadata/adult/STASHDB.md
4. Update ADULT_METADATA.md: Remove duplicate StashDB details, link to integration doc
5. Commit: "docs: consolidate StashDB integration to integrations/metadata/adult/STASHDB.md"

### Step 2.4: Extract Arr Integration Pattern

**Operations** (3 operations):
1. Create `docs/patterns/` directory
2. Read MEDIA_ENHANCEMENTS.md lines 185-210 (Radarr trailer example)
3. Create `patterns/arr_integration.md` with Radarr/Sonarr/Lidarr common patterns
4. Update MEDIA_ENHANCEMENTS.md: Link to pattern instead of embedding example
5. Update ADULT_METADATA.md: Link to pattern for Whisparr comparison
6. Commit: "docs: extract Arr integration pattern to patterns/arr_integration.md"

---

## Phase 3: Create Index Files (9 operations, 4 hours)

### Step 3.1: Create Subdirectory Indexes

**Operations** (9 operations):
1. Create `architecture/INDEX.md` (5 files)
2. Create `operations/INDEX.md` (8 files)
3. Create `research/INDEX.md` (3 files)
4. Create `planning/INDEX.md` (4 files including phases/)
5. Create `features/INDEX.md` (12 files including enhancements/)
6. Create `technical/INDEX.md` (6 files including audio/)
7. Create `integrations/INDEX.md` (66 services across 11 categories)
8. Create `schemas/INDEX.md` (canonical schemas)
9. Create `patterns/INDEX.md` (reusable patterns)
10. Commit: "docs: create INDEX.md files for all subdirectories"

### Step 3.2: Update Master Index

**Operations** (1 operation):
1. Update `docs/INDEX.md` with links to all subdirectory indexes
2. Add quick reference table (file count, status, recent updates)
3. Add search guidance ("Looking for X? Check Y category")
4. Commit: "docs: update master INDEX.md with all subdirectory links"

---

## Phase 4: Extract Schemas & Patterns (20 operations, 1-2 days)

### Step 4.1: Extract Module Schemas

**Operations** (11 operations):
1. Read ARCHITECTURE_V2.md, MODULE_IMPLEMENTATION_TODO.md for schema examples
2. Create `schemas/movie.sql`
3. Create `schemas/tvshow.sql`
4. Create `schemas/music.sql`
5. Create `schemas/book.sql`
6. Create `schemas/audiobook.sql`
7. Create `schemas/podcast.sql`
8. Create `schemas/photo.sql`
9. Create `schemas/livetv.sql`
10. Create `schemas/collection.sql`
11. Create `schemas/comics.sql` (extract from COMICS_MODULE.md)
12. Create `schemas/user_data.sql` (common user data patterns)
13. Update source docs to link to schemas instead of embedding
14. Commit: "docs: extract all module schemas to schemas/ directory"

### Step 4.2: Extract Common Patterns

**Operations** (6 operations):
1. Create `patterns/module_structure.md` (from ARCHITECTURE_V2.md)
2. Create `patterns/api_patterns.md` (from API.md)
3. Create `patterns/trailer_fetching.md` (from MEDIA_ENHANCEMENTS.md)
4. Create `patterns/metadata_enrichment.md` (multi-source pattern)
5. Create `patterns/user_data_patterns.md` (ratings, history, favorites)
6. Create `patterns/transcoding_offload.md` (from OFFLOADING.md)
7. Update source docs to link to patterns
8. Commit: "docs: extract common patterns to patterns/ directory"

### Step 4.3: Create Cross-Reference Map

**Operations** (1 operation):
1. Create `docs/CROSS_REFERENCE_MAP.md`
2. Document primary/secondary sources for each topic
3. Map relationships between files
4. Commit: "docs: create CROSS_REFERENCE_MAP.md for documentation navigation"

---

## Phase 5: Final Cleanup & Validation (5 operations, 2 hours)

### Step 5.1: Update All Cross-References

**Operations** (1 operation):
1. Search all markdown files for old paths (e.g., `EXTERNAL_INTEGRATIONS_TODO.md`)
2. Update links to new paths (e.g., `integrations/servarr/RADARR.md`)
3. Validate no broken links
4. Commit: "docs: update all cross-references to new file paths"

### Step 5.2: Validate File Sizes

**Operations** (1 operation):
1. Run file size check: `Get-ChildItem -Recurse docs\*.md | Where-Object { (Get-Content $_.FullName).Count -gt 500 }`
2. Verify no files >500 lines (except intentional reference docs)
3. Document exceptions in DOCUMENTATION_ANALYSIS.md

### Step 5.3: Validate No Duplication

**Operations** (1 operation):
1. Search for duplicate schemas: `grep -r "CREATE SCHEMA" docs/`
2. Search for duplicate integration info: `grep -r "Whisparr.*API" docs/`
3. Verify single canonical source for each

### Step 5.4: Update AGENTS.md TODO

**Operations** (1 operation):
1. Update AGENTS.md TODO list: Mark documentation restructuring complete
2. Update file paths in "Related Documentation" sections
3. Commit: "docs: update AGENTS.md with restructuring completion"

### Step 5.5: Final Push

**Operations** (1 operation):
1. Run final validation: `git status`, `git diff --stat`
2. Push to develop: `git push origin develop`
3. Create GitHub PR summary (optional)

---

## Success Metrics

After completion, verify:
- [ ] No file >500 lines (except intentional reference docs)
- [ ] No duplicate schemas (canonical in `schemas/`)
- [ ] No duplicate integration info (canonical in `integrations/`)
- [ ] Every subdirectory has INDEX.md
- [ ] Master INDEX.md links to all subdirectories
- [ ] CROSS_REFERENCE_MAP.md exists
- [ ] All cross-references updated
- [ ] AGENTS.md updated
- [ ] All changes committed and pushed

---

## Commit Strategy

**Commits per phase**:
- Phase 1: 5 commits (one per massive file split)
- Phase 2: 4 commits (one per deduplication task)
- Phase 3: 2 commits (indexes + master index)
- Phase 4: 3 commits (schemas, patterns, cross-ref map)
- Phase 5: 2 commits (cross-ref updates, AGENTS.md)

**Total commits**: ~16 structured commits

---

## Rollback Plan

If issues arise:
```bash
# Rollback to pre-restructure state
git checkout develop
git reset --hard bd07b59fa  # Last commit before restructure
git push origin develop --force

# Or revert specific commit
git revert <commit-hash>
```

---

## Execution Order

Execute in strict order:
1. Phase 1.1 → Commit → Push
2. Phase 1.2 → Commit → Push
3. Phase 1.3 → Commit → Push
4. Phase 1.4 → Commit → Push
5. Phase 1.5 → Commit → Push
6. Phase 2 (all steps) → Commit → Push
7. Phase 3 (all steps) → Commit → Push
8. Phase 4 (all steps) → Commit → Push
9. Phase 5 (all steps) → Commit → Push

**DO NOT** skip ahead or reorder steps - dependencies exist between phases.
