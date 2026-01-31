# Comics Module

<!-- SOURCES: anilist, comicvine, fx, myanimelist, ogen, river, sqlc, sqlc-config, svelte-runes, svelte5, sveltekit -->

<!-- DESIGN: features/comics, 01_ARCHITECTURE, 02_DESIGN_PRINCIPLES, 03_METADATA_SYSTEM -->


> Digital comics/manga/graphic novel support with metadata from ComicVine, Marvel API, GCD


<!-- TOC-START -->

## Table of Contents

- [Status](#status)
- [Developer Resources](#developer-resources)
- [Overview](#overview)
- [File Format Support](#file-format-support)
- [Schema Design](#schema-design)
  - [Core Tables](#core-tables)
  - [Publishers](#publishers)
  - [Creators (Writers, Artists, Colorists, etc.)](#creators-writers-artists-colorists-etc)
  - [Genres & Tags](#genres-tags)
  - [User Data (Per-Module Isolation)](#user-data-per-module-isolation)
- [Metadata Sources](#metadata-sources)
  - [1. ComicVine API (Primary)](#1-comicvine-api-primary)
  - [2. Marvel API (Marvel Comics Only)](#2-marvel-api-marvel-comics-only)
  - [3. Grand Comics Database (GCD)](#3-grand-comics-database-gcd)
  - [4. AniList (Manga)](#4-anilist-manga)
  - [5. MyAnimeList (Manga Fallback)](#5-myanimelist-manga-fallback)
- [Reading Experience](#reading-experience)
  - [Web Reader](#web-reader)
- [Folder Structure vs TV Shows](#folder-structure-vs-tv-shows)
- [ComicInfo.xml Parsing](#comicinfoxml-parsing)
- [Relationships to Other Modules](#relationships-to-other-modules)
  - [Cross-Module Links](#cross-module-links)
- [API Endpoints](#api-endpoints)
- [Background Jobs (River)](#background-jobs-river)
- [Go Packages](#go-packages)
- [Summary](#summary)
- [Implementation Checklist](#implementation-checklist)
  - [Phase 1: Core Infrastructure](#phase-1-core-infrastructure)
  - [Phase 2: Database](#phase-2-database)
  - [Phase 3: File Processing](#phase-3-file-processing)
  - [Phase 4: Service Layer](#phase-4-service-layer)
  - [Phase 5: Metadata Integration](#phase-5-metadata-integration)
  - [Phase 6: Background Jobs](#phase-6-background-jobs)
  - [Phase 7: API Integration](#phase-7-api-integration)
- [Sources & Cross-References](#sources-cross-references)
  - [Cross-Reference Indexes](#cross-reference-indexes)
  - [Referenced Sources](#referenced-sources)
- [Related Design Docs](#related-design-docs)
  - [Related Topics](#related-topics)
  - [Indexes](#indexes)
- [Related](#related)

<!-- TOC-END -->

## Status

| Dimension | Status | Notes |
|-----------|--------|-------|
| Design | ✅ | Comprehensive spec with schema, file formats, metadata sources |
| Sources | ✅ | ComicVine, Marvel API, GCD, AniList, MyAnimeList documented |
| Instructions | ✅ | Implementation checklist added |
| Code | 🔴 |  |
| Linting | 🔴 |  |
| Unit Testing | 🔴 |  |
| Integration Testing | 🔴 |  |**Location**: `internal/content/comics/`

---

## Developer Resources

| Source | URL | Purpose |
|--------|-----|---------|
| ComicVine API | [comicvine.gamespot.com/api](https://comicvine.gamespot.com/api/) | Primary comics metadata |
| Marvel API | [developer.marvel.com](https://developer.marvel.com/) | Marvel Comics metadata |
| Grand Comics Database | [comics.org](https://www.comics.org/) | Community comics database |
| AniList GraphQL | [anilist.gitbook.io](https://anilist.gitbook.io/anilist-apiv2-docs/) | Manga metadata |
| MyAnimeList API | [myanimelist.net/apiconfig](https://myanimelist.net/apiconfig/references/api/v2) | Manga fallback |
| ComicInfo.xml Spec | [anansi-project.github.io](https://anansi-project.github.io/docs/comicinfo/documentation) | CBZ metadata format |

---

## Overview

Comics module provides cataloging, reading, and metadata management for digital comics (CBZ, CBR, CB7, PDF).

**Scope**:
- Comics (Western: Marvel, DC, Image, etc.)
- Manga (Japanese: Shonen Jump, Kodansha, etc.)
- Graphic Novels (standalone or collected editions)
- Webcomics (digital-first publications)

**Out of Scope** (separate modules):
- Books (text-only e-books → `book` module)
- Audiobooks → `audiobook` module

---

## File Format Support

| Format | Extension | Support | Notes |
|--------|-----------|---------|-------|
| **Comic Book Archive** | .cbz | ✅ Primary | ZIP of images (JPEG, PNG) |
| **Comic Book RAR** | .cbr | ✅ Primary | RAR of images (requires unrar) |
| **Comic Book 7z** | .cb7 | ✅ Primary | 7z of images (requires p7zip) |
| **PDF** | .pdf | ⚠️ Secondary | Extract images via pdfcpu/poppler |
| **EPUB** | .epub | ❌ No | Text-heavy (book module handles) |

---

## Schema Design

### Core Tables

```sql
CREATE TABLE comics (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Series info
    series_name VARCHAR(500) NOT NULL,
    series_id VARCHAR(100),  -- External ID (ComicVine, Marvel API)
    volume_number INT,       -- Volume 1, Volume 2, etc.
    issue_number VARCHAR(50), -- "1", "1.5", "Annual 2023"

    -- Publisher
    publisher_id UUID REFERENCES comic_publishers(id),

    -- Release info
    release_date DATE,
    release_year INT GENERATED ALWAYS AS (EXTRACT(YEAR FROM release_date)) STORED,

    -- Physical info
    page_count INT,
    cover_image_path VARCHAR(1000),

    -- File info
    file_path VARCHAR(1000) NOT NULL UNIQUE,
    file_size_bytes BIGINT,
    file_format VARCHAR(10) CHECK (file_format IN ('cbz', 'cbr', 'cb7', 'pdf')),

    -- Metadata
    description TEXT,
    metadata_json JSONB,  -- Store extra fields (writers, artists, colorists, etc.)

    -- Timestamps
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    scanned_at TIMESTAMPTZ,

    -- Full-text search
    search_vector tsvector GENERATED ALWAYS AS (
        setweight(to_tsvector('english', coalesce(series_name, '')), 'A') ||
        setweight(to_tsvector('english', coalesce(description, '')), 'B')
    ) STORED
);

CREATE INDEX idx_comics_series ON comics(series_name, volume_number, issue_number);
CREATE INDEX idx_comics_publisher ON comics(publisher_id);
CREATE INDEX idx_comics_release_date ON comics(release_date DESC);
CREATE INDEX idx_comics_search ON comics USING gin(search_vector);
CREATE INDEX idx_comics_metadata ON comics USING gin(metadata_json);
```

### Publishers

```sql
CREATE TABLE comic_publishers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(200) NOT NULL UNIQUE,
    external_id VARCHAR(100), -- ComicVine publisher ID
    logo_path VARCHAR(1000),
    description TEXT,
    founded_year INT,
    website_url VARCHAR(500),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_comic_publishers_name ON comic_publishers(name);
```

### Creators (Writers, Artists, Colorists, etc.)

```sql
CREATE TABLE comic_creators (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(200) NOT NULL,
    external_id VARCHAR(100), -- ComicVine person ID
    role VARCHAR(50), -- 'writer', 'penciller', 'inker', 'colorist', 'letterer', 'cover_artist'
    photo_path VARCHAR(1000),
    bio TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE comic_creator_roles (
    comic_id UUID REFERENCES comics(id) ON DELETE CASCADE,
    creator_id UUID REFERENCES comic_creators(id),
    role VARCHAR(50) NOT NULL,
    PRIMARY KEY (comic_id, creator_id, role)
);

CREATE INDEX idx_comic_creator_roles_comic ON comic_creator_roles(comic_id);
CREATE INDEX idx_comic_creator_roles_creator ON comic_creator_roles(creator_id);
```

### Genres & Tags

```sql
CREATE TABLE comic_genres (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL UNIQUE,
    domain VARCHAR(50) DEFAULT 'comics' CHECK (domain IN ('comics', 'manga'))
);

CREATE TABLE comic_genre_assignments (
    comic_id UUID REFERENCES comics(id) ON DELETE CASCADE,
    genre_id UUID REFERENCES comic_genres(id),
    PRIMARY KEY (comic_id, genre_id)
);

CREATE TABLE comic_tags (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL UNIQUE
);

CREATE TABLE comic_tag_assignments (
    comic_id UUID REFERENCES comics(id) ON DELETE CASCADE,
    tag_id UUID REFERENCES comic_tags(id),
    PRIMARY KEY (comic_id, tag_id)
);
```

### User Data (Per-Module Isolation)

```sql
CREATE TABLE comic_user_ratings (
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    comic_id UUID REFERENCES comics(id) ON DELETE CASCADE,
    rating DECIMAL(3, 1) CHECK (rating >= 0.0 AND rating <= 10.0),
    review TEXT,
    rated_at TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (user_id, comic_id)
);

CREATE TABLE comic_read_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    comic_id UUID REFERENCES comics(id) ON DELETE CASCADE,

    -- Reading progress
    current_page INT DEFAULT 1,
    total_pages INT, -- Cached from comics.page_count
    progress_percent INT GENERATED ALWAYS AS (
        CASE WHEN total_pages > 0 THEN (current_page * 100) / total_pages ELSE 0 END
    ) STORED,

    -- Completion tracking
    completed BOOLEAN DEFAULT FALSE,
    completed_at TIMESTAMPTZ,

    -- Timestamps
    started_at TIMESTAMPTZ DEFAULT NOW(),
    last_read_at TIMESTAMPTZ DEFAULT NOW(),

    UNIQUE (user_id, comic_id)
);

CREATE INDEX idx_comic_read_history_user ON comic_read_history(user_id, last_read_at DESC);
CREATE INDEX idx_comic_read_history_completed ON comic_read_history(user_id) WHERE NOT completed;

CREATE TABLE comic_favorites (
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    comic_id UUID REFERENCES comics(id) ON DELETE CASCADE,
    added_at TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (user_id, comic_id)
);

CREATE TABLE comic_reading_list (
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    comic_id UUID REFERENCES comics(id) ON DELETE CASCADE,
    added_at TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (user_id, comic_id)
);
```

---

## Metadata Sources

### 1. ComicVine API (Primary)

**Provider**: GameSpot (owned by Fandom)
**URL**: https://comicvine.gamespot.com/api/
**Coverage**: Western comics (Marvel, DC, Image, Dark Horse, etc.)
**API Key**: Free (1000 requests/hour)

**Endpoints**:
```
GET /api/issue/{id}       - Single issue details
GET /api/volume/{id}      - Series/volume details
GET /api/publisher/{id}   - Publisher info
GET /api/person/{id}      - Creator info
GET /api/search           - Search by title
```

**Response Example**:
```json
{
  "results": {
    "id": 12345,
    "volume": {
      "id": 5678,
      "name": "The Amazing Spider-Man"
    },
    "issue_number": "1",
    "name": "With Great Power...",
    "cover_date": "1963-03-01",
    "description": "Origin story...",
    "person_credits": [
      {"id": 1, "name": "Stan Lee", "role": "writer"},
      {"id": 2, "name": "Steve Ditko", "role": "penciller"}
    ]
  }
}
```

---

### 2. Marvel API (Marvel Comics Only)

**Provider**: Marvel Entertainment
**URL**: https://developer.marvel.com/
**Coverage**: Marvel Comics only (comprehensive)
**API Key**: Free (3000 requests/day)

**Endpoints**:
```
GET /v1/public/comics/{id}
GET /v1/public/series/{id}
GET /v1/public/creators/{id}
```

**Rate Limiting**: 3000 requests/day, 100 requests/second

---

### 3. Grand Comics Database (GCD)

**Provider**: Community-driven open database
**URL**: https://www.comics.org/
**Coverage**: Comprehensive (100+ years of comics)
**API**: REST API (free, no key required)

**Endpoints**:
```
GET /api/issue/{id}
GET /api/series/{id}
GET /api/publisher/{id}
```

**Notes**: No official API docs, scraping may be required (check robots.txt).

---

### 4. AniList (Manga)

**Provider**: AniList.co
**URL**: https://anilist.gitbook.io/anilist-apiv2-docs/
**Coverage**: Manga/Manhwa/Manhua
**API**: GraphQL (no key required, rate limited)

**Query Example**:
```graphql
query {
  Media(id: 30013, type: MANGA) {
    title {
      romaji
      english
      native
    }
    chapters
    volumes
    description
    coverImage { large }
    genres
    tags { name }
  }
}
```

---

### 5. MyAnimeList (Manga Fallback)

**Provider**: MyAnimeList.net
**URL**: https://myanimelist.net/apiconfig/references/api/v2
**Coverage**: Manga/Light Novels
**API Key**: Required (OAuth2)

---

## Reading Experience

### Web Reader

**Technology**: Canvas-based reader with page preloading

**Features**:
- **Single page mode** (desktop default)
- **Double page spread** (manga/comics with spreads)
- **Continuous scroll** (webtoon-style)
- **Fit to width/height** (zoom options)
- **Page preloading** (prefetch next 3 pages)
- **Keyboard shortcuts** (arrow keys, spacebar)
- **Touch gestures** (swipe, pinch-zoom)
- **Bookmarking** (save page progress)
- **Night mode** (invert colors for dark reading)

**Implementation**:
```typescript
// SvelteKit component
<script lang="ts">
  import { onMount } from 'svelte';

  let currentPage = 1;
  let totalPages = 0;
  let zoomLevel = 'fit-width'; // 'fit-width', 'fit-height', 'original'

  async function loadPage(pageNum: number) {
    const img = await fetch(`/api/comics/${comicId}/page/${pageNum}`);
    // Render to canvas
  }

  function nextPage() {
    if (currentPage < totalPages) {
      currentPage++;
      loadPage(currentPage);
      preloadPages(currentPage + 1, currentPage + 3);
    }
  }
</script>

<canvas bind:this={canvasRef} on:click={nextPage} />
```

---

## Folder Structure vs TV Shows

**TV Show Structure** (hierarchical):
```
/Series Name/
  Season 01/
    S01E01.mkv
    S01E02.mkv
  Season 02/
    S02E01.mkv
```

**Comics Structure** (flat or series-based):
```
/Comics/
  /Marvel/
    /The Amazing Spider-Man/
      The Amazing Spider-Man #001 (1963).cbz
      The Amazing Spider-Man #002 (1963).cbz
  /DC/
    /Batman/
      Batman #001 (1940).cbz
```

**OR Flat**:
```
/Comics/
  The Amazing Spider-Man #001 (1963).cbz
  The Amazing Spider-Man #002 (1963).cbz
  Batman #001 (1940).cbz
```

**Parsing Strategy**:
- Extract series name, issue number, year from filename
- Fallback to metadata inside archive (ComicInfo.xml)
- Match against ComicVine API for enrichment

---

## ComicInfo.xml Parsing

Many CBZ files contain `ComicInfo.xml` metadata:

```xml
<?xml version="1.0"?>
<ComicInfo>
  <Series>The Amazing Spider-Man</Series>
  <Number>1</Number>
  <Volume>1</Volume>
  <Year>1963</Year>
  <Month>3</Month>
  <Writer>Stan Lee</Writer>
  <Penciller>Steve Ditko</Penciller>
  <Publisher>Marvel Comics</Publisher>
  <Genre>Superhero</Genre>
  <PageCount>22</PageCount>
  <Summary>Origin story...</Summary>
</ComicInfo>
```

**Parsing** (Go):
```go
type ComicInfo struct {
    Series     string `xml:"Series"`
    Number     string `xml:"Number"`
    Volume     int    `xml:"Volume"`
    Year       int    `xml:"Year"`
    Month      int    `xml:"Month"`
    Writer     string `xml:"Writer"`
    Penciller  string `xml:"Penciller"`
    Publisher  string `xml:"Publisher"`
    PageCount  int    `xml:"PageCount"`
    Summary    string `xml:"Summary"`
}

func parseComicInfo(cbzPath string) (*ComicInfo, error) {
    rc, err := zip.OpenReader(cbzPath)
    defer rc.Close()

    for _, f := range rc.File {
        if strings.EqualFold(f.Name, "ComicInfo.xml") {
            r, _ := f.Open()
            var ci ComicInfo
            xml.NewDecoder(r).Decode(&ci)
            return &ci, nil
        }
    }
    return nil, errors.New("no ComicInfo.xml found")
}
```

---

## Relationships to Other Modules

### Cross-Module Links

**Comics → Movies/TV**:
- Link comics to adapted movies (e.g., The Amazing Spider-Man #1 → Spider-Man (2002))
- Display "Adapted to Film" badge on comic detail page

**Schema**:
```sql
CREATE TABLE comic_adaptations (
    comic_id UUID REFERENCES comics(id) ON DELETE CASCADE,

    -- Polymorphic link (movie OR tvshow)
    adapted_movie_id UUID REFERENCES movies(id) ON DELETE CASCADE,
    adapted_series_id UUID REFERENCES series(id) ON DELETE CASCADE,

    adaptation_type VARCHAR(50), -- 'direct', 'inspired_by', 'loosely_based'
    notes TEXT,

    CHECK (
        (adapted_movie_id IS NOT NULL AND adapted_series_id IS NULL) OR
        (adapted_movie_id IS NULL AND adapted_series_id IS NOT NULL)
    ),

    PRIMARY KEY (comic_id, COALESCE(adapted_movie_id, adapted_series_id))
);
```

**UI Display**:
```
[Comic Detail Page]
The Amazing Spider-Man #1 (1963)

Related Media:
  ├─ Spider-Man (2002) [Movie]
  ├─ The Amazing Spider-Man (2012) [Movie]
  └─ Spider-Man: The Animated Series (1994) [TV Show]
```

---

## API Endpoints

```
GET    /api/v1/comics                # List all comics
GET    /api/v1/comics/{id}           # Single comic details
POST   /api/v1/comics                # Add comic (manual)
PUT    /api/v1/comics/{id}           # Update comic
DELETE /api/v1/comics/{id}           # Delete comic

GET    /api/v1/comics/{id}/pages     # List all pages (image URLs)
GET    /api/v1/comics/{id}/page/{n}  # Get specific page image
GET    /api/v1/comics/{id}/cover     # Get cover image

GET    /api/v1/comics/{id}/metadata  # Metadata sources (ComicVine, Marvel API)
POST   /api/v1/comics/{id}/refresh   # Re-fetch metadata

GET    /api/v1/comics/series/{name}  # List comics in series
GET    /api/v1/comics/publishers     # List publishers
GET    /api/v1/comics/creators       # List creators

# User data
GET    /api/v1/comics/{id}/progress  # User reading progress
PUT    /api/v1/comics/{id}/progress  # Update progress (page, completed)
POST   /api/v1/comics/{id}/rate      # Rate comic
POST   /api/v1/comics/{id}/favorite  # Add to favorites
```

---

## Background Jobs (River)

```go
// Scan comics library
type ScanComicsLibraryArgs struct {
    LibraryID uuid.UUID `json:"library_id"`
    FullScan  bool      `json:"full_scan"`
}

// Fetch metadata from ComicVine
type FetchComicMetadataArgs struct {
    ComicID uuid.UUID `json:"comic_id"`
    Provider string   `json:"provider"` // "comicvine", "marvel", "gcd"
}

// Extract pages from CBZ/CBR
type ExtractComicPagesArgs struct {
    ComicID uuid.UUID `json:"comic_id"`
}

// Generate thumbnails for pages
type GenerateComicThumbnailsArgs struct {
    ComicID uuid.UUID `json:"comic_id"`
}
```

---

## Go Packages

> See [00_SOURCE_OF_TRUTH.md](../../00_SOURCE_OF_TRUTH.md#go-dependencies-core) for package versions.

Key packages used:
- **archive/zip** - CBZ extraction (stdlib)
- **rardecode** - CBR extraction (RAR)
- **sevenzip** - CB7 extraction (7z)
- **pdfcpu** - PDF page extraction
- **imaging** - Image resizing (thumbnails)

---

## Summary

| Aspect | Value |
|--------|-------|
| **File Formats** | CBZ, CBR, CB7, PDF |
| **Metadata Sources** | ComicVine (primary), Marvel API, GCD, AniList (manga) |
| **Schema** | comics, comic_publishers, comic_creators, comic_creator_roles, genres, tags |
| **User Data** | ratings, read_history (page progress), favorites, reading_list |
| **Web Reader** | Canvas-based, single/double page, continuous scroll, preloading |
| **Folder Structure** | Flat or series-based (NOT hierarchical like TV shows) |
| **Cross-Module** | Link comics → movies/TV adaptations |
| **Background Jobs** | Library scan, metadata fetch, page extraction, thumbnail generation |

**Development Priority**: ⚠️ **Medium** (after core video/music modules, before Live TV)

---

## Implementation Checklist

### Phase 1: Core Infrastructure
- [ ] Create `internal/content/comics/` package structure
- [ ] Define `entity.go` with Comic, Publisher, Creator structs
- [ ] Create `repository.go` interface definition
- [ ] Implement `repository_pg.go` with sqlc queries
- [ ] Add fx module wiring in `module.go`

### Phase 2: Database
- [ ] Create migration `000XXX_create_comics_schema.up.sql`
- [ ] Create `comics` table with all columns
- [ ] Create `comic_publishers` table
- [ ] Create `comic_creators` and `comic_creator_roles` tables
- [ ] Create `comic_genres` and `comic_tags` tables
- [ ] Add user data tables (ratings, read_history, favorites)
- [ ] Add indexes (series_name, publisher_id, release_date, search_vector)
- [ ] Write sqlc queries in `queries/comics/`

### Phase 3: File Processing
- [ ] Implement CBZ extraction (archive/zip)
- [ ] Implement CBR extraction (rardecode)
- [ ] Implement CB7 extraction (sevenzip)
- [ ] Implement PDF page extraction (pdfcpu)
- [ ] Parse ComicInfo.xml metadata
- [ ] Extract cover images for thumbnails

### Phase 4: Service Layer
- [ ] Implement `service.go` with otter caching
- [ ] Add Comic operations (Get, List, Create, Update, Delete)
- [ ] Add Series operations (List by series, group by volume)
- [ ] Add reading progress tracking
- [ ] Implement cache invalidation on mutations

### Phase 5: Metadata Integration
- [ ] Implement ComicVine provider
- [ ] Implement Marvel API provider (optional)
- [ ] Implement AniList provider for manga
- [ ] Add metadata matching logic (title + issue → external ID)
- [ ] Implement image download and caching

### Phase 6: Background Jobs
- [ ] Create River job definitions in `jobs.go`
- [ ] Implement `ScanComicsLibraryJob`
- [ ] Implement `FetchComicMetadataJob`
- [ ] Implement `ExtractComicPagesJob`
- [ ] Implement `GenerateComicThumbnailsJob`

### Phase 7: API Integration
- [ ] Define OpenAPI endpoints for comics
- [ ] Generate ogen handlers
- [ ] Wire handlers to service layer
- [ ] Add page streaming endpoint
- [ ] Add authentication/authorization checks

---


## Related

- [Library Service](../../services/LIBRARY.md) - Library management
- [Movie Module](../video/MOVIE_MODULE.md) - Cross-module adaptations
- [Metadata Service](../../services/METADATA.md) - Provider patterns

