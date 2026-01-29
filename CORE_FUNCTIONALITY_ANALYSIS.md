# Core Functionality Gap Analysis

**Generated**: 2026-01-29
**Analyzer**: GitHub Copilot (Claude Sonnet 4.5)

---

## Executive Summary

Die Codebase hat **kritische Lücken in der Core-Funktionalität**. Während Infrastructure-Layer (DB, Cache, Search, Jobs) existiert, fehlen:
- **Background Workers** (0/5 implemented)
- **Shared Migrations** (5/13 created, 8 missing)
- **Global Services** (0/4 implemented)
- **Session Management** (nur Repository, kein Service)

**Completion Score**: 🔴 **35% / 100%**

---

## 1. Background Workers (River)

### ✅ Infrastructure: Vollständig

**Location**: `internal/infra/jobs/jobs.go`

```go
type Service struct {
    client  *river.Client[pgx.Tx]
    workers *river.Workers
    logger  *slog.Logger
}

func NewWorkers() *river.Workers {
    return river.NewWorkers()
}
```

**Status**: ✅ River Client und Workers Registry vorhanden

### ❌ Workers: KEINE IMPLEMENTIERT

**Expected** (aus ARCHITECTURE_V2.md):

| Worker Type | Purpose | Status |
|-------------|---------|--------|
| `ScanLibraryWorker` | Scan folders for new media | ❌ Missing |
| `FetchMetadataWorker` | TMDb, TheTVDB, MusicBrainz, etc. | ❌ Missing |
| `DownloadImageWorker` | Download posters, generate blurhash | ❌ Missing |
| `IndexSearchWorker` | Update Typesense on changes | ❌ Missing |
| `CleanupWorker` | Remove orphaned files | ❌ Missing |
| `RefreshMetadataWorker` | Re-fetch metadata periodically | ❌ Missing |
| `NotificationWorker` | Webhook calls, email alerts | ❌ Missing |

**Search Results**:
- `AddWorker` in `internal/**/*.go` → **0 matches**
- `river.JobArgs` in `internal/**/*.go` → **2 matches** (nur Interface-Definitionen in jobs.go)

**Current Code**:
```go
// internal/infra/jobs/jobs.go:73-75
func NewWorkers() *river.Workers {
    return river.NewWorkers()  // ❌ Empty registry!
}
```

**Should Have**:
```go
// internal/infra/jobs/workers.go
package jobs

import (
    "context"
    "github.com/riverqueue/river"
    "github.com/lusoris/revenge/internal/service/library"
)

// RegisterWorkers registers all background workers.
func RegisterWorkers(workers *river.Workers, services WorkerServices) {
    // Library scanning
    river.AddWorker(workers, &ScanLibraryWorker{
        scanner: services.LibraryScanner,
    })

    // Metadata fetching
    river.AddWorker(workers, &FetchMetadataWorker{
        providers: services.MetadataProviders,
    })

    // Image processing
    river.AddWorker(workers, &DownloadImageWorker{
        downloader: services.ImageDownloader,
    })

    // Search indexing
    river.AddWorker(workers, &IndexSearchWorker{
        indexer: services.SearchIndexer,
    })

    // Cleanup
    river.AddWorker(workers, &CleanupWorker{
        cleanup: services.CleanupService,
    })
}

type WorkerServices struct {
    LibraryScanner    *library.Scanner
    MetadataProviders *metadata.ProviderRegistry
    ImageDownloader   *images.Downloader
    SearchIndexer     *search.Indexer
    CleanupService    *cleanup.Service
}

// ScanLibraryArgs defines the library scan job.
type ScanLibraryArgs struct {
    LibraryID uuid.UUID `json:"library_id"`
    FullScan  bool      `json:"full_scan"`
}

func (ScanLibraryArgs) Kind() string { return "library.scan" }

// ScanLibraryWorker scans libraries for new media.
type ScanLibraryWorker struct {
    river.WorkerDefaults[ScanLibraryArgs]
    scanner *library.Scanner
}

func (w *ScanLibraryWorker) Work(ctx context.Context, job *river.Job[ScanLibraryArgs]) error {
    return w.scanner.Scan(ctx, job.Args.LibraryID, job.Args.FullScan)
}
```

**Missing Services for Workers**:
- ❌ `internal/service/library/scanner.go` - Library scanner
- ❌ `internal/service/metadata/providers.go` - Metadata provider registry
- ❌ `internal/service/images/downloader.go` - Image downloader
- ❌ `internal/service/search/indexer.go` - Search indexer
- ❌ `internal/service/cleanup/service.go` - Cleanup service

---

## 2. Shared Database Migrations

### ✅ Implemented: 5 Migrations

**Location**: `internal/infra/database/migrations/shared/`

1. ✅ `000001_extensions.up.sql` - PostgreSQL extensions (uuid-ossp, pgcrypto)
2. ✅ `000002_users.up.sql` - Users table
3. ✅ `000003_sessions.up.sql` - Sessions table
4. ✅ `000004_oidc.up.sql` - OIDC providers
5. ✅ `000005_libraries.up.sql` - Libraries table

### ❌ Missing: 8 Migrations

**Expected** (aus MODULE_IMPLEMENTATION_TODO.md):

6. ❌ `000006_api_keys.sql` - External API authentication
   ```sql
   CREATE TABLE api_keys (
       id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
       name VARCHAR(100) NOT NULL,
       service VARCHAR(50) NOT NULL, -- 'tmdb', 'trakt', etc.
       key_value TEXT NOT NULL, -- Encrypted
       user_id UUID REFERENCES users(id),
       created_at TIMESTAMPTZ DEFAULT NOW(),
       last_used_at TIMESTAMPTZ
   );
   ```

7. ❌ `000007_server_settings.sql` - Persisted configuration
   ```sql
   CREATE TABLE server_settings (
       key VARCHAR(100) PRIMARY KEY,
       value JSONB NOT NULL,
       updated_at TIMESTAMPTZ DEFAULT NOW(),
       updated_by UUID REFERENCES users(id)
   );
   ```

8. ❌ `000008_activity_log.sql` - Audit log
   ```sql
   CREATE TABLE activity_log (
       id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
       user_id UUID REFERENCES users(id),
       module VARCHAR(50), -- 'movie', 'tvshow', 'auth', etc.
       action VARCHAR(50) NOT NULL, -- 'create', 'update', 'delete', 'login'
       resource_type VARCHAR(50),
       resource_id UUID,
       details JSONB,
       ip_address INET,
       user_agent TEXT,
       created_at TIMESTAMPTZ DEFAULT NOW()
   );
   CREATE INDEX idx_activity_log_user ON activity_log(user_id, created_at DESC);
   CREATE INDEX idx_activity_log_module ON activity_log(module, created_at DESC);
   ```

9. ❌ `000010_video_playlists.sql` - Shared movie + tvshow playlists
   ```sql
   CREATE TABLE video_playlists (
       id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
       user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
       name VARCHAR(200) NOT NULL,
       description TEXT,
       is_public BOOLEAN DEFAULT false,
       created_at TIMESTAMPTZ DEFAULT NOW(),
       updated_at TIMESTAMPTZ DEFAULT NOW()
   );

   CREATE TABLE video_playlist_items (
       id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
       playlist_id UUID NOT NULL REFERENCES video_playlists(id) ON DELETE CASCADE,
       item_type VARCHAR(20) NOT NULL, -- 'movie', 'episode'
       item_id UUID NOT NULL,
       position INT NOT NULL,
       added_at TIMESTAMPTZ DEFAULT NOW()
   );
   CREATE INDEX idx_video_playlist_items ON video_playlist_items(playlist_id, position);
   ```

10. ❌ `000011_audio_playlists.sql` - Shared music + audiobook + podcast
    ```sql
    CREATE TABLE audio_playlists (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
        user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
        name VARCHAR(200) NOT NULL,
        description TEXT,
        is_public BOOLEAN DEFAULT false,
        created_at TIMESTAMPTZ DEFAULT NOW(),
        updated_at TIMESTAMPTZ DEFAULT NOW()
    );

    CREATE TABLE audio_playlist_items (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
        playlist_id UUID NOT NULL REFERENCES audio_playlists(id) ON DELETE CASCADE,
        item_type VARCHAR(20) NOT NULL, -- 'track', 'audiobook_chapter', 'podcast_episode'
        item_id UUID NOT NULL,
        position INT NOT NULL,
        added_at TIMESTAMPTZ DEFAULT NOW()
    );
    ```

11. ❌ `000012_video_collections.sql` - Shared movie + tvshow collections
    ```sql
    CREATE TABLE video_collections (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
        library_id UUID NOT NULL REFERENCES libraries(id) ON DELETE CASCADE,
        name VARCHAR(200) NOT NULL,
        description TEXT,
        poster_path TEXT,
        backdrop_path TEXT,
        created_at TIMESTAMPTZ DEFAULT NOW()
    );

    CREATE TABLE video_collection_movies (
        collection_id UUID NOT NULL REFERENCES video_collections(id) ON DELETE CASCADE,
        movie_id UUID NOT NULL,
        position INT,
        PRIMARY KEY (collection_id, movie_id)
    );

    CREATE TABLE video_collection_episodes (
        collection_id UUID NOT NULL REFERENCES video_collections(id) ON DELETE CASCADE,
        episode_id UUID NOT NULL,
        position INT,
        PRIMARY KEY (collection_id, episode_id)
    );
    ```

12. ❌ `000013_audio_collections.sql` - Shared music + audiobook
    ```sql
    CREATE TABLE audio_collections (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
        library_id UUID NOT NULL REFERENCES libraries(id) ON DELETE CASCADE,
        name VARCHAR(200) NOT NULL,
        description TEXT,
        cover_path TEXT,
        created_at TIMESTAMPTZ DEFAULT NOW()
    );

    CREATE TABLE audio_collection_tracks (
        collection_id UUID NOT NULL REFERENCES audio_collections(id) ON DELETE CASCADE,
        track_id UUID NOT NULL,
        position INT,
        PRIMARY KEY (collection_id, track_id)
    );

    CREATE TABLE audio_collection_audiobooks (
        collection_id UUID NOT NULL REFERENCES audio_collections(id) ON DELETE CASCADE,
        audiobook_id UUID NOT NULL,
        position INT,
        PRIMARY KEY (collection_id, audiobook_id)
    );
    ```

13. ❌ River migrations - Should be auto-generated by River CLI

---

## 3. Global Services

### ❌ All Missing: 0/4 Implemented

**Expected Services**:

#### 1. Activity Logger Service

**Purpose**: Audit all user actions with module context

**Location**: `internal/service/activity/`

```go
package activity

type Service struct {
    queries *db.Queries
    logger  *slog.Logger
}

type LogEntry struct {
    UserID       uuid.UUID
    Module       string // 'movie', 'tvshow', 'auth', etc.
    Action       string // 'create', 'update', 'delete', 'login'
    ResourceType string
    ResourceID   uuid.UUID
    Details      map[string]any
    IPAddress    netip.Addr
    UserAgent    string
}

func (s *Service) Log(ctx context.Context, entry LogEntry) error {
    // Insert into activity_log table
}

func (s *Service) GetUserActivity(ctx context.Context, userID uuid.UUID, limit int) ([]LogEntry, error) {
    // Query activity_log
}

func (s *Service) GetModuleActivity(ctx context.Context, module string, limit int) ([]LogEntry, error) {
    // Query activity_log filtered by module
}
```

**Usage**:
```go
// In every handler that modifies data:
activityService.Log(ctx, activity.LogEntry{
    UserID:       currentUser.ID,
    Module:       "movie",
    Action:       "update",
    ResourceType: "movie",
    ResourceID:   movieID,
    Details:      map[string]any{"title": newTitle},
    IPAddress:    clientIP,
    UserAgent:    r.UserAgent(),
})
```

**Status**: ❌ Not implemented

---

#### 2. Server Settings Service

**Purpose**: Persisted configuration (overrides config.yaml)

**Location**: `internal/service/settings/`

```go
package settings

type Service struct {
    queries *db.Queries
    cache   *cache.Client
    logger  *slog.Logger
}

func (s *Service) Get(ctx context.Context, key string) (any, error) {
    // Check cache first, then DB
}

func (s *Service) Set(ctx context.Context, key string, value any, userID uuid.UUID) error {
    // Update DB, invalidate cache
}

func (s *Service) GetAll(ctx context.Context) (map[string]any, error) {
    // Return all settings
}
```

**Settings Examples**:
- `transcoding.max_concurrent`
- `scanning.interval_hours`
- `metadata.auto_fetch`
- `maintenance.mode`

**Status**: ❌ Not implemented

---

#### 3. API Key Service

**Purpose**: Manage external API keys (TMDb, Trakt, etc.)

**Location**: `internal/service/apikeys/`

```go
package apikeys

type Service struct {
    queries *db.Queries
    crypto  *crypto.Service // For encryption
    logger  *slog.Logger
}

type APIKey struct {
    ID        uuid.UUID
    Name      string
    Service   string // 'tmdb', 'trakt', 'musicbrainz', etc.
    KeyValue  string // Encrypted in DB
    UserID    uuid.UUID
    CreatedAt time.Time
    LastUsed  time.Time
}

func (s *Service) Get(ctx context.Context, service string) (*APIKey, error) {
    // Decrypt and return
}

func (s *Service) Set(ctx context.Context, service, keyValue string, userID uuid.UUID) error {
    // Encrypt and store
}

func (s *Service) Delete(ctx context.Context, service string) error {
    // Remove key
}
```

**Usage**:
```go
// In TMDb provider:
tmdbKey, err := apiKeyService.Get(ctx, "tmdb")
if err != nil {
    return ErrNoAPIKey
}
client := tmdb.NewClient(tmdbKey.KeyValue)
```

**Status**: ❌ Not implemented

---

#### 4. Notification Service

**Purpose**: Send notifications (webhooks, email, push)

**Location**: `internal/service/notifications/`

```go
package notifications

type Service struct {
    jobs   *jobs.Service
    logger *slog.Logger
}

type Notification struct {
    Type    string // 'webhook', 'email', 'push'
    UserID  uuid.UUID
    Subject string
    Message string
    Data    map[string]any
}

func (s *Service) Send(ctx context.Context, notif Notification) error {
    // Queue notification job
    return s.jobs.Insert(ctx, SendNotificationArgs{
        Type:    notif.Type,
        UserID:  notif.UserID,
        Subject: notif.Subject,
        Message: notif.Message,
        Data:    notif.Data,
    }, nil)
}
```

**Use Cases**:
- Metadata fetch complete
- Library scan complete
- New episode available
- Transcoding failed
- User login from new device

**Status**: ❌ Not implemented

---

## 4. Session Management

### 🟡 Partial: Repository Only

**Current State**:
- ✅ `internal/infra/database/repository/sessions.go` - Repository exists
- ✅ `internal/infra/database/queries/sessions.sql` - Queries exist
- ✅ `internal/infra/database/migrations/shared/000003_sessions.up.sql` - Table exists

**Missing**:
- ❌ `internal/service/session/service.go` - Service layer

**Expected Service**:
```go
// internal/service/session/service.go
package session

type Service struct {
    repo   repository.SessionRepository
    cache  *cache.Client
    logger *slog.Logger
}

func (s *Service) Create(ctx context.Context, userID uuid.UUID, deviceInfo DeviceInfo) (*domain.Session, error) {
    // Create session, store in cache + DB
}

func (s *Service) Validate(ctx context.Context, token string) (*domain.Session, error) {
    // Check cache first, then DB
}

func (s *Service) Invalidate(ctx context.Context, sessionID uuid.UUID) error {
    // Remove from cache, mark as inactive in DB
}

func (s *Service) GetUserSessions(ctx context.Context, userID uuid.UUID) ([]*domain.Session, error) {
    // List all active sessions for user
}

func (s *Service) InvalidateAllUserSessions(ctx context.Context, userID uuid.UUID) error {
    // Logout from all devices
}

func (s *Service) InvalidateOtherSessions(ctx context.Context, userID uuid.UUID, currentSessionID uuid.UUID) error {
    // Logout from all other devices
}

func (s *Service) UpdateActivity(ctx context.Context, sessionID uuid.UUID) error {
    // Update last_activity timestamp
}
```

**Integration**:
```go
// main.go
fx.Module("session",
    fx.Provide(session.NewService),
),
```

**Status**: 🟡 Repository exists, Service missing

---

## 5. Priority Matrix

| Component | Priority | Effort | Blockers | Status |
|-----------|----------|--------|----------|--------|
| **Background Workers** | P0 | 3 days | Content modules, Services | ❌ Missing |
| **Shared Migrations** | P0 | 1 day | None | 🟡 5/13 done |
| **Session Service** | P0 | 1 day | None | 🟡 Repo only |
| **Activity Logger** | P1 | 1 day | Migration #8 | ❌ Missing |
| **API Key Service** | P1 | 1 day | Migration #6 | ❌ Missing |
| **Server Settings Service** | P1 | 1 day | Migration #7 | ❌ Missing |
| **Notification Service** | P2 | 2 days | Workers | ❌ Missing |

---

## 6. Implementation Plan

### Week 1: Core Services & Migrations

**Day 1: Shared Migrations**
- [ ] Create `000006_api_keys.up.sql` + `.down.sql`
- [ ] Create `000007_server_settings.up.sql` + `.down.sql`
- [ ] Create `000008_activity_log.up.sql` + `.down.sql`
- [ ] Run migrations, test rollback

**Day 2: Playlist/Collection Migrations**
- [ ] Create `000010_video_playlists.up.sql` + `.down.sql`
- [ ] Create `000011_audio_playlists.up.sql` + `.down.sql`
- [ ] Create `000012_video_collections.up.sql` + `.down.sql`
- [ ] Create `000013_audio_collections.up.sql` + `.down.sql`

**Day 3: Session Service**
- [ ] Create `internal/service/session/service.go`
- [ ] Implement cache integration
- [ ] Register in main.go
- [ ] Update auth handlers to use session service

**Day 4: Activity Logger**
- [ ] Create `internal/service/activity/service.go`
- [ ] Add sqlc queries for activity_log
- [ ] Register in main.go
- [ ] Add logging to all handlers

**Day 5: API Key Service**
- [ ] Create `internal/service/apikeys/service.go`
- [ ] Implement encryption/decryption
- [ ] Add sqlc queries
- [ ] Register in main.go
- [ ] Create admin endpoints

### Week 2: Background Workers

**Day 1-2: Worker Infrastructure**
- [ ] Create `internal/infra/jobs/workers.go`
- [ ] Define all job args structs
- [ ] Create worker registration function
- [ ] Wire workers to job service

**Day 3: Library Scanner Service**
- [ ] Create `internal/service/library/scanner.go`
- [ ] Implement directory scanning
- [ ] File type detection
- [ ] Create `ScanLibraryWorker`

**Day 4: Metadata Services**
- [ ] Create `internal/service/metadata/providers.go`
- [ ] Provider registry pattern
- [ ] Create `FetchMetadataWorker`

**Day 5: Support Workers**
- [ ] Create `DownloadImageWorker`
- [ ] Create `IndexSearchWorker`
- [ ] Create `CleanupWorker`

---

## 7. Testing Strategy

### Unit Tests Required

- [ ] Session Service: Create, Validate, Invalidate
- [ ] Activity Logger: Log, Query by user, Query by module
- [ ] API Key Service: Encryption, Decryption, CRUD
- [ ] Each Worker: Work() method with mock dependencies

### Integration Tests Required

- [ ] Migrations: Up and Down for all new migrations
- [ ] Workers: End-to-end job execution
- [ ] Session Service: Cache + DB integration
- [ ] Activity Logger: DB queries with complex filters

---

## 8. Summary

**Current State**: 35% Complete
- ✅ Infrastructure (DB, Cache, Search, Jobs) exists
- 🟡 Migrations 5/13 (38%)
- ❌ Workers 0/7 (0%)
- 🟡 Session Management (repo only, no service)
- ❌ Global Services 0/4 (0%)

**Blockers**:
- Background workers need content module services (scanner, metadata providers)
- Global services need migrations

**Critical Path**:
1. Shared migrations (Day 1-2)
2. Session service (Day 3)
3. Global services (Day 4-5)
4. Worker infrastructure (Week 2)
5. Content modules (Week 3+)

**Estimated Time to Complete Core**: 2 weeks

---

**End of Report**

**Confidence**: High (comprehensive analysis performed)
**Methodology**:
- Architecture document review (ARCHITECTURE_V2.md)
- Implementation checklist review (MODULE_IMPLEMENTATION_TODO.md)
- Code structure analysis (file searches, greps)
- Migration directory inspection
- Worker registry inspection
