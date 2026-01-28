---
applyTo: "**/internal/infra/database/**/*.go,**/internal/infra/database/**/*.sql,sqlc.yaml"
---

# sqlc Database Guide

> Type-safe SQL for Go with PostgreSQL/pgx v5

## Installation

```bash
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
go get github.com/jackc/pgx/v5
```

## Configuration (sqlc.yaml)

```yaml
version: "2"
sql:
  - engine: "postgresql"
    queries: "internal/infra/database/queries/"
    schema: "migrations/"
    gen:
      go:
        package: "db"
        out: "internal/infra/database/db"
        sql_package: "pgx/v5"
        emit_json_tags: true
        emit_interface: true
        emit_empty_slices: true
        emit_result_struct_pointers: false
        emit_params_struct_pointers: false
        emit_methods_with_db_argument: false
        json_tags_case_style: "camel"
        output_db_file_name: "db.go"
        output_models_file_name: "models.go"
        output_querier_file_name: "querier.go"
```

## Query Annotations

```sql
-- name: GetUser :one
-- Returns a single row (error if not found)
SELECT * FROM users WHERE id = $1 LIMIT 1;

-- name: ListUsers :many
-- Returns multiple rows (empty slice if none)
SELECT * FROM users ORDER BY created_at DESC;

-- name: CreateUser :one
-- INSERT with RETURNING
INSERT INTO users (name, email, password_hash)
VALUES ($1, $2, $3)
RETURNING *;

-- name: UpdateUser :exec
-- UPDATE without return value
UPDATE users SET name = $2, email = $3 WHERE id = $1;

-- name: UpdateUserReturning :one
-- UPDATE with RETURNING
UPDATE users SET name = $2 WHERE id = $1 RETURNING *;

-- name: DeleteUser :exec
-- DELETE without return value
DELETE FROM users WHERE id = $1;

-- name: DeleteUserReturning :one
-- DELETE with RETURNING
DELETE FROM users WHERE id = $1 RETURNING *;

-- name: CountUsers :one
-- Aggregate query
SELECT count(*) FROM users;

-- name: UserExists :one
-- Boolean check
SELECT EXISTS(SELECT 1 FROM users WHERE email = $1);

-- name: BulkInsertUsers :copyfrom
-- Bulk insert (pgx COPY)
INSERT INTO users (name, email) VALUES ($1, $2);

-- name: BatchGetUsers :batchone
-- Batch query (pgx batch)
SELECT * FROM users WHERE id = $1;
```

## Generated Code Usage

```go
package repository

import (
    "context"
    "github.com/jackc/pgx/v5/pgxpool"
    "revenge/internal/infra/database/db"
)

type UserRepository struct {
    pool    *pgxpool.Pool
    queries *db.Queries
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
    return &UserRepository{
        pool:    pool,
        queries: db.New(pool),
    }
}

// Single row query
func (r *UserRepository) GetByID(ctx context.Context, id int64) (*db.User, error) {
    user, err := r.queries.GetUser(ctx, id)
    if err != nil {
        return nil, fmt.Errorf("get user %d: %w", id, err)
    }
    return &user, nil
}

// Multiple rows
func (r *UserRepository) List(ctx context.Context) ([]db.User, error) {
    return r.queries.ListUsers(ctx)
}

// Insert with returning
func (r *UserRepository) Create(ctx context.Context, params db.CreateUserParams) (*db.User, error) {
    user, err := r.queries.CreateUser(ctx, params)
    if err != nil {
        return nil, fmt.Errorf("create user: %w", err)
    }
    return &user, nil
}
```

## Transactions

```go
func (r *UserRepository) CreateWithProfile(ctx context.Context, user db.CreateUserParams, profile db.CreateProfileParams) error {
    tx, err := r.pool.Begin(ctx)
    if err != nil {
        return fmt.Errorf("begin tx: %w", err)
    }
    defer tx.Rollback(ctx)

    // Use transaction-scoped queries
    qtx := r.queries.WithTx(tx)

    newUser, err := qtx.CreateUser(ctx, user)
    if err != nil {
        return fmt.Errorf("create user: %w", err)
    }

    profile.UserID = newUser.ID
    if _, err := qtx.CreateProfile(ctx, profile); err != nil {
        return fmt.Errorf("create profile: %w", err)
    }

    return tx.Commit(ctx)
}
```

## Complex Queries

### JOINs

```sql
-- name: GetUserWithProfile :one
SELECT
    u.id,
    u.name,
    u.email,
    p.bio,
    p.avatar_url
FROM users u
LEFT JOIN profiles p ON p.user_id = u.id
WHERE u.id = $1;
```

### Custom Struct Results

```sql
-- name: GetUserStats :one
SELECT
    u.id,
    u.name,
    count(m.id) as media_count,
    sum(m.duration) as total_duration
FROM users u
LEFT JOIN media_items m ON m.user_id = u.id
WHERE u.id = $1
GROUP BY u.id, u.name;
```

### Dynamic Filtering

```sql
-- name: ListUsersFiltered :many
SELECT * FROM users
WHERE
    (sqlc.narg('name')::text IS NULL OR name ILIKE '%' || sqlc.narg('name') || '%')
    AND (sqlc.narg('email')::text IS NULL OR email = sqlc.narg('email'))
ORDER BY created_at DESC
LIMIT sqlc.arg('limit')
OFFSET sqlc.arg('offset');
```

### IN Clause

```sql
-- name: GetUsersByIDs :many
SELECT * FROM users WHERE id = ANY($1::bigint[]);
```

```go
users, err := queries.GetUsersByIDs(ctx, []int64{1, 2, 3})
```

### JSON Columns

```sql
-- name: GetUserPreferences :one
SELECT preferences FROM users WHERE id = $1;

-- name: UpdateUserPreferences :exec
UPDATE users SET preferences = $2 WHERE id = $1;
```

```go
// In schema
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    preferences JSONB NOT NULL DEFAULT '{}'
);
```

## UUID Support

```sql
-- Schema
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL
);

-- Query
-- name: GetItem :one
SELECT * FROM items WHERE id = $1;
```

```go
import "github.com/google/uuid"

item, err := queries.GetItem(ctx, uuid.MustParse("..."))
```

## Nullable Fields

```sql
-- name: GetUser :one
SELECT
    id,
    name,
    email,
    bio  -- nullable column
FROM users WHERE id = $1;
```

sqlc generates `pgtype` types for nullable columns:

```go
type User struct {
    ID    int64
    Name  string
    Email string
    Bio   pgtype.Text  // nullable
}

// Check if null
if user.Bio.Valid {
    fmt.Println(user.Bio.String)
}
```

## Enums

```sql
-- Schema
CREATE TYPE user_role AS ENUM ('admin', 'user', 'guest');

CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    role user_role NOT NULL DEFAULT 'user'
);

-- Query
-- name: GetUser :one
SELECT * FROM users WHERE id = $1;
```

sqlc generates Go constants:

```go
type UserRole string

const (
    UserRoleAdmin UserRole = "admin"
    UserRoleUser  UserRole = "user"
    UserRoleGuest UserRole = "guest"
)
```

## Batch Operations (pgx)

```sql
-- name: BatchCreateUsers :batchexec
INSERT INTO users (name, email) VALUES ($1, $2);
```

```go
batch := &pgx.Batch{}
for _, u := range users {
    queries.BatchCreateUsers(batch, db.BatchCreateUsersParams{
        Name:  u.Name,
        Email: u.Email,
    })
}

results := pool.SendBatch(ctx, batch)
defer results.Close()

for range users {
    if _, err := results.Exec(); err != nil {
        return err
    }
}
```

## Database Connection (pgx)

```go
import (
    "context"
    "github.com/jackc/pgx/v5/pgxpool"
)

func NewPool(cfg *config.Database) (*pgxpool.Pool, error) {
    connStr := fmt.Sprintf(
        "postgres://%s:%s@%s:%d/%s?sslmode=%s",
        cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name, cfg.SSLMode,
    )

    poolConfig, err := pgxpool.ParseConfig(connStr)
    if err != nil {
        return nil, err
    }

    poolConfig.MaxConns = int32(cfg.MaxConns)
    poolConfig.MinConns = int32(cfg.MinConns)
    poolConfig.MaxConnLifetime = cfg.MaxConnLifetime
    poolConfig.MaxConnIdleTime = cfg.MaxConnIdleTime

    return pgxpool.NewWithConfig(context.Background(), poolConfig)
}
```

## fx Integration

```go
func NewDatabase(cfg *config.Config) (*pgxpool.Pool, error) {
    return NewPool(&cfg.Database)
}

func NewQueries(pool *pgxpool.Pool) *db.Queries {
    return db.New(pool)
}

// Register
fx.Provide(
    NewDatabase,
    NewQueries,
    NewUserRepository,
)
```

## Error Handling

```go
import (
    "errors"
    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgconn"
)

var ErrNotFound = errors.New("not found")
var ErrDuplicate = errors.New("duplicate entry")

func (r *UserRepository) GetByID(ctx context.Context, id int64) (*db.User, error) {
    user, err := r.queries.GetUser(ctx, id)
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return nil, ErrNotFound
        }
        return nil, fmt.Errorf("get user: %w", err)
    }
    return &user, nil
}

func (r *UserRepository) Create(ctx context.Context, params db.CreateUserParams) (*db.User, error) {
    user, err := r.queries.CreateUser(ctx, params)
    if err != nil {
        var pgErr *pgconn.PgError
        if errors.As(err, &pgErr) {
            if pgErr.Code == "23505" { // unique_violation
                return nil, ErrDuplicate
            }
        }
        return nil, fmt.Errorf("create user: %w", err)
    }
    return &user, nil
}
```

## File Organization

```
internal/infra/database/
├── queries/
│   ├── users.sql
│   ├── media.sql
│   ├── libraries.sql
│   └── sessions.sql
├── db/
│   ├── db.go          # Generated
│   ├── models.go      # Generated
│   ├── querier.go     # Generated
│   ├── users.sql.go   # Generated
│   └── media.sql.go   # Generated
└── repository/
    ├── user.go
    └── media.go
```

## Commands

```bash
# Generate code
sqlc generate

# Verify queries (no generation)
sqlc compile

# Show differences
sqlc diff

# Verify against database
sqlc vet
```
---

## Related

- [INDEX.instructions.md](INDEX.instructions.md) - Main instruction index with all cross-references
- [ARCHITECTURE_V2.md](../../docs/dev/design/architecture/ARCHITECTURE_V2.md) - System architecture
- [BEST_PRACTICES.md](../../docs/dev/design/operations/BEST_PRACTICES.md) - Best practices
