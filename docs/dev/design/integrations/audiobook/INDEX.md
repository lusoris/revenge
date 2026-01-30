# Audiobook & Podcast Module

> Native audiobook and podcast management

---

## Overview

Revenge provides **native audiobook and podcast management** - all functionality is implemented directly without external dependencies:

- **Audiobook Library Management** - Multi-library support, file scanning, metadata
- **Podcast Management** - RSS feed parsing, episode downloads, subscriptions
- **Playback Features** - Progress tracking, bookmarks, sleep timer, chapters
- **Metadata Providers** - Audible/Audnexus, OpenLibrary, Goodreads, etc.
- **Multi-user Support** - Per-user progress, preferences, permissions

---

## Architecture Decision

**Why Native Implementation?**

1. **Simpler Deployment** - No external service dependencies
2. **Unified Experience** - Single UI/API for all media types
3. **Better Integration** - Direct access to all Revenge features (scrobbling, analytics, etc.)
4. **Core Capabilities** - River jobs, caching, auth, API infrastructure already exists
5. **Full Control** - No dependency on third-party release cycles

---

## Features

### Library Management

- Multiple libraries with custom paths
- Auto-detect file changes (fsnotify watcher)
- Flexible directory structures (author/series/disc folders)
- Automated backups + scheduling
- Bulk upload via API/drag-drop

### Audiobook Features

- ID3/M4B metadata extraction (taglib)
- Chapter support (embedded + Audnexus lookup)
- Audio file merging to M4B (FFmpeg)
- Embed metadata/cover in audio files
- Multi-file/multi-disc support
- Track reordering

### Podcast Features

- RSS feed parsing (gofeed library)
- Auto-download new episodes (River scheduled jobs)
- Episode queue management per user
- OPML import/export
- Open RSS feed generation for external apps
- Subscription management

### Playback & Progress

- Stream all audio formats (native + FFmpeg fallback)
- Per-user progress tracking with sync
- Cross-device synchronization
- Bookmarks with notes
- Sleep timer (time-based + end-of-chapter)
- Variable playback speed (0.5x - 3.0x)
- Smart rewind after pause (configurable)
- Continue listening / resume functionality
- Chromecast support

### Metadata Providers

| Provider | Content | Status |
| -------- | ------- | ------ |
| Audnexus | Audiobook chapters, metadata | Planned |
| Audible | Audiobook metadata, covers | Planned |
| OpenLibrary | Book metadata | Planned |
| Goodreads | Ratings, reviews | Planned |
| Google Books | Book metadata | Planned |
| iTunes Podcasts | Podcast search | Planned |
| Podcast Index | Podcast search, chapters | Planned |

### User Features

- Multi-user with RBAC permissions
- Per-user progress, bookmarks, ratings
- Listening session tracking
- OAuth2/OIDC SSO support
- Personal podcast queues

---

## Implementation Status

| Component | Status | Location |
| --------- | ------ | -------- |
| Audiobook Entity | ✅ Done | `internal/content/audiobook/entity.go` |
| Audiobook Repository | ✅ Done | `internal/content/audiobook/repository.go` |
| Audiobook Service | ✅ Done | `internal/content/audiobook/service.go` |
| Audiobook Jobs | ✅ Done | `internal/content/audiobook/jobs.go` |
| Podcast Entity | ✅ Done | `internal/content/podcast/entity.go` |
| Podcast Repository | ✅ Done | `internal/content/podcast/repository.go` |
| Podcast Service | ✅ Done | `internal/content/podcast/service.go` |
| RSS Parser | 🟡 Stub | `internal/content/podcast/rss_parser.go` |
| Audnexus Provider | 🔴 TODO | `internal/service/metadata/audnexus/` |
| Sleep Timer | 🔴 TODO | Playback service |
| Chapter Editor | 🔴 TODO | API + Frontend |

---

## Technical Details

### RSS Feed Parsing

```go
// Using gofeed library
import "github.com/mmcdole/gofeed"

parser := gofeed.NewParser()
feed, err := parser.ParseURL(feedURL)
```

### Episode Downloads

River job with progress tracking:

```go
type DownloadEpisodeArgs struct {
    EpisodeID uuid.UUID `json:"episode_id"`
}

func (w *DownloadWorker) Work(ctx context.Context, job *river.Job[DownloadEpisodeArgs]) error {
    // Download with progress callback
    // Update episode.downloaded = true
    // Store at episode.local_path
}
```

### Chapter Extraction

```go
// From M4B files (embedded)
// From Audnexus API (ASIN lookup)
// From podcast RSS (chapters namespace)
```

---

## Related Documentation

- [Book Metadata Providers](../metadata/books/INDEX.md)
- [Audible/Audnexus](../metadata/books/AUDIBLE.md)
- [Audio Streaming](../../technical/AUDIO_STREAMING.md)
- [Scrobbling](../../features/shared/SCROBBLING.md)
