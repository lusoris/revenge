# Bug #19: Invalid Expire Time with TTL < 1 Second

**Severity**: Medium
**Component**: Cache (L2/Redis)
**Found By**: Integration test `TestTTLAccuracy` and `TestEdgeCases`

## Symptom

Setting cache values with TTL less than 1 second fails with error:
```
L2 cache set failed: invalid expire time in 'set' command
```

### Specific Failures:
1. **500ms TTL**: `Error: L2 cache set failed: invalid expire time in 'set' command`
2. **Zero TTL** (0s): `Error: L2 cache set failed: invalid expire time in 'set' command`
3. **Negative TTL** (-1s): `Error: L2 cache set failed: invalid expire time in 'set' command`

## Test That Reproduces

```go
func TestTTLAccuracy(t *testing.T) {
    c := newTestCache(t)
    defer c.Close()
    ctx := context.Background()

    // Fails with 500ms TTL
    key := "test:ttl:500ms"
    value := []byte("expires")

    err := c.Set(ctx, key, value, 500*time.Millisecond)
    // ERROR: L2 cache set failed: invalid expire time in 'set' command
    require.NoError(t, err)
}
```

## Root Cause Analysis

**Location**: `internal/infra/cache/cache.go` - `Set` method

The issue is in how we convert Go's `time.Duration` to Redis `EX` (seconds) parameter:

```go
func (c *Cache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
    // Store in L1
    c.l1.Set(key, value)

    // Store in L2 if available
    if c.client != nil && c.client.rueidisClient != nil {
        cmd := c.client.rueidisClient.B().Set().Key(key).Value(string(value)).Ex(ttl).Build()
        //                                                                        ^^
        // Problem: Ex() expects duration, but Redis SET EX parameter is in SECONDS
        // When ttl < 1 second, Ex() likely truncates to 0, which is invalid

        if err := c.client.rueidisClient.Do(ctx, cmd).Error(); err != nil {
            return fmt.Errorf("L2 cache set failed: %w", err)
        }
    }

    return nil
}
```

**The problem**: Redis `SET key value EX seconds` requires `seconds` as an integer. When TTL is less than 1 second:
- 500ms → 0 seconds → Invalid
- 0ms → 0 seconds → Invalid
- -1s → negative → Invalid

## Impact Assessment

**Production Impact**: MEDIUM

**When This Matters**:
1. **Short-lived tokens**: Session tokens, CSRF tokens with sub-second TTL
2. **Rate limiting**: Sub-second rate limit windows
3. **Temporary locks**: Distributed locks with millisecond precision
4. **Cache warming**: Quick invalidation scenarios

**Current Behavior**:
- L1 cache works correctly (Otter supports any duration)
- L2 cache fails silently (returns error but data is in L1)
- Inconsistent state: Data in L1 but not L2

## Fix Options

### Option 1: Use PX (milliseconds) for sub-second TTL ✅ RECOMMENDED

```go
func (c *Cache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
    c.l1.Set(key, value)

    if c.client != nil && c.client.rueidisClient != nil {
        var cmd rueidis.Completed

        // Use PX (milliseconds) for sub-second precision
        if ttl < time.Second {
            cmd = c.client.rueidisClient.B().Set().
                Key(key).
                Value(string(value)).
                Px(ttl).  // PX for millisecond precision
                Build()
        } else {
            cmd = c.client.rueidisClient.B().Set().
                Key(key).
                Value(string(value)).
                Ex(ttl).  // EX for second precision
                Build()
        }

        if err := c.client.rueidisClient.Do(ctx, cmd).Error(); err != nil {
            return fmt.Errorf("L2 cache set failed: %w", err)
        }
    }

    return nil
}
```

### Option 2: Validate and reject invalid TTL

```go
func (c *Cache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
    if ttl <= 0 {
        return fmt.Errorf("invalid TTL: must be positive, got %v", ttl)
    }

    if ttl < time.Millisecond {
        return fmt.Errorf("invalid TTL: minimum is 1ms, got %v", ttl)
    }

    // ... rest of code
}
```

### Option 3: Always use PX (milliseconds)

```go
func (c *Cache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
    c.l1.Set(key, value)

    if c.client != nil && c.client.rueidisClient != nil {
        // Always use PX for millisecond precision
        cmd := c.client.rueidisClient.B().Set().
            Key(key).
            Value(string(value)).
            Px(ttl).  // Always use milliseconds
            Build()

        if err := c.client.rueidisClient.Do(ctx, cmd).Error(); err != nil {
            return fmt.Errorf("L2 cache set failed: %w", err)
        }
    }

    return nil
}
```

## Recommended Solution

**Option 1** (conditional PX/EX) provides the best balance:
- Sub-second precision when needed
- Standard seconds for longer TTLs (more efficient)
- No breaking changes to existing code

## Verification Steps

1. Fix the code
2. Run tests:
   ```bash
   go test ./tests/integration/cache -tags=integration -run TestTTLAccuracy -v
   go test ./tests/integration/cache -tags=integration -run TestEdgeCases -v
   ```
3. Verify all TTL tests pass
4. Check that both L1 and L2 have the data
5. Verify expiration works correctly

## Prevention

- Add validation for TTL values at API boundaries
- Document TTL precision limits (1ms minimum)
- Add unit tests for edge case TTLs
- Consider adding TTL constants:
  ```go
  const (
      MinTTL = 1 * time.Millisecond
      MaxTTL = 365 * 24 * time.Hour
  )
  ```

## Related Issues

- Bug #18: Cache eviction async behavior (documented, not a bug)
- This affects only L2 cache, L1 works correctly

## Test Results

Before fix:
```
--- FAIL: TestTTLAccuracy (3.40s)
    --- FAIL: TestTTLAccuracy/500ms (0.00s)
        Error: L2 cache set failed: invalid expire time in 'set' command
    --- FAIL: TestTTLAccuracy/1s (1.20s)
        Error: An error is expected but got nil. (key didn't expire)
    --- FAIL: TestTTLAccuracy/2s (2.20s)
        Error: An error is expected but got nil. (key didn't expire)
```

Expected after fix:
```
--- PASS: TestTTLAccuracy (3.40s)
    --- PASS: TestTTLAccuracy/500ms (0.50s)
    --- PASS: TestTTLAccuracy/1s (1.20s)
    --- PASS: TestTTLAccuracy/2s (2.20s)
```

---

**Next Steps**: Implement Option 1 (conditional PX/EX), run tests, verify fix.
