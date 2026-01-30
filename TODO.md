# Revenge - Development Roadmap

> Modular media server with complete content isolation

**Last Updated**: 2026-01-30
**Build**: `GOEXPERIMENT=greenteagc,jsonv2 go build ./...`

---

## Milestone Approach

Each milestone requires **all three criteria** before moving on:
- **Structure** - All types, interfaces, and packages present per design spec
- **Linted** - `golangci-lint run` passes with no errors
- **Tested** - Unit tests written and passing

```
Legend:
[x] Complete
[~] Partial (in progress)
[ ] Not started
```

---

## Quick Status

```
M1: Infrastructure       [x][x][~]  Structure ✓ | Lint ✓ | Tests ~
M2: Shared Services      [x][x][ ]  Structure ✓ | Lint ✓ | Tests -
M3: Movie Module         [x][x][ ]  Structure ✓ | Lint ✓ | Tests -
M4: TV Shows Module      [x][x][ ]  Structure ✓ | Lint ✓ | Tests -
M5: QAR Module           [x][x][ ]  Structure ✓ | Lint ✓ | Tests -
M6: Music Module         [ ][ ][ ]  Structure - | Lint - | Tests -
M7: Books Module         [ ][ ][ ]  Structure - | Lint - | Tests -
M8: Comics Module        [ ][ ][ ]  Structure - | Lint - | Tests -
M9: Playback & Streaming [~][ ][ ]  Structure ~ | Lint - | Tests -
M10: External Metadata   [~][ ][ ]  Structure ~ | Lint - | Tests -
M11: Integration Tests   [ ][ ][ ]  Structure - | Lint - | Tests -
M12: Frontend            [ ][ ][ ]  Structure - | Lint - | Tests -
```

---

## M1: Infrastructure

> Core infrastructure: database, cache, search, jobs, health, config

### Structure
- [x] **Database** → `internal/infra/database/`
  - [x] PostgreSQL connection with pgx
  - [x] sqlc code generation setup
  - [x] Migration framework (goose)
  - [x] Partitioned activity_log table (000021)
- [x] **Cache** → `internal/infra/cache/`
  - [x] Rueidis client (distributed)
  - [x] Otter local cache (W-TinyLFU)
  - [x] Sturdyc API response coalescing
- [x] **Search** → `internal/infra/search/`
  - [x] Typesense client v4
  - [x] Collection management
  - [x] Isolated collections per module
- [x] **Jobs** → `internal/infra/jobs/`
  - [x] River queue setup
  - [x] Worker registration
  - [x] AuditLogWorker (async audit writes)
- [x] **Health** → `internal/infra/health/`
  - [x] Database health check
  - [x] Cache health check
  - [x] Search health check
  - [x] Jobs health check
- [x] **Config** → `pkg/config/`
  - [x] Koanf setup
  - [x] Environment binding
  - [x] Validation

### Lint
- [x] `golangci-lint run ./internal/infra/...` passes
- [x] `golangci-lint run ./pkg/config/...` passes

### Tests
- [~] Database connection tests
- [ ] Cache integration tests
- [ ] Search integration tests
- [ ] Jobs worker tests
- [ ] Health check tests
- [ ] Config validation tests

---

## M2: Shared Services

> User management, authentication, RBAC, sessions, activity logging

### Structure
- [x] **Users** → `internal/service/user/`
  - [x] User entity
  - [x] Repository (sqlc)
  - [x] Service layer
  - [x] Password hashing (argon2)
- [x] **Profiles** → `internal/service/profile/`
  - [x] Profile entity
  - [x] Repository
  - [x] Service
  - [x] User preferences (auto_play, continue_watching_days, etc.)
- [x] **Sessions** → `internal/service/session/`
  - [x] Session entity
  - [x] Repository
  - [x] Token generation
  - [x] Device tracking
- [x] **RBAC** → `internal/service/rbac/`
  - [x] Casbin setup
  - [x] Role definitions (admin, moderator, user, guest)
  - [x] Permission checks
  - [x] User-role management methods
    - [x] SetUserRole
    - [x] GetUserRole
    - [x] GetUsersForRole
    - [x] CountUsersForRole
    - [x] AddRoleForUser
    - [x] RemoveRoleForUser
- [x] **Grants** → `internal/service/grants/`
  - [x] Resource grants entity
  - [x] Repository
  - [x] HasGrant, CreateGrant, DeleteGrant, DeleteByResource
- [x] **Activity** → `internal/service/activity/`
  - [x] Activity log entity (partitioned)
  - [x] Repository
  - [x] Service with JSON changes
  - [x] Module constants (movie, tvshow, qar, user, library, system)
- [x] **Auth** → `internal/service/auth/`
  - [x] Login/logout
  - [x] Token validation
  - [x] Password reset flow
- [x] **API Keys** → `internal/service/apikeys/`
  - [x] Key generation
  - [x] Key validation
  - [x] Scopes

### Lint
- [x] `golangci-lint run ./internal/service/...` passes

### Tests
- [ ] User service tests
- [ ] Profile service tests
- [ ] Session service tests
- [ ] RBAC service tests
- [ ] Grants service tests
- [ ] Activity service tests
- [ ] Auth service tests
- [ ] API keys service tests

---

## M3: Movie Module

> Movies with metadata, cast, crew, collections, user data

### Structure
- [x] **Entities** → `internal/content/movie/`
  - [x] Movie entity
  - [x] Collection entity
  - [x] Cast/Crew entities
  - [x] User data (ratings, favorites, watchlist)
- [x] **Repository** → `internal/content/movie/repository.go`
  - [x] Full CRUD
  - [x] Relations (cast, crew, collections)
  - [x] User data queries
  - [x] ListResumeableMovies (30-day filter)
- [x] **Service** → `internal/content/movie/service.go`
  - [x] Business logic
  - [x] Metadata enrichment
  - [x] User data management
- [x] **Search** → Typesense collection
  - [x] Movies collection
  - [x] Search handlers
- [x] **API Handlers** → `internal/api/movie.go`
  - [x] CRUD endpoints
  - [x] Relationship endpoints
  - [x] User data endpoints

### Lint
- [x] `golangci-lint run ./internal/content/movie/...` passes

### Tests
- [ ] Movie repository tests
- [ ] Movie service tests
- [ ] Movie handler tests

---

## M4: TV Shows Module

> Series → Seasons → Episodes hierarchy

### Structure
- [x] **Entities** → `internal/content/tvshow/`
  - [x] Series entity
  - [x] Season entity
  - [x] Episode entity
  - [x] Cast/Crew entities
- [x] **Repository** → `internal/content/tvshow/repository.go`
  - [x] Full CRUD for series/seasons/episodes
  - [x] Hierarchical queries
  - [x] ListResumeableEpisodes (30-day filter)
  - [x] ListContinueWatchingSeries
- [x] **Service** → `internal/content/tvshow/service.go`
  - [x] Business logic
  - [x] Season/episode management
- [x] **Search** → Typesense collections
  - [x] Series collection
  - [x] Episodes collection
- [x] **API Handlers** → `internal/api/tvshow.go`
  - [x] Series CRUD
  - [x] Season CRUD
  - [x] Episode CRUD

### Lint
- [x] `golangci-lint run ./internal/content/tvshow/...` passes

### Tests
- [ ] TVShow repository tests
- [ ] TVShow service tests
- [ ] TVShow handler tests

---

## M5: QAR Module (Adult Content)

> Queen Anne's Revenge - isolated adult content system
> Schema: `qar` | Namespace: `/api/v1/qar/`

### Structure
- [x] **Expedition** (Movies) → `internal/content/qar/expedition/`
  - [x] Entity
  - [x] Repository
  - [x] Service
  - [x] ListByCrewID, ListByPortID, ListByFlagID
- [x] **Voyage** (Scenes) → `internal/content/qar/voyage/`
  - [x] Entity
  - [x] Repository
  - [x] Service
- [x] **Crew** (Performers) → `internal/content/qar/crew/`
  - [x] Entity
  - [x] Repository
  - [x] Service
- [x] **Port** (Studios) → `internal/content/qar/port/`
  - [x] Entity
  - [x] Repository
  - [x] Service
- [x] **Flag** (Tags) → `internal/content/qar/flag/`
  - [x] Entity
  - [x] Repository
  - [x] Service
- [x] **Fleet** (Libraries) → `internal/content/qar/fleet/`
  - [x] Entity
  - [x] Repository
  - [x] Service
- [x] **Fingerprinting** → `internal/service/fingerprint/`
  - [x] OSHash
  - [x] pHash
  - [x] MD5
- [x] **API Handlers** → `internal/api/adult.go`
  - [x] ~50 endpoints
  - [x] Converters
  - [x] RBAC checks (adult.browse, adult.stream, adult.metadata.write)
- [x] **Search** → 5 Typesense collections
  - [x] expeditions, voyages, crew, ports, flags

### Lint
- [x] `golangci-lint run ./internal/content/qar/...` passes

### Tests
- [ ] Expedition repository tests
- [ ] Voyage repository tests
- [ ] Crew repository tests
- [ ] Port repository tests
- [ ] Flag repository tests
- [ ] Fleet repository tests
- [ ] Fingerprint service tests
- [ ] QAR handler tests

---

## M6: Music Module

> Artists, albums, tracks with gapless playback

### Structure
- [ ] **Entities** → `internal/content/music/`
  - [ ] Artist entity
  - [ ] Album entity
  - [ ] Track entity
  - [ ] Playlist entity
- [ ] **Repository**
  - [ ] Full CRUD
  - [ ] Artist-Album-Track relationships
- [ ] **Service**
  - [ ] Business logic
  - [ ] Gapless playback support
- [ ] **Metadata Providers**
  - [ ] MusicBrainz client
  - [ ] Last.fm client
- [ ] **Search** → Typesense collections
- [ ] **API Handlers**

### Lint
- [ ] `golangci-lint run ./internal/content/music/...` passes

### Tests
- [ ] Music repository tests
- [ ] Music service tests
- [ ] Music handler tests

---

## M7: Books Module

> E-books and audiobooks

### Structure
- [ ] **Entities** → `internal/content/book/`
  - [ ] Book entity
  - [ ] Author entity
  - [ ] Series entity
  - [ ] Audiobook entity (with chapters)
- [ ] **Repository**
  - [ ] Full CRUD
  - [ ] Series management
- [ ] **Service**
  - [ ] Business logic
  - [ ] Chapter markers
- [ ] **Metadata Providers**
  - [ ] Open Library client
  - [ ] Google Books client
  - [ ] Audible client
- [ ] **Search** → Typesense collections
- [ ] **API Handlers**

### Lint
- [ ] `golangci-lint run ./internal/content/book/...` passes

### Tests
- [ ] Book repository tests
- [ ] Book service tests
- [ ] Book handler tests

---

## M8: Comics Module

> Comics, manga, graphic novels

### Structure
- [ ] **Entities** → `internal/content/comics/`
  - [ ] Comic entity
  - [ ] Series entity
  - [ ] Publisher entity
  - [ ] Character entity
- [ ] **Repository**
  - [ ] Full CRUD
  - [ ] Reading progress
- [ ] **Service**
  - [ ] Business logic
  - [ ] Page tracking
- [ ] **Metadata Providers**
  - [ ] ComicVine client
  - [ ] Marvel API client
- [ ] **Search** → Typesense collections
- [ ] **API Handlers**

### Lint
- [ ] `golangci-lint run ./internal/content/comics/...` passes

### Tests
- [ ] Comics repository tests
- [ ] Comics service tests
- [ ] Comics handler tests

---

## M9: Playback & Streaming

> Session management, progress tracking, up-next queue

### Structure
- [x] **Playback Service** → `internal/service/playback/`
  - [x] PlaybackSession entity
  - [x] StreamInfo types
  - [x] StartPlayback, UpdateProgress, StopPlayback
  - [x] GetActiveSession
  - [x] UpNextProvider interface
  - [x] RegisterUpNextProvider
  - [x] BuildUpNextQueue framework
- [ ] **Provider Implementations**
  - [ ] TV: next episode provider
  - [ ] Movie: similar movies provider
  - [ ] QAR: similar expeditions provider
- [ ] **API Endpoints**
  - [ ] POST /api/playback/start
  - [ ] PUT /api/playback/{sessionId}/progress
  - [ ] POST /api/playback/{sessionId}/stop
  - [ ] GET /api/playback/up-next
- [ ] **Cross-Device Sync**
  - [ ] Polling endpoint `/api/sync/playback?since={ts}`
  - [ ] BroadcastToUser()

### Lint
- [ ] `golangci-lint run ./internal/service/playback/...` passes

### Tests
- [ ] Playback service tests
- [ ] Session management tests
- [ ] Up-next provider tests

---

## M10: External Metadata

> Third-party metadata providers and integrations

### Structure
- [x] **TMDb** → `internal/service/metadata/tmdb/`
  - [x] Client with circuit breaker
  - [x] Movie provider
  - [x] TV provider
- [x] **Radarr** → `internal/service/metadata/radarr/`
  - [x] Client
  - [x] Movie sync
- [x] **Sonarr** → `internal/service/metadata/sonarr/`
  - [x] Client
  - [x] TV sync
- [x] **Whisparr** → `internal/service/metadata/whisparr/`
  - [x] Client with circuit breaker
  - [x] Adult content sync
- [x] **StashDB** → `internal/service/metadata/stashdb/`
  - [x] GraphQL client
  - [x] Scene lookup
  - [x] Performer lookup
- [ ] **Stash-App** → `internal/service/metadata/stash_app/`
  - [ ] types.go - Stash-App GraphQL types
  - [ ] client.go - GraphQL client
  - [ ] provider.go - Metadata provider
  - [ ] module.go - fx wiring
- [ ] **StashDB Search Handlers**
  - [ ] SearchAdultStashDBScenes
  - [ ] GetAdultStashDBScene
  - [ ] SearchAdultStashDBPerformers
  - [ ] GetAdultStashDBPerformer
  - [ ] IdentifyAdultStashDBScene
- [ ] **TPDB Handlers**
  - [ ] SearchAdultTPDBScenes
  - [ ] GetAdultTPDBScene
  - [ ] GetAdultTPDBPerformer
- [ ] **Stash-App Sync Handlers**
  - [ ] SyncAdultStash
  - [ ] ImportAdultStash
  - [ ] GetAdultStashStatus

### Lint
- [ ] `golangci-lint run ./internal/service/metadata/...` passes

### Tests
- [ ] TMDb client tests
- [ ] Radarr client tests
- [ ] Sonarr client tests
- [ ] Whisparr client tests
- [ ] StashDB client tests
- [ ] Stash-App client tests

---

## M11: Integration Tests

> End-to-end workflows and external service integration

### Structure
- [ ] **Database Integration**
  - [ ] Migration up/down tests
  - [ ] Transaction rollback tests
  - [ ] Concurrent access tests
- [ ] **External Service Mocks**
  - [ ] TMDb mock server
  - [ ] Radarr mock server
  - [ ] Sonarr mock server
  - [ ] StashDB mock GraphQL
- [ ] **Workflow Tests**
  - [ ] User registration → login → browse → play
  - [ ] Library scan → metadata fetch → index
  - [ ] Watch progress → continue watching → complete
  - [ ] QAR access control → browse → stream
- [ ] **Performance Benchmarks**
  - [ ] API response time
  - [ ] Database query performance
  - [ ] Cache hit rates

### Lint
- [ ] All test files pass linting

### Tests
- [ ] All integration tests pass

---

## M12: Frontend

> SvelteKit 2 + Tailwind CSS 4 + shadcn-svelte

### Structure
- [ ] **Core Setup** → `web/`
  - [ ] SvelteKit 2 project
  - [ ] Tailwind CSS 4
  - [ ] shadcn-svelte components
- [ ] **Routes**
  - [ ] (auth) - login, register, forgot password
  - [ ] (admin) - admin dashboard
  - [ ] (media) - browse, library, search
  - [ ] (player) - video/audio player
- [ ] **Components**
  - [ ] Media cards
  - [ ] Player controls
  - [ ] Search interface
  - [ ] Navigation
- [ ] **State Management**
  - [ ] User session
  - [ ] Playback state
  - [ ] Theme preferences
- [ ] **RBAC Integration**
  - [ ] Route guards
  - [ ] Permission checks

### Lint
- [ ] ESLint/Prettier passes

### Tests
- [ ] Component tests (Vitest)
- [ ] E2E tests (Playwright)

---

## Pending Items (Not Milestone-Bound)

### QAR Relationship Handlers (Remaining)
- [ ] ListAdultSimilarMovies - shared flags/crew/port logic
- [ ] ListAdultMovieMarkers - needs marker/chapter entity

### QAR Request System
- [ ] SearchAdultRequests
- [ ] ListAdultRequests
- [ ] CreateAdultRequest
- [ ] GetAdultRequest
- [ ] VoteAdultRequest
- [ ] CommentAdultRequest
- [ ] ListAdultAdminRequests
- [ ] ApproveAdultRequest
- [ ] DeclineAdultRequest
- [ ] UpdateAdultRequestQuota
- [ ] ListAdultRequestRules
- [ ] CreateAdultRequestRule
- [ ] UpdateAdultRequestRule
- [ ] DeleteAdultRequestRule

### RBAC Permissions to Seed
- [ ] access.rules.view, access.rules.manage, access.bypass
- [ ] request.* permissions (15 total)
- [ ] adult.request.* permissions (7 total)

### User Preferences API
- [ ] GET /api/users/me/preferences
- [ ] PUT /api/users/me/preferences
- [ ] PIN protection implementation

### QAR Health Check
- [ ] Check adult module enabled status

### Adult Access Audit
- [ ] Log all QAR access (user_id, resource_type, resource_id, action, timestamp, ip)

### Documentation
- [ ] Split QUICKLIST.md into themed files
- [ ] Add cross-references to design docs
- [ ] Update TECH_STACK.md (add casbin, otel)
- [ ] Update CONFIGURATION.md (reflects pkg/config/)
- [ ] Review all docs/dev/design/ for accuracy

---

## Tech Stack Reference

| Component | Package | Notes |
|-----------|---------|-------|
| Cache (distributed) | `github.com/redis/rueidis` | NOT go-redis |
| Cache (local) | `github.com/maypok86/otter` v1.2.4 | W-TinyLFU |
| Cache (API) | `github.com/viccon/sturdyc` v1.1.5 | Request coalescing |
| Search | `github.com/typesense/typesense-go/v4` | NOT v3 |
| Config | `github.com/knadh/koanf/v2` | NOT viper |
| Logging | `log/slog` | NOT zap |
| Jobs | `github.com/riverqueue/river` | PostgreSQL-native |
| RBAC | `github.com/casbin/casbin/v2` | Dynamic roles |
| DI | `go.uber.org/fx` | Dependency injection |
| HTTP client | `github.com/go-resty/resty/v2` | External APIs |

---

## Build Commands

```bash
# With experiments
GOEXPERIMENT=greenteagc,jsonv2 go build -o bin/revenge ./cmd/revenge

# Generate code
sqlc generate
go generate ./api/...

# Lint
golangci-lint run

# Test
go test ./...

# Test with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

---

## Important Notes

**Adult Content** (Queen Anne's Revenge obfuscation):
- Schema: `qar` (isolated PostgreSQL schema)
- API namespace: `/api/v1/qar/*`
- Module location: `internal/content/qar/`
- See [ADULT_CONTENT_SYSTEM.md](docs/dev/design/features/adult/ADULT_CONTENT_SYSTEM.md)

**Design Docs are Source of Truth**:
- Only `docs/dev/design/` is authoritative
- Other documentation may be outdated
- Code must match design, not vice versa

**Milestone Completion Criteria**:
1. All structures present per design spec
2. `golangci-lint run` passes with no errors
3. Unit tests written and passing with >80% coverage
