package scene

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
)

// Service provides adult scene business logic.
type Service struct {
	repo   Repository
	logger *slog.Logger
}

// NewService creates a new adult scene service.
func NewService(repo Repository, logger *slog.Logger) *Service {
	return &Service{
		repo:   repo,
		logger: logger.With("service", "adult_scene"),
	}
}

// GetByID returns an adult scene by ID.
func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*Scene, error) {
	return s.repo.GetByID(ctx, id)
}

// List lists adult scenes.
func (s *Service) List(ctx context.Context, params ListParams) ([]*Scene, error) {
	return s.repo.List(ctx, params)
}

// ListByLibrary lists adult scenes for a library.
func (s *Service) ListByLibrary(ctx context.Context, libraryID uuid.UUID, params ListParams) ([]*Scene, error) {
	return s.repo.ListByLibrary(ctx, libraryID, params)
}

// Create stores a new adult scene.
func (s *Service) Create(ctx context.Context, scene *Scene) error {
	return s.repo.Create(ctx, scene)
}

// Update updates an adult scene.
func (s *Service) Update(ctx context.Context, scene *Scene) error {
	return s.repo.Update(ctx, scene)
}

// Delete removes an adult scene.
func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}
