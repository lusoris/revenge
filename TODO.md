# Revenge - Project TODO

> Modular media server with complete content isolation

## Architecture

See [docs/ARCHITECTURE_V2.md](docs/ARCHITECTURE_V2.md) for the complete modular design.

---

## Implementation Phases

### Phase 1: Core Infrastructure ⬜ IN PROGRESS

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
- [ ] Shared tables (libraries, api_keys, server_settings, activity_log)
- [ ] River job queue setup
- [ ] Typesense search client
- [ ] Dragonfly cache client

### Phase 2: Movie Module ⬜ NOT STARTED

- [ ] Database schema (movies, genres, people, studios, images, streams)
- [ ] Domain entities
- [ ] Repository layer (sqlc)
- [ ] Service layer
- [ ] HTTP handlers (ogen)
- [ ] User data (ratings, history, favorites, watchlist)
- [ ] Radarr integration
- [ ] TMDb fallback provider

### Phase 3: TV Show Module ⬜ NOT STARTED

- [ ] Database schema (series, seasons, episodes)
- [ ] Domain/Repository/Service/Handlers
- [ ] User data
- [ ] Sonarr integration
- [ ] TheTVDB/TMDb fallback

### Phase 4: Music Module ⬜ NOT STARTED

- [ ] Database schema (artists, albums, tracks, music_videos)
- [ ] Domain/Repository/Service/Handlers
- [ ] User data
- [ ] Lidarr integration
- [ ] MusicBrainz/Last.fm fallback

### Phase 5: Playback Service ⬜ NOT STARTED

- [ ] Session management
- [ ] Client capability detection
- [ ] Blackbeard transcoder integration
- [ ] Stream buffering
- [ ] Progress tracking
- [ ] Bandwidth adaptation

### Phase 6: Remaining Content Modules ⬜ NOT STARTED

- [ ] Audiobook module (Audiobookshelf integration)
- [ ] Book module (Audiobookshelf + Chaptarr)
- [ ] Podcast module (Audiobookshelf + RSS)
- [ ] Photo module
- [ ] LiveTV module (PVR backends)
- [ ] Collection module (video + audio pools)

### Phase 7: Adult Modules ⬜ NOT STARTED

- [ ] `c` PostgreSQL schema (isolated)
- [ ] Adult movie module
- [ ] Adult show module
- [ ] Shared performers/studios/tags
- [ ] Adult playlists & collections
- [ ] Whisparr integration
- [ ] Stash/StashDB integration

### Phase 8: Media Enhancements ⬜ NOT STARTED

- [ ] Trailer system (local, Radarr, TMDb, YouTube)
- [ ] Audio themes (Netflix-style hover music)
- [ ] Intro/outro detection (Chromaprint)
- [ ] Trickplay generation
- [ ] Chapter extraction
- [ ] Cinema mode (preroll/postroll)

### Phase 9: External Services ⬜ NOT STARTED

- [ ] Seerr integration + adapter
- [ ] Trakt scrobbling
- [ ] Last.fm scrobbling
- [ ] ListenBrainz scrobbling
- [ ] Import/export ratings

### Phase 10: Frontend ⬜ NOT STARTED

- [ ] SvelteKit 2 setup
- [ ] Tailwind CSS 4 + shadcn-svelte
- [ ] Authentication (JWT + OIDC)
- [ ] Media browser
- [ ] Video player
- [ ] Audio player
- [ ] Admin panel

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

## Documentation Status

### Completed ✅

- [x] ARCHITECTURE_V2.md - Complete modular design
- [x] TECH_STACK.md - Technology choices
- [x] PROJECT_STRUCTURE.md - Directory layout
- [x] METADATA_SYSTEM.md - Servarr-first with Audiobookshelf/Seerr
- [x] AUDIO_STREAMING.md - Progress, bandwidth adaptation
- [x] CLIENT_SUPPORT.md - Chromecast, DLNA, capabilities
- [x] MEDIA_ENHANCEMENTS.md - Trailers, themes, intros, trickplay, Live TV
- [x] SCROBBLING.md - External service sync
- [x] OFFLOADING.md - Blackbeard integration
- [x] BEST_PRACTICES.md - Resilience patterns
- [x] I18N.md - Internationalization

### TODO 📝

- [ ] ADULT_METADATA.md - Stash/StashDB/Whisparr integration
- [ ] CINEMA_MODE.md - Preroll, postroll, intermission
- [ ] API.md - OpenAPI design guidelines
- [ ] REVERSE_PROXY.md - Nginx, Caddy, Traefik configs
- [ ] MOBILE_APPS.md - iOS/Android architecture

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
- [x] Resilience packages (circuit breaker, bulkhead, retry)
- [x] Supervisor/graceful shutdown packages
- [x] Health check system
- [x] Hot reload configuration
- [x] Lazy initialization patterns
- [x] Metrics package
- [x] Playback service architecture (docs)
- [x] Documentation suite
