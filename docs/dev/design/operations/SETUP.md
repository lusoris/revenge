## Table of Contents

- [Production Deployment Setup](#production-deployment-setup)
  - [Status](#status)
  - [Related Documentation](#related-documentation)
    - [Design Documents](#design-documents)
    - [External Sources](#external-sources)

# Production Deployment Setup

<!-- DESIGN: operations, README, test_output_claude, test_output_wiki -->


**Created**: 2026-01-31
**Status**: ✅ Complete
**Category**: operations


> > Production deployment guide for self-hosting Revenge

Multiple deployment options:
- **Docker Compose**: Recommended for homelab, NAS deployments (Easy)
- **Kubernetes (K3s)**: Lightweight Kubernetes for self-hosting (Moderate)
- **Bare Metal**: Direct installation on Linux servers (Advanced)
- **Reverse Proxy**: Traefik, Caddy, nginx for HTTPS and routing


---


## Status

| Dimension | Status | Notes |
|-----------|--------|-------|
| Design | ✅ | Complete production setup guide |
| Sources | ✅ | All deployment tools documented |
| Instructions | ✅ | Generated from design |
| Code | 🔴 | - |
| Linting | 🔴 | - |
| Unit Testing | 🔴 | - |
| Integration Testing | 🔴 | - |

**Overall**: ✅ Complete


## Related Documentation
### Design Documents
- [operations](INDEX.md)
- [TECH_STACK](../technical/TECH_STACK.md)
- [REVERSE_PROXY](../operations/PROXY.md)
- [00_SOURCE_OF_TRUTH](../00_SOURCE_OF_TRUTH.md)

### External Sources
- [Dragonfly Documentation](../../sources/infrastructure/dragonfly.md) - Auto-resolved from dragonfly
- [FFmpeg Documentation](../../sources/media/ffmpeg.md) - Auto-resolved from ffmpeg
- [FFmpeg Codecs](../../sources/media/ffmpeg-codecs.md) - Auto-resolved from ffmpeg-codecs
- [FFmpeg Formats](../../sources/media/ffmpeg-formats.md) - Auto-resolved from ffmpeg-formats
- [go-astiav (FFmpeg bindings)](../../sources/media/go-astiav.md) - Auto-resolved from go-astiav
- [go-astiav GitHub README](../../sources/media/go-astiav-guide.md) - Auto-resolved from go-astiav-docs
- [Go io](../../sources/go/stdlib/io.md) - Auto-resolved from go-io
- [pgx PostgreSQL Driver](../../sources/database/pgx.md) - Auto-resolved from pgx
- [PostgreSQL Arrays](../../sources/database/postgresql-arrays.md) - Auto-resolved from postgresql-arrays
- [PostgreSQL JSON Functions](../../sources/database/postgresql-json.md) - Auto-resolved from postgresql-json
- [River Job Queue](../../sources/tooling/river.md) - Auto-resolved from river
- [rueidis](../../sources/tooling/rueidis.md) - Auto-resolved from rueidis
- [rueidis GitHub README](../../sources/tooling/rueidis-guide.md) - Auto-resolved from rueidis-docs
- [Typesense API](../../sources/infrastructure/typesense.md) - Auto-resolved from typesense
- [Typesense Go Client](../../sources/infrastructure/typesense-go.md) - Auto-resolved from typesense-go

