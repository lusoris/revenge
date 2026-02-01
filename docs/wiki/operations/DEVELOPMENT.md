## Table of Contents

- [Development Environment Setup](#development-environment-setup)
  - [Contents](#contents)
  - [How It Works](#how-it-works)
  - [Features](#features)
  - [Configuration](#configuration)
  - [Related Documentation](#related-documentation)
    - [See Also](#see-also)

---
sources:
- name: Go Installation
    url: https://go.dev/doc/install
    note: Official Go installation guide
- name: air Hot Reload
    url: ../sources/go_dev_tools/air/main.md
    note: Go hot reload tool
- name: Vite
    url: https://vitejs.dev/guide/
    note: Frontend build tool
- name: PostgreSQL Downloads
    url: https://www.postgresql.org/download/
    note: PostgreSQL installation
- name: Node.js Downloads
    url: https://nodejs.org/
    note: Node.js LTS releases
- name: sqlc
    url: ../sources/database/sqlc.md
    note: SQL code generation
- name: Conventional Commits
    url: ../sources/standards/conventional-commits.md
    note: Commit message format
design_refs:
- title: operations
    path: operations/INDEX.md
- title: TECH_STACK
    path: technical/TECH_STACK.md
- title: BEST_PRACTICES
    path: operations/BEST_PRACTICES.md
- title: 00_SOURCE_OF_TRUTH
    path: 00_SOURCE_OF_TRUTH.md
---

# Development Environment Setup

> Get your local development environment running in minutes

This guide walks you through setting up Revenge for local development. You will install the required tools (Go 1.25+, Node.js 20, PostgreSQL 18), configure the database, and start the development servers with hot reload. The backend uses air for instant rebuilds on code changes, while the frontend uses Vite for blazing-fast HMR. Includes troubleshooting tips for common issues and useful Makefile commands to streamline your workflow.

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
