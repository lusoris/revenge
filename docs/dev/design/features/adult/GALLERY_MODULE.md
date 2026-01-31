# Adult Gallery Module (QAR: Treasures)

> Image gallery management for adult content with performer links and Prowlarr integration


<!-- TOC-START -->

## Table of Contents

- [Status](#status)
- [Developer Resources](#developer-resources)
- [QAR Fleet Types](#qar-fleet-types)
- [Overview](#overview)
  - [Key Features](#key-features)
- [QAR Obfuscation](#qar-obfuscation)
  - [Chest Types (Gallery Categories)](#chest-types-gallery-categories)
- [Database Schema](#database-schema)
- [Entity Relationships](#entity-relationships)
- [Go Entities](#go-entities)
- [Prowlarr Integration (Download Search)](#prowlarr-integration-download-search)
  - [Configuration](#configuration)
  - [Prowlarr Client](#prowlarr-client)
  - [River Jobs for Download Processing](#river-jobs-for-download-processing)
- [Search Integration (Typesense)](#search-integration-typesense)
- [API Endpoints](#api-endpoints)
- [Implementation Checklist](#implementation-checklist)
  - [Phase 1: Schema & Repository](#phase-1-schema-repository)
  - [Phase 2: Core Service](#phase-2-core-service)
  - [Phase 3: SABnzbd Integration](#phase-3-sabnzbd-integration)
  - [Phase 4: Image Processing](#phase-4-image-processing)
  - [Phase 5: Search](#phase-5-search)
  - [Phase 6: API](#phase-6-api)
- [Sources & Cross-References](#sources-cross-references)
  - [Cross-Reference Indexes](#cross-reference-indexes)
  - [Referenced Sources](#referenced-sources)
- [Related Design Docs](#related-design-docs)
  - [In This Section](#in-this-section)
  - [Related Topics](#related-topics)
  - [Indexes](#indexes)
- [Related Documentation](#related-documentation)

<!-- TOC-END -->

## Status

| Dimension | Status |
|-----------|--------|
| Design | ✅ |
| Sources | ✅ |
| Instructions | ✅ |
| Code | 🔴 |
| Linting | 🔴 |
| Unit Testing | 🔴 |
| Integration Testing | 🔴 |
**Priority**: 🟡 MEDIUM (Phase 7 - Adult Enhancements)
**Schema**: `qar` (Queen Anne's Revenge isolated schema)
**API Namespace**: `/api/v1/legacy/treasures` (obfuscated)
**Dependencies**: Adult Content System, Prowlarr, Typesense

---

## Developer Resources

> Package versions: [00_SOURCE_OF_TRUTH.md](../../00_SOURCE_OF_TRUTH.md#go-dependencies-core)

| Source | URL | Purpose |
|--------|-----|---------|
| Typesense | [typesense.org/docs](https://typesense.org/docs/) | Gallery search |
| Prowlarr API | [wiki.servarr.com/prowlarr](https://wiki.servarr.com/prowlarr) | Download search |
| SABnzbd API | [sabnzbd.org/wiki](https://sabnzbd.org/wiki/configuration/4.3/api) | Usenet download client |

| Package | Purpose |
|---------|---------|
| goimagehash | pHash image deduplication |
| govips | Thumbnail/WebP generation |
| go-blurhash | Placeholder generation |

---

## QAR Fleet Types

The QAR module supports **three isolated content libraries** (fleets):

| Fleet Type | Content | Real Name | Description |
|------------|---------|-----------|-------------|
| `voyage` | Scenes | Scene Library | Individual adult scenes |
| `expedition` | Movies | Movie Library | Full-length adult movies |
| `treasure` | Galleries | Gallery Library | Image collections/photosets |

Each fleet type is a separate library that users can create and manage independently.

---

## Overview

The Gallery module manages image collections (photosets, behind-the-scenes, promotional images) linked to performers and scenes. Following the QAR obfuscation pattern, galleries are called **Treasures** and individual images are **Doubloons**.

### Key Features

1. **Performer Galleries** - Image collections linked to specific crew members
2. **Scene Galleries** - Behind-the-scenes or promotional images linked to voyages
3. **Standalone Galleries** - Independent photosets (magazine shoots, promotional content)
4. **Search & Discovery** - Full-text search across gallery metadata via Typesense
5. **SABnzbd Integration** - Automated downloading from Usenet sources
6. **Duplicate Detection** - pHash-based image deduplication

---

## QAR Obfuscation

| Real Concept | Obfuscated | Database Table | API Endpoint |
|--------------|------------|----------------|--------------|
| Gallery | **Treasure** | `qar.treasures` | `/legacy/treasures` |
| Gallery Image | **Doubloon** | `qar.doubloons` | `/legacy/treasures/{id}/doubloons` |
| Gallery Type | **Chest Type** | `chest_type` column | - |
| Image Set | **Haul** | - | - |
| Download Queue | **Plunder** | `qar.plunder_queue` | `/legacy/plunder` |

### Chest Types (Gallery Categories)

| Real Type | Obfuscated | Description |
|-----------|------------|-------------|
| Performer Gallery | `crew_chest` | Images of a specific performer |
| Scene Gallery | `voyage_chest` | BTS/promotional for a scene |
| Photoset | `merchant_haul` | Standalone photo collection |
| Magazine | `royal_treasury` | Professional magazine shoots |
| Promotional | `port_display` | Studio promotional content |

---

## Database Schema

```sql
-- Galleries → Treasures
CREATE TABLE qar.treasures (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    fleet_id        UUID NOT NULL REFERENCES qar.fleets(id),  -- Library

    -- Metadata
    title           VARCHAR(500) NOT NULL,
    description     TEXT,
    chest_type      VARCHAR(50) NOT NULL DEFAULT 'merchant_haul',

    -- Relationships (optional - can be standalone)
    voyage_id       UUID REFERENCES qar.voyages(id) ON DELETE SET NULL,      -- Linked scene
    expedition_id   UUID REFERENCES qar.expeditions(id) ON DELETE SET NULL,  -- Linked movie
    port_id         UUID REFERENCES qar.ports(id) ON DELETE SET NULL,        -- Studio

    -- Source info
    source_url      TEXT,                              -- Original source URL
    source_type     VARCHAR(50),                       -- usenet, web, local
    sabnzbd_nzo_id  VARCHAR(100),                      -- SABnzbd job ID

    -- Dates
    release_date    DATE,
    imported_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Stats
    doubloon_count  INT NOT NULL DEFAULT 0,            -- Image count
    total_size      BIGINT NOT NULL DEFAULT 0,         -- Total bytes

    -- Storage
    path            TEXT NOT NULL,                     -- Filesystem path
    cover_path      TEXT,                              -- Cover/thumbnail image

    -- External IDs
    stashdb_id      VARCHAR(100),
    tpdb_id         VARCHAR(100),

    -- Timestamps
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT valid_chest_type CHECK (chest_type IN (
        'crew_chest', 'voyage_chest', 'merchant_haul', 'royal_treasury', 'port_display'
    ))
);

-- Gallery Images → Doubloons
CREATE TABLE qar.doubloons (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    treasure_id     UUID NOT NULL REFERENCES qar.treasures(id) ON DELETE CASCADE,

    -- File info
    filename        VARCHAR(500) NOT NULL,
    path            TEXT NOT NULL,
    position        INT NOT NULL DEFAULT 0,            -- Sort order within gallery

    -- Image metadata
    width           INT,
    height          INT,
    size_bytes      BIGINT,
    format          VARCHAR(20),                       -- jpg, png, webp, gif

    -- Fingerprinting (deduplication)
    phash           VARCHAR(64),                       -- Perceptual hash
    md5             VARCHAR(64),
    blurhash        VARCHAR(100),                      -- For 🔴 Not implemented

    -- Processing
    thumbnail_path  TEXT,                              -- Generated thumbnail
    webp_path       TEXT,                              -- WebP conversion

    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE(treasure_id, path)
);

-- Treasure ↔ Crew (Many-to-Many: Galleries can feature multiple performers)
CREATE TABLE qar.treasure_crew (
    treasure_id     UUID NOT NULL REFERENCES qar.treasures(id) ON DELETE CASCADE,
    crew_id         UUID NOT NULL REFERENCES qar.crew(id) ON DELETE CASCADE,
    is_primary      BOOLEAN NOT NULL DEFAULT FALSE,    -- Primary featured performer
    PRIMARY KEY (treasure_id, crew_id)
);

-- Treasure ↔ Flags (Tags)
CREATE TABLE qar.treasure_flags (
    treasure_id     UUID NOT NULL REFERENCES qar.treasures(id) ON DELETE CASCADE,
    flag_id         UUID NOT NULL REFERENCES qar.flags(id) ON DELETE CASCADE,
    PRIMARY KEY (treasure_id, flag_id)
);

-- User data: Favorites
CREATE TABLE qar.treasure_favorites (
    user_id         UUID NOT NULL,                     -- Reference to shared.users
    treasure_id     UUID NOT NULL REFERENCES qar.treasures(id) ON DELETE CASCADE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, treasure_id)
);

-- User data: Ratings
CREATE TABLE qar.treasure_ratings (
    user_id         UUID NOT NULL,
    treasure_id     UUID NOT NULL REFERENCES qar.treasures(id) ON DELETE CASCADE,
    bounty          INT NOT NULL CHECK (bounty >= 0 AND bounty <= 100),  -- Rating 0-100
    rated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, treasure_id)
);

-- Download Queue → Plunder Queue (SABnzbd integration)
CREATE TABLE qar.plunder_queue (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    fleet_id        UUID NOT NULL REFERENCES qar.fleets(id),

    -- NZB info
    nzb_name        VARCHAR(500) NOT NULL,
    nzb_url         TEXT,
    nzb_file        BYTEA,                             -- Stored NZB content

    -- SABnzbd tracking
    sabnzbd_nzo_id  VARCHAR(100),                      -- SABnzbd job ID
    sabnzbd_status  VARCHAR(50) NOT NULL DEFAULT 'pending',

    -- Metadata (pre-download)
    expected_title  VARCHAR(500),
    expected_crew   UUID[],                            -- Expected performers
    expected_port   UUID REFERENCES qar.ports(id),     -- Expected studio

    -- Result
    treasure_id     UUID REFERENCES qar.treasures(id), -- Created gallery
    error_message   TEXT,

    -- Timestamps
    queued_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    started_at      TIMESTAMPTZ,
    completed_at    TIMESTAMPTZ,

    CONSTRAINT valid_plunder_status CHECK (sabnzbd_status IN (
        'pending', 'queued', 'downloading', 'extracting', 'processing',
        'completed', 'failed', 'cancelled'
    ))
);

-- Indexes
CREATE INDEX idx_treasures_fleet ON qar.treasures(fleet_id);
CREATE INDEX idx_treasures_voyage ON qar.treasures(voyage_id);
CREATE INDEX idx_treasures_port ON qar.treasures(port_id);
CREATE INDEX idx_treasures_release ON qar.treasures(release_date DESC);
CREATE INDEX idx_treasures_type ON qar.treasures(chest_type);
CREATE INDEX idx_doubloons_treasure ON qar.doubloons(treasure_id);
CREATE INDEX idx_doubloons_phash ON qar.doubloons(phash);
CREATE INDEX idx_treasure_crew_crew ON qar.treasure_crew(crew_id);
CREATE INDEX idx_plunder_status ON qar.plunder_queue(sabnzbd_status);
```

---

## Entity Relationships

```
                    ┌─────────────┐
                    │   Fleets    │
                    │ (Libraries) │
                    └──────┬──────┘
                           │
              ┌────────────┼────────────┐
              │            │            │
              ▼            ▼            ▼
        ┌──────────┐ ┌──────────┐ ┌──────────┐
        │ Voyages  │ │Treasures │ │   Crew   │
        │ (Scenes) │ │(Galleries│ │(Performers
        └────┬─────┘ └────┬─────┘ └────┬─────┘
             │            │            │
             │     ┌──────┴──────┐     │
             │     │             │     │
             └─────┤ treasure_   ├─────┘
                   │ crew        │
                   │ (M:M)       │
                   └─────────────┘
                         │
                         ▼
                   ┌───────────┐
                   │ Doubloons │
                   │ (Images)  │
                   └───────────┘
```

---

## Go Entities

```go
// internal/content/qar/gallery/entity.go
package gallery

import (
    "time"
    "github.com/google/uuid"
)

// ChestType represents gallery categories
type ChestType string

const (
    ChestTypeCrew     ChestType = "crew_chest"      // Performer gallery
    ChestTypeVoyage   ChestType = "voyage_chest"    // Scene BTS/promo
    ChestTypeMerchant ChestType = "merchant_haul"   // Standalone photoset
    ChestTypeRoyal    ChestType = "royal_treasury"  // Magazine shoot
    ChestTypePort     ChestType = "port_display"    // Studio promo
)

// Treasure represents a gallery (image collection)
type Treasure struct {
    ID            uuid.UUID  `json:"id"`
    FleetID       uuid.UUID  `json:"fleet_id"`
    Title         string     `json:"title"`
    Description   string     `json:"description,omitempty"`
    ChestType     ChestType  `json:"chest_type"`

    // Relationships
    VoyageID      *uuid.UUID `json:"voyage_id,omitempty"`      // Linked scene
    ExpeditionID  *uuid.UUID `json:"expedition_id,omitempty"`  // Linked movie
    PortID        *uuid.UUID `json:"port_id,omitempty"`        // Studio

    // Source
    SourceURL     string     `json:"source_url,omitempty"`
    SourceType    string     `json:"source_type,omitempty"`

    // Stats
    DoubloonCount int        `json:"doubloon_count"`
    TotalSize     int64      `json:"total_size"`

    // Storage
    Path          string     `json:"path"`
    CoverPath     string     `json:"cover_path,omitempty"`

    // External IDs
    StashDBID     string     `json:"charter,omitempty"`        // QAR: charter = stashdb_id

    // Dates
    ReleaseDate   *time.Time `json:"launch_date,omitempty"`    // QAR: launch_date
    ImportedAt    time.Time  `json:"imported_at"`
    CreatedAt     time.Time  `json:"created_at"`
    UpdatedAt     time.Time  `json:"updated_at"`

    // Populated on fetch
    Crew          []CrewRef  `json:"crew,omitempty"`
    Flags         []FlagRef  `json:"flags,omitempty"`
}

// Doubloon represents a single image in a gallery
type Doubloon struct {
    ID           uuid.UUID `json:"id"`
    TreasureID   uuid.UUID `json:"treasure_id"`
    Filename     string    `json:"filename"`
    Path         string    `json:"path"`
    Position     int       `json:"position"`

    // Dimensions
    Width        int       `json:"width,omitempty"`
    Height       int       `json:"height,omitempty"`
    SizeBytes    int64     `json:"size_bytes,omitempty"`
    Format       string    `json:"format,omitempty"`

    // Fingerprints
    PHash        string    `json:"coordinates,omitempty"`     // QAR: coordinates = phash
    Blurhash     string    `json:"blurhash,omitempty"`

    // Processed paths
    ThumbnailPath string   `json:"thumbnail_path,omitempty"`
    WebPPath      string   `json:"webp_path,omitempty"`

    CreatedAt    time.Time `json:"created_at"`
}

// PlunderStatus represents SABnzbd download status
type PlunderStatus string

const (
    PlunderPending     PlunderStatus = "pending"
    PlunderQueued      PlunderStatus = "queued"
    PlunderDownloading PlunderStatus = "downloading"
    PlunderExtracting  PlunderStatus = "extracting"
    PlunderProcessing  PlunderStatus = "processing"
    PlunderCompleted   PlunderStatus = "completed"
    PlunderFailed      PlunderStatus = "failed"
    PlunderCancelled   PlunderStatus = "cancelled"
)

// PlunderJob represents a download queue item
type PlunderJob struct {
    ID             uuid.UUID      `json:"id"`
    FleetID        uuid.UUID      `json:"fleet_id"`
    NZBName        string         `json:"nzb_name"`
    NZBURL         string         `json:"nzb_url,omitempty"`
    SABnzbdNZOID   string         `json:"sabnzbd_nzo_id,omitempty"`
    Status         PlunderStatus  `json:"status"`

    // Expected metadata
    ExpectedTitle  string         `json:"expected_title,omitempty"`
    ExpectedCrew   []uuid.UUID    `json:"expected_crew,omitempty"`
    ExpectedPortID *uuid.UUID     `json:"expected_port_id,omitempty"`

    // Result
    TreasureID     *uuid.UUID     `json:"treasure_id,omitempty"`
    ErrorMessage   string         `json:"error_message,omitempty"`

    // Timestamps
    QueuedAt       time.Time      `json:"queued_at"`
    StartedAt      *time.Time     `json:"started_at,omitempty"`
    CompletedAt    *time.Time     `json:"completed_at,omitempty"`
}

// CrewRef is a lightweight performer reference
type CrewRef struct {
    ID        uuid.UUID `json:"id"`
    Names     []string  `json:"names"`
    IsPrimary bool      `json:"is_primary"`
}

// FlagRef is a lightweight tag reference
type FlagRef struct {
    ID   uuid.UUID `json:"id"`
    Name string    `json:"name"`
}
```

---

## Prowlarr Integration (Download Search)

Gallery downloads use Prowlarr as the indexer/search aggregator. Prowlarr searches across configured indexers and sends results to download clients.

### Configuration

```yaml
# config.yaml
qar:
  plunder:
    enabled: true

    # Prowlarr for indexer search
    prowlarr:
      url: "http://prowlarr:9696"
      api_key: "${PROWLARR_API_KEY}"
      categories: [6000, 6010]           # Adult categories (Usenet/Torrent)

    # Post-processing
    download_path: "/downloads/qar"      # Completed downloads path
    auto_import: true                    # Auto-import completed downloads
    extract_archives: true               # Extract zip/rar automatically
    delete_archives: true                # Delete archives after extraction

    # Matching
    auto_match_crew: true                # Try to match performers from filename
    auto_match_port: true                # Try to match studio from filename
```

### Prowlarr Client

```go
// internal/service/metadata/prowlarr/client.go
package prowlarr

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "net/url"
)

type Client struct {
    baseURL string
    apiKey  string
    client  *http.Client
}

type SearchResult struct {
    GUID         string   `json:"guid"`
    Title        string   `json:"title"`
    Indexer      string   `json:"indexer"`
    IndexerID    int      `json:"indexerId"`
    Size         int64    `json:"size"`
    PublishDate  string   `json:"publishDate"`
    DownloadURL  string   `json:"downloadUrl"`
    InfoURL      string   `json:"infoUrl"`
    Categories   []int    `json:"categories"`
    Seeders      int      `json:"seeders,omitempty"`
    Leechers     int      `json:"leechers,omitempty"`
}

// Search queries Prowlarr indexers
func (c *Client) Search(ctx context.Context, query string, categories []int) ([]SearchResult, error) {
    params := url.Values{
        "query": {query},
    }
    for _, cat := range categories {
        params.Add("categories", fmt.Sprintf("%d", cat))
    }

    req, err := http.NewRequestWithContext(ctx, "GET",
        fmt.Sprintf("%s/api/v1/search?%s", c.baseURL, params.Encode()), nil)
    if err != nil {
        return nil, err
    }
    req.Header.Set("X-Api-Key", c.apiKey)

    resp, err := c.client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var results []SearchResult
    if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
        return nil, err
    }

    return results, nil
}

// SendToDownloadClient sends a release to the configured download client
func (c *Client) SendToDownloadClient(ctx context.Context, guid string, indexerID int) error {
    body := map[string]interface{}{
        "guid":      guid,
        "indexerId": indexerID,
    }

    // POST to /api/v1/search
    // Prowlarr handles routing to appropriate download client
    return nil
}
```

### River Jobs for Download Processing

```go
// internal/content/qar/gallery/jobs.go
package gallery

import (
    "context"
    "log/slog"
    "time"

    "github.com/google/uuid"
    "github.com/jackc/pgx/v5"
    "github.com/riverqueue/river"
)

// Job kinds
const (
    JobKindPlunderAdd       = "qar.plunder.add"
    JobKindPlunderCheck     = "qar.plunder.check"
    JobKindPlunderProcess   = "qar.plunder.process"
    JobKindTreasureScan     = "qar.treasure.scan"
    JobKindDoubloonProcess  = "qar.doubloon.process"
    JobKindDoubloonDedupe   = "qar.doubloon.dedupe"
)

// PlunderAddArgs - Add NZB to SABnzbd
type PlunderAddArgs struct {
    JobID       uuid.UUID  `json:"job_id"`
    NZBURL      string     `json:"nzb_url,omitempty"`
    NZBContent  []byte     `json:"nzb_content,omitempty"`
    Name        string     `json:"name"`
    FleetID     uuid.UUID  `json:"fleet_id"`
    ExpectedCrew []uuid.UUID `json:"expected_crew,omitempty"`
}

func (PlunderAddArgs) Kind() string { return JobKindPlunderAdd }

func (PlunderAddArgs) InsertOpts() river.InsertOpts {
    return river.InsertOpts{
        Queue:       "downloads",
        MaxAttempts: 3,
    }
}

// PlunderCheckArgs - Check download status
type PlunderCheckArgs struct {
    JobID uuid.UUID `json:"job_id"`
}

func (PlunderCheckArgs) Kind() string { return JobKindPlunderCheck }

func (PlunderCheckArgs) InsertOpts() river.InsertOpts {
    return river.InsertOpts{
        Queue:       "downloads",
        MaxAttempts: 50,  // Keep checking
    }
}

// PlunderProcessArgs - Process completed download
type PlunderProcessArgs struct {
    JobID        uuid.UUID `json:"job_id"`
    DownloadPath string    `json:"download_path"`
}

func (PlunderProcessArgs) Kind() string { return JobKindPlunderProcess }

// TreasureScanArgs - Scan folder for images
type TreasureScanArgs struct {
    TreasureID uuid.UUID `json:"treasure_id"`
    Path       string    `json:"path"`
}

func (TreasureScanArgs) Kind() string { return JobKindTreasureScan }

// DoubloonProcessArgs - Process single image (resize, webp, hash)
type DoubloonProcessArgs struct {
    DoubloonID uuid.UUID `json:"doubloon_id"`
}

func (DoubloonProcessArgs) Kind() string { return JobKindDoubloonProcess }

// DoubloonDedupeArgs - Find and remove duplicates
type DoubloonDedupeArgs struct {
    FleetID uuid.UUID `json:"fleet_id"`
}

func (DoubloonDedupeArgs) Kind() string { return JobKindDoubloonDedupe }
```

---

## Search Integration (Typesense)

```go
// internal/content/qar/gallery/search.go
package gallery

// TreasureSearchDocument for Typesense indexing
type TreasureSearchDocument struct {
    ID            string   `json:"id"`
    Title         string   `json:"title"`
    Description   string   `json:"description"`
    ChestType     string   `json:"chest_type"`
    CrewNames     []string `json:"crew_names"`      // Performer names for search
    PortName      string   `json:"port_name"`       // Studio name
    FlagNames     []string `json:"flag_names"`      // Tag names
    ReleaseDate   int64    `json:"release_date"`    // Unix timestamp
    DoubloonCount int      `json:"doubloon_count"`
    FleetID       string   `json:"fleet_id"`
}

// Typesense collection schema
var TreasureCollectionSchema = map[string]interface{}{
    "name": "qar_treasures",
    "fields": []map[string]interface{}{
        {"name": "id", "type": "string"},
        {"name": "title", "type": "string"},
        {"name": "description", "type": "string", "optional": true},
        {"name": "chest_type", "type": "string", "facet": true},
        {"name": "crew_names", "type": "string[]", "facet": true},
        {"name": "port_name", "type": "string", "facet": true, "optional": true},
        {"name": "flag_names", "type": "string[]", "facet": true},
        {"name": "release_date", "type": "int64", "sort": true, "optional": true},
        {"name": "doubloon_count", "type": "int32", "sort": true},
        {"name": "fleet_id", "type": "string", "facet": true},
    },
    "default_sorting_field": "release_date",
}
```

---

## API Endpoints

> See [00_SOURCE_OF_TRUTH.md](../../00_SOURCE_OF_TRUTH.md#api-namespaces) for API namespace conventions.

All QAR endpoints use `/api/v1/legacy/` namespace externally (obfuscated).

```
# Treasures (Galleries)
GET    /api/v1/legacy/treasures                    # List galleries
POST   /api/v1/legacy/treasures                    # Create gallery
GET    /api/v1/legacy/treasures/{id}               # Get gallery
PUT    /api/v1/legacy/treasures/{id}               # Update gallery
DELETE /api/v1/legacy/treasures/{id}               # Delete gallery

# Doubloons (Images)
GET    /api/v1/legacy/treasures/{id}/doubloons     # List gallery images
POST   /api/v1/legacy/treasures/{id}/doubloons     # Add images
GET    /api/v1/legacy/doubloons/{id}               # Get single image
DELETE /api/v1/legacy/doubloons/{id}               # Delete image

# Relationships
GET    /api/v1/legacy/crew/{id}/treasures          # Galleries for performer
GET    /api/v1/legacy/voyages/{id}/treasures       # Galleries for scene
GET    /api/v1/legacy/ports/{id}/treasures         # Galleries for studio

# Plunder (Downloads)
GET    /api/v1/legacy/plunder                      # List download queue
POST   /api/v1/legacy/plunder                      # Add to download queue
GET    /api/v1/legacy/plunder/{id}                 # Get download status
DELETE /api/v1/legacy/plunder/{id}                 # Cancel download

# Search
GET    /api/v1/legacy/treasures/search?q=...       # Search galleries

# User actions
POST   /api/v1/legacy/treasures/{id}/favorite      # Add to favorites
DELETE /api/v1/legacy/treasures/{id}/favorite      # Remove from favorites
POST   /api/v1/legacy/treasures/{id}/rate          # Rate gallery
```

---

## Implementation Checklist

### Phase 1: Schema & Repository
- [ ] Create migration `000001_qar_gallery.up.sql`
- [ ] Create sqlc queries `queries/qar/gallery.sql`
- [ ] Implement `gallery/repository.go`
- [ ] Implement `gallery/repository_pg.go`

### Phase 2: Core Service
- [ ] Implement `gallery/service.go` with caching
- [ ] Implement `gallery/entity.go`
- [ ] Implement crew/voyage linking

### Phase 3: SABnzbd Integration
- [ ] Implement `download/sabnzbd/client.go`
- [ ] Implement `gallery/plunder_service.go`
- [ ] Create River jobs for download pipeline

### Phase 4: Image Processing
- [ ] Implement pHash generation for deduplication
- [ ] Implement thumbnail generation
- [ ] Implement WebP conversion
- [ ] Implement Blurhash generation

### Phase 5: Search
- [ ] Create Typesense collection schema
- [ ] Implement search indexing on create/update
- [ ] Implement search service

### Phase 6: API
- [ ] Create OpenAPI spec `api/openapi/qar_gallery.yaml`
- [ ] Implement API handlers
- [ ] Add to router with QAR auth middleware

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
| [PostgreSQL Arrays](https://www.postgresql.org/docs/current/arrays.html) | [Local](../../../sources/database/postgresql-arrays.md) |
| [PostgreSQL JSON Functions](https://www.postgresql.org/docs/current/functions-json.html) | [Local](../../../sources/database/postgresql-json.md) |
| [River Job Queue](https://pkg.go.dev/github.com/riverqueue/river) | [Local](../../../sources/tooling/river.md) |
| [Servarr Wiki](https://wiki.servarr.com/) | [Local](../../../sources/apis/servarr-wiki.md) |
| [Typesense API](https://typesense.org/docs/latest/api/) | [Local](../../../sources/infrastructure/typesense.md) |
| [Typesense Go Client](https://github.com/typesense/typesense-go) | [Local](../../../sources/infrastructure/typesense-go.md) |
| [go-blurhash](https://pkg.go.dev/github.com/bbrks/go-blurhash) | [Local](../../../sources/media/go-blurhash.md) |
| [google/uuid](https://pkg.go.dev/github.com/google/uuid) | [Local](../../../sources/tooling/uuid.md) |
| [pgx PostgreSQL Driver](https://pkg.go.dev/github.com/jackc/pgx/v5) | [Local](../../../sources/database/pgx.md) |
| [sqlc](https://docs.sqlc.dev/en/stable/) | [Local](../../../sources/database/sqlc.md) |
| [sqlc Configuration](https://docs.sqlc.dev/en/stable/reference/config.html) | [Local](../../../sources/database/sqlc-config.md) |

<!-- SOURCE-BREADCRUMBS-END -->

<!-- DESIGN-BREADCRUMBS-START -->

## Related Design Docs

> Auto-generated cross-references to related design documentation

**Category**: [Adult](INDEX.md)

### In This Section

- [Revenge - Adult Content System](ADULT_CONTENT_SYSTEM.md)
- [Revenge - Adult Content Metadata System](ADULT_METADATA.md)
- [Adult Data Reconciliation](DATA_RECONCILIATION.md)
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

## Related Documentation

- [Adult Content System](ADULT_CONTENT_SYSTEM.md) - QAR schema isolation
- [Whisparr/StashDB Schema](WHISPARR_STASHDB_SCHEMA.md) - Scene schema
- [Adult Metadata](ADULT_METADATA.md) - Metadata sources
- [SABnzbd Integration](../../integrations/download/SABNZBD.md) - Download client (pending)
