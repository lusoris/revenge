# Audiobook & Podcast Module

> Native audiobook and podcast management

**Source of Truth**: [00_SOURCE_OF_TRUTH.md](../../00_SOURCE_OF_TRUTH.md)

<!-- TOC-START -->

## Table of Contents

- [Status](#status)
- [Overview](#overview)
- [Architecture Decision](#architecture-decision)
- [Features](#features)
  - [Library Management](#library-management)
  - [Audiobook Features](#audiobook-features)
  - [Podcast Features](#podcast-features)
  - [Playback & Progress](#playback--progress)
  - [Metadata Providers](#metadata-providers)
  - [User Features](#user-features)
- [Implementation Status](#implementation-status)
- [Technical Details](#technical-details)
  - [RSS Feed Parsing](#rss-feed-parsing)
  - [Episode Downloads](#episode-downloads)
  - [Chapter Extraction](#chapter-extraction)
- [Sources & Cross-References](#sources--cross-references)
  - [Cross-Reference Indexes](#cross-reference-indexes)
  - [Referenced Sources](#referenced-sources)
- [Related Design Docs](#related-design-docs)
  - [Related Topics](#related-topics)
  - [Indexes](#indexes)
- [Related Documentation](#related-documentation)

<!-- TOC-END -->

## Status

| Dimension | Status | Notes |
|-----------|--------|-------|
| Design | 🟡 | Native implementation approach documented |
| Sources | 🟡 | gofeed, Audnexus, Podcast Index sources available |
| Instructions | 🔴 |  |
| Code | 🔴 |  |
| Linting | 🔴 |  |
| Unit Testing | 🔴 |  |
| Integration Testing | 🔴 |  |

---

## Overview

Revenge provides **native audiobook and podcast management** - all functionality is implemented directly without external dependencies:

- **Audiobook Library Management** - Multi-library support, file scanning, metadata
- **Podcast Management** - RSS feed parsing, episode downloads, subscriptions
- **Playback Features** - Progress tracking, bookmarks, sleep timer, chapters
- **Metadata Providers** - Audible/Audnexus, OpenLibrary, Goodreads, Podcast Index
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
| [Audnexus](../metadata/books/AUDIBLE.md) | Audiobook chapters, metadata, narrators | Planned |
| [OpenLibrary](../metadata/books/OPENLIBRARY.md) | Book metadata | Planned |
| [Goodreads](../metadata/books/GOODREADS.md) | Ratings, reviews | Planned |
| [Hardcover](../metadata/books/HARDCOVER.md) | Book metadata, lists | Planned |
| [Podcast Index](https://podcastindex.org) | Podcast search, chapters | Planned |

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
| Audiobook Entity | 🔴 TODO | `internal/content/audiobook/entity.go` |
| Audiobook Repository | 🔴 TODO | `internal/content/audiobook/repository.go` |
| Audiobook Service | 🔴 TODO | `internal/content/audiobook/service.go` |
| Audiobook Jobs | 🔴 TODO | `internal/content/audiobook/jobs.go` |
| Podcast Entity | 🔴 TODO | `internal/content/podcast/entity.go` |
| Podcast Repository | 🔴 TODO | `internal/content/podcast/repository.go` |
| Podcast Service | 🔴 TODO | `internal/content/podcast/service.go` |
| RSS Parser | 🔴 TODO | `internal/content/podcast/rss_parser.go` |
| Audnexus Provider | 🔴 TODO | `internal/service/metadata/audnexus/` |
| Podcast Index Provider | 🔴 TODO | `internal/service/metadata/podcastindex/` |
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

<!-- SOURCE-BREADCRUMBS-START -->

## Sources & Cross-References

> Auto-generated section linking to external documentation sources

### Cross-Reference Indexes

- [All Sources Index](../../../sources/SOURCES_INDEX.md) - Complete list of external documentation
- [Design ↔ Sources Map](../../../sources/DESIGN_CROSSREF.md) - Which docs reference which sources

### Referenced Sources

| Source | Documentation |
|--------|---------------|
| [gofeed RSS Parser](https://pkg.go.dev/github.com/mmcdole/gofeed) | [Local](../../../sources/tooling/gofeed.md) |
| [gofeed Guide](../../../sources/tooling/gofeed-guide.md) | Usage patterns |
| [River Background Jobs](https://pkg.go.dev/github.com/riverqueue/river) | [Local](../../../sources/tooling/river.md) |
| [fsnotify File Watcher](https://pkg.go.dev/github.com/fsnotify/fsnotify) | [Local](../../../sources/tooling/fsnotify.md) |
| [go-astiav FFmpeg](https://github.com/asticode/go-astiav) | [Local](../../../sources/media/go-astiav.md) |
| [dhowden/tag Audio Tags](https://pkg.go.dev/github.com/dhowden/tag) | [Local](../../../sources/media/dhowden-tag.md) |
| [bogem/id3v2 ID3 Tags](https://pkg.go.dev/github.com/bogem/id3v2/v2) | [Local](../../../sources/media/bogem-id3v2.md) |
| [Podcast Index API](https://podcastindex-org.github.io/docs-api/) | [Local](../../../sources/apis/podcastindex.md) |

<!-- SOURCE-BREADCRUMBS-END -->

<!-- DESIGN-BREADCRUMBS-START -->

## Related Design Docs

> Auto-generated cross-references to related design documentation

**Category**: [Integrations](../INDEX.md) > Audiobook

### Related Topics

- [Audiobook Module Feature](../../features/audiobook/AUDIOBOOK_MODULE.md) _Features_
- [Podcasts Feature](../../features/podcasts/PODCASTS.md) _Features_
- [Audio Streaming](../../technical/AUDIO_STREAMING.md) _Technical_
- [Scrobbling](../../features/shared/SCROBBLING.md) _Shared Features_

### Indexes

- [Design Index](../../DESIGN_INDEX.md) - All design docs by category/topic
- [Source of Truth](../../00_SOURCE_OF_TRUTH.md) - Package versions and status

<!-- DESIGN-BREADCRUMBS-END -->

## Related Documentation

- [Book Metadata Providers](../metadata/books/INDEX.md)
- [Audible/Audnexus](../metadata/books/AUDIBLE.md)
- [Chaptarr Integration](../servarr/CHAPTARR.md)
