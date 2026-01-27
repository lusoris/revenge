---
applyTo: "**/*_test.go,**/tests/**/*.go"
---

# Go Testing Patterns

> Idiomatic Go testing for revenge

## Table-Driven Tests

```go
func TestParseMediaPath(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    *MediaPath
        wantErr bool
    }{
        {
            name:  "valid movie path",
            input: "/movies/The Matrix (1999)/The Matrix.mkv",
            want: &MediaPath{
                Type:  MediaTypeMovie,
                Title: "The Matrix",
                Year:  1999,
            },
        },
        {
            name:  "valid series path",
            input: "/tv/Breaking Bad/Season 01/S01E01.mkv",
            want: &MediaPath{
                Type:    MediaTypeSeries,
                Title:   "Breaking Bad",
                Season:  1,
                Episode: 1,
            },
        },
        {
            name:    "invalid path",
            input:   "",
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := ParseMediaPath(tt.input)

            if (err != nil) != tt.wantErr {
                t.Errorf("ParseMediaPath() error = %v, wantErr %v", err, tt.wantErr)
                return
            }

            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("ParseMediaPath() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

## Benchmarks (Go 1.24)

```go
// Use b.Loop() instead of b.N (Go 1.24+)
func BenchmarkHashPassword(b *testing.B) {
    password := "testpassword123"

    for b.Loop() {
        _, _ = HashPassword(password)
    }
}

// With setup
func BenchmarkGetUser(b *testing.B) {
    svc := setupTestService(b)
    ctx := context.Background()

    b.ResetTimer()
    for b.Loop() {
        _, _ = svc.GetUser(ctx, 1)
    }
}

// Parallel benchmark
func BenchmarkGetUserParallel(b *testing.B) {
    svc := setupTestService(b)

    b.RunParallel(func(pb *testing.PB) {
        ctx := context.Background()
        for pb.Next() {
            _, _ = svc.GetUser(ctx, 1)
        }
    })
}
```

## Test Fixtures

```go
// testdata directory (automatically included)
func TestLoadConfig(t *testing.T) {
    // testdata/ is automatically available
    data, err := os.ReadFile("testdata/config.yaml")
    if err != nil {
        t.Fatal(err)
    }
    // ...
}

// Embed test fixtures
//go:embed testdata/*.json
var testFixtures embed.FS

func loadFixture(t *testing.T, name string) []byte {
    t.Helper()
    data, err := testFixtures.ReadFile("testdata/" + name)
    if err != nil {
        t.Fatalf("load fixture %s: %v", name, err)
    }
    return data
}
```

## Test Helpers

```go
// Use t.Helper() for helper functions
func assertEqual[T comparable](t *testing.T, got, want T) {
    t.Helper()
    if got != want {
        t.Errorf("got %v, want %v", got, want)
    }
}

func assertNoError(t *testing.T, err error) {
    t.Helper()
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
}

func assertError(t *testing.T, err, target error) {
    t.Helper()
    if !errors.Is(err, target) {
        t.Errorf("got error %v, want %v", err, target)
    }
}

// Setup helper
func setupTestDB(t *testing.T) *pgxpool.Pool {
    t.Helper()

    pool, err := pgxpool.New(context.Background(), os.Getenv("TEST_DATABASE_URL"))
    if err != nil {
        t.Fatalf("connect to test db: %v", err)
    }

    t.Cleanup(func() {
        pool.Close()
    })

    return pool
}
```

## Subtests and Parallel

```go
func TestUserService(t *testing.T) {
    // Run subtests in parallel
    t.Run("GetUser", func(t *testing.T) {
        t.Parallel()
        // ...
    })

    t.Run("CreateUser", func(t *testing.T) {
        t.Parallel()
        // ...
    })

    // Nested subtests
    t.Run("Validation", func(t *testing.T) {
        t.Run("EmptyName", func(t *testing.T) {
            t.Parallel()
            // ...
        })
        t.Run("InvalidEmail", func(t *testing.T) {
            t.Parallel()
            // ...
        })
    })
}
```

## Mocks with Interfaces

```go
// Define interface for dependencies
type UserRepository interface {
    GetByID(ctx context.Context, id int64) (*User, error)
    Create(ctx context.Context, params CreateUserParams) (*User, error)
}

// Mock implementation
type mockUserRepo struct {
    getByIDFn func(ctx context.Context, id int64) (*User, error)
    createFn  func(ctx context.Context, params CreateUserParams) (*User, error)
}

func (m *mockUserRepo) GetByID(ctx context.Context, id int64) (*User, error) {
    return m.getByIDFn(ctx, id)
}

func (m *mockUserRepo) Create(ctx context.Context, params CreateUserParams) (*User, error) {
    return m.createFn(ctx, params)
}

// Use in tests
func TestUserService_GetUser(t *testing.T) {
    wantUser := &User{ID: 1, Name: "Alice"}

    repo := &mockUserRepo{
        getByIDFn: func(ctx context.Context, id int64) (*User, error) {
            if id == 1 {
                return wantUser, nil
            }
            return nil, ErrNotFound
        },
    }

    svc := NewUserService(repo)

    got, err := svc.GetUser(context.Background(), 1)
    assertNoError(t, err)
    assertEqual(t, got.Name, wantUser.Name)
}
```

## HTTP Handler Tests

```go
func TestGetUserHandler(t *testing.T) {
    // Setup
    svc := &mockUserService{...}
    handler := NewUserHandler(svc)

    tests := []struct {
        name       string
        userID     string
        wantStatus int
        wantBody   string
    }{
        {
            name:       "valid user",
            userID:     "1",
            wantStatus: http.StatusOK,
        },
        {
            name:       "invalid id",
            userID:     "abc",
            wantStatus: http.StatusBadRequest,
        },
        {
            name:       "not found",
            userID:     "999",
            wantStatus: http.StatusNotFound,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            req := httptest.NewRequest("GET", "/users/"+tt.userID, nil)
            req.SetPathValue("id", tt.userID)

            rec := httptest.NewRecorder()
            handler.GetUser(rec, req)

            if rec.Code != tt.wantStatus {
                t.Errorf("status = %d, want %d", rec.Code, tt.wantStatus)
            }
        })
    }
}
```

## Integration Tests

```go
//go:build integration

package integration_test

import (
    "testing"
)

func TestMain(m *testing.M) {
    // Setup: start containers, migrate DB
    pool := setupTestDatabase()
    defer pool.Close()

    os.Exit(m.Run())
}

func TestCreateUser_Integration(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping integration test")
    }

    ctx := context.Background()

    // Clean state
    _, err := testPool.Exec(ctx, "TRUNCATE users CASCADE")
    assertNoError(t, err)

    // Test
    repo := NewUserRepository(testPool)
    user, err := repo.Create(ctx, CreateUserParams{
        Name:  "Test User",
        Email: "test@example.com",
    })

    assertNoError(t, err)
    assertEqual(t, user.Name, "Test User")

    // Verify
    got, err := repo.GetByID(ctx, user.ID)
    assertNoError(t, err)
    assertEqual(t, got.Email, "test@example.com")
}
```

## fx Test Helpers

```go
import "go.uber.org/fx/fxtest"

func TestApp(t *testing.T) {
    var svc *UserService

    app := fxtest.New(t,
        fx.Provide(
            NewTestConfig,
            NewTestDatabase,
            NewUserRepository,
            NewUserService,
        ),
        fx.Populate(&svc),
    )

    app.RequireStart()
    defer app.RequireStop()

    // Test with real dependencies
    user, err := svc.GetUser(context.Background(), 1)
    // ...
}
```

## Golden Files

```go
func TestRenderTemplate(t *testing.T) {
    got := RenderTemplate(data)

    golden := filepath.Join("testdata", t.Name()+".golden")

    if *update {
        os.WriteFile(golden, []byte(got), 0644)
    }

    want, err := os.ReadFile(golden)
    assertNoError(t, err)

    if got != string(want) {
        t.Errorf("output mismatch:\ngot:\n%s\nwant:\n%s", got, want)
    }
}

// Run with: go test -update
var update = flag.Bool("update", false, "update golden files")
```

## Coverage

```bash
# Run with coverage
go test -cover ./...

# Generate HTML report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html

# Show coverage per function
go tool cover -func=coverage.out
```

## Race Detection

```bash
# Enable race detector
go test -race ./...

# In CI (required)
go test -race -v ./...
```

## Test Organization

```
internal/
├── service/
│   ├── user.go
│   └── user_test.go      # Unit tests
└── infra/
    └── database/
        └── repository/
            ├── user.go
            └── user_test.go
tests/
└── integration/
    ├── user_test.go      # Integration tests
    └── testdata/
        └── fixtures.sql
```

## Test Commands

```bash
# All tests
go test ./...

# Specific package
go test ./internal/service/...

# Verbose
go test -v ./...

# Short (skip slow tests)
go test -short ./...

# Run specific test
go test -run TestUserService_GetUser ./...

# Run matching tests
go test -run "TestUser.*" ./...

# Integration tests only
go test -tags integration ./tests/integration/...
```
