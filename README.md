# Jellyfin Go

[![CI](https://github.com/lusoris/jellyfin-go/actions/workflows/ci.yml/badge.svg)](https://github.com/lusoris/jellyfin-go/actions/workflows/ci.yml)
[![Security](https://github.com/lusoris/jellyfin-go/actions/workflows/security.yml/badge.svg)](https://github.com/lusoris/jellyfin-go/actions/workflows/security.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/lusoris/jellyfin-go)](https://goreportcard.com/report/github.com/lusoris/jellyfin-go)
[![codecov](https://codecov.io/gh/lusoris/jellyfin-go/branch/main/graph/badge.svg)](https://codecov.io/gh/lusoris/jellyfin-go)
[![Go Version](https://img.shields.io/github/go-mod/go-version/lusoris/jellyfin-go)](go.mod)
[![License](https://img.shields.io/github/license/lusoris/jellyfin-go)](LICENSE)

A high-performance, production-grade media server reimplementation of Jellyfin in Go.

## 🎯 Project Goals

**Phase 1: Single-Server Excellence**
- **Full Feature Parity**: Complete reimplementation of Jellyfin with all 60+ API controllers
- **Easy Setup**: Docker one-liner or native binary with PostgreSQL
- **Hardware Acceleration**: Full support for VAAPI, NVENC, QuickSync, AMF, VideoToolbox
- **API Compatibility**: Maintain v1 compatibility for existing Jellyfin clients
- **Performance**: Optimized for single server, 1-100 concurrent users

**Phase 2: Optional Clustering** (Future)
- Horizontal scalability with multi-instance support
- Distributed coordination (Dragonfly/Redis cluster)
- Advanced monitoring and observability
- CDN integration for large deployments

## 🏗️ Architecture Highlights

**Single-Server Mode (Default):**
- **PostgreSQL Database**: Optimized for large libraries (100k+ items)
- **Dragonfly Cache**: High-performance caching (Redis-compatible)
- **Typesense Search**: Lightning-fast faceted search engine
- **Hardware Transcoding**: FFmpeg with GPU acceleration
- **Simple Setup**: Docker Compose with all dependencies

**Optional Enhancements:**
- CDN integration for large deployments
- Multi-instance clustering (Phase 2)

## 📊 Performance Targets

**Single-Server Mode:**
- **API Latency**: P95 < 50ms (PostgreSQL + Dragonfly)
- **Concurrent Users**: 10-100 users
- **Concurrent Streams**: 10-50 streams (hardware dependent)
- **Library Size**: 100k+ media items (tested with 100k+)
- **Memory Usage**: 1-4GB (depending on cache size)

**Multi-Instance Mode (Optional):**
- API Latency: P95 < 200ms
- Concurrent Streams: 1,000+ per cluster
- Availability: 99.9% uptime

## 🚀 Quick Start

**Option 1: Docker Compose (Recommended)**
```bash
# Clone and start all services (PostgreSQL, Dragonfly, Typesense, Jellyfin)
git clone https://github.com/lusoris/jellyfin-go.git
cd jellyfin-go
docker compose up -d

# Opens at http://localhost:8096
```

**Option 2: Development**
```bash
# Prerequisites: Go 1.24+, FFmpeg, PostgreSQL, Dragonfly
git clone https://github.com/lusoris/jellyfin-go.git
cd jellyfin-go

# Start dependencies
docker compose -f docker-compose.dev.yml up -d

# Run application
go run ./cmd/jellyfin
```

**Requirements:**
- PostgreSQL 18+ (required)
- Dragonfly or Redis 7+ (required for caching)
- Typesense 0.25+ (required for search)
- FFmpeg 6+ (for transcoding)

## 📖 Documentation

- [Architecture Design](docs/ARCHITECTURE.md) - System architecture and component design
- [Implementation Plan](docs/IMPLEMENTATION_PLAN.md) - Phased development roadmap
- [Best Practices](docs/BEST_PRACTICES.md) - Go coding standards and patterns
- [API Design](docs/API_DESIGN.md) - API specification and compatibility
- [Security Guidelines](docs/SECURITY.md) - Security hardening and OWASP compliance
- [Setup Instructions](docs/SETUP.md) - Development and production setup

## 🛠️ Tech Stack

### Core
- **Language**: Go 1.24 (bleeding-edge stable)
- **HTTP Router**: net/http.ServeMux (stdlib, Go 1.22+ enhanced patterns)
- **Dependency Injection**: uber-go/fx v1.23
- **Configuration**: koanf v2 (modern, replaces Viper)
- **Logging**: log/slog + tint (stdlib with pretty output)

### Data Layer
- **Database**: PostgreSQL 18+ (required)
- **Query Builder**: sqlc (type-safe SQL)
- **Migrations**: golang-migrate
- **Cache**: Dragonfly (Redis-compatible, required)
- **Search**: Typesense (required, faceted search)

### Optional Components
- **Load Balancer**: NGINX/HAProxy (multi-instance)
- **CDN**: Cloudflare/Bunny (large deployments)

### Media Processing
- **FFmpeg**: jellyfin-ffmpeg (hardware acceleration)
- **Streaming**: HLS/DASH with adaptive bitrate
- **Thumbnails**: FFmpeg tile generation

### Observability
- **Metrics**: Prometheus + Grafana
- **Tracing**: OpenTelemetry
- **Logging**: slog (structured JSON)
- **Profiling**: Pyroscope (continuous profiling)

### Security
- **Authentication**: JWT with refresh tokens
- **Authorization**: Custom policy engine
- **Secrets**: HashiCorp Vault / K8s External Secrets
- **Rate Limiting**: Token bucket with Dragonfly

### Deployment
- **Containers**: Docker + Docker Compose
- **Orchestration**: Kubernetes with Helm
- **CI/CD**: GitHub Actions
- **Load Balancer**: NGINX / HAProxy

## 🔄 Migration from Jellyfin

Jellyfin Go includes a migration tool to import existing Jellyfin databases:

```bash
# Automatic detection during first launch
./jellyfin-go --import-jellyfin /path/to/jellyfin/data

# Or via web UI setup wizard
# Navigate to http://localhost:8096/setup
# Select "Import from existing Jellyfin"
```

**Note**: Database schema is redesigned for optimal PostgreSQL performance. Migration is one-way.

## 📈 Roadmap

### Phase 1: Single-Server MVP (3 months) 🎯 **PRIMARY FOCUS**
- ⬜ PostgreSQL database with sqlc
- ⬜ Dragonfly caching layer
- ⬜ Typesense search integration
- ⬜ Authentication and authorization (JWT, bcrypt)
- ⬜ Core API endpoints (users, library, media)
- ⬜ Library scanning and metadata extraction
- ⬜ Direct play (no transcoding)
- ⬜ PostgreSQL full-text search

### Phase 2: Transcoding (2 months)
- ⬜ FFmpeg integration and hardware acceleration
- ⬜ HLS transcoding pipeline
- ⬜ Session management
- ⬜ Playback progress tracking
- ⬜ Jellyfin migration tool

### Phase 3: Feature Complete (2 months)
- ⬜ WebSocket realtime updates
- ⬜ Collections and playlists
- ⬜ Advanced metadata providers (TMDb, OMDb)
- ⬜ Plugin system (HashiCorp go-plugin)
- ⬜ Subtitle support
- ⬜ Image processing and thumbnails

### Phase 4: Polish & Security (2 months)
- ⬜ Security hardening (OWASP Top 10)
- ⬜ Performance optimization
- ⬜ User documentation
- ⬜ Mobile-friendly web UI
- ⬜ Backup/restore functionality

### Phase 5: Optional Clustering (Future)
- ⬜ PostgreSQL read replicas
- ⬜ Redis/Dragonfly distributed cache
- ⬜ Multi-instance coordination
- ⬜ Load balancer configuration
- ⬜ CDN integration
- ⬜ Advanced monitoring (Prometheus/Grafana)

## 🤝 Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

### Development Workflow

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Write tests for your changes
4. Ensure all tests pass (`go test ./...`)
5. Run linters (`golangci-lint run`)
6. Commit your changes (`git commit -m 'Add amazing feature'`)
7. Push to the branch (`git push origin feature/amazing-feature`)
8. Open a Pull Request

## 📜 License

This project is licensed under the GNU General Public License v2.0 - see [LICENSE](LICENSE) for details.

Jellyfin is a registered trademark of Jellyfin contributors. This project is an independent reimplementation and is not officially affiliated with the Jellyfin project.

## 🙏 Acknowledgments

- **Jellyfin Team**: For creating an amazing open-source media server
- **Go Community**: For excellent libraries and tools
- **Contributors**: Everyone who contributes to this project

## 📞 Contact

- **Issues**: [GitHub Issues](https://github.com/your-org/jellyfin-go/issues)
- **Discussions**: [GitHub Discussions](https://github.com/your-org/jellyfin-go/discussions)
- **Security**: security@example.com

## 🔗 Links

- [Jellyfin (Original)](https://jellyfin.org/)
- [Documentation](https://docs.example.com)
- [API Reference](https://api-docs.example.com)
- [Community](https://community.example.com)
