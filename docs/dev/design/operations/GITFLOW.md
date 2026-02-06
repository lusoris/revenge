## Table of Contents

- [Git Workflow & Branching Strategy](#git-workflow-branching-strategy)
  - [Status](#status)
  - [Related Documentation](#related-documentation)
    - [Design Documents](#design-documents)
    - [External Sources](#external-sources)

# Git Workflow & Branching Strategy

<!-- DESIGN: operations, README, 01_ARCHITECTURE, 02_DESIGN_PRINCIPLES -->


**Created**: 2026-01-31
**Status**: ✅ Complete
**Category**: operations


>   > Gitflow branching strategy and release process

  Branch strategy:
  - **main**: Production-ready code, protected
  - **develop**: Integration branch for next release
  - **feature/**: New features, merged to develop
  - **fix/**: Bug fixes, merged to develop or main (hotfix)
  - **release/**: Release preparation branches

---


## Status

| Dimension | Status | Notes |
|-----------|--------|-------|
| Design | ✅ | Complete Git workflow guide |
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
- [Conventional Commits](../../sources/standards/conventional-commits.md) - Auto-resolved from conventional-commits
- [Git Flow](../../sources/standards/gitflow.md) - Auto-resolved from gitflow
- [Go io](../../sources/go/stdlib/io.md) - Auto-resolved from go-io

