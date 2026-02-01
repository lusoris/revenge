## Table of Contents

- [Revenge - Player Architecture](#revenge-player-architecture)
  - [Contents](#contents)
  - [How It Works](#how-it-works)
  - [Features](#features)
  - [Configuration](#configuration)
  - [Related Documentation](#related-documentation)
    - [See Also](#see-also)



---
sources:
  - name: FFmpeg Documentation
    url: ../../sources/media/ffmpeg.md
    note: Auto-resolved from ffmpeg
  - name: FFmpeg Codecs
    url: ../../sources/media/ffmpeg-codecs.md
    note: Auto-resolved from ffmpeg-codecs
  - name: FFmpeg Formats
    url: ../../sources/media/ffmpeg-formats.md
    note: Auto-resolved from ffmpeg-formats
  - name: go-astiav (FFmpeg bindings)
    url: ../../sources/media/go-astiav.md
    note: Auto-resolved from go-astiav
  - name: go-astiav GitHub README
    url: ../../sources/media/go-astiav-guide.md
    note: Auto-resolved from go-astiav-docs
  - name: gohlslib (HLS)
    url: ../../sources/media/gohlslib.md
    note: Auto-resolved from gohlslib
  - name: M3U8 Extended Format
    url: ../../sources/protocols/m3u8.md
    note: Auto-resolved from m3u8
  - name: Svelte 5 Runes
    url: ../../sources/frontend/svelte-runes.md
    note: Auto-resolved from svelte-runes
  - name: Svelte 5 Documentation
    url: ../../sources/frontend/svelte5.md
    note: Auto-resolved from svelte5
  - name: SvelteKit Documentation
    url: ../../sources/frontend/sveltekit.md
    note: Auto-resolved from sveltekit
design_refs:
  - title: architecture
    path: INDEX.md
  - title: ADULT_CONTENT_SYSTEM
    path: ../features/adult/ADULT_CONTENT_SYSTEM.md
  - title: ADULT_METADATA
    path: ../features/adult/ADULT_METADATA.md
  - title: DATA_RECONCILIATION
    path: ../features/adult/DATA_RECONCILIATION.md
---

# Revenge - Player Architecture




> How Revenge plays your media across all devices

The player architecture handles media playback across web, mobile, and TV clients. Videos stream via HLS for adaptive quality based on your connection. The backend generates HLS manifests and handles transcoding when needed (using FFmpeg). The web player (Vidstack) supports chapters, subtitles, skip intro, and trickplay thumbnails. Chromecast lets you cast to TV, and SyncPlay lets multiple users watch together in sync.

---




## Contents

<!-- TOC will be auto-generated here by markdown-toc -->

---


## How It Works

<!-- User-friendly explanation -->




## Features
<!-- Feature list placeholder -->



## Configuration
<!-- User-friendly configuration guide -->









## Related Documentation
### See Also
<!-- Related wiki pages -->



---

**Need Help?** [Open an issue](https://github.com/revenge-project/revenge/issues) or [Join the discussion](https://github.com/revenge-project/revenge/discussions)