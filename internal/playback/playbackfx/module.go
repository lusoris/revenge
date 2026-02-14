// Package playbackfx provides fx dependency injection for the playback/streaming system.
package playbackfx

import (
	"log/slog"

	"github.com/lusoris/revenge/internal/config"
	"github.com/lusoris/revenge/internal/content/movie"
	"github.com/lusoris/revenge/internal/content/tvshow"
	"github.com/lusoris/revenge/internal/playback"
	"github.com/lusoris/revenge/internal/playback/hls"
	playbackjobs "github.com/lusoris/revenge/internal/playback/jobs"
	"github.com/lusoris/revenge/internal/playback/transcode"
	"github.com/riverqueue/river"
	"go.uber.org/fx"
)

// Module provides playback/streaming dependencies.
// Components are nil when playback is disabled in config.
var Module = fx.Module("playback",
	fx.Provide(
		provideSessionManager,
		providePipelineManager,
		provideStreamHandler,
		providePlaybackService,
		provideCleanupWorker,
	),
	fx.Invoke(registerCleanupWorker),
)

func provideSessionManager(cfg *config.Config, pipeline *transcode.PipelineManager, logger *slog.Logger) (*playback.SessionManager, error) {
	if !cfg.Playback.Enabled {
		return nil, nil
	}

	// When sessions expire or are evicted by the cache, kill their FFmpeg processes.
	var cleanupFn playback.SessionCleanupFunc
	if pipeline != nil {
		cleanupFn = pipeline.StopAllForSession
	}

	return playback.NewSessionManager(
		cfg.Playback.MaxConcurrentSessions,
		cfg.Playback.SessionTimeout,
		logger.With(slog.String("component", "playback.sessions")),
		cleanupFn,
	)
}

func providePipelineManager(cfg *config.Config, logger *slog.Logger) (*transcode.PipelineManager, error) {
	if !cfg.Playback.Enabled {
		return nil, nil
	}
	return transcode.NewPipelineManager(
		cfg.Playback.SegmentDuration,
		logger.With(slog.String("component", "playback.pipeline")),
	)
}

func provideStreamHandler(cfg *config.Config, sessions *playback.SessionManager, logger *slog.Logger) (*hls.StreamHandler, error) {
	if !cfg.Playback.Enabled || sessions == nil {
		return nil, nil
	}
	return hls.NewStreamHandler(
		sessions,
		logger.With(slog.String("component", "playback.hls")),
	)
}

func providePlaybackService(
	cfg *config.Config,
	sessions *playback.SessionManager,
	pipeline *transcode.PipelineManager,
	movieSvc movie.Service,
	tvSvc tvshow.Service,
	logger *slog.Logger,
) (*playback.Service, error) {
	if !cfg.Playback.Enabled || sessions == nil || pipeline == nil {
		return nil, nil
	}
	prober := movie.NewMediaInfoProber()
	return playback.NewService(cfg, sessions, pipeline, prober, movieSvc, tvSvc, logger)
}

func provideCleanupWorker(
	cfg *config.Config,
	sessions *playback.SessionManager,
	pipeline *transcode.PipelineManager,
	logger *slog.Logger,
) *playbackjobs.CleanupWorker {
	if !cfg.Playback.Enabled || sessions == nil || pipeline == nil {
		return nil
	}
	return playbackjobs.NewCleanupWorker(
		sessions,
		pipeline,
		logger.With(slog.String("component", "playback.cleanup")),
	)
}

func registerCleanupWorker(workers *river.Workers, worker *playbackjobs.CleanupWorker) {
	if worker != nil {
		river.AddWorker(workers, worker)
	}
}
