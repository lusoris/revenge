package transcode

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/maypok86/otter/v2"

	"github.com/lusoris/revenge/internal/infra/cache"
)

// FFmpegProcess tracks a running FFmpeg process.
type FFmpegProcess struct {
	Cmd       *exec.Cmd
	SessionID uuid.UUID
	Profile   string
	Done      chan struct{}
	Err       error
}

// PipelineManager manages running FFmpeg processes per session+profile.
type PipelineManager struct {
	ffmpegPath      string
	segmentDuration int
	processes       *cache.L1Cache[string, *FFmpegProcess]
	logger          *slog.Logger
}

// NewPipelineManager creates a new pipeline manager.
func NewPipelineManager(ffmpegPath string, segmentDuration int, logger *slog.Logger) (*PipelineManager, error) {
	pm := &PipelineManager{
		ffmpegPath:      ffmpegPath,
		segmentDuration: segmentDuration,
		logger:          logger,
	}

	// ttl=0: no automatic expiration — processes are managed manually via StopProcess/StopAllForSession.
	// OnDeletion: kill FFmpeg processes evicted by cache size pressure to prevent orphaned children.
	processCache, err := cache.NewL1Cache[string, *FFmpegProcess](1000, 0,
		cache.WithOnDeletion(func(e otter.DeletionEvent[string, *FFmpegProcess]) {
			if !e.WasEvicted() {
				return // manual invalidation already handled in StopProcess
			}
			proc := e.Value
			if proc != nil && proc.Cmd.Process != nil {
				logger.Warn("killing evicted FFmpeg process",
					slog.String("key", e.Key),
					slog.String("session_id", proc.SessionID.String()),
					slog.String("profile", proc.Profile),
				)
				_ = proc.Cmd.Process.Kill()
			}
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create process cache: %w", err)
	}
	pm.processes = processCache

	return pm, nil
}

func processKey(sessionID uuid.UUID, profile string) string {
	return sessionID.String() + ":" + profile
}

// StartVideoSegmenting launches an FFmpeg process to output video-only HLS segments.
// Audio is excluded — each audio track gets its own rendition via StartAudioRendition.
func (pm *PipelineManager) StartVideoSegmenting(ctx context.Context, sessionID uuid.UUID, filePath, segmentDir string, pd ProfileDecision, seekSeconds int) (*FFmpegProcess, error) {
	key := processKey(sessionID, pd.Name)

	profileDir := filepath.Join(segmentDir, pd.Name)
	if err := os.MkdirAll(profileDir, 0o750); err != nil {
		return nil, fmt.Errorf("failed to create segment dir %s: %w", profileDir, err)
	}

	cmd := BuildVideoOnlyCommand(pm.ffmpegPath, filePath, profileDir, pd, pm.segmentDuration, seekSeconds)

	return pm.startProcess(cmd, key, sessionID, pd.Name, pd.NeedsTranscode)
}

// StartAudioRendition launches an FFmpeg process to output audio-only HLS segments
// for a single audio track. Each track is a separate rendition — HLS.js downloads
// only the selected track's segments, preserving original quality and saving bandwidth.
func (pm *PipelineManager) StartAudioRendition(ctx context.Context, sessionID uuid.UUID, filePath, segmentDir string, trackIndex int, codec string, bitrate, seekSeconds int) (*FFmpegProcess, error) {
	renditionName := fmt.Sprintf("audio/%d", trackIndex)
	key := processKey(sessionID, renditionName)

	audioDir := filepath.Join(segmentDir, "audio", fmt.Sprintf("%d", trackIndex))
	if err := os.MkdirAll(audioDir, 0o750); err != nil {
		return nil, fmt.Errorf("failed to create audio rendition dir %s: %w", audioDir, err)
	}

	cmd := BuildAudioRenditionCommand(pm.ffmpegPath, filePath, audioDir, trackIndex, codec, bitrate, pm.segmentDuration, seekSeconds)

	return pm.startProcess(cmd, key, sessionID, renditionName, codec != "copy")
}

func (pm *PipelineManager) startProcess(cmd *exec.Cmd, key string, sessionID uuid.UUID, name string, isTranscode bool) (*FFmpegProcess, error) {
	proc := &FFmpegProcess{
		Cmd:       cmd,
		SessionID: sessionID,
		Profile:   name,
		Done:      make(chan struct{}),
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start FFmpeg for %s/%s: %w", sessionID, name, err)
	}

	pm.processes.Set(key, proc)

	pm.logger.Info("FFmpeg process started",
		slog.String("session_id", sessionID.String()),
		slog.String("profile", name),
		slog.Bool("transcode", isTranscode),
		slog.String("cmd", strings.Join(cmd.Args, " ")),
	)

	go func() {
		defer close(proc.Done)
		proc.Err = cmd.Wait()

		if proc.Err != nil {
			pm.logger.Error("FFmpeg process failed",
				slog.String("session_id", sessionID.String()),
				slog.String("profile", name),
				slog.String("error", proc.Err.Error()),
			)
		} else {
			pm.logger.Info("FFmpeg process completed",
				slog.String("session_id", sessionID.String()),
				slog.String("profile", name),
			)
		}
	}()

	return proc, nil
}

// StopProcess kills an FFmpeg process for a session+profile.
func (pm *PipelineManager) StopProcess(sessionID uuid.UUID, profile string) error {
	key := processKey(sessionID, profile)
	proc, ok := pm.processes.Get(key)
	if !ok {
		return nil
	}

	pm.processes.Delete(key)

	if proc.Cmd.Process != nil {
		if err := proc.Cmd.Process.Kill(); err != nil {
			pm.logger.Warn("failed to kill FFmpeg process",
				slog.String("session_id", sessionID.String()),
				slog.String("profile", profile),
				slog.String("error", err.Error()),
			)
		}
		<-proc.Done
	}

	return nil
}

// StopAllForSession kills all FFmpeg processes for a session (video + audio renditions).
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

// GetProcess returns the FFmpeg process for a session+profile, if running.
func (pm *PipelineManager) GetProcess(sessionID uuid.UUID, profile string) (*FFmpegProcess, bool) {
	return pm.processes.Get(processKey(sessionID, profile))
}

// Close shuts down the pipeline manager and kills all processes.
func (pm *PipelineManager) Close() {
	pm.processes.Close()
}
