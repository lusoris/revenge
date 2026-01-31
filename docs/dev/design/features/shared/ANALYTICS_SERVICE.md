# Tracearr Analytics Service

> Real-time monitoring, analytics, and account sharing detection for Revenge


<!-- TOC-START -->

## Table of Contents

- [Status](#status)
- [Overview](#overview)
- [Core Features (from Tracearr)](#core-features-from-tracearr)
  - [1. Session Monitoring](#1-session-monitoring)
  - [2. Stream Analytics](#2-stream-analytics)
  - [3. Library Analytics](#3-library-analytics)
    - [Overview Page](#overview-page)
    - [Quality Page](#quality-page)
    - [Storage Page](#storage-page)
    - [Watch Page](#watch-page)
  - [4. Account Sharing Detection](#4-account-sharing-detection)
  - [5. Live TV & Music Tracking](#5-live-tv-music-tracking)
  - [6. Stream Map](#6-stream-map)
  - [7. Bulk Actions](#7-bulk-actions)
  - [8. Public API](#8-public-api)
- [Revenge Implementation (Go)](#revenge-implementation-go)
  - [PostgreSQL Schema](#postgresql-schema)
  - [Go Service Structure](#go-service-structure)
  - [River Jobs](#river-jobs)
- [Go Libraries (Open Source)](#go-libraries-open-source)
  - [Geolocation](#geolocation)
  - [Mapping (Frontend)](#mapping-frontend)
  - [Charts (Frontend)](#charts-frontend)
  - [Notifications](#notifications)
- [API Endpoints](#api-endpoints)
  - [Sessions](#sessions)
  - [Violations](#violations)
  - [Trust Scores](#trust-scores)
  - [Library Analytics](#library-analytics)
  - [Sharing Rules](#sharing-rules)
  - [Stream Map](#stream-map)
  - [Bulk Actions](#bulk-actions)
- [Implementation Phases](#implementation-phases)
  - [Phase 1: Session Tracking (Week 1)](#phase-1-session-tracking-week-1)
  - [Phase 2: Sharing Detection (Week 2)](#phase-2-sharing-detection-week-2)
  - [Phase 3: Trust Scoring (Week 2)](#phase-3-trust-scoring-week-2)
  - [Phase 4: Library Analytics (Week 3)](#phase-4-library-analytics-week-3)
  - [Phase 5: Geolocation & Map (Week 3)](#phase-5-geolocation-map-week-3)
  - [Phase 6: Notifications (Week 4)](#phase-6-notifications-week-4)
  - [Phase 7: Frontend Dashboard (Week 4-5)](#phase-7-frontend-dashboard-week-4-5)
- [Implementation Checklist](#implementation-checklist)
  - [Phase 1: Core Infrastructure](#phase-1-core-infrastructure)
  - [Phase 2: Database](#phase-2-database)
  - [Phase 3: Service Layer](#phase-3-service-layer)
  - [Phase 4: Background Jobs](#phase-4-background-jobs)
  - [Phase 5: API Integration](#phase-5-api-integration)
- [Sources & Cross-References](#sources-cross-references)
  - [Cross-Reference Indexes](#cross-reference-indexes)
  - [Referenced Sources](#referenced-sources)
- [Related Design Docs](#related-design-docs)
  - [In This Section](#in-this-section)
  - [Related Topics](#related-topics)
  - [Indexes](#indexes)
- [Related](#related)
- [Notes](#notes)

<!-- TOC-END -->

## Status

| Dimension | Status | Notes |
|-----------|--------|-------|
| Design | ✅ | Full design with PostgreSQL schema, Go service structure, API endpoints |
| Sources | ✅ | Based on Tracearr, extensive Go library references |
| Instructions | ✅ | Implementation checklist complete |
| Code | 🔴 |  |
| Linting | 🔴 |  |
| Unit Testing | 🔴 |  |
| Integration Testing | 🔴 |  |
**Based On**: Tracearr (https://github.com/connorgallopo/Tracearr)
**Priority**: 🟡 HIGH (Phase 9 - External Services)
**License**: AGPL-3.0 (rebuild required to avoid copyleft in Revenge)
**Location**: `internal/service/analytics/`

---

## Overview

**Purpose**: Native analytics and monitoring service for Revenge, replacing the need for external tools like Tautulli/Jellystat/Tracearr.

**Why Build From Scratch**:
- Tracearr is TypeScript/Node.js (React, Fastify, TimescaleDB)
- AGPL-3.0 copyleft license requires sharing modifications
- Native Go integration with Revenge backend
- Single binary deployment (no separate services)
- Optimized for Revenge's architecture (fx, PostgreSQL 18, River)

---

## Core Features (from Tracearr)

### 1. Session Monitoring
- **Real-time stream tracking**: Who's watching, what, where, when, device
- **Session history**: Complete playback logs with geolocation
- **Live session viewer**: Active streams dashboard
- **Playback analytics**: Direct play vs transcoding
- **Bandwidth tracking**: Per-user, per-device, per-content-type

### 2. Stream Analytics
- **Codec breakdowns**: Video/audio codec usage statistics
- **Resolution stats**: 4K vs 1080p vs 720p vs SD distribution
- **Device compatibility scores**: Track client capabilities
- **Transcode efficiency**: % of streams transcoded vs direct play
- **Quality metrics**: Bitrate, resolution, codec per stream
- **Enhanced IP geolocation**: ASN data, continent, postal codes

### 3. Library Analytics

#### Overview Page
- Item counts (movies, TV, music, etc.)
- Storage usage (total GB, per content type)
- Growth charts over time (items added per month)

#### Quality Page
- Resolution distribution (4K/1080p/720p/SD ratios)
- Codec distribution (H.264/H.265/VP9/AV1)
- Track quality ratio changes over time

#### Storage Page
- Usage predictions (growth projections)
- Duplicate detection across libraries
- Stale content identification (never watched)
- ROI analysis (watch hours per GB)

#### Watch Page
- Engagement metrics (completion rates)
- Viewing patterns (hour of day, day of week, month)
- Binge detection (consecutive episodes/seasons)
- Most/least watched content

### 4. Account Sharing Detection

**Six Rule Types**:

1. **Impossible Travel**: Same account in NYC then London 30 minutes later
2. **Simultaneous Locations**: Account streaming from two cities at once
3. **Device Velocity**: Too many unique IPs in short window
4. **Concurrent Streams**: Exceed per-user stream limit
5. **Geo Restrictions**: Block streaming from specific countries
6. **Account Inactivity**: Notify when accounts go dormant

**Trust Scores**:
- Users earn/lose trust based on behavior
- Violations automatically drop scores
- Configurable thresholds for actions (warnings, stream termination, account suspension)

**Real-Time Alerts**:
- Discord webhooks
- Email notifications
- Telegram notifications
- Custom webhook URLs

### 5. Live TV & Music Tracking
- Track live TV sessions (channels, programs)
- Music playback analytics (albums, artists, tracks)
- Content-type-specific dashboards

### 6. Stream Map
- Visualize streams on world map
- Filter by user, server, time period
- Geolocation enrichment (city, region, country, ASN)

### 7. Bulk Actions
- Multi-select operations across tables
- Acknowledge/dismiss violations in bulk
- Reset trust scores
- Enable/disable rules
- Delete session history
- Stream termination

### 8. Public API
- Read-only REST API for third-party integrations
- OpenAPI (Swagger UI) documentation
- API key generation/management
- Integration with Homarr, Home Assistant, etc.

---

## Revenge Implementation (Go)

### PostgreSQL Schema

```sql
-- Session tracking
CREATE TABLE analytics_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    content_type VARCHAR(50) NOT NULL, -- movie, tvshow_episode, music_track, livetv_program
    content_id UUID NOT NULL,          -- Foreign key to content table
    device_name VARCHAR(200),
    device_type VARCHAR(100),          -- Web, Mobile, TV, etc.
    client_name VARCHAR(200),          -- Browser, app name
    client_version VARCHAR(100),
    ip_address INET NOT NULL,
    geo_city VARCHAR(200),
    geo_region VARCHAR(200),
    geo_country VARCHAR(100),
    geo_continent VARCHAR(50),
    geo_postal_code VARCHAR(20),
    geo_asn INT,
    geo_asn_org VARCHAR(200),
    latitude DECIMAL(9,6),
    longitude DECIMAL(9,6),
    stream_type VARCHAR(50),           -- direct_play, transcode
    video_codec VARCHAR(50),
    audio_codec VARCHAR(50),
    resolution VARCHAR(20),            -- 2160p, 1080p, 720p, etc.
    bitrate_kbps INT,
    bandwidth_mbps DECIMAL(10,2),
    started_at TIMESTAMPTZ DEFAULT NOW(),
    stopped_at TIMESTAMPTZ,
    paused_seconds INT DEFAULT 0,
    buffering_seconds INT DEFAULT 0,
    progress_seconds INT,
    duration_seconds INT,
    completed BOOLEAN DEFAULT FALSE
);

-- Sharing violation rules
CREATE TABLE analytics_sharing_rules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(200) NOT NULL,
    type VARCHAR(50) NOT NULL,         -- impossible_travel, simultaneous_locations, device_velocity, concurrent_streams, geo_restriction, inactivity
    enabled BOOLEAN DEFAULT TRUE,
    config JSONB NOT NULL,             -- Rule-specific configuration
    action VARCHAR(50) NOT NULL,       -- warn, terminate_stream, suspend_account
    severity INT CHECK (severity BETWEEN 1 AND 5), -- 1=Critical, 5=Low
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Sharing violations
CREATE TABLE analytics_sharing_violations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    rule_id UUID REFERENCES analytics_sharing_rules(id) ON DELETE CASCADE,
    session_id_1 UUID REFERENCES analytics_sessions(id),
    session_id_2 UUID REFERENCES analytics_sessions(id),
    violation_type VARCHAR(50) NOT NULL,
    details JSONB,
    trust_score_impact INT,           -- Points deducted from trust score
    action_taken VARCHAR(50),         -- warned, stream_terminated, account_suspended
    acknowledged BOOLEAN DEFAULT FALSE,
    detected_at TIMESTAMPTZ DEFAULT NOW(),
    acknowledged_at TIMESTAMPTZ
);

-- User trust scores
CREATE TABLE analytics_user_trust (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    score INT DEFAULT 100 CHECK (score BETWEEN 0 AND 100),
    violations_count INT DEFAULT 0,
    last_violation_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Library analytics (aggregated data)
CREATE TABLE analytics_library_stats (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    content_type VARCHAR(50) NOT NULL,
    stat_type VARCHAR(100) NOT NULL,  -- total_items, total_size_gb, resolution_4k_count, codec_h265_count, etc.
    value BIGINT NOT NULL,
    date DATE NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(content_type, stat_type, date)
);

-- Engagement metrics
CREATE TABLE analytics_engagement (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    content_type VARCHAR(50) NOT NULL,
    content_id UUID NOT NULL,
    watch_count INT DEFAULT 0,
    total_watch_seconds BIGINT DEFAULT 0,
    unique_users INT DEFAULT 0,
    completion_rate DECIMAL(5,2),     -- Percentage (0.00 to 100.00)
    avg_completion_rate DECIMAL(5,2),
    last_watched_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(content_type, content_id)
);

-- Indexes
CREATE INDEX idx_sessions_user_id ON analytics_sessions(user_id);
CREATE INDEX idx_sessions_started_at ON analytics_sessions(started_at DESC);
CREATE INDEX idx_sessions_ip_address ON analytics_sessions(ip_address);
CREATE INDEX idx_sessions_content ON analytics_sessions(content_type, content_id);
CREATE INDEX idx_violations_user_id ON analytics_sharing_violations(user_id);
CREATE INDEX idx_violations_detected_at ON analytics_sharing_violations(detected_at DESC);
CREATE INDEX idx_library_stats_date ON analytics_library_stats(date DESC);
CREATE INDEX idx_engagement_content ON analytics_engagement(content_type, content_id);
```

### Go Service Structure

```
internal/service/analytics/
├── module.go                  # fx module registration
├── service.go                 # Analytics service
├── session_tracker.go         # Real-time session tracking
├── sharing_detector.go        # Account sharing detection
├── trust_scorer.go            # Trust score calculation
├── library_analyzer.go        # Library analytics
├── geolocation.go             # IP geolocation enrichment
├── aggregator.go              # Stats aggregation (River jobs)
└── notifier.go                # Webhook/email/Telegram notifications

internal/api/handlers/analytics/
├── sessions.go                # GET /api/v1/analytics/sessions
├── violations.go              # GET /api/v1/analytics/violations
├── trust.go                   # GET /api/v1/analytics/trust
├── library.go                 # GET /api/v1/analytics/library
├── rules.go                   # CRUD for sharing rules
└── bulk_actions.go            # POST /api/v1/analytics/bulk
```

### River Jobs

```go
// Aggregate library stats (daily)
type AggregateLibraryStatsArgs struct {
    Date time.Time `json:"date"`
}

func (AggregateLibraryStatsArgs) Kind() string { return "analytics.aggregate_library_stats" }

// Detect sharing violations (every 5 minutes)
type DetectSharingViolationsArgs struct{}

func (DetectSharingViolationsArgs) Kind() string { return "analytics.detect_sharing_violations" }

// Cleanup old sessions (daily)
type CleanupOldSessionsArgs struct {
    RetentionDays int `json:"retention_days"`
}

func (CleanupOldSessionsArgs) Kind() string { return "analytics.cleanup_old_sessions" }
```

---

## Go Libraries (Open Source)

### Geolocation
- **github.com/oschwald/geoip2-golang** - MaxMind GeoIP2 reader (city, country, ASN)
- **github.com/ip2location/ip2location-go** - Alternative to MaxMind

### Mapping (Frontend)
- Leaflet.js (already open source)
- OpenStreetMap tiles (free)

### Charts (Frontend)
- Apache ECharts (Apache 2.0 license) - Alternative to Highcharts (commercial)
- Chart.js (MIT license)

### Notifications
- **github.com/bwmarrin/discordgo** - Discord webhooks
- **github.com/go-telegram-bot-api/telegram-bot-api** - Telegram notifications
- **github.com/jordan-wright/email** - SMTP email

---

## API Endpoints

### Sessions
```bash
GET  /api/v1/analytics/sessions             # List all sessions (paginated, filtered)
GET  /api/v1/analytics/sessions/{id}        # Get session detail
GET  /api/v1/analytics/sessions/active      # Get active sessions only
DELETE /api/v1/analytics/sessions/{id}      # Delete session (admin)
POST /api/v1/analytics/sessions/bulk-delete # Bulk delete sessions
```

### Violations
```bash
GET  /api/v1/analytics/violations           # List violations (paginated, filtered)
GET  /api/v1/analytics/violations/{id}      # Get violation detail
PUT  /api/v1/analytics/violations/{id}/acknowledge # Acknowledge violation
POST /api/v1/analytics/violations/bulk-acknowledge # Bulk acknowledge
```

### Trust Scores
```bash
GET  /api/v1/analytics/trust                # List all user trust scores
GET  /api/v1/analytics/trust/{user_id}      # Get user trust score
PUT  /api/v1/analytics/trust/{user_id}/reset # Reset trust score (admin)
```

### Library Analytics
```bash
GET  /api/v1/analytics/library/overview     # Overview stats (item counts, storage)
GET  /api/v1/analytics/library/quality      # Quality stats (resolution, codec)
GET  /api/v1/analytics/library/storage      # Storage stats (predictions, duplicates, ROI)
GET  /api/v1/analytics/library/watch        # Watch stats (engagement, patterns, binge)
```

### Sharing Rules
```bash
GET  /api/v1/analytics/rules                # List all rules
GET  /api/v1/analytics/rules/{id}           # Get rule detail
POST /api/v1/analytics/rules                # Create rule
PUT  /api/v1/analytics/rules/{id}           # Update rule
DELETE /api/v1/analytics/rules/{id}         # Delete rule
```

### Stream Map
```bash
GET  /api/v1/analytics/map                  # Get stream geolocation data
```

### Bulk Actions
```bash
POST /api/v1/analytics/bulk/acknowledge     # Bulk acknowledge violations
POST /api/v1/analytics/bulk/reset-trust     # Bulk reset trust scores
POST /api/v1/analytics/bulk/delete-sessions # Bulk delete sessions
```

---

## Implementation Phases

### Phase 1: Session Tracking (Week 1)
- [ ] PostgreSQL schema (sessions, user_trust)
- [ ] Session tracker service
- [ ] API endpoints (sessions)
- [ ] Real-time session updates (SSE or WebSockets)

### Phase 2: Sharing Detection (Week 2)
- [ ] Sharing rules schema
- [ ] Sharing detector service
- [ ] Rule CRUD endpoints
- [ ] Violation tracking

### Phase 3: Trust Scoring (Week 2)
- [ ] Trust score calculation
- [ ] Automated actions (warn, terminate, suspend)
- [ ] Trust score API endpoints

### Phase 4: Library Analytics (Week 3)
- [ ] Library stats aggregation (River jobs)
- [ ] Analytics API endpoints (overview, quality, storage, watch)
- [ ] Engagement metrics

### Phase 5: Geolocation & Map (Week 3)
- [ ] IP geolocation service (MaxMind/IP2Location)
- [ ] Stream map API endpoint
- [ ] Frontend map component (Leaflet.js)

### Phase 6: Notifications (Week 4)
- [ ] Discord webhook integration
- [ ] Email notifications (SMTP)
- [ ] Telegram notifications
- [ ] Custom webhook support

### Phase 7: Frontend Dashboard (Week 4-5)
- [ ] Session viewer component (Svelte 5)
- [ ] Violation dashboard
- [ ] Library analytics dashboard
- [ ] Stream map component
- [ ] Charts (Apache ECharts)

**Total Estimated Time**: 5-6 weeks

---

---

## Implementation Checklist

### Phase 1: Core Infrastructure
- [ ] Create package structure `internal/service/analytics/`
- [ ] Define `Session` entity with all tracking fields
- [ ] Define `SharingRule` entity with config JSONB
- [ ] Define `SharingViolation` entity
- [ ] Define `UserTrust` entity with score tracking
- [ ] Define `LibraryStat` entity for aggregated data
- [ ] Define `Engagement` entity for content metrics
- [ ] Create repository interfaces for all entities
- [ ] Register fx module `internal/service/analytics/module.go`

### Phase 2: Database
- [ ] Create migration `shared/000XXX_analytics_sessions.up.sql`
- [ ] Create `analytics_sessions` table with geolocation columns
- [ ] Create `analytics_sharing_rules` table with config JSONB
- [ ] Create `analytics_sharing_violations` table
- [ ] Create `analytics_user_trust` table with score constraints
- [ ] Create `analytics_library_stats` table with date partitioning
- [ ] Create `analytics_engagement` table with content metrics
- [ ] Add indexes for user_id, started_at, ip_address, content lookups
- [ ] Generate sqlc queries for session CRUD
- [ ] Generate queries for violation detection (concurrent sessions, geo anomalies)
- [ ] Generate queries for library stat aggregation

### Phase 3: Service Layer
- [ ] Implement `SessionTracker` for real-time session tracking
- [ ] Implement `SharingDetector` with six rule types:
  - [ ] Impossible travel detection
  - [ ] Simultaneous locations detection
  - [ ] Device velocity detection
  - [ ] Concurrent streams detection
  - [ ] Geo restriction enforcement
  - [ ] Account inactivity detection
- [ ] Implement `TrustScorer` with automatic score adjustments
- [ ] Implement `LibraryAnalyzer` for stats aggregation
- [ ] Implement `Geolocation` service (MaxMind GeoLite2 integration)
- [ ] Add caching for active sessions (Redis)
- [ ] Add caching for trust scores (Redis with TTL)

### Phase 4: Background Jobs
- [ ] Create River job `analytics.aggregate_library_stats` (daily)
- [ ] Create River job `analytics.detect_sharing_violations` (every 5 minutes)
- [ ] Create River job `analytics.cleanup_old_sessions` (daily with retention)
- [ ] Create River job for trust score decay/recovery
- [ ] Implement notification dispatch (Discord, email, Telegram webhooks)

### Phase 5: API Integration
- [ ] Add OpenAPI schema for all analytics endpoints
- [ ] Implement sessions endpoints (list, detail, active, delete, bulk-delete)
- [ ] Implement violations endpoints (list, detail, acknowledge, bulk-acknowledge)
- [ ] Implement trust score endpoints (list, detail, reset)
- [ ] Implement library analytics endpoints (overview, quality, storage, watch)
- [ ] Implement sharing rules CRUD endpoints
- [ ] Implement stream map endpoint with geolocation data
- [ ] Implement bulk actions endpoint
- [ ] Add authentication for all endpoints
- [ ] Add RBAC permissions (`analytics.sessions.view`, `analytics.rules.manage`, etc.)

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
| [Go io](https://pkg.go.dev/io) | [Local](../../../sources/go/stdlib/io.md) |
| [PostgreSQL Arrays](https://www.postgresql.org/docs/current/arrays.html) | [Local](../../../sources/database/postgresql-arrays.md) |
| [PostgreSQL JSON Functions](https://www.postgresql.org/docs/current/functions-json.html) | [Local](../../../sources/database/postgresql-json.md) |
| [River Job Queue](https://pkg.go.dev/github.com/riverqueue/river) | [Local](../../../sources/tooling/river.md) |
| [Svelte 5 Documentation](https://svelte.dev/docs/svelte/overview) | [Local](../../../sources/frontend/svelte5.md) |
| [Svelte 5 Runes](https://svelte.dev/docs/svelte/$state) | [Local](../../../sources/frontend/svelte-runes.md) |
| [SvelteKit Documentation](https://svelte.dev/docs/kit/introduction) | [Local](../../../sources/frontend/sveltekit.md) |
| [Uber fx](https://pkg.go.dev/go.uber.org/fx) | [Local](../../../sources/tooling/fx.md) |
| [pgx PostgreSQL Driver](https://pkg.go.dev/github.com/jackc/pgx/v5) | [Local](../../../sources/database/pgx.md) |
| [rueidis](https://pkg.go.dev/github.com/redis/rueidis) | [Local](../../../sources/tooling/rueidis.md) |
| [sqlc](https://docs.sqlc.dev/en/stable/) | [Local](../../../sources/database/sqlc.md) |
| [sqlc Configuration](https://docs.sqlc.dev/en/stable/reference/config.html) | [Local](../../../sources/database/sqlc-config.md) |

<!-- SOURCE-BREADCRUMBS-END -->

<!-- DESIGN-BREADCRUMBS-START -->

## Related Design Docs

> Auto-generated cross-references to related design documentation

**Category**: [Shared](INDEX.md)

### In This Section

- [Time-Based Access Controls](ACCESS_CONTROLS.md)
- [Revenge - Client Support & Device Capabilities](CLIENT_SUPPORT.md)
- [Content Rating System](CONTENT_RATING.md)
- [Revenge - Internationalization (i18n)](I18N.md)
- [Library Types](LIBRARY_TYPES.md)
- [News System](NEWS_SYSTEM.md)
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

- [User Management](../../operations/USER_MANAGEMENT.md) - User roles & permissions (pending)
- [Playback Service](../../architecture/04_PLAYER_ARCHITECTURE.md) - Session tracking integration point
- [Webhook Configuration](../../technical/WEBHOOKS.md) - Webhook patterns (pending)
- [River Job Queue](../../integrations/infrastructure/RIVER.md) - Background job processing (pending)

---

## Notes

- **TimescaleDB not required**: PostgreSQL 18 performance sufficient for analytics workload
- **River for aggregation**: Daily/hourly stats aggregation via background jobs
- **SSE for Plex**: Plex Server-Sent Events for instant session detection (no polling)
- **Websockets for live updates**: Push session updates to frontend in real-time
- **AGPL-3.0 avoidance**: Complete rebuild from scratch to use MIT/Apache-2.0 libraries only
- **Geolocation data**: Requires MaxMind GeoLite2 database (free) or IP2Location database
