# ErsatzTV Integration

> Custom IPTV channel creation from your media library

**Service**: ErsatzTV (https://ersatztv.org)
**API**: REST API (Swagger at `/swagger`)
**Category**: IPTV / Custom Channels
**Priority**: HIGH (Core LiveTV feature)

## Status

| Dimension | Status | Notes |
| --------- | ------ | ----- |
| Design | ✅ | Comprehensive API endpoints, age restrictions, database schema |
| Sources | ✅ | Base URL, Swagger UI, output formats documented |
| Instructions | ✅ | Phased implementation checklist |
| Code | 🔴 | |
| Linting | 🔴 | |
| Unit Testing | 🔴 | |
| Integration Testing | 🔴 | |

---

## Overview

**ErsatzTV** is an open-source platform that transforms your personal media library into live, custom TV channels. It provides scheduling, EPG generation, and hardware-accelerated streaming.

**Key Features**:
- **Custom channel creation**: Design personalized 24/7 TV channels
- **IPTV/EPG output**: M3U playlists + XMLTV guide data
- **Hardware acceleration**: NVENC, QSV, VAAPI, AMF, VideoToolbox
- **Media server integration**: Plex, Jellyfin, Emby support
- **Scheduling**: Blocks, shuffled, scripted schedules
- **Filler content**: Commercials, bumpers, pre-roll

**Use Cases**:
- Create "always-on" channels (Movie Channel, Kids Channel, etc.)
- Simulate traditional TV viewing experience
- Schedule specific content at specific times
- Create themed channels (80s Movies, Documentaries, etc.)

---

## Developer Resources

### API Documentation
- **Base URL**: `http://ersatztv:8409/api`
- **Swagger UI**: `http://ersatztv:8409/swagger`
- **Authentication**: None (local network assumed)
- **Rate Limits**: None defined

### Key Endpoints

| Endpoint | Method | Purpose |
|----------|--------|---------|
| `/api/channels` | GET | List all channels |
| `/api/channels` | POST | Create channel |
| `/api/channels/{id}` | GET/PUT/DELETE | Channel CRUD |
| `/api/channels/{id}/playout` | GET | Get playout schedule |
| `/api/schedules` | GET/POST | Schedule management |
| `/api/media/sources` | GET | List media sources |
| `/api/maintenance/empty_trash` | POST | Empty trash |
| `/api/streaming/sessions` | GET | Active transcoding sessions |

### Output Formats

| Format | URL | Purpose |
|--------|-----|---------|
| M3U Playlist | `/iptv/channels.m3u` | Channel list for IPTV clients |
| XMLTV EPG | `/iptv/guide.xml` | Electronic Program Guide |
| HLS Stream | `/iptv/channel/{number}.m3u8` | Individual channel stream |

---

## Integration Architecture

### Channel Management Flow

```
Admin creates channel in Revenge UI
        ↓
Revenge API → ErsatzTV API (create channel)
        ↓
Admin selects content (movies, shows, playlists)
        ↓
Revenge API → ErsatzTV API (add media items)
        ↓
Admin sets schedule (shuffle, block, scripted)
        ↓
ErsatzTV generates playout + EPG
        ↓
Users watch via IPTV player OR Revenge player
```

### Dual Mode Support

| Mode | Description | Use Case |
|------|-------------|----------|
| **Embedded** | Revenge player streams from ErsatzTV | In-app viewing |
| **External** | Export M3U/EPG for external IPTV apps | VLC, Kodi, etc. |

---

## Age Restriction System

### Channel Age Ratings

Channels inherit the highest rating of their content:

| Rating | Content | Access |
|--------|---------|--------|
| `G` | General | All users |
| `PG` | Parental Guidance | All users |
| `PG-13` | 13+ | Users with birthdate confirming 13+ |
| `R` | Restricted | Users with birthdate confirming 17+ |
| `NC-17` | Adults Only | 18+ verified users |
| `QAR` | Adult (QAR content) | QAR-enabled users with PIN |

### QAR Channel Isolation

Channels containing QAR (adult) content follow strict isolation:

```
QAR Channel Creation:
  1. Admin must have `qar:admin` permission
  2. Channel marked as `age_rating: qar`
  3. Channel stored in `qar.channels` table (isolated schema)
  4. API namespace: `/api/v1/legacy/livetv/channels/*`
  5. Viewing requires:
     - User has `legacy:read` scope
     - Valid QAR PIN entered
     - Session PIN not expired
```

### Channel Visibility

| User Type | G/PG | PG-13/R | NC-17 | QAR |
|-----------|------|---------|-------|-----|
| Guest | | | | |
| Child Profile | | | | |
| Teen Profile | | | | |
| Adult Profile | | | | |
| QAR-Enabled | | | | |

---

## Database Schema

### Regular Channels (`public.iptv_channels`)

```sql
CREATE TABLE public.iptv_channels (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ersatztv_id INTEGER NOT NULL,          -- ErsatzTV channel ID
    name VARCHAR(255) NOT NULL,
    number INTEGER NOT NULL,               -- Channel number
    logo_url TEXT,
    age_rating VARCHAR(10) DEFAULT 'G',
    category VARCHAR(50),                  -- Movie, TV, Music, etc.
    is_enabled BOOLEAN DEFAULT true,
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);
```

### QAR Channels (`qar.channels`)

```sql
CREATE TABLE qar.channels (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ersatztv_id INTEGER NOT NULL,
    name VARCHAR(255) NOT NULL,            -- Obfuscated name in DB
    number INTEGER NOT NULL,
    logo_url TEXT,
    category VARCHAR(50),                  -- voyage, expedition, etc.
    is_enabled BOOLEAN DEFAULT true,
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);
```

---

## Configuration

```yaml
livetv:
  enabled: true

  ersatztv:
    url: "http://ersatztv:8409"

    # Channel sync
    sync:
      interval: "5m"           # Sync channels every 5 min
      auto_import: false       # Don't auto-import channels

    # Streaming
    streaming:
      proxy: true              # Proxy streams through Revenge
      hls_segments: 4          # HLS segment count

    # EPG
    epg:
      cache_hours: 24
      enhance_metadata: true   # Match EPG to TMDb/TVDB

    # Age restrictions
    age_restrictions:
      enforce: true
      default_rating: "G"
      require_pin_for: ["NC-17", "QAR"]

# QAR-specific channel settings
legacy:
  livetv:
    enabled: false             # Separate toggle for QAR channels
    ersatztv_instance: "http://ersatztv-qar:8409"  # Separate instance
```

---

## API Endpoints (Revenge)

### Regular Channels

| Endpoint | Method | Purpose |
|----------|--------|---------|
| `/api/v1/livetv/channels` | GET | List visible channels |
| `/api/v1/livetv/channels/{id}` | GET | Channel details |
| `/api/v1/livetv/channels/{id}/stream` | GET | Get stream URL |
| `/api/v1/livetv/epg` | GET | Get EPG data |
| `/api/v1/livetv/epg/{channel_id}` | GET | Channel EPG |

### Admin Endpoints

| Endpoint | Method | Purpose |
|----------|--------|---------|
| `/api/v1/admin/livetv/channels` | POST | Create channel |
| `/api/v1/admin/livetv/channels/{id}` | PUT/DELETE | Manage channel |
| `/api/v1/admin/livetv/channels/{id}/schedule` | PUT | Set schedule |
| `/api/v1/admin/livetv/sync` | POST | Force sync with ErsatzTV |

### QAR Channels (Isolated)

| Endpoint | Method | Purpose |
|----------|--------|---------|
| `/api/v1/legacy/livetv/channels` | GET | List QAR channels |
| `/api/v1/legacy/livetv/channels/{id}/stream` | GET | QAR stream (PIN required) |

---

## Frontend Integration

### Channel Guide Component

```
┌─────────────────────────────────────────────────────────────┐
│ LIVE TV                                    [Guide] [List]   │
├─────────────────────────────────────────────────────────────┤
│ ▶ 1. Movie Channel      Now: The Matrix     Next: Aliens   │
│   2. Kids Channel       Now: Toy Story      Next: Nemo     │
│   3. Documentaries      Now: Planet Earth   Next: Cosmos   │
│ 🔒 4. Late Night        [PIN Required]                      │
├─────────────────────────────────────────────────────────────┤
│                    CURRENTLY PLAYING                        │
│  ┌──────────┐                                               │
│  │          │  The Matrix (1999)                           │
│  │  [LIVE]  │  Action, Sci-Fi • R                          │
│  │          │  Started: 8:00 PM • Ends: 10:30 PM           │
│  └──────────┘                                               │
└─────────────────────────────────────────────────────────────┘
```

### Applet Support (External Systems)

For integration with external IPTV systems, Revenge provides:

| Export | Format | URL |
|--------|--------|-----|
| M3U Export | M3U8 | `/api/v1/livetv/export/m3u?token={api_token}` |
| EPG Export | XMLTV | `/api/v1/livetv/export/epg?token={api_token}` |
| Kodi Addon | Repository | `/addons/kodi/revenge.livetv/` |

---

## Implementation Checklist

### Phase 1: Basic Integration
- [ ] ErsatzTV REST client
- [ ] Channel sync (ErsatzTV → Revenge DB)
- [ ] Basic channel list API
- [ ] Stream proxy endpoint
- [ ] EPG import and caching

### Phase 2: Channel Management
- [ ] Create/edit channels via Revenge UI
- [ ] Schedule builder UI
- [ ] Content selection (from Revenge libraries)
- [ ] Logo management

### Phase 3: Age Restrictions
- [ ] Age rating detection from content
- [ ] Channel visibility filtering by user profile
- [ ] PIN prompt for restricted channels
- [ ] QAR channel isolation

### Phase 4: External Export
- [ ] M3U export with auth tokens
- [ ] XMLTV export with full metadata
- [ ] Kodi addon (optional)

---

## River Jobs

| Job | Queue | Purpose |
|-----|-------|---------|
| `SyncErsatzTVChannels` | `livetv` | Periodic channel sync |
| `RefreshEPG` | `livetv` | Update EPG cache |
| `EnhanceEPGMetadata` | `metadata` | Match EPG to TMDb/TVDB |

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
| [ErsatzTV Documentation](https://ersatztv.org/docs/) | [Local](../../../sources/livetv/ersatztv-guide.md) |

<!-- SOURCE-BREADCRUMBS-END -->

<!-- DESIGN-BREADCRUMBS-START -->

## Related Design Docs

> Auto-generated cross-references to related design documentation

**Category**: [Livetv](INDEX.md)

### In This Section

- [NextPVR Integration](NEXTPVR.md)
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

- [TVHeadend](TVHEADEND.md) - Traditional PVR backend
- [NextPVR](NEXTPVR.md) - Windows PVR
- [LIVE_TV_DVR.md](../../features/livetv/LIVE_TV_DVR.md) - Feature specification
- [QAR Obfuscation](../../00_SOURCE_OF_TRUTH.md#qar-obfuscation-terminology) - Adult content isolation
- [Age Restrictions](../../features/shared/CONTENT_RATING.md) - Content rating system

---

## Notes

### Separate ErsatzTV Instance for QAR

For maximum isolation, deploy a separate ErsatzTV instance for QAR content:

```yaml
# docker-compose.yml
services:
  ersatztv:
    image: jasongdove/ersatztv:latest
    ports:
      - "8409:8409"
    volumes:
      - ./media:/media:ro

  ersatztv-qar:
    image: jasongdove/ersatztv:latest
    ports:
      - "8410:8409"
    volumes:
      - ./media-qar:/media:ro  # Separate QAR media path
    # Network isolated from main instance
```

### Hardware Transcoding

ErsatzTV supports hardware acceleration:

| Platform | Technology | Config |
|----------|------------|--------|
| NVIDIA | NVENC | `--runtime=nvidia` |
| Intel | QSV | Device passthrough |
| AMD | AMF/VAAPI | Device passthrough |
| Apple | VideoToolbox | Native |

### Channel Number Ranges

| Range | Content |
|-------|---------|
| 1-99 | General entertainment |
| 100-199 | Movies |
| 200-299 | TV Series |
| 300-399 | Kids |
| 400-499 | Music |
| 500-599 | Sports |
| 900-999 | QAR (hidden unless enabled) |
