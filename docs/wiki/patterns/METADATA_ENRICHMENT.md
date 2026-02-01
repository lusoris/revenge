## Table of Contents

- [Metadata Enrichment Pattern](#metadata-enrichment-pattern)
  - [Contents](#contents)
  - [How It Works](#how-it-works)
  - [Features](#features)
  - [Configuration](#configuration)
  - [Related Documentation](#related-documentation)
    - [See Also](#see-also)

---
sources:
- name: River Job Queue
    url: ../sources/tooling/river.md
    note: Background job processing
- name: rueidis
    url: ../sources/tooling/rueidis.md
    note: Distributed cache (L2)
- name: Otter
    url: https://pkg.go.dev/github.com/maypok86/otter
    note: In-memory cache (L1)
- name: Sturdyc
    url: ../sources/tooling/sturdyc-guide.md
    note: Request coalescing cache
design_refs:
- title: patterns
    path: patterns.md
- title: 01_ARCHITECTURE
    path: architecture/01_ARCHITECTURE.md
- title: 02_DESIGN_PRINCIPLES
    path: architecture/02_DESIGN_PRINCIPLES.md
- title: 03_METADATA_SYSTEM
    path: architecture/03_METADATA_SYSTEM.md
---

# Metadata Enrichment Pattern

> Fast, efficient metadata with intelligent caching and background enrichment

The Metadata Enrichment Pattern ensures fast UI response times while maintaining rich, up-to-date metadata. The system uses a five-tier priority chain starting with local cache for instant results, falling back to Arr services (which cache upstream data), then internal sources, external APIs, and finally background enrichment jobs. Multi-tier caching with request coalescing prevents duplicate API calls and reduces external API quota consumption.

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
