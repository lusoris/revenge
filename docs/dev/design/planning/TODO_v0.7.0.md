# TODO v0.7.0 - Media

<!-- DESIGN: planning, README, SCAFFOLD_TEMPLATE, test_output_claude -->

> Additional Content Modules

**Status**: 🔴 Not Started
**Tag**: `v0.7.0`
**Focus**: Audiobooks, Books, Podcasts

**Depends On**: [v0.6.0](TODO_v0.6.0.md) (Playback infrastructure)

---

## Overview

This milestone adds support for audiobooks, ebooks, and podcasts. These content types share some patterns with audio/video but have unique requirements for reading progress, chapter navigation, and RSS feed handling.

---

## Deliverables

### Audiobook Module

- [ ] **Database Schema** (`migrations/`)
  - [ ] `public.audiobooks` table
    - [ ] id, title, subtitle
    - [ ] author_id, narrator_id
    - [ ] series_id, series_position
    - [ ] description, publisher
    - [ ] release_date, duration_ms
    - [ ] cover_path
    - [ ] asin, audnexus_id
  - [ ] `public.audiobook_authors` table
  - [ ] `public.audiobook_narrators` table
  - [ ] `public.audiobook_series` table
  - [ ] `public.audiobook_chapters` table
    - [ ] audiobook_id
    - [ ] chapter_number
    - [ ] title
    - [ ] start_time_ms, end_time_ms
  - [ ] `public.audiobook_files` table
  - [ ] `public.audiobook_progress` table
    - [ ] user_id, audiobook_id
    - [ ] current_chapter, position_ms
    - [ ] finished

- [ ] **Entity** (`internal/content/audiobook/entity.go`)
  - [ ] Audiobook struct
  - [ ] Author struct
  - [ ] Narrator struct
  - [ ] Series struct
  - [ ] Chapter struct
  - [ ] Progress struct

- [ ] **Repository** (`internal/content/audiobook/`)
  - [ ] CRUD operations
  - [ ] Chapter operations
  - [ ] Progress operations

- [ ] **Service** (`internal/content/audiobook/service.go`)
  - [ ] Get audiobook by ID
  - [ ] Get chapters
  - [ ] Get/Set progress
  - [ ] Get continue listening
  - [ ] List by author/narrator/series
  - [ ] Search audiobooks

- [ ] **Library Provider** (`internal/content/audiobook/library_service.go`)
  - [ ] Scan M4B, MP3, etc.
  - [ ] Parse chapter info from files
  - [ ] Match to metadata

- [ ] **Handler** (`internal/api/audiobook_handler.go`)
  - [ ] `GET /api/v1/audiobooks`
  - [ ] `GET /api/v1/audiobooks/:id`
  - [ ] `GET /api/v1/audiobooks/:id/chapters`
  - [ ] `GET /api/v1/audiobooks/:id/stream`
  - [ ] `POST /api/v1/audiobooks/:id/progress`
  - [ ] `GET /api/v1/audiobooks/:id/progress`
  - [ ] `GET /api/v1/audiobook-authors`
  - [ ] `GET /api/v1/audiobook-authors/:id`
  - [ ] `GET /api/v1/audiobook-series`
  - [ ] `GET /api/v1/audiobook-series/:id`

- [ ] **River Jobs**
  - [ ] AudiobookLibraryScanJob
  - [ ] AudiobookMetadataRefreshJob

### Audnexus Integration

- [ ] **Audnexus Client** (`internal/service/metadata/audnexus/client.go`)
  - [ ] API implementation
  - [ ] Rate limiting

- [ ] **Audnexus Service** (`internal/service/metadata/audnexus/service.go`)
  - [ ] Search audiobook
  - [ ] Get audiobook details
  - [ ] Get author details
  - [ ] Get chapter info
  - [ ] Get cover images

### Book (eBook) Module

- [ ] **Database Schema** (`migrations/`)
  - [ ] `public.books` table
    - [ ] id, title, subtitle
    - [ ] author_id
    - [ ] series_id, series_position
    - [ ] description, publisher
    - [ ] publish_date, page_count
    - [ ] cover_path
    - [ ] isbn_10, isbn_13
    - [ ] openlibrary_id, goodreads_id
  - [ ] `public.book_authors` table
  - [ ] `public.book_series` table
  - [ ] `public.book_genres` table
  - [ ] `public.book_files` table
    - [ ] book_id, path
    - [ ] format (epub, pdf, mobi)
  - [ ] `public.book_reading_progress` table
    - [ ] user_id, book_id
    - [ ] current_page, total_pages
    - [ ] percentage
    - [ ] epub_position (CFI)
    - [ ] finished

- [ ] **Entity** (`internal/content/book/entity.go`)
  - [ ] Book struct
  - [ ] Author struct
  - [ ] Series struct
  - [ ] ReadingProgress struct

- [ ] **Repository** (`internal/content/book/`)
  - [ ] CRUD operations
  - [ ] Progress operations

- [ ] **Service** (`internal/content/book/service.go`)
  - [ ] Get book by ID
  - [ ] Get/Set reading progress
  - [ ] Get continue reading
  - [ ] List by author/series/genre
  - [ ] Search books

- [ ] **Library Provider** (`internal/content/book/library_service.go`)
  - [ ] Scan EPUB, PDF, MOBI, CBZ/CBR
  - [ ] Extract metadata from files
  - [ ] Extract cover images
  - [ ] Parse TOC (EPUB)

- [ ] **Handler** (`internal/api/book_handler.go`)
  - [ ] `GET /api/v1/books`
  - [ ] `GET /api/v1/books/:id`
  - [ ] `GET /api/v1/books/:id/download` (file access)
  - [ ] `GET /api/v1/books/:id/read` (reader API)
  - [ ] `POST /api/v1/books/:id/progress`
  - [ ] `GET /api/v1/books/:id/progress`
  - [ ] `GET /api/v1/book-authors`
  - [ ] `GET /api/v1/book-series`

- [ ] **River Jobs**
  - [ ] BookLibraryScanJob
  - [ ] BookMetadataRefreshJob

### OpenLibrary Integration

- [ ] **OpenLibrary Client** (`internal/service/metadata/openlibrary/client.go`)
  - [ ] API implementation
  - [ ] Rate limiting (fair use)

- [ ] **OpenLibrary Service** (`internal/service/metadata/openlibrary/service.go`)
  - [ ] Search book by title
  - [ ] Search by ISBN
  - [ ] Get book details
  - [ ] Get author details
  - [ ] Get cover images

### Goodreads Integration (Limited)

- [ ] **Goodreads Service** (`internal/service/metadata/goodreads/`)
  - [ ] Web scraping (no official API)
  - [ ] Get book ratings
  - [ ] Get similar books
  - [ ] Note: Limited due to no public API

### Podcast Module

- [ ] **Database Schema** (`migrations/`)
  - [ ] `public.podcasts` table
    - [ ] id, title, author
    - [ ] description, summary
    - [ ] feed_url
    - [ ] image_url, cover_path
    - [ ] website_url
    - [ ] language, explicit
    - [ ] last_fetch_at
    - [ ] itunes_id, spotify_id
  - [ ] `public.podcast_episodes` table
    - [ ] id, podcast_id
    - [ ] guid (from RSS)
    - [ ] title, description
    - [ ] published_at
    - [ ] duration_ms
    - [ ] audio_url
    - [ ] file_path (if downloaded)
    - [ ] episode_type (full, trailer, bonus)
    - [ ] season_number, episode_number
  - [ ] `public.podcast_subscriptions` table
    - [ ] user_id, podcast_id
    - [ ] subscribed_at
  - [ ] `public.podcast_episode_progress` table
    - [ ] user_id, episode_id
    - [ ] position_ms
    - [ ] completed

- [ ] **Entity** (`internal/content/podcast/entity.go`)
  - [ ] Podcast struct
  - [ ] Episode struct
  - [ ] Subscription struct
  - [ ] EpisodeProgress struct

- [ ] **Repository** (`internal/content/podcast/`)
  - [ ] Podcast CRUD
  - [ ] Episode CRUD
  - [ ] Subscription operations
  - [ ] Progress operations

- [ ] **Service** (`internal/content/podcast/service.go`)
  - [ ] Add podcast (by feed URL)
  - [ ] Refresh feed
  - [ ] Get podcast by ID
  - [ ] Get episodes (paginated)
  - [ ] Subscribe/Unsubscribe
  - [ ] Get subscriptions
  - [ ] Get/Set episode progress
  - [ ] Get continue listening
  - [ ] Download episode
  - [ ] Search podcasts

- [ ] **RSS Parser** (`internal/content/podcast/rss.go`)
  - [ ] gofeed integration
  - [ ] Parse iTunes extensions
  - [ ] Parse Spotify extensions
  - [ ] Handle enclosures (audio files)
  - [ ] Handle chapter markers (if available)

- [ ] **Handler** (`internal/api/podcast_handler.go`)
  - [ ] `GET /api/v1/podcasts`
  - [ ] `GET /api/v1/podcasts/:id`
  - [ ] `GET /api/v1/podcasts/:id/episodes`
  - [ ] `GET /api/v1/podcasts/episodes/:id`
  - [ ] `GET /api/v1/podcasts/episodes/:id/stream`
  - [ ] `POST /api/v1/podcasts` (add by URL)
  - [ ] `POST /api/v1/podcasts/:id/refresh`
  - [ ] `POST /api/v1/podcasts/:id/subscribe`
  - [ ] `DELETE /api/v1/podcasts/:id/subscribe`
  - [ ] `GET /api/v1/users/me/podcast-subscriptions`
  - [ ] `POST /api/v1/podcasts/episodes/:id/progress`
  - [ ] `GET /api/v1/podcasts/episodes/:id/progress`
  - [ ] `POST /api/v1/podcasts/episodes/:id/download`

- [ ] **River Jobs**
  - [ ] PodcastFeedRefreshJob
  - [ ] PodcastEpisodeDownloadJob
  - [ ] PodcastBatchRefreshJob (scheduled)

### iTunes/Podcast Index Integration

- [ ] **Podcast Search Service** (`internal/service/metadata/podcast/`)
  - [ ] iTunes Search API
  - [ ] Podcast Index API (optional)
  - [ ] Search podcasts
  - [ ] Get podcast details
  - [ ] Category browsing

### Frontend Updates

- [ ] **Audiobooks Section** (`/audiobooks`)
  - [ ] Audiobook grid with covers
  - [ ] Series grouping
  - [ ] Author view
  - [ ] Narrator view

- [ ] **Audiobook Player**
  - [ ] Chapter navigation
  - [ ] Speed control (0.5x - 3x)
  - [ ] Sleep timer
  - [ ] Bookmark support
  - [ ] Progress syncing

- [ ] **Books Section** (`/books`)
  - [ ] Book grid with covers
  - [ ] Series view
  - [ ] Author view
  - [ ] Genre filtering

- [ ] **eBook Reader**
  - [ ] EPUB reader (epub.js or similar)
  - [ ] PDF viewer
  - [ ] Reading progress sync
  - [ ] Font size/family settings
  - [ ] Theme (sepia, dark)
  - [ ] Highlight/notes (optional)

- [ ] **Podcasts Section** (`/podcasts`)
  - [ ] Subscribed podcasts
  - [ ] Discover/Search
  - [ ] Episode list
  - [ ] Download management

- [ ] **Podcast Player**
  - [ ] Episode player
  - [ ] Speed control
  - [ ] Skip forward/back (15s/30s)
  - [ ] Mark as played
  - [ ] Chapter markers (if available)

---

## Verification Checklist

- [ ] Audiobooks scan and display
- [ ] Audiobook player with chapters works
- [ ] Books scan and display
- [ ] eBook reader works (EPUB, PDF)
- [ ] Podcasts can be added by URL
- [ ] RSS feeds refresh automatically
- [ ] Podcast episodes stream
- [ ] Progress syncs across all media types
- [ ] All tests pass (80%+ coverage)
- [ ] CI pipeline passes

---

## Dependencies from SOURCE_OF_TRUTH

| Package | Version | Purpose |
|---------|---------|---------|
| github.com/mmcdole/gofeed | latest | RSS/Atom parsing |
| github.com/wtolson/go-taglib | latest | Audio metadata |

---

## Related Documentation

- [ROADMAP.md](ROADMAP.md) - Full roadmap overview
- [00_SOURCE_OF_TRUTH.md](../00_SOURCE_OF_TRUTH.md) - Authoritative versions
- [AUDIOBOOK_MODULE.md](../features/audiobook/AUDIOBOOK_MODULE.md) - Audiobook design
- [BOOK_MODULE.md](../features/book/BOOK_MODULE.md) - Book design
- [PODCASTS.md](../features/podcasts/PODCASTS.md) - Podcast design
- [AUDIBLE.md](../integrations/metadata/books/AUDIBLE.md) - Audnexus integration
- [OPENLIBRARY.md](../integrations/metadata/books/OPENLIBRARY.md) - OpenLibrary integration
