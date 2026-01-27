# Phase 1 - Detailed Implementation Todo List

> **Goal:** Single-Server MVP with PostgreSQL + OIDC
> **Timeline:** 18 weeks (~4.5 months)
> **Status:** Foundation complete, starting core features

---

## 🏗️ INFRASTRUCTURE (Week 1-2) - ✅ DONE

### 1.1 Project Setup ✅
- [x] Go 1.24 module initialization
- [x] Directory structure (cmd/, internal/, pkg/, docs/, tests/)
- [x] Makefile with common commands
- [x] .editorconfig, .gitignore, .gitattributes
- [x] LICENSE (GPL-2.0)

### 1.2 CI/CD ✅
- [x] GitHub Actions: lint workflow
- [x] GitHub Actions: test workflow
- [x] GitHub Actions: build workflow
- [x] GitHub Actions: release workflow (release-please)
- [x] Branch protection rules
- [x] PR template, issue templates

### 1.3 Development Environment ✅
- [x] Docker Compose (dev): PostgreSQL 18, Dragonfly, Typesense
- [x] Docker Compose (prod): Full stack
- [x] Dockerfile (multi-stage build)
- [x] DevContainer configuration
- [x] VS Code workspace settings
- [x] Setup scripts (bash, PowerShell, fish)

### 1.4 Configuration System ✅
- [x] koanf v2 integration
- [x] YAML config files (defaults.yaml, config.yaml)
- [x] Environment variable support (JELLYFIN_*)
- [x] Config struct with validation
- [x] Environment-specific configs (dev, prod)

### 1.5 Logging ✅
- [x] slog setup with handlers
- [x] Console output (pretty for dev)
- [x] JSON output (for prod)
- [x] Log levels (debug, info, warn, error)
- [x] Context-aware logging

### 1.6 HTTP Server ✅
- [x] net/http.ServeMux with Go 1.22+ patterns
- [x] Graceful shutdown
- [x] Health endpoints (/health/live, /health/ready)
- [x] Basic middleware setup

---

## 🗄️ DATABASE (Week 2-3)

### 2.1 PostgreSQL Setup
- [ ] Connection pool (pgxpool)
- [ ] Connection config (host, port, user, password, database, SSL)
- [ ] Health check query
- [ ] Automatic reconnection

### 2.2 Migrations (golang-migrate)
- [ ] Install migrate tool
- [ ] Migration file structure (migrations/*.sql)
- [ ] Initial schema migration (001_initial.up.sql)
- [ ] Makefile targets: migrate-up, migrate-down, migrate-create
- [ ] Migration in startup (optional auto-migrate)

### 2.3 Schema Design
```sql
-- Core tables needed:
- [ ] users (id, username, email, password_hash, is_admin, created_at, updated_at)
- [ ] sessions (id, user_id, token_hash, refresh_token_hash, expires_at, created_at)
- [ ] oidc_providers (id, name, issuer_url, client_id, client_secret_encrypted, enabled)
- [ ] oidc_user_links (id, user_id, provider_id, subject, created_at)
- [ ] libraries (id, name, type, paths, settings_json, created_at)
- [ ] media_items (id, library_id, type, name, path, metadata_json, created_at)
- [ ] images (id, item_id, type, path, width, height)
- [ ] playback_progress (user_id, item_id, position_ticks, played, updated_at)
```

### 2.4 sqlc Setup
- [ ] sqlc.yaml configuration
- [ ] Query files (queries/*.sql)
- [ ] Generated code (internal/infra/database/db/)
- [ ] Custom types mapping
- [ ] Makefile target: sqlc-generate

### 2.5 Repository Pattern
- [ ] Repository interfaces in internal/domain/
- [ ] PostgreSQL implementations in internal/infra/database/
- [ ] Transaction support
- [ ] fx module for repositories

---

## 🔐 AUTHENTICATION (Week 3-5)

### 3.1 User Entity & Repository
- [ ] User domain entity (internal/domain/user.go)
- [ ] UserRepository interface
- [ ] PostgreSQL UserRepository implementation
- [ ] User validation (email, username rules)

### 3.2 Password Security
- [ ] bcrypt hashing (cost factor 12)
- [ ] Password validation (min length, complexity)
- [ ] Constant-time comparison
- [ ] Password change support

### 3.3 JWT Implementation
- [ ] JWT library selection (golang-jwt/jwt/v5)
- [ ] Access token generation (short-lived, 15min)
- [ ] Refresh token generation (long-lived, 7 days)
- [ ] Token signing (RS256 or EdDSA)
- [ ] Key rotation support
- [ ] Token claims (sub, exp, iat, jti, roles)

### 3.4 Session Management
- [ ] Session entity
- [ ] Session repository
- [ ] Session creation on login
- [ ] Session validation middleware
- [ ] Session revocation (logout)
- [ ] Multi-device session support
- [ ] Session cleanup (expired sessions)

### 3.5 Auth Middleware
- [ ] JWT extraction from Authorization header
- [ ] JWT validation
- [ ] User context injection
- [ ] Role-based access control (RBAC)
- [ ] Admin-only middleware
- [ ] Rate limiting per user

### 3.6 OIDC Integration
- [ ] OIDC discovery (.well-known/openid-configuration)
- [ ] Provider configuration storage
- [ ] Authorization code flow implementation
- [ ] State parameter for CSRF protection
- [ ] PKCE support (code_verifier, code_challenge)
- [ ] ID token validation (signature, claims)
- [ ] User provisioning (auto-create from OIDC)
- [ ] Claim mapping (email, name, groups → roles)
- [ ] Multiple provider support
- [ ] Provider management UI endpoints

### 3.7 Auth API Endpoints
- [ ] `POST /Users/AuthenticateByName` - Local login
- [ ] `POST /Users/New` - Registration (if enabled)
- [ ] `POST /Auth/Logout` - Logout (revoke session)
- [ ] `POST /Auth/Refresh` - Refresh tokens
- [ ] `GET /Auth/OIDC/Providers` - List OIDC providers
- [ ] `GET /Auth/OIDC/Authorize/{providerId}` - Start OIDC flow
- [ ] `GET /Auth/OIDC/Callback` - OIDC callback

### 3.8 User Management API
- [ ] `GET /Users` - List users (admin)
- [ ] `GET /Users/{id}` - Get user
- [ ] `POST /Users/{id}` - Update user
- [ ] `DELETE /Users/{id}` - Delete user (admin)
- [ ] `POST /Users/{id}/Password` - Change password

---

## 📚 LIBRARIES (Week 5-6)

### 4.1 Library Entity
- [ ] Library domain entity
- [ ] Library types enum (Movies, TvShows, Music, Photos, Mixed)
- [ ] Library settings (scan interval, metadata language, etc.)
- [ ] Path validation (exists, readable)

### 4.2 Library Repository
- [ ] LibraryRepository interface
- [ ] PostgreSQL implementation
- [ ] Path storage (array or JSON)

### 4.3 Library Service
- [ ] Create library with validation
- [ ] Update library paths
- [ ] Delete library (cascade media items?)
- [ ] Trigger library scan

### 4.4 Library API Endpoints
- [ ] `GET /Libraries` - List libraries
- [ ] `GET /Libraries/{id}` - Get library details
- [ ] `POST /Libraries` - Create library (admin)
- [ ] `POST /Libraries/{id}` - Update library (admin)
- [ ] `DELETE /Libraries/{id}` - Delete library (admin)
- [ ] `POST /Libraries/{id}/Refresh` - Trigger scan

---

## 🔍 FILE SCANNER (Week 6-7)

### 5.1 Directory Walker
- [ ] Recursive directory scanning
- [ ] Symlink handling
- [ ] Permission error handling
- [ ] Scan cancellation support
- [ ] Progress reporting

### 5.2 File Type Detection
- [ ] Video formats (.mkv, .mp4, .avi, .mov, .wmv, .m4v)
- [ ] Audio formats (.mp3, .flac, .m4a, .wav, .ogg, .opus)
- [ ] Image formats (.jpg, .png, .webp, .gif)
- [ ] Subtitle formats (.srt, .vtt, .ass, .sub)
- [ ] MIME type detection (magic bytes)

### 5.3 Naming Parser
- [ ] Movie parser: "Title (Year).ext", "Title.Year.ext"
- [ ] TV parser: "Show S01E01", "Show 1x01", "Show - 01"
- [ ] Music parser: "Artist - Album - Track - Title"
- [ ] Clean title extraction (remove quality tags, etc.)

### 5.4 Scanner Service
- [ ] Full scan (all files)
- [ ] Incremental scan (new/changed files)
- [ ] Remove deleted files from DB
- [ ] Scan queue (background processing)
- [ ] Concurrent scanning (worker pool)

---

## 🎬 MEDIA ITEMS (Week 7-8)

### 6.1 Media Item Entity
- [ ] Base media item (id, type, name, path, library_id)
- [ ] Movie-specific fields (year, runtime, tagline)
- [ ] Series fields (seasons, episodes)
- [ ] Episode fields (season_number, episode_number, series_id)
- [ ] Music fields (artist, album, track_number, duration)
- [ ] Metadata JSON field for extensibility

### 6.2 Media Repository
- [ ] MediaItemRepository interface
- [ ] PostgreSQL implementation
- [ ] Efficient queries (by library, by type, by parent)
- [ ] Pagination support

### 6.3 Media Service
- [ ] Get item with full details
- [ ] Update item metadata
- [ ] Delete item
- [ ] Get children (episodes for series, tracks for album)

### 6.4 Media API Endpoints
- [ ] `GET /Items` - Browse items (with filters)
- [ ] `GET /Items/{id}` - Get item details
- [ ] `GET /Items/{id}/Similar` - Similar items
- [ ] `POST /Items/{id}` - Update item
- [ ] `DELETE /Items/{id}` - Delete item (admin)
- [ ] `GET /Users/{userId}/Items` - User's view (with progress)
- [ ] `GET /Users/{userId}/Items/Resume` - Continue watching
- [ ] `GET /Users/{userId}/Items/Latest` - Recently added

---

## ▶️ PLAYBACK (Week 9-10)

### 7.1 File Serving
- [ ] Static file handler for media
- [ ] Range request support (HTTP 206)
- [ ] Partial content responses
- [ ] Proper Content-Type headers
- [ ] ETag support for caching
- [ ] Concurrent stream limiting

### 7.2 Playback Info
- [ ] MediaSource info (container, codecs, bitrate)
- [ ] FFprobe integration for media analysis
- [ ] Direct play capability check
- [ ] Transcode requirement detection (for Phase 2)

### 7.3 Streaming Endpoints
- [ ] `GET /Videos/{id}/stream` - Video direct play
- [ ] `GET /Audio/{id}/stream` - Audio direct play
- [ ] `GET /Items/{id}/PlaybackInfo` - Playback capabilities
- [ ] `GET /Videos/{id}/Subtitles/{index}` - Subtitle files

### 7.4 Playback Progress
- [ ] Progress entity (user_id, item_id, position, played)
- [ ] Progress repository
- [ ] `POST /Sessions/Playing` - Start playback
- [ ] `POST /Sessions/Playing/Progress` - Update position
- [ ] `POST /Sessions/Playing/Stopped` - Stop playback
- [ ] `POST /Users/{userId}/PlayedItems/{itemId}` - Mark as played
- [ ] `DELETE /Users/{userId}/PlayedItems/{itemId}` - Mark as unplayed

---

## 🔎 SEARCH (Week 11)

### 8.1 PostgreSQL Full-Text Search
- [ ] tsvector columns on media_items
- [ ] GIN index for fast search
- [ ] Search configuration (language, stemming)
- [ ] Trigger for automatic tsvector update

### 8.2 Typesense Integration (Optional Enhancement)
- [ ] Typesense client setup
- [ ] Index schema definition
- [ ] Index sync on media changes
- [ ] Search with facets

### 8.3 Search Service
- [ ] Text search query
- [ ] Fuzzy matching
- [ ] Result ranking
- [ ] Search suggestions (autocomplete)

### 8.4 Search API
- [ ] `GET /Search/Hints` - Search autocomplete
- [ ] `GET /Items` with searchTerm parameter
- [ ] Search filters (type, library, year range)

### 8.5 Filtering & Sorting
- [ ] Filter by genre
- [ ] Filter by year/year range
- [ ] Filter by media type
- [ ] Filter by library
- [ ] Sort by name, date added, year, rating
- [ ] Pagination (limit, offset / startIndex, limit)

---

## 🖼️ IMAGES (Week 12)

### 9.1 Image Entity
- [ ] Image types (Primary, Backdrop, Logo, Thumb, Banner)
- [ ] Image storage (path, width, height, format)
- [ ] Multiple images per item

### 9.2 Image Extraction
- [ ] FFmpeg thumbnail extraction
- [ ] Embedded artwork extraction (music)
- [ ] Default placeholders

### 9.3 Image Service
- [ ] Get image by item and type
- [ ] Resize on-the-fly (or cache resized)
- [ ] Image format conversion (webp)

### 9.4 Image API
- [ ] `GET /Items/{id}/Images/{type}` - Get image
- [ ] `GET /Items/{id}/Images/{type}/{index}` - Specific image
- [ ] Image parameters (maxWidth, maxHeight, quality)
- [ ] `POST /Items/{id}/Images/{type}` - Upload image (admin)
- [ ] `DELETE /Items/{id}/Images/{type}` - Delete image (admin)

### 9.5 Caching (Dragonfly/Redis)
- [ ] Image cache (binary data or paths)
- [ ] Metadata cache (frequently accessed items)
- [ ] Cache invalidation on updates
- [ ] Cache TTL configuration
- [ ] Cache statistics endpoint

---

## ⚙️ SYSTEM (Week 13)

### 10.1 System Info
- [ ] Server version, build info
- [ ] Go version, OS, architecture
- [ ] Startup time, uptime

### 10.2 System Stats
- [ ] Active sessions count
- [ ] Library item counts
- [ ] Storage usage
- [ ] Active streams

### 10.3 System API
- [ ] `GET /System/Info` - Server info
- [ ] `GET /System/Info/Public` - Public info (no auth)
- [ ] `GET /System/Ping` - Health check
- [ ] `POST /System/Restart` - Restart server (admin)
- [ ] `POST /System/Shutdown` - Shutdown server (admin)
- [ ] `GET /System/Logs` - View logs (admin)
- [ ] `GET /System/Logs/Log` - Download log file (admin)

### 10.4 Configuration API
- [ ] `GET /System/Configuration` - Get config
- [ ] `POST /System/Configuration` - Update config (admin)
- [ ] `GET /Branding/Configuration` - Branding settings
- [ ] `POST /Branding/Configuration` - Update branding (admin)

---

## 🧪 TESTING (Week 14-15)

### 11.1 Unit Tests
- [ ] Domain entity tests
- [ ] Service layer tests
- [ ] Repository mock tests
- [ ] Utility function tests
- [ ] Target: 80% code coverage

### 11.2 Integration Tests
- [ ] API endpoint tests (httptest)
- [ ] Database integration (testcontainers)
- [ ] Full auth flow tests
- [ ] OIDC mock provider tests

### 11.3 Test Infrastructure
- [ ] testcontainers-go setup
- [ ] Test fixtures and factories
- [ ] Test database seeding
- [ ] CI test parallelization

### 11.4 Documentation
- [ ] OpenAPI/Swagger spec generation
- [ ] API documentation site
- [ ] User installation guide
- [ ] Development contributing guide
- [ ] OIDC configuration guide

---

## 🔄 MIGRATION TOOL (Week 16-17)

### 12.1 Jellyfin Database Parser
- [ ] SQLite database reader
- [ ] Schema analysis
- [ ] Data extraction queries

### 12.2 Migration Mapping
- [ ] User migration (preserve IDs, passwords)
- [ ] Library migration
- [ ] Media item migration
- [ ] Image path migration
- [ ] Playback history migration
- [ ] Settings migration

### 12.3 Migration Service
- [ ] Backup existing data
- [ ] Dry-run mode
- [ ] Progress reporting
- [ ] Error handling and rollback
- [ ] Verification after migration

### 12.4 Migration CLI/UI
- [ ] CLI migrate command
- [ ] First-launch wizard detection
- [ ] Migration progress UI (if needed)

---

## 🚀 RELEASE PREP (Week 18)

### 13.1 Polish
- [ ] Error message review
- [ ] Logging consistency check
- [ ] API response consistency
- [ ] Performance profiling (pprof)
- [ ] Memory leak checks

### 13.2 Build & Release
- [ ] Cross-compilation (Linux, macOS, Windows)
- [ ] ARM64 builds
- [ ] Docker image (multi-arch)
- [ ] GitHub release automation
- [ ] Release notes generation

### 13.3 Documentation Final
- [ ] README updates
- [ ] CHANGELOG
- [ ] Installation guide
- [ ] Quick start guide
- [ ] FAQ

---

## 📦 DEPENDENCIES TO ADD

### Core
```go
// go.mod additions needed:
github.com/golang-jwt/jwt/v5      // JWT handling
github.com/coreos/go-oidc/v3      // OIDC client
golang.org/x/oauth2               // OAuth2 flows
github.com/dgraph-io/ristretto/v2 // In-memory cache (optional, have Dragonfly)
github.com/redis/go-redis/v9      // Redis/Dragonfly client
github.com/typesense/typesense-go // Typesense client
```

### Testing
```go
github.com/testcontainers/testcontainers-go  // Integration tests
github.com/stretchr/testify                   // Assertions (optional)
```

### Tools (go:tool directive, Go 1.24)
```go
//go:tool github.com/sqlc-dev/sqlc/cmd/sqlc
//go:tool github.com/golang-migrate/migrate/v4/cmd/migrate
//go:tool github.com/golangci/golangci-lint/cmd/golangci-lint
```

---

## 🎯 PRIORITY ORDER

1. **Database setup** (migrations, sqlc) - Everything depends on this
2. **User & Auth** (including OIDC) - Required for all protected endpoints
3. **Libraries & Scanner** - Core functionality
4. **Media Items** - Browsing content
5. **Playback** - Actually using the media
6. **Search** - Finding content
7. **Images** - Visual experience
8. **System APIs** - Management
9. **Testing** - Quality assurance
10. **Migration** - Adoption path
11. **Release** - Ship it!

---

## ⚠️ BLOCKERS & RISKS

| Risk | Impact | Mitigation |
|------|--------|------------|
| FFmpeg integration complexity | High | Start with basic ffprobe, defer transcoding |
| OIDC provider variations | Medium | Test with multiple providers (Keycloak, Auth0, Authentik) |
| Jellyfin API compatibility | High | Maintain reference client tests |
| Performance with large libraries | Medium | Pagination, indexing, caching from start |
| Migration data integrity | High | Extensive validation, dry-run mode |

---

## 📅 WEEKLY MILESTONES

| Week | Milestone | Validation |
|------|-----------|------------|
| 3 | Database migrations working | Can run migrate up/down |
| 4 | User registration/login | Can authenticate via API |
| 5 | OIDC login working | Can login via external provider |
| 6 | Libraries created | Can add library paths |
| 7 | Scanner running | Files appear in database |
| 8 | Media browsing | Can list media via API |
| 9 | Direct play | Can stream video file |
| 10 | Playback tracking | Progress saved and resumed |
| 11 | Search working | Can find media by name |
| 12 | Images serving | Thumbnails displayed |
| 13 | System APIs | Server info available |
| 15 | 80% test coverage | CI green with coverage |
| 17 | Migration working | Jellyfin data imported |
| 18 | v0.1.0-alpha | Public release |

