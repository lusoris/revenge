## Table of Contents

- [Revenge - Design Principles](#revenge-design-principles)
  - [Contents](#contents)
  - [How It Works](#how-it-works)
  - [Features](#features)
  - [Configuration](#configuration)
  - [Related Documentation](#related-documentation)
    - [See Also](#see-also)



---
sources:
  - name: Dragonfly Documentation
    url: ../sources/infrastructure/dragonfly.md
    note: Auto-resolved from dragonfly
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
  - name: gohlslib (HLS)
    url: ../sources/media/gohlslib.md
    note: Auto-resolved from gohlslib
  - name: M3U8 Extended Format
    url: ../sources/protocols/m3u8.md
    note: Auto-resolved from m3u8
  - name: River Job Queue
    url: ../sources/tooling/river.md
    note: Auto-resolved from river
  - name: shadcn-svelte
    url: ../sources/frontend/shadcn-svelte.md
    note: Auto-resolved from shadcn-svelte
  - name: Svelte 5 Runes
    url: ../sources/frontend/svelte-runes.md
    note: Auto-resolved from svelte-runes
  - name: Svelte 5 Documentation
    url: ../sources/frontend/svelte5.md
    note: Auto-resolved from svelte5
  - name: SvelteKit Documentation
    url: ../sources/frontend/sveltekit.md
    note: Auto-resolved from sveltekit
  - name: TanStack Query
    url: ../sources/frontend/tanstack-query.md
    note: Auto-resolved from tanstack-query
  - name: Typesense API
    url: ../sources/infrastructure/typesense.md
    note: Auto-resolved from typesense
  - name: Typesense Go Client
    url: ../sources/infrastructure/typesense-go.md
    note: Auto-resolved from typesense-go
design_refs:
  - title: architecture
    path: architecture/INDEX.md
  - title: ADULT_CONTENT_SYSTEM
    path: ADULT_CONTENT_SYSTEM.md
  - title: ADULT_METADATA
    path: ADULT_METADATA.md
  - title: DATA_RECONCILIATION
    path: DATA_RECONCILIATION.md
---

# Revenge - Design Principles




> The philosophy behind how Revenge is built

Revenge follows strict design principles for maintainability and performance. PostgreSQL is the only supported database - no SQLite complexity. All code requires 80% minimum test coverage with table-driven tests. Metadata always comes from local cache first, then Arr services, then external APIs. Error handling uses sentinel errors for type safety. Logging is text in development (tint) and JSON in production (zap). These principles ensure consistent, testable, performant code.

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