// Package flag provides adult tag domain models (QAR obfuscation: tags → flags).
package flag

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
)

// Service provides flag (tag) business logic.
type Service struct {
	repo   Repository
	logger *slog.Logger
}

// NewService creates a new flag service.
func NewService(repo Repository, logger *slog.Logger) *Service {
	return &Service{
		repo:   repo,
		logger: logger.With(slog.String("service", "qar.flag")),
	}
}

// GetByID retrieves a flag by ID.
func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*Flag, error) {
	flag, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get flag: %w", err)
	}
	return flag, nil
}

// GetOrCreate retrieves a flag by name, creating it if it doesn't exist.
func (s *Service) GetOrCreate(ctx context.Context, name string) (*Flag, error) {
	flag, err := s.repo.GetByName(ctx, name)
	if err == nil {
		return flag, nil
	}

	// Create new flag
	newFlag := &Flag{
		ID:   uuid.New(),
		Name: name,
	}
	if err := s.repo.Create(ctx, newFlag); err != nil {
		return nil, fmt.Errorf("create flag: %w", err)
	}
	s.logger.Info("flag created", slog.String("name", name))
	return newFlag, nil
}

// List retrieves flags with pagination.
func (s *Service) List(ctx context.Context, limit, offset int) ([]Flag, error) {
	flags, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list flags: %w", err)
	}
	return flags, nil
}

// Create creates a new flag.
func (s *Service) Create(ctx context.Context, flag *Flag) error {
	if err := s.repo.Create(ctx, flag); err != nil {
		return fmt.Errorf("create flag: %w", err)
	}
	s.logger.Info("flag created",
		slog.String("id", flag.ID.String()),
		slog.String("name", flag.Name),
	)
	return nil
}

// Update updates an existing flag.
func (s *Service) Update(ctx context.Context, flag *Flag) error {
	if err := s.repo.Update(ctx, flag); err != nil {
		return fmt.Errorf("update flag: %w", err)
	}
	s.logger.Info("flag updated", slog.String("id", flag.ID.String()))
	return nil
}

// Delete removes a flag.
func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete flag: %w", err)
	}
	s.logger.Info("flag deleted", slog.String("id", id.String()))
	return nil
}

// Search searches flags by name.
func (s *Service) Search(ctx context.Context, query string, limit, offset int) ([]Flag, error) {
	flags, err := s.repo.Search(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("search flags: %w", err)
	}
	return flags, nil
}

// ListExpeditionFlags retrieves flags for an expedition.
func (s *Service) ListExpeditionFlags(ctx context.Context, expeditionID uuid.UUID) ([]Flag, error) {
	flags, err := s.repo.ListExpeditionFlags(ctx, expeditionID)
	if err != nil {
		return nil, fmt.Errorf("list expedition flags: %w", err)
	}
	return flags, nil
}

// ListVoyageFlags retrieves flags for a voyage.
func (s *Service) ListVoyageFlags(ctx context.Context, voyageID uuid.UUID) ([]Flag, error) {
	flags, err := s.repo.ListVoyageFlags(ctx, voyageID)
	if err != nil {
		return nil, fmt.Errorf("list voyage flags: %w", err)
	}
	return flags, nil
}

// SetExpeditionFlags sets the flags for an expedition (replaces existing).
func (s *Service) SetExpeditionFlags(ctx context.Context, expeditionID uuid.UUID, flagIDs []uuid.UUID) error {
	if err := s.repo.ClearExpeditionFlags(ctx, expeditionID); err != nil {
		return fmt.Errorf("clear expedition flags: %w", err)
	}
	for _, flagID := range flagIDs {
		if err := s.repo.AddExpeditionFlag(ctx, expeditionID, flagID); err != nil {
			return fmt.Errorf("add expedition flag: %w", err)
		}
	}
	return nil
}

// SetVoyageFlags sets the flags for a voyage (replaces existing).
func (s *Service) SetVoyageFlags(ctx context.Context, voyageID uuid.UUID, flagIDs []uuid.UUID) error {
	if err := s.repo.ClearVoyageFlags(ctx, voyageID); err != nil {
		return fmt.Errorf("clear voyage flags: %w", err)
	}
	for _, flagID := range flagIDs {
		if err := s.repo.AddVoyageFlag(ctx, voyageID, flagID); err != nil {
			return fmt.Errorf("add voyage flag: %w", err)
		}
	}
	return nil
}
