# Bleeding-Edge Technology Stack

## 🚀 Modern Go Setup (January 2026)

### Core Language
- **Go 1.24** (Latest stable)
  - Enhanced type inference
  - Improved generics performance
  - Better error handling patterns
  - Iterator protocol support

### Standard Library First
- **net/http.ServeMux** (Go 1.22+) - Enhanced routing with method & path patterns
- **log/slog** (Go 1.21+) - Structured logging built-in
- **context** - First-class context support
- **errors** - Modern error handling with `errors.Is` and `errors.As`

### Dependencies (Carefully Selected)

#### Configuration
- **koanf v2** - Modern, type-safe configuration
  - Replaces deprecated Viper
  - Better performance
  - More flexible providers
  - Proper v2 with breaking changes fixed

#### Logging
- **tint** - Beautiful console logging for slog
  - Colored output
  - Human-readable format
  - Minimal overhead

#### Dependency Injection
- **uber-go/fx v1.23** - Latest stable
  - Improved performance
  - Better error messages
  - Lifecycle management

#### Database (Future)
- **modernc.org/sqlite** - Pure Go SQLite
  - No CGO required
  - Better portability
  - Modern API
- **pgx v5** - PostgreSQL driver
  - Better performance than lib/pq
  - Native prepared statements
  - Connection pooling

### Removed/Replaced

❌ **gorilla/mux** → ✅ **net/http (stdlib)**
- Reason: Go 1.22+ has enhanced routing built-in
- Breaking change: Modern pattern matching

❌ **Viper** → ✅ **koanf v2**
- Reason: Viper is in maintenance mode
- Breaking change: Different API, better design

❌ **zap** → ✅ **slog (stdlib)**
- Reason: slog is now standard library
- Breaking change: Different API, but stdlib

## 📦 Dependency Philosophy

1. **Stdlib First** - Use Go standard library when possible
2. **Minimal Dependencies** - Only add when truly needed
3. **Active Maintenance** - No deprecated or abandoned packages
4. **Performance** - Choose performant, well-tested libraries
5. **Type Safety** - Prefer strongly-typed APIs

## 🎯 Modern Go Features Used

### Go 1.22+ Features
- **Enhanced ServeMux Patterns**: `mux.HandleFunc("GET /api/users/{id}", ...)`
- **Ranging over functions**: Iterator support
- **Profile-guided optimization**: Better compiler performance

### Go 1.21+ Features
- **log/slog**: Structured logging
- **min/max built-ins**: No more custom functions
- **clear built-in**: Clear maps and slices

### Go 1.20+ Features
- **Multierror wrapping**: `errors.Join()`
- **HttpOnly cookies**: Better security defaults

### Modern Patterns
- **Context-first APIs**: All functions accept context
- **Generics**: Type-safe collections where appropriate
- **Functional options**: Clean configuration
- **Error wrapping**: Proper error chains with `%w`

## 🔄 Migration from Old Stack

| Old | New | Reason |
|-----|-----|--------|
| gorilla/mux | net/http.ServeMux | Stdlib is enough now |
| Viper | koanf v2 | Modern, maintained |
| zap | slog + tint | Stdlib with pretty output |
| database/sql | pgx v5 | Better PostgreSQL support |
| lib/pq | pgx v5 | More features, better perf |

## ⚡ Performance Improvements

- **Faster routing**: stdlib ServeMux is optimized
- **Less allocations**: slog is allocation-efficient
- **Better GC**: Modern Go runtime improvements
- **Profile-guided optimization**: Compiler optimizations

## 🛡️ Security Updates

- **Latest Go runtime**: Security fixes included
- **No deprecated packages**: Reduced vulnerability surface
- **Modern crypto**: Using latest stdlib crypto
- **Secure defaults**: Better default configurations

## 📝 Code Style

### Modern Go Code
```go
// ✅ Good: Modern Go 1.22+ routing
mux.HandleFunc("GET /users/{id}", handleGetUser)

// ✅ Good: Structured logging with slog
logger.Info("user created", slog.String("id", id), slog.Int("age", age))

// ✅ Good: Error wrapping
return fmt.Errorf("failed to create user: %w", err)

// ✅ Good: Context-first
func GetUser(ctx context.Context, id string) (*User, error)
```

### Deprecated Patterns
```go
// ❌ Bad: Old gorilla/mux
r.HandleFunc("/users/{id}", handleGetUser).Methods("GET")

// ❌ Bad: Printf-style logging
log.Printf("user created: id=%s, age=%d", id, age)

// ❌ Bad: Error messages without wrapping
return errors.New("failed to create user")

// ❌ Bad: No context
func GetUser(id string) (*User, error)
```

## 🔧 Development Tools (Latest)

- **golangci-lint v1.55+**: Latest linters
- **gopls**: Latest Go language server
- **govulncheck**: Security scanning
- **air v1.52+**: Hot reload
- **migrate v4.17+**: Database migrations

## 📚 Resources

- [Go 1.24 Release Notes](https://go.dev/doc/go1.24)
- [Go 1.22 Enhanced Routing](https://go.dev/blog/routing-enhancements)
- [log/slog Package](https://pkg.go.dev/log/slog)
- [koanf Documentation](https://github.com/knadh/koanf)

## ✅ Checklist

- [x] Go 1.24 (latest stable)
- [x] stdlib routing (no gorilla/mux)
- [x] slog for logging (no zap)
- [x] koanf for config (no viper)
- [x] Modern error handling
- [x] Context-first APIs
- [x] Type-safe patterns
- [x] No deprecated dependencies
- [x] Security-focused defaults
- [x] Performance optimized

This is now a **truly modern, bleeding-edge Go project** following 2026 best practices! 🚀
