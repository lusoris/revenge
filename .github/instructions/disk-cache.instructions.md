---
applyTo: "**/internal/service/playback/**/*.go"
---

# Disk Cache Instructions

## Overview

`disk_cache.go` provides persistent transcode caching with quotas. Transcodes survive restarts and can be reused.

## When to Use

- Same content + same profile = cache hit
- Save transcoding resources
- Faster playback start for popular content

## Cache Key

Deterministic key from:

1. Media ID (UUID)
2. Transcode profile ID (string)
3. Source file hash (for invalidation)

```go
cacheKey := sha256(mediaID + profileID + sourceHash)[:32]
```

## Basic Usage

### Initialize

```go
diskCache, err := playback.NewDiskCache(
    playback.DiskCacheConfig{
        BasePath:          "/var/cache/revenge/transcodes",
        MaxSizeBytes:      50 * 1024 * 1024 * 1024, // 50GB
        MaxAgeHours:       72,
        MinFreeSpaceBytes: 10 * 1024 * 1024 * 1024, // 10GB
        PerUserQuotaBytes: 10 * 1024 * 1024 * 1024, // 10GB
        PerMediaQuotaBytes: 5 * 1024 * 1024 * 1024, // 5GB
    },
    logger,
)
if err != nil {
    return err
}

diskCache.Start(ctx)
defer diskCache.Stop()
```

### Check Cache

```go
func (s *Service) StartPlayback(ctx context.Context, req PlaybackRequest) (*PlaybackSession, error) {
    // Generate cache key
    cacheKey := s.diskCache.GenerateCacheKey(
        req.MediaID,
        req.ProfileID,
        req.SourceHash,
    )

    // Check disk cache
    if cached, ok := s.diskCache.GetCached(cacheKey); ok && cached.IsComplete {
        // Serve from disk cache
        return s.playFromDiskCache(ctx, cached)
    }

    // Need to transcode
    return s.transcodeAndCache(ctx, req, cacheKey)
}
```

### Write to Cache

```go
func (s *Service) transcodeAndCache(ctx context.Context, req PlaybackRequest, cacheKey string) error {
    // Create cache entry
    entry, err := s.diskCache.CreateCacheEntry(
        req.MediaID,
        req.UserID,
        req.ProfileID,
        req.SourceHash,
    )
    if err != nil {
        // Quota exceeded - transcode without caching
        return s.transcodeWithoutCache(ctx, req)
    }

    // Start transcode
    stream, err := s.blackbeard.StartTranscode(ctx, req)
    if err != nil {
        return err
    }

    // Write segments as they arrive
    for segment := range stream.Segments {
        // Write to disk cache
        if err := s.diskCache.WriteSegment(cacheKey, segment.Name, segment.Data); err != nil {
            logger.Warn("cache write failed", "error", err)
            // Continue without caching
        }

        // Also cache in memory for immediate use
        s.memoryCache.Store(cacheKey, segment.Name, segment.Data)
    }

    // Write manifest
    s.diskCache.WriteManifest(cacheKey, stream.Manifest)
    s.diskCache.MarkComplete(cacheKey)

    return nil
}
```

### Serve from Cache

```go
func (s *Service) ServeSegment(w http.ResponseWriter, cacheKey, segmentName string) {
    // Try memory cache first
    if data, ok := s.memoryCache.Get(cacheKey, segmentName); ok {
        w.Write(data)
        return
    }

    // Try disk cache
    reader, size, err := s.diskCache.OpenSegmentReader(cacheKey, segmentName)
    if err != nil {
        http.Error(w, "Segment not found", http.StatusNotFound)
        return
    }
    defer reader.Close()

    w.Header().Set("Content-Length", strconv.FormatInt(size, 10))
    io.Copy(w, reader)
}
```

## Quotas

### Global Quota

Total cache size limit:

```go
MaxSizeBytes: 50 * 1024 * 1024 * 1024 // 50GB
```

### Per-User Quota

Prevent one user from filling cache:

```go
PerUserQuotaBytes: 10 * 1024 * 1024 * 1024 // 10GB per user
```

### Per-Media Quota

Limit cache per title (multiple profiles):

```go
PerMediaQuotaBytes: 5 * 1024 * 1024 * 1024 // 5GB per media item
```

## Invalidation

### Source Changed

When source file is modified:

```go
// Recalculate source hash
newHash := calculateFileHash(filePath)

// Invalidate old transcodes
s.diskCache.InvalidateBySourceHash(mediaID, newHash)
```

### Manual Invalidation

```go
// Remove all transcodes for a media item
s.diskCache.InvalidateByMedia(mediaID)
```

### Automatic Expiration

Transcodes older than `MaxAgeHours` are automatically removed.

## Eviction Strategy

1. **Age-based**: Remove transcodes older than `MaxAgeHours`
2. **Incomplete**: Remove incomplete transcodes older than 1 hour
3. **LRU**: When space needed, remove least recently accessed first
4. **Incomplete first**: Prefer evicting incomplete over complete transcodes

## Directory Structure

```
/var/cache/revenge/transcodes/
├── index.json              # Cache index (survives restart)
├── ab/                     # First 2 chars of cache key
│   └── ab1234567890.../
│       ├── master.m3u8     # HLS manifest
│       ├── segment_0.ts
│       ├── segment_1.ts
│       └── ...
└── cd/
    └── cd9876543210.../
        └── ...
```

## Integration with Memory Cache

```
┌─────────────┐    miss     ┌─────────────┐    miss     ┌────────────┐
│   Request   │ ──────────→ │ Memory Cache│ ──────────→ │ Disk Cache │
└─────────────┘             └─────────────┘             └────────────┘
       ↑                           │                          │
       │         hit               │         hit              │
       └───────────────────────────┴──────────────────────────┘
                                                              │
                                                              │ miss
                                                              ▼
                                                     ┌────────────────┐
                                                     │ Blackbeard     │
                                                     │ (transcode)    │
                                                     └────────────────┘
```

## Monitoring

```go
stats := diskCache.Stats()
// DiskCacheStats{
//   TotalEntries: 150,
//   CompleteEntries: 145,
//   IncompleteEntries: 5,
//   TotalSizeBytes: 42949672960,  // 40GB
//   MaxSizeBytes: 53687091200,    // 50GB
//   UsagePercent: 80.0,
//   UniqueMedia: 120,
//   UniqueUsers: 45,
// }
```

## Configuration

```yaml
playback:
  disk_cache:
    enabled: true
    base_path: /var/cache/revenge/transcodes
    max_size_bytes: 53687091200 # 50 GB
    max_age_hours: 72 # 3 days
    eviction_check_interval: 5m
    min_free_space_bytes: 10737418240 # 10 GB free
    per_user_quota_bytes: 10737418240 # 10 GB per user
    per_media_quota_bytes: 5368709120 # 5 GB per media
```

## DO's

- ✅ Check disk cache before transcoding
- ✅ Handle quota exceeded gracefully
- ✅ Invalidate on source file change
- ✅ Use memory cache as L1 for hot segments
- ✅ Monitor cache hit rate

## DON'Ts

- ❌ Block on cache writes
- ❌ Ignore quota errors (fallback to no-cache)
- ❌ Forget to call `MarkComplete()`
- ❌ Store cache on slow storage (SSD recommended)
- ❌ Cache everything (respect quotas)
---

## Related

- [INDEX.instructions.md](INDEX.instructions.md) - Main instruction index with all cross-references
- [ARCHITECTURE_V2.md](../../docs/dev/design/architecture/ARCHITECTURE_V2.md) - System architecture
- [BEST_PRACTICES.md](../../docs/dev/design/operations/BEST_PRACTICES.md) - Best practices
