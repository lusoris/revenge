## Table of Contents

- [Revenge - NSFW Toggle](#revenge-nsfw-toggle)
- [Revenge - NSFW Toggle](#revenge-nsfw-toggle)
  - [Contents](#contents)
  - [How It Works](#how-it-works)
    - [Content Flow](#content-flow)
  - [Features](#features)
  - [Configuration](#configuration)
  - [Related Documentation](#related-documentation)
    - [See Also](#see-also)



---
sources:
  - name: Casbin
    url: ../../../sources/security/casbin.md
    note: Auto-resolved from casbin
  - name: Dragonfly Documentation
    url: ../../../sources/infrastructure/dragonfly.md
    note: Auto-resolved from dragonfly
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
  - name: rueidis
    url: ../../../sources/tooling/rueidis.md
    note: Auto-resolved from rueidis
  - name: rueidis GitHub README
    url: ../../../sources/tooling/rueidis-guide.md
    note: Auto-resolved from rueidis-docs
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

# Revenge - NSFW Toggle




# Revenge - NSFW Toggle

> Show or hide adult libraries with a simple toggle

Control adult content visibility per user. When disabled, adult libraries are completely hidden from the UI. Toggle requires PIN confirmation for security. Quick toggle in the header for users with permission. Server admins control who can enable the toggle. Works with the RBAC system to enforce content access policies.

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