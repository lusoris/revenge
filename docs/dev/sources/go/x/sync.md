# Go x/sync Package

> Source: https://pkg.go.dev/golang.org/x/sync
> Fetched: 2026-01-31
> Content-Hash: auto-generated
> Type: html

---

## Overview

The `golang.org/x/sync` module provides Go concurrency primitives beyond those in the standard `sync` and `sync/atomic` packages.

**Module:** `golang.org/x/sync`
**Version:** v0.19.0
**License:** BSD-3-Clause
**Repository:** https://cs.opensource.google/go/x/sync

## Sub-packages

### 1. **errgroup**

Provides synchronization, error propagation, and Context cancellation for groups of goroutines working on subtasks of a common task.

- Handles coordinating multiple goroutines
- Propagates errors from any goroutine
- Supports context-based cancellation

### 2. **semaphore**

Implements a weighted semaphore for controlling concurrent access to resources.

- Controls how many concurrent operations can proceed
- Supports weighted permits for varying resource costs

### 3. **singleflight**

Provides a duplicate function call suppression mechanism.

- Deduplicates concurrent requests for the same result
- Useful for preventing thundering herd problems
- Returns the same result to multiple callers

### 4. **syncmap**

Offers a concurrent map implementation.

- Thread-safe map operations
- Optimized for specific concurrency patterns
- Alternative to using `sync.RWMutex` with regular maps

## Key Features

- Valid go.mod file
- Redistributable license (BSD-3-Clause)
- Tagged version support
- Pre-v1 (not yet stable)

## Resources

- **Go Reference:** https://pkg.go.dev/golang.org/x/sync
- **Report Issues:** Prefix with "x/sync:" at https://go.dev/issues
- **Git Repository:** https://go.googlesource.com/sync
- **Contribution Guide:** https://go.dev/doc/contribute
