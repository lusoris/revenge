# Stash Integration

> Self-hosted adult media organizer with GraphQL API

**Service**: Stash (https://github.com/stashapp/stash)
**API**: GraphQL API (localhost:9999/graphql)
**Category**: Self-Hosted Media Organizer (Adult Content)
**Priority**: 🟢 MEDIUM (Optional integration for Stash users)
**Status**: 🔴 DESIGN PHASE

---

## Overview

**Stash** is a self-hosted adult media organizer similar to Plex/Jellyfin but specifically for adult content. It provides library management, metadata scraping (via StashDB), scene markers, video streaming, and a web UI.

**Key Features**:
- **Library management**: Scan video files, organize by performer/studio/tags
- **Metadata scraping**: Automatic metadata via StashDB (fingerprinting + API)
- **Scene markers**: Time-based markers (positions, performers, acts)
- **Video streaming**: Web-based video player with HLS transcoding
- **Performer management**: Performer profiles with images, measurements, career info
- **Studio management**: Studio profiles with logos, parent companies
- **Tag system**: Comprehensive tagging (positions, acts, settings, fetishes)
- **GraphQL API**: Full API access for library data, metadata, playback sessions
- **Self-hosted**: Runs on user's server (no cloud dependencies)

**Use Cases**:
- Import existing Stash library into Revenge
- Sync Stash metadata → Revenge database
- Use Revenge as alternative UI for Stash library
- Migrate from Stash to Revenge (one-time import)

**⚠️ CRITICAL: Adult Content Isolation**:
- **Database schema**: `c` schema ONLY (`c.movies`, `c.scenes`, `c.performers`, `c.studios`)
- **API namespace**: `/api/v1/c/integrations/stash/*` (NOT `/api/v1/integrations/stash/*`)
- **Module location**: `internal/content/c/integrations/stash/` (NOT `internal/service/integrations/`)
- **Access control**: Mods/admins can see all data for monitoring, regular users see only their own library

---

## Developer Resources

### API Documentation
- **GraphQL Schema**: https://stash-app.github.io/graphql-api/
- **GraphQL Playground**: http://localhost:9999/playground (Stash instance required)
- **GitHub**: https://github.com/stashapp/stash
- **Docs**: https://docs.stashapp.cc/

### Authentication
- **Method**: API Key (header-based)
- **Header**: `ApiKey: YOUR_API_KEY`
- **API Key location**: Stash Settings → Security → API Key
- **Rate Limits**: None (self-hosted, no rate limiting)

### Data Coverage
- **Scenes**: Unlimited (depends on user's library)
- **Performers**: Unlimited (user-managed)
- **Studios**: Unlimited (user-managed)
- **Tags**: Unlimited (user-managed + StashDB tags)

### Go Client Library
- **Official**: None (use standard GraphQL client)
- **Recommended**: `github.com/machinebox/graphql` OR `github.com/Khan/genqlient`
- **Alternative**: `net/http` + manual GraphQL JSON requests

---

## API Details

### GraphQL Queries

#### List All Scenes
```graphql
query FindScenes($filter: FindFilterType, $scene_filter: SceneFilterType) {
  findScenes(filter: $filter, scene_filter: $scene_filter) {
    count
    scenes {
      id
      title
      details
      date
      rating100
      organized
      file {
        path
        size
        duration
        video_codec
        width
        height
        framerate
        bitrate
      }
      studio {
        id
        name
        url
      }
      performers {
        id
        name
        gender
        birthdate
        measurements
      }
      tags {
        id
        name
      }
      scene_markers {
        id
        title
        seconds
        primary_tag {
          id
          name
        }
      }
      paths {
        screenshot
        preview
        stream
      }
      stash_ids {
        endpoint
        stash_id
      }
    }
  }
}
```

#### Get Scene Details
```graphql
query FindScene($id: ID!) {
  findScene(id: $id) {
    id
    title
    details
    url
    date
    rating100
    organized
    file {
      path
      size
      duration
      video_codec
      audio_codec
      width
      height
      framerate
      bitrate
    }
    studio {
      id
      name
      url
      parent_studio {
        id
        name
      }
      stash_ids {
        endpoint
        stash_id
      }
    }
    performers {
      id
      name
      disambiguation
      gender
      birthdate
      death_date
      ethnicity
      country
      eye_color
      hair_color
      height_cm
      measurements
      tattoos
      piercings
      career_length
      image_path
      stash_ids {
        endpoint
        stash_id
      }
    }
    tags {
      id
      name
      description
    }
    scene_markers {
      id
      title
      seconds
      primary_tag {
        id
        name
      }
      tags {
        id
        name
      }
      screenshot
    }
    paths {
      screenshot
      preview
      stream
      webp
      vtt
      sprite
      funscript
    }
    stash_ids {
      endpoint
      stash_id
    }
  }
}
```

#### List All Performers
```graphql
query FindPerformers($filter: FindFilterType, $performer_filter: PerformerFilterType) {
  findPerformers(filter: $filter, performer_filter: $performer_filter) {
    count
    performers {
      id
      name
      disambiguation
      gender
      birthdate
      death_date
      ethnicity
      country
      eye_color
      hair_color
      height_cm
      measurements
      fake_tits
      tattoos
      piercings
      career_length
      aliases
      image_path
      scene_count
      stash_ids {
        endpoint
        stash_id
      }
    }
  }
}
```

#### Get Performer Details
```graphql
query FindPerformer($id: ID!) {
  findPerformer(id: $id) {
    id
    name
    disambiguation
    gender
    birthdate
    death_date
    ethnicity
    country
    eye_color
    hair_color
    height_cm
    weight
    measurements
    fake_tits
    tattoos
    piercings
    career_length
    aliases
    favorite
    image_path
    scene_count
    stash_ids {
      endpoint
      stash_id
    }
  }
}
```

#### List All Studios
```graphql
query FindStudios($filter: FindFilterType) {
  findStudios(filter: $filter) {
    count
    studios {
      id
      name
      url
      parent_studio {
        id
        name
      }
      child_studios {
        id
        name
      }
      image_path
      scene_count
      stash_ids {
        endpoint
        stash_id
      }
    }
  }
}
```

#### Get Configuration
```graphql
query Configuration {
  configuration {
    general {
      stashes {
        path
      }
      databasePath
      generatedPath
    }
    interface {
      language
    }
  }
}
```

---

## Implementation Checklist

### Phase 1: Core Integration (Optional)
- [ ] GraphQL client setup (`machinebox/graphql` OR `genqlient`)
- [ ] Stash instance configuration (`configs/config.yaml` - `stash.url`, `stash.api_key`)
- [ ] **Adult schema**: Use existing `c.movies`, `c.performers`, `c.studios` tables
- [ ] **API namespace**: `/api/v1/c/integrations/stash/*` endpoints
- [ ] **Module location**: `internal/content/c/integrations/stash/` (isolated)
- [ ] List Stash scenes (GraphQL `findScenes`)
- [ ] Import scene metadata (title, date, performers, studio, tags)
- [ ] Import scene files (copy OR symlink video files)
- [ ] Import performer data (name, measurements, images)
- [ ] Import studio data (name, logos, parent companies)
- [ ] Import tags (positions, acts, settings)
- [ ] Scene markers import (time-based markers)

### Phase 2: One-Time Migration (Optional)
- [ ] Full library migration (Stash → Revenge)
- [ ] File relocation (copy video files to Revenge library paths)
- [ ] Metadata sync (all scenes, performers, studios, tags)
- [ ] Image import (performer images, studio logos, scene screenshots)
- [ ] Watch history migration (play count, last played)
- [ ] Duplicate detection (avoid re-importing existing content)

### Phase 3: Ongoing Sync (Optional)
- [ ] **Job**: `c.integrations.stash.sync_library` (periodic sync)
- [ ] Incremental sync (only new/updated scenes since last sync)
- [ ] Two-way sync (Revenge edits → Stash database)
- [ ] Conflict resolution (prefer Revenge data OR prefer Stash data)

---

## Integration Pattern

### One-Time Library Import Flow
```
User enables Stash integration
        ↓
Configure Stash URL + API Key (configs/config.yaml)
        ↓
Fetch all scenes (GraphQL findScenes)
        ↓
For each scene:
        ↓
        Extract metadata: title, date, performers, studio, tags, file path
        ↓
        Check if scene already exists in Revenge (match by file path OR StashDB ID)
        ↓
        Scene exists? → Skip (avoid duplicates)
        ↓
        Scene NOT exists? → Import scene
                ↓
                Copy OR symlink video file (user preference)
                ↓
                Store in c.movies OR c.scenes
                ↓
                metadata_json.stash_data = full GraphQL response
                ↓
                Import performers (create in c.performers if not exists)
                ↓
                Download performer images
                ↓
                Import studio (create in c.studios if not exists)
                ↓
                Download studio logo
                ↓
                Import tags (create in c.tags if not exists)
                ↓
                Import scene markers (time-based markers)
                ↓
                Update Typesense search index
                ↓
                Notify user: "Imported {count} scenes from Stash"
```

### Incremental Sync Flow (Optional)
```
Scheduled job (daily/weekly)
        ↓
Fetch updated scenes (GraphQL findScenes with updated_at filter)
        ↓
For each updated scene:
        ↓
        Check if scene exists in Revenge
        ↓
        Scene exists? → Update metadata
                ↓
                Compare updated_at timestamps
                ↓
                Stash updated_at > Revenge updated_at? → Sync metadata
                ↓
                Update c.movies.metadata_json.stash_data
                ↓
                Update performers/studio/tags (if changed)
        ↓
        Scene NOT exists? → Import scene (same flow as one-time import)
```

---

## Related Documentation

- [STASHDB.md](./STASHDB.md) - StashDB metadata provider (Stash uses this)
- [THEPORNDB.md](./THEPORNDB.md) - Alternative metadata provider
- [WHISPARR.md](../../servarr/WHISPARR.md) - Adult content management
- [ADULT_METADATA.md](../../../ADULT_METADATA.md) - Adult metadata system architecture

---

## Notes

### Self-Hosted Requirement
- **Stash must be running**: User must have Stash instance running (self-hosted)
- **Local network**: Typically accessed via `http://localhost:9999` OR `http://192.168.x.x:9999`
- **API Key**: User must enable API Key in Stash settings (Security → API Key)

### Use Case: Migration vs Sync
- **One-time migration**: User wants to switch from Stash to Revenge (import library once)
- **Ongoing sync**: User wants to use both Stash + Revenge (keep libraries in sync)
- **Recommended**: One-time migration (Revenge replaces Stash)

### Adult Content Isolation (CRITICAL)
- **Database schema**: `c` schema ONLY
  - `c.movies.metadata_json.stash_data` (JSONB)
  - `c.scenes.metadata_json.stash_data` (JSONB)
  - `c.performers` (shared with StashDB/ThePornDB)
  - `c.studios` (shared with StashDB/ThePornDB)
- **API namespace**: `/api/v1/c/integrations/stash/*` (isolated)
  - `/api/v1/c/integrations/stash/sync` (trigger sync)
  - `/api/v1/c/integrations/stash/import` (one-time import)
  - `/api/v1/c/integrations/stash/status` (sync status)
- **Module location**: `internal/content/c/integrations/stash/` (isolated)
- **Access control**: Mods/admins see all, regular users see only their library

### File Handling
- **Copy files**: Copy video files from Stash library to Revenge library (requires disk space)
- **Symlink files**: Create symlinks to Stash library (saves disk space, but depends on Stash)
- **User preference**: Allow user to choose copy OR symlink

### StashDB IDs
- **Stash stores StashDB IDs**: `stash_ids` field contains StashDB identifiers
- **Use for deduplication**: Match scenes by StashDB ID (avoid re-importing)
- **Sync with StashDB**: If scene has StashDB ID → fetch fresh metadata from StashDB (prefer StashDB over Stash metadata)

### Scene Markers
- **Time-based markers**: Stash supports scene markers (positions, performers, acts at specific timestamps)
- **Import markers**: Store markers in Revenge database (useful for scene navigation)
- **Schema**: `c.scene_markers` table (scene_id, timestamp, tag_id, title)

### Performer Images
- **Stash stores performer images**: `image_path` field (relative to Stash `generatedPath`)
- **Download images**: Fetch performer images from Stash API (serve via Stash web server)
- **Store locally**: Download and store in Revenge media storage

### GraphQL vs REST
- **Stash uses GraphQL**: Same as StashDB (GraphQL API)
- **Flexible queries**: Request only needed fields
- **Batch queries**: Fetch multiple scenes in one query (pagination)

### Rate Limiting
- **No rate limits**: Stash is self-hosted (no rate limiting)
- **Batch processing**: Use batches to avoid overwhelming server (e.g., import 100 scenes at a time)

### Watch History (Optional)
- **Stash tracks watch history**: `o_counter` (play count), `last_played_at`
- **Import watch history**: Optionally import play count/last played into Revenge
- **Schema**: `c.watch_history` table (user_id, scene_id, play_count, last_played_at)

### Two-Way Sync (Advanced)
- **Revenge → Stash**: Optionally sync Revenge edits back to Stash database
- **Use case**: User edits metadata in Revenge → update Stash database
- **Complexity**: High (requires Stash database access OR GraphQL mutations)
- **Recommendation**: One-way sync (Stash → Revenge only)

### Conflict Resolution
- **Timestamp comparison**: Compare `updated_at` timestamps (prefer newer)
- **User preference**: Allow user to choose preferred source (Revenge OR Stash)
- **Default**: Prefer Revenge data (user edited in Revenge = authoritative)

### JSONB Storage
- Store full Stash GraphQL response in `c.movies.metadata_json.stash_data`
- Preserves all Stash-specific fields (scene markers, ratings, etc.)
- Allows querying Stash-specific data via PostgreSQL JSONB operators

### Use Case: Optional Integration
- **Priority**: Medium/Low (most users won't have Stash)
- **Target audience**: Existing Stash users migrating to Revenge
- **Implementation**: Low priority (implement after core adult modules complete)
