# Design Documentation - Clarification Questions

<!-- DESIGN: .complete, README, SCAFFOLD_TEMPLATE, test_output_claude -->


**Purpose**: Gather requirements and clarifications before completing documentation gaps
**Target Audience**: Project maintainers, product owners
**Date**: 2026-01-31

---

## How to Use This Document

Each section contains questions about incomplete documentation. Answer these questions to provide context for completing the docs.

**Priority Levels**:
- 🔴 **CRITICAL** - Blocks implementation or user adoption
- 🟡 **HIGH** - Important for feature completeness
- 🟢 **MEDIUM** - Nice to have, improves UX
- ⚪ **LOW** - Can be deferred

---

## 1. Operations Documentation (CRITICAL 🔴)

### 1.1 DEVELOPMENT.md - Developer Setup Guide

**Context**: New contributors need step-by-step setup instructions.

**Questions**:

1. **Development Environment**:
   - Q: What's the recommended dev environment? (Local, Coder workspace, Docker Compose, or all three?)
   - Q: Do we have a `devcontainer.json` or should we create one?
   - Q: Should we document both Linux and macOS setups? Windows WSL2?

2. **Required Tools**:
   - Q: Minimum required tool versions? (Go 1.25+, Node 20+, PostgreSQL 18+, Docker 27+?)
   - Q: Do devs need local PostgreSQL or Docker Compose only?
   - Q: Is Dragonfly required for local dev or optional?
   - Q: Is Typesense required for local dev or optional?

3. **First-Time Setup**:
   - Q: What's the recommended clone → run sequence?
   - Q: Should we provide a `make dev-setup` command?
   - Q: How do developers seed test data?
   - Q: Is there a sample `.env.example` file?

4. **Hot Reload**:
   - Q: Is `air` configured and working?
   - Q: Does frontend hot reload work with Vite?
   - Q: Should docs mention `--watch` flags?

5. **Troubleshooting**:
   - Q: Common setup errors? (Port conflicts, PATH issues, etc.)
   - Q: How to reset dev database?
   - Q: How to clear cache/search index?

---

### 1.2 SETUP.md - User Setup Guide

**Context**: End users need deployment instructions.

**Questions**:

1. **Deployment Methods**:
   - Q: What deployment methods do we officially support?
     - [ ] Docker Compose (homelab users)
     - [ ] K3s/Kubernetes (advanced users)
     - [ ] Docker Swarm
     - [ ] Bare metal / systemd service
   - Q: Which method should be documented first (most common)?

2. **Hardware Requirements**:
   - Q: Minimum specs? (CPU cores, RAM, disk?)
   - Q: Recommended specs for X users / Y media files?
   - Q: GPU required for transcoding or optional?

3. **Prerequisites**:
   - Q: Can users use managed PostgreSQL (e.g., Supabase, Neon, AWS RDS)?
   - Q: Can users use managed Redis (Upstash, AWS ElastiCache) instead of Dragonfly?
   - Q: Is Typesense mandatory or can search be disabled?

4. **Initial Configuration**:
   - Q: Is there a setup wizard UI or all config via env vars?
   - Q: Do users create admin account via UI or CLI?
   - Q: Can users import existing media libraries on first run?

5. **Upgrade Path**:
   - Q: How do users upgrade versions? (Docker image tags, migrations auto-run?)
   - Q: Do we support downgrade/rollback?
   - Q: Breaking changes between versions?

---

### 1.3 GITFLOW.md - Workflow Guide

**Context**: Contributors need to know branch strategy and PR process.

**Questions**:

1. **Branch Strategy**:
   - Q: Main branch: `main` or `develop`?
   - Q: Feature branches: `feature/*`, `feat/*`, or freeform?
   - Q: Bugfix branches: `fix/*` or `bugfix/*`?
   - Q: Release branches: `release/*` or tags only?

2. **Commit Conventions**:
   - Q: Conventional Commits enforced? (feat, fix, docs, chore, etc.)
   - Q: Are commit hooks configured? (pre-commit, commit-msg?)
   - Q: Do commits require sign-off? GPG signatures?

3. **PR Requirements**:
   - Q: Minimum test coverage for PRs? (80%?)
   - Q: Lint checks required to pass?
   - Q: Who can approve PRs? (Maintainers only, or any contributor?)
   - Q: Squash merge, merge commit, or rebase?

4. **Release Process**:
   - Q: Who triggers releases? (Automated via Release Please, or manual?)
   - Q: Versioning: SemVer strict?
   - Q: Changelog: Auto-generated from commits or manual?
   - Q: Release notes: Where published? (GitHub Releases, docs/wiki?)

5. **Hotfix Process**:
   - Q: Hotfix branches go to `main` or `develop`?
   - Q: How are hotfixes backported?

---

### 1.4 BEST_PRACTICES.md - Go Patterns & Conventions

**Context**: Contributors need coding standards for consistency.

**Questions**:

1. **Code Organization**:
   - Q: Confirmed module structure: `internal/content/{module}/` and `internal/service/{service}/`?
   - Q: Where do shared utilities go? `pkg/` or `internal/utils/`?
   - Q: Confirmed: no `pkg/` for now, all internal?

2. **Error Handling**:
   - Q: Sentinel errors: package-level `var ErrNotFound = errors.New(...)`?
   - Q: Custom error types for business logic errors?
   - Q: How to convert internal errors to API errors? (middleware?)

3. **Testing Standards**:
   - Q: Confirmed: table-driven tests for all logic?
   - Q: Confirmed: testify for assertions?
   - Q: Confirmed: mockery for mocks (not manual mocks)?
   - Q: Test file naming: `*_test.go` only, or also `_integration_test.go`?

4. **Context Usage**:
   - Q: Always pass `context.Context` as first param?
   - Q: When to cancel contexts? (HTTP request scope, job scope?)
   - Q: How to propagate trace IDs in context?

5. **Performance Patterns**:
   - Q: When to use otter (L1 cache) vs rueidis (L2 cache)?
   - Q: When to use sturdyc for request coalescing?
   - Q: Connection pool settings: defaults or custom per-service?

6. **Logging**:
   - Q: Structured logging required? (always `slog.Info(ctx, "msg", "key", val)`?)
   - Q: Log levels: when to use Debug vs Info vs Warn vs Error?
   - Q: Sensitive data: how to redact? (passwords, tokens, etc.)

---

### 1.5 VERSIONING.md - Release Strategy

**Context**: How versions are numbered and released.

**Questions**:

1. **SemVer Rules**:
   - Q: Currently at v0.x or v1.x?
   - Q: Breaking changes: major version bump?
   - Q: New features: minor version bump?
   - Q: Bug fixes: patch version bump?

2. **Pre-Releases**:
   - Q: Alpha/Beta releases: versioned how? (`v1.2.0-alpha.1`, `v1.2.0-beta.1`?)
   - Q: Nightly builds: tagged or not?

3. **Release Cadence**:
   - Q: Release every X weeks, or on-demand?
   - Q: Hotfixes released immediately or batched?

4. **Compatibility**:
   - Q: API compatibility guarantees? (e.g., v1 API stable for 1 year?)
   - Q: Database migration guarantees? (can skip versions, or must upgrade sequentially?)

---

### 1.6 BRANCH_PROTECTION.md - GitHub Branch Rules

**Context**: Branch protection rules for main/develop branches.

**Questions**:

1. **Protected Branches**:
   - Q: `main` and `develop` both protected?
   - Q: `release/*` branches protected?

2. **Required Checks**:
   - Q: All CI checks must pass? (tests, lint, build?)
   - Q: Code coverage check required?
   - Q: Security scans required? (CodeQL, Trivy?)

3. **Review Requirements**:
   - Q: How many approvals required? (1, 2, or more?)
   - Q: Can authors approve own PRs? (should be No)
   - Q: Stale review dismissal on new commits?

4. **Force Push**:
   - Q: Force push allowed for anyone? (should be No)
   - Q: Admins only can force push?

---

### 1.7 REVERSE_PROXY.md - Deployment Best Practices

**Context**: Users deploying behind reverse proxies (Traefik, Nginx, Caddy).

**Questions**:

1. **Recommended Proxy**:
   - Q: Do we recommend Traefik (K8s-native), Nginx (traditional), or Caddy (auto HTTPS)?
   - Q: Should we provide config examples for all three?

2. **HTTPS/TLS**:
   - Q: Do we expect TLS termination at proxy? (yes, recommended)
   - Q: Should Revenge backend support TLS natively? (optional)

3. **Reverse Proxy Headers**:
   - Q: Do we read `X-Forwarded-For`, `X-Real-IP`, `X-Forwarded-Proto`?
   - Q: Do we trust proxy headers by default or require config?

4. **Subdomain vs Subpath**:
   - Q: Can Revenge run at `example.com/revenge` or requires `revenge.example.com`?
   - Q: Base path configuration? (env var `REVENGE_BASE_PATH=/revenge`?)

5. **WebSocket Proxying**:
   - Q: Special config needed for WebSocket upgrades?
   - Q: Tested with Traefik/Nginx/Caddy?

---

### 1.8 DATABASE_AUTO_HEALING.md - Self-Healing Patterns

**Context**: How the app recovers from database issues.

**Questions**:

1. **Connection Pool Healing**:
   - Q: pgxpool already handles reconnects?
   - Q: Custom health checks beyond pgxpool defaults?

2. **Consistency Checks**:
   - Q: Are there background jobs to verify data integrity?
   - Q: Examples: orphaned records, invalid foreign keys?

3. **Migration Failures**:
   - Q: What happens if migration fails mid-way?
   - Q: Can app auto-retry migrations on startup?

4. **Corruption Detection**:
   - Q: Any checksums or validation for critical tables?
   - Q: How to handle corrupted data?

---

## 2. Technical Documentation (HIGH 🟡)

### 2.1 API.md - OpenAPI Reference

**Context**: API documentation for developers.

**Questions**:

1. **OpenAPI Spec Location**:
   - Q: Where is the OpenAPI spec? (`api/openapi/spec.yaml`?)
   - Q: Is it hand-written or generated from Go code (ogen)?

2. **API Versioning**:
   - Q: API version in URL: `/api/v1/*`?
   - Q: Future versions: `/api/v2/*`?
   - Q: How long is v1 supported after v2 release?

3. **Authentication**:
   - Q: All endpoints require Bearer token except `/api/v1/auth/*`?
   - Q: API keys supported? (header `X-API-Key`?)

4. **Rate Limiting**:
   - Q: Per-user rate limits?
   - Q: Per-IP rate limits?
   - Q: Headers returned: `X-RateLimit-Limit`, `X-RateLimit-Remaining`?

5. **Error Responses**:
   - Q: Standardized error format? (RFC 7807 Problem Details?)
   - Q: Example error response:
     ```json
     {
       "error": "not_found",
       "message": "Movie with ID 123 not found",
       "details": {}
     }
     ```

---

### 2.2 CONFIGURATION.md - koanf Usage Guide

**Context**: How to configure Revenge via YAML/env vars.

**Questions**:

1. **Config File Location**:
   - Q: Default config path: `./config.yaml`, `/etc/revenge/config.yaml`, or both?
   - Q: Can user override via `--config` flag?

2. **Environment Variable Precedence**:
   - Q: Confirmed: env vars override config file?
   - Q: Prefix: `REVENGE_*` for all vars?

3. **Hot Reload**:
   - Q: Which config changes can be hot-reloaded?
     - Log level: Yes?
     - Database URL: No (requires restart)?
     - Cache URL: No (requires restart)?

4. **Validation**:
   - Q: Config validated on startup?
   - Q: App exits if invalid config?

5. **Secrets Management**:
   - Q: Support for env var expansion? (e.g., `${JWT_SECRET}`?)
   - Q: Support for file-based secrets? (e.g., `jwt_secret_file: /run/secrets/jwt`?)

---

### 2.3 FRONTEND.md - SvelteKit Architecture

**Context**: Frontend structure for contributors.

**Questions**:

1. **Component Organization**:
   - Q: Confirmed structure:
     ```
     src/
     ├── routes/           # SvelteKit pages
     ├── lib/
     │   ├── components/   # Reusable components
     │   ├── stores/       # Svelte stores
     │   └── utils/        # Utilities
     ```

2. **State Management**:
   - Q: TanStack Query for server state?
   - Q: Svelte stores for UI state?
   - Q: Any global state library? (Zustand, Jotai, or just stores?)

3. **Styling**:
   - Q: Tailwind CSS 4 config location?
   - Q: Custom theme? (colors, fonts?)
   - Q: Dark mode implementation? (class-based, media query?)

4. **API Client**:
   - Q: Is there a generated API client from OpenAPI spec?
   - Q: Fetch wrapper or axios?

5. **Forms**:
   - Q: Superforms for form handling?
   - Q: Zod for validation?

6. **Authentication**:
   - Q: JWT stored in httpOnly cookie or localStorage?
   - Q: Token refresh logic?

---

### 2.4 AUDIO_STREAMING.md - Streaming Protocols

**Context**: How audio is streamed to clients.

**Questions**:

1. **Streaming Protocol**:
   - Q: HLS (HTTP Live Streaming) for audio too?
   - Q: Progressive download supported?
   - Q: Adaptive bitrate for audio?

2. **Audio Formats**:
   - Q: Supported codecs: AAC, MP3, FLAC, Opus?
   - Q: Transcoding: when does it happen? (on-demand, pre-generated?)

3. **Progress Tracking**:
   - Q: How often does client report playback position? (every 10s, 30s?)
   - Q: Resume playback: stored per-user?

4. **Offline Support**:
   - Q: Can users download for offline? (mobile apps?)
   - Q: DRM or plain files?

---

## 3. Content Modules (MEDIUM 🟡)

### 3.1 MUSIC_MODULE.md - Music Library

**Questions**:

1. **Lidarr Integration**:
   - Q: Sync from Lidarr: one-way or two-way?
   - Q: Metadata priority: Lidarr → MusicBrainz → Last.fm?

2. **Music Formats**:
   - Q: Supported: MP3, FLAC, AAC, OGG, ALAC?
   - Q: Lossless vs lossy handling?

3. **Metadata**:
   - Q: Embedded tags (ID3, Vorbis) vs external metadata?
   - Q: Album art: embedded or separate files?

4. **Playlists**:
   - Q: M3U playlists supported?
   - Q: User-created playlists stored in DB?

5. **Music Organization**:
   - Q: Library scan: recursive folder scan?
   - Q: Expected structure: `Artist/Album/Track.mp3`?

---

### 3.2 AUDIOBOOK_MODULE.md - Audiobook Library

**Questions**:

1. **Chaptarr Integration**:
   - Q: Chaptarr is fork of Readarr for audiobooks?
   - Q: Metadata sync: one-way or two-way?

2. **Audiobook Formats**:
   - Q: Supported: M4B (chapters), MP3 (multi-file)?
   - Q: Chapter markers: extracted from M4B metadata?

3. **Progress Tracking**:
   - Q: Resume per-user, per-audiobook?
   - Q: Sync progress to Chaptarr?

4. **Metadata**:
   - Q: Primary source: Audnexus API?
   - Q: Fallback: OpenLibrary?

---

### 3.3 BOOK_MODULE.md - eBook Library

**Questions**:

1. **eBook Formats**:
   - Q: Supported: EPUB, PDF, MOBI, AZW3?
   - Q: Does Revenge include eBook reader UI? (or just library management?)

2. **Reading Progress**:
   - Q: Track page/chapter progress?
   - Q: Sync across devices?

3. **Metadata**:
   - Q: Primary: OpenLibrary?
   - Q: Fallback: Goodreads?

4. **Chaptarr Integration**:
   - Q: Same Chaptarr instance as audiobooks, or separate?

---

### 3.4 PODCASTS.md - Podcast Library

**Questions**:

1. **RSS Feeds**:
   - Q: Users add podcast by RSS URL?
   - Q: Auto-refresh interval? (daily, hourly?)

2. **Episode Download**:
   - Q: Auto-download new episodes?
   - Q: Delete after playback?

3. **Metadata**:
   - Q: Podcast Index API for discovery?
   - Q: iTunes API for artwork?

---

## 4. Patterns (MEDIUM 🟡)

### 4.1 ARR_INTEGRATION.md - Arr Pattern

**Questions**:

1. **Webhook Handling**:
   - Q: Arr services push webhooks to `/api/v1/webhooks/radarr`?
   - Q: Events: `Download`, `Upgrade`, `Rename`, `Delete`?

2. **Metadata Sync**:
   - Q: On webhook event, fetch full metadata from Arr API?
   - Q: Store Arr ID for future lookups?

3. **Conflict Resolution**:
   - Q: If local metadata conflicts with Arr, which wins?
   - Q: User can choose priority?

---

### 4.2 METADATA_ENRICHMENT.md - Enrichment Pattern

**Questions**:

1. **Priority Chain**:
   - Q: Confirmed:
     ```
     1. Local cache
     2. Arr services
     3. Internal (Stash-App)
     4. External APIs
     5. Background enrichment
     ```

2. **Cache Strategy**:
   - Q: Cache TTL per data type? (metadata: 24h, images: 7d?)

3. **Background Jobs**:
   - Q: River jobs for enrichment?
   - Q: Job priority: low (don't block user requests)?

---

### 4.3 WEBHOOK_PATTERNS.md - Webhook Events

**Questions**:

1. **Event Catalog**:
   - Q: Which events do we emit?
     - [ ] `playback.started`, `playback.stopped`, `playback.paused`
     - [ ] `library.scan.started`, `library.scan.completed`
     - [ ] `library.item.added`, `library.item.updated`, `library.item.deleted`
     - [ ] `user.created`, `user.updated`, `user.deleted`
     - [ ] `transcode.started`, `transcode.completed`, `transcode.failed`

2. **Webhook Configuration**:
   - Q: Users configure webhooks in UI?
   - Q: Retry on failure?

3. **Security**:
   - Q: HMAC signatures for webhook payloads?

---

## 5. Missing Documentation (MEDIUM 🟡)

### 5.1 COLLECTIONS.md (New Feature)

**Questions**:

1. **Feature Scope**:
   - Q: User-created collections: yes?
   - Q: Cross-content-type collections? (movies + TV shows in one collection?)
   - Q: Smart collections based on filters? (e.g., "All 4K movies")

2. **Metadata**:
   - Q: Collection poster/backdrop?
   - Q: Collection description?

---

### 5.2 TRANSCODING.md (New Service Doc)

**Questions**:

1. **Service Scope**:
   - Q: Is transcoding a service or part of playback service?
   - Q: On-demand transcoding only, or pre-transcode jobs?

2. **Offloading**:
   - Q: Use Blackbeard for offloading?
   - Q: Local transcoding with FFmpeg also supported?

---

### 5.3 EPG.md (New Service Doc - LiveTV)

**Questions**:

1. **EPG Source**:
   - Q: Fetch XMLTV from ErsatzTV/TVHeadend?
   - Q: Parse and store in database?

2. **UI**:
   - Q: EPG grid view in frontend?

---

### 5.4 OBSERVABILITY.md (New Pattern Doc)

**Questions**:

1. **Metrics**:
   - Q: Prometheus `/metrics` endpoint?
   - Q: Key metrics: request rate, error rate, duration (RED)?

2. **Tracing**:
   - Q: OpenTelemetry traces to Jaeger?
   - Q: Trace all HTTP requests?

3. **Logging**:
   - Q: Log aggregation: Loki?

---

### 5.5 TESTING.md (New Pattern Doc)

**Questions**:

1. **Testing Patterns**:
   - Q: Confirmed: table-driven tests?
   - Q: Test naming: `TestServiceName_MethodName`?

2. **Mocking**:
   - Q: Confirmed: mockery for generating mocks?

3. **Integration Tests**:
   - Q: Confirmed: testcontainers for PostgreSQL, Dragonfly?

---

## 6. Research Documentation (LOW ⚪)

### 6.1 USER_PAIN_POINTS_RESEARCH.md

**Questions**:

1. **Research Scope**:
   - Q: What pain points are we researching?
     - [ ] Plex issues (slow, closed-source)
     - [ ] Jellyfin issues (UI/UX, performance)
     - [ ] Emby issues (licensing, features)

2. **Sources**:
   - Q: Reddit threads, GitHub issues, forum posts?

---

### 6.2 UX_UI_RESOURCES.md

**Questions**:

1. **Resource Types**:
   - Q: Design systems to follow? (Material Design, Radix UI?)
   - Q: UI kits to reference? (Figma community files?)

---

## 7. Wiki Content (MEDIUM 🟡)

### Wiki-Specific Questions

**Context**: YAML files have `wiki_overview: "PLACEHOLDER: User-friendly overview"` for many docs.

**Questions**:

1. **Wiki Audience**:
   - Q: Target audience: end users (non-technical)?
   - Q: Tone: casual, friendly, or professional?

2. **Wiki Content**:
   - Q: For each integration, what do users need to know?
     - [ ] Why use this integration?
     - [ ] How to set it up? (API keys, config)
     - [ ] Common issues / FAQ

3. **Wiki vs Design Docs**:
   - Q: Design docs = technical implementation details
   - Q: Wiki docs = user-facing setup guides
   - Q: Correct understanding?

---

## Next Steps

1. **Answer Questions**: Project maintainer reviews and answers applicable questions
2. **Document Answers**: Add answers as inline responses or separate doc
3. **Update TODO.md**: Refine action items based on answers
4. **Begin Documentation**: Use answers to fill gaps systematically

---

**Related Files**:
- [ANALYSIS.md](ANALYSIS.md) - Gap analysis
- [TODO.md](TODO.md) - Action items
