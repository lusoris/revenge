# Metadata Enrichment Pattern

<!-- SOURCES: river, rueidis -->

<!-- DESIGN: patterns, 01_ARCHITECTURE, 02_DESIGN_PRINCIPLES, 03_METADATA_SYSTEM -->


> Patterns for enriching content metadata from external providers

**Source of Truth**: [00_SOURCE_OF_TRUTH.md](../00_SOURCE_OF_TRUTH.md)

---

## Status

| Dimension | Status | Notes |
|-----------|--------|-------|
| Design | 🟡 | Scaffold |
| Sources | 🔴 |  |
| Instructions | 🔴 |  |
| Code | 🔴 |  |
| Linting | 🔴 |  |
| Unit Testing | 🔴 |  |
| Integration Testing | 🔴 |  |
---

## Overview

Metadata enrichment follows the priority chain defined in SOURCE_OF_TRUTH:

```
Priority Order (ALWAYS):
1. LOCAL CACHE     → First, instant UI display
2. ARR SERVICES    → Radarr, Sonarr, Whisparr (cached metadata)
3. INTERNAL        → Stash-App (if connected)
4. EXTERNAL        → TMDb, StashDB.org, MusicBrainz, etc.
5. ENRICHMENT      → Background jobs, lower priority, seamless
```

---

## Enrichment Job Pattern

### Job Definition

```go
type MetadataEnrichArgs struct {
    ContentType string    `json:"content_type"` // movie, tvshow, music, etc.
    ContentID   uuid.UUID `json:"content_id"`
    Providers   []string  `json:"providers"`    // tmdb, musicbrainz, etc.
}

func (w *EnrichWorker) Work(ctx context.Context, job *river.Job[MetadataEnrichArgs]) error {
    for _, provider := range job.Args.Providers {
        if err := w.enrichFromProvider(ctx, job.Args, provider); err != nil {
            slog.Warn("enrichment failed", "provider", provider, "error", err)
            continue // Try next provider
        }
    }
    return nil
}
```

### Provider Interface

```go
type MetadataProvider interface {
    Name() string
    Search(ctx context.Context, query string) ([]SearchResult, error)
    GetDetails(ctx context.Context, id string) (*Metadata, error)
    SupportsContentType(contentType string) bool
}
```

---

## Caching Strategy

### Three-Layer Cache

```go
func (s *MetadataService) GetMetadata(ctx context.Context, id uuid.UUID) (*Metadata, error) {
    // L1: In-memory (otter)
    if cached, ok := s.localCache.Get(id.String()); ok {
        return cached.(*Metadata), nil
    }

    // L2: Distributed (rueidis)
    if cached, err := s.redisCache.Get(ctx, id.String()); err == nil {
        s.localCache.Set(id.String(), cached)
        return cached, nil
    }

    // L3: Database
    metadata, err := s.repo.GetByID(ctx, id)
    if err != nil {
        return nil, err
    }

    // Populate caches
    s.localCache.Set(id.String(), metadata)
    s.redisCache.Set(ctx, id.String(), metadata, time.Hour)

    return metadata, nil
}
```

---

## Background Enrichment

### Trigger Conditions

| Trigger | Action |
|---------|--------|
| New content added | Queue enrichment job |
| User views content | Check if enrichment needed |
| Scheduled refresh | Re-enrich stale metadata |
| Manual request | Immediate enrichment |

### Staleness Check

```go
func (s *MetadataService) needsEnrichment(m *Metadata) bool {
    if m.EnrichedAt.IsZero() {
        return true
    }
    return time.Since(m.EnrichedAt) > 7*24*time.Hour
}
```

---

## Related

- [Metadata System](../architecture/03_METADATA_SYSTEM.md)
- [TMDb Integration](../integrations/metadata/video/TMDB.md)
- [TheTVDB Integration](../integrations/metadata/video/THETVDB.md)
- [MusicBrainz Integration](../integrations/metadata/music/MUSICBRAINZ.md)
