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

### Features
- [Metadata System](docs/dev/design/architecture/METADATA_SYSTEM.md) - Servarr-first metadata with fallback providers
- [Watch Next / Continue Watching](docs/dev/design/features/WATCH_NEXT_CONTINUE_WATCHING.md) - Playback continuation system
- [Release Calendar](docs/dev/design/features/RELEASE_CALENDAR.md) - Upcoming releases via Servarr
- [Content Rating](docs/dev/design/features/CONTENT_RATING.md) - Module-specific age restriction systems
- [Audio Streaming](docs/dev/design/technical/AUDIO_STREAMING.md) - Progress tracking, bandwidth adaptation
- [Scrobbling](docs/dev/design/features/SCROBBLING.md) - Trakt, Last.fm, ListenBrainz sync

### Integrations
- [Servarr Index](docs/dev/design/integrations/servarr/INDEX.md) - Radarr, Sonarr, Lidarr, Whisparr, Chaptarr
- [Metadata Providers](docs/dev/design/integrations/metadata/INDEX.md) - TMDb, MusicBrainz, AniList, etc.

### Operations
- [Setup Guide](docs/dev/design/operations/SETUP.md) - Production deployment
- [Development Guide](docs/dev/design/operations/DEVELOPMENT.md) - Development environment
- [Best Practices](docs/dev/design/operations/BEST_PRACTICES.md) - Resilience, observability patterns

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
