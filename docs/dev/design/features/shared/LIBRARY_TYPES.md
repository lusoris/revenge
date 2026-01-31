# Library Types

> Per-module library architecture and supported content types

## Status

| Dimension | Status | Notes |
|-----------|--------|-------|
| Design | ✅ | Full design with per-module tables, polymorphic permissions |
| Sources | 🟡 | Architecture references |
| Instructions | ✅ | Implementation checklist complete |
| Code | 🔴 | |
| Linting | 🔴 | |
| Unit Testing | 🔴 | |
| Integration Testing | 🔴 | |

**Location**: `internal/module/` (per-module library implementations)
**Current Migration**: `shared/000005_libraries.up.sql` - **TO BE SPLIT**

---

## Architecture Overview

Revenge uses **per-module library tables** for full module isolation:

```
┌─────────────────────────────────────────────────────────────────────────┐
│                           Module Boundaries                              │
├─────────────────┬─────────────────┬─────────────────┬───────────────────┤
│  movie module   │  tvshow module  │  music module   │   adult module    │
│                 │                 │                 │   (qar schema)    │
│ ┌─────────────┐ │ ┌─────────────┐ │ ┌─────────────┐ │ ┌───────────────┐ │
│ │movie_       │ │ │tv_          │ │ │music_       │ │ │qar.fleets     │ │
│ │libraries    │ │ │libraries    │ │ │libraries    │ │ │               │ │
│ └──────┬──────┘ │ └──────┬──────┘ │ └──────┬──────┘ │ └───────┬───────┘ │
│        │        │        │        │        │        │         │         │
│        ▼        │        ▼        │        ▼        │         ▼         │
│ ┌─────────────┐ │ ┌─────────────┐ │ ┌─────────────┐ │ ┌───────────────┐ │
│ │   movies    │ │ │   series    │ │ │   albums    │ │ │qar.expeditions│ │
│ └─────────────┘ │ │   seasons   │ │ │   tracks    │ │ │qar.voyages    │ │
│                 │ │   episodes  │ │ │   artists   │ │ │qar.crew       │ │
│                 │ └─────────────┘ │ └─────────────┘ │ └───────────────┘ │
└─────────────────┴─────────────────┴─────────────────┴───────────────────┘
```

### Why Per-Module Libraries?

1. **Full module isolation** - Each module owns its complete data model
2. **Module-specific settings** - Library settings baked into each schema
3. **No shared enum** - Modules can be added without modifying shared code
4. **Independent deployment** - Modules can be enabled/disabled cleanly
5. **Schema-level isolation** - Adult content fully contained in `qar` schema (Queen Anne's Revenge)

### Shared Components

Only truly cross-cutting concerns remain shared:
- `shared.users` - User accounts
- `shared.profiles` - User profiles
- `shared.sessions` - Auth sessions
- `shared.resource_grants` - Polymorphic access control (see [RBAC_CASBIN.md](RBAC_CASBIN.md))

---

## Per-Module Library Tables

Each module defines its own library table with module-specific settings:

| Module | Library Table | Content Tables | Migration |
|--------|---------------|----------------|-----------|
| movie | `movie_libraries` | `movies`, `movie_collections` | `movie/000001_*.sql` |
| tvshow | `tv_libraries` | `series`, `seasons`, `episodes` | `tvshow/000001_*.sql` |
| music | `music_libraries` | `albums`, `tracks`, `artists` | `music/000001_*.sql` |
| audiobook | `audiobook_libraries` | `audiobooks`, `chapters` | `audiobook/000001_*.sql` |
| book | `book_libraries` | `books` | `book/000001_*.sql` |
| podcast | `podcast_libraries` | `podcasts`, `podcast_episodes` | `podcast/000001_*.sql` |
| photo | `photo_libraries` | `photos`, `photo_albums` | `photo/000001_*.sql` |
| livetv | `livetv_sources` | `channels`, `programs`, `recordings` | `livetv/000001_*.sql` |
| comics | `comic_libraries` | `comics`, `issues` | `comics/000001_*.sql` |
| adult | `qar.fleets` | `qar.expeditions`, `qar.voyages`, `qar.crew` | `qar/000001_*.sql` |

---

## Example: Movie Library Schema

Each module defines its own library table with module-specific settings:

```sql
-- movie/000001_movie_libraries.up.sql
CREATE TABLE movie_libraries (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name            VARCHAR(255) NOT NULL,
    paths           TEXT[] NOT NULL,

    -- Movie-specific settings
    scan_enabled        BOOLEAN NOT NULL DEFAULT true,
    scan_interval_hours INT NOT NULL DEFAULT 24,
    last_scan_at        TIMESTAMPTZ,

    -- Metadata settings
    preferred_language      VARCHAR(10) DEFAULT 'en',
    tmdb_enabled            BOOLEAN NOT NULL DEFAULT true,
    imdb_enabled            BOOLEAN NOT NULL DEFAULT true,
    download_trailers       BOOLEAN NOT NULL DEFAULT false,
    download_backdrops      BOOLEAN NOT NULL DEFAULT true,

    -- Access control (simple, module handles it)
    is_private          BOOLEAN NOT NULL DEFAULT false,
    owner_user_id       UUID REFERENCES shared.users(id) ON DELETE SET NULL,

    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Content references library
CREATE TABLE movies (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    library_id  UUID NOT NULL REFERENCES movie_libraries(id) ON DELETE CASCADE,
    -- ... movie fields
);
```

---

## Example: Adult Library Schema (Isolated - QAR Obfuscation)

Adult module lives entirely in `qar` schema (Queen Anne's Revenge themed):

```sql
-- qar/000001_fleets.up.sql
SET search_path TO qar;

-- Libraries → Fleets
CREATE TABLE fleets (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name            VARCHAR(255) NOT NULL,
    fleet_type      VARCHAR(20) NOT NULL CHECK (fleet_type IN ('expedition', 'voyage')),
    paths           TEXT[] NOT NULL,

    -- Adult-specific settings
    stashdb_endpoint    TEXT DEFAULT 'https://stashdb.org/graphql',
    tpdb_enabled        BOOLEAN NOT NULL DEFAULT true,
    whisparr_sync       BOOLEAN NOT NULL DEFAULT false,
    auto_tag_crew       BOOLEAN NOT NULL DEFAULT true,  -- performers → crew

    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Movies → Expeditions
CREATE TABLE expeditions (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    fleet_id    UUID NOT NULL REFERENCES fleets(id) ON DELETE CASCADE,
    -- ... expedition (movie) fields
);

-- Scenes → Voyages
CREATE TABLE voyages (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    fleet_id    UUID NOT NULL REFERENCES fleets(id) ON DELETE CASCADE,
    -- ... voyage (scene) fields
);
```

See [ADULT_CONTENT_SYSTEM.md](../adult/ADULT_CONTENT_SYSTEM.md) for full obfuscation mapping.

---

## Cross-Module Access Control (Polymorphic)

**No central registry.** Permissions are polymorphic - each permission knows its target:

```sql
-- shared/000014_permissions.up.sql
CREATE TABLE permissions (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    -- Polymorphic reference (permission knows what it's for)
    resource_type   VARCHAR(50) NOT NULL,   -- 'movie_library', 'tv_library', 'qar.fleet'
    resource_id     UUID NOT NULL,          -- UUID of the actual resource

    -- Permission level
    permission      VARCHAR(50) NOT NULL,   -- 'view', 'manage', 'admin'

    granted_by      UUID REFERENCES users(id),
    granted_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE (user_id, resource_type, resource_id, permission)
);

CREATE INDEX idx_permissions_user ON permissions(user_id);
CREATE INDEX idx_permissions_resource ON permissions(resource_type, resource_id);
```

### Why Polymorphic?

1. **No registry to maintain** - Modules don't register/unregister
2. **Permission owns the reference** - Self-describing, no joins needed
3. **Works for any resource** - Libraries, items, playlists, etc.
4. **Module isolation** - Each module handles its own FK validation
5. **No giant join table** - Smaller, faster queries

### Usage Pattern

```go
// Check permission - module validates resource exists
func (s *MovieModule) CanAccess(ctx context.Context, userID, libraryID uuid.UUID) (bool, error) {
    // Check if user owns the library
    lib, err := s.repo.GetLibrary(ctx, libraryID)
    if err != nil {
        return false, err
    }
    if lib.OwnerUserID == userID {
        return true, nil
    }

    // Check polymorphic permission
    return s.permissions.HasPermission(ctx, userID, "movie_library", libraryID, "view")
}

// List accessible libraries - query module table + filter by permissions
func (s *MovieModule) ListLibraries(ctx context.Context, userID uuid.UUID) ([]Library, error) {
    // Get all libraries user owns OR has permission for
    return s.repo.ListAccessibleLibraries(ctx, userID)
}
```

### Listing All Libraries (UI)

For unified listing, query each enabled module in parallel:

```go
func (s *LibraryService) ListAllLibraries(ctx context.Context, userID uuid.UUID) ([]LibraryInfo, error) {
    var results []LibraryInfo
    var mu sync.Mutex
    g, ctx := errgroup.WithContext(ctx)

    for _, module := range s.enabledModules {
        module := module
        g.Go(func() error {
            libs, err := module.ListLibraries(ctx, userID)
            if err != nil {
                return err
            }
            mu.Lock()
            results = append(results, libs...)
            mu.Unlock()
            return nil
        })
    }

    if err := g.Wait(); err != nil {
        return nil, err
    }
    return results, nil
}
```

---

## Module Implementation

Each module implements a `LibraryProvider` interface:

```go
// internal/module/interface.go
type LibraryProvider interface {
    // List libraries for this module
    ListLibraries(ctx context.Context, userID uuid.UUID) ([]LibraryInfo, error)

    // Create a new library
    CreateLibrary(ctx context.Context, req CreateLibraryRequest) (*LibraryInfo, error)

    // Delete a library
    DeleteLibrary(ctx context.Context, libraryID uuid.UUID) error

    // Scan a library
    ScanLibrary(ctx context.Context, libraryID uuid.UUID) error
}

// LibraryInfo is the common interface for all library types
type LibraryInfo struct {
    ID        uuid.UUID
    Module    string   // "movie", "tvshow", "music", "qar"
    Name      string
    Paths     []string
    IsAdult   bool
    ItemCount int64
}
```

---

## Migration Required

⚠️ **STATUS**: Pending refactor - current migrations use shared library table

### Current State (To Be Changed)

The current implementation uses a shared `libraries` table:
- `shared/000005_libraries.up.sql` - Creates shared `libraries` table with `library_type` enum
- `movie/000001_movie_core.up.sql` - References `libraries(id)` (should be `movie_libraries(id)`)
- `tvshow/000001_tvshow_core.up.sql` - References `libraries(id)` (should be `tv_libraries(id)`)

### Migration Plan

1. **Create per-module library tables**:
   - `movie/000005_movie_libraries.up.sql` → `movie_libraries` table
   - `tvshow/000005_tv_libraries.up.sql` → `tv_libraries` table
   - `music/000001_music_libraries.up.sql` → `music_libraries` table
   - (Similar for audiobook, book, podcast, photo, livetv, comics)
   - `qar/000001_fleets.up.sql` → `qar.fleets` table (adult)

2. **Update content table FKs**:
   - `movies.library_id` → `REFERENCES movie_libraries(id)`
   - `series.library_id` → `REFERENCES tv_libraries(id)`

3. **Deprecate shared library table**:
   - Add `shared/000020_deprecate_libraries.up.sql` - Migration to move data
   - Eventually remove `shared/000005_libraries.up.sql`

4. **Update polymorphic permissions**:
   - `resource_type` values: `'movie_library'`, `'tv_library'`, `'qar.fleet'`, etc.

### Why This Change?

The per-module approach provides:
- Full module isolation (no shared enum to modify)
- Module-specific library settings (e.g., `tmdb_enabled` for movies, `lidarr_sync` for music)
- Independent module deployment
- Schema-level isolation for adult content (`qar` schema)

See: [MODULE_IMPLEMENTATION_TODO.md](../../planning/MODULE_IMPLEMENTATION_TODO.md)

---

---

## Implementation Checklist

### Phase 1: Core Infrastructure
- [ ] Create package structure `internal/module/` for module interfaces
- [ ] Define `LibraryProvider` interface with standard methods
- [ ] Define `LibraryInfo` common struct for cross-module listing
- [ ] Define per-module library entities:
  - [ ] `MovieLibrary` with TMDB/IMDB settings
  - [ ] `TVLibrary` with series-specific settings
  - [ ] `MusicLibrary` with Lidarr/MusicBrainz settings
  - [ ] `AudiobookLibrary` with chapter detection settings
  - [ ] `BookLibrary` with metadata settings
  - [ ] `PodcastLibrary` with RSS feed settings
  - [ ] `PhotoLibrary` with album settings
  - [ ] `LiveTVSource` with EPG settings
  - [ ] `ComicLibrary` with ComicVine settings
  - [ ] `qar.Fleet` for adult content (isolated schema)
- [ ] Create repository interfaces for each module library

### Phase 2: Database
- [ ] Create migration `movie/000005_movie_libraries.up.sql`
- [ ] Create migration `tvshow/000005_tv_libraries.up.sql`
- [ ] Create migration `music/000001_music_libraries.up.sql`
- [ ] Create migration `audiobook/000001_audiobook_libraries.up.sql`
- [ ] Create migration `book/000001_book_libraries.up.sql`
- [ ] Create migration `podcast/000001_podcast_libraries.up.sql`
- [ ] Create migration `photo/000001_photo_libraries.up.sql`
- [ ] Create migration `livetv/000001_livetv_sources.up.sql`
- [ ] Create migration `comics/000001_comic_libraries.up.sql`
- [ ] Create migration `qar/000001_fleets.up.sql` (adult, isolated schema)
- [ ] Update content table FKs to reference per-module library tables
- [ ] Create `shared/000020_deprecate_libraries.up.sql` data migration
- [ ] Create polymorphic `permissions` table for cross-module access
- [ ] Add indexes on owner_user_id, paths for efficient queries
- [ ] Generate sqlc queries for each module library CRUD

### Phase 3: Service Layer
- [ ] Implement `MovieModule.LibraryProvider` interface
- [ ] Implement `TVShowModule.LibraryProvider` interface
- [ ] Implement `MusicModule.LibraryProvider` interface
- [ ] Implement remaining module LibraryProvider interfaces
- [ ] Implement `LibraryService.ListAllLibraries()` with parallel module queries
- [ ] Implement permission checking with polymorphic grants
- [ ] Add caching for library metadata (Redis)
- [ ] Implement library scan triggering per module

### Phase 4: Background Jobs
- [ ] Create River job per module for library scanning
- [ ] Create River job for library metadata refresh
- [ ] Create job for storage usage calculation
- [ ] Create job for stale content detection

### Phase 5: API Integration
- [ ] Add OpenAPI schema for unified library endpoints
- [ ] Implement `GET /api/v1/libraries` (unified listing across modules)
- [ ] Implement `GET /api/v1/libraries/:module/:id` (module-specific detail)
- [ ] Implement `POST /api/v1/libraries/:module` (create library)
- [ ] Implement `PUT /api/v1/libraries/:module/:id` (update library)
- [ ] Implement `DELETE /api/v1/libraries/:module/:id` (delete library)
- [ ] Implement `POST /api/v1/libraries/:module/:id/scan` (trigger scan)
- [ ] Add authentication and authorization (owner + permission grants)
- [ ] Add RBAC permissions (`library.create`, `library.manage`, `library.scan`)

---


<!-- SOURCE-BREADCRUMBS-START -->

## Sources & Cross-References

> Auto-generated section linking to external documentation sources

### Cross-Reference Indexes

- [All Sources Index](../../../sources/SOURCES_INDEX.md) - Complete list of external documentation
- [Design ↔ Sources Map](../../../sources/DESIGN_CROSSREF.md) - Which docs reference which sources

### Referenced Sources

| Source | Documentation |
|--------|---------------|
| [StashDB GraphQL API](https://stashdb.org/graphql) | [Local](../../../sources/apis/stashdb-schema.graphql) |

<!-- SOURCE-BREADCRUMBS-END -->

<!-- DESIGN-BREADCRUMBS-START -->

## Related Design Docs

> Auto-generated cross-references to related design documentation

**Category**: [Shared](INDEX.md)

### In This Section

- [Time-Based Access Controls](ACCESS_CONTROLS.md)
- [Tracearr Analytics Service](ANALYTICS_SERVICE.md)
- [Revenge - Client Support & Device Capabilities](CLIENT_SUPPORT.md)
- [Content Rating System](CONTENT_RATING.md)
- [Revenge - Internationalization (i18n)](I18N.md)
- [News System](NEWS_SYSTEM.md)
- [Revenge - NSFW Toggle](NSFW_TOGGLE.md)
- [Dynamic RBAC with Casbin](RBAC_CASBIN.md)

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

## Related

- [CONTENT_RATING.md](CONTENT_RATING.md) - Age restriction and rating systems
- [Adult Content System](../adult/ADULT_CONTENT_SYSTEM.md) - Adult module isolation
- [Architecture](../../architecture/01_ARCHITECTURE.md) - System architecture
- [RBAC with Casbin](RBAC_CASBIN.md) - Permission system
