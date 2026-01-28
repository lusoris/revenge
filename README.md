# Revenge

[![CI](https://github.com/lusoris/jellyfin-go/actions/workflows/ci.yml/badge.svg?branch=develop)](https://github.com/lusoris/jellyfin-go/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/lusoris/jellyfin-go)](https://goreportcard.com/report/github.com/lusoris/jellyfin-go)
[![Go Version](https://img.shields.io/github/go-mod/go-version/lusoris/jellyfin-go)](go.mod)
[![License](https://img.shields.io/badge/license-GPL--2.0-blue)](LICENSE)

> Modular media server with complete content isolation

A ground-up media server built in Go with a fully modular architecture. Each content type (movies, shows, music, etc.) is completely isolated with its own tables, services, and handlers.

## ✨ Features

- **11 Content Modules**: Movies, TV Shows, Music, Audiobooks, Books, Podcasts, Photos, Live TV, Collections, and Adult content (isolated schema)
- **External Transcoding**: Delegates to "Blackbeard" service - Revenge stays lightweight
- **Servarr Integration**: Radarr, Sonarr, Lidarr, Whisparr, Chaptarr as primary metadata sources
- **Modern Stack**: Go 1.25, PostgreSQL 18, Dragonfly (Redis-compatible), Typesense, River job queue
- **OIDC/SSO**: Full SSO support with external providers

## 🏗️ Architecture

| Component | Technology |
|-----------|------------|
| Language | Go 1.25+ |
| Database | PostgreSQL 18+ |
| Cache | Dragonfly (Redis-compatible) |
| Search | Typesense |
| Job Queue | River (PostgreSQL-native) |
| API | ogen (OpenAPI spec-first) |
| DI | uber-go/fx |
| Config | koanf v2 |

See [docs/ARCHITECTURE_V2.md](docs/ARCHITECTURE_V2.md) for the complete design.

## 🚀 Quick Start

### Docker Compose (Recommended)

```bash
git clone https://github.com/lusoris/jellyfin-go.git
cd jellyfin-go
docker compose up -d

# Opens at http://localhost:8096
```

### Development

```bash
# Prerequisites: Go 1.25+, PostgreSQL, Dragonfly
git clone https://github.com/lusoris/jellyfin-go.git
cd jellyfin-go

# Start dependencies
docker compose -f docker-compose.dev.yml up -d

# Run
go run ./cmd/revenge
```

## 📖 Documentation

### Architecture & Design
- [Architecture V2](docs/ARCHITECTURE_V2.md) - Complete modular architecture
- [Tech Stack](docs/TECH_STACK.md) - Technology choices
- [Project Structure](docs/PROJECT_STRUCTURE.md) - Directory layout

### Features
- [Metadata System](docs/METADATA_SYSTEM.md) - Servarr-first metadata with fallback
- [Audio Streaming](docs/AUDIO_STREAMING.md) - Progress tracking, bandwidth adaptation
- [Media Enhancements](docs/MEDIA_ENHANCEMENTS.md) - Trailers, themes, intros, trickplay, Live TV
- [Client Support](docs/CLIENT_SUPPORT.md) - Chromecast, DLNA, device capabilities
- [Scrobbling](docs/SCROBBLING.md) - Trakt, Last.fm, ListenBrainz sync

### Operations
- [Setup Guide](docs/SETUP.md) - Production deployment
- [Development Guide](docs/DEVELOPMENT.md) - Development environment
- [Best Practices](docs/BEST_PRACTICES.md) - Resilience, observability patterns

## 🛠️ Development

```bash
# Build
go build -o bin/revenge ./cmd/revenge

# Test
go test ./...

# Lint
golangci-lint run

# Generate (after schema/query changes)
sqlc generate
go generate ./api/...
```

## 📋 Project Status

See [TODO.md](TODO.md) for current implementation status.

**Current Phase**: Core Infrastructure

## 🤝 Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

1. Fork the repository
2. Create a feature branch
3. Write tests
4. Run `go test ./...` and `golangci-lint run`
5. Open a Pull Request

## 📜 License

GNU General Public License v2.0 - see [LICENSE](LICENSE).

## 📞 Contact

- **Issues**: [GitHub Issues](https://github.com/lusoris/jellyfin-go/issues)
- **Security**: See [SECURITY.md](SECURITY.md)
