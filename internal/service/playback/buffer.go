// Package playback provides stream buffering for stable playback.
package playback

import (
	"context"
	"errors"
	"io"
	"sync"
	"time"
)

// BufferConfig holds buffer configuration.
type BufferConfig struct {
	// SegmentBufferSize is how many HLS/DASH segments to buffer ahead.
	SegmentBufferSize int `koanf:"segment_buffer_size"`
	// MinBufferDuration is minimum buffered duration before starting playback.
	MinBufferDuration time.Duration `koanf:"min_buffer_duration"`
	// MaxBufferDuration is maximum buffered duration to prevent memory issues.
	MaxBufferDuration time.Duration `koanf:"max_buffer_duration"`
	// RecoveryTimeout is how long to wait for recovery on errors.
	RecoveryTimeout time.Duration `koanf:"recovery_timeout"`
}

// DefaultBufferConfig returns sensible defaults.
func DefaultBufferConfig() BufferConfig {
	return BufferConfig{
		SegmentBufferSize: 5, // 5 segments ahead
		MinBufferDuration: 10 * time.Second,
		MaxBufferDuration: 60 * time.Second,
		RecoveryTimeout:   30 * time.Second,
	}
}

// StreamBuffer buffers transcoded stream segments for stable delivery.
type StreamBuffer struct {
	mu       sync.RWMutex
	config   BufferConfig
	segments map[int]*BufferedSegment

	// Key-based storage (for segment path strings)
	segmentsByKey map[string]*BufferedSegment

	// Tracking
	lastRequestedSeq int
	lastBufferedSeq  int
	lastAccess       time.Time

	// Stats
	bufferUnderruns  int
	recoveryAttempts int
}

// BufferedSegment holds a cached stream segment.
type BufferedSegment struct {
	Sequence   int
	Key        string // segment path for key-based access
	Data       []byte
	Duration   time.Duration
	CreatedAt  time.Time
	AccessedAt time.Time
	MimeType   string
}

// NewStreamBuffer creates a new stream buffer.
func NewStreamBuffer(config BufferConfig) *StreamBuffer {
	return &StreamBuffer{
		config:        config,
		segments:      make(map[int]*BufferedSegment),
		segmentsByKey: make(map[string]*BufferedSegment),
		lastAccess:    time.Now(),
	}
}

// PutSegment adds a segment to the buffer.
func (b *StreamBuffer) PutSegment(seq int, data []byte, duration time.Duration, mimeType string) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.segments[seq] = &BufferedSegment{
		Sequence:   seq,
		Data:       data,
		Duration:   duration,
		CreatedAt:  time.Now(),
		AccessedAt: time.Now(),
		MimeType:   mimeType,
	}

	if seq > b.lastBufferedSeq {
		b.lastBufferedSeq = seq
	}

	// Evict old segments beyond buffer size
	b.evictOldSegments()
}

// GetSegment retrieves a segment from the buffer.
func (b *StreamBuffer) GetSegment(seq int) (*BufferedSegment, bool) {
	b.mu.Lock()
	defer b.mu.Unlock()

	seg, ok := b.segments[seq]
	if ok {
		seg.AccessedAt = time.Now()
		b.lastRequestedSeq = seq
		b.lastAccess = time.Now()
	}
	return seg, ok
}

// Get retrieves a segment by key (segment path).
func (b *StreamBuffer) Get(key string) (*BufferedSegment, bool) {
	b.mu.Lock()
	defer b.mu.Unlock()

	seg, ok := b.segmentsByKey[key]
	if ok {
		seg.AccessedAt = time.Now()
		b.lastAccess = time.Now()
	}
	return seg, ok
}

// Add adds a segment by key (segment path).
func (b *StreamBuffer) Add(key string, data []byte, mimeType string) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.segmentsByKey[key] = &BufferedSegment{
		Key:        key,
		Data:       data,
		MimeType:   mimeType,
		CreatedAt:  time.Now(),
		AccessedAt: time.Now(),
	}
	b.lastAccess = time.Now()

	// Evict old segments if buffer is too large
	b.evictOldKeySegments()
}

// LastAccess returns the last time this buffer was accessed.
func (b *StreamBuffer) LastAccess() time.Time {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.lastAccess
}

// evictOldSegments removes segments that are too far behind playback.
func (b *StreamBuffer) evictOldSegments() {
	// Keep segments from (lastRequestedSeq - 2) to (lastRequestedSeq + bufferSize)
	minSeq := b.lastRequestedSeq - 2
	maxSeq := b.lastRequestedSeq + b.config.SegmentBufferSize + 2

	for seq := range b.segments {
		if seq < minSeq || seq > maxSeq {
			delete(b.segments, seq)
		}
	}
}

// evictOldKeySegments removes oldest segments by key when buffer is full.
func (b *StreamBuffer) evictOldKeySegments() {
	maxSegments := b.config.SegmentBufferSize
	if maxSegments == 0 {
		maxSegments = 10 // Default
	}

	for len(b.segmentsByKey) > maxSegments {
		// Find oldest segment
		var oldestKey string
		var oldestTime time.Time
		first := true
		for key, seg := range b.segmentsByKey {
			if first || seg.AccessedAt.Before(oldestTime) {
				oldestKey = key
				oldestTime = seg.AccessedAt
				first = false
			}
		}
		if oldestKey != "" {
			delete(b.segmentsByKey, oldestKey)
		}
	}
}

// BufferedDuration returns total buffered duration ahead of playback.
func (b *StreamBuffer) BufferedDuration() time.Duration {
	b.mu.RLock()
	defer b.mu.RUnlock()

	var total time.Duration
	for seq, seg := range b.segments {
		if seq > b.lastRequestedSeq {
			total += seg.Duration
		}
	}
	return total
}

// IsReady returns true if minimum buffer is filled.
func (b *StreamBuffer) IsReady() bool {
	return b.BufferedDuration() >= b.config.MinBufferDuration
}

// Stats returns buffer statistics.
func (b *StreamBuffer) Stats() BufferStats {
	b.mu.RLock()
	defer b.mu.RUnlock()

	var totalBytes int64
	for _, seg := range b.segments {
		totalBytes += int64(len(seg.Data))
	}

	return BufferStats{
		SegmentCount:     len(b.segments),
		BufferedDuration: b.BufferedDuration(),
		TotalBytes:       totalBytes,
		LastRequestedSeq: b.lastRequestedSeq,
		LastBufferedSeq:  b.lastBufferedSeq,
		BufferUnderruns:  b.bufferUnderruns,
		RecoveryAttempts: b.recoveryAttempts,
	}
}

// BufferStats holds buffer statistics.
type BufferStats struct {
	SegmentCount     int           `json:"segment_count"`
	BufferedDuration time.Duration `json:"buffered_duration"`
	TotalBytes       int64         `json:"total_bytes"`
	LastRequestedSeq int           `json:"last_requested_seq"`
	LastBufferedSeq  int           `json:"last_buffered_seq"`
	BufferUnderruns  int           `json:"buffer_underruns"`
	RecoveryAttempts int           `json:"recovery_attempts"`
}

// StreamProxy handles buffered proxying from Blackbeard to client.
type StreamProxy struct {
	buffer     *StreamBuffer
	transcoder *TranscoderClient
	config     BufferConfig
}

// NewStreamProxy creates a new stream proxy.
func NewStreamProxy(transcoder *TranscoderClient, config BufferConfig) *StreamProxy {
	return &StreamProxy{
		buffer:     NewStreamBuffer(config),
		transcoder: transcoder,
		config:     config,
	}
}

// ProxySegment proxies a segment with buffering and recovery.
func (p *StreamProxy) ProxySegment(ctx context.Context, transcodeID string, segmentSeq int, w io.Writer) error {
	// Check buffer first
	if seg, ok := p.buffer.GetSegment(segmentSeq); ok {
		_, err := w.Write(seg.Data)
		return err
	}

	// Fetch from Blackbeard with retry
	var lastErr error
	for attempt := 0; attempt < 3; attempt++ {
		data, mimeType, err := p.fetchSegmentFromBlackbeard(ctx, transcodeID, segmentSeq)
		if err != nil {
			lastErr = err
			// Wait before retry
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(time.Duration(attempt+1) * 500 * time.Millisecond):
				continue
			}
		}

		// Buffer the segment
		p.buffer.PutSegment(segmentSeq, data, 6*time.Second, mimeType) // Assume 6s segments

		// Write to client
		_, err = w.Write(data)
		return err
	}

	return lastErr
}

// fetchSegmentFromBlackbeard fetches a segment from Blackbeard.
func (p *StreamProxy) fetchSegmentFromBlackbeard(ctx context.Context, transcodeID string, segmentSeq int) ([]byte, string, error) {
	// This would use the TranscoderClient to fetch the segment
	// For now, placeholder implementation
	return nil, "", errors.New("not implemented - use TranscoderClient")
}

// StartPrefetch starts prefetching segments ahead of playback.
func (p *StreamProxy) StartPrefetch(ctx context.Context, transcodeID string, startSeq int) {
	go func() {
		currentSeq := startSeq
		for {
			select {
			case <-ctx.Done():
				return
			default:
				// Check if we need more segments
				if p.buffer.BufferedDuration() < p.config.MaxBufferDuration {
					// Prefetch next segment if not buffered
					if _, ok := p.buffer.GetSegment(currentSeq); !ok {
						data, mimeType, err := p.fetchSegmentFromBlackbeard(ctx, transcodeID, currentSeq)
						if err == nil {
							p.buffer.PutSegment(currentSeq, data, 6*time.Second, mimeType)
						}
					}
					currentSeq++
				} else {
					// Buffer full, wait a bit
					time.Sleep(1 * time.Second)
				}
			}
		}
	}()
}
