# DLNA/UPnP Integration

> Universal Plug and Play streaming to compatible devices

**Status**: 🟡 PLANNED
**Priority**: 🟢 LOW (Phase 6 - Casting)
**Type**: UPnP/DLNA server + control point

---

## Overview

DLNA (Digital Living Network Alliance) integration enables streaming to TVs, game consoles, and media players that support UPnP/DLNA. Revenge acts as both:
- **Media Server (DMS)**: Expose media library to DLNA clients
- **Control Point (DMC)**: Control DLNA renderers for casting

**Integration Points**:
- **SSDP Discovery**: Find DLNA devices on network
- **ContentDirectory**: Expose media library structure
- **AVTransport**: Control playback on renderers
- **ConnectionManager**: Media format negotiation

---

## Developer Resources

- 📚 **UPnP Spec**: http://www.upnp.org/specs/av/UPnP-av-AVArchitecture-v1.pdf
- 🔗 **DLNA Guidelines**: https://spirespark.com/dlna/guidelines
- 🔗 **Go Library**: `github.com/anacrolix/dms`
- 🔗 **Go UPnP**: `github.com/huin/goupnp`

---

## Technical Details

### DLNA Device Classes

| Class | Code | Description |
|-------|------|-------------|
| Digital Media Server | DMS | Serves content (Revenge) |
| Digital Media Player | DMP | Plays content (clients) |
| Digital Media Renderer | DMR | Receives/plays pushed content |
| Digital Media Controller | DMC | Controls DMR devices (Revenge) |

### UPnP Services

| Service | Purpose |
|---------|---------|
| `ContentDirectory:1` | Browse/search media library |
| `ConnectionManager:1` | Connection/protocol negotiation |
| `AVTransport:1` | Playback control |
| `RenderingControl:1` | Volume, mute, etc. |

### Discovery (SSDP)

Multicast address: `239.255.255.250:1900`

```http
M-SEARCH * HTTP/1.1
HOST: 239.255.255.250:1900
MAN: "ssdp:discover"
MX: 3
ST: urn:schemas-upnp-org:device:MediaRenderer:1
```

---

## Implementation Checklist

### DLNA Server (DMS)

- [ ] **SSDP Announcements** (`internal/service/dlna/ssdp.go`)
  - [ ] Device announcement
  - [ ] Service announcement
  - [ ] Response to M-SEARCH
  - [ ] Periodic alive broadcasts

- [ ] **ContentDirectory** (`internal/service/dlna/content_directory.go`)
  - [ ] Browse action
  - [ ] Search action
  - [ ] Object metadata (DIDL-Lite)
  - [ ] Container hierarchy (movies, TV, music)
  - [ ] Pagination support

- [ ] **ConnectionManager** (`internal/service/dlna/connection_manager.go`)
  - [ ] GetProtocolInfo
  - [ ] PrepareForConnection
  - [ ] ConnectionComplete

- [ ] **HTTP Streaming** (`internal/service/dlna/streaming.go`)
  - [ ] Range request support
  - [ ] DLNA.ORG headers
  - [ ] Time-seek support
  - [ ] Transcoding on-the-fly

### DLNA Control Point (DMC)

- [ ] **Device Discovery** (`internal/service/dlna/discovery.go`)
  - [ ] SSDP search
  - [ ] Device description parsing
  - [ ] Service URL extraction
  - [ ] Device caching

- [ ] **AVTransport Control** (`internal/service/dlna/av_transport.go`)
  - [ ] SetAVTransportURI
  - [ ] Play, Pause, Stop
  - [ ] Seek
  - [ ] GetPositionInfo
  - [ ] GetTransportInfo

- [ ] **RenderingControl** (`internal/service/dlna/rendering_control.go`)
  - [ ] GetVolume, SetVolume
  - [ ] GetMute, SetMute

---

## Configuration

```yaml
# configs/config.yaml
integrations:
  dlna:
    enabled: true

    server:
      enabled: true
      friendly_name: "Revenge Media Server"
      uuid: ""  # Auto-generated if empty
      port: 1900

      # Transcoding for DLNA clients
      transcoding:
        enabled: true
        profile: "dlna"

    control_point:
      enabled: true
      discovery_interval: "60s"

    # Network settings
    network:
      interfaces: []  # Empty = all interfaces
      bind_address: ""
```

### DLNA Transcode Profile

```yaml
profiles:
  dlna:
    video:
      codec: "h264"
      profile: "main"
      level: "4.0"
      max_width: 1920
      max_height: 1080
    audio:
      codec: "aac"
      channels: 2
      sample_rate: 48000
    container: "mpegts"  # Or mp4
```

---

## Database Schema

```sql
-- Discovered DLNA devices
CREATE TABLE dlna_devices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    udn VARCHAR(100) NOT NULL UNIQUE,  -- Unique Device Name
    friendly_name VARCHAR(255) NOT NULL,
    device_type VARCHAR(100) NOT NULL,  -- MediaRenderer, etc.
    manufacturer VARCHAR(255),
    model_name VARCHAR(255),
    location_url TEXT NOT NULL,  -- Device description URL
    services JSONB NOT NULL DEFAULT '[]',
    capabilities JSONB,
    last_seen_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_dlna_devices_type ON dlna_devices(device_type);
CREATE INDEX idx_dlna_devices_seen ON dlna_devices(last_seen_at);

-- Active DLNA sessions
CREATE TABLE dlna_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    device_id UUID NOT NULL REFERENCES dlna_devices(id),
    user_id UUID NOT NULL REFERENCES users(id),
    media_item_id UUID NOT NULL,
    media_item_type VARCHAR(20) NOT NULL,
    transport_state VARCHAR(20) DEFAULT 'STOPPED',
    current_time_seconds INTEGER DEFAULT 0,
    duration_seconds INTEGER,
    volume INTEGER DEFAULT 100,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

---

## ContentDirectory Structure

DLNA clients browse a tree structure:

```
Root (0)
├── Movies
│   ├── Recently Added
│   ├── By Genre
│   │   ├── Action
│   │   ├── Comedy
│   │   └── ...
│   └── All Movies
├── TV Shows
│   ├── Show Name
│   │   ├── Season 1
│   │   │   ├── Episode 1
│   │   │   └── ...
│   │   └── ...
│   └── ...
├── Music
│   ├── Artists
│   ├── Albums
│   └── Playlists
└── Photos
```

---

## DIDL-Lite Response

Content is described using DIDL-Lite XML:

```xml
<DIDL-Lite xmlns="urn:schemas-upnp-org:metadata-1-0/DIDL-Lite/"
           xmlns:dc="http://purl.org/dc/elements/1.1/"
           xmlns:upnp="urn:schemas-upnp-org:metadata-1-0/upnp/">
  <item id="movie-123" parentID="movies" restricted="1">
    <dc:title>Movie Title</dc:title>
    <dc:date>2024-01-01</dc:date>
    <upnp:class>object.item.videoItem.movie</upnp:class>
    <upnp:genre>Action</upnp:genre>
    <res protocolInfo="http-get:*:video/mp4:DLNA.ORG_PN=AVC_MP4_MP_SD_AAC_MULT5"
         duration="1:45:30.000"
         resolution="1920x1080"
         size="4500000000">
      http://revenge:8096/dlna/stream/movie-123
    </res>
    <upnp:albumArtURI>http://revenge:8096/dlna/art/movie-123</upnp:albumArtURI>
  </item>
</DIDL-Lite>
```

---

## DLNA.ORG Headers

Required HTTP headers for DLNA streaming:

```http
# Protocol info
contentFeatures.dlna.org: DLNA.ORG_PN=AVC_MP4_MP_SD_AAC_MULT5;DLNA.ORG_OP=01;DLNA.ORG_FLAGS=01500000000000000000000000000000

# Transfer mode
transferMode.dlna.org: Streaming

# Content type
Content-Type: video/mp4
```

### DLNA Flags

| Flag | Meaning |
|------|---------|
| `DLNA.ORG_OP=01` | Range seeking supported |
| `DLNA.ORG_OP=10` | Time seeking supported |
| `DLNA.ORG_OP=11` | Both seeking modes |
| `DLNA.ORG_CI=0` | Not transcoded |
| `DLNA.ORG_CI=1` | Transcoded |

---

## Casting Flow

```
┌────────┐     ┌─────────┐     ┌──────────────┐
│ Client │     │ Revenge │     │ DLNA Renderer│
└───┬────┘     └────┬────┘     └──────┬───────┘
    │               │                 │
    │ Cast to device│                 │
    │──────────────>│                 │
    │               │                 │
    │               │ SetAVTransportURI
    │               │────────────────>│
    │               │                 │
    │               │ Play            │
    │               │────────────────>│
    │               │                 │
    │               │ 200 OK          │
    │               │<────────────────│
    │               │                 │
    │ Session info  │                 │
    │<──────────────│                 │
    │               │                 │
    │               │ [Renderer fetches stream]
    │               │<────────────────│
    │               │                 │
    │               │ Video data      │
    │               │────────────────>│
    │               │                 │
```

---

## Supported DLNA Devices

| Device | Protocol Support | Notes |
|--------|------------------|-------|
| Samsung TVs | Good | DLNA+ extensions |
| LG TVs | Good | WebOS has native support |
| Sony TVs | Good | Standard DLNA |
| Xbox | Limited | Requires specific profiles |
| PlayStation | Limited | Strict format requirements |
| Roku | None | Use Chromecast instead |

---

## Error Handling

| Error | Cause | Solution |
|-------|-------|----------|
| Device not found | SSDP failed | Check firewall, multicast |
| 401 Unauthorized | Device auth required | Add credentials if supported |
| 501 Not Implemented | Unsupported action | Check device capabilities |
| Playback failed | Format not supported | Transcode to compatible format |

---

## Related Documentation

- [Chromecast Integration](CHROMECAST.md)
- [Client Support](../../features/CLIENT_SUPPORT.md)
- [Player Architecture](../../architecture/PLAYER_ARCHITECTURE.md)
