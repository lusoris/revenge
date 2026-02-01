## Table of Contents

- [Testing Patterns and Practices](#testing-patterns-and-practices)
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
  - name: Go Testing Package
    url: ../sources/go/stdlib/testing.md
    note: Standard library testing
  - name: testify
    url: ../sources/testing/testify.md
    note: Assertion and mocking toolkit
  - name: mockery
    url: ../sources/testing/mockery-guide.md
    note: Mock generation from interfaces
  - name: testcontainers-go
    url: https://pkg.go.dev/github.com/testcontainers/testcontainers-go
    note: Docker containers for integration tests
  - name: embedded-postgres
    url: ../sources/testing/embedded-postgres.md
    note: Embedded PostgreSQL for testing
design_refs:
  - title: technical
    path: technical.md
  - title: 01_ARCHITECTURE
    path: architecture/01_ARCHITECTURE.md
  - title: BEST_PRACTICES
    path: operations/BEST_PRACTICES.md
---

# Testing Patterns and Practices

<!-- DESIGN: technical, README, SCAFFOLD_TEMPLATE, test_output_claude -->


**Created**: 2026-01-31
**Status**: ✅ Complete
**Category**: technical


> > Comprehensive testing strategy with unit, integration, and end-to-end tests

Testing approach:
- **Unit Tests**: Fast, isolated tests with mocks (testify + mockery)
- **Integration Tests**: Real dependencies via testcontainers
- **Table-Driven Tests**: Go testing pattern for multiple scenarios
- **Coverage**: Minimum 80% code coverage required
- **CI/CD**: Automated testing on every PR and commit
- **Best Practices**: Test-driven development, clear test names, arrange-act-assert

---


## Status

| Dimension | Status | Notes |
|-----------|--------|-------|
| Design | ✅ | Complete testing patterns documentation |
| Sources | ✅ | All testing tools documented |
| Instructions | ✅ | Generated from design |
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
- [technical](technical.md)
- [01_ARCHITECTURE](../architecture/01_ARCHITECTURE.md)
- [BEST_PRACTICES](../operations/BEST_PRACTICES.md)

### External Sources
- [Go Testing Package](../sources/go/stdlib/testing.md) - Standard library testing
- [testify](../sources/testing/testify.md) - Assertion and mocking toolkit
- [mockery](../sources/testing/mockery-guide.md) - Mock generation from interfaces
- [testcontainers-go](https://pkg.go.dev/github.com/testcontainers/testcontainers-go) - Docker containers for integration tests
- [embedded-postgres](../sources/testing/embedded-postgres.md) - Embedded PostgreSQL for testing

