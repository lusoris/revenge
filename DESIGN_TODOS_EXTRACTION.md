# Design Documentation TODO Extraction

**Generated**: 2026-01-29
**Source**: Complete analysis of all design documents

---

## Methodology

Systematically reviewed all design documents to extract missing implementation requirements:
- `docs/architecture/*.md`
- `docs/features/*.md`
- `docs/integrations/**/*.md`
- `docs/technical/*.md`

Compared against current codebase to identify gaps.

---

## 1. Frontend (Complete Implementation Missing)

### From ARCHITECTURE_V2.md + FRONTEND.md

**Status**: ❌ **0% - No frontend exists**

**Required Structure**:
```
web/
  src/
    lib/
      components/
        ui/              # shadcn-svelte components
        media/           # Media cards, players
        admin/           # Admin components
      stores/            # Svelte stores (auth, theme, playback)
      api/               # Generated API client
      utils/
    routes/
      (app)/
        (admin)/         # Admin panel
        (media)/         # Media browsing
        (player)/        # Video/audio player
      (auth)/            # Login, register, OIDC
```

**Technologies**:
- SvelteKit 2
- Tailwind CSS 4 + shadcn-svelte
- TanStack Query
- JWT + OIDC
- Shaka Player (video)
- Web Audio API (audio)

**Features to Implement**:
- [ ] Gapless audio (30s prefetch)
- [ ] Crossfade (5s overlap, dual gain nodes)
- [ ] Synced lyrics (LRC format)
- [ ] Visualizations (Canvas frequency bars)
- [ ] Quality switching (WebSocket)
- [ ] WebVTT subtitles
- [ ] RBAC (admin, moderator, user, guest)
- [ ] Light/Dark theme system
- [ ] PWA with offline support

**Priority**: P1 (Week 4-8)

---

## 2. API Layer (ogen/OpenAPI)

### From ARCHITECTURE_V2.md

**Status**: ❌ **0% - No OpenAPI specs exist**

**Required Structure**:
```
api/
  openapi/
    revenge.yaml      # Main spec
    movies.yaml       # Movie endpoints
    shows.yaml        # TV endpoints
    music.yaml        # Music endpoints
    ...
  generated/          # ogen-generated handlers
```

**Implementation Steps**:
- [ ] Create OpenAPI 3.1 specs per module
- [ ] Configure ogen code generation
- [ ] Implement generated handler interfaces
- [ ] Add ogen to go.mod
- [ ] Wire generated routers to main.go

**Priority**: P0 (Week 1)

---

## 3. Content Rating System

### From CONTENT_RATING.md

**Status**: ❌ **Migration + Service missing**

**Required Migration**: `000009_content_ratings.up.sql`

```sql
CREATE TABLE content_ratings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    system VARCHAR(20) NOT NULL,  -- 'mpaa', 'fsk', 'pegi', etc.
    code VARCHAR(20) NOT NULL,     -- 'PG-13', 'FSK16', etc.
    min_age INT NOT NULL,
    description TEXT,
    icon_path TEXT,
    UNIQUE(system, code)
);

CREATE TABLE movie_ratings (
    movie_id UUID PRIMARY KEY REFERENCES movies(id) ON DELETE CASCADE,
    content_rating_id UUID REFERENCES content_ratings(id)
);

CREATE TABLE show_ratings (
    show_id UUID PRIMARY KEY REFERENCES shows(id) ON DELETE CASCADE,
    content_rating_id UUID REFERENCES content_ratings(id)
);
```

**Required Service**: `internal/service/rating/content_rating.go`

```go
type ContentRatingService struct {
    queries *db.Queries
    cache   *cache.Client
}

func (s *ContentRatingService) Get(ctx context.Context, system, code string) (*ContentRating, error)
func (s *ContentRatingService) List(ctx context.Context, system string) ([]*ContentRating, error)
func (s *ContentRatingService) SetMovieRating(ctx context.Context, movieID uuid.UUID, ratingID uuid.UUID) error
```

**Seed Data**: MPAA, FSK, PEGI, BBFC ratings

**Priority**: P1 (Week 3)

---

## 4. Internationalization (i18n)

### From I18N.md

**Status**: ❌ **Not implemented**

**Required Layers**:

1. **UI Translation** (SvelteKit built-in)
   ```
   src/lib/i18n/
     en.json
     de.json
     fr.json
   ```

2. **Metadata Translation** (PostgreSQL JSONB)
   ```sql
   ALTER TABLE movies ADD COLUMN translations JSONB;
   -- { "de": { "title": "...", "overview": "..." }, ... }
   ```

3. **Audio/Subtitle Track Language**
   ```sql
   ALTER TABLE video_streams ADD COLUMN language VARCHAR(10);
   ALTER TABLE audio_streams ADD COLUMN language VARCHAR(10);
   ALTER TABLE subtitle_streams ADD COLUMN language VARCHAR(10);
   ```

**Priority**: P2 (Week 4)

---

## 5. Analytics Service

### From ANALYTICS_SERVICE.md

**Status**: ❌ **Not implemented**

**Required Tables**:
```sql
-- 000014_analytics.up.sql
CREATE TABLE analytics_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id),
    event_type VARCHAR(50) NOT NULL,  -- 'play', 'pause', 'seek', 'stop'
    module VARCHAR(50) NOT NULL,       -- 'movie', 'tvshow', 'music'
    resource_id UUID NOT NULL,
    metadata JSONB,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_analytics_user_date ON analytics_events(user_id, created_at DESC);
CREATE INDEX idx_analytics_module ON analytics_events(module, created_at DESC);

CREATE TABLE year_in_review (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    year INT NOT NULL,
    data JSONB NOT NULL,  -- Precomputed stats
    generated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(user_id, year)
);
```

**Required Service**: `internal/service/analytics/service.go`

```go
type Service struct {
    queries *db.Queries
    logger  *slog.Logger
}

func (s *Service) TrackEvent(ctx context.Context, event AnalyticsEvent) error
func (s *Service) GetUserStats(ctx context.Context, userID uuid.UUID, period Period) (*Stats, error)
func (s *Service) GenerateYearInReview(ctx context.Context, userID uuid.UUID, year int) error
```

**Features**:
- Watch time tracking
- Most watched content
- Genre preferences
- Time-of-day patterns
- Year in Review generation

**Priority**: P2 (Week 5)

---

## 6. Request System

### From REQUEST_SYSTEM.md

**Status**: ❌ **Not implemented**

**Required Migration**: `000015_requests.up.sql`

```sql
CREATE TABLE content_requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    module VARCHAR(50) NOT NULL,  -- 'movie', 'tvshow', 'music'
    title VARCHAR(500) NOT NULL,
    tmdb_id INT,
    tvdb_id INT,
    description TEXT,
    status VARCHAR(20) DEFAULT 'pending',  -- pending, approved, rejected, completed
    votes INT DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE request_votes (
    request_id UUID REFERENCES content_requests(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (request_id, user_id)
);
```

**Integration**: Radarr/Sonarr/Lidarr auto-add on approval

**Priority**: P3 (Week 6)

---

## 7. Ticketing System

### From TICKETING_SYSTEM.md

**Status**: ❌ **Not implemented**

**Required Migration**: `000016_tickets.up.sql`

```sql
CREATE TABLE tickets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    module VARCHAR(50),
    resource_id UUID,
    category VARCHAR(50) NOT NULL,  -- 'bug', 'feature', 'metadata', 'playback'
    priority VARCHAR(20) DEFAULT 'normal',
    status VARCHAR(20) DEFAULT 'open',
    title VARCHAR(200) NOT NULL,
    description TEXT NOT NULL,
    attachments JSONB,
    assigned_to UUID REFERENCES users(id),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE ticket_comments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ticket_id UUID NOT NULL REFERENCES tickets(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id),
    comment TEXT NOT NULL,
    is_staff_reply BOOLEAN DEFAULT false,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
```

**Priority**: P3 (Week 7)

---

## 8. Comics Module

### From COMICS_MODULE.md

**Status**: ❌ **Not implemented**

**Module Structure**:
```
internal/content/comics/
  domain.go
  service.go
  repository.go
  handlers.go
```

**Migration**: `internal/infra/database/migrations/comics/000001_comics.up.sql`

```sql
CREATE TABLE comics (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    library_id UUID NOT NULL REFERENCES libraries(id),
    title VARCHAR(500) NOT NULL,
    series VARCHAR(500),
    issue_number INT,
    volume INT,
    publisher VARCHAR(200),
    writers JSONB,      -- Array of names
    artists JSONB,      -- Array of names
    release_date DATE,
    page_count INT,
    format VARCHAR(50), -- 'cbz', 'cbr', 'pdf'
    file_path TEXT NOT NULL,
    cover_path TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE comic_pages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    comic_id UUID NOT NULL REFERENCES comics(id) ON DELETE CASCADE,
    page_number INT NOT NULL,
    image_path TEXT NOT NULL,
    UNIQUE(comic_id, page_number)
);
```

**Reader Features**:
- CBZ/CBR/PDF support
- Page-by-page navigation
- Double-page mode
- Continuous scroll mode
- Reading progress tracking

**Priority**: P3 (Week 8+)

---

## 9. External Integration Services

### From EXTERNAL_INTEGRATIONS_TODO.md

**Status**: ❌ **0/40+ integrations**

**High Priority Integrations** (P1):

#### Servarr Ecosystem
- [ ] **Radarr** - Movie management
  - `internal/service/metadata/radarr_client.go`
  - Webhook handler for imports
  - Metadata sync

- [ ] **Sonarr** - TV show management
  - `internal/service/metadata/sonarr_client.go`
  - Episode import events
  - Series metadata sync

- [ ] **Lidarr** - Music management
  - `internal/service/metadata/lidarr_client.go`
  - Artist/album imports
  - MusicBrainz IDs

#### Metadata Providers (P1)
- [ ] **TMDb** - Movie/TV metadata
  - `internal/service/metadata/tmdb_client.go`
  - Rate limiting (40 req/10s)

- [ ] **TheTVDB** - TV show metadata
  - `internal/service/metadata/tvdb_client.go`
  - JWT authentication

- [ ] **MusicBrainz** - Music metadata
  - `internal/service/metadata/musicbrainz_client.go`
  - Rate limiting (1 req/s)

#### Scrobbling (P2)
- [ ] **Trakt** - Movie/TV scrobbling
  - `internal/service/scrobble/trakt_client.go`
  - OAuth2 flow
  - Watch history sync

- [ ] **Last.fm** - Music scrobbling
  - `internal/service/scrobble/lastfm_client.go`
  - API key authentication

- [ ] **ListenBrainz** - Music scrobbling
  - `internal/service/scrobble/listenbrainz_client.go`
  - User token authentication

**Priority**: P1-P2 (Week 3-6)

---

## 10. Adult Content System

### From ADULT_CONTENT_SYSTEM.md, WHISPARR_STASHDB_SCHEMA.md

**Status**: ❌ **Complete system missing**

**Schema**: PostgreSQL schema `c` (obscured namespace)

**Integrations Required**:
- [ ] **Whisparr-v3** - Acquisition proxy
- [ ] **StashDB** - Performer/scene database
- [ ] **Stash-App** - Local organization
- [ ] **TPDB** - Fallback metadata
- [ ] **pHash** - Scene fingerprinting

**Migration**: `internal/infra/database/migrations/c/000001_c_schema.up.sql`

**Modules**:
- `internal/content/c/movie/` - Adult movies
- `internal/content/c/show/` - Adult series

**API Namespace**: `/c/*` (obscured)

**Priority**: P2 (Week 6+)

---

## 11. LiveTV & DVR

### From MEDIA_ENHANCEMENTS.md

**Status**: ❌ **Not implemented**

**Required Tables**:
```sql
CREATE TABLE livetv_channels (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    library_id UUID NOT NULL REFERENCES libraries(id),
    name VARCHAR(200) NOT NULL,
    number INT NOT NULL,
    stream_url TEXT NOT NULL,
    epg_id VARCHAR(100),
    logo_path TEXT,
    enabled BOOLEAN DEFAULT true
);

CREATE TABLE livetv_programs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    channel_id UUID NOT NULL REFERENCES livetv_channels(id) ON DELETE CASCADE,
    title VARCHAR(500) NOT NULL,
    description TEXT,
    start_time TIMESTAMPTZ NOT NULL,
    end_time TIMESTAMPTZ NOT NULL,
    season INT,
    episode INT
);

CREATE TABLE dvr_recordings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    program_id UUID REFERENCES livetv_programs(id),
    channel_id UUID NOT NULL REFERENCES livetv_channels(id),
    user_id UUID NOT NULL REFERENCES users(id),
    title VARCHAR(500) NOT NULL,
    file_path TEXT NOT NULL,
    start_time TIMESTAMPTZ NOT NULL,
    end_time TIMESTAMPTZ NOT NULL,
    status VARCHAR(20) DEFAULT 'scheduled'  -- scheduled, recording, completed, failed
);
```

**Integration**: XMLTV EPG support

**Priority**: P3 (Week 9+)

---

## 12. Media Enhancements

### From MEDIA_ENHANCEMENTS.md

**Features Missing**:

#### Trickplay (BIF Thumbnails)
- [ ] Generate thumbnails every 10s
- [ ] Store in WebP format
- [ ] BIF file format support
- [ ] Integration with video player

#### Intro Detection
- [ ] Chromaprint fingerprinting
- [ ] Season-wide intro matching
- [ ] "Skip Intro" button
- [ ] Auto-skip preference

#### Chapter Markers
- [ ] FFmpeg chapter extraction
- [ ] Manual chapter editing
- [ ] Chapter thumbnails
- [ ] Navigation UI

#### Theme Music
- [ ] Extract theme from episodes
- [ ] Background audio on details page
- [ ] Fade in/out

**Priority**: P2 (Week 5-6)

---

## 13. RBAC & Permissions

### From ARCHITECTURE_V2.md

**Status**: 🟡 **Partial - Basic auth exists, RBAC missing**

**Required Enhancement**: `000017_rbac.up.sql`

```sql
ALTER TABLE users ADD COLUMN role VARCHAR(20) DEFAULT 'user';
-- Roles: 'admin', 'moderator', 'user', 'guest'

CREATE TABLE permissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    role VARCHAR(20) NOT NULL,
    resource VARCHAR(100) NOT NULL,  -- 'libraries', 'users', 'settings', etc.
    action VARCHAR(50) NOT NULL,     -- 'read', 'write', 'delete', 'manage'
    UNIQUE(role, resource, action)
);

-- Seed default permissions
INSERT INTO permissions (role, resource, action) VALUES
    ('admin', '*', '*'),
    ('moderator', 'libraries', 'manage'),
    ('moderator', 'metadata', 'manage'),
    ('user', 'media', 'read'),
    ('user', 'playlists', 'manage'),
    ('guest', 'media', 'read');
```

**Service Enhancement**: `internal/service/auth/rbac.go`

```go
func (s *AuthService) CheckPermission(ctx context.Context, userID uuid.UUID, resource, action string) error
func (s *AuthService) GetUserPermissions(ctx context.Context, userID uuid.UUID) ([]Permission, error)
```

**Priority**: P1 (Week 2)

---

## 14. Profiles (Netflix-Style)

### From ARCHITECTURE_V2.md

**Status**: ❌ **Not implemented**

**Required Migration**: `000018_profiles.up.sql`

```sql
CREATE TABLE profiles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    avatar_path TEXT,
    is_kids_profile BOOLEAN DEFAULT false,
    pin VARCHAR(10),  -- Optional profile lock
    settings JSONB,   -- Theme, language, parental restrictions
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_profiles_user_name ON profiles(user_id, name);

-- Profile-specific watch history
ALTER TABLE watch_history ADD COLUMN profile_id UUID REFERENCES profiles(id);
ALTER TABLE favorites ADD COLUMN profile_id UUID REFERENCES profiles(id);
```

**Service**: `internal/service/user/profiles.go`

**Priority**: P2 (Week 4)

---

## Summary - Missing Components by Priority

### P0 - Critical (Week 1-2)
- [ ] Background Workers (7 workers)
- [ ] Shared Migrations (8 missing)
- [ ] Global Services (4 services)
- [ ] Session Service
- [ ] OpenAPI specs + ogen integration
- [ ] Module registration in main.go

### P1 - High (Week 2-4)
- [ ] Content Modules (11 modules)
- [ ] Content Rating System
- [ ] RBAC enhancements
- [ ] Servarr integrations (Radarr, Sonarr, Lidarr)
- [ ] TMDb/TheTVDB/MusicBrainz clients
- [ ] Frontend (SvelteKit)

### P2 - Medium (Week 4-8)
- [ ] i18n system
- [ ] Analytics Service
- [ ] Scrobbling (Trakt, Last.fm, ListenBrainz)
- [ ] Media Enhancements (Trickplay, Intro Detection, Chapters)
- [ ] Profiles system
- [ ] Adult Content System

### P3 - Low (Week 8+)
- [ ] Request System
- [ ] Ticketing System
- [ ] Comics Module
- [ ] LiveTV & DVR
- [ ] Additional integrations (40+ remaining)

---

## Implementation Dependencies

```
P0 (Core) ──┬──> P1 (Modules) ──┬──> P2 (Features)
            │                   │
            │                   └──> P2 (Scrobbling)
            │
            └──> P1 (Frontend) ──> P2 (Media Enhancements)
```

**Critical Path**:
1. P0: Core infrastructure (workers, services, migrations)
2. P1: Content modules + OpenAPI
3. P1: Frontend + Servarr integrations
4. P2: Features (Analytics, i18n, Profiles)
5. P3: Advanced features (Requests, Tickets, Comics)

---

**End of Extraction**

**Total Identified Gaps**: 100+ components/features
**Current Implementation**: ~10-15% of designed system
**Estimated Time to MVP**: 12-16 weeks (3-4 developers)
