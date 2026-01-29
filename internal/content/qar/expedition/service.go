// Package expedition provides adult movie domain models (QAR obfuscation: movies → expeditions).
package expedition

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
)

// Service provides expedition (adult movie) business logic.
type Service struct {
	repo   Repository
	logger *slog.Logger
}

// NewService creates a new expedition service.
func NewService(repo Repository, logger *slog.Logger) *Service {
	return &Service{
		repo:   repo,
		logger: logger.With(slog.String("service", "qar.expedition")),
	}
}

// GetByID retrieves an expedition by ID.
func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*Expedition, error) {
	expedition, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get expedition: %w", err)
	}
	return expedition, nil
}

// List retrieves expeditions with pagination.
func (s *Service) List(ctx context.Context, limit, offset int) ([]Expedition, error) {
	expeditions, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list expeditions: %w", err)
	}
	return expeditions, nil
}

// ListByFleet retrieves expeditions for a specific fleet.
func (s *Service) ListByFleet(ctx context.Context, fleetID uuid.UUID, limit, offset int) ([]Expedition, error) {
	expeditions, err := s.repo.ListByFleet(ctx, fleetID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list expeditions by fleet: %w", err)
	}
	return expeditions, nil
}

// Create creates a new expedition.
func (s *Service) Create(ctx context.Context, expedition *Expedition) error {
	if err := s.repo.Create(ctx, expedition); err != nil {
		return fmt.Errorf("create expedition: %w", err)
	}
	s.logger.Info("expedition created",
		slog.String("id", expedition.ID.String()),
		slog.String("title", expedition.Title),
	)
	return nil
}

// Update updates an existing expedition.
func (s *Service) Update(ctx context.Context, expedition *Expedition) error {
	if err := s.repo.Update(ctx, expedition); err != nil {
		return fmt.Errorf("update expedition: %w", err)
	}
	s.logger.Info("expedition updated", slog.String("id", expedition.ID.String()))
	return nil
}

// Delete removes an expedition.
func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete expedition: %w", err)
	}
	s.logger.Info("expedition deleted", slog.String("id", id.String()))
	return nil
}

// Search searches expeditions by title.
func (s *Service) Search(ctx context.Context, query string, limit, offset int) ([]Expedition, error) {
	expeditions, err := s.repo.Search(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("search expeditions: %w", err)
	}
	return expeditions, nil
}

// MatchByCoordinates finds an expedition by perceptual hash (fingerprint).
func (s *Service) MatchByCoordinates(ctx context.Context, coordinates string) (*Expedition, error) {
	expedition, err := s.repo.GetByCoordinates(ctx, coordinates)
	if err != nil {
		return nil, fmt.Errorf("match by coordinates: %w", err)
	}
	return expedition, nil
}
