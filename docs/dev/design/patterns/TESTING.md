# Testing Patterns

> How to write tests in this codebase: database setup, mocking, table-driven tests, integration tests, and CI. Written from code as of 2026-02-06.

---

## Test Infrastructure

### Fast Database Tests (`testutil.NewFastTestDB`)

The core innovation: **template database pattern** via testcontainers.

```go
func setupTestRepo(t *testing.T) (*RepositoryPG, testutil.DB) {
    t.Helper()
    testDB := testutil.NewFastTestDB(t)  // ~10ms per test
    queries := db.New(testDB.Pool())
    return NewRepositoryPG(queries), testDB
}
```

How it works:
1. First call starts a shared PostgreSQL container (testcontainers)
2. Runs all migrations once on a template database
3. Each test gets a **cloned copy** of the template (`CREATE DATABASE ... TEMPLATE ...`)
4. Tests run in full isolation — no cleanup needed
5. Container is shared across the entire test run via `TestMain`

```go
func TestMain(m *testing.M) {
    code := m.Run()
    testutil.StopSharedPostgres()
    os.Exit(code)
}
```

**Performance**: ~10ms per test database vs 3-5s for individual container startup.

### Dragonfly Container

For cache integration tests:

```go
df := testutil.NewDragonflyContainer(t)
defer df.Close()

cfg := &config.Config{
    Cache: config.CacheConfig{Enabled: true, URL: df.URL},
}
```

### Typesense Container

For search integration tests:

```go
ts := testutil.NewTypesenseContainer(t)
```

---

## Three Test Layers

### 1. Unit Tests (Service Layer — Mocked)

Mock repository interfaces, test business logic in isolation.

```go
func TestService_Register(t *testing.T) {
    t.Parallel()

    t.Run("success", func(t *testing.T) {
        repo := NewMockAuthRepository(t)
        tokenMgr := NewMockTokenManager(t)
        svc := auth.NewService(repo, tokenMgr, ...)

        repo.EXPECT().
            CreateUser(ctx, mock.AnythingOfType("CreateUserParams")).
            Return(expectedUser, nil).
            Once()

        user, err := svc.Register(ctx, req)
        require.NoError(t, err)
        assert.Equal(t, req.Username, user.Username)
    })

    t.Run("duplicate username", func(t *testing.T) {
        repo := NewMockAuthRepository(t)
        repo.EXPECT().
            CreateUser(ctx, mock.Anything).
            Return(nil, errors.New("unique constraint")).
            Once()

        _, err := svc.Register(ctx, req)
        require.Error(t, err)
    })
}
```

### 2. Repository Tests (Database Layer — Real DB)

Test SQL queries against a real PostgreSQL instance.

```go
func TestRepositoryPG_CreateUser(t *testing.T) {
    t.Parallel()
    repo, _ := setupTestRepo(t)
    ctx := context.Background()

    t.Run("valid user", func(t *testing.T) {
        user, err := repo.CreateUser(ctx, CreateUserParams{
            Username:     "testuser",
            Email:        "test@example.com",
            PasswordHash: "hashedpassword123",
        })
        require.NoError(t, err)
        assert.NotEqual(t, uuid.Nil, user.ID)
    })

    t.Run("duplicate username fails", func(t *testing.T) {
        // First create succeeds
        _, err := repo.CreateUser(ctx, CreateUserParams{
            Username: "duplicate", Email: "a@b.com", PasswordHash: "h",
        })
        require.NoError(t, err)

        // Second create fails on constraint
        _, err = repo.CreateUser(ctx, CreateUserParams{
            Username: "duplicate", Email: "c@d.com", PasswordHash: "h",
        })
        require.Error(t, err)
    })
}
```

### 3. Integration Tests (Full Stack)

Build tag: `//go:build integration`. Located in `tests/integration/`.

```go
//go:build integration

func TestAuthService_Login(t *testing.T) {
    authSvc, _, _, cleanup := setupAuthService(t)  // Real repos, real DB
    defer cleanup()

    // Register user first
    user, err := authSvc.Register(ctx, registerReq)
    require.NoError(t, err)

    // Login with real password hashing
    resp, err := authSvc.Login(ctx, user.Username, password, nil, nil, nil, nil)
    require.NoError(t, err)
    assert.NotEmpty(t, resp.AccessToken)
}
```

---

## Mock Generation (Mockery)

Configuration in `.mockery.yaml`:

```yaml
with-expecter: true
mockname: "Mock{{.InterfaceName}}"
filename: "mock_{{.InterfaceName | snakecase}}_test.go"
outpkg: "{{.PackageName}}_test"
dir: "{{.InterfaceDir}}"
```

Generate mocks: `go generate ./...` or `mockery --all --recursive`

### Two Mock API Styles

**EXPECT() style** (newer, preferred):
```go
repo := NewMockAuthRepository(t)
repo.EXPECT().
    GetUserByID(ctx, userID).
    Return(expectedUser, nil).
    Once()
```

**On() style** (older, still in movie module):
```go
repo := new(MockMovieRepository)
repo.On("GetMovie", ctx, movie.ID).Return(movie, nil)
// ... test ...
repo.AssertExpectations(t)
```

Both work. New code should use `EXPECT()`.

---

## Table-Driven Tests

Standard pattern for testing multiple inputs:

```go
func TestParseMovieFilename(t *testing.T) {
    tests := []struct {
        name          string
        filename      string
        expectedTitle string
        expectedYear  *int
    }{
        {
            name:          "Title (YEAR).ext",
            filename:      "The Matrix (1999).mkv",
            expectedTitle: "The Matrix",
            expectedYear:  intPtr(1999),
        },
        {
            name:          "Title with quality markers",
            filename:      "The.Matrix.1999.1080p.BluRay.x264-GROUP.mkv",
            expectedTitle: "The Matrix",
            expectedYear:  intPtr(1999),
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            title, year := parseMovieFilename(tt.filename)
            assert.Equal(t, tt.expectedTitle, title)
        })
    }
}
```

---

## Test Helpers

### Custom Assertions (`testutil/assertions.go`)

```go
testutil.AssertTimeEqual(t, expected, actual)       // 1 second tolerance
testutil.AssertRecentTime(t, tm)                    // Within 1 minute
testutil.AssertValidUUID(t, uuidString)
testutil.AssertNotZeroUUID(t, id)
testutil.AssertSliceContains(t, slice, element)
```

### Fixtures (`testutil/fixtures.go`)

```go
user := testutil.CreateUser(t, pool, testutil.User{
    Username: "testuser",
    Email:    "test@example.com",
    IsAdmin:  true,
})

testutil.DefaultUser()   // Regular user preset
testutil.AdminUser()     // Admin preset
testutil.QARUser()       // QAR-enabled preset
```

### Common Helpers

```go
func ptr[T any](v T) *T { return &v }
func intPtr(v int) *int { return &v }
```

---

## Test Naming

Format: `Test{Receiver}_{Method}_{Condition}`

```go
func TestRepositoryPG_CreateUser(t *testing.T)           // Repository
func TestService_Register_ErrorCreatingUser(t *testing.T) // Service
func TestUnit_GetUser(t *testing.T)                       // Unit
func TestIntegration_L2Cache(t *testing.T)                // Integration
func TestParseMovieFilename(t *testing.T)                 // Algorithm
```

Subtests describe what's being tested:
```go
t.Run("success", ...)
t.Run("not found", ...)
t.Run("duplicate username fails", ...)
```

---

## Conventions

- **Always** use `t.Parallel()` when safe (551+ parallel tests in codebase)
- **Always** use `t.Helper()` in setup functions
- **Always** use `require.NoError(t, err)` for must-pass checks (fails test immediately)
- **Use** `assert.Equal(t, ...)` for non-fatal checks (continues test)
- **Mocks** for service-layer tests, **real DB** for repository tests
- **Never** mix mocks and real database in the same test
- **No manual cleanup** — template database pattern handles isolation
- **No embedded Postgres** — testcontainers PostgreSQL only

---

## Running Tests

```bash
# Unit tests (same as CI)
make test                  # -race -coverprofile=coverage.out

# Fast unit tests
make test-short            # Skips slow tests

# Integration tests (requires Docker)
make test-integration      # -tags=integration

# All tests
make test-all              # Unit + integration

# Specific test
go test -run TestRepositoryPG_CreateUser ./internal/service/auth

# Specific subtest
go test -run TestRepositoryPG_CreateUser/duplicate -v ./internal/service/auth

# Coverage report
make test-coverage         # Generates HTML report
```

CI runs: `go test -race -coverprofile=coverage.out -covermode=atomic -count=1 ./...`

---

## Quick Reference: New Test

### New repository method
```go
func TestRepositoryPG_MyMethod(t *testing.T) {
    t.Parallel()
    repo, _ := setupTestRepo(t)
    ctx := context.Background()

    t.Run("success", func(t *testing.T) { /* ... */ })
    t.Run("error", func(t *testing.T) { /* ... */ })
}
```

### New service method
```go
func TestService_MyMethod(t *testing.T) {
    t.Parallel()
    mockRepo := NewMockRepository(t)
    svc := NewService(mockRepo, ...)

    mockRepo.EXPECT().SomeCall(ctx, mock.Anything).Return(val, nil)

    result, err := svc.MyMethod(ctx, param)
    require.NoError(t, err)
    assert.Equal(t, expected, result)
}
```

### New integration test
```go
//go:build integration

func TestIntegration_MyFlow(t *testing.T) {
    if testing.Short() { t.Skip("skipping integration test") }
    svc, cleanup := setupService(t)
    defer cleanup()

    result, err := svc.FullOperation(ctx, input)
    require.NoError(t, err)
    assert.NotNil(t, result)
}
```

---

## Key Files

| File | Purpose |
|------|---------|
| `internal/testutil/pgtestdb.go` | `NewFastTestDB` (testcontainers + template) |
| `internal/testutil/containers.go` | PostgreSQL, Dragonfly, Typesense containers |
| `internal/testutil/fixtures.go` | Test data factories |
| `internal/testutil/assertions.go` | Custom testify assertions |
| `.mockery.yaml` | Mock generation config |

---

## Related Documentation

- [New Service Checklist](NEW_SERVICE.md) — Includes test setup steps
- [Error Handling](ERROR_HANDLING.md) — How errors flow to test assertions
- [Database Transactions](DATABASE_TRANSACTIONS.md) — Testing transactional code
