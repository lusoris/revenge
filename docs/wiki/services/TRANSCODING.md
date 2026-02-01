## Table of Contents

- [Transcoding Service](#transcoding-service)
- [Transcoding Service](#transcoding-service)
  - [Contents](#contents)
  - [How It Works](#how-it-works)
  - [Features](#features)
  - [Configuration](#configuration)
  - [Related Documentation](#related-documentation)
    - [See Also](#see-also)

---
sources:
- name: go-astiav (FFmpeg)
    url: ../sources/media/go-astiav.md
    note: FFmpeg Go bindings
- name: gohlslib
    url: ../sources/media/gohlslib.md
    note: HLS streaming library
- name: River Job Queue
    url: ../sources/tooling/river.md
    note: Background job processing
- name: Uber fx
    url: ../sources/tooling/fx.md
    note: Dependency injection
design_refs:
- title: services
    path: services.md
- title: 01_ARCHITECTURE
    path: architecture/01_ARCHITECTURE.md
- title: OFFLOADING
    path: technical/OFFLOADING.md
- title: AUDIO_STREAMING
    path: technical/AUDIO_STREAMING.md
---

# Transcoding Service

# Transcoding Service

> High-performance media transcoding with hardware acceleration

The Transcoding Service converts media files on-demand to ensure compatibility across all devices. Primary transcoding offloads to Blackbeard service for maximum performance, with local FFmpeg fallback. Hardware acceleration via NVENC, QSV, or VAAPI dramatically reduces CPU usage. Generates HLS adaptive streams with multiple quality levels, caching results for faster subsequent playback.

---

## Contents

<!-- TOC will be auto-generated here by markdown-toc -->

---

## How It Works

<!-- How it works -->

## Features
<!-- Feature list placeholder -->

## Configuration
<!-- User-friendly configuration guide -->

## Related Documentation
### See Also
<!-- Related wiki pages -->

---

**Need Help?** [Open an issue](https://github.com/revenge-project/revenge/issues) or [Join the discussion](https://github.com/revenge-project/revenge/discussions)
