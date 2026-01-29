# Revenge

[![CI](https://github.com/lusoris/revenge/actions/workflows/ci.yml/badge.svg?branch=develop)](https://github.com/lusoris/revenge/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/lusoris/revenge)](https://goreportcard.com/report/github.com/lusoris/revenge)
[![Go Version](https://img.shields.io/github/go-mod-go-version/lusoris/revenge)](go.mod)
[![License](https://img.shields.io/badge/license-AGPL--3.0-blue)](LICENSE)

> Modular media server with complete content isolation

A ground-up media server built in Go with a fully modular architecture. Each content type (movies, shows, music, etc.) is completely isolated with its own tables, services, and handlers.

**WARNING: EARLY DEVELOPMENT** - Core infrastructure ready, content modules in development. See [TODO.md](TODO.md) for roadmap.

---

## Design Principles

- **Performance First** - UX never blocked by backend tasks
- **Client Agnostic** - Native mobile/TV apps with Revenge API, DLNA compatibility
- **Privacy by Default** - Encrypted local storage, opt-in tracking only
- **Bleeding Edge Stable** - Latest stable Go/PostgreSQL, no alpha deps
- **Optional ML** - Ollama integration for recommendations, not required

See [DESIGN_PRINCIPLES.md](docs/dev/design/architecture/DESIGN_PRINCIPLES.md) for full details.

---

## Features

### Content Modules
- **12 Content Types**: Movies, TV Shows, Music, Audiobooks, Books, Podcasts, Photos, Live TV, Comics, Collections, and Adult content (isolated `c` schema)
- **Module-Specific Age Restrictions**: Separate rating systems per module (video MPAA/BBFC, music parental advisory, books age ranges, comics publisher ratings)

### Playback
- **Watch Next / Continue Watching**: Smart progress tracking with 5%-90% thresholds, series navigation, cross-device sync
- **External Transcoding**: Delegates to "Blackbeard" service - Revenge stays lightweight
- **Audio Streaming**: Gapless playback, bandwidth adaptation, format negotiation

### Integrations
- **Servarr Integration**: Radarr, Sonarr, Lidarr, Whisparr v3 (eros), Chaptarr as primary metadata sources
- **Release Calendar**: Unified calendar from all Servarr instances with upcoming/recent views
- **Scrobbling**: Trakt, Last.fm, ListenBrainz, Letterboxd sync
- **OIDC/SSO**: Full SSO support with external providers (Authelia, Authentik, Keycloak)

---

## Architecture

| Component   | Technology          | Notes                            |
|-------------|---------------------|----------------------------------|
| Language    | Go 1.25+            | `GOEXPERIMENT=greenteagc,jsonv2` |
| Database    | PostgreSQL 18+      | sqlc for type-safe queries       |
| Cache       | Dragonfly + rueidis | Redis-compatible, auto-pipelining|
| Local Cache | otter v1.2.4        | W-TinyLFU eviction               |
| Search      | Typesense           | typesense-go/v4 client           |
| Job Queue   | River               | PostgreSQL-native                |
| API         | ogen                | OpenAPI spec-first               |
| DI          | uber-go/fx          | Dependency injection             |
| Config      | koanf v2            | Multi-source configuration       |
| Logging     | slog                | Structured logging               |
| Resilience  | failsafe-go         | Circuit breakers, retries        |

See [ARCHITECTURE_V2.md](docs/dev/design/architecture/ARCHITECTURE_V2.md) for the complete design.

---

## Quick Start

### Docker Compose (Recommended)

```bash
git clone https://github.com/lusoris/revenge.git
cd revenge
docker compose up -d

# Opens at http://localhost:8096
```

### Development

```bash
# Prerequisites: Go 1.25+, PostgreSQL, Dragonfly
git clone https://github.com/lusoris/revenge.git
cd revenge

# Start dependencies
docker compose -f docker-compose.dev.yml up -d

# Run with experiments enabled
GOEXPERIMENT=greenteagc,jsonv2 go run ./cmd/revenge
```

---

## Documentation

All documentation lives in [docs/dev/](docs/dev/INDEX.md).

### Core Design
- [Architecture V2](docs/dev/design/architecture/ARCHITECTURE_V2.md) - Complete modular architecture
- [Tech Stack](docs/dev/design/technical/TECH_STACK.md) - Technology choices with rationale
- [Design Principles](docs/dev/design/architecture/DESIGN_PRINCIPLES.md) - Guiding principles
- [Configuration Reference](docs/dev/design/technical/CONFIGURATION.md) - koanf configuration options

### Services
- [Services Index](docs/dev/design/services/INDEX.md) - Service layer overview
- [Auth Service](docs/dev/design/services/AUTH.md) - Authentication and login flows
- [User Service](docs/dev/design/services/USER.md) - User management and roles
- [Session Service](docs/dev/design/services/SESSION.md) - Token and session management
- [Library Service](docs/dev/design/services/LIBRARY.md) - Library management and scanning
- [RBAC Service](docs/dev/design/services/RBAC.md) - Casbin role-based access control
- [OIDC Service](docs/dev/design/services/OIDC.md) - SSO provider integration
- [Settings Service](docs/dev/design/services/SETTINGS.md) - Server settings persistence

### Content Modules
- [Movie Module](docs/dev/design/features/video/MOVIE_MODULE.md) - Movie content management
- [TV Show Module](docs/dev/design/features/video/TVSHOW_MODULE.md) - Series, seasons, episodes

### Features
- [Metadata System](docs/dev/design/architecture/METADATA_SYSTEM.md) - Servarr-first metadata with fallback providers
- [Watch Next / Continue Watching](docs/dev/design/features/playback/WATCH_NEXT_CONTINUE_WATCHING.md) - Playback continuation system
- [Release Calendar](docs/dev/design/features/playback/RELEASE_CALENDAR.md) - Upcoming releases via Servarr
- [Request System](docs/dev/design/features/shared/REQUEST_SYSTEM.md) - Content requests with polls, voting, RBAC rules
- [Content Rating](docs/dev/design/features/shared/CONTENT_RATING.md) - Module-specific age restriction systems
- [Audio Streaming](docs/dev/design/technical/AUDIO_STREAMING.md) - Progress tracking, bandwidth adaptation
- [Scrobbling](docs/dev/design/features/shared/SCROBBLING.md) - Trakt, Last.fm, ListenBrainz sync
- [RBAC with Casbin](docs/dev/design/features/shared/RBAC_CASBIN.md) - Dynamic role-based access control

### Integrations
- [Servarr Index](docs/dev/design/integrations/servarr/INDEX.md) - Radarr, Sonarr, Lidarr, Whisparr, Chaptarr
- [Metadata Providers](docs/dev/design/integrations/metadata/INDEX.md) - TMDb, MusicBrainz, AniList, etc.
- [Authentication Providers](docs/dev/design/integrations/auth/INDEX.md) - Authelia, Authentik, Keycloak

### Operations
- [Setup Guide](docs/dev/design/operations/SETUP.md) - Production deployment
- [Development Guide](docs/dev/design/operations/DEVELOPMENT.md) - Development environment
- [Best Practices](docs/dev/design/operations/BEST_PRACTICES.md) - Resilience, observability patterns

### Development Instructions
- [Instructions Index](.github/instructions/INDEX.instructions.md) - AI-assisted development guidelines

---

## Development

```bash
# Build
GOEXPERIMENT=greenteagc,jsonv2 go build -o bin/revenge ./cmd/revenge

# Test
go test ./...

# Lint
golangci-lint run

# Generate (after schema/query changes)
sqlc generate
go generate ./api/...
```

---

## Project Status

See [TODO.md](TODO.md) for current implementation status.

**Current Phase**: Content Modules

---

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

1. Fork the repository
2. Create a feature branch
3. Write tests
4. Run `go test ./...` and `golangci-lint run`
5. Open a Pull Request

---

## License

GNU Affero General Public License v3.0 - see [LICENSE](LICENSE).

---

## Contact

- **Issues**: [GitHub Issues](https://github.com/lusoris/revenge/issues)
- **Security**: See [SECURITY.md](SECURITY.md)
