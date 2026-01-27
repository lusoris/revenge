// Package playback provides playback session management and streaming.
package playback

import (
	"context"
	"log/slog"
	"net/http"

	"go.uber.org/fx"
)

// Module provides all playback components for fx dependency injection.
var Module = fx.Module("playback",
	fx.Provide(
		NewClientDetector,
		NewTranscoderClient,
		NewSessionManager,
		NewMediaFileServer,
		NewStreamHandlerFx,
		NewService,
	),
)

// Service is the main entry point for playback functionality.
// It orchestrates all playback components.
type Service struct {
	Sessions      *SessionManager
	StreamHandler *StreamHandler
	FileServer    *MediaFileServer
	Transcoder    *TranscoderClient
	Detector      *ClientDetector
	logger        *slog.Logger
}

// ServiceParams contains dependencies for NewService.
type ServiceParams struct {
	fx.In

	Sessions      *SessionManager
	StreamHandler *StreamHandler
	FileServer    *MediaFileServer
	Transcoder    *TranscoderClient
	Detector      *ClientDetector
	Logger        *slog.Logger
}

// NewService creates a new playback service.
func NewService(params ServiceParams) *Service {
	return &Service{
		Sessions:      params.Sessions,
		StreamHandler: params.StreamHandler,
		FileServer:    params.FileServer,
		Transcoder:    params.Transcoder,
		Detector:      params.Detector,
		logger:        params.Logger.With(slog.String("service", "playback")),
	}
}

// RegisterRoutes registers all playback-related HTTP routes.
func (s *Service) RegisterRoutes(mux *http.ServeMux) {
	s.StreamHandler.RegisterRoutes(mux)
}

// Start starts background workers (buffer cleanup, etc.).
func (s *Service) Start(ctx context.Context) {
	go s.StreamHandler.CleanupIdleBuffers(ctx)
	s.logger.Info("playback service started")
}

// LifecycleHook provides fx lifecycle integration.
func LifecycleHook(lc fx.Lifecycle, service *Service) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			service.Start(ctx)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			service.logger.Info("playback service stopped")
			return nil
		},
	})
}

// NewStreamHandlerFx creates a new stream handler with all dependencies.
// This is the fx-compatible constructor.
func NewStreamHandlerFx(
	transcoder *TranscoderClient,
	fileServer *MediaFileServer,
	sessions *SessionManager,
	logger *slog.Logger,
) *StreamHandler {
	config := DefaultStreamHandlerConfig()
	return NewStreamHandler(transcoder, fileServer, sessions, logger, config)
}
