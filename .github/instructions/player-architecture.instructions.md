---
applyTo: "**/web/**/*.ts,**/web/**/*.svelte,**/internal/service/playback/**/*.go"
alwaysApply: false
---

# Player Architecture Instructions

> Unified web player for video and audio with native streaming and transcode fallback.

## Core Principles

1. **Unified Player**: One player manager, switches between video/audio modes
2. **Protocol Failover**: HLS primary → DASH fallback → Progressive last resort
3. **Intelligent Prefetch**: 30s before track end + minimal adjacent tracks
4. **True Crossfade**: Overlapping audio streams with gain nodes
5. **Custom Controls**: Full accessibility, no reliance on browser defaults

## Player Stack

### Video Mode
```typescript
// Libraries
import shaka from 'shaka-player';  // DASH support
import Hls from 'hls.js';          // HLS for non-Safari

// Protocol selection priority
1. Safari: Native HLS (video.src = m3u8)
2. Modern codecs (AV1/HEVC): DASH via Shaka
3. Fallback: HLS via hls.js
4. Last resort: Progressive download
```

### Audio Mode
```typescript
// Libraries
import { Howl, Howler } from 'howler';  // High-level wrapper

// Use Web Audio API for:
- Gapless playback (seamless track transitions)
- Crossfade (dual gain nodes)
- Visualizations (AnalyserNode)

// NOT Howler.js for:
- iOS Safari (Web Audio API restrictions)
```

## Streaming Decision Tree

```
Playback Request
    │
    ▼
Check Client Capabilities
    │
    ├─ Codecs Match (H.264/AAC) → DIRECT PLAY (Progressive HTTP)
    ├─ Container Mismatch → DIRECT STREAM (Remux via Blackbeard)
    └─ Codec Mismatch → TRANSCODE (HLS/DASH via Blackbeard)
```

## Code Patterns

### DO
- ✅ Use capability detection before choosing protocol
- ✅ Implement error handlers with protocol failover
- ✅ Prefetch next track 30s before current ends
- ✅ Use Web Workers for audio analysis/processing
- ✅ Cache tracks in IndexedDB for offline playback
- ✅ Implement graceful degradation for iOS Safari
- ✅ Use WebSocket for seamless quality switching

### DON'T
- ❌ Hardcode protocol choice (always detect capabilities)
- ❌ Use autoplay without user gesture (iOS restrictions)
- ❌ Rely on browser default controls
- ❌ Block UI thread for audio processing
- ❌ Assume Web Audio API works on all browsers
- ❌ Transcode in frontend (Blackbeard handles all)

## Gapless & Crossfade

### Gapless (No Overlap)
```typescript
// Prefetch next track
const checkInterval = 1000; // 1 second
const prefetchThreshold = 30; // 30 seconds

setInterval(() => {
    const remaining = duration - currentTime;
    if (remaining <= prefetchThreshold && !nextTrackLoaded) {
        preloadNextTrack();
    }
}, checkInterval);

// Instant switch on track end
onTrackEnd(() => {
    currentTrack = nextTrack;
    nextTrack = null;
    currentTrack.play();
});
```

### Crossfade (True Overlap)
```typescript
const crossfadeDuration = 5000; // 5 seconds
const steps = 50;
const stepTime = crossfadeDuration / steps;

// Start next track at 0 volume
nextTrack.volume(0);
nextTrack.play();

// Fade both tracks simultaneously
for (let i = 0; i <= steps; i++) {
    const progress = i / steps;
    
    currentTrack.volume(1 - progress);  // Fade out
    nextTrack.volume(progress);          // Fade in
    
    await sleep(stepTime);
}
```

## Synced Lyrics

### Data Model
```typescript
interface LyricLine {
    timeMs: number;  // Milliseconds from start
    text: string;    // Lyric text
}

// Parse LRC format
// [00:12.50] First line
// [00:18.75] Second line
```

### Display Strategy
```typescript
// Find current line (binary search for performance)
function findCurrentLine(lyrics: LyricLine[], timeMs: number): number {
    let left = 0, right = lyrics.length - 1;
    
    while (left <= right) {
        const mid = Math.floor((left + right) / 2);
        const current = lyrics[mid];
        const next = lyrics[mid + 1];
        
        if (current.timeMs <= timeMs && (!next || next.timeMs > timeMs)) {
            return mid;
        }
        
        if (current.timeMs > timeMs) {
            right = mid - 1;
        } else {
            left = mid + 1;
        }
    }
    
    return -1;
}

// Update every 100ms for smooth transitions
setInterval(() => {
    const currentLine = findCurrentLine(lyrics, currentTimeMs);
    updateLyricsDisplay(currentLine);
}, 100);
```

## Audio Visualization

### Frequency Bars (Canvas API)
```typescript
const analyser = audioContext.createAnalyser();
analyser.fftSize = 2048;
const bufferLength = analyser.frequencyBinCount;
const dataArray = new Uint8Array(bufferLength);

function draw() {
    requestAnimationFrame(draw);
    analyser.getByteFrequencyData(dataArray);
    
    // Draw bars
    for (let i = 0; i < bufferLength; i++) {
        const barHeight = (dataArray[i] / 255) * canvasHeight;
        ctx.fillRect(x, canvasHeight - barHeight, barWidth, barHeight);
    }
}
```

### Fanart-based (Alternative)
```typescript
// Pulse/rotate album art with beat
const average = dataArray.reduce((a, b) => a + b) / dataArray.length;
const scale = 1 + (average / 255) * 0.3;  // 1.0 to 1.3x
const rotation = (average / 255) * 10;    // 0-10 degrees

fanartElement.style.transform = `scale(${scale}) rotate(${rotation}deg)`;
```

## Subtitle Support

### External Subtitles (WebVTT)
```typescript
// Blackbeard endpoint
GET /api/v1/playback/{mediaId}/subtitles/{trackIndex}.vtt

// Add to video element
const track = document.createElement('track');
track.kind = 'subtitles';
track.label = language;
track.srclang = langCode;
track.src = subtitleUrl;
videoElement.appendChild(track);
```

### Internal Subtitles (from Container)
```typescript
// Blackbeard extracts without transcoding
// Supports: SRT, ASS, SSA, PGS, VobSub
// Returns WebVTT format for browser compatibility

// Request specific subtitle track
GET /api/v1/playback/{mediaId}/subtitles/{streamIndex}.vtt?extract=true
```

## Quality Switching (Seamless)

### Client → Revenge → Blackbeard Flow
```typescript
// 1. Client detects bandwidth change or user selection
function changeQuality(newQuality: string) {
    const currentTime = videoElement.currentTime;
    
    // 2. Send to Revenge via WebSocket
    ws.send(JSON.stringify({
        type: 'quality_change',
        sessionId: currentSessionId,
        quality: newQuality,
        position: currentTime,
    }));
}

// 3. Revenge forwards to Blackbeard
// 4. Blackbeard starts new transcode with overlap
// 5. Returns new manifest URL
// 6. Client switches at segment boundary (no playback interruption)

ws.onmessage = (event) => {
    const msg = JSON.parse(event.data);
    if (msg.type === 'quality_changed') {
        // hls.js or Shaka handles seamless switch
        player.updateManifest(msg.newManifestUrl);
    }
};
```

## iOS Safari Limitations

| Feature | Desktop | iOS Safari | Implementation |
|---------|---------|------------|----------------|
| Gapless | ✅ Full | ⚠️ ~0.1s gap | Acceptable, use instant switch |
| Crossfade | ✅ Full | ❌ Not supported | Disable on iOS, use gapless |
| Web Audio API | ✅ Full | ⚠️ Requires user gesture | Delay init until play button |
| Multiple <video> | ✅ Yes | ❌ One at a time | Block multiple videos, show warning |
| Autoplay | ✅ Yes | ❌ No | Require explicit play button |
| Background audio | ✅ Yes | ❌ Tab suspend | Use Media Session API |

```typescript
// Feature detection
const isIOS = /iPad|iPhone|iPod/.test(navigator.userAgent);
const isSafari = /^((?!chrome|android).)*safari/i.test(navigator.userAgent);

const features = {
    crossfade: !(isIOS && isSafari),
    webAudio: !(isIOS && isSafari) || userGestureReceived,
    gapless: true, // Works everywhere, just small gap on iOS
};

// Adjust player
if (!features.crossfade) {
    audioPlayer.disableCrossfade();
}
```

## Performance

### 1. Request Prioritization
```typescript
// High: Current track/video (play ASAP)
// Medium: Next track (prefetch 30s before end)
// Low: Adjacent tracks (prefetch 5s of data)

class PriorityFetcher {
    private queues = {
        high: [] as Request[],
        medium: [] as Request[],
        low: [] as Request[],
    };
    
    async fetch(url: string, priority: 'high' | 'medium' | 'low') {
        this.queues[priority].push(fetch(url));
        await this.processQueues();
    }
    
    private async processQueues() {
        // Process high first, then medium, then low
        while (this.queues.high.length) {
            await this.queues.high.shift();
        }
        // ... etc
    }
}
```

### 2. Web Worker for Processing
```typescript
// audio-worker.ts
self.addEventListener('message', (e) => {
    if (e.data.type === 'analyze') {
        const analysis = analyzeAudio(e.data.buffer);
        self.postMessage({ type: 'result', data: analysis });
    }
});

// main.ts
const worker = new Worker('audio-worker.ts');
worker.postMessage({ type: 'analyze', buffer: audioBuffer });
```

### 3. IndexedDB Caching
```typescript
// Cache downloaded tracks
async function cacheTrack(trackId: string, blob: Blob) {
    const db = await openDB('revenge-audio', 1);
    const tx = db.transaction('tracks', 'readwrite');
    await tx.store.put({ id: trackId, data: blob, cached: Date.now() });
}

// Retrieve cached track
async function getCachedTrack(trackId: string): Promise<Blob | null> {
    const db = await openDB('revenge-audio', 1);
    const track = await db.get('tracks', trackId);
    return track?.data || null;
}
```

## Testing Requirements

### Browser Compatibility Matrix
```typescript
// Must test:
- Chrome/Edge (latest + 2 versions back)
- Firefox (latest + 2 versions back)
- Safari (latest + 1 version back)
- iOS Safari (latest iOS)
- Android Chrome (latest)

// Test scenarios:
1. Direct play (native codecs)
2. HLS streaming (transcoded)
3. DASH streaming (fallback)
4. Protocol failover (kill Blackbeard mid-stream)
5. Quality switching (seamless)
6. Gapless playback (album play)
7. Crossfade (playlist mode)
8. Synced lyrics (karaoke mode)
9. Subtitles (multiple languages)
10. Offline playback (cached tracks)
```

## Integration with Blackbeard

### Stream Request
```typescript
// Client → Revenge → Blackbeard
POST /api/v1/playback/start
{
    "mediaId": "uuid",
    "protocol": "hls",        // or "dash"
    "quality": "1080p",
    "startPosition": 0,
    "audioTrack": 0,
    "subtitleTrack": null,
    "clientCapabilities": {
        "videoCodecs": ["h264", "vp9"],
        "audioCodecs": ["aac", "opus"],
        "maxResolution": "1080p",
        "maxBitrate": 8000  // kbps
    }
}

// Response
{
    "sessionId": "uuid",
    "protocol": "hls",
    "manifestUrl": "http://revenge:8096/stream/session-id/master.m3u8",
    "websocketUrl": "ws://revenge:8096/playback/session-id"
}
```

### WebSocket Events
```typescript
// Client → Server
{
    "type": "quality_change",
    "quality": "720p"
}

{
    "type": "position_update",
    "position": 123.45  // seconds
}

// Server → Client
{
    "type": "quality_changed",
    "newManifestUrl": "http://..."
}

{
    "type": "buffer_warning",
    "message": "Network slow, reducing quality"
}
```

## Summary

```yaml
Player Architecture:
  Unified Manager: Video + Audio modes
  Video Stack: Shaka (DASH) + hls.js (HLS) + Native (Safari/Direct Play)
  Audio Stack: Web Audio API + Howler.js
  
  Features:
    - Gapless: 30s prefetch + instant switch
    - Crossfade: True overlap (5s default)
    - Lyrics: LRC format, binary search lookup
    - Visualizer: Canvas frequency bars OR fanart pulse
    - Subtitles: WebVTT external + container extraction
    - Quality: Seamless via WebSocket
    
  iOS Safari:
    - Gapless: ⚠️ Works (small gap)
    - Crossfade: ❌ Disable
    - Web Audio: ⚠️ Requires user gesture
    - Autoplay: ❌ Requires play button
    
  Performance:
    - Priority fetch: High > Medium > Low
    - Web Workers: Audio processing
    - IndexedDB: Offline caching
```
