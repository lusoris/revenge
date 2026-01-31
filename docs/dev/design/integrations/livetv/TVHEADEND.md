# TVHeadend Integration

> Open-source TV streaming server and DVR

**Priority**: 🟢 LOW (Phase 6 - LiveTV Module)
**Type**: HTTP API + HTSP client

## Status

| Dimension | Status | Notes |
| --------- | ------ | ----- |
| Design | ✅ | Comprehensive HTTP/HTSP API spec, data mapping, database schema |
| Sources | ✅ | Docs, API reference, HTSP protocol, GitHub linked |
| Instructions | ✅ | Detailed implementation checklist |
| Code | 🔴 | |
| Linting | 🔴 | |
| Unit Testing | 🔴 | |
| Integration Testing | 🔴 | |

---

## Overview

TVHeadend is a popular open-source TV streaming server and DVR that supports various input sources (DVB, IPTV, SAT>IP). Revenge integrates with TVHeadend for:
- Live TV streaming
- Electronic Program Guide (EPG)
- DVR recording management
- Channel management

**Integration Points**:
- **HTTP API**: Configuration, channel info, EPG
- **HTSP Protocol**: Streaming and real-time data
- **Stream URLs**: Direct MPEG-TS/HLS streams
- **DVR**: Recording scheduling and management

---

## Developer Resources

- 📚 **Docs**: https://docs.tvheadend.org/
- 🔗 **API Reference**: https://docs.tvheadend.org/development/json-api/
- 🔗 **HTSP Protocol**: https://docs.tvheadend.org/development/htsp/
- 🔗 **GitHub**: https://github.com/tvheadend/tvheadend

---

## API Details

### HTTP API

**Base URL**: `http://tvheadend:9981/api`
**Authentication**: HTTP Basic Auth or Digest Auth

| Endpoint | Method | Purpose |
|----------|--------|---------|
| `/api/serverinfo` | GET | Server information |
| `/api/channel/grid` | GET | List channels |
| `/api/channel/list` | GET | Simple channel list |
| `/api/epg/events/grid` | GET | EPG events |
| `/api/epg/events/load` | GET | Load specific EPG event |
| `/api/dvr/entry/grid` | GET | DVR entries |
| `/api/dvr/entry/create` | POST | Create DVR entry |
| `/api/dvr/entry/cancel` | POST | Cancel recording |
| `/api/dvr/entry/remove` | POST | Remove recording |
| `/api/status/inputs` | GET | Input status |
| `/api/status/connections` | GET | Active connections |

### Stream URLs

```bash
# MPEG-TS stream
http://user:pass@tvheadend:9981/stream/channel/{uuid}

# HLS stream (if enabled)
http://user:pass@tvheadend:9981/stream/channel/{uuid}/playlist.m3u8

# Recording playback
http://user:pass@tvheadend:9981/dvrfile/{uuid}
```

### HTSP Protocol (Optional)

Binary protocol on port 9982 for:
- Low-latency streaming
- Real-time EPG updates
- Subscription management
- DVR notifications

---

## Data Mapping

### Channel Mapping

| TVHeadend Field | Revenge Field | Notes |
|-----------------|---------------|-------|
| `uuid` | `tvheadend_channel_id` | Channel UUID |
| `name` | `name` | Channel name |
| `number` | `channel_number` | Channel number |
| `icon_public_url` | `logo_url` | Channel logo |
| `tags[]` | `tags[]` | Channel tags/categories |
| `enabled` | `enabled` | Channel enabled status |

### EPG Mapping

| TVHeadend Field | Revenge Field | Notes |
|-----------------|---------------|-------|
| `eventId` | `tvheadend_event_id` | EPG event ID |
| `channelUuid` | `channel_id` | Channel reference |
| `title` | `title` | Program title |
| `subtitle` | `subtitle` | Episode title |
| `description` | `overview` | Program description |
| `start` | `start_time` | Unix timestamp |
| `stop` | `end_time` | Unix timestamp |
| `genre[]` | `genres[]` | DVB genre codes |
| `episodeNumber` | `episode_number` | Episode number |
| `seasonNumber` | `season_number` | Season number |

### Recording Mapping

| TVHeadend Field | Revenge Field | Notes |
|-----------------|---------------|-------|
| `uuid` | `tvheadend_recording_id` | Recording UUID |
| `channel` | `channel_id` | Channel reference |
| `title` | `title` | Recording title |
| `start` | `start_time` | Scheduled start |
| `stop` | `end_time` | Scheduled end |
| `status` | `status` | scheduled, recording, completed, failed |
| `filename` | `file_path` | Output file path |
| `filesize` | `file_size` | File size in bytes |

---

## Implementation Checklist

- [ ] **HTTP Client** (`internal/service/livetv/provider_tvheadend.go`)
  - [ ] Authentication (Basic/Digest)
  - [ ] Channel fetching
  - [ ] EPG fetching with pagination
  - [ ] Stream URL generation
  - [ ] Error handling & retries

- [ ] **HTSP Client** (Optional) (`internal/service/livetv/htsp_client.go`)
  - [ ] Connection management
  - [ ] Authentication handshake
  - [ ] Subscription handling
  - [ ] Event streaming

- [ ] **Channel Sync** (`internal/service/sync/tvheadend_channels.go`)
  - [ ] Initial channel import
  - [ ] Logo downloading
  - [ ] Channel ordering
  - [ ] Tag mapping to categories

- [ ] **EPG Sync** (`internal/service/sync/tvheadend_epg.go`)
  - [ ] Periodic EPG import
  - [ ] Delta updates
  - [ ] Genre mapping
  - [ ] Series/episode linking

- [ ] **DVR Integration** (`internal/service/livetv/dvr.go`)
  - [ ] List recordings
  - [ ] Schedule recording
  - [ ] Cancel/delete recording
  - [ ] Recording playback

---

## Configuration

```yaml
# configs/config.yaml
integrations:
  tvheadend:
    enabled: true
    base_url: "http://tvheadend:9981"
    username: "${REVENGE_TVHEADEND_USER}"
    password: "${REVENGE_TVHEADEND_PASS}"

    # Optional HTSP
    htsp:
      enabled: false
      host: "tvheadend"
      port: 9982

    sync:
      channels:
        enabled: true
        interval: "1h"
      epg:
        enabled: true
        interval: "15m"
        days_ahead: 7

    streaming:
      profile: "pass"  # TVHeadend stream profile
      prefer_hls: false
      proxy_streams: false  # Proxy through Revenge

    dvr:
      enabled: true
      default_profile: "Default Profile"
      pre_padding_minutes: 5
      post_padding_minutes: 10
```

---

## Database Schema

```sql
-- TVHeadend source configuration
CREATE TABLE livetv_sources (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    provider_type VARCHAR(20) NOT NULL,  -- tvheadend, nextpvr
    base_url TEXT NOT NULL,
    config JSONB NOT NULL DEFAULT '{}',
    enabled BOOLEAN NOT NULL DEFAULT true,
    last_sync_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Channels
CREATE TABLE livetv_channels (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    source_id UUID NOT NULL REFERENCES livetv_sources(id) ON DELETE CASCADE,
    external_id VARCHAR(100) NOT NULL,  -- TVHeadend UUID
    name VARCHAR(255) NOT NULL,
    number INTEGER,
    logo_url TEXT,
    logo_cached_path TEXT,
    tags TEXT[],
    enabled BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(source_id, external_id)
);

-- EPG Programs
CREATE TABLE livetv_programs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    channel_id UUID NOT NULL REFERENCES livetv_channels(id) ON DELETE CASCADE,
    external_id VARCHAR(100) NOT NULL,
    title VARCHAR(500) NOT NULL,
    subtitle VARCHAR(500),
    overview TEXT,
    start_time TIMESTAMPTZ NOT NULL,
    end_time TIMESTAMPTZ NOT NULL,
    genres TEXT[],
    season_number INTEGER,
    episode_number INTEGER,
    image_url TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(channel_id, external_id)
);

CREATE INDEX idx_livetv_programs_time ON livetv_programs(channel_id, start_time, end_time);
CREATE INDEX idx_livetv_programs_start ON livetv_programs(start_time);

-- Recordings
CREATE TABLE livetv_recordings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    source_id UUID NOT NULL REFERENCES livetv_sources(id),
    channel_id UUID REFERENCES livetv_channels(id),
    external_id VARCHAR(100) NOT NULL,
    title VARCHAR(500) NOT NULL,
    overview TEXT,
    start_time TIMESTAMPTZ NOT NULL,
    end_time TIMESTAMPTZ NOT NULL,
    status VARCHAR(20) NOT NULL,  -- scheduled, recording, completed, failed
    file_path TEXT,
    file_size BIGINT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(source_id, external_id)
);
```

---

## Stream Handling

### Direct Stream URL

```go
func (p *TVHeadendProvider) GetStreamURL(channelID string) string {
    // Authenticated stream URL
    return fmt.Sprintf(
        "%s/stream/channel/%s?profile=%s",
        p.baseURL, channelID, p.streamProfile,
    )
}
```

### Proxied Stream

```go
func (h *LiveTVHandler) StreamChannel(w http.ResponseWriter, r *http.Request) {
    channelID := r.PathValue("id")

    // Get upstream URL
    streamURL := h.provider.GetStreamURL(channelID)

    // Proxy the stream
    resp, err := h.client.Get(streamURL)
    if err != nil {
        http.Error(w, "Stream unavailable", 503)
        return
    }
    defer resp.Body.Close()

    w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
    io.Copy(w, resp.Body)
}
```

---

## Error Handling

| Error | Cause | Solution |
|-------|-------|----------|
| 401 Unauthorized | Invalid credentials | Check username/password |
| 404 Not Found | Invalid channel UUID | Re-sync channels |
| 503 Service Unavailable | No tuners available | Inform user, try later |
| Connection refused | TVHeadend offline | Use circuit breaker |

---


<!-- SOURCE-BREADCRUMBS-START -->

## Sources & Cross-References

> Auto-generated section linking to external documentation sources

### Cross-Reference Indexes

- [All Sources Index](../../../sources/SOURCES_INDEX.md) - Complete list of external documentation
- [Design ↔ Sources Map](../../../sources/DESIGN_CROSSREF.md) - Which docs reference which sources

<!-- SOURCE-BREADCRUMBS-END -->

<!-- DESIGN-BREADCRUMBS-START -->

## Related Design Docs

> Auto-generated cross-references to related design documentation

**Category**: [Livetv](INDEX.md)

### In This Section

- [ErsatzTV Integration](ERSATZTV.md)
- [NextPVR Integration](NEXTPVR.md)

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

## Related Documentation

- [NextPVR Integration](NEXTPVR.md)
- [Live TV Module](../../features/LIBRARY_TYPES.md)
- [Media Enhancements - Live TV](../../features/MEDIA_ENHANCEMENTS.md)
