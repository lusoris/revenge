# Go 1.24 Quick Reference

> New features in Go 1.24 (February 2025) - USE THESE!

## Generic Type Aliases

```go
// Type aliases can now be generic
type Set[T comparable] = map[T]struct{}
type Result[T any] = func() (T, error)
type Pair[A, B any] = struct{ First A; Second B }

// Usage
var users Set[string] = make(Set[string])
users["alice"] = struct{}{}
```

## Tool Directive in go.mod

```go
// go.mod - No more tools.go workaround!
module example.com/myproject

go 1.24

// Tool dependencies are now first-class
tool (
    golang.org/x/tools/cmd/stringer
    github.com/sqlc-dev/sqlc/cmd/sqlc
    github.com/golangci/golangci-lint/cmd/golangci-lint
)

require (
    // regular dependencies
)
```

```bash
# Run tools
go tool stringer -type=State
go tool sqlc generate
```

## testing.B.Loop (Benchmarks)

```go
// OLD (don't use)
func BenchmarkOld(b *testing.B) {
    for i := 0; i < b.N; i++ {
        doWork()
    }
}

// NEW (use this)
func BenchmarkNew(b *testing.B) {
    for b.Loop() {
        doWork()
    }
}
```

## runtime.AddCleanup

```go
// OLD (don't use)
runtime.SetFinalizer(obj, func(o *Object) {
    o.Close()
})

// NEW (use this) - more predictable
runtime.AddCleanup(obj, func(ptr *Object) {
    ptr.Close()
})
```

## encoding/json `omitzero`

```go
type Config struct {
    // omitempty: omits zero value AND empty strings/slices
    Name string `json:"name,omitempty"`

    // omitzero: omits ONLY zero value (new!)
    Timeout time.Duration `json:"timeout,omitzero"`
    Count   int           `json:"count,omitzero"`

    // Useful for: time.Time, time.Duration, custom types
}

// With omitzero:
// - zero Duration → omitted
// - empty string → included (unlike omitempty)
```

## os.Root (Sandboxed Filesystem)

```go
// Create a root that can't escape directory
root, err := os.OpenRoot("/var/data")
if err != nil {
    return err
}
defer root.Close()

// These can't access files outside /var/data
f, _ := root.Open("config.yaml")         // /var/data/config.yaml
f, _ := root.Open("../etc/passwd")       // ERROR: escapes root

// Perfect for:
// - Media file access (libraries)
// - Plugin sandboxing
// - User upload directories
```

## weak Package

```go
import "weak"

// Create weak pointer (doesn't prevent GC)
ptr := weak.Make(&myObject)

// Later: check if still alive
if obj := ptr.Value(); obj != nil {
    // object is still alive
    use(obj)
}
// obj is nil if garbage collected
```

## crypto/mlkem (Post-Quantum)

```go
import "crypto/mlkem"

// ML-KEM-768 (recommended)
dk, ek := mlkem.GenerateKey768()

// Encapsulate
ciphertext, sharedKey := ek.Encapsulate()

// Decapsulate
sharedKey2 := dk.Decapsulate(ciphertext)
```

## Iterator Functions (stdlib)

```go
import (
    "bytes"
    "strings"
)

// strings.Lines - iterate over lines
for line := range strings.Lines(text) {
    fmt.Println(line)
}

// bytes.Lines
for line := range bytes.Lines(data) {
    process(line)
}

// strings.SplitSeq
for part := range strings.SplitSeq(csv, ",") {
    values = append(values, part)
}
```

## Swiss Tables Map Implementation

```go
// Automatic - 2-3% CPU improvement for map-heavy code
// No code changes needed, just recompile with Go 1.24

// If issues occur, disable with:
// GOEXPERIMENT=noswissmap go build
```

## FIPS 140-3 Mode

```bash
# Enable FIPS-compliant crypto
GOFIPS140=latest go build

# Use specific FIPS version
GOFIPS140=v1.0.0 go build
```

## encoding Interfaces

```go
// New interfaces for append-style encoding
type BinaryAppender interface {
    AppendBinary([]byte) ([]byte, error)
}

type TextAppender interface {
    AppendText([]byte) ([]byte, error)
}

// Implemented by: time.Time, netip.Addr, etc.
// Reduces allocations
buf = t.AppendBinary(buf[:0])
```

## Cgo Annotations

```go
// #cgo noescape funcName
// → tells Go that funcName doesn't let Go pointers escape

// #cgo nocallback funcName
// → tells Go that funcName doesn't call back into Go
```
