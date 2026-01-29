# Revenge Documentation Index

> Central navigation for all project documentation

**Last Updated**: 2026-01-29

## Quick Start

- **New to Revenge?** Start with [Setup Guide](operations/SETUP.md)
- **Development?** See [Development Guide](operations/DEVELOPMENT.md)
- **Architecture?** Read [Overall Architecture](architecture/ARCHITECTURE_V2.md)
- **External APIs?** Check [External Integrations](EXTERNAL_INTEGRATIONS_TODO.md)
- **Current Status?** See [Analysis Reports](#analysis-reports)

---

## Analysis Reports

**Current codebase analysis and status reports.**

- [Architecture Compliance Analysis](../ARCHITECTURE_COMPLIANCE_ANALYSIS.md) - 65% compliance score
- [Advanced Features Integration](../ADVANCED_FEATURES_INTEGRATION_ANALYSIS.md) - 10% integration score
- [Core Functionality Analysis](../CORE_FUNCTIONALITY_ANALYSIS.md) - Missing workers, services, migrations
- [Documentation Cleanup Report](../DOCUMENTATION_CLEANUP_REPORT.md) - Archived outdated TODOs

**Action Items**: See [`TODO.md`](../TODO.md)

---

## Architecture

Core system design and architectural decisions.

- [Overall Architecture V2](architecture/ARCHITECTURE_V2.md) - Complete system architecture
- [Design Principles](architecture/DESIGN_PRINCIPLES.md) - Guiding principles
- [Player Architecture](architecture/PLAYER_ARCHITECTURE.md) - Media player design
- [Metadata System](architecture/METADATA_SYSTEM.md) - Metadata management
- [Plugin Architecture Decision](architecture/PLUGIN_ARCHITECTURE_DECISION.md) - Native vs plugins

## Operations

Setup, deployment, and operational guides.

- [Setup Guide](operations/SETUP.md) - Initial setup instructions
- [Development Guide](operations/DEVELOPMENT.md) - Development workflow
- [Database Auto-Healing](operations/DATABASE_AUTO_HEALING.md) - PostgreSQL corruption detection/repair
- [Reverse Proxy](operations/REVERSE_PROXY.md) - Nginx/Caddy configuration
- [Git Workflow](operations/GITFLOW.md) - Branch strategy
- [Branch Protection](operations/BRANCH_PROTECTION.md) - Repository protection rules
- [Best Practices](operations/BEST_PRACTICES.md) - Coding standards
- [Upstream Sync](operations/UPSTREAM_SYNC.md) - Syncing with upstream

## Research

Analysis, user research, and technology evaluations.

- [User Pain Points Research](research/USER_PAIN_POINTS_RESEARCH.md) - Jellyfin/Plex/Emby issues
- [Go Packages Research](research/GO_PACKAGES_RESEARCH.md) - Awesome-Go analysis (500+ packages)

**Archived**: See `archive/reports/` for historical analysis (2026-01-28 snapshots)

**Current Analysis** (2026-01-29):
- [Architecture Compliance](../ARCHITECTURE_COMPLIANCE_ANALYSIS.md) - 65% compliance
- [Advanced Features Integration](../ADVANCED_FEATURES_INTEGRATION_ANALYSIS.md) - 10% integration
- [Core Functionality Analysis](../CORE_FUNCTIONALITY_ANALYSIS.md) - Missing workers/services
- [Documentation Cleanup Report](../DOCUMENTATION_CLEANUP_REPORT.md) - Archived 264+ outdated TODOs
- [Design TODOs Extraction](../DESIGN_TODOS_EXTRACTION.md) - 100+ missing components

## Planning

Roadmaps, versioning, and implementation tracking.

- [VERSION_POLICY.md](planning/VERSION_POLICY.md) - Semantic versioning policy
- [VERSIONING.md](planning/VERSIONING.md) - Version management strategy

**Current Progress**: See [`TODO.md`](../TODO.md) for active tasks

**Archived Planning**: See `archive/planning/` for historical roadmaps

- [Module Implementation Roadmap](planning/MODULE_IMPLEMENTATION_TODO.md) - Module rollout plan
- [Version Policy](planning/VERSION_POLICY.md) - Bleeding edge/latest stable policy
- [Versioning Strategy](planning/VERSIONING.md) - Semantic versioning approach
- [Preparation Master Plan](PREPARATION_MASTER_PLAN.md) - Pre-implementation checklist (2737 lines, SPLIT PENDING)

## Integrations

External service integrations and API documentation.

- [External Integrations TODO](EXTERNAL_INTEGRATIONS_TODO.md) - 66 services (1103 lines, SPLIT PENDING)
- **Servarr**: Radarr, Sonarr, Lidarr, Whisparr, Readarr
- **Metadata**: TMDb, TheTVDB, MusicBrainz, StashDB, ComicVine
- **Scrobbling**: Trakt, Last.fm, ListenBrainz, Letterboxd, Simkl
- **Auth**: Authelia, Authentik, Keycloak, OIDC

## Features

Content modules and feature designs.

- [Adult Content System](features/ADULT_CONTENT_SYSTEM.md) - Adult module isolation (`c` schema)
- [Adult Metadata](features/ADULT_METADATA.md) - Whisparr/StashDB integration
- [Comics Module](features/COMICS_MODULE.md) - Digital comics/manga support
- [Scrobbling](features/SCROBBLING.md) - Trakt/Last.fm/ListenBrainz
- [User Experience Features](features/USER_EXPERIENCE_FEATURES.md) - UX enhancements
- [Media Enhancements](features/MEDIA_ENHANCEMENTS.md) - Media-specific features
- [Content Rating](features/CONTENT_RATING.md) - Rating systems (MPAA/PEGI/etc.)
- [Client Support](features/CLIENT_SUPPORT.md) - Client compatibility
- [Internationalization](features/I18N.md) - i18n/l10n support
- [Library Types](features/LIBRARY_TYPES.md) - Library management

## Technical

API design, frontend architecture, and technical documentation.

- [API Documentation](technical/API.md) - REST API design
- [Frontend Architecture](technical/FRONTEND.md) - Svelte 5 + shadcn-svelte
- [Tech Stack](technical/TECH_STACK.md) - Technologies used
- [Audio Streaming](technical/AUDIO_STREAMING.md) - Audio streaming architecture
- [Offloading](technical/OFFLOADING.md) - Task offloading patterns

---

## Documentation Status

| Category | Files | Status |
|----------|-------|--------|
| Architecture | 5 | ✅ Complete |
| Operations | 8 | ✅ Complete |
| Research | 3 | 🟡 Plex/Emby data pending |
| Planning | 4 | 🟡 Split large files pending |
| Integrations | 2 | 🔴 Split EXTERNAL_INTEGRATIONS_TODO.md |
| Features | 10 | 🟡 4 more docs needed |
| Technical | 5 | ✅ Complete |

**Total**: 37 files | **Next**: Split PREPARATION_MASTER_PLAN.md, EXTERNAL_INTEGRATIONS_TODO.md

---

## Related Resources

- [Agent Instructions](../AGENTS.md) - Automated coding agent rules
- [Copilot Instructions](../.github/copilot-instructions.md) - GitHub Copilot rules
- [Instruction Files](../.github/instructions/) - 23 pattern-specific instructions
- [TODO List](../TODO.md) - Project backlog
- [Contributing Guide](../CONTRIBUTING.md) - Contribution guidelines
