## Table of Contents

- [Metadata Service](#metadata-service)
- [Metadata Service](#metadata-service)
  - [Contents](#contents)
  - [How It Works](#how-it-works)
  - [Features](#features)
  - [Configuration](#configuration)
  - [Related Documentation](#related-documentation)
    - [See Also](#see-also)



---
sources:
  - name: Uber fx
    url: ../../sources/tooling/fx.md
    note: Auto-resolved from fx
  - name: Last.fm API
    url: ../../sources/apis/lastfm.md
    note: Auto-resolved from lastfm-api
  - name: River Job Queue
    url: ../../sources/tooling/river.md
    note: Auto-resolved from river
design_refs:
  - title: services
    path: INDEX.md
  - title: 01_ARCHITECTURE
    path: ../architecture/01_ARCHITECTURE.md
  - title: 02_DESIGN_PRINCIPLES
    path: ../architecture/02_DESIGN_PRINCIPLES.md
  - title: 03_METADATA_SYSTEM
    path: ../architecture/03_METADATA_SYSTEM.md
---

# Metadata Service




# Metadata Service

> External metadata providers for media enrichment

The Metadata service fetches information about your media from external sources like TMDb, TheTVDB, and MusicBrainz. Configure which providers to use per library. Automatic matching uses filenames and folder structure. Manual matching lets you fix incorrect matches. Metadata includes titles, descriptions, posters, cast, crew, and ratings. Background jobs keep metadata fresh.

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