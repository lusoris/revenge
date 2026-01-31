# Podcasts

<!-- SOURCES: fx, ogen, river, sqlc, sqlc-config -->

<!-- DESIGN: features/podcasts, 01_ARCHITECTURE, 02_DESIGN_PRINCIPLES, 03_METADATA_SYSTEM -->


> RSS podcast subscription and playback


<!-- TOC-START -->

## Table of Contents

- [Status](#status)
- [Developer Resources](#developer-resources)
- [Overview](#overview)
- [Features](#features)
- [Go Packages](#go-packages)
- [Database Schema](#database-schema)
- [River Jobs](#river-jobs)
- [Go Implementation](#go-implementation)
- [API Endpoints](#api-endpoints)
- [Configuration](#configuration)
- [RBAC Permissions](#rbac-permissions)
- [Implementation Checklist](#implementation-checklist)
  - [Phase 1: Core Infrastructure](#phase-1-core-infrastructure)
  - [Phase 2: Database](#phase-2-database)
  - [Phase 3: Feed Parsing](#phase-3-feed-parsing)
  - [Phase 4: Service Layer](#phase-4-service-layer)
  - [Phase 5: Download Management](#phase-5-download-management)
  - [Phase 6: Background Jobs](#phase-6-background-jobs)
  - [Phase 7: API Integration](#phase-7-api-integration)
- [Sources & Cross-References](#sources-cross-references)
  - [Cross-Reference Indexes](#cross-reference-indexes)
  - [Referenced Sources](#referenced-sources)
- [Related Design Docs](#related-design-docs)
  - [Related Topics](#related-topics)
  - [Indexes](#indexes)
- [Related](#related)

<!-- TOC-END -->

## Status

| Dimension | Status | Notes |
|-----------|--------|-------|
| Design | ✅ | Comprehensive spec with features, schema, jobs |
| Sources | ✅ | Go package URLs listed (gofeed, gorilla/feeds) |
| Instructions | ✅ | Code examples exist, no explicit checklist |
| Code | 🔴 |  |
| Linting | 🔴 |  |
| Unit Testing | 🔴 |  |
| Integration Testing | 🔴 |  |**Priority**: LOW (Nice to have - Plex has this)
**Inspired By**: Plex Podcasts, Apple Podcasts
**Location**: `internal/content/podcasts/`

---

## Developer Resources

| Source | URL | Purpose |
|--------|-----|---------|
| RSS 2.0 Spec | [cyber.harvard.edu/rss/rss.html](https://cyber.harvard.edu/rss/rss.html) | RSS feed format |
| Podcast RSS | [podcasters.apple.com/support/823](https://podcasters.apple.com/support/823-podcast-requirements) | Apple podcast RSS extensions |
| OPML Spec | [opml.org/spec2.opml](http://opml.org/spec2.opml) | Subscription import/export format |
| iTunes Search | [developer.apple.com/library/archive/documentation/AudioVideo/Conceptual/iTuneSearchAPI](https://developer.apple.com/library/archive/documentation/AudioVideo/Conceptual/iTuneSearchAPI/) | Podcast directory search |

---

## Overview

Podcast support through RSS feed subscription, allowing users to subscribe, download, and listen to their favorite podcasts.

---

## Features

| Feature | Description |
|---------|-------------|
| RSS Subscription | Subscribe to podcast feeds |
| Auto-Download | Automatically download new episodes |
| Episode Tracking | Track played/unplayed episodes |
| Offline Playback | Download for offline listening |
| Playback Speed | Variable speed playback (0.5x - 3x) |
| Skip Silence | Optional silence skipping |
| Sleep Timer | Auto-stop after duration |
| Queue Management | Episode queue/playlist |

---

## Go Packages

> See [00_SOURCE_OF_TRUTH.md](../../00_SOURCE_OF_TRUTH.md#go-dependencies-media-processing--planned) for package versions.

Key packages used:
- **gofeed** - RSS feed parsing
- **gorilla/feeds** - Feed generation (for OPML export)

---

## Database Schema

```sql
-- Podcasts (shows)
CREATE TABLE podcasts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    feed_url TEXT NOT NULL UNIQUE,

    -- Metadata
    title VARCHAR(500) NOT NULL,
    author VARCHAR(300),
    description TEXT,
    summary TEXT,
    link TEXT,
    language VARCHAR(10),
    copyright TEXT,

    -- Images
    image_url TEXT,
    image_path TEXT, -- Cached locally

    -- Categories/tags
    categories TEXT[],

    -- Feed info
    feed_type VARCHAR(20), -- rss, atom
    last_fetched_at TIMESTAMPTZ,
    last_episode_at TIMESTAMPTZ,
    episode_count INT DEFAULT 0,

    -- Status
    is_explicit BOOLEAN DEFAULT false,
    is_complete BOOLEAN DEFAULT false, -- No more episodes

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Episodes
CREATE TABLE podcast_episodes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    podcast_id UUID REFERENCES podcasts(id) ON DELETE CASCADE,
    guid VARCHAR(500) NOT NULL,

    -- Metadata
    title VARCHAR(500) NOT NULL,
    description TEXT,
    summary TEXT,
    link TEXT,
    image_url TEXT,

    -- Audio
    audio_url TEXT NOT NULL,
    audio_size_bytes BIGINT,
    audio_type VARCHAR(100), -- audio/mpeg
    duration_seconds INT,

    -- Publishing
    published_at TIMESTAMPTZ NOT NULL,
    season_number INT,
    episode_number INT,

    -- Local file
    downloaded BOOLEAN DEFAULT false,
    download_path TEXT,

    -- Flags
    is_explicit BOOLEAN DEFAULT false,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE(podcast_id, guid)
);

-- User subscriptions
CREATE TABLE podcast_subscriptions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    podcast_id UUID REFERENCES podcasts(id) ON DELETE CASCADE,

    -- Settings
    auto_download BOOLEAN DEFAULT true,
    keep_episodes INT, -- NULL = keep all
    notification_enabled BOOLEAN DEFAULT true,

    subscribed_at TIMESTAMPTZ DEFAULT NOW(),

    UNIQUE(user_id, podcast_id)
);

-- Episode playback state
CREATE TABLE podcast_episode_state (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    episode_id UUID REFERENCES podcast_episodes(id) ON DELETE CASCADE,

    -- Playback
    position_seconds INT DEFAULT 0,
    duration_seconds INT,
    played BOOLEAN DEFAULT false,
    play_count INT DEFAULT 0,

    -- Timestamps
    started_at TIMESTAMPTZ,
    finished_at TIMESTAMPTZ,
    last_played_at TIMESTAMPTZ,

    UNIQUE(user_id, episode_id)
);

-- User queue
CREATE TABLE podcast_queue (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    episode_id UUID REFERENCES podcast_episodes(id) ON DELETE CASCADE,
    sort_order INT NOT NULL,
    added_at TIMESTAMPTZ DEFAULT NOW(),

    UNIQUE(user_id, episode_id)
);

-- Indexes
CREATE INDEX idx_podcast_episodes_podcast ON podcast_episodes(podcast_id);
CREATE INDEX idx_podcast_episodes_published ON podcast_episodes(published_at DESC);
CREATE INDEX idx_podcast_subscriptions_user ON podcast_subscriptions(user_id);
CREATE INDEX idx_podcast_episode_state_user ON podcast_episode_state(user_id);
CREATE INDEX idx_podcast_queue_user ON podcast_queue(user_id, sort_order);
```

---

## River Jobs

```go
const (
    JobKindRefreshPodcast     = "podcast.refresh"
    JobKindRefreshAllPodcasts = "podcast.refresh_all"
    JobKindDownloadEpisode    = "podcast.download_episode"
    JobKindCleanupDownloads   = "podcast.cleanup_downloads"
)

type RefreshPodcastArgs struct {
    PodcastID uuid.UUID `json:"podcast_id"`
}

type DownloadEpisodeArgs struct {
    EpisodeID uuid.UUID `json:"episode_id"`
    UserID    uuid.UUID `json:"user_id"`
}
```

---

## Go Implementation

```go
// internal/content/podcasts/

type Service struct {
    podcasts     PodcastRepository
    episodes     EpisodeRepository
    subscriptions SubscriptionRepository
    parser       *gofeed.Parser
    river        *river.Client[pgx.Tx]
}

func (s *Service) Subscribe(ctx context.Context, userID uuid.UUID, feedURL string) (*Podcast, error) {
    // Check if podcast already exists
    podcast, err := s.podcasts.GetByFeedURL(ctx, feedURL)
    if errors.Is(err, ErrNotFound) {
        // Fetch and parse feed
        feed, err := s.parser.ParseURLWithContext(feedURL, ctx)
        if err != nil {
            return nil, fmt.Errorf("parse feed: %w", err)
        }

        // Create podcast
        podcast = &Podcast{
            FeedURL:     feedURL,
            Title:       feed.Title,
            Author:      feed.Author.Name,
            Description: feed.Description,
            ImageURL:    feed.Image.URL,
            Language:    feed.Language,
        }
        podcast, err = s.podcasts.Create(ctx, podcast)
        if err != nil {
            return nil, err
        }

        // Import episodes
        for _, item := range feed.Items {
            episode := &Episode{
                PodcastID:   podcast.ID,
                GUID:        item.GUID,
                Title:       item.Title,
                Description: item.Description,
                AudioURL:    findAudioEnclosure(item),
                PublishedAt: *item.PublishedParsed,
            }
            s.episodes.Create(ctx, episode)
        }
    }

    // Create subscription
    sub := &Subscription{
        UserID:    userID,
        PodcastID: podcast.ID,
    }
    _, err = s.subscriptions.Create(ctx, sub)
    if err != nil {
        return nil, err
    }

    return podcast, nil
}

func (s *Service) RefreshPodcast(ctx context.Context, podcastID uuid.UUID) error {
    podcast, err := s.podcasts.GetByID(ctx, podcastID)
    if err != nil {
        return err
    }

    feed, err := s.parser.ParseURLWithContext(podcast.FeedURL, ctx)
    if err != nil {
        return fmt.Errorf("parse feed: %w", err)
    }

    // Import new episodes only
    for _, item := range feed.Items {
        exists, _ := s.episodes.ExistsByGUID(ctx, podcastID, item.GUID)
        if !exists {
            episode := &Episode{
                PodcastID:   podcastID,
                GUID:        item.GUID,
                Title:       item.Title,
                AudioURL:    findAudioEnclosure(item),
                PublishedAt: *item.PublishedParsed,
            }
            s.episodes.Create(ctx, episode)
        }
    }

    // Update last_fetched_at
    return s.podcasts.UpdateLastFetched(ctx, podcastID)
}

func findAudioEnclosure(item *gofeed.Item) string {
    for _, enc := range item.Enclosures {
        if strings.HasPrefix(enc.Type, "audio/") {
            return enc.URL
        }
    }
    return ""
}
```

---

## API Endpoints

```
# Podcasts
GET  /api/v1/podcasts                       # List all podcasts
GET  /api/v1/podcasts/:id                   # Get podcast
GET  /api/v1/podcasts/:id/episodes          # Get episodes
POST /api/v1/podcasts/search?q=...          # Search podcast directories

# Subscriptions
GET  /api/v1/podcasts/subscriptions         # My subscriptions
POST /api/v1/podcasts/subscribe             # Subscribe to podcast
DELETE /api/v1/podcasts/subscriptions/:id   # Unsubscribe
PUT  /api/v1/podcasts/subscriptions/:id     # Update settings

# Episodes
GET  /api/v1/podcasts/episodes/:id          # Get episode
GET  /api/v1/podcasts/episodes/:id/stream   # Stream audio
POST /api/v1/podcasts/episodes/:id/download # Download episode
DELETE /api/v1/podcasts/episodes/:id/download # Remove download

# Playback state
GET  /api/v1/podcasts/episodes/:id/state    # Get playback state
PUT  /api/v1/podcasts/episodes/:id/state    # Update progress
POST /api/v1/podcasts/episodes/:id/played   # Mark as played

# Queue
GET  /api/v1/podcasts/queue                 # Get queue
POST /api/v1/podcasts/queue                 # Add to queue
DELETE /api/v1/podcasts/queue/:id           # Remove from queue
PUT  /api/v1/podcasts/queue/reorder         # Reorder queue

# OPML
GET  /api/v1/podcasts/export/opml           # Export subscriptions
POST /api/v1/podcasts/import/opml           # Import subscriptions
```

---

## Configuration

```yaml
podcasts:
  enabled: true

  refresh:
    interval: 1h
    on_startup: true

  downloads:
    path: "/data/podcasts"
    auto_download: true
    keep_episodes: 10  # Per subscription

  playback:
    default_speed: 1.0
    skip_silence: false
    skip_intros: false  # If silence detected
```

---

## RBAC Permissions

| Permission | Description |
|------------|-------------|
| `podcasts.view` | View/listen to podcasts |
| `podcasts.subscribe` | Subscribe to podcasts |
| `podcasts.download` | Download episodes |
| `podcasts.manage` | Add/remove podcasts globally |

---

## Implementation Checklist

### Phase 1: Core Infrastructure
- [ ] Create `internal/content/podcasts/` package structure
- [ ] Define `entity.go` with Podcast, Episode, Subscription structs
- [ ] Create `repository.go` interface definition
- [ ] Implement `repository_pg.go` with sqlc queries
- [ ] Add fx module wiring in `module.go`

### Phase 2: Database
- [ ] Create migration `000XXX_create_podcasts_schema.up.sql`
- [ ] Create `podcasts` table
- [ ] Create `podcast_episodes` table
- [ ] Create `podcast_subscriptions` table
- [ ] Create `podcast_episode_state` table
- [ ] Create `podcast_queue` table
- [ ] Add indexes (feed_url, published_at, user subscriptions)
- [ ] Write sqlc queries in `queries/podcasts/`

### Phase 3: Feed Parsing
- [ ] Implement RSS feed parser (gofeed)
- [ ] Extract podcast metadata (title, author, image)
- [ ] Extract episode metadata (audio URL, duration)
- [ ] Handle feed variations (RSS 2.0, Atom, iTunes extensions)
- [ ] Implement OPML import/export

### Phase 4: Service Layer
- [ ] Implement `service.go` with otter caching
- [ ] Add Podcast operations (Get, List, Search)
- [ ] Add Subscription operations (Subscribe, Unsubscribe, Update)
- [ ] Add Episode operations (Get, List, UpdateState)
- [ ] Add Queue operations (Add, Remove, Reorder)
- [ ] Implement cache invalidation

### Phase 5: Download Management
- [ ] Implement episode download
- [ ] Add download progress tracking
- [ ] Implement auto-download for subscriptions
- [ ] Add download cleanup (keep N episodes)

### Phase 6: Background Jobs
- [ ] Create River job definitions in `jobs.go`
- [ ] Implement `RefreshPodcastJob`
- [ ] Implement `RefreshAllPodcastsJob`
- [ ] Implement `DownloadEpisodeJob`
- [ ] Implement `CleanupDownloadsJob`

### Phase 7: API Integration
- [ ] Define OpenAPI endpoints for podcasts
- [ ] Generate ogen handlers
- [ ] Wire handlers to service layer
- [ ] Add audio streaming endpoint
- [ ] Add OPML import/export endpoints
- [ ] Add authentication/authorization checks

---


## Related

- [Music Module](../audio/MUSIC_MODULE.md) - Audio playback patterns
- [Scrobbling](../shared/SCROBBLING.md) - Playback tracking
- [News System](../shared/NEWS_SYSTEM.md) - Similar RSS patterns
