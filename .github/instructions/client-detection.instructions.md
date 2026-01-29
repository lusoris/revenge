---
applyTo: "**/internal/service/playback/**/*.go,**/internal/service/session/**/*.go"
---

# Client Detection & Capabilities Instructions

> Guidelines for detecting client capabilities and adapting streaming quality.

## Client Capability Model

### Required Fields

```go
type ClientCapabilities struct {
    // Required identification
    DeviceID       string `json:"device_id"`
    DeviceName     string `json:"device_name"`
    DeviceType     string `json:"device_type"` // "browser", "mobile", "tv", "desktop"

    // Video capabilities
    SupportedVideoCodecs []string `json:"supported_video_codecs"`
    MaxVideoWidth        int      `json:"max_video_width"`
    MaxVideoHeight       int      `json:"max_video_height"`
    SupportsHDR          bool     `json:"supports_hdr"`

    // Audio capabilities
    SupportedAudioCodecs []string `json:"supported_audio_codecs"`
    MaxAudioChannels     int      `json:"max_audio_channels"`

    // Network
    IsRemote         bool `json:"is_remote"`
    BandwidthLimit   int  `json:"bandwidth_limit"` // kbps, user setting
}
```

## User-Agent Detection

### Device Detection

```go
// Good: Hierarchical detection
func DetectDevice(userAgent string) DeviceProfile {
    ua := strings.ToLower(userAgent)

    // Check TV platforms first (most specific)
    switch {
    case strings.Contains(ua, "tizen"):
        return profiles["samsung_tv"]
    case strings.Contains(ua, "webos"):
        return profiles["lg_tv"]
    case strings.Contains(ua, "roku"):
        return profiles["roku"]
    case strings.Contains(ua, "firetv"):
        return profiles["fire_tv"]
    case strings.Contains(ua, "androidtv"):
        return profiles["android_tv"]
    }

    // Check mobile
    switch {
    case strings.Contains(ua, "iphone"):
        return profiles["iphone"]
    case strings.Contains(ua, "ipad"):
        return profiles["ipad"]
    case strings.Contains(ua, "android"):
        return profiles["android"]
    }

    // Check browsers
    switch {
    case strings.Contains(ua, "safari") && !strings.Contains(ua, "chrome"):
        return profiles["safari"]
    case strings.Contains(ua, "firefox"):
        return profiles["firefox"]
    case strings.Contains(ua, "chrome"):
        return profiles["chrome"]
    }

    return profiles["generic"]
}
```

### Override with Client Report

```go
// Good: Trust client-reported capabilities over UA detection
func MergeCapabilities(detected, reported *ClientCapabilities) *ClientCapabilities {
    result := *detected

    if reported.MaxVideoWidth > 0 {
        result.MaxVideoWidth = min(detected.MaxVideoWidth, reported.MaxVideoWidth)
    }
    if len(reported.SupportedVideoCodecs) > 0 {
        result.SupportedVideoCodecs = intersect(
            detected.SupportedVideoCodecs,
            reported.SupportedVideoCodecs,
        )
    }
    // ... similar for other fields

    return &result
}
```

## Bandwidth Adaptation

### External Client Detection

```go
// Good: Detect external clients
func IsExternalClient(clientIP net.IP, localNetworks []net.IPNet) bool {
    for _, network := range localNetworks {
        if network.Contains(clientIP) {
            return false
        }
    }
    return true
}
```

### Bandwidth Measurement

```go
// Good: Rolling average with jitter tracking
type BandwidthTracker struct {
    samples    []BandwidthSample
    windowSize int
}

func (b *BandwidthTracker) SafeBandwidth() int {
    if len(b.samples) < 3 {
        return 0 // Not enough data
    }

    avg := b.average()
    jitter := b.standardDeviation()

    // Conservative: 70% of average minus jitter
    return int(float64(avg)*0.7) - jitter
}
```

### Quality Tier Selection

```go
// Good: Map bandwidth to quality tier
func SelectQualityTier(bandwidth int, isVideo bool) QualityTier {
    if isVideo {
        switch {
        case bandwidth >= 25000: return QualityTier4K
        case bandwidth >= 10000: return QualityTier1080p
        case bandwidth >= 5000:  return QualityTier720p
        case bandwidth >= 2000:  return QualityTier480p
        default:                 return QualityTier360p
        }
    }

    // Audio
    switch {
    case bandwidth >= 1500: return QualityTierLossless
    case bandwidth >= 400:  return QualityTier320k
    case bandwidth >= 200:  return QualityTier192k
    case bandwidth >= 150:  return QualityTier128k
    default:                return QualityTier64k
    }
}
```

## Session Management

### Session Creation

```go
// Good: Full capability tracking per session
type PlaybackSession struct {
    ID           uuid.UUID
    UserID       uuid.UUID
    DeviceID     string
    Capabilities *ClientCapabilities

    // Current playback
    MediaID      uuid.UUID
    PositionMs   int64
    State        PlaybackState

    // Quality tracking
    CurrentProfile   string
    BandwidthTracker *BandwidthTracker
    QualitySwitches  int
}
```

### Cross-Device Transfer

```go
// Good: Transfer playback between devices
func (s *SessionService) TransferPlayback(ctx context.Context, userID uuid.UUID, from, to string) error {
    fromSession, err := s.GetSession(ctx, userID, from)
    if err != nil {
        return err
    }

    // Pause source
    s.SendCommand(ctx, from, CommandPause)

    // Create target session with same position
    toSession := &PlaybackSession{
        UserID:     userID,
        DeviceID:   to,
        MediaID:    fromSession.MediaID,
        PositionMs: fromSession.PositionMs,
    }

    // Start on target
    return s.StartPlayback(ctx, toSession)
}
```

## Cast Protocols

### Chromecast Support

```go
// Good: Chromecast-specific profile
var chromecastProfile = ClientCapabilities{
    DeviceType:           "cast",
    SupportedVideoCodecs: []string{"h264"},
    SupportedAudioCodecs: []string{"aac", "mp3"},
    MaxVideoWidth:        1920,
    MaxVideoHeight:       1080,
    MaxAudioChannels:     2,
}
```

### DLNA Profiles

```go
// Good: DLNA profile negotiation
var dlnaProfiles = map[string]TranscodeProfile{
    "DLNA.ORG_PN=AVC_MP4_HP_HD_AAC": {
        VideoCodec: "h264",
        Profile:    "high",
        Level:      "4.1",
        MaxWidth:   1920,
        MaxHeight:  1080,
    },
}
```

## Testing

```go
// Good: Test different device profiles
func TestDeviceDetection(t *testing.T) {
    tests := []struct {
        userAgent string
        expected  string
    }{
        {"Mozilla/5.0 (SMART-TV; Linux; Tizen 5.0)", "samsung_tv"},
        {"Mozilla/5.0 (iPhone; CPU iPhone OS 15_0)", "iphone"},
        {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/91", "chrome"},
    }

    for _, tt := range tests {
        t.Run(tt.expected, func(t *testing.T) {
            profile := DetectDevice(tt.userAgent)
            if profile.Name != tt.expected {
                t.Errorf("got %s, want %s", profile.Name, tt.expected)
            }
        })
    }
}
```
---

## Related

- [INDEX.instructions.md](INDEX.instructions.md) - Main instruction index with all cross-references
- [ARCHITECTURE_V2.md](../../docs/dev/design/architecture/ARCHITECTURE_V2.md) - System architecture
- [BEST_PRACTICES.md](../../docs/dev/design/operations/BEST_PRACTICES.md) - Best practices
