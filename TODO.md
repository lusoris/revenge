# Revenge - Project TODO

> Modular media server with complete content isolation

**Last Updated**: 2026-01-28

---

## Current Status

| Area | Status | Notes |
|------|--------|-------|
| Core Infrastructure | 🟢 85% | fx, koanf, sqlc, health, graceful shutdown |
| Shared Services | 🟢 80% | auth, user, session, oidc, library, genre, rating |
| Playback Service | 🟢 100% | Complete architecture, Blackbeard integration |
| Content Modules | 🔴 0% | Not started - all empty folders |
| Frontend | 🔴 0% | Not started |

---

## Implementation Phases

### Phase 1: Core Infrastructure ✅ MOSTLY COMPLETE

- [x] Project setup (Go 1.25, fx, koanf, sqlc)
- [x] CI/CD (GitHub Actions, release-please)
- [x] Docker Compose (dev + prod)
- [x] Configuration system (REVENGE_* env vars)
- [x] Logging (slog)
- [x] HTTP server with graceful shutdown
- [x] Health endpoints
- [x] Basic auth middleware
- [x] User/Session/OIDC tables + migrations
- [x] Genre domain separation
- [x] Resilience packages (circuit breaker, bulkhead, retry)
- [x] Supervisor/graceful shutdown packages
- [x] Health check system
- [x] Hot reload configuration
- [x] Lazy initialization patterns
- [x] Metrics package
- [ ] River job queue integration in main.go
- [ ] Typesense search client integration
- [ ] Dragonfly cache client (migrate to rueidis)

### Phase 2: Movie Module ⬜ NOT STARTED

- [ ] Database migrations (movies, genres, people, studios, images, streams)
- [ ] sqlc queries
- [ ] Domain entities
- [ ] Repository layer
- [ ] Service layer
- [ ] HTTP handlers (ogen)
- [ ] User data (ratings, history, favorites, watchlist)
- [ ] Radarr integration
- [ ] TMDb fallback provider
- [ ] River jobs (scan, metadata, images)
- [ ] Typesense indexing

### Phase 3: TV Show Module ⬜ NOT STARTED

- [ ] Database migrations (series, seasons, episodes)
- [ ] sqlc queries + Domain/Repository/Service/Handlers
- [ ] User data
- [ ] Sonarr integration
- [ ] TheTVDB/TMDb fallback

### Phase 4: Music Module ⬜ NOT STARTED

- [ ] Database migrations (artists, albums, tracks, music_videos)
- [ ] Full stack implementation
- [ ] Lidarr integration
- [ ] MusicBrainz/Last.fm fallback

### Phase 5: Remaining Content Modules ⬜ NOT STARTED

- [ ] Audiobook module (Audiobookshelf integration)
- [ ] Book module
- [ ] Podcast module
- [ ] Photo module
- [ ] LiveTV module (TVHeadend/NextPVR backends)
- [ ] Comics module (ComicVine/Marvel/GCD)
- [ ] Collection module (cross-module pools)

### Phase 6: Adult Modules ⬜ NOT STARTED

- [ ] `c` PostgreSQL schema (isolated)
- [ ] Adult movie module
- [ ] Adult show module
- [ ] Shared performers/studios/tags
- [ ] Whisparr integration
- [ ] StashDB/ThePornDB metadata

### Phase 7: Media Enhancements ⬜ NOT STARTED

- [ ] Trailer system
- [ ] Audio themes
- [ ] Intro/outro detection
- [ ] Trickplay generation
- [ ] Chapter extraction

### Phase 8: External Scrobbling ⬜ NOT STARTED

- [ ] Trakt sync
- [ ] Last.fm scrobbling
- [ ] ListenBrainz scrobbling

### Phase 9: Frontend ⬜ NOT STARTED

- [ ] SvelteKit 2 setup
- [ ] Tailwind CSS 4 + shadcn-svelte
- [ ] Media browser
- [ ] Video/Audio players
- [ ] Admin panel

---

## Go 1.25 Features to Adopt

- [ ] `sync.WaitGroup.Go` - Replace manual wg.Add/Done patterns
- [ ] `testing/synctest` - Concurrent code testing
- [ ] `net/http.CrossOriginProtection` - Replace custom CSRF
- [ ] Built-in container CPU/memory awareness (replaces automaxprocs)

---

## References

- [Architecture](docs/dev/design/architecture/ARCHITECTURE_V2.md)
- [Module Implementation TODO](docs/dev/design/planning/MODULE_IMPLEMENTATION_TODO.md)
- [Go Packages Research](docs/dev/design/research/GO_PACKAGES_RESEARCH.md)
