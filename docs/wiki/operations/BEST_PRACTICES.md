## Table of Contents

- [Development Best Practices](#development-best-practices)
  - [Contents](#contents)
  - [How It Works](#how-it-works)
  - [Features](#features)
  - [Configuration](#configuration)
  - [Related Documentation](#related-documentation)
    - [See Also](#see-also)



---
sources:
  - name: Dragonfly Documentation
    url: ../../sources/infrastructure/dragonfly.md
    note: Auto-resolved from dragonfly
  - name: gohlslib (HLS)
    url: ../../sources/media/gohlslib.md
    note: Auto-resolved from gohlslib
  - name: koanf
    url: ../../sources/tooling/koanf.md
    note: Auto-resolved from koanf
  - name: M3U8 Extended Format
    url: ../../sources/protocols/m3u8.md
    note: Auto-resolved from m3u8
  - name: otter Cache
    url: ../../sources/tooling/otter.md
    note: Auto-resolved from otter
  - name: pgx PostgreSQL Driver
    url: ../../sources/database/pgx.md
    note: Auto-resolved from pgx
  - name: PostgreSQL Arrays
    url: ../../sources/database/postgresql-arrays.md
    note: Auto-resolved from postgresql-arrays
  - name: PostgreSQL JSON Functions
    url: ../../sources/database/postgresql-json.md
    note: Auto-resolved from postgresql-json
  - name: River Job Queue
    url: ../../sources/tooling/river.md
    note: Auto-resolved from river
  - name: rueidis
    url: ../../sources/tooling/rueidis.md
    note: Auto-resolved from rueidis
  - name: rueidis GitHub README
    url: ../../sources/tooling/rueidis-guide.md
    note: Auto-resolved from rueidis-docs
  - name: sturdyc
    url: ../../sources/tooling/sturdyc.md
    note: Auto-resolved from sturdyc
  - name: sturdyc GitHub README
    url: ../../sources/tooling/sturdyc-guide.md
    note: Auto-resolved from sturdyc-docs
design_refs:
  - title: operations
    path: INDEX.md
  - title: 01_ARCHITECTURE
    path: ../architecture/01_ARCHITECTURE.md
  - title: 02_DESIGN_PRINCIPLES
    path: ../architecture/02_DESIGN_PRINCIPLES.md
  - title: 03_METADATA_SYSTEM
    path: ../architecture/03_METADATA_SYSTEM.md
---

# Development Best Practices




> Coding standards and patterns for contributing to Revenge

Follow these guidelines when contributing code to Revenge. Use gofmt for formatting and golangci-lint for linting. All code requires 80% minimum test coverage with table-driven tests. Use the repository pattern with fx dependency injection. Handle errors with sentinel errors and wrap with context using %w. Cache expensive operations with otter (L1) and rueidis (L2). See code examples for common patterns.

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