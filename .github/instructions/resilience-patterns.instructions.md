---
applyTo: "**/pkg/resilience/**/*.go,**/internal/**/*.go"
---

# Resilience Patterns Instructions

## Overview

Use `pkg/resilience` for fault-tolerance patterns:

- Circuit breaker for external service calls
- Bulkhead for isolation
- Rate limiting for API protection
- Retry for transient failures

## Circuit Breaker

### When to Use

- External API calls (TMDb, TheTVDB, MusicBrainz)
- Blackbeard transcoding service
- Any service that can fail

### Pattern

```go
cb := resilience.NewCircuitBreaker(resilience.CircuitBreakerConfig{
    Name:        "tmdb",
    MaxFailures: 5,             // Open after 5 failures
    Timeout:     30 * time.Second, // Wait before half-open
})

// Execute with protection
err := cb.Execute(func() error {
    return tmdb.FetchMetadata(ctx, id)
})

// Check state
if errors.Is(err, resilience.ErrCircuitOpen) {
    // Use cached data or return degraded response
}
```

### Configuration Guidelines

| Service Type        | MaxFailures | Timeout |
| ------------------- | ----------- | ------- |
| Fast API (metadata) | 5           | 30s     |
| Slow API (search)   | 3           | 60s     |
| Internal service    | 10          | 10s     |

### States

```
Closed (normal) → failures exceed threshold → Open (rejecting)
                                                    ↓
Half-Open (testing) ← timeout expires ←────────────┘
        ↓
successes → Closed
failures → Open
```

## Bulkhead

### When to Use

- Limit concurrent transcodes
- Limit concurrent metadata fetches
- Isolate resource-heavy operations

### Pattern

```go
bulkhead := resilience.NewBulkhead(resilience.BulkheadConfig{
    Name:          "transcoding",
    MaxConcurrent: 10,              // Max parallel operations
    MaxWait:       5 * time.Second, // Wait for slot
    QueueSize:     100,             // Max waiting requests
})

err := bulkhead.Execute(func() error {
    return transcoder.Start(ctx, req)
})

if errors.Is(err, resilience.ErrBulkheadFull) {
    // Return 503 Service Unavailable
}
```

### Sizing Guidelines

| Operation      | MaxConcurrent | QueueSize |
| -------------- | ------------- | --------- |
| Transcoding    | CPU cores     | 100       |
| Metadata fetch | 50            | 500       |
| Image download | 20            | 200       |

## Rate Limiting

### When to Use

- API endpoints (per-user, per-IP)
- External API calls (respect rate limits)
- Abuse prevention

### Per-Key Rate Limiter (Users/IPs)

```go
limiter := resilience.NewPerKeyLimiter(
    resilience.RateLimiterConfig{
        Rate:  100, // requests per second
        Burst: 200, // allow bursts
    },
    time.Minute, // cleanup unused limiters
)

// In middleware
if !limiter.Allow(userID) {
    http.Error(w, "Rate limited", http.StatusTooManyRequests)
    return
}
```

### Blocking Rate Limiter

```go
// Wait for token (for batch operations)
if err := limiter.Wait(ctx, userID); err != nil {
    return err // context cancelled or rate limited
}
```

### Rate Limit Guidelines

| Endpoint Type | Rate  | Burst |
| ------------- | ----- | ----- |
| Read API      | 100/s | 200   |
| Write API     | 20/s  | 50    |
| Search        | 10/s  | 20    |
| Stream start  | 5/s   | 10    |

## Retry

### When to Use

- Transient network errors
- Database connection issues
- Any idempotent operation

### Pattern

```go
retry := resilience.Retry{
    MaxAttempts: 3,
    InitialWait: 100 * time.Millisecond,
    MaxWait:     10 * time.Second,
    Multiplier:  2.0,
    Jitter:      0.1, // ±10% randomness
}

err := retry.DoWithContext(ctx, func(ctx context.Context) error {
    return db.Query(ctx, query)
})
```

### With Circuit Breaker

```go
err := retry.WithCircuitBreaker(cb, func() error {
    return externalAPI.Call()
})
```

### Retry Guidelines

| Operation    | MaxAttempts | InitialWait |
| ------------ | ----------- | ----------- |
| Database     | 3           | 100ms       |
| External API | 3           | 500ms       |
| File I/O     | 2           | 50ms        |

## DO's

- ✅ Use circuit breakers for ALL external services
- ✅ Use bulkheads for resource-intensive operations
- ✅ Use rate limiting at API boundaries
- ✅ Log state changes for debugging
- ✅ Expose metrics for monitoring

## DON'Ts

- ❌ Retry non-idempotent operations
- ❌ Use very short circuit breaker timeouts
- ❌ Forget to handle `ErrCircuitOpen`
- ❌ Set bulkhead too small (causes unnecessary rejections)
- ❌ Rate limit internal service calls

## Registry Pattern

Use registries for dynamic management:

```go
// Create once
registry := resilience.NewCircuitBreakerRegistry(
    resilience.DefaultCircuitBreakerConfig("default"),
)

// Get/create by name
cb := registry.Get("tmdb")
cb := registry.Get("musicbrainz")

// Monitor all
for _, stats := range registry.Stats() {
    logger.Info("circuit breaker",
        "name", stats.Name,
        "state", stats.State,
        "failures", stats.Failures,
    )
}
```
