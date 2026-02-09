# HLS Streaming Implementation

## Status: Phase 1 Complete (MVP)

Built: 2026-02-07

## Architecture

```
internal/playback/
  types.go              # Session, MediaType, StartPlaybackRequest, responses
  config.go             # PlaybackConfig, TranscodeConfig (in internal/config)
  service.go            # PlaybackService: start/stop/validate sessions
  session.go            # SessionManager: L1Cache-backed session store
  session_test.go
  hls/
    manifest.go         # Generate master playlist, read media playlists
    manifest_test.go
    handler.go          # Raw HTTP handler for .m3u8/.ts/.vtt serving
    handler_test.go
  transcode/
    decision.go         # Analyze MediaInfo -> remux vs transcode per profile
    decision_test.go
    profiles.go         # Quality profiles: original, 1080p, 720p, 480p
    profiles_test.go
    ffmpeg.go           # FFmpeg exec wrapper: video-only + audio rendition commands
    ffmpeg_test.go
    pipeline.go         # Manage running FFmpeg processes per session
    pipeline_test.go
  subtitle/
    extract.go          # Extract subtitle tracks -> WebVTT via FFmpeg
  jobs/
    cleanup.go          # River job: session health monitoring
  playbackfx/
    module.go           # fx.Module wiring (separate package to avoid import cycle)
```

## Key Design Decisions

### Separate Audio Renditions (NOT Muxed)

Each audio track gets its own HLS segment stream. Video is segmented separately.
This is critical for bandwidth efficiency with high-quality audio (TrueHD, DTS-HD MA).

```
Video: -map 0:v:0 -an  (video only, no audio)
Audio: -map 0:a:{N} -vn (single audio track, no video)
```

HLS.js downloads ONLY the active audio track's segments. Switching is instant
because it just starts fetching from a different rendition URL.

**Why not mux all tracks into every segment?**
- TrueHD Atmos: 18 Mbps, DTS-HD MA: 24.5 Mbps
- Muxing all tracks would multiply bandwidth usage with no benefit
- Most users only listen to one audio track at a time

### HLS-Compatible Audio Codecs (Zero Transcode)

These codecs are copied directly at original quality:
- AAC (most common)
- MP3
- AC-3 (Dolby Digital 5.1)
- E-AC-3 (Dolby Digital Plus / Atmos)

Incompatible codecs transcoded to AAC 256kbps:
- DTS, DTS-HD, DTS-HD MA
- TrueHD, TrueHD Atmos
- FLAC
- PCM

### Subtitle Handling

Subtitles are NOT segmented. Each text subtitle is extracted as a complete
WebVTT file served in full. HLS.js switches client-side with zero server
interaction. Bitmap subtitles (PGS, VobSub) are skipped for MVP.

### Transcode Decision Logic

Server generates ALL enabled profiles. Client (HLS.js) selects via ABR:
- H.264 + compatible audio -> remux (copy streams) to HLS
- H.265/AV1/VP9 -> must transcode video to H.264
- Source smaller than profile -> use source dimensions (no upscale)
- Each profile independently decides video+audio codec

### Session-as-Token Auth

Session UUID (v7) validates stream requests. No JWT in query params needed.
The session ID itself is the bearer token for the stream endpoints.

### Zero-Copy Segment Serving

`http.ServeFile()` for .ts segments = kernel-level sendfile(2), zero-copy.
Immutable cache headers (`Cache-Control: public, max-age=31536000, immutable`).

### Caching Strategy

| What | Cache Type | TTL |
|------|-----------|-----|
| Session lookup | L1Cache (otter) | session_timeout (30m) |
| FFmpeg processes | L1Cache (otter) | manual cleanup |
| Media probe results | L1Cache (otter) | 1h |
| Master playlist | L1Cache (otter) | 30m |
| Media playlist | L1Cache (otter) | 2s |
| Segment serving | http.ServeFile (OS page cache) | kernel-managed |

## Stream URL Structure

```
/api/v1/playback/stream/{sessionId}/master.m3u8             # Master playlist
/api/v1/playback/stream/{sessionId}/{profile}/index.m3u8     # Video media playlist
/api/v1/playback/stream/{sessionId}/{profile}/seg-NNNNN.ts   # Video segment
/api/v1/playback/stream/{sessionId}/audio/{track}/index.m3u8 # Audio rendition playlist
/api/v1/playback/stream/{sessionId}/audio/{track}/seg-NNNNN.ts # Audio segment
/api/v1/playback/stream/{sessionId}/subs/{track}.vtt         # Subtitle file
```

## API Endpoints (ogen-managed)

```
POST   /api/v1/playback/sessions              # Start session
GET    /api/v1/playback/sessions/{sessionId}   # Get session info
DELETE /api/v1/playback/sessions/{sessionId}   # Stop session
```

## Master Playlist Example

```m3u8
#EXTM3U

#EXT-X-MEDIA:TYPE=AUDIO,GROUP-ID="audio",NAME="English 5.1",LANGUAGE="eng",DEFAULT=YES,AUTOSELECT=YES,CHANNELS="6",URI="audio/0/index.m3u8"
#EXT-X-MEDIA:TYPE=AUDIO,GROUP-ID="audio",NAME="German Stereo",LANGUAGE="deu",DEFAULT=NO,AUTOSELECT=YES,CHANNELS="2",URI="audio/1/index.m3u8"

#EXT-X-MEDIA:TYPE=SUBTITLES,GROUP-ID="subs",NAME="English",LANGUAGE="eng",DEFAULT=YES,URI="subs/0.vtt"
#EXT-X-MEDIA:TYPE=SUBTITLES,GROUP-ID="subs",NAME="German",LANGUAGE="deu",DEFAULT=NO,URI="subs/1.vtt"

#EXT-X-STREAM-INF:BANDWIDTH=5192000,RESOLUTION=1920x1080,AUDIO="audio",SUBTITLES="subs"
original/index.m3u8
#EXT-X-STREAM-INF:BANDWIDTH=2928000,RESOLUTION=1280x720,AUDIO="audio",SUBTITLES="subs"
720p/index.m3u8
#EXT-X-STREAM-INF:BANDWIDTH=1496000,RESOLUTION=854x480,AUDIO="audio",SUBTITLES="subs"
480p/index.m3u8
```

## Files Modified (Existing)

| File | Change |
|------|--------|
| `internal/config/config.go` | Added PlaybackConfig + TranscodeConfig |
| `internal/api/handler.go` | Added playbackService field |
| `internal/api/server.go` | Added StreamHandler mux, PlaybackService wiring |
| `internal/app/module.go` | Added playbackfx.Module |
| `api/openapi/openapi.yaml` | Added playback tag, 3 endpoints, 5 schemas |
| `internal/infra/database/migrations/shared/000001_create_schemas.up.sql` | Added pg_trgm extension |

## Files Created

| File | Purpose |
|------|---------|
| `internal/playback/types.go` | Session, request/response types |
| `internal/playback/session.go` | L1Cache-backed session manager |
| `internal/playback/session_test.go` | Session manager tests |
| `internal/playback/service.go` | PlaybackService (start/stop/get) |
| `internal/playback/transcode/profiles.go` | Quality profile definitions |
| `internal/playback/transcode/profiles_test.go` | Profile tests |
| `internal/playback/transcode/decision.go` | Transcode/remux decision engine |
| `internal/playback/transcode/decision_test.go` | Decision logic tests |
| `internal/playback/transcode/ffmpeg.go` | FFmpeg command builders |
| `internal/playback/transcode/ffmpeg_test.go` | FFmpeg command tests |
| `internal/playback/transcode/pipeline.go` | FFmpeg process manager |
| `internal/playback/transcode/pipeline_test.go` | Pipeline tests |
| `internal/playback/hls/manifest.go` | HLS manifest generation |
| `internal/playback/hls/manifest_test.go` | Manifest tests |
| `internal/playback/hls/handler.go` | HTTP stream handler |
| `internal/playback/hls/handler_test.go` | Handler tests |
| `internal/playback/subtitle/extract.go` | WebVTT extraction |
| `internal/playback/jobs/cleanup.go` | River cleanup worker |
| `internal/playback/playbackfx/module.go` | fx DI module |
| `internal/api/handler_playback.go` | Ogen API handler |

## What's NOT Done Yet (Future Work)

### Phase 2: Smart Streaming & QoS

#### Client Bandwidth Detection
- Real-time bandwidth measurement via segment download timing
- Server-side bandwidth estimation endpoint: client reports download speeds
- Store per-session bandwidth history in L1Cache for trend analysis
- Initial bandwidth probe: serve a small test segment to estimate connection speed
- Use existing `cache.L1Cache` for per-user bandwidth profiles

#### Server-Side ABR (Adaptive Bitrate)
- Don't just rely on HLS.js client-side ABR — supplement with server intelligence
- Server suggests initial quality based on user's historical bandwidth
- API endpoint: `GET /api/v1/playback/sessions/{id}/quality-hint` returns recommended profile
- Profile filtering: only generate profiles the client can realistically play
- If client bandwidth < profile bitrate for 3+ consecutive segments, server can stop generating that profile

#### Client Capability Detection
- User-Agent parsing for device type (TV, mobile, desktop, tablet)
- Screen resolution detection via client-reported `devicePixelRatio` and viewport
- Codec support detection via `MediaSource.isTypeSupported()` (reported by client)
- Store device profiles in user session for profile filtering
- Don't generate 4K profiles for a 720p mobile screen

#### QoS Prioritization
- Stream segment requests (`/api/v1/playback/stream/`) get highest priority
- Implement request priority in middleware: segments > playlists > API > background jobs
- Rate limit non-streaming API requests under load, never rate limit segment serving
- FFmpeg process nice values: active streams get CPU priority over pre-buffering
- Dedicated goroutine pool for segment serving (separate from API handlers)

#### User Quality Preferences
- Per-user settings: `max_quality` (original/1080p/720p/480p), `prefer_original` (bool)
- Per-session override: StartPlaybackRequest gains `maxQuality` field
- Admin settings: global max quality, max concurrent transcodes, disable profiles
- Bandwidth cap per user (admin-configurable): prevents single user from consuming all bandwidth

#### Admin Transcoding Controls
- Max concurrent transcode sessions (global and per-user)
- Transcode queue with priority (resume > new session)
- GPU allocation limits for hardware acceleration
- Scheduled transcode windows (e.g., pre-transcode popular content overnight)

### Phase 3: Advanced Features
- [ ] Hardware acceleration (VAAPI, NVENC, QSV) — detect available GPU, auto-select
- [ ] Custom FFmpeg build (revenge-ffmpeg with all codecs) — Dockerfile + CI
- [ ] Bitmap subtitle support (PGS/VobSub -> text via OCR or burn-in)
- [ ] Direct play for native HLS sources (already segmented content)
- [ ] Seek optimization (keyframe-aligned seeking, segment pre-computation)
- [ ] Multi-server transcoding (distributed FFmpeg workers via River jobs)
- [ ] HDR -> SDR tone mapping for non-HDR displays
- [ ] Dolby Vision profile conversion

### Phase 4: Analytics & Observability
- [ ] Per-stream bandwidth tracking (Prometheus metrics per session)
- [ ] Transcode performance metrics (encode speed, queue depth, GPU utilization)
- [ ] Buffer underrun detection (client reports stalls)
- [ ] Quality switch frequency tracking (how often ABR changes quality)
- [ ] User experience score (composite: stalls + quality + startup time)
- [ ] Dashboard: active streams, bandwidth usage, transcode queue, error rates
