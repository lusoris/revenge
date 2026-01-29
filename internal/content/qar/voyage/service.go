// Package voyage provides adult scene domain models (QAR obfuscation: scenes → voyages).
package voyage

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
)

// Service provides voyage (adult scene) business logic.
type Service struct {
	repo   Repository
	logger *slog.Logger
}

// NewService creates a new voyage service.
func NewService(repo Repository, logger *slog.Logger) *Service {
	return &Service{
		repo:   repo,
		logger: logger.With(slog.String("service", "qar.voyage")),
	}
}

// GetByID retrieves a voyage by ID.
func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*Voyage, error) {
	voyage, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get voyage: %w", err)
	}
	return voyage, nil
}

// List retrieves voyages with pagination.
func (s *Service) List(ctx context.Context, limit, offset int) ([]Voyage, error) {
	voyages, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list voyages: %w", err)
	}
	return voyages, nil
}

// ListByFleet retrieves voyages for a specific fleet.
func (s *Service) ListByFleet(ctx context.Context, fleetID uuid.UUID, limit, offset int) ([]Voyage, error) {
	voyages, err := s.repo.ListByFleet(ctx, fleetID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list voyages by fleet: %w", err)
	}
	return voyages, nil
}

// ListByPort retrieves voyages for a specific port (studio).
func (s *Service) ListByPort(ctx context.Context, portID uuid.UUID, limit, offset int) ([]Voyage, error) {
	voyages, err := s.repo.ListByPort(ctx, portID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list voyages by port: %w", err)
	}
	return voyages, nil
}

// Create creates a new voyage.
func (s *Service) Create(ctx context.Context, voyage *Voyage) error {
	if err := s.repo.Create(ctx, voyage); err != nil {
		return fmt.Errorf("create voyage: %w", err)
	}
	s.logger.Info("voyage created",
		slog.String("id", voyage.ID.String()),
		slog.String("title", voyage.Title),
	)
	return nil
}

// Update updates an existing voyage.
func (s *Service) Update(ctx context.Context, voyage *Voyage) error {
	if err := s.repo.Update(ctx, voyage); err != nil {
		return fmt.Errorf("update voyage: %w", err)
	}
	s.logger.Info("voyage updated", slog.String("id", voyage.ID.String()))
	return nil
}

// Delete removes a voyage.
func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete voyage: %w", err)
	}
	s.logger.Info("voyage deleted", slog.String("id", id.String()))
	return nil
}

// Search searches voyages by title.
func (s *Service) Search(ctx context.Context, query string, limit, offset int) ([]Voyage, error) {
	voyages, err := s.repo.Search(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("search voyages: %w", err)
	}
	return voyages, nil
}

// MatchByFingerprint finds a voyage by fingerprint (oshash or phash).
func (s *Service) MatchByFingerprint(ctx context.Context, oshash, coordinates string) (*Voyage, error) {
	// Try coordinates (phash) first - more reliable
	if coordinates != "" {
		voyage, err := s.repo.GetByCoordinates(ctx, coordinates)
		if err == nil && voyage != nil {
			return voyage, nil
		}
	}

	// Fallback to oshash
	if oshash != "" {
		voyage, err := s.repo.GetByOshash(ctx, oshash)
		if err == nil && voyage != nil {
			return voyage, nil
		}
	}

	return nil, fmt.Errorf("no matching voyage found")
}
