# TODO v0.6.0 - Playback

<!-- DESIGN: planning, README, SCAFFOLD_TEMPLATE, test_output_claude -->

> Playback Features

**Status**: 🔴 Not Started
**Tag**: `v0.6.0`
**Focus**: Advanced Playback Features

**Depends On**:
- [v0.5.0](TODO_v0.5.0.md) (Content modules needed)
- [v0.1.0](TODO_v0.1.0.md) (Shared playback infrastructure)

---

## Overview

This milestone adds advanced playback features: trickplay (preview thumbnails), skip intro/credits detection, watch next logic, SyncPlay (synchronized group playback), media enhancements (audio boost, subtitles), and casting support.

---

## Deliverables

### Trickplay (Preview Thumbnails)

- [ ] **Database Schema** (`migrations/`)
  - [ ] `public.trickplay_manifests` table
    - [ ] media_id, media_type
    - [ ] interval_ms (e.g., 10000 for 10s)
    - [ ] tile_width, tile_height
    - [ ] columns_per_image
    - [ ] image_count
    - [ ] status (pending, processing, complete, failed)
  - [ ] `public.trickplay_images` table
    - [ ] manifest_id
    - [ ] image_index
    - [ ] path
    - [ ] start_time_ms, end_time_ms

- [ ] **Service** (`internal/service/playback/trickplay/service.go`)
  - [ ] Generate trickplay for media
  - [ ] Get trickplay manifest
  - [ ] Get trickplay image
  - [ ] Check trickplay status
  - [ ] Delete trickplay data

- [ ] **Generator** (`internal/service/playback/trickplay/generator.go`)
  - [ ] FFmpeg frame extraction
  - [ ] Thumbnail generation (govips)
  - [ ] Sprite sheet creation
  - [ ] WebVTT manifest generation
  - [ ] Configurable intervals

- [ ] **Handler** (`internal/api/trickplay_handler.go`)
  - [ ] `GET /api/v1/media/:type/:id/trickplay`
  - [ ] `GET /api/v1/media/:type/:id/trickplay/manifest.vtt`
  - [ ] `GET /api/v1/media/:type/:id/trickplay/tiles/:index.jpg`
  - [ ] `POST /api/v1/admin/trickplay/generate/:type/:id`
  - [ ] `DELETE /api/v1/admin/trickplay/:type/:id`

- [ ] **River Jobs** (`internal/service/playback/trickplay/jobs.go`)
  - [ ] TrickplayGenerateJob - Generate for single media
  - [ ] TrickplayBatchJob - Generate for library

- [ ] **Configuration**
  - [ ] Interval setting (default 10s)
  - [ ] Thumbnail dimensions
  - [ ] Quality settings
  - [ ] Auto-generate on library scan

- [ ] **Frontend Integration**
  - [ ] Player scrubbing preview
  - [ ] Load on hover
  - [ ] Position calculation

- [ ] **Tests**
  - [ ] Unit tests
  - [ ] Integration tests with FFmpeg

### Skip Intro/Credits

- [ ] **Database Schema** (`migrations/`)
  - [ ] `public.media_chapters` table
    - [ ] media_id, media_type
    - [ ] chapter_type (intro, credits, recap, etc.)
    - [ ] start_time_ms, end_time_ms
    - [ ] confidence (for auto-detected)
    - [ ] source (auto, manual, external)

- [ ] **Service** (`internal/service/playback/chapters/service.go`)
  - [ ] Get chapters for media
  - [ ] Add chapter marker
  - [ ] Update chapter marker
  - [ ] Delete chapter marker
  - [ ] Auto-detect intro (audio fingerprinting)
  - [ ] Auto-detect credits (scene analysis)

- [ ] **Detection** (`internal/service/playback/chapters/detection.go`)
  - [ ] Audio fingerprint matching across episodes
  - [ ] Black frame detection (credits)
  - [ ] Credit text detection (optional)
  - [ ] Confidence scoring

- [ ] **Handler** (`internal/api/chapters_handler.go`)
  - [ ] `GET /api/v1/media/:type/:id/chapters`
  - [ ] `POST /api/v1/media/:type/:id/chapters`
  - [ ] `PATCH /api/v1/media/:type/:id/chapters/:chapterId`
  - [ ] `DELETE /api/v1/media/:type/:id/chapters/:chapterId`
  - [ ] `POST /api/v1/admin/chapters/detect/:type/:id`

- [ ] **River Jobs**
  - [ ] ChapterDetectionJob - Detect for single media
  - [ ] ChapterBatchDetectionJob - Detect for series

- [ ] **Frontend Integration**
  - [ ] "Skip Intro" button (appears at intro start)
  - [ ] "Skip Credits" button
  - [ ] Auto-skip setting (per user)
  - [ ] Chapter markers on timeline

- [ ] **Tests**
  - [ ] Unit tests
  - [ ] Integration tests

### Watch Next / Continue Watching

- [ ] **Database Schema** (`migrations/`)
  - [ ] `public.user_watch_queue` table
    - [ ] user_id, media_id, media_type
    - [ ] position (order in queue)
    - [ ] added_at, source (auto, manual)

- [ ] **Service** (`internal/service/playback/watchnext/service.go`)
  - [ ] Get continue watching list
  - [ ] Get up next (what to play after current)
  - [ ] Add to watch queue
  - [ ] Remove from watch queue
  - [ ] Reorder watch queue
  - [ ] Auto-populate from watch history

- [ ] **Logic** (`internal/service/playback/watchnext/logic.go`)
  - [ ] Calculate next episode in series
  - [ ] Handle season boundaries
  - [ ] Handle series completion
  - [ ] Movie recommendations after movie
  - [ ] Resume position calculation
  - [ ] "Mark as watched" threshold (90%)

- [ ] **Handler** (`internal/api/watchnext_handler.go`)
  - [ ] `GET /api/v1/users/me/continue-watching`
  - [ ] `GET /api/v1/users/me/up-next`
  - [ ] `GET /api/v1/media/:type/:id/up-next`
  - [ ] `POST /api/v1/users/me/watch-queue`
  - [ ] `DELETE /api/v1/users/me/watch-queue/:id`
  - [ ] `PATCH /api/v1/users/me/watch-queue/reorder`

- [ ] **Frontend Integration**
  - [ ] Continue watching carousel
  - [ ] Up next overlay (before video ends)
  - [ ] Auto-play next setting
  - [ ] Countdown timer (10s)

- [ ] **Tests**
  - [ ] Unit tests
  - [ ] Integration tests

### SyncPlay (Group Playback)

- [ ] **Database Schema** (`migrations/`)
  - [ ] `public.syncplay_groups` table
    - [ ] id, name, created_by
    - [ ] current_media_id, media_type
    - [ ] playback_state (playing, paused, stopped)
    - [ ] position_ms
    - [ ] created_at
  - [ ] `public.syncplay_participants` table
    - [ ] group_id, user_id
    - [ ] is_ready
    - [ ] joined_at

- [ ] **Service** (`internal/service/playback/syncplay/service.go`)
  - [ ] Create group
  - [ ] Join group
  - [ ] Leave group
  - [ ] Set group media
  - [ ] Play/Pause/Seek group
  - [ ] Get group state
  - [ ] Broadcast state changes
  - [ ] Sync position (latency compensation)

- [ ] **WebSocket Handler** (`internal/api/syncplay_ws.go`)
  - [ ] `WS /api/v1/syncplay/:groupId`
  - [ ] Join message
  - [ ] Leave message
  - [ ] Play/Pause/Seek commands
  - [ ] Position updates
  - [ ] Ready state
  - [ ] Chat messages (optional)

- [ ] **Sync Algorithm** (`internal/service/playback/syncplay/sync.go`)
  - [ ] NTP-like time synchronization
  - [ ] Latency measurement
  - [ ] Position correction
  - [ ] Buffering coordination

- [ ] **Handler** (`internal/api/syncplay_handler.go`)
  - [ ] `POST /api/v1/syncplay/groups` (create)
  - [ ] `GET /api/v1/syncplay/groups/:id`
  - [ ] `DELETE /api/v1/syncplay/groups/:id`
  - [ ] `POST /api/v1/syncplay/groups/:id/join`
  - [ ] `POST /api/v1/syncplay/groups/:id/leave`

- [ ] **Frontend Integration**
  - [ ] Create/Join group UI
  - [ ] Participant list
  - [ ] Ready check system
  - [ ] Synced playback controls
  - [ ] Chat sidebar (optional)

- [ ] **Tests**
  - [ ] Unit tests
  - [ ] Integration tests (WebSocket)

### Media Enhancements

- [ ] **Audio Boost**
  - [ ] Normalize audio levels
  - [ ] Boost quiet dialogue
  - [ ] Settings per user
  - [ ] Dynamic range compression

- [ ] **Subtitle Handling** (`internal/service/playback/subtitles/`)
  - [ ] Extract embedded subtitles
  - [ ] Parse SRT, ASS, WebVTT
  - [ ] Convert formats (FFmpeg)
  - [ ] Subtitle preferences (language, format)
  - [ ] External subtitle support

- [ ] **Handler** (`internal/api/subtitles_handler.go`)
  - [ ] `GET /api/v1/media/:type/:id/subtitles`
  - [ ] `GET /api/v1/media/:type/:id/subtitles/:trackId`
  - [ ] `POST /api/v1/media/:type/:id/subtitles` (upload external)

- [ ] **Frontend Integration**
  - [ ] Subtitle track selection
  - [ ] Subtitle styling (size, color, background)
  - [ ] Audio track selection
  - [ ] Audio normalization toggle

### Chromecast Support

- [ ] **Chromecast Service** (`internal/service/playback/cast/chromecast.go`)
  - [ ] Device discovery (mDNS)
  - [ ] Cast session management
  - [ ] Media loading
  - [ ] Playback control
  - [ ] Status monitoring

- [ ] **Handler** (`internal/api/cast_handler.go`)
  - [ ] `GET /api/v1/cast/devices` (discover)
  - [ ] `POST /api/v1/cast/sessions` (start casting)
  - [ ] `DELETE /api/v1/cast/sessions/:id` (stop)
  - [ ] `POST /api/v1/cast/sessions/:id/play`
  - [ ] `POST /api/v1/cast/sessions/:id/pause`
  - [ ] `POST /api/v1/cast/sessions/:id/seek`

- [ ] **Frontend Integration**
  - [ ] Cast button in player
  - [ ] Device picker
  - [ ] Cast controls

### DLNA Support

- [ ] **DLNA Service** (`internal/service/playback/cast/dlna.go`)
  - [ ] UPnP device discovery
  - [ ] DMR (renderer) control
  - [ ] AV Transport service
  - [ ] Media serving (DMS)

- [ ] **Handler** (`internal/api/dlna_handler.go`)
  - [ ] `GET /api/v1/dlna/devices`
  - [ ] `POST /api/v1/dlna/play`
  - [ ] Control endpoints

- [ ] **Frontend Integration**
  - [ ] DLNA device list
  - [ ] Cast to DLNA device

---

## Verification Checklist

- [ ] Trickplay thumbnails generate and display
- [ ] Skip intro button appears at correct time
- [ ] Continue watching shows correct items
- [ ] Up next plays automatically (if enabled)
- [ ] SyncPlay syncs multiple clients
- [ ] Subtitles display correctly
- [ ] Chromecast casting works
- [ ] DLNA discovery and playback works
- [ ] All tests pass (80%+ coverage)
- [ ] CI pipeline passes

---

## Dependencies from SOURCE_OF_TRUTH

| Package | Version | Purpose |
|---------|---------|---------|
| github.com/asticode/go-astiav | latest | FFmpeg bindings |
| github.com/davidbyttow/govips/v2 | latest | Image processing |
| github.com/gobwas/ws | latest | WebSocket (SyncPlay) |

---

## Related Documentation

- [ROADMAP.md](ROADMAP.md) - Full roadmap overview
- [00_SOURCE_OF_TRUTH.md](../00_SOURCE_OF_TRUTH.md) - Authoritative versions
- [TRICKPLAY.md](../features/playback/TRICKPLAY.md) - Trickplay design
- [SKIP_INTRO.md](../features/playback/SKIP_INTRO.md) - Skip intro design
- [WATCH_NEXT_CONTINUE_WATCHING.md](../features/playback/WATCH_NEXT_CONTINUE_WATCHING.md) - Watch next design
- [SYNCPLAY.md](../features/playback/SYNCPLAY.md) - SyncPlay design
- [CHROMECAST.md](../integrations/casting/CHROMECAST.md) - Chromecast integration
- [DLNA.md](../integrations/casting/DLNA.md) - DLNA integration
