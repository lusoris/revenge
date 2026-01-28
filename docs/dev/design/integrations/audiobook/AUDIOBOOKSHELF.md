# Audiobookshelf Integration

> Self-hosted audiobook and podcast server

**Status**: 🟡 PLANNED
**Priority**: 🟡 MEDIUM (Phase 6 - Audiobook/Podcast/Book Modules)
**Type**: REST API client + metadata provider

---

## Overview

Audiobookshelf is a self-hosted audiobook and podcast server that serves as both a metadata provider and primary source for audiobook/podcast/ebook content. Revenge integrates with Audiobookshelf as:
- **Primary metadata source** for audiobooks, podcasts, ebooks
- **Library sync** - import libraries from existing Audiobookshelf instances
- **Progress sync** - bidirectional playback progress
- **Optional playback** - use Audiobookshelf as playback backend

**Integration Points**:
- **REST API**: Query libraries, books, podcasts
- **Socket.io**: Real-time progress updates
- **Audio streaming**: Direct or proxied playback
- **Metadata**: Cover art, descriptions, chapters

---

## Developer Resources

- 📚 **API Docs**: https://api.audiobookshelf.org/
- 🔗 **GitHub**: https://github.com/advplyr/audiobookshelf
- 🔗 **OpenAPI Spec**: https://github.com/advplyr/audiobookshelf/blob/master/docs/openapi.json
- 🔗 **Socket Events**: https://github.com/advplyr/audiobookshelf/wiki/Socket-Events

---

## API Details

**Base URL**: `https://audiobookshelf.example.com/api`
**Authentication**: API key via header `Authorization: Bearer {token}`
**Rate Limits**: None (self-hosted)
**Websocket**: `wss://audiobookshelf.example.com/socket.io/`

### Key Endpoints

| Endpoint | Method | Purpose |
|----------|--------|---------|
| `/libraries` | GET | List all libraries |
| `/libraries/{id}/items` | GET | List items in library |
| `/items/{id}` | GET | Get item details (book/podcast) |
| `/items/{id}/cover` | GET | Get cover image |
| `/items/{id}/play` | POST | Start playback session |
| `/me/progress/{id}` | GET | Get user's progress |
| `/me/progress/{id}` | PATCH | Update progress |
| `/podcasts/{id}/episodes` | GET | List podcast episodes |
| `/search/library/{id}` | GET | Search within library |

### Authentication

```bash
# Login to get token
POST /login
{
  "username": "user",
  "password": "pass"
}
# Response: { "user": {...}, "token": "..." }

# Use token in subsequent requests
Authorization: Bearer {token}
```

---

## Data Mapping

### Audiobookshelf → Revenge

**Audiobook Mapping**:
| Audiobookshelf Field | Revenge Field | Notes |
|---------------------|---------------|-------|
| `id` | `audiobookshelf_id` | ABS identifier |
| `media.metadata.title` | `title` | Book title |
| `media.metadata.subtitle` | `subtitle` | Subtitle |
| `media.metadata.authors[].name` | `authors[]` | Author names |
| `media.metadata.narrators[]` | `narrators[]` | Narrator names |
| `media.metadata.series[].name` | `series` | Series name |
| `media.metadata.series[].sequence` | `series_position` | Position in series |
| `media.metadata.description` | `overview` | Description |
| `media.metadata.publishedYear` | `release_year` | Publication year |
| `media.metadata.genres[]` | `genres[]` | Genre list |
| `media.metadata.language` | `language` | Language code |
| `media.duration` | `duration` | Total duration (seconds) |
| `media.chapters[]` | `chapters[]` | Chapter markers |
| `media.audioFiles[]` | `audio_files[]` | Audio track info |

**Podcast Mapping**:
| Audiobookshelf Field | Revenge Field | Notes |
|---------------------|---------------|-------|
| `id` | `audiobookshelf_id` | ABS identifier |
| `media.metadata.title` | `title` | Podcast title |
| `media.metadata.author` | `author` | Podcast author |
| `media.metadata.description` | `overview` | Description |
| `media.metadata.feedUrl` | `rss_feed_url` | RSS feed |
| `media.metadata.imageUrl` | `poster_url` | Podcast art |
| `media.episodes[]` | `episodes[]` | Episode list |

**Progress Mapping**:
| Audiobookshelf Field | Revenge Field | Notes |
|---------------------|---------------|-------|
| `currentTime` | `position_seconds` | Current position |
| `progress` | `progress_percent` | 0.0 - 1.0 |
| `isFinished` | `completed` | Completion status |
| `lastUpdate` | `updated_at` | Last sync time |

---

## Implementation Checklist

- [ ] **API Client** (`internal/service/audiobook/provider_audiobookshelf.go`)
  - [ ] HTTP client with auth token
  - [ ] Library listing/sync
  - [ ] Book/podcast metadata fetching
  - [ ] Search functionality
  - [ ] Error handling & retries

- [ ] **Library Sync** (`internal/service/sync/audiobookshelf_library.go`)
  - [ ] Initial full library import
  - [ ] Incremental sync (new/updated items)
  - [ ] Handle deletions
  - [ ] Cover art caching
  - [ ] Chapter extraction

- [ ] **Progress Sync** (`internal/service/sync/audiobookshelf_progress.go`)
  - [ ] Fetch user progress from ABS
  - [ ] Push Revenge progress to ABS
  - [ ] Conflict resolution (newer wins)
  - [ ] Real-time sync via Socket.io (optional)

- [ ] **Playback Integration** (`internal/service/playback/audiobookshelf.go`)
  - [ ] Direct stream URL generation
  - [ ] Proxied playback (for transcoding)
  - [ ] Chapter navigation
  - [ ] Sleep timer support

---

## Configuration

```yaml
# configs/config.yaml
integrations:
  audiobookshelf:
    enabled: true
    base_url: "https://audiobookshelf.example.com"
    api_token: "${REVENGE_AUDIOBOOKSHELF_TOKEN}"

    sync:
      enabled: true
      interval: "1h"           # How often to sync library
      libraries:               # Which libraries to import
        - all                  # Or specific library IDs
      direction: "import"      # import, export, bidirectional

    progress_sync:
      enabled: true
      realtime: false          # Use Socket.io for real-time
      interval: "5m"           # Polling interval if not realtime
      direction: "bidirectional"

    playback:
      mode: "direct"           # direct, proxy
      transcode: false         # Let ABS handle transcoding

    content_types:
      audiobooks: true
      podcasts: true
      ebooks: false            # If using ABS for ebooks
```

---

## Database Schema

```sql
-- Audiobookshelf library mapping
CREATE TABLE audiobookshelf_libraries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    abs_library_id VARCHAR(100) NOT NULL UNIQUE,
    revenge_library_id UUID REFERENCES libraries(id),
    library_type VARCHAR(20) NOT NULL,  -- audiobook, podcast, ebook
    name VARCHAR(255) NOT NULL,
    last_sync_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Item mapping (audiobooks, podcasts)
CREATE TABLE audiobookshelf_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    abs_item_id VARCHAR(100) NOT NULL UNIQUE,
    abs_library_id VARCHAR(100) NOT NULL,
    revenge_item_id UUID,  -- FK to audiobooks/podcasts table
    revenge_item_type VARCHAR(20) NOT NULL,  -- audiobook, podcast
    abs_updated_at TIMESTAMPTZ,
    last_sync_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_abs_items_library ON audiobookshelf_items(abs_library_id);

-- Progress mapping
CREATE TABLE audiobookshelf_progress (
    user_id UUID NOT NULL REFERENCES users(id),
    abs_item_id VARCHAR(100) NOT NULL,
    abs_episode_id VARCHAR(100),  -- For podcasts
    position_seconds REAL NOT NULL DEFAULT 0,
    progress_percent REAL NOT NULL DEFAULT 0,
    completed BOOLEAN NOT NULL DEFAULT FALSE,
    abs_updated_at TIMESTAMPTZ,
    revenge_updated_at TIMESTAMPTZ,
    last_sync_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, abs_item_id, COALESCE(abs_episode_id, ''))
);
```

---

## Sync Strategies

### Initial Import

1. Fetch all libraries from Audiobookshelf
2. For each library, fetch all items (paginated)
3. Create Revenge audiobook/podcast entries
4. Download and cache cover art
5. Extract chapter information
6. Record mapping in `audiobookshelf_items`

### Incremental Sync

1. Fetch items updated since `last_sync_at`
2. Update existing Revenge entries
3. Create new entries for new items
4. Handle deleted items (mark unavailable)

### Progress Sync

**Bidirectional flow**:
```
User plays in Revenge:
1. Update Revenge progress
2. Push to Audiobookshelf API
3. Record sync timestamp

User plays in Audiobookshelf:
1. Poll/receive progress update
2. Compare timestamps
3. Update Revenge if ABS is newer
```

---

## Socket.io Integration (Optional)

For real-time progress sync:

```go
// Connect to Audiobookshelf Socket.io
socket := socketio.Connect("wss://abs.example.com/socket.io/")

// Authenticate
socket.Emit("auth", token)

// Listen for progress updates
socket.On("user_item_progress_updated", func(data ProgressUpdate) {
    // Update Revenge progress
})

// Send progress updates
socket.Emit("set_progress", ProgressPayload{
    LibraryItemId: "...",
    Progress: 0.5,
    CurrentTime: 3600,
})
```

---

## Playback Modes

### Direct Mode

Revenge redirects to Audiobookshelf stream URL:
```
GET /api/v1/audiobooks/{id}/stream
→ 302 Redirect to https://abs.example.com/api/items/{abs_id}/play
```

### Proxy Mode

Revenge proxies the stream (for custom auth, logging):
```
GET /api/v1/audiobooks/{id}/stream
→ Revenge fetches from ABS and streams to client
```

---

## Error Handling

| Error | Cause | Solution |
|-------|-------|----------|
| 401 Unauthorized | Invalid/expired token | Re-authenticate, get new token |
| 404 Not Found | Item deleted from ABS | Mark as unavailable in Revenge |
| 500 Server Error | ABS issue | Retry with backoff |
| Connection refused | ABS offline | Queue for retry, show offline status |

---

## Related Documentation

- [Audiobook Module](../../features/LIBRARY_TYPES.md)
- [Podcast Module](../../features/LIBRARY_TYPES.md)
- [Scrobbling - Audiobooks](../scrobbling/INDEX.md)
