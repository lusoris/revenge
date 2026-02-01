## Table of Contents

- [Dynamic RBAC with Casbin](#dynamic-rbac-with-casbin)
- [Dynamic RBAC with Casbin](#dynamic-rbac-with-casbin)
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
    url: ../sources/security/casbin.md
    note: Auto-resolved from casbin
  - name: Casbin Documentation
    url: ../sources/security/casbin-guide.md
    note: Auto-resolved from casbin-docs
  - name: Casbin pgx Adapter
    url: ../sources/security/casbin-pgx.md
    note: Auto-resolved from casbin-pgx-adapter
  - name: Uber fx
    url: ../sources/tooling/fx.md
    note: Auto-resolved from fx
  - name: ogen OpenAPI Generator
    url: ../sources/tooling/ogen.md
    note: Auto-resolved from ogen
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
  - name: sqlc
    url: ../sources/database/sqlc.md
    note: Auto-resolved from sqlc
  - name: sqlc Configuration
    url: ../sources/database/sqlc-config.md
    note: Auto-resolved from sqlc-config
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

# Dynamic RBAC with Casbin




# Dynamic RBAC with Casbin

> Fine-grained permissions for users, roles, and content

Define exactly what each user can do with role-based access control. Create custom roles (Admin, Moderator, Family, Guest) with specific permissions. Control access to libraries, features, and administrative functions. Permissions update instantly without requiring logout. Built on Casbin for enterprise-grade policy enforcement with PostgreSQL persistence.

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