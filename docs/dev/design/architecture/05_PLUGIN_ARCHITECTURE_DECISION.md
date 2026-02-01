## Table of Contents

- [Plugin Architecture Decision](#plugin-architecture-decision)
  - [Status](#status)
  - [Related Documentation](#related-documentation)
    - [Design Documents](#design-documents)
    - [External Sources](#external-sources)

# Plugin Architecture Decision


**Created**: 2026-01-31
**Status**: ✅ Complete
**Category**: architecture


> > ADR: Decision to use integrations over plugins

Decision rationale:
- **No Plugin System**: Revenge uses direct integrations instead of plugins
- **Why**: Simpler maintenance, better security, faster development
- **Integrations**: First-class support for common services (Arr stack, metadata providers)
- **Webhooks**: External systems can integrate via webhooks
- **Future**: May add scripting for power users (Lua or Starlark)


---


## Status

| Dimension | Status | Notes |
|-----------|--------|-------|
| Design | ✅ | - |
| Sources | ⚪ | - |
| Instructions | ⚪ | - |
| Code | ⚪ | - |
| Linting | ⚪ | - |
| Unit Testing | ⚪ | - |
| Integration Testing | ⚪ | - |

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
- [pgx PostgreSQL Driver](../../sources/database/pgx.md) - Auto-resolved from pgx
- [PostgreSQL Arrays](../../sources/database/postgresql-arrays.md) - Auto-resolved from postgresql-arrays
- [PostgreSQL JSON Functions](../../sources/database/postgresql-json.md) - Auto-resolved from postgresql-json
- [River Job Queue](../../sources/tooling/river.md) - Auto-resolved from river
- [rueidis](../../sources/tooling/rueidis.md) - Auto-resolved from rueidis
- [rueidis GitHub README](../../sources/tooling/rueidis-guide.md) - Auto-resolved from rueidis-docs
- [Typesense API](../../sources/infrastructure/typesense.md) - Auto-resolved from typesense
- [Typesense Go Client](../../sources/infrastructure/typesense-go.md) - Auto-resolved from typesense-go

