# Revenge - Client Support & Device Capabilities

> Multi-platform client support with intelligent capability detection.

## Supported Clients Overview

### Native Clients (Planned)

| Platform | Framework | Status | Features |
|----------|-----------|--------|----------|
| Web | SvelteKit | ✅ Primary | Full features |
| iOS/iPadOS | Swift/SwiftUI | 🔶 Planned | Native player, offline |
| Android | Kotlin/Jetpack | 🔶 Planned | Native player, offline |
| Android TV | Kotlin/Leanback | 🔶 Planned | D-pad navigation |
| Apple TV | Swift/tvOS | 🔶 Planned | Native player |
| Desktop | Tauri (Rust+Web) | 🔶 Planned | Cross-platform |

### Third-Party Clients

| Protocol | Clients | Support Level |
|----------|---------|---------------|
| Jellyfin API | Jellyfin apps, Infuse, Swiftfin | ✅ Compatible |
| DLNA/UPnP | Smart TVs, receivers | ✅ Supported |
| Chromecast | Google TV, Chromecast devices | ✅ Supported |
| AirPlay | Apple TV, HomePod | 🔶 Planned |

---

## Client Capability Detection

### Capability Model

```go
type ClientCapabilities struct {
    // Device identification
    DeviceID       string   `json:"device_id"`
    DeviceName     string   `json:"device_name"`
    DeviceType     string   `json:"device_type"`     // "browser", "mobile", "tv", "desktop"
    AppName        string   `json:"app_name"`
    AppVersion     string   `json:"app_version"`

    // Video capabilities
    MaxVideoWidth      int      `json:"max_video_width"`
    MaxVideoHeight     int      `json:"max_video_height"`
    MaxVideoBitrate    int      `json:"max_video_bitrate"`    // kbps
    SupportedVideoCodecs []string `json:"supported_video_codecs"` // h264, hevc, av1, vp9
    SupportsHDR        bool     `json:"supports_hdr"`
    SupportsHDR10      bool     `json:"supports_hdr10"`
    SupportsDolbyVision bool    `json:"supports_dolby_vision"`

    // Audio capabilities
    MaxAudioChannels   int      `json:"max_audio_channels"`
    MaxAudioBitrate    int      `json:"max_audio_bitrate"`    // kbps
    SupportedAudioCodecs []string `json:"supported_audio_codecs"` // aac, opus, flac, ac3, eac3, truehd, dts
    SupportsDolbyAtmos bool     `json:"supports_dolby_atmos"`
    SupportsDTS        bool     `json:"supports_dts"`

    // Container support
    SupportedContainers []string `json:"supported_containers"` // mp4, mkv, webm, hls, dash

    // Subtitle support
    SupportsTextSubtitles   bool     `json:"supports_text_subtitles"`   // SRT, VTT
    SupportsBitmapSubtitles bool     `json:"supports_bitmap_subtitles"` // PGS, VOBSUB
    SupportsEmbeddedSubtitles bool   `json:"supports_embedded_subtitles"`
    SupportedSubtitleFormats []string `json:"supported_subtitle_formats"`

    // Streaming
    SupportsDirectPlay  bool   `json:"supports_direct_play"`
    SupportsDirectStream bool  `json:"supports_direct_stream"`
    SupportsTranscoding bool   `json:"supports_transcoding"`
    SupportsHLS         bool   `json:"supports_hls"`
    SupportsDASH        bool   `json:"supports_dash"`

    // Network
    IsRemote           bool   `json:"is_remote"`          // Outside local network
    NetworkType        string `json:"network_type"`       // "wifi", "ethernet", "cellular"
    BandwidthLimit     int    `json:"bandwidth_limit"`    // User-set limit (kbps)

    // Features
    SupportsSync        bool   `json:"supports_sync"`      // Offline download
    SupportsChromecast  bool   `json:"supports_chromecast"`
    SupportsAirPlay     bool   `json:"supports_airplay"`
    SupportsWebSocket   bool   `json:"supports_websocket"`
}
```

### User-Agent Detection

```go
var deviceProfiles = map[string]ClientCapabilities{
    // Browsers
    "chrome": {
        DeviceType: "browser",
        MaxVideoWidth: 3840, MaxVideoHeight: 2160,
        SupportedVideoCodecs: []string{"h264", "vp9", "av1"},
        SupportedAudioCodecs: []string{"aac", "opus", "mp3", "flac"},
        SupportsHLS: true, SupportsDASH: true,
        SupportsTextSubtitles: true,
    },
    "firefox": {
        DeviceType: "browser",
        MaxVideoWidth: 1920, MaxVideoHeight: 1080, // No HEVC
        SupportedVideoCodecs: []string{"h264", "vp9", "av1"},
        SupportedAudioCodecs: []string{"aac", "opus", "mp3", "flac"},
        SupportsHLS: true, SupportsDASH: true,
    },
    "safari": {
        DeviceType: "browser",
        MaxVideoWidth: 3840, MaxVideoHeight: 2160,
        SupportedVideoCodecs: []string{"h264", "hevc"},
        SupportedAudioCodecs: []string{"aac", "alac", "mp3"},
        SupportsHLS: true, SupportsDASH: false,
        SupportsHDR: true,
    },

    // Mobile
    "iphone": {
        DeviceType: "mobile",
        MaxVideoWidth: 1920, MaxVideoHeight: 1080,
        SupportedVideoCodecs: []string{"h264", "hevc"},
        SupportedAudioCodecs: []string{"aac", "alac"},
        SupportsHLS: true,
        SupportsAirPlay: true,
    },
    "android": {
        DeviceType: "mobile",
        MaxVideoWidth: 1920, MaxVideoHeight: 1080,
        SupportedVideoCodecs: []string{"h264", "hevc", "vp9"},
        SupportedAudioCodecs: []string{"aac", "opus"},
        SupportsHLS: true, SupportsDASH: true,
        SupportsChromecast: true,
    },

    // TV
    "tizen": { // Samsung
        DeviceType: "tv",
        MaxVideoWidth: 3840, MaxVideoHeight: 2160,
        SupportedVideoCodecs: []string{"h264", "hevc"},
        SupportedAudioCodecs: []string{"aac", "ac3", "eac3"},
        SupportsHLS: true,
        SupportsHDR: true, SupportsHDR10: true,
    },
    "webos": { // LG
        DeviceType: "tv",
        MaxVideoWidth: 3840, MaxVideoHeight: 2160,
        SupportedVideoCodecs: []string{"h264", "hevc", "vp9"},
        SupportedAudioCodecs: []string{"aac", "ac3", "eac3", "dts"},
        SupportsHLS: true, SupportsDASH: true,
        SupportsHDR: true, SupportsDolbyVision: true,
    },
    "roku": {
        DeviceType: "tv",
        MaxVideoWidth: 3840, MaxVideoHeight: 2160,
        SupportedVideoCodecs: []string{"h264", "hevc"},
        SupportedAudioCodecs: []string{"aac", "ac3"},
        SupportsHLS: true,
    },
    "fire_tv": { // Amazon
        DeviceType: "tv",
        MaxVideoWidth: 3840, MaxVideoHeight: 2160,
        SupportedVideoCodecs: []string{"h264", "hevc"},
        SupportedAudioCodecs: []string{"aac", "ac3", "eac3"},
        SupportsHLS: true, SupportsDASH: true,
        SupportsHDR: true, SupportsHDR10: true,
    },
    "android_tv": {
        DeviceType: "tv",
        MaxVideoWidth: 3840, MaxVideoHeight: 2160,
        SupportedVideoCodecs: []string{"h264", "hevc", "vp9", "av1"},
        SupportedAudioCodecs: []string{"aac", "ac3", "eac3", "dts"},
        SupportsHLS: true, SupportsDASH: true,
        SupportsChromecast: true,
    },
}

func DetectClientCapabilities(userAgent string, reportedCaps *ClientCapabilities) *ClientCapabilities {
    // Start with detected profile
    profile := detectProfileFromUserAgent(userAgent)

    // Override with client-reported capabilities
    if reportedCaps != nil {
        mergeCapabilities(profile, reportedCaps)
    }

    return profile
}
```

---

## Chromecast Integration

### Cast SDK Integration (Frontend)

```typescript
// lib/cast/chromecast.ts
import { CastContext, CastSession } from 'chromecast-caf-sender';

class ChromecastManager {
    private context: CastContext;
    private session: CastSession | null = null;

    async initialize() {
        // Load Cast SDK
        await this.loadCastSDK();

        this.context = cast.framework.CastContext.getInstance();
        this.context.setOptions({
            receiverApplicationId: chrome.cast.media.DEFAULT_MEDIA_RECEIVER_APP_ID,
            autoJoinPolicy: chrome.cast.AutoJoinPolicy.ORIGIN_SCOPED,
        });

        // Listen for session changes
        this.context.addEventListener(
            cast.framework.CastContextEventType.SESSION_STATE_CHANGED,
            this.onSessionStateChanged.bind(this)
        );
    }

    async castMedia(mediaUrl: string, metadata: MediaMetadata) {
        if (!this.session) {
            await this.requestSession();
        }

        const mediaInfo = new chrome.cast.media.MediaInfo(mediaUrl, 'video/mp4');
        mediaInfo.metadata = new chrome.cast.media.MovieMediaMetadata();
        mediaInfo.metadata.title = metadata.title;
        mediaInfo.metadata.images = [{ url: metadata.posterUrl }];

        const request = new chrome.cast.media.LoadRequest(mediaInfo);
        request.currentTime = metadata.resumePosition || 0;

        await this.session.loadMedia(request);
    }

    async stopCasting() {
        if (this.session) {
            await this.session.endSession(true);
        }
    }
}

export const chromecast = new ChromecastManager();
```

### Backend Support

```go
// Chromecast-optimized streaming endpoint
func (h *PlaybackHandler) GetCastStream(w http.ResponseWriter, r *http.Request) {
    mediaID := r.PathValue("id")

    // Chromecast prefers MP4 container with H.264/AAC
    profile := TranscodeProfile{
        VideoCodec:     "h264",
        AudioCodec:     "aac",
        Container:      "mp4",
        MaxWidth:       1920,
        MaxHeight:      1080,
        VideoBitrate:   8000,  // 8 Mbps
        AudioBitrate:   192,   // 192 kbps
        AudioChannels:  2,     // Stereo
    }

    streamURL := h.playback.GetTranscodedStreamURL(mediaID, profile)

    json.NewEncoder(w).Encode(map[string]string{
        "url": streamURL,
        "type": "video/mp4",
    })
}
```

---

## DLNA/UPnP Support

### DLNA Server Implementation

```go
type DLNAServer struct {
    name        string
    uuid        string
    port        int
    mediaServer *MediaServer
    logger      *slog.Logger
}

func NewDLNAServer(name string, mediaServer *MediaServer) *DLNAServer {
    return &DLNAServer{
        name:        name,
        uuid:        generateUUID(),
        port:        1900,
        mediaServer: mediaServer,
    }
}

func (s *DLNAServer) Start(ctx context.Context) error {
    // SSDP discovery
    go s.runSSDPServer(ctx)

    // HTTP server for device description and content
    go s.runHTTPServer(ctx)

    return nil
}

// SSDP Discovery (UDP multicast)
func (s *DLNAServer) runSSDPServer(ctx context.Context) error {
    addr, _ := net.ResolveUDPAddr("udp", "239.255.255.250:1900")
    conn, err := net.ListenMulticastUDP("udp", nil, addr)
    if err != nil {
        return err
    }
    defer conn.Close()

    buf := make([]byte, 8192)
    for {
        select {
        case <-ctx.Done():
            return nil
        default:
            n, remoteAddr, err := conn.ReadFromUDP(buf)
            if err != nil {
                continue
            }
            s.handleSSDPRequest(buf[:n], remoteAddr)
        }
    }
}

// Device description XML
func (s *DLNAServer) deviceDescriptionXML() string {
    return fmt.Sprintf(`<?xml version="1.0"?>
<root xmlns="urn:schemas-upnp-org:device-1-0">
    <specVersion><major>1</major><minor>0</minor></specVersion>
    <device>
        <deviceType>urn:schemas-upnp-org:device:MediaServer:1</deviceType>
        <friendlyName>%s</friendlyName>
        <manufacturer>Revenge</manufacturer>
        <modelName>Revenge Media Server</modelName>
        <UDN>uuid:%s</UDN>
        <serviceList>
            <service>
                <serviceType>urn:schemas-upnp-org:service:ContentDirectory:1</serviceType>
                <serviceId>urn:upnp-org:serviceId:ContentDirectory</serviceId>
                <SCPDURL>/dlna/ContentDirectory.xml</SCPDURL>
                <controlURL>/dlna/control/ContentDirectory</controlURL>
                <eventSubURL>/dlna/event/ContentDirectory</eventSubURL>
            </service>
        </serviceList>
    </device>
</root>`, s.name, s.uuid)
}

// Content Directory Service (Browse)
func (s *DLNAServer) handleBrowse(objectID string, flag string) (string, error) {
    switch objectID {
    case "0":
        return s.buildRootContainer(), nil
    case "movies":
        return s.buildMoviesContainer(), nil
    case "music":
        return s.buildMusicContainer(), nil
    default:
        // Individual item
        return s.buildItemMetadata(objectID), nil
    }
}
```

### DLNA Profiles

```go
// DLNA-compliant transcoding profiles
var dlnaProfiles = map[string]TranscodeProfile{
    "DLNA.ORG_PN=AVC_MP4_BL_CIF15_AAC_520": {
        VideoCodec: "h264", Profile: "baseline", Level: "3.0",
        MaxWidth: 352, MaxHeight: 288,
        AudioCodec: "aac", AudioBitrate: 64,
    },
    "DLNA.ORG_PN=AVC_MP4_MP_SD_AAC_MULT5": {
        VideoCodec: "h264", Profile: "main", Level: "3.1",
        MaxWidth: 720, MaxHeight: 576,
        AudioCodec: "aac", AudioBitrate: 192,
    },
    "DLNA.ORG_PN=AVC_MP4_HP_HD_AAC": {
        VideoCodec: "h264", Profile: "high", Level: "4.1",
        MaxWidth: 1920, MaxHeight: 1080,
        AudioCodec: "aac", AudioBitrate: 320,
    },
}
```

---

## Bandwidth Adaptation

### Adaptive Bitrate for Audio

```go
type AudioBandwidthAdapter struct {
    samples       []BandwidthSample
    windowSize    int
    currentBitrate int
}

func (a *AudioBandwidthAdapter) UpdateSample(bytesReceived int64, duration time.Duration) {
    kbps := int(float64(bytesReceived*8) / duration.Seconds() / 1000)

    a.samples = append(a.samples, BandwidthSample{
        Timestamp: time.Now(),
        Kbps:      kbps,
    })

    // Keep last N samples
    if len(a.samples) > a.windowSize {
        a.samples = a.samples[1:]
    }
}

func (a *AudioBandwidthAdapter) RecommendedBitrate() int {
    if len(a.samples) == 0 {
        return 320 // Default to high quality
    }

    // Calculate average and jitter
    var sum, variance int
    for _, s := range a.samples {
        sum += s.Kbps
    }
    avg := sum / len(a.samples)

    for _, s := range a.samples {
        variance += (s.Kbps - avg) * (s.Kbps - avg)
    }
    jitter := int(math.Sqrt(float64(variance / len(a.samples))))

    // Conservative: 70% of average minus jitter
    safe := int(float64(avg)*0.7) - jitter

    // Map to audio quality tiers
    switch {
    case safe >= 320:
        return 320
    case safe >= 256:
        return 256
    case safe >= 192:
        return 192
    case safe >= 128:
        return 128
    case safe >= 96:
        return 96
    default:
        return 64
    }
}
```

### Client-Reported Bandwidth

```typescript
// Frontend bandwidth measurement
class BandwidthMonitor {
    private samples: number[] = [];
    private readonly windowSize = 10;

    async measureBandwidth(): Promise<number> {
        const start = performance.now();
        const response = await fetch('/api/v1/bandwidth-test', {
            cache: 'no-store',
        });
        const data = await response.arrayBuffer();
        const duration = performance.now() - start;

        const kbps = (data.byteLength * 8) / duration; // kbps
        this.samples.push(kbps);

        if (this.samples.length > this.windowSize) {
            this.samples.shift();
        }

        return this.getAverageBandwidth();
    }

    getAverageBandwidth(): number {
        if (this.samples.length === 0) return 0;
        return this.samples.reduce((a, b) => a + b, 0) / this.samples.length;
    }

    // Report to server
    async reportCapabilities() {
        const bandwidth = await this.measureBandwidth();
        const connection = (navigator as any).connection;

        await fetch('/api/v1/session/capabilities', {
            method: 'POST',
            body: JSON.stringify({
                measured_bandwidth_kbps: bandwidth,
                network_type: connection?.effectiveType || 'unknown',
                save_data: connection?.saveData || false,
            }),
        });
    }
}
```

---

## Session Management

### Multi-Device Sync

```go
type PlaybackSession struct {
    ID           uuid.UUID        `json:"id"`
    UserID       uuid.UUID        `json:"user_id"`
    DeviceID     string           `json:"device_id"`
    DeviceType   string           `json:"device_type"`

    // Current playback
    MediaID      uuid.UUID        `json:"media_id"`
    MediaType    string           `json:"media_type"`
    PositionMs   int64            `json:"position_ms"`
    State        PlaybackState    `json:"state"`

    // Transfer info
    CanTransfer  bool             `json:"can_transfer"`
    TransferFrom *string          `json:"transfer_from,omitempty"`

    StartedAt    time.Time        `json:"started_at"`
    UpdatedAt    time.Time        `json:"updated_at"`
}

// Transfer playback to another device
func (s *SessionService) TransferPlayback(ctx context.Context, userID uuid.UUID, fromDeviceID, toDeviceID string) error {
    // Get current session
    fromSession, err := s.GetActiveSession(ctx, userID, fromDeviceID)
    if err != nil {
        return err
    }

    // Pause on source device
    s.SendCommand(ctx, fromDeviceID, PlaybackCommand{Action: "pause"})

    // Create session on target device
    toSession := &PlaybackSession{
        ID:         uuid.New(),
        UserID:     userID,
        DeviceID:   toDeviceID,
        MediaID:    fromSession.MediaID,
        MediaType:  fromSession.MediaType,
        PositionMs: fromSession.PositionMs,
        State:      PlaybackStatePaused,
        TransferFrom: &fromDeviceID,
    }

    // Notify target device via WebSocket
    s.SendCommand(ctx, toDeviceID, PlaybackCommand{
        Action:    "load",
        MediaID:   fromSession.MediaID.String(),
        Position:  fromSession.PositionMs,
        AutoPlay:  true,
    })

    return s.SaveSession(ctx, toSession)
}
```

---

## Configuration

```yaml
# configs/config.yaml
clients:
  # Capability detection
  detection:
    use_user_agent: true
    trust_client_report: true
    cache_ttl: 24h

  # Chromecast
  chromecast:
    enabled: true
    receiver_app_id: ""  # Empty = default receiver

  # DLNA
  dlna:
    enabled: true
    server_name: "Revenge Media Server"
    advertise_interval: 30s

  # AirPlay (future)
  airplay:
    enabled: false

  # Bandwidth
  bandwidth:
    test_enabled: true
    test_file_size: 1048576  # 1MB
    measurement_interval: 30s
    adapt_quality: true

  # Default profiles by device type
  default_profiles:
    browser:
      max_video_bitrate: 20000
      max_audio_bitrate: 320
    mobile:
      max_video_bitrate: 8000
      max_audio_bitrate: 256
    tv:
      max_video_bitrate: 40000
      max_audio_bitrate: 640
```

---

## Summary

| Client Type | Direct Play | Transcoding | Casting | Offline |
|-------------|-------------|-------------|---------|---------|
| Web (Chrome) | ✅ | ✅ | Chromecast | ❌ |
| Web (Safari) | ✅ | ✅ | AirPlay | ❌ |
| iOS App | ✅ | ✅ | AirPlay | ✅ |
| Android App | ✅ | ✅ | Chromecast | ✅ |
| Smart TV | ✅ | ✅ | N/A | ❌ |
| DLNA | Transcode only | ✅ | N/A | ❌ |
