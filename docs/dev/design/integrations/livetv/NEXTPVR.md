# NextPVR Integration

> Windows/Linux DVR software with IPTV support

**Priority**: 🟢 LOW (Phase 6 - LiveTV Module)
**Type**: REST API client

## Status

| Dimension | Status | Notes |
| --------- | ------ | ----- |
| Design | ✅ | Comprehensive API endpoints, data mapping, session management |
| Sources | ✅ | Wiki, API reference, GitHub linked |
| Instructions | ✅ | Detailed implementation checklist |
| Code | 🔴 | |
| Linting | 🔴 | |
| Unit Testing | 🔴 | |
| Integration Testing | 🔴 | |

---

## Overview

NextPVR is a personal video recorder (PVR) software for Windows and Linux that supports various TV sources. Revenge integrates with NextPVR for:
- Live TV streaming
- Electronic Program Guide (EPG)
- DVR recording management
- IPTV/HDHR support

**Integration Points**:
- **REST API**: Channel, EPG, recording management
- **Stream URLs**: Direct HTTP streaming
- **Recording playback**: Completed recording access

---

## Developer Resources

- 📚 **Wiki**: https://github.com/sub3/NextPVR/wiki
- 🔗 **API Reference**: https://github.com/sub3/NextPVR/wiki/API
- 🔗 **GitHub**: https://github.com/sub3/NextPVR

---

## API Details

**Base URL**: `http://nextpvr:8866/`
**Authentication**: API key via query parameter `?sid={session_id}` or PIN

### Authentication Flow

```bash
# 1. Initiate session
GET /service?method=session.initiate&ver=1.0&device=revenge

# Response: <sid>SESSION_ID</sid>

# 2. Login with PIN
GET /service?method=session.login&md5={md5(":PIN:")}&sid={sid}

# 3. Use sid in subsequent requests
```

### Key Endpoints

| Endpoint | Purpose |
|----------|---------|
| `session.initiate` | Start session |
| `session.login` | Authenticate |
| `channel.list` | List channels |
| `channel.icon` | Get channel icon |
| `guide.list` | Get EPG data |
| `recording.list` | List recordings |
| `recording.schedule` | Schedule recording |
| `recording.delete` | Delete recording |
| `live.m3u8` | HLS stream |

### Example Requests

```bash
# List channels
GET /service?method=channel.list&sid={sid}

# Get EPG (next 24 hours)
GET /service?method=guide.list&channel_id={id}&sid={sid}

# Get live stream
GET /live?channel={oid}&client=revenge&sid={sid}

# HLS stream
GET /live.m3u8?channel={oid}&sid={sid}
```

---

## Data Mapping

### Channel Mapping

| NextPVR Field | Revenge Field | Notes |
|---------------|---------------|-------|
| `channel_id` | `nextpvr_channel_id` | Channel ID |
| `channel_oid` | `external_id` | Unique OID |
| `channel_name` | `name` | Display name |
| `channel_number` | `channel_number` | Channel number |
| `channel_minor` | `subchannel` | Minor channel number |
| `icon` | `logo_url` | Channel icon |

### EPG Mapping

| NextPVR Field | Revenge Field | Notes |
|---------------|---------------|-------|
| `id` | `nextpvr_event_id` | EPG event ID |
| `name` | `title` | Program title |
| `desc` | `overview` | Description |
| `start` | `start_time` | Start timestamp |
| `end` | `end_time` | End timestamp |
| `subtitle` | `subtitle` | Episode title |
| `season` | `season_number` | Season |
| `episode` | `episode_number` | Episode |
| `genres` | `genres[]` | Genre list |

### Recording Mapping

| NextPVR Field | Revenge Field | Notes |
|---------------|---------------|-------|
| `id` | `nextpvr_recording_id` | Recording ID |
| `name` | `title` | Recording title |
| `desc` | `overview` | Description |
| `start_time` | `start_time` | Start time |
| `duration` | `duration_seconds` | Duration |
| `status` | `status` | Recording status |
| `file` | `file_path` | File location |

---

## Implementation Checklist

- [ ] **API Client** (`internal/service/livetv/provider_nextpvr.go`)
  - [ ] Session management (initiate, login)
  - [ ] Session refresh
  - [ ] Channel fetching
  - [ ] EPG fetching
  - [ ] Stream URL generation
  - [ ] Error handling

- [ ] **Channel Sync** (`internal/service/sync/nextpvr_channels.go`)
  - [ ] Initial channel import
  - [ ] Icon downloading
  - [ ] Channel ordering

- [ ] **EPG Sync** (`internal/service/sync/nextpvr_epg.go`)
  - [ ] Periodic EPG import
  - [ ] Multi-day scheduling
  - [ ] Genre mapping

- [ ] **DVR Integration** (`internal/service/livetv/nextpvr_dvr.go`)
  - [ ] List recordings
  - [ ] Schedule recording
  - [ ] Delete recording
  - [ ] Recording playback

---

## Configuration

```yaml
# configs/config.yaml
integrations:
  nextpvr:
    enabled: true
    base_url: "http://nextpvr:8866"
    pin: "${REVENGE_NEXTPVR_PIN}"

    sync:
      channels:
        enabled: true
        interval: "2h"
      epg:
        enabled: true
        interval: "30m"
        days_ahead: 7

    streaming:
      prefer_hls: true
      proxy_streams: false

    dvr:
      enabled: true
      pre_padding_minutes: 2
      post_padding_minutes: 5
```

---

## Database Schema

Uses shared Live TV tables from [TVHeadend Integration](TVHEADEND.md#database-schema).

---

## Session Management

NextPVR sessions expire. Handle with:

```go
type NextPVRClient struct {
    baseURL   string
    pin       string
    sessionID string
    expiresAt time.Time
    mu        sync.RWMutex
}

func (c *NextPVRClient) ensureSession(ctx context.Context) error {
    c.mu.Lock()
    defer c.mu.Unlock()

    if c.sessionID != "" && time.Now().Before(c.expiresAt) {
        return nil
    }

    // Initiate new session
    sid, err := c.initiateSession(ctx)
    if err != nil {
        return err
    }

    // Login
    if err := c.login(ctx, sid); err != nil {
        return err
    }

    c.sessionID = sid
    c.expiresAt = time.Now().Add(1 * time.Hour)
    return nil
}
```

---

## Stream Handling

### HLS Stream

```go
func (p *NextPVRProvider) GetStreamURL(channelOID string) string {
    return fmt.Sprintf(
        "%s/live.m3u8?channel=%s&client=revenge&sid=%s",
        p.baseURL, channelOID, p.sessionID,
    )
}
```

### Recording Playback

```go
func (p *NextPVRProvider) GetRecordingURL(recordingID string) string {
    return fmt.Sprintf(
        "%s/service?method=recording.stream&recording_id=%s&sid=%s",
        p.baseURL, recordingID, p.sessionID,
    )
}
```

---

## Error Handling

| Error | Cause | Solution |
|-------|-------|----------|
| Invalid session | Session expired | Re-authenticate |
| Invalid PIN | Wrong PIN | Check configuration |
| Channel not found | Invalid channel OID | Re-sync channels |
| No tuners | All tuners busy | Inform user |

---

## NextPVR vs TVHeadend

| Feature | NextPVR | TVHeadend |
|---------|---------|-----------|
| Platform | Windows/Linux | Linux |
| UI | Web + Desktop | Web |
| IPTV support | Good | Excellent |
| DVB support | Limited | Excellent |
| API | REST/XML | REST/JSON + HTSP |
| Community | Medium | Large |
| Resource usage | Medium | Low |

**Recommendation**: Use TVHeadend for Linux/DVB setups, NextPVR for Windows users.

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
| [NextPVR Documentation](https://github.com/sub3/NextPVR) | [Local](../../../sources/livetv/nextpvr.md) |

<!-- SOURCE-BREADCRUMBS-END -->

<!-- DESIGN-BREADCRUMBS-START -->

## Related Design Docs

> Auto-generated cross-references to related design documentation

**Category**: [Livetv](INDEX.md)

### In This Section

- [ErsatzTV Integration](ERSATZTV.md)
- [TVHeadend Integration](TVHEADEND.md)

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

- [TVHeadend Integration](TVHEADEND.md)
- [Live TV Module](../../features/LIBRARY_TYPES.md)
