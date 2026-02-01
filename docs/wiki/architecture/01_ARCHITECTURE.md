## Table of Contents

- [Revenge - Architecture v2](#revenge-architecture-v2)
  - [Contents](#contents)
  - [How It Works](#how-it-works)
  - [Features](#features)
  - [Configuration](#configuration)
  - [Related Documentation](#related-documentation)
    - [Related Pages](#related-pages)
    - [Learn More](#learn-more)



---
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
### Related Pages
<!-- Related wiki pages -->

### Learn More

Official documentation and guides:
- [Dragonfly Documentation](../../sources/infrastructure/dragonfly.md)
- [FFmpeg Documentation](../../sources/media/ffmpeg.md)
- [FFmpeg Codecs](../../sources/media/ffmpeg-codecs.md)
- [FFmpeg Formats](../../sources/media/ffmpeg-formats.md)
- [Uber fx](../../sources/tooling/fx.md)
- [go-astiav (FFmpeg bindings)](../../sources/media/go-astiav.md)
- [go-astiav GitHub README](../../sources/media/go-astiav-guide.md)
- [go-blurhash](../../sources/media/go-blurhash.md)
- [gohlslib (HLS)](../../sources/media/gohlslib.md)
- [koanf](../../sources/tooling/koanf.md)
- [Last.fm API](../../sources/apis/lastfm.md)
- [M3U8 Extended Format](../../sources/protocols/m3u8.md)
- [ogen OpenAPI Generator](../../sources/tooling/ogen.md)
- [pgx PostgreSQL Driver](../../sources/database/pgx.md)
- [PostgreSQL Arrays](../../sources/database/postgresql-arrays.md)
- [PostgreSQL JSON Functions](../../sources/database/postgresql-json.md)
- [River Job Queue](../../sources/tooling/river.md)
- [rueidis](../../sources/tooling/rueidis.md)
- [rueidis GitHub README](../../sources/tooling/rueidis-guide.md)
- [shadcn-svelte](../../sources/frontend/shadcn-svelte.md)
- [sqlc](../../sources/database/sqlc.md)
- [sqlc Configuration](../../sources/database/sqlc-config.md)
- [Svelte 5 Runes](../../sources/frontend/svelte-runes.md)
- [Svelte 5 Documentation](../../sources/frontend/svelte5.md)
- [SvelteKit Documentation](../../sources/frontend/sveltekit.md)
- [TanStack Query](../../sources/frontend/tanstack-query.md)
- [Typesense API](../../sources/infrastructure/typesense.md)
- [Typesense Go Client](../../sources/infrastructure/typesense-go.md)



---

**Need Help?** [Open an issue](https://github.com/revenge-project/revenge/issues) or [Join the discussion](https://github.com/revenge-project/revenge/discussions)