# Revenge - Project TODO

> Modular media server with complete content isolation

## Architecture

See [docs/ARCHITECTURE_V2.md](docs/ARCHITECTURE_V2.md) for the complete modular design.
See [docs/MODULE_IMPLEMENTATION_TODO.md](docs/MODULE_IMPLEMENTATION_TODO.md) for detailed implementation phases.

## Implementation Phases

### Phase 1: Core Infrastructure ⬜ IN PROGRESS
- [ ] Shared tables (users, sessions, libraries, api_keys, server_settings)
- [ ] Video playlists (movie + tvshow pool)
- [ ] Audio playlists (music + audiobook + podcast pool)
- [ ] Video collections
- [ ] Audio collections

### Phase 2: Movie Module ⬜ NOT STARTED
- [ ] Database schema (movies, genres, people, studios, images)
- [ ] Domain entities
- [ ] Repository layer (sqlc)
- [ ] Service layer
- [ ] HTTP handlers
- [ ] User data (ratings, history, favorites)

### Phase 3: TV Show Module ⬜ NOT STARTED
- [ ] Database schema (series, seasons, episodes)
- [ ] Domain/Repository/Service/Handlers
- [ ] User data

### Phase 4: Music Module ⬜ NOT STARTED
- [ ] Database schema (artists, albums, tracks, music_videos)
- [ ] Domain/Repository/Service/Handlers
- [ ] User data

### Phase 5: Remaining Modules ⬜ NOT STARTED
- [ ] Audiobook module
- [ ] Book module
- [ ] Podcast module
- [ ] Photo module
- [ ] LiveTV module

### Phase 6: Adult Modules ⬜ NOT STARTED
- [ ] `adult` PostgreSQL schema
- [ ] Adult movie module
- [ ] Adult show module
- [ ] Shared performers/studios/tags
- [ ] Adult playlists & collections

---

## Go 1.25 Features to Adopt

- [ ] `sync.WaitGroup.Go` - Replace manual wg.Add/Done patterns
- [ ] `testing/synctest` - Concurrent code testing
- [ ] `net/http.CrossOriginProtection` - Replace custom CSRF
- [ ] `slog.GroupAttrs` - Grouped logging
- [ ] `runtime/trace.FlightRecorder` - Observability
- [ ] `reflect.TypeAssert` - Zero-allocation type assertions

## Experimental (Evaluate)

- [ ] `GOEXPERIMENT=greenteagc` - New GC (10-40% reduction)
- [ ] `GOEXPERIMENT=jsonv2` - Faster JSON

---

## Completed ✅

- [x] Project setup (Go 1.25, fx, koanf, sqlc)
- [x] CI/CD (GitHub Actions, release-please)
- [x] Docker Compose (dev + prod)
- [x] Configuration system (REVENGE_* env vars)
- [x] Logging (slog)
- [x] HTTP server with graceful shutdown
- [x] Health endpoints
- [x] Basic auth middleware
- [x] User/Session/OIDC tables
- [x] Genre domain separation
- [x] Rename to Revenge

