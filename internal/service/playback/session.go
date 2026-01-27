// Package playback provides playback session management.
package playback

import (
	"context"
	"errors"
	"log/slog"
	"sync"
	"time"

	"github.com/google/uuid"
)

// State represents the current state of playback.
type State string

const (
	StateBuffering State = "buffering"
	StatePlaying   State = "playing"
	StatePaused    State = "paused"
	StateStopped   State = "stopped"
	StateError     State = "error"
)

// Session represents an active playback session.
type Session struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	MediaID   uuid.UUID
	MediaType string // "movie", "episode", etc.

	// Client info
	Client    *ClientInfo
	Bandwidth *BandwidthMonitor

	// Transcoding
	TranscodeID      string
	IsTranscoding    bool
	TranscodeProfile string

	// Playback state
	State          State
	Position       time.Duration
	Duration       time.Duration
	StartedAt      time.Time
	LastActivityAt time.Time

	// Quality tracking
	CurrentBitrate  int
	QualitySwitches int
	BufferingEvents int
	TotalBufferTime time.Duration
}

// SessionManager manages active playback sessions.
type SessionManager struct {
	mu         sync.RWMutex
	sessions   map[uuid.UUID]*Session
	byUser     map[uuid.UUID][]uuid.UUID // userID -> sessionIDs
	transcoder *TranscoderClient
	detector   *ClientDetector
	logger     *slog.Logger
}

// NewSessionManager creates a new session manager.
func NewSessionManager(transcoder *TranscoderClient, detector *ClientDetector, logger *slog.Logger) *SessionManager {
	return &SessionManager{
		sessions:   make(map[uuid.UUID]*Session),
		byUser:     make(map[uuid.UUID][]uuid.UUID),
		transcoder: transcoder,
		detector:   detector,
		logger:     logger.With(slog.String("component", "playback")),
	}
}

// StartSession creates a new playback session.
func (m *SessionManager) StartSession(ctx context.Context, params StartSessionParams) (*Session, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	session := &Session{
		ID:             uuid.New(),
		UserID:         params.UserID,
		MediaID:        params.MediaID,
		MediaType:      params.MediaType,
		Client:         params.Client,
		Bandwidth:      NewBandwidthMonitor(params.Client.IsExternal),
		State:          StateBuffering,
		Position:       params.StartPosition,
		Duration:       params.Duration,
		StartedAt:      time.Now(),
		LastActivityAt: time.Now(),
	}

	// Determine if transcoding is needed
	if params.NeedsTranscoding {
		req := m.buildTranscodeRequest(session, params)
		resp, err := m.transcoder.StartTranscode(ctx, req)
		if err != nil {
			return nil, err
		}
		session.TranscodeID = resp.TranscodeID
		session.IsTranscoding = true
		session.TranscodeProfile = params.TranscodeProfile
		session.CurrentBitrate = resp.EstimatedBitrate
	}

	m.sessions[session.ID] = session
	m.byUser[session.UserID] = append(m.byUser[session.UserID], session.ID)

	m.logger.Info("playback session started",
		slog.String("session_id", session.ID.String()),
		slog.String("user_id", session.UserID.String()),
		slog.String("media_id", session.MediaID.String()),
		slog.Bool("transcoding", session.IsTranscoding),
		slog.Bool("external", session.Client.IsExternal),
	)

	return session, nil
}

// StartSessionParams contains parameters for starting a session.
type StartSessionParams struct {
	UserID           uuid.UUID
	MediaID          uuid.UUID
	MediaType        string
	Client           *ClientInfo
	StartPosition    time.Duration
	Duration         time.Duration
	NeedsTranscoding bool
	TranscodeProfile string
	SourceURL        string
	SourceToken      string
}

// buildTranscodeRequest creates a TranscodeRequest from session params.
func (m *SessionManager) buildTranscodeRequest(session *Session, params StartSessionParams) *TranscodeRequest {
	req := &TranscodeRequest{
		MediaID:       params.MediaID,
		SourceURL:     params.SourceURL,
		SourceToken:   params.SourceToken,
		StartPosition: params.StartPosition,
		SessionID:     session.ID,
		IsExternal:    session.Client.IsExternal,
	}

	// Apply client capabilities
	if caps := session.Client.Capabilities; caps != nil {
		req.MaxWidth = caps.MaxVideoWidth
		req.MaxHeight = caps.MaxVideoHeight
		if len(caps.SupportedCodecs) > 0 {
			req.TargetCodec = caps.SupportedCodecs[0] // Prefer first supported
		}
		req.AudioChannels = caps.MaxAudioChannels
		if len(caps.SupportedAudioCodecs) > 0 {
			req.AudioCodec = caps.SupportedAudioCodecs[0]
		}
	}

	// Apply bandwidth constraints for external clients
	if session.Client.IsExternal {
		estimate := session.Bandwidth.GetEstimate()
		if estimate.IsReliable {
			req.BandwidthKbps = estimate.AverageKbps
			req.JitterKbps = estimate.JitterKbps
			req.MaxBitrate = estimate.RecommendedKbps
		}
	}

	return req
}

// UpdateProgress updates the playback position.
func (m *SessionManager) UpdateProgress(sessionID uuid.UUID, position time.Duration, state State) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	session, ok := m.sessions[sessionID]
	if !ok {
		return ErrSessionNotFound
	}

	session.Position = position
	session.State = state
	session.LastActivityAt = time.Now()

	return nil
}

// RecordBandwidth records a bandwidth sample for a session.
func (m *SessionManager) RecordBandwidth(sessionID uuid.UUID, bytesSent int64, duration, latency time.Duration) {
	m.mu.RLock()
	session, ok := m.sessions[sessionID]
	m.mu.RUnlock()

	if !ok || session.Bandwidth == nil {
		return
	}

	session.Bandwidth.AddSample(bytesSent, duration, latency)
}

// RecordBuffering records a buffering event.
func (m *SessionManager) RecordBuffering(sessionID uuid.UUID, duration time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()

	session, ok := m.sessions[sessionID]
	if !ok {
		return
	}

	session.BufferingEvents++
	session.TotalBufferTime += duration
}

// StopSession ends a playback session.
func (m *SessionManager) StopSession(ctx context.Context, sessionID uuid.UUID) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	session, ok := m.sessions[sessionID]
	if !ok {
		return ErrSessionNotFound
	}

	// Stop transcoding if active
	if session.IsTranscoding && session.TranscodeID != "" {
		if err := m.transcoder.StopTranscode(ctx, session.TranscodeID); err != nil {
			m.logger.Warn("failed to stop transcode",
				slog.String("transcode_id", session.TranscodeID),
				slog.Any("error", err),
			)
		}
	}

	session.State = StateStopped

	// Remove from maps
	delete(m.sessions, sessionID)

	// Remove from user's sessions
	userSessions := m.byUser[session.UserID]
	for i, id := range userSessions {
		if id == sessionID {
			m.byUser[session.UserID] = append(userSessions[:i], userSessions[i+1:]...)
			break
		}
	}

	m.logger.Info("playback session stopped",
		slog.String("session_id", sessionID.String()),
		slog.Duration("watched", session.Position),
		slog.Int("buffering_events", session.BufferingEvents),
	)

	return nil
}

// GetSession returns a session by ID.
func (m *SessionManager) GetSession(sessionID uuid.UUID) (*Session, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	session, ok := m.sessions[sessionID]
	return session, ok
}

// GetUserSessions returns all sessions for a user.
func (m *SessionManager) GetUserSessions(userID uuid.UUID) []*Session {
	m.mu.RLock()
	defer m.mu.RUnlock()

	sessionIDs := m.byUser[userID]
	sessions := make([]*Session, 0, len(sessionIDs))
	for _, id := range sessionIDs {
		if session, ok := m.sessions[id]; ok {
			sessions = append(sessions, session)
		}
	}
	return sessions
}

// CleanupStale removes sessions with no activity.
func (m *SessionManager) CleanupStale(ctx context.Context, maxIdle time.Duration) int {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()
	var stale []uuid.UUID

	for id, session := range m.sessions {
		if now.Sub(session.LastActivityAt) > maxIdle {
			stale = append(stale, id)
		}
	}

	for _, id := range stale {
		session := m.sessions[id]
		if session.IsTranscoding && session.TranscodeID != "" {
			_ = m.transcoder.StopTranscode(ctx, session.TranscodeID)
		}
		delete(m.sessions, id)
	}

	return len(stale)
}

// ErrSessionNotFound is returned when a session is not found.
var ErrSessionNotFound = errors.New("playback session not found")
