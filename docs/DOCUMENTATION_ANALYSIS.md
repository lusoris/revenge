# Documentation Analysis & Restructuring Plan

> Deep analysis of documentation for deduplication, modularity, and discoverability

**Date**: 2026-01-28
**Status**: 🔴 CRITICAL - Major restructuring needed
**Total Files**: 39 markdown files (113,832 lines total)

---

## Executive Summary

### Problems Identified

1. **MASSIVE FILES** (Top 5 consume 50% of content):
   - `EXTERNAL_INTEGRATIONS_TODO.md` - 1445 lines (93KB) ← **MUST SPLIT**
   - `MEDIA_ENHANCEMENTS.md` - 1319 lines (78KB) ← **MUST SPLIT**
   - `PLAYER_ARCHITECTURE.md` - 1293 lines (75KB) ← **NEEDS MODULAR BREAKDOWN**
   - `PREPARATION_MASTER_PLAN.md` - 1292 lines (79KB) ← **MUST SPLIT**
   - `AUDIO_STREAMING.md` - 1195 lines (73KB) ← **NEEDS MODULAR BREAKDOWN**

2. **DUPLICATION DETECTED**:
   - **Adult content schema** (`CREATE SCHEMA c;`) duplicated across:
     * `WHISPARR_STASHDB_SCHEMA.md` (full schema)
     * `ADULT_CONTENT_SYSTEM.md` (full schema)
     * `MODULE_IMPLEMENTATION_TODO.md` (example)
     * `ARCHITECTURE_V2.md` (migration reference)

   - **Whisparr/StashDB info** scattered across:
     * `WHISPARR_STASHDB_SCHEMA.md` (primary, 375 lines)
     * `ADULT_METADATA.md` (partial, 1001 lines)
     * `PREPARATION_MASTER_PLAN.md` (TODO item)
     * `EXTERNAL_INTEGRATIONS_TODO.md` (integration stub)

   - **Radarr/Sonarr/Lidarr integration** mentioned in:
     * `PREPARATION_MASTER_PLAN.md` (API status)
     * `EXTERNAL_INTEGRATIONS_TODO.md` (integration details)
     * `MEDIA_ENHANCEMENTS.md` (trailer fetching example)
     * `ADULT_METADATA.md` (comparison)

3. **NO INDEX FILES** per subdirectory:
   - `architecture/` - 5 files, no index
   - `operations/` - 8 files, no index
   - `research/` - 3 files, no index
   - `planning/` - 3 files, no index
   - `features/` - 11 files, no index
   - `technical/` - 5 files, no index

4. **POOR MODULARITY**:
   - Schema definitions embedded in feature docs instead of separate schema files
   - API integration info mixed with high-level architecture
   - Examples not extracted to reusable patterns
   - No cross-referencing between related content

---

## File Size Distribution

| Category | Files | Total Lines | Avg Lines/File | Issues |
|----------|-------|-------------|----------------|--------|
| **CRITICAL** (>1000 lines) | 5 | 6,544 | 1,309 | Must split |
| **LARGE** (500-1000 lines) | 12 | 10,157 | 846 | Consider splitting |
| **MEDIUM** (200-500 lines) | 14 | 4,722 | 337 | Good size |
| **SMALL** (<200 lines) | 8 | 1,409 | 176 | Good size |

**Files >1000 lines** (CRITICAL):
```
1445 lines - EXTERNAL_INTEGRATIONS_TODO.md (66 services, MUST split by category)
1319 lines - features/MEDIA_ENHANCEMENTS.md (10+ features, MUST split)
1293 lines - architecture/PLAYER_ARCHITECTURE.md (needs modular breakdown)
1292 lines - PREPARATION_MASTER_PLAN.md (5 phases, MUST split)
1195 lines - technical/AUDIO_STREAMING.md (multiple codecs/protocols, needs breakdown)
```

---

## Content Duplication Matrix

### Schema Definitions

| Schema | Primary Source | Duplicates | Action |
|--------|---------------|------------|--------|
| `CREATE SCHEMA c;` | `features/WHISPARR_STASHDB_SCHEMA.md` | 3 files | Extract to `schemas/adult_content.sql`, reference only |
| Adult tables | `features/WHISPARR_STASHDB_SCHEMA.md` | 2 files | Keep primary, link others |
| Movie/TV schemas | `architecture/ARCHITECTURE_V2.md` | 2 files | Extract to `schemas/` directory |

### Integration Info

| Service | Primary Source | Duplicates | Action |
|---------|---------------|------------|--------|
| Whisparr | `EXTERNAL_INTEGRATIONS_TODO.md` (stub) | `WHISPARR_STASHDB_SCHEMA.md` (full), `ADULT_METADATA.md` (partial) | Move full info to `integrations/WHISPARR.md`, link others |
| StashDB | `EXTERNAL_INTEGRATIONS_TODO.md` (stub) | `WHISPARR_STASHDB_SCHEMA.md` (full) | Move full info to `integrations/STASHDB.md` |
| Radarr/Sonarr/Lidarr | `EXTERNAL_INTEGRATIONS_TODO.md` (full) | 3 files (examples) | Keep integration doc, extract examples to `patterns/` |

### Architecture Patterns

| Pattern | Primary Source | Duplicates | Action |
|---------|---------------|------------|--------|
| Module structure | `architecture/ARCHITECTURE_V2.md` | 3 files | Extract to `patterns/module_structure.md` |
| API design | `technical/API.md` | 2 files | Extract examples to `patterns/api_patterns.md` |
| Transcoding offload | `technical/OFFLOADING.md` | 2 files | Keep single source |

---

## Restructuring Plan

### Phase 1: Split Massive Files (URGENT)

#### 1.1 Split EXTERNAL_INTEGRATIONS_TODO.md (1445 lines → 17 files)

**Target structure**: `docs/integrations/`

```
integrations/
├── INDEX.md                           # Master navigation (66 services)
├── servarr/
│   ├── INDEX.md
│   ├── RADARR.md                      # Already researched
│   ├── SONARR.md                      # Already researched
│   ├── LIDARR.md                      # Already researched
│   ├── WHISPARR.md                    # Pending (move from WHISPARR_STASHDB_SCHEMA.md)
│   └── READARR.md                     # Pending research
├── metadata/
│   ├── INDEX.md
│   ├── video/
│   │   ├── TMDB.md                    # Already researched
│   │   ├── THETVDB.md                 # Already researched
│   │   ├── OMDB.md                    # Pending research
│   │   └── THEPOSTERDB.md             # Pending research
│   ├── music/
│   │   ├── MUSICBRAINZ.md             # Already researched
│   │   ├── LASTFM.md                  # Already researched
│   │   ├── SPOTIFY.md                 # Pending research
│   │   └── DISCOGS.md                 # Pending research
│   ├── books/
│   │   ├── GOODREADS.md               # Pending research
│   │   ├── OPENLIBRARY.md             # Pending research
│   │   ├── AUDIBLE.md                 # Pending research
│   │   └── HARDCOVER.md               # Pending research
│   ├── comics/
│   │   ├── COMICVINE.md               # From COMICS_MODULE.md
│   │   ├── MARVEL_API.md              # From COMICS_MODULE.md
│   │   └── GRAND_COMICS_DATABASE.md   # From COMICS_MODULE.md
│   └── adult/
│       ├── STASHDB.md                 # Move from WHISPARR_STASHDB_SCHEMA.md
│       ├── THEPORNDB.md               # Pending research
│       └── STASH.md                   # Pending research
├── scrobbling/
│   ├── INDEX.md
│   ├── TRAKT.md                       # Already researched
│   ├── LASTFM.md                      # Already researched (duplicate with metadata)
│   ├── LISTENBRAINZ.md                # Pending research
│   ├── LETTERBOXD.md                  # Pending research
│   └── SIMKL.md                       # Pending research
├── anime/
│   ├── INDEX.md
│   ├── ANILIST.md                     # Already researched
│   ├── MYANIMELIST.md                 # Already researched
│   └── KITSU.md                       # Already researched
├── auth/
│   ├── INDEX.md
│   ├── AUTHELIA.md                    # Pending research
│   ├── AUTHENTIK.md                   # Pending research
│   ├── KEYCLOAK.md                    # Pending research
│   └── GENERIC_OIDC.md                # Generic pattern
├── request/
│   ├── INDEX.md
│   ├── OVERSEERR.md                   # Already researched
│   └── JELLYSEERR.md                  # Pending research
├── audiobook/
│   ├── INDEX.md
│   └── AUDIOBOOKSHELF.md              # Already researched
├── transcoding/
│   ├── INDEX.md
│   └── BLACKBEARD.md                  # Internal service
├── livetv/
│   ├── INDEX.md
│   ├── TVHEADEND.md                   # Pending research
│   └── NEXTPVR.md                     # Pending research
├── casting/
│   ├── INDEX.md
│   ├── CHROMECAST.md                  # Pending research
│   └── DLNA.md                        # Pending research
└── infrastructure/
    ├── INDEX.md
    ├── POSTGRESQL.md                  # Database
    ├── DRAGONFLY.md                   # Cache
    ├── TYPESENSE.md                   # Search
    └── RIVER.md                       # Job queue
```

**Action Items**:
- [ ] Create `integrations/` subdirectory structure (11 categories)
- [ ] Extract each service from EXTERNAL_INTEGRATIONS_TODO.md to dedicated file
- [ ] Create category INDEX.md files (11 files)
- [ ] Create master `integrations/INDEX.md`
- [ ] Delete EXTERNAL_INTEGRATIONS_TODO.md after split
- [ ] Update cross-references in all docs

#### 1.2 Split PREPARATION_MASTER_PLAN.md (1292 lines → 8 files)

**Target structure**: `docs/planning/phases/`

```
planning/
├── PREPARATION_MASTER_PLAN.md         # Keep as high-level overview (200 lines max)
└── phases/
    ├── INDEX.md                       # Phase navigation
    ├── PHASE_01_CORE_INFRASTRUCTURE.md   # Week 1-2: PostgreSQL, Dragonfly, Typesense, River, Echo
    ├── PHASE_02_MOVIE_MODULE.md          # Week 2-3: Radarr integration, TMDb metadata, transcoding
    ├── PHASE_03_TV_MODULE.md             # Week 3-4: Sonarr integration, TheTVDB metadata, episode management
    ├── PHASE_04_MUSIC_MODULE.md          # Week 4-6: Lidarr integration, MusicBrainz, Last.fm scrobbling
    ├── PHASE_05_FRONTEND_UX.md           # Week 9-11: Svelte 5 runes, shadcn-svelte, accessibility
    └── APPENDIX_A_UX_UI_RESOURCES.md     # 22 UX/UI sources
```

**Action Items**:
- [ ] Create `planning/phases/` subdirectory
- [ ] Extract each phase to dedicated file (5 files)
- [ ] Extract appendix to dedicated file
- [ ] Create `phases/INDEX.md` navigation
- [ ] Condense PREPARATION_MASTER_PLAN.md to 200-line executive summary
- [ ] Update cross-references

#### 1.3 Split MEDIA_ENHANCEMENTS.md (1319 lines → 10 files)

**Target structure**: `docs/features/enhancements/`

```
features/
├── MEDIA_ENHANCEMENTS.md              # Keep as overview (200 lines max)
└── enhancements/
    ├── INDEX.md                       # Feature navigation
    ├── TRAILERS.md                    # Trailer fetching (Radarr, TMDb, YouTube)
    ├── EXTRAS.md                      # Behind-the-scenes, deleted scenes, interviews
    ├── SUBTITLES.md                   # Subtitle search (OpenSubtitles, Subscene)
    ├── LYRICS.md                      # Lyrics fetching (Musixmatch, Genius)
    ├── ALBUM_ART.md                   # Album art sources (Last.fm, MusicBrainz, Spotify)
    ├── POSTERS_FANART.md              # Poster/fanart sources (Fanart.tv, ThePosterDB)
    ├── PREVIEWS_TRICKPLAY.md          # Video previews, trickplay thumbnails
    ├── CREDITS_CHAPTERS.md            # Credits detection, chapter markers
    ├── RECOMMENDATIONS.md             # Content recommendations (Trakt, TMDb, collaborative filtering)
    └── COLLECTIONS.md                 # Smart collections, user-defined collections
```

**Action Items**:
- [ ] Create `features/enhancements/` subdirectory
- [ ] Extract each feature to dedicated file (10 files)
- [ ] Create `enhancements/INDEX.md` navigation
- [ ] Condense MEDIA_ENHANCEMENTS.md to 200-line overview
- [ ] Update cross-references

#### 1.4 Modularize PLAYER_ARCHITECTURE.md (1293 lines → keep single file but extract examples)

**Target structure**: `docs/architecture/player/`

```
architecture/
├── PLAYER_ARCHITECTURE.md             # High-level design (400 lines max)
└── player/
    ├── INDEX.md                       # Component navigation
    ├── FORMATS.md                     # Supported codecs/containers
    ├── PROTOCOLS.md                   # HLS, DASH, Direct Play
    ├── TRANSCODING.md                 # Blackbeard integration
    ├── BANDWIDTH.md                   # Adaptive bitrate, bandwidth detection
    ├── CLIENT_DETECTION.md            # Client capabilities, codec support
    └── EXAMPLES.md                    # Code examples, integration patterns
```

**Action Items**:
- [ ] Create `architecture/player/` subdirectory
- [ ] Extract sub-topics to dedicated files (6 files)
- [ ] Create `player/INDEX.md` navigation
- [ ] Keep PLAYER_ARCHITECTURE.md as high-level design (400 lines)
- [ ] Update cross-references

#### 1.5 Modularize AUDIO_STREAMING.md (1195 lines → keep single file but extract codec details)

**Target structure**: `docs/technical/audio/`

```
technical/
├── AUDIO_STREAMING.md                 # High-level design (400 lines max)
└── audio/
    ├── INDEX.md                       # Codec/protocol navigation
    ├── CODECS.md                      # AAC, MP3, FLAC, Opus, Vorbis, DTS, TrueHD
    ├── PROTOCOLS.md                   # HLS, DASH, Direct Play
    ├── TRANSCODING.md                 # FFmpeg integration, Blackbeard
    ├── BITRATE_PROFILES.md            # Low/Medium/High/Lossless quality tiers
    └── EXAMPLES.md                    # Code examples, integration patterns
```

**Action Items**:
- [ ] Create `technical/audio/` subdirectory
- [ ] Extract sub-topics to dedicated files (5 files)
- [ ] Create `audio/INDEX.md` navigation
- [ ] Keep AUDIO_STREAMING.md as high-level design (400 lines)
- [ ] Update cross-references

---

### Phase 2: Eliminate Duplication

#### 2.1 Schema Definitions

**Problem**: `CREATE SCHEMA c;` + adult tables duplicated across 4 files.

**Solution**: Extract to canonical source.

**Action**:
- [ ] Create `docs/schemas/` directory
- [ ] Create `schemas/adult_content.sql` (canonical schema)
- [ ] Update 4 files to reference `schemas/adult_content.sql` instead of embedding
- [ ] Pattern: "See [Adult Content Schema](../schemas/adult_content.sql) for full schema"

**Files to deduplicate**:
- `features/WHISPARR_STASHDB_SCHEMA.md` (line 97: `CREATE SCHEMA IF NOT EXISTS c;` + full tables)
- `features/ADULT_CONTENT_SYSTEM.md` (line 414: `CREATE SCHEMA IF NOT EXISTS c;` + full tables)
- `planning/MODULE_IMPLEMENTATION_TODO.md` (line 337: example schema)
- `architecture/ARCHITECTURE_V2.md` (line 623: migration reference)

#### 2.2 Integration Info

**Problem**: Whisparr/StashDB info scattered across 4 files.

**Solution**: Consolidate to dedicated integration docs.

**Action**:
- [ ] Create `integrations/servarr/WHISPARR.md` (move from WHISPARR_STASHDB_SCHEMA.md API sections)
- [ ] Create `integrations/metadata/adult/STASHDB.md` (move from WHISPARR_STASHDB_SCHEMA.md StashDB sections)
- [ ] Update `features/WHISPARR_STASHDB_SCHEMA.md` to focus on Revenge schema design (link to integrations)
- [ ] Update `features/ADULT_METADATA.md` to link to integration docs instead of duplicating

**Files to deduplicate**:
- `features/WHISPARR_STASHDB_SCHEMA.md` (lines 25-70: Whisparr API, lines 269-290: StashDB API)
- `features/ADULT_METADATA.md` (lines 1-100: Whisparr overview, lines 500-600: StashDB details)
- `EXTERNAL_INTEGRATIONS_TODO.md` (lines 104-120: Whisparr stub, lines 606-630: StashDB stub)
- `PREPARATION_MASTER_PLAN.md` (line 68: Whisparr TODO, line 83: StashDB TODO)

#### 2.3 Radarr/Sonarr/Lidarr Examples

**Problem**: Integration examples scattered across 4 files.

**Solution**: Extract to reusable patterns.

**Action**:
- [ ] Create `docs/patterns/` directory
- [ ] Create `patterns/arr_integration.md` (Radarr/Sonarr/Lidarr common patterns)
- [ ] Create `patterns/trailer_fetching.md` (example from MEDIA_ENHANCEMENTS.md)
- [ ] Update 4 files to link to patterns instead of duplicating examples

**Files to deduplicate**:
- `features/MEDIA_ENHANCEMENTS.md` (lines 185-210: Radarr trailer fetching example)
- `features/ADULT_METADATA.md` (lines 80-90: Whisparr comparison to Radarr)
- `PREPARATION_MASTER_PLAN.md` (lines 53-55: Radarr/Sonarr/Lidarr API status)
- `EXTERNAL_INTEGRATIONS_TODO.md` (lines 17-60: Radarr integration details)

---

### Phase 3: Create Index Files

#### 3.1 Per-Subdirectory Indexes

Create `INDEX.md` in each subdirectory for navigation.

**Template**:
```markdown
# {Category} Documentation

> {Brief description}

**Quick Links**: [{Topic 1}]({FILE1}.md) | [{Topic 2}]({FILE2}.md) | [{Topic 3}]({FILE3}.md)

---

## Files in this Category

- **[{File 1}]({FILE1}.md)** - {One-line description}
- **[{File 2}]({FILE2}.md)** - {One-line description}
- **[{File 3}]({FILE3}.md)** - {One-line description}

---

## Related Categories

- [{Other Category}](../{other}/) - {Why related}
```

**Action Items**:
- [ ] Create `architecture/INDEX.md` (5 files)
- [ ] Create `operations/INDEX.md` (8 files)
- [ ] Create `research/INDEX.md` (3 files)
- [ ] Create `planning/INDEX.md` (3 files)
- [ ] Create `features/INDEX.md` (11 files)
- [ ] Create `technical/INDEX.md` (5 files)
- [ ] Create `integrations/INDEX.md` (66 services, post-split)
- [ ] Create `schemas/INDEX.md` (canonical schemas, post-creation)
- [ ] Create `patterns/INDEX.md` (reusable patterns, post-creation)

#### 3.2 Update Master Index

Update `docs/INDEX.md` with:
- Links to all subdirectory indexes
- Quick reference table (file count, status, recent updates)
- Search guidance ("Looking for X? Check Y category")

---

### Phase 4: Improve Modularity

#### 4.1 Extract Schemas

Create `docs/schemas/` directory for canonical SQL schemas:

```
schemas/
├── INDEX.md
├── adult_content.sql              # Full schema c.* (from WHISPARR_STASHDB_SCHEMA.md)
├── movie.sql                      # Full movie module schema
├── tvshow.sql                     # Full TV show module schema
├── music.sql                      # Full music module schema
├── book.sql                       # Full book module schema
├── audiobook.sql                  # Full audiobook module schema
├── podcast.sql                    # Full podcast module schema
├── photo.sql                      # Full photo module schema
├── livetv.sql                     # Full live TV module schema
├── collection.sql                 # Full collection module schema
├── comics.sql                     # Full comics module schema (from COMICS_MODULE.md)
└── user_data.sql                  # User preferences, ratings, history (shared patterns)
```

**Action**:
- [ ] Create `schemas/` directory
- [ ] Extract schemas from feature docs
- [ ] Create schema-specific files
- [ ] Update feature docs to reference schemas (don't duplicate)

#### 4.2 Extract Patterns

Create `docs/patterns/` directory for reusable implementation patterns:

```
patterns/
├── INDEX.md
├── module_structure.md            # Standard module layout (from ARCHITECTURE_V2.md)
├── api_patterns.md                # REST API design patterns (from API.md)
├── arr_integration.md             # Radarr/Sonarr/Lidarr integration patterns
├── trailer_fetching.md            # Trailer fetching pattern (from MEDIA_ENHANCEMENTS.md)
├── metadata_enrichment.md         # Multi-source metadata enrichment
├── user_data_patterns.md          # Ratings, watch history, favorites (common across modules)
└── transcoding_offload.md         # Blackbeard integration pattern (from OFFLOADING.md)
```

**Action**:
- [ ] Create `patterns/` directory
- [ ] Extract common patterns from multiple docs
- [ ] Create pattern-specific files
- [ ] Update docs to reference patterns instead of duplicating examples

#### 4.3 Cross-Reference Map

Create `docs/CROSS_REFERENCE_MAP.md` to document relationships:

```markdown
# Cross-Reference Map

## Adult Content (`c` schema)

**Primary**: `features/ADULT_CONTENT_SYSTEM.md`
**Schema**: `schemas/adult_content.sql`
**Integration**: `integrations/servarr/WHISPARR.md`, `integrations/metadata/adult/STASHDB.md`
**UI/UX**: `features/WHISPARR_STASHDB_SCHEMA.md`
**Related**: `features/ADULT_METADATA.md`, `architecture/ARCHITECTURE_V2.md`

## Movie Module

**Primary**: `architecture/ARCHITECTURE_V2.md` (section)
**Schema**: `schemas/movie.sql`
**Integration**: `integrations/servarr/RADARR.md`, `integrations/metadata/video/TMDB.md`
**Patterns**: `patterns/arr_integration.md`, `patterns/trailer_fetching.md`
**Enhancements**: `features/enhancements/TRAILERS.md`, `features/enhancements/EXTRAS.md`
```

**Action**:
- [ ] Create `docs/CROSS_REFERENCE_MAP.md`
- [ ] Document primary/secondary sources for each topic
- [ ] Map relationships between files
- [ ] Update monthly as docs evolve

---

## Estimated Effort

| Phase | Tasks | Estimated Time |
|-------|-------|----------------|
| Phase 1: Split massive files | 40 files to create | 2-3 days |
| Phase 2: Eliminate duplication | 15 files to update | 1 day |
| Phase 3: Create indexes | 9 index files | 4 hours |
| Phase 4: Improve modularity | 20 files to create/update | 1-2 days |
| **TOTAL** | **~85 file operations** | **4-6 days** |

---

## Success Criteria

- ✅ No file >500 lines (except intentional reference docs)
- ✅ No schema duplication (canonical source in `schemas/`)
- ✅ No integration info duplication (canonical source in `integrations/`)
- ✅ Every subdirectory has `INDEX.md`
- ✅ Master `docs/INDEX.md` links to all subdirectories
- ✅ `CROSS_REFERENCE_MAP.md` documents all relationships
- ✅ AI agents can find info in <3 navigation steps
- ✅ Humans can find info in <5 clicks

---

## Next Steps

1. **IMMEDIATE**: Split `EXTERNAL_INTEGRATIONS_TODO.md` (1445 lines → 66 service files)
2. **URGENT**: Split `PREPARATION_MASTER_PLAN.md` (1292 lines → 8 phase files)
3. **HIGH PRIORITY**: Split `MEDIA_ENHANCEMENTS.md` (1319 lines → 11 enhancement files)
4. **MEDIUM PRIORITY**: Extract schemas to `schemas/` directory
5. **MEDIUM PRIORITY**: Extract patterns to `patterns/` directory
6. **LOW PRIORITY**: Create all `INDEX.md` files
7. **FINAL**: Create `CROSS_REFERENCE_MAP.md`

**Commit strategy**: Commit after each phase to preserve git history.
