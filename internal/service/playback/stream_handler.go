// Package playback provides playback session management and streaming.
package playback

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

// StreamHandler manages HLS/DASH stream delivery with buffering.
// It coordinates between Blackbeard (transcoder) and clients.
type StreamHandler struct {
	transcoder *TranscoderClient
	fileServer *MediaFileServer
	sessions   *SessionManager
	buffers    map[string]*StreamBuffer // transcodeID -> buffer
	buffersMu  sync.RWMutex
	logger     *slog.Logger
	config     StreamHandlerConfig
}

// StreamHandlerConfig configures the stream handler.
type StreamHandlerConfig struct {
	// Buffer settings
	MaxSegments       int           // Max segments per buffer
	MinBufferDuration time.Duration // Minimum buffer duration
	MaxBufferDuration time.Duration // Maximum buffer duration
	PrefetchCount     int           // Segments to prefetch ahead
	MaxRetries        int           // Retries per segment

	// Cleanup
	IdleTimeout time.Duration // Remove buffer after idle
}

// DefaultStreamHandlerConfig returns sensible defaults.
func DefaultStreamHandlerConfig() StreamHandlerConfig {
	return StreamHandlerConfig{
		MaxSegments:       8,
		MinBufferDuration: 15 * time.Second,
		MaxBufferDuration: 60 * time.Second,
		PrefetchCount:     3,
		MaxRetries:        3,
		IdleTimeout:       5 * time.Minute,
	}
}

// NewStreamHandler creates a new stream handler.
func NewStreamHandler(
	transcoder *TranscoderClient,
	fileServer *MediaFileServer,
	sessions *SessionManager,
	logger *slog.Logger,
	config StreamHandlerConfig,
) *StreamHandler {
	return &StreamHandler{
		transcoder: transcoder,
		fileServer: fileServer,
		sessions:   sessions,
		buffers:    make(map[string]*StreamBuffer),
		logger:     logger,
		config:     config,
	}
}

// ServeManifest serves the HLS/DASH manifest for a playback session.
func (h *StreamHandler) ServeManifest(w http.ResponseWriter, r *http.Request) {
	sessionID := r.PathValue("sessionID")
	if sessionID == "" {
		http.Error(w, "session ID required", http.StatusBadRequest)
		return
	}

	sessionUUID, err := uuid.Parse(sessionID)
	if err != nil {
		http.Error(w, "invalid session ID", http.StatusBadRequest)
		return
	}

	session, ok := h.sessions.GetSession(sessionUUID)
	if !ok {
		h.logger.Error("session not found", "sessionID", sessionID)
		http.Error(w, "session not found", http.StatusNotFound)
		return
	}

	// Get or start transcode
	transcodeID := session.TranscodeID
	if transcodeID == "" {
		// No transcode needed - direct play
		h.serveDirectManifest(w, r, session)
		return
	}

	manifest, contentType, err := h.transcoder.GetManifest(r.Context(), transcodeID)
	if err != nil {
		h.logger.Error("failed to get manifest", "transcodeID", transcodeID, "error", err)
		http.Error(w, "failed to get manifest", http.StatusBadGateway)
		return
	}

	// Rewrite segment URLs to go through Revenge
	manifest = h.rewriteManifest(manifest, sessionID)

	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	if _, err := w.Write(manifest); err != nil {
		h.logger.Error("failed to write manifest", "error", err)
	}
}

// ServeSegment serves a single segment with buffering.
func (h *StreamHandler) ServeSegment(w http.ResponseWriter, r *http.Request) {
	sessionID := r.PathValue("sessionID")
	segmentPath := r.PathValue("segment")

	if sessionID == "" || segmentPath == "" {
		http.Error(w, "session ID and segment required", http.StatusBadRequest)
		return
	}

	sessionUUID, err := uuid.Parse(sessionID)
	if err != nil {
		http.Error(w, "invalid session ID", http.StatusBadRequest)
		return
	}

	session, ok := h.sessions.GetSession(sessionUUID)
	if !ok {
		http.Error(w, "session not found", http.StatusNotFound)
		return
	}

	transcodeID := session.TranscodeID
	if transcodeID == "" {
		http.Error(w, "no transcode for direct play", http.StatusBadRequest)
		return
	}

	// Get or create buffer for this transcode
	buffer := h.getOrCreateBuffer(transcodeID)

	// Try to get from buffer first
	if segment, bufOK := buffer.Get(segmentPath); bufOK {
		h.logger.Debug("serving from buffer", "segment", segmentPath)
		w.Header().Set("Content-Type", segment.MimeType)
		w.Header().Set("X-Served-From", "buffer")
		if _, err := w.Write(segment.Data); err != nil {
			h.logger.Error("failed to write segment from buffer", "error", err)
		}
		h.triggerPrefetch(transcodeID, segmentPath, buffer)
		return
	}

	// Fetch with retry
	data, mimeType, fetchErr := h.fetchWithRetry(r.Context(), transcodeID, segmentPath)
	if fetchErr != nil {
		h.logger.Error("failed to fetch segment", "segment", segmentPath, "error", fetchErr)
		http.Error(w, "failed to fetch segment", http.StatusBadGateway)
		return
	}

	// Store in buffer
	buffer.Add(segmentPath, data, mimeType)

	// Serve
	w.Header().Set("Content-Type", mimeType)
	w.Header().Set("X-Served-From", "origin")
	if _, err := w.Write(data); err != nil {
		h.logger.Error("failed to write segment", "error", err)
	}

	// Trigger prefetch of next segments
	h.triggerPrefetch(transcodeID, segmentPath, buffer)
}

// fetchWithRetry fetches a segment with retry logic.
func (h *StreamHandler) fetchWithRetry(ctx context.Context, transcodeID, segmentPath string) ([]byte, string, error) {
	var lastErr error

	for attempt := 0; attempt <= h.config.MaxRetries; attempt++ {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return nil, "", ctx.Err()
			case <-time.After(time.Duration(attempt) * 200 * time.Millisecond):
			}
		}

		data, mimeType, err := h.transcoder.FetchSegment(ctx, transcodeID, segmentPath)
		if err != nil {
			lastErr = err
			h.logger.Warn("segment fetch failed, retrying",
				"segment", segmentPath,
				"attempt", attempt+1,
				"error", err)
			continue
		}

		return data, mimeType, nil
	}

	return nil, "", fmt.Errorf("failed after %d attempts: %w", h.config.MaxRetries+1, lastErr)
}

// getOrCreateBuffer gets or creates a buffer for a transcode.
func (h *StreamHandler) getOrCreateBuffer(transcodeID string) *StreamBuffer {
	h.buffersMu.RLock()
	buf, ok := h.buffers[transcodeID]
	h.buffersMu.RUnlock()

	if ok {
		return buf
	}

	h.buffersMu.Lock()
	defer h.buffersMu.Unlock()

	// Double-check
	if buf, ok = h.buffers[transcodeID]; ok {
		return buf
	}

	buf = NewStreamBuffer(BufferConfig{
		SegmentBufferSize: h.config.MaxSegments,
		MinBufferDuration: h.config.MinBufferDuration,
		MaxBufferDuration: h.config.MaxBufferDuration,
	})
	h.buffers[transcodeID] = buf

	return buf
}

// triggerPrefetch prefetches upcoming segments in background.
func (h *StreamHandler) triggerPrefetch(transcodeID, currentSegment string, buffer *StreamBuffer) {
	// Extract segment number and predict next segments
	nextSegments := h.predictNextSegments(currentSegment, h.config.PrefetchCount)

	for _, seg := range nextSegments {
		if _, ok := buffer.Get(seg); ok {
			continue // Already buffered
		}

		// Fetch in background
		go func(segment string) {
			fetchCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			data, mimeType, err := h.transcoder.FetchSegment(fetchCtx, transcodeID, segment)
			if err != nil {
				h.logger.Debug("prefetch failed", "segment", segment, "error", err)
				return
			}

			buffer.Add(segment, data, mimeType)
			h.logger.Debug("prefetched segment", "segment", segment)
		}(seg)
	}
}

// predictNextSegments predicts the next N segment names.
// Supports both HLS (.ts) and DASH (.m4s) naming patterns.
func (h *StreamHandler) predictNextSegments(current string, count int) []string {
	// Common patterns:
	// HLS: segment_0.ts, segment_1.ts, ...
	// DASH: chunk-stream0-00001.m4s, chunk-stream0-00002.m4s, ...

	segNumRe := regexp.MustCompile(`(\d+)\.(ts|m4s)$`)
	matches := segNumRe.FindStringSubmatch(current)
	if matches == nil {
		return nil
	}

	numStr := matches[1]
	ext := matches[2]

	var num int
	fmt.Sscanf(numStr, "%d", &num)

	// Determine padding
	padding := len(numStr)

	prefix := current[:len(current)-len(matches[0])]

	result := make([]string, 0, count)
	for i := 1; i <= count; i++ {
		next := fmt.Sprintf("%s%0*d.%s", prefix, padding, num+i, ext)
		result = append(result, next)
	}

	return result
}

// rewriteManifest rewrites segment URLs in manifest to route through Revenge.
func (h *StreamHandler) rewriteManifest(manifest []byte, sessionID string) []byte {
	content := string(manifest)

	// Replace relative segment URLs with session-scoped URLs
	// Example: segment_0.ts -> /api/v1/playback/{sessionID}/segment/segment_0.ts

	lines := strings.Split(content, "\n")
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// This is a segment URL
		if strings.HasSuffix(line, ".ts") || strings.HasSuffix(line, ".m4s") ||
			strings.Contains(line, "segment") || strings.Contains(line, "chunk") {
			lines[i] = fmt.Sprintf("/api/v1/playback/%s/segment/%s", sessionID, line)
		}
	}

	return []byte(strings.Join(lines, "\n"))
}

// serveDirectManifest serves manifest for direct play (no transcode).
func (h *StreamHandler) serveDirectManifest(w http.ResponseWriter, r *http.Request, session *Session) {
	// For direct play, we generate a simple manifest pointing to the file
	// This allows seeking with byte-range requests

	// TODO: Implement direct play manifest generation
	http.Error(w, "direct play not implemented", http.StatusNotImplemented)
}

// CleanupIdleBuffers removes buffers that haven't been accessed recently.
func (h *StreamHandler) CleanupIdleBuffers(ctx context.Context) {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			h.buffersMu.Lock()
			for id, buf := range h.buffers {
				if time.Since(buf.LastAccess()) > h.config.IdleTimeout {
					delete(h.buffers, id)
					h.logger.Debug("cleaned up idle buffer", "transcodeID", id)
				}
			}
			h.buffersMu.Unlock()
		}
	}
}

// RegisterRoutes registers stream handling routes.
// Pattern: /api/v1/playback/{sessionID}/...
func (h *StreamHandler) RegisterRoutes(mux *http.ServeMux) {
	// Manifest
	mux.HandleFunc("GET /api/v1/playback/{sessionID}/manifest", h.ServeManifest)
	mux.HandleFunc("GET /api/v1/playback/{sessionID}/master.m3u8", h.ServeManifest)

	// Segments
	mux.HandleFunc("GET /api/v1/playback/{sessionID}/segment/{segment...}", h.ServeSegment)

	// Internal routes for Blackbeard
	mux.HandleFunc("GET /internal/stream/{mediaID}", h.fileServer.InternalStreamHandler())
	mux.HandleFunc("GET /internal/probe/{mediaID}", h.fileServer.InternalProbeHandler())
}
