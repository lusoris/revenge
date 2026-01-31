# Fingerprint Service

<!-- SOURCES: ffmpeg, ffmpeg-codecs, ffmpeg-formats, fx, go-astiav, go-astiav-docs, pgx, postgresql-arrays, postgresql-json, river -->

<!-- DESIGN: services, 01_ARCHITECTURE, 02_DESIGN_PRINCIPLES, 03_METADATA_SYSTEM -->


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
| Integration Testing | 🔴 |## Developer Resources

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


## Related Documents

- [Metadata Service](METADATA.md) - External metadata matching
- [Library Service](LIBRARY.md) - Library scanning integration
- [Search Service](SEARCH.md) - Duplicate detection queries
- [00_SOURCE_OF_TRUTH.md](../00_SOURCE_OF_TRUTH.md) - Service inventory
