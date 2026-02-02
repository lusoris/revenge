## Table of Contents

- [Testing Patterns](#testing-patterns)
  - [Status](#status)
  - [Related Documentation](#related-documentation)
    - [Design Documents](#design-documents)
    - [External Sources](#external-sources)

# Testing Patterns

<!-- DESIGN: patterns, README, test_output_claude, test_output_wiki -->


**Created**: 2026-01-31
**Status**: ✅ Complete
**Category**: pattern


> > Table-driven tests, mocking patterns, and integration testing with testcontainers

Standard testing approach for Revenge project:
- **Table-Driven Tests**: All logic tests use table-driven pattern
- **Testify**: Assertions and test suites
- **Mockery**: Auto-generated mocks from interfaces
- **Embedded Postgres**: Fast unit tests without containers
- **Testcontainers**: Integration tests with real PostgreSQL, Dragonfly


---


## Status

| Dimension | Status | Notes |
|-----------|--------|-------|
| Design | ✅ | - |
| Sources | ✅ | - |
| Instructions | ✅ | Generated from design |
| Code | 🔴 | - |
| Linting | 🔴 | - |
| Unit Testing | 🔴 | - |
| Integration Testing | 🔴 | - |

**Overall**: ✅ Complete


## Related Documentation
### Design Documents
- [patterns](INDEX.md)
- [BEST_PRACTICES](../operations/BEST_PRACTICES.md)
- [02_DESIGN_PRINCIPLES](../architecture/02_DESIGN_PRINCIPLES.md)

### External Sources
- [Go Testing](../../sources/go/stdlib/testing.md) - Standard library testing
- [Testify](https://github.com/stretchr/testify) - Assertions and test suites
- [Mockery](../../sources/testing/mockery-guide.md) - Mock generation
- [Testcontainers Go](../../sources/testing/testcontainers.md) - Integration testing with containers
- [Embedded Postgres](../../sources/testing/embedded-postgres-guide.md) - Fast PostgreSQL for unit tests

