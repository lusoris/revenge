# Live TV Providers

> PVR backend integration for live TV and DVR

---

## Overview

Live TV integration provides:
- Live channel streaming
- EPG (Electronic Program Guide)
- DVR recording
- Timeshift playback

---

## Providers

| Provider | Type | Status |
|----------|------|--------|
| [TVHeadend](TVHEADEND.md) | Full PVR | 🟡 Planned |
| [NextPVR](NEXTPVR.md) | Windows PVR | 🟡 Planned |

---

## Provider Details

### TVHeadend
**Full-featured PVR backend**

- ✅ DVB/ATSC/IPTV support
- ✅ Full EPG management
- ✅ Recording and series recording
- ✅ Timeshift
- ✅ HTTP/HTSP streaming
- ✅ Free and open source

### NextPVR
**Windows-focused PVR**

- ✅ Windows-native
- ✅ Good hardware support
- ✅ EPG import
- ✅ Recording
- ⚠️ Windows only

---

## Integration Modes

### Pass-through
Revenge acts as a frontend to the PVR backend:

```
User → Revenge → TVHeadend → Tuner → Content
```

### Metadata Enhancement
Revenge enriches EPG data with additional metadata:

```
TVHeadend EPG
    ↓
Revenge matches to TMDB/TVDB
    ↓
Enhanced program info displayed
```

---

## Configuration

```yaml
livetv:
  enabled: true
  provider: tvheadend

  tvheadend:
    url: "http://tvheadend:9981"
    username: "${TVH_USERNAME}"
    password: "${TVH_PASSWORD}"

    # Streaming settings
    streaming:
      profile: "pass"  # or transcode profile

    # EPG settings
    epg:
      enhance_metadata: true
      cache_hours: 24
```

---

## Data Flow

```
TVHeadend provides:
  - Channel list
  - EPG data
  - Stream URLs
  - Recording management
    ↓
Revenge provides:
  - Unified UI
  - Enhanced metadata
  - Cross-device resume
  - Watch history
```

---

## Related Documentation

- [Video Metadata](../metadata/video/INDEX.md)
- [Transcoding](../transcoding/INDEX.md)
