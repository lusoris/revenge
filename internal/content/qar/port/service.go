// Package port provides adult studio domain models (QAR obfuscation: studios → ports).
package port

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
)

// Service provides port (studio) business logic.
type Service struct {
	repo   Repository
	logger *slog.Logger
}

// NewService creates a new port service.
func NewService(repo Repository, logger *slog.Logger) *Service {
	return &Service{
		repo:   repo,
		logger: logger.With(slog.String("service", "qar.port")),
	}
}

// GetByID retrieves a port by ID.
func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*Port, error) {
	port, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get port: %w", err)
	}
	return port, nil
}

// List retrieves ports with pagination.
func (s *Service) List(ctx context.Context, limit, offset int) ([]Port, error) {
	ports, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list ports: %w", err)
	}
	return ports, nil
}

// Create creates a new port.
func (s *Service) Create(ctx context.Context, port *Port) error {
	if err := s.repo.Create(ctx, port); err != nil {
		return fmt.Errorf("create port: %w", err)
	}
	s.logger.Info("port created",
		slog.String("id", port.ID.String()),
		slog.String("name", port.Name),
	)
	return nil
}

// Update updates an existing port.
func (s *Service) Update(ctx context.Context, port *Port) error {
	if err := s.repo.Update(ctx, port); err != nil {
		return fmt.Errorf("update port: %w", err)
	}
	s.logger.Info("port updated", slog.String("id", port.ID.String()))
	return nil
}

// Delete removes a port.
func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete port: %w", err)
	}
	s.logger.Info("port deleted", slog.String("id", id.String()))
	return nil
}

// Search searches ports by name.
func (s *Service) Search(ctx context.Context, query string, limit, offset int) ([]Port, error) {
	ports, err := s.repo.Search(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("search ports: %w", err)
	}
	return ports, nil
}

// GetWithChildren retrieves a port with its child ports.
func (s *Service) GetWithChildren(ctx context.Context, id uuid.UUID) (*Port, []Port, error) {
	port, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, nil, fmt.Errorf("get port: %w", err)
	}
	children, err := s.repo.ListChildren(ctx, id)
	if err != nil {
		return nil, nil, fmt.Errorf("list port children: %w", err)
	}
	return port, children, nil
}
