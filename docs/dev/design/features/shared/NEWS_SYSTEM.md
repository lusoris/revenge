# News System

> External news aggregation and internal announcements

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

| Package | Purpose | URL |
|---------|---------|-----|
| **gofeed** | RSS/Atom/JSON feed parsing | github.com/mmcdole/gofeed |
| **gorilla/feeds** | Feed generation | github.com/gorilla/feeds |

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
| **gorilla/websocket** | WebSocket connections |
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

## Related Documentation

- [Notification System](NOTIFICATIONS.md)
- [RBAC Permissions](RBAC_CASBIN.md)
- [River Job Queue Patterns](../../SOURCE_OF_TRUTH.md#river-job-queue-patterns)
