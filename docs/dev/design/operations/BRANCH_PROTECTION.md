## Table of Contents

- [Branch Protection Rules](#branch-protection-rules)
  - [Status](#status)
  - [Related Documentation](#related-documentation)
    - [Design Documents](#design-documents)
    - [External Sources](#external-sources)

# Branch Protection Rules

<!-- DESIGN: operations, README, 01_ARCHITECTURE, 02_DESIGN_PRINCIPLES -->


**Created**: 2026-01-31
**Status**: ✅ Complete
**Category**: operations


>   > GitHub branch protection rules for main and develop

  Protection rules:
  - **Required Reviews**: 1 approval for PRs to main/develop
  - **Status Checks**: CI must pass (tests, lint, coverage)
  - **No Force Push**: Prevent history rewriting on protected branches
  - **Linear History**: Require merge commits or squash
  - **Up-to-date**: Branch must be current with base before merge

---


## Status

| Dimension | Status | Notes |
|-----------|--------|-------|
| Design | ✅ | Complete branch protection guide |
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
- [01_ARCHITECTURE](../architecture/01_ARCHITECTURE.md)
- [02_DESIGN_PRINCIPLES](../architecture/02_DESIGN_PRINCIPLES.md)
- [03_METADATA_SYSTEM](../architecture/03_METADATA_SYSTEM.md)

### External Sources
- [Conventional Commits](../../sources/standards/conventional-commits.md) - Auto-resolved from conventional-commits
- [Git Flow](../../sources/standards/gitflow.md) - Auto-resolved from gitflow

