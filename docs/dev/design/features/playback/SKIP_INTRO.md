# Skip Intro / Credits Detection

> Automatic intro and credits detection with one-click skip


<!-- TOC-START -->

## Table of Contents

- [Status](#status)
- [Developer Resources](#developer-resources)
- [Overview](#overview)
- [Features](#features)
- [Detection Methods](#detection-methods)
  - [1. Audio Fingerprinting (Primary)](#1-audio-fingerprinting-primary)
  - [2. Silence Detection](#2-silence-detection)
  - [3. Black Frame Detection](#3-black-frame-detection)
  - [4. Template Matching (Credits)](#4-template-matching-credits)
- [Architecture](#architecture)
- [Go Packages](#go-packages)
- [Database Schema](#database-schema)
- [River Jobs](#river-jobs)
- [Go Implementation](#go-implementation)
- [API Endpoints](#api-endpoints)
- [Client Integration](#client-integration)
  - [Skip Button UI](#skip-button-ui)
  - [JavaScript Implementation](#javascript-implementation)
- [Configuration](#configuration)
- [RBAC Permissions](#rbac-permissions)
- [Sources & Cross-References](#sources-cross-references)
  - [Cross-Reference Indexes](#cross-reference-indexes)
  - [Referenced Sources](#referenced-sources)
- [Related Design Docs](#related-design-docs)
  - [In This Section](#in-this-section)
  - [Related Topics](#related-topics)
  - [Indexes](#indexes)
- [Related Documentation](#related-documentation)

<!-- TOC-END -->

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
**Priority**: 🟢 HIGH (Critical Gap - Plex/Jellyfin have this)
**Inspired By**: Plex Skip Intro, Jellyfin Intro Skipper plugin
**Location**: `internal/feature/skipintro/`

---

## Developer Resources

| Source                 | URL                                                                                                     | Purpose                       |
| ---------------------- | ------------------------------------------------------------------------------------------------------- | ----------------------------- |
| Chromaprint            | [acoustid.org/chromaprint](https://acoustid.org/chromaprint)                                            | Audio fingerprinting library  |
| FFmpeg Filters         | [ffmpeg.org/ffmpeg-filters.html](https://ffmpeg.org/ffmpeg-filters.html)                                | Silence/black frame detection |
| Jellyfin Intro Skipper | [github.com/ConfusedPolarBear/intro-skipper](https://github.com/ConfusedPolarBear/intro-skipper)        | Reference implementation      |

---

## Overview

Automatically detect intro sequences and credits in video content, allowing users to skip them with a single click. Supports TV show intros, movie credits, and recap sequences.

---

## Features

| Feature | Description |
|---------|-------------|
| Intro Detection | Detect opening sequences in TV episodes |
| Credits Detection | Detect end credits |
| Recap Detection | Detect "Previously on..." sequences |
| Skip Button | One-click skip during playback |
| Auto-Skip | Optional automatic skipping |
| Chapter Markers | Add detected segments as chapters |

---

## Detection Methods

### 1. Audio Fingerprinting (Primary)

Compare audio signatures across episodes to find common intro music:

```
Episode 1: [Intro: 0:00-1:30] [Content] [Credits: 42:00-43:00]
Episode 2: [Intro: 0:00-1:30] [Content] [Credits: 41:30-42:30]
Episode 3: [Intro: 0:00-1:30] [Content] [Credits: 42:15-43:15]
              ↑ Same audio     matches across episodes
```

### 2. Silence Detection

Detect silence patterns that typically surround intros:

```
[Pre-intro content] [SILENCE] [INTRO MUSIC] [SILENCE] [Main content]
```

### 3. Black Frame Detection

Detect black frames that often bookend intros:

```
[Content] → [BLACK FRAMES] → [Intro] → [BLACK FRAMES] → [Content]
```

### 4. Template Matching (Credits)

Detect credit roll patterns:
- Text scrolling upward
- Dark background with light text
- Consistent pace/speed

---

## Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                  Intro/Credits Detection Pipeline               │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  ┌─────────────┐     ┌─────────────┐     ┌─────────────┐      │
│  │   Audio     │     │   Video     │     │   Silence   │      │
│  │ Fingerprint │     │   Analysis  │     │  Detection  │      │
│  └──────┬──────┘     └──────┬──────┘     └──────┬──────┘      │
│         │                   │                   │              │
│         └───────────────────┴───────────────────┘              │
│                             │                                   │
│                    ┌────────▼────────┐                         │
│                    │    Combiner     │                         │
│                    │  (Confidence)   │                         │
│                    └────────┬────────┘                         │
│                             │                                   │
│                    ┌────────▼────────┐                         │
│                    │   Segment DB    │                         │
│                    └─────────────────┘                         │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

---

## Go Packages

> Package versions: [00_SOURCE_OF_TRUTH.md](../../00_SOURCE_OF_TRUTH.md#go-dependencies-core)

| Package | Purpose |
|---------|---------|
| go-astiav | Audio extraction, frame analysis (FFmpeg bindings) |
| chromaprint | Audio fingerprinting (via FFmpeg) |
| imaging | Frame analysis |

---

## Database Schema

```sql
CREATE TYPE segment_type AS ENUM ('intro', 'credits', 'recap', 'preview');

CREATE TABLE media_segments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    content_type VARCHAR(50) NOT NULL,
    content_id UUID NOT NULL,

    -- Segment info
    segment_type segment_type NOT NULL,
    start_ms BIGINT NOT NULL,
    end_ms BIGINT NOT NULL,

    -- Detection metadata
    detection_method VARCHAR(50), -- audio_fingerprint, silence, black_frame, manual
    confidence DECIMAL(5,4), -- 0.0000 to 1.0000
    fingerprint_hash TEXT, -- Audio fingerprint for matching

    -- User override
    is_verified BOOLEAN DEFAULT false,
    verified_by UUID REFERENCES users(id),

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE(content_type, content_id, segment_type)
);

-- Audio fingerprints for cross-episode matching
CREATE TABLE audio_fingerprints (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    series_id UUID NOT NULL,
    season_number INT,

    -- Fingerprint data
    fingerprint_hash TEXT NOT NULL,
    duration_ms INT NOT NULL,

    -- Where this fingerprint appears
    occurrences INT DEFAULT 1,
    first_seen_content_id UUID,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- User preferences
CREATE TABLE user_skip_preferences (
    user_id UUID PRIMARY KEY REFERENCES users(id),
    auto_skip_intro BOOLEAN DEFAULT false,
    auto_skip_credits BOOLEAN DEFAULT false,
    auto_skip_recap BOOLEAN DEFAULT false,
    skip_button_duration_seconds INT DEFAULT 10
);

CREATE INDEX idx_media_segments_content ON media_segments(content_type, content_id);
CREATE INDEX idx_audio_fingerprints_series ON audio_fingerprints(series_id);
CREATE INDEX idx_audio_fingerprints_hash ON audio_fingerprints(fingerprint_hash);
```

---

## River Jobs

```go
const (
    JobKindDetectIntro       = "skip.detect_intro"
    JobKindDetectCredits     = "skip.detect_credits"
    JobKindDetectSeriesIntro = "skip.detect_series_intro"
    JobKindFingerprintAudio  = "skip.fingerprint_audio"
)

type DetectIntroArgs struct {
    ContentType string    `json:"content_type"`
    ContentID   uuid.UUID `json:"content_id"`
    VideoPath   string    `json:"video_path"`
    SeriesID    uuid.UUID `json:"series_id,omitempty"` // For cross-episode matching
}

type DetectSeriesIntroArgs struct {
    SeriesID uuid.UUID `json:"series_id"`
    SeasonNumber int   `json:"season_number,omitempty"` // Optional: specific season
}
```

---

## Go Implementation

```go
// internal/service/skipdetect/

type Service struct {
    repo   SegmentRepository
    finger AudioFingerprinter
    river  *river.Client[pgx.Tx]
}

type AudioFingerprinter struct {
    ffmpegPath string
}

// Extract audio fingerprint using FFmpeg with chromaprint
func (f *AudioFingerprinter) Fingerprint(ctx context.Context, videoPath string, startSec, durationSec int) (string, error) {
    // Extract audio segment
    cmd := exec.CommandContext(ctx, f.ffmpegPath,
        "-i", videoPath,
        "-ss", fmt.Sprintf("%d", startSec),
        "-t", fmt.Sprintf("%d", durationSec),
        "-ac", "1",
        "-ar", "22050",
        "-f", "chromaprint",
        "-fp_format", "compressed",
        "-",
    )

    output, err := cmd.Output()
    if err != nil {
        return "", fmt.Errorf("chromaprint: %w", err)
    }

    return string(output), nil
}

// Detect silence in audio track
func (s *Service) DetectSilence(ctx context.Context, videoPath string) ([]SilenceSegment, error) {
    // Use FFmpeg silencedetect filter
    cmd := exec.CommandContext(ctx, "ffmpeg",
        "-i", videoPath,
        "-af", "silencedetect=noise=-50dB:d=0.5",
        "-f", "null", "-",
    )

    // Parse stderr for silence_start and silence_end
    stderr, _ := cmd.StderrPipe()
    cmd.Start()

    var segments []SilenceSegment
    scanner := bufio.NewScanner(stderr)
    for scanner.Scan() {
        line := scanner.Text()
        // Parse: [silencedetect @ ...] silence_start: 0.000
        // Parse: [silencedetect @ ...] silence_end: 1.234
        segments = append(segments, parseSilenceLine(line)...)
    }

    cmd.Wait()
    return segments, nil
}

// Find common intro across series episodes
func (s *Service) FindSeriesIntro(ctx context.Context, seriesID uuid.UUID) (*IntroSegment, error) {
    // Get fingerprints for first 3 minutes of each episode
    episodes, _ := s.repo.GetSeriesEpisodes(ctx, seriesID)

    fingerprints := make(map[string][]uuid.UUID) // hash -> episode IDs

    for _, ep := range episodes {
        fp, _ := s.finger.Fingerprint(ctx, ep.VideoPath, 0, 180) // First 3 minutes
        fingerprints[fp] = append(fingerprints[fp], ep.ID)
    }

    // Find most common fingerprint (appears in >50% of episodes)
    threshold := len(episodes) / 2
    for hash, eps := range fingerprints {
        if len(eps) >= threshold {
            // Found common intro
            return &IntroSegment{
                Hash:       hash,
                Confidence: float64(len(eps)) / float64(len(episodes)),
            }, nil
        }
    }

    return nil, ErrNoCommonIntro
}
```

---

## API Endpoints

```
# Get segments for content
GET /api/v1/segments/:content_type/:content_id

# Trigger detection (admin)
POST /api/v1/segments/:content_type/:content_id/detect

# Detect for entire series
POST /api/v1/segments/series/:series_id/detect

# Manual segment (admin/user with permission)
POST /api/v1/segments/:content_type/:content_id
PUT  /api/v1/segments/:id
DELETE /api/v1/segments/:id

# User preferences
GET  /api/v1/users/me/skip-preferences
PUT  /api/v1/users/me/skip-preferences
```

---

## Client Integration

### Skip Button UI

```
┌─────────────────────────────────────────┐
│                                         │
│   ┌─────────────────┐                   │
│   │  SKIP INTRO ►►  │ ← Appears at      │
│   │    (8 sec)      │   intro start     │
│   └─────────────────┘                   │
│                                         │
│  ═══════════════════════════════════   │
│  0:15 ▶──────●────────────────── 43:00 │
└─────────────────────────────────────────┘
```

### JavaScript Implementation

```typescript
interface Segment {
    type: 'intro' | 'credits' | 'recap';
    startMs: number;
    endMs: number;
}

function handleSkipButton(player: VideoPlayer, segments: Segment[]) {
    player.on('timeupdate', (currentTime) => {
        const intro = segments.find(s => s.type === 'intro');
        if (intro && currentTime >= intro.startMs && currentTime < intro.endMs) {
            showSkipButton('Skip Intro', intro.endMs - currentTime);
        }

        const credits = segments.find(s => s.type === 'credits');
        if (credits && currentTime >= credits.startMs) {
            showSkipButton('Skip Credits', credits.endMs - currentTime);
            // Optionally: show "Next Episode" button
        }
    });
}
```

---

## Configuration

```yaml
skip_detection:
  enabled: true
  auto_detect_on_scan: true

  intro:
    enabled: true
    max_duration_seconds: 180  # Max 3 minutes
    min_confidence: 0.75
    methods:
      - audio_fingerprint
      - silence

  credits:
    enabled: true
    min_confidence: 0.70
    methods:
      - black_frame
      - template_match

  recap:
    enabled: true
    max_duration_seconds: 60

  user_defaults:
    auto_skip_intro: false
    auto_skip_credits: false
    skip_button_duration: 10  # seconds
```

---

## RBAC Permissions

| Permission | Description |
|------------|-------------|
| `segments.view` | View detected segments |
| `segments.edit` | Manually add/edit segments |
| `segments.delete` | Delete segments |
| `segments.detect` | Trigger detection |

---


<!-- SOURCE-BREADCRUMBS-START -->

## Sources & Cross-References

> Auto-generated section linking to external documentation sources

### Cross-Reference Indexes

- [All Sources Index](../../../sources/SOURCES_INDEX.md) - Complete list of external documentation
- [Design ↔ Sources Map](../../../sources/DESIGN_CROSSREF.md) - Which docs reference which sources

### Referenced Sources

| Source | Documentation |
|--------|---------------|
| [Chromaprint/AcoustID](https://acoustid.org/chromaprint) | [Local](../../../sources/standards/chromaprint.md) |
| [FFmpeg Codecs](https://ffmpeg.org/ffmpeg-codecs.html) | [Local](../../../sources/media/ffmpeg-codecs.md) |
| [FFmpeg Documentation](https://ffmpeg.org/ffmpeg.html) | [Local](../../../sources/media/ffmpeg.md) |
| [FFmpeg Formats](https://ffmpeg.org/ffmpeg-formats.html) | [Local](../../../sources/media/ffmpeg-formats.md) |
| [River Job Queue](https://pkg.go.dev/github.com/riverqueue/river) | [Local](../../../sources/tooling/river.md) |
| [go-astiav (FFmpeg bindings)](https://pkg.go.dev/github.com/asticode/go-astiav) | [Local](../../../sources/media/go-astiav.md) |

<!-- SOURCE-BREADCRUMBS-END -->

<!-- DESIGN-BREADCRUMBS-START -->

## Related Design Docs

> Auto-generated cross-references to related design documentation

**Category**: [Playback](INDEX.md)

### In This Section

- [Revenge - Media Enhancement Features](MEDIA_ENHANCEMENTS.md)
- [Release Calendar System](RELEASE_CALENDAR.md)
- [SyncPlay (Watch Together)](SYNCPLAY.md)
- [Trickplay (Timeline Thumbnails)](TRICKPLAY.md)
- [Watch Next & Continue Watching System](WATCH_NEXT_CONTINUE_WATCHING.md)

### Related Topics

- [Revenge - Architecture v2](../../architecture/01_ARCHITECTURE.md) _Architecture_
- [Revenge - Design Principles](../../architecture/02_DESIGN_PRINCIPLES.md) _Architecture_
- [Revenge - Metadata System](../../architecture/03_METADATA_SYSTEM.md) _Architecture_
- [Revenge - Player Architecture](../../architecture/04_PLAYER_ARCHITECTURE.md) _Architecture_
- [Plugin Architecture Decision](../../architecture/05_PLUGIN_ARCHITECTURE_DECISION.md) _Architecture_

### Indexes

- [Design Index](../../DESIGN_INDEX.md) - All design docs by category/topic
- [Source of Truth](../../00_SOURCE_OF_TRUTH.md) - Package versions and status

<!-- DESIGN-BREADCRUMBS-END -->

## Related Documentation

- [Trickplay](TRICKPLAY.md)
- [Media Enhancements](MEDIA_ENHANCEMENTS.md)
- [River Job Queue Patterns](../../00_SOURCE_OF_TRUTH.md#river-job-queue-patterns)
