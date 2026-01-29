# Revenge - Development Roadmap

> Modular media server with complete content isolation

**Last Updated**: 2026-01-29
**Current Phase**: Content Modules
**Build**: `GOEXPERIMENT=greenteagc,jsonv2 go build ./...`

---

## Quick Status

```
Foundation (Week 1-2)     ████████████████████████ 100%
Movie Module              ████████████████████████ 100%
TV Shows Module           ██████████████████████░░  90% (missing: API handlers)
Design Docs               ████████████████████████ 100%
Music Module              ░░░░░░░░░░░░░░░░░░░░░░░░   0%
Books Module              ░░░░░░░░░░░░░░░░░░░░░░░░   0%
Comics Module             ░░░░░░░░░░░░░░░░░░░░░░░░   0%
Frontend                  ░░░░░░░░░░░░░░░░░░░░░░░░   0%
```

---

## Current Sprint

### In Progress: TV Shows Module Completion

- [x] Database migrations (series, seasons, episodes, credits)
- [x] sqlc queries (7 files, 100+ queries)
- [x] Entity definitions
- [x] Repository (PostgreSQL implementation)
- [x] Service layer
- [x] **module.go** - FX dependency injection registration
- [x] **jobs.go** - River workers (metadata enrichment)
- [x] **metadata_provider.go** - Provider interface + adapters
- [ ] API handlers

### Next: Music Module

- [ ] Database migrations (`music_artists`, `music_albums`, `music_tracks`)
- [ ] sqlc queries
- [ ] Entity/Repository/Service
- [ ] Lidarr integration
- [ ] MusicBrainz fallback

---

## Completed

### Foundation (Week 1-2)
- PostgreSQL + sqlc type-safe queries
- Dragonfly cache (rueidis client, 14x faster than go-redis)
- Typesense search (typesense-go/v4)
- River job queue (PostgreSQL-native)
- uber-go/fx dependency injection
- koanf v2 configuration
- slog structured logging
- failsafe-go resilience (circuit breakers, retries)
- otter local cache (W-TinyLFU)
- Health checks + graceful shutdown
- RBAC with Casbin (dynamic roles)
- Session management + OIDC support
- OpenAPI spec + ogen code generation

### Movie Module
- Full CRUD with relations (genres, cast, crew, studios)
- User data (ratings, favorites, watchlist, watch history)
- TMDb metadata provider
- Radarr integration ready
- River jobs for metadata enrichment
- 3-tier caching (local, API, distributed)

### Design Documentation
- [DESIGN_DOC_TEMPLATE.md](docs/dev/design/DESIGN_DOC_TEMPLATE.md) - Standard template
- [WATCH_NEXT_CONTINUE_WATCHING.md](docs/dev/design/features/WATCH_NEXT_CONTINUE_WATCHING.md) - Playback continuation
- [RELEASE_CALENDAR.md](docs/dev/design/features/RELEASE_CALENDAR.md) - Servarr calendar integration
- [METADATA_SYSTEM.md](docs/dev/design/architecture/METADATA_SYSTEM.md) - Servarr-first with fallbacks
- [CONTENT_RATING.md](docs/dev/design/features/CONTENT_RATING.md) - Module-specific age restrictions
- All docs audited for consistency (rueidis, typesense-go/v4, Chaptarr)

---

## Content Modules Roadmap

| Module | Status | People Table | Servarr |
|--------|--------|--------------|---------|
| Movies | Done | `video_people` (shared) | Radarr |
| TV Shows | 90% | `video_people` (shared) | Sonarr |
| Music | Pending | `music_artists` | Lidarr |
| Audiobooks | Pending | `book_authors` | Chaptarr |
| Books | Pending | `book_authors` | Chaptarr |
| Podcasts | Pending | - | RSS |
| Comics | Pending | `comic_creators` | - |
| Photos | Pending | - | - |
| LiveTV | Pending | - | TVHeadend/NextPVR |
| Adult | Scaffolded | `c.performers` | Whisparr v3 (eros) |

---

## External Integrations

### Metadata Providers
- [ ] TMDb (movies, TV) - partially done
- [ ] MusicBrainz (music)
- [ ] AniList/MyAnimeList (anime)
- [ ] OpenLibrary/Hardcover (books)
- [ ] ComicVine (comics)
- [ ] StashDB (adult, schema `c`)

### Servarr Ecosystem
- [ ] Radarr - Movie management
- [ ] Sonarr - TV show management
- [ ] Lidarr - Music management
- [ ] Chaptarr - Books & audiobooks (Readarr API)
- [ ] Whisparr v3 (eros) - Adult content (schema `c`)

### Scrobbling
- [ ] Trakt - Movie/TV sync
- [ ] Last.fm - Music scrobbling
- [ ] ListenBrainz - Music scrobbling
- [ ] Letterboxd - Import only

---

## Feature Enhancements (P2)

- [ ] Watch Next / Continue Watching (design done)
- [ ] Release Calendar (design done)
- [ ] Request System with Polls (design done)
- [ ] i18n System
- [ ] Analytics Service (Year in Review)
- [ ] Profiles System (Netflix-style)
- [ ] Media Enhancements (trickplay, intro detection, chapters)

---

## Frontend (P3)

- [ ] SvelteKit 2 + Tailwind CSS 4
- [ ] shadcn-svelte components
- [ ] TanStack Query
- [ ] Video player (Shaka + hls.js)
- [ ] Audio player (Howler.js, gapless)
- [ ] Admin panel

---

## Tech Stack Reference

| Component | Package | Notes |
|-----------|---------|-------|
| Cache (distributed) | `github.com/redis/rueidis` | NOT go-redis |
| Cache (local) | `github.com/maypok86/otter` v1.2.4 | W-TinyLFU |
| Search | `github.com/typesense/typesense-go/v4` | NOT v3 |
| Config | `github.com/knadh/koanf/v2` | NOT viper |
| Logging | `log/slog` | NOT zap |
| Jobs | `github.com/riverqueue/river` | PostgreSQL-native |
| Resilience | `github.com/failsafe-go/failsafe-go` | Circuit breakers |
| DI | `go.uber.org/fx` | Dependency injection |

---

## Important Notes

**Adult Content**:
- Schema: `c` (isolated PostgreSQL schema)
- API namespace: `/c/*`
- Module location: `internal/content/c/`

**External Transcoding**:
- Blackbeard service handles all transcoding
- Revenge proxies streams only

**Build Commands**:
```bash
# With experiments
GOEXPERIMENT=greenteagc,jsonv2 go build -o bin/revenge ./cmd/revenge

# Generate code
sqlc generate
go generate ./api/...

# Lint
golangci-lint run
```
