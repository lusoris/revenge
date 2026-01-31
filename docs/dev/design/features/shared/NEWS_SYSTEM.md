# News System

> External news aggregation and internal announcements


<!-- TOC-START -->

## Table of Contents

- [Status](#status)
- [Overview](#overview)
- [Architecture](#architecture)
- [External News (RSS Aggregation)](#external-news-rss-aggregation)
  - [Go Packages](#go-packages)
  - [Feed Categories](#feed-categories)
    - [Standard Content](#standard-content)
    - [Adult Content (Isolated)](#adult-content-isolated)
  - [Database Schema](#database-schema)
  - [River Jobs](#river-jobs)
  - [Go Implementation Example](#go-implementation-example)
- [Internal News (Announcements)](#internal-news-announcements)
  - [Announcement Types](#announcement-types)
  - [Database Schema](#database-schema)
  - [RBAC Permissions](#rbac-permissions)
- [Real-Time Notifications](#real-time-notifications)
  - [Go Packages](#go-packages)
  - [SSE for Announcements](#sse-for-announcements)
- [Configuration](#configuration)
- [API Endpoints](#api-endpoints)
  - [External News](#external-news)
  - [Internal Announcements](#internal-announcements)
- [UI/UX Integration](#uiux-integration)
  - [News Dashboard Widget](#news-dashboard-widget)
  - [Announcement Banner](#announcement-banner)
- [Implementation Checklist](#implementation-checklist)
  - [Phase 1: Core Infrastructure](#phase-1-core-infrastructure)
  - [Phase 2: Database](#phase-2-database)
  - [Phase 3: Service Layer](#phase-3-service-layer)
  - [Phase 4: Background Jobs](#phase-4-background-jobs)
  - [Phase 5: Real-Time Notifications](#phase-5-real-time-notifications)
  - [Phase 6: API Integration](#phase-6-api-integration)
- [Sources & Cross-References](#sources-cross-references)
  - [Cross-Reference Indexes](#cross-reference-indexes)
  - [Referenced Sources](#referenced-sources)
- [Related Design Docs](#related-design-docs)
  - [In This Section](#in-this-section)
  - [Related Topics](#related-topics)
  - [Indexes](#indexes)
- [Related](#related)

<!-- TOC-END -->

## Status

| Dimension | Status | Notes |
|-----------|--------|-------|
| Design | ✅ | Full design with RSS aggregation, announcements, SSE |
| Sources | ✅ | gofeed, gorilla/feeds, r3labs/sse documented |
| Instructions | ✅ | Implementation checklist added |
| Code | 🔴 |  |
| Linting | 🔴 |  |
| Unit Testing | 🔴 |  |
| Integration Testing | 🔴 |  |
---

## Overview

The News System provides two distinct functionalities:
1. **External News**: RSS/Atom feed aggregation for media industry news
2. **Internal News**: Announcements from admins/mods/devs to users

---

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                      News System                            │
├─────────────────────────┬───────────────────────────────────┤
│    External News        │       Internal News               │
│    (RSS Aggregation)    │       (Announcements)             │
├─────────────────────────┼───────────────────────────────────┤
│ • Media industry news   │ • System announcements            │
│ • Release calendars     │ • Feature updates                 │
│ • Review aggregation    │ • Maintenance notices             │
│ • Streaming updates     │ • Community updates               │
│ • Adult news (isolated) │ • Admin broadcasts                │
└─────────────────────────┴───────────────────────────────────┘
```

---

## External News (RSS Aggregation)

### Go Packages

> Package versions: [00_SOURCE_OF_TRUTH.md](../../00_SOURCE_OF_TRUTH.md#go-dependencies-core)

| Package | Purpose |
|---------|---------|
| gofeed | RSS/Atom/JSON feed parsing |
| gorilla/feeds | Feed generation |

### Feed Categories

#### Standard Content

| Category | Example Sources |
|----------|-----------------|
| Movies | Collider, Screen Rant, IGN Movies |
| TV Shows | TV Line, Deadline, TVGuide |
| Music | Pitchfork, Billboard, NME |
| Books | Publishers Weekly, BookRiot |
| Gaming | Kotaku, IGN, GameSpot |

#### Adult Content (Isolated)

| Category | Example Sources |
|----------|-----------------|
| Industry News | AVN, XBIZ |
| Release Calendar | Studio RSS feeds |

> **Isolation**: Adult news stored in `c.news_feeds` and `c.news_articles`

### Database Schema

```sql
-- Standard news
CREATE TABLE news.feeds (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    url TEXT NOT NULL UNIQUE,
    category VARCHAR(100),
    enabled BOOLEAN DEFAULT true,
    last_fetched_at TIMESTAMPTZ,
    fetch_interval_minutes INT DEFAULT 60,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE news.articles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    feed_id UUID REFERENCES news.feeds(id),
    external_id TEXT,
    title TEXT NOT NULL,
    description TEXT,
    content TEXT,
    url TEXT NOT NULL,
    image_url TEXT,
    published_at TIMESTAMPTZ,
    fetched_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(feed_id, external_id)
);

-- Adult news (isolated)
CREATE TABLE c.news_feeds (
    -- Same structure as news.feeds
);

CREATE TABLE c.news_articles (
    -- Same structure as news.articles
);
```

### River Jobs

```go
// Job kinds
const (
    JobKindFetchFeed     = "news.fetch_feed"
    JobKindFetchAllFeeds = "news.fetch_all_feeds"
    JobKindCleanupOld    = "news.cleanup_old_articles"
)

// FetchFeedArgs for individual feed fetching
type FetchFeedArgs struct {
    FeedID uuid.UUID `json:"feed_id"`
}

// Scheduled via River periodic jobs
// Every hour: fetch all enabled feeds
// Daily: cleanup articles older than retention period
```

### Go Implementation Example

```go
import "github.com/mmcdole/gofeed"

type FeedService struct {
    parser *gofeed.Parser
    repo   FeedRepository
}

func (s *FeedService) FetchFeed(ctx context.Context, feedURL string) ([]Article, error) {
    feed, err := s.parser.ParseURLWithContext(feedURL, ctx)
    if err != nil {
        return nil, fmt.Errorf("parse feed: %w", err)
    }

    articles := make([]Article, 0, len(feed.Items))
    for _, item := range feed.Items {
        articles = append(articles, Article{
            ExternalID:  item.GUID,
            Title:       item.Title,
            Description: item.Description,
            Content:     item.Content,
            URL:         item.Link,
            ImageURL:    extractImage(item),
            PublishedAt: item.PublishedParsed,
        })
    }
    return articles, nil
}
```

---

## Internal News (Announcements)

### Announcement Types

| Type | Audience | Example |
|------|----------|---------|
| `system` | All users | Maintenance windows |
| `feature` | All users | New feature releases |
| `admin` | Admins only | Admin-specific updates |
| `mod` | Mods + Admins | Moderation updates |
| `dev` | Developers | API changes |
| `community` | All users | Community events |

### Database Schema

```sql
CREATE TABLE announcements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(500) NOT NULL,
    content TEXT NOT NULL,
    content_html TEXT, -- Rendered markdown
    type VARCHAR(50) NOT NULL,
    priority INT DEFAULT 0, -- Higher = more important
    target_roles TEXT[], -- NULL = all users
    pinned BOOLEAN DEFAULT false,
    published_at TIMESTAMPTZ,
    expires_at TIMESTAMPTZ,
    author_id UUID REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE announcement_reads (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    announcement_id UUID REFERENCES announcements(id),
    user_id UUID REFERENCES users(id),
    read_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(announcement_id, user_id)
);

-- Adult announcements (isolated)
CREATE TABLE c.announcements (
    -- Same structure
);
```

### RBAC Permissions

| Permission | Description |
|------------|-------------|
| `news.announcements.create` | Create announcements |
| `news.announcements.edit` | Edit announcements |
| `news.announcements.delete` | Delete announcements |
| `news.announcements.pin` | Pin announcements |
| `news.feeds.manage` | Manage RSS feeds |

---

## Real-Time Notifications

### Go Packages

| Package | Purpose |
|---------|---------|
| **gobwas/ws** | WebSocket connections |
| **r3labs/sse/v2** | Server-Sent Events |

### SSE for Announcements

```go
import "github.com/r3labs/sse/v2"

type NotificationService struct {
    server *sse.Server
}

func (s *NotificationService) BroadcastAnnouncement(announcement *Announcement) {
    data, _ := json.Marshal(announcement)
    s.server.Publish("announcements", &sse.Event{
        Data: data,
    })
}

// Client subscription
// GET /api/v1/news/stream
```

---

## Configuration

```yaml
news:
  external:
    enabled: true
    fetch_interval: 60m
    retention_days: 30
    max_articles_per_feed: 100

  internal:
    enabled: true
    default_expiry_days: 30

  # Adult news isolation
  adult:
    enabled: false
    schema: c
```

---

## API Endpoints

### External News

```
GET  /api/v1/news/feeds              # List feeds
POST /api/v1/news/feeds              # Add feed (admin)
GET  /api/v1/news/articles           # List articles
GET  /api/v1/news/articles/:id       # Get article

# Adult (isolated)
GET  /api/v1/legacy/news/feeds
GET  /api/v1/legacy/news/articles
```

### Internal Announcements

```
GET  /api/v1/announcements           # List announcements
POST /api/v1/announcements           # Create (admin/mod)
PUT  /api/v1/announcements/:id       # Update
POST /api/v1/announcements/:id/read  # Mark as read
GET  /api/v1/announcements/unread    # Get unread count

# Real-time
GET  /api/v1/news/stream             # SSE endpoint
```

---

## UI/UX Integration

### News Dashboard Widget

- Collapsible news feed
- Category filters
- Read/unread indicators
- "Mark all as read"

### Announcement Banner

- Pinned announcements at top
- Dismissible notifications
- Priority-based styling (warning, info, success)
- Unread badge in navigation

---

## Implementation Checklist

**Location**: `internal/service/news/`

### Phase 1: Core Infrastructure
- [ ] Create package structure `internal/service/news/`
- [ ] Define entities: `Feed`, `Article`, `Announcement`, `AnnouncementRead`
- [ ] Create repository interface `NewsRepository`
- [ ] Implement fx module `news.Module`
- [ ] Add configuration struct for news settings

### Phase 2: Database
- [ ] Create migration `xxx_news_system.up.sql`
- [ ] Create `news.feeds` table with feed metadata
- [ ] Create `news.articles` table with article content
- [ ] Create `announcements` table with types and targeting
- [ ] Create `announcement_reads` table for read tracking
- [ ] Create `c.news_feeds` and `c.news_articles` for adult content isolation
- [ ] Create `c.announcements` for adult content announcements
- [ ] Add indexes for feed fetching and article queries
- [ ] Generate sqlc queries for CRUD operations

### Phase 3: Service Layer
- [ ] Implement `FeedService` with gofeed parser integration
- [ ] Implement RSS/Atom/JSON feed parsing with `github.com/mmcdole/gofeed`
- [ ] Implement `AnnouncementService` for internal news
- [ ] Add caching layer for frequently accessed articles
- [ ] Implement content sanitization for feed content
- [ ] Add markdown rendering for announcements (`content_html`)
- [ ] Implement role-based announcement targeting

### Phase 4: Background Jobs
- [ ] Create `FetchFeedWorker` for individual feed fetching
- [ ] Create `FetchAllFeedsWorker` for periodic batch fetching
- [ ] Create `CleanupOldArticlesWorker` for retention policy
- [ ] Configure River periodic jobs (hourly fetch, daily cleanup)
- [ ] Implement retry logic for failed feed fetches
- [ ] Add feed health monitoring (track consecutive failures)

### Phase 5: Real-Time Notifications
- [ ] Implement SSE endpoint using `github.com/r3labs/sse/v2`
- [ ] Create `NotificationService` for announcement broadcasts
- [ ] Add SSE channel management for authenticated users
- [ ] Implement announcement push on publish

### Phase 6: API Integration
- [ ] Define OpenAPI spec for news endpoints
- [ ] Generate ogen handlers for feed management
- [ ] Implement `GET /api/v1/news/feeds` - list feeds
- [ ] Implement `POST /api/v1/news/feeds` - add feed (admin)
- [ ] Implement `GET /api/v1/news/articles` - list articles
- [ ] Implement `GET /api/v1/news/articles/:id` - get article
- [ ] Implement `GET /api/v1/announcements` - list announcements
- [ ] Implement `POST /api/v1/announcements` - create (admin/mod)
- [ ] Implement `PUT /api/v1/announcements/:id` - update
- [ ] Implement `POST /api/v1/announcements/:id/read` - mark as read
- [ ] Implement `GET /api/v1/announcements/unread` - unread count
- [ ] Implement `GET /api/v1/news/stream` - SSE endpoint
- [ ] Add RBAC permission checks for admin operations
- [ ] Implement adult content endpoints under `/api/v1/legacy/news/`

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
| [Casbin](https://pkg.go.dev/github.com/casbin/casbin/v2) | [Local](../../../sources/security/casbin.md) |
| [River Job Queue](https://pkg.go.dev/github.com/riverqueue/river) | [Local](../../../sources/tooling/river.md) |
| [Uber fx](https://pkg.go.dev/go.uber.org/fx) | [Local](../../../sources/tooling/fx.md) |
| [gofeed GitHub README](https://github.com/mmcdole/gofeed) | [Local](../../../sources/tooling/gofeed-guide.md) |
| [mmcdole/gofeed](https://pkg.go.dev/github.com/mmcdole/gofeed) | [Local](../../../sources/tooling/gofeed.md) |
| [ogen OpenAPI Generator](https://pkg.go.dev/github.com/ogen-go/ogen) | [Local](../../../sources/tooling/ogen.md) |
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
- [Revenge - NSFW Toggle](NSFW_TOGGLE.md)
- [Dynamic RBAC with Casbin](RBAC_CASBIN.md)

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

- [Notification System](NOTIFICATIONS.md)
- [RBAC Permissions](RBAC_CASBIN.md)
- [River Job Queue Patterns](../../00_SOURCE_OF_TRUTH.md#river-job-queue-patterns)
