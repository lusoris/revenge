// Package playback provides raw media file streaming to Blackbeard.
package playback

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

// MediaFileServer serves raw media files to Blackbeard with HTTP Range support.
type MediaFileServer struct {
	logger     *slog.Logger
	authTokens map[string]*StreamToken // Internal auth tokens
}

// StreamToken authorizes Blackbeard to fetch a specific file.
type StreamToken struct {
	Token     string
	MediaID   uuid.UUID
	FilePath  string
	ExpiresAt time.Time
	SessionID uuid.UUID
}

// NewMediaFileServer creates a new media file server.
func NewMediaFileServer(logger *slog.Logger) *MediaFileServer {
	return &MediaFileServer{
		logger:     logger.With(slog.String("component", "media_file_server")),
		authTokens: make(map[string]*StreamToken),
	}
}

// CreateStreamToken creates a temporary token for Blackbeard to access a file.
func (s *MediaFileServer) CreateStreamToken(mediaID uuid.UUID, filePath string, sessionID uuid.UUID, ttl time.Duration) *StreamToken {
	token := &StreamToken{
		Token:     uuid.New().String(),
		MediaID:   mediaID,
		FilePath:  filePath,
		ExpiresAt: time.Now().Add(ttl),
		SessionID: sessionID,
	}
	s.authTokens[token.Token] = token
	return token
}

// ValidateToken checks if a token is valid.
func (s *MediaFileServer) ValidateToken(token string) (*StreamToken, bool) {
	st, ok := s.authTokens[token]
	if !ok {
		return nil, false
	}
	if time.Now().After(st.ExpiresAt) {
		delete(s.authTokens, token)
		return nil, false
	}
	return st, true
}

// RevokeToken removes a token.
func (s *MediaFileServer) RevokeToken(token string) {
	delete(s.authTokens, token)
}

// ServeFile serves a media file with HTTP Range support for streaming.
// This is the handler for Blackbeard to fetch raw files.
func (s *MediaFileServer) ServeFile(w http.ResponseWriter, r *http.Request) {
	// Extract token from Authorization header or query param
	token := extractToken(r)
	if token == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	st, ok := s.ValidateToken(token)
	if !ok {
		http.Error(w, "invalid or expired token", http.StatusUnauthorized)
		return
	}

	// Open file
	file, err := os.Open(st.FilePath)
	if err != nil {
		s.logger.Error("failed to open file",
			slog.String("path", st.FilePath),
			slog.Any("error", err))
		http.Error(w, "file not found", http.StatusNotFound)
		return
	}
	defer file.Close()

	// Get file info
	stat, err := file.Stat()
	if err != nil {
		http.Error(w, "failed to stat file", http.StatusInternalServerError)
		return
	}

	fileSize := stat.Size()
	fileName := filepath.Base(st.FilePath)

	// Set content type based on extension
	contentType := getContentType(fileName)
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Accept-Ranges", "bytes")

	// Handle Range request
	rangeHeader := r.Header.Get("Range")
	if rangeHeader != "" {
		s.serveRangeRequest(w, file, fileSize, rangeHeader)
		return
	}

	// Full file request
	w.Header().Set("Content-Length", strconv.FormatInt(fileSize, 10))
	w.WriteHeader(http.StatusOK)

	// Stream file
	s.streamFile(r.Context(), file, w, 0, fileSize)
}

// serveRangeRequest handles HTTP Range requests for seeking.
func (s *MediaFileServer) serveRangeRequest(w http.ResponseWriter, file *os.File, fileSize int64, rangeHeader string) {
	// Parse Range header: "bytes=start-end" or "bytes=start-"
	rangeHeader = strings.TrimPrefix(rangeHeader, "bytes=")
	parts := strings.Split(rangeHeader, "-")

	if len(parts) != 2 {
		http.Error(w, "invalid range", http.StatusRequestedRangeNotSatisfiable)
		return
	}

	var start, end int64
	var err error

	if parts[0] != "" {
		start, err = strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			http.Error(w, "invalid range start", http.StatusRequestedRangeNotSatisfiable)
			return
		}
	}

	if parts[1] != "" {
		end, err = strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			http.Error(w, "invalid range end", http.StatusRequestedRangeNotSatisfiable)
			return
		}
	} else {
		// Open-ended range: bytes=1000-
		end = fileSize - 1
	}

	// Validate range
	if start < 0 || start >= fileSize || end >= fileSize || start > end {
		w.Header().Set("Content-Range", fmt.Sprintf("bytes */%d", fileSize))
		http.Error(w, "range not satisfiable", http.StatusRequestedRangeNotSatisfiable)
		return
	}

	contentLength := end - start + 1

	w.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, fileSize))
	w.Header().Set("Content-Length", strconv.FormatInt(contentLength, 10))
	w.WriteHeader(http.StatusPartialContent)

	// Seek to start position
	if _, err := file.Seek(start, io.SeekStart); err != nil {
		s.logger.Error("failed to seek", slog.Any("error", err))
		return
	}

	// Stream the range
	s.streamFile(context.Background(), file, w, start, contentLength)
}

// streamFile streams file content with chunked transfer.
func (s *MediaFileServer) streamFile(ctx context.Context, file *os.File, w io.Writer, start, length int64) {
	const chunkSize = 64 * 1024 // 64KB chunks

	buf := make([]byte, chunkSize)
	remaining := length

	for remaining > 0 {
		select {
		case <-ctx.Done():
			return
		default:
		}

		readSize := chunkSize
		if remaining < int64(chunkSize) {
			readSize = int(remaining)
		}

		n, err := file.Read(buf[:readSize])
		if err != nil {
			if err != io.EOF {
				s.logger.Error("read error", slog.Any("error", err))
			}
			return
		}

		written, err := w.Write(buf[:n])
		if err != nil {
			// Client disconnected
			return
		}

		remaining -= int64(written)

		// Flush if writer supports it (for streaming)
		if flusher, ok := w.(http.Flusher); ok {
			flusher.Flush()
		}
	}
}

// extractToken extracts auth token from request.
func extractToken(r *http.Request) string {
	// Check Authorization header first
	auth := r.Header.Get("Authorization")
	if strings.HasPrefix(auth, "Bearer ") {
		return strings.TrimPrefix(auth, "Bearer ")
	}

	// Check query param
	return r.URL.Query().Get("token")
}

// getContentType returns MIME type based on file extension.
func getContentType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))

	switch ext {
	// Video
	case ".mp4":
		return "video/mp4"
	case ".mkv":
		return "video/x-matroska"
	case ".avi":
		return "video/x-msvideo"
	case ".webm":
		return "video/webm"
	case ".mov":
		return "video/quicktime"
	case ".m4v":
		return "video/x-m4v"
	case ".ts":
		return "video/mp2t"
	case ".m2ts":
		return "video/mp2t"

	// Audio
	case ".mp3":
		return "audio/mpeg"
	case ".flac":
		return "audio/flac"
	case ".m4a":
		return "audio/mp4"
	case ".aac":
		return "audio/aac"
	case ".ogg":
		return "audio/ogg"
	case ".opus":
		return "audio/opus"
	case ".wav":
		return "audio/wav"

	// Subtitles
	case ".srt":
		return "text/plain"
	case ".ass", ".ssa":
		return "text/plain"
	case ".vtt":
		return "text/vtt"

	default:
		return "application/octet-stream"
	}
}

// --- Internal API Endpoints for Blackbeard ---

// InternalStreamHandler handles the internal streaming endpoint for Blackbeard.
// Route: GET /internal/stream/{mediaID}
func (s *MediaFileServer) InternalStreamHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.ServeFile(w, r)
	}
}

// InternalProbeHandler returns media file metadata for Blackbeard to analyze.
// Route: GET /internal/probe/{mediaID}
func (s *MediaFileServer) InternalProbeHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := extractToken(r)
		if token == "" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		st, ok := s.ValidateToken(token)
		if !ok {
			http.Error(w, "invalid or expired token", http.StatusUnauthorized)
			return
		}

		// Return file info for Blackbeard
		stat, err := os.Stat(st.FilePath)
		if err != nil {
			http.Error(w, "file not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"path":"%s","size":%d,"media_id":"%s"}`,
			st.FilePath, stat.Size(), st.MediaID.String())
	}
}
