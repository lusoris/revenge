package movie

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
)

// Service provides adult movie business logic.
type Service struct {
	repo   Repository
	logger *slog.Logger
}

// NewService creates a new adult movie service.
func NewService(repo Repository, logger *slog.Logger) *Service {
	return &Service{
		repo:   repo,
		logger: logger.With("service", "adult_movie"),
	}
}

// GetByID returns an adult movie by ID.
func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*Movie, error) {
	return s.repo.GetByID(ctx, id)
}

// List lists adult movies.
func (s *Service) List(ctx context.Context, params ListParams) ([]*Movie, error) {
	return s.repo.List(ctx, params)
}

// ListByLibrary lists adult movies for a library.
func (s *Service) ListByLibrary(ctx context.Context, libraryID uuid.UUID, params ListParams) ([]*Movie, error) {
	return s.repo.ListByLibrary(ctx, libraryID, params)
}

// Create stores a new adult movie.
func (s *Service) Create(ctx context.Context, movie *Movie) error {
	return s.repo.Create(ctx, movie)
}

// Update updates an adult movie.
func (s *Service) Update(ctx context.Context, movie *Movie) error {
	return s.repo.Update(ctx, movie)
}

// Delete removes an adult movie.
func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}
