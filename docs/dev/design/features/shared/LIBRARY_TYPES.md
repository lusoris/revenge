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
│                 │                 │                 │    (c schema)     │
│ ┌─────────────┐ │ ┌─────────────┐ │ ┌─────────────┐ │ ┌───────────────┐ │
│ │movie_       │ │ │tv_          │ │ │music_       │ │ │c.adult_       │ │
│ │libraries    │ │ │libraries    │ │ │libraries    │ │ │libraries      │ │
│ └──────┬──────┘ │ └──────┬──────┘ │ └──────┬──────┘ │ └───────┬───────┘ │
│        │        │        │        │        │        │         │         │
│        ▼        │        ▼        │        ▼        │         ▼         │
│ ┌─────────────┐ │ ┌─────────────┐ │ ┌─────────────┐ │ ┌───────────────┐ │
│ │   movies    │ │ │   series    │ │ │   albums    │ │ │c.adult_movies │ │
│ └─────────────┘ │ │   seasons   │ │ │   tracks    │ │ │c.adult_scenes │ │
│                 │ │   episodes  │ │ │   artists   │ │ │c.performers   │ │
│                 │ └─────────────┘ │ └─────────────┘ │ └───────────────┘ │
└─────────────────┴─────────────────┴─────────────────┴───────────────────┘
```

### Why Per-Module Libraries?

1. **Full module isolation** - Each module owns its complete data model
2. **Module-specific settings** - Library settings baked into each schema
3. **No shared enum** - Modules can be added without modifying shared code
4. **Independent deployment** - Modules can be enabled/disabled cleanly
5. **Schema-level isolation** - Adult content fully contained in `c` schema

### Shared Components

Only truly cross-cutting concerns remain shared:
- `shared.users` - User accounts
- `shared.profiles` - User profiles
- `shared.sessions` - Auth sessions
- `shared.library_access` - Cross-module access control (references module libraries by UUID + type)

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
| adult | `c.adult_libraries` | `c.adult_movies`, `c.adult_scenes` | `c/000001_*.sql` |

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

## Example: Adult Library Schema (Isolated)

Adult module lives entirely in `c` schema:

```sql
-- c/000001_adult_libraries.up.sql
SET search_path TO c;

CREATE TABLE adult_libraries (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name            VARCHAR(255) NOT NULL,
    library_type    VARCHAR(20) NOT NULL CHECK (library_type IN ('movie', 'scene')),
    paths           TEXT[] NOT NULL,

    -- Adult-specific settings
    stashdb_endpoint    TEXT DEFAULT 'https://stashdb.org/graphql',
    tpdb_enabled        BOOLEAN NOT NULL DEFAULT true,
    whisparr_sync       BOOLEAN NOT NULL DEFAULT false,
    auto_tag_performers BOOLEAN NOT NULL DEFAULT true,

    -- Always adult content
    content_rating      VARCHAR(10) NOT NULL DEFAULT 'XXX',

    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Content in same schema
CREATE TABLE adult_movies (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    library_id  UUID NOT NULL REFERENCES adult_libraries(id) ON DELETE CASCADE,
    -- ... adult movie fields
);
```

---

## Cross-Module Access Control

For unified library listing in UI, use a shared registry:

```sql
-- shared/000020_library_registry.up.sql
-- Lightweight registry for cross-module access control
CREATE TABLE library_registry (
    id              UUID PRIMARY KEY,           -- Same as module library ID
    module          VARCHAR(50) NOT NULL,       -- 'movie', 'tvshow', 'c.adult'
    name            VARCHAR(255) NOT NULL,      -- Cached for listing
    is_adult        BOOLEAN NOT NULL DEFAULT false,
    owner_user_id   UUID REFERENCES users(id),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_library_registry_module ON library_registry(module);
CREATE INDEX idx_library_registry_owner ON library_registry(owner_user_id);

-- Access grants (who can see what)
CREATE TABLE library_access (
    library_id      UUID NOT NULL REFERENCES library_registry(id) ON DELETE CASCADE,
    user_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    can_manage      BOOLEAN NOT NULL DEFAULT false,
    PRIMARY KEY (library_id, user_id)
);
```

Modules register/unregister on library create/delete via triggers or service calls.

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
    Module    string   // "movie", "tvshow", "music", "c.adult"
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
1. **Removed** from shared migrations
2. **Split** into per-module library tables
3. **Add** lightweight `library_registry` for cross-module access

See: [MODULE_IMPLEMENTATION_TODO.md](../../planning/MODULE_IMPLEMENTATION_TODO.md)

---

## See Also

- [CONTENT_RATING.md](CONTENT_RATING.md) - Age restriction and rating systems
- [Adult Content System](../adult/ADULT_CONTENT_SYSTEM.md) - Adult module isolation
- [Architecture V2](../../architecture/ARCHITECTURE_V2.md) - System architecture
