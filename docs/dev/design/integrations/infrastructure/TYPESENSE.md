# Typesense Integration

> Fast, typo-tolerant search engine


<!-- TOC-START -->

## Table of Contents

- [Status](#status)
- [Overview](#overview)
- [Developer Resources](#developer-resources)
- [Connection Details](#connection-details)
- [Configuration](#configuration)
- [Collection Schemas](#collection-schemas)
  - [Movies Collection](#movies-collection)
  - [TV Shows Collection](#tv-shows-collection)
  - [Music Collection](#music-collection)
  - [People Collection](#people-collection)
- [Search Operations](#search-operations)
  - [Basic Search](#basic-search)
  - [Faceted Search](#faceted-search)
  - [Multi-Search](#multi-search)
  - [Vector Search (Similarity)](#vector-search-similarity)
- [Indexing](#indexing)
  - [Index Document](#index-document)
  - [Batch Indexing](#batch-indexing)
  - [Delete from Index](#delete-from-index)
- [Implementation Checklist](#implementation-checklist)
- [Docker Compose](#docker-compose)
- [Health Checks](#health-checks)
- [Monitoring](#monitoring)
  - [Key Metrics](#key-metrics)
  - [Typesense Stats](#typesense-stats)
- [Error Handling](#error-handling)
- [Sources & Cross-References](#sources-cross-references)
  - [Cross-Reference Indexes](#cross-reference-indexes)
  - [Referenced Sources](#referenced-sources)
- [Related Design Docs](#related-design-docs)
  - [In This Section](#in-this-section)
  - [Related Topics](#related-topics)
  - [Indexes](#indexes)
- [Related Documentation](#related-documentation)

<!-- TOC-END -->

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
**Priority**: 🟡 MEDIUM (Phase 1 - Core Infrastructure)
**Type**: Search engine

---

## Overview

Typesense is Revenge's primary search engine, providing:
- Full-text search across all media
- Typo-tolerant queries
- Faceted search & filtering
- Geo-search (for location-based features)
- Vector search (for AI recommendations)
- Sub-millisecond latency

**Why Typesense?**:
- Simple to deploy (single binary)
- Built-in typo tolerance
- RAM-optimized for speed
- Open source alternative to Algolia
- REST API + Go SDK

---

## Developer Resources

- 📚 **Docs**: https://typesense.org/docs/
- 🔗 **GitHub**: https://github.com/typesense/typesense
- 🔗 **Go SDK**: https://github.com/typesense/typesense-go
- 📖 **API Reference**: https://typesense.org/docs/latest/api/

---

## Connection Details

**Default Settings**:
| Setting | Value |
|---------|-------|
| Host | `localhost` |
| Port | `8108` |
| Protocol | `http` |
| API Key | (required) |

---

## Configuration

```yaml
# configs/config.yaml
search:
  enabled: true
  driver: "typesense"

  connection:
    host: "${REVENGE_SEARCH_HOST:localhost}"
    port: ${REVENGE_SEARCH_PORT:8108}
    protocol: "http"  # or "https"
    api_key: "${REVENGE_SEARCH_API_KEY:}"

  # Connection pool
  pool:
    connection_timeout: "5s"
    max_retries: 3

  # Index settings
  indexing:
    batch_size: 100
    workers: 4
```

---

## Collection Schemas

### Movies Collection

```go
moviesSchema := &api.CollectionSchema{
    Name: "movies",
    Fields: []api.Field{
        {Name: "id", Type: "string"},
        {Name: "title", Type: "string"},
        {Name: "original_title", Type: "string", Optional: ptr(true)},
        {Name: "overview", Type: "string", Optional: ptr(true)},
        {Name: "tagline", Type: "string", Optional: ptr(true)},
        {Name: "year", Type: "int32", Facet: ptr(true)},
        {Name: "release_date", Type: "int64"},  // Unix timestamp
        {Name: "runtime", Type: "int32", Optional: ptr(true)},
        {Name: "genres", Type: "string[]", Facet: ptr(true)},
        {Name: "cast", Type: "string[]"},
        {Name: "director", Type: "string[]", Facet: ptr(true)},
        {Name: "rating", Type: "float", Optional: ptr(true)},
        {Name: "vote_count", Type: "int32", Optional: ptr(true)},
        {Name: "poster_path", Type: "string", Optional: ptr(true), Index: ptr(false)},
        {Name: "library_id", Type: "string", Facet: ptr(true)},
        {Name: "content_rating", Type: "string", Facet: ptr(true), Optional: ptr(true)},
        // Vector field for similarity search
        {Name: "embedding", Type: "float[]", NumDim: ptr(int32(768)), Optional: ptr(true)},
    },
    DefaultSortingField: ptr("release_date"),
}
```

### TV Shows Collection

```go
tvShowsSchema := &api.CollectionSchema{
    Name: "tvshows",
    Fields: []api.Field{
        {Name: "id", Type: "string"},
        {Name: "title", Type: "string"},
        {Name: "original_title", Type: "string", Optional: ptr(true)},
        {Name: "overview", Type: "string", Optional: ptr(true)},
        {Name: "year", Type: "int32", Facet: ptr(true)},
        {Name: "first_air_date", Type: "int64"},
        {Name: "status", Type: "string", Facet: ptr(true)},  // Ended, Continuing
        {Name: "genres", Type: "string[]", Facet: ptr(true)},
        {Name: "cast", Type: "string[]"},
        {Name: "network", Type: "string[]", Facet: ptr(true)},
        {Name: "season_count", Type: "int32"},
        {Name: "episode_count", Type: "int32"},
        {Name: "rating", Type: "float", Optional: ptr(true)},
        {Name: "library_id", Type: "string", Facet: ptr(true)},
    },
    DefaultSortingField: ptr("first_air_date"),
}
```

### Music Collection

```go
musicSchema := &api.CollectionSchema{
    Name: "music",
    Fields: []api.Field{
        {Name: "id", Type: "string"},
        {Name: "type", Type: "string", Facet: ptr(true)},  // artist, album, track
        {Name: "title", Type: "string"},
        {Name: "artist", Type: "string[]"},
        {Name: "album", Type: "string", Optional: ptr(true)},
        {Name: "year", Type: "int32", Facet: ptr(true), Optional: ptr(true)},
        {Name: "genres", Type: "string[]", Facet: ptr(true)},
        {Name: "duration", Type: "int32", Optional: ptr(true)},
        {Name: "track_number", Type: "int32", Optional: ptr(true)},
        {Name: "library_id", Type: "string", Facet: ptr(true)},
    },
}
```

### People Collection

```go
peopleSchema := &api.CollectionSchema{
    Name: "people",
    Fields: []api.Field{
        {Name: "id", Type: "string"},
        {Name: "name", Type: "string"},
        {Name: "known_for", Type: "string", Facet: ptr(true)},  // Acting, Directing
        {Name: "known_for_titles", Type: "string[]"},
        {Name: "birth_year", Type: "int32", Optional: ptr(true)},
        {Name: "popularity", Type: "float"},
        {Name: "profile_path", Type: "string", Optional: ptr(true), Index: ptr(false)},
    },
    DefaultSortingField: ptr("popularity"),
}
```

---

## Search Operations

### Basic Search

```go
func (s *SearchService) SearchMovies(ctx context.Context, query string, page int) (*SearchResult, error) {
    params := &api.SearchCollectionParams{
        Q:       query,
        QueryBy: "title,original_title,overview,cast,director",
        Page:    ptr(page),
        PerPage: ptr(20),
        // Prioritize title matches
        QueryByWeights: ptr("4,2,1,1,1"),
    }

    result, err := s.client.Collection("movies").Documents().Search(ctx, params)
    if err != nil {
        return nil, fmt.Errorf("search failed: %w", err)
    }

    return mapSearchResult(result), nil
}
```

### Faceted Search

```go
func (s *SearchService) SearchWithFilters(ctx context.Context, req SearchRequest) (*SearchResult, error) {
    filterParts := []string{}

    if req.Year > 0 {
        filterParts = append(filterParts, fmt.Sprintf("year:=%d", req.Year))
    }
    if len(req.Genres) > 0 {
        filterParts = append(filterParts, fmt.Sprintf("genres:[%s]", strings.Join(req.Genres, ",")))
    }
    if req.MinRating > 0 {
        filterParts = append(filterParts, fmt.Sprintf("rating:>=%f", req.MinRating))
    }

    params := &api.SearchCollectionParams{
        Q:         req.Query,
        QueryBy:   "title,overview",
        FilterBy:  ptr(strings.Join(filterParts, " && ")),
        FacetBy:   ptr("genres,year,director"),
        SortBy:    ptr(req.SortBy),  // e.g., "rating:desc"
        Page:      ptr(req.Page),
        PerPage:   ptr(req.PageSize),
    }

    return s.search(ctx, "movies", params)
}
```

### Multi-Search

```go
func (s *SearchService) GlobalSearch(ctx context.Context, query string) (*GlobalSearchResult, error) {
    multiSearchParams := api.MultiSearchParams{}

    searches := api.MultiSearchSearchesParameter{
        Searches: []api.MultiSearchCollectionParameters{
            {Collection: "movies", Q: &query, QueryBy: ptr("title,overview")},
            {Collection: "tvshows", Q: &query, QueryBy: ptr("title,overview")},
            {Collection: "music", Q: &query, QueryBy: ptr("title,artist,album")},
            {Collection: "people", Q: &query, QueryBy: ptr("name")},
        },
    }

    results, err := s.client.MultiSearch.Perform(ctx, &multiSearchParams, searches)
    if err != nil {
        return nil, err
    }

    return mapGlobalResults(results), nil
}
```

### Vector Search (Similarity)

```go
func (s *SearchService) FindSimilar(ctx context.Context, movieID string) ([]Movie, error) {
    // Get embedding for source movie
    movie, _ := s.client.Collection("movies").Document(movieID).Retrieve(ctx)
    embedding := movie["embedding"].([]float64)

    params := &api.SearchCollectionParams{
        Q:              "*",
        VectorQuery:    ptr(fmt.Sprintf("embedding:([%s], k:10)", floatsToString(embedding))),
        ExcludeFields:  ptr("embedding"),
    }

    return s.search(ctx, "movies", params)
}
```

---

## Indexing

### Index Document

```go
func (i *Indexer) IndexMovie(ctx context.Context, movie *Movie) error {
    doc := map[string]interface{}{
        "id":             movie.ID.String(),
        "title":          movie.Title,
        "original_title": movie.OriginalTitle,
        "overview":       movie.Overview,
        "year":           movie.Year,
        "release_date":   movie.ReleaseDate.Unix(),
        "genres":         movie.GenreNames(),
        "cast":           movie.CastNames(),
        "director":       movie.DirectorNames(),
        "rating":         movie.Rating,
        "library_id":     movie.LibraryID.String(),
    }

    _, err := i.client.Collection("movies").Documents().Upsert(ctx, doc)
    return err
}
```

### Batch Indexing

```go
func (i *Indexer) IndexMoviesBatch(ctx context.Context, movies []*Movie) error {
    var documents []interface{}
    for _, movie := range movies {
        documents = append(documents, movieToDoc(movie))
    }

    params := &api.ImportDocumentsParams{
        Action: ptr("upsert"),
    }

    _, err := i.client.Collection("movies").Documents().Import(ctx, documents, params)
    return err
}
```

### Delete from Index

```go
func (i *Indexer) DeleteMovie(ctx context.Context, movieID string) error {
    _, err := i.client.Collection("movies").Document(movieID).Delete(ctx)
    return err
}

// Delete by filter
func (i *Indexer) DeleteByLibrary(ctx context.Context, libraryID string) error {
    params := &api.DeleteDocumentsParams{
        FilterBy: fmt.Sprintf("library_id:=%s", libraryID),
    }
    _, err := i.client.Collection("movies").Documents().Delete(ctx, params)
    return err
}
```

---

## Implementation Checklist

- [ ] **Search Client** (`internal/infra/search/client.go`)
  - [ ] Connection management
  - [ ] API key rotation
  - [ ] Health checks

- [ ] **Collection Schemas** (`internal/infra/search/schemas/`)
  - [ ] movies.go
  - [ ] tvshows.go
  - [ ] music.go
  - [ ] audiobooks.go
  - [ ] books.go
  - [ ] people.go

- [ ] **Indexer** (`internal/infra/search/indexer.go`)
  - [ ] Single document upsert
  - [ ] Batch import
  - [ ] Delete operations
  - [ ] Full reindex

- [ ] **Search Service** (`internal/service/search/`)
  - [ ] Basic search
  - [ ] Filtered search
  - [ ] Faceted search
  - [ ] Global search
  - [ ] Autocomplete

- [ ] **River Jobs** (`internal/infra/jobs/search/`)
  - [ ] IndexItemJob
  - [ ] ReindexLibraryJob
  - [ ] SyncIndexJob

---

## Docker Compose

```yaml
services:
  typesense:
    image: typesense/typesense:27.0
    container_name: revenge-typesense
    ports:
      - "8108:8108"
    volumes:
      - typesense_data:/data
    environment:
      TYPESENSE_API_KEY: ${TYPESENSE_API_KEY:-revenge_search_key}
      TYPESENSE_DATA_DIR: /data
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8108/health"]
      interval: 10s
      timeout: 5s
      retries: 5

volumes:
  typesense_data:
```

---

## Health Checks

```go
func (c *SearchClient) HealthCheck(ctx context.Context) error {
    ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
    defer cancel()

    health, err := c.client.Health(ctx, 2*time.Second)
    if err != nil {
        return fmt.Errorf("typesense health check failed: %w", err)
    }
    if !health {
        return errors.New("typesense unhealthy")
    }
    return nil
}
```

---

## Monitoring

### Key Metrics

```go
var (
    searchLatency = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "revenge_search_duration_seconds",
            Help:    "Search latency",
            Buckets: []float64{.001, .005, .01, .025, .05, .1, .25},
        },
        []string{"collection"},
    )

    searchResults = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "revenge_search_results_total",
            Help:    "Number of search results",
            Buckets: []float64{0, 1, 5, 10, 20, 50, 100},
        },
        []string{"collection"},
    )

    indexOperations = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "revenge_search_index_operations_total",
            Help: "Index operations",
        },
        []string{"collection", "operation"},
    )
)
```

### Typesense Stats

```bash
# Get cluster stats
curl http://localhost:8108/stats.json -H "X-TYPESENSE-API-KEY: ${API_KEY}"

# Get collection stats
curl http://localhost:8108/collections/movies -H "X-TYPESENSE-API-KEY: ${API_KEY}"
```

---

## Error Handling

| Error | Cause | Solution |
|-------|-------|----------|
| 401 Unauthorized | Invalid API key | Check TYPESENSE_API_KEY |
| 404 Not Found | Collection/doc missing | Create collection first |
| 409 Conflict | Schema mismatch | Update or recreate schema |
| 503 Service Unavailable | Server overloaded | Scale or optimize queries |

---


<!-- SOURCE-BREADCRUMBS-START -->

## Sources & Cross-References

> Auto-generated section linking to external documentation sources

### Cross-Reference Indexes

- [All Sources Index](../../../sources/SOURCES_INDEX.md) - Complete list of external documentation
- [Design ↔ Sources Map](../../../sources/DESIGN_CROSSREF.md) - Which docs reference which sources

### Referenced Sources

| Source | Documentation |
|--------|---------------|
| [Dragonfly Documentation](https://www.dragonflydb.io/docs) | [Local](../../../sources/infrastructure/dragonfly.md) |
| [PostgreSQL Arrays](https://www.postgresql.org/docs/current/arrays.html) | [Local](../../../sources/database/postgresql-arrays.md) |
| [PostgreSQL JSON Functions](https://www.postgresql.org/docs/current/functions-json.html) | [Local](../../../sources/database/postgresql-json.md) |
| [Prometheus Go Client](https://pkg.go.dev/github.com/prometheus/client_golang/prometheus) | [Local](../../../sources/observability/prometheus.md) |
| [Prometheus Metric Types](https://prometheus.io/docs/concepts/metric_types/) | [Local](../../../sources/observability/prometheus-metrics.md) |
| [River Job Queue](https://pkg.go.dev/github.com/riverqueue/river) | [Local](../../../sources/tooling/river.md) |
| [Typesense API](https://typesense.org/docs/latest/api/) | [Local](../../../sources/infrastructure/typesense.md) |
| [Typesense Go Client](https://github.com/typesense/typesense-go) | [Local](../../../sources/infrastructure/typesense-go.md) |
| [pgx PostgreSQL Driver](https://pkg.go.dev/github.com/jackc/pgx/v5) | [Local](../../../sources/database/pgx.md) |

<!-- SOURCE-BREADCRUMBS-END -->

<!-- DESIGN-BREADCRUMBS-START -->

## Related Design Docs

> Auto-generated cross-references to related design documentation

**Category**: [Infrastructure](INDEX.md)

### In This Section

- [Dragonfly Integration](DRAGONFLY.md)
- [PostgreSQL Integration](POSTGRESQL.md)
- [River Integration](RIVER.md)

### Related Topics

- [Revenge - Architecture v2](../../architecture/01_ARCHITECTURE.md) _Architecture_
- [Revenge - Design Principles](../../architecture/02_DESIGN_PRINCIPLES.md) _Architecture_
- [Revenge - Metadata System](../../architecture/03_METADATA_SYSTEM.md) _Architecture_
- [Revenge - Player Architecture](../../architecture/04_PLAYER_ARCHITECTURE.md) _Architecture_
- [Plugin Architecture Decision](../../architecture/05_PLUGIN_ARCHITECTURE_DECISION.md) _Architecture_

### Indexes

- [Design Index](../../DESIGN_INDEX.md) - All design docs by category/topic
- [Source of Truth](../../00_SOURCE_OF_TRUTH.md) - Package versions and status

<!-- DESIGN-BREADCRUMBS-END -->

## Related Documentation

- [PostgreSQL](POSTGRESQL.md)
- [Dragonfly](DRAGONFLY.md)
- [River](RIVER.md)
