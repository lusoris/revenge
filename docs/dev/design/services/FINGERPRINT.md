# Fingerprint Service

> Media file identification via perceptual hashing and acoustic fingerprinting


<!-- TOC-START -->

## Table of Contents

- [Status](#status)
- [Developer Resources](#developer-resources)
- [Overview](#overview)
- [Goals](#goals)
- [Non-Goals](#non-goals)
- [Technical Design](#technical-design)
  - [Fingerprint Types](#fingerprint-types)
  - [Repository Interface](#repository-interface)
  - [Service Layer](#service-layer)
- [Database Schema](#database-schema)
- [River Jobs](#river-jobs)
- [Configuration](#configuration)
- [Implementation Files](#implementation-files)
- [Checklist](#checklist)
- [Sources & Cross-References](#sources-cross-references)
  - [Cross-Reference Indexes](#cross-reference-indexes)
  - [Referenced Sources](#referenced-sources)
- [Related Design Docs](#related-design-docs)
  - [In This Section](#in-this-section)
  - [Related Topics](#related-topics)
  - [Indexes](#indexes)
- [Related Documents](#related-documents)

<!-- TOC-END -->

**Module**: `internal/service/fingerprint`
**Dependencies**: [00_SOURCE_OF_TRUTH.md](../00_SOURCE_OF_TRUTH.md#go-dependencies-media-processing)

## Status

| Dimension | Status |
|-----------|--------|
| Design | ✅ |
| Sources | ✅ |
| Instructions | ✅ |
| Code | 🔴 |
| Linting | 🔴 |
| Unit Testing | 🔴 |
| Integration Testing | 🔴 |
## Developer Resources

> Package versions: [00_SOURCE_OF_TRUTH.md](../00_SOURCE_OF_TRUTH.md#go-dependencies-media-processing--planned)

| Package | Purpose |
|---------|---------|
| goimagehash | Perceptual image hashing |
| go-astiav | FFmpeg bindings |
| Chromaprint | Audio fingerprinting |
| AcoustID API | Audio matching service |
| StashDB API | Video fingerprint matching |

---

## Overview

The Fingerprint service provides content identification through multiple hashing techniques:
- **Perceptual hashing** (pHash) for images and video frames
- **Acoustic fingerprinting** for audio content (Chromaprint/AcoustID)
- **Video fingerprinting** for scene detection and matching

Primary use cases:
- Matching local content against external databases (StashDB, AcoustID)
- Duplicate detection within libraries
- Scene/chapter detection for skip intro/credits

## Goals

- Generate fingerprints for all media types (video, audio, images)
- Support external matching services (StashDB, AcoustID, MusicBrainz)
- Enable duplicate detection across libraries
- Provide scene boundary detection for playback features

## Non-Goals

- Real-time fingerprinting during playback
- Fingerprint-based DRM or content protection
- User-facing fingerprint data (internal service only)

---

## Technical Design

### Fingerprint Types

| Type | Algorithm | Use Case | Package |
|------|-----------|----------|---------|
| Image pHash | Perceptual hash | Poster/artwork matching | `goimagehash` |
| Video pHash | Frame sampling + pHash | Scene matching, StashDB | `goimagehash` + `go-astiav` |
| Audio Chromaprint | Acoustic fingerprint | MusicBrainz/AcoustID | External CLI or CGo |

### Repository Interface

```go
type FingerprintRepository interface {
    // Store fingerprints
    StoreImageHash(ctx context.Context, mediaID uuid.UUID, hash uint64, algorithm string) error
    StoreVideoHash(ctx context.Context, mediaID uuid.UUID, hashes []FrameHash, algorithm string) error
    StoreAudioFingerprint(ctx context.Context, mediaID uuid.UUID, fingerprint []byte, duration int) error

    // Query by similarity
    FindSimilarImages(ctx context.Context, hash uint64, threshold int) ([]Match, error)
    FindSimilarVideos(ctx context.Context, hashes []FrameHash, threshold int) ([]Match, error)

    // External matching
    MatchStashDB(ctx context.Context, videoHashes []FrameHash) (*StashDBMatch, error)
    MatchAcoustID(ctx context.Context, fingerprint []byte, duration int) (*AcoustIDMatch, error)
}

type FrameHash struct {
    Offset   time.Duration
    Hash     uint64
    Algorithm string
}

type Match struct {
    MediaID    uuid.UUID
    Similarity float64
    Source     string // "local", "stashdb", "acoustid"
}
```

### Service Layer

```go
type FingerprintService struct {
    repo     FingerprintRepository
    stashdb  *stashdb.Client
    acoustid *acoustid.Client
    ffmpeg   *astiav.Context
}

func (s *FingerprintService) GenerateImageHash(ctx context.Context, imagePath string) (uint64, error)
func (s *FingerprintService) GenerateVideoHashes(ctx context.Context, videoPath string, sampleCount int) ([]FrameHash, error)
func (s *FingerprintService) GenerateAudioFingerprint(ctx context.Context, audioPath string) ([]byte, int, error)
func (s *FingerprintService) FindDuplicates(ctx context.Context, mediaID uuid.UUID) ([]Match, error)
func (s *FingerprintService) MatchExternal(ctx context.Context, mediaID uuid.UUID) (*ExternalMatch, error)
```

---

## Database Schema

```sql
-- Image/video perceptual hashes
CREATE TABLE media_fingerprints (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    media_id UUID NOT NULL,
    media_type VARCHAR(20) NOT NULL, -- 'image', 'video', 'audio'
    algorithm VARCHAR(20) NOT NULL,  -- 'phash', 'ahash', 'dhash', 'chromaprint'
    hash_value BYTEA NOT NULL,
    frame_offset INTERVAL,           -- For video frame hashes
    duration_ms INT,                 -- For audio fingerprints
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE(media_id, algorithm, frame_offset)
);

CREATE INDEX idx_fingerprints_hash ON media_fingerprints USING hash (hash_value);
CREATE INDEX idx_fingerprints_media ON media_fingerprints (media_id);
```

---

## River Jobs

```go
type GenerateFingerprintArgs struct {
    MediaID   uuid.UUID `json:"media_id"`
    MediaType string    `json:"media_type"` // video, audio, image
    FilePath  string    `json:"file_path"`
}

func (GenerateFingerprintArgs) Kind() string { return "fingerprint.generate" }

type MatchExternalArgs struct {
    MediaID   uuid.UUID `json:"media_id"`
    Services  []string  `json:"services"` // ["stashdb", "acoustid"]
}

func (MatchExternalArgs) Kind() string { return "fingerprint.match_external" }
```

---

## Configuration

```yaml
fingerprint:
  enabled: true
  video:
    sample_count: 10        # Frames to sample per video
    sample_interval: 60s    # Or sample every N seconds
  audio:
    chromaprint_path: ""    # Path to fpcalc binary (optional)
  external:
    stashdb_enabled: true
    acoustid_enabled: true
    acoustid_api_key: ""
```

---

## Implementation Files

| File | Action | Description |
|------|--------|-------------|
| `internal/service/fingerprint/service.go` | CREATE | Core fingerprint service |
| `internal/service/fingerprint/repository.go` | CREATE | Repository interface |
| `internal/service/fingerprint/repository_pg.go` | CREATE | PostgreSQL implementation |
| `internal/service/fingerprint/jobs.go` | CREATE | River job workers |
| `internal/service/fingerprint/module.go` | CREATE | fx module |
| `migrations/shared/000XXX_fingerprints.up.sql` | CREATE | Database schema |

---

## Checklist

- [ ] Database migration created
- [ ] Repository interface defined
- [ ] PostgreSQL repository implemented
- [ ] Service layer with hash generation
- [ ] StashDB matching integration
- [ ] AcoustID matching integration
- [ ] River jobs for background processing
- [ ] Tests written
- [ ] Documentation updated

---


<!-- SOURCE-BREADCRUMBS-START -->

## Sources & Cross-References

> Auto-generated section linking to external documentation sources

### Cross-Reference Indexes

- [All Sources Index](../../sources/SOURCES_INDEX.md) - Complete list of external documentation
- [Design ↔ Sources Map](../../sources/DESIGN_CROSSREF.md) - Which docs reference which sources

### Referenced Sources

| Source | Documentation |
|--------|---------------|
| [FFmpeg Codecs](https://ffmpeg.org/ffmpeg-codecs.html) | [Local](../../sources/media/ffmpeg-codecs.md) |
| [FFmpeg Documentation](https://ffmpeg.org/ffmpeg.html) | [Local](../../sources/media/ffmpeg.md) |
| [FFmpeg Formats](https://ffmpeg.org/ffmpeg-formats.html) | [Local](../../sources/media/ffmpeg-formats.md) |
| [PostgreSQL Arrays](https://www.postgresql.org/docs/current/arrays.html) | [Local](../../sources/database/postgresql-arrays.md) |
| [PostgreSQL JSON Functions](https://www.postgresql.org/docs/current/functions-json.html) | [Local](../../sources/database/postgresql-json.md) |
| [River Job Queue](https://pkg.go.dev/github.com/riverqueue/river) | [Local](../../sources/tooling/river.md) |
| [Uber fx](https://pkg.go.dev/go.uber.org/fx) | [Local](../../sources/tooling/fx.md) |
| [go-astiav (FFmpeg bindings)](https://pkg.go.dev/github.com/asticode/go-astiav) | [Local](../../sources/media/go-astiav.md) |
| [pgx PostgreSQL Driver](https://pkg.go.dev/github.com/jackc/pgx/v5) | [Local](../../sources/database/pgx.md) |

<!-- SOURCE-BREADCRUMBS-END -->

<!-- DESIGN-BREADCRUMBS-START -->

## Related Design Docs

> Auto-generated cross-references to related design documentation

**Category**: [Services](INDEX.md)

### In This Section

- [Activity Service](ACTIVITY.md)
- [Analytics Service](ANALYTICS.md)
- [API Keys Service](APIKEYS.md)
- [Auth Service](AUTH.md)
- [Grants Service](GRANTS.md)
- [Library Service](LIBRARY.md)
- [Metadata Service](METADATA.md)
- [Notification Service](NOTIFICATION.md)

### Related Topics

- [Revenge - Architecture v2](../architecture/01_ARCHITECTURE.md) _Architecture_
- [Revenge - Design Principles](../architecture/02_DESIGN_PRINCIPLES.md) _Architecture_
- [Revenge - Metadata System](../architecture/03_METADATA_SYSTEM.md) _Architecture_
- [Revenge - Player Architecture](../architecture/04_PLAYER_ARCHITECTURE.md) _Architecture_
- [Plugin Architecture Decision](../architecture/05_PLUGIN_ARCHITECTURE_DECISION.md) _Architecture_

### Indexes

- [Design Index](../DESIGN_INDEX.md) - All design docs by category/topic
- [Source of Truth](../00_SOURCE_OF_TRUTH.md) - Package versions and status

<!-- DESIGN-BREADCRUMBS-END -->

## Related Documents

- [Metadata Service](METADATA.md) - External metadata matching
- [Library Service](LIBRARY.md) - Library scanning integration
- [Search Service](SEARCH.md) - Duplicate detection queries
- [00_SOURCE_OF_TRUTH.md](../00_SOURCE_OF_TRUTH.md) - Service inventory
