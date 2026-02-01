

---
sources:
  - name: Chromaprint/AcoustID
    url: ../../../sources/standards/chromaprint.md
    note: Auto-resolved from chromaprint-acoustid
  - name: FFmpeg Documentation
    url: ../../../sources/media/ffmpeg.md
    note: Auto-resolved from ffmpeg
  - name: FFmpeg Codecs
    url: ../../../sources/media/ffmpeg-codecs.md
    note: Auto-resolved from ffmpeg-codecs
  - name: FFmpeg Formats
    url: ../../../sources/media/ffmpeg-formats.md
    note: Auto-resolved from ffmpeg-formats
  - name: go-astiav (FFmpeg bindings)
    url: ../../../sources/media/go-astiav.md
    note: Auto-resolved from go-astiav
  - name: go-astiav GitHub README
    url: ../../../sources/media/go-astiav-guide.md
    note: Auto-resolved from go-astiav-docs
  - name: River Job Queue
    url: ../../../sources/tooling/river.md
    note: Auto-resolved from river
design_refs:
  - title: 01_ARCHITECTURE
    path: ../../architecture/01_ARCHITECTURE.md
  - title: 02_DESIGN_PRINCIPLES
    path: ../../architecture/02_DESIGN_PRINCIPLES.md
  - title: 03_METADATA_SYSTEM
    path: ../../architecture/03_METADATA_SYSTEM.md
---

## Table of Contents

- [Skip Intro / Credits Detection](#skip-intro-credits-detection)
  - [Status](#status)
  - [Architecture](#architecture)
    - [Database Schema](#database-schema)
    - [Module Structure](#module-structure)
    - [Component Interaction](#component-interaction)
  - [Implementation](#implementation)
    - [File Structure](#file-structure)
    - [Key Interfaces](#key-interfaces)
    - [Dependencies](#dependencies)
  - [Configuration](#configuration)
    - [Environment Variables](#environment-variables)
- [Feature toggle](#feature-toggle)
- [Detection settings](#detection-settings)
- [Chromaprint settings](#chromaprint-settings)
- [Credits detection](#credits-detection)
- [Worker](#worker)
- [User defaults](#user-defaults)
    - [Config Keys](#config-keys)
  - [API Endpoints](#api-endpoints)
    - [Content Management](#content-management)
      - [POST /api/v1/skip-intro/detect](#post-apiv1skip-introdetect)
      - [GET /api/v1/skip-intro/markers](#get-apiv1skip-intromarkers)
      - [POST /api/v1/skip-intro/markers](#post-apiv1skip-intromarkers)
      - [PUT /api/v1/skip-intro/markers/:id](#put-apiv1skip-intromarkersid)
      - [DELETE /api/v1/skip-intro/markers/:id](#delete-apiv1skip-intromarkersid)
      - [POST /api/v1/skip-intro/markers/:id/verify](#post-apiv1skip-intromarkersidverify)
      - [GET /api/v1/skip-intro/preferences](#get-apiv1skip-intropreferences)
      - [PUT /api/v1/skip-intro/preferences](#put-apiv1skip-intropreferences)
      - [POST /api/v1/skip-intro/skip](#post-apiv1skip-introskip)
      - [POST /api/v1/skip-intro/bulk-detect](#post-apiv1skip-introbulk-detect)
      - [GET /api/v1/skip-intro/series/:id/patterns](#get-apiv1skip-introseriesidpatterns)
  - [Testing Strategy](#testing-strategy)
    - [Unit Tests](#unit-tests)
    - [Integration Tests](#integration-tests)
    - [Test Coverage](#test-coverage)
  - [Related Documentation](#related-documentation)
    - [Design Documents](#design-documents)
    - [External Sources](#external-sources)


# Skip Intro / Credits Detection


**Created**: 2026-01-31
**Status**: ✅ Complete
**Category**: feature


> Content module for 

> Automatic intro and credits detection with one-click skip

---


## Status

| Dimension | Status | Notes |
|-----------|--------|-------|
| Design | ✅ | - |
| Sources | ✅ | - |
| Instructions | ✅ | - |
| Code | 🔴 | - |
| Linting | 🔴 | - |
| Unit Testing | 🔴 | - |
| Integration Testing | 🔴 | - |

**Overall**: ✅ Complete



---


## Architecture

### Database Schema

**Schema**: `public`

<!-- Schema diagram -->

### Module Structure

```
internal/content/skip_intro_/_credits_detection/
├── module.go              # fx module definition
├── repository.go          # Database operations
├── service.go             # Business logic
├── handler.go             # HTTP handlers (ogen)
├── types.go               # Domain types
└── skip_intro_/_credits_detection_test.go
```

### Component Interaction

<!-- Component interaction diagram -->


## Implementation

### File Structure

```
internal/playback/skipintro/
├── module.go                    # fx module registration
├── repository.go                # Database operations (sqlc)
├── queries.sql                  # SQL queries for sqlc
├── service.go                   # Business logic
├── handler.go                   # HTTP handlers (ogen-generated)
├── types.go                     # Domain types
├── detector.go                  # Detection orchestration
├── audio_detector.go            # Audio fingerprinting (Chromaprint)
├── visual_detector.go           # Visual analysis (black frames, silence)
├── pattern_matcher.go           # Series pattern matching
├── chromaprint.go               # Chromaprint integration
├── ffmpeg.go                    # FFmpeg audio/video extraction
├── queue.go                     # Detection queue management
└── cache.go                     # Caching layer (otter)

cmd/server/
└── main.go                      # Server entry point with fx

migrations/
├── 032_skip_intro.up.sql        # Skip intro tables
└── 032_skip_intro.down.sql      # Rollback

api/openapi/
└── skipintro.yaml               # OpenAPI spec for skip intro

web/src/lib/components/player/
├── SkipIntroButton.svelte       # Skip intro button overlay
└── SkipPreferences.svelte       # User preferences UI
```


### Key Interfaces

```go
// Repository interface for skip intro database operations
type Repository interface {
    // Markers
    CreateMarker(ctx context.Context, params CreateMarkerParams) (*IntroCreditsMarker, error)
    GetMarker(ctx context.Context, id uuid.UUID) (*IntroCreditsMarker, error)
    GetMarkersForContent(ctx context.Context, contentType string, contentID uuid.UUID) ([]*IntroCreditsMarker, error)
    UpdateMarker(ctx context.Context, id uuid.UUID, params UpdateMarkerParams) (*IntroCreditsMarker, error)
    DeleteMarker(ctx context.Context, id uuid.UUID) error
    VerifyMarker(ctx context.Context, markerID, userID uuid.UUID) error

    // Series patterns
    CreateSeriesPattern(ctx context.Context, params CreatePatternParams) (*SeriesIntroPattern, error)
    GetSeriesPatterns(ctx context.Context, seriesID uuid.UUID, patternType string) ([]*SeriesIntroPattern, error)
    UpdatePatternStats(ctx context.Context, id uuid.UUID, matchedCount int, avgConfidence float64) error
    DeactivatePattern(ctx context.Context, id uuid.UUID) error

    // Detection queue
    EnqueueDetection(ctx context.Context, contentType string, contentID uuid.UUID, priority int) error
    GetQueuedItems(ctx context.Context, limit int) ([]*DetectionQueueItem, error)
    UpdateQueueStatus(ctx context.Context, id uuid.UUID, status, errorMsg string) error

    // Preferences
    GetUserPreferences(ctx context.Context, userID uuid.UUID) (*IntroSkipPreferences, error)
    UpsertUserPreferences(ctx context.Context, userID uuid.UUID, prefs PreferencesParams) (*IntroSkipPreferences, error)

    // Statistics
    RecordSkip(ctx context.Context, params SkipStatsParams) error
    GetSkipStats(ctx context.Context, markerID uuid.UUID) (*SkipStatistics, error)
}

// Service interface for skip intro operations
type Service interface {
    // Detection
    DetectIntroCredits(ctx context.Context, contentType string, contentID uuid.UUID) ([]*IntroCreditsMarker, error)
    GetMarkers(ctx context.Context, contentType string, contentID uuid.UUID) ([]*IntroCreditsMarker, error)
    QueueBulkDetection(ctx context.Context, contentRefs []ContentReference, priority int) error

    // Manual management
    CreateManualMarker(ctx context.Context, req CreateManualMarkerRequest) (*IntroCreditsMarker, error)
    UpdateMarker(ctx context.Context, markerID uuid.UUID, updates MarkerUpdates) (*IntroCreditsMarker, error)
    DeleteMarker(ctx context.Context, markerID uuid.UUID) error
    VerifyMarker(ctx context.Context, markerID, userID uuid.UUID) error

    // Preferences
    GetPreferences(ctx context.Context, userID uuid.UUID) (*IntroSkipPreferences, error)
    UpdatePreferences(ctx context.Context, userID uuid.UUID, updates PreferencesUpdate) (*IntroSkipPreferences, error)

    // Statistics
    RecordSkipAction(ctx context.Context, userID uuid.UUID, markerID uuid.UUID, skipped bool) error
}

// Detector interface for intro/credits detection
type Detector interface {
    DetectIntro(ctx context.Context, contentType string, contentID uuid.UUID, videoPath string) (*DetectionResult, error)
    DetectCredits(ctx context.Context, contentType string, contentID uuid.UUID, videoPath string) (*DetectionResult, error)
    DetectRecap(ctx context.Context, contentType string, contentID uuid.UUID, videoPath string) (*DetectionResult, error)
}

// AudioDetector interface for audio fingerprinting
type AudioDetector interface {
    ExtractFingerprint(ctx context.Context, audioPath string, startSec, durationSec float64) (*AudioFingerprint, error)
    CompareFingerprints(fp1, fp2 *AudioFingerprint) (similarity float64, err error)
    FindMatchingSegment(ctx context.Context, haystack, needle *AudioFingerprint) (*MatchResult, error)
}

// VisualDetector interface for visual analysis
type VisualDetector interface {
    DetectBlackFrames(ctx context.Context, videoPath string) ([]BlackFrameSegment, error)
    DetectSilence(ctx context.Context, audioPath string) ([]SilenceSegment, error)
    DetectSceneChanges(ctx context.Context, videoPath string) ([]SceneChange, error)
}

// PatternMatcher interface for series-wide matching
type PatternMatcher interface {
    CreatePattern(ctx context.Context, seriesID, episodeID uuid.UUID, marker *IntroCreditsMarker) (*SeriesIntroPattern, error)
    MatchAgainstPatterns(ctx context.Context, seriesID uuid.UUID, episodeID uuid.UUID, fingerprint *AudioFingerprint) ([]*PatternMatch, error)
    UpdatePatternStatistics(ctx context.Context, patternID uuid.UUID) error
}

// ChromaprintClient interface for Chromaprint integration
type ChromaprintClient interface {
    GenerateFingerprint(audioData []byte, sampleRate int, channels int) (string, error)
    CompareFinger prints(fp1, fp2 string) (float64, error)
}

// FFmpegClient interface for audio/video extraction
type FFmpegClient interface {
    ExtractAudio(ctx context.Context, videoPath, outputPath string, startSec, durationSec float64) error
    ExtractFrames(ctx context.Context, videoPath, outputDir string, interval float64) error
    DetectBlackFrames(ctx context.Context, videoPath string, threshold float64) ([]TimeRange, error)
    DetectSilence(ctx context.Context, audioPath string, threshold float64, minDuration float64) ([]TimeRange, error)
}
```


### Dependencies

**Go Packages**:
```go
require (
    // Core
    github.com/google/uuid v1.6.0
    go.uber.org/fx v1.23.0

    // Database
    github.com/jackc/pgx/v5 v5.7.2
    github.com/sqlc-dev/sqlc v1.28.0

    // API
    github.com/ogen-go/ogen v1.7.0

    // Caching
    github.com/maypok86/otter v1.2.4

    // FFmpeg
    github.com/asticode/go-astiav v0.23.0  // FFmpeg bindings

    // Chromaprint (audio fingerprinting)
    // Note: Use CGO bindings to libchromaprint or subprocess to fpcalc
    github.com/acoustid/chromaprint v1.0.0  // Hypothetical Go bindings

    // Job queue
    github.com/riverqueue/river v0.15.0

    // Testing
    github.com/stretchr/testify v1.10.0
    github.com/testcontainers/testcontainers-go v0.35.0
)
```

**External Dependencies**:
- **FFmpeg 7.1+**: Audio/video extraction, black frame detection, silence detection
- **Chromaprint (fpcalc)**: Audio fingerprinting
- **PostgreSQL 18+**: Database






## Configuration
### Environment Variables

```bash
# Feature toggle
SKIP_INTRO_ENABLED=true
SKIP_INTRO_AUTO_DETECT=true              # Auto-detect on library scan

# Detection settings
SKIP_INTRO_MIN_DURATION_SEC=15           # Min intro duration
SKIP_INTRO_MAX_DURATION_SEC=180          # Max intro duration (3 minutes)
SKIP_INTRO_MIN_CONFIDENCE=70             # Min confidence % to accept
SKIP_INTRO_USE_AUDIO_FINGERPRINT=true   # Use Chromaprint
SKIP_INTRO_USE_VISUAL_ANALYSIS=true     # Use black frame/silence detection
SKIP_INTRO_USE_SERIES_PATTERNS=true     # Match against series patterns

# Chromaprint settings
CHROMAPRINT_BIN_PATH=/usr/bin/fpcalc     # Path to fpcalc binary
CHROMAPRINT_SAMPLE_DURATION_SEC=60       # Sample duration for fingerprinting

# Credits detection
SKIP_CREDITS_ENABLED=true
SKIP_CREDITS_START_OFFSET_SEC=300        # Start checking X seconds before end

# Worker
SKIP_INTRO_WORKER_CONCURRENCY=2          # Concurrent detection jobs
SKIP_INTRO_WORKER_PRIORITY=5             # Job priority (1-10)
SKIP_INTRO_DETECTION_TIMEOUT_MIN=10      # Timeout for detection

# User defaults
SKIP_INTRO_DEFAULT_AUTO_SKIP=false       # Default auto-skip setting
SKIP_INTRO_BUTTON_DURATION_SEC=10        # Show button for X seconds
```


### Config Keys

```yaml
skip_intro:
  # Feature toggle
  enabled: true
  auto_detect: true                 # Auto-detect on library scan

  # Detection settings
  detection:
    intro:
      min_duration_seconds: 15
      max_duration_seconds: 180     # 3 minutes
      min_confidence: 70
    credits:
      enabled: true
      start_offset_seconds: 300     # Check last 5 minutes
    recap:
      enabled: true
      max_duration_seconds: 120     # 2 minutes

    # Detection methods
    methods:
      audio_fingerprint:
        enabled: true
        sample_duration_seconds: 60
      visual_analysis:
        enabled: true
        black_frame_threshold: 0.98  # 98% black
        silence_threshold_db: -40
        silence_duration_seconds: 2
      series_patterns:
        enabled: true
        min_match_confidence: 80

  # Chromaprint
  chromaprint:
    bin_path: /usr/bin/fpcalc
    sample_rate: 16000
    channels: 1

  # Worker
  worker:
    concurrency: 2
    priority: 5
    timeout_minutes: 10

  # User defaults
  defaults:
    auto_skip_intro: false
    auto_skip_credits: false
    auto_skip_recap: false
    show_skip_button: true
    button_display_duration_seconds: 10

  # Cache
  cache:
    ttl_markers: 1h
    ttl_series_patterns: 24h
```



## API Endpoints

### Content Management
#### POST /api/v1/skip-intro/detect

Detect intro/credits for content

---
#### GET /api/v1/skip-intro/markers

Get markers for content

---
#### POST /api/v1/skip-intro/markers

Create manual marker

---
#### PUT /api/v1/skip-intro/markers/:id

Update marker

---
#### DELETE /api/v1/skip-intro/markers/:id

Delete marker

---
#### POST /api/v1/skip-intro/markers/:id/verify

Verify marker accuracy

---
#### GET /api/v1/skip-intro/preferences

Get user skip preferences

---
#### PUT /api/v1/skip-intro/preferences

Update user skip preferences

---
#### POST /api/v1/skip-intro/skip

Record skip action (analytics)

---
#### POST /api/v1/skip-intro/bulk-detect

Queue bulk detection for series

---
#### GET /api/v1/skip-intro/series/:id/patterns

Get series intro patterns

---


## Testing Strategy

### Unit Tests

<!-- Unit test strategy -->

### Integration Tests

<!-- Integration test strategy -->

### Test Coverage

Target: **80% minimum**







## Related Documentation
### Design Documents
- [01_ARCHITECTURE](../../architecture/01_ARCHITECTURE.md)
- [02_DESIGN_PRINCIPLES](../../architecture/02_DESIGN_PRINCIPLES.md)
- [03_METADATA_SYSTEM](../../architecture/03_METADATA_SYSTEM.md)

### External Sources
- [Chromaprint/AcoustID](../../../sources/standards/chromaprint.md) - Auto-resolved from chromaprint-acoustid
- [FFmpeg Documentation](../../../sources/media/ffmpeg.md) - Auto-resolved from ffmpeg
- [FFmpeg Codecs](../../../sources/media/ffmpeg-codecs.md) - Auto-resolved from ffmpeg-codecs
- [FFmpeg Formats](../../../sources/media/ffmpeg-formats.md) - Auto-resolved from ffmpeg-formats
- [go-astiav (FFmpeg bindings)](../../../sources/media/go-astiav.md) - Auto-resolved from go-astiav
- [go-astiav GitHub README](../../../sources/media/go-astiav-guide.md) - Auto-resolved from go-astiav-docs
- [River Job Queue](../../../sources/tooling/river.md) - Auto-resolved from river

