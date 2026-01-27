---
applyTo: "**/internal/service/playback/**/*.go,**/internal/infra/streaming/**/*.go"
---

# Streaming Best Practices

> Guidelines for implementing smooth, reliable media streaming.

## Architecture

### Stream Flow

```
Client → Revenge → Buffer → Blackbeard (transcoder)
                     ↓
              Transcode Cache
```

**Key Principle:** Stream always flows through Revenge for:

- Access control validation
- Progress tracking
- Bandwidth monitoring
- Quality adaptation
- Analytics

## Buffering

### Segment Buffer

| Setting         | Recommended  | Why                     |
| --------------- | ------------ | ----------------------- |
| Buffer ahead    | 5-8 segments | ~30-48s of content      |
| Min before play | 2 segments   | ~12s initial buffer     |
| Prefetch        | 3 segments   | Smooth playback         |
| Retry count     | 3            | Handle transient errors |

```go
// Good: Configurable buffer with sensible defaults
type BufferConfig struct {
    SegmentBufferSize int           `koanf:"segment_buffer_size"` // 5
    MinBufferDuration time.Duration `koanf:"min_buffer_duration"` // 10s
    MaxBufferDuration time.Duration `koanf:"max_buffer_duration"` // 60s
}

// Good: Prefetch upcoming segments
func (h *StreamHandler) triggerPrefetch(transcodeID, currentSegment string, buffer *StreamBuffer) {
    nextSegments := h.predictNextSegments(currentSegment, 3)
    for _, seg := range nextSegments {
        go h.prefetchSegment(transcodeID, seg, buffer)
    }
}
```

### Transcode Cache

| Setting            | Recommended    | Why                         |
| ------------------ | -------------- | --------------------------- |
| Max memory         | 25% system RAM | Balance with other services |
| Min retention      | 30s            | Allow seeking back          |
| Eviction trigger   | 80% full       | Prevent OOM                 |
| Emergency eviction | 95% full       | Aggressive cleanup          |

```go
// Good: Memory-pressure-aware caching
type TranscodeCacheConfig struct {
    MaxMemoryBytes          int64   // 0 = auto (25% RAM)
    HighMemoryThreshold     float64 // 0.8
    CriticalMemoryThreshold float64 // 0.95
    MinRetentionTime        time.Duration // 30s
}
```

## HLS/DASH Best Practices

### Manifest Handling

```go
// Good: Rewrite URLs to route through Revenge
func (h *StreamHandler) rewriteManifest(manifest []byte, sessionID string) []byte {
    // Replace: segment_0.ts
    // With: /api/v1/playback/{sessionID}/segment/segment_0.ts
}

// Good: Set proper cache headers
w.Header().Set("Content-Type", "application/vnd.apple.mpegurl")
w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
```

### Segment Delivery

```go
// Good: Track segment delivery source
w.Header().Set("X-Served-From", "buffer") // or "origin", "cache"

// Good: Flush immediately for streaming
if flusher, ok := w.(http.Flusher); ok {
    flusher.Flush()
}
```

## Bandwidth Adaptation

### Measurement

```go
// Good: Rolling window for stability
type BandwidthMonitor struct {
    samples    []BandwidthSample
    windowSize int // 10 samples
}

// Good: Calculate jitter for stability assessment
func (m *BandwidthMonitor) GetEstimate() BandwidthEstimate {
    avg := calculateAverage(m.samples)
    jitter := maxKbps - minKbps
    return BandwidthEstimate{
        AverageKbps: avg,
        JitterKbps:  jitter,
        IsStable:    jitter < avg/4, // <25% variance
    }
}
```

### Profile Selection

```go
// Good: Conservative bitrate selection
targetBitrate := bandwidth * 0.8 // 80% of measured
if jitter > bandwidth/4 {
    targetBitrate = bandwidth * 0.6 // More conservative if unstable
}

// Good: Match to available profiles
profile := selectProfileForBitrate(targetBitrate)
```

## Error Handling

### Retry Logic

```go
// Good: Exponential backoff with jitter
func (h *StreamHandler) fetchWithRetry(ctx context.Context, transcodeID, segment string) ([]byte, error) {
    for attempt := 0; attempt <= maxRetries; attempt++ {
        if attempt > 0 {
            backoff := time.Duration(attempt) * 200 * time.Millisecond
            jitter := time.Duration(rand.Int63n(100)) * time.Millisecond
            select {
            case <-ctx.Done():
                return nil, ctx.Err()
            case <-time.After(backoff + jitter):
            }
        }

        data, err := h.transcoder.FetchSegment(ctx, transcodeID, segment)
        if err == nil {
            return data, nil
        }
        lastErr = err
    }
    return nil, fmt.Errorf("failed after %d attempts: %w", maxRetries+1, lastErr)
}
```

### Graceful Degradation

```go
// Good: Fall back to lower quality on errors
func (h *StreamHandler) handleSegmentError(session *Session, err error) {
    session.ErrorCount++

    if session.ErrorCount > 3 {
        // Switch to lower quality profile
        lowerProfile := h.getLowerQualityProfile(session.TranscodeProfile)
        h.requestProfileSwitch(session, lowerProfile)
        session.ErrorCount = 0
    }
}
```

## Session Management

### State Tracking

```go
// Good: Track comprehensive session state
type Session struct {
    // Identity
    ID        uuid.UUID
    UserID    uuid.UUID
    MediaID   uuid.UUID

    // Playback
    Position  time.Duration
    State     State // buffering, playing, paused, stopped

    // Quality
    CurrentBitrate  int
    QualitySwitches int
    BufferingEvents int

    // Tracking
    LastActivityAt time.Time
}
```

### Cleanup

```go
// Good: Clean up idle sessions
func (m *SessionManager) CleanupIdleSessions(ctx context.Context) {
    ticker := time.NewTicker(time.Minute)
    defer ticker.Stop()

    for {
        select {
        case <-ctx.Done():
            return
        case <-ticker.C:
            m.mu.Lock()
            for id, session := range m.sessions {
                if time.Since(session.LastActivityAt) > 30*time.Minute {
                    m.stopSession(ctx, id)
                }
            }
            m.mu.Unlock()
        }
    }
}
```

## Raw File Serving (to Blackbeard)

### HTTP Range Requests

```go
// Good: Full range request support
func (s *MediaFileServer) ServeFile(w http.ResponseWriter, r *http.Request, filePath string) {
    rangeHeader := r.Header.Get("Range")
    if rangeHeader != "" {
        // Parse: bytes=start-end
        // Return 206 Partial Content
        w.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, fileSize))
        w.WriteHeader(http.StatusPartialContent)
    } else {
        w.WriteHeader(http.StatusOK)
    }
}
```

### Chunked Streaming

```go
// Good: Stream in chunks, respect cancellation
func (s *MediaFileServer) streamFile(ctx context.Context, file *os.File, w io.Writer, length int64) {
    buf := make([]byte, 64*1024) // 64KB chunks
    remaining := length

    for remaining > 0 {
        select {
        case <-ctx.Done():
            return // Client disconnected
        default:
        }

        toRead := min(int64(len(buf)), remaining)
        n, err := file.Read(buf[:toRead])
        if err != nil {
            return
        }

        w.Write(buf[:n])
        if flusher, ok := w.(http.Flusher); ok {
            flusher.Flush()
        }

        remaining -= int64(n)
    }
}
```

## Performance

### Avoid

```go
// Bad: Loading entire file into memory
data, _ := os.ReadFile(filePath)
w.Write(data)

// Bad: No context cancellation
func streamForever(w io.Writer, file *os.File) {
    io.Copy(w, file) // Never stops if client disconnects
}

// Bad: No buffering
for {
    segment := fetchSegment() // Blocking on every segment
    w.Write(segment.Data)
}
```

### Prefer

```go
// Good: Stream without full memory load
file, _ := os.Open(filePath)
io.Copy(w, file)

// Good: Respect cancellation
func streamWithContext(ctx context.Context, w io.Writer, file *os.File) {
    reader := &contextReader{ctx: ctx, r: file}
    io.Copy(w, reader)
}

// Good: Buffer ahead
buffer := h.getOrCreateBuffer(transcodeID)
go h.prefetchSegments(transcodeID, currentPos, buffer)
```

## Monitoring

### Key Metrics

```go
// Track these for every stream
type StreamMetrics struct {
    BufferingEvents    int           // Count of rebuffers
    TotalBufferTime    time.Duration // Time spent buffering
    QualitySwitches    int           // ABR switches
    AverageBitrate     int           // Kbps delivered
    SegmentsFetched    int           // From origin
    SegmentsFromCache  int           // Cache hits
    ErrorCount         int           // Fetch failures
}
```

### Logging

```go
// Good: Structured logging with context
h.logger.Info("segment served",
    "session_id", session.ID,
    "segment", segmentPath,
    "source", "buffer", // or "origin", "cache"
    "latency_ms", latency.Milliseconds(),
)

// Good: Log quality switches
h.logger.Info("quality switch",
    "session_id", session.ID,
    "from_profile", oldProfile,
    "to_profile", newProfile,
    "reason", "bandwidth_drop",
)
```
