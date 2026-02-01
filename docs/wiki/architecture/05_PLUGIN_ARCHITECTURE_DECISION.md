## Table of Contents

- [Plugin Architecture Decision](#plugin-architecture-decision)
  - [Contents](#contents)
  - [How It Works](#how-it-works)
  - [Features](#features)
  - [Configuration](#configuration)
  - [Related Documentation](#related-documentation)
    - [See Also](#see-also)

---
sources:
- name: Dragonfly Documentation
    url: ../sources/infrastructure/dragonfly.md
    note: Auto-resolved from dragonfly
- name: pgx PostgreSQL Driver
    url: ../sources/database/pgx.md
    note: Auto-resolved from pgx
- name: PostgreSQL Arrays
    url: ../sources/database/postgresql-arrays.md
    note: Auto-resolved from postgresql-arrays
- name: PostgreSQL JSON Functions
    url: ../sources/database/postgresql-json.md
    note: Auto-resolved from postgresql-json
- name: River Job Queue
    url: ../sources/tooling/river.md
    note: Auto-resolved from river
- name: rueidis
    url: ../sources/tooling/rueidis.md
    note: Auto-resolved from rueidis
- name: rueidis GitHub README
    url: ../sources/tooling/rueidis-guide.md
    note: Auto-resolved from rueidis-docs
- name: Typesense API
    url: ../sources/infrastructure/typesense.md
    note: Auto-resolved from typesense
- name: Typesense Go Client
    url: ../sources/infrastructure/typesense-go.md
    note: Auto-resolved from typesense-go
design_refs:
- title: architecture
    path: architecture/INDEX.md
- title: ADULT_CONTENT_SYSTEM
    path: ADULT_CONTENT_SYSTEM.md
- title: ADULT_METADATA
    path: ADULT_METADATA.md
- title: DATA_RECONCILIATION
    path: DATA_RECONCILIATION.md
---

# Plugin Architecture Decision

> Why Revenge uses integrations instead of plugins

Revenge deliberately chose not to implement a plugin system. Instead, common integrations (Radarr, Sonarr, TMDb, etc.) are built directly into the codebase with first-class support. This means faster development, better security (no arbitrary code execution), and simpler maintenance. External systems can still integrate via webhooks and the REST API. For power users who need custom automation, scripting support (Lua or Starlark) may be added in the future.

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
