## Table of Contents

- [Reverse Proxy Configuration](#reverse-proxy-configuration)
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
  - name: gohlslib (HLS)
    url: ../sources/media/gohlslib.md
    note: Auto-resolved from gohlslib
  - name: M3U8 Extended Format
    url: ../sources/protocols/m3u8.md
    note: Auto-resolved from m3u8
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
  - name: rueidis
    url: ../sources/tooling/rueidis.md
    note: Auto-resolved from rueidis
  - name: rueidis GitHub README
    url: ../sources/tooling/rueidis-guide.md
    note: Auto-resolved from rueidis-docs
design_refs:
  - title: operations
    path: operations/INDEX.md
  - title: 01_ARCHITECTURE
    path: architecture/01_ARCHITECTURE.md
  - title: 02_DESIGN_PRINCIPLES
    path: architecture/02_DESIGN_PRINCIPLES.md
  - title: 03_METADATA_SYSTEM
    path: architecture/03_METADATA_SYSTEM.md
---

# Reverse Proxy Configuration




> Configure Traefik, Caddy, or nginx for HTTPS access to Revenge

A reverse proxy sits between the internet and Revenge, handling HTTPS encryption and routing. Traefik is recommended for Docker setups with automatic SSL certificates. Caddy offers the simplest configuration with automatic HTTPS. nginx provides the best performance but requires manual SSL setup. All options support WebSocket connections required for real-time features like SyncPlay. Includes complete configuration examples for each proxy.

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