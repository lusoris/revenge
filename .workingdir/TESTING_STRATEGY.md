# Revenge Testing Strategy

**Created**: 2026-02-03
**Status**: Active
**Infrastructure**: Docker Compose (PostgreSQL 18, Dragonfly, Typesense)

---

## Philosophy

**Test with REAL infrastructure, not mocks.** Integration tests catch real bugs that unit tests miss.

### Testing Pyramid (Inverted for Infrastructure Projects)

```
        /\
       /  \      E2E Tests (10%)
      /----\
     /      \    Integration Tests (60%) ← FOCUS HERE
    /--------\
   /          \  Unit Tests (30%)
  /------------\
```

**Why inverted?** Infrastructure bugs happen at integration points:
- Database connection pool exhaustion
- Cache eviction under load
- Transaction deadlocks
- Network failures
- Memory leaks over time

---

## Test Layers

### Layer 1: Infrastructure Tests (REAL services via Docker)

**Goal**: Verify our code works with real PostgreSQL, Redis, Typesense

**Approach**:
1. Start Docker services
2. Run integration tests against them
3. Verify behavior under various conditions
4. Tear down

**Components to Test**:
- Database connection pool
- Cache layers (L1 + L2 with real Redis)
- Search indexing
- Job queue with real worker
- Health checks with real services

**Test Types**:
- ✅ Happy path
- ⚠️ Failure scenarios (service down, network issues)
- 🔥 Load testing (concurrent access, large data)
- 💥 Chaos (kill services mid-operation)

---

### Layer 2: Service Layer Tests

**Goal**: Verify business logic with real database

**Components**:
- User service (create, update, delete with real DB)
- Settings service
- Content service
- RBAC enforcement

**What to Test**:
- Transactions (commit/rollback)
- Concurrent modifications
- Unique constraint violations
- Foreign key cascades
- NULL handling (pgtype.UUID, pgtype.Text)

---

### Layer 3: API Layer Tests

**Goal**: Verify HTTP handlers with full stack

**Components**:
- Auth endpoints
- CRUD endpoints
- Error responses
- Rate limiting
- CORS

**What to Test**:
- Request validation
- Response serialization
- Error handling
- Authentication/authorization
- Context propagation

---

## Test Execution Strategy

### Phase 1: Infrastructure Layer (NOW)

**Priority**: HIGH - Foundation must be solid

#### 1.1 Database Integration Tests
```bash
# Start PostgreSQL
docker-compose -f docker-compose.dev.yml up -d postgres

# Run database tests
go test ./internal/infra/database/... -tags=integration -v

# Test scenarios:
# - Connection pool exhaustion (acquire all connections)
# - Transaction rollback on error
# - Prepared statement caching
# - Connection leak detection
# - Query timeout handling
# - pgtype nullable types (UUID, Text, Int4, etc.)
```

**Tests to Create**:
- `database_pool_test.go` - Connection pool behavior
- `database_transaction_test.go` - Transaction handling
- `database_concurrency_test.go` - Concurrent access patterns
- `database_failover_test.go` - Connection loss and recovery

#### 1.2 Cache Integration Tests (L2 with Real Redis)
```bash
# Start Dragonfly
docker-compose -f docker-compose.dev.yml up -d dragonfly

# Run cache tests with real Redis
REDIS_ADDR=localhost:6379 go test ./internal/infra/cache/... -tags=integration -v

# Test scenarios:
# - L1 + L2 coordination
# - Redis disconnection and reconnection
# - Pattern-based invalidation with SCAN
# - Pipeline operations
# - Large value storage (>1MB)
# - TTL accuracy
# - Pub/Sub functionality
```

**Tests to Create**:
- `cache_l2_integration_test.go` - L2 cache with real Redis
- `cache_unified_integration_test.go` - L1+L2 together
- `cache_failover_test.go` - Redis disconnect scenarios
- `cache_load_test.go` - High concurrency (1000+ goroutines)

#### 1.3 Search Integration Tests
```bash
# Start Typesense
docker-compose -f docker-compose.dev.yml up -d typesense

# Run search tests
TYPESENSE_HOST=localhost:8108 go test ./internal/infra/search/... -tags=integration -v

# Test scenarios:
# - Index creation and deletion
# - Document indexing and search
# - Bulk operations
# - Search relevance
# - Faceting and filtering
```

---

### Phase 2: Service Layer Integration Tests

**Goal**: Test business logic with real database

```bash
# All infrastructure running
docker-compose -f docker-compose.dev.yml up -d postgres dragonfly typesense

# Run service tests
go test ./internal/service/... -tags=integration -v
```

**Critical Service Tests**:

#### User Service
- Create user with all fields (including nullable)
- Update user (concurrent updates)
- Delete user (cascade to sessions, tokens)
- Unique email constraint
- Password hashing verification

#### Settings Service
- Server settings CRUD
- User settings with defaults
- User preferences with JSON
- Concurrent setting updates

#### Session Service
- Session creation and validation
- Session expiration
- Concurrent session access
- Session cleanup (expired sessions)

#### RBAC Service
- Role assignment
- Permission checking
- Casbin policy enforcement
- Role hierarchy

---

### Phase 3: Full Stack Integration Tests

**Goal**: Test entire request flow from HTTP to database

```bash
# Start full application
docker-compose -f docker-compose.dev.yml up -d

# Run E2E tests
go test ./tests/e2e/... -v
```

**E2E Scenarios**:
1. User registration → email verification → login → access protected endpoint
2. Create content → search → retrieve → update → delete
3. Admin creates user → assigns role → user accesses based on permissions
4. Session expiration → re-login → new session created

---

## Test Organization

### Directory Structure
```
tests/
├── integration/          # Integration tests (require Docker)
│   ├── database/        # Database integration tests
│   ├── cache/           # Cache integration tests (L2 with real Redis)
│   ├── search/          # Search integration tests
│   └── services/        # Service layer integration tests
├── e2e/                 # End-to-end tests (full application)
├── load/                # Load and performance tests (k6 scripts)
└── chaos/               # Chaos engineering tests
```

### Test Tags
```go
//go:build integration

// For tests requiring Docker services
```

```go
//go:build unit

// For pure unit tests (no external dependencies)
```

### Running Tests

**Unit tests only**:
```bash
go test ./... -tags=unit -short
```

**Integration tests** (requires Docker):
```bash
make test-integration
# OR
docker-compose -f docker-compose.dev.yml up -d postgres dragonfly typesense
go test ./tests/integration/... -v
```

**All tests**:
```bash
make test-all
```

---

## Test Data Management

### Database Test Data

**Use TestDB for isolation**:
```go
func TestUserService_Create(t *testing.T) {
    db := testutil.NewTestDB(t) // Creates isolated test database
    defer db.Close()

    service := user.NewService(db)
    // Test with clean database
}
```

**Fixtures for complex scenarios**:
```sql
-- tests/fixtures/users.sql
INSERT INTO core.users (id, email, username, password_hash) VALUES
    ('user-1', 'test@example.com', 'testuser', '$2a$...'),
    ('user-2', 'admin@example.com', 'admin', '$2a$...');
```

### Cache Test Data

**Clear cache between tests**:
```go
func TestCache(t *testing.T) {
    cache := setupTestCache(t)
    defer cache.Close()

    t.Run("scenario1", func(t *testing.T) {
        cache.Clear() // Isolated test
        // ...
    })
}
```

---

## Performance Testing

### Load Tests (k6)

**Goal**: Establish performance baselines

```javascript
// tests/load/api_load.js
import http from 'k6/http';
import { check } from 'k6';

export let options = {
  stages: [
    { duration: '30s', target: 100 },  // Ramp up
    { duration: '1m', target: 100 },   // Sustained
    { duration: '10s', target: 0 },    // Ramp down
  ],
};

export default function () {
  let res = http.get('http://localhost:8096/api/health');
  check(res, { 'status is 200': (r) => r.status === 200 });
}
```

**Run load tests**:
```bash
docker-compose -f docker-compose.dev.yml up -d
k6 run tests/load/api_load.js
```

**Metrics to Track**:
- p50, p95, p99 latency
- Requests per second
- Error rate
- Database connection pool usage
- Cache hit ratio
- Memory usage

---

## Chaos Testing

### Scenarios

**1. Database Connection Loss**
```bash
# Mid-test, kill database
docker-compose -f docker-compose.dev.yml stop postgres
# Verify: Application degrades gracefully, doesn't crash

# Restart database
docker-compose -f docker-compose.dev.yml start postgres
# Verify: Application reconnects automatically
```

**2. Cache Unavailable**
```bash
# Kill Dragonfly
docker-compose -f docker-compose.dev.yml stop dragonfly
# Verify: Falls back to L1 cache only, no errors

# Restart
docker-compose -f docker-compose.dev.yml start dragonfly
# Verify: L2 cache re-enabled automatically
```

**3. Memory Pressure**
```bash
# Fill cache to capacity
# Verify: Eviction works correctly, no OOM
```

**4. Network Partitions**
```bash
# Use tc (traffic control) to introduce latency/packet loss
sudo tc qdisc add dev lo root netem delay 100ms 20ms loss 5%
# Verify: Timeouts work, retries succeed

# Remove
sudo tc qdisc del dev lo root
```

---

## CI/CD Integration

### GitHub Actions Workflow

```yaml
name: Integration Tests

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest

    services:
      postgres:
        image: postgres:18-alpine
        env:
          POSTGRES_PASSWORD: test
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5

      dragonfly:
        image: docker.dragonflydb.io/dragonflydb/dragonfly:latest
        options: >-
          --health-cmd "redis-cli ping"
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5

    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-go@v5
        with:
          go-version: '1.25'

      - name: Run integration tests
        run: |
          go test ./tests/integration/... -v -race -coverprofile=coverage.out

      - name: Upload coverage
        uses: codecov/codecov-action@v4
```

---

## Coverage Goals

### By Component

| Component | Unit Coverage | Integration Coverage | Total Goal |
|-----------|--------------|---------------------|------------|
| Cache (L1) | 70% | 90% | **90%** |
| Cache (L2) | N/A | 90% | **90%** |
| Database | 50% | 90% | **85%** |
| Services | 60% | 80% | **80%** |
| API Handlers | 40% | 70% | **70%** |
| Health Checks | 30% | 90% | **80%** |
| Job Queue | 50% | 80% | **75%** |
| Search | 40% | 80% | **75%** |

**Priority**: Integration coverage > Unit coverage for infrastructure code

---

## Bug Tracking

### Bug Report Template

When integration tests find bugs:

```markdown
## Bug #XX: [Short Description]

**Severity**: Critical/High/Medium/Low
**Component**: Cache/Database/API/etc.
**Found By**: Integration test / Manual testing / Production

### Symptom
What fails, error message, unexpected behavior

### Test That Reproduces
```go
func TestBug_XX(t *testing.T) {
    // Minimal reproduction
}
```

### Root Cause
Analysis of why it happens

### Fix
Code changes made

### Verification
- [ ] Test passes after fix
- [ ] Related tests still pass
- [ ] No regressions introduced

### Prevention
How to avoid similar bugs in future
```

---

## Current Status

### Infrastructure Ready ✅
- [x] Docker Compose configured
- [x] PostgreSQL 18 running (healthy)
- [x] Dragonfly running (healthy)
- [x] Typesense running
- [x] Docker image built

### Tests Created
- [x] L1 Cache comprehensive tests (20 tests, all passing)
- [ ] L2 Cache integration tests (NEXT)
- [ ] Database integration tests
- [ ] Service layer integration tests
- [ ] E2E tests

### Next Immediate Steps

1. **Create L2 cache integration tests** with real Dragonfly
2. **Create database integration tests** with real PostgreSQL
3. **Run tests and find bugs** (expect to find real issues)
4. **Fix bugs** and document them
5. **Add load tests** to find performance issues
6. **Chaos testing** to verify resilience

---

## Success Criteria

- ✅ All integration tests pass with real infrastructure
- ✅ Performance meets baselines (p95 < 100ms for API calls)
- ✅ Zero connection leaks after 1000+ requests
- ✅ Graceful degradation when services fail
- ✅ No memory leaks over 1 hour sustained load
- ✅ 85%+ code coverage across all components
- ✅ CI/CD pipeline runs all tests on every PR

---

**Remember**: Real infrastructure finds real bugs. Integration tests > Unit tests for infrastructure code.
