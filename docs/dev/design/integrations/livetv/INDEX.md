# Live TV Providers

> PVR backend and custom IPTV channel integration

---

## Overview

Live TV integration provides:
- **Custom IPTV channels** from your media library (via ErsatzTV)
- Live channel streaming from PVR backends
- EPG (Electronic Program Guide)
- DVR recording
- Timeshift playback
- **Age-restricted channels** (including QAR isolation)

---

## Providers

| Provider | Type | Status | Priority |
|----------|------|--------|----------|
| [ErsatzTV](ERSATZTV.md) | Custom IPTV | 🟡 Planned | HIGH |
| [TVHeadend](TVHEADEND.md) | Full PVR | 🟡 Planned | Medium |
| [NextPVR](NEXTPVR.md) | Windows PVR | 🟡 Planned | Low |

---

## Provider Details

### ErsatzTV (PRIMARY)
**Custom IPTV channel creation from media library**

- ✅ Create custom 24/7 channels from your media
- ✅ Scheduling (shuffle, block, scripted)
- ✅ Hardware transcoding (NVENC, QSV, VAAPI)
- ✅ M3U/XMLTV export for external apps
- ✅ Plex/Jellyfin/Emby media source support
- ✅ Age-restricted channels (including QAR isolation)
- ✅ Free and open source

**Use Cases:**
- "Movie Channel" playing random movies 24/7
- "Kids Channel" with age-appropriate content only
- "80s Night" scheduled programming
- QAR channels (isolated, PIN-protected)

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
