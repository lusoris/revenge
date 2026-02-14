package playback

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/lusoris/revenge/internal/config"
	"github.com/lusoris/revenge/internal/content/movie"
	"github.com/lusoris/revenge/internal/content/tvshow"
	"github.com/lusoris/revenge/internal/infra/cache"
	"github.com/lusoris/revenge/internal/infra/observability"
	"github.com/lusoris/revenge/internal/playback/subtitle"
	"github.com/lusoris/revenge/internal/playback/transcode"
)

// Service manages playback sessions and streaming pipelines.
type Service struct {
	cfg        *config.Config
	sessions   *SessionManager
	pipeline   *transcode.PipelineManager
	prober     movie.Prober
	movieSvc   movie.Service
	tvSvc      tvshow.Service
	profiles   []transcode.QualityProfile
	probeCache *cache.L1Cache[uuid.UUID, *movie.MediaInfo]
	logger     *slog.Logger
}

// NewService creates a new playback service.
func NewService(
	cfg *config.Config,
	sessions *SessionManager,
	pipeline *transcode.PipelineManager,
	prober movie.Prober,
	movieSvc movie.Service,
	tvSvc tvshow.Service,
	logger *slog.Logger,
) (*Service, error) {
	profiles := transcode.GetEnabledProfiles(cfg.Playback.Transcode.Profiles)

	probeCache, err := cache.NewL1Cache[uuid.UUID, *movie.MediaInfo](500, 1*time.Hour)
	if err != nil {
		return nil, fmt.Errorf("failed to create probe cache: %w", err)
	}

	return &Service{
		cfg:        cfg,
		sessions:   sessions,
		pipeline:   pipeline,
		prober:     prober,
		movieSvc:   movieSvc,
		tvSvc:      tvSvc,
		profiles:   profiles,
		probeCache: probeCache,
		logger:     logger,
	}, nil
}

// StartSession creates a new playback session with FFmpeg pipeline.
func (s *Service) StartSession(ctx context.Context, userID uuid.UUID, req *StartPlaybackRequest) (*Session, error) {
	// 1. Resolve the file path
	filePath, fileID, err := s.resolveFilePath(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve file: %w", err)
	}

	// 2. Probe media file (cached)
	info, err := s.probeFile(fileID, filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to probe media: %w", err)
	}

	// 3. Analyze transcode decision
	decision := transcode.AnalyzeMedia(info, s.profiles)

	// 4. Create session
	sessionID := uuid.Must(uuid.NewV7())
	segmentDir := filepath.Join(s.cfg.Playback.SegmentDir, sessionID.String())

	sess := &Session{
		ID:                sessionID,
		UserID:            userID,
		MediaType:         req.MediaType,
		MediaID:           req.MediaID,
		FileID:            fileID,
		FilePath:          filePath,
		SegmentDir:        segmentDir,
		TranscodeDecision: decision,
		ActiveProfiles:    profileNames(decision.Profiles),
		AudioTrack:        req.AudioTrack,
		SubtitleTrack:     req.SubtitleTrack,
		StartPosition:     req.StartPosition,
		DurationSeconds:   info.DurationSeconds,
		AudioTracks:       AudioTracksFromMediaInfo(info),
		SubtitleTracks:    SubtitleTracksFromMediaInfo(info, sessionID),
	}

	if err := s.sessions.Create(sess); err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	// 5. Create segment directory
	if err := os.MkdirAll(segmentDir, 0o750); err != nil {
		s.sessions.Delete(sessionID)
		return nil, fmt.Errorf("failed to create segment dir: %w", err)
	}

	// 6. Start video-only FFmpeg pipelines per quality profile (async, non-blocking)
	for _, pd := range decision.Profiles {
		if _, err := s.pipeline.StartVideoSegmenting(ctx, sessionID, filePath, segmentDir, pd, req.StartPosition); err != nil {
			s.logger.Error("failed to start video segmenting",
				slog.String("session_id", sessionID.String()),
				slog.String("profile", pd.Name),
				slog.String("error", err.Error()),
			)
		}
	}

	// 7. Start separate audio renditions — one per audio track.
	// Each track is segmented independently so HLS.js only downloads the active track.
	// HLS-compatible codecs (AAC, AC-3, E-AC-3) are copied, others transcoded to AAC.
	for _, as := range info.AudioStreams {
		codec, bitrate := audioRenditionCodec(as.Codec)
		if _, err := s.pipeline.StartAudioRendition(ctx, sessionID, filePath, segmentDir, as.Index, codec, bitrate, req.StartPosition); err != nil {
			s.logger.Error("failed to start audio rendition",
				slog.String("session_id", sessionID.String()),
				slog.Int("track_index", as.Index),
				slog.String("error", err.Error()),
			)
		}
	}

	// 8. Extract subtitles (async, non-blocking) — full VTT files, not segmented
	go s.extractSubtitles(sessionID, filePath, segmentDir, info)

	// Record playback start metrics
	quality := "direct"
	if len(decision.Profiles) > 0 {
		quality = decision.Profiles[0].Name
	}
	observability.RecordPlaybackStart(string(req.MediaType), quality)

	s.logger.Info("playback session started",
		slog.String("session_id", sessionID.String()),
		slog.String("user_id", userID.String()),
		slog.String("media_type", string(req.MediaType)),
		slog.Bool("can_remux", decision.CanRemux),
		slog.Int("profiles", len(decision.Profiles)),
		slog.Int("audio_tracks", len(info.AudioStreams)),
		slog.Int("subtitle_tracks", len(sess.SubtitleTracks)),
	)

	return sess, nil
}

// GetSession returns a session by ID.
func (s *Service) GetSession(sessionID uuid.UUID) (*Session, bool) {
	return s.sessions.Get(sessionID)
}

// HeartbeatSession keeps a playback session alive and optionally updates the
// playback position. Returns the updated session or false if the session doesn't exist.
func (s *Service) HeartbeatSession(sessionID uuid.UUID, positionSeconds *int) (*Session, bool) {
	session, ok := s.sessions.Get(sessionID)
	if !ok {
		return nil, false
	}

	now := time.Now()
	session.LastAccessedAt = now
	session.ExpiresAt = now.Add(s.sessions.timeout)

	if positionSeconds != nil {
		session.StartPosition = *positionSeconds
	}

	s.sessions.cache.Set(sessionID, session)
	return session, true
}

// StopSession terminates a playback session and cleans up resources.
func (s *Service) StopSession(sessionID uuid.UUID) error {
	sess := s.sessions.Delete(sessionID)
	if sess == nil {
		return fmt.Errorf("session %s not found", sessionID)
	}

	// Record playback end metrics
	duration := time.Since(sess.CreatedAt).Seconds()
	observability.RecordPlaybackEnd(string(sess.MediaType), duration)

	// Stop all FFmpeg processes
	s.pipeline.StopAllForSession(sessionID)

	// Clean up segment directory
	go func() {
		if err := os.RemoveAll(sess.SegmentDir); err != nil {
			s.logger.Warn("failed to clean up segment dir",
				slog.String("session_id", sessionID.String()),
				slog.String("dir", sess.SegmentDir),
				slog.String("error", err.Error()),
			)
		}
	}()

	s.logger.Info("playback session stopped",
		slog.String("session_id", sessionID.String()),
	)

	return nil
}

// resolveFilePath finds the media file path for a playback request.
func (s *Service) resolveFilePath(ctx context.Context, req *StartPlaybackRequest) (string, uuid.UUID, error) {
	switch req.MediaType {
	case MediaTypeMovie:
		files, err := s.movieSvc.GetMovieFiles(ctx, req.MediaID)
		if err != nil {
			return "", uuid.Nil, fmt.Errorf("movie files not found: %w", err)
		}
		if len(files) == 0 {
			return "", uuid.Nil, fmt.Errorf("no files available for movie %s", req.MediaID)
		}

		// If specific file requested, find it
		if req.FileID != nil {
			for _, f := range files {
				if f.ID == *req.FileID {
					return f.FilePath, f.ID, nil
				}
			}
			return "", uuid.Nil, fmt.Errorf("file %s not found for movie %s", req.FileID, req.MediaID)
		}

		// Default: first file
		return files[0].FilePath, files[0].ID, nil

	case MediaTypeEpisode:
		if s.tvSvc == nil {
			return "", uuid.Nil, fmt.Errorf("TV show service not available")
		}

		// If specific file requested
		if req.FileID != nil {
			file, err := s.tvSvc.GetEpisodeFile(ctx, *req.FileID)
			if err != nil {
				return "", uuid.Nil, fmt.Errorf("episode file not found: %w", err)
			}
			return file.FilePath, file.ID, nil
		}

		// Get files for episode
		files, err := s.tvSvc.ListEpisodeFiles(ctx, req.MediaID)
		if err != nil {
			return "", uuid.Nil, fmt.Errorf("episode files not found: %w", err)
		}
		if len(files) == 0 {
			return "", uuid.Nil, fmt.Errorf("no files available for episode %s", req.MediaID)
		}
		return files[0].FilePath, files[0].ID, nil

	default:
		return "", uuid.Nil, fmt.Errorf("unsupported media type: %s", req.MediaType)
	}
}

// probeFile probes a media file, using L1 cache for repeated lookups.
func (s *Service) probeFile(fileID uuid.UUID, filePath string) (*movie.MediaInfo, error) {
	if info, ok := s.probeCache.Get(fileID); ok {
		return info, nil
	}

	info, err := s.prober.Probe(filePath)
	if err != nil {
		return nil, err
	}

	s.probeCache.Set(fileID, info)
	return info, nil
}

// extractSubtitles extracts all text subtitle tracks to WebVTT files.
func (s *Service) extractSubtitles(sessionID uuid.UUID, filePath, segmentDir string, info *movie.MediaInfo) {
	for _, sub := range info.SubtitleStreams {
		if isBitmapSubtitle(sub.Codec) {
			continue
		}

		_, err := subtitle.ExtractToWebVTT(
			context.Background(),
			filePath,
			segmentDir,
			sub.Index,
		)
		if err != nil {
			s.logger.Warn("failed to extract subtitle",
				slog.String("session_id", sessionID.String()),
				slog.Int("track_index", sub.Index),
				slog.String("error", err.Error()),
			)
		}
	}
}

// SessionToResponse converts a Session to a PlaybackSessionResponse.
func SessionToResponse(sess *Session) *PlaybackSessionResponse {
	profiles := make([]ProfileInfo, 0, len(sess.TranscodeDecision.Profiles))
	for _, pd := range sess.TranscodeDecision.Profiles {
		profiles = append(profiles, ProfileInfo{
			Name:       pd.Name,
			Width:      pd.Width,
			Height:     pd.Height,
			Bitrate:    pd.VideoBitrate,
			IsOriginal: pd.VideoCodec == "copy" && pd.AudioCodec == "copy",
		})
	}

	return &PlaybackSessionResponse{
		SessionID:         sess.ID,
		MasterPlaylistURL: "/api/v1/playback/stream/" + sess.ID.String() + "/master.m3u8",
		DurationSeconds:   sess.DurationSeconds,
		Profiles:          profiles,
		AudioTracks:       sess.AudioTracks,
		SubtitleTracks:    sess.SubtitleTracks,
		CreatedAt:         sess.CreatedAt,
		ExpiresAt:         sess.ExpiresAt,
	}
}

// audioRenditionCodec determines the output codec and bitrate for an audio rendition.
// HLS-compatible codecs are copied at original quality, others transcoded to AAC.
func audioRenditionCodec(sourceCodec string) (codec string, bitrate int) {
	switch sourceCodec {
	case "aac", "mp3", "ac3", "eac3":
		return "copy", 0
	default:
		// DTS, TrueHD, FLAC, etc. → transcode to AAC at good quality
		return "aac", 256
	}
}

func profileNames(profiles []transcode.ProfileDecision) []string {
	names := make([]string, len(profiles))
	for i, p := range profiles {
		names[i] = p.Name
	}
	return names
}

// Close shuts down the service.
func (s *Service) Close() {
	s.probeCache.Close()
}
