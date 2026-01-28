# Casting Protocols

> Stream to external devices

---

## Overview

Casting support enables playback on:
- Smart TVs
- Streaming devices (Chromecast, Fire TV)
- Game consoles
- DLNA-compatible devices

---

## Protocols

| Protocol | Type | Status |
|----------|------|--------|
| [Chromecast](CHROMECAST.md) | Google Cast | 🟡 Planned |
| [DLNA](DLNA.md) | UPnP/DLNA | 🟡 Planned |

---

## Protocol Details

### Chromecast (Google Cast)
**Google's casting protocol**

- ✅ Chromecast devices
- ✅ Android TV
- ✅ Google Home displays
- ✅ Chrome browser casting
- ⚠️ Requires Cast SDK

### DLNA/UPnP
**Universal standard**

- ✅ Smart TVs (Samsung, LG, Sony)
- ✅ Game consoles (Xbox, PlayStation)
- ✅ Media players
- ✅ No proprietary SDK needed

---

## Feature Comparison

| Feature | Chromecast | DLNA |
|---------|------------|------|
| Discovery | mDNS/DIAL | SSDP |
| Control | Sender/Receiver | UPnP AV |
| Protocols | HTTP/HTTPS | HTTP |
| Subtitles | WebVTT/TTML | SRT/WebVTT |
| Transcoding | Often needed | Device-dependent |
| Auth | OAuth | None |

---

## Architecture

```
User selects device
    ↓
Protocol-specific discovery
    ↓
Revenge sends media URL to device
    ↓
Device requests media directly
    ↓
Revenge handles playback control
```

---

## Configuration

```yaml
casting:
  enabled: true

  chromecast:
    enabled: true
    app_id: "${CAST_APP_ID}"

  dlna:
    enabled: true
    server_name: "Revenge Media Server"

    # Auto-discovery
    discovery:
      enabled: true
      interval: "30s"
```

---

## Transcoding Considerations

Cast devices have limited codec support:

| Device | Supported | May Need Transcode |
|--------|-----------|-------------------|
| Chromecast | H.264, VP8, AAC | HEVC, DTS |
| DLNA TV | Varies | Check device |
| Xbox | H.264, HEVC | AV1 |

---

## Related Documentation

- [Transcoding](../transcoding/INDEX.md)
- [Client Support](../../features/CLIENT_SUPPORT.md)
