# Native Request System

> Content request management for all modules - replaces Overseerr/Jellyseerr

**Status**: 🔴 DESIGN PHASE
**Priority**: 🟡 HIGH (Phase 9 - External Services)
**Replaces**: Overseerr, Jellyseerr (NO integration - native only)

---

## Design Decision

**NO Overseerr/Jellyseerr Integration**: Build native request system optimized for Revenge's modular architecture.

**Why Native**:
- Overseerr only supports movies/TV
- Jellyseerr only supports Jellyfin (movies/TV)
- Neither supports music, books, audiobooks, podcasts, comics, adult content
- Native integration with Revenge modules = better UX
- Direct integration with Audiobookshelf, Radarr, Sonarr, Lidarr, etc.
- Unified request workflow across ALL content types
- **Inline UI**: Integrated directly in Revenge UI (no external apps)
- **Advanced automation**: Intelligent request rules based on watch history
- **Modular architecture**: Per-content-type request modules (movie, tvshow, music, audiobook, book, podcast, comic, adult)
- **Deep integrations**: Ticketing system, rating system, analytics, storage quotas

---

## Core Features (from Overseerr/Jellyseerr)

### User-Facing Features
- **Request content**: Users can request movies, TV shows, music, books, audiobooks, podcasts, comics, adult content
- **Search**: Integrated search across TMDb, TheTVDB, MusicBrainz, Goodreads, ComicVine, StashDB
- **Availability checking**: Show if content already exists in library
- **Request tracking**: Status updates (Pending, Approved, Processing, Available, Declined)
- **Notifications**: Email/Discord/Telegram when requests approved/available
- **Voting**: Users can upvote requests (priority queue)
- **Comments**: Discussion on requests

### Admin Features
- **Approval workflow**: Auto-approve OR manual review
- **Quota management**: Per-user request limits (daily/weekly/monthly) + **disk space quotas** (per content type, per user, global)
- **Request rules**: Auto-approve based on user role/trust score/watch history/storage capacity
- **Batch approval**: Approve multiple requests at once
- **Integration triggers**: Automatically add to Radarr/Sonarr/Lidarr on approval
- **Priority management**: Admin can prioritize requests
- **Request history**: Full audit trail
- **Storage analytics**: Real-time disk usage per content type, projected storage needs
- **Automated cleanup**: Auto-decline requests if storage quota exceeded

### Advanced Automation Features
- **Intelligent season requests** (TV shows):
  - User watching S1 → auto-request S2 (configurable rule)
  - Nobody watched episode yet → only fetch S1, wait for engagement before S2
  - User completed 80% of S1 → pre-approve S2
- **Watch-based priority**:
  - Frequently requested content = higher priority
  - Abandoned content (no one watching) = lower priority
- **Storage-aware rules**:
  - Auto-decline 4K requests if disk space <100GB
  - Suggest lower quality if storage constrained
- **User behavior analysis**:
  - User never watches horror → auto-decline horror requests
  - User binge-watches sci-fi → auto-approve sci-fi requests
- **Content lifecycle management**:
  - Auto-delete unwatched content after 90 days (free up space for new requests)
  - Keep frequently re-watched content indefinitely

---

## Per-Module Request Handling

### Movies (Radarr Integration)
```
User searches TMDb → Selects movie → Submits request
                                          ↓
                               Admin approves (OR auto-approve)
                                          ↓
                      Revenge adds to Radarr via API
                                          ↓
                      Radarr downloads → Imports → Revenge notified
                                          ↓
                      Request status: Available → Notify user
```

### TV Shows (Sonarr Integration)
```
User searches TheTVDB → Selects show → Selects seasons → Submits request
                                                              ↓
                                               Admin approves
                                                              ↓
                                 Revenge adds to Sonarr (seasons configured)
                                                              ↓
                                 Sonarr downloads → Imports → Notify user
```

### Music (Lidarr Integration)
```
User searches MusicBrainz → Selects artist/album → Submits request
                                                          ↓
                                               Admin approves
                                                          ↓
                              Revenge adds to Lidarr
                                                          ↓
                              Lidarr downloads → Imports → Notify user
```

### Audiobooks (Audiobookshelf Integration)
```
User searches Audible/Goodreads → Selects audiobook → Submits request
                                                              ↓
                                               Admin approves
                                                              ↓
                                  Admin manually downloads (OR script integration)
                                                              ↓
                                  Add to Audiobookshelf library
                                                              ↓
                                  Request status: Available → Notify user
```

### Books (Readarr Integration)
```
User searches Goodreads → Selects book → Submits request
                                              ↓
                                   Admin approves
                                              ↓
                          Revenge adds to Readarr
                                              ↓
                          Readarr downloads → Imports → Notify user
```

### Podcasts (Audiobookshelf Integration)
```
User searches by RSS feed OR podcast name → Submits request
                                                    ↓
OPTION 1: Scene request
User searches StashDB → Selects scene → Submits request
                                            ↓
                                 Admin approves
                                            ↓
                        Revenge adds to Whisparr
                                            ↓
                        Whisparr downloads → Imports → Notify user

OPTION 2: Studio request (all content from studio)
User searches StashDB → Selects studio (e.g., "Studio XYZ") → Submits request
                                                                      ↓
                                                           Admin approves
                                                                      ↓
                              Revenge adds ALL studio scenes to Whisparr (monitored)
                                                                      ↓
                              Whisparr downloads new releases automatically → Notify user

OPTION 3: Performer request (all content with performer)
User searches StashDB → Selects performer (e.g., "Performer ABC") → Submits request
                                                                          ↓
                                                               Admin approves
                                                                          ↓
                          Revenge adds ALL performer scenes to Whisparr (monitored)
                                                                          ↓
                          Whisparr downloads → Imports → Notify user

OPTION 4: Tag/genre combination (e.g., "VR + POV")
User selects tags (VR, POV, etc.) → Submits request
                                            ↓
                                 Admin approves
                                            ↓
            Revenge searches StashDB for matching scenes → Adds to Whisparr
                                            ↓

### Comics (Mylar3 Integration - Future)
```
User searches ComicVine → Selects series/issue → Submits request
                                                      ↓
                                           Admin approves
                                                      ↓
                          Revenge adds to Mylar3 (future integration)
                                                      ↓
                          Mylar3 downloads → Imports → Notify user
```

---

## Database Schema

### Public Schema (Non-Adult Content)

```sql
-- Main requests table (public content only - NO adult content)
CREATE TABLE requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    content_type VARCHAR(50) NOT NULL CHECK (content_type IN ('movie', 'tvshow', 'tvshow_season', 'music_album', 'music_artist', 'audiobook', 'book', 'podcast', 'comic')),
    content_subtype VARCHAR(50),              -- movie, tvshow, tvshow_season, music_album, music_artist, audiobook, book, podcast, comic
    external_id VARCHAR(200),                 -- TMDb ID, TheTVDB ID, MusicBrainz ID, etc. (NOT StashDB)
    title VARCHAR(500) NOT NULL,
    release_year INT,
    metadata_json JSONB,                      -- Content-specific metadata (season selection, quality preference - NO adult tags)
    status VARCHAR(50) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'approved', 'processing', 'available', 'declined', 'on_hold')),
    auto_approved BOOLEAN DEFAULT FALSE,
    auto_rule_id UUID REFERENCES request_rules(id) ON DELETE SET NULL,  -- Which rule triggered auto-approval
    approved_by_user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    approved_at TIMESTAMPTZ,
    declined_reason TEXT,
    priority INT DEFAULT 0,                   -- Higher = more important
    votes_count INT DEFAULT 0,
    integration_id VARCHAR(200),              -- Radarr/Sonarr/Lidarr ID (after approval)
    integration_status VARCHAR(100),          -- Radarr/Sonarr status
    estimated_size_gb DECIMAL(10,2),          -- Estimated disk space required
    actual_size_gb DECIMAL(10,2),             -- Actual disk space used (after download)
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    available_at TIMESTAMPTZ,
    triggered_by_automation BOOLEAN DEFAULT FALSE,  -- Auto-requested by automation (e.g., user watching S1 → request S2)
    parent_request_id UUID REFERENCES requests(id) ON DELETE SET NULL  -- Link to parent request (e.g., S2 request triggered by S1)
);

-- Request votes (upvoting)
CREATE TABLE request_votes (
    request_id UUID REFERENCES requests(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    voted_at TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (request_id, user_id)
);

-- Request comments (integrated with ticketing system)
CREATE TABLE request_comments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    request_id UUID REFERENCES requests(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    comment TEXT NOT NULL,
    is_admin BOOLEAN DEFAULT FALSE,           -- Admin comment (highlighted in UI)
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Request quotas (per user)
CREATE TABLE request_quotas (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    daily_limit INT DEFAULT 5,
    weekly_limit INT DEFAULT 20,
    monthly_limit INT DEFAULT 50,
    daily_used INT DEFAULT 0,
    weekly_used INT DEFAULT 0,
    monthly_used INT DEFAULT 0,
    -- Storage quotas per content type
    storage_quota_movies_gb DECIMAL(10,2) DEFAULT 500,
    storage_quota_tvshows_gb DECIMAL(10,2) DEFAULT 1000,
    storage_quota_music_gb DECIMAL(10,2) DEFAULT 200,
    storage_quota_audiobooks_gb DECIMAL(10,2) DEFAULT 100,
    storage_quota_books_gb DECIMAL(10,2) DEFAULT 50,
    storage_quota_podcasts_gb DECIMAL(10,2) DEFAULT 100,
    storage_quota_comics_gb DECIMAL(10,2) DEFAULT 50,
    storage_quota_adult_gb DECIMAL(10,2) DEFAULT 500,
    -- Current storage usage per content type
    storage_used_movies_gb DECIMAL(10,2) DEFAULT 0,
    storage_used_tvshows_gb DECIMAL(10,2) DEFAULT 0,
    storage_used_music_gb DECIMAL(10,2) DEFAULT 0,
    storage_used_audiobooks_gb DECIMAL(10,2) DEFAULT 0,
    storage_used_books_gb DECIMAL(10,2) DEFAULT 0,
    storage_used_podcasts_gb DECIMAL(10,2) DEFAULT 0,
    storage_used_comics_gb DECIMAL(10,2) DEFAULT 0,
    storage_used_adult_gb DECIMAL(10,2) DEFAULT 0,
    -- Reset timestamps
    last_reset_daily DATE DEFAULT CURRENT_DATE,
    last_reset_weekly DATE DEFAULT CURRENT_DATE,
    last_reset_monthly DATE DEFAULT CURRENT_DATE
);

-- Global storage quotas (server-wide)
CREATE TABLE global_storage_quotas (
    id INT PRIMARY KEY DEFAULT 1,
    total_quota_gb DECIMAL(10,2) DEFAULT 10000,
    total_used_gb DECIMAL(10,2) DEFAULT 0,
    quota_movies_gb DECIMAL(10,2) DEFAULT 3000,
    quota_tvshows_gb DECIMAL(10,2) DEFAULT 4000,
    quota_music_gb DECIMAL(10,2) DEFAULT 1000,
    quota_audiobooks_gb DECIMAL(10,2) DEFAULT 500,
    quota_books_gb DECIMAL(10,2) DEFAULT 200,
    quota_podcasts_gb DECIMAL(10,2) DEFAULT 300,
    quota_comics_gb DECIMAL(10,2) DEFAULT 200,
    quota_adult_gb DECIMAL(10,2) DEFAULT 800,
    used_movies_gb DECIMAL(10,2) DEFAULT 0,
    used_tvshows_gb DECIMAL(10,2) DEFAULT 0,
    used_music_gb DECIMAL(10,2) DEFAULT 0,
    used_audiobooks_gb DECIMAL(10,2) DEFAULT 0,
    used_books_gb DECIMAL(10,2) DEFAULT 0,
    used_podcasts_gb DECIMAL(10,2) DEFAULT 0,
    used_comics_gb DECIMAL(10,2) DEFAULT 0,
    used_adult_gb DECIMAL(10,2) DEFAULT 0,
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Request rules (auto-approval + automation)
CREATE TABLE request_rules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(200) NOT NULL,
    content_type VARCHAR(50),                 -- NULL = all content types
    condition_type VARCHAR(50) NOT NULL,      -- user_role, trust_score, release_year, watch_history, storage_available, user_genre_preference, etc.
    condition_value JSONB NOT NULL,
    action VARCHAR(50) NOT NULL DEFAULT 'auto_approve', -- auto_approve, require_approval, decline, on_hold
    enabled BOOLEAN DEFAULT TRUE,
    priority INT DEFAULT 0,                   -- Higher priority rules checked first
    automation_trigger VARCHAR(50),           -- NULL (manual rule) OR "season_completed", "user_watching", "storage_low"
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_requests_user_id ON requests(user_id);
CREATE INDEX idx_requests_status ON requests(status);
CREATE INDEX idx_requests_content_type ON requests(content_type);
CREATE INDEX idx_requests_created_at ON requests(created_at DESC);
CREATE INDEX idx_requests_priority ON requests(priority DESC);
CREATE INDEX idx_requests_parent_id ON requests(parent_request_id);
CREATE INDEX idx_request_votes_request_id ON request_votes(request_id);
CREATE INDEX idx_request_comments_request_id ON request_comments(request_id);
CREATE INDEX idx_request_rules_automation ON request_rules(automation_trigger) WHERE automation_trigger IS NOT NULL
    weekly_used INT DEFAULT 0,
    monthly_used INT DEFAULT 0,
    last_reset_daily DATE DEFAULT CURRENT_DATE,
    last_reset_weekly DATE DEFAULT CURRENT_DATE,
    last_reset_monthly DATE DEFAULT CURRENT_DATE
);

-- Request rules (auto-approval)
CREATE TABLE request_rules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(200) NOT NULL,
    content_type VARCHAR(50),                 -- NULL = all content types
    condition_type VARCHAR(50) NOT NULL,      -- user_role, trust_score, release_year, etc.
    condition_value JSONB NOT NULL,
    action VARCHAR(50) NOT NULL DEFAULT 'auto_approve', -- auto_approve, require_approval, decline
    enabled BOOLEAN DEFAULT TRUE,
    priority INT DEFAULT 0,                   -- Higher priority rules checked first
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_requests_user_id ON requests(user_id);
CREATE INDEX idx_requests_status ON requests(status);
CREATE INDEX idx_requests_content_type ON requests(content_type);
CREATE INDEX idx_requests_created_at ON requests(created_at DESC);
CREATE INDEX idx_requests_priority ON requests(priority DESC);
CREATE INDEX idx_request_votes_request_id ON request_votes(request_id);
CREATE INDEX idx_request_comments_request_id ON request_comments(request_id);
```

### Adult Content Schema (Isolated in `c` schema)

```sql
-- Adult requests table (isolated in c schema)
CREATE TABLE c.adult_requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    content_type VARCHAR(50) NOT NULL CHECK (content_type IN ('adult_movie', 'adult_scene')),
    request_subtype VARCHAR(50),              -- "scene", "studio", "performer", "tag_combination"
    external_id VARCHAR(200),                 -- StashDB ID (NULL for tag combinations)
    title VARCHAR(500) NOT NULL,
    release_year INT,
    metadata_json JSONB,                      -- Adult-specific metadata (tags, performer IDs, studio IDs)
    status VARCHAR(50) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'approved', 'processing', 'available', 'declined', 'on_hold')),
    auto_approved BOOLEAN DEFAULT FALSE,
    auto_rule_id UUID REFERENCES c.adult_request_rules(id) ON DELETE SET NULL,
    approved_by_user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    approved_at TIMESTAMPTZ,
    declined_reason TEXT,
    priority INT DEFAULT 0,
    votes_count INT DEFAULT 0,
    integration_id VARCHAR(200),              -- Whisparr ID (after approval)
    integration_status VARCHAR(100),          -- Whisparr status
    estimated_size_gb DECIMAL(10,2),
    actual_size_gb DECIMAL(10,2),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    available_at TIMESTAMPTZ,
    triggered_by_automation BOOLEAN DEFAULT FALSE,
    parent_request_id UUID REFERENCES c.adult_requests(id) ON DELETE SET NULL
);

-- Adult request votes
CREATE TABLE c.adult_request_votes (
    request_id UUID REFERENCES c.adult_requests(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    voted_at TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (request_id, user_id)
);

-- Adult request comments
CREATE TABLE c.adult_request_comments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    request_id UUID REFERENCES c.adult_requests(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    comment TEXT NOT NULL,
    is_admin BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Adult request quotas (per user)
CREATE TABLE c.adult_request_quotas (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    daily_limit INT DEFAULT 5,
    weekly_limit INT DEFAULT 20,
    monthly_limit INT DEFAULT 50,
    daily_used INT DEFAULT 0,
    weekly_used INT DEFAULT 0,
    monthly_used INT DEFAULT 0,
    storage_quota_adult_gb DECIMAL(10,2) DEFAULT 500,
    storage_used_adult_gb DECIMAL(10,2) DEFAULT 0,
    last_reset_daily DATE DEFAULT CURRENT_DATE,
    last_reset_weekly DATE DEFAULT CURRENT_DATE,
    last_reset_monthly DATE DEFAULT CURRENT_DATE
);

-- Adult request rules (auto-approval + automation)
CREATE TABLE c.adult_request_rules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(200) NOT NULL,
    content_type VARCHAR(50),                 -- 'adult_movie', 'adult_scene', or NULL for all
    condition_type VARCHAR(50) NOT NULL,
    condition_value JSONB NOT NULL,
    action VARCHAR(50) NOT NULL DEFAULT 'auto_approve',
    enabled BOOLEAN DEFAULT TRUE,
    priority INT DEFAULT 0,
    automation_trigger VARCHAR(50),
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_adult_requests_user_id ON c.adult_requests(user_id);
CREATE INDEX idx_adult_requests_status ON c.adult_requests(status);
CREATE INDEX idx_adult_requests_content_type ON c.adult_requests(content_type);
CREATE INDEX idx_adult_requests_created_at ON c.adult_requests(created_at DESC);
CREATE INDEX idx_adult_requests_priority ON c.adult_requests(priority DESC);
CREATE INDEX idx_adult_request_votes_request_id ON c.adult_request_votes(request_id);
CREATE INDEX idx_adult_request_comments_request_id ON c.adult_request_comments(request_id);
```

---

## API Endpoints

### User Endpoints
```bash
# Search content
GET  /api/v1/requests/search?type=movie&query=Matrix
GET  /api/v1/requests/search?type=tvshow&query=Breaking Bad
GET  /api/v1/requests/search?type=music&query=Radiohead

# Submit request
POST /api/v1/requests
{
  "content_type": "movie",
  "external_id": "603",  // TMDb ID
  "title": "The Matrix",
  "release_year": 1999
}

# List user's requests
GET  /api/v1/requests?user_id=me&status=pending

# Get request detail
GET  /api/v1/requests/{id}

# Vote on request
POST /api/v1/requests/{id}/vote

# Comment on request
POST /api/v1/requests/{id}/comments
```

### Adult Content Endpoints (Isolated - `/api/v1/c/` namespace)

**⚠️ CRITICAL: Adult requests use separate API namespace `/api/v1/c/`**

```bash
# Search adult content (StashDB)
GET  /api/v1/c/requests/search?type=scene&query=...
GET  /api/v1/c/requests/search?type=studio&query=...
GET  /api/v1/c/requests/search?type=performer&query=...

# Submit adult request
POST /api/v1/c/requests
{
  "content_type": "adult_movie",
  "request_subtype": "scene",  // "scene", "studio", "performer", "tag_combination"
  "external_id": "stashdb-uuid",
  "title": "Scene Title",
  "metadata_json": {
    "tags": ["VR", "POV"],
    "performer_ids": ["uuid1", "uuid2"],
    "studio_id": "studio-uuid"
  }
}

# List user's adult requests
GET  /api/v1/c/requests?user_id=me&status=pending

# Get adult request detail
GET  /api/v1/c/requests/{id}

# Vote on adult request
POST /api/v1/c/requests/{id}/vote

# Comment on adult request
POST /api/v1/c/requests/{id}/comments
```

### Admin Endpoints (Adult: `/api/v1/c/admin/`)

```bash
# List all adult requests (isolated)
GET  /api/v1/c/admin/requests?status=pending

# Approve adult request
PUT  /api/v1/c/admin/requests/{id}/approve

# Decline adult request
PUT  /api/v1/c/admin/requests/{id}/decline

# Manage adult quotas
PUT  /api/v1/c/admin/users/{user_id}/quota

# Manage adult request rules
GET  /api/v1/c/admin/request-rules
POST /api/v1/c/admin/request-rules
PUT  /api/v1/c/admin/request-rules/{id}
DEL  /api/v1/c/admin/request-rules/{id}
```

### Admin Endpoints (Non-Adult)
```bash
# List all requests (with filters)
GET  /api/v1/admin/requests?status=pending&content_type=movie

# Approve request
PUT  /api/v1/admin/requests/{id}/approve

# Decline request
PUT  /api/v1/admin/requests/{id}/decline
{
  "reason": "Not available in region"
}

# Set priority
PUT  /api/v1/admin/requests/{id}/priority
{
  "priority": 10
}

# Batch approve
POST /api/v1/admin/requests/batch-approve
{
  "request_ids": ["uuid1", "uuid2", "uuid3"]
}

# Manage quotas
PUT  /api/v1/admin/users/{user_id}/quota
{
  "daily_limit": 10,
  "weekly_limit": 50,
  "monthly_limit": 200
}

# Manage rules
GET  /api/v1/admin/request-rules
POST /api/v1/admin/request-rules
PUT  /api/v1/admin/request-rules/{id}
DEL

### Rule 5: Auto-Request Next Season (Watch-Based Automation)
```json
{
  "name": "Auto-request S2 when user watching S1",
  "content_type": "tvshow",
  "condition_type": "watch_history",
  "condition_value": {
    "season_completed_percentage": 80,
    "trigger": "auto_request_next_season"
  },
  "action": "auto_approve",
  "automation_trigger": "season_completed",
  "priority": 70
}
```

### Rule 6: Hold Requests if Storage Low
```json
{
  "name": "Hold requests if storage <100GB",
  "content_type": null,
  "condition_type": "storage_available",
  "condition_value": {"min_free_gb": 100},
  "action": "on_hold",
  "automation_trigger": "storage_low",
  "priority": 100
}
```

### Rule 7: Auto-Decline if User Never Watches Genre
```json
{
  "name": "Decline horror requests for users who never watch horror",
  "content_type": "movie",
  "condition_type": "user_genre_preference",
  "condition_value": {
    "genre": "Horror",
    "watch_count": 0,
    "requests_declined": 3
  },
  "action": "decline",
  "priority": 60
}
```

### Rule 8: Auto-Approve if User Frequently Watches Genre
```json
{
  "name": "Auto-approve sci-fi for sci-fi fans",
  "content_type": "movie",
  "condition_type": "user_genre_preference",
  "condition_value": {
    "genre": "Science Fiction",
    "watch_count_min": 20,
    "completion_rate_min": 0.8
  },
  "action": "auto_approve",
  "Modular Architecture

### Core Request System (`internal/service/requests/`)
- `core.go`: Base request service (create, approve, decline, quota validation)
- `rule_engine.go`: Rule evaluation engine (condition matching, priority sorting)
- `automation.go`: Automation triggers (watch history analysis, storage monitoring)
- `storage.go`: Storage quota management (disk usage tracking, projection)

### Per-Content-Type Request Modules (`internal/service/requests/modules/`)
- `movie.go`: Movie request module (TMDb search, Radarr integration, size estimation)
- `tvshow.go`: TV show request module (TheTVDB search, Sonarr integration, season selection, intelligent season automation)
- `music.go`: Music request module (MusicBrainz search, Lidarr integration, artist/album requests)
- `audiobook.go`: Audiobook request module (Audible search, Audiobookshelf integration)
- `book.go`: Book request module (Goodreads search, Readarr integration)
- `podcast.go`: Podcast request module (RSS feed lookup, Audiobookshelf API)
- `comic.go`: Comic request module (ComicVine search, Mylar3 integration)

### Adult Content Request Module (ISOLATED)
**Location**: `internal/content/c/requests/` (NOT in `internal/service/requests/modules/`)
**Database**: `c` schema only (`c.adult_requests`, `c.adult_request_votes`, etc.)
**API**: `/api/v1/c/requests/*` namespace

- `adult.go`: Adult content request module (StashDB search, Whisparr integration, studio/performer/tag requests)
- Complete isolation from public request system
- Separate quota management, rule engine, automation triggers

### Integration with Other Systems
- **Ticketing System**: Link requests to support tickets (user feedback, issues, feature requests)
- **Rating System**: Use user ratings to prioritize requests (high-rated content = higher priority)
- **Analytics Service**: Track request patterns, popular content, storage trends
- **Notification Service**: Multi-channel notifications (Email, Discord, Telegram, in-app)

---

## Implementation Phases

### Phase 1: Core Request System (Week 1)
- [ ] PostgreSQL schema (requests, votes, comments, quotas, rules, global storage)
- [ ] Core request service (create, approve, decline, on_hold)
- [ ] Quota enforcement (request limits + storage quotas)
- [ ] API endpoints (user + admin)
- [ ] Storage quota tracking (real-time disk usage)

### Phase 2: Content Search Integration (Week 2)
- [ ] Movie module: TMDb search
- [ ] TV show module: TheTVDB search + season selection
- [ ] Music module: MusicBrainz search (artist/album)
- [ ] Book/Audiobook modules: Goodreads/Audible search
- [ ] Podcast module: RSS feed lookup
- [ ] Adult module: StashDB search (scene/studio/performer/tag)

### Phase 3: Arr Integration (Week 3)
- [ ] Movie module: Radarr integration (add movie on approval)
- [ ] TV show module: Sonarr integration (add show + seasons on approval)
- [ ] Music module: Lidarr integration (add artist/album on approval)
- [ ] Book module: Readarr integration (add book on approval)
- [ ] Adult module: Whisparr integration (add scene/studio/performer on approval)

### Phase 4: Audiobookshelf Integration (Week 3)
- [ ] Podcast module: Audiobookshelf API (add podcast on approval)
- [ ] Audiobook module: Audiobookshelf integration (manual workflow OR automated)

### Phase 5: Rule Engine + Automation (Week 4)
- [ ] Rule engine (condition evaluation, priority sorting)
- [ ] Watch-based automation (S1 completed → request S2)
- [ ] Storage-aware rules (auto-decline if quota exceeded)
- [ ] User behavior analysis (genre preferences, watch history)
- [ ] Content lifecycle management (auto-delete unwatched content)

### Phase 6: Notifications (Week 4)
- [ ] Email notifications (request approved/available/declined)
- [ ] Discord webhooks
- [ ] Telegram notifications
- [ ] In-app notifications (Svelte 5 UI)

### Phase 7: Frontend (Week 5-6)
- [ ] **Inline UI** (integrated in Revenge UI, not external)
- [ ] Request submission form (per content type, dynamic fields)
- [ ] Adult request UI (studio/performer/tag selection)
- [ ] Request list (user view, filter by status/content type)
- [ ] Request detail page (comments, votes, admin actions)
- [ ] Admin approval dashboard (batch approval, priority management)
- [ ] Rule management UI (create/edit automation rules)
- [ ] Quota management UI (per-user storage quotas, global quotas)
- [ ] Storage analytics dashboard (real-time usage, projections)

### Phase 8: Ticketing + Rating Integration (Week 6)
- [ ] Link requests to ticketing system (user feedback, issues)
- [ ] Integrate user rating system (high-rated content = higher priority)
- [ ] Request-to-ticket conversion (declined request → create ticket)
- [ ] Rating-based auto-approval (highly-rated content auto-approved)

### Phase 9: Analytics Integration (Week 7)
- [ ] Request pattern analysis (popular content, trends)
- [ ] Storage trend analysis (growth projections)
- [ ] User engagement metrics (request approval rate, watch rate)
- [ ] Admin dashboard (request analytics, storage health)

**Total Estimated Time**: 7Restricted
```json
{
  "name": "Decline requests from blocked countries",
  "content_type": null,
  "condition_type": "user_country",
  "condition_value": {"blocked_countries": ["XX", "YY"]},
  "action": "decline",
  "priority": 95
}
```

---

## Integration Flow

### Radarr (Movies)
```go
func (s *RequestService) ApproveMovieRequest(ctx context.Context, requestID uuid.UUID) error {
    req, _ := s.GetRequest(ctx, requestID)

    // Add to Radarr
    radarrMovieID, err := s.radarrClient.AddMovie(ctx, &radarr.Movie{
        TmdbID: req.ExternalID,
        Title: req.Title,
---

## Advanced Automation Examples

### Example 1: Intelligent Season Automation (TV Shows)

**Scenario**: User watching Breaking Bad S1

```
1. User starts watching Breaking Bad S1E01
   ↓
2. Analytics service tracks watch progress
   ↓
3. User completes 80% of S1 (watched 10/13 episodes)
   ↓
4. Automation rule triggered: "Auto-request S2 when S1 80% complete"
   ↓
5. Request service creates automated request for S2
   ↓
6. Rule engine evaluates request:
   - User has high trust score (90) → Auto-approve
   - Storage quota available (500GB free) → Proceed
   - Nobody else watching Breaking Bad S2 → Lower priority
   ↓
7. Request auto-approved, added to Sonarr
   ↓
8. Sonarr downloads S2
   ↓
9. User notified: "Breaking Bad S2 is now available!"
```

**Scenario**: Nobody watching show yet

```
1. User requests entire series (all 5 seasons)
   ↓
2. Rule engine evaluates:
   - Nobody has watched any episode yet
   - Rule: "Only fetch S1 if no watch history"
   ↓
3. Request S1: Approved
4. Requests S2-S5: On Hold (wait for S1 engagement)
   ↓
5. User watches S1
   ↓
6. Automation triggers: Release S2 from hold → Approve
```

### Example 2: Adult Content Studio Request

**Scenario**: User wants all content from specific studio

```
1. User navigates to adult request UI
   ↓
2. Selects "Request by Studio"
   ↓
3. Searches StashDB: "Studio XYZ"
   ↓
4. Submits request: "All scenes from Studio XYZ"
   ↓
5. Admin approves
   ↓
6. Request service queries StashDB:
   - Finds 50 scenes from Studio XYZ
   - Estimates total size: 250GB
   ↓
7. Storage check:
   - User quota: 500GB adult content (200GB used)
   - Global quota: 800GB adult content (400GB used)
   - Available: 300GB (user), 400GB (global) → Proceed
   ↓
8. Adds all 50 scenes to Whisparr (monitored)
   ↓
9. Whisparr downloads scenes automatically
   ↓
10. User notified as each scene becomes available
```

### Example 3: Storage-Aware Request Management

**Scenario**: Low disk space, smart request handling

```
1. User requests 4K movie (estimated 80GB)
   ↓
2. Rule engine evaluates:
   - Global storage: 50GB free (< 100GB threshold)
   - Rule: "Hold requests if storage < 100GB"
   ↓
3. Request status: On Hold
   ↓
4. Admin notified: "Low storage, request on hold"
   ↓
5. Background job runs: Content lifecycle cleanup
   - Identifies unwatched movies (> 90 days old)
   - Deletes 3 unwatched movies (200GB freed)
   ↓
6. Storage now: 250GB free
   ↓
7. Automation releases request from hold → Approved
   ↓
8. User notified: "Your request is now approved!"
```

---

## UI Design (Inline in Revenge)

### User Request Flow (Inline UI)

```
Main Navigation → "Requests" → Request Dashboard
                                      ↓
                            ┌─────────┴─────────┐
                            │                   │
                      My Requests         Submit Request
                            │                   │
                            │         ┌─────────┴─────────┐
                            │         │                   │
                            │    Content Type      Advanced Options
                            │         │                   │
                            │    [Movies]          [Storage: 200GB/500GB]
                            │    [TV Shows]        [Quality: Auto/1080p/4K]
                            │    [Music]           [Priority: Normal]
                            │    [Adult ▼]         [Season Selection]
                            │         │                   │
                            │    ┌────┴────┐             │
                            │    │         │             │
                            │  Scene   Studio/Performer  │
                            │    │         │             │
                            │    Search   [Studio XYZ ▼] │
                            │    TMDb     [Performer ABC]│
                            │              [Tags: VR+POV]│
                            │                   │
                            │              Submit Request
                            │
                      Request List
                            │
                      [#1234] Breaking Bad S2
                      Status: Approved, Processing
                      Priority: High (15 votes)
                      Storage: 12GB estimated
                      ETA: 2 hours
                      [View Details] [Vote] [Comment]
```

### Admin Approval Dashboard (Inline UI)

```
Admin Panel → "Requests" → Approval Queue
                                  ↓
                    ┌─────────────┴─────────────┐
                    │                           │
              Pending Requests          Storage Analytics
                    │                           │
        [Filter: All Types ▼]       [Movies: 2.8TB / 3TB]
        [Sort: Priority ▼]          [TV Shows: 3.5TB / 4TB]
                    │                [Music: 0.8TB / 1TB]
        ┌───────────┴───────────┐   [Adult: 0.5TB / 0.8TB]
        │                       │   [Total: 7.6TB / 10TB]
   Bulk Actions          Request Card
        │                       │
   [Select All]      [#1234] The Matrix (1999)
   [Approve]         Type: Movie (4K)
   [Decline]         User: john.doe (Trust: 85)
   [Set Priority]    Estimated: 80GB
                     Votes: 12
                     Rule: "Auto-approve sci-fi fans"
                     [Approve] [Decline] [On Hold] [Details]
```

---

## Notes

- **NO Overseerr/Jellyseerr integration** - native system only
- **Inline UI**: Fully integrated in Revenge UI (Svelte 5 runes), no external apps
- **Modular architecture**: Core system + per-content-type modules (movie, tvshow, music, audiobook, book, podcast, comic, adult)
- **Advanced automation**: Watch history analysis, intelligent season requests, storage-aware rules
- **Adult content flexibility**: Request by scene, studio, performer, tag combinations
- **Deep integrations**: Ticketing system, rating system, analytics, storage quotas
- **Storage quotas**: Per-user quotas per content type, global quotas, real-time tracking
- **Intelligent automation**: S1 completed → auto-request S2, nobody watching → hold S2-S5
- Auto-approval rules provide flexibility (trust-based, role-based, content-based, watch-based, storage-based)
- Voting system creates priority queue (community-driven)
- All content types supported (not just movies/TV like Overseerr)
- Podcast requests integrate directly with Audiobookshelf API
- Content lifecycle management: Auto-delete unwatched content to free space for new requests

    // Update request
    req.Status = "processing"
    req.IntegrationID = fmt.Sprintf("%d", radarrMovieID)
    s.UpdateRequest(ctx, req)

    return nil
}
```

### Audiobookshelf (Podcasts)
```go
func (s *RequestService) ApprovePodcastRequest(ctx context.Context, requestID uuid.UUID) error {
    req, _ := s.GetRequest(ctx, requestID)

    // Add to Audiobookshelf
    podcastID, err := s.audiobookshelfClient.AddPodcast(ctx, &audiobookshelf.Podcast{
        FeedURL: req.MetadataJSON["feed_url"].(string),
        LibraryID: s.config.PodcastLibraryID,
        AutoDownloadEpisodes: true,
    })

    req.Status = "available"  // Podcast added = immediately available
    req.IntegrationID = podcastID
    req.AvailableAt = time.Now()
    s.UpdateRequest(ctx, req)

    // Notify user
    s.notifier.NotifyRequestAvailable(ctx, req)

    return nil
}
```

---

## Implementation Phases

### Phase 1: Core Request System (Week 1)
- [ ] PostgreSQL schema (requests, votes, comments, quotas, rules)
- [ ] Request service (create, approve, decline)
- [ ] API endpoints (user + admin)
- [ ] Quota enforcement

### Phase 2: Content Search Integration (Week 2)
- [ ] TMDb search (movies)
- [ ] TheTVDB search (TV shows)
- [ ] MusicBrainz search (music)
- [ ] Goodreads/Audible search (audiobooks/books)
- [ ] Podcast search (RSS feed lookup)

### Phase 3: Arr Integration (Week 3)
- [ ] Radarr integration (add movie on approval)
- [ ] Sonarr integration (add TV show on approval)
- [ ] Lidarr integration (add music on approval)
- [ ] Readarr integration (add book on approval)
- [ ] Whisparr integration (add adult content on approval)

### Phase 4: Audiobookshelf Integration (Week 3)
- [ ] Podcast addition via API
- [ ] Audiobook manual workflow (OR script integration)

### Phase 5: Notifications (Week 4)
- [ ] Email notifications (request approved/available)
- [ ] Discord webhooks
- [ ] Telegram notifications
- [ ] In-app notifications

### Phase 6: Frontend (Week 4-5)
- [ ] Request submission form (per content type)
- [ ] Request list (user view)
- [ ] Request detail page (with comments/votes)
- [ ] Admin approval dashboard
- [ ] Rule management UI
- [ ] Quota management UI

### Phase 7: Auto-Approval Rules (Week 5)
- [ ] Rule engine
- [ ] Rule CRUD
- [ ] Condition evaluation (role, trust score, release year, country, etc.)

**Total Estimated Time**: 5-6 weeks

---

## Related Documentation

- [Radarr Integration](../../integrations/servarr/RADARR.md)
- [Sonarr Integration](../../integrations/servarr/SONARR.md)
- [Lidarr Integration](../../integrations/servarr/LIDARR.md)
- [Whisparr Integration](../../integrations/servarr/WHISPARR.md)
- [Readarr Integration](../../integrations/servarr/READARR.md)
- [Audiobookshelf Integration](../../integrations/audiobook/AUDIOBOOKSHELF.md)
- [User Management](../../operations/USER_MANAGEMENT.md) - User roles & permissions
- [Notifications](../../technical/NOTIFICATIONS.md) - Email/Discord/Telegram

---

## Notes

- **NO Overseerr/Jellyseerr integration** - native system only
- Podcast requests integrate directly with Audiobookshelf API
- Auto-approval rules provide flexibility (trust-based, role-based, content-based)
- Quota system prevents abuse
- Voting system creates priority queue (community-driven)
- All content types supported (not just movies/TV like Overseerr)
