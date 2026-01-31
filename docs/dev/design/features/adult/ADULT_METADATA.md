# Revenge - Adult Content Metadata System

<!-- SOURCES: ffmpeg, ffmpeg-codecs, ffmpeg-formats, go-astiav, go-astiav-docs, pgx, postgresql-arrays, postgresql-json, river, stashdb, theporndb -->

<!-- DESIGN: features/adult, 01_ARCHITECTURE, 02_DESIGN_PRINCIPLES, 03_METADATA_SYSTEM -->


> ⚠️ **DEPRECATED**: This document has been merged into [ADULT_CONTENT_SYSTEM.md](ADULT_CONTENT_SYSTEM.md).
> See that document for the complete adult content architecture including metadata, privacy, API endpoints, and database schema.
> This file is kept for reference during migration only.


<!-- TOC-START -->

## Table of Contents

- [Status](#status)
- [Design Philosophy](#design-philosophy)
- [Architecture Overview](#architecture-overview)
- [Data Sources Priority](#data-sources-priority)
  - [Adult Movies (Scenes)](#adult-movies-scenes)
  - [Adult Scenes (Series)](#adult-scenes-series)
  - [Performers](#performers)
  - [Studios](#studios)
- [Whisparr Integration (Primary)](#whisparr-integration-primary)
  - [Client Implementation](#client-implementation)
- [Stash App Integration (Enrichment)](#stash-app-integration-enrichment)
  - [GraphQL Client](#graphql-client)
  - [Stash Scene Markers](#stash-scene-markers)
- [StashDB Integration (Community Database)](#stashdb-integration-community-database)
  - [GraphQL Client](#graphql-client)
  - [Fingerprint Generation](#fingerprint-generation)
- [TPDB Integration (Fallback)](#tpdb-integration-fallback)
- [Database Schema (`qar` Schema)](#database-schema-qar-schema)
- [Metadata Service](#metadata-service)
- [Configuration](#configuration)
- [Sources & Cross-References](#sources-cross-references)
  - [Cross-Reference Indexes](#cross-reference-indexes)
  - [Referenced Sources](#referenced-sources)
- [Related Design Docs](#related-design-docs)
  - [In This Section](#in-this-section)
  - [Related Topics](#related-topics)
  - [Indexes](#indexes)
- [Summary](#summary)

<!-- TOC-END -->

## Status

| Dimension | Status | Notes |
|-----------|--------|-------|
| Design | ✅ | Deprecated, merged into ADULT_CONTENT_SYSTEM |
| Sources | ✅ |  |
| Instructions | ✅ | See ADULT_CONTENT_SYSTEM.md |
| Code | 🔴 |  |
| Linting | 🔴 |  |
| Unit Testing | 🔴 |  |
| Integration Testing | 🔴 |  |---

> Complete metadata integration for adult content modules using Whisparr, Stash, and StashDB.

## Design Philosophy

1. **Whisparr as Primary Source** - Curated, cached metadata (like Radarr for movies)
2. **Stash App as Enrichment** - Local Stash instance provides additional metadata, organization
3. **StashDB as Fallback** - Community database for scenes, performers, studios
4. **Complete Isolation** - All data in `qar` PostgreSQL schema
5. **User Privacy** - No external calls without explicit consent, all data stays local

---

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    Adult Metadata Flow                                   │
└─────────────────────────────────────────────────────────────────────────┘

┌──────────┐     ┌───────────────────┐     ┌─────────────────┐     ┌──────────┐
│  Library │ ──→ │  Local Importers  │ ──→ │ Metadata Store  │ ←── │   UI     │
│   Scan   │     │    (Primary)      │     │   (c schema)    │     │ Request  │
└──────────┘     └───────────────────┘     └─────────────────┘     └──────────┘
                       │                          │
        ┌──────────────┼──────────────┐           │
        ▼              ▼              ▼           ▼
┌─────────────┐ ┌─────────────┐ ┌──────────┐ ┌─────────────────┐
│  Whisparr   │ │ Stash App   │ │ StashDB  │ │ Manual/Local    │
│  (v3 API)   │ │ (GraphQL)   │ │ (GraphQL)│ │ NFO Files       │
└─────────────┘ └─────────────┘ └──────────┘ └─────────────────┘
        │              │              │
        └──────────────┼──────────────┘
                       ▼
                ┌──────────────┐
                │   Missing?   │
                │   TPDB API   │
                └──────────────┘
```

---

## Data Sources Priority

### Adult Movies (Scenes)

| Priority | Source | Data Provided |
|----------|--------|---------------|
| 1 | **Whisparr** | Title, studio, performers, tags, cover, release date |
| 2 | **Stash App** | Local organization, user tags, O-counter, markers |
| 3 | **StashDB** | Community metadata, performer links, scene fingerprints |
| 4 | **TPDB** | ThePornDB API - extended metadata |
| 5 | **NFO Files** | Local metadata sidecar files |

### Adult Scenes (Series)

| Priority | Source | Data Provided |
|----------|--------|---------------|
| 1 | **Whisparr** | Series info, episodes, performers |
| 2 | **Stash App** | Episode organization, tags |
| 3 | **StashDB** | Series database |
| 4 | **TPDB** | Series metadata |

### Performers

| Priority | Source | Data Provided |
|----------|--------|---------------|
| 1 | **Stash App** | Local performer data with user edits |
| 2 | **StashDB** | Comprehensive performer database |
| 3 | **Whisparr** | Basic performer info |
| 4 | **TPDB** | Performer profiles, measurements, social links |
| 5 | **FreeOnes** | Additional performer data |

### Studios

| Priority | Source | Data Provided |
|----------|--------|---------------|
| 1 | **Whisparr** | Studio list from downloads |
| 2 | **StashDB** | Studio database with parent/network info |
| 3 | **TPDB** | Studio profiles |

---

## Whisparr Integration (Primary)

Whisparr v3 provides the same API structure as Radarr, making integration straightforward.

### Client Implementation

```go
// WhisparrClient manages communication with Whisparr v3 API
type WhisparrClient struct {
    baseURL string
    apiKey  string
    client  *http.Client
    logger  *slog.Logger
}

type WhisparrMovie struct {
    ID              int       `json:"id"`
    Title           string    `json:"title"`
    SortTitle       string    `json:"sortTitle"`
    Overview        string    `json:"overview"`
    Studio          string    `json:"studio"`
    StudioForeignId string    `json:"studioForeignId"`
    Year            int       `json:"year"`
    ReleaseDate     string    `json:"releaseDate"`
    Runtime         int       `json:"runtime"`
    Genres          []string  `json:"genres"`
    Tags            []int     `json:"tags"`
    Monitored       bool      `json:"monitored"`
    HasFile         bool      `json:"hasFile"`
    Path            string    `json:"path"`
    MovieFile       *MovieFile `json:"movieFile,omitempty"`

    // Performer info
    Credits         []WhisparrCredit `json:"credits"`

    // Images
    Images          []WhisparrImage `json:"images"`

    // External IDs
    TpdbId          int    `json:"tpdbId,omitempty"`
    StashId         string `json:"stashId,omitempty"`
}

type WhisparrCredit struct {
    PersonName      string `json:"personName"`
    Character       string `json:"character"` // Role in scene
    ForeignId       string `json:"foreignId"` // StashDB/TPDB ID
    Type            string `json:"type"`      // "performer", "director"
    Images          []WhisparrImage `json:"images"`
}

// Sync from Whisparr
func (c *WhisparrClient) SyncMovies(ctx context.Context) ([]WhisparrMovie, error) {
    resp, err := c.get(ctx, "/api/v3/movie")
    if err != nil {
        return nil, fmt.Errorf("whisparr sync: %w", err)
    }

    var movies []WhisparrMovie
    if err := json.Unmarshal(resp, &movies); err != nil {
        return nil, fmt.Errorf("parse whisparr response: %w", err)
    }

    return movies, nil
}

// Webhook handler for real-time updates
func (c *WhisparrClient) HandleWebhook(event WebhookEvent) error {
    switch event.EventType {
    case "Download":
        return c.handleDownload(event)
    case "Rename":
        return c.handleRename(event)
    case "MovieDelete":
        return c.handleDelete(event)
    }
    return nil
}
```

---

## Stash App Integration (Enrichment)

Stash is a self-hosted app for organizing adult content. We use its GraphQL API for enrichment.

### GraphQL Client

```go
// StashClient communicates with local Stash instance
type StashClient struct {
    endpoint string
    apiKey   string
    client   *graphql.Client
    logger   *slog.Logger
}

// Stash GraphQL types
type StashScene struct {
    ID           string           `json:"id"`
    Title        string           `json:"title"`
    Details      string           `json:"details"`
    URL          string           `json:"url"`
    Date         string           `json:"date"`
    Rating100    *int             `json:"rating100"`
    OCounter     int              `json:"o_counter"`
    Organized    bool             `json:"organized"`

    // File info
    Files        []StashFile      `json:"files"`

    // Relationships
    Studio       *StashStudio     `json:"studio"`
    Performers   []StashPerformer `json:"performers"`
    Tags         []StashTag       `json:"tags"`
    Markers      []StashMarker    `json:"scene_markers"`

    // StashDB link
    StashIDs     []StashID        `json:"stash_ids"`
}

type StashPerformer struct {
    ID           string   `json:"id"`
    Name         string   `json:"name"`
    Disambiguation string `json:"disambiguation"`
    Aliases      string   `json:"aliases"`
    Gender       string   `json:"gender"`
    Birthdate    string   `json:"birthdate"`
    Ethnicity    string   `json:"ethnicity"`
    Country      string   `json:"country"`
    HairColor    string   `json:"hair_color"`
    EyeColor     string   `json:"eye_color"`
    Height       int      `json:"height_cm"`
    Weight       int      `json:"weight"`
    Measurements string   `json:"measurements"`
    FakeTits     string   `json:"fake_tits"`
    Tattoos      string   `json:"tattoos"`
    Piercings    string   `json:"piercings"`
    CareerLength string   `json:"career_length"`
    Details      string   `json:"details"`
    DeathDate    string   `json:"death_date"`

    // Images
    ImagePath    string   `json:"image_path"`

    // StashDB link
    StashIDs     []StashID `json:"stash_ids"`
}

type StashMarker struct {
    ID           string `json:"id"`
    Title        string `json:"title"`
    Seconds      float64 `json:"seconds"`
    PrimaryTag   StashTag `json:"primary_tag"`
}

type StashID struct {
    Endpoint     string `json:"endpoint"`
    StashID      string `json:"stash_id"`
}

// Query scenes by file path
const findSceneByPathQuery = `
query FindSceneByPath($path: String!) {
    findScenes(scene_filter: { path: { value: $path, modifier: EQUALS } }) {
        scenes {
            id
            title
            details
            date
            rating100
            o_counter
            studio { id name }
            performers { id name gender }
            tags { id name }
            scene_markers { id title seconds primary_tag { name } }
            stash_ids { endpoint stash_id }
        }
    }
}
`

func (c *StashClient) GetSceneByPath(ctx context.Context, path string) (*StashScene, error) {
    var resp struct {
        FindScenes struct {
            Scenes []StashScene `json:"scenes"`
        } `json:"findScenes"`
    }

    err := c.client.Query(ctx, findSceneByPathQuery, map[string]interface{}{
        "path": path,
    }, &resp)

    if err != nil {
        return nil, err
    }

    if len(resp.FindScenes.Scenes) == 0 {
        return nil, ErrNotFound
    }

    return &resp.FindScenes.Scenes[0], nil
}

// Sync all performers
func (c *StashClient) GetAllPerformers(ctx context.Context) ([]StashPerformer, error) {
    const query = `
    query AllPerformers {
        allPerformers {
            id name disambiguation aliases gender birthdate
            ethnicity country hair_color eye_color height_cm
            weight measurements fake_tits tattoos piercings
            career_length details death_date image_path
            stash_ids { endpoint stash_id }
        }
    }
    `

    var resp struct {
        AllPerformers []StashPerformer `json:"allPerformers"`
    }

    err := c.client.Query(ctx, query, nil, &resp)
    return resp.AllPerformers, err
}
```

### Stash Scene Markers

Stash markers provide chapter-like functionality for adult content:

```go
// Convert Stash markers to our chapter format
func (s *AdultMetadataService) ConvertMarkersToChapters(markers []StashMarker) []Chapter {
    chapters := make([]Chapter, len(markers))

    for i, m := range markers {
        chapters[i] = Chapter{
            Index:     i,
            Title:     m.Title,
            StartTime: m.Seconds,
            Tag:       m.PrimaryTag.Name,
            Source:    "stash",
        }
    }

    // Sort by time
    sort.Slice(chapters, func(i, j int) bool {
        return chapters[i].StartTime < chapters[j].StartTime
    })

    return chapters
}
```

---

## StashDB Integration (Community Database)

StashDB is a community-maintained database with scene fingerprints and comprehensive metadata.

### GraphQL Client

```go
// StashDBClient communicates with StashDB instances
type StashDBClient struct {
    endpoint string // https://stashdb.org/graphql
    apiKey   string
    client   *graphql.Client
    logger   *slog.Logger
}

// Scene fingerprint matching
type StashDBFingerprint struct {
    Algorithm   string `json:"algorithm"`   // "PHASH", "OSHASH", "MD5"
    Hash        string `json:"hash"`
    Duration    int    `json:"duration"`
}

type StashDBScene struct {
    ID           string               `json:"id"`
    Title        string               `json:"title"`
    Details      string               `json:"details"`
    Date         string               `json:"date"`
    Duration     int                  `json:"duration"`
    Director     string               `json:"director"`
    Code         string               `json:"code"`        // Scene code/ID from studio

    URLs         []StashDBURL         `json:"urls"`
    Studio       *StashDBStudio       `json:"studio"`
    Performers   []StashDBPerformer   `json:"performers"`
    Tags         []StashDBTag         `json:"tags"`
    Fingerprints []StashDBFingerprint `json:"fingerprints"`
    Images       []StashDBImage       `json:"images"`
}

// Query scene by fingerprint (most reliable matching)
const findSceneByFingerprintQuery = `
query FindSceneByFingerprint($fingerprints: [FingerprintInput!]!) {
    findScenesByFingerprints(fingerprints: $fingerprints) {
        id
        title
        details
        date
        duration
        studio { id name }
        performers { performer { id name gender } }
        tags { id name }
        images { id url width height }
    }
}
`

func (c *StashDBClient) FindByFingerprint(ctx context.Context, fp Fingerprint) (*StashDBScene, error) {
    var resp struct {
        FindScenes []StashDBScene `json:"findScenesByFingerprints"`
    }

    input := []map[string]interface{}{{
        "algorithm": fp.Algorithm,
        "hash":      fp.Hash,
        "duration":  fp.Duration,
    }}

    err := c.client.Query(ctx, findSceneByFingerprintQuery, map[string]interface{}{
        "fingerprints": input,
    }, &resp)

    if err != nil {
        return nil, err
    }

    if len(resp.FindScenes) == 0 {
        return nil, ErrNotFound
    }

    return &resp.FindScenes[0], nil
}

// Get performer by StashDB ID
func (c *StashDBClient) GetPerformer(ctx context.Context, id string) (*StashDBPerformer, error) {
    const query = `
    query GetPerformer($id: ID!) {
        findPerformer(id: $id) {
            id
            name
            disambiguation
            aliases
            gender
            birth_date
            ethnicity
            country
            hair_color
            eye_color
            height
            measurements { cup_size band_size waist hip }
            breast_type
            tattoos { location description }
            piercings { location description }
            career_start_year
            career_end_year
            urls { url type }
            images { id url }
        }
    }
    `

    var resp struct {
        FindPerformer StashDBPerformer `json:"findPerformer"`
    }

    err := c.client.Query(ctx, query, map[string]interface{}{"id": id}, &resp)
    return &resp.FindPerformer, err
}
```

### Fingerprint Generation

```go
// FingerprintService generates scene fingerprints for StashDB matching
type FingerprintService struct {
    ffmpeg string
    logger *slog.Logger
}

type Fingerprint struct {
    Algorithm string
    Hash      string
    Duration  int
}

// Generate OSHASH (Stash/StashDB standard)
func (s *FingerprintService) GenerateOSHash(path string) (string, error) {
    file, err := os.Open(path)
    if err != nil {
        return "", err
    }
    defer file.Close()

    stat, _ := file.Stat()
    size := stat.Size()

    if size < 65536*2 {
        return "", errors.New("file too small for oshash")
    }

    // Read first and last 64KB
    head := make([]byte, 65536)
    tail := make([]byte, 65536)

    file.Read(head)
    file.Seek(-65536, io.SeekEnd)
    file.Read(tail)

    // Calculate hash
    var hash uint64 = uint64(size)
    for i := 0; i < 65536; i += 8 {
        hash += binary.LittleEndian.Uint64(head[i:])
        hash += binary.LittleEndian.Uint64(tail[i:])
    }

    return fmt.Sprintf("%016x", hash), nil
}

// Generate PHASH (perceptual hash) via FFmpeg
func (s *FingerprintService) GeneratePHash(ctx context.Context, path string) (string, error) {
    // Extract frame at 10% position
    cmd := exec.CommandContext(ctx, s.ffmpeg,
        "-ss", "10%",
        "-i", path,
        "-vframes", "1",
        "-f", "image2pipe",
        "-vcodec", "png",
        "-",
    )

    output, err := cmd.Output()
    if err != nil {
        return "", err
    }

    // Calculate perceptual hash
    img, _ := png.Decode(bytes.NewReader(output))
    hash := phash.DCT(img)

    return fmt.Sprintf("%016x", hash), nil
}

// Generate all fingerprints for a file
func (s *FingerprintService) GenerateAll(ctx context.Context, path string) ([]Fingerprint, error) {
    var fps []Fingerprint

    // OSHASH (fast, file-based)
    if oshash, err := s.GenerateOSHash(path); err == nil {
        duration := s.getDuration(ctx, path)
        fps = append(fps, Fingerprint{"OSHASH", oshash, duration})
    }

    // PHASH (perceptual, content-based)
    if phash, err := s.GeneratePHash(ctx, path); err == nil {
        duration := s.getDuration(ctx, path)
        fps = append(fps, Fingerprint{"PHASH", phash, duration})
    }

    // MD5 (slow but definitive)
    if md5hash, err := s.GenerateMD5(path); err == nil {
        duration := s.getDuration(ctx, path)
        fps = append(fps, Fingerprint{"MD5", md5hash, duration})
    }

    return fps, nil
}
```

---

## TPDB Integration (Fallback)

ThePornDB provides extensive metadata when other sources fail.

```go
// TPDBClient communicates with ThePornDB API
type TPDBClient struct {
    baseURL string
    apiKey  string
    client  *http.Client
    logger  *slog.Logger
}

type TPDBScene struct {
    ID          string   `json:"id"`
    Title       string   `json:"title"`
    Description string   `json:"description"`
    Date        string   `json:"date"`
    Duration    int      `json:"duration"`
    Site        TPDBSite `json:"site"`
    Performers  []TPDBPerformer `json:"performers"`
    Tags        []string `json:"tags"`
    Image       string   `json:"image"`
    Trailer     string   `json:"trailer,omitempty"`
}

type TPDBPerformer struct {
    ID           string   `json:"id"`
    Name         string   `json:"name"`
    Bio          string   `json:"bio"`
    Gender       string   `json:"gender"`
    Birthdate    string   `json:"birthdate"`
    Birthplace   string   `json:"birthplace"`
    Nationality  string   `json:"nationality"`
    Ethnicity    string   `json:"ethnicity"`
    HairColor    string   `json:"hair_color"`
    EyeColor     string   `json:"eye_color"`
    Height       int      `json:"height"`
    Weight       int      `json:"weight"`
    Measurements string   `json:"measurements"`
    Tattoos      string   `json:"tattoos"`
    Piercings    string   `json:"piercings"`
    Image        string   `json:"image"`
    Aliases      []string `json:"aliases"`

    // Social links
    Twitter      string   `json:"twitter,omitempty"`
    Instagram    string   `json:"instagram,omitempty"`
    OnlyFans     string   `json:"onlyfans,omitempty"`
}

// Search by title and performers
func (c *TPDBClient) SearchScene(ctx context.Context, title string, performers []string) (*TPDBScene, error) {
    query := url.Values{}
    query.Set("q", title)
    if len(performers) > 0 {
        query.Set("performers", strings.Join(performers, ","))
    }

    resp, err := c.get(ctx, "/scenes?"+query.Encode())
    if err != nil {
        return nil, err
    }

    var result struct {
        Data []TPDBScene `json:"data"`
    }
    json.Unmarshal(resp, &result)

    if len(result.Data) == 0 {
        return nil, ErrNotFound
    }

    return &result.Data[0], nil
}

// Get performer by name with fuzzy matching
func (c *TPDBClient) SearchPerformer(ctx context.Context, name string) (*TPDBPerformer, error) {
    resp, err := c.get(ctx, "/performers?q="+url.QueryEscape(name))
    if err != nil {
        return nil, err
    }

    var result struct {
        Data []TPDBPerformer `json:"data"`
    }
    json.Unmarshal(resp, &result)

    if len(result.Data) == 0 {
        return nil, ErrNotFound
    }

    return &result.Data[0], nil
}
```

---

## Database Schema (`qar` Schema)

All adult content data is stored in the isolated `qar` PostgreSQL schema (Queen Anne's Revenge).

```sql
-- Voyages (scenes - main content)
CREATE TABLE qar.voyages (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    fleet_id        UUID NOT NULL,               -- library
    title           VARCHAR(500) NOT NULL,
    sort_title      VARCHAR(500),
    overview        TEXT,
    launch_date     DATE,                        -- release_date (obfuscated)
    distance        INT,                         -- runtime_minutes (obfuscated)
    port_id         UUID REFERENCES qar.ports(id), -- studio

    -- External IDs (obfuscated)
    whisparr_id     INT,
    stash_id        VARCHAR(100),
    charter         VARCHAR(100),                -- stashdb_id
    registry        VARCHAR(100),                -- tpdb_id

    -- File info
    path            TEXT NOT NULL,
    size_bytes      BIGINT,
    video_codec     VARCHAR(50),
    audio_codec     VARCHAR(50),
    resolution      VARCHAR(20),

    -- Fingerprints for matching (obfuscated)
    oshash          VARCHAR(32),
    coordinates     VARCHAR(32),                 -- phash
    md5             VARCHAR(64),

    -- Images
    cover_path      TEXT,

    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_qar_voyages_fleet ON qar.voyages(fleet_id);
CREATE INDEX idx_qar_voyages_port ON qar.voyages(port_id);
CREATE INDEX idx_qar_voyages_oshash ON qar.voyages(oshash);
CREATE INDEX idx_qar_voyages_charter ON qar.voyages(charter);

-- Crew (performers)
CREATE TABLE qar.crew (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name            VARCHAR(255) NOT NULL,
    disambiguation  VARCHAR(255),
    gender          VARCHAR(50),
    christening     DATE,                        -- birthdate (obfuscated)
    death_date      DATE,
    birth_city      VARCHAR(255),
    origin          VARCHAR(100),                -- ethnicity (obfuscated)
    nationality     VARCHAR(100),
    rigging         VARCHAR(50),                 -- hair_color (obfuscated)
    compass         VARCHAR(50),                 -- eye_color (obfuscated)
    height_cm       INT,
    weight_kg       INT,
    cargo           JSONB,                       -- measurements (obfuscated)
    cup_size        VARCHAR(10),
    breast_type     VARCHAR(50),
    markings        TEXT[],                      -- tattoos (obfuscated)
    anchors         TEXT[],                      -- piercings (obfuscated)
    maiden_voyage   INT,                         -- career_start (obfuscated)
    last_port       INT,                         -- career_end (obfuscated)
    bio             TEXT,

    -- External IDs (obfuscated)
    stash_id        VARCHAR(100),
    charter         VARCHAR(100),                -- stashdb_id
    registry        VARCHAR(100),                -- tpdb_id
    manifest        VARCHAR(100),                -- freeones_id

    -- Social links (encrypted)
    twitter         TEXT,
    instagram       TEXT,

    -- Images
    image_path      TEXT,

    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_qar_crew_name ON qar.crew(name);
CREATE INDEX idx_qar_crew_charter ON qar.crew(charter);

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

-- Voyage-Crew relationship (scene-performer)
CREATE TABLE qar.voyage_crew (
    voyage_id       UUID REFERENCES qar.voyages(id) ON DELETE CASCADE,
    crew_id         UUID REFERENCES qar.crew(id) ON DELETE CASCADE,
    role            VARCHAR(100), -- Character name if applicable
    PRIMARY KEY (voyage_id, crew_id)
);

-- Ports (studios)
CREATE TABLE qar.ports (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name            VARCHAR(255) NOT NULL,
    parent_id       UUID REFERENCES qar.ports(id), -- For port networks
    url             TEXT,

    -- External IDs (obfuscated)
    charter         VARCHAR(100),                -- stashdb_id
    registry        VARCHAR(100),                -- tpdb_id

    -- Images
    logo_path       TEXT,

    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);

-- Flags (tags)
CREATE TABLE qar.flags (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name            VARCHAR(255) NOT NULL UNIQUE,
    description     TEXT,
    parent_id       UUID REFERENCES qar.flags(id), -- Hierarchical flags
    charter         VARCHAR(100),                -- stashdb_id
    created_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE qar.voyage_flags (
    voyage_id       UUID REFERENCES qar.voyages(id) ON DELETE CASCADE,
    flag_id         UUID REFERENCES qar.flags(id) ON DELETE CASCADE,
    PRIMARY KEY (voyage_id, flag_id)
);

-- Voyage markers (scene markers - chapters/positions)
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

CREATE INDEX idx_qar_markers_voyage ON qar.voyage_markers(voyage_id);

-- User data (per-module, in qar schema)
CREATE TABLE qar.user_voyage_data (
    user_id         UUID REFERENCES users(id) ON DELETE CASCADE,
    voyage_id       UUID REFERENCES qar.voyages(id) ON DELETE CASCADE,

    -- Watch progress
    position_ms     BIGINT DEFAULT 0,
    watch_count     INT DEFAULT 0,
    last_watched    TIMESTAMPTZ,

    -- Rating (bounty)
    bounty          SMALLINT CHECK (bounty >= 1 AND bounty <= 10),

    -- Stash-style O-counter
    o_counter       INT DEFAULT 0,

    -- Favorite
    is_favorite     BOOLEAN DEFAULT false,

    -- Organization
    is_organized    BOOLEAN DEFAULT false,

    PRIMARY KEY (user_id, voyage_id)
);

-- User crew favorites (performer favorites)
CREATE TABLE qar.user_crew_favorites (
    user_id         UUID REFERENCES users(id) ON DELETE CASCADE,
    crew_id         UUID REFERENCES qar.crew(id) ON DELETE CASCADE,
    added_at        TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (user_id, crew_id)
);
```

---

## Metadata Service

```go
// AdultMetadataService coordinates all adult metadata sources
type AdultMetadataService struct {
    whisparr    *WhisparrClient
    stash       *StashClient
    stashdb     *StashDBClient
    tpdb        *TPDBClient
    fingerprint *FingerprintService
    sceneRepo   SceneRepository
    performerRepo PerformerRepository
    logger      *slog.Logger
}

// EnrichScene fetches metadata from all sources
func (s *AdultMetadataService) EnrichScene(ctx context.Context, sceneID uuid.UUID) error {
    scene, err := s.sceneRepo.GetByID(ctx, sceneID)
    if err != nil {
        return err
    }

    // 1. Try Stash (local)
    if stashScene, err := s.stash.GetSceneByPath(ctx, scene.Path); err == nil {
        s.mergeStashData(scene, stashScene)

        // Get markers as chapters
        if len(stashScene.Markers) > 0 {
            chapters := s.ConvertMarkersToChapters(stashScene.Markers)
            s.sceneRepo.SaveChapters(ctx, sceneID, chapters)
        }
    }

    // 2. Try StashDB (community)
    if scene.OsHash != "" {
        fp := Fingerprint{Algorithm: "OSHASH", Hash: scene.OsHash, Duration: scene.Runtime * 60}
        if stashdbScene, err := s.stashdb.FindByFingerprint(ctx, fp); err == nil {
            s.mergeStashDBData(scene, stashdbScene)
        }
    }

    // 3. Try TPDB (fallback)
    if scene.Overview == "" {
        performers := s.getPerformerNames(ctx, sceneID)
        if tpdbScene, err := s.tpdb.SearchScene(ctx, scene.Title, performers); err == nil {
            s.mergeTPDBData(scene, tpdbScene)
        }
    }

    return s.sceneRepo.Update(ctx, scene)
}

// SyncFromWhisparr imports all scenes from Whisparr
func (s *AdultMetadataService) SyncFromWhisparr(ctx context.Context) error {
    movies, err := s.whisparr.SyncMovies(ctx)
    if err != nil {
        return err
    }

    for _, m := range movies {
        if !m.HasFile {
            continue
        }

        scene := s.whisparrToScene(m)

        // Generate fingerprints for StashDB matching
        fps, _ := s.fingerprint.GenerateAll(ctx, scene.Path)
        if len(fps) > 0 {
            for _, fp := range fps {
                switch fp.Algorithm {
                case "OSHASH":
                    scene.OsHash = fp.Hash
                case "PHASH":
                    scene.PHash = fp.Hash
                case "MD5":
                    scene.MD5 = fp.Hash
                }
            }
        }

        // Upsert scene
        if err := s.sceneRepo.Upsert(ctx, scene); err != nil {
            s.logger.Error("failed to upsert scene", "title", scene.Title, "error", err)
            continue
        }

        // Queue enrichment job
        s.queueEnrichment(ctx, scene.ID)
    }

    return nil
}
```

---

## Configuration

```yaml
# configs/config.yaml

adult:
  enabled: true

  whisparr:
    enabled: true
    url: "http://whisparr:6969"
    api_key: "${WHISPARR_API_KEY}"
    sync_interval: "1h"

  stash:
    enabled: true
    url: "http://stash:9999"
    api_key: "${STASH_API_KEY}"
    # Sync user data (o-counter, organized, markers)
    sync_user_data: true

  stashdb:
    enabled: true
    endpoint: "https://stashdb.org/graphql"
    api_key: "${STASHDB_API_KEY}"
    # Match scenes by fingerprint
    fingerprint_matching: true

  tpdb:
    enabled: true
    url: "https://api.theporndb.net"
    api_key: "${TPDB_API_KEY}"

  # Privacy settings
  privacy:
    # Never call external APIs without user consent
    require_consent: true
    # Cache all images locally
    local_images: true
    # No external analytics
    no_tracking: true

  # Fingerprinting
  fingerprints:
    generate_oshash: true
    generate_phash: true
    generate_md5: false  # Slow, optional
```


---

## Summary

| Source | Role | Data Provided |
|--------|------|---------------|
| Whisparr | Primary (proxy) | Scenes, studios, basic performer info |
| Stash App | Enrichment | Local organization, markers, O-counter |
| StashDB | Community DB | Fingerprint matching, detailed metadata |
| TPDB | Fallback | Extended metadata when others fail |
| NFO Files | Local | User-created metadata sidecar files |

| Feature | Implementation |
|---------|----------------|
| Scene Matching | Fingerprints (OSHASH, PHASH) → StashDB |
| Performer Data | StashDB + TPDB (comprehensive profiles) |
| Chapter Markers | Stash markers → chapters |
| User Data | O-counter, favorites, ratings (isolated) |
| Privacy | All data stays local, consent required |
