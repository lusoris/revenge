# TV Show Content Module

<!-- DESIGN: features/video -->

**Package**: `internal/content/tvshow`
**fx Module**: `tvshow.Module` + `tvshowjobs.Module`

> Hierarchical TV show management (Series > Seasons > Episodes > Files) with TMDb metadata, library scanning, file matching, and watch progress tracking

---

## Module Structure

```
internal/content/tvshow/
├── service.go             # Service interface (89 methods) + tvService implementation
├── repository.go          # Repository interface (106 methods)
├── repository_postgres.go # PostgreSQL implementation (pgxpool + sqlc tvshowdb.Queries)
├── types.go               # Domain types: Series, Season, Episode, EpisodeFile, credits, etc.
├── metadata_provider.go   # MetadataProvider interface (8 methods)
├── module.go              # fx.Module("tvshow") wiring
├── adapters/
│   ├── metadata_adapter.go  # TMDb TV client setup (4 req/sec, burst 10, 17 genre mappings)
│   └── scanner_adapter.go   # TVShowFileParser (5 regex patterns: SxxExx, daily, legacy)
├── jobs/
│   ├── module.go            # fx.Module("tvshowjobs"), RegisterWorkers
│   └── jobs.go              # 5 workers: library_scan, metadata_refresh, file_match, search_index, series_refresh
└── db/                      # sqlc-generated (tvshowdb package)
```

## Domain Types

### Series

Hierarchical root entity with relations:

| Field Group | Fields |
|-------------|--------|
| Core | ID, Title, OriginalTitle, OriginalLanguage, Tagline, Overview |
| External IDs | TMDbID, TVDbID, IMDbID, SonarrID |
| Show Info | Status (Returning/Ended/Canceled), Type (Scripted/Reality/Documentary), FirstAirDate, LastAirDate |
| Media | PosterPath, BackdropPath, VoteAverage, VoteCount, Popularity, TrailerURL, Homepage |
| Stats | TotalSeasons, TotalEpisodes |
| Timestamps | CreatedAt, UpdatedAt, MetadataUpdatedAt |
| i18n | TitlesI18n, TaglinesI18n, OverviewsI18n (`map[string]string`), AgeRatings (`map[string]map[string]string`) |
| Relations | Seasons, Genres, Credits, Networks (populated on demand) |

Methods: `GetTitle(lang)`, `GetTagline(lang)`, `GetOverview(lang)`, `GetAgeRating(country, system)`, `GetAvailableLanguages()`, `GetAvailableAgeRatingCountries()`, `IsEnded()`

### Season

| Field Group | Fields |
|-------------|--------|
| Core | ID, SeriesID, TMDbID, SeasonNumber, Name, Overview, EpisodeCount |
| Media | PosterPath, AirDate, VoteAverage |
| i18n | NamesI18n, OverviewsI18n |
| Relations | Episodes (populated on demand) |

Methods: `GetName(lang)`, `GetOverview(lang)`, `IsSpecials()` (season 0)

### Episode

| Field Group | Fields |
|-------------|--------|
| Core | ID, SeriesID, SeasonID, SeasonNumber, EpisodeNumber, Title, Overview |
| External IDs | TMDbID, TVDbID, IMDbID |
| Media | AirDate, Runtime, VoteAverage, VoteCount, StillPath, ProductionCode |
| i18n | TitlesI18n, OverviewsI18n |
| Relations | Files, Credits |

Methods: `GetTitle(lang)`, `GetOverview(lang)`, `EpisodeCode()` (e.g. "S01E05"), `HasAired()`

### Supporting Types

| Type | Key Fields |
|------|-----------|
| EpisodeFile | FilePath, Resolution, VideoCodec, AudioCodec, BitrateKbps, DurationSeconds, AudioLanguages, SubtitleLanguages, SonarrFileID |
| SeriesCredit | TMDbPersonID, Name, CreditType (cast/crew), Character, Job, Department, CastOrder |
| EpisodeCredit | Same as SeriesCredit but for individual episodes (guest stars) |
| SeriesGenre | TMDbGenreID, Name |
| Network | TMDbID, Name, LogoPath, OriginCountry |
| EpisodeWatched | UserID, EpisodeID, ProgressSeconds, DurationSeconds, IsCompleted, WatchCount |
| SeriesWatchStats | WatchedCount, InProgressCount, TotalWatches, TotalEpisodes + `CompletionPercent()` |
| UserTVStats | SeriesCount, EpisodesWatched, EpisodesInProgress, TotalWatches |
| ContinueWatchingItem | Series + last episode info + progress |
| EpisodeWithSeriesInfo | Episode + SeriesTitle + SeriesPosterPath |
| SeasonWithEpisodeCount | Season + ActualEpisodeCount |

## Service Interface

89 exported methods on the `Service` interface, organized by entity:

| Category | Count | Key Methods |
|----------|-------|-------------|
| Series CRUD | 14 | Get, GetByTMDbID/TVDbID/SonarrID, List, Count, Search, SearchAnyLanguage, RecentlyAdded, ByGenre, ByNetwork, ByStatus, Create, Update, Delete |
| Seasons | 9 | Get, GetByNumber, ListBySeries, ListWithEpisodeCount, Create, Upsert, Update, Delete, DeleteBySeries |
| Episodes | 16 | Get, GetByTMDbID, GetByNumber, ListBySeries/Season/SeasonNumber, Recent, Upcoming, Count, Create, Upsert, Update, Delete, DeleteBySeason/Series |
| Files | 8 | Get, GetByPath/SonarrID, ListByEpisode, Create, Update, Delete, DeleteByEpisode |
| Credits | 8 | CreateSeries/EpisodeCredit, ListSeriesCast/Crew, ListEpisodeGuestStars/Crew, DeleteSeries/EpisodeCredits |
| Genres | 3 | Add, List, Delete |
| Networks | 6 | Create, Get, GetByTMDbID, ListBySeries, AddToSeries, DeleteFromSeries |
| Watch Progress | 12 | CreateOrUpdate, MarkWatched, Get, Delete, DeleteBySeries, ContinueWatching, WatchedBySeries/User, SeriesStats, UserStats, NextUnwatched |
| Metadata | 3 | RefreshSeriesMetadata, RefreshSeasonMetadata, RefreshEpisodeMetadata |

Implementation: `tvService` struct with `repo Repository` + `metadataProvider MetadataProvider`.

## Repository Interface

106 exported methods on the `Repository` interface. Same categories as Service but includes additional database-level operations (UpdateSeriesStats, Upsert variants). Implementation: `postgresRepository` wrapping `pgxpool.Pool` + sqlc `tvshowdb.Queries`.

Conversion helpers: `dbSeriesToSeries()`, `dbSeasonToSeason()`, `dbEpisodeToEpisode()`, `dbEpisodeFileToEpisodeFile()`, etc. Same JSON marshaling patterns as movie module for i18n fields.

## MetadataProvider Interface

```go
type MetadataProvider interface {
    SearchSeries(ctx, query, year) ([]*Series, error)
    EnrichSeries(ctx, series, opts...) error
    EnrichSeason(ctx, season, seriesTMDbID, opts...) error
    EnrichEpisode(ctx, episode, seriesTMDbID, opts...) error
    GetSeriesCredits(ctx, seriesID, tmdbID) ([]SeriesCredit, error)
    GetSeriesGenres(ctx, seriesID, tmdbID) ([]SeriesGenre, error)
    GetSeriesNetworks(ctx, tmdbID) ([]Network, error)
    ClearCache()
}
```

Implementation injected via `metadatafx` module as `TVShowMetadataAdapter`. Uses shared `metadata.BaseClient` with TMDb TV endpoints.

### TMDb Integration

Append-to-response queries for efficient fetching:
- **Series**: `credits,images,content_ratings,external_ids,translations,alternative_titles`
- **Season**: `credits,images,translations`
- **Episode**: `credits,images,translations`

17 TV genre mappings (Action & Adventure=10759, Animation=16, Comedy=35, Crime=80, Documentary=99, Drama=18, etc.)

Status mappings: "Returning Series", "Ended", "Canceled", "In Production", "Planned"
Type mappings: "Scripted", "Reality", "Documentary", "Miniseries", "News", "Talk Show"

## File Parser

`TVShowFileParser` implements `scanner.FileParser` with 5 regex patterns:

| Pattern | Example | Extracts |
|---------|---------|----------|
| SxxExx | `Breaking.Bad.S01E05.mkv` | season, episode, episode_title, series_year |
| Season x Episode x | `Season 1 Episode 5` | season, episode |
| Daily show | `2024.01.15` or `2024-01-15` | air_year, air_month, air_day |
| Legacy x.xx | `1.05` or `1-05` | season, episode |
| Multi-episode | `S01E05E06` | season, episode, end_episode |

Also supports `ParseFromPath(filePath)` to extract series title from parent directory structure.

## Background Workers

5 River workers registered via `tvshowjobs.RegisterWorkers()`:

| Worker | Kind | Queue | Timeout | Purpose |
|--------|------|-------|---------|---------|
| LibraryScanWorker | `tvshow_library_scan` | bulk | 30m | Scan directories, auto-create series/seasons/episodes |
| MetadataRefreshWorker | `tvshow_metadata_refresh` | default | 15m | Refresh series/season/episode metadata (batch or targeted) |
| FileMatchWorker | `tvshow_file_match` | default | 5m | Match individual file to episode, auto-create if needed |
| SearchIndexWorker | `tvshow_search_index` | bulk | 10m | Stub (TV search not yet implemented) |
| SeriesRefreshWorker | `tvshow_series_refresh` | default | 10m | Cascading refresh: series > seasons > episodes |

### Library Scan Flow

1. TVShowFileParser scans paths for video files
2. For each file: check if already matched, skip if exists (unless Force)
3. If AutoCreate enabled: search existing series > search TMDb > create series/season/episode/file
4. Reports progress via River job client

### Series Refresh Flow

1. Refresh series metadata from TMDb
2. If RefreshSeasons: iterate all seasons, refresh each
3. If RefreshEpisodes: for each season, refresh all episodes
4. Progress reported per season

## API Endpoints

25 endpoints in `internal/api/tvshow_handlers.go` (ogen-generated types):

**Series** (10):

| Endpoint | Handler |
|----------|---------|
| `GET /tv` | ListTVShows |
| `GET /tv/search` | SearchTVShows |
| `GET /tv/recently-added` | GetRecentlyAddedTVShows |
| `GET /tv/continue-watching` | GetTVContinueWatching |
| `GET /tv/stats` | GetUserTVStats |
| `GET /tv/recent-episodes` | GetRecentEpisodes |
| `GET /tv/upcoming-episodes` | GetUpcomingEpisodes |
| `GET /tv/{id}` | GetTVShow |
| `GET /tv/{id}/seasons` | GetTVShowSeasons |
| `GET /tv/{id}/episodes` | GetTVShowEpisodes |

**Series Metadata & Relations** (5):

| Endpoint | Handler |
|----------|---------|
| `GET /tv/{id}/cast` | GetTVShowCast |
| `GET /tv/{id}/crew` | GetTVShowCrew |
| `GET /tv/{id}/genres` | GetTVShowGenres |
| `GET /tv/{id}/networks` | GetTVShowNetworks |
| `POST /tv/{id}/refresh` | RefreshTVShowMetadata |

**Watch Stats** (2):

| Endpoint | Handler |
|----------|---------|
| `GET /tv/{id}/watch-stats` | GetTVShowWatchStats |
| `GET /tv/{id}/next-episode` | GetTVShowNextEpisode |

**Seasons & Episodes** (8):

| Endpoint | Handler |
|----------|---------|
| `GET /tv/{id}/seasons/{seasonNumber}` | GetTVSeason |
| `GET /tv/{id}/seasons/{seasonNumber}/episodes` | GetTVSeasonEpisodes |
| `GET /tv/{id}/seasons/{sn}/episodes/{en}` | GetTVEpisode |
| `GET /tv/{id}/seasons/{sn}/episodes/{en}/files` | GetTVEpisodeFiles |
| `GET /tv/{id}/seasons/{sn}/episodes/{en}/progress` | GetTVEpisodeProgress |
| `PUT /tv/{id}/seasons/{sn}/episodes/{en}/progress` | UpdateTVEpisodeProgress |
| `DELETE /tv/{id}/seasons/{sn}/episodes/{en}/progress` | DeleteTVEpisodeProgress |
| `POST /tv/{id}/seasons/{sn}/episodes/{en}/watched` | MarkTVEpisodeWatched |

Converter functions in `tvshow_converters.go` bridge domain types to ogen API types.

## Dependencies

- `github.com/jackc/pgx/v5/pgxpool` - PostgreSQL (via repository)
- `github.com/imroc/req/v3` - HTTP client for TMDb (via shared metadata.BaseClient)
- `github.com/riverqueue/river` - Background job processing
- `github.com/google/uuid` - UUID generation
- `github.com/shopspring/decimal` - Decimal types for ratings
- `go.uber.org/zap` - Structured logging
- `go.uber.org/fx` - Dependency injection
- Shared packages: `content/shared/scanner`, `content/shared/matcher`, `content/shared/metadata`, `content/shared/library`, `content/shared/jobs`

## fx Wiring

```go
// tvshow.Module provides:
fx.Provide(NewPostgresRepository)  // → Repository
fx.Provide(provideService)         // → Service (repo + metadataProvider)

// tvshowjobs.Module provides:
fx.Provide(provideLibraryScanWorker, provideMetadataRefreshWorker,
           provideFileMatchWorker, provideSearchIndexWorker, provideSeriesRefreshWorker)
```

## Related Documentation

- [MOVIE_MODULE.md](MOVIE_MODULE.md) - Movie content module (similar architecture)
- [../../architecture/METADATA_SYSTEM.md](../../architecture/METADATA_SYSTEM.md) - Provider chain and caching
- [../../infrastructure/JOBS.md](../../infrastructure/JOBS.md) - River job queue setup
- [../../infrastructure/CACHE.md](../../infrastructure/CACHE.md) - L1/L2 caching infrastructure
- [../../services/LIBRARY.md](../../services/LIBRARY.md) - Library management service
