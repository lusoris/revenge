---
applyTo: "**/*.go"
---

# Go 1.25 Quick Reference

> New features in Go 1.25 (August 2025) - USE THESE!

## Go 1.25 New Features

### sync.WaitGroup.Go (NEW!)

```go
// OLD (don't use)
var wg sync.WaitGroup
wg.Add(1)
go func() {
    defer wg.Done()
    doWork()
}()

// NEW (use this) - much cleaner
var wg sync.WaitGroup
wg.Go(func() {
    doWork()
})
wg.Wait()
```

### testing/synctest (NEW!)

```go
import "testing/synctest"

func TestConcurrent(t *testing.T) {
    synctest.Test(t, func(t *testing.T) {
        // Time is virtualized
        // Goroutines are tracked
        // Race conditions detected
    })
}
```

### net/http.CrossOriginProtection (NEW!)

```go
mux := http.NewServeMux()

// Built-in CSRF protection
protection := http.CrossOriginProtection{}
mux.Handle("/api/", protection.Handler(apiHandler))
```

### slog.GroupAttrs (NEW!)

```go
// Cleaner grouped logging
slog.Info("request", slog.GroupAttrs("http",
    slog.String("method", r.Method),
    slog.String("path", r.URL.Path),
)...)
```

### runtime/trace.FlightRecorder (NEW!)

```go
import "runtime/trace"

fr := &trace.FlightRecorder{}
fr.Start()

// On error, snapshot last few seconds:
fr.WriteTo(file)
```

### Container-Aware GOMAXPROCS (Automatic!)

```go
// No more automaxprocs needed!
// Go 1.25 respects cgroup CPU limits automatically
// Also updates dynamically when limits change
```

## Go 1.24 Features (Still Valid)

### Generic Type Aliases

```go
type Set[T comparable] = map[T]struct{}
type Result[T any] = func() (T, error)
```

### Tool Directive in go.mod

```go
// go.mod
tool (
    golang.org/x/tools/cmd/stringer
    github.com/sqlc-dev/sqlc/cmd/sqlc
)
```

### testing.B.Loop (Benchmarks)

```go
func BenchmarkNew(b *testing.B) {
    for b.Loop() {  // Not b.N!
        doWork()
    }
}
```

### runtime.AddCleanup

```go
// Now runs concurrently and in parallel (1.25)
runtime.AddCleanup(obj, func(ptr *Object) {
    ptr.Close()
})
```

### encoding/json `omitzero`

```go
type Config struct {
    Timeout time.Duration `json:"timeout,omitzero"`
}
```

### os.Root (Extended in 1.25)

```go
root, _ := os.OpenRoot("/var/data")
// 1.25 adds: Chmod, Chown, Link, MkdirAll, ReadFile, Symlink
data, _ := root.ReadFile("config.yaml")
```

## Experimental Features (Test These)

```bash
# New GC with 10-40% overhead reduction
GOEXPERIMENT=greenteagc go build

# New JSON implementation (faster)
GOEXPERIMENT=jsonv2 go build
```

## New go vet Analyzers

```go
// waitgroup - catches misplaced Add calls
var wg sync.WaitGroup
go func() {
    wg.Add(1)  // vet error: Add after goroutine start
    // ...
}()

// hostport - catches IPv6-unsafe address formatting
addr := host + ":" + port  // vet error: use net.JoinHostPort
addr := net.JoinHostPort(host, port)  // correct
```

## Performance Improvements (Automatic)

- Faster slices (more stack allocations)
- DWARF5 debug info (smaller binaries)
- 3x faster crypto/rsa key generation
- 2x faster crypto/sha1 (amd64 with SHA-NI)
- 2x faster crypto/sha3 (Apple Silicon)
