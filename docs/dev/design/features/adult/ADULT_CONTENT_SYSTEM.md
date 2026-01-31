# Revenge - Adult Content System

> Complete adult content management with Stash ecosystem integration.
> All adult content isolated in PostgreSQL schema `qar` with "Queen Anne's Revenge" themed obfuscation.

## Status

| Dimension           | Status | Notes |
| ------------------- | ------ | ----- |
| Design              | ✅     |       |
| Sources             | ✅     |       |
| Instructions        | ✅     |       |
| Code                | 🔴     |       |
| Linting             | 🔴     |       |
| Unit Testing        | 🔴     |       |
| Integration Testing | 🔴     |       |

## Design Principles

1. **Complete Isolation** - Separate schema, namespace, and storage
2. **Whisparr as Proxy** - Whisparr-v3 for acquisition, local metadata enrichment
3. **Stash as Enrichment** - StashDB, Stash-App for comprehensive performer/scene data
4. **Perceptual Hashing** - pHash fingerprinting for scene identification
5. **Privacy First** - Fully obfuscated terminology, optional analytics exclusion, audit logging

---

## Queen Anne's Revenge Obfuscation

> See [00_SOURCE_OF_TRUTH.md](../../00_SOURCE_OF_TRUTH.md#qar-obfuscation-terminology) for complete QAR terminology reference.

All adult content uses pirate ship terminology themed after Blackbeard's famous vessel. The obfuscation covers:

- **Schema & API Namespace**: `qar` schema internally, `/api/v1/legacy/` externally
- **Entity Obfuscation**: Performers → Crew, Scenes → Voyages, Movies → Expeditions, etc.
- **Field Obfuscation**: birth_date → christening, tattoos → markings, etc.
- **URL Obfuscation**: External URLs never reveal content type

### Example API Response

```json
// GET /api/v1/legacy/crew/123
// Note: External URL uses /legacy/, internal schema is qar.*
{
  "id": "123",
  "title": "Anne Bonny",
  "names": ["Anne Bonny", "Anne Cormac"],
  "callsign": "pirate-queen",
  "christening": "1697-03-08",
  "voyage_count": 142,
  "charter": "stashdb-uuid-here"
}
```

### Database Schema Example

```sql
-- Schema named after Queen Anne's Revenge
CREATE SCHEMA qar;

-- Performers → Crew
CREATE TABLE qar.crew (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    names           TEXT[] NOT NULL,              -- aliases
    christening     DATE,                         -- birth_date
    origin          VARCHAR(50),                  -- ethnicity
    maiden_voyage   INT,                          -- career_start year
    last_port       INT,                          -- career_end year
    cargo           JSONB,                        -- measurements
    markings        TEXT[],                       -- tattoos
    anchors         TEXT[],                       -- piercings
    rigging         VARCHAR(20),                  -- hair_color
    compass         VARCHAR(20),                  -- eye_color
    cutlass         VARCHAR(20),                  -- penis_size (male/trans only)
    figurehead      BOOLEAN,                      -- has_breasts (trans)
    keel            VARCHAR(20),                  -- genitalia_type (male/female/both)
    refit           VARCHAR(20),                  -- surgical_status (pre-op/post-op/non-op)
    voyage_count    INT NOT NULL DEFAULT 0,       -- scene_count
    charter         VARCHAR(100),                 -- stashdb_id (obfuscated)
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Scenes → Voyages
CREATE TABLE qar.voyages (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    fleet_id        UUID NOT NULL REFERENCES qar.fleets(id),
    title           VARCHAR(500) NOT NULL,
    port_id         UUID REFERENCES qar.ports(id),
    distance        INT,                          -- duration_seconds
    bounty          DECIMAL(3,1),                 -- rating
    launch_date     DATE,                         -- release_date
    coordinates     VARCHAR(100),                 -- phash fingerprint
    oshash          VARCHAR(32),
    flags           UUID[],                       -- tag IDs
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Junction: Voyages ↔ Crew
CREATE TABLE qar.voyage_crew (
    voyage_id       UUID NOT NULL REFERENCES qar.voyages(id) ON DELETE CASCADE,
    crew_id         UUID NOT NULL REFERENCES qar.crew(id) ON DELETE CASCADE,
    role            VARCHAR(50),                  -- position in scene
    PRIMARY KEY (voyage_id, crew_id)
);

-- Studios → Ports
CREATE TABLE qar.ports (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name            VARCHAR(255) NOT NULL,
    parent_id       UUID REFERENCES qar.ports(id), -- network/parent studio
    logo_path       TEXT,
    stashdb_id      VARCHAR(100),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Tags → Flags
CREATE TABLE qar.flags (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name            VARCHAR(100) NOT NULL UNIQUE,
    category        VARCHAR(50),                  -- waters reference
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Libraries → Fleets
CREATE TABLE qar.fleets (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name            VARCHAR(255) NOT NULL,
    fleet_type      VARCHAR(20) NOT NULL CHECK (fleet_type IN ('expedition', 'voyage')),
    paths           TEXT[] NOT NULL,
    stashdb_endpoint TEXT DEFAULT 'https://stashdb.org/graphql',
    owner_user_id   UUID REFERENCES shared.users(id),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

### Code Mapping Layer

```go
// internal/content/qar/mapping.go
// Thin mapping layer for developer sanity - internally use real terms

type Performer = Crew
type Scene = Voyage
type Movie = Expedition
type Studio = Port
type Tag = Flag
type Library = Fleet

// Repository methods use obfuscated names externally
func (r *Repository) GetCrew(ctx context.Context, id uuid.UUID) (*Crew, error)
func (r *Repository) ListVoyages(ctx context.Context, fleetID uuid.UUID) ([]Voyage, error)
func (r *Repository) GetPort(ctx context.Context, id uuid.UUID) (*Port, error)
```

---

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      Adult Content Data Flow                                 │
└─────────────────────────────────────────────────────────────────────────────┘

┌──────────┐     ┌───────────────────┐     ┌─────────────────┐
│  Library │ ──→ │  Whisparr-v3      │ ──→ │  qar.* Schema   │
│   Scan   │     │  (Acquisition)    │     │  (PostgreSQL)   │
└──────────┘     └───────────────────┘     └─────────────────┘
                        │                          │
         ┌──────────────┼──────────────┐           │
         ▼              ▼              ▼           ▼
  ┌─────────────┐ ┌─────────────┐ ┌──────────┐ ┌──────────────────┐
  │   StashDB   │ │  Stash-App  │ │  TPDB    │ │  Fingerprint     │
  │  (Public)   │ │  (Private)  │ │ (Legacy) │ │  Matching        │
  └─────────────┘ └─────────────┘ └──────────┘ └──────────────────┘
         │              │              │                │
         └──────────────┴──────────────┘                │
                        │                               │
                        ▼                               ▼
              ┌──────────────────────┐      ┌────────────────────┐
              │   Performer Data     │      │   Scene Identity   │
              │   (Photos, Bio, etc) │      │   (pHash, oshash)  │
              └──────────────────────┘      └────────────────────┘
```

---

## Data Sources

### Priority Order

| Priority | Source | Type | Data Provided |
|----------|--------|------|---------------|
| 1 | **Whisparr-v3** | Proxy | Acquisition, basic metadata, studio, tags |
| 2 | **StashDB** | Enrichment | Performer details, scene fingerprints, community data |
| 3 | **Stash-App** | Enrichment | Local organization, custom tags, user collections |
| 4 | **TPDB** | Fallback | Extended metadata, legacy performer data |
| 5 | **IAD** | Fallback | European content metadata |

### Whisparr-v3 (Acquisition Proxy)

Whisparr handles content acquisition and provides initial metadata:

```go
type WhisparrClient struct {
    baseURL    string
    apiKey     string
    client     *http.Client
}

type WhisparrMovie struct {
    ID              int       `json:"id"`
    Title           string    `json:"title"`
    SortTitle       string    `json:"sortTitle"`
    Studio          Studio    `json:"studio"`
    Performers      []Person  `json:"credits"`
    Genres          []string  `json:"genres"`
    Tags            []int     `json:"tags"`
    Year            int       `json:"year"`
    ReleaseDate     string    `json:"releaseDate"`
    Runtime         int       `json:"runtime"`
    Overview        string    `json:"overview"`
    StashID         string    `json:"stashId,omitempty"` // StashDB reference
    Images          []Image   `json:"images"`
    HasFile         bool      `json:"hasFile"`
    Path            string    `json:"path"`
}

// Import from Whisparr to Revenge
func (c *WhisparrClient) ImportMovie(ctx context.Context, whisparrID int) (*AdultMovie, error) {
    movie, err := c.GetMovie(ctx, whisparrID)
    if err != nil {
        return nil, err
    }

    // Convert to Revenge domain model
    adultMovie := &AdultMovie{
        ID:          uuid.New(),
        WhisparrID:  &movie.ID,
        Title:       movie.Title,
        StudioID:    c.ensureStudio(ctx, movie.Studio),
        ReleaseDate: parseDate(movie.ReleaseDate),
        Runtime:     time.Duration(movie.Runtime) * time.Minute,
        Overview:    movie.Overview,
    }

    // Queue StashDB enrichment
    return adultMovie, nil
}
```

### StashDB Integration (Community Database)

StashDB is the canonical open-source adult content metadata database:

```go
type StashDBClient struct {
    endpoint   string
    apiKey     string
    client     *graphql.Client
}

// StashDB GraphQL Schema (simplified)
type StashDBScene struct {
    ID          string           `json:"id"`
    Title       string           `json:"title"`
    Details     string           `json:"details"`
    Date        string           `json:"date"`
    Duration    int              `json:"duration"`
    Director    string           `json:"director"`
    Code        string           `json:"code"`        // Scene code (e.g., "BZ-123")
    Studio      *StashDBStudio   `json:"studio"`
    Performers  []ScenePerformer `json:"performers"`
    Fingerprints []Fingerprint   `json:"fingerprints"`
    Tags        []StashDBTag     `json:"tags"`
    Images      []StashDBImage   `json:"images"`
    URLs        []URL            `json:"urls"`
}

type StashDBPerformer struct {
    ID              string           `json:"id"`
    Name            string           `json:"name"`
    Disambiguation  string           `json:"disambiguation"` // For same-name performers
    Aliases         []string         `json:"aliases"`
    Gender          string           `json:"gender"`
    Birthdate       *string          `json:"birthdate"`
    BirthCity       *string          `json:"birth_city"`
    Ethnicity       *string          `json:"ethnicity"`
    Height          *int             `json:"height"` // cm
    Measurements    *Measurements    `json:"measurements"`
    TattooDesc      *string          `json:"tattoo_description"`
    PiercingDesc    *string          `json:"piercing_description"`
    Career          *CareerInfo      `json:"career"`
    Images          []StashDBImage   `json:"images"`
    URLs            []URL            `json:"urls"`
}

type Fingerprint struct {
    Hash      string `json:"hash"`
    Algorithm string `json:"algorithm"` // "phash", "oshash", "md5"
    Duration  int    `json:"duration"`
}

// Query scenes by fingerprint
func (c *StashDBClient) FindSceneByFingerprint(ctx context.Context, fp Fingerprint) (*StashDBScene, error) {
    query := `
        query FindScenesByFingerprint($fingerprint: FingerprintInput!) {
            findScenesByFingerprints(fingerprints: [$fingerprint]) {
                id
                title
                date
                studio { id name }
                performers { performer { id name } }
                fingerprints { hash algorithm duration }
            }
        }
    `
    // Execute GraphQL query
    return c.queryScene(ctx, query, map[string]interface{}{
        "fingerprint": fp,
    })
}

// Query performer details
func (c *StashDBClient) GetPerformer(ctx context.Context, stashID string) (*StashDBPerformer, error) {
    query := `
        query GetPerformer($id: ID!) {
            findPerformer(id: $id) {
                id name disambiguation aliases gender
                birthdate birth_city ethnicity height
                measurements { bust waist hip cup_size }
                tattoo_description piercing_description
                career { start end }
                images { url type }
                urls { url type }
            }
        }
    `
    return c.queryPerformer(ctx, query, map[string]interface{}{"id": stashID})
}
```

### Stash-App Integration (Private Instance)

For users running Stash locally - pull their organized data:

```go
type StashAppClient struct {
    baseURL string
    apiKey  string
    client  *graphql.Client
}

// Stash-App provides local organization
type StashAppScene struct {
    ID           string         `json:"id"`
    Title        string         `json:"title"`
    Details      string         `json:"details"`
    Date         string         `json:"date"`
    Rating100    *int           `json:"rating100"`  // User rating 0-100
    OCounter     int            `json:"o_counter"`  // Play count
    Organized    bool           `json:"organized"`
    Interactive  bool           `json:"interactive"` // Supports funscript
    Performers   []Performer    `json:"performers"`
    Studio       *Studio        `json:"studio"`
    Tags         []Tag          `json:"tags"`
    Galleries    []Gallery      `json:"galleries"`
    Movies       []SceneMovie   `json:"movies"`
    StashIDs     []StashID      `json:"stash_ids"` // References to StashDB
    Paths        ScenePaths     `json:"paths"`
    SceneMarkers []SceneMarker  `json:"scene_markers"` // Timestamped tags
}

type SceneMarker struct {
    ID         string    `json:"id"`
    Title      string    `json:"title"`
    Seconds    float64   `json:"seconds"`
    PrimaryTag Tag       `json:"primary_tag"`
    Tags       []Tag     `json:"tags"`
    Screenshot string    `json:"screenshot"`
}

// Sync from user's Stash instance
func (c *StashAppClient) SyncLibrary(ctx context.Context, libraryID uuid.UUID) error {
    // Fetch all scenes from Stash
    scenes, err := c.getAllScenes(ctx)
    if err != nil {
        return err
    }

    for _, scene := range scenes {
        // Match by fingerprint or path
        // Import metadata, markers, ratings
    }
    return nil
}

// Import scene markers as chapters
func (c *StashAppClient) ImportMarkers(ctx context.Context, sceneID string) ([]Chapter, error) {
    scene, err := c.getScene(ctx, sceneID)
    if err != nil {
        return nil, err
    }

    chapters := make([]Chapter, 0, len(scene.SceneMarkers))
    for _, marker := range scene.SceneMarkers {
        chapters = append(chapters, Chapter{
            Title:       marker.Title,
            StartTicks:  int64(marker.Seconds * 10_000_000),
            ImagePath:   marker.Screenshot,
            Tags:        extractTags(marker.Tags),
        })
    }
    return chapters, nil
}
```

---

## Perceptual Hashing (Scene Identification)

Scene identification via fingerprinting - match files to metadata databases:

### Hash Types

| Algorithm | Description | Use Case |
|-----------|-------------|----------|
| **pHash** | Perceptual hash of video frames | StashDB scene matching |
| **oshash** | OpenSubtitles hash (file-based) | Fast local matching |
| **MD5** | File content hash | Exact duplicate detection |
| **VideoHash** | Combined frame + audio hash | Fuzzy scene matching |

### Implementation

```go
type FingerprintService struct {
    ffprobePath string
    stashDB     *StashDBClient
}

type VideoFingerprint struct {
    PHash      string
    OsHash     string
    MD5        string
    Duration   int // seconds
    Resolution string
    Codec      string
}

// Generate fingerprints for a video file
func (s *FingerprintService) Fingerprint(ctx context.Context, path string) (*VideoFingerprint, error) {
    // Run ffprobe for technical metadata
    probe, err := s.probeFile(ctx, path)
    if err != nil {
        return nil, err
    }

    fp := &VideoFingerprint{
        Duration:   probe.Duration,
        Resolution: fmt.Sprintf("%dx%d", probe.Width, probe.Height),
        Codec:      probe.VideoCodec,
    }

    // Calculate oshash (OpenSubtitles algorithm)
    fp.OsHash, err = s.calculateOsHash(path)
    if err != nil {
        return nil, err
    }

    // Calculate pHash (requires frame extraction)
    fp.PHash, err = s.calculatePHash(ctx, path)
    if err != nil {
        // pHash may fail, continue with oshash
        slog.Warn("pHash calculation failed", "path", path, "error", err)
    }

    return fp, nil
}

// OpenSubtitles hash algorithm
func (s *FingerprintService) calculateOsHash(path string) (string, error) {
    file, err := os.Open(path)
    if err != nil {
        return "", err
    }
    defer file.Close()

    fi, err := file.Stat()
    if err != nil {
        return "", err
    }

    if fi.Size() < 65536*2 {
        return "", errors.New("file too small for oshash")
    }

    // oshash = first 64KB + last 64KB + filesize
    buf := make([]byte, 65536)
    var hash uint64 = uint64(fi.Size())

    // First 64KB
    _, err = file.Read(buf)
    if err != nil {
        return "", err
    }
    for i := 0; i < 65536; i += 8 {
        hash += binary.LittleEndian.Uint64(buf[i:])
    }

    // Last 64KB
    _, err = file.Seek(-65536, io.SeekEnd)
    if err != nil {
        return "", err
    }
    _, err = file.Read(buf)
    if err != nil {
        return "", err
    }
    for i := 0; i < 65536; i += 8 {
        hash += binary.LittleEndian.Uint64(buf[i:])
    }

    return fmt.Sprintf("%016x", hash), nil
}

// Match file to StashDB scene
func (s *FingerprintService) MatchScene(ctx context.Context, fp *VideoFingerprint) (*StashDBScene, error) {
    // Try pHash first (most reliable)
    if fp.PHash != "" {
        scene, err := s.stashDB.FindSceneByFingerprint(ctx, Fingerprint{
            Hash:      fp.PHash,
            Algorithm: "phash",
            Duration:  fp.Duration,
        })
        if err == nil && scene != nil {
            return scene, nil
        }
    }

    // Fallback to oshash
    return s.stashDB.FindSceneByFingerprint(ctx, Fingerprint{
        Hash:      fp.OsHash,
        Algorithm: "oshash",
        Duration:  fp.Duration,
    })
}
```

---

## Database Schema

All tables in isolated `qar` schema:

```sql
CREATE SCHEMA IF NOT EXISTS qar;

-- Studios → Ports
CREATE TABLE qar.ports (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name            VARCHAR(255) NOT NULL,
    parent_id       UUID REFERENCES qar.ports(id),
    stashdb_id      VARCHAR(100),          -- StashDB studio ID
    tpdb_id         VARCHAR(100),          -- TPDB studio ID
    url             TEXT,
    logo_path       TEXT,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(stashdb_id),
    UNIQUE(tpdb_id)
);

-- Performers → Crew (shared between expeditions and voyages)
CREATE TABLE qar.crew (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name            VARCHAR(255) NOT NULL,
    disambiguation  VARCHAR(255),          -- For same-name performers
    stash_id        VARCHAR(100),
    charter         VARCHAR(100),          -- stashdb_id (obfuscated)
    registry        VARCHAR(100),          -- tpdb_id (obfuscated)
    manifest        VARCHAR(100),          -- freeones_id (obfuscated)

    -- Demographics
    gender          VARCHAR(50),
    christening     DATE,                  -- birthdate (obfuscated)
    death_date      DATE,
    birth_city      VARCHAR(255),
    origin          VARCHAR(100),          -- ethnicity (obfuscated)
    nationality     VARCHAR(100),
    height_cm       INT,
    weight_kg       INT,

    -- Measurements → Cargo
    cargo           JSONB,                 -- measurements (obfuscated)
    cup_size        VARCHAR(10),
    breast_type     VARCHAR(50),
    cutlass         VARCHAR(20),           -- penis_size (obfuscated) - male/trans performers

    -- Trans-specific (obfuscated)
    figurehead      BOOLEAN,               -- has_breasts: trans breast presence
    keel            VARCHAR(20),           -- genitalia_type: male/female/both
    refit           VARCHAR(20),           -- surgical_status: pre-op/post-op/non-op

    -- Appearance
    rigging         VARCHAR(50),           -- hair_color (obfuscated)
    compass         VARCHAR(50),           -- eye_color (obfuscated)
    markings        TEXT[],                -- tattoos (obfuscated)
    anchors         TEXT[],                -- piercings (obfuscated)

    -- Career
    maiden_voyage   INT,                   -- career_start year (obfuscated)
    last_port       INT,                   -- career_end year (obfuscated) (NULL = active)
    bio             TEXT,

    -- Social
    twitter         TEXT,
    instagram       TEXT,

    -- Images
    image_path      TEXT,                  -- Primary image

    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(charter),
    UNIQUE(registry)
);

-- Crew aliases (performer aliases)
CREATE TABLE qar.crew_aliases (
    crew_id         UUID REFERENCES qar.crew(id) ON DELETE CASCADE,
    alias           VARCHAR(255) NOT NULL,
    PRIMARY KEY (crew_id, alias)
);

-- Crew images (performer images - additional)
CREATE TABLE qar.crew_images (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    crew_id         UUID REFERENCES qar.crew(id) ON DELETE CASCADE,
    path            TEXT NOT NULL,
    type            VARCHAR(50) DEFAULT 'photo', -- photo, headshot, full
    source          VARCHAR(50),                  -- stashdb, tpdb, local
    primary_image   BOOLEAN DEFAULT FALSE,
    created_at      TIMESTAMPTZ DEFAULT NOW()
);

-- Movies → Expeditions
CREATE TABLE qar.expeditions (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    fleet_id        UUID NOT NULL,              -- Reference to qar.fleets (library)
    whisparr_id     INT,
    charter         VARCHAR(100),               -- stashdb_id (obfuscated)
    registry        VARCHAR(100),               -- tpdb_id (obfuscated)

    -- Metadata
    title           VARCHAR(500) NOT NULL,
    sort_title      VARCHAR(500),
    original_title  VARCHAR(500),
    overview        TEXT,
    launch_date     DATE,                       -- release_date (obfuscated)
    runtime_ticks   BIGINT,                     -- 10,000,000 ticks = 1 second
    port_id         UUID REFERENCES qar.ports(id), -- studio_id (obfuscated)
    director        VARCHAR(255),
    series          VARCHAR(255),

    -- File info
    path            TEXT NOT NULL,
    size_bytes      BIGINT,
    container       VARCHAR(50),
    video_codec     VARCHAR(50),
    audio_codec     VARCHAR(50),
    resolution      VARCHAR(50),

    -- Fingerprints → Coordinates
    coordinates     VARCHAR(64),                -- phash (obfuscated)
    oshash          VARCHAR(64),

    -- Status
    has_file        BOOLEAN DEFAULT TRUE,
    is_hdr          BOOLEAN DEFAULT FALSE,
    is_3d           BOOLEAN DEFAULT FALSE,

    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(charter),
    UNIQUE(path)
);

-- Expedition crew (movie performers)
CREATE TABLE qar.expedition_crew (
    expedition_id   UUID REFERENCES qar.expeditions(id) ON DELETE CASCADE,
    crew_id         UUID REFERENCES qar.crew(id) ON DELETE CASCADE,
    character_name  VARCHAR(255),               -- Optional role name
    PRIMARY KEY (expedition_id, crew_id)
);

-- Tags → Flags
CREATE TABLE qar.flags (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name            VARCHAR(255) NOT NULL UNIQUE,
    charter         VARCHAR(36),                -- stashdb_id (obfuscated)
    parent_id       UUID REFERENCES qar.flags(id),
    description     TEXT,
    created_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE qar.expedition_flags (
    expedition_id   UUID REFERENCES qar.expeditions(id) ON DELETE CASCADE,
    flag_id         UUID REFERENCES qar.flags(id) ON DELETE CASCADE,
    PRIMARY KEY (expedition_id, flag_id)
);

-- Scenes → Voyages (primary adult scene content)
CREATE TABLE qar.voyages (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    fleet_id        UUID NOT NULL REFERENCES qar.fleets(id), -- library
    title           VARCHAR(500) NOT NULL,
    sort_title      VARCHAR(500),
    overview        TEXT,
    launch_date     DATE,                       -- release_date (obfuscated)
    distance        INT,                        -- runtime_minutes (obfuscated)
    port_id         UUID REFERENCES qar.ports(id), -- studio

    whisparr_id     INT,
    stash_id        VARCHAR(100),
    charter         VARCHAR(100),               -- stashdb_id (obfuscated)
    registry        VARCHAR(100),               -- tpdb_id (obfuscated)

    path            TEXT NOT NULL,
    size_bytes      BIGINT,
    video_codec     VARCHAR(50),
    audio_codec     VARCHAR(50),
    resolution      VARCHAR(20),

    oshash          VARCHAR(32),
    coordinates     VARCHAR(32),                -- phash (obfuscated)
    md5             VARCHAR(64),

    cover_path      TEXT,

    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(path)
);

CREATE TABLE qar.voyage_crew (
    voyage_id       UUID REFERENCES qar.voyages(id) ON DELETE CASCADE,
    crew_id         UUID REFERENCES qar.crew(id) ON DELETE CASCADE,
    PRIMARY KEY (voyage_id, crew_id)
);

CREATE TABLE qar.voyage_flags (
    voyage_id       UUID REFERENCES qar.voyages(id) ON DELETE CASCADE,
    flag_id         UUID REFERENCES qar.flags(id) ON DELETE CASCADE,
    PRIMARY KEY (voyage_id, flag_id)
);

-- Voyage markers (scene markers - timestamped tags from Stash)
CREATE TABLE qar.voyage_markers (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    voyage_id       UUID REFERENCES qar.voyages(id) ON DELETE CASCADE,
    title           VARCHAR(255),
    start_seconds   FLOAT NOT NULL,
    end_seconds     FLOAT,
    flag_id         UUID REFERENCES qar.flags(id),
    stash_marker_id VARCHAR(100),
    created_at      TIMESTAMPTZ DEFAULT NOW()
);

-- Expedition images (movie images)
CREATE TABLE qar.expedition_images (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    expedition_id   UUID REFERENCES qar.expeditions(id) ON DELETE CASCADE,
    type            VARCHAR(50) NOT NULL,       -- poster, backdrop, screenshot
    path            TEXT NOT NULL,
    source          VARCHAR(50),                -- stashdb, whisparr, local
    primary_image   BOOLEAN DEFAULT FALSE,
    created_at      TIMESTAMPTZ DEFAULT NOW()
);

-- Galleries → Treasures (image sets)
CREATE TABLE qar.treasures (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    expedition_id   UUID REFERENCES qar.expeditions(id) ON DELETE SET NULL,
    title           VARCHAR(500) NOT NULL,
    path            TEXT,
    image_count     INT,
    created_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE qar.treasure_images (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    treasure_id     UUID REFERENCES qar.treasures(id) ON DELETE CASCADE,
    path            TEXT NOT NULL,
    position        INT,
    width           INT,
    height          INT
);

-- User data (per-module)
CREATE TABLE qar.user_bounties (
    user_id         UUID NOT NULL,              -- Reference to public.users (no FK)
    expedition_id   UUID REFERENCES qar.expeditions(id) ON DELETE CASCADE,
    bounty          INT CHECK (bounty >= 0 AND bounty <= 100), -- rating (obfuscated)
    rated_at        TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (user_id, expedition_id)
);

CREATE TABLE qar.user_favorites (
    user_id         UUID NOT NULL,
    expedition_id   UUID REFERENCES qar.expeditions(id) ON DELETE CASCADE,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (user_id, expedition_id)
);

CREATE TABLE qar.ship_log (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL,
    expedition_id   UUID REFERENCES qar.expeditions(id) ON DELETE CASCADE,
    logged_at       TIMESTAMPTZ DEFAULT NOW(),  -- watched_at (obfuscated)
    position_ticks  BIGINT,                     -- Last position
    completed       BOOLEAN DEFAULT FALSE
);

-- Indexes
CREATE INDEX idx_qar_expeditions_fleet ON qar.expeditions(fleet_id);
CREATE INDEX idx_qar_expeditions_port ON qar.expeditions(port_id);
CREATE INDEX idx_qar_expeditions_coordinates ON qar.expeditions(coordinates);
CREATE INDEX idx_qar_expeditions_oshash ON qar.expeditions(oshash);
CREATE INDEX idx_qar_crew_name ON qar.crew(name);
CREATE INDEX idx_qar_crew_charter ON qar.crew(charter);
CREATE INDEX idx_qar_ship_log_user ON qar.ship_log(user_id, logged_at DESC);
CREATE INDEX idx_qar_voyages_fleet ON qar.voyages(fleet_id);
CREATE INDEX idx_qar_voyages_port ON qar.voyages(port_id);
CREATE INDEX idx_qar_voyages_coordinates ON qar.voyages(coordinates);
```

---

## Voyage-First Model (Whisparr v3)

Whisparr models scenes as episodes under a series/site. Revenge keeps a voyage-first model:

- No `qar.series`, `qar.seasons`, or `qar.episodes` tables.
- Voyages (scenes) are stored directly in `qar.voyages` and linked to `qar.crew`, `qar.ports`, and `qar.flags`.
- Series/site grouping is represented via port/site metadata and UI grouping, not separate tables.

---

## Features Specific to Adult Content

### Voyage Markers (Timestamped Tags)

Import Stash scene markers as navigable chapters:

```go
type VoyageMarkerService struct {
    db    *pgxpool.Pool
    stash *StashAppClient
}

type VoyageMarker struct {
    ID          uuid.UUID
    VoyageID    uuid.UUID
    Title       string
    StartTicks  int64
    EndTicks    *int64
    PrimaryFlag *Flag
    Flags       []Flag
    Thumbnail   string
}

// Import markers from Stash
func (s *VoyageMarkerService) ImportFromStash(ctx context.Context, voyageID uuid.UUID, stashSceneID string) error {
    markers, err := s.stash.GetSceneMarkers(ctx, stashSceneID)
    if err != nil {
        return err
    }

    for _, marker := range markers {
        _, err := s.db.Exec(ctx, `
            INSERT INTO qar.voyage_markers (voyage_id, title, start_ticks, primary_flag_id, thumbnail_path)
            VALUES ($1, $2, $3, $4, $5)
        `, voyageID, marker.Title, int64(marker.Seconds*10_000_000),
           s.ensureFlag(ctx, marker.PrimaryTag), marker.Screenshot)
        if err != nil {
            return err
        }
    }
    return nil
}
```

### Performer Filtering

Advanced performer search with attribute filtering:

```go
type PerformerFilter struct {
    Name        *string
    Gender      *string
    Ethnicity   *string
    HairColor   *string
    EyeColor    *string
    MinHeight   *int
    MaxHeight   *int
    HasTattoos  *bool
    HasPiercings *bool
    ActiveOnly  *bool
    Tags        []uuid.UUID
    Studios     []uuid.UUID
}

func (r *PerformerRepository) Search(ctx context.Context, filter PerformerFilter, page, pageSize int) ([]Performer, int, error) {
    // Build dynamic query with filters
}
```

### Studio Networks

Track studio ownership and networks:

```go
type StudioService struct {
    db *pgxpool.Pool
}

// Get studio with parent chain
func (s *StudioService) GetWithNetwork(ctx context.Context, studioID uuid.UUID) (*StudioNetwork, error) {
    // Returns studio with parent hierarchy
    // e.g., Reality Kings -> Team Skeet -> step-sibling brand
}
```

### Interactive Content

Support for interactive content (funscripts):

```sql
-- Interactive content support
ALTER TABLE qar.expeditions ADD COLUMN interactive BOOLEAN DEFAULT FALSE;
ALTER TABLE qar.expeditions ADD COLUMN funscript_path TEXT;

CREATE TABLE qar.interactive_scripts (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    expedition_id   UUID REFERENCES qar.expeditions(id) ON DELETE CASCADE,
    script_path     TEXT NOT NULL,
    script_type     VARCHAR(50) DEFAULT 'funscript',
    created_at      TIMESTAMPTZ DEFAULT NOW()
);
```

---

## Security Hardening

### 1. Complete Schema Isolation

```sql
-- Separate schema with no cross-references to public schema
CREATE SCHEMA qar;

-- QAR tables NEVER reference public.* tables directly
-- User IDs are stored but no foreign keys to shared.users
-- This prevents accidental data leakage through JOINs

-- Row-level security (RLS) for additional protection
ALTER TABLE qar.crew ENABLE ROW LEVEL SECURITY;
ALTER TABLE qar.voyages ENABLE ROW LEVEL SECURITY;
ALTER TABLE qar.expeditions ENABLE ROW LEVEL SECURITY;

-- Only users with legacy_access role can query
CREATE POLICY qar_crew_policy ON qar.crew
    FOR ALL TO authenticated
    USING (current_setting('app.legacy_access', true)::boolean = true);
```

### 2. Access Control (RBAC)

```go
// Obfuscated scope names - no "adult" in any token/claim
const (
    ScopeLegacyRead  = "legacy:read"   // Read access
    ScopeLegacyWrite = "legacy:write"  // Write access
    ScopeLegacyAdmin = "legacy:admin"  // Admin access
)

// Middleware validates scope AND explicit user opt-in
func LegacyAuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        claims := auth.ClaimsFromContext(r.Context())

        // 1. Check scope in token
        if !claims.HasScope(ScopeLegacyRead) {
            http.Error(w, "Not Found", http.StatusNotFound) // 404, not 403
            return
        }

        // 2. Check user has explicitly enabled legacy content
        if !claims.LegacyEnabled {
            http.Error(w, "Not Found", http.StatusNotFound)
            return
        }

        // 3. Check PIN if required
        if config.Legacy.RequirePIN && !validatePIN(r, claims.UserID) {
            http.Error(w, "PIN Required", http.StatusUnauthorized)
            return
        }

        // Set database session variable for RLS
        ctx := context.WithValue(r.Context(), "legacy_access", true)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

### 3. URL Obfuscation

```go
// Router configuration - external URLs never reveal content type
router.Route("/api/v1/legacy", func(r chi.Router) {
    r.Use(LegacyAuthMiddleware)

    // All routes look like generic "legacy" API
    r.Get("/crew/{id}", handlers.GetCrew)
    r.Get("/voyages/{id}", handlers.GetVoyage)
    r.Get("/expeditions/{id}", handlers.GetExpedition)
})

// Response headers - no hints about content type
func SetSecurityHeaders(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Remove any identifying headers
        w.Header().Del("X-Content-Type")
        w.Header().Set("X-Content-Type-Options", "nosniff")
        w.Header().Set("Cache-Control", "private, no-store")
        next.ServeHTTP(w, r)
    })
}
```

### 4. Audit Logging (Isolated)

```go
// Audit logs stored in QAR schema, not main activity_log
func (l *AuditLogger) LogLegacyAccess(ctx context.Context, userID uuid.UUID, action, resource string) {
    _, err := l.db.Exec(ctx, `
        INSERT INTO qar.ship_log (sailor_id, action, cargo, bearing, vessel, logged_at)
        VALUES ($1, $2, $3, $4, $5, NOW())
    `, userID, action, resource, getIP(ctx), getUserAgent(ctx))
    if err != nil {
        slog.Error("ship log failed", "error", err)
    }
}

// Audit log table uses QAR terminology
// sailor_id = user_id, cargo = resource, bearing = ip_address, vessel = user_agent
```

### 5. Network Security

```yaml
# Reverse proxy configuration (nginx/caddy)
# Legacy endpoints should:
# 1. Not appear in access logs (or use separate log file)
# 2. Use separate rate limiting
# 3. Have no referrer leakage

location /api/v1/legacy/ {
    # Disable access logging for privacy
    access_log off;

    # Separate rate limit bucket
    limit_req zone=legacy burst=10 nodelay;

    # No referrer leakage
    add_header Referrer-Policy "no-referrer" always;

    # Prevent caching at proxy level
    proxy_no_cache 1;
    proxy_cache_bypass 1;

    proxy_pass http://backend;
}
```

### 6. Data Encryption

```go
// Sensitive fields encrypted at rest
type EncryptedField struct {
    Ciphertext []byte
    Nonce      []byte
}

// Crew names and aliases can be encrypted
func (r *Repository) GetCrew(ctx context.Context, id uuid.UUID) (*Crew, error) {
    var crew Crew
    err := r.db.QueryRow(ctx, `
        SELECT id, pgp_sym_decrypt(names_encrypted, $2) as names, ...
        FROM qar.crew WHERE id = $1
    `, id, r.encryptionKey).Scan(&crew)
    return &crew, err
}
```

### 7. Privacy Controls

```yaml
# config.yaml - obfuscated section name
legacy:
  enabled: true

  privacy:
    # Exclude from all cross-module features
    exclude_from_analytics: true
    exclude_from_recommendations: true
    exclude_from_continue_watching: true
    exclude_from_search_history: true

    # Access controls
    require_pin: true
    pin_timeout: 30m
    auto_lock_on_idle: 5m

    # Audit settings
    audit_all_access: true
    audit_retention_days: 90

    # Network privacy
    disable_external_requests: false  # Set true to block StashDB calls
    proxy_all_images: true            # Proxy images through server
```

### 8. Error Response Obfuscation

```go
// Never reveal that legacy content exists in error messages
func LegacyErrorHandler(err error, w http.ResponseWriter, r *http.Request) {
    switch {
    case errors.Is(err, ErrCrewNotFound):
        http.Error(w, "Not Found", http.StatusNotFound)
    case errors.Is(err, ErrUnauthorized):
        http.Error(w, "Not Found", http.StatusNotFound) // 404, not 401
    case errors.Is(err, ErrForbidden):
        http.Error(w, "Not Found", http.StatusNotFound) // 404, not 403
    default:
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
    }
}
```

### 9. Database Connection Isolation

```go
// Separate connection pool for QAR queries
func NewQARPool(cfg config.Database) (*pgxpool.Pool, error) {
    config, _ := pgxpool.ParseConfig(cfg.URL)

    // Set search_path to only qar schema
    config.ConnConfig.RuntimeParams["search_path"] = "qar"

    // Use separate pool with lower limits
    config.MaxConns = 5
    config.MinConns = 1

    return pgxpool.NewWithConfig(context.Background(), config)
}
```

### Security Checklist

- [ ] Schema isolation verified (no FK to public.*)
- [ ] RLS policies enabled on all QAR tables
- [ ] URL obfuscation tested (no "adult" in any URL/response)
- [ ] Error messages don't reveal content type
- [ ] Audit logging to separate table
- [ ] Access logging disabled for /legacy/ endpoints
- [ ] PIN protection implemented
- [ ] Session timeout configured
- [ ] Image proxying enabled
- [ ] Encryption at rest for sensitive fields
- [ ] RBAC scopes use obfuscated names

---

## API Endpoints

All endpoints under innocuous `/legacy/` namespace (looks like deprecated API):

```
# Expeditions (Movies)
POST   /api/v1/legacy/expeditions                 # Create expedition
GET    /api/v1/legacy/expeditions                 # List expeditions
GET    /api/v1/legacy/expeditions/{id}            # Get expedition
PUT    /api/v1/legacy/expeditions/{id}            # Update expedition
DELETE /api/v1/legacy/expeditions/{id}            # Delete expedition
GET    /api/v1/legacy/expeditions/{id}/crew       # Get expedition crew
GET    /api/v1/legacy/expeditions/{id}/markers    # Get voyage markers
GET    /api/v1/legacy/expeditions/{id}/similar    # Get similar expeditions

# Voyages (Scenes)
GET    /api/v1/legacy/voyages                     # List voyages
GET    /api/v1/legacy/voyages/{id}                # Get voyage
GET    /api/v1/legacy/voyages/{id}/crew           # Get voyage crew
GET    /api/v1/legacy/voyages/{id}/markers        # Get voyage markers

# Crew (Performers)
GET    /api/v1/legacy/crew                        # List crew
GET    /api/v1/legacy/crew/{id}                   # Get crew member
GET    /api/v1/legacy/crew/{id}/voyages           # Get crew voyages
GET    /api/v1/legacy/crew/{id}/expeditions       # Get crew expeditions
GET    /api/v1/legacy/crew/{id}/treasures         # Get crew galleries

# Ports (Studios)
GET    /api/v1/legacy/ports                       # List ports
GET    /api/v1/legacy/ports/{id}                  # Get port
GET    /api/v1/legacy/ports/{id}/voyages          # Get port voyages
GET    /api/v1/legacy/ports/{id}/expeditions      # Get port expeditions

# Flags (Tags)
GET    /api/v1/legacy/flags                       # List flags
GET    /api/v1/legacy/flags/{id}                  # Get flag
GET    /api/v1/legacy/flags/{id}/voyages          # Get flagged voyages

# Treasures (Galleries)
GET    /api/v1/legacy/treasures                   # List treasures
GET    /api/v1/legacy/treasures/{id}              # Get treasure
GET    /api/v1/legacy/treasures/{id}/doubloons    # Get treasure images

# Fleets (Libraries)
GET    /api/v1/legacy/fleets                      # List fleets
GET    /api/v1/legacy/fleets/{id}                 # Get fleet
POST   /api/v1/legacy/fleets/{id}/scan            # Trigger fleet scan

# Identification
POST   /api/v1/legacy/match                       # Match file by fingerprint
POST   /api/v1/legacy/identify                    # Submit for identification

# User Data (per-user stats)
GET    /api/v1/legacy/logbook                     # User watch history
GET    /api/v1/legacy/prized                      # User favorites
POST   /api/v1/legacy/voyages/{id}/bounty         # Rate voyage
POST   /api/v1/legacy/expeditions/{id}/bounty     # Rate expedition
```

---

## Configuration

```yaml
# config.yaml
# Note: Section named "legacy" for obfuscation - no "adult" in config keys
legacy:
  enabled: true

  # Whisparr integration (Arr for acquisition)
  whisparr:
    url: "http://whisparr:6969"
    api_key: "${WHISPARR_API_KEY}"

  # StashDB integration (metadata enrichment)
  stashdb:
    endpoint: "https://stashdb.org/graphql"
    api_key: "${STASHDB_API_KEY}"

  # Stash-App integration (optional, for local instance)
  stash_app:
    url: "http://stash:9999"
    api_key: "${STASH_API_KEY}"

  # ThePornDB fallback
  tpdb:
    url: "https://api.metadataapi.net/api"
    api_key: "${TPDB_API_KEY}"

  # Fingerprinting
  fingerprint:
    generate_phash: true
    generate_oshash: true
    auto_match: true

  # Privacy & Security
  privacy:
    exclude_from_analytics: true
    exclude_from_recommendations: true
    exclude_from_continue_watching: true
    exclude_from_search_history: true
    require_pin: true
    pin_timeout: 30m
    auto_lock_on_idle: 5m
    audit_all_access: true
    audit_retention_days: 90
    proxy_all_images: true

  # Storage (paths use "qar" internally, never exposed)
  storage:
    base_path: "/data/qar"
    images_path: "/data/qar/images"
    thumbnails_path: "/data/qar/thumbs"
    cache_path: "/data/qar/cache"

  # Database
  database:
    schema: "qar"
    separate_pool: true
    max_connections: 5
    encryption_key: "${QAR_ENCRYPTION_KEY}"
```

---

## River Jobs

```go
// Fingerprint and match scenes
type FingerprintSceneArgs struct {
    MovieID uuid.UUID `json:"movie_id"`
}

func (FingerprintSceneArgs) Kind() string { return "adult.fingerprint_scene" }

type FingerprintSceneWorker struct {
    river.WorkerDefaults[FingerprintSceneArgs]
    fingerprinter *FingerprintService
    db           *pgxpool.Pool
}

func (w *FingerprintSceneWorker) Work(ctx context.Context, job *river.Job[FingerprintSceneArgs]) error {
    // Get expedition path
    var path string
    err := w.db.QueryRow(ctx, `SELECT path FROM qar.expeditions WHERE id = $1`, job.Args.MovieID).Scan(&path)
    if err != nil {
        return err
    }

    // Generate fingerprints
    fp, err := w.fingerprinter.Fingerprint(ctx, path)
    if err != nil {
        return err
    }

    // Store fingerprints
    _, err = w.db.Exec(ctx, `
        UPDATE qar.expeditions SET coordinates = $1, oshash = $2 WHERE id = $3
    `, fp.PHash, fp.OsHash, job.Args.MovieID)
    if err != nil {
        return err
    }

    // Try to match with StashDB
    scene, err := w.fingerprinter.MatchScene(ctx, fp)
    if err == nil && scene != nil {
        // Queue metadata enrichment
        return w.river.Insert(ctx, EnrichFromStashDBArgs{
            MovieID:  job.Args.MovieID,
            SceneID:  scene.ID,
        })
    }

    return nil
}

// Enrich from StashDB
type EnrichFromStashDBArgs struct {
    MovieID uuid.UUID `json:"movie_id"`
    SceneID string    `json:"scene_id"`
}

func (EnrichFromStashDBArgs) Kind() string { return "adult.enrich_stashdb" }

// Sync with Stash-App
type SyncStashAppArgs struct {
    LibraryID uuid.UUID `json:"library_id"`
}

func (SyncStashAppArgs) Kind() string { return "adult.sync_stash_app" }
```

---

## Implementation Checklist

### Phase 1: Schema & Database Setup
- [ ] Create `qar` schema in PostgreSQL
- [ ] Create migration files for all tables:
  - [ ] `qar.fleets` (libraries)
  - [ ] `qar.crew` (performers) with aliases and images
  - [ ] `qar.crew_aliases` and `qar.crew_images`
  - [ ] `qar.ports` (studios)
  - [ ] `qar.flags` (tags)
  - [ ] `qar.expeditions` (movies)
  - [ ] `qar.expedition_crew` and `qar.expedition_flags`
  - [ ] `qar.expedition_images`
  - [ ] `qar.voyages` (scenes)
  - [ ] `qar.voyage_crew` and `qar.voyage_flags`
  - [ ] `qar.voyage_markers` (timestamped chapters)
  - [ ] `qar.treasures` (galleries) and `qar.doubloons` (images)
  - [ ] `qar.user_bounties`, `qar.user_favorites`, `qar.ship_log` (user data)
  - [ ] `qar.plunder_queue` (download queue)
- [ ] Create all indexes
- [ ] Implement Row-Level Security (RLS) policies
- [ ] Create `sqlc` queries file `queries/qar/adult.sql`

### Phase 2: Core Repository & Service Layer
- [ ] Implement `internal/content/qar/repository.go` interface
- [ ] Implement `internal/content/qar/repository_pg.go` using sqlc
- [ ] Implement `internal/content/qar/service.go` with business logic
- [ ] Implement entity structs in `internal/content/qar/entity.go`
- [ ] Create mapping layer for QAR terminology

### Phase 3: Whisparr Integration (Acquisition)
- [ ] Create `internal/integrations/whisparr/client.go`
- [ ] Implement movie sync from Whisparr v3 API
- [ ] Implement webhook handler for real-time updates (Download, Rename, Delete)
- [ ] Queue fingerprinting jobs after import
- [ ] Handle metadata mapping (studio, performers, genres)

### Phase 4: Fingerprinting & Identification
- [ ] Implement `internal/content/qar/fingerprint/service.go`
- [ ] Implement OSHASH generation (file-based)
- [ ] Implement pHash generation (perceptual hash)
- [ ] Implement MD5 hash generation
- [ ] Create River jobs for fingerprinting pipeline
- [ ] Integrate StashDB matching

### Phase 5: StashDB Enrichment
- [ ] Create `internal/integrations/stashdb/client.go` (GraphQL)
- [ ] Implement performer lookup by fingerprint
- [ ] Implement scene data enrichment
- [ ] Implement performer detail fetching
- [ ] Create River jobs for enrichment pipeline
- [ ] Cache StashDB responses

### Phase 6: Stash-App Integration (Optional/Private)
- [ ] Create `internal/integrations/stash/client.go` (GraphQL)
- [ ] Implement scene/marker sync from local Stash instance
- [ ] Implement performer data sync
- [ ] Implement scene marker → chapter conversion

### Phase 7: TPDB Fallback
- [ ] Create `internal/integrations/tpdb/client.go`
- [ ] Implement scene search by title/performers
- [ ] Implement performer search
- [ ] Use as fallback when StashDB has no results

### Phase 8: Access Control & Security
- [ ] Create RBAC scopes: `legacy:read`, `legacy:write`, `legacy:admin`
- [ ] Implement LegacyAuthMiddleware with scope + opt-in checks
- [ ] Implement PIN protection if configured
- [ ] Implement session timeout and auto-lock
- [ ] Create separate connection pool for QAR queries
- [ ] Implement audit logging to `qar.ship_log`

### Phase 9: API Endpoints
- [ ] Create OpenAPI spec `api/openapi/qar_adult.yaml`
- [ ] Implement expedition endpoints (`/api/v1/legacy/expeditions/*`)
- [ ] Implement voyage endpoints (`/api/v1/legacy/voyages/*`)
- [ ] Implement crew endpoints (`/api/v1/legacy/crew/*`)
- [ ] Implement port endpoints (`/api/v1/legacy/ports/*`)
- [ ] Implement flag endpoints (`/api/v1/legacy/flags/*`)
- [ ] Implement fleet endpoints (`/api/v1/legacy/fleets/*`)
- [ ] Implement user data endpoints (logbook, favorites, ratings)
- [ ] Implement fingerprint matching endpoint `/api/v1/legacy/match`

### Phase 10: River Jobs & Background Processing
- [ ] Implement job handlers for fingerprinting
- [ ] Implement job handlers for StashDB enrichment
- [ ] Implement job handlers for Stash-App sync
- [ ] Implement job handlers for TPDB fallback
- [ ] Create job scheduling for periodic reconciliation

### Phase 11: Testing & Validation
- [ ] Write unit tests for repository layer
- [ ] Write unit tests for service layer
- [ ] Write integration tests for Whisparr client
- [ ] Write integration tests for StashDB client
- [ ] Write API tests for all endpoints
- [ ] Test RLS policies with different user roles
- [ ] Test audit logging

### Phase 12: Deployment & Documentation
- [ ] Create environment variable documentation
- [ ] Document Whisparr webhook configuration
- [ ] Document StashDB API key setup
- [ ] Create migration guide for fresh installs
- [ ] Document obfuscation terminology mapping


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

**Category**: [Adult](INDEX.md)

### In This Section

- [Revenge - Adult Content Metadata System](ADULT_METADATA.md)
- [Adult Data Reconciliation](DATA_RECONCILIATION.md)
- [Adult Gallery Module (QAR: Treasures)](GALLERY_MODULE.md)
- [Whisparr v3 & StashDB Schema Integration](WHISPARR_STASHDB_SCHEMA.md)

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

---

## Summary

| Aspect | Implementation |
|--------|----------------|
| Schema | Isolated `qar` PostgreSQL schema |
| API Namespace | Obscured `/qar/` path |
| Primary Source | Whisparr-v3 for acquisition |
| Enrichment | StashDB (public), Stash-App (private) |
| Identification | pHash (coordinates) + oshash fingerprinting |
| Privacy | Audit logging, analytics exclusion, PIN protection |
| Voyage Markers | Timestamped flags from Stash |
| Crew Data | StashDB with aliases, demographics, cargo |
