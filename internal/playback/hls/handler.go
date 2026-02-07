package hls

import (
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lusoris/revenge/internal/infra/cache"
	"github.com/lusoris/revenge/internal/playback"
)

// StreamHandler serves HLS manifests, segments, and subtitles via HTTP.
// Registered at /api/v1/playback/stream/{sessionId}/...
type StreamHandler struct {
	sessions        *playback.SessionManager
	masterCache     *cache.L1Cache[uuid.UUID, string]   // session → master playlist
	mediaCache      *cache.L1Cache[string, mediaEntry]   // session:profile → media playlist
	logger          *slog.Logger
}

type mediaEntry struct {
	content   string
	cachedAt  time.Time
}

// NewStreamHandler creates a new HLS stream handler.
func NewStreamHandler(sessions *playback.SessionManager, logger *slog.Logger) (*StreamHandler, error) {
	masterCache, err := cache.NewL1Cache[uuid.UUID, string](1000, 30*time.Minute)
	if err != nil {
		return nil, err
	}

	mediaCache, err := cache.NewL1Cache[string, mediaEntry](5000, 2*time.Second)
	if err != nil {
		return nil, err
	}

	return &StreamHandler{
		sessions:    sessions,
		masterCache: masterCache,
		mediaCache:  mediaCache,
		logger:      logger,
	}, nil
}

// ServeHTTP routes stream requests:
//
//	GET .../master.m3u8                          → master playlist
//	GET .../{profile}/index.m3u8                 → video media playlist
//	GET .../{profile}/seg-NNNNN.ts               → video segment (zero-copy)
//	GET .../audio/{track}/index.m3u8             → audio rendition playlist
//	GET .../audio/{track}/seg-NNNNN.ts           → audio rendition segment (zero-copy)
//	GET .../subs/{track}.vtt                     → subtitle track (full file)
func (h *StreamHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Set CORS headers for HLS.js
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")

	// Parse path: /api/v1/playback/stream/{sessionId}/...
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/playback/stream/")
	slashIdx := strings.IndexByte(path, '/')
	if slashIdx < 0 {
		http.NotFound(w, r)
		return
	}

	sessionID, err := uuid.Parse(path[:slashIdx])
	if err != nil {
		http.Error(w, "invalid session ID", http.StatusBadRequest)
		return
	}

	session, ok := h.sessions.Get(sessionID)
	if !ok {
		http.Error(w, "session not found or expired", http.StatusNotFound)
		return
	}

	// Touch session (keep alive) — non-blocking
	go h.sessions.Touch(sessionID)

	remaining := path[slashIdx+1:]

	switch {
	case remaining == "master.m3u8":
		h.serveMasterPlaylist(w, r, session)

	case strings.HasPrefix(remaining, "subs/"):
		h.serveSubtitle(w, r, session, remaining)

	case strings.HasPrefix(remaining, "audio/"):
		// Audio rendition: audio/{track}/index.m3u8 or audio/{track}/seg-NNNNN.ts
		h.serveAudioRendition(w, r, session, strings.TrimPrefix(remaining, "audio/"))

	case strings.HasSuffix(remaining, "/index.m3u8"):
		profile := strings.TrimSuffix(remaining, "/index.m3u8")
		h.serveMediaPlaylist(w, r, session, profile)

	case strings.Contains(remaining, "/seg-") && strings.HasSuffix(remaining, ".ts"):
		slashPos := strings.IndexByte(remaining, '/')
		if slashPos > 0 {
			h.serveSegment(w, r, session, remaining[:slashPos], remaining[slashPos+1:])
		} else {
			http.NotFound(w, r)
		}

	default:
		http.NotFound(w, r)
	}
}

func (h *StreamHandler) serveMasterPlaylist(w http.ResponseWriter, _ *http.Request, session *playback.Session) {
	// Check master playlist cache
	if cached, ok := h.masterCache.Get(session.ID); ok {
		w.Header().Set("Content-Type", "application/vnd.apple.mpegurl")
		w.Header().Set("Cache-Control", "no-cache")
		_, _ = w.Write([]byte(cached))
		return
	}

	// Generate master playlist
	profiles := make([]ProfileVariant, 0, len(session.TranscodeDecision.Profiles))
	for _, pd := range session.TranscodeDecision.Profiles {
		bw := estimateBandwidth(pd, 0, 0) // use defaults for unknown source bitrate
		profiles = append(profiles, ProfileVariant{
			Name:      pd.Name,
			Width:     pd.Width,
			Height:    pd.Height,
			Bandwidth: bw,
		})
	}

	audioVariants := make([]AudioVariant, 0, len(session.AudioTracks))
	for _, at := range session.AudioTracks {
		audioVariants = append(audioVariants, AudioVariant{
			Index:     at.Index,
			Name:      audioDisplayName(at),
			Language:  at.Language,
			Channels:  at.Channels,
			IsDefault: at.IsDefault,
		})
	}

	subtitles := make([]SubtitleVariant, 0, len(session.SubtitleTracks))
	for _, st := range session.SubtitleTracks {
		subtitles = append(subtitles, SubtitleVariant{
			Index:     st.Index,
			Name:      subtitleDisplayName(st),
			Language:  st.Language,
			IsDefault: st.Index == 0,
		})
	}

	playlist := GenerateMasterPlaylist(profiles, audioVariants, subtitles)
	h.masterCache.Set(session.ID, playlist)

	w.Header().Set("Content-Type", "application/vnd.apple.mpegurl")
	w.Header().Set("Cache-Control", "no-cache")
	_, _ = w.Write([]byte(playlist))
}

func (h *StreamHandler) serveMediaPlaylist(w http.ResponseWriter, _ *http.Request, session *playback.Session, profile string) {
	cacheKey := session.ID.String() + ":" + profile

	// Check media playlist cache (1s TTL)
	if entry, ok := h.mediaCache.Get(cacheKey); ok {
		w.Header().Set("Content-Type", "application/vnd.apple.mpegurl")
		w.Header().Set("Cache-Control", "no-cache")
		_, _ = w.Write([]byte(entry.content))
		return
	}

	// Read from disk (FFmpeg-generated)
	content, err := ReadMediaPlaylist(session.SegmentDir, profile)
	if err != nil {
		h.logger.Warn("media playlist not available",
			slog.String("session_id", session.ID.String()),
			slog.String("profile", profile),
			slog.String("error", err.Error()),
		)
		http.Error(w, "media playlist not ready", http.StatusServiceUnavailable)
		return
	}

	h.mediaCache.Set(cacheKey, mediaEntry{content: content, cachedAt: time.Now()})

	w.Header().Set("Content-Type", "application/vnd.apple.mpegurl")
	w.Header().Set("Cache-Control", "no-cache")
	_, _ = w.Write([]byte(content))
}

func (h *StreamHandler) serveAudioRendition(w http.ResponseWriter, r *http.Request, session *playback.Session, remaining string) {
	// remaining = "{track}/index.m3u8" or "{track}/seg-NNNNN.ts"
	slashPos := strings.IndexByte(remaining, '/')
	if slashPos < 0 {
		http.NotFound(w, r)
		return
	}

	trackStr := remaining[:slashPos]
	file := remaining[slashPos+1:]

	trackIndex, err := strconv.Atoi(trackStr)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	if file == "index.m3u8" {
		// Audio rendition playlist — reuse media playlist caching
		h.serveMediaPlaylist(w, r, session, "audio/"+trackStr)
	} else if strings.HasPrefix(file, "seg-") && strings.HasSuffix(file, ".ts") {
		// Audio rendition segment — zero-copy serve
		segPath := AudioRenditionSegmentPath(session.SegmentDir, trackIndex, file)
		w.Header().Set("Content-Type", "video/mp2t")
		w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		http.ServeFile(w, r, segPath)
	} else {
		http.NotFound(w, r)
	}
}

func (h *StreamHandler) serveSegment(w http.ResponseWriter, r *http.Request, session *playback.Session, profile, segmentFile string) {
	segPath := SegmentPath(session.SegmentDir, profile, segmentFile)

	w.Header().Set("Content-Type", "video/mp2t")
	w.Header().Set("Cache-Control", "public, max-age=31536000, immutable") // Segments are immutable
	// http.ServeFile handles If-Modified-Since, Range requests, and sendfile(2) zero-copy
	http.ServeFile(w, r, segPath)
}

func (h *StreamHandler) serveSubtitle(w http.ResponseWriter, r *http.Request, session *playback.Session, remaining string) {
	// remaining = "subs/0.vtt"
	trackStr := strings.TrimPrefix(remaining, "subs/")
	trackStr = strings.TrimSuffix(trackStr, ".vtt")
	trackIndex, err := strconv.Atoi(trackStr)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	vttPath := SubtitlePath(session.SegmentDir, trackIndex)

	w.Header().Set("Content-Type", "text/vtt")
	w.Header().Set("Cache-Control", "public, max-age=3600")
	http.ServeFile(w, r, vttPath)
}

func audioDisplayName(at playback.AudioTrackInfo) string {
	if at.Title != "" {
		return at.Title
	}
	if at.Language != "" {
		return at.Language
	}
	return "Track " + strconv.Itoa(at.Index)
}

func subtitleDisplayName(st playback.SubtitleTrackInfo) string {
	if st.Title != "" {
		return st.Title
	}
	if st.Language != "" {
		return st.Language
	}
	return "Track " + strconv.Itoa(st.Index)
}

// Close shuts down caches.
func (h *StreamHandler) Close() {
	h.masterCache.Close()
	h.mediaCache.Close()
}
