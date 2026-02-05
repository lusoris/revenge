# TODO A10: Shared Abstractions & Code Deduplication

**Phase**: A10
**Priority**: P1 (High - enables faster TV development)
**Effort**: 40-60 hours
**Status**: Pending
**Dependencies**: A9 (Multi-Language)
**Created**: 2026-02-05

---

## Overview

**Goal**: Extract shared functionality from Movie module to enable fast TV module implementation.

**Problem**: Without deduplication, each new content module requires ~3-4 weeks and duplicates 60-70% of code.

**Solution**: Create shared abstractions with adapter pattern - modules stay independent, tables stay separate, but code is shared.

**Impact**:
- 60-70% code reduction
- New modules in <1 week (vs 3-4 weeks)
- Centralized bug fixes benefit all modules

**Source**: [REPORT_4_MODULE_ARCHITECTURE.md](REPORT_4_MODULE_ARCHITECTURE.md)

---

## Architecture Principles

**✅ SHARE**:
- Algorithms (Levenshtein, fuzzy matching)
- File scanner framework (with adapters)
- HTTP provider framework (with adapters)
- Background job boilerplate
- Library orchestration patterns

**❌ DON'T SHARE** (would create monolith):
- Domain models (Movie, TVShow stay separate)
- Repositories (each module has own)
- Database tables (separate schemas/tables per module)
- API endpoints (each module has own handlers)

---

## Tasks

### A10.1: Scanner Framework 🔴 CRITICAL

**Priority**: P0
**Effort**: 12-16h
**Location**: `internal/content/shared/scanner/`

**Create**:
```
internal/content/shared/scanner/
├── scanner.go       # Generic directory walker
├── parser.go        # FileParser interface
├── patterns.go      # Quality markers, release groups
└── walker.go        # Filesystem traversal
```

**Interfaces**:
```go
type FileParser interface {
    Parse(filename string) (*ParseResult, error)
    GetExtensions() []string
}

type ParseResult struct {
    Title    string
    Metadata map[string]any  // Flexible for year, season, episode, etc.
}
```

**Adapters**:
- `internal/content/movie/adapters/scanner_adapter.go` - Parse "Title (Year).mkv"
- `internal/content/tvshow/adapters/scanner_adapter.go` - Parse "Series.S01E05.mkv"

**Subtasks**:
- [ ] Extract directory traversal from movie scanner
- [ ] Extract quality/release group patterns
- [ ] Define FileParser interface
- [ ] Create MovieFileParser adapter
- [ ] Refactor movie module to use shared scanner
- [ ] Write tests (90% coverage target)

---

### A10.2: Matcher Framework 🔴 CRITICAL

**Priority**: P0
**Effort**: 12-16h
**Location**: `internal/content/shared/matcher/`

**Create**:
```
internal/content/shared/matcher/
├── matcher.go       # Generic matcher with strategy pattern
├── fuzzy.go         # Levenshtein distance, confidence scoring
├── strategy.go      # MatchStrategy interface
└── normalizer.go    # Title normalization
```

**Interfaces**:
```go
type MatchStrategy[T ContentItem] interface {
    FindExisting(ctx context.Context, parse ParseResult) (T, error)
    SearchExternal(ctx context.Context, parse ParseResult) ([]T, error)
    CalculateConfidence(parse ParseResult, candidate T) float64
}
```

**Adapters**:
- `internal/content/movie/adapters/matcher_adapter.go` - Title + Year matching
- `internal/content/tvshow/adapters/matcher_adapter.go` - Series + S##E## matching

**Subtasks**:
- [ ] Extract Levenshtein from movie matcher
- [ ] Extract confidence scoring
- [ ] Extract title normalization
- [ ] Define MatchStrategy interface
- [ ] Create MovieMatchStrategy adapter
- [ ] Refactor movie module to use shared matcher
- [ ] Write tests

---

### A10.3: Metadata Provider Framework 🟠 HIGH

**Priority**: P1
**Effort**: 12-16h
**Location**: `internal/content/shared/metadata/`

**Create**:
```
internal/content/shared/metadata/
├── provider.go      # Provider interface + HTTPProvider base
├── cache.go         # Metadata caching wrapper
├── mapper.go        # Mapper interface
└── client.go        # HTTP client with rate limiting
```

**Interfaces**:
```go
type Provider[T ContentItem] interface {
    Search(ctx context.Context, query string, filters map[string]any) ([]T, error)
    GetByID(ctx context.Context, id any) (T, error)
    Enrich(ctx context.Context, item T) error
}

type HTTPProvider[T ContentItem] struct {
    client   *http.Client
    limiter  *rate.Limiter
    cache    *cache.Cache
    baseURL  string
    apiKey   string
    mapper   Mapper[T]
}
```

**Adapters**:
- `internal/content/movie/adapters/tmdb_adapter.go` - TMDb movies
- `internal/content/tvshow/adapters/tmdb_tv_adapter.go` - TMDb TV shows

**Subtasks**:
- [ ] Extract HTTP client setup from tmdb_client.go
- [ ] Extract rate limiting logic
- [ ] Define Provider interface
- [ ] Create HTTPProvider base class
- [ ] Create TMDbMovieProvider adapter
- [ ] Refactor movie module to use shared provider
- [ ] Write tests

---

### A10.4: Library Service Framework 🟡 MEDIUM

**Priority**: P2
**Effort**: 8-12h
**Location**: `internal/content/shared/library/`

**Create**:
```
internal/content/shared/library/
├── service.go       # BaseLibraryService[T]
├── types.go         # ScanResult, MatchResult
└── orchestrator.go  # Scan orchestration logic
```

**Generic Service**:
```go
type BaseLibraryService[T ContentItem] struct {
    repo     Repository[T]
    scanner  *scanner.Scanner
    matcher  *matcher.Matcher[T]
    provider metadata.Provider[T]
    jobs     *jobs.JobQueue
}

func (s *BaseLibraryService[T]) ScanLibrary(ctx context.Context, libraryID uuid.UUID) (*ScanSummary, error) {
    // Generic scan orchestration
    // Modules override/extend as needed
}
```

**Subtasks**:
- [ ] Extract common library patterns
- [ ] Define Repository[T] interface
- [ ] Create BaseLibraryService
- [ ] Refactor movie library service to extend base
- [ ] Write tests

---

### A10.5: Background Jobs Framework 🟡 MEDIUM

**Priority**: P2
**Effort**: 6-8h
**Location**: `internal/content/shared/jobs/`

**Create**:
```
internal/content/shared/jobs/
├── base_worker.go   # Base River worker
├── scan_job.go      # Generic scan job
└── metadata_job.go  # Generic metadata refresh job
```

**Base Worker**:
```go
type BaseWorker[T ContentItem] struct {
    logger *slog.Logger
    // Common fields
}

func (w *BaseWorker[T]) Work(ctx context.Context, job *river.Job[Args]) error {
    // Common job patterns
}
```

**Subtasks**:
- [ ] Extract common job patterns
- [ ] Create base worker types
- [ ] Refactor movie jobs to extend base
- [ ] Write tests

---

## Migration Strategy

### Week 1-2: Extract & Define
1. Create `internal/content/shared/` package structure
2. Extract scanner, matcher, provider frameworks
3. Define all interfaces
4. Write comprehensive tests for shared code

### Week 3-4: Refactor Movie Module
5. Create movie adapters (scanner, matcher, provider)
6. Refactor movie module to use shared code + adapters
7. Ensure all movie tests still pass
8. Verify no functionality regression

### Week 5-6: Prepare for TV
9. Document adapter pattern
10. Create example templates for TV adapters
11. Validate that TV implementation is <1 week effort

---

## Code Savings Estimate

| Component | Before (per module) | After (shared + adapter) | Savings |
|-----------|---------------------|--------------------------|---------|
| Scanner | 300 lines | 250 shared + 50 adapter | 83% |
| Matcher | 385 lines | 300 shared + 85 adapter | 78% |
| Provider | 428 lines | 200 shared + 228 adapter | 47% |
| Library | 315 lines | 200 shared + 115 adapter | 63% |
| Jobs | 150 lines | 100 shared + 50 adapter | 67% |
| **Total** | **1,578 lines** | **1,050 shared + 528 adapter** | **67%** |

**With 11 content modules**: Save ~10,500 lines of code!

---

## Testing Requirements

- [ ] Shared scanner: 90% coverage
- [ ] Shared matcher: 90% coverage
- [ ] Shared provider: 85% coverage
- [ ] All movie tests still pass
- [ ] No performance regression (<5% overhead)
- [ ] Integration tests with adapters

---

## Documentation

- [ ] Architecture decision record (ADR)
- [ ] Adapter pattern guide
- [ ] How to create new content module
- [ ] Code examples for each adapter type
- [ ] Migration notes for future modules

---

## Verification Checklist

- [ ] All shared code extracted
- [ ] Movie module refactored to use adapters
- [ ] Movie functionality unchanged
- [ ] Movie tests passing (100%)
- [ ] Code duplication reduced (measure with tools)
- [ ] Documentation complete
- [ ] Ready for TV module implementation

---

**Completion Criteria**:
✅ Shared abstractions complete and tested
✅ Movie module refactored successfully
✅ No functionality regression
✅ 60%+ code reduction measured
✅ TV module can be implemented in <1 week
