package api

import (
	"context"

	"log/slog"

	"github.com/lusoris/revenge/internal/api/ogen"
	"github.com/lusoris/revenge/internal/playback"
)

// ============================================================================
// Playback / HLS Streaming Endpoints
// ============================================================================

// StartPlaybackSession creates a new HLS playback session.
// POST /api/v1/playback/sessions
func (h *Handler) StartPlaybackSession(ctx context.Context, req *ogen.StartPlaybackRequest) (ogen.StartPlaybackSessionRes, error) {
	if h.playbackService == nil {
		return &ogen.StartPlaybackSessionBadRequest{
			Code:    400,
			Message: "Playback is not enabled",
		}, nil
	}

	userID, ok := h.getUserID(ctx)
	if !ok {
		return &ogen.StartPlaybackSessionUnauthorized{
			Code:    401,
			Message: "Authentication required",
		}, nil
	}

	// Convert ogen request to internal request
	pbReq := &playback.StartPlaybackRequest{
		MediaType: playback.MediaType(req.MediaType),
		MediaID:   req.MediaID,
	}
	if req.FileID.Set {
		id := req.FileID.Value
		pbReq.FileID = &id
	}
	if req.AudioTrack.Set {
		pbReq.AudioTrack = req.AudioTrack.Value
	}
	if req.SubtitleTrack.Set {
		v := req.SubtitleTrack.Value
		pbReq.SubtitleTrack = &v
	}
	if req.StartPosition.Set {
		pbReq.StartPosition = req.StartPosition.Value
	}

	sess, err := h.playbackService.StartSession(ctx, userID, pbReq)
	if err != nil {
		h.logger.Error("failed to start playback session",
			slog.Any("error",err),
			slog.String("user_id", userID.String()),
		)
		return &ogen.StartPlaybackSessionNotFound{
			Code:    404,
			Message: err.Error(),
		}, nil
	}

	return sessionToOgen(sess), nil
}

// GetPlaybackSession returns metadata for an active playback session.
// GET /api/v1/playback/sessions/{sessionId}
func (h *Handler) GetPlaybackSession(ctx context.Context, params ogen.GetPlaybackSessionParams) (ogen.GetPlaybackSessionRes, error) {
	if h.playbackService == nil {
		return &ogen.GetPlaybackSessionNotFound{
			Code:    404,
			Message: "Playback is not enabled",
		}, nil
	}

	_, ok := h.getUserID(ctx)
	if !ok {
		return &ogen.GetPlaybackSessionUnauthorized{
			Code:    401,
			Message: "Authentication required",
		}, nil
	}

	sess, ok := h.playbackService.GetSession(params.SessionId)
	if !ok {
		return &ogen.GetPlaybackSessionNotFound{
			Code:    404,
			Message: "Session not found or expired",
		}, nil
	}

	return sessionToOgen(sess), nil
}

// StopPlaybackSession terminates a playback session and cleans up.
// DELETE /api/v1/playback/sessions/{sessionId}
func (h *Handler) StopPlaybackSession(ctx context.Context, params ogen.StopPlaybackSessionParams) (ogen.StopPlaybackSessionRes, error) {
	if h.playbackService == nil {
		return &ogen.StopPlaybackSessionNotFound{
			Code:    404,
			Message: "Playback is not enabled",
		}, nil
	}

	_, ok := h.getUserID(ctx)
	if !ok {
		return &ogen.StopPlaybackSessionUnauthorized{
			Code:    401,
			Message: "Authentication required",
		}, nil
	}

	if err := h.playbackService.StopSession(params.SessionId); err != nil {
		return &ogen.StopPlaybackSessionNotFound{
			Code:    404,
			Message: "Session not found",
		}, nil
	}

	return &ogen.StopPlaybackSessionNoContent{}, nil
}

// sessionToOgen converts an internal Session to the ogen PlaybackSession response.
func sessionToOgen(sess *playback.Session) *ogen.PlaybackSession {
	resp := playback.SessionToResponse(sess)

	profiles := make([]ogen.PlaybackProfile, len(resp.Profiles))
	for i, p := range resp.Profiles {
		profiles[i] = ogen.PlaybackProfile{
			Name:       p.Name,
			Width:      p.Width,
			Height:     p.Height,
			Bitrate:    p.Bitrate,
			IsOriginal: p.IsOriginal,
		}
	}

	audioTracks := make([]ogen.PlaybackAudioTrack, len(resp.AudioTracks))
	for i, a := range resp.AudioTracks {
		at := ogen.PlaybackAudioTrack{
			Index:     a.Index,
			Channels:  a.Channels,
			Codec:     a.Codec,
			IsDefault: a.IsDefault,
		}
		if a.Language != "" {
			at.Language = ogen.NewOptString(a.Language)
		}
		if a.Title != "" {
			at.Title = ogen.NewOptString(a.Title)
		}
		if a.Layout != "" {
			at.Layout = ogen.NewOptString(a.Layout)
		}
		audioTracks[i] = at
	}

	subtitleTracks := make([]ogen.PlaybackSubtitleTrack, len(resp.SubtitleTracks))
	for i, s := range resp.SubtitleTracks {
		st := ogen.PlaybackSubtitleTrack{
			Index:    s.Index,
			Codec:    s.Codec,
			URL:      s.URL,
			IsForced: s.IsForced,
		}
		if s.Language != "" {
			st.Language = ogen.NewOptString(s.Language)
		}
		if s.Title != "" {
			st.Title = ogen.NewOptString(s.Title)
		}
		subtitleTracks[i] = st
	}

	return &ogen.PlaybackSession{
		SessionID:         resp.SessionID,
		MasterPlaylistURL: resp.MasterPlaylistURL,
		DurationSeconds:   resp.DurationSeconds,
		Profiles:          profiles,
		AudioTracks:       audioTracks,
		SubtitleTracks:    subtitleTracks,
		CreatedAt:         resp.CreatedAt,
		ExpiresAt:         resp.ExpiresAt,
	}
}
