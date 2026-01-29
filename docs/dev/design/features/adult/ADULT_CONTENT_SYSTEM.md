# Revenge - Adult Content System

> Complete adult content management with Stash ecosystem integration.
> All adult content isolated in PostgreSQL schema `qar` with "Queen Anne's Revenge" themed obfuscation.

## Design Principles

1. **Complete Isolation** - Separate schema, namespace, and storage
2. **Whisparr as Proxy** - Whisparr-v3 for acquisition, local metadata enrichment
3. **Stash as Enrichment** - StashDB, Stash-App for comprehensive performer/scene data
4. **Perceptual Hashing** - pHash fingerprinting for scene identification
5. **Privacy First** - Fully obfuscated terminology, optional analytics exclusion, audit logging

---

## Queen Anne's Revenge Obfuscation

All adult content uses pirate ship terminology themed after Blackbeard's famous vessel.

### Schema & API Namespace

| Real Concept | Obfuscated | Reason |
|--------------|------------|--------|
| Adult schema | `qar` | Queen Anne's Revenge |
| API namespace | `/api/v1/qar/` | Matches schema |
| Storage path | `/media/qar/` | Consistent theming |

### Entity Obfuscation

| Real Entity | Obfuscated | Database Table | API Endpoint |
|-------------|------------|----------------|--------------|
| Performers | **Crew** | `qar.crew` | `/qar/crew` |
| Scenes | **Voyages** | `qar.voyages` | `/qar/voyages` |
| Movies | **Expeditions** | `qar.expeditions` | `/qar/expeditions` |
| Studios | **Ports** | `qar.ports` | `/qar/ports` |
| Tags | **Flags** | `qar.flags` | `/qar/flags` |
| Categories | **Waters** | `qar.waters` | `/qar/waters` |
| Galleries/Images | **Treasures** | `qar.treasures` | `/qar/treasures` |
| Libraries | **Fleets** | `qar.fleets` | `/qar/fleets` |

### Field Obfuscation

| Real Field | Obfuscated | Used In |
|------------|------------|---------|
| measurements | `cargo` | crew |
| aliases | `names` | crew |
| tattoos | `markings` | crew |
| piercings | `anchors` | crew |
| career_start | `maiden_voyage` | crew |
| career_end | `last_port` | crew |
| birth_date | `christening` | crew |
| ethnicity | `origin` | crew |
| hair_color | `rigging` | crew |
| eye_color | `compass` | crew |
| scene_count | `voyage_count` | crew |
| penis_size | `cutlass` | crew (male/trans) |
| has_breasts | `figurehead` | crew (trans) |
| genitalia_type | `keel` | crew (trans: male/female/both) |
| surgical_status | `refit` | crew (trans: pre-op/post-op/non-op) |
| gallery | `chest` | crew |
| stashdb_id | `charter` | all entities |
| tpdb_id | `registry` | all entities |
| freeones_id | `manifest` | crew |
| rating | `bounty` | voyages, expeditions |
| duration | `distance` | voyages, expeditions |
| release_date | `launch_date` | expeditions |
| fingerprint | `coordinates` | voyages |

### Example API Responses

```json
// GET /api/v1/qar/crew/123
{
  "id": "123",
  "names": ["Anne Bonny", "Anne Cormac"],
  "christening": "1697-03-08",
  "origin": "Irish",
  "maiden_voyage": "1719",
  "last_port": null,
  "cargo": {
    "height": 165,
    "measurements": "34-24-35"
  },
  "voyage_count": 142,
  "chest": ["/qar/treasures/123/1.jpg"]
}

// GET /api/v1/qar/voyages/456
{
  "id": "456",
  "title": "Caribbean Adventure",
  "port": {"id": "789", "name": "Port Royal"},
  "crew": [{"id": "123", "name": "Anne Bonny"}],
  "flags": ["adventure", "caribbean"],
  "distance": 1847,
  "bounty": 4.5,
  "coordinates": "phash:abc123def456"
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
│  Library │ ──→ │  Whisparr-v3      │ ──→ │  c.* Schema     │
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

All tables in isolated `c` schema:

```sql
CREATE SCHEMA IF NOT EXISTS c;

-- Studios
CREATE TABLE c.studios (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name            VARCHAR(255) NOT NULL,
    parent_id       UUID REFERENCES c.studios(id),
    stashdb_id      VARCHAR(100),          -- StashDB studio ID
    tpdb_id         VARCHAR(100),          -- TPDB studio ID
    url             TEXT,
    logo_path       TEXT,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(stashdb_id),
    UNIQUE(tpdb_id)
);

-- Performers (shared between movies and scenes)
CREATE TABLE c.performers (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name            VARCHAR(255) NOT NULL,
    disambiguation  VARCHAR(255),          -- For same-name performers
    stash_id        VARCHAR(100),
    stashdb_id      VARCHAR(100),
    tpdb_id         VARCHAR(100),
    freeones_id     VARCHAR(100),

    -- Demographics
    gender          VARCHAR(50),
    birthdate       DATE,
    death_date      DATE,
    birth_city      VARCHAR(255),
    ethnicity       VARCHAR(100),
    nationality     VARCHAR(100),
    height_cm       INT,
    weight_kg       INT,

    -- Measurements
    measurements    VARCHAR(50),           -- e.g., "34D-24-34"
    cup_size        VARCHAR(10),
    breast_type     VARCHAR(50),
    penis_size      VARCHAR(20),           -- male/trans performers (Stash-App source)

    -- Trans-specific (obfuscated: figurehead, keel, refit)
    has_breasts     BOOLEAN,               -- trans: breast presence
    genitalia_type  VARCHAR(20),           -- trans: male/female/both
    surgical_status VARCHAR(20),           -- trans: pre-op/post-op/non-op

    -- Appearance
    hair_color      VARCHAR(50),
    eye_color       VARCHAR(50),
    tattoos         TEXT,
    piercings       TEXT,

    -- Career
    career_start    INT,                   -- Year
    career_end      INT,                   -- Year (NULL = active)
    bio             TEXT,

    -- Social
    twitter         TEXT,
    instagram       TEXT,

    -- Images
    image_path      TEXT,                  -- Primary image

    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(stashdb_id),
    UNIQUE(tpdb_id)
);

-- Performer aliases
CREATE TABLE c.performer_aliases (
    performer_id    UUID REFERENCES c.performers(id) ON DELETE CASCADE,
    alias           VARCHAR(255) NOT NULL,
    PRIMARY KEY (performer_id, alias)
);

-- Performer images (additional)
CREATE TABLE c.performer_images (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    performer_id    UUID REFERENCES c.performers(id) ON DELETE CASCADE,
    path            TEXT NOT NULL,
    type            VARCHAR(50) DEFAULT 'photo', -- photo, headshot, full
    source          VARCHAR(50),                  -- stashdb, tpdb, local
    primary_image   BOOLEAN DEFAULT FALSE,
    created_at      TIMESTAMPTZ DEFAULT NOW()
);

-- Movies
CREATE TABLE c.movies (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    library_id      UUID NOT NULL,              -- Reference to public.libraries (no FK)
    whisparr_id     INT,
    stashdb_id      VARCHAR(100),
    tpdb_id         VARCHAR(100),

    -- Metadata
    title           VARCHAR(500) NOT NULL,
    sort_title      VARCHAR(500),
    original_title  VARCHAR(500),
    overview        TEXT,
    release_date    DATE,
    runtime_ticks   BIGINT,                     -- 10,000,000 ticks = 1 second
    studio_id       UUID REFERENCES c.studios(id),
    director        VARCHAR(255),
    series          VARCHAR(255),

    -- File info
    path            TEXT NOT NULL,
    size_bytes      BIGINT,
    container       VARCHAR(50),
    video_codec     VARCHAR(50),
    audio_codec     VARCHAR(50),
    resolution      VARCHAR(50),

    -- Fingerprints
    phash           VARCHAR(64),
    oshash          VARCHAR(64),

    -- Status
    has_file        BOOLEAN DEFAULT TRUE,
    is_hdr          BOOLEAN DEFAULT FALSE,
    is_3d           BOOLEAN DEFAULT FALSE,

    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(stashdb_id),
    UNIQUE(path)
);

-- Movie performers
CREATE TABLE c.movie_performers (
    movie_id        UUID REFERENCES c.movies(id) ON DELETE CASCADE,
    performer_id    UUID REFERENCES c.performers(id) ON DELETE CASCADE,
    character_name  VARCHAR(255),               -- Optional role name
    PRIMARY KEY (movie_id, performer_id)
);

-- Movie tags
CREATE TABLE c.tags (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name            VARCHAR(255) NOT NULL UNIQUE,
    stashdb_id      VARCHAR(36),
    parent_id       UUID REFERENCES c.tags(id),
    description     TEXT,
    created_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE c.movie_tags (
    movie_id        UUID REFERENCES c.movies(id) ON DELETE CASCADE,
    tag_id          UUID REFERENCES c.tags(id) ON DELETE CASCADE,
    PRIMARY KEY (movie_id, tag_id)
);

-- Scenes (primary adult scene content)
CREATE TABLE c.scenes (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    library_id      UUID NOT NULL,
    title           VARCHAR(500) NOT NULL,
    sort_title      VARCHAR(500),
    overview        TEXT,
    release_date    DATE,
    runtime_minutes INT,
    studio_id       UUID REFERENCES c.studios(id),

    whisparr_id     INT,
    stash_id        VARCHAR(100),
    stashdb_id      VARCHAR(100),
    tpdb_id         VARCHAR(100),

    path            TEXT NOT NULL,
    size_bytes      BIGINT,
    video_codec     VARCHAR(50),
    audio_codec     VARCHAR(50),
    resolution      VARCHAR(20),

    oshash          VARCHAR(32),
    phash           VARCHAR(32),
    md5             VARCHAR(64),

    cover_path      TEXT,

    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(path)
);

CREATE TABLE c.scene_performers (
    scene_id        UUID REFERENCES c.scenes(id) ON DELETE CASCADE,
    performer_id    UUID REFERENCES c.performers(id) ON DELETE CASCADE,
    PRIMARY KEY (scene_id, performer_id)
);

CREATE TABLE c.scene_tags (
    scene_id        UUID REFERENCES c.scenes(id) ON DELETE CASCADE,
    tag_id          UUID REFERENCES c.tags(id) ON DELETE CASCADE,
    PRIMARY KEY (scene_id, tag_id)
);

-- Scene markers (timestamped tags from Stash)
CREATE TABLE c.scene_markers (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    scene_id        UUID REFERENCES c.scenes(id) ON DELETE CASCADE,
    title           VARCHAR(255),
    start_seconds   FLOAT NOT NULL,
    end_seconds     FLOAT,
    tag_id          UUID REFERENCES c.tags(id),
    stash_marker_id VARCHAR(100),
    created_at      TIMESTAMPTZ DEFAULT NOW()
);

-- Movie images
CREATE TABLE c.movie_images (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    movie_id        UUID REFERENCES c.movies(id) ON DELETE CASCADE,
    type            VARCHAR(50) NOT NULL,       -- poster, backdrop, screenshot
    path            TEXT NOT NULL,
    source          VARCHAR(50),                -- stashdb, whisparr, local
    primary_image   BOOLEAN DEFAULT FALSE,
    created_at      TIMESTAMPTZ DEFAULT NOW()
);

-- Galleries (image sets)
CREATE TABLE c.galleries (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    movie_id        UUID REFERENCES c.movies(id) ON DELETE SET NULL,
    title           VARCHAR(500) NOT NULL,
    path            TEXT,
    image_count     INT,
    created_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE c.gallery_images (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    gallery_id      UUID REFERENCES c.galleries(id) ON DELETE CASCADE,
    path            TEXT NOT NULL,
    position        INT,
    width           INT,
    height          INT
);

-- User data (per-module)
CREATE TABLE c.user_ratings (
    user_id         UUID NOT NULL,              -- Reference to public.users (no FK)
    movie_id        UUID REFERENCES c.movies(id) ON DELETE CASCADE,
    rating          INT CHECK (rating >= 0 AND rating <= 100),
    rated_at        TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (user_id, movie_id)
);

CREATE TABLE c.user_favorites (
    user_id         UUID NOT NULL,
    movie_id        UUID REFERENCES c.movies(id) ON DELETE CASCADE,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (user_id, movie_id)
);

CREATE TABLE c.watch_history (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL,
    movie_id        UUID REFERENCES c.movies(id) ON DELETE CASCADE,
    watched_at      TIMESTAMPTZ DEFAULT NOW(),
    position_ticks  BIGINT,                     -- Last position
    completed       BOOLEAN DEFAULT FALSE
);

-- Indexes
CREATE INDEX idx_c_movies_library ON c.movies(library_id);
CREATE INDEX idx_c_movies_studio ON c.movies(studio_id);
CREATE INDEX idx_c_movies_phash ON c.movies(phash);
CREATE INDEX idx_c_movies_oshash ON c.movies(oshash);
CREATE INDEX idx_c_performers_name ON c.performers(name);
CREATE INDEX idx_c_performers_stashdb ON c.performers(stashdb_id);
CREATE INDEX idx_c_watch_history_user ON c.watch_history(user_id, watched_at DESC);
```

---

## Scene-First Model (Whisparr v3)

Whisparr models scenes as episodes under a series/site. Revenge keeps a scene-first model:

- No `c.series`, `c.seasons`, or `c.episodes` tables.
- Scenes are stored directly in `c.scenes` and linked to `c.performers`, `c.studios`, and `c.tags`.
- Series/site grouping is represented via studio/site metadata and UI grouping, not separate tables.

---

## Features Specific to Adult Content

### Scene Markers (Timestamped Tags)

Import Stash scene markers as navigable chapters:

```go
type SceneMarkerService struct {
    db    *pgxpool.Pool
    stash *StashAppClient
}

type SceneMarker struct {
    ID          uuid.UUID
    MovieID     uuid.UUID
    Title       string
    StartTicks  int64
    EndTicks    *int64
    PrimaryTag  *Tag
    Tags        []Tag
    Thumbnail   string
}

// Import markers from Stash
func (s *SceneMarkerService) ImportFromStash(ctx context.Context, movieID uuid.UUID, stashSceneID string) error {
    markers, err := s.stash.GetSceneMarkers(ctx, stashSceneID)
    if err != nil {
        return err
    }

    for _, marker := range markers {
        _, err := s.db.Exec(ctx, `
            INSERT INTO c.scene_markers (movie_id, title, start_ticks, primary_tag_id, thumbnail_path)
            VALUES ($1, $2, $3, $4, $5)
        `, movieID, marker.Title, int64(marker.Seconds*10_000_000),
           s.ensureTag(ctx, marker.PrimaryTag), marker.Screenshot)
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
ALTER TABLE c.movies ADD COLUMN interactive BOOLEAN DEFAULT FALSE;
ALTER TABLE c.movies ADD COLUMN funscript_path TEXT;

CREATE TABLE c.interactive_scripts (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    movie_id        UUID REFERENCES c.movies(id) ON DELETE CASCADE,
    script_path     TEXT NOT NULL,
    script_type     VARCHAR(50) DEFAULT 'funscript',
    created_at      TIMESTAMPTZ DEFAULT NOW()
);
```

---

## Privacy & Security

### Access Control

```go
// Adult content requires special auth scope
const ScopeAdultContent = "adult:read"
const ScopeAdultWrite   = "adult:write"

func AdultAuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        claims := auth.ClaimsFromContext(r.Context())
        if !claims.HasScope(ScopeAdultContent) {
            http.Error(w, "Forbidden", http.StatusForbidden)
            return
        }
        next.ServeHTTP(w, r)
    })
}
```

### Audit Logging

All adult content access is logged:

```go
type AuditLogger struct {
    db *pgxpool.Pool
}

func (l *AuditLogger) LogAccess(ctx context.Context, userID uuid.UUID, action, resource string) {
    _, err := l.db.Exec(ctx, `
        INSERT INTO activity_log (user_id, action, resource, module, ip_address, user_agent, created_at)
        VALUES ($1, $2, $3, 'adult', $4, $5, NOW())
    `, userID, action, resource, getIP(ctx), getUserAgent(ctx))
    if err != nil {
        slog.Error("audit log failed", "error", err)
    }
}
```

### Analytics Exclusion

Option to exclude adult viewing from analytics:

```yaml
# config.yaml
adult:
  exclude_from_analytics: true
  exclude_from_recommendations: true
  require_pin: true           # Require PIN to access
  pin_timeout: 30m            # Re-prompt after timeout
```

---

## API Endpoints

All endpoints under obscured `/c/` namespace:

```
POST   /api/v1/c/movies                    # Create movie
GET    /api/v1/c/movies                    # List movies
GET    /api/v1/c/movies/{id}               # Get movie
PUT    /api/v1/c/movies/{id}               # Update movie
DELETE /api/v1/c/movies/{id}               # Delete movie

GET    /api/v1/c/movies/{id}/performers    # Get movie performers
GET    /api/v1/c/movies/{id}/markers       # Get scene markers
GET    /api/v1/c/movies/{id}/similar       # Get similar movies

GET    /api/v1/c/performers                # List performers
GET    /api/v1/c/performers/{id}           # Get performer
GET    /api/v1/c/performers/{id}/movies    # Get performer's movies

GET    /api/v1/c/studios                   # List studios
GET    /api/v1/c/studios/{id}              # Get studio
GET    /api/v1/c/studios/{id}/movies       # Get studio's movies

GET    /api/v1/c/tags                      # List tags
GET    /api/v1/c/tags/{id}/movies          # Get movies with tag

GET    /api/v1/c/scenes                    # List scenes
GET    /api/v1/c/scenes/{id}               # Get scene

POST   /api/v1/c/match                     # Match file by fingerprint
POST   /api/v1/c/identify                  # Submit for identification
```

---

## Configuration

```yaml
# config.yaml
adult:
  enabled: true

  # Whisparr integration
  whisparr:
    url: "http://whisparr:6969"
    api_key: "${WHISPARR_API_KEY}"

  # StashDB integration
  stashdb:
    endpoint: "https://stashdb.org/graphql"
    api_key: "${STASHDB_API_KEY}"

  # Stash-App integration (optional, for local instance)
  stash_app:
    url: "http://stash:9999"
    api_key: "${STASH_API_KEY}"

  # TPDB fallback
  tpdb:
    url: "https://api.metadataapi.net/api"
    api_key: "${TPDB_API_KEY}"

  # Fingerprinting
  fingerprint:
    generate_phash: true
    generate_oshash: true
    auto_match: true          # Auto-match on scan

  # Privacy
  privacy:
    exclude_from_analytics: true
    exclude_from_recommendations: true
    require_pin: false
    audit_all_access: true

  # Storage
  storage:
    images_path: "/data/adult/images"
    thumbnails_path: "/data/adult/thumbnails"
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
    // Get movie path
    var path string
    err := w.db.QueryRow(ctx, `SELECT path FROM c.movies WHERE id = $1`, job.Args.MovieID).Scan(&path)
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
        UPDATE c.movies SET phash = $1, oshash = $2 WHERE id = $3
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

## Summary

| Aspect | Implementation |
|--------|----------------|
| Schema | Isolated `c` PostgreSQL schema |
| API Namespace | Obscured `/c/` path |
| Primary Source | Whisparr-v3 for acquisition |
| Enrichment | StashDB (public), Stash-App (private) |
| Identification | pHash + oshash fingerprinting |
| Privacy | Audit logging, analytics exclusion, PIN protection |
| Scene Markers | Timestamped tags from Stash |
| Performer Data | StashDB with aliases, demographics, measurements |
