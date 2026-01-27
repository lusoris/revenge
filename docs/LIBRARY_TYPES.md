# Library Types

> Extended library types for Jellyfin Go, including support for books, audiobooks, podcasts, and adult content.

## Overview

Jellyfin Go extends the original Jellyfin library types to support additional content categories:

- **Books & Audiobooks** - For Calibre/Audiobookshelf-like functionality
- **Podcasts** - RSS-based podcast management
- **Adult Content** - Whisparr/Stash integration for adult movies and shows

## Library Type Enum

```sql
CREATE TYPE library_type AS ENUM (
    -- Original Jellyfin types
    'movies',           -- Movie library
    'tvshows',          -- TV series library
    'music',            -- Music library
    'musicvideos',      -- Music video library
    'photos',           -- Photo library
    'homevideos',       -- Home video library
    'boxsets',          -- Movie collections/box sets
    'livetv',           -- Live TV & DVR
    'playlists',        -- Playlist container
    'mixed',            -- Mixed content library

    -- Extended types (NEW)
    'books',            -- E-books (epub, pdf, mobi)
    'audiobooks',       -- Audiobooks (m4b, mp3)
    'podcasts',         -- Podcast feeds
    'adult_movies',     -- Adult movies (Whisparr integration)
    'adult_shows'       -- Adult series (Stash integration)
);
```

## Media Type Enum

```sql
CREATE TYPE media_type AS ENUM (
    -- Video types
    'movie',
    'episode',
    'musicvideo',
    'trailer',
    'homevideo',

    -- Audio types
    'audio',            -- Music track
    'audiobook_chapter',-- Audiobook chapter
    'podcast_episode',  -- Podcast episode

    -- Image types
    'photo',

    -- Document types
    'book',             -- E-book (NEW)

    -- Collection types
    'series',           -- TV series
    'season',           -- TV season
    'album',            -- Music album
    'audiobook',        -- Audiobook (container)
    'podcast',          -- Podcast feed (container)
    'artist',           -- Music artist
    'boxset',           -- Collection
    'playlist',
    'folder',

    -- Live TV
    'channel',
    'program',
    'recording'
);
```

## Library Type Details

### Standard Libraries

| Type | Content | Source | Metadata Providers |
|------|---------|--------|-------------------|
| `movies` | Movies | Local files | TMDB, IMDB, TVDB |
| `tvshows` | TV Series | Local files | TMDB, TVDB, TheTVDB |
| `music` | Music tracks | Local files | MusicBrainz, Last.fm |
| `musicvideos` | Music videos | Local files | TMDB, MusicBrainz |
| `photos` | Images | Local files | EXIF data |
| `homevideos` | Personal videos | Local files | EXIF/file metadata |
| `boxsets` | Collections | Virtual | TMDB Collections |
| `livetv` | Live channels | EPG providers | XMLTV, SchedulesDirect |
| `playlists` | Playlists | User-created | - |
| `mixed` | Multiple types | Local files | All providers |

### Extended Libraries

#### Books (`books`)

E-book library supporting:
- **Formats**: EPUB, PDF, MOBI, AZW3, CBZ, CBR
- **Metadata Sources**: OpenLibrary, Google Books, Calibre DB
- **Features**:
  - Reading progress sync (OPDS)
  - Series tracking
  - Author management
  - Cover extraction

```yaml
library:
  type: books
  paths:
    - /media/books
  settings:
    parseSubdirectoriesAsSeries: true
    extractCovers: true
    opdsEnabled: true
```

#### Audiobooks (`audiobooks`)

Audiobook library supporting:
- **Formats**: M4B, MP3, M4A, AAC, FLAC
- **Metadata Sources**: Audible, OpenLibrary, Audnexus
- **Features**:
  - Chapter detection
  - Listening progress sync
  - Narrator tracking
  - Series management

```yaml
library:
  type: audiobooks
  paths:
    - /media/audiobooks
  settings:
    detectChapters: true
    audibleIntegration: true
    syncProgress: true
```

#### Podcasts (`podcasts`)

Podcast library supporting:
- **Source**: RSS feeds
- **Formats**: MP3, M4A, AAC
- **Features**:
  - Auto-download new episodes
  - Episode filtering
  - Listen progress tracking
  - OPML import/export

```yaml
library:
  type: podcasts
  settings:
    autoDownload: true
    retentionDays: 30
    maxEpisodes: 50
```

#### Adult Movies (`adult_movies`)

Adult movie library supporting:
- **Integration**: Whisparr, Stash
- **Metadata Sources**: TPDB (ThePornDB), Stash-Box
- **Rating**: Always `normalized_level: 100` (Adult)
- **Features**:
  - Performer tracking
  - Studio management
  - Scene-level metadata

```yaml
library:
  type: adult_movies
  paths:
    - /media/adult/movies
  settings:
    provider: stash-box  # or 'tpdb'
    stashBoxEndpoint: "https://stashdb.org/graphql"
    autoOrganize: true
```

#### Adult Shows (`adult_shows`)

Adult series library supporting:
- **Integration**: Stash, Whisparr
- **Metadata Sources**: Stash-Box, TPDB
- **Rating**: Always `normalized_level: 100` (Adult)
- **Features**:
  - Site/Studio as "Network"
  - Performer tracking

```yaml
library:
  type: adult_shows
  paths:
    - /media/adult/series
  settings:
    provider: stash-box
```

## Adult Content Handling

### Access Control

```go
func (s *LibraryService) ListUserLibraries(ctx context.Context, userID uuid.UUID) ([]Library, error) {
    user := s.getUser(ctx, userID)

    libraries := s.getAllLibraries(ctx)

    // Filter based on user's adult_content_enabled flag
    var allowed []Library
    for _, lib := range libraries {
        if !lib.IsAdultLibrary() || user.AdultContentEnabled {
            allowed = append(allowed, lib)
        }
    }

    return allowed, nil
}

func (l *Library) IsAdultLibrary() bool {
    return l.Type == LibraryTypeAdultMovies || l.Type == LibraryTypeAdultShows
}
```

### Metadata Provider Selection

```go
func (s *MetadataService) GetProvidersForLibrary(libType LibraryType) []MetadataProvider {
    switch libType {
    case LibraryTypeMovies, LibraryTypeTvShows:
        return []MetadataProvider{TMDB, IMDB, TVDB}
    case LibraryTypeMusic:
        return []MetadataProvider{MusicBrainz, LastFM}
    case LibraryTypeBooks:
        return []MetadataProvider{OpenLibrary, GoogleBooks}
    case LibraryTypeAudiobooks:
        return []MetadataProvider{Audnexus, Audible, OpenLibrary}
    case LibraryTypePodcasts:
        return []MetadataProvider{PodcastIndex, ITunes}
    case LibraryTypeAdultMovies, LibraryTypeAdultShows:
        return []MetadataProvider{StashBox, TPDB}
    default:
        return []MetadataProvider{Local}
    }
}
```

### Database Indexes

```sql
-- Efficient library type filtering
CREATE INDEX idx_libraries_type ON libraries(type);
CREATE INDEX idx_libraries_type_adult ON libraries(type)
    WHERE type IN ('adult_movies', 'adult_shows');

-- Efficient content filtering for adult libraries
CREATE INDEX idx_media_items_library_type ON media_items(library_id, media_type);
```

## Migration from Original Jellyfin

When migrating an existing Jellyfin installation:

1. Standard library types map directly
2. Collections → `boxsets`
3. Music Videos → `musicvideos`
4. Home Videos → `homevideos`
5. Adult content requires:
   - Create new `adult_movies` or `adult_shows` library
   - Move content from `movies`/`tvshows`
   - Re-scan with adult metadata providers

## See Also

- [CONTENT_RATING.md](CONTENT_RATING.md) - Age restriction and rating systems
- [ARCHITECTURE.md](ARCHITECTURE.md) - System architecture
- [PHASE1_CHECKLIST.md](PHASE1_CHECKLIST.md) - Implementation status
