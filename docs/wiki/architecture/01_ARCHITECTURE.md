## Table of Contents

- [Revenge - Architecture v2](#revenge-architecture-v2)
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
  - name: Uber fx
    url: ../sources/tooling/fx.md
    note: Auto-resolved from fx
  - name: go-astiav (FFmpeg bindings)
    url: ../sources/media/go-astiav.md
    note: Auto-resolved from go-astiav
  - name: go-astiav GitHub README
    url: ../sources/media/go-astiav-guide.md
    note: Auto-resolved from go-astiav-docs
  - name: go-blurhash
    url: ../sources/media/go-blurhash.md
    note: Auto-resolved from go-blurhash
  - name: gohlslib (HLS)
    url: ../sources/media/gohlslib.md
    note: Auto-resolved from gohlslib
  - name: koanf
    url: ../sources/tooling/koanf.md
    note: Auto-resolved from koanf
  - name: Last.fm API
    url: ../sources/apis/lastfm.md
    note: Auto-resolved from lastfm-api
  - name: M3U8 Extended Format
    url: ../sources/protocols/m3u8.md
    note: Auto-resolved from m3u8
  - name: ogen OpenAPI Generator
    url: ../sources/tooling/ogen.md
    note: Auto-resolved from ogen
  - name: pgx PostgreSQL Driver
    url: ../sources/database/pgx.md
    note: Auto-resolved from pgx
  - name: PostgreSQL Arrays
    url: ../sources/database/postgresql-arrays.md
    note: Auto-resolved from postgresql-arrays
  - name: PostgreSQL JSON Functions
    url: ../sources/database/postgresql-json.md
    note: Auto-resolved from postgresql-json
  - name: River Job Queue
    url: ../sources/tooling/river.md
    note: Auto-resolved from river
  - name: rueidis
    url: ../sources/tooling/rueidis.md
    note: Auto-resolved from rueidis
  - name: rueidis GitHub README
    url: ../sources/tooling/rueidis-guide.md
    note: Auto-resolved from rueidis-docs
  - name: shadcn-svelte
    url: ../sources/frontend/shadcn-svelte.md
    note: Auto-resolved from shadcn-svelte
  - name: sqlc
    url: ../sources/database/sqlc.md
    note: Auto-resolved from sqlc
  - name: sqlc Configuration
    url: ../sources/database/sqlc-config.md
    note: Auto-resolved from sqlc-config
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

# Revenge - Architecture v2




> How Revenge is built - the technical foundation

Revenge is built with a Go backend and SvelteKit frontend. PostgreSQL stores all data (no SQLite), with Dragonfly providing fast caching and Typesense powering search. Background jobs run through River queue. The backend uses fx for dependency injection and ogen for type-safe API generation. The frontend uses Svelte 5 with shadcn-svelte components. All components are designed for self-hosting with Docker or bare metal deployment.

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