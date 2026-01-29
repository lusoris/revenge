# Library Types

> Per-module library architecture and supported content types

**Status**: 🟡 DESIGNED (needs migration update)
**Current Migration**: `shared/000005_libraries.up.sql` ⚠️ **TO BE SPLIT**

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

⚠️ **TODO**: Update `shared/000005_libraries.up.sql`

The current migration creates a shared `libraries` table with type enum. This needs to be:
1. **Removed** from shared migrations (delete `shared/000005_libraries.up.sql`)
2. **Split** into per-module library tables (e.g., `movie/000001_movie_libraries.up.sql`)
3. **Add** `resource_grants` table for polymorphic access control (already defined in `shared/000019_resource_grants.up.sql`)

See: [MODULE_IMPLEMENTATION_TODO.md](../../planning/MODULE_IMPLEMENTATION_TODO.md)

---

## See Also

- [CONTENT_RATING.md](CONTENT_RATING.md) - Age restriction and rating systems
- [Adult Content System](../adult/ADULT_CONTENT_SYSTEM.md) - Adult module isolation
- [Architecture V2](../../architecture/ARCHITECTURE_V2.md) - System architecture
