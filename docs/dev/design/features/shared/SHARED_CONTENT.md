# Shared Content Packages

<!-- DESIGN: features/shared -->

**Package**: `internal/content/shared`

> Reusable frameworks for content modules: filesystem scanning, content matching, metadata clients, library interfaces, and job utilities

---

## Package Structure

```
internal/content/shared/
├── scanner/               # Filesystem scanning framework
│   ├── scanner.go         # FilesystemScanner + FileParser interface
│   ├── extensions.go      # Video/audio/image/subtitle extension maps
│   └── patterns.go        # Quality markers, release groups, title cleaning
├── matcher/               # Generic content matching
│   ├── matcher.go         # Matcher[T] + MatchStrategy[T] interface
│   └── fuzzy.go           # Levenshtein distance, title similarity, ConfidenceScore
├── metadata/              # Metadata provider utilities
│   ├── client.go          # BaseClient (rate limiter + L1Cache/otter + req HTTP)
│   ├── types.go           # SearchResult, SearchOptions, ClientConfig, shared types
│   ├── images.go          # ImageURLBuilder, ImageDownloader, TMDb image sizes
│   └── maputil.go         # Date parsing, optional values, age ratings, language conversion
├── library/               # Library operation interfaces
│   └── types.go           # LibraryScanner, ContentMatcher[T], MediaProber, ScanSummary
└── jobs/                  # Shared job types
    └── types.go           # JobResult, JobContext, arg types, action constants
```

## Scanner Package

Framework for discovering media files on the filesystem.

### Interfaces

```go
type FileParser interface {
    Parse(filename string) (title string, metadata map[string]any)
    GetExtensions() []string
    ContentType() string  // "movie", "tvshow"
}
```

Implementations: `MovieFileParser` (adapters/scanner_adapter.go), `TVShowFileParser` (adapters/scanner_adapter.go)

### FilesystemScanner

```go
type FilesystemScanner struct {
    paths      []string
    parser     FileParser
    options    ScanOptions
    extensions map[string]bool
}
```

Methods: `Scan(ctx)`, `ScanWithSummary(ctx)`, `ScanPath(ctx, path)`

Returns `[]ScanResult` with: FilePath, FileName, FileSize, ParsedTitle, Metadata (content-specific), IsMedia, Error.

### ScanOptions

```go
DefaultScanOptions() → {FollowSymlinks: false, MaxDepth: 0,
    ExcludePatterns: [".Trash*", ".recycle*", "@eaDir", ".DS_Store"],
    IncludeHidden: false}
```

### Supported Extensions

| Category | Count | Examples |
|----------|-------|---------|
| Video | 15 | .mp4, .mkv, .avi, .mov, .wmv, .ts, .m2ts |
| Audio | 13 | .mp3, .flac, .wav, .aac, .ogg, .opus |
| Image | 8 | .jpg, .png, .gif, .webp, .svg |
| Subtitle | 9 | .srt, .ass, .ssa, .sub, .vtt, .pgs |

Helper functions: `IsVideoFile()`, `IsAudioFile()`, `IsImageFile()`, `IsSubtitleFile()`, `MergeExtensions()`

### Pattern Cleaning

80+ quality markers removed during title cleaning (resolutions, sources, codecs, HDR, release types). 50+ release group names. 12 streaming service markers (AMZN, NF, DSNP, etc.).

Key functions:
- `CleanTitle(title)` - Remove markers, replace dots/underscores, normalize whitespace
- `NormalizeTitle(title)` - Lowercase, remove articles ("the"/"a"/"an"), remove punctuation
- `ExtractYear(text)` - Regex: `\b(19\d{2}|20\d{2})\b`
- `ExtractResolution(text)` - Find 2160p/1080p/720p/480p/4K/UHD
- `ExtractSource(text)` - Map "bluray" to "BluRay", "web-dl" to "WEB-DL", etc.

## Matcher Package

Generic content matching framework using Go generics.

### Interfaces

```go
type MatchStrategy[T any] interface {
    FindExisting(ctx, scanResult) (*T, float64, error)
    SearchExternal(ctx, scanResult) ([]*T, error)
    CalculateConfidence(scanResult, candidate *T) float64
    CreateContent(ctx, candidate *T) (*T, error)
}
```

### Matcher[T]

```go
type Matcher[T any] struct { strategy MatchStrategy[T] }
```

Match logic:
1. `FindExisting` with confidence >= 0.8 → `MatchTypeTitle`
2. `SearchExternal` providers
3. Calculate confidence: < 0.7 → `MatchTypeFuzzy`, else → `MatchTypeTitle`
4. `CreateContent` from top candidate

### Match Types

`exact` (ID match), `title` (title+year, >=0.8), `fuzzy` (similarity, >=0.5), `manual` (user-matched), `unmatched` (no match)

### Fuzzy Matching

- `LevenshteinDistance(s1, s2)` - Edit distance (insertions, deletions, substitutions)
- `TitleSimilarity(t1, t2)` - Normalized after removing articles and punctuation
- `YearMatch(y1, y2)` - Exact = 1.0, +/-1 year = 0.5, else = 0.0

### ConfidenceScore Builder

```go
score := NewConfidenceScore().
    Add(titleSimilarity, 0.6).  // 60% weight
    Add(yearMatch, 0.4).        // 40% weight
    AddBonus(popularityBonus).   // Extra bonus
    Calculate()                  // → [0.0, 1.0]
```

## Metadata Package

Shared HTTP client infrastructure for external metadata providers (TMDb).

### BaseClient

```go
type BaseClient struct {
    client      *req.Client     // imroc/req with retries
    apiKey      string
    rateLimiter *rate.Limiter   // golang.org/x/time/rate
    cache       *cache.L1Cache[string, any]  // otter W-TinyLFU, bounded, TTL-based
    cacheTTL    time.Duration   // Default 24h
    baseURL     string
}
```

Methods: `WaitForRateLimit(ctx)`, `GetFromCache(key)`, `SetCache(key, data)`, `SetCacheWithTTL(key, data, ttl)`, `ClearCache()`, `Request()`, `GetClient()`

### ClientConfig

```go
type ClientConfig struct {
    BaseURL    string          // API base URL
    APIKey     string          // Auth key
    RateLimit  rate.Limit      // Default 4.0 req/sec
    RateBurst  int             // Default 10
    CacheTTL   time.Duration   // Default 24h
    Timeout    time.Duration   // Default 30s
    RetryCount int             // Default 3
    ProxyURL   string          // Optional HTTP proxy
}
```

### Image Handling

`ImageURLBuilder` constructs TMDb image URLs:

| Method | Default Size |
|--------|-------------|
| GetPosterURL | w500 |
| GetBackdropURL | w1280 |
| GetProfileURL | w185 |
| GetLogoURL | w154 |
| GetStillURL | w300 |

`ImageDownloader` downloads images via BaseClient with rate limiting.

### Utilities

- Age rating systems: MPAA, FSK, BBFC, CNC, Eirin, KMRB, DJCTQ, ACB
- Language conversion: `LanguageToISO("en-US") → "en"`, `ISOToLanguage("en") → "en-US"`
- Date parsing: `ParseReleaseDate()`, `ExtractYearFromDate()`
- Safe type conversion: `SafeIntToInt32()`, `ParseOptionalString()`, etc.

## Library Package

Interface definitions for library operations (implemented by content modules):

```go
type LibraryScanner interface {
    Scan(ctx) ([]ScanItem, error)
}

type ContentMatcher[T any] interface {
    MatchFile(ctx, item) MatchResult[T]
    MatchFiles(ctx, items) []MatchResult[T]
}

type MediaProber interface {
    Probe(filePath) (*MediaFileInfo, error)
}

type ContentFileRepository[T any] interface {
    GetFileByPath(ctx, path) (*T, error)
    CreateFile(ctx, file) (*T, error)
    UpdateFile(ctx, file) (*T, error)
    DeleteFile(ctx, id) error
}
```

`MediaFileInfo` contains: Path, Size, Container, Resolution, VideoCodec, AudioCodec, BitrateKbps, DurationSeconds, Framerate, DynamicRange, ColorSpace, AudioChannels, Languages, SubtitleLangs.

## Jobs Package

Shared types for River background job workers:

### JobResult

```go
type JobResult struct {
    Success, ItemsProcessed, ItemsFailed int
    Duration time.Duration
    Errors   []error
    Message  string
}
```

Methods: `AddError(err)`, `HasErrors()`, `LogSummary(logger, jobKind)`, `LogErrors(logger, maxErrors)`

### JobContext

Wraps `context.Context` with logger, job ID, kind, and start time. Methods: `Elapsed()`, `LogStart()`, `LogComplete()`, `LogError()`.

### Shared Arg Types

| Type | Fields |
|------|--------|
| LibraryScanArgs | Paths, Force |
| FileMatchArgs | FilePath, ForceRematch |
| MetadataRefreshArgs | ContentID, Force |
| SearchIndexArgs | ContentID, FullReindex |

### JobKind Function

`JobKind("movie", "library_scan") → "movie_library_scan"`

Action constants: `library_scan`, `file_match`, `metadata_refresh`, `search_index`, `media_probe`

## Dependencies

- `github.com/imroc/req/v3` - HTTP client (metadata package)
- `golang.org/x/time/rate` - Rate limiting (metadata package)
- `github.com/google/uuid` - UUID types (library, jobs packages)
- `go.uber.org/zap` - Logging (jobs package)

## Related Documentation

- [../video/MOVIE_MODULE.md](../video/MOVIE_MODULE.md) - Movie module (uses all shared packages)
- [../video/TVSHOW_MODULE.md](../video/TVSHOW_MODULE.md) - TV show module (uses all shared packages)
- [../../architecture/METADATA_SYSTEM.md](../../architecture/METADATA_SYSTEM.md) - Metadata provider chain
- [../../infrastructure/JOBS.md](../../infrastructure/JOBS.md) - River job queue infrastructure
