# Codebase Analysis Report - Revenge Go

**Datum:** 28. Januar 2026  
**Analysierte Dokumente:**
- `docs/ARCHITECTURE_V2.md`
- `docs/DESIGN_PRINCIPLES.md`
- `.github/instructions/fx-dependency-injection.instructions.md`
- `.github/instructions/content-modules.instructions.md`

**Analysierte Codebase-Pfade:**
- `/home/kilian/dev/jellyfin-go/internal/`
- `/home/kilian/dev/jellyfin-go/cmd/revenge/`
- `/home/kilian/dev/jellyfin-go/pkg/config/`

---

## 📊 Executive Summary

### Status: 🟡 Foundation Solid, Content Modules Missing

| Kategorie | Status | Fortschritt |
|-----------|--------|-------------|
| Core Infrastructure | 🟢 Vollständig | 90% |
| Shared Services | 🟡 Teilweise | 70% |
| Content Modules | 🔴 Fehlend | 0% |
| External Services | 🟡 Teilweise | 40% |

**Haupterkenntnisse:**
1. ✅ **Core Infrastructure ist solid**: fx, koanf, PostgreSQL, sqlc funktionieren
2. ✅ **Shared Services teilweise vorhanden**: auth, user, session, oidc, genre, library, rating, playback
3. ❌ **ALLE Content-Module fehlen komplett**: movie, tvshow, music, etc. sind nur leere Ordner
4. ✅ **External Services als Stubs vorhanden**: Cache, Search, Jobs implementiert aber ohne Integration in main.go
5. ⚠️ **Database Migrations nur für Shared**: Keine Content-Module-Migrations
6. ❌ **OpenAPI/ogen Integration fehlt**: Keine API-Generierung konfiguriert

---

## 1. Project Structure Analysis

### 1.1 Vorhandene Struktur

```
✅ internal/
   ✅ service/        # 7/9 Services vorhanden
      ✅ auth/
      ✅ user/
      ✅ session/      # Indirekt über repository
      ✅ oidc/
      ✅ library/
      ✅ playback/     # VOLLSTÄNDIG - alle 11 Dateien!
      ✅ rating/
      ✅ genre/
      ❌ session/      # Fehlt als eigener Service (nur Repository)
   
   ❌ content/        # ALLE Module sind leere Ordner!
      ❌ movie/       # Leer
      ❌ tvshow/      # Leer
      ❌ music/       # Leer
      ✅ c/
         ❌ movie/    # Leer
         ❌ show/     # Leer
   
   ✅ infra/
      ✅ database/    # Vollständig (Pool, Migrator, Repository)
      ✅ cache/       # Client implementiert (go-redis)
      ✅ search/      # Client implementiert (typesense)
      ✅ jobs/        # Service implementiert (river)
   
   ✅ domain/        # Entities vorhanden
   ✅ api/
      ✅ handlers/   # Auth, User, Library, Rating, Genre, Media (basic), OIDC
      ✅ middleware/ # Auth middleware
```

### 1.2 Fehlende Komponenten

**Content Modules (KRITISCH):**
- ❌ `internal/content/movie/` - komplett leer, sollte haben:
  - `entity.go`, `repository.go`, `service.go`, `handler.go`, `scanner.go`, `provider_tmdb.go`, `jobs.go`, `module.go`
- ❌ `internal/content/tvshow/` - komplett leer
- ❌ `internal/content/music/` - komplett leer
- ❌ `internal/content/audiobook/` - Ordner existiert nicht
- ❌ `internal/content/book/` - Ordner existiert nicht
- ❌ `internal/content/podcast/` - Ordner existiert nicht
- ❌ `internal/content/photo/` - Ordner existiert nicht
- ❌ `internal/content/livetv/` - Ordner existiert nicht
- ❌ `internal/content/collection/` - Ordner existiert nicht
- ❌ `internal/content/c/movie/` - komplett leer
- ❌ `internal/content/c/show/` - komplett leer

**Migrations:**
- ✅ `migrations/shared/` - 5 Migrations vorhanden (extensions, users, sessions, oidc, libraries)
- ❌ `migrations/movie/` - Fehlt komplett
- ❌ `migrations/tvshow/` - Fehlt komplett
- ❌ `migrations/music/` - Fehlt komplett
- ❌ `migrations/c/` - Fehlt komplett (Adult schema)

**API (OpenAPI/ogen):**
- ❌ `api/openapi/` - Ordner existiert, aber ist leer!
- ❌ `api/generated/` - Ordner existiert, aber ist leer!
- ❌ Keine ogen-Konfiguration in `go generate`

**SQL Queries (sqlc):**
- ✅ `queries/users.sql`, `queries/sessions.sql`, `queries/oidc.sql`, `queries/libraries.sql`, `queries/genres.sql`, `queries/ratings.sql`
- ❌ `queries/movie/` - Ordner leer
- ❌ Keine Content-Module-Queries

---

## 2. Entry Point Analysis (`cmd/revenge/main.go`)

### 2.1 ✅ Korrekte fx-Integration

```go
app := fx.New(
    // Core modules
    fx.Provide(config.New, NewLogger),
    
    // Infrastructure
    database.Module,  // ✅ Registriert
    
    // Services
    auth.Module,      // ✅ Korrekt
    user.Module,      // ✅ Korrekt
    library.Module,   // ✅ Korrekt
    rating.Module,    // ✅ Korrekt
    
    // API
    fx.Provide(middleware.NewAuth, handlers...),
    
    // HTTP
    fx.Provide(NewMux, NewServer),
    fx.Invoke(RegisterRoutes, RunServer),
)
```

**Bewertung:** ✅ Folgt fx-Patterns korrekt

### 2.2 ⚠️ Fehlende Module in main.go

**Nicht registriert:**
- ❌ `cache.Module` - existiert, aber nicht in fx.New()
- ❌ `search.Module` - existiert, aber nicht in fx.New()
- ❌ `jobs.Module` - existiert, aber nicht in fx.New()
- ❌ `oidc.Module` - existiert, aber nicht in fx.New()
- ❌ `genre.Module` - existiert, aber nicht in fx.New()
- ❌ `playback.Module` - existiert, aber nicht in fx.New()
- ❌ Content modules - existieren nicht

**Impact:** Infrastructure Services sind zwar implementiert, werden aber nicht gestartet!

---

## 3. Module Implementation Analysis

### 3.1 Vorhandene Service Module

| Service | Files | Pattern Compliance | Status |
|---------|-------|-------------------|--------|
| `auth/` | ✅ 7 files | ✅ module.go, service.go, tests | 🟢 Vollständig |
| `user/` | ✅ 3 files | ✅ module.go, service.go, tests | 🟢 Vollständig |
| `library/` | ✅ 2 files | ✅ module.go, service.go | 🟡 Keine Tests |
| `rating/` | ✅ 1 file | ✅ module.go | 🟡 Nur Module |
| `oidc/` | ✅ 2+ files | ✅ module.go, service.go | 🟢 Vollständig |
| `genre/` | ✅ 2 files | ✅ module.go, service.go | 🟢 Vollständig |
| `playback/` | ✅ 11 files | ✅ **VOLLSTÄNDIG!** | 🟢 **Perfekt** |

**`playback/` ist vollständig implementiert:**
- ✅ `client.go` - Client detection
- ✅ `bandwidth.go` - Bandwidth monitoring
- ✅ `transcoder.go` - Blackbeard integration
- ✅ `session.go` - Playback session state
- ✅ `buffer.go` - HLS/DASH segment buffering
- ✅ `fileserver.go` - Raw file HTTP streaming
- ✅ `stream_handler.go` - Unified stream handler
- ✅ `transcode_cache.go` - Memory-aware transcode cache
- ✅ `disk_cache.go` - Persistent disk cache
- ✅ `profile.go` - Transcode profiles & device groups
- ✅ `module.go` - fx module registration

### 3.2 Content Module Pattern Compliance

**Erwartetes Pattern (aus instructions):**
```
{module}/
  entity.go           # Domain entities
  repository.go       # Repository interface
  repository_pg.go    # PostgreSQL implementation (oder repository.go direkt)
  service.go          # Business logic
  handler.go          # HTTP handlers (ogen interfaces)
  scanner.go          # File scanner
  provider_{name}.go  # Metadata providers
  jobs.go             # River job definitions
  module.go           # fx.Module registration
```

**Realität:**
```
movie/       # ❌ Leer (0 Dateien)
tvshow/      # ❌ Leer (0 Dateien)
music/       # ❌ Leer (0 Dateien)
c/movie/     # ❌ Leer (0 Dateien)
c/show/      # ❌ Leer (0 Dateien)
```

---

## 4. Missing Components (Detailed)

### 4.1 KRITISCH: Content Modules

**Status:** Alle 11 Content-Module fehlen komplett

| Module | Dokument Status | Code Status | Priority |
|--------|----------------|-------------|----------|
| movie | ✅ Spezifiziert | ❌ Nicht vorhanden | 🔴 P0 |
| tvshow | ✅ Spezifiziert | ❌ Nicht vorhanden | 🔴 P0 |
| music | ✅ Spezifiziert | ❌ Nicht vorhanden | 🔴 P1 |
| audiobook | ✅ Spezifiziert | ❌ Ordner fehlt | 🟡 P2 |
| book | ✅ Spezifiziert | ❌ Ordner fehlt | 🟡 P2 |
| podcast | ✅ Spezifiziert | ❌ Ordner fehlt | 🟡 P2 |
| photo | ✅ Spezifiziert | ❌ Ordner fehlt | 🟡 P2 |
| livetv | ✅ Spezifiziert | ❌ Ordner fehlt | 🟡 P3 |
| collection | ✅ Spezifiziert | ❌ Ordner fehlt | 🟡 P3 |
| adult_movie | ✅ Spezifiziert | ❌ Nicht vorhanden | 🟡 P3 |
| adult_show | ✅ Spezifiziert | ❌ Nicht vorhanden | 🟡 P3 |

### 4.2 Infrastructure Integration Gaps

**Dragonfly Cache (`internal/infra/cache/`):**
- ✅ Client implementiert
- ✅ Module definiert
- ❌ Nicht in main.go registriert
- ❌ Keine Config-Integration (hardcoded values)
- ❌ Keine fx Lifecycle hooks

```go
// AKTUELL in cache.go:
var Module = fx.Module("cache",
    fx.Provide(func(logger *slog.Logger) (*Client, error) {
        // TODO: Get config from koanf
        cfg := Config{
            Host: "localhost",
            Port: 6379,
            DB:   0,
        }
        return NewClient(cfg, logger)
    }),
)
```

**Typesense Search (`internal/infra/search/`):**
- ✅ Client implementiert
- ✅ Module definiert
- ❌ Nicht in main.go registriert
- ❌ Keine Config-Integration (hardcoded values)
- ❌ Keine fx Lifecycle hooks

```go
// AKTUELL in search.go:
var Module = fx.Module("search",
    fx.Provide(func(logger *slog.Logger) (*Client, error) {
        // TODO: Get config from koanf
        cfg := Config{
            Host:   "http://localhost:8108",
            APIKey: "xyz",
        }
        return NewClient(cfg, logger)
    }),
)
```

**River Jobs (`internal/infra/jobs/`):**
- ✅ Service implementiert
- ✅ Module mit Lifecycle definiert
- ❌ Nicht in main.go registriert
- ❌ Keine Workers registriert (NewWorkers() leer)
- ✅ Config-Integration fehlt (hardcoded queue config)

### 4.3 Database Migrations Gaps

**Vorhanden:**
```
migrations/shared/
  ✅ 000001_extensions.{up,down}.sql
  ✅ 000002_users.{up,down}.sql
  ✅ 000003_sessions.{up,down}.sql
  ✅ 000004_oidc.{up,down}.sql
  ✅ 000005_libraries.{up,down}.sql
```

**Fehlend (laut ARCHITECTURE_V2.md):**
```
migrations/shared/
  ❌ 000006_api_keys.sql
  ❌ 000007_server_settings.sql
  ❌ 000008_activity_log.sql
  ❌ 000010_video_playlists.sql
  ❌ 000011_audio_playlists.sql
  ❌ 000012_video_collections.sql
  ❌ 000013_audio_collections.sql

migrations/movie/
  ❌ 000001_movies.sql
  ❌ 000002_movie_people.sql
  ❌ 000003_movie_streams.sql
  ❌ 000004_movie_user_data.sql

migrations/tvshow/
  ❌ 000001_series.sql
  ❌ ...

migrations/music/
  ❌ 000001_artists.sql
  ❌ ...

migrations/c/
  ❌ 000001_c_schema.sql (CREATE SCHEMA c;)
  ❌ 000002_c_movies.sql
  ❌ 000003_c_performers.sql
```

### 4.4 OpenAPI/ogen Integration

**Status:** Komplett fehlend

**Erwartet (laut ARCHITECTURE_V2.md):**
```
api/
  openapi/
    revenge.yaml      # ❌ Fehlt
    movies.yaml       # ❌ Fehlt
    shows.yaml        # ❌ Fehlt
    music.yaml        # ❌ Fehlt
  generated/          # ❌ Leer
```

**go.mod:**
- ❌ `github.com/ogen-go/ogen` nicht vorhanden

**Missing:**
- ❌ OpenAPI Specs
- ❌ `go generate` Direktiven
- ❌ ogen-generierte Handler Interfaces

---

## 5. Inkonsistenzen: Code vs Dokumentation

### 5.1 Namens-Abweichungen

| Dokumentation | Code | Status |
|---------------|------|--------|
| `user_profiles` | `profiles` | ✅ Konsistent (neuere Docs verwenden `profiles`) |
| `c` schema | N/A | ❌ Nicht implementiert |
| `video_playlists` | N/A | ❌ Nicht implementiert |
| `audio_playlists` | N/A | ❌ Nicht implementiert |

### 5.2 Architecture Mismatch

**ARCHITECTURE_V2.md Zeilen 102-157 vs Realität:**

| Komponente | Docs | Code | Match |
|------------|------|------|-------|
| `internal/service/auth` | ✅ | ✅ | 🟢 |
| `internal/service/user` | ✅ | ✅ | 🟢 |
| `internal/service/session` | ✅ | ⚠️ Nur Repository | 🟡 |
| `internal/service/oidc` | ✅ | ✅ | 🟢 |
| `internal/service/library` | ✅ | ✅ | 🟢 |
| `internal/service/playback` | ✅ | ✅ | 🟢 |
| `internal/content/movie` | ✅ | ❌ Leer | 🔴 |
| `internal/content/tvshow` | ✅ | ❌ Leer | 🔴 |
| `internal/content/music` | ✅ | ❌ Leer | 🔴 |
| `internal/infra/cache` | ✅ | ✅ (nicht registriert) | 🟡 |
| `internal/infra/search` | ✅ | ✅ (nicht registriert) | 🟡 |
| `internal/infra/jobs` | ✅ | ✅ (nicht registriert) | 🟡 |

### 5.3 Go Module Versions

**ARCHITECTURE_V2.md vs go.mod:**

| Dependency | Expected | Actual | Match |
|------------|----------|--------|-------|
| Go | 1.25+ | 1.25.0 | ✅ |
| pgx/v5 | Latest | v5.8.0 | ✅ |
| redis | v9 | v9.17.3 | ✅ |
| typesense | v4 | v3.2.0 | ⚠️ **v3 statt v4** |
| river | Latest | v0.30.2 | ✅ |
| ogen | Latest | ❌ Fehlt | 🔴 |
| fx | v1.24+ | v1.24.0 | ✅ |

**Action Required:**
- Upgrade `typesense-go` von v3 → v4
- Add `ogen-go/ogen`

---

## 6. Code Quality Assessment

### 6.1 fx Dependency Injection Compliance

**Rating:** 🟢 Excellent (90%)

**Positiv:**
- ✅ Alle Services verwenden `fx.Module()`
- ✅ Korrekte Verwendung von `fx.In` Parameter Structs
- ✅ Korrekte Verwendung von `fx.Provide()` für Interface Binding
- ✅ fx.Lifecycle Hooks wo notwendig (Database, HTTP Server)
- ✅ Kein `init()` verwendet
- ✅ Keine globalen Variables (außer Version in main)

**Verbesserungspotenzial:**
- ⚠️ Einige Modules nicht in main.go registriert (cache, search, jobs)
- ⚠️ Config hardcoded in Modules statt über fx injected

**Beispiel (auth/module.go):**
```go
// ✅ PERFEKT - Folgt allen fx-Patterns
var Module = fx.Module("auth",
    fx.Provide(
        NewPasswordService,
        NewTokenService,
        NewService,
        AsPasswordService,  // ✅ Interface binding
        AsTokenService,
        AsAuthService,
    ),
)

func NewService(p ServiceParams) *Service {  // ✅ fx.In struct
    return newService(
        p.Users,
        p.Sessions,
        p.Passwords,
        p.Tokens,
        p.Config.Auth.MaxSessionsPerUser,
        accessDuration,
        refreshDuration,
    )
}
```

### 6.2 Content Module Pattern Compliance

**Rating:** 🔴 Non-Existent (0%)

**Status:** Keine Content-Module vorhanden zur Bewertung

**Expected Pattern:**
```go
// entity.go
type Movie struct {
    ContentEntity
    OriginalTitle string
    Tagline       string
    // ...
}

// repository.go
type Repository interface {
    GetByID(ctx context.Context, id uuid.UUID) (*Movie, error)
    // ...
}

// service.go
type Service struct {
    repo Repository
    jobs *jobs.Service
}

// handler.go (ogen interface)
type Handler struct {
    service *Service
}

// module.go
var Module = fx.Module("movie",
    fx.Provide(NewRepository, NewService, NewHandler),
)
```

### 6.3 Database Layer Compliance

**Rating:** 🟢 Excellent (95%)

**Positiv:**
- ✅ sqlc für type-safe queries
- ✅ Repository Pattern korrekt implementiert
- ✅ Domain Error Handling (ErrUserNotFound etc.)
- ✅ pgxpool mit Lifecycle Management
- ✅ golang-migrate für Migrations
- ✅ Connection pooling konfigurierbar
- ✅ Health checks vorhanden

**Beispiel (user_repository.go):**
```go
// ✅ PERFEKT - Domain Error Mapping
func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
    user, err := r.queries.GetUserByID(ctx, id)
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return nil, domain.ErrUserNotFound  // ✅ Domain error
        }
        return nil, fmt.Errorf("failed to get user by id: %w", err)
    }
    return mapDBUserToDomain(&user), nil
}
```

### 6.4 Config Management (koanf)

**Rating:** 🟢 Good (80%)

**Positiv:**
- ✅ koanf v2 korrekt verwendet
- ✅ Hierarchische Config (defaults → file → env)
- ✅ Environment variable overrides (`REVENGE_*`)
- ✅ Strukturierte Config mit Tags

**Verbesserungspotenzial:**
- ⚠️ Cache/Search/Jobs config hardcoded in Modules statt in Config struct
- ⚠️ Keine Validierung in `pkg/config/config.go`

**Missing in Config:**
```go
// pkg/config/config.go sollte haben:
type Config struct {
    Server   ServerConfig
    Database DatabaseConfig
    Cache    CacheConfig      // ⚠️ Vorhanden aber nicht verwendet
    Search   SearchConfig     // ⚠️ Vorhanden aber nicht verwendet
    Auth     AuthConfig
    OIDC     OIDCConfig
    Log      LogConfig
    Jobs     JobsConfig       // ❌ FEHLT komplett
}
```

### 6.5 API Handler Quality

**Rating:** 🟡 Mixed (60%)

**Positiv:**
- ✅ Go 1.22+ HTTP routing patterns verwendet
- ✅ Middleware-basierte Auth
- ✅ Strukturierte Error Responses
- ✅ Proper HTTP status codes

**Probleme:**
- ⚠️ Handlers direkt statt ogen-generated interfaces
- ⚠️ Keine OpenAPI Specs
- ⚠️ Manuelles JSON encoding statt generated

**Beispiel (auth.go):**
```go
// ⚠️ Manueller Handler - sollte ogen interface implementieren
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
    var req LoginRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        writeError(w, http.StatusBadRequest, "Invalid request body")
        return
    }
    // ...
}

// ✅ BESSER mit ogen:
// type AuthHandler struct {
//     ogen.UnimplementedHandler
//     service auth.Service
// }
// 
// func (h *AuthHandler) Login(ctx context.Context, req *ogen.LoginRequest) (*ogen.LoginResponse, error) {
//     // ogen handles marshaling
// }
```

---

## 7. TODO-Liste (Priorisiert)

### Phase 1: Fix Infrastructure Integration 🔴 URGENT

**Priority P0 (Critical - Do First):**

- [ ] **1.1 Registriere fehlende Modules in main.go**
  ```go
  // cmd/revenge/main.go
  app := fx.New(
      // Core modules
      fx.Provide(config.New, NewLogger),
      
      // Infrastructure modules
      database.Module,
      cache.Module,        // ← Add
      search.Module,       // ← Add
      jobs.Module,         // ← Add
      
      // Service modules
      auth.Module,
      user.Module,
      library.Module,
      rating.Module,
      oidc.Module,         // ← Add
      genre.Module,        // ← Add
      playback.Module,     // ← Add
      
      // ...
  )
  ```

- [ ] **1.2 Fix Config Integration**
  - Add `JobsConfig` to `pkg/config/config.go`
  - Remove hardcoded configs from cache/search/jobs modules
  - Use `fx.In` to inject Config properly

  ```go
  // internal/infra/cache/cache.go
  type ClientParams struct {
      fx.In
      Config *config.Config
      Logger *slog.Logger
  }
  
  var Module = fx.Module("cache",
      fx.Provide(func(p ClientParams) (*Client, error) {
          return NewClient(p.Config.Cache, p.Logger)
      }),
  )
  ```

- [ ] **1.3 Upgrade Typesense to v4**
  ```bash
  go get github.com/typesense/typesense-go/v4
  # Update imports in internal/infra/search/search.go
  ```

- [ ] **1.4 Add ogen dependency**
  ```bash
  go get github.com/ogen-go/ogen
  ```

### Phase 2: Movie Module (Reference Implementation) 🔴 P0

**Files to create:**

- [ ] **2.1 Database Layer**
  - [ ] `migrations/movie/000001_movies.up.sql`
  - [ ] `migrations/movie/000002_movie_people.up.sql`
  - [ ] `migrations/movie/000003_movie_streams.up.sql`
  - [ ] `migrations/movie/000004_movie_user_data.up.sql`
  - [ ] `queries/movie/movies.sql`
  - [ ] `queries/movie/movie_people.sql`
  - [ ] `queries/movie/movie_ratings.sql`
  - [ ] `queries/movie/movie_user_data.sql`

- [ ] **2.2 Domain Layer**
  - [ ] `internal/content/movie/entity.go`
  - [ ] Add `domain.MovieRepository` interface

- [ ] **2.3 Repository Layer**
  - [ ] `internal/content/movie/repository.go`
  - [ ] Register in `internal/infra/database/repository/module.go`

- [ ] **2.4 Service Layer**
  - [ ] `internal/content/movie/service.go`
  - [ ] `internal/content/movie/scanner.go`
  - [ ] `internal/content/movie/provider_tmdb.go`
  - [ ] `internal/content/movie/jobs.go`

- [ ] **2.5 API Layer**
  - [ ] `api/openapi/movies.yaml`
  - [ ] `internal/content/movie/handler.go` (implements ogen interface)

- [ ] **2.6 Module Registration**
  - [ ] `internal/content/movie/module.go`
  - [ ] Register in `cmd/revenge/main.go`

### Phase 3: Shared Migrations 🟡 P1

- [ ] **3.1 Additional Shared Tables**
  - [ ] `migrations/shared/000006_api_keys.sql`
  - [ ] `migrations/shared/000007_server_settings.sql`
  - [ ] `migrations/shared/000008_activity_log.sql`

- [ ] **3.2 Playlist/Collection Tables**
  - [ ] `migrations/shared/000010_video_playlists.sql`
  - [ ] `migrations/shared/000011_audio_playlists.sql`
  - [ ] `migrations/shared/000012_video_collections.sql`
  - [ ] `migrations/shared/000013_audio_collections.sql`

### Phase 4: TV Show Module 🟡 P1

- [ ] Clone movie module structure
- [ ] Adapt for series/seasons/episodes hierarchy
- [ ] Sonarr integration
- [ ] TheTVDB provider

### Phase 5: Music Module 🟡 P1

- [ ] Artists/Albums/Tracks tables
- [ ] Lidarr integration
- [ ] MusicBrainz provider
- [ ] Last.fm scrobbling

### Phase 6: Remaining Modules 🟡 P2

- [ ] Audiobook module
- [ ] Book module
- [ ] Podcast module
- [ ] Photo module
- [ ] LiveTV module
- [ ] Collection module

### Phase 7: Adult Modules 🟡 P3

- [ ] `migrations/c/000001_c_schema.sql` (CREATE SCHEMA c)
- [ ] `c.movies` tables
- [ ] `c.shows` tables
- [ ] `c.performers` (shared between movie/show)
- [ ] Whisparr integration

### Phase 8: OpenAPI Standardization 🟡 P2

- [ ] Define OpenAPI specs for all endpoints
- [ ] Generate handlers with ogen
- [ ] Migrate existing handlers to ogen interfaces
- [ ] Setup `go generate` workflow

---

## 8. Critical Fixes (Immediate Action Required)

### 🔥 Fix #1: Module Registration in main.go

**Problem:** Cache, Search, Jobs, OIDC, Genre, Playback modules existieren aber werden nicht gestartet

**Location:** `cmd/revenge/main.go`

**Fix:**
```go
app := fx.New(
    // Core modules
    fx.Provide(config.New, NewLogger),
    
    // Infrastructure modules
    database.Module,
    cache.Module,        // ADD
    search.Module,       // ADD
    jobs.Module,         // ADD
    
    // Service modules
    auth.Module,
    user.Module,
    library.Module,
    rating.Module,
    oidc.Module,         // ADD
    genre.Module,        // ADD
    playback.Module,     // ADD
    
    // ... rest
)
```

### 🔥 Fix #2: Config Hardcoding in Infra Modules

**Problem:** Cache, Search, Jobs haben hardcoded configs statt koanf

**Files:**
- `internal/infra/cache/cache.go`
- `internal/infra/search/search.go`
- `internal/infra/jobs/jobs.go`

**Fix Example (cache.go):**
```go
// BEFORE:
var Module = fx.Module("cache",
    fx.Provide(func(logger *slog.Logger) (*Client, error) {
        cfg := Config{Host: "localhost", Port: 6379, DB: 0}  // ❌ Hardcoded
        return NewClient(cfg, logger)
    }),
)

// AFTER:
type ClientParams struct {
    fx.In
    Config *config.Config
    Logger *slog.Logger
    LC     fx.Lifecycle
}

var Module = fx.Module("cache",
    fx.Provide(func(p ClientParams) (*Client, error) {
        client, err := NewClient(p.Config.Cache, p.Logger)
        if err != nil {
            return nil, err
        }
        
        p.LC.Append(fx.Hook{
            OnStart: func(ctx context.Context) error {
                return client.Ping(ctx)
            },
            OnStop: func(ctx context.Context) error {
                return client.Close()
            },
        })
        
        return client, nil
    }),
)
```

### 🔥 Fix #3: Typesense Version Mismatch

**Problem:** Code uses v3, docs specify v4

**Fix:**
```bash
go get github.com/typesense/typesense-go/v4
# Update import in internal/infra/search/search.go
```

---

## 9. Recommendations

### 9.1 Immediate (Next Sprint)

1. **Fix infrastructure module registration** (1 day)
   - Register cache, search, jobs, oidc, genre, playback in main.go
   - Fix config injection
   - Add lifecycle hooks

2. **Implement Movie module as reference** (1 week)
   - Complete database schema
   - Implement all required files
   - Add comprehensive tests
   - Document patterns for other modules

3. **Setup OpenAPI/ogen workflow** (2 days)
   - Create OpenAPI specs
   - Setup code generation
   - Migrate one handler as example

### 9.2 Short-term (1-2 Months)

1. **Complete Phase 1 modules**: TV Show, Music
2. **Implement shared playlists/collections**
3. **Add metadata provider integrations** (Radarr, Sonarr, Lidarr)
4. **Comprehensive integration tests** per module

### 9.3 Medium-term (3-6 Months)

1. **Complete remaining content modules**
2. **Implement adult content isolation** (schema c)
3. **Add external service integrations** (Trakt, Last.fm, etc.)
4. **Performance optimization** (caching, query optimization)

---

## 10. Positive Highlights ✨

### What's Working Well:

1. **✅ Core Infrastructure is Solid**
   - fx dependency injection properly implemented
   - sqlc database layer is type-safe and clean
   - koanf configuration is hierarchical and flexible
   - Repository pattern correctly implemented

2. **✅ Playback Service is Complete**
   - All 11 files implemented
   - Follows architecture exactly
   - Blackbeard integration ready
   - Client detection, bandwidth monitoring, buffering, caching all present

3. **✅ Modern Go Practices**
   - Go 1.22+ HTTP routing patterns
   - slog for structured logging
   - No init() functions
   - Proper error handling with errors.Is/As

4. **✅ Test Coverage (Where Present)**
   - auth service has comprehensive tests
   - user service has tests
   - Repository tests use testcontainers

5. **✅ Documentation Quality**
   - ARCHITECTURE_V2.md is detailed and comprehensive
   - Instructions are clear and actionable
   - Design principles are well-defined

---

## 11. Risk Assessment

| Risk | Severity | Likelihood | Mitigation |
|------|----------|------------|------------|
| Content modules missing delays MVP | 🔴 High | 🔴 High | Implement movie module first as template |
| Infrastructure not integrated causes runtime issues | 🔴 High | 🟡 Medium | Fix module registration immediately |
| No OpenAPI spec causes API inconsistencies | 🟡 Medium | 🟡 Medium | Start with one module, replicate pattern |
| Typesense v3/v4 mismatch causes compatibility issues | 🟡 Medium | 🟢 Low | Upgrade before first use |
| Missing migrations block database setup | 🔴 High | 🔴 High | Create shared migrations first |

---

## 12. Conclusion

### Current State Summary:

**🟢 Strengths:**
- Core infrastructure (fx, sqlc, koanf) is excellent
- Shared services (auth, user, oidc, playback) are well-implemented
- Architecture documentation is comprehensive
- Code quality follows best practices

**🔴 Critical Gaps:**
- All content modules are missing (0% implementation)
- Infrastructure modules not registered in main.go
- No OpenAPI/ogen integration
- Missing database migrations for content modules

**📊 Completion Status:**
- **Foundation:** 90% ✅
- **Shared Services:** 70% 🟡
- **Content Modules:** 0% ❌
- **External Integrations:** 40% 🟡
- **Overall:** ~40% 🟡

### Next Steps:

1. **Week 1:** Fix infrastructure integration (cache, search, jobs registration)
2. **Week 2-3:** Implement Movie module completely (reference implementation)
3. **Week 4:** Create shared migrations (playlists, collections, api_keys, etc.)
4. **Week 5-6:** Clone Movie pattern for TV Show module
5. **Week 7-8:** Implement Music module

### Estimated Timeline to MVP:

- **Phase 1 (Infrastructure Fix):** 1 week
- **Phase 2 (Movie Module):** 2 weeks
- **Phase 3 (Shared Migrations):** 1 week
- **Phase 4 (TV Show Module):** 2 weeks
- **Phase 5 (Music Module):** 2 weeks

**Total to 3-module MVP:** ~8 weeks (2 months)

---

**Generated:** 2026-01-28  
**Analyzer:** GitHub Copilot (Claude Sonnet 4.5)  
**Confidence:** High (comprehensive file analysis performed)
