// Package playback provides playback session management and progress tracking.
package playback

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/google/uuid"
)

var (
	// ErrSessionNotFound indicates the playback session was not found.
	ErrSessionNotFound = errors.New("playback session not found")
	// ErrMediaNotFound indicates the media item was not found.
	ErrMediaNotFound = errors.New("media not found")
)

// Service provides playback session management.
type Service struct {
	logger *slog.Logger

	// In-memory session store (for MVP - later back with Redis/DB)
	sessions   map[uuid.UUID]*PlaybackSession
	sessionsMu sync.RWMutex

	// User to session mapping (user can have one session per device)
	userSessions   map[uuid.UUID]map[string]uuid.UUID // userID -> deviceID -> sessionID
	userSessionsMu sync.RWMutex

	// Configuration
	autoMarkWatchedPercent float64       // Mark as watched when this percent is reached
	continueWatchingDays   int           // How many days to keep in continue watching
	sessionTimeout         time.Duration // Timeout for inactive sessions
}

// NewService creates a new playback service.
func NewService(logger *slog.Logger) *Service {
	return &Service{
		logger:                 logger.With(slog.String("service", "playback")),
		sessions:               make(map[uuid.UUID]*PlaybackSession),
		userSessions:           make(map[uuid.UUID]map[string]uuid.UUID),
		autoMarkWatchedPercent: 90.0,
		continueWatchingDays:   30,
		sessionTimeout:         30 * time.Minute,
	}
}

// StartPlayback starts a new playback session or resumes an existing one.
func (s *Service) StartPlayback(ctx context.Context, params StartPlaybackParams) (*PlaybackSession, error) {
	s.sessionsMu.Lock()
	defer s.sessionsMu.Unlock()

	// Check if user already has a session for this device
	deviceID := "default"
	if params.DeviceID != nil {
		deviceID = *params.DeviceID
	}

	s.userSessionsMu.Lock()
	if userDevices, ok := s.userSessions[params.UserID]; ok {
		if existingSessionID, ok := userDevices[deviceID]; ok {
			// Return existing session
			if session, ok := s.sessions[existingSessionID]; ok {
				s.userSessionsMu.Unlock()
				session.LastActivityAt = time.Now()
				return session, nil
			}
		}
	}
	s.userSessionsMu.Unlock()

	// Create new session
	now := time.Now()
	session := &PlaybackSession{
		ID:             uuid.New(),
		UserID:         params.UserID,
		MediaID:        params.MediaID,
		MediaType:      params.MediaType,
		PositionTicks:  params.PositionTicks,
		RuntimeTicks:   params.RuntimeTicks,
		IsPaused:       false,
		StartedAt:      now,
		LastActivityAt: now,
		DeviceID:       params.DeviceID,
		DeviceName:     params.DeviceName,
		ClientName:     params.ClientName,
	}

	if params.RuntimeTicks > 0 {
		session.PlayedPercent = float64(params.PositionTicks) / float64(params.RuntimeTicks) * 100
	}

	s.sessions[session.ID] = session

	// Track user -> device -> session mapping
	s.userSessionsMu.Lock()
	if s.userSessions[params.UserID] == nil {
		s.userSessions[params.UserID] = make(map[string]uuid.UUID)
	}
	s.userSessions[params.UserID][deviceID] = session.ID
	s.userSessionsMu.Unlock()

	s.logger.Info("Playback started",
		slog.String("session_id", session.ID.String()),
		slog.String("user_id", params.UserID.String()),
		slog.String("media_id", params.MediaID.String()),
		slog.String("media_type", string(params.MediaType)),
	)

	return session, nil
}

// UpdateProgress updates the playback progress for a session.
func (s *Service) UpdateProgress(ctx context.Context, params UpdateProgressParams) error {
	s.sessionsMu.Lock()
	defer s.sessionsMu.Unlock()

	session, ok := s.sessions[params.SessionID]
	if !ok {
		return ErrSessionNotFound
	}

	session.PositionTicks = params.PositionTicks
	session.IsPaused = params.IsPaused
	session.LastActivityAt = time.Now()

	if session.RuntimeTicks > 0 {
		session.PlayedPercent = float64(params.PositionTicks) / float64(session.RuntimeTicks) * 100
	}

	return nil
}

// StopPlayback stops an active playback session.
func (s *Service) StopPlayback(ctx context.Context, sessionID uuid.UUID) error {
	s.sessionsMu.Lock()
	defer s.sessionsMu.Unlock()

	session, ok := s.sessions[sessionID]
	if !ok {
		return ErrSessionNotFound
	}

	// Remove from user sessions mapping
	s.userSessionsMu.Lock()
	deviceID := "default"
	if session.DeviceID != nil {
		deviceID = *session.DeviceID
	}
	if userDevices, ok := s.userSessions[session.UserID]; ok {
		delete(userDevices, deviceID)
		if len(userDevices) == 0 {
			delete(s.userSessions, session.UserID)
		}
	}
	s.userSessionsMu.Unlock()

	delete(s.sessions, sessionID)

	s.logger.Info("Playback stopped",
		slog.String("session_id", sessionID.String()),
		slog.Float64("played_percent", session.PlayedPercent),
	)

	// TODO: Persist final position to database for continue watching
	// This would call into movie/episode service to update watch history

	return nil
}

// GetSession retrieves a playback session by ID.
func (s *Service) GetSession(ctx context.Context, sessionID uuid.UUID) (*PlaybackSession, error) {
	s.sessionsMu.RLock()
	defer s.sessionsMu.RUnlock()

	session, ok := s.sessions[sessionID]
	if !ok {
		return nil, ErrSessionNotFound
	}

	return session, nil
}

// GetActiveSession retrieves the active session for a user and media item.
func (s *Service) GetActiveSession(ctx context.Context, userID, mediaID uuid.UUID) (*PlaybackSession, error) {
	s.sessionsMu.RLock()
	defer s.sessionsMu.RUnlock()

	for _, session := range s.sessions {
		if session.UserID == userID && session.MediaID == mediaID {
			return session, nil
		}
	}

	return nil, ErrSessionNotFound
}

// GetUserSessions retrieves all active sessions for a user.
func (s *Service) GetUserSessions(ctx context.Context, userID uuid.UUID) ([]*PlaybackSession, error) {
	s.sessionsMu.RLock()
	defer s.sessionsMu.RUnlock()

	var result []*PlaybackSession
	for _, session := range s.sessions {
		if session.UserID == userID {
			result = append(result, session)
		}
	}

	return result, nil
}

// BuildUpNextQueue builds the up-next queue for the current media.
// For TV episodes: next episode in series
// For movies: similar movies or next in collection
// For adult: similar content
func (s *Service) BuildUpNextQueue(ctx context.Context, userID, currentMediaID uuid.UUID, mediaType MediaType) (*UpNextQueue, error) {
	// TODO: Implement based on media type
	// For now, return empty queue
	return &UpNextQueue{
		Items: []UpNextItem{},
	}, nil
}

// CleanupInactiveSessions removes sessions that have been inactive for too long.
func (s *Service) CleanupInactiveSessions(ctx context.Context) error {
	s.sessionsMu.Lock()
	defer s.sessionsMu.Unlock()

	cutoff := time.Now().Add(-s.sessionTimeout)
	var toRemove []uuid.UUID

	for id, session := range s.sessions {
		if session.LastActivityAt.Before(cutoff) {
			toRemove = append(toRemove, id)
		}
	}

	for _, id := range toRemove {
		session := s.sessions[id]

		// Remove from user sessions mapping
		s.userSessionsMu.Lock()
		deviceID := "default"
		if session.DeviceID != nil {
			deviceID = *session.DeviceID
		}
		if userDevices, ok := s.userSessions[session.UserID]; ok {
			delete(userDevices, deviceID)
			if len(userDevices) == 0 {
				delete(s.userSessions, session.UserID)
			}
		}
		s.userSessionsMu.Unlock()

		delete(s.sessions, id)
	}

	if len(toRemove) > 0 {
		s.logger.Info("Cleaned up inactive sessions",
			slog.Int("count", len(toRemove)),
		)
	}

	return nil
}

// SetAutoMarkWatchedPercent sets the percentage at which content is marked as watched.
func (s *Service) SetAutoMarkWatchedPercent(percent float64) {
	s.autoMarkWatchedPercent = percent
}

// SetSessionTimeout sets the timeout for inactive sessions.
func (s *Service) SetSessionTimeout(timeout time.Duration) {
	s.sessionTimeout = timeout
}

// IsWatched checks if the content should be marked as watched.
func (s *Service) IsWatched(session *PlaybackSession) bool {
	return session.PlayedPercent >= s.autoMarkWatchedPercent
}

// String implements fmt.Stringer for MediaType.
func (m MediaType) String() string {
	return string(m)
}

// ValidateMediaType checks if a media type string is valid.
func ValidateMediaType(t string) (MediaType, error) {
	switch MediaType(t) {
	case MediaTypeMovie, MediaTypeEpisode, MediaTypeAdult:
		return MediaType(t), nil
	default:
		return "", fmt.Errorf("invalid media type: %s", t)
	}
}
