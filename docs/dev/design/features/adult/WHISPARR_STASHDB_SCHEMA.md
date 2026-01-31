# Whisparr v3 & StashDB Schema Integration

> Custom UI/UX approach for adult content scenes using Whisparr cache

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

**Last Updated**: 2026-01-28
**Dependencies**: Whisparr v3 ("eros" branch), StashDB API, PostgreSQL schema `qar`

---

## Developer Resources

> Package versions: [00_SOURCE_OF_TRUTH.md](../../00_SOURCE_OF_TRUTH.md#go-dependencies-core)

| Source | URL | Purpose |
|--------|-----|---------|
| StashDB GraphQL | [docs.stashapp.cc](https://docs.stashapp.cc/) | Adult metadata API |
| Whisparr API | [whisparr.com](https://whisparr.com/) | Scene management |

| Package | Purpose |
|---------|---------|
| genqlient | Type-safe GraphQL client |
| goimagehash | pHash fingerprinting |
| pgx | PostgreSQL driver |

---

## Executive Summary

**Problem**: Whisparr v3 uses Sonarr codebase for scenes but the folder structure is NOT like TV shows (no Network/Show/Season/Episode hierarchy). All metadata is stored in Whisparr cache.

**Solution**: Custom UI/UX approach that handles scene-based content without forcing it into a TV show paradigm. Revenge will design its own schema and presentation optimized for adult content.

**Key Constraints**:
- Folder structure: NOT `Network/Show/Season/Episode`
- Data source: Whisparr cache (not direct API calls)
- Codebase: Sonarr-based (adapted for scenes, not episodes)
- UI/UX: Complete freedom to design custom approach

---

## Whisparr v3 Context

### Branch Information
- **Previous branch name**: `eros` (may have changed)
- **Codebase**: Uses Sonarr code for scene management
- **API Version**: v3 (endpoints: `/api/v3/series`, `/api/v3/episode`, `/api/v3/episodefile`)
- **OpenAPI**: `https://raw.githubusercontent.com/Whisparr/Whisparr/develop/src/Whisparr.Api.V3/openapi.json`

### API Endpoints (Sonarr-based)
```bash
# Series = Studios/Networks/Sites (NOT TV shows)
GET /api/v3/series
GET /api/v3/series/{id}
POST /api/v3/series
PUT /api/v3/series/{id}
DELETE /api/v3/series/{id}

# Episode = Scenes (NOT TV episodes)
GET /api/v3/episode
GET /api/v3/episode/{id}
PUT /api/v3/episode/{id}

# Episode File = Scene Files
GET /api/v3/episodefile
GET /api/v3/episodefile/{id}
DELETE /api/v3/episodefile/{id}

# Other
GET /api/v3/calendar
GET /api/v3/history
GET /api/v3/queue
GET /api/v3/tag
GET /api/v3/autotagging
```

### Key Differences from TV Shows
| Aspect | TV Shows (Sonarr) | Adult Scenes (Whisparr v3) |
|--------|-------------------|----------------------------|
| **Hierarchy** | Network → Show → Season → Episode | Studio/Network → Site/Series → Scene |
| **Folder Structure** | `Network/ShowName/Season XX/EpisodeFile` | **NOT** TV-like (custom Whisparr structure) |
| **Data Source** | Live API calls + local cache | **Whisparr cache only** (all metadata stored locally) |
| **Metadata** | TheTVDB, TMDb | StashDB, ThePornDB |
| **Relationships** | Actors play characters in episodes | Performers appear in scenes |

---

## Whisparr Cache Structure

**Status**: 🔴 NOT YET ANALYZED

Whisparr stores all metadata in its local cache (likely SQLite or JSON). Revenge needs to:

1. **Locate Whisparr cache** (e.g., `~/.config/Whisparr/whisparr.db`)
2. **Analyze schema**: Tables/collections for scenes, performers, studios, tags
3. **Extract metadata**: Scene titles, release dates, performer relationships, tags
4. **Sync strategy**: Periodic polling OR webhook triggers OR one-time import

**Action Items**:
- [ ] Install Whisparr v3 locally
- [ ] Inspect cache database schema
- [ ] Document scene/performer/studio/tag structure
- [ ] Design Revenge import strategy

---

## Revenge Schema Design

### PostgreSQL Schema: `qar` (Adult Content Isolation - Queen Anne's Revenge)

All adult content tables use isolated schema `qar` (see [Adult Content System](ADULT_CONTENT_SYSTEM.md)).

```sql
CREATE SCHEMA IF NOT EXISTS qar;

-- Studios/Networks/Sites → Ports (Whisparr "Series")
CREATE TABLE qar.ports (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    charter VARCHAR(50) UNIQUE,              -- StashDB studio ID (obfuscated)
    whisparr_series_id INT UNIQUE,           -- Whisparr "series" ID
    name VARCHAR(500) NOT NULL,
    url VARCHAR(1000),                       -- Official site URL
    armada_id UUID REFERENCES qar.armadas(id), -- Parent network (optional)
    logo_path VARCHAR(1000),
    metadata_json JSONB,                     -- Flexible metadata storage
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Networks → Armadas (e.g., MindGeek, Aylo)
CREATE TABLE qar.armadas (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    charter VARCHAR(50) UNIQUE,              -- stashdb_id (obfuscated)
    name VARCHAR(500) NOT NULL,
    description TEXT,
    logo_path VARCHAR(1000),
    metadata_json JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Scenes → Voyages (Whisparr "Episodes")
CREATE TABLE qar.voyages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    charter VARCHAR(50) UNIQUE,              -- StashDB scene ID (obfuscated)
    whisparr_episode_id INT UNIQUE,          -- Whisparr "episode" ID
    port_id UUID REFERENCES qar.ports(id),   -- Production studio
    title VARCHAR(500) NOT NULL,
    launch_date DATE,                        -- release_date (obfuscated)
    release_year INT,
    distance INT,                            -- duration_seconds (obfuscated)
    description TEXT,
    file_path VARCHAR(1000),                 -- Actual video file path
    metadata_json JSONB,                     -- Flexible metadata
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Voyage Fingerprints (pHash, oshash, MD5 for identification)
CREATE TABLE qar.voyage_coordinates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    voyage_id UUID NOT NULL REFERENCES qar.voyages(id) ON DELETE CASCADE,
    algorithm VARCHAR(20) NOT NULL,          -- 'phash', 'oshash', 'md5'
    hash VARCHAR(64) NOT NULL,               -- Hash value
    distance INT,                            -- Duration at time of hash (for validation)
    source VARCHAR(50),                      -- 'stashdb', 'stashapp', 'local'
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(voyage_id, algorithm),
    UNIQUE(algorithm, hash)                  -- Same hash = same voyage
);

-- Index for fast fingerprint lookups
CREATE INDEX idx_qar_voyage_coordinates_hash ON qar.voyage_coordinates(algorithm, hash);

-- Performers → Crew (actors/actresses)
CREATE TABLE qar.crew (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    charter VARCHAR(50) UNIQUE,              -- StashDB performer ID (obfuscated)
    name VARCHAR(500) NOT NULL,
    names TEXT[],                            -- aliases (obfuscated)
    christening DATE,                        -- birthdate (obfuscated)
    gender VARCHAR(50),                      -- Male, Female, Non-binary, etc.
    origin VARCHAR(100),                     -- ethnicity (obfuscated)
    country VARCHAR(100),
    image_path VARCHAR(1000),
    metadata_json JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Voyage-Crew relationships (many-to-many)
CREATE TABLE qar.voyage_crew (
    voyage_id UUID REFERENCES qar.voyages(id) ON DELETE CASCADE,
    crew_id UUID REFERENCES qar.crew(id) ON DELETE CASCADE,
    role VARCHAR(100),                       -- Optional role (e.g., "lead", "supporting")
    PRIMARY KEY (voyage_id, crew_id)
);

-- Tags → Flags (genres, categories)
CREATE TABLE qar.flags (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    charter VARCHAR(50) UNIQUE,              -- stashdb_id (obfuscated)
    name VARCHAR(200) NOT NULL UNIQUE,
    waters VARCHAR(100),                     -- category (obfuscated) - Genre, Action, Position, etc.
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Voyage-Flag relationships (many-to-many)
CREATE TABLE qar.voyage_flags (
    voyage_id UUID REFERENCES qar.voyages(id) ON DELETE CASCADE,
    flag_id UUID REFERENCES qar.flags(id) ON DELETE CASCADE,
    PRIMARY KEY (voyage_id, flag_id)
);

-- User data (ratings, watch history, favorites)
CREATE TABLE qar.voyage_bounties (
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    voyage_id UUID REFERENCES qar.voyages(id) ON DELETE CASCADE,
    bounty DECIMAL(2,1) CHECK (bounty >= 0 AND bounty <= 10), -- rating (obfuscated)
    rated_at TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (user_id, voyage_id)
);

CREATE TABLE qar.ship_log (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    voyage_id UUID REFERENCES qar.voyages(id) ON DELETE CASCADE,
    logged_at TIMESTAMPTZ DEFAULT NOW(),     -- watched_at (obfuscated)
    progress_seconds INT,                    -- Playback position
    completed BOOLEAN DEFAULT FALSE
);

CREATE TABLE qar.voyage_favorites (
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    voyage_id UUID REFERENCES qar.voyages(id) ON DELETE CASCADE,
    favorited_at TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (user_id, voyage_id)
);

-- Indexes for performance
CREATE INDEX idx_qar_voyages_port ON qar.voyages(port_id);
CREATE INDEX idx_qar_voyages_launch_date ON qar.voyages(launch_date DESC);
CREATE INDEX idx_qar_voyage_crew_crew ON qar.voyage_crew(crew_id);
CREATE INDEX idx_qar_voyage_flags_flag ON qar.voyage_flags(flag_id);
CREATE INDEX idx_qar_ship_log_user ON qar.ship_log(user_id);
```

---

## UI/UX Approach (Custom Design)

### No TV Show Paradigm

**Don't force scenes into Season/Episode structure**. Instead:

1. **Grid View** (Primary): Scene thumbnails with metadata overlays
   - Thumbnail: Poster/screenshot from scene
   - Overlay: Title, studio logo, release date, duration, rating
   - Sorting: Release date (newest first), rating, studio, duration
   - Filtering: Studio, network, performer, tags, release year

2. **Studio/Network View**: Browse by production studio or network
   - Studio cards: Logo, name, scene count
   - Clicking studio → Grid view of all scenes from that studio
   - Network hierarchy: Network → Studios → Scenes (3-level drill-down)

3. **Performer View**: Browse by performer
   - Performer cards: Photo, name, scene count
   - Clicking performer → Grid view of all scenes featuring that performer
   - Performer details: Aliases, metadata, filmography timeline

4. **Tag/Genre View**: Browse by tags/genres
   - Tag cloud OR tag list with scene counts
   - Clicking tag → Grid view of all tagged scenes
   - Multi-tag filtering (AND/OR logic)

5. **Search**: Full-text search across scenes, performers, studios, tags
   - Search bar: Auto-complete with suggestions
   - Results: Scenes, performers, studios (separate sections)

6. **Scene Detail Page**:
   - Large thumbnail/poster
   - Title, release date, duration, rating
   - Studio/network badges (clickable)
   - Performers (photos + names, clickable)
   - Tags (clickable)
   - Description
   - File info (resolution, codec, file size)
   - User actions: Play, Rate, Favorite, Mark as watched

### URL Structure

```
/qar/voyages            # All voyages (scenes - grid view)
/qar/voyages/{id}       # Voyage detail page
/qar/ports              # All ports (studios)
/qar/ports/{id}         # Port detail + voyages
/qar/armadas            # All armadas (networks)
/qar/armadas/{id}       # Armada detail + ports + voyages
/qar/crew               # All crew (performers)
/qar/crew/{id}          # Crew detail + voyages
/qar/flags              # All flags (tags)
/qar/flags/{id}         # Flag detail + voyages
/qar/search?q=...       # Search results
```

### NSFW Toggle Integration

- **Global Toggle**: User setting (`user_preferences.nsfw_enabled`)
- **Default**: OFF (explicit opt-in required)
- **Behavior when OFF**:
  - `/qar/*` routes return 404
  - Adult content hidden from search
  - Sidebar hides adult module
  - Dashboard excludes adult activity

See [NSFW Toggle](NSFW_TOGGLE.md) for component details.

---

## StashDB Integration

### StashDB API
- **URL**: `https://stashdb.org/graphql`
- **Auth**: API key required
- **Data**: Performers, studios, scenes, tags

### Metadata Enrichment Flow
1. **Import from Whisparr cache**: Extract scene/performer/studio data
2. **Match with StashDB**: Use StashDB IDs OR fuzzy matching (name + date)
3. **Enrich metadata**: Download StashDB-provided metadata (descriptions, images, aliases)
4. **Store in Revenge**: Save to PostgreSQL schema `qar`

### StashDB Schema Mapping
| StashDB Entity | Revenge Table | ID Column |
|----------------|---------------|-----------|
| Scene | `qar.voyages` | `charter` |
| Performer | `qar.crew` | `charter` |
| Studio | `qar.ports` | `charter` |
| Tag | `qar.flags` | `charter` |

---

## Perceptual Hashing (Scene Identification)

### Hash Types

| Algorithm | Purpose | Source |
|-----------|---------|--------|
| **pHash** (Perceptual Hash) | Visual fingerprint, tolerant to re-encoding | StashDB, Stash-App, local generation |
| **oshash** | OpenSubtitles hash - fast file-based hash | Stash-App, local generation |
| **MD5** | File integrity check | Local generation |

### Fingerprint Fetch Flow

```
┌──────────────────────────────────────────────────────────────┐
│                  Scene Fingerprint Flow                       │
└──────────────────────────────────────────────────────────────┘

1. Library Scan
       │
       ▼
2. Check Stash-App (if configured)
   └─→ Fetch existing fingerprints via GraphQL
       │
       ▼
3. If no Stash-App OR missing hashes
   └─→ Query StashDB by title/studio/duration
       │
       ▼
4. If still missing
   └─→ Generate locally (pHash, oshash, MD5)
       │
       ▼
5. Store in c.scene_fingerprints
```

### Fingerprint Matching

```go
// FingerprintService handles scene identification via hashes
type FingerprintService struct {
    db          *pgxpool.Pool
    stashClient *StashAppClient
    stashDB     *StashDBClient
}

// IdentifyScene tries to identify a scene by its fingerprints
func (s *FingerprintService) IdentifyScene(ctx context.Context, filePath string) (*Scene, error) {
    // 1. Generate local fingerprints
    phash, _ := s.generatePHash(filePath)
    oshash, _ := s.generateOsHash(filePath)

    // 2. Check local database
    if scene, err := s.findByFingerprint(ctx, "phash", phash); err == nil {
        return scene, nil
    }

    // 3. Check Stash-App (if configured)
    if s.stashClient != nil {
        if match, err := s.stashClient.FindSceneByFingerprint(ctx, phash, oshash); err == nil {
            return s.importFromStashApp(ctx, match)
        }
    }

    // 4. Check StashDB
    if match, err := s.stashDB.FindSceneByFingerprint(ctx, phash); err == nil {
        return s.importFromStashDB(ctx, match)
    }

    return nil, ErrSceneNotIdentified
}
```

### pHash Generation

```go
// GeneratePHash creates a perceptual hash from video frames
func (s *FingerprintService) generatePHash(filePath string) (string, error) {
    // Extract keyframes using ffmpeg
    frames, err := extractKeyframes(filePath, 8) // 8 frames
    if err != nil {
        return "", err
    }

    // Generate pHash from combined frames
    // Uses same algorithm as StashDB for compatibility
    hash := goimagehash.PerceptionHash(frames)
    return hash.ToString(), nil
}
```

---

## Implementation Phases

### Phase 1: Whisparr Cache Analysis (Week 1)
- [ ] Install Whisparr v3 locally
- [ ] Locate cache database
- [ ] Document schema (scenes, performers, studios, tags, relationships)
- [ ] Extract sample data for testing

### Phase 2: Revenge Schema Creation (Week 1)
- [ ] Create PostgreSQL schema `qar`
- [ ] Implement tables: `ports`, `armadas`, `voyages`, `crew`, `voyage_crew`, `flags`, `voyage_flags`
- [ ] Implement user data tables: `voyage_bounties`, `ship_log`, `voyage_favorites`
- [ ] Create indexes for performance

### Phase 3: Whisparr Import Service (Week 2)
- [ ] Build Go service to read Whisparr cache
- [ ] Extract voyages, crew, ports, flags
- [ ] Import into Revenge schema `qar`
- [ ] Handle duplicates (upsert by `whisparr_episode_id`)

### Phase 4: StashDB Enrichment (Week 2)
- [ ] Integrate StashDB GraphQL API
- [ ] Match Revenge scenes with StashDB (by ID OR fuzzy matching)
- [ ] Download metadata (descriptions, images, aliases)
- [ ] Update Revenge tables with enriched data

### Phase 5: API Endpoints (Week 3)
- [ ] `/api/v1/qar/voyages` (GET, POST, PUT, DELETE)
- [ ] `/api/v1/qar/crew` (GET, POST, PUT, DELETE)
- [ ] `/api/v1/qar/ports` (GET, POST, PUT, DELETE)
- [ ] `/api/v1/qar/flags` (GET)
- [ ] User data endpoints (bounties, ship_log, favorites)

### Phase 6: Frontend UI/UX (Week 4-5)
- [ ] Grid view component (Svelte 5 + shadcn-svelte)
- [ ] Studio/Network drill-down
- [ ] Performer browser
- [ ] Tag/Genre browser
- [ ] Search functionality
- [ ] Scene detail page
- [ ] NSFW toggle integration

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
| [Whisparr OpenAPI Spec](https://raw.githubusercontent.com/Whisparr/Whisparr/develop/src/Whisparr.Api.V3/openapi.json) | [Local](../../../sources/apis/whisparr-openapi.json) |

<!-- SOURCE-BREADCRUMBS-END -->

<!-- DESIGN-BREADCRUMBS-START -->

## Related Design Docs

> Auto-generated cross-references to related design documentation

**Category**: [Adult](INDEX.md)

### In This Section

- [Revenge - Adult Content System](ADULT_CONTENT_SYSTEM.md)
- [Revenge - Adult Content Metadata System](ADULT_METADATA.md)
- [Adult Data Reconciliation](DATA_RECONCILIATION.md)
- [Adult Gallery Module (QAR: Treasures)](GALLERY_MODULE.md)

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

- [Adult Content System](ADULT_CONTENT_SYSTEM.md) - Schema `qar` isolation
- [Adult Metadata](ADULT_METADATA.md) - StashDB/ThePornDB integration
- [NSFW Toggle](NSFW_TOGGLE.md) - User preference component (pending creation)
- [Whisparr Integration](../integrations/WHISPARR.md) - API documentation (pending creation)
- [StashDB Integration](../integrations/STASHDB.md) - GraphQL API (pending creation)

---

## Design Decisions (Resolved 2026-01-30)

| Decision | Choice | Rationale |
|----------|--------|-----------|
| **Sync Method** | Periodic API Polling | Konsistent mit allen Arr-Services (Radarr, Sonarr, Lidarr) |
| **Tag Ontology** | Mapping | StashDB Tags → eigene Revenge/QAR Kategorien mappen |
| **Performer Images** | StashDB URLs | Keine lokale Kopie - direkt von StashDB referenzieren |
| **Folder Structure** | Whisparr Managed | Wir lesen nur was Whisparr API uns gibt |
| **Rate Limits** | Implement on demand | Bei Implementation rate limiting einbauen falls nötig |

### Sync Architecture

```
┌─────────────┐    Polling     ┌─────────────┐
│   Revenge   │ ────────────►  │  Whisparr   │
│   (River)   │   every 15min  │   API v3    │
└─────────────┘                └─────────────┘
       │
       │  On new content
       ▼
┌─────────────┐    GraphQL     ┌─────────────┐
│  Metadata   │ ────────────►  │   StashDB   │
│  Enrichment │                │             │
└─────────────┘                └─────────────┘
```

### Tag Mapping Strategy

StashDB Tags werden auf Revenge Flag-Kategorien (`waters`) gemappt:

| StashDB Tag Category | QAR Flag `waters` | Example Tags |
|---------------------|-------------------|--------------|
| Genre | `genre` | Anal, Oral, Group |
| Position | `position` | Cowgirl, Doggy, Missionary |
| Action | `action` | Creampie, Facial, Swallow |
| Attribute | `attribute` | Big Tits, Tattoo, Blonde |
| Scene Type | `scene_type` | POV, VR, Gonzo |

Unmapped tags werden als `waters = 'other'` importiert.
