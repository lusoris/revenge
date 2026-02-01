## Table of Contents

- [Testing Patterns](#testing-patterns)
  - [Contents](#contents)
  - [How It Works](#how-it-works)
  - [Features](#features)
  - [Configuration](#configuration)
  - [Related Documentation](#related-documentation)
    - [See Also](#see-also)

---
sources:
- name: Go Testing
    url: ../sources/go/stdlib/testing.md
    note: Standard library testing
- name: Testify
    url: https://github.com/stretchr/testify
    note: Assertions and test suites
- name: Mockery
    url: ../sources/testing/mockery-guide.md
    note: Mock generation
- name: Testcontainers Go
    url: ../sources/testing/testcontainers.md
    note: Integration testing with containers
- name: Embedded Postgres
    url: ../sources/testing/embedded-postgres-guide.md
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

> Testing standards and patterns for code quality

Revenge maintains high code quality through comprehensive testing. All new code requires 80% test coverage minimum. Unit tests use table-driven patterns for clarity and completeness, with mocks generated automatically via Mockery. Integration tests run against real PostgreSQL using embedded-postgres for speed or testcontainers for full fidelity. The CI pipeline runs all tests on every pull request, blocking merges that reduce coverage or introduce failing tests.

---

## Contents

<!-- TOC will be auto-generated here by markdown-toc -->

---

## How It Works

<!-- User-friendly explanation -->

## Features
<!-- Feature list placeholder -->

## Configuration
<!-- User-friendly configuration guide -->

## Related Documentation
### See Also
<!-- Related wiki pages -->

---

**Need Help?** [Open an issue](https://github.com/revenge-project/revenge/issues) or [Join the discussion](https://github.com/revenge-project/revenge/discussions)
