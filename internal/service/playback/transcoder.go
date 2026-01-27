// Package playback provides Blackbeard transcoder integration.
package playback

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/google/uuid"
)

// TranscoderConfig holds Blackbeard configuration.
type TranscoderConfig struct {
	BaseURL        string        `koanf:"base_url"`
	APIKey         string        `koanf:"api_key"`
	Timeout        time.Duration `koanf:"timeout"`
	InternalAPIURL string        `koanf:"internal_api_url"` // URL Blackbeard uses to fetch raw files
}

// TranscodeRequest describes what transcoding is needed.
type TranscodeRequest struct {
	// Source identification
	MediaID     uuid.UUID `json:"media_id"`
	StreamIndex int       `json:"stream_index"` // Video stream index in file

	// Source file info (for Blackbeard to fetch)
	SourceURL   string `json:"source_url"`   // Internal URL to raw file
	SourceToken string `json:"source_token"` // Auth token for fetching

	// Client constraints (from capabilities)
	MaxWidth    int    `json:"max_width"`
	MaxHeight   int    `json:"max_height"`
	TargetCodec string `json:"target_codec"` // "h264", "hevc", "av1"

	// Bandwidth constraints (for external clients)
	MaxBitrate    int  `json:"max_bitrate_kbps"`
	BandwidthKbps int  `json:"bandwidth_kbps"` // Measured bandwidth
	JitterKbps    int  `json:"jitter_kbps"`    // Bandwidth variance
	IsExternal    bool `json:"is_external"`

	// Audio settings
	AudioStreamIndex int    `json:"audio_stream_index"`
	AudioCodec       string `json:"audio_codec"` // "aac", "ac3", etc.
	AudioChannels    int    `json:"audio_channels"`

	// Subtitles
	SubtitleStreamIndex *int `json:"subtitle_stream_index,omitempty"`
	BurnSubtitles       bool `json:"burn_subtitles"`

	// Playback
	StartPosition time.Duration `json:"start_position"`

	// Session
	SessionID uuid.UUID `json:"session_id"`
}

// TranscodeResponse from Blackbeard.
type TranscodeResponse struct {
	TranscodeID      string `json:"transcode_id"`
	ManifestURL      string `json:"manifest_url"` // HLS master playlist
	StreamType       string `json:"stream_type"`  // "hls" or "dash"
	EstimatedBitrate int    `json:"estimated_bitrate_kbps"`
}

// TranscoderClient communicates with Blackbeard.
type TranscoderClient struct {
	config     TranscoderConfig
	httpClient *http.Client
}

// NewTranscoderClient creates a new Blackbeard client.
func NewTranscoderClient(cfg TranscoderConfig) *TranscoderClient {
	timeout := cfg.Timeout
	if timeout == 0 {
		timeout = 30 * time.Second
	}

	return &TranscoderClient{
		config: cfg,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

// StartTranscode initiates a transcoding session.
func (c *TranscoderClient) StartTranscode(ctx context.Context, req *TranscodeRequest) (*TranscodeResponse, error) {
	endpoint := fmt.Sprintf("%s/transcode/start", c.config.BaseURL)

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint,
		io.NopCloser(jsonReader(body)))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.config.APIKey)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("blackbeard error %d: %s", resp.StatusCode, string(body))
	}

	var result TranscodeResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return &result, nil
}

// StopTranscode stops a transcoding session.
func (c *TranscoderClient) StopTranscode(ctx context.Context, transcodeID string) error {
	endpoint := fmt.Sprintf("%s/transcode/%s", c.config.BaseURL, url.PathEscape(transcodeID))

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodDelete, endpoint, nil)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+c.config.APIKey)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("blackbeard error %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// ProxyStream proxies a segment/manifest from Blackbeard to the client.
// This maintains access control and allows bandwidth measurement.
func (c *TranscoderClient) ProxyStream(ctx context.Context, transcodeID, path string, w http.ResponseWriter) error {
	endpoint := fmt.Sprintf("%s/transcode/%s/%s", c.config.BaseURL,
		url.PathEscape(transcodeID), path)

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+c.config.APIKey)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	// Copy headers
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}
	w.WriteHeader(resp.StatusCode)

	// Stream body to client
	_, err = io.Copy(w, resp.Body)
	return err
}

// GetManifest fetches the HLS/DASH manifest.
func (c *TranscoderClient) GetManifest(ctx context.Context, transcodeID string) ([]byte, string, error) {
	endpoint := fmt.Sprintf("%s/transcode/%s/master.m3u8", c.config.BaseURL,
		url.PathEscape(transcodeID))

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, "", fmt.Errorf("create request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+c.config.APIKey)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, "", fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, "", fmt.Errorf("blackbeard error %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", fmt.Errorf("read body: %w", err)
	}

	contentType := resp.Header.Get("Content-Type")
	return body, contentType, nil
}

// FetchSegment fetches a single segment from Blackbeard.
// Returns the segment data, mime type, and any error.
func (c *TranscoderClient) FetchSegment(ctx context.Context, transcodeID string, segmentPath string) ([]byte, string, error) {
	endpoint := fmt.Sprintf("%s/transcode/%s/%s", c.config.BaseURL,
		url.PathEscape(transcodeID), segmentPath)

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, "", fmt.Errorf("create request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+c.config.APIKey)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, "", fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, "", fmt.Errorf("blackbeard error %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", fmt.Errorf("read body: %w", err)
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "video/mp2t" // Default for HLS segments
	}

	return body, contentType, nil
}

// StreamSegmentBuffered fetches a segment with retry and writes to output.
// This is the primary method for proxying segments with recovery support.
func (c *TranscoderClient) StreamSegmentBuffered(ctx context.Context, transcodeID, segmentPath string, w io.Writer, maxRetries int) (int64, error) {
	var lastErr error

	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			// Wait before retry with exponential backoff
			select {
			case <-ctx.Done():
				return 0, ctx.Err()
			case <-time.After(time.Duration(attempt) * 500 * time.Millisecond):
			}
		}

		data, _, err := c.FetchSegment(ctx, transcodeID, segmentPath)
		if err != nil {
			lastErr = err
			continue
		}

		n, err := w.Write(data)
		if err != nil {
			return int64(n), fmt.Errorf("write to client: %w", err)
		}

		// Flush if possible
		if flusher, ok := w.(http.Flusher); ok {
			flusher.Flush()
		}

		return int64(n), nil
	}

	return 0, fmt.Errorf("failed after %d attempts: %w", maxRetries+1, lastErr)
}

// jsonReader creates an io.Reader from JSON bytes.
func jsonReader(data []byte) io.Reader {
	return &jsonBytesReader{data: data}
}

type jsonBytesReader struct {
	data []byte
	pos  int
}

func (r *jsonBytesReader) Read(p []byte) (n int, err error) {
	if r.pos >= len(r.data) {
		return 0, io.EOF
	}
	n = copy(p, r.data[r.pos:])
	r.pos += n
	return n, nil
}
