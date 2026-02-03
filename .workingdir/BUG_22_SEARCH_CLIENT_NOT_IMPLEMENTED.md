# Bug #22: Search Client Not Implemented - Stub Only

**Status**: 🔴 CRITICAL - FOUND MISSING FEATURE
**Severity**: High (Feature Gap)
**Component**: Search Infrastructure (`internal/infra/search`)
**Discovered**: 2025-01-XX during comprehensive integration testing

## Symptom

Tests attempting to use search functionality fail to compile because the search client is a stub with no methods implemented:

```
client.CreateIndex undefined
client.DeleteIndex undefined
client.IndexDocument undefined
client.Search undefined
```

## Root Cause

The search client at `internal/infra/search/module.go` is a **placeholder stub** left from v0.1.0 skeleton:

```go
// Client represents the Typesense search client.
// This is a placeholder stub for v0.1.0 skeleton.
type Client struct {
    config *config.Config
    logger *slog.Logger
}
```

**No actual Typesense integration exists**:
- No client initialization
- No index management methods
- No document indexing methods
- No search query methods
- No collection management

## Impact

- **Critical**: Search functionality completely non-functional
- Typesense service running but not being used
- Configuration exists (`REVENGE_SEARCH_URL`, `REVENGE_SEARCH_ENABLED`) but ignored
- Application starts successfully but search features don't work
- Silent failure - no errors, just missing functionality

## Current State

**What Exists**:
- ✅ Docker service running (Typesense 0.25.2)
- ✅ Configuration structure in place
- ✅ Module registration in DI container
- ✅ Lifecycle hooks (start/stop)
- ✅ Healthcheck working

**What's Missing**:
- ❌ Actual Typesense client library integration
- ❌ Collection/index creation methods
- ❌ Document CRUD operations
- ❌ Search query methods
- ❌ Type mappings and schema definitions
- ❌ Error handling for search operations
- ❌ Connection pooling/management

## Expected Implementation

The search client should provide:

### Index Management
```go
func (c *Client) CreateIndex(ctx context.Context, name string, schema interface{}) error
func (c *Client) DeleteIndex(ctx context.Context, name string) error
func (c *Client) GetIndex(ctx context.Context, name string) (*Collection, error)
func (c *Client) ListIndexes(ctx context.Context) ([]Collection, error)
```

### Document Operations
```go
func (c *Client) IndexDocument(ctx context.Context, collection string, doc interface{}) error
func (c *Client) UpdateDocument(ctx context.Context, collection, id string, doc interface{}) error
func (c *Client) DeleteDocument(ctx context.Context, collection, id string) error
func (c *Client) GetDocument(ctx context.Context, collection, id string) (interface{}, error)
```

### Search Operations
```go
func (c *Client) Search(ctx context.Context, collection, query string, params *SearchParams) (*SearchResults, error)
func (c *Client) MultiSearch(ctx context.Context, queries []SearchQuery) (*MultiSearchResults, error)
```

## Recommended Solution

### Option 1: Use Official Typesense Go Client ✅ RECOMMENDED
Use `github.com/typesense/typesense-go`:

```go
import "github.com/typesense/typesense-go/typesense"

type Client struct {
    client *typesense.Client
    config *config.Config
    logger *slog.Logger
}

func NewClient(cfg *config.Config, logger *slog.Logger) (*Client, error) {
    client := typesense.NewClient(
        typesense.WithServer(cfg.Search.URL),
        typesense.WithAPIKey(cfg.Search.APIKey),
    )
    return &Client{
        client: client,
        config: cfg,
        logger: logger,
    }, nil
}
```

### Option 2: Custom HTTP Client
Implement Typesense REST API directly using HTTP client:
- More control over requests
- No external dependencies
- More work to maintain

### Option 3: Defer Implementation
Document as "coming soon" and disable search features:
- Not recommended - infrastructure already set up
- Wastes existing Typesense service
- Configuration already in place

## Implementation Checklist

1. ✅ Typesense service running
2. ✅ Configuration structure
3. ❌ Add `typesense-go` dependency: `go get github.com/typesense/typesense-go`
4. ❌ Implement `NewClient()` with actual Typesense client
5. ❌ Add collection management methods
6. ❌ Add document indexing methods
7. ❌ Add search query methods
8. ❌ Add error handling and retries
9. ❌ Update lifecycle hooks with real connection logic
10. ❌ Create integration tests
11. ❌ Document search API usage

## Priority

**HIGH** - This is infrastructure that's already deployed and configured but completely non-functional. The application pretends to have search capabilities but doesn't.

## Next Steps

1. Document this as missing feature (DONE)
2. Add `typesense-go` dependency
3. Implement full client with all CRUD operations
4. Create comprehensive integration tests
5. Test with real data indexing
6. Performance test with large datasets

## Related Code

- Stub: `internal/infra/search/module.go`
- Config: `internal/config/config.go` (search section)
- Docker: `docker-compose.dev.yml` (typesense service)
- Tests Attempted: `tests/integration/search/search_test.go` (compile failures)

## Design Considerations

When implementing, consider:
- **Connection Management**: Connection pooling, retries, timeouts
- **Schema Versioning**: How to handle schema changes
- **Batch Operations**: Bulk indexing for efficiency
- **Error Handling**: Graceful degradation if Typesense unavailable
- **Type Safety**: Strong typing for documents and queries
- **Testing**: Mock client for unit tests, real client for integration tests

## Conclusion

This is a **feature gap**, not a bug - search functionality was planned but never implemented beyond the stub. The infrastructure is ready (Typesense running, config in place), just needs the actual client implementation.

**Recommendation**: Implement immediately using Option 1 (official client library) as it's the fastest path to working search functionality.
