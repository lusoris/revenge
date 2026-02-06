## Table of Contents

- [Semantic Versioning & Releases](#semantic-versioning-releases)
  - [Status](#status)
  - [Related Documentation](#related-documentation)
    - [Design Documents](#design-documents)
    - [External Sources](#external-sources)

# Semantic Versioning & Releases

<!-- DESIGN: operations, README, test_output_claude, test_output_wiki -->


**Created**: 2026-01-31
**Status**: ✅ Complete
**Category**: operations


>   > Semantic versioning (semver) and automated release process

  Versioning strategy:
  - **Format**: vMAJOR.MINOR.PATCH (e.g., v1.2.3)
  - **MAJOR**: Breaking changes to API or behavior
  - **MINOR**: New features, backwards compatible
  - **PATCH**: Bug fixes, security patches
  - **Automation**: Release Please for changelog and version bumps

---


## Status

| Dimension | Status | Notes |
|-----------|--------|-------|
| Design | ✅ | Complete versioning guide |
| Sources | 🔴 | - |
| Instructions | ✅ | Generated from design |
| Code | 🔴 | - |
| Linting | 🔴 | - |
| Unit Testing | 🔴 | - |
| Integration Testing | 🔴 | - |

**Overall**: ✅ Complete


## Related Documentation
### Design Documents
- [operations](INDEX.md)
- [01_ARCHITECTURE](../architecture/ARCHITECTURE.md)
- [02_DESIGN_PRINCIPLES](../architecture/DESIGN_PRINCIPLES.md)
- [03_METADATA_SYSTEM](../architecture/METADATA_SYSTEM.md)

### External Sources
- [pgx PostgreSQL Driver](../../sources/database/pgx.md) - Auto-resolved from pgx
- [PostgreSQL Arrays](../../sources/database/postgresql-arrays.md) - Auto-resolved from postgresql-arrays
- [PostgreSQL JSON Functions](../../sources/database/postgresql-json.md) - Auto-resolved from postgresql-json
- [River Job Queue](../../sources/tooling/river.md) - Auto-resolved from river
- [Semantic Versioning](../../sources/standards/semver.md) - Auto-resolved from semver

