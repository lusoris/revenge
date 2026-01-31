# Example Feature

<!-- BREADCRUMB: [Design Index](../DESIGN_INDEX.md) > [Features - Example](../features/example/INDEX.md) > Example Feature -->

<!-- SOURCES: fx, pgx, sqlc -->

<!-- DESIGN: 01_ARCHITECTURE, 02_DESIGN_PRINCIPLES -->


---

## Status

| Dimension | Status | Notes |
|-----------|--------|-------|
| Design | 🟡 | Scaffold - needs completion |
| Sources | ✅ | All sources documented |
| Instructions | 🔴 | Implementation checklist needed |
| Code | 🔴 | Not started |
| Linting | 🔴 | N/A - no code yet |
| Unit Testing | 🔴 | N/A - no code yet |
| Integration Testing | 🔴 | N/A - no code yet |

**Location**: `internal/example/`

**Last Updated**: 2026-01-31

---

## Overview


Technical example feature for testing template generation.

**Purpose**: Demonstrate template usage and validation

**Architecture Pattern**: Repository pattern with service layer


---

## Architecture

**Reference**: See [00_SOURCE_OF_TRUTH.md](../../00_SOURCE_OF_TRUTH.md) for versions and dependencies.

### Component Structure

```
┌─────────────────────────────────┐
│        API Layer (ogen)         │
└───────────┬─────────────────────┘
            │
┌───────────▼─────────────────────┐
│    Example Service (cache)      │
└───────────┬─────────────────────┘
            │
┌───────────▼─────────────────────┐
│  Repository (PostgreSQL+sqlc)   │
└─────────────────────────────────┘

```

### Key Components

- **ExampleService**: Main business logic
  - Location: `internal/example/service.go`
  - Dependencies: Repository, Cache
- **ExampleRepository**: Data access layer
  - Location: `internal/example/repository.go`
  - Dependencies: PostgreSQL

### Design Patterns

- **Repository Pattern**: Separates business logic from data access
- **Dependency Injection**: Uses fx for DI

### Data Flow

```
User Request → API Handler → Service → Repository → Database

```

---

## Database Schema

**Reference**: All schemas defined in SOURCE_OF_TRUTH.md database section.

### Tables

#### `examples`

Example entities

**Columns**:
```sql
CREATE TABLE examples (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name TEXT NOT NULL,
  created_at TIMESTAMPTZ DEFAULT NOW()
);

```

**Indexes**:
- `idx_examples_name`

**Relationships**:
- No foreign keys (simple table)

### Migrations

**Location**: `internal/database/migrations/`

**Migration Files**:
- `20260131000000_create_examples.sql` - Create examples table

---

## API Endpoints

**Reference**: All endpoints defined in OpenAPI specs (see API.md).

**Namespace**: `/api/v1/examples`

### Endpoints

#### GET `/api/v1/examples`

List all examples

**Request**:
```json
GET /api/v1/examples?limit=10

```

**Response**:
```json
{
  "examples": [
    {"id": "uuid", "name": "Example 1"}
  ]
}

```

**Authentication**: Bearer token

**RBAC Scope**: `example:read`

**Rate Limit**: 100 req/min

#### POST `/api/v1/examples`

Create new example

**Request**:
```json
{
  "name": "New Example"
}

```

**Response**:
```json
{
  "id": "uuid",
  "name": "New Example"
}

```

**Authentication**: Bearer token

**RBAC Scope**: `example:write`

**Rate Limit**: 10 req/min


---

## External Integrations

### Example API

**Purpose**: Fetch external data

**API Documentation**: See [EXAMPLE_API.md](../../integrations/example/EXAMPLE_API.md)

**Integration Points**:
- Webhook endpoint for updates
- REST API for data sync

**Data Sync Strategy**: Real-time via webhooks

**Error Handling**: Retry with exponential backoff

**Rate Limiting**: 100 req/min with token bucket


---

## Business Logic

### Core Rules

- Example names must be unique
- Examples cannot be deleted if in use

### Validation

- **name**: Required, max 255 chars

### State Transitions

```
PENDING → ACTIVE → ARCHIVED

```

---

## Caching Strategy

**L1 Cache** (otter - in-memory):
- Cache full example entities by ID
- TTL: 5m

**L2 Cache** (rueidis - distributed):
- Cache list queries results
- TTL: 1h

**Cache Invalidation**:
- On example create/update/delete
- On manual cache clear

---

## Testing Strategy

**Coverage Target**: 80%+ (per SOURCE_OF_TRUTH.md standards)

### Unit Tests

**Location**: `internal/example/service_test.go`

**Test Cases**:
- TestExampleService_Create
- TestExampleService_Get
- TestExampleRepository_List

**Mocking**:
- Mock Repository for service tests
- Mock Database for repository tests

### Integration Tests

**Location**: `internal/example/integration_test.go`

**Test Scenarios**:
- Create example via API
- List examples with pagination
- Delete example cascades correctly

**Dependencies**:
- testcontainers/postgres:18

### E2E Tests


---

## Security Considerations

**Authentication**: Bearer token required for all endpoints

**Authorization**: RBAC scopes: example:read, example:write, example:delete

**RBAC Scopes**:
- `example:read`: Read examples
- `example:write`: Create/update examples
- `example:delete`: Delete examples (admin only)

**Sensitive Data**:
- name: Public field, no special handling

**Input Validation**:
- Name: max 255 chars, no special chars
- ID: valid UUID format

**Security Best Practices**:
- Always validate user input
- Use parameterized queries (sqlc)
- Rate limit all endpoints

---

## Performance Considerations

**Expected Load**: 1000 req/min peak

**Query Optimization**:
- Index on name column for search
- Use LIMIT/OFFSET for pagination

**Monitoring Metrics**:
- example_requests_total: Total API requests
- example_cache_hit_ratio: L1/L2 cache effectiveness

**Potential Bottlenecks**:
- Database queries for large lists: Pagination + caching

---

## Implementation Checklist

**Reference**: See SOURCE_OF_TRUTH.md for all dependency versions.

### Phase 1: Core Implementation
- [ ] Define entity struct
- [ ] Create repository interface
- [ ] Implement PostgreSQL repository
- [ ] Create migrations

### Phase 2: Integration
- [ ] Implement service layer
- [ ] Add caching (otter + rueidis)
- [ ] Create API handlers (ogen)
- [ ] Add RBAC middleware

### Phase 3: Testing & Polish
- [ ] Write unit tests (80%+ coverage)
- [ ] Write integration tests
- [ ] Add monitoring metrics
- [ ] Document API in OpenAPI

**Go Packages Required**:
- `go.uber.org/fx` - Dependency injection (version: v1.23.0 per SOT)
- `github.com/jackc/pgx/v5` - PostgreSQL driver (version: v5.7.2 per SOT)

---

## Dependencies

**Depends on**:
- [Auth Service](../../services/AUTH.md) - Requires authentication
- [Database](../../00_SOURCE_OF_TRUTH.md) - Uses PostgreSQL 18+

**Blocks**:

**Related to**:
- [Architecture](../../architecture/01_ARCHITECTURE.md)

---


## Source Documentation

**From `docs/dev/sources/`**:
- [fx Documentation](../../../../sources/tooling/fx.md) - Dependency injection patterns
- [pgx Documentation](../../../../sources/database/pgx.md) - PostgreSQL driver usage

---

## Cross-References

**Design Docs**:
- [Repository Pattern](../../patterns/REPOSITORY_PATTERN.md)

**External Sources**:
- [fx GitHub](https://github.com/uber-go/fx)

---

## Implementation Notes

### Caching Strategy

Use otter for L1 (in-memory), rueidis for L2 (distributed). See SOURCE_OF_TRUTH for TTL values.
### Testing Approach

Use testcontainers for integration tests. Mock repository for unit tests.

---

## Open Questions

- Should we support bulk operations?
- What's the maximum pagination limit?

---

**Status**: 🟡 SCAFFOLD - Needs completion
**Next Step**: Complete database schema and API endpoints
