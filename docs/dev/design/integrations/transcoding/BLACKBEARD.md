# Blackbeard Integration

> External transcoding service for Revenge

**Status**: 🟡 PLANNED
**Priority**: 🔴 CRITICAL (Phase 5 - Playback Service)
**Type**: gRPC/REST service client

---

## Overview

Blackbeard is Revenge's external transcoding service, designed to offload all video/audio transcoding to a dedicated service. This architecture ensures:
- **Separation of concerns**: Revenge handles metadata/UI, Blackbeard handles media processing
- **Scalability**: Run multiple Blackbeard instances for high load
- **Hardware flexibility**: Deploy on GPU-enabled machines
- **Resource isolation**: Transcoding doesn't impact Revenge's responsiveness

**Integration Points**:
- **Transcode requests**: Send transcode jobs to Blackbeard
- **Stream handling**: Receive transcoded streams
- **Progress tracking**: Monitor transcode progress
- **Health checks**: Monitor Blackbeard availability

---

## Architecture

```
┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐
│     Client      │     │     Revenge     │     │   Blackbeard    │
│ (Browser/App)   │     │  (Media Server) │     │  (Transcoder)   │
└────────┬────────┘     └────────┬────────┘     └────────┬────────┘
         │                       │                       │
         │ Request stream        │                       │
         │──────────────────────>│                       │
         │                       │                       │
         │                       │ Check client caps     │
         │                       │ Direct play possible? │
         │                       │                       │
         │                       │ [If transcode needed] │
         │                       │ Request transcode     │
         │                       │──────────────────────>│
         │                       │                       │
         │                       │ Stream URL/Manifest   │
         │                       │<──────────────────────│
         │                       │                       │
         │ Stream manifest       │                       │
         │<──────────────────────│                       │
         │                       │                       │
         │ Fetch segments        │                       │
         │─────────────────────────────────────────────>│
         │                       │                       │
         │ Video/audio chunks    │                       │
         │<─────────────────────────────────────────────│
         │                       │                       │
```

---

## API Details

### gRPC API (Primary)

```protobuf
syntax = "proto3";
package blackbeard.v1;

service TranscodeService {
  // Start a new transcode session
  rpc StartSession(StartSessionRequest) returns (StartSessionResponse);

  // Get session status
  rpc GetSession(GetSessionRequest) returns (SessionStatus);

  // Stop/cancel a session
  rpc StopSession(StopSessionRequest) returns (StopSessionResponse);

  // Health check
  rpc Health(HealthRequest) returns (HealthResponse);

  // Probe media file
  rpc ProbeMedia(ProbeRequest) returns (ProbeResponse);
}

message StartSessionRequest {
  string session_id = 1;
  string input_path = 2;  // File path or URL
  TranscodeProfile profile = 3;
  int64 start_time_ms = 4;  // Seek position
  map<string, string> metadata = 5;
}

message TranscodeProfile {
  VideoProfile video = 1;
  AudioProfile audio = 2;
  SubtitleProfile subtitle = 3;
  ContainerFormat container = 4;
}

message VideoProfile {
  string codec = 1;  // h264, h265, av1, vp9
  int32 width = 2;
  int32 height = 3;
  int32 bitrate_kbps = 4;
  string preset = 5;  // ultrafast, fast, medium, slow
  string hw_accel = 6;  // none, nvenc, qsv, vaapi
}

message AudioProfile {
  string codec = 1;  // aac, opus, flac
  int32 channels = 2;
  int32 bitrate_kbps = 3;
  int32 sample_rate = 4;
}

message SubtitleProfile {
  bool burn_in = 1;
  string track_index = 2;
  string external_path = 3;
}

enum ContainerFormat {
  HLS = 0;
  DASH = 1;
  MP4 = 2;
  MKV = 3;
  WEBM = 4;
}

message StartSessionResponse {
  string session_id = 1;
  string manifest_url = 2;  // HLS/DASH manifest
  string direct_url = 3;    // For progressive download
}

message SessionStatus {
  string session_id = 1;
  SessionState state = 2;
  float progress = 3;  // 0.0 - 1.0
  int64 transcoded_ms = 4;
  int64 total_ms = 5;
  string error = 6;
}

enum SessionState {
  PENDING = 0;
  PREPARING = 1;
  TRANSCODING = 2;
  READY = 3;
  COMPLETED = 4;
  ERROR = 5;
  CANCELLED = 6;
}
```

### REST API (Alternative)

| Endpoint | Method | Purpose |
|----------|--------|---------|
| `/api/v1/sessions` | POST | Start transcode session |
| `/api/v1/sessions/{id}` | GET | Get session status |
| `/api/v1/sessions/{id}` | DELETE | Stop session |
| `/api/v1/health` | GET | Health check |
| `/api/v1/probe` | POST | Probe media file |
| `/streams/{session_id}/master.m3u8` | GET | HLS manifest |
| `/streams/{session_id}/{segment}.ts` | GET | HLS segment |

---

## Implementation Checklist

- [ ] **gRPC Client** (`internal/service/playback/transcoder.go`)
  - [ ] Connection management (with retries)
  - [ ] Session lifecycle (start, status, stop)
  - [ ] Stream URL construction
  - [ ] Error handling & fallbacks

- [ ] **Session Management** (`internal/service/playback/session.go`)
  - [ ] Map user sessions to transcode sessions
  - [ ] Track active transcodes
  - [ ] Cleanup stale sessions
  - [ ] Handle session handoff (seek, quality change)

- [ ] **Profile Selection** (`internal/service/playback/profile.go`)
  - [ ] Client capability detection
  - [ ] Bandwidth-based profile selection
  - [ ] Quality ladder generation
  - [ ] Hardware acceleration detection

- [ ] **Health Monitoring** (`internal/service/playback/health.go`)
  - [ ] Periodic health checks
  - [ ] Circuit breaker for failures
  - [ ] Fallback to direct play if Blackbeard down

---

## Configuration

```yaml
# configs/config.yaml
integrations:
  blackbeard:
    enabled: true

    # Connection settings
    grpc:
      address: "blackbeard:50051"
      tls: true
      cert_file: "/etc/revenge/certs/blackbeard.crt"
      timeout: "30s"
      retry:
        max_attempts: 3
        backoff: "1s"

    # Alternative REST endpoint
    rest:
      base_url: "http://blackbeard:8080"
      timeout: "30s"

    # Session settings
    sessions:
      max_concurrent: 10
      idle_timeout: "5m"
      cleanup_interval: "1m"

    # Transcoding defaults
    defaults:
      container: "hls"
      segment_duration: 6
      video:
        codec: "h264"
        preset: "fast"
        hw_accel: "auto"  # auto, nvenc, qsv, vaapi, none
      audio:
        codec: "aac"
        channels: 2

    # Quality profiles
    profiles:
      - name: "4k"
        width: 3840
        height: 2160
        video_bitrate: 25000
        audio_bitrate: 384
      - name: "1080p"
        width: 1920
        height: 1080
        video_bitrate: 8000
        audio_bitrate: 192
      - name: "720p"
        width: 1280
        height: 720
        video_bitrate: 4000
        audio_bitrate: 128
      - name: "480p"
        width: 854
        height: 480
        video_bitrate: 1500
        audio_bitrate: 128
      - name: "360p"
        width: 640
        height: 360
        video_bitrate: 800
        audio_bitrate: 96
```

---

## Database Schema

```sql
-- Active transcode sessions
CREATE TABLE transcode_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    playback_session_id UUID REFERENCES playback_sessions(id),
    media_item_id UUID NOT NULL,
    media_item_type VARCHAR(20) NOT NULL,  -- movie, episode, etc.

    -- Blackbeard session
    blackbeard_session_id VARCHAR(100) UNIQUE,
    state VARCHAR(20) NOT NULL DEFAULT 'pending',

    -- Transcode settings
    input_path TEXT NOT NULL,
    profile_name VARCHAR(50),
    profile_settings JSONB,

    -- Progress
    start_position_ms BIGINT DEFAULT 0,
    transcoded_ms BIGINT DEFAULT 0,
    total_ms BIGINT,

    -- URLs
    manifest_url TEXT,

    -- Timestamps
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    started_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    last_activity_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_transcode_sessions_user ON transcode_sessions(user_id);
CREATE INDEX idx_transcode_sessions_state ON transcode_sessions(state);
CREATE INDEX idx_transcode_sessions_activity ON transcode_sessions(last_activity_at);
```

---

## Transcode Decision Flow

```go
func (s *PlaybackService) GetStream(ctx context.Context, req StreamRequest) (*StreamResponse, error) {
    // 1. Get media info
    media, err := s.mediaRepo.Get(ctx, req.MediaID)
    if err != nil {
        return nil, err
    }

    // 2. Get client capabilities
    client := s.detectClient(req.UserAgent, req.ClientHints)

    // 3. Check if direct play possible
    if s.canDirectPlay(media, client) {
        return &StreamResponse{
            Type: StreamTypeDirect,
            URL:  s.getDirectURL(media),
        }, nil
    }

    // 4. Check if direct stream possible (remux only)
    if s.canDirectStream(media, client) {
        return &StreamResponse{
            Type: StreamTypeDirectStream,
            URL:  s.getDirectStreamURL(media, client.Container),
        }, nil
    }

    // 5. Need transcoding - select profile
    profile := s.selectProfile(media, client, req.Bandwidth)

    // 6. Start Blackbeard session
    session, err := s.blackbeard.StartSession(ctx, BlackbeardRequest{
        InputPath:   media.FilePath,
        Profile:     profile,
        StartTimeMs: req.StartPosition,
    })
    if err != nil {
        // Fallback: try lower quality or return error
        return nil, fmt.Errorf("transcode failed: %w", err)
    }

    return &StreamResponse{
        Type:        StreamTypeTranscode,
        ManifestURL: session.ManifestURL,
        SessionID:   session.ID,
    }, nil
}
```

---

## Error Handling

| Error | Cause | Action |
|-------|-------|--------|
| `UNAVAILABLE` | Blackbeard offline | Use circuit breaker, fallback to direct play |
| `RESOURCE_EXHAUSTED` | Too many concurrent transcodes | Queue request, return 503 |
| `INVALID_ARGUMENT` | Unsupported format | Return 415, suggest alternative |
| `INTERNAL` | FFmpeg error | Log, retry with different settings |
| `DEADLINE_EXCEEDED` | Transcode too slow | Lower quality, increase timeout |

---

## Monitoring

### Metrics to Expose

```go
var (
    transcodeSessionsActive = prometheus.NewGauge(
        prometheus.GaugeOpts{
            Name: "revenge_transcode_sessions_active",
            Help: "Number of active transcode sessions",
        },
    )

    transcodeSessionDuration = prometheus.NewHistogram(
        prometheus.HistogramOpts{
            Name:    "revenge_transcode_session_duration_seconds",
            Help:    "Duration of transcode sessions",
            Buckets: []float64{1, 5, 15, 30, 60, 120, 300},
        },
    )

    blackbeardHealth = prometheus.NewGauge(
        prometheus.GaugeOpts{
            Name: "revenge_blackbeard_healthy",
            Help: "Blackbeard health status (1=healthy, 0=unhealthy)",
        },
    )
)
```

---

## Related Documentation

- [Player Architecture](../../architecture/PLAYER_ARCHITECTURE.md)
- [Streaming Best Practices](../../.github/instructions/streaming-best-practices.instructions.md)
- [Client Detection](../../.github/instructions/client-detection.instructions.md)
- [Offloading Patterns](../../technical/OFFLOADING.md)
