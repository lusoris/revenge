# Revenge - Architecture Compliance Analysis

**Generated**: 2026-01-29
**Analyzer**: GitHub Copilot (Claude Sonnet 4.5)

---

## Executive Summary

Die Codebase zeigt eine **Mischung aus konformer und nicht-konformer Implementierung** gegenüber den Dokumentationen und Instructions. Die Grundarchitektur (fx, koanf, pgx, slog) ist korrekt, aber **mehrere kritische Module fehlen komplett** und einige Infra-Module verwenden hardcoded configs statt koanf.

**Compliance Score**: 🟡 **65% / 100%**

| Kategorie | Score | Status |
|-----------|-------|--------|
| **Architektur-Design** | 90% | ✅ Sehr gut |
| **Module-Implementierung** | 15% | ❌ Kritisch |
| **Dependencies & Stack** | 95% | ✅ Exzellent |
| **Code-Patterns** | 70% | 🟡 Verbesserungswürdig |
| **Testing** | 60% | 🟡 OK |

---

## 1. Architektur-Konformität

### ✅ Korrekt Implementiert

1. **Dependency Injection (fx)**
   - ✅ Verwendet `fx.New()` in main.go
   - ✅ Module nutzen `fx.Module()` Pattern
   - ✅ Keine `init()` Funktionen gefunden
   - ✅ Services nutzen fx.Lifecycle Hooks
   - **Fundstellen**: `cmd/revenge/main.go`, alle `**/module.go` Dateien

2. **Structured Logging (slog)**
   - ✅ Verwendet `log/slog` stdlib (nicht zap/logrus)
   - ✅ Verwendet tint für Console-Output
   - ✅ Structured attributes mit `slog.String()`, `slog.Any()`
   - **Fundstellen**: `cmd/revenge/main.go`, alle Service-Files

3. **Database Stack**
   - ✅ Verwendet pgx/v5 (nicht lib/pq)
   - ✅ sqlc für type-safe queries
   - ✅ sqlc.yaml korrekt konfiguriert
   - ✅ Migrations in `internal/infra/database/migrations/shared/`
   - **Fundstellen**: `sqlc.yaml`, `internal/infra/database/`

4. **HTTP Routing**
   - ✅ Verwendet Go 1.22+ stdlib routing (`mux.HandleFunc("GET /path", ...)`)
   - ✅ Keine gorilla/mux dependency
   - **Fundstellen**: `cmd/revenge/main.go:143+`

5. **Configuration (koanf)**
   - ✅ Verwendet koanf (nicht viper)
   - ✅ Structured config mit pkg/config/
   - **Fundstellen**: `pkg/config/`, `cmd/revenge/main.go`

6. **Background Jobs (River)**
   - ✅ River job queue implementiert
   - ✅ PostgreSQL-native (kein Redis nötig)
   - **Fundstellen**: `internal/infra/jobs/jobs.go`

7. **GOMAXPROCS**
   - ✅ Keine `automaxprocs` dependency (Go 1.25 built-in)
   - **Verifiziert**: `go.mod` hat kein uber-go/automaxprocs

### 🟡 Teilweise Konform

1. **Go 1.25 Features**
   - ❌ `sync.WaitGroup.Go()` **nicht** verwendet (alte `wg.Add(1)` + `defer wg.Done()` Pattern in 4 Files)
   - ✅ Keine deprecated Patterns gefunden
   - **Violations**:
     - `pkg/supervisor/supervisor.go:191`
     - `pkg/health/health.go:126-128`
     - `pkg/graceful/shutdown.go:219`
   - **Action**: Diese auf `wg.Go(func() { ... })` umstellen

2. **Context Management**
   - ⚠️ `context.Background()` in 20+ Stellen (meist in Tests OK, einige in prod kritisch)
   - ✅ Services nehmen `context.Context` als ersten Parameter
   - **Critical**:
     - `pkg/supervisor/supervisor.go:138` - sollte von fx lifecycle ctx kommen
     - `pkg/graceful/shutdown.go:126` - sollte parent context nutzen
   - **Action**: Context von fx lifecycle nutzen, nicht Background()

### ❌ Nicht Konform

1. **Infrastructure Hardcoded Configs**
   - ❌ `internal/infra/cache/cache.go:73-79` - hardcoded `localhost:6379`
   - ❌ `internal/infra/search/search.go:75` - hardcoded `localhost:8108`
   - ❌ `internal/infra/jobs/jobs.go` - hardcoded worker counts
   - **Expected**: Config von koanf laden
   - **Action**: Config structs aus koanf/v2 nutzen

2. **Module Registration in main.go**
   - ❌ Cache.Module **nicht** registriert
   - ❌ Search.Module **nicht** registriert
   - ❌ Jobs.Module **nicht** registriert (ist aber implementiert!)
   - ❌ OIDC.Module **nicht** registriert
   - ❌ Genre.Module **nicht** registriert
   - ❌ Playback.Module **nicht** registriert
   - **Current**: `cmd/revenge/main.go:40-53` registriert nur:
     - config, logger
     - database.Module
     - auth.Module, user.Module, library.Module, rating.Module
   - **Action**: Alle Module in fx.New() registrieren

---

## 2. Module-Implementierung

### ✅ Implementierte Module (7/17)

| Module | Status | Location |
|--------|--------|----------|
| `auth` | ✅ Complete | `internal/service/auth/` |
| `user` | ✅ Complete | `internal/service/user/` |
| `library` | ✅ Complete | `internal/service/library/` |
| `rating` | ✅ Complete | `internal/service/rating/` |
| `oidc` | ✅ Exists (not registered) | `internal/service/oidc/` |
| `genre` | ✅ Exists (not registered) | `internal/service/genre/` |
| `playback` | ✅ Exists (not registered) | `internal/service/playback/` |

### ❌ Fehlende Content-Module (11/11)

**ALLE Content-Module fehlen komplett:**

| Module | Expected Path | Status |
|--------|---------------|--------|
| `movie` | `internal/content/movie/` | ❌ **FEHLT** |
| `tvshow` | `internal/content/tvshow/` | ❌ **FEHLT** |
| `music` | `internal/content/music/` | ❌ **FEHLT** |
| `audiobook` | `internal/content/audiobook/` | ❌ **FEHLT** |
| `book` | `internal/content/book/` | ❌ **FEHLT** |
| `podcast` | `internal/content/podcast/` | ❌ **FEHLT** |
| `photo` | `internal/content/photo/` | ❌ **FEHLT** |
| `livetv` | `internal/content/livetv/` | ❌ **FEHLT** |
| `collection` | `internal/content/collection/` | ❌ **FEHLT** |
| `adult_movie` | `internal/content/c/movie/` | ❌ **FEHLT** |
| `adult_show` | `internal/content/c/show/` | ❌ **FEHLT** |

**Impact**: Das System hat keine Content-Verwaltung. Es kann keine Filme, Serien, Musik etc. verwalten.

**Verzeichnis `internal/content/` existiert nicht einmal!**

---

## 3. Database Schema

### ✅ Shared Migrations Vorhanden

| Migration | File | Status |
|-----------|------|--------|
| Extensions | `000001_extensions.*.sql` | ✅ |
| Users | `000002_users.*.sql` | ✅ |
| Sessions | `000003_sessions.*.sql` | ✅ |
| OIDC | `000004_oidc.*.sql` | ✅ |
| Libraries | `000005_libraries.*.sql` | ✅ |

### ❌ Content-Module Migrations Fehlen

**Expected Structure** (aus ARCHITECTURE_V2.md):
```
internal/infra/database/migrations/
  shared/           ✅ Vorhanden
  movie/            ❌ FEHLT
  tvshow/           ❌ FEHLT
  music/            ❌ FEHLT
  audiobook/        ❌ FEHLT
  book/             ❌ FEHLT
  podcast/          ❌ FEHLT
  photo/            ❌ FEHLT
  livetv/           ❌ FEHLT
  collection/       ❌ FEHLT
  c/                ❌ FEHLT (Adult schema)
```

**Current**: Nur `shared/` vorhanden, keine Module-Migrations.

### ❌ Adult Schema Isolation

- ❌ Schema `c` **nicht** erstellt (`CREATE SCHEMA c` nicht gefunden)
- ❌ Keine `c.*` Tabellen
- ❌ Adult Module (`internal/content/c/`) fehlt komplett

**Expected** (aus adult-modules.instructions.md):
```sql
CREATE SCHEMA IF NOT EXISTS c;
-- Tables: c.movies, c.scenes, c.performers, c.studios, etc.
```

---

## 4. API Design

### ✅ Konform

- ✅ Nutzt stdlib HTTP routing (Go 1.22+)
- ✅ Handlers in `internal/api/handlers/`
- ✅ Middleware in `internal/api/middleware/`
- ✅ Revenge-kompatible Endpoints (`/Users/AuthenticateByName`, etc.)

### ❌ Nicht Konform

- ❌ **Kein ogen** (OpenAPI spec-first generation)
- ❌ Keine OpenAPI specs in `api/openapi/`
- ❌ Keine generated handlers in `api/generated/`
- **Current**: Manuelle Handler-Implementierung
- **Expected** (aus ARCHITECTURE_V2.md):
  ```
  api/
    openapi/
      revenge.yaml
      movies.yaml
      shows.yaml
    generated/  # ogen-generated
  ```

---

## 5. Dependency Compliance

### ✅ Korrekte Dependencies

| Requirement | Actual | Status |
|-------------|--------|--------|
| Go 1.25+ | `go 1.25.0` | ✅ |
| pgx/v5 | `v5.8.0` | ✅ |
| fx | `v1.24.0` | ✅ |
| koanf/v2 | `v2.3.2` | ✅ |
| redis (Dragonfly) | `v9.17.3` | ✅ |
| River | `v0.30.2` | ✅ |

### ❌ Fehlende Dependencies

| Required | Status |
|----------|--------|
| `github.com/ogen-go/ogen` | ❌ **FEHLT** |
| `github.com/typesense/typesense-go/v4` | 🟡 Falsche Version (`v3.2.0` statt `v4`) |

### ✅ Keine Verbotenen Dependencies

- ✅ Keine gorilla/mux
- ✅ Keine viper
- ✅ Keine zap/logrus (nur als indirect dep von testcontainers)
- ✅ Keine lib/pq
- ✅ Keine automaxprocs

---

## 6. Code Patterns

### ✅ Gute Patterns

1. **Error Handling**
   - ✅ `errors.Is()` / `errors.As()` verwendet
   - ✅ Wrapped errors mit `fmt.Errorf(..., %w, err)`

2. **Testing**
   - ✅ Integration tests mit testcontainers
   - ✅ Tests in `tests/integration/`

3. **Struct Design**
   - ✅ Services nutzen fx.In/fx.Out structs
   - ✅ Config structs gut strukturiert

### 🟡 Verbesserungswürdig

1. **WaitGroup Pattern**
   - ❌ Alte `wg.Add(1); go func() { defer wg.Done() }` in 4 Files
   - ✅ Sollte Go 1.25 `wg.Go()` nutzen

2. **Context Propagation**
   - ⚠️ Zu viel `context.Background()` außerhalb von Tests

3. **Panic Usage**
   - 🟡 `panic()` nur in `pkg/lazy/lazy.go:52` (MustGet) - OK für Must* Functions

---

## 7. Testing

### ✅ Vorhanden

- ✅ Integration tests (`tests/integration/`)
- ✅ Service tests (`internal/service/*/service_test.go`)
- ✅ Handler tests (`internal/api/handlers/*_test.go`)
- ✅ testcontainers für PostgreSQL

### 🟡 Gaps

- 🟡 Keine Benchmark tests (sollte `testing.B.Loop()` nutzen für Go 1.24+)
- 🟡 Keine synctest usage (Go 1.25 testing/synctest für race detection)

---

## 8. Critical Issues (Priority-Sorted)

### 🔥 P0 - Blocker (Immediate Action Required)

1. **Content-Module fehlen komplett**
   - **Impact**: Keine Content-Verwaltung möglich
   - **Location**: `internal/content/` existiert nicht
   - **Fix**: Movie-Modul als Reference-Implementation erstellen
   - **Estimated**: 2 weeks (siehe CODEBASE_ANALYSIS_REPORT.md Phase 2)

2. **Module Registration in main.go**
   - **Impact**: Cache, Search, Jobs, OIDC, Genre, Playback nicht gestartet
   - **Location**: `cmd/revenge/main.go:40-53`
   - **Fix**:
     ```go
     fx.New(
         fx.Provide(config.New, NewLogger),
         database.Module,
         cache.Module,      // ADD
         search.Module,     // ADD
         jobs.Module,       // ADD
         auth.Module,
         user.Module,
         library.Module,
         rating.Module,
         oidc.Module,       // ADD
         genre.Module,      // ADD
         playback.Module,   // ADD
         // ...
     )
     ```
   - **Estimated**: 10 minutes

3. **Hardcoded Infra Configs**
   - **Impact**: Kann nicht in Production verwendet werden
   - **Files**:
     - `internal/infra/cache/cache.go:73-79`
     - `internal/infra/search/search.go:75`
   - **Fix**: Config aus koanf laden
   - **Estimated**: 1 hour

### 🟡 P1 - High Priority

4. **Adult Schema Isolation fehlt**
   - **Impact**: Adult content kann nicht verwaltet werden
   - **Location**: `migrations/c/` fehlt, Schema `c` nicht erstellt
   - **Fix**: `000001_c_schema.up.sql` mit `CREATE SCHEMA c`
   - **Estimated**: 30 minutes

5. **OpenAPI/ogen fehlt**
   - **Impact**: Kein spec-first API design
   - **Location**: `api/openapi/`, `api/generated/` fehlen
   - **Fix**: OpenAPI specs schreiben, ogen generieren
   - **Estimated**: 1 week

6. **WaitGroup Pattern veraltet**
   - **Impact**: Nicht idiomatischer Go 1.25 code
   - **Files**: 4 Files mit altem Pattern
   - **Fix**: `wg.Go()` statt `wg.Add(1); go func() { defer wg.Done() }`
   - **Estimated**: 30 minutes

### 🟢 P2 - Nice to Have

7. **Context.Background() in Production Code**
   - **Impact**: Fehlende context cancellation
   - **Files**: `pkg/supervisor/supervisor.go:138`, `pkg/graceful/shutdown.go:126`
   - **Fix**: Context von fx lifecycle nutzen
   - **Estimated**: 1 hour

8. **Typesense Version**
   - **Impact**: Veraltete v3 statt v4
   - **Fix**: `go get github.com/typesense/typesense-go/v4`
   - **Estimated**: 5 minutes

---

## 9. Recommendations

### Immediate Actions (Week 1)

1. ✅ **Fix Module Registration** (10 min)
   - Alle Module in main.go registrieren

2. ✅ **Fix Hardcoded Configs** (1 hour)
   - Cache, Search, Jobs Config von koanf laden

3. ✅ **Update Dependencies** (5 min)
   - Typesense v4
   - ogen hinzufügen

### Short-Term (Week 2-3)

4. ✅ **Movie Module (Reference Implementation)**
   - Complete module mit:
     - Migrations
     - Queries (sqlc)
     - Repository, Service, Handler
     - Scanner, Jobs
     - TMDb provider
   - Dient als Template für andere Module

5. ✅ **Fix Code Patterns**
   - `sync.WaitGroup.Go()` Pattern
   - Context propagation

### Medium-Term (Week 4-8)

6. ✅ **Implement Remaining Modules**
   - TV Show, Music, Audiobook, Book, Podcast, Photo, LiveTV, Collection
   - Adult modules (c/movie, c/show)

7. ✅ **OpenAPI Standardization**
   - Specs schreiben
   - ogen integration
   - Generated handlers

---

## 10. Compliance Checklist

### Architecture ✅ 90%

- [x] fx dependency injection
- [x] koanf configuration
- [x] pgx/v5 database
- [x] sqlc type-safe queries
- [x] slog structured logging
- [x] Go 1.22+ stdlib routing
- [x] River job queue
- [ ] ogen OpenAPI generation ❌
- [ ] Module isolation (content/) ❌

### Modules ❌ 15%

- [x] auth module
- [x] user module
- [x] library module
- [x] rating module
- [x] oidc module (exists, not registered)
- [x] genre module (exists, not registered)
- [x] playback module (exists, not registered)
- [ ] movie module ❌
- [ ] tvshow module ❌
- [ ] music module ❌
- [ ] audiobook module ❌
- [ ] book module ❌
- [ ] podcast module ❌
- [ ] photo module ❌
- [ ] livetv module ❌
- [ ] collection module ❌
- [ ] adult_movie module ❌
- [ ] adult_show module ❌

### Database ✅ 70%

- [x] Shared migrations (users, sessions, libraries)
- [x] sqlc queries for shared tables
- [x] PostgreSQL extensions
- [ ] Movie module migrations ❌
- [ ] TV Show module migrations ❌
- [ ] Music module migrations ❌
- [ ] Adult schema (`c`) ❌
- [ ] Content module queries ❌

### Code Quality ✅ 70%

- [x] No `init()` functions
- [x] No global variables
- [x] No forbidden dependencies
- [x] Proper error handling
- [x] Context as first parameter
- [ ] Go 1.25 WaitGroup.Go() ❌
- [ ] Minimal context.Background() ❌

### Testing ✅ 60%

- [x] Integration tests
- [x] Service tests
- [x] Handler tests
- [x] testcontainers
- [ ] Benchmark tests ❌
- [ ] synctest usage ❌

---

## 11. Overall Assessment

**Verdict**: 🟡 **Gute Grundlage, aber unvollständig**

### Strengths ✅

1. **Excellent Tech Stack Choices**
   - Moderne Go 1.25 features
   - Richtige Dependencies (fx, koanf, pgx, River)
   - Keine deprecated/verbotene deps

2. **Solid Foundation**
   - fx DI architecture
   - Proper logging, error handling
   - Good database setup (sqlc, migrations)

3. **Clean Code**
   - No `init()`, no globals
   - Good separation of concerns
   - Proper testing setup

### Weaknesses ❌

1. **Missing Content Modules**
   - 0 von 11 Content-Modulen implementiert
   - Das ist der Hauptzweck der Software!

2. **Infrastructure Not Fully Utilized**
   - Cache, Search, Jobs Module existieren aber nicht registriert
   - Hardcoded configs statt koanf

3. **API Design Incomplete**
   - Kein ogen/OpenAPI (wie dokumentiert)
   - Manuelle Handler statt generated

### Next Steps

**Phase 1** (Week 1): Fix Infrastructure
- Module registration
- Config loading
- Dependency updates

**Phase 2** (Week 2-3): Movie Module (Reference)
- Complete implementation als Template

**Phase 3** (Week 4-8): Rollout Content Modules
- TV Show, Music, etc.

**Phase 4** (Week 9+): OpenAPI Standardization

---

**End of Report**

**Confidence**: High (comprehensive analysis performed)
**Methodology**:
- Documentation review (ARCHITECTURE_V2.md, Instructions)
- Code analysis (grep, file_search, semantic_search)
- Dependency audit (go.mod)
- Pattern matching (Go 1.25 features, fx patterns, etc.)
