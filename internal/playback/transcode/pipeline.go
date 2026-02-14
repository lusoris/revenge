package transcode

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/maypok86/otter/v2"

	"github.com/lusoris/revenge/internal/infra/cache"
	"github.com/lusoris/revenge/internal/infra/observability"
)

// PipelineManager manages running transcode jobs per session+profile.
// Jobs run in-process using astiav (libavcodec/libavformat/libavfilter)
// instead of spawning FFmpeg child processes.
type PipelineManager struct {
	segmentDuration int
	jobs            *cache.L1Cache[string, *TranscodeJob]
	logger          *slog.Logger
}

// NewPipelineManager creates a new pipeline manager.
func NewPipelineManager(segmentDuration int, logger *slog.Logger) (*PipelineManager, error) {
	pm := &PipelineManager{
		segmentDuration: segmentDuration,
		logger:          logger,
	}

	// ttl=0: no automatic expiration — jobs are managed manually via StopProcess/StopAllForSession.
	// OnDeletion: stop transcode jobs evicted by cache size pressure to prevent orphaned goroutines.
	jobCache, err := cache.NewL1Cache[string, *TranscodeJob](1000, 0,
		cache.WithOnDeletion(func(e otter.DeletionEvent[string, *TranscodeJob]) {
			if !e.WasEvicted() {
				return // manual invalidation already handled in StopProcess
			}
			job := e.Value
			if job != nil {
				logger.Warn("stopping evicted transcode job",
					slog.String("key", e.Key),
					slog.String("session_id", job.SessionID),
					slog.String("profile", job.Profile),
				)
				job.Stop()
			}
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create job cache: %w", err)
	}
	pm.jobs = jobCache

	return pm, nil
}

func processKey(sessionID uuid.UUID, profile string) string {
	return sessionID.String() + ":" + profile
}

// StartVideoSegmenting launches an in-process transcode job to output video-only HLS segments.
// Audio is excluded — each audio track gets its own rendition via StartAudioRendition.
func (pm *PipelineManager) StartVideoSegmenting(ctx context.Context, sessionID uuid.UUID, filePath, segmentDir string, pd ProfileDecision, seekSeconds int) (*TranscodeJob, error) {
	key := processKey(sessionID, pd.Name)

	profileDir := filepath.Join(segmentDir, pd.Name)
	if err := os.MkdirAll(profileDir, 0o750); err != nil {
		return nil, fmt.Errorf("failed to create segment dir %s: %w", profileDir, err)
	}

	job := NewTranscodeJob(TranscodeJobConfig{
		InputFile:        filePath,
		OutputDir:        profileDir,
		SessionID:        sessionID.String(),
		Profile:          pd.Name,
		VideoCodec:       pd.VideoCodec,
		AudioCodec:       "", // no audio — separate renditions
		Width:            pd.Width,
		Height:           pd.Height,
		VideoBitrate:     pd.VideoBitrate,
		SegmentDuration:  pm.segmentDuration,
		VideoStreamIndex: 0,  // first video stream
		AudioStreamIndex: -1, // disable audio
		SeekSeconds:      seekSeconds,
	})

	return pm.startJob(ctx, job, key, sessionID, pd.Name, pd.VideoCodec, pd.NeedsTranscode)
}

// StartAudioRendition launches an in-process transcode job to output audio-only HLS segments
// for a single audio track. Each track is a separate rendition — HLS.js downloads
// only the selected track's segments, preserving original quality and saving bandwidth.
func (pm *PipelineManager) StartAudioRendition(ctx context.Context, sessionID uuid.UUID, filePath, segmentDir string, trackIndex int, codec string, bitrate, seekSeconds int) (*TranscodeJob, error) {
	renditionName := fmt.Sprintf("audio/%d", trackIndex)
	key := processKey(sessionID, renditionName)

	audioDir := filepath.Join(segmentDir, "audio", fmt.Sprintf("%d", trackIndex))
	if err := os.MkdirAll(audioDir, 0o750); err != nil {
		return nil, fmt.Errorf("failed to create audio rendition dir %s: %w", audioDir, err)
	}

	job := NewTranscodeJob(TranscodeJobConfig{
		InputFile:        filePath,
		OutputDir:        audioDir,
		SessionID:        sessionID.String(),
		Profile:          renditionName,
		VideoCodec:       "",     // no video
		AudioCodec:       codec,
		AudioBitrate:     bitrate,
		SegmentDuration:  pm.segmentDuration,
		VideoStreamIndex: -1, // disable video
		AudioStreamIndex: trackIndex,
		SeekSeconds:      seekSeconds,
	})

	return pm.startJob(ctx, job, key, sessionID, renditionName, codec, codec != "copy")
}

func (pm *PipelineManager) startJob(_ context.Context, job *TranscodeJob, key string, sessionID uuid.UUID, name, codec string, isTranscode bool) (*TranscodeJob, error) {
	pm.jobs.Set(key, job)

	// Record transcoding start metric
	if isTranscode {
		observability.RecordTranscodingStart()
	}

	pm.logger.Info("transcode job started",
		slog.String("session_id", sessionID.String()),
		slog.String("profile", name),
		slog.Bool("transcode", isTranscode),
		slog.String("input", job.InputFile),
		slog.String("output", job.OutputFile),
	)

	startTime := time.Now()

	// Use a detached context — transcode jobs must outlive the HTTP request.
	// Cancellation is handled by job.Stop() (called from StopProcess/StopAllForSession).
	go func() {
		defer close(job.Done)
		job.Err = job.Run(context.Background())

		// Record transcoding end metric
		if job.IsTranscode {
			duration := time.Since(startTime).Seconds()
			resolution := resolveResolutionLabel(job.Profile)
			observability.RecordTranscodingEnd(codec, resolution, duration)
		}

		if job.Err != nil {
			// Context cancellation is expected during cleanup (Stop() was called)
			if strings.Contains(job.Err.Error(), "context canceled") {
				pm.logger.Debug("transcode job cancelled (expected cleanup)",
					slog.String("session_id", sessionID.String()),
					slog.String("profile", name),
				)
			} else {
				pm.logger.Error("transcode job failed",
					slog.String("session_id", sessionID.String()),
					slog.String("profile", name),
					slog.String("error", job.Err.Error()),
				)
			}
		} else {
			pm.logger.Info("transcode job completed",
				slog.String("session_id", sessionID.String()),
				slog.String("profile", name),
			)
		}
	}()

	return job, nil
}

// StopProcess stops a transcode job for a session+profile.
func (pm *PipelineManager) StopProcess(sessionID uuid.UUID, profile string) error {
	key := processKey(sessionID, profile)
	job, ok := pm.jobs.Get(key)
	if !ok {
		return nil
	}

	pm.jobs.Delete(key)

	if job != nil {
		job.Stop()
		<-job.Done
	}

	return nil
}

// StopAllForSession stops all transcode jobs for a session (video + audio renditions).
func (pm *PipelineManager) StopAllForSession(sessionID uuid.UUID) {
	profiles := []string{"original", "1080p", "720p", "480p"}
	for _, p := range profiles {
		_ = pm.StopProcess(sessionID, p)
	}
	// Stop audio renditions (up to 16 tracks)
	for i := range 16 {
		_ = pm.StopProcess(sessionID, fmt.Sprintf("audio/%d", i))
	}
}

// GetProcess returns the transcode job for a session+profile, if running.
func (pm *PipelineManager) GetProcess(sessionID uuid.UUID, profile string) (*TranscodeJob, bool) {
	return pm.jobs.Get(processKey(sessionID, profile))
}

// Close shuts down the pipeline manager and stops all jobs.
func (pm *PipelineManager) Close() {
	pm.jobs.Close()
}

// resolveResolutionLabel extracts a resolution label from a profile name.
func resolveResolutionLabel(profile string) string {
	switch {
	case strings.Contains(profile, "1080"):
		return "1080p"
	case strings.Contains(profile, "720"):
		return "720p"
	case strings.Contains(profile, "480"):
		return "480p"
	case profile == "original":
		return "original"
	default:
		return "unknown"
	}
}
