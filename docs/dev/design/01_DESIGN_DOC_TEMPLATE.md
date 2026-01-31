# Design Document Template

> Use this template for all feature and integration design documents.

**Key Principle**: Never duplicate information from [00_SOURCE_OF_TRUTH.md](00_SOURCE_OF_TRUTH.md). Reference it instead.

---

## Document Structure

Every design doc should follow this structure for consistency and actionability:

```
# Feature Name

> One-line description

## Status

| Dimension | Status | Notes |
|-----------|--------|-------|
| Design | 🔴 | |
| Sources | 🔴 | |
| Instructions | 🔴 | |
| Code | 🔴 | |
| Linting | 🔴 | |
| Unit Testing | 🔴 | |
| Integration Testing | 🔴 | |

**Priority**: 🔴 HIGH | 🟡 MEDIUM | 🟢 LOW
**Module**: `internal/content/{module}` or `internal/infra/{component}`
**Dependencies**: [Link to related docs]

---

## Overview
Brief explanation of what this feature does and why it exists.

## Goals
- Clear, measurable objectives
- What success looks like

## Non-Goals
- Explicitly state what this feature does NOT do
- Prevents scope creep

---

## Technical Design

### Database Schema
```sql
-- Complete, runnable SQL
CREATE TABLE ...
```

### Repository Interface
```go
// Interface definition
type Repository interface {
    Method(ctx context.Context, ...) (Result, error)
}
```

### Service Layer
```go
// Service struct and key methods
type Service struct { ... }
func (s *Service) DoThing(ctx context.Context, ...) error { ... }
```

### API Endpoints
```
METHOD /api/path
Request: { ... }
Response: { ... }
```

---

## Implementation

### Files to Create/Modify
| File | Action | Description |
|------|--------|-------------|
| `path/to/file.go` | CREATE | Service implementation |
| `path/to/file.sql` | CREATE | Database migration |

### SQL Queries (sqlc)
```sql
-- name: QueryName :one
SELECT ... FROM ... WHERE ...
```

### River Jobs (if applicable)
```go
type JobArgs struct { ... }
func (JobArgs) Kind() string { return "module.job_name" }
```

---

## Configuration
```yaml
feature:
  enabled: true
  option: value
```

---

## Testing Strategy
- Unit tests: ...
- Integration tests: ...
- E2E tests: ...

---

## Migration Path
How to upgrade from current state (if applicable).

---

## Checklist
- [ ] Database migration created
- [ ] sqlc queries written
- [ ] Repository implemented
- [ ] Service implemented
- [ ] API handlers created
- [ ] Tests written
- [ ] Documentation updated

---


<!-- SOURCE-BREADCRUMBS-START -->

## Sources & Cross-References

> Auto-generated section linking to external documentation sources

### Cross-Reference Indexes

- [All Sources Index](../sources/SOURCES_INDEX.md) - Complete list of external documentation
- [Design ↔ Sources Map](../sources/DESIGN_CROSSREF.md) - Which docs reference which sources

<!-- SOURCE-BREADCRUMBS-END -->

## Related Documents
- [Link to related doc](path/to/doc.md)
```

---

## Key Principles

### 1. Code-First Examples
Always provide runnable code, not pseudocode:
- Complete SQL schemas with indexes
- Full Go interface definitions
- Real API request/response examples

### 2. File Locations
Explicitly state where code should go:
```
internal/
  content/
    {module}/
      entity.go        # Domain types
      repository.go    # Interface
      repository_pg.go # PostgreSQL implementation
      service.go       # Business logic
      jobs.go          # River workers
      module.go        # fx module
```

### 3. Implementation Order
Documents should guide implementation sequence:
1. Database schema → Migration
2. sqlc queries → Repository
3. Repository → Service
4. Service → API Handlers
5. API Handlers → Tests

### 4. Status Tracking

Use the multi-dimensional status table (see [00_SOURCE_OF_TRUTH.md](00_SOURCE_OF_TRUTH.md#status-system)):

| Emoji | Meaning |
|-------|---------|
| ✅ | Complete |
| 🟡 | Partial |
| 🔴 | Not Started |
| ⚪ | N/A |

**Dimensions tracked:**
- **Design** - Feature spec, architecture, DB schema, API endpoints
- **Sources** - External docs fetched (API specs, GraphQL schemas)
- **Instructions** - Claude Code implementation instructions
- **Code** - Go implementation
- **Linting** - golangci-lint rules
- **Unit Testing** - Unit tests with embedded-postgres
- **Integration Testing** - Integration tests with testcontainers

### 5. Dependency Links
Always link to:
- Related feature docs
- Relevant integration docs
- Shared component docs

---

## Example: Minimal Feature Doc

```markdown
# User Favorites

> Allow users to mark content as favorites

## Status

| Dimension | Status | Notes |
|-----------|--------|-------|
| Design | ✅ | Schema, API, Repository defined |
| Sources | ⚪ | N/A - internal feature |
| Instructions | 🔴 | |
| Code | 🔴 | |
| Linting | 🔴 | |
| Unit Testing | 🔴 | |
| Integration Testing | 🔴 | |

**Module**: `internal/content/movie`, `internal/content/tvshow`

---

## Database

```sql
CREATE TABLE user_favorites (
    user_id UUID NOT NULL REFERENCES users(id),
    content_id UUID NOT NULL,
    content_type VARCHAR(20) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, content_id, content_type)
);
```

## Repository

```go
type FavoritesRepository interface {
    IsFavorite(ctx context.Context, userID, contentID uuid.UUID, contentType string) (bool, error)
    AddFavorite(ctx context.Context, userID, contentID uuid.UUID, contentType string) error
    RemoveFavorite(ctx context.Context, userID, contentID uuid.UUID, contentType string) error
    ListFavorites(ctx context.Context, userID uuid.UUID, contentType string, limit, offset int) ([]uuid.UUID, error)
}
```

## API

```
POST /api/users/{userId}/favorites
{ "content_id": "uuid", "content_type": "movie" }

DELETE /api/users/{userId}/favorites/{contentId}

GET /api/users/{userId}/favorites?type=movie&limit=20
```

## Checklist
- [ ] Migration: `000X_user_favorites.sql`
- [ ] sqlc: `queries/shared/favorites.sql`
- [ ] Repository: per-module implementation
- [ ] Service: per-module favorites methods
- [ ] API: `/api/users/{userId}/favorites`
```

---

## Anti-Patterns to Avoid

1. **Vague descriptions** - Be specific about behavior
2. **Missing error handling** - Document error cases
3. **No code examples** - Always provide runnable code
4. **Orphan docs** - Always link to related docs
5. **Stale status** - Update status as work progresses
6. **No checklist** - Always include implementation checklist
7. **Duplicating SOT** - Never copy package versions, module lists, or config keys from 00_SOURCE_OF_TRUTH.md - reference it instead

## Deduplication Rules

| Information Type | Where It Lives | In Design Docs |
|-----------------|----------------|----------------|
| Package versions | 00_SOURCE_OF_TRUTH.md | Link to SOT |
| Module/service list | 00_SOURCE_OF_TRUTH.md | Link to SOT |
| Config keys | 00_SOURCE_OF_TRUTH.md | Link to SOT |
| API namespaces | 00_SOURCE_OF_TRUTH.md | Link to SOT |
| QAR terminology | 00_SOURCE_OF_TRUTH.md | Link to SOT |
| Design rationale | Design docs | Full detail |
| Implementation patterns | Design docs | Full detail |
| Business logic | Design docs | Full detail |
