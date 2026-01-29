// Package fleet provides adult library domain models (QAR obfuscation: libraries → fleets).
package fleet

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
)

// Service provides fleet (adult library) business logic.
type Service struct {
	repo   Repository
	logger *slog.Logger
}

// NewService creates a new fleet service.
func NewService(repo Repository, logger *slog.Logger) *Service {
	return &Service{
		repo:   repo,
		logger: logger.With(slog.String("service", "qar.fleet")),
	}
}

// GetByID retrieves a fleet by ID.
func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*Fleet, error) {
	fleet, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get fleet: %w", err)
	}
	return fleet, nil
}

// List retrieves fleets with pagination.
func (s *Service) List(ctx context.Context, limit, offset int) ([]Fleet, error) {
	fleets, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list fleets: %w", err)
	}
	return fleets, nil
}

// ListByOwner retrieves fleets owned by a specific user.
func (s *Service) ListByOwner(ctx context.Context, ownerID uuid.UUID) ([]Fleet, error) {
	fleets, err := s.repo.ListByOwner(ctx, ownerID)
	if err != nil {
		return nil, fmt.Errorf("list fleets by owner: %w", err)
	}
	return fleets, nil
}

// Create creates a new fleet.
func (s *Service) Create(ctx context.Context, fleet *Fleet) error {
	if err := s.repo.Create(ctx, fleet); err != nil {
		return fmt.Errorf("create fleet: %w", err)
	}
	s.logger.Info("fleet created",
		slog.String("id", fleet.ID.String()),
		slog.String("name", fleet.Name),
		slog.String("type", string(fleet.FleetType)),
	)
	return nil
}

// Update updates an existing fleet.
func (s *Service) Update(ctx context.Context, fleet *Fleet) error {
	if err := s.repo.Update(ctx, fleet); err != nil {
		return fmt.Errorf("update fleet: %w", err)
	}
	s.logger.Info("fleet updated", slog.String("id", fleet.ID.String()))
	return nil
}

// Delete removes a fleet.
func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete fleet: %w", err)
	}
	s.logger.Info("fleet deleted", slog.String("id", id.String()))
	return nil
}

// GetStats returns statistics for a fleet.
func (s *Service) GetStats(ctx context.Context, id uuid.UUID) (*FleetStats, error) {
	stats, err := s.repo.GetStats(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get fleet stats: %w", err)
	}
	return stats, nil
}
