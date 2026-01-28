# Dragonfly Cache Instructions

> Source: https://www.dragonflydb.io/docs, https://github.com/redis/rueidis
> Updated: 2026-01-28

Apply to: `**/internal/infra/cache/**/*.go`, `**/internal/service/**/cache*.go`

## Overview

Dragonfly is a Redis-compatible in-memory datastore. Use `rueidis` client (14x faster than go-redis).

**Why rueidis over go-redis:**

- 14x faster throughput (auto-pipelining)
- Server-assisted client-side caching (RESP3)
- Built-in retry and reconnection
- Zero allocation in hot paths

## Installation

```bash
go get github.com/redis/rueidis
```

## Client Setup

### Basic Client

```go
import "github.com/redis/rueidis"

client, err := rueidis.NewClient(rueidis.ClientOption{
    InitAddress: []string{"localhost:6379"},
})
if err != nil {
    return err
}
defer client.Close()
```

### With Server-Assisted Client Caching

```go
client, err := rueidis.NewClient(rueidis.ClientOption{
    InitAddress: []string{"localhost:6379"},
    ClientTrackingOptions: []string{
        "PREFIX", "cache:", "BCAST",  // Track keys with "cache:" prefix
    },
})
```

### With TLS

```go
client, err := rueidis.NewClient(rueidis.ClientOption{
    InitAddress: []string{"localhost:6379"},
    TLSConfig:   &tls.Config{InsecureSkipVerify: false},
})
```

## Basic Operations

### String Commands

```go
ctx := context.Background()

// SET with expiration
err := client.Do(ctx, client.B().Set().Key("key").Value("value").Ex(10).Build()).Error()

// GET
result, err := client.Do(ctx, client.B().Get().Key("key").Build()).ToString()
if rueidis.IsRedisNil(err) {
    // Key does not exist
}

// SETNX (set if not exists)
ok, err := client.Do(ctx, client.B().Setnx().Key("key").Value("value").Build()).AsBool()
```

### Hash Commands

```go
// HSET
err := client.Do(ctx, client.B().Hset().Key("hash").FieldValue().
    FieldValue("field1", "value1").
    FieldValue("field2", "value2").
    Build()).Error()

// HGET
val, err := client.Do(ctx, client.B().Hget().Key("hash").Field("field1").Build()).ToString()

// HGETALL
vals, err := client.Do(ctx, client.B().Hgetall().Key("hash").Build()).AsStrMap()
```

### List Commands

```go
// LPUSH
err := client.Do(ctx, client.B().Lpush().Key("list").Element("value1", "value2").Build()).Error()

// LRANGE
vals, err := client.Do(ctx, client.B().Lrange().Key("list").Start(0).Stop(-1).Build()).AsStrSlice()

// LPOP
val, err := client.Do(ctx, client.B().Lpop().Key("list").Build()).ToString()
```

### Set Commands

```go
// SADD
err := client.Do(ctx, client.B().Sadd().Key("set").Member("member1", "member2").Build()).Error()

// SMEMBERS
members, err := client.Do(ctx, client.B().Smembers().Key("set").Build()).AsStrSlice()

// SISMEMBER
isMember, err := client.Do(ctx, client.B().Sismember().Key("set").Member("member1").Build()).AsBool()
```

### Sorted Set Commands

```go
// ZADD
err := client.Do(ctx, client.B().Zadd().Key("zset").ScoreMember().ScoreMember(1, "one").Build()).Error()

// ZRANGEBYSCORE with scores
vals, err := client.Do(ctx, client.B().Zrangebyscore().Key("zset").Min("-inf").Max("+inf").Limit(0, 10).Build()).AsStrSlice()
```

## Auto-Pipelining (Automatic!)

rueidis automatically pipelines commands. No manual batching needed:

```go
// These commands are automatically pipelined
go func() { client.Do(ctx, client.B().Set().Key("key1").Value("val1").Build()) }()
go func() { client.Do(ctx, client.B().Set().Key("key2").Value("val2").Build()) }()
go func() { client.Do(ctx, client.B().Get().Key("key1").Build()) }()
```

### Manual Pipeline (for specific ordering)

```go
cmds := make(rueidis.Commands, 0, 3)
cmds = append(cmds, client.B().Set().Key("key1").Value("value1").Build())
cmds = append(cmds, client.B().Set().Key("key2").Value("value2").Build())
cmds = append(cmds, client.B().Get().Key("key1").Build())

results := client.DoMulti(ctx, cmds...)
for _, res := range results {
    if err := res.Error(); err != nil {
        // handle error
    }
}
```

## Transactions (MULTI/EXEC)

```go
// Optimistic locking with WATCH
results, err := client.DoMulti(ctx,
    client.B().Watch().Key("key").Build(),
    client.B().Multi().Build(),
    client.B().Set().Key("key").Value("newvalue").Build(),
    client.B().Exec().Build(),
)
```

## Pub/Sub

```go
// Subscribe
err := client.Receive(ctx, client.B().Subscribe().Channel("channel").Build(), func(msg rueidis.PubSubMessage) {
    fmt.Println(msg.Channel, msg.Message)
})

// Publish
err := client.Do(ctx, client.B().Publish().Channel("channel").Message("hello").Build()).Error()
```

## Lua Scripting

```go
script := rueidis.NewLuaScript(`
    return redis.call('set', KEYS[1], ARGV[1])
`)

result, err := script.Exec(ctx, client, []string{"key"}, []string{"value"})
```

## Client-Side Caching (HUGE Performance Win)

```go
// Enable client-side caching for read-heavy keys
client, _ := rueidis.NewClient(rueidis.ClientOption{
    InitAddress: []string{"localhost:6379"},
    ClientTrackingOptions: []string{"PREFIX", "cache:", "BCAST"},
})

// Use DoCache for cached reads (10 second client-side TTL)
result, err := client.DoCache(ctx,
    client.B().Get().Key("cache:user:123").Cache(),
    10*time.Second,
).ToString()
```

## Error Handling

```go
// Check for nil (key not found)
if rueidis.IsRedisNil(err) {
    // Key does not exist
}

// Check for specific Redis errors
var redisErr *rueidis.RedisError
if errors.As(err, &redisErr) {
    if redisErr.IsLoading() { /* loading */ }
    if redisErr.IsReadOnly() { /* read-only */ }
    if redisErr.IsOOM() { /* out of memory */ }
}
```

## OpenTelemetry Instrumentation

rueidis has built-in OpenTelemetry support via hooks:

```go
import (
    "github.com/redis/rueidis"
    "github.com/redis/rueidis/rueidisotel"
)

client, err := rueidis.NewClient(rueidis.ClientOption{
    InitAddress: []string{"localhost:6379"},
})
if err != nil {
    log.Fatal(err)
}

// Wrap client with OpenTelemetry instrumentation
client = rueidisotel.WithClient(client, rueidisotel.MetricAttrs(
    attribute.String("service.name", "revenge"),
))
```

## Revenge Cache Patterns

### Cache Key Conventions

```
session:{token}           # User sessions (24h TTL)
user:{id}                 # User profiles (5m TTL)
libraries:{user_id}       # Library lists (1m TTL)
search:{module}:{hash}    # Search results (30s TTL)
meta:{provider}:{id}      # Metadata cache (1h TTL)
playback:{session_id}     # Playback state (30m TTL)
```

### Cache Service Interface

```go
type CacheService interface {
    Get(ctx context.Context, key string) ([]byte, error)
    Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
    Delete(ctx context.Context, key string) error
    Exists(ctx context.Context, key string) (bool, error)
}
```

### Session Cache Example (rueidis)

```go
func (c *CacheService) GetSession(ctx context.Context, token string) (*Session, error) {
    key := fmt.Sprintf("session:%s", token)

    // Use DoCache for client-side caching (reduces round trips)
    resp := c.rdb.DoCache(ctx, c.rdb.B().Get().Key(key).Cache(), 5*time.Minute)
    data, err := resp.AsBytes()
    if rueidis.IsRedisNil(err) {
        return nil, ErrSessionNotFound
    }
    if err != nil {
        return nil, fmt.Errorf("cache get failed: %w", err)
    }

    var session Session
    if err := json.Unmarshal(data, &session); err != nil {
        return nil, fmt.Errorf("unmarshal session: %w", err)
    }
    return &session, nil
}

func (c *CacheService) SetSession(ctx context.Context, session *Session) error {
    data, err := json.Marshal(session)
    if err != nil {
        return err
    }
    key := fmt.Sprintf("session:%s", session.Token)
    return c.rdb.Do(ctx, c.rdb.B().Set().Key(key).Value(string(data)).Ex(24*time.Hour).Build()).Error()
}
```

### Cache-Aside Pattern (rueidis)

```go
func (s *UserService) GetUser(ctx context.Context, id uuid.UUID) (*User, error) {
    // Try cache first (with client-side caching)
    key := fmt.Sprintf("user:%s", id)
    resp := s.cache.DoCache(ctx, s.cache.B().Get().Key(key).Cache(), 5*time.Minute)
    data, err := resp.AsBytes()
    if err == nil {
        var user User
        json.Unmarshal(data, &user)
        return &user, nil
    }

    // Cache miss - load from database
    user, err := s.repo.GetByID(ctx, id)
    if err != nil {
        return nil, err
    }

    // Store in cache
    data, _ := json.Marshal(user)
    s.cache.Do(ctx, s.cache.B().Set().Key(key).Value(string(data)).Ex(5*time.Minute).Build())

    return user, nil
}
```

## TTL Guidelines

| Data Type      | TTL    | Rationale                 |
| -------------- | ------ | ------------------------- |
| Sessions       | 24h    | Long-lived user auth      |
| User profiles  | 5m     | Moderate change frequency |
| Library lists  | 1m     | User may modify           |
| Search results | 30s    | Fresh results preferred   |
| Metadata       | 1h     | External API cache        |
| Playback state | 30m    | Active session            |
| Rate limits    | varies | Per endpoint              |

## DO's and DON'Ts

### DO

- ✅ Use `context.Context` for all operations
- ✅ Use `rueidis.IsRedisNil()` for missing keys
- ✅ Use `DoCache()` for read-heavy data (client-side caching)
- ✅ Use auto-pipelining (default in rueidis)
- ✅ Set appropriate TTLs for all cached data
- ✅ Use consistent key naming conventions
- ✅ Use JSON for complex objects
- ✅ Instrument with OpenTelemetry (rueidisotel)

### DON'T

- ❌ Use go-redis/v9 - use rueidis instead (14x faster)
- ❌ Store large objects (>1MB) in cache
- ❌ Use cache as primary data store
- ❌ Forget to handle cache misses gracefully
- ❌ Use blocking operations without timeouts
- ❌ Share keys between modules without namespacing
- ❌ Cache sensitive data without encryption consideration

---

## Related

- [INDEX.instructions.md](INDEX.instructions.md) - Main instruction index
- [otter-local-cache.instructions.md](otter-local-cache.instructions.md) - Local cache (Tier 1)
- [sturdyc-api-cache.instructions.md](sturdyc-api-cache.instructions.md) - API cache (Tier 3)
- [rueidis.md](../../docs/dev/sources/tooling/rueidis.md) - Live rueidis docs
- [DRAGONFLY.md](../../docs/dev/design/integrations/infrastructure/DRAGONFLY.md) - Dragonfly setup
