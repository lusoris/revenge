## Table of Contents

- [Revenge - Design Principles](#revenge-design-principles)
  - [Status](#status)
  - [Related Documentation](#related-documentation)
    - [Design Documents](#design-documents)
    - [External Sources](#external-sources)

# Revenge - Design Principles

<!-- DESIGN: architecture, README, test_output_claude, test_output_wiki -->


**Created**: 2026-01-31
**Status**: ✅ Complete
**Category**: architecture


> > Core design principles and patterns guiding Revenge development

Key principles:
- **PostgreSQL Only**: No SQLite - simplifies codebase, pgx is best driver
- **80% Test Coverage**: Required minimum for all packages
- **Local-First Metadata**: Always prefer cached/local data over external APIs
- **Table-Driven Tests**: testify + mockery for consistent testing
- **Sentinel Errors**: Type-safe error handling with wrapped contexts


---


## Status

| Dimension | Status | Notes |
|-----------|--------|-------|
| Design | ✅ | - |
| Sources | ⚪ | - |
| Instructions | ✅ | Generated from design |
| Code | 🔴 | - |
| Linting | 🔴 | - |
| Unit Testing | 🔴 | - |
| Integration Testing | 🔴 | - |

**Overall**: ✅ Complete


## Related Documentation
### Design Documents
- [architecture](INDEX.md)
- [01_ARCHITECTURE](ARCHITECTURE.md)
- [03_METADATA_SYSTEM](METADATA_SYSTEM.md)
- [METADATA (service)](../services/METADATA.md)
- [HTTP_CLIENT (service)](../services/HTTP.md)
- [ADULT_CONTENT_SYSTEM](../features/adult/ADULT_CONTENT_SYSTEM.md)
- [ADULT_METADATA](../features/adult/ADULT_METADATA.md)
- [DATA_RECONCILIATION](../features/adult/DATA_RECONCILIATION.md)

### External Sources
- [Dragonfly Documentation](../../sources/infrastructure/dragonfly.md) - Auto-resolved from dragonfly
- [FFmpeg Documentation](../../sources/media/ffmpeg.md) - Auto-resolved from ffmpeg
- [FFmpeg Codecs](../../sources/media/ffmpeg-codecs.md) - Auto-resolved from ffmpeg-codecs
- [FFmpeg Formats](../../sources/media/ffmpeg-formats.md) - Auto-resolved from ffmpeg-formats
- [go-astiav (FFmpeg bindings)](../../sources/media/go-astiav.md) - Auto-resolved from go-astiav
- [go-astiav GitHub README](../../sources/media/go-astiav-guide.md) - Auto-resolved from go-astiav-docs
- [gohlslib (HLS)](../../sources/media/gohlslib.md) - Auto-resolved from gohlslib
- [M3U8 Extended Format](../../sources/protocols/m3u8.md) - Auto-resolved from m3u8
- [River Job Queue](../../sources/tooling/river.md) - Auto-resolved from river
- [shadcn-svelte](../../sources/frontend/shadcn-svelte.md) - Auto-resolved from shadcn-svelte
- [Svelte 5 Runes](../../sources/frontend/svelte-runes.md) - Auto-resolved from svelte-runes
- [Svelte 5 Documentation](../../sources/frontend/svelte5.md) - Auto-resolved from svelte5
- [SvelteKit Documentation](../../sources/frontend/sveltekit.md) - Auto-resolved from sveltekit
- [TanStack Query](../../sources/frontend/tanstack-query.md) - Auto-resolved from tanstack-query
- [Typesense API](../../sources/infrastructure/typesense.md) - Auto-resolved from typesense
- [Typesense Go Client](../../sources/infrastructure/typesense-go.md) - Auto-resolved from typesense-go

