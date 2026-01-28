# Whisparr v3 (eros) - Adult Content Structure Analysis

> Analysis of Whisparr v3 (eros) codebase for adult movie/show schema structure

**Source**: https://github.com/Whisparr/Whisparr/tree/eros
**API Docs**: https://github.com/Whisparr/Whisparr/tree/eros-api-docs
**Purpose**: Understanding adult content metadata structure for Revenge integration

---

## Key Findings

### 1. Metadata Providers

**Primary**: ThePornDB (TPDb)
- Frontend: `MetadataProvidedBy ThePornDB` attribution
- API: `https://api.whisparr.com/v3/{route}` (Whisparr metadata proxy)
- ID Fields: `TpdbId`, `ExternalId` (TPDb scene ID)

**Secondary**: StashDB
- Used by Stash integration (GraphQL endpoint: `https://theporndb.net/graphql`)
- Stash notification plugin for metadata sync
- Perceptual hashing (phash) support

### 2. Database Schema

**Episode (Scene) Structure**:
```csharp
public class Episode : ModelBase
{
    public int SeriesId { get; set; }
    public int TvdbId { get; set; }           // TPDb ID (NOT TheTVDB)
    public string ExternalId { get; set; }    // TPDb scene ID (e.g., "SSIS-123")
    public int SeasonNumber { get; set; }     // Year (e.g., 2023)
    public int EpisodeNumber { get; set; }
    public int? AbsoluteEpisodeNumber { get; set; }
    public string Title { get; set; }
    public string Overview { get; set; }
    public string AirDate { get; set; }       // Release date
    public DateTime? AirDateUtc { get; set; }
    public int Runtime { get; set; }
    public List<Actor> Actors { get; set; }   // Performers
    public Ratings Ratings { get; set; }
}
```

**Actor (Performer) Structure**:
```csharp
public class Actor : IEmbeddedDocument
{
    public int TpdbId { get; set; }           // TPDb performer ID
    public string Name { get; set; }
    public string Character { get; set; }
    public Gender Gender { get; set; }        // Female, Male, Other
    public List<MediaCover> Images { get; set; }
}

public enum Gender
{
    Female,
    Male,
    Other
}
```

**Series (Studio/Site) Structure**:
```csharp
public class Series : ModelBase
{
    public int TvdbId { get; set; }           // TPDb site ID
    public string Title { get; set; }         // Studio/site name
    public string Network { get; set; }       // Parent studio
    public SeriesTypes SeriesType { get; set; }
    public string Overview { get; set; }
    public int Runtime { get; set; }
    public List<MediaCover> Images { get; set; }
}

public enum SeriesTypes
{
    Standard,
    Jav    // Japanese Adult Video
}
```

### 3. Important Differences from TV Shows

**Seasons = Years**: Whisparr uses `SeasonNumber` as YEAR (e.g., 2023, 2024) NOT traditional TV seasons.

**TvdbId ≠ TheTVDB**: Field name is `TvdbId` but contains TPDb IDs (NOT TheTVDB IDs).

**ExternalId = TPDb Scene ID**: Scene identifier (e.g., "SSIS-123" for JAV, scene UUID for Western).

**No Traditional Episodes**: Each scene is an "episode" with `AbsoluteEpisodeNumber` (sequential per site).

**Network = Parent Studio**: `Network` field stores parent studio/network (e.g., "MindGeek" for Brazzers).

### 4. Metadata Export (XBMC/Kodi/Emby)

**Episode NFO** (.nfo files):
```xml
<episodedetails>
  <title>Scene Title</title>
  <season>2023</season>
  <aired>2023-08-15</aired>
  <plot>Scene description...</plot>
  <studio>Studio Name</studio>

  <!-- IDs -->
  <uniqueid type="tvdb" default="true">12345</uniqueid>
  <uniqueid type="whisparr">uuid</uniqueid>
  <uniqueid type="tpdb">scene-uuid</uniqueid>

  <!-- Performers -->
  <actor>
    <name>Performer Name</name>
    <role>Performer Name</role>
    <type>Performer</type>
    <thumb>https://path/to/image.jpg</thumb>
    <order>1</order>
  </actor>

  <rating>8.5</rating>
  <watched>false</watched>
  <thumb>https://scene-cover.jpg</thumb>
</episodedetails>
```

**Series NFO** (tvshow.nfo):
```xml
<tvshow>
  <title>Studio Name</title>
  <plot>Studio description...</plot>
  <uniqueid type="tvdb" default="true">site-id</uniqueid>
  <studio>Parent Studio</studio>
  <premiered>2020-01-01</premiered>
  <status>Continuing</status>
</tvshow>

<!-- OR URL-only mode -->
https://theporndb.net/sites/{site-id}
```

### 5. File Naming Tokens

**Whisparr supports adult-specific tokens**:
```
{Episode PerformersFemale}  // Female performers (max 4)
{Episode PerformersMale}    // Male performers (max 4)
{Episode PerformersOther}   // Other gender performers
{TpdbSceneId}              // TPDb scene ID (ExternalId)
```

Example: `Brazzers - 2023-08-15 - Scene Title [Performer1, Performer2].mp4`

### 6. TPDb Import Lists

**Performer Import**:
- API: `https://api.whisparr.com/v3/performer/{performer_id}/scenes`
- Returns: List of scenes with SiteId, EpisodeId, SiteName
- Auto-add all scenes from specific performer

**Site Import** (similar):
- Import all scenes from specific site/studio

### 7. Stash Integration

**Stash Notification Plugin**:
- GraphQL API: `http://localhost:9999/graphql`
- Actions: `metadataScan`, `metadataIdentify`
- Features:
  - Generate covers, previews, sprites, phashes
  - Identify scenes via StashDB (`https://theporndb.net/graphql`)
  - Merge metadata (studio, performers, tags)
  - Path mapping (Whisparr → Stash path translation)

**Metadata Identify**:
```graphql
mutation {
  metadataIdentify(
    input: {
      sources: [
        {source: {stash_box_endpoint: "https://theporndb.net/graphql"}},
        {source: {scraper_id: "builtin_autotag"}}
      ],
      options: {
        includeMalePerformers: true,
        setCoverImage: true,
        setOrganized: true,
        fieldOptions: [
          {field: "studio", strategy: MERGE, createMissing: true},
          {field: "performers", strategy: MERGE, createMissing: true},
          {field: "tags", strategy: MERGE, createMissing: true}
        ]
      },
      paths: ["/path/to/scene.mp4"]
    }
  )
}
```

---

## Schema Mapping for Revenge

### c.adult_movies Table

```sql
CREATE TABLE c.adult_movies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- External IDs
    tpdb_id INT,                          -- ThePornDB scene ID (Whisparr: TvdbId)
    tpdb_scene_id VARCHAR(200),           -- TPDb scene identifier (Whisparr: ExternalId, e.g., "SSIS-123")
    stashdb_id UUID,                      -- StashDB scene UUID
    tmdb_id INT,                          -- TMDb ID (if available)

    -- Basic Info
    title VARCHAR(500) NOT NULL,
    release_date DATE,
    year INT,                             -- Whisparr: SeasonNumber
    duration INT,                         -- Runtime in seconds
    description TEXT,                     -- Whisparr: Overview

    -- Studio/Site
    studio_id UUID REFERENCES c.studios(id),
    site_id UUID,                         -- TPDb site ID (Whisparr: Series.TvdbId)

    -- Metadata
    metadata_json JSONB,                  -- Full metadata (tpdb_data, stashdb_data, stash_data)

    -- Media Files
    file_path VARCHAR(1000),
    file_size BIGINT,

    -- Indexing
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);
```

### c.performers Table

```sql
CREATE TABLE c.performers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- External IDs
    tpdb_id INT UNIQUE,                   -- ThePornDB performer ID (Whisparr: Actor.TpdbId)
    stashdb_id UUID,                      -- StashDB performer UUID

    -- Basic Info
    name VARCHAR(200) NOT NULL,
    gender VARCHAR(20) CHECK (gender IN ('Female', 'Male', 'Other')),

    -- Metadata
    metadata_json JSONB,                  -- Full performer data (measurements, bio, etc.)

    -- Images
    image_url VARCHAR(500),

    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);
```

### c.studios Table

```sql
CREATE TABLE c.studios (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- External IDs
    tpdb_id INT UNIQUE,                   -- ThePornDB site ID (Whisparr: Series.TvdbId)
    stashdb_id UUID,                      -- StashDB studio UUID

    -- Basic Info
    name VARCHAR(200) NOT NULL,
    parent_studio_id UUID REFERENCES c.studios(id),  -- Whisparr: Network

    -- Metadata
    metadata_json JSONB,

    -- Images
    logo_url VARCHAR(500),

    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);
```

### c.adult_movie_performers (Junction Table)

```sql
CREATE TABLE c.adult_movie_performers (
    movie_id UUID REFERENCES c.adult_movies(id) ON DELETE CASCADE,
    performer_id UUID REFERENCES c.performers(id) ON DELETE CASCADE,
    "order" INT,                          -- Display order (Whisparr: Actor order)
    PRIMARY KEY (movie_id, performer_id)
);
```

---

## Key Takeaways for Revenge

1. **TPDb is primary**: Use ThePornDB API (NOT StashDB) as primary metadata source (Whisparr uses TPDb extensively).

2. **Year-based organization**: Use year (NOT seasons) for scene organization (Whisparr: `SeasonNumber = Year`).

3. **Multiple ID support**: Store TPDb ID, StashDB ID, TMDb ID (scenes can have multiple external IDs).

4. **Performer gender**: Track performer gender (Female, Male, Other) for filtering/search.

5. **Studio hierarchy**: Support parent/subsidiary studios (Network field in Whisparr).

6. **NFO export compatibility**: Export Whisparr-compatible NFOs (XBMC/Kodi/Emby format) with `<uniqueid type="tpdb">`.

7. **Stash integration**: Optional Stash GraphQL integration for users migrating from Stash.

8. **JAV support**: Recognize JAV series type (Japanese Adult Video) with `SeriesType = Jav`.

9. **File naming**: Support Whisparr-style file naming tokens (`{Episode PerformersFemale}`, `{TpdbSceneId}`).

10. **API namespace**: Use `/api/v1/c/` namespace for all adult endpoints (isolation from public API).

---

## Related Documentation

- [STASHDB.md](./STASHDB.md) - StashDB integration (supplementary)
- [THEPORNDB.md](./THEPORNDB.md) - ThePornDB integration (primary)
- [STASH.md](./STASH.md) - Stash self-hosted organizer
- [WHISPARR.md](../../servarr/WHISPARR.md) - Whisparr integration
- [ADULT_METADATA.md](../../../ADULT_METADATA.md) - Adult metadata system architecture
