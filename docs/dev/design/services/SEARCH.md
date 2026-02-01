## Table of Contents

- [Search Service](#search-service)
  - [Status](#status)
  - [Architecture](#architecture)
    - [Service Structure](#service-structure)
    - [Dependencies](#dependencies)
    - [Provides](#provides)
    - [Component Diagram](#component-diagram)
  - [Implementation](#implementation)
    - [File Structure](#file-structure)
    - [Key Interfaces](#key-interfaces)
    - [Dependencies](#dependencies)
  - [Configuration](#configuration)
    - [Environment Variables](#environment-variables)
    - [Config Keys](#config-keys)
  - [API Endpoints](#api-endpoints)
- [Search](#search)
- [Admin (indexing)](#admin-indexing)
  - [Testing Strategy](#testing-strategy)
    - [Unit Tests](#unit-tests)
    - [Integration Tests](#integration-tests)
    - [Test Coverage](#test-coverage)
  - [Related Documentation](#related-documentation)
    - [Design Documents](#design-documents)
    - [External Sources](#external-sources)

# Search Service


**Created**: 2026-01-31
**Status**: ✅ Complete
**Category**: service


> > Full-text search via Typesense with per-module collections

**Package**: `internal/service/search`
**fx Module**: `search.Module`

---


## Status

| Dimension | Status | Notes |
|-----------|--------|-------|
| Design | ✅ | - |
| Sources | ✅ | - |
| Instructions | ✅ | - |
| Code | 🔴 | - |
| Linting | 🔴 | - |
| Unit Testing | 🔴 | - |
| Integration Testing | 🔴 | - |

**Overall**: ✅ Complete



---


## Architecture

### Service Structure

```
internal/service/search/
├── module.go              # fx module definition
├── service.go             # Service implementation
├── repository.go          # Data access (if needed)
├── handler.go             # HTTP handlers (if exposed)
├── middleware.go          # Middleware (if needed)
├── types.go               # Domain types
└── service_test.go        # Tests
```

### Dependencies
**Go Packages**:
- `github.com/google/uuid`
- `github.com/jackc/pgx/v5`
- `github.com/typesense/typesense-go/typesense` - Typesense client
- `github.com/riverqueue/river` - Background sync jobs
- `go.uber.org/fx`

**External Services**:
- Typesense server (https://typesense.org/)


### Provides
<!-- Service provides -->

### Component Diagram

<!-- Component diagram -->


## Implementation

### File Structure

<!-- File structure -->

### Key Interfaces

```go
type SearchService interface {
  // Search
  Search(ctx context.Context, query SearchQuery) (*SearchResults, error)
  MultiSearch(ctx context.Context, queries []SearchQuery) ([]SearchResults, error)

  // Indexing
  IndexDocument(ctx context.Context, collection string, document interface{}) error
  UpdateDocument(ctx context.Context, collection, id string, document interface{}) error
  DeleteDocument(ctx context.Context, collection, id string) error
  BulkIndex(ctx context.Context, collection string, documents []interface{}) error

  // Collections
  CreateCollection(ctx context.Context, schema CollectionSchema) error
  DeleteCollection(ctx context.Context, name string) error
  ListCollections(ctx context.Context) ([]string, error)

  // Sync
  SyncAll(ctx context.Context) error
  SyncCollection(ctx context.Context, collectionName string) error
}

type SearchQuery struct {
  Query          string            `json:"q"`
  Collections    []string          `json:"collections"`    // ['movies', 'tvshows']
  QueryBy        []string          `json:"query_by"`       // ['title', 'overview']
  FilterBy       string            `json:"filter_by"`      // 'release_year:>2020'
  SortBy         string            `json:"sort_by"`        // 'rating:desc'
  Page           int               `json:"page"`
  PerPage        int               `json:"per_page"`
  UserID         uuid.UUID         `json:"-"`              // For permission filtering
}

type SearchResults struct {
  Hits          []SearchHit       `json:"hits"`
  Found         int               `json:"found"`
  Page          int               `json:"page"`
  SearchTimeMS  int               `json:"search_time_ms"`
}

type SearchHit struct {
  Document      map[string]interface{} `json:"document"`
  Highlights    map[string]interface{} `json:"highlights"`
  TextMatch     int64                  `json:"text_match"`
}
```


### Dependencies
**Go Packages**:
- `github.com/google/uuid`
- `github.com/jackc/pgx/v5`
- `github.com/typesense/typesense-go/typesense` - Typesense client
- `github.com/riverqueue/river` - Background sync jobs
- `go.uber.org/fx`

**External Services**:
- Typesense server (https://typesense.org/)






## Configuration
### Environment Variables

```bash
TYPESENSE_HOST=localhost
TYPESENSE_PORT=8108
TYPESENSE_API_KEY=your_api_key
TYPESENSE_PROTOCOL=http
SEARCH_SYNC_INTERVAL=5m
```


### Config Keys
```yaml
search:
  typesense:
    host: localhost
    port: 8108
    api_key: your_api_key
    protocol: http
  sync:
    interval: 5m
    batch_size: 100
  query:
    max_results: 50
    typo_tolerance: true
```



## API Endpoints
```
# Search
GET    /api/v1/search?q=query&collections=movies,tvshows  # Multi-collection search
POST   /api/v1/search/multi                               # Advanced multi-search

# Admin (indexing)
POST   /api/v1/search/sync                                # Trigger full sync
POST   /api/v1/search/sync/:collection                    # Sync specific collection
POST   /api/v1/search/reindex/:collection                 # Drop and recreate collection
```

**Example Search Request**:
```json
GET /api/v1/search?q=inception&collections=movies&filter_by=release_year:>2000&sort_by=rating:desc
```

**Example Search Response**:
```json
{
  "hits": [
    {
      "document": {
        "id": "27205",
        "title": "Inception",
        "overview": "A thief who steals corporate secrets...",
        "release_year": 2010,
        "rating": 8.4,
        "poster_url": "https://..."
      },
      "highlights": {
        "title": {
          "matched_tokens": ["Inception"],
          "snippet": "<mark>Inception</mark>"
        }
      },
      "text_match": 578934906667
    }
  ],
  "found": 1,
  "page": 1,
  "search_time_ms": 12
}
```

**Typesense Collection Schema** (Movies):
```json
{
  "name": "movies",
  "fields": [
    {"name": "id", "type": "string"},
    {"name": "title", "type": "string"},
    {"name": "overview", "type": "string"},
    {"name": "release_year", "type": "int32", "facet": true},
    {"name": "rating", "type": "float"},
    {"name": "genres", "type": "string[]", "facet": true},
    {"name": "library_id", "type": "string"},
    {"name": "poster_url", "type": "string", "optional": true}
  ],
  "default_sorting_field": "rating"
}
```



## Testing Strategy

### Unit Tests

<!-- Unit test strategy -->

### Integration Tests

<!-- Integration test strategy -->

### Test Coverage

Target: **80% minimum**







## Related Documentation
### Design Documents
- [services](INDEX.md)
- [01_ARCHITECTURE](../architecture/01_ARCHITECTURE.md)
- [02_DESIGN_PRINCIPLES](../architecture/02_DESIGN_PRINCIPLES.md)
- [03_METADATA_SYSTEM](../architecture/03_METADATA_SYSTEM.md)

### External Sources
- [Uber fx](../../sources/tooling/fx.md) - Auto-resolved from fx
- [River Job Queue](../../sources/tooling/river.md) - Auto-resolved from river
- [Typesense API](../../sources/infrastructure/typesense.md) - Auto-resolved from typesense
- [Typesense Go Client](../../sources/infrastructure/typesense-go.md) - Auto-resolved from typesense-go

