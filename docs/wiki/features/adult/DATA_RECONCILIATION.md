## Table of Contents

- [Adult Data Reconciliation](#adult-data-reconciliation)
- [Adult Data Reconciliation](#adult-data-reconciliation)
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
  - name: River Job Queue
    url: ../../../sources/tooling/river.md
    note: Auto-resolved from river
  - name: sqlc
    url: ../../../sources/database/sqlc.md
    note: Auto-resolved from sqlc
  - name: sqlc Configuration
    url: ../../../sources/database/sqlc-config.md
    note: Auto-resolved from sqlc-config
design_refs:
  - title: 01_ARCHITECTURE
    path: ../../architecture/01_ARCHITECTURE.md
  - title: 02_DESIGN_PRINCIPLES
    path: ../../architecture/02_DESIGN_PRINCIPLES.md
  - title: 03_METADATA_SYSTEM
    path: ../../architecture/03_METADATA_SYSTEM.md
---

# Adult Data Reconciliation




# Adult Data Reconciliation

> Smart matching when metadata sources disagree

When multiple metadata sources provide conflicting information, the reconciliation system resolves differences. Uses fuzzy string matching and confidence scoring to pick the best data. Manually override any field when automatic matching is wrong. Track which source provided each piece of metadata. Particularly useful for adult content where naming is inconsistent.

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