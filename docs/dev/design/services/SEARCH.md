# Search Service

> Full-text search via Typesense with per-module collections


<!-- TOC-START -->

## Table of Contents

- [Status](#status)
- [Developer Resources](#developer-resources)
- [Overview](#overview)
- [Goals](#goals)
- [Non-Goals](#non-goals)
- [Technical Design](#technical-design)
  - [Collections](#collections)
  - [Collection Schema Example](#collection-schema-example)
  - [Repository Interface](#repository-interface)
  - [Service Layer](#service-layer)
  - [Access Control](#access-control)
- [River Jobs](#river-jobs)
- [API Endpoints](#api-endpoints)
- [Configuration](#configuration)
- [Implementation Files](#implementation-files)
- [Checklist](#checklist)
- [Sources & Cross-References](#sources-cross-references)
  - [Cross-Reference Indexes](#cross-reference-indexes)
  - [Referenced Sources](#referenced-sources)
- [Related Design Docs](#related-design-docs)
  - [In This Section](#in-this-section)
  - [Related Topics](#related-topics)
  - [Indexes](#indexes)
- [Related Documents](#related-documents)

<!-- TOC-END -->

**Module**: `internal/service/search`
**Dependencies**: [00_SOURCE_OF_TRUTH.md](../00_SOURCE_OF_TRUTH.md#infrastructure-components)

## Status

| Dimension | Status |
|-----------|--------|
| Design | ✅ |
| Sources | ✅ |
| Instructions | ✅ |
| Code | 🔴 |
| Linting | 🔴 |
| Unit Testing | 🔴 |
| Integration Testing | 🔴 |
## Developer Resources

> Package versions: [00_SOURCE_OF_TRUTH.md](../00_SOURCE_OF_TRUTH.md#infrastructure-components)

| Package | Purpose |
|---------|---------|
| Typesense | Search engine |
| typesense-go | Go client library |
| River | Background indexing jobs |

---

## Overview

The Search service provides fast, typo-tolerant full-text search across all content modules using Typesense. Each module has its own collection with optimized schema.

Key features:
- Per-module search collections
- Faceted filtering (genre, year, rating)
- Typo tolerance and fuzzy matching
- Real-time index updates via River jobs
- Access-controlled search results

## Goals

- Sub-50ms search response times
- Typo-tolerant matching
- Faceted filtering for discovery
- Real-time index synchronization
- Respect user permissions in results

## Non-Goals

- Replace database queries for exact lookups
- Full-text search within media files (subtitles, lyrics)
- Semantic/AI-powered search (future consideration)

---

## Technical Design

### Collections

> See [00_SOURCE_OF_TRUTH.md](../00_SOURCE_OF_TRUTH.md#infrastructure-components) for Typesense version.

| Collection | Primary Fields | Facets |
|------------|---------------|--------|
| `movies` | title, overview, tagline | genres, year, rating |
| `series` | title, overview | genres, year, networks |
| `episodes` | title, overview | series_id, season |
| `tracks` | title, artist, album | genres, year |
| `audiobooks` | title, author, narrator | genres |
| `qar_expeditions` | title, crew, port | flags, year |
| `qar_voyages` | title, crew, expedition_id | flags, port |
| `qar_crew` | name, aliases | ports, flags |
| `qar_ports` | name | parent_port |
| `qar_treasures` | title, crew | port, flags |

> **QAR Terminology**: expeditions=movies, voyages=scenes, crew=performers, ports=studios, treasures=galleries, flags=tags. See [00_SOURCE_OF_TRUTH.md](../00_SOURCE_OF_TRUTH.md#qar-obfuscation-terminology).

### Collection Schema Example

```json
{
  "name": "movies",
  "fields": [
    {"name": "id", "type": "string"},
    {"name": "title", "type": "string"},
    {"name": "original_title", "type": "string", "optional": true},
    {"name": "overview", "type": "string", "optional": true},
    {"name": "tagline", "type": "string", "optional": true},
    {"name": "year", "type": "int32", "facet": true},
    {"name": "genres", "type": "string[]", "facet": true},
    {"name": "rating", "type": "float", "facet": true},
    {"name": "popularity", "type": "float"},
    {"name": "poster_url", "type": "string", "optional": true},
    {"name": "library_id", "type": "string", "facet": true}
  ],
  "default_sorting_field": "popularity"
}
```

### Repository Interface

```go
type SearchRepository interface {
    // Index management
    EnsureCollection(ctx context.Context, name string, schema CollectionSchema) error
    DeleteCollection(ctx context.Context, name string) error

    // Document operations
    IndexDocument(ctx context.Context, collection string, doc interface{}) error
    IndexDocuments(ctx context.Context, collection string, docs []interface{}) error
    DeleteDocument(ctx context.Context, collection string, id string) error

    // Search
    Search(ctx context.Context, req SearchRequest) (*SearchResult, error)
    MultiSearch(ctx context.Context, reqs []SearchRequest) ([]SearchResult, error)
}

type SearchRequest struct {
    Collection  string
    Query       string
    QueryBy     []string          // Fields to search
    FilterBy    string            // Typesense filter syntax
    FacetBy     []string
    SortBy      string
    Page        int
    PerPage     int
    GroupBy     string            // Optional grouping
}

type SearchResult struct {
    Hits       []SearchHit
    TotalHits  int
    Facets     map[string][]FacetCount
    SearchTime time.Duration
}
```

### Service Layer

```go
type SearchService struct {
    repo   SearchRepository
    grants *grants.Service
}

func (s *SearchService) Search(ctx context.Context, userID uuid.UUID, query string, opts SearchOptions) (*SearchResult, error)
func (s *SearchService) SearchModule(ctx context.Context, userID uuid.UUID, module string, query string, opts SearchOptions) (*SearchResult, error)
func (s *SearchService) Reindex(ctx context.Context, module string) error
func (s *SearchService) IndexItem(ctx context.Context, module string, item interface{}) error
func (s *SearchService) RemoveItem(ctx context.Context, module string, id uuid.UUID) error
```

### Access Control

Search results are filtered based on user permissions:

```go
func (s *SearchService) buildAccessFilter(ctx context.Context, userID uuid.UUID) (string, error) {
    // Get libraries user can access
    libraries, err := s.grants.GetAccessibleLibraries(ctx, userID)
    if err != nil {
        return "", err
    }

    // Build Typesense filter
    ids := make([]string, len(libraries))
    for i, lib := range libraries {
        ids[i] = lib.ID.String()
    }
    return fmt.Sprintf("library_id:[%s]", strings.Join(ids, ",")), nil
}
```

---

## River Jobs

```go
type IndexDocumentArgs struct {
    Collection string    `json:"collection"`
    DocumentID uuid.UUID `json:"document_id"`
    Action     string    `json:"action"` // "index", "delete"
}

func (IndexDocumentArgs) Kind() string { return "search.index_document" }

type ReindexCollectionArgs struct {
    Collection string `json:"collection"`
}

func (ReindexCollectionArgs) Kind() string { return "search.reindex_collection" }
```

---

## API Endpoints

```
GET /api/v1/search?q=query&modules=movies,series&genres=action&year=2024
    → Global search across specified modules

GET /api/v1/movies/search?q=query&genres=action
    → Module-specific search

GET /api/v1/search/suggest?q=par
    → Autocomplete suggestions
```

---

## Configuration

```yaml
search:
  url: "http://typesense:8108"
  api_key: "${TYPESENSE_API_KEY}"
  connection_timeout: 5s

  # Per-collection settings
  collections:
    movies:
      typo_tolerance: true
      num_typos: 2
    tracks:
      typo_tolerance: true
      num_typos: 1
```

---

## Implementation Files

| File | Action | Description |
|------|--------|-------------|
| `internal/service/search/service.go` | CREATE | Core search service |
| `internal/service/search/repository.go` | CREATE | Repository interface |
| `internal/service/search/typesense.go` | CREATE | Typesense implementation |
| `internal/service/search/schemas.go` | CREATE | Collection schemas |
| `internal/service/search/jobs.go` | CREATE | River indexing jobs |
| `internal/service/search/module.go` | CREATE | fx module |

---

## Checklist

- [ ] Typesense client wrapper created
- [ ] Collection schemas defined
- [ ] Repository interface defined
- [ ] Typesense repository implemented
- [ ] Service layer with access control
- [ ] River jobs for real-time indexing
- [ ] Reindex job for full rebuild
- [ ] API handlers created
- [ ] Tests written
- [ ] Documentation updated

---


<!-- SOURCE-BREADCRUMBS-START -->

## Sources & Cross-References

> Auto-generated section linking to external documentation sources

### Cross-Reference Indexes

- [All Sources Index](../../sources/SOURCES_INDEX.md) - Complete list of external documentation
- [Design ↔ Sources Map](../../sources/DESIGN_CROSSREF.md) - Which docs reference which sources

### Referenced Sources

| Source | Documentation |
|--------|---------------|
| [River Job Queue](https://pkg.go.dev/github.com/riverqueue/river) | [Local](../../sources/tooling/river.md) |
| [Typesense API](https://typesense.org/docs/latest/api/) | [Local](../../sources/infrastructure/typesense.md) |
| [Typesense Go Client](https://github.com/typesense/typesense-go) | [Local](../../sources/infrastructure/typesense-go.md) |
| [Uber fx](https://pkg.go.dev/go.uber.org/fx) | [Local](../../sources/tooling/fx.md) |

<!-- SOURCE-BREADCRUMBS-END -->

<!-- DESIGN-BREADCRUMBS-START -->

## Related Design Docs

> Auto-generated cross-references to related design documentation

**Category**: [Services](INDEX.md)

### In This Section

- [Activity Service](ACTIVITY.md)
- [Analytics Service](ANALYTICS.md)
- [API Keys Service](APIKEYS.md)
- [Auth Service](AUTH.md)
- [Fingerprint Service](FINGERPRINT.md)
- [Grants Service](GRANTS.md)
- [Library Service](LIBRARY.md)
- [Metadata Service](METADATA.md)

### Related Topics

- [Revenge - Architecture v2](../architecture/01_ARCHITECTURE.md) _Architecture_
- [Revenge - Design Principles](../architecture/02_DESIGN_PRINCIPLES.md) _Architecture_
- [Revenge - Metadata System](../architecture/03_METADATA_SYSTEM.md) _Architecture_
- [Revenge - Player Architecture](../architecture/04_PLAYER_ARCHITECTURE.md) _Architecture_
- [Plugin Architecture Decision](../architecture/05_PLUGIN_ARCHITECTURE_DECISION.md) _Architecture_

### Indexes

- [Design Index](../DESIGN_INDEX.md) - All design docs by category/topic
- [Source of Truth](../00_SOURCE_OF_TRUTH.md) - Package versions and status

<!-- DESIGN-BREADCRUMBS-END -->

## Related Documents

- [Grants Service](GRANTS.md) - Access control for search results
- [Library Service](LIBRARY.md) - Library-based filtering
- [Metadata Service](METADATA.md) - Content indexed by search
- [00_SOURCE_OF_TRUTH.md](../00_SOURCE_OF_TRUTH.md) - Service inventory
