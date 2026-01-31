# Go x/time/rate Package

> Source: https://pkg.go.dev/golang.org/x/time/rate
> Fetched: 2026-01-31
> Content-Hash: auto-generated
> Type: html

---

## Overview

Package `rate` provides a **rate limiter** implementation using a token bucket algorithm. It controls how frequently events are allowed to happen.

**Module:** `golang.org/x/time`
**Version:** v0.14.0
**License:** BSD-3-Clause
**Imported by:** 13,352 packages

## Core Types

### Limit

- **Type**: `float64`
- **Definition**: Maximum frequency of events (events per second)
- **Special values**:
  - `Inf`: Infinite rate limit (allows all events)
  - `0`: Allows no events

### Limiter

A rate limiter controlling event frequency with:
- **Token bucket** of size `b` (burst), initially full
- **Refill rate** of `r` tokens per second
- **Thread-safe** for concurrent use

#### Key Methods

| Method | Description |
|--------|-------------|
| `Allow()` / `AllowN(t, n)` | Returns `false` if no tokens available |
| `Reserve()` / `ReserveN(t, n)` | Returns a `Reservation` with delay information |
| `Wait(ctx)` / `WaitN(ctx, n)` | Blocks until tokens available or context canceled |
| `Tokens()` / `TokensAt(t)` | Query available tokens |
| `SetLimit(newLimit)` / `SetBurst(newBurst)` | Dynamically adjust rates |

### Reservation

Holds information about permitted events after a delay:
- `OK()` - Whether tokens can be granted
- `Delay()` / `DelayFrom(t)` - Wait duration before acting
- `Cancel()` / `CancelAt(t)` - Cancel reservation, refund tokens

### Sometimes

Performs actions occasionally based on:
- `First: N` - First N calls run the function
- `Every: M` - Every Mth call runs the function
- `Interval` - Run if elapsed time since last execution

## Usage Examples

```go
// Create a limiter: 100 events/sec with burst of 10
limiter := rate.NewLimiter(100, 10)

// Or using Every() for interval-based limits
limiter := rate.NewLimiter(rate.Every(100*time.Millisecond), 1)

// Simple check (drop events)
if limiter.Allow() {
    // Handle event
}

// Reserve with delay
r := limiter.Reserve()
if r.OK() {
    time.Sleep(r.Delay())
    // Handle event
}

// Wait with context
if err := limiter.Wait(ctx); err == nil {
    // Handle event
}
```

## Constants

- `Inf` - Infinite rate limit
- `InfDuration` - Duration returned when reservation cannot be satisfied

## Resources

- **Go Reference:** https://pkg.go.dev/golang.org/x/time/rate
- **Repository:** cs.opensource.google/go/x/time
