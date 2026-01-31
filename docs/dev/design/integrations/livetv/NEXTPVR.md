# NextPVR Integration

<!-- SOURCES: gohlslib, m3u8, nextpvr -->

<!-- DESIGN: integrations/livetv, 01_ARCHITECTURE, 02_DESIGN_PRINCIPLES, 03_METADATA_SYSTEM -->


> Windows/Linux DVR software with IPTV support


<!-- TOC-START -->

## Table of Contents

- [Status](#status)
- [Overview](#overview)
- [Developer Resources](#developer-resources)
- [API Details](#api-details)
  - [Authentication Flow](#authentication-flow)
  - [Key Endpoints](#key-endpoints)
  - [Example Requests](#example-requests)
- [Data Mapping](#data-mapping)
  - [Channel Mapping](#channel-mapping)
  - [EPG Mapping](#epg-mapping)
  - [Recording Mapping](#recording-mapping)
- [Implementation Checklist](#implementation-checklist)
- [Configuration](#configuration)
- [Database Schema](#database-schema)
- [Session Management](#session-management)
- [Stream Handling](#stream-handling)
  - [HLS Stream](#hls-stream)
  - [Recording Playback](#recording-playback)
- [Error Handling](#error-handling)
- [NextPVR vs TVHeadend](#nextpvr-vs-tvheadend)
- [Sources & Cross-References](#sources-cross-references)
  - [Cross-Reference Indexes](#cross-reference-indexes)
  - [Referenced Sources](#referenced-sources)
- [Related Design Docs](#related-design-docs)
  - [In This Section](#in-this-section)
  - [Related Topics](#related-topics)
  - [Indexes](#indexes)
- [Related Documentation](#related-documentation)

<!-- TOC-END -->

**Priority**: 🟢 LOW (Phase 6 - LiveTV Module)
**Type**: REST API client

## Status

| Dimension | Status | Notes |
|-----------|--------|-------|
| Design | ✅ | Comprehensive API endpoints, data mapping, session management |
| Sources | ✅ | Wiki, API reference, GitHub linked |
| Instructions | ✅ | Detailed implementation checklist |
| Code | 🔴 |  |
| Linting | 🔴 |  |
| Unit Testing | 🔴 |  |
| Integration Testing | 🔴 |  |---

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


## Related Documentation

- [TVHeadend Integration](TVHEADEND.md)
- [Live TV Module](../../features/LIBRARY_TYPES.md)
