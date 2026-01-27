// Package playback provides transcode caching with memory-aware eviction.
package playback

import (
	"container/list"
	"context"
	"log/slog"
	"runtime"
	"sync"
	"time"

	"github.com/google/uuid"
)

// TranscodeCacheConfig configures the transcode cache.
type TranscodeCacheConfig struct {
	// MaxMemoryBytes is the maximum memory to use for caching.
	// If 0, uses 25% of system memory.
	MaxMemoryBytes int64 `koanf:"max_memory_bytes"`

	// MaxSegmentsPerTranscode limits segments per transcode session.
	MaxSegmentsPerTranscode int `koanf:"max_segments_per_transcode"`

	// EvictionCheckInterval is how often to check memory pressure.
	EvictionCheckInterval time.Duration `koanf:"eviction_check_interval"`

	// MinRetentionTime is minimum time to keep a segment before eviction.
	MinRetentionTime time.Duration `koanf:"min_retention_time"`

	// HighMemoryThreshold triggers aggressive eviction (0.0-1.0).
	HighMemoryThreshold float64 `koanf:"high_memory_threshold"`

	// CriticalMemoryThreshold triggers emergency eviction (0.0-1.0).
	CriticalMemoryThreshold float64 `koanf:"critical_memory_threshold"`
}

// DefaultTranscodeCacheConfig returns sensible defaults.
func DefaultTranscodeCacheConfig() TranscodeCacheConfig {
	return TranscodeCacheConfig{
		MaxMemoryBytes:          0,    // Auto-detect (25% of system RAM)
		MaxSegmentsPerTranscode: 50,   // ~5 minutes at 6s segments
		EvictionCheckInterval:   5 * time.Second,
		MinRetentionTime:        30 * time.Second,
		HighMemoryThreshold:     0.8,  // 80% - start evicting old segments
		CriticalMemoryThreshold: 0.95, // 95% - aggressive eviction
	}
}

// CachedSegment represents a cached transcode segment.
type CachedSegment struct {
	TranscodeID string
	SegmentKey  string
	Data        []byte
	MimeType    string
	Duration    time.Duration
	CreatedAt   time.Time
	LastAccess  time.Time
	AccessCount int64
	Size        int64

	// LRU tracking
	lruElement *list.Element
}

// TranscodeSession tracks cached segments for a transcode.
type TranscodeSession struct {
	ID          string
	MediaID     uuid.UUID
	UserID      uuid.UUID
	ProfileID   string
	CreatedAt   time.Time
	LastAccess  time.Time
	Segments    map[string]*CachedSegment
	TotalSize   int64
	IsActive    bool // Currently being watched
	Priority    int  // Higher = keep longer (1=low, 2=normal, 3=high)
}

// TranscodeCache manages cached transcoded segments with memory-aware eviction.
// Segments are only evicted when memory pressure requires it or newer
// transcodes need the space.
type TranscodeCache struct {
	mu       sync.RWMutex
	config   TranscodeCacheConfig
	sessions map[string]*TranscodeSession // transcodeID -> session

	// LRU tracking for segments (global across all sessions)
	lru       *list.List
	lruIndex  map[string]*list.Element // segmentKey -> element

	// Memory tracking
	currentBytes int64
	maxBytes     int64

	// Stats
	stats CacheStats

	logger *slog.Logger
	stopCh chan struct{}
}

// CacheStats tracks cache performance.
type CacheStats struct {
	Hits            int64 `json:"hits"`
	Misses          int64 `json:"misses"`
	Evictions       int64 `json:"evictions"`
	EvictionsByAge  int64 `json:"evictions_by_age"`
	EvictionsBySize int64 `json:"evictions_by_size"`
	CurrentBytes    int64 `json:"current_bytes"`
	MaxBytes        int64 `json:"max_bytes"`
	SessionCount    int   `json:"session_count"`
	SegmentCount    int   `json:"segment_count"`
}

// NewTranscodeCache creates a new transcode cache.
func NewTranscodeCache(config TranscodeCacheConfig, logger *slog.Logger) *TranscodeCache {
	maxBytes := config.MaxMemoryBytes
	if maxBytes == 0 {
		// Use 25% of system memory
		var memStats runtime.MemStats
		runtime.ReadMemStats(&memStats)
		maxBytes = int64(memStats.Sys) / 4
		if maxBytes < 512*1024*1024 { // Minimum 512MB
			maxBytes = 512 * 1024 * 1024
		}
	}

	tc := &TranscodeCache{
		config:    config,
		sessions:  make(map[string]*TranscodeSession),
		lru:       list.New(),
		lruIndex:  make(map[string]*list.Element),
		maxBytes:  maxBytes,
		logger:    logger.With(slog.String("component", "transcode-cache")),
		stopCh:    make(chan struct{}),
	}

	return tc
}

// Start begins background memory monitoring.
func (tc *TranscodeCache) Start(ctx context.Context) {
	go tc.memoryMonitor(ctx)
	tc.logger.Info("transcode cache started",
		"max_bytes", tc.maxBytes,
		"max_mb", tc.maxBytes/(1024*1024))
}

// Stop stops the cache.
func (tc *TranscodeCache) Stop() {
	close(tc.stopCh)
}

// CreateSession creates or returns an existing transcode session.
func (tc *TranscodeCache) CreateSession(transcodeID string, mediaID, userID uuid.UUID, profileID string) *TranscodeSession {
	tc.mu.Lock()
	defer tc.mu.Unlock()

	if session, ok := tc.sessions[transcodeID]; ok {
		session.LastAccess = time.Now()
		session.IsActive = true
		return session
	}

	session := &TranscodeSession{
		ID:         transcodeID,
		MediaID:    mediaID,
		UserID:     userID,
		ProfileID:  profileID,
		CreatedAt:  time.Now(),
		LastAccess: time.Now(),
		Segments:   make(map[string]*CachedSegment),
		IsActive:   true,
		Priority:   2, // Normal priority
	}

	tc.sessions[transcodeID] = session
	tc.stats.SessionCount++

	return session
}

// GetSession returns a session if it exists.
func (tc *TranscodeCache) GetSession(transcodeID string) (*TranscodeSession, bool) {
	tc.mu.RLock()
	defer tc.mu.RUnlock()

	session, ok := tc.sessions[transcodeID]
	if ok {
		session.LastAccess = time.Now()
	}
	return session, ok
}

// MarkSessionInactive marks a session as no longer actively watched.
// Inactive sessions have lower eviction priority than active ones.
func (tc *TranscodeCache) MarkSessionInactive(transcodeID string) {
	tc.mu.Lock()
	defer tc.mu.Unlock()

	if session, ok := tc.sessions[transcodeID]; ok {
		session.IsActive = false
		session.Priority = 1 // Lower priority for eviction
	}
}

// SetSessionPriority sets eviction priority (1=low, 2=normal, 3=high).
func (tc *TranscodeCache) SetSessionPriority(transcodeID string, priority int) {
	tc.mu.Lock()
	defer tc.mu.Unlock()

	if session, ok := tc.sessions[transcodeID]; ok {
		session.Priority = priority
	}
}

// PutSegment adds a segment to the cache.
func (tc *TranscodeCache) PutSegment(transcodeID, segmentKey string, data []byte, mimeType string, duration time.Duration) bool {
	tc.mu.Lock()
	defer tc.mu.Unlock()

	session, ok := tc.sessions[transcodeID]
	if !ok {
		return false
	}

	// Check if already cached
	if _, exists := session.Segments[segmentKey]; exists {
		return true
	}

	segmentSize := int64(len(data))

	// Check if we need to evict to make room
	if tc.currentBytes+segmentSize > tc.maxBytes {
		tc.evictForSpace(segmentSize)
	}

	// Check segment limit per session
	if len(session.Segments) >= tc.config.MaxSegmentsPerTranscode {
		tc.evictOldestFromSession(session)
	}

	// Create cached segment
	segment := &CachedSegment{
		TranscodeID: transcodeID,
		SegmentKey:  segmentKey,
		Data:        data,
		MimeType:    mimeType,
		Duration:    duration,
		CreatedAt:   time.Now(),
		LastAccess:  time.Now(),
		Size:        segmentSize,
	}

	// Add to LRU
	globalKey := transcodeID + ":" + segmentKey
	element := tc.lru.PushFront(segment)
	segment.lruElement = element
	tc.lruIndex[globalKey] = element

	// Add to session
	session.Segments[segmentKey] = segment
	session.TotalSize += segmentSize
	session.LastAccess = time.Now()

	// Update totals
	tc.currentBytes += segmentSize
	tc.stats.SegmentCount++
	tc.stats.CurrentBytes = tc.currentBytes

	return true
}

// GetSegment retrieves a segment from the cache.
func (tc *TranscodeCache) GetSegment(transcodeID, segmentKey string) (*CachedSegment, bool) {
	tc.mu.Lock()
	defer tc.mu.Unlock()

	session, ok := tc.sessions[transcodeID]
	if !ok {
		tc.stats.Misses++
		return nil, false
	}

	segment, ok := session.Segments[segmentKey]
	if !ok {
		tc.stats.Misses++
		return nil, false
	}

	// Update access tracking
	segment.LastAccess = time.Now()
	segment.AccessCount++
	session.LastAccess = time.Now()

	// Move to front of LRU
	if segment.lruElement != nil {
		tc.lru.MoveToFront(segment.lruElement)
	}

	tc.stats.Hits++
	return segment, true
}

// evictForSpace evicts segments until we have enough space.
func (tc *TranscodeCache) evictForSpace(needed int64) {
	target := tc.maxBytes - needed

	// First, evict from inactive sessions with low priority
	for tc.currentBytes > target && tc.lru.Len() > 0 {
		// Get oldest from LRU
		element := tc.lru.Back()
		if element == nil {
			break
		}

		segment := element.Value.(*CachedSegment)

		// Skip if recently accessed and not under critical pressure
		if time.Since(segment.LastAccess) < tc.config.MinRetentionTime {
			memoryPressure := float64(tc.currentBytes) / float64(tc.maxBytes)
			if memoryPressure < tc.config.CriticalMemoryThreshold {
				// Can't evict recent segments, stop trying
				break
			}
		}

		tc.evictSegment(segment)
		tc.stats.EvictionsBySize++
	}
}

// evictOldestFromSession evicts the oldest segment from a session.
func (tc *TranscodeCache) evictOldestFromSession(session *TranscodeSession) {
	var oldest *CachedSegment
	for _, seg := range session.Segments {
		if oldest == nil || seg.LastAccess.Before(oldest.LastAccess) {
			oldest = seg
		}
	}

	if oldest != nil {
		tc.evictSegment(oldest)
		tc.stats.EvictionsByAge++
	}
}

// evictSegment removes a segment from the cache.
func (tc *TranscodeCache) evictSegment(segment *CachedSegment) {
	session, ok := tc.sessions[segment.TranscodeID]
	if !ok {
		return
	}

	// Remove from session
	delete(session.Segments, segment.SegmentKey)
	session.TotalSize -= segment.Size

	// Remove from LRU
	globalKey := segment.TranscodeID + ":" + segment.SegmentKey
	if element, ok := tc.lruIndex[globalKey]; ok {
		tc.lru.Remove(element)
		delete(tc.lruIndex, globalKey)
	}

	// Update totals
	tc.currentBytes -= segment.Size
	tc.stats.SegmentCount--
	tc.stats.Evictions++
	tc.stats.CurrentBytes = tc.currentBytes

	// Remove empty sessions
	if len(session.Segments) == 0 && !session.IsActive {
		delete(tc.sessions, session.ID)
		tc.stats.SessionCount--
	}
}

// memoryMonitor periodically checks memory pressure and evicts if needed.
func (tc *TranscodeCache) memoryMonitor(ctx context.Context) {
	ticker := time.NewTicker(tc.config.EvictionCheckInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-tc.stopCh:
			return
		case <-ticker.C:
			tc.checkMemoryPressure()
		}
	}
}

// checkMemoryPressure evicts segments if memory pressure is high.
func (tc *TranscodeCache) checkMemoryPressure() {
	tc.mu.Lock()
	defer tc.mu.Unlock()

	memoryPressure := float64(tc.currentBytes) / float64(tc.maxBytes)

	if memoryPressure >= tc.config.CriticalMemoryThreshold {
		// Emergency eviction - evict 20% of cache
		target := int64(float64(tc.maxBytes) * 0.8)
		tc.evictToTarget(target, true)
		tc.logger.Warn("critical memory pressure, emergency eviction",
			"pressure", memoryPressure,
			"evicted_to_bytes", tc.currentBytes)
	} else if memoryPressure >= tc.config.HighMemoryThreshold {
		// Normal eviction - evict inactive session segments
		tc.evictInactiveSegments()
		tc.logger.Debug("high memory pressure, evicting inactive",
			"pressure", memoryPressure)
	}
}

// evictToTarget evicts until we're below target bytes.
func (tc *TranscodeCache) evictToTarget(target int64, force bool) {
	for tc.currentBytes > target && tc.lru.Len() > 0 {
		element := tc.lru.Back()
		if element == nil {
			break
		}

		segment := element.Value.(*CachedSegment)

		// If not forcing, respect minimum retention
		if !force && time.Since(segment.LastAccess) < tc.config.MinRetentionTime {
			session, ok := tc.sessions[segment.TranscodeID]
			if ok && session.IsActive {
				break // Don't evict active session segments
			}
		}

		tc.evictSegment(segment)
	}
}

// evictInactiveSegments evicts old segments from inactive sessions.
func (tc *TranscodeCache) evictInactiveSegments() {
	now := time.Now()

	for _, session := range tc.sessions {
		if session.IsActive {
			continue
		}

		// Evict segments older than retention time from inactive sessions
		for key, segment := range session.Segments {
			if now.Sub(segment.LastAccess) > tc.config.MinRetentionTime {
				tc.evictSegment(segment)
				tc.logger.Debug("evicted inactive segment",
					"transcode_id", session.ID,
					"segment", key)
			}
		}
	}
}

// CleanupSession removes all segments for a transcode session.
func (tc *TranscodeCache) CleanupSession(transcodeID string) {
	tc.mu.Lock()
	defer tc.mu.Unlock()

	session, ok := tc.sessions[transcodeID]
	if !ok {
		return
	}

	for _, segment := range session.Segments {
		tc.evictSegment(segment)
	}

	delete(tc.sessions, transcodeID)
	tc.stats.SessionCount--

	tc.logger.Debug("cleaned up session", "transcode_id", transcodeID)
}

// Stats returns cache statistics.
func (tc *TranscodeCache) Stats() CacheStats {
	tc.mu.RLock()
	defer tc.mu.RUnlock()

	return CacheStats{
		Hits:            tc.stats.Hits,
		Misses:          tc.stats.Misses,
		Evictions:       tc.stats.Evictions,
		EvictionsByAge:  tc.stats.EvictionsByAge,
		EvictionsBySize: tc.stats.EvictionsBySize,
		CurrentBytes:    tc.currentBytes,
		MaxBytes:        tc.maxBytes,
		SessionCount:    len(tc.sessions),
		SegmentCount:    tc.lru.Len(),
	}
}

// MemoryUsage returns current memory usage as a percentage.
func (tc *TranscodeCache) MemoryUsage() float64 {
	tc.mu.RLock()
	defer tc.mu.RUnlock()
	return float64(tc.currentBytes) / float64(tc.maxBytes)
}
