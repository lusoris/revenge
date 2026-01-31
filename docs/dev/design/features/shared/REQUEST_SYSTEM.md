# Native Request System

> Content request management for all modules - replaces Overseerr/Jellyseerr


<!-- TOC-START -->

## Table of Contents

- [Status](#status)
- [Design Decision](#design-decision)
- [Core Features (from Overseerr/Jellyseerr)](#core-features-from-overseerrjellyseerr)
  - [User-Facing Features](#user-facing-features)
  - [Polls System](#polls-system)
  - [Admin Features](#admin-features)
  - [Advanced Automation Features](#advanced-automation-features)
- [Per-Module Request Handling](#per-module-request-handling)
  - [Movies (Radarr Integration)](#movies-radarr-integration)
  - [TV Shows (Sonarr Integration)](#tv-shows-sonarr-integration)
  - [Music (Lidarr Integration)](#music-lidarr-integration)
  - [Audiobooks (Audiobookshelf Integration)](#audiobooks-audiobookshelf-integration)
  - [Books (Chaptarr Integration)](#books-chaptarr-integration)
  - [Podcasts (Audiobookshelf Integration)](#podcasts-audiobookshelf-integration)
  - [Rule 5: Auto-Request Next Season (Watch-Based Automation)](#rule-5-auto-request-next-season-watch-based-automation)
  - [Rule 6: Hold Requests if Storage Low](#rule-6-hold-requests-if-storage-low)
  - [Rule 7: Auto-Decline if User Never Watches Genre](#rule-7-auto-decline-if-user-never-watches-genre)
  - [Rule 8: Auto-Approve if User Frequently Watches Genre](#rule-8-auto-approve-if-user-frequently-watches-genre)
- [Advanced Automation Examples](#advanced-automation-examples)
  - [Example 1: Intelligent Season Automation (TV Shows)](#example-1-intelligent-season-automation-tv-shows)
  - [Example 2: Adult Content Studio Request](#example-2-adult-content-studio-request)
  - [Example 3: Storage-Aware Request Management](#example-3-storage-aware-request-management)
- [UI Design (Inline in Revenge)](#ui-design-inline-in-revenge)
  - [User Request Flow (Inline UI)](#user-request-flow-inline-ui)
  - [Admin Approval Dashboard (Inline UI)](#admin-approval-dashboard-inline-ui)
- [Notes](#notes)
- [Implementation Phases](#implementation-phases)
  - [Phase 1: Core Request System (Week 1)](#phase-1-core-request-system-week-1)
  - [Phase 2: Content Search Integration (Week 2)](#phase-2-content-search-integration-week-2)
  - [Phase 3: Arr Integration (Week 3)](#phase-3-arr-integration-week-3)
  - [Phase 4: Audiobookshelf Integration (Week 3)](#phase-4-audiobookshelf-integration-week-3)
  - [Phase 5: Notifications (Week 4)](#phase-5-notifications-week-4)
  - [Phase 6: Frontend (Week 4-5)](#phase-6-frontend-week-4-5)
  - [Phase 7: Auto-Approval Rules (Week 5)](#phase-7-auto-approval-rules-week-5)
- [Implementation Checklist](#implementation-checklist)
  - [Phase 1: Core Infrastructure](#phase-1-core-infrastructure)
  - [Phase 2: Database](#phase-2-database)
  - [Phase 3: Service Layer - Core](#phase-3-service-layer---core)
  - [Phase 4: Content Search Modules](#phase-4-content-search-modules)
  - [Phase 5: Arr Integration](#phase-5-arr-integration)
  - [Phase 6: Audiobookshelf Integration](#phase-6-audiobookshelf-integration)
  - [Phase 7: Background Jobs - Automation](#phase-7-background-jobs---automation)
  - [Phase 8: Polls System](#phase-8-polls-system)
  - [Phase 9: Notifications](#phase-9-notifications)
  - [Phase 10: API Integration](#phase-10-api-integration)
- [Sources & Cross-References](#sources-cross-references)
  - [Cross-Reference Indexes](#cross-reference-indexes)
  - [Referenced Sources](#referenced-sources)
- [Related Design Docs](#related-design-docs)
  - [In This Section](#in-this-section)
  - [Related Topics](#related-topics)
  - [Indexes](#indexes)
- [Related](#related)
- [Notes](#notes)

<!-- TOC-END -->

## Status

| Dimension | Status | Notes |
|-----------|--------|-------|
| Design | ✅ | Full design with DB schema, automation rules, polls system |
| Sources | 🟡 | Replaces Overseerr/Jellyseerr |
| Instructions | ✅ | Implementation checklist added |
| Code | 🔴 |  |
| Linting | 🔴 |  |
| Unit Testing | 🔴 |  |
| Integration Testing | 🔴 |  |
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
- **Polls**: Community polls to decide what content to add next

### Polls System

Polls allow admins/mods to create community votes on content decisions:

**Poll Types**:
- **Manual Polls**: Admin creates poll with specific options ("Which Marvel series should we add next?")
- **Rule-Based Polls**: Auto-generated based on conditions ("Top 5 most requested movies this month")
- **Tie-Breaker Polls**: When multiple requests have similar priority

**Features**:
- Multiple voting options (single-choice or ranked voting)
- Time-limited polls (e.g., vote closes in 7 days)
- Minimum participation threshold (e.g., at least 10 users must vote)
- Results visibility (hidden until poll ends, or real-time)
- Notification when poll opens/closes
- Auto-approve winning option(s)

**Database Schema (Polls)**:
```sql
-- Polls table
CREATE TABLE request_polls (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(500) NOT NULL,
    description TEXT,
    content_type VARCHAR(50),                -- NULL = mixed content types
    poll_type VARCHAR(50) NOT NULL DEFAULT 'manual', -- manual, rule_based, tie_breaker
    voting_style VARCHAR(50) NOT NULL DEFAULT 'single', -- single, ranked, multi_select
    min_votes INT DEFAULT 1,                 -- Minimum votes required for valid result
    max_selections INT DEFAULT 1,            -- For multi_select: how many can user pick
    starts_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ends_at TIMESTAMPTZ NOT NULL,
    show_results BOOLEAN DEFAULT FALSE,      -- Show results before poll ends
    auto_approve_winner BOOLEAN DEFAULT TRUE,
    status VARCHAR(50) NOT NULL DEFAULT 'active', -- draft, active, closed, cancelled
    created_by_user_id UUID REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Poll options (content items to vote on)
CREATE TABLE request_poll_options (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    poll_id UUID REFERENCES request_polls(id) ON DELETE CASCADE,
    request_id UUID REFERENCES requests(id) ON DELETE CASCADE, -- Link to existing request
    title VARCHAR(500) NOT NULL,             -- Display title
    description TEXT,
    external_id VARCHAR(200),                -- TMDb/TVDB/etc. ID
    metadata_json JSONB,                     -- Additional metadata
    display_order INT DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Poll votes
CREATE TABLE request_poll_votes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    poll_id UUID REFERENCES request_polls(id) ON DELETE CASCADE,
    option_id UUID REFERENCES request_poll_options(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    rank INT DEFAULT 1,                      -- For ranked voting (1 = first choice)
    voted_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(poll_id, user_id, option_id)      -- Prevent duplicate votes on same option
);

CREATE INDEX idx_request_polls_status ON request_polls(status, ends_at);
CREATE INDEX idx_poll_votes_poll_id ON request_poll_votes(poll_id);
CREATE INDEX idx_poll_votes_user_id ON request_poll_votes(user_id);
```

**Rule-Based Poll Generation**:
```json
{
  "name": "Monthly Top Requested Movies Poll",
  "trigger": "scheduled",
  "schedule": "0 0 1 * *",  // First of each month
  "content_type": "movie",
  "selection_criteria": {
    "status": "pending",
    "min_votes": 3,
    "order_by": "votes_count DESC",
    "limit": 5
  },
  "poll_config": {
    "title": "Top Movie Requests - {month} {year}",
    "voting_style": "ranked",
    "duration_days": 7,
    "auto_approve_winner": true,
    "approve_top_n": 2  // Approve top 2 winners
  }
}
```

**QAR Polls (Schema `qar`)**:
Adult content has separate poll tables in schema `qar` with identical structure but additional fields for crew/port/flag filtering. See [00_SOURCE_OF_TRUTH.md](../../00_SOURCE_OF_TRUTH.md#qar-obfuscation-terminology) for terminology.

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

### Books (Chaptarr Integration)
```
User searches Goodreads → Selects book → Submits request
                                              ↓
                                   Admin approves
                                              ↓
                          Revenge adds to Chaptarr
                                              ↓
                          Chaptarr downloads → Imports → Notify user
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
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
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
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
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
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
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
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
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
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
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
-- QAR requests table (isolated in qar schema)
CREATE TABLE qar.requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    content_type VARCHAR(50) NOT NULL CHECK (content_type IN ('expedition', 'voyage')),
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
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
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
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
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

-- QAR request rules (auto-approval + automation)
CREATE TABLE qar.request_rules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(200) NOT NULL,
    content_type VARCHAR(50),                 -- 'expedition', 'voyage', or NULL for all
    condition_type VARCHAR(50) NOT NULL,
    condition_value JSONB NOT NULL,
    action VARCHAR(50) NOT NULL DEFAULT 'auto_approve',
    enabled BOOLEAN DEFAULT TRUE,
    priority INT DEFAULT 0,
    automation_trigger VARCHAR(50),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
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

### Adult Content Endpoints (Isolated - `/api/v1/legacy/` namespace)

**⚠️ CRITICAL: Adult requests use separate API namespace `/api/v1/legacy/`**

```bash
# Search adult content (StashDB)
GET  /api/v1/legacy/requests/search?type=scene&query=...
GET  /api/v1/legacy/requests/search?type=studio&query=...
GET  /api/v1/legacy/requests/search?type=performer&query=...

# Submit adult request
POST /api/v1/legacy/requests
{
  "content_type": "expedition",
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
GET  /api/v1/legacy/requests?user_id=me&status=pending

# Get adult request detail
GET  /api/v1/legacy/requests/{id}

# Vote on adult request
POST /api/v1/legacy/requests/{id}/vote

# Comment on adult request
POST /api/v1/legacy/requests/{id}/comments
```

### Admin Endpoints (Adult: `/api/v1/legacy/admin/`)

```bash
# List all adult requests (isolated)
GET  /api/v1/legacy/admin/requests?status=pending

# Approve adult request
PUT  /api/v1/legacy/admin/requests/{id}/approve

# Decline adult request
PUT  /api/v1/legacy/admin/requests/{id}/decline

# Manage adult quotas
PUT  /api/v1/legacy/admin/users/{user_id}/quota

# Manage adult request rules
GET  /api/v1/legacy/admin/request-rules
POST /api/v1/legacy/admin/request-rules
PUT  /api/v1/legacy/admin/request-rules/{id}
DEL  /api/v1/legacy/admin/request-rules/{id}
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
- `book.go`: Book request module (Goodreads search, Chaptarr integration)
- `podcast.go`: Podcast request module (RSS feed lookup, Audiobookshelf API)
- `comic.go`: Comic request module (ComicVine search, Mylar3 integration)

### Adult Content Request Module (ISOLATED)
**Location**: `internal/content/c/requests/` (NOT in `internal/service/requests/modules/`)
**Database**: `c` schema only (`c.adult_requests`, `c.adult_request_votes`, etc.)
**API**: `/api/v1/legacy/requests/*` namespace

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
- [ ] Book module: Chaptarr integration (add book on approval)
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
- [ ] Chaptarr integration (add book on approval)
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

## Implementation Checklist

**Location**: `internal/service/requests/`

### Phase 1: Core Infrastructure
- [ ] Create package structure `internal/service/requests/`
- [ ] Create modules directory `internal/service/requests/modules/`
- [ ] Define entities: `Request`, `RequestVote`, `RequestComment`, `RequestQuota`, `RequestRule`
- [ ] Define poll entities: `RequestPoll`, `PollOption`, `PollVote`
- [ ] Create repository interface `RequestRepository`
- [ ] Implement fx module `requests.Module`
- [ ] Add configuration struct for request settings

### Phase 2: Database
- [ ] Create migration `xxx_request_system.up.sql`
- [ ] Create `requests` table with status, priority, integration fields
- [ ] Create `request_votes` table for upvoting
- [ ] Create `request_comments` table for discussions
- [ ] Create `request_quotas` table with storage limits per content type
- [ ] Create `global_storage_quotas` table for server-wide limits
- [ ] Create `request_rules` table for auto-approval logic
- [ ] Create `request_polls` table for community voting
- [ ] Create `request_poll_options` table for poll choices
- [ ] Create `request_poll_votes` table for user votes
- [ ] Create QAR schema tables: `qar.requests`, `qar.request_votes`, etc.
- [ ] Create QAR poll tables in `qar` schema
- [ ] Add indexes for status, content_type, priority, created_at
- [ ] Generate sqlc queries for all CRUD operations

### Phase 3: Service Layer - Core
- [ ] Implement `RequestService` (create, approve, decline, on_hold)
- [ ] Implement quota enforcement (request limits + storage quotas)
- [ ] Implement `RuleEngine` for condition matching and priority sorting
- [ ] Implement storage tracking (real-time disk usage)
- [ ] Implement vote aggregation and priority calculation
- [ ] Add caching for frequently accessed requests

### Phase 4: Content Search Modules
- [ ] Implement `movie.go` with TMDb search integration
- [ ] Implement `tvshow.go` with TheTVDB search + season selection
- [ ] Implement `music.go` with MusicBrainz search (artist/album)
- [ ] Implement `audiobook.go` with Audible search
- [ ] Implement `book.go` with Goodreads search
- [ ] Implement `podcast.go` with RSS feed lookup
- [ ] Implement `comic.go` with ComicVine search
- [ ] Create adult module in `internal/content/c/requests/adult.go` (isolated)
- [ ] Implement StashDB search (scene/studio/performer/tag)

### Phase 5: Arr Integration
- [ ] Implement Radarr client for movie additions
- [ ] Implement Sonarr client for TV show additions (with season config)
- [ ] Implement Lidarr client for music additions
- [ ] Implement Chaptarr client for book additions
- [ ] Implement Whisparr client for adult content (in c/ module)
- [ ] Add webhook handlers for arr status updates
- [ ] Update request status when content becomes available

### Phase 6: Audiobookshelf Integration
- [ ] Implement Audiobookshelf client for podcast additions
- [ ] Implement RSS feed addition via Audiobookshelf API
- [ ] Handle audiobook manual workflow OR automated script

### Phase 7: Background Jobs - Automation
- [ ] Create `RequestApprovalWorker` for auto-approval rule processing
- [ ] Create `QuotaResetWorker` for daily/weekly/monthly quota resets
- [ ] Create `StorageSyncWorker` for disk usage tracking
- [ ] Create `SeasonAutomationWorker` for S1 completed -> request S2
- [ ] Create `ContentLifecycleWorker` for auto-delete unwatched content
- [ ] Create `PollClosingWorker` for closing polls and auto-approving winners
- [ ] Configure River periodic jobs for automation triggers

### Phase 8: Polls System
- [ ] Implement `PollService` for poll CRUD operations
- [ ] Implement voting logic (single, ranked, multi_select)
- [ ] Implement rule-based poll generation
- [ ] Implement tie-breaker poll creation
- [ ] Implement auto-approve winner on poll close
- [ ] Add poll notifications (open/close events)

### Phase 9: Notifications
- [ ] Integrate with notification service for request events
- [ ] Implement email notifications (approved/available/declined)
- [ ] Implement Discord webhook notifications
- [ ] Implement Telegram notifications
- [ ] Implement in-app notifications

### Phase 10: API Integration
- [ ] Define OpenAPI spec for request endpoints
- [ ] Generate ogen handlers for user and admin APIs
- [ ] Implement `GET /api/v1/requests/search` - content search
- [ ] Implement `POST /api/v1/requests` - submit request
- [ ] Implement `GET /api/v1/requests` - list user requests
- [ ] Implement `GET /api/v1/requests/:id` - request detail
- [ ] Implement `POST /api/v1/requests/:id/vote` - upvote
- [ ] Implement `POST /api/v1/requests/:id/comments` - add comment
- [ ] Implement admin approval endpoints (approve, decline, priority, batch)
- [ ] Implement quota management endpoints
- [ ] Implement rule management endpoints
- [ ] Implement poll endpoints (create, vote, results)
- [ ] Add adult content endpoints under `/api/v1/legacy/requests/`
- [ ] Add RBAC permission checks for all operations

---


<!-- SOURCE-BREADCRUMBS-START -->

## Sources & Cross-References

> Auto-generated section linking to external documentation sources

### Cross-Reference Indexes

- [All Sources Index](../../../sources/SOURCES_INDEX.md) - Complete list of external documentation
- [Design ↔ Sources Map](../../../sources/DESIGN_CROSSREF.md) - Which docs reference which sources

### Referenced Sources

| Source | Documentation |
|--------|---------------|
| [PostgreSQL Arrays](https://www.postgresql.org/docs/current/arrays.html) | [Local](../../../sources/database/postgresql-arrays.md) |
| [PostgreSQL JSON Functions](https://www.postgresql.org/docs/current/functions-json.html) | [Local](../../../sources/database/postgresql-json.md) |
| [River Job Queue](https://pkg.go.dev/github.com/riverqueue/river) | [Local](../../../sources/tooling/river.md) |
| [Svelte 5 Documentation](https://svelte.dev/docs/svelte/overview) | [Local](../../../sources/frontend/svelte5.md) |
| [Svelte 5 Runes](https://svelte.dev/docs/svelte/$state) | [Local](../../../sources/frontend/svelte-runes.md) |
| [SvelteKit Documentation](https://svelte.dev/docs/kit/introduction) | [Local](../../../sources/frontend/sveltekit.md) |
| [Uber fx](https://pkg.go.dev/go.uber.org/fx) | [Local](../../../sources/tooling/fx.md) |
| [ogen OpenAPI Generator](https://pkg.go.dev/github.com/ogen-go/ogen) | [Local](../../../sources/tooling/ogen.md) |
| [pgx PostgreSQL Driver](https://pkg.go.dev/github.com/jackc/pgx/v5) | [Local](../../../sources/database/pgx.md) |
| [sqlc](https://docs.sqlc.dev/en/stable/) | [Local](../../../sources/database/sqlc.md) |
| [sqlc Configuration](https://docs.sqlc.dev/en/stable/reference/config.html) | [Local](../../../sources/database/sqlc-config.md) |

<!-- SOURCE-BREADCRUMBS-END -->

<!-- DESIGN-BREADCRUMBS-START -->

## Related Design Docs

> Auto-generated cross-references to related design documentation

**Category**: [Shared](INDEX.md)

### In This Section

- [Time-Based Access Controls](ACCESS_CONTROLS.md)
- [Tracearr Analytics Service](ANALYTICS_SERVICE.md)
- [Revenge - Client Support & Device Capabilities](CLIENT_SUPPORT.md)
- [Content Rating System](CONTENT_RATING.md)
- [Revenge - Internationalization (i18n)](I18N.md)
- [Library Types](LIBRARY_TYPES.md)
- [News System](NEWS_SYSTEM.md)
- [Revenge - NSFW Toggle](NSFW_TOGGLE.md)

### Related Topics

- [Revenge - Architecture v2](../../architecture/01_ARCHITECTURE.md) _Architecture_
- [Revenge - Design Principles](../../architecture/02_DESIGN_PRINCIPLES.md) _Architecture_
- [Revenge - Metadata System](../../architecture/03_METADATA_SYSTEM.md) _Architecture_
- [Revenge - Player Architecture](../../architecture/04_PLAYER_ARCHITECTURE.md) _Architecture_
- [Plugin Architecture Decision](../../architecture/05_PLUGIN_ARCHITECTURE_DECISION.md) _Architecture_

### Indexes

- [Design Index](../../DESIGN_INDEX.md) - All design docs by category/topic
- [Source of Truth](../../00_SOURCE_OF_TRUTH.md) - Package versions and status

<!-- DESIGN-BREADCRUMBS-END -->

## Related

- [Radarr Integration](../../integrations/servarr/RADARR.md)
- [Sonarr Integration](../../integrations/servarr/SONARR.md)
- [Lidarr Integration](../../integrations/servarr/LIDARR.md)
- [Whisparr Integration](../../integrations/servarr/WHISPARR.md)
- [Chaptarr Integration](../../integrations/servarr/CHAPTARR.md)
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
