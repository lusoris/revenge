# Adult Content Module Development Instructions

> Special instructions for developing adult content modules.
> All adult content is completely isolated in a separate PostgreSQL schema `c` (obscured name).

## Schema Isolation

All adult content tables MUST be in the `c` PostgreSQL schema:

```sql
CREATE SCHEMA IF NOT EXISTS c;

-- All tables in c schema
c.movies
c.scenes
c.performers
c.studios
c.tags
-- etc.
```

## API Namespace

Adult content uses obscured API namespace `/c/`:

```
/api/v1/c/movies
/api/v1/c/movies/{id}
/api/v1/c/shows
/api/v1/c/performers
```

> **Security:** `/c/` endpoints require special auth scope, are not listed in public API docs, have separate rate limiting, and all access is audit-logged.

## Why Full Isolation?

1. **Legal compliance** - Clear data separation for regulations
2. **Backup flexibility** - `pg_dump -n c` or exclude from backups
3. **Access control** - PostgreSQL GRANT per schema
4. **Easy purge** - `DROP SCHEMA c CASCADE` removes everything
5. **No data leakage** - No FK references to public schema
6. **Separate images** - Adult images completely isolated
7. **Obscured namespace** - `/c/` and schema `c` for discretion

## Module Structure

```
internal/
  content/
    c/                      # Obscured directory name
      movie/
        entity.go
        repository.go
        repository_pg.go
        service.go
        handler.go
        scanner.go
        jobs.go             # River job definitions
        module.go
      show/
        entity.go
        repository.go
        ...
      shared/
        performer.go       # Shared between c/movie and c/show
        studio.go          # Shared between c/movie and c/show
        repository_performer.go
        repository_studio.go
```

## Database Tables

### Adult Movies

```sql
-- Core
c.movies (
    id UUID PRIMARY KEY,
    library_id UUID NOT NULL,  -- FK to public.libraries but not enforced
    title VARCHAR(500) NOT NULL,
    path TEXT NOT NULL,
    runtime_ticks BIGINT,
    release_date DATE,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ
)

c.scenes (
    id UUID PRIMARY KEY,
    movie_id UUID REFERENCES c.movies(id) ON DELETE CASCADE,
    title VARCHAR(500),
    start_ticks BIGINT,
    end_ticks BIGINT
)
```

### Performers (shared c/movie + c/show)

```sql
c.performers (
    id UUID PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    aliases TEXT[],
    gender VARCHAR(20),
    birth_date DATE,
    ethnicity VARCHAR(50),
    height_cm INT,
    measurements VARCHAR(20),
    tattoos TEXT,
    piercings TEXT,
    career_start INT,
    career_end INT,
    bio TEXT,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ
)

c.movie_performers (
    movie_id UUID REFERENCES c.movies(id) ON DELETE CASCADE,
    performer_id UUID REFERENCES c.performers(id) ON DELETE CASCADE,
    PRIMARY KEY (movie_id, performer_id)
)

c.scene_performers (
    scene_id UUID REFERENCES c.scenes(id) ON DELETE CASCADE,
    performer_id UUID REFERENCES c.performers(id) ON DELETE CASCADE,
    PRIMARY KEY (scene_id, performer_id)
)
```

### Studios

```sql
c.studios (
    id UUID PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    logo_path TEXT,
    website TEXT,
    created_at TIMESTAMPTZ
)

c.movie_studios (
    movie_id UUID REFERENCES c.movies(id) ON DELETE CASCADE,
    studio_id UUID REFERENCES c.studios(id) ON DELETE CASCADE,
    PRIMARY KEY (movie_id, studio_id)
)
```

### Tags (own taxonomy, not genres)

```sql
c.tags (
    id UUID PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    category VARCHAR(50)  -- 'act', 'attribute', 'setting', etc.
)

c.movie_tags (
    movie_id UUID REFERENCES c.movies(id) ON DELETE CASCADE,
    tag_id UUID REFERENCES c.tags(id) ON DELETE CASCADE,
    PRIMARY KEY (movie_id, tag_id)
)

c.scene_tags (
    scene_id UUID REFERENCES c.scenes(id) ON DELETE CASCADE,
    tag_id UUID REFERENCES c.tags(id) ON DELETE CASCADE,
    PRIMARY KEY (scene_id, tag_id)
)
```

### User Data (all isolated)

```sql
c.movie_user_ratings (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,  -- References public.users but not enforced
    movie_id UUID REFERENCES c.movies(id) ON DELETE CASCADE,
    score DECIMAL(3,1) NOT NULL,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    UNIQUE(user_id, movie_id)
)

c.scene_user_ratings (...)
c.performer_user_ratings (...)
c.studio_user_ratings (...)

c.movie_favorites (user_id, movie_id, added_at)
c.scene_favorites (...)
c.performer_favorites (...)
c.studio_favorites (...)

c.movie_history (user_id, movie_id, position_ticks, completed, watched_at)
```

### Images (isolated)

```sql
c.images (
    id UUID PRIMARY KEY,
    item_type VARCHAR(20) NOT NULL,  -- 'movie', 'scene', 'performer', 'studio'
    item_id UUID NOT NULL,
    image_type VARCHAR(20) NOT NULL, -- 'poster', 'backdrop', 'profile'
    path TEXT NOT NULL,
    width INT,
    height INT,
    blurhash TEXT,
    created_at TIMESTAMPTZ
)
```

### Playlists & Collections (isolated)

```sql
c.playlists (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    name VARCHAR(200) NOT NULL,
    is_public BOOLEAN DEFAULT false,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ
)

c.playlist_items (
    id UUID PRIMARY KEY,
    playlist_id UUID REFERENCES c.playlists(id) ON DELETE CASCADE,
    item_type VARCHAR(20) NOT NULL,  -- 'movie', 'scene', 'episode'
    item_id UUID NOT NULL,
    position INT NOT NULL
)

c.collections (...)
c.collection_items (...)
```

## Entity Design

```go
// content/c/movie/entity.go
package movie

type Movie struct {
    ID           uuid.UUID
    LibraryID    uuid.UUID
    Title        string
    Path         string
    RuntimeTicks int64
    ReleaseDate  *time.Time
    CreatedAt    time.Time
    UpdatedAt    time.Time

    // Relationships
    Scenes     []Scene
    Performers []Performer
    Studios    []Studio
    Tags       []Tag
}

type Scene struct {
    ID         uuid.UUID
    MovieID    uuid.UUID
    Title      string
    StartTicks int64
    EndTicks   int64

    Performers []Performer
    Tags       []Tag
}
```

```go
// content/c/shared/performer.go
package shared

type Performer struct {
    ID           uuid.UUID
    Name         string
    Aliases      []string
    Gender       string
    BirthDate    *time.Time
    Ethnicity    string
    HeightCm     *int
    Measurements string
    Tattoos      string
    Piercings    string
    CareerStart  *int
    CareerEnd    *int
    Bio          string
    CreatedAt    time.Time
    UpdatedAt    time.Time
}
```

## Access Control

EVERY adult handler MUST verify user has adult content access:

```go
// content/c/movie/handler.go
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    user := middleware.UserFromContext(ctx)

    // REQUIRED: Check adult access scope
    if !user.HasScope("adult:read") {
        handlers.Forbidden(w, "Adult content access not enabled")
        return
    }

    // Proceed with handler logic
}
```

Consider middleware for adult routes:

```go
func (h *Handler) RegisterRoutes(mux *http.ServeMux, auth *middleware.Auth) {
    adultAuth := auth.ScopeRequired("adult:read")

    mux.Handle("GET /api/v1/c/movies", auth.Required(adultAuth(http.HandlerFunc(h.List))))
    mux.Handle("GET /api/v1/c/movies/{id}", auth.Required(adultAuth(http.HandlerFunc(h.Get))))
    // ...
}
```

## API Routes

Adult content uses obscured `/c/` namespace:

```
/api/v1/c/
  /movies
  /movies/{id}
  /movies/{id}/scenes
  /movies/{id}/performers
  /shows
  /shows/{id}
  /shows/{id}/seasons
  /shows/{id}/seasons/{seasonId}/episodes
  /performers
  /performers/{id}
  /performers/{id}/movies
  /studios
  /studios/{id}
  /tags
  /playlists
```

## No External Ratings

Adult modules do NOT have external ratings (no IMDb, etc.):

- Only user ratings
- No sync services
- Privacy by default

## No Content Ratings

Adult content has no age restriction ratings:

- Implicit 18+/adult-only
- Access controlled by scope
- No MPAA/FSK/etc. needed

## Search Isolation

Adult content is NOT included in unified search:

- Separate Typesense collections: `c_movies`, `c_series`
- Separate search endpoint: `/api/v1/c/search`
- Requires `adult:read` scope

## Testing

Tests for adult modules should:

1. Test access control (user without adult scope = 403)
2. Test data isolation (no cross-schema queries)
3. Use separate test fixtures

```go
func TestAdultMovieHandler_RequiresAdultScope(t *testing.T) {
    // Create user WITHOUT adult scope
    user := &domain.User{Scopes: []string{"read"}}  // No "adult:read"

    req := httptest.NewRequest("GET", "/api/v1/c/movies", nil)
    req = req.WithContext(middleware.ContextWithUser(req.Context(), user))

    rr := httptest.NewRecorder()
    handler.List(rr, req)

    assert.Equal(t, http.StatusForbidden, rr.Code)
}
```
---

## Related

- [INDEX.instructions.md](INDEX.instructions.md) - Main instruction index with all cross-references
- [ARCHITECTURE_V2.md](../../docs/dev/design/architecture/ARCHITECTURE_V2.md) - System architecture
- [BEST_PRACTICES.md](../../docs/dev/design/operations/BEST_PRACTICES.md) - Best practices
