// Package crew provides adult performer domain models (QAR obfuscation: performers → crew).
package crew

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
)

// Service provides crew (performer) business logic.
type Service struct {
	repo   Repository
	logger *slog.Logger
}

// NewService creates a new crew service.
func NewService(repo Repository, logger *slog.Logger) *Service {
	return &Service{
		repo:   repo,
		logger: logger.With(slog.String("service", "qar.crew")),
	}
}

// GetByID retrieves a crew member by ID.
func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*Crew, error) {
	crew, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get crew: %w", err)
	}
	return crew, nil
}

// List retrieves crew members with pagination.
func (s *Service) List(ctx context.Context, limit, offset int) ([]Crew, error) {
	crew, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list crew: %w", err)
	}
	return crew, nil
}

// Create creates a new crew member.
func (s *Service) Create(ctx context.Context, crew *Crew) error {
	if err := s.repo.Create(ctx, crew); err != nil {
		return fmt.Errorf("create crew: %w", err)
	}
	s.logger.Info("crew created",
		slog.String("id", crew.ID.String()),
		slog.String("name", crew.Name),
	)
	return nil
}

// Update updates an existing crew member.
func (s *Service) Update(ctx context.Context, crew *Crew) error {
	if err := s.repo.Update(ctx, crew); err != nil {
		return fmt.Errorf("update crew: %w", err)
	}
	s.logger.Info("crew updated", slog.String("id", crew.ID.String()))
	return nil
}

// Delete removes a crew member.
func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete crew: %w", err)
	}
	s.logger.Info("crew deleted", slog.String("id", id.String()))
	return nil
}

// Search searches crew members by name.
func (s *Service) Search(ctx context.Context, query string, limit, offset int) ([]Crew, error) {
	crew, err := s.repo.Search(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("search crew: %w", err)
	}
	return crew, nil
}

// GetWithNames retrieves a crew member with all their names/aliases.
func (s *Service) GetWithNames(ctx context.Context, id uuid.UUID) (*Crew, []CrewName, error) {
	crew, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, nil, fmt.Errorf("get crew: %w", err)
	}
	names, err := s.repo.ListNames(ctx, id)
	if err != nil {
		return nil, nil, fmt.Errorf("list crew names: %w", err)
	}
	return crew, names, nil
}

// AddName adds a name/alias to a crew member.
func (s *Service) AddName(ctx context.Context, crewID uuid.UUID, name string) error {
	if err := s.repo.AddName(ctx, crewID, name); err != nil {
		return fmt.Errorf("add crew name: %w", err)
	}
	return nil
}

// ListExpeditionCrew retrieves crew for an expedition (movie).
func (s *Service) ListExpeditionCrew(ctx context.Context, expeditionID uuid.UUID) ([]Crew, error) {
	crew, err := s.repo.ListExpeditionCrew(ctx, expeditionID)
	if err != nil {
		return nil, fmt.Errorf("list expedition crew: %w", err)
	}
	return crew, nil
}

// ListVoyageCrew retrieves crew for a voyage (scene).
func (s *Service) ListVoyageCrew(ctx context.Context, voyageID uuid.UUID) ([]Crew, error) {
	crew, err := s.repo.ListVoyageCrew(ctx, voyageID)
	if err != nil {
		return nil, fmt.Errorf("list voyage crew: %w", err)
	}
	return crew, nil
}
