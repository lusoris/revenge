# Design Documentation Reorganization Proposal

> Restructure docs for better human and AI readability

**Status**: 📋 PROPOSAL
**Date**: 2026-01-30

---

## Current Issues

1. **Mixed concerns** - Backend/frontend/ops mixed in same directories
2. **Deep nesting** - Hard to navigate (integrations/metadata/books/GOODREADS.md)
3. **No quick reference** - No single-page overview of entire system
4. **Inconsistent naming** - Some files use INDEX.md, others don't
5. **Missing cross-references** - Hard to understand dependencies

---

## Proposed Structure

```
docs/dev/design/
├── README.md                     # Quick start + navigation
├── 00_SOURCE_OF_TRUTH.md            # Single-page system overview (tables)
├── PACKAGES.md                   # All Go packages with versions
│
├── 01-architecture/              # System design (read first)
│   ├── README.md                 # Architecture overview
│   ├── PRINCIPLES.md             # Design principles
│   ├── DATA_FLOW.md              # How data flows through system
│   ├── SECURITY.md               # Security architecture
│   └── DECISIONS.md              # ADRs (Architecture Decision Records)
│
├── 02-backend/                   # Go backend implementation
│   ├── README.md                 # Backend overview
│   ├── modules/                  # Content modules
│   │   ├── README.md             # Module pattern overview
│   │   ├── movie.md              # Movie module spec
│   │   ├── tv.md                 # TV module spec
│   │   ├── music.md              # Music module spec
│   │   ├── audiobook.md          # Audiobook module spec
│   │   ├── book.md               # Book module spec
│   │   ├── podcast.md            # Podcast module spec
│   │   ├── photo.md              # Photo module spec
│   │   ├── comics.md             # Comics module spec
│   │   ├── livetv.md             # LiveTV module spec
│   │   └── qar.md                # QAR module spec (obfuscated)
│   ├── services/                 # Shared services
│   │   ├── README.md
│   │   ├── auth.md
│   │   ├── user.md
│   │   ├── playback.md
│   │   ├── metadata.md
│   │   ├── search.md
│   │   └── jobs.md
│   ├── database/                 # Database design
│   │   ├── README.md
│   │   ├── schemas.md            # All schema definitions
│   │   ├── migrations.md         # Migration strategy
│   │   └── queries.md            # Query patterns
│   └── api/                      # API design
│       ├── README.md
│       ├── openapi.md            # OpenAPI conventions
│       ├── auth.md               # Authentication
│       └── errors.md             # Error handling
│
├── 03-frontend/                  # SvelteKit frontend
│   ├── README.md                 # Frontend overview
│   ├── components/               # Component library
│   ├── pages/                    # Page structure
│   ├── state/                    # State management
│   └── player/                   # Media player
│
├── 04-integrations/              # External services
│   ├── README.md                 # Integration patterns
│   ├── arr/                      # *arr ecosystem
│   │   ├── radarr.md
│   │   ├── sonarr.md
│   │   ├── lidarr.md
│   │   ├── whisparr.md
│   │   └── prowlarr.md
│   ├── metadata/                 # Metadata providers
│   │   ├── tmdb.md
│   │   ├── musicbrainz.md
│   │   ├── stashdb.md
│   │   └── ...
│   ├── scrobbling/               # Scrobbling services
│   │   ├── trakt.md
│   │   ├── lastfm.md
│   │   └── ...
│   └── auth/                     # Auth providers
│       ├── oidc.md
│       └── ...
│
├── 05-operations/                # Deployment & ops
│   ├── README.md
│   ├── deployment.md
│   ├── monitoring.md
│   ├── backup.md
│   └── security.md
│
├── 06-features/                  # Feature specifications
│   ├── README.md
│   ├── playback/                 # Playback features
│   │   ├── syncplay.md
│   │   ├── trickplay.md
│   │   └── skip-intro.md
│   ├── discovery/                # Content discovery
│   │   ├── recommendations.md
│   │   ├── search.md
│   │   └── collections.md
│   └── social/                   # Social features
│       ├── sharing.md
│       └── activity.md
│
└── 99-reference/                 # Quick reference tables
    ├── GLOSSARY.md               # Term definitions
    ├── PACKAGES.md               # Package versions
    ├── API_ENDPOINTS.md          # All endpoints
    ├── DATABASE_TABLES.md        # All tables
    └── ENV_VARS.md               # All config options
```

---

## Key Changes

### 1. Numbered Directories
- Forces reading order: architecture → backend → frontend → integrations
- Easier to navigate alphabetically

### 2. README.md in Every Directory
- Each folder has overview + links to children
- AI can read README first to understand context

### 3. Flat Module Structure
- `02-backend/modules/movie.md` instead of nested `features/movies/MOVIE_MODULE.md`
- All modules at same level for easy comparison

### 4. Source of Truth Tables
- Single document with all modules, packages, versions
- Easy to scan, verify, update

### 5. Reference Section
- Quick lookup tables for common queries
- API endpoints, database tables, env vars all in one place

---

## Migration Plan

### Phase 1: Create new structure (empty)
1. Create new directory tree
2. Create README.md stubs

### Phase 2: Migrate content
1. Copy existing docs to new locations
2. Update internal links
3. Consolidate duplicate information

### Phase 3: Create reference tables
1. Generate 00_SOURCE_OF_TRUTH.md
2. Generate PACKAGES.md
3. Generate API_ENDPOINTS.md

### Phase 4: Cleanup
1. Remove old directories
2. Update root INDEX.md
3. Update any external links

---

## Benefits

| Benefit | Before | After |
|---------|--------|-------|
| Find module spec | 4+ clicks | 2 clicks |
| Understand reading order | Unclear | Numbered dirs |
| Check package version | Search multiple files | Single PACKAGES.md |
| AI context loading | Load many files | Load README + specific doc |
| Add new module | Copy pattern from scattered files | Copy from modules/README.md |

---

## Questions for Review

1. Should QAR (adult) docs be in a separate, encrypted location?
2. Should we version the design docs (v1, v2)?
3. Should frontend/backend have completely separate repos?
4. How often should 00_SOURCE_OF_TRUTH.md be regenerated?
