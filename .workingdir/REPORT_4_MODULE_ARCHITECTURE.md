# Revenge Module Architecture: Shared Functionality & Adapter Pattern Strategy

**Generated**: 2026-02-05
**Focus**: Eliminating code duplication across content modules
**Goal**: 60-70% code reduction through shared abstractions

---

## Executive Summary

The Revenge project has a **complete movie module** (50 files, 41% coverage) that serves as an excellent reference implementation. However, implementing 10 more content modules (TV, music, books, etc.) with the current approach would result in **massive code duplication** (~70% of code would be duplicated).

This document proposes a **shared abstractions + adapter pattern** architecture that would:
- **Reduce code duplication by 60-70%**
- **Accelerate new module development** from 3-4 weeks to <1 week
- **Maintain module independence** while sharing infrastructure
- **Improve maintainability** through centralized bug fixes

---

## 1. Current Module Structure

### Movie Module (Reference Implementation)

**Location**: `internal\content\movie\`

```
movie/
├── module.go                    # fx dependency injection
├── types.go                     # Domain models (Movie, MovieFile, etc.)
├── repository.go                # Interface definitions
├── repository_postgres.go       # PostgreSQL implementation (526 lines)
├── service.go                   # Business logic (315 lines)
├── cached_service.go            # Cache wrapper
├── handler.go                   # HTTP/API handlers
├── library_service.go           # Library management (315 lines)
├── library_scanner.go           # File system scanning (300 lines)
├── library_matcher.go           # File-to-content matching (385 lines)
├── metadata_service.go          # Metadata fetching (202 lines)
├── tmdb_client.go               # TMDb API client (428 lines)
├── tmdb_mapper.go               # TMDb → Domain mapping (248 lines)
├── tmdb_types.go                # TMDb API types (613 lines)
├── mediainfo.go                 # Media file probing (87 lines)
├── mediainfo_types.go           # MediaInfo types (129 lines)
├── db/                          # sqlc generated code (9 files)
│   ├── models.go
│   ├── movies.sql.go
│   ├── movie_files.sql.go
│   ├── movie_credits.sql.go
│   └── ...
└── moviejobs/                   # River background jobs
    ├── library_scan.go          # Library scanning job
    ├── file_match.go            # File matching job
    ├── metadata_refresh.go      # Metadata refresh job
    └── search_index.go          # Search indexing job
```

**Total**: ~50 files, ~4,500 lines of code

---

### TV Show Module (Skeleton)

**Location**: `internal\content\tvshow\`

**Status**: Only placeholder exists
```
tvshow/
└── db/
    └── placeholder.sql.go       # Minimal placeholder query
```

**Missing**: Everything else (repository, service, handlers, scanner, metadata, jobs)

---

### Expected Pattern (Without Refactoring)

If we implement TV shows the same way as movies:

```
tvshow/
├── types.go                     # TVShow, Season, Episode, TVFile
├── repository.go
├── repository_postgres.go       # 500+ lines (similar to movie)
├── service.go                   # 300+ lines
├── library_service.go           # 300+ lines (70% duplicate)
├── library_scanner.go           # 280+ lines (90% duplicate)
├── library_matcher.go           # 350+ lines (80% duplicate)
├── metadata_service.go          # 180+ lines (70% duplicate)
├── thetvdb_client.go            # 400+ lines (similar to tmdb_client.go)
├── thetvdb_mapper.go            # 250+ lines
├── thetvdb_types.go             # 600+ lines
└── tvjobs/
    ├── library_scan.go          # 90% duplicate
    ├── file_match.go            # 90% duplicate
    ├── metadata_refresh.go      # 90% duplicate
    └── search_index.go          # 90% duplicate
```

**Problem**: ~70% of this code would be duplicated from the movie module!

---

## 2. Code Duplication Analysis

### High Duplication (90%+ Similar)

#### A. Library Scanner (300 lines)

**Shared Logic** (90%):
- Recursive directory traversal with `filepath.Walk`
- File extension filtering (video files: `.mp4`, `.mkv`, `.avi`)
- Progress tracking and context cancellation
- Error aggregation
- Quality marker removal (`1080p`, `BluRay`, `x264`, `HEVC`)
- Release group detection (`SPARKS`, `RARBG`, `YTS`, `YIFY`)

**Content-Specific** (10%):
- Filename parsing:
  - Movies: `Title (Year).ext` → Extract title + year
  - TV: `Series.S01E05.Title.ext` → Extract series + season + episode
  - Music: `Artist - Album - Track.ext` → Extract artist + album + track

**Example Duplication**:
```go
// internal/content/movie/library_scanner.go:45-78
// This exact code would be duplicated in tvshow, music, etc.
func (s *Scanner) walkDirectory(ctx context.Context, path string) ([]ScanResult, error) {
    var results []ScanResult
    err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if info.IsDir() {
            return nil
        }
        if !s.isVideoFile(filePath) {
            return nil
        }
        // ... 30 more lines ...
    })
    return results, err
}
```

---

#### B. Fuzzy Matcher (385 lines)

**Shared Logic** (80%):
- Levenshtein distance calculation (80 lines)
- Fuzzy matching algorithms
- Confidence scoring (thresholds: 0.95 excellent, 0.85 good, 0.70 fair)
- Database lookup patterns
- Title normalization (lowercase, remove "the", remove punctuation)

**Content-Specific** (20%):
- Matching criteria:
  - Movies: Title + Year
  - TV: Series + Season + Episode
  - Music: Artist + Album + Track Number

**Example Duplication**:
```go
// internal/content/movie/library_matcher.go:173-251
// Levenshtein distance - would be duplicated in every module
func levenshteinDistance(a, b string) int {
    // 78 lines of algorithm
    // Identical across all content types
}

func calculateConfidence(scan, candidate string) float64 {
    // 25 lines of scoring logic
    // Identical across all content types
}
```

---

#### C. Background Jobs (River)

**Shared Logic** (90%):
- River worker boilerplate
- Progress logging
- Error handling
- Job argument patterns
- Transaction support

**Content-Specific** (10%):
- Domain models (Movie vs TVShow vs Album)
- Repository calls

**Example Duplication**:
```go
// internal/content/movie/moviejobs/library_scan.go:35-85
// This pattern would be duplicated for every content type
func (w *LibraryScanWorker) Work(ctx context.Context, job *river.Job[LibraryScanArgs]) error {
    logger := riverx.LoggerFromContext(ctx)
    logger.Info("starting library scan", "library_id", job.Args.LibraryID)

    // ... 50 lines of boilerplate ...
}
```

---

### Medium Duplication (70% Similar)

#### D. Metadata Provider (HTTP Client + Mapper)

**Shared Logic** (70%):
- HTTP client setup with rate limiting
- Request/response handling
- Error parsing (TMDb returns same error format as TheTVDB)
- Retry logic
- Cache integration
- Null handling patterns

**Content-Specific** (30%):
- API endpoints (TMDb: `/movie/{id}` vs TheTVDB: `/series/{id}`)
- Authentication methods (API key vs OAuth)
- Data models (Movie vs TVShow)
- Mapping logic (different fields)

**Example Duplication**:
```go
// internal/content/movie/tmdb_client.go:78-145
// HTTP client setup - identical for all providers
func newHTTPClient(apiKey string) *http.Client {
    return &http.Client{
        Timeout: 30 * time.Second,
        Transport: &http.Transport{
            MaxIdleConns:        100,
            MaxIdleConnsPerHost: 10,
            IdleConnTimeout:     90 * time.Second,
        },
    }
}

// Rate limiter - identical for all providers
func (c *Client) rateLimitedRequest(ctx context.Context, url string) (*http.Response, error) {
    // 40 lines of rate limiting logic
}
```

---

#### E. Search Indexing

**Shared Logic** (70%):
- Document conversion
- Bulk indexing operations
- Collection management
- Error handling

**Content-Specific** (30%):
- Schema definition (different fields for movies vs TV)
- Facet definitions
- Search weights

---

### Low Duplication (30% Similar)

#### F. Repository Layer

**Shared Logic** (30%):
- Database conversion utilities (`pgDateToTimePtr`, `pgNumericToDecimalPtr`)
- Error wrapping patterns
- Query parameter construction

**Content-Specific** (70%):
- Domain models (Movie vs TVShow vs Album)
- sqlc-generated queries
- Table schemas

**Note**: Low duplication here is acceptable. Repository pattern provides good separation.

---

## 3. Proposed Architecture

### Package Structure

```
internal/
├── content/
│   ├── shared/                          # NEW: Shared abstractions
│   │   ├── scanner/
│   │   │   ├── scanner.go              # Generic file scanner
│   │   │   ├── parser.go               # File parsing interface
│   │   │   ├── patterns.go             # Quality markers, release groups
│   │   │   └── walker.go               # Directory traversal
│   │   ├── matcher/
│   │   │   ├── matcher.go              # Generic matcher
│   │   │   ├── fuzzy.go                # Levenshtein, scoring algorithms
│   │   │   ├── strategy.go             # Match strategy interface
│   │   │   └── normalizer.go           # Title normalization
│   │   ├── metadata/
│   │   │   ├── provider.go             # Generic HTTP provider
│   │   │   ├── cache.go                # Metadata caching
│   │   │   ├── mapper.go               # Mapping interface
│   │   │   └── client.go               # HTTP client with rate limiting
│   │   ├── library/
│   │   │   ├── service.go              # Generic library manager
│   │   │   ├── types.go                # ScanResult, MatchResult
│   │   │   └── orchestrator.go         # Scan orchestration
│   │   ├── prober/
│   │   │   ├── interface.go            # Media probing interface
│   │   │   ├── video.go                # FFmpeg/mediainfo wrapper
│   │   │   ├── audio.go                # ID3/Vorbis tag reader
│   │   │   └── document.go             # PDF/EPUB metadata
│   │   ├── indexer/
│   │   │   ├── typesense.go            # Generic Typesense indexer
│   │   │   ├── schema.go               # Schema builder
│   │   │   └── document.go             # Document mapper interface
│   │   └── jobs/
│   │       ├── base_worker.go          # Base River worker
│   │       ├── scan_job.go             # Generic scan job
│   │       ├── match_job.go            # Generic match job
│   │       └── metadata_job.go         # Generic metadata job
│   ├── movie/                           # Content-specific implementations
│   │   ├── module.go
│   │   ├── types.go                    # Movie, MovieFile (domain models)
│   │   ├── repository.go
│   │   ├── repository_postgres.go
│   │   ├── service.go
│   │   ├── handler.go
│   │   ├── adapters/                   # NEW: Adapters for shared code
│   │   │   ├── scanner_adapter.go      # Movie-specific parsing
│   │   │   ├── matcher_adapter.go      # Title+year matching
│   │   │   ├── tmdb_adapter.go         # TMDb provider adapter
│   │   │   ├── prober_adapter.go       # Video probing adapter
│   │   │   └── search_adapter.go       # Movie search schema
│   │   ├── db/                         # sqlc generated
│   │   └── jobs/                       # River jobs (use shared base)
│   ├── tvshow/                          # Similar structure
│   │   ├── types.go                    # TVShow, Season, Episode
│   │   ├── repository.go
│   │   ├── repository_postgres.go
│   │   ├── service.go
│   │   ├── handler.go
│   │   ├── adapters/
│   │   │   ├── scanner_adapter.go      # S##E## parsing
│   │   │   ├── matcher_adapter.go      # Series+season+episode
│   │   │   ├── thetvdb_adapter.go      # TheTVDB provider
│   │   │   ├── prober_adapter.go       # Video probing (reuse)
│   │   │   └── search_adapter.go       # TV search schema
│   │   └── ...
│   ├── music/
│   │   ├── adapters/
│   │   │   ├── scanner_adapter.go      # Artist/Album/Track parsing
│   │   │   ├── matcher_adapter.go      # Fuzzy artist/album match
│   │   │   ├── musicbrainz_adapter.go  # MusicBrainz provider
│   │   │   ├── prober_adapter.go       # ID3 tag reader
│   │   │   └── search_adapter.go       # Music search schema
│   │   └── ...
│   └── qar/                             # Adult content
│       ├── adapters/
│       │   ├── scanner_adapter.go      # Scene/performer parsing
│       │   ├── matcher_adapter.go      # Scene ID matching
│       │   ├── stashdb_adapter.go      # StashDB provider
│       │   └── ...
│       └── ...
├── service/                             # Cross-cutting services
│   ├── library/                         # Generic library management
│   ├── search/                          # Search infrastructure
│   └── ...
└── infra/                               # Infrastructure
    ├── cache/
    ├── database/
    └── ...
```

---

## 4. Core Abstractions & Interfaces

### A. ContentItem (Base Abstraction)

```go
// internal/content/shared/types.go
package shared

type ContentItem interface {
    GetID() uuid.UUID
    GetTitle() string
    GetType() ContentType
    GetMetadata() map[string]any
}

type ContentType string

const (
    ContentTypeMovie     ContentType = "movie"
    ContentTypeTVShow    ContentType = "tvshow"
    ContentTypeMusic     ContentType = "music"
    ContentTypeBook      ContentType = "book"
    ContentTypeAudiobook ContentType = "audiobook"
    ContentTypeComic     ContentType = "comic"
    ContentTypeLiveTV    ContentType = "livetv"
    ContentTypeQAR       ContentType = "qar"
)
```

---

### B. Scanner Framework

```go
// internal/content/shared/scanner/scanner.go
package scanner

// FileParser interface - content-specific implementation
type FileParser interface {
    Parse(filename string) (*ParseResult, error)
    GetExtensions() []string
    NormalizeTitle(title string) string
}

type ParseResult struct {
    Title    string
    Metadata map[string]any  // Year, Season, Episode, Artist, Album, etc.
}

// Scanner - generic implementation
type Scanner struct {
    parser FileParser
    config Config
}

type Config struct {
    ExcludePatterns []string
    FollowSymlinks  bool
    IgnoreHidden    bool
}

func (s *Scanner) Scan(ctx context.Context, paths []string) ([]ScanResult, error) {
    // Generic directory traversal
    // Calls parser.Parse() for each file
}

type ScanResult struct {
    Path      string
    FileInfo  os.FileInfo
    ParsedTitle string
    Metadata  map[string]any
}
```

**Movie Adapter**:
```go
// internal/content/movie/adapters/scanner_adapter.go
package adapters

type MovieFileParser struct{}

func (p *MovieFileParser) Parse(filename string) (*scanner.ParseResult, error) {
    // Extract title and year from "Movie Title (2023).mkv"
    title, year := extractTitleYear(filename)
    return &scanner.ParseResult{
        Title: title,
        Metadata: map[string]any{
            "year": year,
        },
    }, nil
}

func (p *MovieFileParser) GetExtensions() []string {
    return []string{".mp4", ".mkv", ".avi", ".m4v"}
}
```

**TV Adapter**:
```go
// internal/content/tvshow/adapters/scanner_adapter.go
package adapters

type TVFileParser struct{}

func (p *TVFileParser) Parse(filename string) (*scanner.ParseResult, error) {
    // Extract series, season, episode from "Series.S01E05.Title.mkv"
    series, season, episode := extractTVInfo(filename)
    return &scanner.ParseResult{
        Title: series,
        Metadata: map[string]any{
            "season":  season,
            "episode": episode,
        },
    }, nil
}
```

---

### C. Matcher Framework

```go
// internal/content/shared/matcher/matcher.go
package matcher

// MatchStrategy interface - content-specific implementation
type MatchStrategy[T shared.ContentItem] interface {
    FindExisting(ctx context.Context, parse scanner.ParseResult) (T, error)
    SearchExternal(ctx context.Context, parse scanner.ParseResult) ([]T, error)
    CalculateConfidence(parse scanner.ParseResult, candidate T) float64
}

// Matcher - generic implementation with shared algorithms
type Matcher[T shared.ContentItem] struct {
    strategy MatchStrategy[T]
}

func (m *Matcher[T]) Match(ctx context.Context, scan scanner.ScanResult) (*MatchResult[T], error) {
    // 1. Try to find existing in database
    existing, err := m.strategy.FindExisting(ctx, scan)
    if err == nil {
        return &MatchResult[T]{
            Item: existing,
            Confidence: 1.0,
            Source: "database",
        }, nil
    }

    // 2. Search external provider
    candidates, err := m.strategy.SearchExternal(ctx, scan)
    if err != nil {
        return nil, err
    }

    // 3. Find best match using confidence scoring
    var best *MatchResult[T]
    for _, candidate := range candidates {
        confidence := m.strategy.CalculateConfidence(scan, candidate)
        if best == nil || confidence > best.Confidence {
            best = &MatchResult[T]{
                Item: candidate,
                Confidence: confidence,
                Source: "external",
            }
        }
    }

    return best, nil
}

type MatchResult[T shared.ContentItem] struct {
    Item       T
    Confidence float64
    Source     string
}
```

**Movie Adapter**:
```go
// internal/content/movie/adapters/matcher_adapter.go
package adapters

type MovieMatchStrategy struct {
    repo     *movie.Repository
    provider metadata.Provider[*movie.Movie]
}

func (s *MovieMatchStrategy) FindExisting(ctx context.Context, parse scanner.ParseResult) (*movie.Movie, error) {
    title := parse.Title
    year := parse.Metadata["year"].(int)
    return s.repo.FindByTitleYear(ctx, title, year)
}

func (s *MovieMatchStrategy) SearchExternal(ctx context.Context, parse scanner.ParseResult) ([]*movie.Movie, error) {
    query := parse.Title
    filters := map[string]any{
        "year": parse.Metadata["year"],
    }
    return s.provider.Search(ctx, query, filters)
}

func (s *MovieMatchStrategy) CalculateConfidence(parse scanner.ParseResult, candidate *movie.Movie) float64 {
    // Use shared fuzzy matching from matcher.Levenshtein()
    titleSimilarity := matcher.LevenshteinNormalized(parse.Title, candidate.Title)

    year := parse.Metadata["year"].(int)
    yearMatch := candidate.Year == year

    if yearMatch {
        return titleSimilarity
    }
    return titleSimilarity * 0.8 // Penalty for year mismatch
}
```

**TV Adapter**:
```go
// internal/content/tvshow/adapters/matcher_adapter.go
package adapters

type TVMatchStrategy struct {
    repo     *tvshow.Repository
    provider metadata.Provider[*tvshow.Episode]
}

func (s *TVMatchStrategy) CalculateConfidence(parse scanner.ParseResult, candidate *tvshow.Episode) float64 {
    seriesSimilarity := matcher.LevenshteinNormalized(parse.Title, candidate.Series.Title)

    season := parse.Metadata["season"].(int)
    episode := parse.Metadata["episode"].(int)

    if candidate.Season == season && candidate.Episode == episode {
        return seriesSimilarity
    }
    return 0.0 // Must match season + episode
}
```

---

### D. Metadata Provider Framework

```go
// internal/content/shared/metadata/provider.go
package metadata

// Provider interface - content-specific implementation
type Provider[T shared.ContentItem] interface {
    Search(ctx context.Context, query string, filters map[string]any) ([]T, error)
    GetByID(ctx context.Context, id any) (T, error)
    Enrich(ctx context.Context, item T) error
}

// HTTPProvider - generic HTTP client with rate limiting
type HTTPProvider[T shared.ContentItem] struct {
    client   *http.Client
    limiter  *rate.Limiter
    cache    *cache.Cache
    baseURL  string
    apiKey   string
    mapper   Mapper[T]
}

type Mapper[T shared.ContentItem] interface {
    MapSearchResult(data any) (T, error)
    MapDetailResult(data any) (T, error)
}

func (p *HTTPProvider[T]) Search(ctx context.Context, query string, filters map[string]any) ([]T, error) {
    // Generic HTTP request with rate limiting
    // Calls mapper.MapSearchResult() for each result
}

func (p *HTTPProvider[T]) GetByID(ctx context.Context, id any) (T, error) {
    // Generic HTTP request with caching
    // Calls mapper.MapDetailResult()
}
```

**Movie Adapter (TMDb)**:
```go
// internal/content/movie/adapters/tmdb_adapter.go
package adapters

type TMDbProvider struct {
    metadata.HTTPProvider[*movie.Movie]
}

func NewTMDbProvider(apiKey string, cache *cache.Cache) *TMDbProvider {
    return &TMDbProvider{
        HTTPProvider: metadata.HTTPProvider[*movie.Movie]{
            BaseURL: "https://api.themoviedb.org/3",
            APIKey:  apiKey,
            Limiter: rate.NewLimiter(rate.Every(250*time.Millisecond), 40),
            Cache:   cache,
            Mapper:  &TMDbMapper{},
        },
    }
}

type TMDbMapper struct{}

func (m *TMDbMapper) MapSearchResult(data any) (*movie.Movie, error) {
    tmdbMovie := data.(TMDbMovie)
    return &movie.Movie{
        TMDbID:   tmdbMovie.ID,
        Title:    tmdbMovie.Title,
        Year:     tmdbMovie.ReleaseDate.Year(),
        Overview: tmdbMovie.Overview,
    }, nil
}
```

**TV Adapter (TheTVDB)**:
```go
// internal/content/tvshow/adapters/thetvdb_adapter.go
package adapters

type TheTVDBProvider struct {
    metadata.HTTPProvider[*tvshow.Series]
}

func NewTheTVDBProvider(apiKey string, cache *cache.Cache) *TheTVDBProvider {
    return &TheTVDBProvider{
        HTTPProvider: metadata.HTTPProvider[*tvshow.Series]{
            BaseURL: "https://api4.thetvdb.com/v4",
            APIKey:  apiKey,
            Limiter: rate.NewLimiter(rate.Every(1*time.Second), 10),
            Cache:   cache,
            Mapper:  &TheTVDBMapper{},
        },
    }
}
```

---

### E. Library Service Framework

```go
// internal/content/shared/library/service.go
package library

// LibraryManager interface - implemented by BaseLibraryService
type LibraryManager[T shared.ContentItem] interface {
    ScanLibrary(ctx context.Context, libraryID uuid.UUID) (*ScanSummary, error)
    MatchFile(ctx context.Context, path string) (*MatchResult[T], error)
    RefreshMetadata(ctx context.Context, id uuid.UUID) error
}

// BaseLibraryService - generic implementation
type BaseLibraryService[T shared.ContentItem] struct {
    repo     Repository[T]
    scanner  *scanner.Scanner
    matcher  *matcher.Matcher[T]
    provider metadata.Provider[T]
    jobs     *jobs.JobQueue
}

type Repository[T shared.ContentItem] interface {
    Create(ctx context.Context, item T) (T, error)
    Update(ctx context.Context, item T) error
    FindByID(ctx context.Context, id uuid.UUID) (T, error)
}

func (s *BaseLibraryService[T]) ScanLibrary(ctx context.Context, libraryID uuid.UUID) (*ScanSummary, error) {
    // 1. Get library paths
    lib, err := s.libraryRepo.GetLibrary(ctx, libraryID)
    if err != nil {
        return nil, err
    }

    // 2. Scan filesystem
    scanResults, err := s.scanner.Scan(ctx, lib.Paths)
    if err != nil {
        return nil, err
    }

    // 3. Queue match jobs for each file
    for _, result := range scanResults {
        s.jobs.EnqueueMatch(ctx, result)
    }

    return &ScanSummary{
        FilesFound: len(scanResults),
        Status:     "queued",
    }, nil
}
```

**Movie Implementation**:
```go
// internal/content/movie/library_service.go
package movie

type LibraryService struct {
    library.BaseLibraryService[*Movie]
    prober Prober  // FFmpeg for video metadata
}

func NewLibraryService(
    repo *Repository,
    scanner *scanner.Scanner,
    matcher *matcher.Matcher[*Movie],
    provider metadata.Provider[*Movie],
    prober Prober,
) *LibraryService {
    return &LibraryService{
        BaseLibraryService: library.BaseLibraryService[*Movie]{
            Repo:     repo,
            Scanner:  scanner,
            Matcher:  matcher,
            Provider: provider,
        },
        prober: prober,
    }
}

// Override or extend base methods as needed
func (s *LibraryService) MatchFile(ctx context.Context, path string) (*library.MatchResult[*Movie], error) {
    // 1. Call base implementation
    result, err := s.BaseLibraryService.MatchFile(ctx, path)
    if err != nil {
        return nil, err
    }

    // 2. Add movie-specific logic: probe video file
    mediaInfo, err := s.prober.Probe(ctx, path)
    if err != nil {
        return nil, err
    }

    result.Item.Duration = mediaInfo.Duration
    result.Item.VideoCodec = mediaInfo.VideoCodec
    result.Item.AudioCodec = mediaInfo.AudioCodec

    return result, nil
}
```

---

## 5. Benefits of Proposed Architecture

### Code Reuse Metrics

| Component | Current (Movie) | Proposed (Shared) | Reuse % | Savings |
|-----------|-----------------|-------------------|---------|---------|
| Scanner | 300 lines/module | 250 lines shared + 50 lines adapter | 83% | 250 lines × 10 modules = 2,500 lines |
| Matcher | 385 lines/module | 300 lines shared + 85 lines adapter | 78% | 300 lines × 10 modules = 3,000 lines |
| Metadata Provider | 428 lines/module | 200 lines shared + 228 lines adapter | 47% | 200 lines × 10 modules = 2,000 lines |
| Library Service | 315 lines/module | 200 lines shared + 115 lines adapter | 63% | 200 lines × 10 modules = 2,000 lines |
| Background Jobs | 150 lines/module | 100 lines shared + 50 lines adapter | 67% | 100 lines × 10 modules = 1,000 lines |
| **Total** | **~1,578 lines/module** | **~1,050 shared + 528 adapter** | **67%** | **10,500 lines saved** |

**With 11 content modules** (movie, TV, music, audiobook, book, podcast, photo, comic, livetv, qar):
- **Without refactoring**: 1,578 lines × 11 = **17,358 lines**
- **With refactoring**: 1,050 shared + (528 × 11) = **6,858 lines**
- **Savings**: 10,500 lines (60% reduction)

---

### Development Time Reduction

| Metric | Current | Proposed | Improvement |
|--------|---------|----------|-------------|
| Time to implement new module | 3-4 weeks | <1 week | 75% faster |
| Test coverage effort | Full suite per module | Test shared once + adapter tests | 60% less |
| Bug fixes | Fix in each module | Fix in shared code, all modules benefit | N/A |
| Onboarding time | Understand each module | Understand shared abstractions once | 70% faster |

---

### Consistency Benefits

1. **Uniform Behavior**: All content types scan, match, and fetch metadata the same way
2. **Common Patterns**: Developers see familiar code across modules
3. **Centralized Improvements**: Performance optimizations benefit all modules
4. **Easier Testing**: Shared components tested once, trust across all modules

---

### Maintainability Benefits

1. **Single Source of Truth**: Bug fixes in shared code fix all modules
2. **Reduced Surface Area**: Less code = fewer bugs
3. **Clear Boundaries**: Adapter interfaces make content-specific logic explicit
4. **Easy Refactoring**: Change scanner implementation once, all modules benefit

---

## 6. Migration Strategy

### Phase 1: Extract Utilities (Week 1-2)

**Goal**: Extract non-controversial shared utilities

**Tasks**:
1. Create `internal/content/shared/` package
2. Extract scanner utilities:
   - Quality markers (1080p, BluRay, etc.)
   - Release groups (SPARKS, RARBG, etc.)
   - Directory traversal logic
3. Extract fuzzy matching:
   - Levenshtein distance
   - Confidence scoring
   - Title normalization
4. Update movie module to use shared utilities
5. Ensure all tests pass

**Verification**:
- Movie module tests: 100% passing
- No functionality regression
- Code coverage maintained

---

### Phase 2: Generalize Interfaces (Week 3-4)

**Goal**: Define adapter interfaces and refactor movie module

**Tasks**:
1. Define `FileParser` interface
2. Create `MovieFileParser` adapter (refactor from `library_scanner.go`)
3. Define `MatchStrategy` interface
4. Create `MovieMatchStrategy` adapter (refactor from `library_matcher.go`)
5. Define `Provider` interface
6. Create `TMDbProvider` adapter (refactor from `tmdb_client.go`)
7. Update movie module to use adapters
8. Ensure all tests pass

**Verification**:
- Movie module tests: 100% passing
- Same functionality, new abstractions
- Code coverage maintained

---

### Phase 3: Implement TV Module (Week 5-6)

**Goal**: Validate architecture by implementing TV shows

**Tasks**:
1. Create `TVFileParser` (S##E## parsing)
2. Create `TVMatchStrategy` (series+season+episode matching)
3. Create `TheTVDBProvider` adapter
4. Implement TV repository (PostgreSQL)
5. Implement TV service (using `BaseLibraryService`)
6. Add TV jobs (scan, match, metadata refresh, search index)
7. Add TV API handlers
8. Write tests (reuse shared test patterns)

**Verification**:
- TV module fully functional
- Reused 60%+ of movie code
- Development time: <1 week (vs 3-4 weeks without refactoring)

---

### Phase 4: Implement Music Module (Week 7-8)

**Goal**: Further validate with different content type

**Tasks**:
1. Create `MusicFileParser` (artist/album/track)
2. Create `MusicMatchStrategy` (fuzzy artist/album)
3. Create `MusicBrainzProvider` + `LastFmProvider` adapters
4. Create `ID3Prober` adapter (audio tags)
5. Implement music repository
6. Implement music service (using `BaseLibraryService`)
7. Add music jobs
8. Add music API handlers
9. Write tests

**Verification**:
- Music module fully functional
- Adapters handle different content type well
- Shared code unchanged

---

### Phase 5: Generalize Library Service (Week 9-10)

**Goal**: Extract more common patterns discovered during TV/music implementation

**Tasks**:
1. Extract common library orchestration
2. Create `BaseLibraryService[T]` generic
3. Refactor movie, TV, music to use base service
4. Extract search indexer framework
5. Consolidate background job patterns

**Verification**:
- All modules refactored
- Code reduction: 60%+
- No functionality regression

---

## 7. Dependency Injection Pattern

### Module Registration (fx)

```go
// internal/content/movie/module.go
package movie

var Module = fx.Module("movie",
    fx.Provide(
        // Shared components with adapters
        fx.Annotate(
            NewMovieFileParser,
            fx.As(new(scanner.FileParser)),
        ),
        fx.Annotate(
            NewMovieMatchStrategy,
            fx.As(new(matcher.MatchStrategy[*Movie])),
        ),
        fx.Annotate(
            NewTMDbProvider,
            fx.As(new(metadata.Provider[*Movie])),
        ),
        fx.Annotate(
            NewVideoProber,
            fx.As(new(prober.Prober)),
        ),

        // Content-specific components
        NewPostgresRepository,
        NewService,
        NewLibraryService,  // Uses adapters via interfaces
        NewHandler,

        // Background jobs
        NewLibraryScanWorker,
        NewFileMatchWorker,
        NewMetadataRefreshWorker,
        NewSearchIndexWorker,
    ),
)
```

**Benefits**:
- Adapters injected as interfaces
- Content-specific code isolated
- Easy to test with mocks

---

## 8. Testing Strategy

### Shared Component Tests

Test shared components once with generic test suites:

```go
// internal/content/shared/matcher/matcher_test.go
package matcher_test

func TestLevenshteinDistance(t *testing.T) {
    tests := []struct {
        a, b     string
        expected int
    }{
        {"kitten", "sitting", 3},
        {"saturday", "sunday", 3},
        // ... more test cases
    }

    for _, tt := range tests {
        t.Run(fmt.Sprintf("%s-%s", tt.a, tt.b), func(t *testing.T) {
            got := matcher.LevenshteinDistance(tt.a, tt.b)
            assert.Equal(t, tt.expected, got)
        })
    }
}
```

---

### Adapter Tests

Test adapters with content-specific test cases:

```go
// internal/content/movie/adapters/scanner_adapter_test.go
package adapters_test

func TestMovieFileParser_Parse(t *testing.T) {
    parser := adapters.NewMovieFileParser()

    tests := []struct {
        filename      string
        expectedTitle string
        expectedYear  int
    }{
        {"The Matrix (1999).mkv", "The Matrix", 1999},
        {"Inception.2010.1080p.BluRay.x264.mkv", "Inception", 2010},
        {"Interstellar (2014) [1080p].mp4", "Interstellar", 2014},
    }

    for _, tt := range tests {
        t.Run(tt.filename, func(t *testing.T) {
            result, err := parser.Parse(tt.filename)
            require.NoError(t, err)
            assert.Equal(t, tt.expectedTitle, result.Title)
            assert.Equal(t, tt.expectedYear, result.Metadata["year"])
        })
    }
}
```

---

### Integration Tests

Test full workflow with real implementations:

```go
// internal/content/movie/library_service_test.go
func TestLibraryService_ScanLibrary(t *testing.T) {
    // Setup: real scanner + mock repository + mock provider
    scanner := scanner.NewScanner(adapters.NewMovieFileParser(), scanner.Config{})
    mockRepo := &MockRepository{}
    mockProvider := &MockTMDbProvider{}

    service := movie.NewLibraryService(mockRepo, scanner, ..., mockProvider, ...)

    // Test
    summary, err := service.ScanLibrary(ctx, libraryID)

    // Verify
    require.NoError(t, err)
    assert.Greater(t, summary.FilesFound, 0)
}
```

---

## 9. Trade-offs & Considerations

### Pros

✅ **Massive code reduction** (60-70%)
✅ **Faster development** (75% faster for new modules)
✅ **Consistent behavior** across all content types
✅ **Centralized bug fixes** benefit all modules
✅ **Easier testing** (test shared components once)
✅ **Better documentation** (understand abstractions once)

### Cons

⚠️ **More abstraction layers** (learning curve)
⚠️ **Generic syntax** can be verbose (Go generics)
⚠️ **Migration effort** (4-6 weeks initial investment)
⚠️ **Potential over-engineering** (if only 2-3 modules planned)

### Mitigations

1. **Documentation**: Comprehensive examples for each adapter type
2. **Type aliases**: Reduce generic syntax verbosity
3. **Phased migration**: Validate architecture with TV module before full rollout
4. **Clear boundaries**: Explicit adapter interfaces prevent confusion

---

## 10. Success Metrics

### Code Metrics

- [ ] Scanner code reduction: 60%+ (target: 2,500 lines saved across 11 modules)
- [ ] Matcher code reduction: 70%+ (target: 3,000 lines saved)
- [ ] Metadata provider reduction: 50%+ (target: 2,000 lines saved)
- [ ] Total code reduction: 60%+ (target: 10,500 lines saved)

### Quality Metrics

- [ ] Test coverage maintained: 80%+ (same as before)
- [ ] No performance regression: <5% overhead from abstractions
- [ ] Bug fix propagation: Shared code fix benefits all modules

### Development Metrics

- [ ] New module implementation time: <1 week (vs 3-4 weeks)
- [ ] Developer onboarding time: 70% faster
- [ ] Code review time: 50% faster (reviewers understand shared patterns)

---

## Conclusion

The proposed architecture provides a **clear path forward** for implementing the remaining 10 content modules efficiently. By extracting shared abstractions and using adapter patterns, we can:

1. **Reduce code by 60-70%** (10,500+ lines saved)
2. **Accelerate development by 75%** (new modules in <1 week)
3. **Improve maintainability** through centralized bug fixes
4. **Maintain module independence** via clear adapter boundaries

**Recommendation**: Begin with Phase 1 (extract utilities) and Phase 2 (define interfaces) to validate the approach without disrupting existing movie module functionality. Once proven, proceed with TV module implementation to validate end-to-end workflow.

**Total effort**: 10 weeks (2.5 months)
- Weeks 1-4: Refactor movie module
- Weeks 5-6: Implement TV module (validation)
- Weeks 7-8: Implement music module (validation)
- Weeks 9-10: Generalize remaining patterns

**ROI**: After 10 weeks of investment, each new content module takes <1 week vs 3-4 weeks, saving ~30 weeks across remaining modules.
