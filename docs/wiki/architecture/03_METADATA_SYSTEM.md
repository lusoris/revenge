## Table of Contents

- [Revenge - Metadata System](#revenge-metadata-system)
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
  - name: go-blurhash
    url: ../../sources/media/go-blurhash.md
    note: Auto-resolved from go-blurhash
  - name: Last.fm API
    url: ../../sources/apis/lastfm.md
    note: Auto-resolved from lastfm-api
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

# Revenge - Metadata System




> How Revenge finds and stores information about your media

The metadata system gathers information about your media from multiple sources. It always checks local cache first for instant display, then queries Arr services (Radarr, Sonarr) which already have metadata, then external APIs like TMDb or MusicBrainz. Background jobs enrich media with additional data like cast info, thumbnails, and blurhash previews. Two-tier caching (memory + distributed) ensures fast lookups even for large libraries.

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