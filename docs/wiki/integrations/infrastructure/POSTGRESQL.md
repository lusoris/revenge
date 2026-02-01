## Table of Contents

- [PostgreSQL](#postgresql)
- [PostgreSQL](#postgresql)
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
- name: Prometheus Go Client
    url: ../sources/observability/prometheus.md
    note: Auto-resolved from prometheus
- name: Prometheus Metric Types
    url: ../sources/observability/prometheus-metrics.md
    note: Auto-resolved from prometheus-metrics
- name: River Job Queue
    url: ../sources/tooling/river.md
    note: Auto-resolved from river
- name: River Documentation
    url: ../sources/tooling/river-guide.md
    note: Auto-resolved from river-docs
- name: sqlc
    url: ../sources/database/sqlc.md
    note: Auto-resolved from sqlc
- name: sqlc Configuration
    url: ../sources/database/sqlc-config.md
    note: Auto-resolved from sqlc-config
- name: Typesense API
    url: ../sources/infrastructure/typesense.md
    note: Auto-resolved from typesense
- name: Typesense Go Client
    url: ../sources/infrastructure/typesense-go.md
    note: Auto-resolved from typesense-go
design_refs:
- title: integrations/infrastructure
    path: integrations/infrastructure.md
- title: 01_ARCHITECTURE
    path: architecture/01_ARCHITECTURE.md
- title: 02_DESIGN_PRINCIPLES
    path: architecture/02_DESIGN_PRINCIPLES.md
- title: 03_METADATA_SYSTEM
    path: architecture/03_METADATA_SYSTEM.md
---

# PostgreSQL

# PostgreSQL

> PostgreSQL stores all your media data

PostgreSQL is the only supported database for Revenge (no SQLite). Stores all metadata, user data, watch history, and configuration. Version 18+ required for modern features. Supports automatic schema migrations on upgrade. Backup your database regularly using standard PostgreSQL tools.

---

## Contents

<!-- TOC will be auto-generated here by markdown-toc -->

---

## How It Works

<!-- How it works -->

## Features
<!-- Feature list placeholder -->

## Configuration
<!-- User-friendly configuration guide -->

## Related Documentation
### See Also
<!-- Related wiki pages -->

---

**Need Help?** [Open an issue](https://github.com/revenge-project/revenge/issues) or [Join the discussion](https://github.com/revenge-project/revenge/discussions)
