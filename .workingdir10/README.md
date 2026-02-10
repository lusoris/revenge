# Working Directory 10 - Leader Election & Cache Integration

## Focus Areas
1. Leader Election (Raft) implementation
2. Rueidis/Otter cache integration fixes
3. Dragonfly cluster emulation configuration
4. Test coverage improvements

## Status
- [x] Raft leader election review
- [x] Rueidis/Otter integration fixes
- [x] Dragonfly cluster emulation fix
- [x] Unit tests passing
- [x] Integration tests passing
- [x] Lint clean
- [x] gosec clean
- [x] trivy clean

---

## Analysis Report

### 1. Raft Leader Election ✅
The Raft implementation in `internal/infra/raft/` looks solid:
- Uses HashiCorp Raft library correctly
- Proper FSM implementation (simpleFSM) for leader-election-only use case
- Good hclog adapter for slog integration
- Lifecycle management via fx hooks

**No major issues found. Implementation is correct for its purpose.**

### 2. Rueidis/Otter Cache Integration ❌ CRITICAL BUGS

#### Bug #1: Not Using Server-Assisted Client-Side Caching
**Current:** Uses `client.Do()` for all Redis operations
**Should use:** `client.DoCache()` for read operations

rueidis has BUILT-IN server-assisted client-side caching:
```go
// CURRENT (wrong)
cmd := client.B().Get().Key(key).Build()
resp := client.Do(ctx, cmd)

// CORRECT - leverages server-assisted invalidation
resp := client.DoCache(ctx, client.B().Get().Key(key).Cache(), ttl)
```

**Impact:**
- No automatic invalidation when data changes
- L1 (otter) cache serves stale data across instances
- Missing the primary benefit of rueidis

#### Bug #2: L1 Cache Not Being Invalidated Across Instances
When instance A writes to Redis, instance B's L1 otter cache is NOT invalidated.
With `DoCache()`, Redis pushes invalidations to ALL connected clients automatically.

#### Bug #3: Redundant Caching Layers
rueidis with `DoCache()` already maintains a client-side cache with:
- Server-assisted invalidation via RESP3 protocol
- Automatic TTL management
- Per-connection cache (configurable via `CacheSizeEachConn`)

The manual L1 otter cache duplicates this without the invalidation benefits.

### 3. Dragonfly Cluster Emulation ❌ BUG

**Bug:** Dragonfly container missing `--cluster_mode=emulated`

rueidis auto-detects cluster mode. Without this flag, commands may fail or behave unexpectedly.

**Fix needed in:**
- `internal/testutil/containers.go` - test containers
- `docker-compose.yml` - production deployment

### 4. Chosen Architecture

**Keep both layers, fix L2 to use DoCache():**
```
Request → L1 (otter, ~100ns) → L2 (rueidis DoCache(), ~200ns) → Dragonfly
                                      ↑
                        Server-assisted invalidation (RESP3)
```

- L1 otter: ultra-fast process-local cache
- L2 DoCache(): leverages rueidis built-in client-side cache WITH automatic invalidation
- When Dragonfly invalidates a key, rueidis pushes to all clients → next L1 miss goes to fresh L2

---

## Bugs Found

### Critical
1. **Cache not using DoCache()** - No server-assisted invalidation
2. **Dragonfly missing cluster_mode** - Commands may fail

### Medium
3. **L1 cache stale across instances** - No multi-instance invalidation
4. **Redundant caching layers** - otter duplicates rueidis client-side cache

### Low
5. **hcLogAdapter args not passed** - Log messages lose context

---

## Full Package Audit

### rueidis (v1.0.71)

| Severity | Issue | Location |
|----------|-------|----------|
| **CRITICAL** | `Do()` used instead of `DoCache()` for reads — no server-assisted invalidation | `cache.go` Get() |
| **MEDIUM** | Lua rate limit uses raw `EVAL` instead of `rueidis.NewLuaScript()` — misses EVALSHA optimization | `middleware/ratelimit_redis.go` |
| **MEDIUM** | No `rueidis.IsRedisNil()` check — cache miss treated as error | `cache.go:66` |
| **LOW** | No `rueidisotel` integration — manual metrics miss built-in `do_cache_miss/hits` | `module.go` |
| **LOW** | `DisableAutoPipelining: false` comment says "disable" but value is false (misleading) | `module.go:107` |
| **INFO** | No retry configuration tuning — uses defaults (fine) | `module.go` |
| **INFO** | No health-check/fallback in cache layer (rate limiter has one, cache doesn't) | `cache.go` |

### otter (v2.3.0)

| Severity | Issue | Location |
|----------|-------|----------|
| **MEDIUM** | 14+ metadata/integration clients never call `Close()`/`StopAllGoroutines()` — goroutine leaks | All provider `client.go` files |
| **MEDIUM** | Transcode pipeline cache (`ttl=0`) falls through to 5min default — likely unintended | `playback/transcode/pipeline.go` |
| **MEDIUM** | Transcode pipeline has no `OnDeletion` handler — evicted FFmpeg processes not killed | `playback/transcode/pipeline.go` |
| **LOW** | `ExpiryWriting` used everywhere — `ExpiryAccessing` better for read-heavy metadata caches | All provider clients |
| **LOW** | `L1Cache[string, any]` loses type safety in 14 provider clients | All provider clients |
| **INFO** | `Keys()` iteration in `DeleteByPrefix` is O(n) — acceptable for current cache sizes | `otter.go:93` |

### HashiCorp Raft (v1.7.3)

| Severity | Issue | Location |
|----------|-------|----------|
| **HIGH** | No cluster join workflow — `AddVoter` never called — multi-node doesn't work | `election.go` |
| **MEDIUM** | Transport/snapshot errors go to `os.Stderr` bypassing structured logging | `election.go:76,81` |
| **MEDIUM** | `hcLogAdapter` drops ALL structured key-value args — cluster debugging impossible | `election.go:218-223` |
| **MEDIUM** | BoltDB stores never closed on shutdown — resource leak | `election.go:87-96` |
| **LOW** | Two separate BoltDB files — one shared store is sufficient | `election.go:87-96` |

### River Queue (v0.30.2)

| Severity | Issue | Location |
|----------|-------|----------|
| **HIGH** | Double leader election (Raft + River) — jobs may silently not run if they disagree | `cleanup_job.go:97-106` |
| **MEDIUM** | No `UniqueOpts` on notification jobs — duplicate sends possible | `notification_job.go:56-62` |
| **LOW** | Cleanup jobs retry 25x (global default) — excessive for idempotent ops | `river.go:67` |
| **LOW** | Event subscription goroutine leak potential on context cancellation | `river.go:91-113` |

### Casbin (v2.135.0)

| Severity | Issue | Location |
|----------|-------|----------|
| **MEDIUM** | `SavePolicy` truncates all rules then re-inserts — brief auth gap on multi-instance | `adapter.go:91` |
| **MEDIUM** | `RemovePolicy` skips empty fields in WHERE — could over-delete | `adapter.go:134-138` |
| **LOW** | No deny policies supported (model uses `allow` only) | `casbin_model.conf` |

### pgx/v5 (v5.8.0)

| Severity | Issue | Location |
|----------|-------|----------|
| **NONE** | Pool config, transactions, `rows.Close()` all correct | — |
| **INFO** | No explicit statement cache tuning (defaults are fine) | `pool.go` |

### ogen (v1.18.0)

| Severity | Issue | Location |
|----------|-------|----------|
| **LOW** | 400 validation errors returned as `*ogen.AdminListUsersForbidden` (403 type) | `handler_admin_users.go:54-61` |
| **NONE** | No manual modifications to generated code | — |

### req/v3 (v3.57.0)

| Severity | Issue | Location |
|----------|-------|----------|
| **MEDIUM** | Ad-hoc `req.C()` in image downloads — no timeout, can hang indefinitely | `images.go:172`, `tmdb/client.go:1238` |
| **LOW** | No `SetCommonRetryCondition` — 4xx errors retried unnecessarily | All 17 clients |

### govips/v2

| Severity | Issue | Location |
|----------|-------|----------|
| **MEDIUM** | `vips.Shutdown()` never called — memory leak on app shutdown | `image_utils.go:19` |
| **LOW** | `vips.Startup(nil)` — no concurrency/cache config for production | `image_utils.go:19` |

### OIDC + OAuth2 (go-oidc v3.17.0)

| Severity | Issue | Location |
|----------|-------|----------|
| **HIGH** | No `nonce` validation in ID token — replay attack vector | `oidc/service.go` |
| **MEDIUM** | OIDC discovery performed on every callback — add caching | `oidc/service.go:320` |
| **MEDIUM** | Hardcoded fallback endpoints (`/authorize`, `/token`) won't work for many IdPs | `oidc/service.go:506-511` |

### OTP/TOTP (pquerna/otp v1.5.0)

| Severity | Issue | Location |
|----------|-------|----------|
| **HIGH** | No TOTP replay protection — same code reusable within ~90s window | `mfa/totp.go:132` |
| **MEDIUM** | Re-enrollment allowed without verifying current TOTP | `mfa/totp.go:97-114` |

### WebAuthn (go-webauthn v0.15.0)

| Severity | Issue | Location |
|----------|-------|----------|
| **MEDIUM** | Clone detection false-positives on authenticators that always report counter=0 | `mfa/webauthn.go:449` |

### AWS SDK v2 (S3)

| Severity | Issue | Location |
|----------|-------|----------|
| **MEDIUM** | `isNotFoundError` uses fragile string matching instead of SDK error types | `storage/s3.go:180-186` |
| **MEDIUM** | No multipart upload — large files buffered entirely in memory | `storage/s3.go:89-95` |
| **LOW** | Custom `contains` reimplements `strings.Contains` | `storage/s3.go:188-196` |

### Typesense (typesense-go v2.0.0)

| Severity | Issue | Location |
|----------|-------|----------|
| **MEDIUM** | Bulk index silently swallows individual document failures | `search/movie_service.go:140-158` |
| **LOW** | Manual URL parsing instead of `net/url.Parse` | `search/module.go:44-69` |
| **LOW** | No timeout/circuit-breaker on client | `search/module.go:78` |

### Prometheus (client_golang v1.23.2)

| Severity | Issue | Location |
|----------|-------|----------|
| **LOW** | `dragonflyCommandsProcessed` is a Gauge for a monotonic value — should be Counter | `collector.go:93` |
| **LOW** | Duplicate River queue metrics (`JobsQueueSize` + `riverQueueSize`) | `metrics.go:144`, `collector.go:19` |

### golang-migrate/migrate (v4.19.1)

| Severity | Issue | Location |
|----------|-------|----------|
| **HIGH** | Dirty database state not handled — blocks all future migrations | `migrate.go:34,68` |

### argon2id (v1.0.0)

| Severity | Issue | Location |
|----------|-------|----------|
| **NONE** | Params align with OWASP recommendations, usage correct | — |
| **INFO** | Test uses p=4, prod uses p=2 — minor inconsistency | `service_unit_test.go:70` |

---

## Summary: All HIGH Issues

1. **OIDC nonce** — No nonce validation in ID token verification
2. **TOTP replay** — Same TOTP code reusable within ~90s validity window
3. **Migrate dirty state** — Dirty database not handled, blocks future migrations
4. **Raft join** — No AddVoter/join workflow — multi-node clustering non-functional
5. **River double leader** — Raft + River leader election conflict — jobs may silently skip
6. **rueidis DoCache()** — Read operations don't use server-assisted caching

---

## Questions

1. Do we want to keep otter L1 for ultra-low-latency reads, or rely solely on rueidis DoCache()?
   - **Decision:** Keep both layers. Use `DoCache()` for L2 reads → server-assisted invalidation. Otter L1 stays for process-local speed.

2. Should we use broadcast mode (`ClientTrackingOptions: BCAST`) for specific key prefixes?
   - **Recommendation:** Not initially. Opt-in mode (default) is safer.

---

## Changes Made

### Dragonfly Cluster Emulation Fix ✅
Added `--cluster_mode=emulated` to Dragonfly in all locations:
- `internal/testutil/containers.go` - test containers
- `docker-compose.yml` - production
- `docker-compose.dev.yml` - development

### Cache DoCache() Migration (PENDING)
**Fix:** Change `client.Do()` → `client.DoCache()` for GET operations in `internal/infra/cache/cache.go`

---

## Next Steps
1. Decide which issues to fix in this working dir
2. Implement fixes
3. Run tests
4. Lint
5. gosec/trivy
3. Run tests
4. Lint
5. gosec/trivy
