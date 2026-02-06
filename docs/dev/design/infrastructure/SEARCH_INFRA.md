# Search Infrastructure

<!-- DESIGN: infrastructure -->

**Package**: `internal/infra/search`
**fx Module**: Part of `search.Module` (shared with service layer)

> Typesense client wrapper with graceful degradation and retry-based startup

---

## Service Structure

```
internal/infra/search/
└── module.go              # Client wrapper, collection/document ops, health, fx hooks
```

## Client Interface

```go
type Client struct {
    client *typesense.Client   // typesense-go SDK
    config config.SearchConfig
    logger *slog.Logger
}

func NewClient(cfg config.SearchConfig, logger *slog.Logger) *Client
func (c *Client) IsEnabled() bool
```

**Nil-safe design**: All operations check `IsEnabled()` first. When search is disabled, methods return gracefully without errors.

## Collection Operations

```go
func (c *Client) CreateCollection(ctx context.Context, schema *api.CollectionSchema) error
func (c *Client) DeleteCollection(ctx context.Context, name string) error
func (c *Client) GetCollection(ctx context.Context, name string) (*api.Collection, error)
func (c *Client) ListCollections(ctx context.Context) ([]*api.Collection, error)
```

## Document Operations

```go
func (c *Client) IndexDocument(ctx context.Context, collection string, document interface{}) error
func (c *Client) UpdateDocument(ctx context.Context, collection, id string, document interface{}) error
func (c *Client) DeleteDocument(ctx context.Context, collection, id string) error
func (c *Client) BulkImport(ctx context.Context, collection string, documents []interface{}, batchSize int) error
```

## Search Operations

```go
func (c *Client) Search(ctx context.Context, collection string, params *api.SearchCollectionParams) (*api.SearchResult, error)
func (c *Client) MultiSearch(ctx context.Context, params *api.MultiSearchParams) (*api.MultiSearchResult, error)
```

## Startup Behavior

fx lifecycle hooks handle startup with retry:
1. Attempts health check with 10-second timeout
2. Retries up to 5 times with exponential backoff
3. Logs warnings but does **not** fail startup if Typesense is unavailable
4. Service degrades gracefully (search disabled, other features work)

## Configuration

From `config.go` `SearchConfig` (koanf namespace `search.*`):
```yaml
search:
  url: http://localhost:8108    # Typesense server URL
  api_key: ""                   # Typesense API key
  enabled: false                # Disabled by default
```

## Dependencies

- `github.com/typesense/typesense-go` - Typesense Go client SDK
- `log/slog` - Structured logging

## Related Documentation

- [../services/SEARCH.md](../services/SEARCH.md) - MovieSearchService (uses this client)
- [CACHE.md](CACHE.md) - CachedMovieSearchService caches search results
