---
applyTo: "**"
---

# Revenge Development Instructions Index

> Master reference for all development guidelines. Start here.

## Quick Links

| Need to...               | Go to                                                                           |
| ------------------------ | ------------------------------------------------------------------------------- |
| Understand architecture  | [ARCHITECTURE_V2.md](../../docs/dev/design/architecture/ARCHITECTURE_V2.md)     |
| Follow design principles | [DESIGN_PRINCIPLES.md](../../docs/dev/design/architecture/DESIGN_PRINCIPLES.md) |
| Check best practices     | [BEST_PRACTICES.md](../../docs/dev/design/operations/BEST_PRACTICES.md)         |
| See all external sources | [SOURCES.yaml](../../docs/dev/sources/SOURCES.yaml)                             |
| Find package docs        | [docs/dev/sources/](../../docs/dev/sources/)                                    |

---

## 🏗️ Core Patterns

| Instruction                                                        | Package                     | Live Docs                                                      |
| ------------------------------------------------------------------ | --------------------------- | -------------------------------------------------------------- |
| [fx-dependency-injection](fx-dependency-injection.instructions.md) | `go.uber.org/fx`            | [fx.md](../../docs/dev/sources/tooling/fx.md)                  |
| [koanf-configuration](koanf-configuration.instructions.md)         | `github.com/knadh/koanf/v2` | [koanf.md](../../docs/dev/sources/tooling/koanf.md)            |
| [go-features](go-features.instructions.md)                         | Go 1.25 stdlib              | [release-notes.md](../../docs/dev/sources/go/release-notes.md) |
| [testing-patterns](testing-patterns.instructions.md)               | `testing`, `testify`        | [testify.md](../../docs/dev/sources/testing/testify.md)        |

---

## 💾 Data & Storage

| Instruction                                            | Package          | Live Docs                                                                                            | Design Doc                                                                     |
| ------------------------------------------------------ | ---------------- | ---------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------ |
| [sqlc-database](sqlc-database.instructions.md)         | `sqlc`, `pgx/v5` | [sqlc.md](../../docs/dev/sources/database/sqlc.md), [pgx.md](../../docs/dev/sources/database/pgx.md) | —                                                                              |
| [migrations](migrations.instructions.md)               | `golang-migrate` | [migrations.md](../../docs/dev/sources/database/migrations.md)                                       | —                                                                              |
| [dragonfly-cache](dragonfly-cache.instructions.md)     | `rueidis`        | [rueidis.md](../../docs/dev/sources/tooling/rueidis.md)                                              | [DRAGONFLY.md](../../docs/dev/design/integrations/infrastructure/DRAGONFLY.md) |
| [otter-local-cache](otter-local-cache.instructions.md) | `otter`          | [otter.md](../../docs/dev/sources/tooling/otter.md)                                                  | —                                                                              |
| [sturdyc-api-cache](sturdyc-api-cache.instructions.md) | `sturdyc`        | [sturdyc.md](../../docs/dev/sources/tooling/sturdyc.md)                                              | —                                                                              |
| [typesense-search](typesense-search.instructions.md)   | `typesense-go`   | [typesense-go.md](../../docs/dev/sources/infrastructure/typesense-go.md)                             | [TYPESENSE.md](../../docs/dev/design/integrations/infrastructure/TYPESENSE.md) |

### Caching Strategy (Three-Tier)

```
┌─────────────────────────────────────────────────────────────────┐
│ Tier 1: otter (local)     → Hot data, µs latency, per-instance  │
│ Tier 2: rueidis (remote)  → Shared data, ms latency, distributed│
│ Tier 3: sturdyc (API)     → External APIs, request coalescing   │
└─────────────────────────────────────────────────────────────────┘
```

See [BEST_PRACTICES.md#caching](../../docs/dev/design/operations/BEST_PRACTICES.md) for details.

---

## 🌐 API & HTTP

| Instruction                                                            | Package           | Live Docs                                                   | Design Doc                                       |
| ---------------------------------------------------------------------- | ----------------- | ----------------------------------------------------------- | ------------------------------------------------ |
| [ogen-api](ogen-api.instructions.md)                                   | `ogen`            | [ogen.md](../../docs/dev/sources/tooling/ogen.md)           | [API.md](../../docs/dev/design/technical/API.md) |
| [resty-http-client](resty-http-client.instructions.md)                 | `resty/v3`        | [resty.md](../../docs/dev/sources/tooling/resty.md)         | —                                                |
| [websocket](websocket.instructions.md)                                 | `coder/websocket` | [websocket.md](../../docs/dev/sources/tooling/websocket.md) | —                                                |
| [revenge-api-compatibility](revenge-api-compatibility.instructions.md) | —                 | —                                                           | [API.md](../../docs/dev/design/technical/API.md) |

---

## 📦 Content Modules

| Instruction                                                      | Design Docs                                                                                                                                                                |
| ---------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [content-modules](content-modules.instructions.md)               | [LIBRARY_TYPES.md](../../docs/dev/design/features/LIBRARY_TYPES.md)                                                                                                        |
| [adult-modules](adult-modules.instructions.md)                   | [ADULT_CONTENT_SYSTEM.md](../../docs/dev/design/features/ADULT_CONTENT_SYSTEM.md), [WHISPARR_STASHDB_SCHEMA.md](../../docs/dev/design/features/WHISPARR_STASHDB_SCHEMA.md) |
| [metadata-providers](metadata-providers.instructions.md)         | [METADATA_SYSTEM.md](../../docs/dev/design/architecture/METADATA_SYSTEM.md)                                                                                                |
| [fsnotify-file-watching](fsnotify-file-watching.instructions.md) | [fsnotify.md](../../docs/dev/sources/tooling/fsnotify.md)                                                                                                                  |

### Metadata Provider Matrix

| Content Type | Primary     | Fallback     | Live Docs                                                                                                          |
| ------------ | ----------- | ------------ | ------------------------------------------------------------------------------------------------------------------ |
| Movies       | TMDb        | OMDb         | [tmdb.md](../../docs/dev/sources/apis/tmdb.md), [omdb.md](../../docs/dev/sources/apis/omdb.md)                     |
| TV Shows     | TMDb        | TheTVDB      | [tmdb.md](../../docs/dev/sources/apis/tmdb.md)                                                                     |
| Music        | MusicBrainz | Discogs      | [musicbrainz.md](../../docs/dev/sources/apis/musicbrainz.md), [discogs.md](../../docs/dev/sources/apis/discogs.md) |
| Books        | OpenLibrary | Google Books | [openlibrary.md](../../docs/dev/sources/apis/openlibrary.md)                                                       |
| Comics       | ComicVine   | Marvel API   | [comicvine.md](../../docs/dev/sources/apis/comicvine.md)                                                           |
| Adult        | StashDB     | ThePornDB    | [stashdb-schema.graphql](../../docs/dev/sources/apis/stashdb-schema.graphql)                                       |

See [docs/dev/design/integrations/metadata/](../../docs/dev/design/integrations/metadata/) for full provider docs.

---

## ⚙️ Services & Jobs

| Instruction                                                | Package | Live Docs                                                    | Design Doc                                                                     |
| ---------------------------------------------------------- | ------- | ------------------------------------------------------------ | ------------------------------------------------------------------------------ |
| [river-job-queue](river-job-queue.instructions.md)         | `river` | [river.md](../../docs/dev/sources/tooling/river.md)          | [RIVER.md](../../docs/dev/design/integrations/infrastructure/RIVER.md)         |
| [oidc-authentication](oidc-authentication.instructions.md) | —       | [oidc-core.md](../../docs/dev/sources/security/oidc-core.md) | [docs/dev/design/integrations/auth/](../../docs/dev/design/integrations/auth/) |
| [external-services](external-services.instructions.md)     | —       | —                                                            | [SCROBBLING.md](../../docs/dev/design/features/SCROBBLING.md)                  |

### Scrobbling Services

| Service      | Design Doc                                                                             | Live Docs                                                      |
| ------------ | -------------------------------------------------------------------------------------- | -------------------------------------------------------------- |
| Trakt        | [TRAKT.md](../../docs/dev/design/integrations/scrobbling/TRAKT.md)                     | [trakt.md](../../docs/dev/sources/apis/trakt.md)               |
| Last.fm      | [LASTFM_SCROBBLE.md](../../docs/dev/design/integrations/scrobbling/LASTFM_SCROBBLE.md) | [lastfm.md](../../docs/dev/sources/apis/lastfm.md)             |
| ListenBrainz | [LISTENBRAINZ.md](../../docs/dev/design/integrations/scrobbling/LISTENBRAINZ.md)       | [listenbrainz.md](../../docs/dev/sources/apis/listenbrainz.md) |

---

## 🎬 Playback & Streaming

| Instruction                                                          | Design Doc                                                                                                                                    |
| -------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------- |
| [streaming-best-practices](streaming-best-practices.instructions.md) | [AUDIO_STREAMING.md](../../docs/dev/design/technical/AUDIO_STREAMING.md)                                                                      |
| [player-architecture](player-architecture.instructions.md)           | [PLAYER_ARCHITECTURE.md](../../docs/dev/design/architecture/PLAYER_ARCHITECTURE.md)                                                           |
| [client-detection](client-detection.instructions.md)                 | [CLIENT_SUPPORT.md](../../docs/dev/design/features/CLIENT_SUPPORT.md)                                                                         |
| [disk-cache](disk-cache.instructions.md)                             | [BEST_PRACTICES.md#disk-cache](../../docs/dev/design/operations/BEST_PRACTICES.md)                                                            |
| [offloading-patterns](offloading-patterns.instructions.md)           | [OFFLOADING.md](../../docs/dev/design/technical/OFFLOADING.md), [BLACKBEARD.md](../../docs/dev/design/integrations/transcoding/BLACKBEARD.md) |
| [media-processing](media-processing.instructions.md)                 | go-astiav (FFmpeg), bimg, dhowden/tag, go-astisub, go-blurhash                                                                                |

---

## 🛡️ Resilience & Operations

| Instruction                                                | Design Doc                                                                            |
| ---------------------------------------------------------- | ------------------------------------------------------------------------------------- |
| [resilience-patterns](resilience-patterns.instructions.md) | [BEST_PRACTICES.md#resilience](../../docs/dev/design/operations/BEST_PRACTICES.md)    |
| [self-healing](self-healing.instructions.md)               | [BEST_PRACTICES.md#self-healing](../../docs/dev/design/operations/BEST_PRACTICES.md)  |
| [health-checks](health-checks.instructions.md)             | [BEST_PRACTICES.md#health-checks](../../docs/dev/design/operations/BEST_PRACTICES.md) |
| [hotreload](hotreload.instructions.md)                     | [BEST_PRACTICES.md#hot-reload](../../docs/dev/design/operations/BEST_PRACTICES.md)    |
| [lazy-initialization](lazy-initialization.instructions.md) | —                                                                                     |
| [observability](observability.instructions.md)             | —                                                                                     |

---

## 🎨 Frontend

| Instruction                                                    | Live Docs                                                                                                              |
| -------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------- |
| [frontend-architecture](frontend-architecture.instructions.md) | [svelte5.md](../../docs/dev/sources/frontend/svelte5.md), [sveltekit.md](../../docs/dev/sources/frontend/sveltekit.md) |

See [FRONTEND.md](../../docs/dev/design/technical/FRONTEND.md) for full frontend architecture.

---

## 📚 Reference Links

### Architecture & Design

- [ARCHITECTURE_V2.md](../../docs/dev/design/architecture/ARCHITECTURE_V2.md) - System architecture
- [DESIGN_PRINCIPLES.md](../../docs/dev/design/architecture/DESIGN_PRINCIPLES.md) - Core principles
- [TECH_STACK.md](../../docs/dev/design/technical/TECH_STACK.md) - Technology choices

### Operations

- [BEST_PRACTICES.md](../../docs/dev/design/operations/BEST_PRACTICES.md) - Patterns & practices
- [DEVELOPMENT.md](../../docs/dev/design/operations/DEVELOPMENT.md) - Dev environment
- [SETUP.md](../../docs/dev/design/operations/SETUP.md) - Production setup

### External Integrations

- [Servarr](../../docs/dev/design/integrations/servarr/) - Radarr, Sonarr, Lidarr, Whisparr
- [Metadata](../../docs/dev/design/integrations/metadata/) - TMDb, MusicBrainz, etc.
- [Scrobbling](../../docs/dev/design/integrations/scrobbling/) - Trakt, Last.fm, etc.
- [Auth/OIDC](../../docs/dev/design/integrations/auth/) - Authelia, Authentik, Keycloak

### Live Documentation Sources

All external docs are fetched to `docs/dev/sources/`:

| Category  | Path         | Contents                                        |
| --------- | ------------ | ----------------------------------------------- |
| Go stdlib | `go/stdlib/` | context, slog, net/http, testing                |
| Tooling   | `tooling/`   | fx, koanf, ogen, river, rueidis, otter, sturdyc |
| Database  | `database/`  | sqlc, pgx, PostgreSQL                           |
| APIs      | `apis/`      | TMDb, MusicBrainz, OpenAPI specs                |
| Frontend  | `frontend/`  | Svelte 5, SvelteKit, shadcn                     |
| Security  | `security/`  | OIDC, OAuth2, JWT RFCs                          |
| Protocols | `protocols/` | HLS, DASH, HTTP Range                           |

Run `python scripts/fetch-sources.py` to update.

---

## DO's and DON'Ts Summary

### DO ✅

- Use `context.Context` as first parameter
- Use `slog` for logging
- Use Go 1.22+ HTTP routing patterns
- Use `sync.WaitGroup.Go` (Go 1.25)
- Keep modules isolated
- Use River for background jobs
- Use ogen for API handlers
- Follow three-tier caching (otter → rueidis → sturdyc)

### DON'T ❌

- Use `init()` - use fx constructors
- Use global variables - inject dependencies
- Use `panic` for errors
- Use deprecated packages:
  - gorilla/mux, gorilla/websocket → use stdlib, coder/websocket
  - viper → use koanf
  - go-redis/v9 → use rueidis
  - ristretto → use otter
  - zap, logrus → use slog
- Share tables between content modules
- Transcode internally - use Blackbeard

---

## Version Reference

| Component  | Version | Notes                         |
| ---------- | ------- | ----------------------------- |
| Go         | 1.25    | Container-aware, WaitGroup.Go |
| PostgreSQL | 18+     | Latest stable                 |
| Dragonfly  | 1.36+   | Redis-compatible              |
| Typesense  | 30.0    | Full-text search              |

See [VERSION_POLICY.md](../../docs/dev/design/planning/VERSION_POLICY.md) for update policy.
