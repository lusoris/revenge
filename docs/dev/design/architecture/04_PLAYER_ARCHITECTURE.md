# Revenge - Player Architecture

> Unified web player for video and audio with native streaming and transcode fallback.


<!-- TOC-START -->

## Table of Contents

- [Status](#status)
- [Implementation Checklist](#implementation-checklist)
  - [Phase 1: Playback Core](#phase-1-playback-core)
  - [Phase 2: Streaming Protocol Support](#phase-2-streaming-protocol-support)
  - [Phase 3: Audio Features](#phase-3-audio-features)
  - [Phase 4: Video Features](#phase-4-video-features)
- [Overview](#overview)
- [Design Decisions](#design-decisions)
- [Player Components](#player-components)
- [Streaming Strategy](#streaming-strategy)
  - [Decision Tree](#decision-tree)
  - [Capability Detection](#capability-detection)
- [MKV Direct Play Support](#mkv-direct-play-support)
  - [Browser Compatibility](#browser-compatibility)
  - [Why MKV Direct Play?](#why-mkv-direct-play)
  - [Implementation Strategy](#implementation-strategy)
  - [MKV MIME Types](#mkv-mime-types)
  - [Video Element Setup](#video-element-setup)
  - [Fallback Strategy](#fallback-strategy)
  - [Testing MKV Support](#testing-mkv-support)
  - [Performance Benefits](#performance-benefits)
- [HLS vs DASH](#hls-vs-dash)
  - [Comparison](#comparison)
  - [Strategy: HLS Primary, DASH Fallback](#strategy-hls-primary-dash-fallback)
  - [Error Handling & Failover](#error-handling-failover)
- [Video Player Implementation](#video-player-implementation)
  - [Library: Shaka Player + hls.js](#library-shaka-player-hlsjs)
- [Audio Player Implementation](#audio-player-implementation)
  - [Library: Web Audio API + Howler.js](#library-web-audio-api-howlerjs)
- [Synced Lyrics](#synced-lyrics)
  - [Data Model](#data-model)
  - [LRC Format Support](#lrc-format-support)
  - [Lyrics Display Component](#lyrics-display-component)
- [Audio Visualization](#audio-visualization)
  - [Canvas-based Visualizer](#canvas-based-visualizer)
  - [Fanart-based Visualization (Alternative)](#fanart-based-visualization-alternative)
- [Subtitle Support](#subtitle-support)
  - [External WebVTT](#external-webvtt)
  - [Internal Subtitle Tracks (from Container)](#internal-subtitle-tracks-from-container)
  - [Styling](#styling)
- [Quality Switching](#quality-switching)
  - [User Settings](#user-settings)
  - [Seamless Quality Switch](#seamless-quality-switch)
- [iOS Safari Limitations & Graceful Degradation](#ios-safari-limitations-graceful-degradation)
- [Performance Optimizations](#performance-optimizations)
  - [1. Web Workers for Audio Processing](#1-web-workers-for-audio-processing)
  - [2. IndexedDB for Offline Caching](#2-indexeddb-for-offline-caching)
  - [3. Request Prioritization](#3-request-prioritization)
- [Custom UI Controls](#custom-ui-controls)
  - [Control Bar Components](#control-bar-components)
- [Architecture Summary](#architecture-summary)
- [Next Steps](#next-steps)
- [Sources & Cross-References](#sources-cross-references)
  - [Cross-Reference Indexes](#cross-reference-indexes)
  - [Referenced Sources](#referenced-sources)
- [Related Design Docs](#related-design-docs)
  - [In This Section](#in-this-section)
  - [Related Topics](#related-topics)
  - [Indexes](#indexes)
- [Cross-References](#cross-references)

<!-- TOC-END -->

## Status

| Dimension | Status | Notes |
|-----------|--------|-------|
| Design | ✅ | Complete player architecture specification |
| Sources | ⚪ | N/A - internal design doc |
| Instructions | 🔴 |  |
| Code | 🔴 | Reset to template |
| Linting | 🔴 |  |
| Unit Testing | 🔴 |  |
| Integration Testing | 🔴 |  |
**Priority**: 🟡 MEDIUM
**Module**: `web/src/lib/player`
**Dependencies**: [00_SOURCE_OF_TRUTH.md](../00_SOURCE_OF_TRUTH.md)

---

## Implementation Checklist

### Phase 1: Playback Core
- [ ] Define playback session types (video, audio, live)
- [ ] Implement session management and state tracking
- [ ] Add progress tracking and history recording
- [ ] Set up capability detection (codecs, containers, features)

### Phase 2: Streaming Protocol Support
- [ ] Implement direct play detection
- [ ] Add direct stream support (remuxing)
- [ ] Implement HLS playback (hls.js + Safari native)
- [ ] Implement DASH playback (Shaka Player)
- [ ] Set up protocol fallover (HLS → DASH → Progressive)

### Phase 3: Audio Features
- [ ] Implement Web Audio API integration
- [ ] Add gapless playback (30s prefetch)
- [ ] Implement true crossfade (overlapping gain nodes)
- [ ] Add synced lyrics (LRC format)
- [ ] Implement audio visualization (Canvas + fanart)

### Phase 4: Video Features
- [ ] Add subtitle support (WebVTT + container extraction)
- [ ] Implement quality switching (seamless via WebSocket)
- [ ] Add fullscreen support
- [ ] Implement iOS graceful degradation

---

## Overview

Revenge's WebUI uses a custom-built unified player that handles:
- **Video**: Movies, TV episodes, Live TV
- **Audio**: Music tracks, audiobooks, podcasts
- **Direct Play**: Native browser playback when codecs match
- **Adaptive Streaming**: HLS/DASH with quality switching
- **Advanced Features**: Gapless, crossfade, synced lyrics, visualizations

---

## Design Decisions

| Decision | Choice | Rationale |
|----------|--------|-----------|
| **Player Architecture** | Unified player with mode switching | One codebase, consistent UI, optimal for each content type |
| **Streaming Protocol** | HLS primary, DASH fallback | HLS = Safari native, DASH = modern codec support |
| **Quality Switching** | Seamless on-the-fly via Blackbeard | Best UX, no playback interruption |
| **Gapless Audio** | Intelligent prefetch + Web Audio API | Preload 30s before end, cache adjacent tracks |
| **Crossfade** | True overlapping crossfade | Dual gain nodes, professional audio experience |
| **Subtitles** | External WebVTT + internal container tracks | Blackbeard extracts, browser renders |
| **Controls** | Custom UI (SvelteKit 2 components) | Full control, accessibility built-in |
| **iOS Safari** | Graceful degradation | Core features work, advanced features degrade |

---

## Player Components

```
┌─────────────────────────────────────────────────────────────────┐
│                      Revenge Player Stack                        │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ┌───────────────────────────────────────────────────────────┐  │
│  │                   Player Manager                           │  │
│  │  (Unified interface, mode switching)                      │  │
│  └───────────────────────────────────────────────────────────┘  │
│         │                                │                       │
│         ▼                                ▼                       │
│  ┌──────────────┐                ┌─────────────────┐           │
│  │ Video Mode   │                │   Audio Mode    │           │
│  │ - Shaka      │                │ - Web Audio API │           │
│  │ - hls.js     │                │ - Howler.js     │           │
│  └──────────────┘                └─────────────────┘           │
│         │                                │                       │
│         └────────────┬───────────────────┘                       │
│                      ▼                                           │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │              Stream Manager                              │   │
│  │  - Protocol selection (HLS/DASH/Progressive)            │   │
│  │  - Direct Play detection                                │   │
│  │  - Blackbeard communication                             │   │
│  └─────────────────────────────────────────────────────────┘   │
│                      │                                           │
│                      ▼                                           │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │              UI Components                               │   │
│  │  - Custom controls, seek bar, volume                    │   │
│  │  - Lyrics overlay, visualizations                       │   │
│  │  - Quality selector, subtitle selector                  │   │
│  └─────────────────────────────────────────────────────────┘   │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

---

## Streaming Strategy

### Decision Tree

```
┌─────────────────────────────────────────────────────────────────┐
│              Playback Request                                    │
└─────────────────────────────────────────────────────────────────┘
                         │
                         ▼
        ┌────────────────────────────────┐
        │ Check Client Capabilities      │
        │ (codecs, containers, features) │
        └────────────────────────────────┘
                         │
        ┌────────────────┴────────────────┐
        ▼                                 ▼
┌──────────────────┐            ┌──────────────────┐
│ Codecs Match?    │            │ Container Match? │
│ H.264/VP9/AV1    │            │ MP4/WebM/MKV     │
│ AAC/Opus         │            │                  │
└──────────────────┘            └──────────────────┘
        │                                 │
   YES  │  NO                        YES  │  NO
        ▼                                 ▼
┌──────────────────┐            ┌──────────────────┐
│ DIRECT PLAY      │            │ DIRECT STREAM    │
│ - No Blackbeard  │            │ - Remux only     │
│ - Progressive    │            │ - No transcode   │
│ - Seek: Range    │            │ - Fast start     │
│ - MKV supported! │            │                  │
└──────────────────┘            └──────────────────┘
                         │
                         │ (Codec/Container mismatch)
                         ▼
                ┌──────────────────┐
                │ TRANSCODE        │
                │ - via Blackbeard │
                │ - HLS/DASH       │
                │ - Quality ladder │
                └──────────────────┘
```

### Capability Detection

```typescript
// internal/service/playback/capabilities.go
type ClientCapabilities struct {
    // Video codecs
    VideoCodecs []string // ["h264", "hevc", "vp9", "av1"]

    // Audio codecs
    AudioCodecs []string // ["aac", "mp3", "opus", "flac"]

    // Containers
    Containers  []string // ["mp4", "webm", "mkv"]

    // Features
    Features struct {
        HLS         bool
        DASH        bool
        MSE         bool  // Media Source Extensions
        WebAudio    bool
        PiP         bool
        Chromecast  bool
    }

    // Quality preferences
    MaxResolution  string  // "4k", "1080p", etc.
    MaxBitrate     int     // kbps
    PreferredCodec string  // User preference
}

// Detect from User-Agent + explicit probing
func DetectCapabilities(ctx context.Context, userAgent string) ClientCapabilities {
    caps := ClientCapabilities{}

    // Browser detection
    browser := parseUserAgent(userAgent)

    switch {
    case browser.IsSafari():
        caps.VideoCodecs = []string{"h264"}  // No HEVC in Safari (yet)
        caps.AudioCodecs = []string{"aac", "mp3"}
        caps.Containers = []string{"mp4"}
        caps.Features.HLS = true  // Safari native HLS

    case browser.IsChrome() || browser.IsEdge():
        caps.VideoCodecs = []string{"h264", "vp9", "av1"}
        caps.AudioCodecs = []string{"aac", "mp3", "opus"}
        caps.Containers = []string{"mp4", "webm", "mkv"}  // MKV Direct Play supported!
        caps.Features.HLS = false  // Use hls.js
        caps.Features.DASH = true
        caps.Features.MSE = true

    case browser.IsFirefox():
        caps.VideoCodecs = []string{"h264", "vp9", "av1"}
        caps.AudioCodecs = []string{"aac", "mp3", "opus", "flac"}
        caps.Containers = []string{"mp4", "webm", "mkv"}  // MKV Direct Play supported!
        caps.Features.DASH = true
        caps.Features.MSE = true
    }

    return caps
}
```

---

## MKV Direct Play Support

### Browser Compatibility

Modern browsers (Chrome, Firefox, Edge) **support MKV container natively** via Media Source Extensions (MSE):

| Browser | MKV Support | Codecs | Notes |
|---------|-------------|--------|-------|
| **Chrome** | ✅ Full | VP9/AV1 + Opus | Best MKV support |
| **Firefox** | ✅ Full | VP9/AV1 + Opus/FLAC | Native support |
| **Edge** | ✅ Full | VP9/AV1 + Opus | Chromium-based |
| **Safari** | ❌ No | - | Only MP4 supported |

### Why MKV Direct Play?

**Benefits:**
- No remuxing overhead (instant playback)
- Preserve original quality (lossless)
- Support for modern codecs (VP9, AV1, Opus)
- Multiple audio/subtitle tracks in one file

**Common MKV Scenarios:**
1. **Anime**: MKV with H.264 + AAC → Direct Play in Chrome/Firefox
2. **4K HDR**: MKV with VP9/AV1 + Opus → Direct Play (no transcoding!)
3. **Archival**: MKV with FLAC audio → Direct Play in Firefox

### Implementation Strategy

```typescript
// internal/service/playback/capabilities.go
func CanDirectPlay(file MediaFile, caps ClientCapabilities) bool {
    // Check container support
    containerSupported := contains(caps.Containers, file.Container)
    
    // Check video codec support
    videoSupported := contains(caps.VideoCodecs, file.VideoCodec)
    
    // Check audio codec support
    audioSupported := contains(caps.AudioCodecs, file.AudioCodec)
    
    return containerSupported && videoSupported && audioSupported
}

// Example: MKV with VP9 + Opus on Chrome
file := MediaFile{
    Container: "mkv",
    VideoCodec: "vp9",
    AudioCodec: "opus",
}

caps := ClientCapabilities{
    Containers: []string{"mp4", "webm", "mkv"},
    VideoCodecs: []string{"h264", "vp9", "av1"},
    AudioCodecs: []string{"aac", "opus"},
}

canPlay := CanDirectPlay(file, caps)  // ✅ true
```

### MKV MIME Types

```typescript
// Correct MIME type for MKV
const mimeTypes = {
    'mkv': 'video/x-matroska',
    'webm': 'video/webm',
    'mp4': 'video/mp4',
};

// Server response headers
Content-Type: video/x-matroska
Accept-Ranges: bytes
Content-Length: 1234567890
```

### Video Element Setup

```typescript
// Direct Play MKV file
const video = document.createElement('video');
video.src = '/api/media/stream/movie-123.mkv';
video.type = 'video/x-matroska';
video.play();

// With Media Source Extensions (for seeking)
const mediaSource = new MediaSource();
video.src = URL.createObjectURL(mediaSource);

mediaSource.addEventListener('sourceopen', () => {
    const sourceBuffer = mediaSource.addSourceBuffer('video/x-matroska; codecs="vp9,opus"');
    // Fetch and append video chunks...
});
```

### Fallback Strategy

If MKV Direct Play fails (e.g., Safari), fallback to remux:

```
MKV (VP9+Opus) on Chrome/Firefox → DIRECT PLAY ✅
MKV (VP9+Opus) on Safari         → REMUX to MP4 (Safari can't play MKV)
MKV (HEVC+DTS) on Chrome          → TRANSCODE (Chrome doesn't support HEVC/DTS)
```

```typescript
function selectPlaybackMethod(file: MediaFile, caps: ClientCapabilities): PlaybackMethod {
    if (canDirectPlay(file, caps)) {
        return PlaybackMethod.DirectPlay;  // ✅ MKV works!
    }
    
    if (canDirectStream(file, caps)) {
        return PlaybackMethod.DirectStream;  // Remux MKV → MP4
    }
    
    return PlaybackMethod.Transcode;  // Last resort via Blackbeard
}
```

### Testing MKV Support

```typescript
// Feature detection
function supportsFormat(mimeType: string): boolean {
    const video = document.createElement('video');
    return video.canPlayType(mimeType) !== '';
}

// Test MKV support
const mkvSupported = supportsFormat('video/x-matroska; codecs="vp9,opus"');
console.log('MKV Direct Play:', mkvSupported ? 'Supported' : 'Not supported');
```

### Performance Benefits

| Scenario | Without MKV Support | With MKV Support |
|----------|---------------------|------------------|
| **Anime Library** (MKV+H.264) | Remux to MP4 (~5s delay) | Direct Play (instant) ✅ |
| **4K VP9 Content** | Transcode to H.264 (CPU intensive) | Direct Play (no CPU) ✅ |
| **Multiple Audio Tracks** | Extract + remux each track | Direct Play all tracks ✅ |

**Result:** ~80% of modern content can Direct Play without any server-side processing!

---

## HLS vs DASH

### Comparison

| Aspect | HLS | DASH |
|--------|-----|------|
| **Browser Support** | Safari native, others need hls.js | Chrome/Firefox native, Safari needs Shaka |
| **Codec Support** | H.264/AAC primary | H.264/HEVC/VP9/AV1/Opus |
| **Manifest** | .m3u8 (simple text) | .mpd (XML, complex) |
| **Segment Format** | .ts or fMP4 | fMP4 (fragmented MP4) |
| **Quality Switching** | Bandwidth-based | More flexible (bitrate/resolution/codec) |
| **DRM** | FairPlay (Apple) | Widevine (Google), PlayReady (MS) |
| **Maturity** | Older, battle-tested | Newer, more features |

### Strategy: HLS Primary, DASH Fallback

```typescript
// Playback protocol selection
function selectProtocol(capabilities: ClientCapabilities, media: MediaInfo): Protocol {
    // Safari: Always HLS (native)
    if (capabilities.Features.HLS && isSafari()) {
        return Protocol.HLS;
    }

    // Modern codecs (AV1, HEVC): Prefer DASH
    if (media.codec === 'av1' || media.codec === 'hevc') {
        if (capabilities.Features.DASH) {
            return Protocol.DASH;
        }
    }

    // Default: HLS with hls.js
    return Protocol.HLS;
}
```

### Error Handling & Failover

```typescript
class StreamManager {
    private protocols = [Protocol.HLS, Protocol.DASH, Protocol.Progressive];
    private currentProtocolIndex = 0;

    async playMedia(mediaId: string): Promise<void> {
        const protocol = this.protocols[this.currentProtocolIndex];

        try {
            const streamUrl = await this.getStreamUrl(mediaId, protocol);
            await this.player.load(streamUrl, protocol);
        } catch (error) {
            console.warn(`Failed to play with ${protocol}, trying fallback`, error);

            // Try next protocol
            this.currentProtocolIndex++;
            if (this.currentProtocolIndex < this.protocols.length) {
                return this.playMedia(mediaId);  // Recursive retry
            }

            throw new Error('All streaming protocols failed');
        }
    }
}
```

---

## Video Player Implementation

### Library: Shaka Player + hls.js

**Shaka Player** (v4.8+): DASH support, DRM, advanced features
**hls.js** (v1.5+): HLS support for non-Safari browsers

```typescript
// lib/player/VideoPlayer.ts
import shaka from 'shaka-player';
import Hls from 'hls.js';

export class VideoPlayer {
    private video: HTMLVideoElement;
    private shakaPlayer: shaka.Player | null = null;
    private hlsPlayer: Hls | null = null;
    private currentProtocol: Protocol;

    constructor(videoElement: HTMLVideoElement) {
        this.video = videoElement;
    }

    async loadStream(url: string, protocol: Protocol): Promise<void> {
        this.cleanup();

        switch (protocol) {
            case Protocol.HLS:
                if (this.video.canPlayType('application/vnd.apple.mpegurl')) {
                    // Safari native HLS
                    this.video.src = url;
                    await this.video.play();
                } else {
                    // hls.js for other browsers
                    await this.loadHLS(url);
                }
                break;

            case Protocol.DASH:
                await this.loadDASH(url);
                break;

            case Protocol.Progressive:
                this.video.src = url;
                await this.video.play();
                break;
        }
    }

    private async loadHLS(url: string): Promise<void> {
        if (!Hls.isSupported()) {
            throw new Error('HLS not supported');
        }

        this.hlsPlayer = new Hls({
            maxBufferLength: 30,      // 30 second buffer
            maxMaxBufferLength: 600,  // Max 10 min
            startLevel: -1,           // Auto quality
            enableWorker: true,       // Use Web Worker
        });

        this.hlsPlayer.loadSource(url);
        this.hlsPlayer.attachMedia(this.video);

        // Quality change handling
        this.hlsPlayer.on(Hls.Events.MANIFEST_PARSED, () => {
            this.video.play();
        });

        // Error handling with retry
        this.hlsPlayer.on(Hls.Events.ERROR, (event, data) => {
            if (data.fatal) {
                this.handleHLSError(data);
            }
        });
    }

    private async loadDASH(url: string): Promise<void> {
        this.shakaPlayer = new shaka.Player(this.video);

        // Configure
        this.shakaPlayer.configure({
            streaming: {
                bufferingGoal: 30,
                rebufferingGoal: 2,
                bufferBehind: 30,
            },
            abr: {
                enabled: true,
                defaultBandwidthEstimate: 5000000, // 5 Mbps
            },
        });

        await this.shakaPlayer.load(url);
        await this.video.play();
    }

    // Quality switching
    async setQuality(qualityId: string): Promise<void> {
        if (this.shakaPlayer) {
            const tracks = this.shakaPlayer.getVariantTracks();
            const track = tracks.find(t => t.id === qualityId);
            if (track) {
                this.shakaPlayer.selectVariantTrack(track, true);
            }
        } else if (this.hlsPlayer) {
            const levelIndex = parseInt(qualityId);
            this.hlsPlayer.currentLevel = levelIndex;
        }
    }

    private cleanup(): void {
        if (this.shakaPlayer) {
            this.shakaPlayer.destroy();
            this.shakaPlayer = null;
        }
        if (this.hlsPlayer) {
            this.hlsPlayer.destroy();
            this.hlsPlayer = null;
        }
    }
}
```

---

## Audio Player Implementation

### Library: Web Audio API + Howler.js

**Web Audio API**: Low-level, gapless, crossfade, visualization
**Howler.js** (v2.2+): High-level wrapper, format support

```typescript
// lib/player/AudioPlayer.ts
import { Howl, Howler } from 'howler';

export class AudioPlayer {
    private audioContext: AudioContext;
    private currentTrack: Howl | null = null;
    private nextTrack: Howl | null = null;
    private queue: Track[] = [];
    private currentIndex = 0;

    // Crossfade
    private crossfadeDuration = 5000; // 5 seconds

    // Gapless prefetch
    private prefetchThreshold = 30; // Prefetch 30s before end

    constructor() {
        this.audioContext = new AudioContext();
    }

    async playTrack(track: Track, startTime = 0): Promise<void> {
        // Stop current if playing
        if (this.currentTrack) {
            this.currentTrack.stop();
        }

        this.currentTrack = new Howl({
            src: [track.url],
            format: [track.format], // 'mp3', 'flac', 'opus'
            html5: false,  // Use Web Audio API
            volume: 1.0,
            onload: () => this.onTrackLoaded(),
            onplay: () => this.onTrackPlay(),
            onend: () => this.onTrackEnd(),
        });

        if (startTime > 0) {
            this.currentTrack.seek(startTime);
        }

        this.currentTrack.play();

        // Start prefetch timer
        this.startPrefetchTimer();
    }

    // Gapless: Prefetch next track
    private startPrefetchTimer(): void {
        const checkInterval = 1000; // Check every second

        const timer = setInterval(() => {
            if (!this.currentTrack) {
                clearInterval(timer);
                return;
            }

            const duration = this.currentTrack.duration();
            const position = this.currentTrack.seek() as number;
            const remaining = duration - position;

            if (remaining <= this.prefetchThreshold && !this.nextTrack) {
                this.prefetchNextTrack();
            }

            // Crossfade trigger
            if (remaining <= (this.crossfadeDuration / 1000) && this.nextTrack) {
                this.startCrossfade();
            }
        }, checkInterval);
    }

    private async prefetchNextTrack(): Promise<void> {
        const nextInQueue = this.queue[this.currentIndex + 1];
        if (!nextInQueue) return;

        console.log('Prefetching next track:', nextInQueue.title);

        this.nextTrack = new Howl({
            src: [nextInQueue.url],
            format: [nextInQueue.format],
            html5: false,
            volume: 0, // Start silent
            preload: true,
        });

        // Also prefetch track after next (minimal)
        this.prefetchAdjacentTracks();
    }

    // Intelligent caching: Fetch minimal data for track+2 and track-1
    private async prefetchAdjacentTracks(): Promise<void> {
        const prevTrack = this.queue[this.currentIndex - 1];
        const nextNext = this.queue[this.currentIndex + 2];

        // Prefetch just the first 5 seconds
        if (prevTrack) {
            await this.partialPrefetch(prevTrack.url, 0, 5);
        }
        if (nextNext) {
            await this.partialPrefetch(nextNext.url, 0, 5);
        }
    }

    private async partialPrefetch(url: string, start: number, duration: number): Promise<void> {
        // Use HTTP Range request to fetch only partial audio
        const response = await fetch(url, {
            headers: {
                'Range': `bytes=0-${start + duration * 128000 / 8}`, // Rough estimate
            },
        });
        const blob = await response.blob();
        // Cache in memory/IndexedDB
        await cacheAudioChunk(url, blob);
    }

    // True overlapping crossfade
    private startCrossfade(): void {
        if (!this.currentTrack || !this.nextTrack) return;

        const steps = 50;
        const stepTime = this.crossfadeDuration / steps;
        let currentStep = 0;

        // Start next track at 0 volume
        this.nextTrack.play();

        const interval = setInterval(() => {
            currentStep++;
            const progress = currentStep / steps;

            // Fade out current
            this.currentTrack!.volume(1 - progress);

            // Fade in next
            this.nextTrack!.volume(progress);

            if (currentStep >= steps) {
                clearInterval(interval);

                // Swap references
                this.currentTrack?.stop();
                this.currentTrack = this.nextTrack;
                this.nextTrack = null;
                this.currentIndex++;

                // Start prefetch timer for new current
                this.startPrefetchTimer();
            }
        }, stepTime);
    }

    // No crossfade - instant gapless
    private onTrackEnd(): void {
        if (this.nextTrack) {
            this.currentTrack = this.nextTrack;
            this.nextTrack = null;
            this.currentIndex++;

            this.currentTrack.volume(1.0);
            this.currentTrack.play();

            this.startPrefetchTimer();
        } else {
            // Queue ended
            this.onQueueEnd();
        }
    }
}
```

---

## Synced Lyrics

### Data Model

```sql
CREATE TABLE track_lyrics (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    track_id UUID NOT NULL REFERENCES tracks(id),
    language VARCHAR(10) DEFAULT 'en',
    sync_type VARCHAR(20) NOT NULL,  -- 'line', 'word', 'unsynced'
    lyrics_json JSONB NOT NULL,      -- [{time_ms: 1000, text: "Line"}, ...]
    source VARCHAR(50),               -- 'lrclib', 'embedded', 'manual'
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

### LRC Format Support

```typescript
// Parse LRC file format
interface LyricLine {
    timeMs: number;
    text: string;
}

function parseLRC(lrcContent: string): LyricLine[] {
    const lines = lrcContent.split('\n');
    const lyrics: LyricLine[] = [];

    const timeRegex = /\[(\d{2}):(\d{2})\.(\d{2,3})\]/g;

    for (const line of lines) {
        const matches = Array.from(line.matchAll(timeRegex));
        if (matches.length === 0) continue;

        const text = line.replace(timeRegex, '').trim();

        for (const match of matches) {
            const minutes = parseInt(match[1]);
            const seconds = parseInt(match[2]);
            const centiseconds = parseInt(match[3].padEnd(3, '0'));

            const timeMs = (minutes * 60 + seconds) * 1000 + centiseconds;

            lyrics.push({ timeMs, text });
        }
    }

    return lyrics.sort((a, b) => a.timeMs - b.timeMs);
}
```

### Lyrics Display Component

```svelte
<!-- components/LyricsOverlay.svelte -->
<script lang="ts">
    import { onMount, onDestroy } from 'svelte';
    import { fade } from 'svelte/transition';

    export let lyrics: LyricLine[];
    export let currentTime: number; // In seconds

    let currentLineIndex = -1;
    let previousLineIndex = -1;

    $: {
        // Find current line based on time
        const timeMs = currentTime * 1000;
        const index = lyrics.findIndex((line, i) => {
            const nextLine = lyrics[i + 1];
            return line.timeMs <= timeMs && (!nextLine || nextLine.timeMs > timeMs);
        });

        if (index !== currentLineIndex) {
            previousLineIndex = currentLineIndex;
            currentLineIndex = index;
        }
    }

    function getLineClass(index: number): string {
        if (index === currentLineIndex) return 'current';
        if (index === previousLineIndex) return 'previous';
        if (index === currentLineIndex + 1) return 'next';
        return 'future';
    }
</script>

<div class="lyrics-container">
    {#each lyrics as line, i}
        <div
            class="lyric-line {getLineClass(i)}"
            class:active={i === currentLineIndex}
            in:fade={{ duration: 300 }}
        >
            {line.text}
        </div>
    {/each}
</div>

<style>
    .lyrics-container {
        position: absolute;
        bottom: 80px;
        left: 0;
        right: 0;
        text-align: center;
        pointer-events: none;
    }

    .lyric-line {
        font-size: 1.2rem;
        margin: 0.5rem 0;
        opacity: 0.4;
        transition: all 0.3s ease;
    }

    .lyric-line.current {
        font-size: 1.5rem;
        opacity: 1;
        font-weight: 600;
        color: var(--accent-color);
    }

    .lyric-line.previous {
        opacity: 0.6;
    }

    .lyric-line.next {
        opacity: 0.5;
    }
</style>
```

---

## Audio Visualization

### Canvas-based Visualizer

```typescript
// lib/player/AudioVisualizer.ts
export class AudioVisualizer {
    private canvas: HTMLCanvasElement;
    private ctx: CanvasRenderingContext2D;
    private analyser: AnalyserNode;
    private dataArray: Uint8Array;
    private animationId: number | null = null;

    constructor(canvas: HTMLCanvasElement, audioContext: AudioContext, source: MediaElementAudioSourceNode) {
        this.canvas = canvas;
        this.ctx = canvas.getContext('2d')!;

        // Create analyser
        this.analyser = audioContext.createAnalyser();
        this.analyser.fftSize = 2048;
        const bufferLength = this.analyser.frequencyBinCount;
        this.dataArray = new Uint8Array(bufferLength);

        // Connect audio
        source.connect(this.analyser);
        this.analyser.connect(audioContext.destination);
    }

    start(): void {
        this.draw();
    }

    stop(): void {
        if (this.animationId) {
            cancelAnimationFrame(this.animationId);
            this.animationId = null;
        }
    }

    private draw = (): void => {
        this.animationId = requestAnimationFrame(this.draw);

        this.analyser.getByteFrequencyData(this.dataArray);

        const { width, height } = this.canvas;
        this.ctx.clearRect(0, 0, width, height);

        // Draw bars
        const barWidth = (width / this.dataArray.length) * 2.5;
        let x = 0;

        for (let i = 0; i < this.dataArray.length; i++) {
            const barHeight = (this.dataArray[i] / 255) * height;

            // Gradient color based on frequency
            const hue = (i / this.dataArray.length) * 360;
            this.ctx.fillStyle = `hsl(${hue}, 70%, 50%)`;

            this.ctx.fillRect(x, height - barHeight, barWidth, barHeight);
            x += barWidth + 1;
        }
    };
}
```

### Fanart-based Visualization (Alternative)

```svelte
<!-- components/FanartVisualizer.svelte -->
<script lang="ts">
    export let fanartUrl: string;
    export let audioData: Uint8Array; // From analyser

    let scale = 1;
    let rotation = 0;

    $: if (audioData) {
        // Calculate average volume
        const average = audioData.reduce((a, b) => a + b, 0) / audioData.length;
        scale = 1 + (average / 255) * 0.3; // Pulse with beat
        rotation += average / 1000; // Slow rotate
    }
</script>

<div class="fanart-visualizer">
    <img
        src={fanartUrl}
        alt="Album art"
        style="
            transform: scale({scale}) rotate({rotation}deg);
            filter: blur({Math.max(0, scale - 1) * 20}px);
        "
    />
    <div class="overlay" />
</div>

<style>
    .fanart-visualizer {
        position: absolute;
        inset: 0;
        overflow: hidden;
    }

    img {
        width: 100%;
        height: 100%;
        object-fit: cover;
        transition: transform 0.1s ease-out, filter 0.1s ease-out;
    }

    .overlay {
        position: absolute;
        inset: 0;
        background: linear-gradient(to top, rgba(0,0,0,0.9), transparent);
    }
</style>
```

---

## Subtitle Support

### External WebVTT

```typescript
// Blackbeard extracts subtitles from container
GET /api/v1/playback/{mediaId}/subtitles/{trackIndex}.vtt

// Add to video element
const track = document.createElement('track');
track.kind = 'subtitles';
track.label = 'English';
track.srclang = 'en';
track.src = subtitleUrl;
videoElement.appendChild(track);
```

### Internal Subtitle Tracks (from Container)

```go
// Blackbeard extracts subtitle streams without transcoding
type SubtitleTrack struct {
    Index    int    `json:"index"`
    Codec    string `json:"codec"`    // "srt", "ass", "pgs"
    Language string `json:"language"`
    Title    string `json:"title"`
    Default  bool   `json:"default"`
    Forced   bool   `json:"forced"`
}

// Extract with FFmpeg
ffmpeg -i input.mkv -map 0:s:0 -c copy subtitle.srt
```

### Styling

```css
/* Custom subtitle styling */
::cue {
    font-family: var(--font-sans);
    font-size: 1.2rem;
    color: white;
    background-color: rgba(0, 0, 0, 0.8);
    padding: 0.2em 0.5em;
}

::cue(.large) {
    font-size: 1.5rem;
}
```

---

## Quality Switching

### User Settings

```typescript
interface QualitySettings {
    // Video
    maxVideoQuality: '4k' | '1080p' | '720p' | '480p' | '360p';
    preferredCodec: 'h264' | 'hevc' | 'av1' | 'auto';

    // Audio
    maxAudioQuality: 'lossless' | 'high' | 'normal' | 'low';

    // Bandwidth
    maxBitrate: number; // kbps, 0 = unlimited
    adaptiveBitrate: boolean;

    // Network
    wifiQuality: 'high' | 'auto';
    cellularQuality: 'low' | 'off';
}
```

### Seamless Quality Switch

```typescript
// Client sends quality change request
async function changeQuality(newQuality: string): Promise<void> {
    const currentTime = videoPlayer.getCurrentTime();

    // Send request to Blackbeard via WebSocket
    ws.send(JSON.stringify({
        type: 'quality_change',
        mediaId: currentMediaId,
        quality: newQuality,
        position: currentTime,
    }));

    // Blackbeard responds with new stream URL
    // Switch happens without stopping playback
}

// WebSocket handler in Blackbeard
func (s *StreamService) HandleQualityChange(req QualityChangeRequest) {
    // Get current stream session
    session := s.sessions[req.SessionID]

    // Start new transcode with new quality
    newStream := s.StartTranscode(req.MediaID, req.Quality, req.Position)

    // Seamless handoff: Buffer overlap
    session.AddStream(newStream, req.Position)

    // Response with new segment URLs
    s.SendResponse(session.ClientID, &QualityChangeResponse{
        NewManifestURL: newStream.ManifestURL,
        SwitchAtSegment: calculateSwitchPoint(req.Position),
    })
}
```

---

## iOS Safari Limitations & Graceful Degradation

| Feature | Desktop | iOS Safari | Degradation |
|---------|---------|------------|-------------|
| **Gapless Playback** | ✅ Full | ⚠️ Limited | Small gap acceptable |
| **Crossfade** | ✅ Full | ❌ Not supported | Disable, use instant switch |
| **Picture-in-Picture** | ✅ Full | ✅ Native | Works |
| **Background Playback** | ✅ Full | ❌ Requires app | Pause when tab inactive |
| **Multiple Streams** | ✅ Full | ❌ One at a time | Block multiple videos |
| **Autoplay** | ✅ With user gesture | ⚠️ Restricted | Require explicit play button |
| **Custom Controls** | ✅ Full | ✅ Full | Works |

```typescript
// Feature detection
const features = {
    gapless: !isSafari() || !isIOS(),
    crossfade: !isSafari() || !isIOS(),
    backgroundPlay: !isIOS(),
    autoplay: checkAutoplaySupport(),
};

// Adjust player behavior
if (!features.crossfade) {
    audioPlayer.setCrossfade(false);
}
```

---

## Performance Optimizations

### 1. Web Workers for Audio Processing

```typescript
// audio-worker.ts
self.addEventListener('message', (e) => {
    const { type, data } = e.data;

    switch (type) {
        case 'analyze':
            const analysis = analyzeAudioBuffer(data);
            self.postMessage({ type: 'analysis', data: analysis });
            break;
    }
});
```

### 2. IndexedDB for Offline Caching

```typescript
class AudioCache {
    private db: IDBDatabase;

    async cacheTrack(trackId: string, blob: Blob): Promise<void> {
        const tx = this.db.transaction('tracks', 'readwrite');
        const store = tx.objectStore('tracks');

        await store.put({ id: trackId, data: blob, cached: Date.now() });
    }

    async getTrack(trackId: string): Promise<Blob | null> {
        const tx = this.db.transaction('tracks', 'readonly');
        const store = tx.objectStore('tracks');

        const result = await store.get(trackId);
        return result?.data || null;
    }
}
```

### 3. Request Prioritization

```typescript
// High priority: Current track
// Medium: Next track
// Low: Adjacent tracks (prefetch)

class RequestQueue {
    private highPriorityQueue: Request[] = [];
    private mediumPriorityQueue: Request[] = [];
    private lowPriorityQueue: Request[] = [];

    async fetch(priority: 'high' | 'medium' | 'low'): Promise<void> {
        const queue = this[`${priority}PriorityQueue`];
        // Process queues in order
    }
}
```

---

## Custom UI Controls

### Control Bar Components

```svelte
<!-- components/PlayerControls.svelte -->
<script lang="ts">
    export let isPlaying: boolean;
    export let currentTime: number;
    export let duration: number;
    export let volume: number;
    export let showLyrics: boolean;
    export let showVisualizer: boolean;

    export let onPlayPause: () => void;
    export let onSeek: (time: number) => void;
    export let onVolumeChange: (vol: number) => void;
    export let onToggleLyrics: () => void;
    export let onToggleVisualizer: () => void;
</script>

<div class="player-controls">
    <!-- Play/Pause -->
    <button on:click={onPlayPause} aria-label={isPlaying ? 'Pause' : 'Play'}>
        {#if isPlaying}
            <PauseIcon />
        {:else}
            <PlayIcon />
        {/if}
    </button>

    <!-- Seek Bar -->
    <div class="seek-bar">
        <span class="time">{formatTime(currentTime)}</span>
        <input
            type="range"
            min="0"
            max={duration}
            value={currentTime}
            on:input={(e) => onSeek(e.target.value)}
            aria-label="Seek"
        />
        <span class="time">{formatTime(duration)}</span>
    </div>

    <!-- Volume -->
    <div class="volume-control">
        <button aria-label="Mute">
            <VolumeIcon />
        </button>
        <input
            type="range"
            min="0"
            max="1"
            step="0.01"
            value={volume}
            on:input={(e) => onVolumeChange(e.target.value)}
            aria-label="Volume"
        />
    </div>

    <!-- Audio-specific controls -->
    {#if audioMode}
        <button on:click={onToggleLyrics} class:active={showLyrics}>
            <LyricsIcon />
        </button>
        <button on:click={onToggleVisualizer} class:active={showVisualizer}>
            <VisualizerIcon />
        </button>
    {/if}

    <!-- Video-specific controls -->
    {#if videoMode}
        <button on:click={onToggleSubtitles}>
            <SubtitlesIcon />
        </button>
        <button on:click={onToggleFullscreen}>
            <FullscreenIcon />
        </button>
    {/if}

    <!-- Quality Selector -->
    <QualityMenu bind:currentQuality />
</div>
```

---

## Architecture Summary

```yaml
Player Architecture:
  Core:
    - Unified player manager (video + audio modes)
    - Custom UI controls (full accessibility)
    - Seamless protocol failover (HLS → DASH → Progressive)

  Video:
    - Shaka Player (DASH)
    - hls.js (HLS for non-Safari)
    - Native <video> (Safari HLS, Direct Play)

  Audio:
    - Web Audio API (gapless, crossfade)
    - Howler.js (high-level wrapper)
    - Intelligent prefetch (30s before end + adjacent tracks)

  Streaming:
    - Direct Play: Native browser (H.264/AAC in MP4)
    - Direct Stream: Remux only (via Blackbeard)
    - Transcode: HLS/DASH multi-quality (via Blackbeard)

  Features:
    - Synced lyrics (LRC format)
    - Audio visualization (Canvas API + fanart effects)
    - Subtitle support (WebVTT external + container extraction)
    - Quality switching (seamless via WebSocket)
    - iOS graceful degradation

  Performance:
    - Web Workers for processing
    - IndexedDB for caching
    - Request prioritization (current > next > adjacent)
```

---

## Next Steps

1. **Implement Core Player** - Video + Audio managers
2. **Stream Protocol Handler** - HLS/DASH with failover
3. **UI Components** - Controls, lyrics, visualizer
4. **Blackbeard Integration** - WebSocket quality switching
5. **Testing** - Browser compatibility matrix
6. **Documentation** - Client integration guide for Jellyfin/etc.


<!-- SOURCE-BREADCRUMBS-START -->

## Sources & Cross-References

> Auto-generated section linking to external documentation sources

### Cross-Reference Indexes

- [All Sources Index](../../sources/SOURCES_INDEX.md) - Complete list of external documentation
- [Design ↔ Sources Map](../../sources/DESIGN_CROSSREF.md) - Which docs reference which sources

### Referenced Sources

| Source | Documentation |
|--------|---------------|
| [FFmpeg Codecs](https://ffmpeg.org/ffmpeg-codecs.html) | [Local](../../sources/media/ffmpeg-codecs.md) |
| [FFmpeg Documentation](https://ffmpeg.org/ffmpeg.html) | [Local](../../sources/media/ffmpeg.md) |
| [FFmpeg Formats](https://ffmpeg.org/ffmpeg-formats.html) | [Local](../../sources/media/ffmpeg-formats.md) |
| [M3U8 Extended Format](https://datatracker.ietf.org/doc/html/rfc8216) | [Local](../../sources/protocols/m3u8.md) |
| [Svelte 5 Documentation](https://svelte.dev/docs/svelte/overview) | [Local](../../sources/frontend/svelte5.md) |
| [Svelte 5 Runes](https://svelte.dev/docs/svelte/$state) | [Local](../../sources/frontend/svelte-runes.md) |
| [SvelteKit Documentation](https://svelte.dev/docs/kit/introduction) | [Local](../../sources/frontend/sveltekit.md) |
| [go-astiav (FFmpeg bindings)](https://pkg.go.dev/github.com/asticode/go-astiav) | [Local](../../sources/media/go-astiav.md) |
| [gohlslib (HLS)](https://pkg.go.dev/github.com/bluenviron/gohlslib/v2) | [Local](../../sources/media/gohlslib.md) |

<!-- SOURCE-BREADCRUMBS-END -->

<!-- DESIGN-BREADCRUMBS-START -->

## Related Design Docs

> Auto-generated cross-references to related design documentation

**Category**: [Architecture](INDEX.md)

### In This Section

- [Revenge - Architecture v2](01_ARCHITECTURE.md)
- [Revenge - Design Principles](02_DESIGN_PRINCIPLES.md)
- [Revenge - Metadata System](03_METADATA_SYSTEM.md)
- [Plugin Architecture Decision](05_PLUGIN_ARCHITECTURE_DECISION.md)

### Related Topics

- [Revenge - Adult Content System](../features/adult/ADULT_CONTENT_SYSTEM.md) _Adult_
- [Revenge - Adult Content Metadata System](../features/adult/ADULT_METADATA.md) _Adult_
- [Adult Data Reconciliation](../features/adult/DATA_RECONCILIATION.md) _Adult_
- [Adult Gallery Module (QAR: Treasures)](../features/adult/GALLERY_MODULE.md) _Adult_
- [Whisparr v3 & StashDB Schema Integration](../features/adult/WHISPARR_STASHDB_SCHEMA.md) _Adult_

### Indexes

- [Design Index](../DESIGN_INDEX.md) - All design docs by category/topic
- [Source of Truth](../00_SOURCE_OF_TRUTH.md) - Package versions and status

<!-- DESIGN-BREADCRUMBS-END -->

---

## Cross-References

| Related Document | Relationship |
|------------------|--------------|
| [00_SOURCE_OF_TRUTH.md](../00_SOURCE_OF_TRUTH.md) | Frontend stack, package versions |
| [01_ARCHITECTURE.md](01_ARCHITECTURE.md) | Blackbeard transcoding integration |
| [02_DESIGN_PRINCIPLES.md](02_DESIGN_PRINCIPLES.md) | WebUI player principles |
| [03_METADATA_SYSTEM.md](03_METADATA_SYSTEM.md) | Metadata for player display |
