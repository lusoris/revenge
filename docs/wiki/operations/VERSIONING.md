## Table of Contents

- [Semantic Versioning & Releases](#semantic-versioning-releases)
  - [Contents](#contents)
  - [How It Works](#how-it-works)
  - [Features](#features)
  - [Configuration](#configuration)
  - [Related Documentation](#related-documentation)
    - [See Also](#see-also)



---
sources:
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
  - name: Semantic Versioning
    url: ../sources/standards/semver.md
    note: Auto-resolved from semver
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

# Semantic Versioning & Releases




> How version numbers work and when releases happen

Revenge follows Semantic Versioning (semver). Version numbers are formatted as vMAJOR.MINOR.PATCH. MAJOR bumps indicate breaking changes, MINOR adds new features without breaking existing functionality, and PATCH is for bug fixes. Releases are automated via Release Please - when commits merge to main, a release PR is created automatically with updated changelog. Pre-release versions (alpha, beta, rc) are used for testing before stable releases.

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