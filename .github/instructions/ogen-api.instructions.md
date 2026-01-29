# ogen API Instructions

> Source: https://ogen.dev/docs/intro, https://github.com/ogen-go/ogen

Apply to: `**/api/**/*.go`, `**/api/openapi/**/*.yaml`, `**/api/generated/**/*.go`

## Overview

ogen is an OpenAPI v3 code generator for Go. It generates type-safe clients and servers from OpenAPI specifications.

**Key features:**

- No reflection or `interface{}`
- Code-generated JSON encoding (using go-faster/jx)
- Code-generated validation
- Static radix router
- Generated Optional[T], Nullable[T] wrappers
- Sum types for oneOf
- OpenTelemetry support

## Installation

```bash
go install -v github.com/ogen-go/ogen/cmd/ogen@latest
```

## Code Generation

### Setup Generator

Create `generate.go`:

```go
package api

//go:generate go run github.com/ogen-go/ogen/cmd/ogen --target generated --package api --clean ../api/openapi/revenge.yaml
```

Run generation:

```bash
go generate ./api/...
```

### Generated Files

```
api/generated/
├── oas_cfg_gen.go           # Configuration
├── oas_client_gen.go        # HTTP client
├── oas_handlers_gen.go      # Handler interfaces
├── oas_interfaces_gen.go    # Interfaces
├── oas_json_gen.go          # JSON encoding
├── oas_middleware_gen.go    # Middleware
├── oas_parameters_gen.go    # Parameter types
├── oas_request_decoders_gen.go
├── oas_request_encoders_gen.go
├── oas_response_decoders_gen.go
├── oas_response_encoders_gen.go
├── oas_router_gen.go        # Static router
├── oas_schemas_gen.go       # Schema types
├── oas_server_gen.go        # Server
├── oas_unimplemented_gen.go # Stub handler
├── oas_validators_gen.go    # Validation
```

## OpenAPI Specification

### Basic Structure

```yaml
openapi: 3.1.0
info:
  title: Revenge API
  version: 1.0.0
servers:
  - url: /api/v1
paths:
  /movies:
    get:
      operationId: listMovies
      responses:
        "200":
          description: Success
          content:
            application/json:
              schema:
                type: array
                items:
                  $ref: "#/components/schemas/Movie"
components:
  schemas:
    Movie:
      type: object
      required:
        - id
        - title
      properties:
        id:
          type: string
          format: uuid
        title:
          type: string
```

## Server Implementation

### Handler Interface

ogen generates a `Handler` interface:

```go
// Handler handles operations described by OpenAPI v3 specification.
type Handler interface {
    // ListMovies implements listMovies operation.
    // GET /movies
    ListMovies(ctx context.Context) ([]Movie, error)

    // GetMovie implements getMovie operation.
    // GET /movies/{id}
    GetMovie(ctx context.Context, params GetMovieParams) (GetMovieRes, error)

    // CreateMovie implements createMovie operation.
    // POST /movies
    CreateMovie(ctx context.Context, req *CreateMovieReq) (*Movie, error)
}
```

### Implement Handler

```go
type MovieHandler struct {
    service *movie.Service
    logger  *slog.Logger
}

func (h *MovieHandler) ListMovies(ctx context.Context) ([]api.Movie, error) {
    movies, err := h.service.List(ctx)
    if err != nil {
        return nil, fmt.Errorf("list movies: %w", err)
    }

    result := make([]api.Movie, len(movies))
    for i, m := range movies {
        result[i] = toAPIMovie(m)
    }
    return result, nil
}

func (h *MovieHandler) GetMovie(ctx context.Context, params api.GetMovieParams) (api.GetMovieRes, error) {
    movie, err := h.service.Get(ctx, params.ID)
    if errors.Is(err, domain.ErrMovieNotFound) {
        return &api.GetMovieNotFound{}, nil
    }
    if err != nil {
        return nil, err
    }

    result := toAPIMovie(movie)
    return &result, nil
}
```

### Create Server

```go
func main() {
    handler := &MovieHandler{...}

    srv, err := api.NewServer(handler)
    if err != nil {
        log.Fatal(err)
    }

    if err := http.ListenAndServe(":8080", srv); err != nil {
        log.Fatal(err)
    }
}
```

## Client Usage

```go
client, err := api.NewClient("http://localhost:8080")
if err != nil {
    log.Fatal(err)
}

// Type-safe API call
movie, err := client.GetMovie(ctx, api.GetMovieParams{ID: movieID})
if err != nil {
    log.Fatal(err)
}

switch v := movie.(type) {
case *api.Movie:
    fmt.Println(v.Title)
case *api.GetMovieNotFound:
    fmt.Println("Movie not found")
}
```

## Optional and Nullable Types

ogen generates wrapper types instead of pointers:

```go
// OptNilString is optional nullable string
type OptNilString struct {
    Value string
    Set   bool
    Null  bool
}

// Helper methods
func (o OptNilString) Get() (v string, ok bool)
func (o OptNilString) IsNull() bool
func (o OptNilString) IsSet() bool

// Constructor
func NewOptNilString(v string) OptNilString
```

### Usage

```go
// Check if optional value is set
if params.Name.IsSet() {
    name := params.Name.Value
}

// Create optional value
req := &api.UpdateMovieReq{
    Title: api.NewOptString("New Title"),
}
```

## Sum Types (oneOf)

For oneOf schemas:

```go
type ID struct {
    Type   IDType
    String string
    Int    int
}

// Constructors
func NewStringID(v string) ID
func NewIntID(v int) ID
```

### Discriminator Inference

ogen automatically infers discrimination:

1. **Type-based**: Different JSON types
2. **Explicit discriminator**: `propertyName` defined
3. **Field-based**: Unique fields per variant
4. **Value-based**: Different enum values

## Extension Properties

### Custom Type Name

```yaml
components:
  schemas:
    MySchema:
      x-ogen-name: CustomTypeName
      type: object
```

### Custom Field Name

```yaml
components:
  schemas:
    Node:
      type: object
      properties:
        parent:
          $ref: "#/components/schemas/Node"
      x-ogen-properties:
        parent:
          name: "Prev"
```

### Operation Groups

```yaml
paths:
  /movies:
    x-ogen-operation-group: Movies
    get:
      operationId: listMovies
```

Generates separate handler interface:

```go
type MoviesHandler interface {
    ListMovies(ctx context.Context) ([]Movie, error)
}

type Handler interface {
    MoviesHandler
    // Other grouped handlers...
}
```

### Extra Struct Tags

```yaml
properties:
  id:
    type: integer
    x-oapi-codegen-extra-tags:
      gorm: primaryKey
      valid: customValidator
```

### Streaming JSON

```yaml
requestBody:
  content:
    application/json:
      x-ogen-json-streaming: true
      schema:
        type: array
        items:
          type: number
```

## Revenge Patterns

### Handler Structure

```go
// internal/api/handlers/movie.go
type MovieHandler struct {
    service *movie.Service
    logger  *slog.Logger
}

func NewMovieHandler(service *movie.Service, logger *slog.Logger) *MovieHandler {
    return &MovieHandler{service: service, logger: logger}
}
```

### Composite Handler

```go
// internal/api/handlers/handler.go
type Handler struct {
    *MovieHandler
    *TVShowHandler
    *MusicHandler
    // ... other module handlers
}

func NewHandler(
    movie *MovieHandler,
    tvshow *TVShowHandler,
    music *MusicHandler,
) *Handler {
    return &Handler{
        MovieHandler:  movie,
        TVShowHandler: tvshow,
        MusicHandler:  music,
    }
}
```

### Error Responses

```go
func (h *MovieHandler) GetMovie(ctx context.Context, params api.GetMovieParams) (api.GetMovieRes, error) {
    movie, err := h.service.Get(ctx, params.ID)

    // Domain error → API response type
    if errors.Is(err, domain.ErrNotFound) {
        return &api.GetMovieNotFound{}, nil
    }

    // Unexpected error → returns 500
    if err != nil {
        h.logger.Error("get movie failed", "error", err, "id", params.ID)
        return nil, err
    }

    return toAPIMovie(movie), nil
}
```

### OpenAPI File Structure

```
api/openapi/
├── revenge.yaml          # Main spec (refs other files)
├── movies.yaml           # Movie endpoints
├── tvshows.yaml          # TV show endpoints
├── music.yaml            # Music endpoints
├── auth.yaml             # Authentication
├── components/
│   ├── schemas.yaml      # Shared schemas
│   ├── parameters.yaml   # Shared parameters
│   └── responses.yaml    # Shared responses
```

## DO's and DON'Ts

### DO

- ✅ Use `operationId` for all endpoints
- ✅ Use `x-ogen-operation-group` for logical grouping
- ✅ Handle all response types in switch statements
- ✅ Use Optional[T] helpers (`IsSet()`, `Get()`)
- ✅ Keep OpenAPI specs modular with `$ref`
- ✅ Regenerate after spec changes: `go generate ./api/...`
- ✅ Return domain errors as typed responses (404, 400, etc.)

### DON'T

- ❌ Edit generated files (they will be overwritten)
- ❌ Use `interface{}` in handlers
- ❌ Ignore the generated validation
- ❌ Mix business logic in handlers
- ❌ Return generic errors for expected failures
- ❌ Use pointers when Optional[T] is available
---

## Related

- [INDEX.instructions.md](INDEX.instructions.md) - Main instruction index with all cross-references
- [ARCHITECTURE_V2.md](../../docs/dev/design/architecture/ARCHITECTURE_V2.md) - System architecture
- [BEST_PRACTICES.md](../../docs/dev/design/operations/BEST_PRACTICES.md) - Best practices
