---
sources:
  - name: FFmpeg Documentation
    url: ../../sources/media/ffmpeg.md
    note: Auto-resolved from ffmpeg
  - name: FFmpeg Codecs
    url: ../../sources/media/ffmpeg-codecs.md
    note: Auto-resolved from ffmpeg-codecs
  - name: FFmpeg Formats
    url: ../../sources/media/ffmpeg-formats.md
    note: Auto-resolved from ffmpeg-formats
  - name: go-astiav (FFmpeg bindings)
    url: ../../sources/media/go-astiav.md
    note: Auto-resolved from go-astiav
  - name: go-astiav GitHub README
    url: ../../sources/media/go-astiav-guide.md
    note: Auto-resolved from go-astiav-docs
  - name: gohlslib (HLS)
    url: ../../sources/media/gohlslib.md
    note: Auto-resolved from gohlslib
  - name: M3U8 Extended Format
    url: ../../sources/protocols/m3u8.md
    note: Auto-resolved from m3u8
  - name: Svelte 5 Runes
    url: ../../sources/frontend/svelte-runes.md
    note: Auto-resolved from svelte-runes
  - name: Svelte 5 Documentation
    url: ../../sources/frontend/svelte5.md
    note: Auto-resolved from svelte5
  - name: SvelteKit Documentation
    url: ../../sources/frontend/sveltekit.md
    note: Auto-resolved from sveltekit
design_refs:
  - title: architecture
    path: INDEX.md
  - title: ADULT_CONTENT_SYSTEM
    path: ../features/adult/ADULT_CONTENT_SYSTEM.md
  - title: ADULT_METADATA
    path: ../features/adult/ADULT_METADATA.md
  - title: DATA_RECONCILIATION
    path: ../features/adult/DATA_RECONCILIATION.md
---

## Table of Contents

- [Revenge - Player Architecture](#revenge-player-architecture)
  - [Status](#status)
  - [Architecture](#architecture)
- [Player Architecture](#player-architecture)
    - [Components](#components)
  - [Implementation](#implementation)
    - [File Structure](#file-structure)
    - [Key Interfaces](#key-interfaces)
    - [Dependencies](#dependencies)
  - [Configuration](#configuration)
    - [Environment Variables](#environment-variables)
    - [Config Keys](#config-keys)
  - [Testing Strategy](#testing-strategy)
    - [Unit Tests](#unit-tests)
    - [Integration Tests](#integration-tests)
    - [Test Coverage](#test-coverage)
  - [Related Documentation](#related-documentation)
    - [Design Documents](#design-documents)
    - [External Sources](#external-sources)


# Revenge - Player Architecture


**Created**: 2026-01-31
**Status**: ✅ Complete
**Category**: architecture


> > Media playback system with HLS streaming and Vidstack player

Player components:
- **Backend**: gohlslib for HLS manifest generation, FFmpeg for transcoding
- **Frontend**: Vidstack player with HLS.js for adaptive streaming
- **Features**: Skip intro, trickplay thumbnails, chapter markers, subtitles
- **Casting**: Chromecast and DLNA support
- **Sync**: SyncPlay for watching together remotely


---


## Status

| Dimension | Status | Notes |
|-----------|--------|-------|
| Design | ✅ | - |
| Sources | ⚪ | - |
| Instructions | ✅ | Generated from design |
| Code | 🔴 | - |
| Linting | 🔴 | - |
| Unit Testing | 🔴 | - |
| Integration Testing | 🔴 | - |

**Overall**: ✅ Complete



---


## Architecture

# Player Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                         CLIENT LAYER                             │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐         │
│  │ Web Player   │  │ Mobile App   │  │ TV App       │         │
│  │ (Vidstack)   │  │ (React Native)│  │ (Android TV) │         │
│  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘         │
│         │                 │                  │                  │
│         └─────────────────┼──────────────────┘                  │
│                           │                                     │
│                    ┌──────▼───────┐                            │
│                    │   HLS.js     │                            │
│                    │ (Adaptive)   │                            │
│                    └──────┬───────┘                            │
└───────────────────────────┼─────────────────────────────────────┘
                            │
┌───────────────────────────▼─────────────────────────────────────┐
│                      STREAMING LAYER                             │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │              HLS Manifest Generator (gohlslib)           │  │
│  │  - Master Playlist (.m3u8)                               │  │
│  │  - Media Playlists (video/audio/subs)                    │  │
│  │  - Adaptive Bitrate Profiles                             │  │
│  └──────────────────┬───────────────────────────────────────┘  │
│                     │                                           │
│  ┌──────────────────▼───────────────────────────────────────┐  │
│  │           Transcoding Engine (FFmpeg/go-astiav)          │  │
│  │  - Codec Conversion (H.264, H.265, AV1)                  │  │
│  │  - Audio Transcoding (AAC, Opus)                         │  │
│  │  - Subtitle Burning/Extract                              │  │
│  │  - Hardware Acceleration (VAAPI, NVENC, QSV)             │  │
│  └──────────────────┬───────────────────────────────────────┘  │
│                     │                                           │
└─────────────────────┼───────────────────────────────────────────┘
                      │
┌─────────────────────▼───────────────────────────────────────────┐
│                    FEATURE LAYER                                 │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐         │
│  │ Skip Intro   │  │  Trickplay   │  │  Chapters    │         │
│  │ - Timeline   │  │  - Thumbnails│  │  - Markers   │         │
│  │ - Detection  │  │  - Scrubbing │  │  - Navigation│         │
│  └──────────────┘  └──────────────┘  └──────────────┘         │
│                                                                  │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐         │
│  │  SyncPlay    │  │  Casting     │  │  Subtitles   │         │
│  │ - WebSocket  │  │  - Chromecast│  │  - SSA/ASS   │         │
│  │ - Sync State │  │  - DLNA      │  │  - WebVTT    │         │
│  └──────────────┘  └──────────────┘  └──────────────┘         │
│                                                                  │
└──────────────────────────────────────────────────────────────────┘

┌──────────────────────────────────────────────────────────────────┐
│                       STORAGE LAYER                               │
├──────────────────────────────────────────────────────────────────┤
│                                                                   │
│  ┌──────────────────┐  ┌──────────────────┐  ┌────────────────┐ │
│  │ Media Files      │  │ Metadata DB      │  │ Cache          │ │
│  │ - Video Streams  │  │ - Playback State │  │ - Transcodes   │ │
│  │ - Audio Tracks   │  │ - User Progress  │  │ - Trickplay    │ │
│  │ - Subtitle Files │  │ - Watch History  │  │ - Segments     │ │
│  └──────────────────┘  └──────────────────┘  └────────────────┘ │
│                                                                   │
└───────────────────────────────────────────────────────────────────┘
```


### Components

<!-- Component description -->


## Implementation

### File Structure

```
internal/
  playback/
    player/
      service.go           # Player service
      repository.go        # Playback state persistence
      repository_test.go
    hls/
      manifest.go          # HLS manifest generator
      segmenter.go         # Video segmentation
      segmenter_test.go
    transcode/
      engine.go            # FFmpeg wrapper
      profiles.go          # Quality profiles
      hwaccel.go           # Hardware acceleration
      queue.go             # River job integration
      engine_test.go
    features/
      skipintro/
        detector.go        # Intro detection
        fingerprint.go     # Audio fingerprinting
        detector_test.go
      trickplay/
        generator.go       # Thumbnail generation
        sprites.go         # Sprite sheet creation
        generator_test.go
      chapters/
        extractor.go       # Chapter metadata extraction
        service.go         # Chapter management
      syncplay/
        room.go            # SyncPlay room management
        state.go           # State synchronization
        websocket.go       # WebSocket handler
        room_test.go
      casting/
        chromecast.go      # Chromecast integration
        dlna.go            # DLNA/UPnP support
      subtitles/
        converter.go       # Format conversion
        renderer.go        # Subtitle rendering
        ocr.go             # PGS OCR

web/
  src/
    lib/
      player/
        Player.svelte      # Main player component
        controls/
          Controls.svelte  # Player controls
          Timeline.svelte  # Scrub bar with trickplay
          Volume.svelte
          Settings.svelte
        overlays/
          SkipIntro.svelte # Skip intro button
          Chapters.svelte  # Chapter menu
          Subtitles.svelte # Subtitle selector
        hooks/
          usePlayback.ts   # Playback state management
          useHLS.ts        # HLS.js integration
          useSyncPlay.ts   # SyncPlay client

migrations/
  playback/
    001_playback_state.sql    # User playback progress
    002_intro_markers.sql     # Skip intro data
    003_trickplay.sql          # Thumbnail metadata
    004_chapters.sql           # Chapter markers
    005_syncplay.sql           # SyncPlay sessions
```


### Key Interfaces

```go
// Player Service Interface
type PlayerService interface {
    // Start a playback session
    StartSession(ctx context.Context, req *StartSessionRequest) (*Session, error)

    // Get HLS manifest for media
    GetHLSManifest(ctx context.Context, mediaID string, profile QualityProfile) (*HLSManifest, error)

    // Update playback position
    UpdatePosition(ctx context.Context, sessionID string, position time.Duration) error

    // Get playback state
    GetState(ctx context.Context, sessionID string) (*PlaybackState, error)

    // Stop session
    StopSession(ctx context.Context, sessionID string) error
}

type StartSessionRequest struct {
    MediaID      string
    UserID       string
    QualityProfile QualityProfile
    StartPosition  time.Duration
    AudioTrack     int
    SubtitleTrack  *int
}

type Session struct {
    ID            string
    MediaID       string
    UserID        string
    HLSManifestURL string
    Created       time.Time
    ExpiresAt     time.Time
}

type PlaybackState struct {
    SessionID     string
    Position      time.Duration
    Duration      time.Duration
    State         string  // playing, paused, buffering
    QualityProfile QualityProfile
    UpdatedAt     time.Time
}

// HLS Manifest Generator Interface
type HLSGenerator interface {
    // Generate master playlist
    GenerateMasterPlaylist(ctx context.Context, mediaID string, profiles []QualityProfile) (*MasterPlaylist, error)

    // Generate media playlist for specific quality
    GenerateMediaPlaylist(ctx context.Context, mediaID string, profile QualityProfile) (*MediaPlaylist, error)

    // Get segment URL
    GetSegmentURL(ctx context.Context, mediaID string, profile QualityProfile, segmentIndex int) (string, error)
}

type QualityProfile struct {
    Name       string
    Resolution Resolution
    Bitrate    int
    VideoCodec string
    AudioCodec string
}

type Resolution struct {
    Width  int
    Height int
}

type MasterPlaylist struct {
    Variants []Variant
    AudioGroups []AudioGroup
    SubtitleGroups []SubtitleGroup
}

// Transcoding Engine Interface
type TranscodeEngine interface {
    // Start real-time transcode
    StartTranscode(ctx context.Context, req *TranscodeRequest) (*TranscodeSession, error)

    // Stop transcode
    StopTranscode(ctx context.Context, sessionID string) error

    // Queue offline transcode job
    QueueTranscode(ctx context.Context, mediaID string, profiles []QualityProfile) ([]string, error)

    // Get transcode progress
    GetProgress(ctx context.Context, jobID string) (*TranscodeProgress, error)
}

type TranscodeRequest struct {
    MediaID    string
    Profile    QualityProfile
    StartTime  time.Duration
    HWAccel    *HWAccelConfig
    Priority   int
}

type TranscodeSession struct {
    ID         string
    MediaID    string
    Profile    QualityProfile
    OutputURL  string
    Created    time.Time
}

type HWAccelConfig struct {
    Type       HWAccelType  // vaapi, nvenc, qsv, videotoolbox
    Device     string
    MaxStreams int
}

type HWAccelType string

const (
    HWAccelVAAPI        HWAccelType = "vaapi"
    HWAccelNVENC        HWAccelType = "nvenc"
    HWAccelQSV          HWAccelType = "qsv"
    HWAccelVideoToolbox HWAccelType = "videotoolbox"
)

// Skip Intro Detection Interface
type SkipIntroDetector interface {
    // Detect intro in episode
    DetectIntro(ctx context.Context, episodeID string) (*IntroMarker, error)

    // Batch detect intros for series
    BatchDetect(ctx context.Context, seriesID string) ([]IntroMarker, error)

    // Get intro marker
    GetIntroMarker(ctx context.Context, episodeID string) (*IntroMarker, error)
}

type IntroMarker struct {
    EpisodeID        string
    IntroStart       time.Duration
    IntroEnd         time.Duration
    Confidence       float64
    DetectionMethod  string
    CreatedAt        time.Time
}

// Trickplay Generator Interface
type TrickplayGenerator interface {
    // Generate thumbnails for media
    Generate(ctx context.Context, mediaID string, opts *TrickplayOptions) (*TrickplayMetadata, error)

    // Get trickplay metadata
    GetMetadata(ctx context.Context, mediaID string) (*TrickplayMetadata, error)
}

type TrickplayOptions struct {
    Interval   time.Duration  // Thumbnail every N seconds
    Width      int
    Height     int
    TileWidth  int  // Sprites per row
    TileHeight int  // Sprites per column
}

type TrickplayMetadata struct {
    MediaID       string
    VTTFile       string
    SpriteSheets  []string
    Interval      time.Duration
    ThumbnailSize Resolution
    GeneratedAt   time.Time
}

// SyncPlay Room Manager Interface
type SyncPlayManager interface {
    // Create room
    CreateRoom(ctx context.Context, req *CreateRoomRequest) (*Room, error)

    // Join room
    JoinRoom(ctx context.Context, roomID string, userID string) error

    // Leave room
    LeaveRoom(ctx context.Context, roomID string, userID string) error

    // Sync playback command
    SyncCommand(ctx context.Context, roomID string, cmd *PlaybackCommand) error

    // Get room state
    GetRoomState(ctx context.Context, roomID string) (*SyncPlayState, error)
}

type CreateRoomRequest struct {
    Name     string
    MediaID  string
    LeaderID string
    Password *string
}

type Room struct {
    ID        string
    Name      string
    MediaID   string
    LeaderID  string
    Members   []string
    CreatedAt time.Time
}

type PlaybackCommand struct {
    Type      CommandType
    Position  *time.Duration
    Rate      *float64
    Timestamp time.Time
}

type CommandType string

const (
    CmdPlay  CommandType = "play"
    CmdPause CommandType = "pause"
    CmdSeek  CommandType = "seek"
    CmdSpeed CommandType = "speed"
)

// Subtitle Converter Interface
type SubtitleConverter interface {
    // Convert subtitle to WebVTT
    ConvertToWebVTT(ctx context.Context, input io.Reader, format SubtitleFormat) (io.Reader, error)

    // Extract PGS subtitles via OCR
    ExtractPGS(ctx context.Context, videoFile string, trackIndex int) (io.Reader, error)
}

type SubtitleFormat string

const (
    FormatSSA    SubtitleFormat = "ssa"
    FormatASS    SubtitleFormat = "ass"
    FormatSRT    SubtitleFormat = "srt"
    FormatWebVTT SubtitleFormat = "vtt"
    FormatPGS    SubtitleFormat = "pgs"
)
```


### Dependencies

{'backend': [{'name': 'github.com/bluenviron/gohlslib/v2', 'version': 'v2.0.0', 'purpose': 'HLS manifest generation and serving'}, {'name': 'github.com/asticode/go-astiav', 'version': 'v0.20.0', 'purpose': 'FFmpeg bindings for transcoding'}, {'name': 'github.com/riverqueue/river', 'version': 'v0.14.2', 'purpose': 'Transcode job queue'}, {'name': 'github.com/gorilla/websocket', 'version': 'v1.5.3', 'purpose': 'WebSocket for SyncPlay'}, {'name': 'github.com/pion/webrtc/v4', 'version': 'v4.0.5', 'purpose': 'WebRTC for future P2P streaming'}], 'frontend': [{'name': 'vidstack', 'version': '^1.12.0', 'purpose': 'Web video player framework'}, {'name': 'hls.js', 'version': '^1.5.0', 'purpose': 'HLS adaptive streaming'}, {'name': '@vidstack/react', 'version': '^1.12.0', 'purpose': 'React bindings for Vidstack (mobile)'}]}





## Configuration
### Environment Variables

<!-- Environment variables -->

### Config Keys

{'playback': [{'key': 'playback.hls.segment_duration', 'type': 'int', 'default': 6, 'description': 'HLS segment duration in seconds'}, {'key': 'playback.hls.playlist_size', 'type': 'int', 'default': 5, 'description': 'Number of segments in media playlist'}, {'key': 'playback.transcode.enabled', 'type': 'bool', 'default': True, 'description': 'Enable transcoding for unsupported formats'}, {'key': 'playback.transcode.hw_accel', 'type': 'string', 'default': 'auto', 'description': 'Hardware acceleration (auto, vaapi, nvenc, qsv, videotoolbox, none)'}, {'key': 'playback.transcode.max_concurrent', 'type': 'int', 'default': 3, 'description': 'Maximum concurrent transcode sessions'}, {'key': 'playback.transcode.profiles', 'type': 'string[]', 'default': ['1080p', '720p', '480p'], 'description': 'Quality profiles for transcoding'}, {'key': 'playback.skip_intro.enabled', 'type': 'bool', 'default': True, 'description': 'Enable skip intro detection'}, {'key': 'playback.skip_intro.auto_skip', 'type': 'bool', 'default': False, 'description': 'Automatically skip intros (user preference)'}, {'key': 'playback.trickplay.enabled', 'type': 'bool', 'default': True, 'description': 'Enable trickplay thumbnail generation'}, {'key': 'playback.trickplay.interval', 'type': 'int', 'default': 10, 'description': 'Thumbnail interval in seconds'}, {'key': 'playback.syncplay.enabled', 'type': 'bool', 'default': True, 'description': 'Enable SyncPlay feature'}, {'key': 'playback.casting.chromecast_enabled', 'type': 'bool', 'default': True, 'description': 'Enable Chromecast support'}, {'key': 'playback.casting.dlna_enabled', 'type': 'bool', 'default': True, 'description': 'Enable DLNA support'}]}




## Testing Strategy

### Unit Tests

<!-- Unit test strategy -->

### Integration Tests

<!-- Integration test strategy -->

### Test Coverage

Target: **80% minimum**







## Related Documentation
### Design Documents
- [architecture](INDEX.md)
- [ADULT_CONTENT_SYSTEM](../features/adult/ADULT_CONTENT_SYSTEM.md)
- [ADULT_METADATA](../features/adult/ADULT_METADATA.md)
- [DATA_RECONCILIATION](../features/adult/DATA_RECONCILIATION.md)

### External Sources
- [FFmpeg Documentation](../../sources/media/ffmpeg.md) - Auto-resolved from ffmpeg
- [FFmpeg Codecs](../../sources/media/ffmpeg-codecs.md) - Auto-resolved from ffmpeg-codecs
- [FFmpeg Formats](../../sources/media/ffmpeg-formats.md) - Auto-resolved from ffmpeg-formats
- [go-astiav (FFmpeg bindings)](../../sources/media/go-astiav.md) - Auto-resolved from go-astiav
- [go-astiav GitHub README](../../sources/media/go-astiav-guide.md) - Auto-resolved from go-astiav-docs
- [gohlslib (HLS)](../../sources/media/gohlslib.md) - Auto-resolved from gohlslib
- [M3U8 Extended Format](../../sources/protocols/m3u8.md) - Auto-resolved from m3u8
- [Svelte 5 Runes](../../sources/frontend/svelte-runes.md) - Auto-resolved from svelte-runes
- [Svelte 5 Documentation](../../sources/frontend/svelte5.md) - Auto-resolved from svelte5
- [SvelteKit Documentation](../../sources/frontend/sveltekit.md) - Auto-resolved from sveltekit

