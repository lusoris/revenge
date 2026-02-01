## Table of Contents

- [Trickplay (Timeline Thumbnails)](#trickplay-timeline-thumbnails)
- [Trickplay (Timeline Thumbnails)](#trickplay-timeline-thumbnails)
  - [Contents](#contents)
  - [How It Works](#how-it-works)
    - [Content Flow](#content-flow)
  - [Features](#features)
  - [Configuration](#configuration)
  - [Related Documentation](#related-documentation)
    - [See Also](#see-also)

---
sources:
- name: Roku BIF Format
    url: ../sources/protocols/bif.md
    note: Auto-resolved from bif-spec
- name: FFmpeg Documentation
    url: ../sources/media/ffmpeg.md
    note: Auto-resolved from ffmpeg
- name: FFmpeg Codecs
    url: ../sources/media/ffmpeg-codecs.md
    note: Auto-resolved from ffmpeg-codecs
- name: FFmpeg Formats
    url: ../sources/media/ffmpeg-formats.md
    note: Auto-resolved from ffmpeg-formats
- name: go-astiav (FFmpeg bindings)
    url: ../sources/media/go-astiav.md
    note: Auto-resolved from go-astiav
- name: go-astiav GitHub README
    url: ../sources/media/go-astiav-guide.md
    note: Auto-resolved from go-astiav-docs
- name: Jellyfin Trickplay
    url: ../sources/apis/jellyfin-trickplay.md
    note: Auto-resolved from jellyfin-trickplay
- name: River Job Queue
    url: ../sources/tooling/river.md
    note: Auto-resolved from river
- name: WebVTT Specification
    url: ../sources/protocols/webvtt.md
    note: Auto-resolved from webvtt
design_refs:
- title: features/playback
    path: features/playback.md
- title: 01_ARCHITECTURE
    path: architecture/01_ARCHITECTURE.md
- title: 02_DESIGN_PRINCIPLES
    path: architecture/02_DESIGN_PRINCIPLES.md
- title: 03_METADATA_SYSTEM
    path: architecture/03_METADATA_SYSTEM.md
---

# Trickplay (Timeline Thumbnails)

# Trickplay (Timeline Thumbnails)

> Thumbnail previews on video seek bar

See where you are going when scrubbing through videos. Trickplay generates thumbnail previews that appear when hovering over the seek bar. Thumbnails are created during library scan using FFmpeg. Supports BIF format for Roku compatibility. Configurable interval (default 10 seconds) and quality settings per library.

---

## Contents

<!-- TOC will be auto-generated here by markdown-toc -->

---

## How It Works

<!-- Wiki how it works description pending -->

### Content Flow

1. **Add Content**: Import or add manually
2. **Metadata**: Automatically fetched
3. **Library**: Organized by folder structure
4. **Playback**: Stream directly or use external player

## Features
<!-- Feature list placeholder -->

## Configuration
<!-- User-friendly configuration guide -->

## Related Documentation
### See Also
<!-- Related wiki pages -->

---

**Need Help?** [Open an issue](https://github.com/revenge-project/revenge/issues) or [Join the discussion](https://github.com/revenge-project/revenge/discussions)
