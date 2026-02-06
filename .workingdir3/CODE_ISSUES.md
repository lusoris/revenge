# Code Issues Found During Doc Rewrite

Issues discovered while reading code for documentation alignment. These are architecture deviations or bugs, not doc problems.

---

## Known Issues

### 1. metadata.BaseClient uses sync.Map instead of proper cache

**Location**: `internal/content/shared/metadata/client.go`
**Severity**: Medium (unbounded memory growth)

The shared `BaseClient` used by all metadata providers (movie, tvshow) uses a raw `sync.Map` as its in-process cache. Problems:
- No size limit - entries accumulate without bound
- No eviction policy (no LRU/LFU)
- "TTL" is only checked on read - expired entries linger in memory until accessed
- Only cleared via explicit `ClearCache()` call

The proper cache infrastructure already exists: `internal/infra/cache/` provides otter L1 (W-TinyLFU, bounded) + Dragonfly L2 (distributed). The movie module's `CachedService` uses it correctly. The metadata `BaseClient` should use it too.

**Fix**: Replace `sync.Map` + `CacheEntry` with the existing `cache.Cache` (otter + Dragonfly), or at minimum use a bounded in-memory cache like otter directly.

### 2. CI/CD workflows need fixing

**Location**: `.github/workflows/`
**Severity**: Medium (broken or misconfigured pipelines)

CI/CD workflows need review and fixing. Specifics TBD after deeper audit.

---

## TODO: Audit for more architecture deviations

During the doc rewrite, only the code structure and interfaces were reviewed. A deeper audit should check for:

- [ ] Other uses of `sync.Map` as cache (should use otter/Dragonfly)
- [ ] Inconsistent error handling patterns across modules
- [ ] Services bypassing the repository layer (direct DB access)
- [ ] Missing cache invalidation (writes that don't invalidate related reads)
- [ ] Hardcoded values that should come from config (timeouts, limits, URLs)
- [ ] Duplicate logic between movie and tvshow modules that should be in shared/
- [ ] Workers missing progress reporting or proper error handling
- [ ] Rate limiters with wrong values or missing entirely
- [ ] Context propagation gaps (functions that don't accept/pass context)
- [ ] Type conversions that silently lose data (int64 → int32, etc.)
