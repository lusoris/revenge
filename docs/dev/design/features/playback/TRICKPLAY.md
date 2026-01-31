# Trickplay (Timeline Thumbnails)

> Thumbnail previews on video seek bar


<!-- TOC-START -->

## Table of Contents

- [Status](#status)
- [Developer Resources](#developer-resources)
- [Overview](#overview)
- [Features](#features)
- [Thumbnail Formats](#thumbnail-formats)
  - [Sprite Sheets (Recommended)](#sprite-sheets-recommended)
  - [BIF (Base Index Frames)](#bif-base-index-frames)
  - [WebVTT](#webvtt)
- [Architecture](#architecture)
- [Go Packages](#go-packages)
- [Database Schema](#database-schema)
- [River Jobs](#river-jobs)
- [Go Implementation](#go-implementation)
- [API Endpoints](#api-endpoints)
- [Client Integration](#client-integration)
  - [HTML5 Video Player](#html5-video-player)
  - [Video.js Plugin](#videojs-plugin)
- [Configuration](#configuration)
- [Priority Queue](#priority-queue)
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
**Priority**: 🟢 HIGH (Critical Gap - All competitors have this)
**Inspired By**: Jellyfin Trickplay, Plex Timeline Preview
**Location**: `internal/feature/trickplay/`

---

## Developer Resources

| Source             | URL                                                                                                                                                     | Purpose                           |
| ------------------ | ------------------------------------------------------------------------------------------------------------------------------------------------------- | --------------------------------- |
| FFmpeg             | [ffmpeg.org/documentation.html](https://ffmpeg.org/documentation.html)                                                                                  | Thumbnail extraction              |
| BIF Spec           | [sdkdocs.roku.com/display/sdkdoc/Trick+Mode+Support](https://developer.roku.com/docs/developer-program/media-playback/trick-mode/bif-file-creation.md) | Roku BIF format specification     |
| WebVTT             | [w3.org/TR/webvtt1](https://www.w3.org/TR/webvtt1/)                                                                                                     | Chapter/thumbnail metadata format |
| Jellyfin Trickplay | [jellyfin.org/docs/general/server/media/trickplay](https://jellyfin.org/docs/general/server/media/trickplay/)                                           | Reference implementation          |

---

## Overview

Trickplay generates thumbnail images at regular intervals throughout a video. When users hover over the seek bar, they see a preview of that point in the video, making navigation easier.

---

## Features

| Feature | Description |
|---------|-------------|
| Thumbnail Generation | Extract frames at configurable intervals |
| Sprite Sheets | Combine thumbnails into single images for efficiency |
| BIF Support | Roku BIF format compatibility |
| WebVTT Chapters | Standard chapter format with thumbnails |
| On-Demand Generation | Generate when needed, cache results |
| Priority Queue | Prioritize recently added/watched content |

---

## Thumbnail Formats

### Sprite Sheets (Recommended)

Single image containing grid of thumbnails:

```
┌─────┬─────┬─────┬─────┬─────┐
│ 0:00│ 0:10│ 0:20│ 0:30│ 0:40│
├─────┼─────┼─────┼─────┼─────┤
│ 0:50│ 1:00│ 1:10│ 1:20│ 1:30│
├─────┼─────┼─────┼─────┼─────┤
│ 1:40│ 1:50│ 2:00│ 2:10│ 2:20│
└─────┴─────┴─────┴─────┴─────┘

Sprite: 320x180 per thumbnail, 5 columns
Total: 1600x540 for 15 thumbnails
```

### BIF (Base Index Frames)

Roku format - single binary file with all thumbnails:

```
Header: 64 bytes (magic, version, count, interval)
Index: 8 bytes per frame (timestamp + offset)
Data: JPEG thumbnails concatenated
```

### WebVTT

Standard chapter/thumbnail format:

```vtt
WEBVTT

00:00:00.000 --> 00:00:10.000
thumbnails/sprite.jpg#xywh=0,0,320,180

00:00:10.000 --> 00:00:20.000
thumbnails/sprite.jpg#xywh=320,0,320,180
```

---

## Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                    Trickplay Pipeline                           │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  Video File ──► FFmpeg ──► Frames ──► Sprite Sheet ──► Storage │
│                                                                 │
│  ┌───────────┐    ┌───────────┐    ┌───────────┐              │
│  │  River    │───►│  FFmpeg   │───►│  govips   │              │
│  │  Job      │    │  Extract  │    │  Sprite   │              │
│  └───────────┘    └───────────┘    └───────────┘              │
│        │                                   │                    │
│        │         ┌─────────────────────────┘                   │
│        │         ▼                                              │
│        │    ┌───────────┐    ┌───────────┐                     │
│        │    │  WebVTT   │    │ Filesystem│                     │
│        └───►│  Generate │───►│  / S3     │                     │
│             └───────────┘    └───────────┘                     │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

---

## Go Packages

> Package versions: [00_SOURCE_OF_TRUTH.md](../../00_SOURCE_OF_TRUTH.md#go-dependencies-core)

| Package | Purpose |
|---------|---------|
| go-astiav | FFmpeg bindings |
| govips | Image processing (sprite sheets) |
| imaging | Pure Go imaging (fallback) |

---

## Database Schema

```sql
CREATE TABLE trickplay_data (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    content_type VARCHAR(50) NOT NULL, -- movie, tvshow_episode
    content_id UUID NOT NULL,

    -- Generation settings
    interval_seconds INT NOT NULL DEFAULT 10,
    thumbnail_width INT NOT NULL DEFAULT 320,
    thumbnail_height INT NOT NULL DEFAULT 180,

    -- Output
    format VARCHAR(20) NOT NULL, -- sprite, bif, individual
    sprite_columns INT,
    total_thumbnails INT NOT NULL,

    -- Files
    file_path TEXT NOT NULL, -- Path to sprite/BIF file
    webvtt_path TEXT, -- Path to WebVTT file
    file_size_bytes BIGINT,

    -- Status
    status VARCHAR(20) DEFAULT 'pending', -- pending, generating, complete, failed
    error_message TEXT,
    progress_percent INT DEFAULT 0,

    -- Timestamps
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    completed_at TIMESTAMPTZ,

    UNIQUE(content_type, content_id)
);

CREATE INDEX idx_trickplay_content ON trickplay_data(content_type, content_id);
CREATE INDEX idx_trickplay_status ON trickplay_data(status);
```

---

## River Jobs

```go
const (
    JobKindGenerateTrickplay     = "trickplay.generate"
    JobKindGenerateTrickplayBulk = "trickplay.generate_bulk"
    JobKindCleanupOrphaned       = "trickplay.cleanup_orphaned"
)

type GenerateTrickplayArgs struct {
    ContentType string    `json:"content_type"`
    ContentID   uuid.UUID `json:"content_id"`
    VideoPath   string    `json:"video_path"`
    Interval    int       `json:"interval"` // seconds
    Width       int       `json:"width"`
    Height      int       `json:"height"`
    Format      string    `json:"format"` // sprite, bif
}

func (GenerateTrickplayArgs) Kind() string {
    return JobKindGenerateTrickplay
}
```

---

## Go Implementation

```go
// internal/service/trickplay/

type Service struct {
    repo   TrickplayRepository
    river  *river.Client[pgx.Tx]
    config *Config
}

type Generator struct {
    ffmpegPath string
    outputDir  string
}

func (g *Generator) ExtractFrames(ctx context.Context, videoPath string, interval int) ([]string, error) {
    // Use ffmpeg to extract frames at interval
    outputPattern := filepath.Join(g.outputDir, "frame_%04d.jpg")

    err := ffmpeg.Input(videoPath).
        Filter("fps", ffmpeg.Args{fmt.Sprintf("1/%d", interval)}).
        Filter("scale", ffmpeg.Args{fmt.Sprintf("%d:%d", g.config.Width, g.config.Height)}).
        Output(outputPattern, ffmpeg.KwArgs{"q:v": "5"}).
        OverWriteOutput().
        Run()

    if err != nil {
        return nil, fmt.Errorf("ffmpeg extract: %w", err)
    }

    // Return list of generated frame paths
    return filepath.Glob(filepath.Join(g.outputDir, "frame_*.jpg"))
}

func (g *Generator) CreateSpriteSheet(frames []string, columns int) (string, error) {
    // Use govips to combine frames into sprite sheet
    // Each row has `columns` thumbnails

    images := make([]*vips.ImageRef, len(frames))
    for i, frame := range frames {
        img, err := vips.NewImageFromFile(frame)
        if err != nil {
            return "", err
        }
        images[i] = img
    }

    // Calculate dimensions
    rows := (len(frames) + columns - 1) / columns
    width := g.config.Width * columns
    height := g.config.Height * rows

    // Create sprite sheet (using imaging library for composition)
    // ... composition logic ...

    spritePath := filepath.Join(g.outputDir, "sprite.jpg")
    return spritePath, nil
}

func (g *Generator) GenerateWebVTT(frames []string, interval int, spritePath string, columns int) (string, error) {
    var buf bytes.Buffer
    buf.WriteString("WEBVTT\n\n")

    for i := range frames {
        start := i * interval
        end := (i + 1) * interval

        // Calculate position in sprite
        col := i % columns
        row := i / columns
        x := col * g.config.Width
        y := row * g.config.Height

        fmt.Fprintf(&buf, "%s --> %s\n", formatTime(start), formatTime(end))
        fmt.Fprintf(&buf, "%s#xywh=%d,%d,%d,%d\n\n",
            spritePath, x, y, g.config.Width, g.config.Height)
    }

    vttPath := filepath.Join(g.outputDir, "thumbnails.vtt")
    return vttPath, os.WriteFile(vttPath, buf.Bytes(), 0644)
}

func formatTime(seconds int) string {
    h := seconds / 3600
    m := (seconds % 3600) / 60
    s := seconds % 60
    return fmt.Sprintf("%02d:%02d:%02d.000", h, m, s)
}
```

---

## API Endpoints

```
# Get trickplay data for content
GET /api/v1/trickplay/:content_type/:content_id

# Get sprite sheet image
GET /api/v1/trickplay/:content_type/:content_id/sprite.jpg

# Get WebVTT file
GET /api/v1/trickplay/:content_type/:content_id/thumbnails.vtt

# Get BIF file (for Roku)
GET /api/v1/trickplay/:content_type/:content_id/index.bif

# Trigger generation (admin)
POST /api/v1/trickplay/:content_type/:content_id/generate

# Bulk generation (admin)
POST /api/v1/trickplay/generate-bulk
```

---

## Client Integration

### HTML5 Video Player

```typescript
// Fetch WebVTT on video load
const trickplayUrl = `/api/v1/trickplay/movie/${movieId}/thumbnails.vtt`;

// Parse WebVTT and show thumbnails on hover
player.on('timeupdate', (position) => {
    const thumbnail = getThumbnailForPosition(position);
    showThumbnailPreview(thumbnail);
});
```

### Video.js Plugin

```javascript
videojs.registerPlugin('trickplay', function(options) {
    const player = this;
    // Load and display trickplay thumbnails
});
```

---

## Configuration

```yaml
trickplay:
  enabled: true
  auto_generate: true  # Generate on library scan
  interval_seconds: 10
  thumbnail:
    width: 320
    height: 180
  sprite:
    columns: 10  # 10 thumbnails per row
  formats:
    - sprite
    - webvtt
    # - bif  # Uncomment for Roku support
  storage:
    path: "/data/trickplay"
    # s3_bucket: "revenge-trickplay"  # Optional S3 storage
  priority:
    new_content: true
    recently_watched: true
```

---

## Priority Queue

Content is prioritized for trickplay generation:

1. **Currently watching** - Generate immediately
2. **Recently added** - Within 24 hours
3. **Popular content** - High watch count
4. **Remainder** - Background generation

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
| [FFmpeg Codecs](https://ffmpeg.org/ffmpeg-codecs.html) | [Local](../../../sources/media/ffmpeg-codecs.md) |
| [FFmpeg Documentation](https://ffmpeg.org/ffmpeg.html) | [Local](../../../sources/media/ffmpeg.md) |
| [FFmpeg Formats](https://ffmpeg.org/ffmpeg-formats.html) | [Local](../../../sources/media/ffmpeg-formats.md) |
| [Jellyfin Trickplay](https://jellyfin.org/docs/general/server/media/trickplay/) | [Local](../../../sources/apis/jellyfin-trickplay.md) |
| [River Job Queue](https://pkg.go.dev/github.com/riverqueue/river) | [Local](../../../sources/tooling/river.md) |
| [Roku BIF Format](https://developer.roku.com/docs/developer-program/media-playback/trick-mode/bif-file-creation.md) | [Local](../../../sources/protocols/bif.md) |
| [WebVTT Specification](https://www.w3.org/TR/webvtt1/) | [Local](../../../sources/protocols/webvtt.md) |
| [go-astiav (FFmpeg bindings)](https://pkg.go.dev/github.com/asticode/go-astiav) | [Local](../../../sources/media/go-astiav.md) |

<!-- SOURCE-BREADCRUMBS-END -->

<!-- DESIGN-BREADCRUMBS-START -->

## Related Design Docs

> Auto-generated cross-references to related design documentation

**Category**: [Playback](INDEX.md)

### In This Section

- [Revenge - Media Enhancement Features](MEDIA_ENHANCEMENTS.md)
- [Release Calendar System](RELEASE_CALENDAR.md)
- [Skip Intro / Credits Detection](SKIP_INTRO.md)
- [SyncPlay (Watch Together)](SYNCPLAY.md)
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

- [Media Enhancements](MEDIA_ENHANCEMENTS.md)
- [River Job Queue Patterns](../../00_SOURCE_OF_TRUTH.md#river-job-queue-patterns)
- [Go Packages](../architecture/GO_PACKAGES.md)
