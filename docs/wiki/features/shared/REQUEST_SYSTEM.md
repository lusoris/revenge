## Table of Contents

- [Native Request System](#native-request-system)
- [Native Request System](#native-request-system)
  - [Contents](#contents)
  - [How It Works](#how-it-works)
    - [Content Flow](#content-flow)
  - [Features](#features)
  - [Configuration](#configuration)
  - [Related Documentation](#related-documentation)
    - [See Also](#see-also)



---
sources:
  - name: Uber fx
    url: ../../../sources/tooling/fx.md
    note: Auto-resolved from fx
  - name: ogen OpenAPI Generator
    url: ../../../sources/tooling/ogen.md
    note: Auto-resolved from ogen
  - name: pgx PostgreSQL Driver
    url: ../../../sources/database/pgx.md
    note: Auto-resolved from pgx
  - name: PostgreSQL Arrays
    url: ../../../sources/database/postgresql-arrays.md
    note: Auto-resolved from postgresql-arrays
  - name: PostgreSQL JSON Functions
    url: ../../../sources/database/postgresql-json.md
    note: Auto-resolved from postgresql-json
  - name: River Job Queue
    url: ../../../sources/tooling/river.md
    note: Auto-resolved from river
  - name: sqlc
    url: ../../../sources/database/sqlc.md
    note: Auto-resolved from sqlc
  - name: sqlc Configuration
    url: ../../../sources/database/sqlc-config.md
    note: Auto-resolved from sqlc-config
  - name: Svelte 5 Runes
    url: ../../../sources/frontend/svelte-runes.md
    note: Auto-resolved from svelte-runes
  - name: Svelte 5 Documentation
    url: ../../../sources/frontend/svelte5.md
    note: Auto-resolved from svelte5
  - name: SvelteKit Documentation
    url: ../../../sources/frontend/sveltekit.md
    note: Auto-resolved from sveltekit
design_refs:
  - title: features/shared
    path: features/shared.md
  - title: 01_ARCHITECTURE
    path: architecture/01_ARCHITECTURE.md
  - title: 02_DESIGN_PRINCIPLES
    path: architecture/02_DESIGN_PRINCIPLES.md
  - title: 03_METADATA_SYSTEM
    path: architecture/03_METADATA_SYSTEM.md
---

# Native Request System




# Native Request System

> Let users request movies, shows, and music directly

Built-in request system replaces Overseerr/Jellyseerr. Users browse TMDB, TheTVDB, or MusicBrainz and request content they want. Requests route to the appropriate Arr service (Radarr, Sonarr, Lidarr) for automatic download. Admins can approve/deny requests or set auto-approval rules. Track request status and get notified when content arrives.

---




## Contents

<!-- TOC will be auto-generated here by markdown-toc -->

---


## How It Works

<!-- Wiki how it works description pending -->

### Content Flow

1. **Add Content**: Import or add manually
2. **Metadata**: Automatically fetched
3. **Library**: Organized by folder structure
4. **Playback**: Stream directly or use external player




## Features
<!-- Feature list placeholder -->



## Configuration
<!-- User-friendly configuration guide -->









## Related Documentation
### See Also
<!-- Related wiki pages -->



---

**Need Help?** [Open an issue](https://github.com/revenge-project/revenge/issues) or [Join the discussion](https://github.com/revenge-project/revenge/discussions)