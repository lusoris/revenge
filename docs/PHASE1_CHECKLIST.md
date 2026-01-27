# Jellyfin Go - Phase 1 Implementation Checklist

## Goal: Single-Server MVP with PostgreSQL

**Timeline:** 3 months
**Focus:** Core functionality, Jellyfin API compatibility, easy setup

---

## Week 1-2: Project Foundation

### Setup
- [x] Initialize Go module
- [x] Create directory structure (cmd/, internal/, pkg/, docs/)
- [x] Setup CI/CD (GitHub Actions: lint, test, build)
- [x] Docker Compose for development (PostgreSQL optional)
- [x] Configuration system (koanf v2, env variables)
- [x] Logging (slog + tint with pretty console output)
- [x] Basic HTTP server (net/http.ServeMux, Go 1.22+ patterns)
- [x] Health check endpoints (/health/live, /health/ready)

### Database
- [x] PostgreSQL 18+ (required database)
- [ ] Setup golang-migrate
- [ ] Initial schema (users, sessions, libraries, media_items)
- [ ] Setup sqlc for type-safe queries
- [ ] Repository interfaces

**Deliverable:** Running HTTP server with health checks, database migrations working

---

## Week 3-4: Authentication & Users

### User Management
- [ ] User entity and repository
- [ ] Password hashing (bcrypt)
- [ ] User CRUD operations
- [ ] Admin vs regular user roles

### Authentication
- [ ] JWT token generation
- [ ] JWT validation
- [ ] Refresh token mechanism
- [ ] Session management (DB-backed)
- [ ] Auth middleware

### OIDC/SSO Integration
- [ ] OIDC discovery endpoint support
- [ ] OIDC provider configuration (Keycloak, Auth0, Authentik, etc.)
- [ ] Authorization code flow
- [ ] Token validation (ID token, access token)
- [ ] User provisioning from OIDC claims
- [ ] Role/group mapping from OIDC
- [ ] Multiple provider support
- [ ] Fallback to local auth when OIDC unavailable

### API Endpoints
- [ ] POST /api/Users/AuthenticateByName
- [ ] GET /api/Users/{userId}
- [ ] POST /api/Users/New
- [ ] POST /api/Auth/Logout
- [ ] GET /api/Users (list users, admin only)
- [ ] DELETE /api/Users/{userId} (admin only)
- [ ] GET /api/Auth/OIDC/Providers (list configured OIDC providers)
- [ ] GET /api/Auth/OIDC/Authorize/{providerId} (initiate OIDC flow)
- [ ] GET /api/Auth/OIDC/Callback (OIDC callback handler)

**Deliverable:** User registration, login, JWT auth, OIDC SSO working

---

## Week 5-6: Library Management

### Libraries
- [ ] Library entity (Movies, TV Shows, Music, Photos)
- [ ] Library repository
- [ ] Path validation

### API Endpoints
- [ ] GET /api/Libraries
- [ ] GET /api/Libraries/{id}
- [ ] POST /api/Libraries (create library)
- [ ] DELETE /api/Libraries/{id}
- [ ] POST /api/Libraries/{id}/Refresh (trigger scan)

### File System Scanner
- [ ] Basic directory walker
- [ ] File type detection (video, audio, image)
- [ ] Store file paths in database

**Deliverable:** Library creation, basic file scanning

---

## Week 7-8: Media Items & Metadata

### Media Items
- [ ] Media item entity
- [ ] Media repository
- [ ] Basic metadata extraction (filename parsing)

### Naming Parser
- [ ] Movie parser (e.g., "Movie Title (2024).mkv")
- [ ] TV episode parser (e.g., "Show S01E01.mkv")
- [ ] Music file parser

### API Endpoints
- [ ] GET /api/Items (browse media)
- [ ] GET /api/Items/{id}
- [ ] GET /api/Users/{userId}/Items (user-specific view)
- [ ] POST /api/Items/{id} (update metadata)
- [ ] DELETE /api/Items/{id}

**Deliverable:** Media browsing, basic metadata display

---

## Week 9-10: Direct Play & Streaming

### File Serving
- [ ] Static file serving for media
- [ ] Range request support (seeking)
- [ ] MIME type detection
- [ ] Direct play (no transcoding)

### API Endpoints
- [ ] GET /api/Items/{id}/PlaybackInfo
- [ ] GET /api/Videos/{id}/stream (direct play)
- [ ] GET /api/Audio/{id}/stream
- [ ] POST /api/Sessions/Playing (playback start)
- [ ] POST /api/Sessions/Playing/Progress (update progress)
- [ ] POST /api/Sessions/Playing/Stopped

### Playback Progress
- [ ] Track playback position per user
- [ ] Mark as played/unplayed
- [ ] Resume from last position

**Deliverable:** Direct play working, progress tracking

---

## Week 11: Search & Filtering

### PostgreSQL Full-Text Search
- [ ] Add tsvector columns for search
- [ ] Search index maintenance
- [ ] Search query builder

### API Endpoints
- [ ] GET /api/Search/Hints
- [ ] GET /api/Items with filters (genre, year, type)
- [ ] GET /api/Items with sorting

### Filtering
- [ ] By genre
- [ ] By year
- [ ] By media type
- [ ] By library

**Deliverable:** Basic search and filtering

---

## Week 12: Image Serving & Caching

### Images
- [ ] Image entity (thumbnails, posters, backdrops)
- [ ] Image serving endpoint
- [ ] Basic thumbnail extraction (FFmpeg)

### Ristretto Cache
- [ ] In-memory cache setup
- [ ] Cache metadata queries
- [ ] Cache images
- [ ] Cache eviction policy

### API Endpoints
- [ ] GET /api/Items/{id}/Images/{type}
- [ ] POST /api/Items/{id}/Images (upload)
- [ ] DELETE /api/Items/{id}/Images/{type}

**Deliverable:** Image serving, basic caching

---

## Week 13: System & Configuration

### System Info
- [ ] Server info (version, OS, architecture)
- [ ] System stats (CPU, memory, disk)

### API Endpoints
- [ ] GET /api/System/Info
- [ ] GET /api/System/Ping
- [ ] POST /api/System/Restart (admin only)
- [ ] POST /api/System/Shutdown (admin only)
- [ ] GET /api/System/Logs (admin only)

### Configuration
- [ ] Server settings API
- [ ] Branding customization
- [ ] Network settings

**Deliverable:** System management endpoints

---

## Week 14-15: Testing & Documentation

### Testing
- [ ] Unit tests (80%+ coverage)
- [ ] Integration tests (API endpoints)
- [ ] Test database setup (testcontainers)
- [ ] Mock repository tests

### Documentation
- [ ] API documentation (Swagger/OpenAPI)
- [ ] User guide updates
- [ ] Development guide
- [ ] Deployment guide

### Bug Fixes
- [ ] Address known issues
- [ ] Performance optimization
- [ ] Memory leak checks

**Deliverable:** Well-tested, documented MVP

---

## Week 16-17: Jellyfin Migration Tool

### Migration
- [ ] SQLite/PostgreSQL parser for original jellyfin database
- [ ] Schema mapper (Jellyfin → Jellyfin Go)
- [ ] User migration (preserve passwords)
- [ ] Library migration
- [ ] Media item migration
- [ ] Playback history migration

### First-Launch Wizard
- [ ] Detect existing Jellyfin installation
- [ ] Migration UI (or CLI)
- [ ] Progress tracking
- [ ] Backup creation

**Deliverable:** Working migration from Jellyfin C# to Jellyfin Go

---

## Week 18: Polish & Release Prep

### Polish
- [ ] Error message improvements
- [ ] Logging improvements
- [ ] Performance tuning
- [ ] UI/UX refinements (if web UI modified)

### Release
- [ ] Binary builds (Linux, macOS, Windows)
- [ ] Docker images
- [ ] Release notes
- [ ] GitHub release

**Deliverable:** v0.1.0-alpha release

---

## Success Criteria

At the end of Phase 1, users should be able to:

- ✅ Install Jellyfin Go (binary, Docker, or from source)
- ✅ Create user accounts and login
- ✅ Add media libraries (Movies, TV, Music)
- ✅ Scan libraries and see media items
- ✅ Browse and search media
- ✅ Play media files (direct play)
- ✅ Track playback progress
- ✅ View images (posters, thumbnails)
- ✅ Migrate from existing Jellyfin installation

## Non-Goals for Phase 1

- ❌ Transcoding (Phase 2)
- ❌ Hardware acceleration (Phase 2)
- ❌ Advanced metadata providers (Phase 3)
- ❌ Plugins (Phase 4)
- ❌ Multi-instance clustering (Phase 5)
- ❌ CDN integration (Phase 5)

---

## Next: Phase 2 - Transcoding & Streaming

After Phase 1 completion, move to implementing FFmpeg transcoding, HLS streaming, and hardware acceleration support.
