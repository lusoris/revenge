// Package playback provides persistent transcode caching on disk.
package playback

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"github.com/google/uuid"
)

// DiskCacheConfig configures the disk-based transcode cache.
type DiskCacheConfig struct {
	// BasePath is the directory for cached transcodes.
	BasePath string `koanf:"base_path"`

	// MaxSizeBytes is the maximum total cache size.
	// 0 = unlimited (use quota instead).
	MaxSizeBytes int64 `koanf:"max_size_bytes"`

	// MaxAgeHours is how long to keep cached transcodes.
	// 0 = forever (until evicted by size).
	MaxAgeHours int `koanf:"max_age_hours"`

	// EvictionCheckInterval is how often to check for eviction.
	EvictionCheckInterval time.Duration `koanf:"eviction_check_interval"`

	// MinFreeSpaceBytes ensures this much space remains free.
	MinFreeSpaceBytes int64 `koanf:"min_free_space_bytes"`

	// PerUserQuotaBytes limits cache per user (0 = no limit).
	PerUserQuotaBytes int64 `koanf:"per_user_quota_bytes"`

	// PerMediaQuotaBytes limits cache per media item (0 = no limit).
	PerMediaQuotaBytes int64 `koanf:"per_media_quota_bytes"`
}

// DefaultDiskCacheConfig returns sensible defaults.
func DefaultDiskCacheConfig() DiskCacheConfig {
	return DiskCacheConfig{
		BasePath:              "/var/cache/revenge/transcodes",
		MaxSizeBytes:          50 * 1024 * 1024 * 1024, // 50GB
		MaxAgeHours:           72,                       // 3 days
		EvictionCheckInterval: 5 * time.Minute,
		MinFreeSpaceBytes:     10 * 1024 * 1024 * 1024, // 10GB free
		PerUserQuotaBytes:     10 * 1024 * 1024 * 1024, // 10GB per user
		PerMediaQuotaBytes:    5 * 1024 * 1024 * 1024,  // 5GB per media
	}
}

// CachedTranscode represents a cached transcode on disk.
type CachedTranscode struct {
	// Identity
	ID        string    `json:"id"`
	MediaID   uuid.UUID `json:"media_id"`
	UserID    uuid.UUID `json:"user_id"`
	ProfileID string    `json:"profile_id"`

	// Cache key components
	CacheKey string `json:"cache_key"`

	// File info
	BasePath     string `json:"base_path"`
	ManifestPath string `json:"manifest_path"`
	SegmentCount int    `json:"segment_count"`
	TotalSize    int64  `json:"total_size"`

	// Metadata
	CreatedAt    time.Time `json:"created_at"`
	LastAccess   time.Time `json:"last_access"`
	AccessCount  int64     `json:"access_count"`
	IsComplete   bool      `json:"is_complete"`
	SourceHash   string    `json:"source_hash"` // Hash of source file for invalidation

	// Quotas
	UserQuotaUsed  int64 `json:"user_quota_used"`
	MediaQuotaUsed int64 `json:"media_quota_used"`
}

// DiskCache manages persistent transcode caching.
type DiskCache struct {
	mu     sync.RWMutex
	config DiskCacheConfig
	logger *slog.Logger

	// Index
	transcodes map[string]*CachedTranscode // cacheKey -> transcode
	byMedia    map[uuid.UUID][]string      // mediaID -> cacheKeys
	byUser     map[uuid.UUID][]string      // userID -> cacheKeys

	// Quota tracking
	userUsage  map[uuid.UUID]int64
	mediaUsage map[uuid.UUID]int64
	totalUsage int64

	// Background
	stopCh chan struct{}
}

// NewDiskCache creates a new disk cache.
func NewDiskCache(config DiskCacheConfig, logger *slog.Logger) (*DiskCache, error) {
	// Ensure base path exists
	if err := os.MkdirAll(config.BasePath, 0755); err != nil {
		return nil, fmt.Errorf("create cache directory: %w", err)
	}

	dc := &DiskCache{
		config:     config,
		logger:     logger.With(slog.String("component", "disk-cache")),
		transcodes: make(map[string]*CachedTranscode),
		byMedia:    make(map[uuid.UUID][]string),
		byUser:     make(map[uuid.UUID][]string),
		userUsage:  make(map[uuid.UUID]int64),
		mediaUsage: make(map[uuid.UUID]int64),
		stopCh:     make(chan struct{}),
	}

	// Load existing index
	if err := dc.loadIndex(); err != nil {
		logger.Warn("failed to load cache index, starting fresh", "error", err)
	}

	return dc, nil
}

// Start begins background maintenance.
func (dc *DiskCache) Start(ctx context.Context) {
	go dc.maintenanceLoop(ctx)
	dc.logger.Info("disk cache started",
		"path", dc.config.BasePath,
		"max_size_gb", dc.config.MaxSizeBytes/(1024*1024*1024),
		"current_size_gb", dc.totalUsage/(1024*1024*1024))
}

// Stop stops the cache.
func (dc *DiskCache) Stop() {
	close(dc.stopCh)
	dc.saveIndex()
}

// GenerateCacheKey creates a deterministic cache key for a transcode.
func (dc *DiskCache) GenerateCacheKey(mediaID uuid.UUID, profileID string, sourceHash string) string {
	h := sha256.New()
	h.Write([]byte(mediaID.String()))
	h.Write([]byte(profileID))
	h.Write([]byte(sourceHash))
	return hex.EncodeToString(h.Sum(nil))[:32]
}

// GetCached returns a cached transcode if available and valid.
func (dc *DiskCache) GetCached(cacheKey string) (*CachedTranscode, bool) {
	dc.mu.Lock()
	defer dc.mu.Unlock()

	tc, ok := dc.transcodes[cacheKey]
	if !ok {
		return nil, false
	}

	// Verify files still exist
	if _, err := os.Stat(tc.ManifestPath); err != nil {
		dc.removeCachedLocked(cacheKey)
		return nil, false
	}

	// Update access stats
	tc.LastAccess = time.Now()
	tc.AccessCount++

	dc.logger.Debug("cache hit",
		"cache_key", cacheKey,
		"media_id", tc.MediaID,
		"access_count", tc.AccessCount)

	return tc, true
}

// CreateCacheEntry creates a new cache entry for an upcoming transcode.
func (dc *DiskCache) CreateCacheEntry(mediaID, userID uuid.UUID, profileID, sourceHash string) (*CachedTranscode, error) {
	dc.mu.Lock()
	defer dc.mu.Unlock()

	cacheKey := dc.GenerateCacheKey(mediaID, profileID, sourceHash)

	// Check if already exists
	if existing, ok := dc.transcodes[cacheKey]; ok {
		existing.LastAccess = time.Now()
		return existing, nil
	}

	// Check quotas
	if err := dc.checkQuotasLocked(userID, mediaID, 0); err != nil {
		return nil, err
	}

	// Create directory structure
	cachePath := filepath.Join(dc.config.BasePath, cacheKey[:2], cacheKey)
	if err := os.MkdirAll(cachePath, 0755); err != nil {
		return nil, fmt.Errorf("create cache directory: %w", err)
	}

	tc := &CachedTranscode{
		ID:           cacheKey,
		MediaID:      mediaID,
		UserID:       userID,
		ProfileID:    profileID,
		CacheKey:     cacheKey,
		BasePath:     cachePath,
		ManifestPath: filepath.Join(cachePath, "master.m3u8"),
		CreatedAt:    time.Now(),
		LastAccess:   time.Now(),
		SourceHash:   sourceHash,
	}

	dc.transcodes[cacheKey] = tc
	dc.byMedia[mediaID] = append(dc.byMedia[mediaID], cacheKey)
	dc.byUser[userID] = append(dc.byUser[userID], cacheKey)

	return tc, nil
}

// WriteSegment writes a segment to the cache.
func (dc *DiskCache) WriteSegment(cacheKey string, segmentName string, data []byte) error {
	dc.mu.Lock()
	tc, ok := dc.transcodes[cacheKey]
	if !ok {
		dc.mu.Unlock()
		return errors.New("cache entry not found")
	}

	segmentPath := filepath.Join(tc.BasePath, segmentName)
	segmentSize := int64(len(data))

	// Check if adding this would exceed quotas
	if err := dc.checkQuotasLocked(tc.UserID, tc.MediaID, segmentSize); err != nil {
		dc.mu.Unlock()
		// Evict and retry
		dc.evictForSpace(segmentSize)
		dc.mu.Lock()
		if err := dc.checkQuotasLocked(tc.UserID, tc.MediaID, segmentSize); err != nil {
			dc.mu.Unlock()
			return err
		}
	}
	dc.mu.Unlock()

	// Write file
	if err := os.WriteFile(segmentPath, data, 0644); err != nil {
		return fmt.Errorf("write segment: %w", err)
	}

	// Update tracking
	dc.mu.Lock()
	defer dc.mu.Unlock()

	tc.SegmentCount++
	tc.TotalSize += segmentSize
	tc.LastAccess = time.Now()

	dc.totalUsage += segmentSize
	dc.userUsage[tc.UserID] += segmentSize
	dc.mediaUsage[tc.MediaID] += segmentSize

	return nil
}

// WriteManifest writes the manifest file.
func (dc *DiskCache) WriteManifest(cacheKey string, data []byte) error {
	dc.mu.RLock()
	tc, ok := dc.transcodes[cacheKey]
	dc.mu.RUnlock()

	if !ok {
		return errors.New("cache entry not found")
	}

	if err := os.WriteFile(tc.ManifestPath, data, 0644); err != nil {
		return fmt.Errorf("write manifest: %w", err)
	}

	return nil
}

// MarkComplete marks a transcode as complete.
func (dc *DiskCache) MarkComplete(cacheKey string) {
	dc.mu.Lock()
	defer dc.mu.Unlock()

	if tc, ok := dc.transcodes[cacheKey]; ok {
		tc.IsComplete = true
		dc.logger.Info("transcode cached",
			"cache_key", cacheKey,
			"segments", tc.SegmentCount,
			"size_mb", tc.TotalSize/(1024*1024))
	}
}

// ReadSegment reads a segment from cache.
func (dc *DiskCache) ReadSegment(cacheKey, segmentName string) ([]byte, error) {
	dc.mu.RLock()
	tc, ok := dc.transcodes[cacheKey]
	dc.mu.RUnlock()

	if !ok {
		return nil, errors.New("cache entry not found")
	}

	segmentPath := filepath.Join(tc.BasePath, segmentName)
	data, err := os.ReadFile(segmentPath)
	if err != nil {
		return nil, fmt.Errorf("read segment: %w", err)
	}

	// Update access time async
	go func() {
		dc.mu.Lock()
		defer dc.mu.Unlock()
		if tc, ok := dc.transcodes[cacheKey]; ok {
			tc.LastAccess = time.Now()
			tc.AccessCount++
		}
	}()

	return data, nil
}

// ReadManifest reads the manifest from cache.
func (dc *DiskCache) ReadManifest(cacheKey string) ([]byte, error) {
	dc.mu.RLock()
	tc, ok := dc.transcodes[cacheKey]
	dc.mu.RUnlock()

	if !ok {
		return nil, errors.New("cache entry not found")
	}

	return os.ReadFile(tc.ManifestPath)
}

// checkQuotasLocked checks if adding size would exceed quotas (must hold lock).
func (dc *DiskCache) checkQuotasLocked(userID, mediaID uuid.UUID, addSize int64) error {
	// Global quota
	if dc.config.MaxSizeBytes > 0 && dc.totalUsage+addSize > dc.config.MaxSizeBytes {
		return fmt.Errorf("global cache quota exceeded (%d/%d bytes)",
			dc.totalUsage+addSize, dc.config.MaxSizeBytes)
	}

	// Per-user quota
	if dc.config.PerUserQuotaBytes > 0 {
		if dc.userUsage[userID]+addSize > dc.config.PerUserQuotaBytes {
			return fmt.Errorf("user cache quota exceeded (%d/%d bytes)",
				dc.userUsage[userID]+addSize, dc.config.PerUserQuotaBytes)
		}
	}

	// Per-media quota
	if dc.config.PerMediaQuotaBytes > 0 {
		if dc.mediaUsage[mediaID]+addSize > dc.config.PerMediaQuotaBytes {
			return fmt.Errorf("media cache quota exceeded (%d/%d bytes)",
				dc.mediaUsage[mediaID]+addSize, dc.config.PerMediaQuotaBytes)
		}
	}

	return nil
}

// evictForSpace evicts old transcodes to make room.
func (dc *DiskCache) evictForSpace(needed int64) {
	dc.mu.Lock()
	defer dc.mu.Unlock()

	target := dc.config.MaxSizeBytes - needed
	if target < 0 {
		target = dc.config.MaxSizeBytes / 2 // Evict to 50% if needed > max
	}

	dc.evictToTargetLocked(target)
}

// evictToTargetLocked evicts until under target size (must hold lock).
func (dc *DiskCache) evictToTargetLocked(target int64) {
	// Sort by last access (oldest first)
	type tcEntry struct {
		key string
		tc  *CachedTranscode
	}

	entries := make([]tcEntry, 0, len(dc.transcodes))
	for key, tc := range dc.transcodes {
		entries = append(entries, tcEntry{key, tc})
	}

	sort.Slice(entries, func(i, j int) bool {
		// Incomplete transcodes first (they're probably abandoned)
		if !entries[i].tc.IsComplete && entries[j].tc.IsComplete {
			return true
		}
		if entries[i].tc.IsComplete && !entries[j].tc.IsComplete {
			return false
		}
		// Then by last access
		return entries[i].tc.LastAccess.Before(entries[j].tc.LastAccess)
	})

	for _, entry := range entries {
		if dc.totalUsage <= target {
			break
		}
		dc.removeCachedLocked(entry.key)
	}
}

// removeCachedLocked removes a cached transcode (must hold lock).
func (dc *DiskCache) removeCachedLocked(cacheKey string) {
	tc, ok := dc.transcodes[cacheKey]
	if !ok {
		return
	}

	// Update quotas
	dc.totalUsage -= tc.TotalSize
	dc.userUsage[tc.UserID] -= tc.TotalSize
	dc.mediaUsage[tc.MediaID] -= tc.TotalSize

	// Remove from indices
	delete(dc.transcodes, cacheKey)
	dc.byMedia[tc.MediaID] = dc.removeFromSlice(dc.byMedia[tc.MediaID], cacheKey)
	dc.byUser[tc.UserID] = dc.removeFromSlice(dc.byUser[tc.UserID], cacheKey)

	// Remove files async
	go func() {
		if err := os.RemoveAll(tc.BasePath); err != nil {
			dc.logger.Error("failed to remove cached files", "path", tc.BasePath, "error", err)
		}
	}()

	dc.logger.Debug("evicted cached transcode",
		"cache_key", cacheKey,
		"size_mb", tc.TotalSize/(1024*1024),
		"age", time.Since(tc.CreatedAt))
}

func (dc *DiskCache) removeFromSlice(slice []string, val string) []string {
	for i, v := range slice {
		if v == val {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}

// InvalidateByMedia removes all cached transcodes for a media item.
func (dc *DiskCache) InvalidateByMedia(mediaID uuid.UUID) {
	dc.mu.Lock()
	defer dc.mu.Unlock()

	keys := dc.byMedia[mediaID]
	for _, key := range keys {
		dc.removeCachedLocked(key)
	}
	delete(dc.byMedia, mediaID)
}

// InvalidateBySourceHash removes transcodes with outdated source hash.
func (dc *DiskCache) InvalidateBySourceHash(mediaID uuid.UUID, newHash string) {
	dc.mu.Lock()
	defer dc.mu.Unlock()

	for _, key := range dc.byMedia[mediaID] {
		if tc, ok := dc.transcodes[key]; ok && tc.SourceHash != newHash {
			dc.removeCachedLocked(key)
		}
	}
}

// maintenanceLoop runs periodic maintenance.
func (dc *DiskCache) maintenanceLoop(ctx context.Context) {
	ticker := time.NewTicker(dc.config.EvictionCheckInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-dc.stopCh:
			return
		case <-ticker.C:
			dc.runMaintenance()
		}
	}
}

// runMaintenance performs cache maintenance.
func (dc *DiskCache) runMaintenance() {
	dc.mu.Lock()
	defer dc.mu.Unlock()

	now := time.Now()
	maxAge := time.Duration(dc.config.MaxAgeHours) * time.Hour

	// Remove expired
	if dc.config.MaxAgeHours > 0 {
		for key, tc := range dc.transcodes {
			if now.Sub(tc.LastAccess) > maxAge {
				dc.removeCachedLocked(key)
			}
		}
	}

	// Remove incomplete transcodes older than 1 hour
	for key, tc := range dc.transcodes {
		if !tc.IsComplete && now.Sub(tc.CreatedAt) > time.Hour {
			dc.removeCachedLocked(key)
		}
	}

	// Check disk space
	dc.ensureFreeSpace()

	// Save index
	dc.saveIndexLocked()
}

// ensureFreeSpace evicts if disk space is low.
func (dc *DiskCache) ensureFreeSpace() {
	// Get free space (platform-specific, simplified here)
	// In production, use syscall.Statfs on Linux
	if dc.config.MinFreeSpaceBytes > 0 {
		// Simplified: evict 10% if over 90% of quota
		if float64(dc.totalUsage)/float64(dc.config.MaxSizeBytes) > 0.9 {
			target := int64(float64(dc.config.MaxSizeBytes) * 0.8)
			dc.evictToTargetLocked(target)
		}
	}
}

// loadIndex loads the cache index from disk.
func (dc *DiskCache) loadIndex() error {
	indexPath := filepath.Join(dc.config.BasePath, "index.json")
	data, err := os.ReadFile(indexPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // Fresh start
		}
		return err
	}

	var index struct {
		Transcodes []*CachedTranscode `json:"transcodes"`
	}

	if err := json.Unmarshal(data, &index); err != nil {
		return err
	}

	// Rebuild in-memory state
	for _, tc := range index.Transcodes {
		// Verify files exist
		if _, err := os.Stat(tc.BasePath); err != nil {
			continue // Skip missing
		}

		dc.transcodes[tc.CacheKey] = tc
		dc.byMedia[tc.MediaID] = append(dc.byMedia[tc.MediaID], tc.CacheKey)
		dc.byUser[tc.UserID] = append(dc.byUser[tc.UserID], tc.CacheKey)
		dc.totalUsage += tc.TotalSize
		dc.userUsage[tc.UserID] += tc.TotalSize
		dc.mediaUsage[tc.MediaID] += tc.TotalSize
	}

	dc.logger.Info("loaded cache index",
		"entries", len(dc.transcodes),
		"total_size_gb", dc.totalUsage/(1024*1024*1024))

	return nil
}

// saveIndex saves the cache index to disk.
func (dc *DiskCache) saveIndex() {
	dc.mu.RLock()
	defer dc.mu.RUnlock()
	dc.saveIndexLocked()
}

func (dc *DiskCache) saveIndexLocked() {
	transcodes := make([]*CachedTranscode, 0, len(dc.transcodes))
	for _, tc := range dc.transcodes {
		transcodes = append(transcodes, tc)
	}

	index := struct {
		Transcodes []*CachedTranscode `json:"transcodes"`
		SavedAt    time.Time          `json:"saved_at"`
	}{
		Transcodes: transcodes,
		SavedAt:    time.Now(),
	}

	data, err := json.MarshalIndent(index, "", "  ")
	if err != nil {
		dc.logger.Error("failed to marshal index", "error", err)
		return
	}

	indexPath := filepath.Join(dc.config.BasePath, "index.json")
	if err := os.WriteFile(indexPath, data, 0644); err != nil {
		dc.logger.Error("failed to save index", "error", err)
	}
}

// Stats returns cache statistics.
func (dc *DiskCache) Stats() DiskCacheStats {
	dc.mu.RLock()
	defer dc.mu.RUnlock()

	var completeCount, incompleteCount int
	for _, tc := range dc.transcodes {
		if tc.IsComplete {
			completeCount++
		} else {
			incompleteCount++
		}
	}

	return DiskCacheStats{
		TotalEntries:     len(dc.transcodes),
		CompleteEntries:  completeCount,
		IncompleteEntries: incompleteCount,
		TotalSizeBytes:   dc.totalUsage,
		MaxSizeBytes:     dc.config.MaxSizeBytes,
		UsagePercent:     float64(dc.totalUsage) / float64(dc.config.MaxSizeBytes) * 100,
		UniqueMedia:      len(dc.byMedia),
		UniqueUsers:      len(dc.byUser),
	}
}

// DiskCacheStats contains cache statistics.
type DiskCacheStats struct {
	TotalEntries      int     `json:"total_entries"`
	CompleteEntries   int     `json:"complete_entries"`
	IncompleteEntries int     `json:"incomplete_entries"`
	TotalSizeBytes    int64   `json:"total_size_bytes"`
	MaxSizeBytes      int64   `json:"max_size_bytes"`
	UsagePercent      float64 `json:"usage_percent"`
	UniqueMedia       int     `json:"unique_media"`
	UniqueUsers       int     `json:"unique_users"`
}

// OpenSegmentReader opens a segment for streaming.
func (dc *DiskCache) OpenSegmentReader(cacheKey, segmentName string) (io.ReadCloser, int64, error) {
	dc.mu.RLock()
	tc, ok := dc.transcodes[cacheKey]
	dc.mu.RUnlock()

	if !ok {
		return nil, 0, errors.New("cache entry not found")
	}

	segmentPath := filepath.Join(tc.BasePath, segmentName)
	file, err := os.Open(segmentPath)
	if err != nil {
		return nil, 0, err
	}

	stat, err := file.Stat()
	if err != nil {
		file.Close()
		return nil, 0, err
	}

	return file, stat.Size(), nil
}
