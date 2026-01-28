# Typesense Search Instructions

> Source: https://github.com/typesense/typesense-go (v3), https://typesense.org/docs/

Apply to: `**/internal/infra/search/**/*.go`, `**/internal/content/**/search*.go`

## Installation

```bash
go get github.com/typesense/typesense-go/v3/typesense
```

## Client Setup

### Basic Client

```go
import "github.com/typesense/typesense-go/v3/typesense"

client := typesense.NewClient(
    typesense.WithServer("http://localhost:8108"),
    typesense.WithAPIKey("<API_KEY>"),
)
```

### Client with Advanced Options

```go
client := typesense.NewClient(
    typesense.WithServer("http://localhost:8108"),
    typesense.WithAPIKey("<API_KEY>"),
    typesense.WithConnectionTimeout(5*time.Second),
    typesense.WithCircuitBreakerMaxRequests(50),
    typesense.WithCircuitBreakerInterval(2*time.Minute),
    typesense.WithCircuitBreakerTimeout(1*time.Minute),
)
```

### Multi-Node Configuration

```go
client := typesense.NewClient(
    typesense.WithNearestNode("https://xxx.a1.typesense.net:443"),
    typesense.WithNodes([]string{
        "https://xxx-1.a1.typesense.net:443",
        "https://xxx-2.a1.typesense.net:443",
        "https://xxx-3.a1.typesense.net:443",
    }),
    typesense.WithAPIKey("<API_KEY>"),
    typesense.WithNumRetries(5),
    typesense.WithRetryInterval(1*time.Second),
    typesense.WithHealthcheckInterval(2*time.Minute),
)
```

## Collection Schema

### Create Collection

```go
schema := &api.CollectionSchema{
    Name: "movies",
    Fields: []api.Field{
        {Name: "title", Type: "string"},
        {Name: "overview", Type: "string"},
        {Name: "year", Type: "int32"},
        {Name: "genres", Type: "string[]", Facet: pointer.True()},
        {Name: "rating", Type: "float"},
    },
    DefaultSortingField: pointer.String("year"),
}

client.Collections().Create(context.Background(), schema)
```

### Field Types

| Type             | Description          |
| ---------------- | -------------------- |
| `string`         | Text field           |
| `string[]`       | Array of strings     |
| `int32`, `int64` | Integer fields       |
| `float`          | Floating point       |
| `bool`           | Boolean              |
| `geopoint`       | Lat/long coordinates |
| `auto`           | Auto-detect type     |

### Faceted Fields

Set `Facet: pointer.True()` for filterable/facetable fields:

```go
{Name: "genres", Type: "string[]", Facet: pointer.True()}
```

## Document Operations

### Index Document

```go
document := struct {
    ID       string `json:"id"`
    Title    string `json:"title"`
    Overview string `json:"overview"`
    Year     int    `json:"year"`
}{
    ID:       "123",
    Title:    "The Matrix",
    Overview: "A computer hacker...",
    Year:     1999,
}

client.Collection("movies").Documents().Create(ctx, document)
```

### Upsert Document

```go
client.Collection("movies").Documents().Upsert(ctx, document)
```

### Import Documents (Bulk)

```go
documents := []interface{}{doc1, doc2, doc3}
params := &api.ImportDocumentsParams{
    Action:    pointer.String("upsert"),
    BatchSize: pointer.Int(40),
}

client.Collection("movies").Documents().Import(ctx, documents, params)
```

Import actions: `create`, `upsert`, `update`

### Typed Document Operations (v2.0.0+)

```go
type MovieDocument struct {
    ID       string `json:"id"`
    Title    string `json:"title"`
    Overview string `json:"overview"`
    Year     int    `json:"year"`
}

// Typed retrieval
doc, err := typesense.GenericCollection[*MovieDocument](client, "movies").
    Document("123").Retrieve(ctx)
```

## Search

### Basic Search

```go
searchParams := &api.SearchCollectionParams{
    Q:       pointer.String("matrix"),
    QueryBy: pointer.String("title,overview"),
}

result, err := client.Collection("movies").Documents().Search(ctx, searchParams)
```

### Search with Filters

```go
searchParams := &api.SearchCollectionParams{
    Q:        pointer.String("action"),
    QueryBy:  pointer.String("title,overview"),
    FilterBy: pointer.String("year:>2000 && genres:=Action"),
    SortBy:   &([]string{"rating:desc"}),
}
```

### Multiple QueryBy Fields

```go
searchParams := &api.SearchCollectionParams{
    Q:       pointer.String("query"),
    QueryBy: pointer.String("title, overview, cast"),  // comma-separated
}
```

### Pagination

```go
searchParams := &api.SearchCollectionParams{
    Q:       pointer.String("*"),
    QueryBy: pointer.String("title"),
    Page:    pointer.Int(1),
    PerPage: pointer.Int(20),
}
```

## Delete Operations

### Delete Single Document

```go
client.Collection("movies").Document("123").Delete(ctx)
```

### Delete by Filter

```go
filter := &api.DeleteDocumentsParams{
    FilterBy:  "year:<2000",
    BatchSize: 100,
}
client.Collection("movies").Documents().Delete(ctx, filter)
```

### Drop Collection

```go
client.Collection("movies").Delete(ctx)
```

## Per-Module Collections

Each content module has its own collection:

| Module      | Collection   | Key Fields                          |
| ----------- | ------------ | ----------------------------------- |
| movie       | `movies`     | title, overview, year, genres, cast |
| tvshow      | `series`     | title, overview, year, genres, cast |
| music       | `tracks`     | title, artist, album, genre         |
| audiobook   | `audiobooks` | title, author, narrator             |
| book        | `books`      | title, author, publisher            |
| podcast     | `podcasts`   | title, author, description          |
| adult_movie | `c_movies`   | title, performers, studio, tags     |
| adult_show  | `c_series`   | title, performers, studio, tags     |

## Revenge Patterns

### Search Service Interface

```go
type SearchService interface {
    Index(ctx context.Context, collection string, doc any) error
    Upsert(ctx context.Context, collection string, doc any) error
    Delete(ctx context.Context, collection string, id string) error
    Search(ctx context.Context, collection string, params SearchParams) (*SearchResult, error)
}
```

### Index on Change

Index documents via River jobs on create/update:

```go
type IndexMovieArgs struct {
    MovieID uuid.UUID `json:"movie_id"`
}

func (IndexMovieArgs) Kind() string { return "movie.index_search" }

func (w *IndexMovieWorker) Work(ctx context.Context, job *river.Job[IndexMovieArgs]) error {
    movie, err := w.repo.Get(ctx, job.Args.MovieID)
    if err != nil {
        return err
    }
    return w.search.Upsert(ctx, "movies", movie.ToSearchDoc())
}
```

### Error Handling

```go
result, err := client.Collection("movies").Documents().Search(ctx, params)
if err != nil {
    // Check for specific errors
    var httpErr *typesense.HTTPError
    if errors.As(err, &httpErr) {
        if httpErr.Status == 404 {
            return nil, ErrCollectionNotFound
        }
    }
    return nil, fmt.Errorf("search failed: %w", err)
}
```

## DO's and DON'Ts

### DO

- ✅ Use `upsert` for idempotent indexing
- ✅ Use batch imports for bulk operations
- ✅ Set `Facet: true` for filterable fields
- ✅ Use typed document operations for type safety
- ✅ Index documents asynchronously via River jobs
- ✅ One collection per content module
- ✅ Use circuit breaker options in production

### DON'T

- ❌ Index synchronously in HTTP handlers
- ❌ Share collections between modules
- ❌ Use `*` queries without filters in production
- ❌ Forget to handle 404 errors for missing documents
- ❌ Create collections without schema validation
---

## Related

- [INDEX.instructions.md](INDEX.instructions.md) - Main instruction index with all cross-references
- [ARCHITECTURE_V2.md](../../docs/dev/design/architecture/ARCHITECTURE_V2.md) - System architecture
- [BEST_PRACTICES.md](../../docs/dev/design/operations/BEST_PRACTICES.md) - Best practices
