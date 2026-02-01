## Table of Contents

- [Revenge - Metadata System](#revenge-metadata-system)
  - [Status](#status)
  - [Related Documentation](#related-documentation)
    - [Design Documents](#design-documents)
    - [External Sources](#external-sources)

# Revenge - Metadata System


**Created**: 2026-01-31
**Status**: ✅ Complete
**Category**: architecture


> > Multi-source metadata system with caching and priority chain

Metadata handling:
- **Priority Chain**: Local cache → Arr services → Internal (Stash) → External APIs
- **Providers**: TMDb, TheTVDB, MusicBrainz, StashDB, and many more
- **Caching**: Two-tier with otter (L1 memory) and rueidis (L2 distributed)
- **Enrichment**: Background jobs for additional metadata, thumbnails, blurhash
- **Matching**: Fingerprinting for audio, hash matching for media


---


## Status

| Dimension | Status | Notes |
|-----------|--------|-------|
| Design | ✅ | - |
| Sources | 🟡 | - |
| Instructions | ✅ | Generated from design |
| Code | 🔴 | - |
| Linting | 🔴 | - |
| Unit Testing | 🔴 | - |
| Integration Testing | 🔴 | - |

**Overall**: ✅ Complete



---



















## Related Documentation
### Design Documents
- [architecture](INDEX.md)
- [ADULT_CONTENT_SYSTEM](../features/adult/ADULT_CONTENT_SYSTEM.md)
- [ADULT_METADATA](../features/adult/ADULT_METADATA.md)
- [DATA_RECONCILIATION](../features/adult/DATA_RECONCILIATION.md)

### External Sources
- [Dragonfly Documentation](../../sources/infrastructure/dragonfly.md) - Auto-resolved from dragonfly
- [go-blurhash](../../sources/media/go-blurhash.md) - Auto-resolved from go-blurhash
- [Last.fm API](../../sources/apis/lastfm.md) - Auto-resolved from lastfm-api
- [pgx PostgreSQL Driver](../../sources/database/pgx.md) - Auto-resolved from pgx
- [PostgreSQL Arrays](../../sources/database/postgresql-arrays.md) - Auto-resolved from postgresql-arrays
- [PostgreSQL JSON Functions](../../sources/database/postgresql-json.md) - Auto-resolved from postgresql-json
- [River Job Queue](../../sources/tooling/river.md) - Auto-resolved from river
- [rueidis](../../sources/tooling/rueidis.md) - Auto-resolved from rueidis
- [rueidis GitHub README](../../sources/tooling/rueidis-guide.md) - Auto-resolved from rueidis-docs

