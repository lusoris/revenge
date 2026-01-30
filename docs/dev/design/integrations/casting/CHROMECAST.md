# Chromecast Integration

> Google Cast protocol for streaming to Chromecast devices

**Status**: 🟡 PLANNED
**Priority**: 🟢 LOW (Phase 6 - Casting)
**Type**: CAST protocol client

---

## Overview

Chromecast integration enables streaming media from Revenge to Google Cast-compatible devices including Chromecast dongles, Android TV, Google TV, and third-party Cast-enabled devices.

**Integration Points**:
- **Device discovery**: mDNS/DIAL protocol
- **Session management**: Cast session lifecycle
- **Media playback**: Load, play, pause, seek
- **Queue management**: Play next, add to queue
- **Receiver status**: Track playback state

---

## Developer Resources

- 📚 **Cast SDK**: https://developers.google.com/cast/docs/developers
- 🔗 **Receiver SDK**: https://developers.google.com/cast/docs/caf_receiver
- 🔗 **Protocol**: https://github.com/nickoala/pychromecast (unofficial docs)
- 🔗 **Go Library**: `github.com/vishen/go-chromecast`

---

## Technical Details

### Discovery

Chromecast devices advertise via mDNS:
- Service type: `_googlecast._tcp.local`
- Port: 8009 (TLS)

### Cast Protocol

1. **TLS Connection**: Connect to device on port 8009
2. **Protobuf Messages**: Communicate via Protocol Buffers
3. **Namespaces**: Different message types for different features
4. **Channels**: Virtual channels for communication

### Key Namespaces

| Namespace | Purpose |
|-----------|---------|
| `urn:x-cast:com.google.cast.tp.connection` | Connection management |
| `urn:x-cast:com.google.cast.tp.heartbeat` | Keep-alive |
| `urn:x-cast:com.google.cast.receiver` | Receiver control |
| `urn:x-cast:com.google.cast.media` | Media playback |

---

## Implementation Checklist

- [ ] **Device Discovery** (`internal/service/casting/chromecast_discovery.go`)
  - [ ] mDNS scanner
  - [ ] Device list caching
  - [ ] Device status monitoring
  - [ ] Periodic re-discovery

- [ ] **Cast Client** (`internal/service/casting/chromecast_client.go`)
  - [ ] TLS connection management
  - [ ] Protobuf serialization
  - [ ] Heartbeat handling
  - [ ] Session management
  - [ ] Reconnection handling

- [ ] **Media Controller** (`internal/service/casting/chromecast_media.go`)
  - [ ] Load media (video, audio, photos)
  - [ ] Playback control (play, pause, stop)
  - [ ] Seek to position
  - [ ] Volume control
  - [ ] Queue management

- [ ] **Web Receiver** (`web/cast-receiver/`)
  - [ ] Custom receiver app (optional)
  - [ ] Styled Cast UI
  - [ ] Subtitle support
  - [ ] Chapter display

---

## Configuration

```yaml
# configs/config.yaml
integrations:
  chromecast:
    enabled: true

    discovery:
      enabled: true
      interval: "30s"
      networks: []  # Empty = all networks

    # Optional: Custom receiver app
    receiver:
      app_id: "CC1AD845"  # Default media receiver
      # app_id: "${REVENGE_CAST_APP_ID}"  # Custom receiver

    # Stream settings for Cast
    streaming:
      prefer_direct: true
      transcode_profile: "chromecast"  # See profiles below
      subtitle_format: "vtt"
```

### Transcode Profiles for Chromecast

```yaml
# Chromecast-compatible profile
profiles:
  chromecast:
    video:
      codec: "h264"
      profile: "high"
      level: "4.1"
      max_width: 1920
      max_height: 1080
      max_bitrate: 20000
    audio:
      codec: "aac"
      channels: 6
      max_bitrate: 512
    container: "mp4"

  chromecast_ultra:
    video:
      codec: "h265"  # Or VP9
      max_width: 3840
      max_height: 2160
      max_bitrate: 40000
    audio:
      codec: "aac"
      channels: 6
    container: "mp4"
```

---

## Database Schema

```sql
-- Discovered Cast devices
CREATE TABLE cast_devices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    device_id VARCHAR(100) NOT NULL UNIQUE,  -- Chromecast UUID
    name VARCHAR(255) NOT NULL,
    model VARCHAR(100),
    address INET NOT NULL,
    port INTEGER NOT NULL DEFAULT 8009,
    capabilities JSONB,  -- Supported features
    last_seen_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_cast_devices_seen ON cast_devices(last_seen_at);

-- Active Cast sessions
CREATE TABLE cast_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    device_id UUID NOT NULL REFERENCES cast_devices(id),
    user_id UUID NOT NULL REFERENCES users(id),
    media_item_id UUID NOT NULL,
    media_item_type VARCHAR(20) NOT NULL,
    session_id VARCHAR(100),  -- Chromecast session ID
    state VARCHAR(20) NOT NULL DEFAULT 'idle',
    current_time_ms BIGINT DEFAULT 0,
    duration_ms BIGINT,
    volume REAL DEFAULT 1.0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_cast_sessions_user ON cast_sessions(user_id);
CREATE INDEX idx_cast_sessions_device ON cast_sessions(device_id);
```

---

## API Endpoints

```yaml
# Cast API
GET  /api/v1/cast/devices          # List discovered devices
POST /api/v1/cast/devices/refresh  # Force re-discovery
GET  /api/v1/cast/devices/{id}     # Get device info

POST /api/v1/cast/play             # Start casting
{
  "device_id": "uuid",
  "media_id": "uuid",
  "media_type": "movie",
  "start_position": 0
}

POST /api/v1/cast/control          # Control playback
{
  "session_id": "uuid",
  "action": "pause|play|stop|seek",
  "seek_to": 120000  # Optional, for seek
}

GET  /api/v1/cast/sessions         # List active sessions
DELETE /api/v1/cast/sessions/{id}  # Stop casting
```

---

## Cast Flow

```
┌────────┐     ┌─────────┐     ┌────────────┐     ┌────────────┐
│ Client │     │ Revenge │     │ Blackbeard │     │ Chromecast │
└───┬────┘     └────┬────┘     └─────┬──────┘     └─────┬──────┘
    │               │                │                  │
    │ Cast to device│                │                  │
    │──────────────>│                │                  │
    │               │                │                  │
    │               │ Check direct play capability       │
    │               │────────────────────────────────────>
    │               │                │                  │
    │               │ [If transcode needed]             │
    │               │ Request transcode                 │
    │               │───────────────>│                  │
    │               │                │                  │
    │               │ Stream URL     │                  │
    │               │<───────────────│                  │
    │               │                │                  │
    │               │ LOAD media command                │
    │               │────────────────────────────────────>
    │               │                │                  │
    │               │ Session started│                  │
    │               │<────────────────────────────────────
    │               │                │                  │
    │ Session info  │                │                  │
    │<──────────────│                │                  │
    │               │                │                  │
    │               │ [Chromecast fetches stream]       │
    │               │                │                  │
    │               │                │<─────────────────│
    │               │                │                  │
```

---

## Media Load Message

```json
{
  "type": "LOAD",
  "media": {
    "contentId": "https://revenge.example.com/api/v1/stream/movie/123",
    "contentType": "video/mp4",
    "streamType": "BUFFERED",
    "metadata": {
      "metadataType": 1,
      "title": "Movie Title",
      "subtitle": "2024",
      "images": [
        {
          "url": "https://revenge.example.com/api/v1/movies/123/poster"
        }
      ]
    },
    "tracks": [
      {
        "trackId": 1,
        "type": "TEXT",
        "subtype": "SUBTITLES",
        "contentId": "https://revenge.example.com/api/v1/movies/123/subtitles/en.vtt",
        "language": "en"
      }
    ]
  },
  "currentTime": 0,
  "autoplay": true
}
```

---

## Chromecast Capabilities

| Device | Max Resolution | HDR | Codecs |
|--------|----------------|-----|--------|
| Chromecast (1st gen) | 1080p | No | H.264 |
| Chromecast (2nd gen) | 1080p | No | H.264 |
| Chromecast (3rd gen) | 1080p | No | H.264, VP8 |
| Chromecast Ultra | 4K | Yes | H.264, H.265, VP9 |
| Chromecast with Google TV | 4K | Yes | H.264, H.265, VP9, AV1 |

---

## Error Handling

| Error | Cause | Solution |
|-------|-------|----------|
| Device not found | mDNS failed | Check network, retry discovery |
| Connection failed | Device offline/busy | Retry connection |
| Media load failed | Incompatible format | Transcode to compatible format |
| Session lost | Network interruption | Attempt reconnection |

---

## Related Documentation

- [DLNA Integration](DLNA.md)
- [Client Support](../../features/CLIENT_SUPPORT.md)
- [Player Architecture](../../architecture/04_PLAYER_ARCHITECTURE.md)
