## Table of Contents

- [Development Environment Setup](#development-environment-setup)
  - [Status](#status)
  - [Related Documentation](#related-documentation)
    - [Design Documents](#design-documents)
    - [External Sources](#external-sources)

# Development Environment Setup

<!-- DESIGN: operations, README, test_output_claude, test_output_wiki -->


**Created**: 2026-01-31
**Status**: ✅ Complete
**Category**: operations


> > Local development environment setup with Go, Node.js, PostgreSQL, and hot reload

Complete guide for setting up Revenge development environment:
- **Prerequisites**: Go 1.25+, Node.js 20.x, PostgreSQL 18+, Python 3.12
- **Hot Reload**: air for Go backend, Vite for frontend
- **Database**: PostgreSQL only (no SQLite support)
- **Build Flags**: GOEXPERIMENT=greenteagc,jsonv2


---


## Status

| Dimension | Status | Notes |
|-----------|--------|-------|
| Design | ✅ | Complete developer setup guide |
| Sources | ✅ | All tools documented |
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
- [BEST_PRACTICES](../operations/CODING_STANDARDS.md)
- [00_SOURCE_OF_TRUTH](../00_SOURCE_OF_TRUTH.md)

### External Sources
- [Go Installation](https://go.dev/doc/install) - Official Go installation guide
- [air Hot Reload](../../sources/go_dev_tools/air/main.md) - Go hot reload tool
- [Vite](https://vitejs.dev/guide/) - Frontend build tool
- [PostgreSQL Downloads](https://www.postgresql.org/download/) - PostgreSQL installation
- [Node.js Downloads](https://nodejs.org/) - Node.js LTS releases
- [sqlc](../../sources/database/sqlc.md) - SQL code generation
- [Conventional Commits](../../sources/standards/conventional-commits.md) - Commit message format

