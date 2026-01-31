# Navigation Map

> Comprehensive navigation hub for design documentation
>
> Auto-generated: 2026-01-31 01:53
>
> **Master Reference**: [00_SOURCE_OF_TRUTH.md](00_SOURCE_OF_TRUTH.md)



<!-- TOC-START -->

## Table of Contents

- [Status](#status)
- [Quick Start](#quick-start)
- [Categories Overview](#categories-overview)
- [🏗️ Architecture](#architecture)
  - [Documents](#documents)
- [✨ Features](#features)
  - [Subcategories](#subcategories)
- [⚙️ Services](#services)
  - [Documents](#documents)
- [🔌 Integrations](#integrations)
  - [Subcategories](#subcategories)
- [🔧 Technical](#technical)
  - [Documents](#documents)
- [🚀 Operations](#operations)
  - [Documents](#documents)
- [🔬 Research](#research)
  - [Documents](#documents)
- [📋 Planning](#planning)
  - [Documents](#documents)
- [Deep Directory Shortcuts](#deep-directory-shortcuts)
- [Cross-References](#cross-references)
- [Sources & Cross-References](#sources-cross-references)
  - [Cross-Reference Indexes](#cross-reference-indexes)
  - [Referenced Sources](#referenced-sources)

<!-- TOC-END -->

## Status

| Dimension | Status | Notes |
|-----------|--------|-------|
| Design | 🔴 |  |
| Sources | 🔴 |  |
| Instructions | 🔴 |  |
| Code | 🔴 |  |
| Linting | 🔴 |  |
| Unit Testing | 🔴 |  |
| Integration Testing | 🔴 |  |

---

---

## Quick Start

| If you want to... | Start here |
|-------------------|------------|
| Understand the system | [Architecture](architecture/INDEX.md) → [01_ARCHITECTURE.md](architecture/01_ARCHITECTURE.md) |
| Learn about a content type | [Features](features/INDEX.md) → Pick a module |
| Implement a service | [Services](services/INDEX.md) → Pick a service |
| Add an integration | [Integrations](integrations/INDEX.md) → Pick a provider |
| Deploy the system | [Operations](operations/INDEX.md) → [SETUP.md](operations/SETUP.md) |
| Check package versions | [00_SOURCE_OF_TRUTH.md](00_SOURCE_OF_TRUTH.md) |

---

## Categories Overview

| Category | Description | Docs |
|----------|-------------|------|
| 🏗️ [Architecture](#architecture) | System design, principles, and core architecture | 5 |
| ✨ [Features](#features) | Content modules and feature specifications | 35 |
| ⚙️ [Services](#services) | Backend services and business logic | 15 |
| 🔌 [Integrations](#integrations) | External APIs, metadata providers, and Arr stack | 58 |
| 🔧 [Technical](#technical) | API design, frontend, and configuration | 6 |
| 🚀 [Operations](#operations) | Deployment, setup, and best practices | 7 |
| 🔬 [Research](#research) | User research and UX/UI resources | 2 |
| 📋 [Planning](#planning) | Project planning and versioning | 1 |

---

## 🏗️ Architecture

> System design, principles, and core architecture

**Index**: [architecture/INDEX.md](architecture/INDEX.md)

### Documents

| Document | Status | Description |
|----------|--------|-------------|
| [Plugin Architecture Decision](architecture/05_PLUGIN_ARCHITECTURE_DECISION.md) | ✅ | Should Revenge use plugins or native integration? |
| [Revenge - Architecture v2](architecture/01_ARCHITECTURE.md) | ✅ | Complete modular architecture for a next-generation media se... |
| [Revenge - Design Principles](architecture/02_DESIGN_PRINCIPLES.md) | ✅ | Non-negotiable architecture principles for the entire projec... |
| [Revenge - Metadata System](architecture/03_METADATA_SYSTEM.md) | ✅ | Servarr-first metadata with intelligent fallback and multi-l... |
| [Revenge - Player Architecture](architecture/04_PLAYER_ARCHITECTURE.md) | ✅ | Unified web player for video and audio with native streaming... |

---

## ✨ Features

> Content modules and feature specifications

**Index**: [features/INDEX.md](features/INDEX.md)

### Subcategories

| Subcategory | Documents |
|-------------|-----------|
| [Adult Content Module](features/adult/INDEX.md) | 6 |
| [Audiobook Features](features/audiobook/INDEX.md) | 2 |
| [Book Features](features/book/INDEX.md) | 2 |
| [Comics Module](features/comics/INDEX.md) | 2 |
| [Live TV Module](features/livetv/INDEX.md) | 2 |
| [Music Features](features/music/INDEX.md) | 2 |
| [Photos Module](features/photos/INDEX.md) | 2 |
| [Playback Features](features/playback/INDEX.md) | 7 |
| [Podcasts Module](features/podcasts/INDEX.md) | 2 |
| [Shared Features](features/shared/INDEX.md) | 16 |
| [Video Module](features/video/INDEX.md) | 3 |

---

## ⚙️ Services

> Backend services and business logic

**Index**: [services/INDEX.md](services/INDEX.md)

### Documents

| Document | Status | Description |
|----------|--------|-------------|
| [API Keys Service](services/APIKEYS.md) | ✅ | API key generation, validation, and management |
| [Activity Service](services/ACTIVITY.md) | ✅ | Audit logging and event tracking |
| [Analytics Service](services/ANALYTICS.md) | ✅ | Usage analytics, playback statistics, and library insights |
| [Auth Service](services/AUTH.md) | ✅ | Authentication, registration, and password management |
| [Fingerprint Service](services/FINGERPRINT.md) | ✅ | Media file identification via perceptual hashing and acousti... |
| [Grants Service](services/GRANTS.md) | ✅ | Polymorphic resource access grants for fine-grained sharing |
| [Library Service](services/LIBRARY.md) | ✅ | Library management and access control |
| [Metadata Service](services/METADATA.md) | ✅ | External metadata providers for media enrichment |
| [Notification Service](services/NOTIFICATION.md) | ✅ | Multi-channel notifications for users and admins |
| [OIDC Service](services/OIDC.md) | ✅ | OpenID Connect / SSO provider management |
| [RBAC Service](services/RBAC.md) | ✅ | Role-based access control with Casbin |
| [Search Service](services/SEARCH.md) | ✅ | Full-text search via Typesense with per-module collections |
| [Session Service](services/SESSION.md) | ✅ | Session token management and device tracking |
| [Settings Service](services/SETTINGS.md) | ✅ | Server settings persistence and retrieval |
| [User Service](services/USER.md) | ✅ | User account management and authentication |

---

## 🔌 Integrations

> External APIs, metadata providers, and Arr stack

**Index**: [integrations/INDEX.md](integrations/INDEX.md)

### Subcategories

| Subcategory | Documents |
|-------------|-----------|
| [Anime Integration](integrations/anime/INDEX.md) | 4 |
| [Audiobook & Podcast Module](integrations/audiobook/INDEX.md) | 1 |
| [Authentication Providers](integrations/auth/INDEX.md) | 5 |
| [Casting Protocols](integrations/casting/INDEX.md) | 3 |
| [External Services](integrations/external/INDEX.md) | 1 |
| [Infrastructure Components](integrations/infrastructure/INDEX.md) | 5 |
| [Live TV Integration](integrations/livetv/INDEX.md) | 4 |
| [Metadata Providers](integrations/metadata/INDEX.md) | 1 |
| [Scrobbling Services](integrations/scrobbling/INDEX.md) | 6 |
| [Servarr Stack](integrations/servarr/INDEX.md) | 6 |
| [Transcoding Services](integrations/transcoding/INDEX.md) | 2 |
| [Wiki Integration](integrations/wiki/INDEX.md) | 4 |

---

## 🔧 Technical

> API design, frontend, and configuration

**Index**: [technical/INDEX.md](technical/INDEX.md)

### Documents

| Document | Status | Description |
|----------|--------|-------------|
| [API Reference](technical/API.md) | ⚪ |  |
| [Configuration Reference](technical/CONFIGURATION.md) | ✅ | Complete configuration options for Revenge |
| [Revenge - Advanced Offloading Architecture](technical/OFFLOADING.md) | ⚪ | Keep only essential services hot, offload everything else wi... |
| [Revenge - Audio Streaming & Progress Tracking](technical/AUDIO_STREAMING.md) | ✅ | Complete audio streaming with progress persistence and sessi... |
| [Revenge - Frontend Architecture](technical/FRONTEND.md) | ⚪ | Modern, responsive web interface with full RBAC and theme su... |
| [Revenge - Technology Stack](technical/TECH_STACK.md) | ⚪ |  |

---

## 🚀 Operations

> Deployment, setup, and best practices

**Index**: [operations/INDEX.md](operations/INDEX.md)

### Documents

| Document | Status | Description |
|----------|--------|-------------|
| [Advanced Patterns & Best Practices](operations/BEST_PRACTICES.md) | ⚪ | Comprehensive guide for professional-grade implementations i... |
| [Branch Protection Rules](operations/BRANCH_PROTECTION.md) | ⚪ |  |
| [Clone repository](operations/DEVELOPMENT.md) | ⚪ |  |
| [Database Auto-Healing & Consistency Restoration](operations/DATABASE_AUTO_HEALING.md) | ⚪ | Automatic recovery strategies for PostgreSQL database corrup... |
| [GitFlow Workflow Guide](operations/GITFLOW.md) | ⚪ |  |
| [Revenge - Reverse Proxy & Deployment Best Practices](operations/REVERSE_PROXY.md) | ⚪ | Production deployment with nginx, Caddy, Traefik, and Docker... |
| [revenge - Setup Guide](operations/SETUP.md) | ⚪ |  |

---

## 🔬 Research

> User research and UX/UI resources

**Index**: [research/INDEX.md](research/INDEX.md)

### Documents

| Document | Status | Description |
|----------|--------|-------------|
| [UX/UI Design & Frontend Resources](research/UX_UI_RESOURCES.md) | ⚪ | **Status**: ✅ ALL 14 SOURCES FETCHED (2026-01-28) |
| [User Pain Points Research - Existing Media Servers](research/USER_PAIN_POINTS_RESEARCH.md) | ⚪ | Analysis of what users most complain about (and love) in Ple... |

---

## 📋 Planning

> Project planning and versioning

**Index**: [planning/INDEX.md](planning/INDEX.md)

### Documents

| Document | Status | Description |
|----------|--------|-------------|
| [Versioning Strategy](operations/VERSIONING.md) | ⚪ |  |

---

## Deep Directory Shortcuts

> Direct links to deeply nested documentation (depth 3+)

| Path | Description |
|------|-------------|
| [integrations/metadata/adult/](integrations/metadata/adult/INDEX.md) | Adult metadata providers (StashDB, TPDB) |
| [integrations/metadata/books/](integrations/metadata/books/INDEX.md) | Book metadata (OpenLibrary, Audible) |
| [integrations/metadata/comics/](integrations/metadata/comics/INDEX.md) | Comics metadata (ComicVine) |
| [integrations/metadata/music/](integrations/metadata/music/INDEX.md) | Music metadata (MusicBrainz, Last.fm) |
| [integrations/metadata/video/](integrations/metadata/video/INDEX.md) | Video metadata (TMDb, TheTVDB) |
| [integrations/wiki/adult/](integrations/wiki/adult/INDEX.md) | Adult wiki sources (IAFD, Babepedia) |
| [integrations/metadata/adult/](integrations/metadata/adult/INDEX.md) | Adult external sources |

---

## Cross-References

| Resource | Description |
|----------|-------------|
| [00_SOURCE_OF_TRUTH.md](00_SOURCE_OF_TRUTH.md) | Package versions, module status, config keys |
| [DESIGN_INDEX.md](../sources/DESIGN_INDEX.md) | Auto-generated index of all design docs |
| [SOURCES_INDEX.md](../sources/SOURCES_INDEX.md) | Index of external documentation sources |
| [DESIGN_CROSSREF.md](../sources/DESIGN_CROSSREF.md) | Design ↔ Sources cross-reference map |


<!-- SOURCE-BREADCRUMBS-START -->

## Sources & Cross-References

> Auto-generated section linking to external documentation sources

### Cross-Reference Indexes

- [All Sources Index](../sources/SOURCES_INDEX.md) - Complete list of external documentation
- [Design ↔ Sources Map](../sources/DESIGN_CROSSREF.md) - Which docs reference which sources

### Referenced Sources

| Source | Documentation |
|--------|---------------|
| [Casbin](https://pkg.go.dev/github.com/casbin/casbin/v2) | [Local](../sources/security/casbin.md) |
| [Last.fm API](https://www.last.fm/api/intro) | [Local](../sources/apis/lastfm.md) |
| [PostgreSQL Arrays](https://www.postgresql.org/docs/current/arrays.html) | [Local](../sources/database/postgresql-arrays.md) |
| [PostgreSQL JSON Functions](https://www.postgresql.org/docs/current/functions-json.html) | [Local](../sources/database/postgresql-json.md) |
| [Typesense API](https://typesense.org/docs/latest/api/) | [Local](../sources/infrastructure/typesense.md) |
| [Typesense Go Client](https://github.com/typesense/typesense-go) | [Local](../sources/infrastructure/typesense-go.md) |
| [pgx PostgreSQL Driver](https://pkg.go.dev/github.com/jackc/pgx/v5) | [Local](../sources/database/pgx.md) |

<!-- SOURCE-BREADCRUMBS-END -->

---

*This file is auto-generated by `scripts/generate-navigation-map.py`*
