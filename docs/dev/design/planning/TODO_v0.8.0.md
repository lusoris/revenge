# TODO v0.8.0 - Intelligence

<!-- DESIGN: planning, README, test_output_claude, test_output_wiki -->


<!-- TOC-START -->

## Table of Contents

- [Overview](#overview)
- [Deliverables](#deliverables)
  - [Scrobbling Service](#scrobbling-service)
  - [Analytics Service](#analytics-service)
  - [Notification Service](#notification-service)
  - [Request System](#request-system)
  - [Fingerprint Service](#fingerprint-service)
  - [Grants Service (Sharing)](#grants-service-sharing)
  - [i18n Support](#i18n-support)
- [Verification Checklist](#verification-checklist)
- [Dependencies](#dependencies)
- [Related Documentation](#related-documentation)

<!-- TOC-END -->


> Advanced Features

**Status**: 🔴 Not Started
**Tag**: `v0.8.0`
**Focus**: Scrobbling, Analytics, Notifications, Requests, i18n

**Depends On**: [v0.7.0](TODO_v0.7.0.md) (All content modules needed)

---

## Overview

This milestone adds cross-cutting intelligence features: scrobbling to external services (Trakt, Last.fm, ListenBrainz), usage analytics, notification systems, content request workflows, media fingerprinting, sharing via grants, and internationalization.

---

## Deliverables

### Scrobbling Service

- [ ] **Database Schema** (`migrations/`)
  - [ ] `shared.scrobble_connections` table
    - [ ] user_id, provider (trakt, lastfm, listenbrainz)
    - [ ] access_token, refresh_token
    - [ ] expires_at
    - [ ] enabled
    - [ ] scrobble_movies, scrobble_tv, scrobble_music
  - [ ] `shared.scrobble_queue` table
    - [ ] user_id, provider
    - [ ] media_type, media_id
    - [ ] action (start, stop, scrobble)
    - [ ] timestamp
    - [ ] status (pending, sent, failed)
    - [ ] retry_count

- [ ] **Scrobble Service** (`internal/service/scrobbling/service.go`)
  - [ ] Connect provider (OAuth)
  - [ ] Disconnect provider
  - [ ] Queue scrobble
  - [ ] Process scrobble queue
  - [ ] Sync watch history (import from provider)
  - [ ] Check connection status

- [ ] **Trakt Integration** (`internal/service/scrobbling/trakt/`)
  - [ ] OAuth 2.0 authentication
  - [ ] Scrobble movie/episode (start, pause, stop)
  - [ ] Sync watched history
  - [ ] Sync watchlist
  - [ ] Sync ratings
  - [ ] Get recommendations

- [ ] **Last.fm Integration** (`internal/service/scrobbling/lastfm/`)
  - [ ] API key + session authentication
  - [ ] Scrobble track
  - [ ] Update now playing
  - [ ] Love/Unlove track
  - [ ] Get similar tracks/artists

- [ ] **ListenBrainz Integration** (`internal/service/scrobbling/listenbrainz/`)
  - [ ] Token authentication
  - [ ] Submit listen
  - [ ] Submit now playing
  - [ ] Get recommendations

- [ ] **Handler** (`internal/api/scrobbling_handler.go`)
  - [ ] `GET /api/v1/users/me/scrobbling/connections`
  - [ ] `POST /api/v1/users/me/scrobbling/connect/:provider`
  - [ ] `GET /api/v1/users/me/scrobbling/callback/:provider`
  - [ ] `DELETE /api/v1/users/me/scrobbling/disconnect/:provider`
  - [ ] `PATCH /api/v1/users/me/scrobbling/:provider/settings`
  - [ ] `POST /api/v1/users/me/scrobbling/:provider/sync`

- [ ] **River Jobs**
  - [ ] ScrobbleProcessJob - Process queue
  - [ ] ScrobbleSyncJob - Sync history

### Analytics Service

- [ ] **Database Schema** (`migrations/`)
  - [ ] `shared.play_statistics` table (aggregated)
    - [ ] media_type, media_id
    - [ ] total_plays, unique_users
    - [ ] total_duration_ms
    - [ ] first_played, last_played
  - [ ] `shared.user_statistics` table
    - [ ] user_id
    - [ ] period (day, week, month, year, all)
    - [ ] movies_watched, episodes_watched
    - [ ] music_tracks_played
    - [ ] total_watch_time_ms

- [ ] **Analytics Service** (`internal/service/analytics/service.go`)
  - [ ] Record play event
  - [ ] Get media statistics
  - [ ] Get user statistics
  - [ ] Get server statistics
  - [ ] Get popular content
  - [ ] Get trending content
  - [ ] Generate reports

- [ ] **Handler** (`internal/api/analytics_handler.go`)
  - [ ] `GET /api/v1/users/me/stats`
  - [ ] `GET /api/v1/users/me/stats/history`
  - [ ] `GET /api/v1/admin/analytics/overview`
  - [ ] `GET /api/v1/admin/analytics/popular`
  - [ ] `GET /api/v1/admin/analytics/users`
  - [ ] `GET /api/v1/admin/analytics/media/:type/:id`

- [ ] **River Jobs**
  - [ ] AnalyticsAggregationJob - Daily/weekly rollup

### Notification Service

- [ ] **Database Schema** (`migrations/`)
  - [ ] `shared.notification_channels` table
    - [ ] user_id, channel_type (email, push, webhook)
    - [ ] enabled, config (JSON)
  - [ ] `shared.notifications` table
    - [ ] user_id
    - [ ] type (new_content, request_update, system)
    - [ ] title, message
    - [ ] data (JSON)
    - [ ] read_at, created_at
  - [ ] `shared.notification_preferences` table
    - [ ] user_id, notification_type
    - [ ] email, push, in_app

- [ ] **Notification Service** (`internal/service/notification/service.go`)
  - [ ] Send notification
  - [ ] Get user notifications
  - [ ] Mark as read
  - [ ] Clear notifications
  - [ ] Get/Set preferences
  - [ ] Register push token (FCM)

- [ ] **Email Provider** (`internal/service/notification/email.go`)
  - [ ] go-mail integration
  - [ ] Template rendering
  - [ ] Queue emails

- [ ] **Push Provider** (`internal/service/notification/push.go`)
  - [ ] FCM integration (go-fcm)
  - [ ] Device token management
  - [ ] Send push notification

- [ ] **Webhook Provider** (`internal/service/notification/webhook.go`)
  - [ ] Custom webhook endpoints
  - [ ] Payload templates
  - [ ] Retry logic

- [ ] **Handler** (`internal/api/notification_handler.go`)
  - [ ] `GET /api/v1/users/me/notifications`
  - [ ] `POST /api/v1/users/me/notifications/:id/read`
  - [ ] `POST /api/v1/users/me/notifications/read-all`
  - [ ] `DELETE /api/v1/users/me/notifications/:id`
  - [ ] `GET /api/v1/users/me/notification-preferences`
  - [ ] `PATCH /api/v1/users/me/notification-preferences`
  - [ ] `POST /api/v1/users/me/push-token`

- [ ] **River Jobs**
  - [ ] EmailSendJob
  - [ ] PushSendJob
  - [ ] WebhookSendJob
  - [ ] NewContentNotificationJob

### Request System

- [ ] **Database Schema** (`migrations/`)
  - [ ] `shared.content_requests` table
    - [ ] id, user_id
    - [ ] media_type (movie, tv, music)
    - [ ] external_id (tmdb_id, etc.)
    - [ ] title, year, overview
    - [ ] poster_path
    - [ ] status (pending, approved, denied, available)
    - [ ] admin_note
    - [ ] requested_at, processed_at
    - [ ] processed_by

- [ ] **Request Service** (`internal/service/request/service.go`)
  - [ ] Create request
  - [ ] Get user requests
  - [ ] Get all requests (admin)
  - [ ] Approve request
  - [ ] Deny request
  - [ ] Check if content now available
  - [ ] Notify user on status change

- [ ] **Handler** (`internal/api/request_handler.go`)
  - [ ] `GET /api/v1/requests`
  - [ ] `POST /api/v1/requests`
  - [ ] `GET /api/v1/requests/:id`
  - [ ] `DELETE /api/v1/requests/:id`
  - [ ] `GET /api/v1/admin/requests`
  - [ ] `POST /api/v1/admin/requests/:id/approve`
  - [ ] `POST /api/v1/admin/requests/:id/deny`

- [ ] **River Jobs**
  - [ ] RequestAvailabilityCheckJob - Check if requested content arrived

### Fingerprint Service

- [ ] **Database Schema** (`migrations/`)
  - [ ] `public.media_fingerprints` table
    - [ ] media_id, media_type
    - [ ] fingerprint_type (chromaprint, phash)
    - [ ] fingerprint_data
    - [ ] duration_ms

- [ ] **Fingerprint Service** (`internal/service/fingerprint/service.go`)
  - [ ] Generate audio fingerprint (Chromaprint)
  - [ ] Generate image hash (pHash)
  - [ ] Match fingerprint
  - [ ] Detect duplicates

- [ ] **Handler** (`internal/api/fingerprint_handler.go`)
  - [ ] `POST /api/v1/admin/fingerprint/generate/:type/:id`
  - [ ] `GET /api/v1/admin/fingerprint/duplicates`

- [ ] **River Jobs**
  - [ ] FingerprintGenerateJob
  - [ ] DuplicateDetectionJob

### Grants Service (Sharing)

- [ ] **Database Schema** (`migrations/`)
  - [ ] `shared.media_grants` table
    - [ ] id, created_by
    - [ ] media_type, media_id
    - [ ] token (unique share URL)
    - [ ] expires_at
    - [ ] max_views, current_views
    - [ ] requires_password, password_hash
    - [ ] created_at

- [ ] **Grants Service** (`internal/service/grants/service.go`)
  - [ ] Create grant (share link)
  - [ ] Validate grant
  - [ ] Revoke grant
  - [ ] List user's grants
  - [ ] Track view count

- [ ] **Handler** (`internal/api/grants_handler.go`)
  - [ ] `GET /api/v1/users/me/grants`
  - [ ] `POST /api/v1/users/me/grants`
  - [ ] `DELETE /api/v1/users/me/grants/:id`
  - [ ] `GET /api/v1/share/:token` (public access)
  - [ ] `GET /api/v1/share/:token/stream`

### i18n Support

- [ ] **Translation Files** (`locales/`)
  - [ ] `en.json` - English (default)
  - [ ] `de.json` - German
  - [ ] `es.json` - Spanish
  - [ ] `fr.json` - French
  - [ ] (Add more as contributed)

- [ ] **Backend i18n**
  - [ ] Error message localization
  - [ ] Email template localization
  - [ ] Date/time formatting

- [ ] **Frontend i18n**
  - [ ] i18n library setup (svelte-i18n)
  - [ ] Language switcher
  - [ ] Browser language detection
  - [ ] User language preference

- [ ] **Handler** (`internal/api/i18n_handler.go`)
  - [ ] `GET /api/v1/i18n/:locale`
  - [ ] `GET /api/v1/i18n/locales`

---

## Verification Checklist

- [ ] Trakt scrobbling works (movies, TV)
- [ ] Last.fm scrobbling works (music)
- [ ] ListenBrainz scrobbling works
- [ ] Analytics dashboard shows statistics
- [ ] Notifications deliver (email, push)
- [ ] Content requests workflow works
- [ ] Share links work for external users
- [ ] Multiple languages selectable
- [ ] All tests pass (80%+ coverage)
- [ ] CI pipeline passes

---

## Dependencies

| Package | Version | Purpose |
|---------|---------|---------|
| github.com/wneessen/go-mail | v0.6.2 | Email sending |
| github.com/appleboy/go-fcm | v0.2.1 | Push notifications |
| github.com/corona10/goimagehash | latest | Image hashing |

---

## Related Documentation

- [ROADMAP.md](ROADMAP.md) - Full roadmap overview
- [DESIGN_INDEX.md](../DESIGN_INDEX.md) - Full design documentation index
- [SCROBBLING.md](../features/shared/SCROBBLING.md) - Scrobbling design
- [ANALYTICS_SERVICE.md](../features/shared/ANALYTICS_SERVICE.md) - Analytics design
- [NOTIFICATION.md](../services/NOTIFICATION.md) - Notification service design
- [REQUEST_SYSTEM.md](../features/shared/REQUEST_SYSTEM.md) - Request system design
- [GRANTS.md](../services/GRANTS.md) - Grants service design
- [I18N.md](../features/shared/I18N.md) - i18n design
- [TRAKT.md](../integrations/scrobbling/TRAKT.md) - Trakt integration
- [LASTFM_SCROBBLE.md](../integrations/scrobbling/LASTFM_SCROBBLE.md) - Last.fm integration
- [LISTENBRAINZ.md](../integrations/scrobbling/LISTENBRAINZ.md) - ListenBrainz integration
