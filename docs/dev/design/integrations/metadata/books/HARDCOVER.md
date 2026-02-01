## Table of Contents

- [Hardcover](#hardcover)
  - [Status](#status)
  - [Architecture](#architecture)
    - [Integration Structure](#integration-structure)
    - [Data Flow](#data-flow)
    - [Provides](#provides)
  - [Implementation](#implementation)
    - [Key Interfaces](#key-interfaces)
    - [Dependencies](#dependencies)
  - [Configuration](#configuration)
    - [Environment Variables](#environment-variables)
    - [Config Keys](#config-keys)
  - [Related Documentation](#related-documentation)
    - [Design Documents](#design-documents)
    - [External Sources](#external-sources)

# Hardcover


**Created**: 2026-01-31
**Status**: ✅ Complete
**Category**: integration


> Integration with Hardcover

> Book reading tracker and scrobbling - Goodreads alternative with GraphQL API
**API Base URL**: `https://api.hardcover.app/v1/graphql`
**Authentication**: oauth2

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

```
┌──────────────────┐                    ┌──────────────────┐
│  Revenge         │                    │   Hardcover      │
│  Book Library    │◀──── Two-Way ────▶│   Account        │
└────────┬─────────┘       Sync        └────────┬─────────┘
         │                                       │
         │ Reading Events                        │
         ▼                                       ▼
┌──────────────────┐                    ┌──────────────────┐
│   River Queue    │                    │   User Shelves   │
│   (background)   │                    │   - Want to Read │
└────────┬─────────┘                    │   - Reading      │
         │                              │   - Read         │
         ▼                              └──────────────────┘
┌──────────────────┐
│  Hardcover API   │
│  (GraphQL)       │
└──────────────────┘

Sync Types:
- Reading progress → Hardcover
- Book completion → Hardcover
- Shelf changes ← Hardcover
- Ratings ← Hardcover (enrichment)
```


### Integration Structure

```
internal/integration/hardcover/
├── client.go              # API client
├── types.go               # Response types
├── mapper.go              # Map external → internal types
├── cache.go               # Response caching
└── client_test.go         # Tests
```

### Data Flow

<!-- Data flow diagram -->

### Provides
<!-- Data provided by integration -->


## Implementation

### Key Interfaces

```go
// Hardcover scrobbler
type HardcoverScrobbler struct {
  client      *graphql.Client
  tokenStore  TokenStore
  queue       *river.Client
}

// Book scrobbling interface
type BookScrobbler interface {
  Connect(ctx context.Context, userID uuid.UUID, authCode string) error
  Disconnect(ctx context.Context, userID uuid.UUID) error
  SyncProgress(ctx context.Context, userID uuid.UUID, bookID uuid.UUID, progress float64) error
  MarkAsRead(ctx context.Context, userID uuid.UUID, bookID uuid.UUID, rating *int) error
  ImportShelves(ctx context.Context, userID uuid.UUID) ([]*BookShelf, error)
  ExportToShelf(ctx context.Context, userID uuid.UUID, bookID uuid.UUID, shelf string) error
}

// OAuth token storage
type TokenStore interface {
  GetToken(ctx context.Context, userID uuid.UUID) (*oauth2.Token, error)
  SaveToken(ctx context.Context, userID uuid.UUID, token *oauth2.Token) error
  DeleteToken(ctx context.Context, userID uuid.UUID) error
}
```


### Dependencies
**Go Packages**:
- `github.com/Khan/genqlient` - Type-safe GraphQL client
- `golang.org/x/oauth2` - OAuth 2.0
- `github.com/riverqueue/river` - Background sync jobs
- `github.com/jackc/pgx/v5` - PostgreSQL
- `go.uber.org/fx` - DI

**External**:
- Hardcover API (OAuth 2.0 required)







## Configuration

### Environment Variables

```bash
HARDCOVER_CLIENT_ID=your_client_id
HARDCOVER_CLIENT_SECRET=your_client_secret
HARDCOVER_REDIRECT_URI=https://your-revenge-server/api/v1/integrations/hardcover/callback
```


### Config Keys
```yaml
scrobbling:
  hardcover:
    enabled: true
    client_id: ${HARDCOVER_CLIENT_ID}
    client_secret: ${HARDCOVER_CLIENT_SECRET}
    redirect_uri: ${HARDCOVER_REDIRECT_URI}
    sync:
      interval: 30m
      direction: bidirectional    # 'to_hardcover', 'from_hardcover', 'bidirectional'
```










## Related Documentation
### Design Documents
- [03_METADATA_SYSTEM](../../../architecture/03_METADATA_SYSTEM.md)
- [BOOK_MODULE](../../../features/book/BOOK_MODULE.md)
- [SCROBBLING](../../../features/shared/SCROBBLING.md)
- [CHAPTARR](../../servarr/CHAPTARR.md)

### External Sources
- [Hardcover API](../../../../sources/apis/hardcover.md) - Auto-resolved from hardcover
- [Khan/genqlient](../../../../sources/tooling/genqlient.md) - GraphQL client
- [golang.org/x/oauth2](https://pkg.go.dev/golang.org/x/oauth2) - OAuth 2.0

