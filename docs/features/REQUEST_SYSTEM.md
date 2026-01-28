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
- **Quota management**: Per-user request limits (daily/weekly/monthly)
- **Request rules**: Auto-approve based on user role/trust score
- **Batch approval**: Approve multiple requests at once
- **Integration triggers**: Automatically add to Radarr/Sonarr/Lidarr on approval
- **Priority management**: Admin can prioritize requests
- **Request history**: Full audit trail

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
                                         Admin approves
                                                    ↓
                      Revenge adds podcast to Audiobookshelf via API
                                                    ↓
                      Audiobookshelf fetches episodes → Request: Available
```

### Comics (Mylar3 Integration - Future)
```
User searches ComicVine → Selects series/issue → Submits request
                                                      ↓
                                           Admin approves
                                                      ↓
                              Revenge adds to Mylar3 (OR manual)
                                                      ↓
                              Downloads → Imports → Notify user
```

### Adult Content (Whisparr Integration)
```
User searches StashDB → Selects scene → Submits request
                                            ↓
                                 Admin approves
                                            ↓
                        Revenge adds to Whisparr
                                            ↓
                        Whisparr downloads → Imports → Notify user
```

---

## PostgreSQL Schema

```sql
-- Requests table
CREATE TABLE requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    request_number SERIAL UNIQUE NOT NULL,    -- Human-readable ID (e.g., #1234)
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    content_type VARCHAR(50) NOT NULL,        -- movie, tvshow, music_album, audiobook, book, podcast, comic, adult_scene
    external_id VARCHAR(200) NOT NULL,        -- TMDb ID, TheTVDB ID, MusicBrainz ID, etc.
    title VARCHAR(500) NOT NULL,
    release_year INT,
    metadata_json JSONB,                      -- Content-specific metadata
    status VARCHAR(50) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'approved', 'processing', 'available', 'declined')),
    auto_approved BOOLEAN DEFAULT FALSE,
    approved_by_user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    approved_at TIMESTAMPTZ,
    declined_reason TEXT,
    priority INT DEFAULT 0,                   -- Higher = more important
    votes_count INT DEFAULT 0,
    integration_id VARCHAR(200),              -- Radarr/Sonarr/Lidarr ID (after approval)
    integration_status VARCHAR(100),          -- Radarr/Sonarr status
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    available_at TIMESTAMPTZ
);

-- Request votes (upvoting)
CREATE TABLE request_votes (
    request_id UUID REFERENCES requests(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    voted_at TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (request_id, user_id)
);

-- Request comments
CREATE TABLE request_comments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    request_id UUID REFERENCES requests(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    comment TEXT NOT NULL,
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

### Admin Endpoints
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
DELETE /api/v1/admin/request-rules/{id}
```

---

## Auto-Approval Rules Examples

### Rule 1: Auto-Approve Admins
```json
{
  "name": "Auto-approve all admin requests",
  "content_type": null,
  "condition_type": "user_role",
  "condition_value": {"role": "admin"},
  "action": "auto_approve",
  "priority": 100
}
```

### Rule 2: Auto-Approve High Trust Users
```json
{
  "name": "Auto-approve users with trust score >80",
  "content_type": null,
  "condition_type": "trust_score",
  "condition_value": {"min": 80},
  "action": "auto_approve",
  "priority": 90
}
```

### Rule 3: Require Approval for Old Movies
```json
{
  "name": "Require approval for movies before 1980",
  "content_type": "movie",
  "condition_type": "release_year",
  "condition_value": {"max": 1979},
  "action": "require_approval",
  "priority": 50
}
```

### Rule 4: Auto-Decline Geo-Restricted
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
        Year: req.ReleaseYear,
        QualityProfileID: s.config.DefaultQualityProfile,
        RootFolderPath: s.config.DefaultRootFolder,
        Monitored: true,
        AddOptions: &radarr.AddOptions{
            SearchForMovie: true,  // Trigger search immediately
        },
    })
    
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
