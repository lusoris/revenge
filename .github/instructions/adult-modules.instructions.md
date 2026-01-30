# Adult Content Module Development Instructions

> Special instructions for developing adult content modules.
> All adult content is completely isolated in a separate PostgreSQL schema `qar` (Queen Anne's Revenge obfuscation).

## Queen Anne's Revenge (QAR) Terminology

All adult content uses pirate ship terminology themed after Blackbeard's famous vessel:

| Real Entity | QAR Term | Database Table | API Endpoint |
|-------------|----------|----------------|--------------|
| Libraries | **Fleets** | `qar.fleets` | `/qar/fleets` |
| Movies | **Expeditions** | `qar.expeditions` | `/qar/expeditions` |
| Scenes | **Voyages** | `qar.voyages` | `/qar/voyages` |
| Performers | **Crew** | `qar.crew` | `/qar/crew` |
| Studios | **Ports** | `qar.ports` | `/qar/ports` |
| Tags | **Flags** | `qar.flags` | `/qar/flags` |
| Categories | **Waters** | Field in flags | - |

### Field Obfuscation

| Real Field | QAR Term | Used In |
|------------|----------|---------|
| birth_date | `christening` | crew |
| ethnicity | `origin` | crew |
| hair_color | `rigging` | crew |
| eye_color | `compass` | crew |
| career_start | `maiden_voyage` | crew |
| career_end | `last_port` | crew |
| tattoos | `markings` | crew |
| piercings | `anchors` | crew |
| stashdb_id | `charter` | all entities |
| tpdb_id | `registry` | all entities |
| freeones_id | `manifest` | crew |
| release_date | `launch_date` | expeditions, voyages |
| runtime | `distance` | expeditions, voyages |
| phash | `coordinates` | voyages |

## Schema Isolation

All adult content tables MUST be in the `qar` PostgreSQL schema:

```sql
CREATE SCHEMA IF NOT EXISTS qar;

-- All tables in qar schema
qar.fleets        -- Libraries
qar.expeditions   -- Movies
qar.voyages       -- Scenes
qar.crew          -- Performers
qar.ports         -- Studios
qar.flags         -- Tags
-- etc.
```

## API Namespace

Adult content uses the QAR namespace `/qar/`:

```
/api/v1/qar/fleets
/api/v1/qar/expeditions
/api/v1/qar/expeditions/{id}
/api/v1/qar/voyages
/api/v1/qar/crew
/api/v1/qar/ports
/api/v1/qar/flags
```

> **Security:** `/qar/` endpoints require special auth scope, are not listed in public API docs, have separate rate limiting, and all access is audit-logged.

## Why Full Isolation?

1. **Legal compliance** - Clear data separation for regulations
2. **Backup flexibility** - `pg_dump -n qar` or exclude from backups
3. **Access control** - PostgreSQL GRANT per schema
4. **Easy purge** - `DROP SCHEMA qar CASCADE` removes everything
5. **No data leakage** - No FK references to public schema
6. **Separate images** - Adult images completely isolated
7. **Obscured namespace** - `/qar/` and schema `qar` for discretion

## Module Structure

```
internal/
  content/
    qar/                     # Queen Anne's Revenge namespace
      expedition/            # Movies
        entity.go
        repository.go
        repository_sqlc.go
        service.go
        handler.go
        scanner.go
        jobs.go
        module.go
      voyage/                # Scenes
        entity.go
        repository.go
        ...
      crew/                  # Performers
        entity.go
        repository.go
        ...
      port/                  # Studios
        entity.go
        repository.go
        ...
      flag/                  # Tags
        entity.go
        repository.go
        ...
      fleet/                 # Libraries
        entity.go
        repository.go
        ...
      shared/
        content_entity.go    # Common embedded struct
        interfaces.go        # Shared interfaces
```

## Database Tables

### Fleets (Libraries)

```sql
qar.fleets (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    fleet_type VARCHAR(20) NOT NULL,  -- 'expedition', 'voyage'
    paths TEXT[] NOT NULL,
    stashdb_endpoint TEXT DEFAULT 'https://stashdb.org/graphql',
    owner_user_id UUID,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ
)
```

### Expeditions (Movies)

```sql
qar.expeditions (
    id UUID PRIMARY KEY,
    fleet_id UUID NOT NULL REFERENCES qar.fleets(id),
    title VARCHAR(500) NOT NULL,
    sort_title VARCHAR(500),
    overview TEXT,
    launch_date DATE,       -- release_date
    distance INT,           -- runtime_minutes
    port_id UUID REFERENCES qar.ports(id),
    charter VARCHAR(100),   -- stashdb_id
    registry VARCHAR(100),  -- tpdb_id
    whisparr_id INT,
    cover_path TEXT,
    path TEXT NOT NULL,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ
)
```

### Voyages (Scenes)

```sql
qar.voyages (
    id UUID PRIMARY KEY,
    fleet_id UUID NOT NULL REFERENCES qar.fleets(id),
    title VARCHAR(500) NOT NULL,
    overview TEXT,
    launch_date DATE,       -- release_date
    distance INT,           -- runtime_minutes
    port_id UUID REFERENCES qar.ports(id),
    coordinates VARCHAR(64), -- phash
    oshash VARCHAR(32),
    md5 VARCHAR(64),
    charter VARCHAR(100),   -- stashdb_id
    registry VARCHAR(100),  -- tpdb_id
    stash_id VARCHAR(100),
    whisparr_id INT,
    cover_path TEXT,
    path TEXT NOT NULL,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ
)
```

### Crew (Performers)

```sql
qar.crew (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    disambiguation VARCHAR(255),
    gender VARCHAR(50),
    christening DATE,       -- birth_date
    death_date DATE,
    birth_city VARCHAR(255),
    origin VARCHAR(100),    -- ethnicity
    nationality VARCHAR(100),
    rigging VARCHAR(50),    -- hair_color
    compass VARCHAR(50),    -- eye_color
    height_cm INT,
    weight_kg INT,
    measurements VARCHAR(50),
    cup_size VARCHAR(10),
    breast_type VARCHAR(50),
    markings TEXT,          -- tattoos
    anchors TEXT,           -- piercings
    maiden_voyage INT,      -- career_start
    last_port INT,          -- career_end
    bio TEXT,
    stash_id VARCHAR(100),
    charter VARCHAR(100),   -- stashdb_id
    registry VARCHAR(100),  -- tpdb_id
    manifest VARCHAR(100),  -- freeones_id
    twitter TEXT,
    instagram TEXT,
    image_path TEXT,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ
)

qar.crew_names (
    crew_id UUID REFERENCES qar.crew(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    PRIMARY KEY (crew_id, name)
)
```

### Ports (Studios)

```sql
qar.ports (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    parent_id UUID REFERENCES qar.ports(id),  -- Network/parent
    stashdb_id VARCHAR(100),
    tpdb_id VARCHAR(100),
    url TEXT,
    logo_path TEXT,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ
)
```

### Flags (Tags)

```sql
qar.flags (
    id UUID PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    parent_id UUID REFERENCES qar.flags(id),
    stashdb_id VARCHAR(36),
    waters VARCHAR(50),    -- category
    created_at TIMESTAMPTZ
)
```

### Junction Tables

```sql
qar.expedition_crew (
    expedition_id UUID REFERENCES qar.expeditions(id) ON DELETE CASCADE,
    crew_id UUID REFERENCES qar.crew(id) ON DELETE CASCADE,
    PRIMARY KEY (expedition_id, crew_id)
)

qar.voyage_crew (
    voyage_id UUID REFERENCES qar.voyages(id) ON DELETE CASCADE,
    crew_id UUID REFERENCES qar.crew(id) ON DELETE CASCADE,
    PRIMARY KEY (voyage_id, crew_id)
)

qar.expedition_flags (
    expedition_id UUID REFERENCES qar.expeditions(id) ON DELETE CASCADE,
    flag_id UUID REFERENCES qar.flags(id) ON DELETE CASCADE,
    PRIMARY KEY (expedition_id, flag_id)
)

qar.voyage_flags (
    voyage_id UUID REFERENCES qar.voyages(id) ON DELETE CASCADE,
    flag_id UUID REFERENCES qar.flags(id) ON DELETE CASCADE,
    PRIMARY KEY (voyage_id, flag_id)
)
```

## Entity Design

```go
// content/qar/voyage/entity.go
package voyage

import (
    "github.com/lusoris/revenge/internal/content/shared"
)

// Voyage represents an adult scene (obfuscated as "voyage").
type Voyage struct {
    shared.ContentEntity
    FleetID     uuid.UUID   // Library reference
    LaunchDate  *time.Time  // release_date
    Distance    int         // runtime_minutes
    Overview    string
    PortID      *uuid.UUID  // studio_id
    Coordinates string      // phash
    Oshash      string
    MD5         string
    CoverPath   string
    Charter     string      // stashdb_id
    Registry    string      // tpdb_id
    StashID     string
    WhisparrID  *int
}
```

```go
// content/qar/crew/entity.go
package crew

// Crew represents an adult performer (obfuscated as "crew").
type Crew struct {
    ID             uuid.UUID
    Name           string
    Disambiguation string
    Gender         string
    Christening    *time.Time  // birth_date
    DeathDate      *time.Time
    BirthCity      string
    Origin         string      // ethnicity
    Nationality    string
    Rigging        string      // hair_color
    Compass        string      // eye_color
    HeightCM       *int
    WeightKG       *int
    Measurements   string
    CupSize        string
    BreastType     string
    Markings       string      // tattoos
    Anchors        string      // piercings
    MaidenVoyage   *int        // career_start
    LastPort       *int        // career_end
    Bio            string
    StashID        string
    Charter        string      // stashdb_id
    Registry       string      // tpdb_id
    Manifest       string      // freeones_id
    Twitter        string
    Instagram      string
    ImagePath      string
    CreatedAt      time.Time
    UpdatedAt      time.Time
}
```

## Access Control

EVERY adult handler MUST verify user has adult content access:

```go
// content/qar/voyage/handler.go
func (h *Handler) Get(ctx context.Context, params api.GetVoyageParams) (*api.Voyage, error) {
    // REQUIRED: Check adult access via handler helper
    user, err := h.handler.requireAdultBrowse(ctx)
    if err != nil {
        return nil, err
    }

    // Proceed with handler logic
    voyage, err := h.service.GetByID(ctx, params.ID)
    // ...
}
```

RBAC permissions for adult content:

```go
// Permission constants
const (
    PermAdultBrowse            = "adult:browse"
    PermAdultStream            = "adult:stream"
    PermAdultMetadataWrite     = "adult:metadata:write"
    PermAdultRequestsSubmit    = "adult:requests:submit"
    PermAdultRequestsViewOwn   = "adult:requests:view_own"
    PermAdultRequestsVote      = "adult:requests:vote"
    PermAdultRequestsApprove   = "adult:requests:approve"
    PermAdultRequestsDecline   = "adult:requests:decline"
    // ... more permissions
)
```

## Fingerprinting

Adult content uses fingerprinting for scene identification:

| Algorithm | Description | Use Case |
|-----------|-------------|----------|
| **oshash** | OpenSubtitles hash | Fast file-based matching |
| **pHash** | Perceptual hash | StashDB scene matching |
| **MD5** | File content hash | Exact duplicate detection |

```go
// Fingerprint service generates hashes
type Fingerprinter interface {
    GenerateFingerprints(ctx context.Context, filePath string) (*FingerprintResult, error)
}

type FingerprintResult struct {
    Coordinates string  // phash
    Oshash      string
    MD5         string
}
```

## Metadata Sources

Adult content uses specialized metadata sources:

1. **Whisparr** - Primary acquisition proxy (Radarr fork)
2. **StashDB** - Community database for scenes/performers
3. **Stash-App** - Local instance sync
4. **TPDB** - Fallback metadata source

## Search Isolation

Adult content is NOT included in unified search:

- Separate Typesense collections: `qar_expeditions`, `qar_voyages`, `qar_crew`
- Separate search endpoint: `/api/v1/qar/search`
- Requires `adult:browse` scope

## Testing

Tests for adult modules should:

1. Test access control (user without adult scope = 403)
2. Test data isolation (no cross-schema queries)
3. Use separate test fixtures

```go
func TestVoyageHandler_RequiresAdultScope(t *testing.T) {
    // Create user WITHOUT adult scope
    user := &domain.User{AdultEnabled: false}

    req := httptest.NewRequest("GET", "/api/v1/qar/voyages", nil)
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
- [ADULT_CONTENT_SYSTEM.md](../../docs/dev/design/features/adult/ADULT_CONTENT_SYSTEM.md) - Full adult content design
- [BEST_PRACTICES.md](../../docs/dev/design/operations/BEST_PRACTICES.md) - Best practices
