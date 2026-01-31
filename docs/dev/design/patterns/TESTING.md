## Table of Contents

- [Testing Patterns](#testing-patterns)
  - [Status](#status)
  - [Architecture](#architecture)
    - [Components](#components)
  - [Implementation](#implementation)
    - [File Structure](#file-structure)
    - [Key Interfaces](#key-interfaces)
    - [Dependencies](#dependencies)
  - [Configuration](#configuration)
    - [Environment Variables](#environment-variables)
    - [Config Keys](#config-keys)
  - [Testing Strategy](#testing-strategy)
    - [Unit Tests](#unit-tests)
    - [Integration Tests](#integration-tests)
    - [Test Coverage](#test-coverage)
  - [Related Documentation](#related-documentation)
    - [Design Documents](#design-documents)
    - [External Sources](#external-sources)



---
sources:
  - name: Go Testing
    url: https://pkg.go.dev/testing
    note: Standard library testing
  - name: Testify
    url: https://github.com/stretchr/testify
    note: Assertions and test suites
  - name: Mockery
    url: https://vektra.github.io/mockery/
    note: Mock generation
  - name: Testcontainers Go
    url: https://golang.testcontainers.org/
    note: Integration testing with containers
  - name: Embedded Postgres
    url: https://github.com/fergusstrange/embedded-postgres
    note: Fast PostgreSQL for unit tests
design_refs:
  - title: patterns
    path: patterns/INDEX.md
  - title: BEST_PRACTICES
    path: operations/BEST_PRACTICES.md
  - title: 02_DESIGN_PRINCIPLES
    path: architecture/02_DESIGN_PRINCIPLES.md
---

# Testing Patterns


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
| Instructions | 🔴 | - |
| Code | 🔴 | - |
| Linting | 🔴 | - |
| Unit Testing | 🔴 | - |
| Integration Testing | 🔴 | - |

**Overall**: ✅ Complete



---


## Architecture

<!-- Architecture diagram placeholder -->

### Components

<!-- Component description -->


## Implementation

### File Structure

<!-- File structure -->

### Key Interfaces

<!-- Interface definitions -->

### Dependencies

<!-- Dependency list -->





## Configuration
### Environment Variables

<!-- Environment variables -->

### Config Keys

<!-- Configuration keys -->




## Testing Strategy

### Unit Tests

<!-- Unit test strategy -->

### Integration Tests

<!-- Integration test strategy -->

### Test Coverage

Target: **80% minimum**







## Related Documentation
### Design Documents
- [patterns](patterns/INDEX.md)
- [BEST_PRACTICES](operations/BEST_PRACTICES.md)
- [02_DESIGN_PRINCIPLES](architecture/02_DESIGN_PRINCIPLES.md)

### External Sources
- [Go Testing](https://pkg.go.dev/testing) - Standard library testing
- [Testify](https://github.com/stretchr/testify) - Assertions and test suites
- [Mockery](https://vektra.github.io/mockery/) - Mock generation
- [Testcontainers Go](https://golang.testcontainers.org/) - Integration testing with containers
- [Embedded Postgres](https://github.com/fergusstrange/embedded-postgres) - Fast PostgreSQL for unit tests

