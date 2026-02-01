## Table of Contents

- [RBAC Service](#rbac-service)
- [RBAC Service](#rbac-service)
  - [Contents](#contents)
  - [How It Works](#how-it-works)
  - [Features](#features)
  - [Configuration](#configuration)
  - [Related Documentation](#related-documentation)
    - [See Also](#see-also)

---
sources:
- name: Casbin
    url: ../sources/security/casbin.md
    note: Auto-resolved from casbin
- name: Uber fx
    url: ../sources/tooling/fx.md
    note: Auto-resolved from fx
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
design_refs:
- title: services
    path: services/INDEX.md
- title: 01_ARCHITECTURE
    path: architecture/01_ARCHITECTURE.md
- title: 02_DESIGN_PRINCIPLES
    path: architecture/02_DESIGN_PRINCIPLES.md
- title: 03_METADATA_SYSTEM
    path: architecture/03_METADATA_SYSTEM.md
---

# RBAC Service

# RBAC Service

> Role-based access control with Casbin

The RBAC service controls who can do what in Revenge. Built on Casbin, it supports roles (Admin, User, Guest) and fine-grained permissions. Admins define policies like "Users can view movies" or "Guests cannot access adult content". Policies are stored in PostgreSQL and cached for fast authorization checks on every request.

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
