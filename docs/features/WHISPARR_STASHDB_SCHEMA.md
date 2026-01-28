# Whisparr v3 & StashDB Schema Integration

> Custom UI/UX approach for adult content scenes using Whisparr cache

**Status**: 🟡 Research Complete, Schema Design Pending
**Last Updated**: 2026-01-28
**Dependencies**: Whisparr v3 ("eros" branch), StashDB API, PostgreSQL schema `c`

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

### PostgreSQL Schema: `c` (Adult Content Isolation)

All adult content tables use isolated schema `c` (see [Adult Content System](ADULT_CONTENT_SYSTEM.md)).

```sql
CREATE SCHEMA IF NOT EXISTS c;

-- Studios/Networks/Sites (Whisparr "Series")
CREATE TABLE c.studios (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    stashdb_id VARCHAR(50) UNIQUE,           -- StashDB studio ID
    whisparr_series_id INT UNIQUE,           -- Whisparr "series" ID
    name VARCHAR(500) NOT NULL,
    url VARCHAR(1000),                       -- Official site URL
    network_id UUID REFERENCES c.networks(id), -- Parent network (optional)
    logo_path VARCHAR(1000),
    metadata_json JSONB,                     -- Flexible metadata storage
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Networks (e.g., MindGeek, Aylo)
CREATE TABLE c.networks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    stashdb_id VARCHAR(50) UNIQUE,
    name VARCHAR(500) NOT NULL,
    description TEXT,
    logo_path VARCHAR(1000),
    metadata_json JSONB,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Scenes (Whisparr "Episodes")
CREATE TABLE c.scenes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    stashdb_id VARCHAR(50) UNIQUE,           -- StashDB scene ID
    whisparr_episode_id INT UNIQUE,          -- Whisparr "episode" ID
    studio_id UUID REFERENCES c.studios(id), -- Production studio
    title VARCHAR(500) NOT NULL,
    release_date DATE,
    release_year INT,
    duration_seconds INT,
    description TEXT,
    file_path VARCHAR(1000),                 -- Actual video file path
    metadata_json JSONB,                     -- Flexible metadata
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Performers (actors/actresses)
CREATE TABLE c.performers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    stashdb_id VARCHAR(50) UNIQUE,           -- StashDB performer ID
    name VARCHAR(500) NOT NULL,
    aliases TEXT[],                          -- Alternative names
    birthdate DATE,
    gender VARCHAR(50),                      -- Male, Female, Non-binary, etc.
    ethnicity VARCHAR(100),
    country VARCHAR(100),
    image_path VARCHAR(1000),
    metadata_json JSONB,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Scene-Performer relationships (many-to-many)
CREATE TABLE c.scene_performers (
    scene_id UUID REFERENCES c.scenes(id) ON DELETE CASCADE,
    performer_id UUID REFERENCES c.performers(id) ON DELETE CASCADE,
    role VARCHAR(100),                       -- Optional role (e.g., "lead", "supporting")
    PRIMARY KEY (scene_id, performer_id)
);

-- Tags (genres, categories)
CREATE TABLE c.tags (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    stashdb_id VARCHAR(50) UNIQUE,
    name VARCHAR(200) NOT NULL UNIQUE,
    category VARCHAR(100),                   -- Genre, Action, Position, etc.
    description TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Scene-Tag relationships (many-to-many)
CREATE TABLE c.scene_tags (
    scene_id UUID REFERENCES c.scenes(id) ON DELETE CASCADE,
    tag_id UUID REFERENCES c.tags(id) ON DELETE CASCADE,
    PRIMARY KEY (scene_id, tag_id)
);

-- User data (ratings, watch history, favorites)
CREATE TABLE c.scene_user_ratings (
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    scene_id UUID REFERENCES c.scenes(id) ON DELETE CASCADE,
    rating DECIMAL(2,1) CHECK (rating >= 0 AND rating <= 10),
    rated_at TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (user_id, scene_id)
);

CREATE TABLE c.scene_watch_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    scene_id UUID REFERENCES c.scenes(id) ON DELETE CASCADE,
    watched_at TIMESTAMPTZ DEFAULT NOW(),
    progress_seconds INT,                    -- Playback position
    completed BOOLEAN DEFAULT FALSE
);

CREATE TABLE c.scene_favorites (
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    scene_id UUID REFERENCES c.scenes(id) ON DELETE CASCADE,
    favorited_at TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (user_id, scene_id)
);

-- Indexes for performance
CREATE INDEX idx_scenes_studio_id ON c.scenes(studio_id);
CREATE INDEX idx_scenes_release_date ON c.scenes(release_date DESC);
CREATE INDEX idx_scene_performers_performer_id ON c.scene_performers(performer_id);
CREATE INDEX idx_scene_tags_tag_id ON c.scene_tags(tag_id);
CREATE INDEX idx_scene_watch_history_user_id ON c.scene_watch_history(user_id);
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
/c/scenes               # All scenes (grid view)
/c/scenes/{id}          # Scene detail page
/c/studios              # All studios
/c/studios/{id}         # Studio detail + scenes
/c/networks             # All networks
/c/networks/{id}        # Network detail + studios + scenes
/c/performers           # All performers
/c/performers/{id}      # Performer detail + scenes
/c/tags                 # All tags
/c/tags/{id}            # Tag detail + scenes
/c/search?q=...         # Search results
```

### NSFW Toggle Integration

- **Global Toggle**: User setting (`user_preferences.nsfw_enabled`)
- **Default**: OFF (explicit opt-in required)
- **Behavior when OFF**:
  - `/c/*` routes return 404
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
4. **Store in Revenge**: Save to PostgreSQL schema `c`

### StashDB Schema Mapping
| StashDB Entity | Revenge Table | ID Column |
|----------------|---------------|-----------|
| Scene | `c.scenes` | `stashdb_id` |
| Performer | `c.performers` | `stashdb_id` |
| Studio | `c.studios` | `stashdb_id` |
| Tag | `c.tags` | `stashdb_id` |

---

## Implementation Phases

### Phase 1: Whisparr Cache Analysis (Week 1)
- [ ] Install Whisparr v3 locally
- [ ] Locate cache database
- [ ] Document schema (scenes, performers, studios, tags, relationships)
- [ ] Extract sample data for testing

### Phase 2: Revenge Schema Creation (Week 1)
- [ ] Create PostgreSQL schema `c`
- [ ] Implement tables: `studios`, `networks`, `scenes`, `performers`, `scene_performers`, `tags`, `scene_tags`
- [ ] Implement user data tables: `scene_user_ratings`, `scene_watch_history`, `scene_favorites`
- [ ] Create indexes for performance

### Phase 3: Whisparr Import Service (Week 2)
- [ ] Build Go service to read Whisparr cache
- [ ] Extract scenes, performers, studios, tags
- [ ] Import into Revenge schema `c`
- [ ] Handle duplicates (upsert by `whisparr_episode_id`)

### Phase 4: StashDB Enrichment (Week 2)
- [ ] Integrate StashDB GraphQL API
- [ ] Match Revenge scenes with StashDB (by ID OR fuzzy matching)
- [ ] Download metadata (descriptions, images, aliases)
- [ ] Update Revenge tables with enriched data

### Phase 5: API Endpoints (Week 3)
- [ ] `/api/v1/c/scenes` (GET, POST, PUT, DELETE)
- [ ] `/api/v1/c/performers` (GET, POST, PUT, DELETE)
- [ ] `/api/v1/c/studios` (GET, POST, PUT, DELETE)
- [ ] `/api/v1/c/tags` (GET)
- [ ] User data endpoints (ratings, watch history, favorites)

### Phase 6: Frontend UI/UX (Week 4-5)
- [ ] Grid view component (Svelte 5 + shadcn-svelte)
- [ ] Studio/Network drill-down
- [ ] Performer browser
- [ ] Tag/Genre browser
- [ ] Search functionality
- [ ] Scene detail page
- [ ] NSFW toggle integration

---

## Related Documentation

- [Adult Content System](ADULT_CONTENT_SYSTEM.md) - Schema `c` isolation
- [Adult Metadata](ADULT_METADATA.md) - StashDB/ThePornDB integration
- [NSFW Toggle](NSFW_TOGGLE.md) - User preference component (pending creation)
- [Whisparr Integration](../integrations/WHISPARR.md) - API documentation (pending creation)
- [StashDB Integration](../integrations/STASHDB.md) - GraphQL API (pending creation)

---

## Open Questions

1. **Whisparr cache format**: SQLite? JSON? Custom binary?
2. **Sync frequency**: Real-time webhooks OR periodic polling (hourly/daily)?
3. **Folder structure**: How does Whisparr organize scene files on disk?
4. **StashDB rate limits**: Are there API rate limits we need to respect?
5. **Performer images**: StashDB-hosted OR external URLs?
6. **Tag ontology**: StashDB tags OR custom Revenge taxonomy?

**Next Steps**: Install Whisparr v3, analyze cache, answer open questions.
