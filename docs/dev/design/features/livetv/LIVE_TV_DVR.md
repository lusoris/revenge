# Live TV & DVR

> Live television streaming and digital video recording

**Status**: 🔴 PLANNING
**Priority**: 🟢 HIGH (Critical Gap - All competitors have this)
**Inspired By**: Jellyfin Live TV, Plex DVR, Emby Live TV

---

## Overview

Live TV & DVR provides live television streaming, electronic program guide (EPG), and recording capabilities through integration with TV tuners and IPTV sources.

---

## Features

### Live TV

| Feature | Description |
|---------|-------------|
| Channel Streaming | Watch live TV channels |
| EPG (Program Guide) | Browse current and upcoming programs |
| Channel Groups | Organize channels into groups |
| Channel Logos | Automatic logo fetching |
| Multiple Tuners | Support for multiple tuner sources |

### DVR (Recording)

| Feature | Description |
|---------|-------------|
| Manual Recording | Record specific time slot |
| Series Recording | Record all episodes of a series |
| Season Pass | Record specific seasons |
| Conflict Resolution | Handle tuner conflicts |
| Recording Rules | Automatic recording based on rules |
| Post-Processing | Commercial detection/removal |

---

## Supported Sources

### Tuners

| Source | Protocol | Status |
|--------|----------|--------|
| HDHomeRun | HTTP/UDP | 🟢 Primary |
| TVHeadend | HTTP | 🟢 Supported |
| NextPVR | API | 🟡 Planned |
| Plex DVR Tuner | - | ❌ N/A |

### IPTV

| Source | Format | Status |
|--------|--------|--------|
| M3U Playlists | HTTP | 🟢 Supported |
| Xtream Codes | API | 🟡 Planned |

### EPG Sources

| Source | Format | Status |
|--------|--------|--------|
| XMLTV | XML | 🟢 Supported |
| Schedules Direct | JSON | 🟡 Planned |
| TVHeadend EPG | - | 🟢 Via integration |

---

## Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                    Live TV / DVR System                         │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  ┌───────────┐   ┌───────────┐   ┌───────────┐                │
│  │ HDHomeRun │   │ TVHeadend │   │   IPTV    │                │
│  │  Tuner    │   │  Server   │   │   M3U     │                │
│  └─────┬─────┘   └─────┬─────┘   └─────┬─────┘                │
│        │               │               │                       │
│        └───────────────┴───────────────┘                       │
│                        │                                        │
│               ┌────────▼────────┐                              │
│               │  Tuner Manager  │                              │
│               └────────┬────────┘                              │
│                        │                                        │
│  ┌─────────────────────┼─────────────────────┐                 │
│  │                     │                     │                 │
│  ▼                     ▼                     ▼                 │
│ ┌──────────┐   ┌──────────────┐   ┌──────────────┐            │
│ │ Live     │   │    EPG       │   │    DVR       │            │
│ │ Stream   │   │   Manager    │   │  Scheduler   │            │
│ └──────────┘   └──────────────┘   └──────────────┘            │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

---

## Go Packages

| Package | Purpose | URL |
|---------|---------|-----|
| **gorilla/websocket** | Live updates | github.com/gorilla/websocket |
| **u2takey/ffmpeg-go** | Stream transcoding | github.com/u2takey/ffmpeg-go |
| **encoding/xml** | XMLTV parsing | stdlib |
| **robfig/cron** | Recording scheduler | github.com/robfig/cron/v3 |

---

## Database Schema

```sql
-- Tuner sources
CREATE TABLE livetv_tuners (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(200) NOT NULL,
    type VARCHAR(50) NOT NULL, -- hdhr, tvheadend, iptv
    url TEXT NOT NULL,
    api_key TEXT,

    -- Capabilities
    channel_count INT,
    tuner_count INT DEFAULT 1,

    -- Status
    is_enabled BOOLEAN DEFAULT true,
    last_scan_at TIMESTAMPTZ,
    status VARCHAR(50) DEFAULT 'unknown',

    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Channels
CREATE TABLE livetv_channels (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tuner_id UUID REFERENCES livetv_tuners(id) ON DELETE CASCADE,
    external_id VARCHAR(100) NOT NULL,

    name VARCHAR(200) NOT NULL,
    number VARCHAR(20),
    logo_url TEXT,
    logo_path TEXT, -- Local cached logo

    -- Stream info
    stream_url TEXT NOT NULL,
    stream_type VARCHAR(20), -- hls, dash, udp, http

    -- Organization
    group_name VARCHAR(100),
    is_favorite BOOLEAN DEFAULT false,
    is_hidden BOOLEAN DEFAULT false,
    sort_order INT DEFAULT 0,

    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),

    UNIQUE(tuner_id, external_id)
);

-- EPG sources
CREATE TABLE livetv_epg_sources (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(200) NOT NULL,
    type VARCHAR(50) NOT NULL, -- xmltv, schedules_direct
    url TEXT,
    api_key TEXT,

    refresh_interval_hours INT DEFAULT 24,
    last_refresh_at TIMESTAMPTZ,

    is_enabled BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Program guide
CREATE TABLE livetv_programs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    channel_id UUID REFERENCES livetv_channels(id) ON DELETE CASCADE,
    epg_source_id UUID REFERENCES livetv_epg_sources(id) ON DELETE SET NULL,
    external_id VARCHAR(200),

    title VARCHAR(500) NOT NULL,
    subtitle VARCHAR(500),
    description TEXT,
    category VARCHAR(100),
    episode_title VARCHAR(500),

    -- Episode info
    season_number INT,
    episode_number INT,

    -- Timing
    start_time TIMESTAMPTZ NOT NULL,
    end_time TIMESTAMPTZ NOT NULL,

    -- Media
    image_url TEXT,

    -- Flags
    is_new BOOLEAN DEFAULT false,
    is_live BOOLEAN DEFAULT false,
    is_premiere BOOLEAN DEFAULT false,
    is_finale BOOLEAN DEFAULT false,

    created_at TIMESTAMPTZ DEFAULT NOW(),

    UNIQUE(channel_id, start_time)
);

-- Recordings
CREATE TABLE livetv_recordings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    channel_id UUID REFERENCES livetv_channels(id) ON DELETE SET NULL,
    program_id UUID REFERENCES livetv_programs(id) ON DELETE SET NULL,

    -- What to record
    title VARCHAR(500) NOT NULL,
    description TEXT,

    -- Timing
    start_time TIMESTAMPTZ NOT NULL,
    end_time TIMESTAMPTZ NOT NULL,
    start_padding_minutes INT DEFAULT 1,
    end_padding_minutes INT DEFAULT 3,

    -- Recording settings
    quality VARCHAR(20) DEFAULT 'original',
    tuner_id UUID REFERENCES livetv_tuners(id),

    -- Status
    status VARCHAR(50) DEFAULT 'scheduled', -- scheduled, recording, completed, failed, cancelled
    file_path TEXT,
    file_size_bytes BIGINT,
    error_message TEXT,

    -- Post-processing
    commercial_detect BOOLEAN DEFAULT false,
    commercial_skip BOOLEAN DEFAULT false,
    transcode_to VARCHAR(20), -- h264, hevc, etc.

    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Series recordings (Season Pass)
CREATE TABLE livetv_series_recordings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,

    -- Match criteria
    series_name VARCHAR(500) NOT NULL,
    channel_id UUID REFERENCES livetv_channels(id), -- NULL = any channel

    -- Options
    record_new_only BOOLEAN DEFAULT true,
    keep_episodes INT, -- NULL = keep all
    priority INT DEFAULT 0,

    -- Recording settings
    quality VARCHAR(20) DEFAULT 'original',
    start_padding_minutes INT DEFAULT 1,
    end_padding_minutes INT DEFAULT 3,

    is_enabled BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_livetv_channels_tuner ON livetv_channels(tuner_id);
CREATE INDEX idx_livetv_programs_channel ON livetv_programs(channel_id);
CREATE INDEX idx_livetv_programs_time ON livetv_programs(start_time, end_time);
CREATE INDEX idx_livetv_programs_title ON livetv_programs USING gin(to_tsvector('english', title));
CREATE INDEX idx_livetv_recordings_status ON livetv_recordings(status);
CREATE INDEX idx_livetv_recordings_time ON livetv_recordings(start_time);
```

---

## River Jobs

```go
const (
    JobKindRefreshEPG       = "livetv.refresh_epg"
    JobKindScanTuners       = "livetv.scan_tuners"
    JobKindStartRecording   = "livetv.start_recording"
    JobKindStopRecording    = "livetv.stop_recording"
    JobKindPostProcess      = "livetv.post_process"
    JobKindCleanupRecordings = "livetv.cleanup_recordings"
)

type StartRecordingArgs struct {
    RecordingID uuid.UUID `json:"recording_id"`
}

type RefreshEPGArgs struct {
    SourceID uuid.UUID `json:"source_id,omitempty"` // Specific source or all
}
```

---

## Go Implementation

```go
// internal/content/livetv/

type Service struct {
    tuners    TunerRepository
    channels  ChannelRepository
    programs  ProgramRepository
    recordings RecordingRepository
    river     *river.Client[pgx.Tx]
}

type HDHomeRunClient struct {
    baseURL string
    client  *http.Client
}

func (c *HDHomeRunClient) Discover() (*TunerInfo, error) {
    resp, err := c.client.Get(c.baseURL + "/discover.json")
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var info TunerInfo
    json.NewDecoder(resp.Body).Decode(&info)
    return &info, nil
}

func (c *HDHomeRunClient) GetChannels() ([]Channel, error) {
    resp, err := c.client.Get(c.baseURL + "/lineup.json")
    // Parse channels...
}

func (c *HDHomeRunClient) GetStreamURL(channelNumber string) string {
    return fmt.Sprintf("%s/auto/v%s", c.baseURL, channelNumber)
}

type RecordingManager struct {
    activeRecordings map[uuid.UUID]*ActiveRecording
    mu               sync.RWMutex
}

type ActiveRecording struct {
    RecordingID uuid.UUID
    FFmpegCmd   *exec.Cmd
    OutputPath  string
    StartedAt   time.Time
}

func (m *RecordingManager) StartRecording(ctx context.Context, rec *Recording, streamURL string) error {
    outputPath := filepath.Join(m.recordingsDir, rec.ID.String()+".ts")

    cmd := exec.CommandContext(ctx, "ffmpeg",
        "-i", streamURL,
        "-c", "copy",
        "-f", "mpegts",
        outputPath,
    )

    if err := cmd.Start(); err != nil {
        return err
    }

    m.mu.Lock()
    m.activeRecordings[rec.ID] = &ActiveRecording{
        RecordingID: rec.ID,
        FFmpegCmd:   cmd,
        OutputPath:  outputPath,
        StartedAt:   time.Now(),
    }
    m.mu.Unlock()

    return nil
}
```

---

## API Endpoints

```
# Tuners
GET  /api/v1/livetv/tuners           # List tuners
POST /api/v1/livetv/tuners           # Add tuner
GET  /api/v1/livetv/tuners/:id       # Get tuner
DELETE /api/v1/livetv/tuners/:id     # Remove tuner
POST /api/v1/livetv/tuners/:id/scan  # Scan for channels

# Channels
GET  /api/v1/livetv/channels         # List channels
GET  /api/v1/livetv/channels/:id     # Get channel
PUT  /api/v1/livetv/channels/:id     # Update channel (favorite, hidden)
GET  /api/v1/livetv/channels/:id/stream # Get stream URL

# EPG
GET  /api/v1/livetv/epg              # Get program guide
GET  /api/v1/livetv/epg/now          # Currently airing
GET  /api/v1/livetv/programs/:id     # Get program details
POST /api/v1/livetv/epg/refresh      # Refresh EPG

# Recordings
GET  /api/v1/livetv/recordings       # List recordings
POST /api/v1/livetv/recordings       # Schedule recording
GET  /api/v1/livetv/recordings/:id   # Get recording
DELETE /api/v1/livetv/recordings/:id # Cancel/delete recording

# Series recordings
GET  /api/v1/livetv/series           # List series recordings
POST /api/v1/livetv/series           # Create series recording
PUT  /api/v1/livetv/series/:id       # Update series recording
DELETE /api/v1/livetv/series/:id     # Delete series recording

# Guide data sources
GET  /api/v1/livetv/epg-sources      # List EPG sources
POST /api/v1/livetv/epg-sources      # Add EPG source
```

---

## Configuration

```yaml
livetv:
  enabled: true

  tuners:
    scan_on_startup: true
    auto_refresh_interval: 24h

  epg:
    refresh_interval: 12h
    days_ahead: 14

  recordings:
    path: "/data/recordings"
    default_quality: original
    default_padding:
      start_minutes: 1
      end_minutes: 3
    max_concurrent: 2  # Limited by tuner count

  post_processing:
    commercial_detection: false
    transcoding: false
```

---

## RBAC Permissions

| Permission | Description |
|------------|-------------|
| `livetv.watch` | Watch live TV |
| `livetv.record` | Schedule recordings |
| `livetv.manage` | Manage tuners, channels, EPG |
| `livetv.delete_recordings` | Delete recordings |

---

## Related Documentation

- [Integrations: Live TV](../integrations/livetv/INDEX.md)
- [Transcoding](../integrations/transcoding/INDEX.md)
- [Go Packages](../architecture/GO_PACKAGES.md)
