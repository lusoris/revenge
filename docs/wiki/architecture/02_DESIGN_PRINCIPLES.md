## Table of Contents

- [Revenge - Design Principles](#revenge-design-principles)
  - [Contents](#contents)
  - [How It Works](#how-it-works)
  - [Features](#features)
  - [Configuration](#configuration)
  - [Related Documentation](#related-documentation)
    - [Related Pages](#related-pages)
    - [Learn More](#learn-more)



---
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
### Related Pages
<!-- Related wiki pages -->

### Learn More

Official documentation and guides:
- [Dragonfly Documentation](../../sources/infrastructure/dragonfly.md)
- [FFmpeg Documentation](../../sources/media/ffmpeg.md)
- [FFmpeg Codecs](../../sources/media/ffmpeg-codecs.md)
- [FFmpeg Formats](../../sources/media/ffmpeg-formats.md)
- [go-astiav (FFmpeg bindings)](../../sources/media/go-astiav.md)
- [go-astiav GitHub README](../../sources/media/go-astiav-guide.md)
- [gohlslib (HLS)](../../sources/media/gohlslib.md)
- [M3U8 Extended Format](../../sources/protocols/m3u8.md)
- [River Job Queue](../../sources/tooling/river.md)
- [shadcn-svelte](../../sources/frontend/shadcn-svelte.md)
- [Svelte 5 Runes](../../sources/frontend/svelte-runes.md)
- [Svelte 5 Documentation](../../sources/frontend/svelte5.md)
- [SvelteKit Documentation](../../sources/frontend/sveltekit.md)
- [TanStack Query](../../sources/frontend/tanstack-query.md)
- [Typesense API](../../sources/infrastructure/typesense.md)
- [Typesense Go Client](../../sources/infrastructure/typesense-go.md)



---

**Need Help?** [Open an issue](https://github.com/revenge-project/revenge/issues) or [Join the discussion](https://github.com/revenge-project/revenge/discussions)