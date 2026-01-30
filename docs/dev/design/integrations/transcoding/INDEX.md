# Transcoding Services

> External media transcoding

---

## Overview

Revenge delegates transcoding to external services:
- No internal transcoding (by design)
- Hardware acceleration offloaded
- Scalable architecture
- Separation of concerns

---

## Services

| Service | Type | Status |
|---------|------|--------|
| [Blackbeard](BLACKBEARD.md) | gRPC Transcoder | 🟡 Planned |

---

## Why External Transcoding?

**Design Decision**: Revenge does not transcode internally.

| Internal Transcoding | External Transcoding |
|---------------------|---------------------|
| ❌ High CPU/GPU usage | ✅ Offloaded |
| ❌ Complex FFmpeg management | ✅ Specialized service |
| ❌ Hardware detection | ✅ Pre-configured |
| ❌ Single point of failure | ✅ Scalable |
| ❌ Resource contention | ✅ Isolated |

---

## Architecture

```
Client requests playback
    ↓
Revenge checks client capabilities
    ↓
Direct stream? → Serve file directly
    ↓
Transcode needed? → Request from Blackbeard
    ↓
Blackbeard transcodes (HW accel)
    ↓
Stream to client
```

---

## Blackbeard Service

**Purpose-built transcoding service for Revenge**

- gRPC API for low latency
- Hardware acceleration (NVENC, QSV, VAAPI)
- Adaptive bitrate streaming
- Session management
- Subtitle burning

---

## Configuration

```yaml
playback:
  # Prefer direct play
  direct_play:
    enabled: true

  # Transcoding via Blackbeard
  transcoding:
    enabled: true
    service: blackbeard

    blackbeard:
      url: "grpc://blackbeard:50051"

      # Quality profiles
      profiles:
        - name: "1080p"
          max_bitrate: 8000
          max_width: 1920
          max_height: 1080

        - name: "720p"
          max_bitrate: 4000
          max_width: 1280
          max_height: 720
```

---

## Direct Play Priority

Revenge prioritizes direct play:

1. **Direct Play** - Native playback, no processing
2. **Direct Stream** - Remux only (container change)
3. **Transcode** - Full transcoding via Blackbeard

---

## Related Documentation

- [Playback Architecture](../../architecture/04_PLAYER_ARCHITECTURE.md)
- [Client Support](../../features/CLIENT_SUPPORT.md)
