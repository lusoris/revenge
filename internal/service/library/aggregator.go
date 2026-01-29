// Package library provides aggregated library management across all content modules.
package library

import (
	"context"
	"errors"
	"log/slog"
	"sync"

	"github.com/google/uuid"

	"github.com/lusoris/revenge/internal/content/shared"
)

var (
	// ErrLibraryNotFound indicates the library was not found.
	ErrLibraryNotFound = errors.New("library not found")
	// ErrAccessDenied indicates the user doesn't have access to the library.
	ErrAccessDenied = errors.New("access denied")
	// ErrInvalidLibraryType indicates an invalid library type.
	ErrInvalidLibraryType = errors.New("invalid library type")
	// ErrModuleNotFound indicates the module for a library was not found.
	ErrModuleNotFound = errors.New("module not found")
)

// Service aggregates library operations across all content modules.
// It uses the LibraryProvider interface to delegate to module-specific implementations.
type Service struct {
	providers map[string]shared.LibraryProvider
	logger    *slog.Logger
	mu        sync.RWMutex
}

// NewService creates a new library aggregator service.
func NewService(providers []shared.LibraryProvider, logger *slog.Logger) *Service {
	providerMap := make(map[string]shared.LibraryProvider, len(providers))
	for _, p := range providers {
		providerMap[p.ModuleName()] = p
	}

	return &Service{
		providers: providerMap,
		logger:    logger.With(slog.String("service", "library_aggregator")),
	}
}

// WithAccess contains a library and access info.
type WithAccess struct {
	Library   shared.LibraryInfo
	CanManage bool
}

// ListAll returns all libraries from all modules (admin only).
func (s *Service) ListAll(ctx context.Context) ([]WithAccess, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []WithAccess

	for _, provider := range s.providers {
		libs, err := provider.ListLibraries(ctx, uuid.Nil) // Nil UUID = list all
		if err != nil {
			s.logger.Warn("failed to list libraries from module",
				"module", provider.ModuleName(),
				"error", err,
			)
			continue
		}

		for _, lib := range libs {
			result = append(result, WithAccess{
				Library:   lib,
				CanManage: true, // Admin can manage all
			})
		}
	}

	return result, nil
}

// ListForUser returns all libraries accessible by a user.
func (s *Service) ListForUser(ctx context.Context, userID uuid.UUID) ([]WithAccess, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []WithAccess

	for _, provider := range s.providers {
		libs, err := provider.ListLibraries(ctx, userID)
		if err != nil {
			s.logger.Warn("failed to list libraries from module",
				"module", provider.ModuleName(),
				"user_id", userID,
				"error", err,
			)
			continue
		}

		for _, lib := range libs {
			result = append(result, WithAccess{
				Library:   lib,
				CanManage: false, // TODO: Check actual permissions
			})
		}
	}

	return result, nil
}

// GetByID retrieves a library by ID, searching all modules.
func (s *Service) GetByID(ctx context.Context, libraryID uuid.UUID) (*shared.LibraryInfo, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, provider := range s.providers {
		lib, err := provider.GetLibrary(ctx, libraryID)
		if err == nil && lib != nil {
			return lib, nil
		}
	}

	return nil, ErrLibraryNotFound
}

// UserCanAccess checks if a user can access a library.
func (s *Service) UserCanAccess(ctx context.Context, userID, libraryID uuid.UUID) (bool, error) {
	// Get library to find its module
	lib, err := s.GetByID(ctx, libraryID)
	if err != nil {
		return false, err
	}

	// Get provider for this module
	s.mu.RLock()
	provider, ok := s.providers[lib.Module]
	s.mu.RUnlock()

	if !ok {
		return false, ErrModuleNotFound
	}

	// Check if user has access by listing their libraries
	libs, err := provider.ListLibraries(ctx, userID)
	if err != nil {
		return false, err
	}

	for _, l := range libs {
		if l.ID == libraryID {
			return true, nil
		}
	}

	return false, nil
}

// CreateParams contains parameters for creating a library.
type CreateParams struct {
	Name        string
	LibraryType string // Module type: "movie", "tvshow", "qar"
	Paths       []string
	Settings    any // Module-specific settings
}

// Create creates a new library in the appropriate module.
func (s *Service) Create(ctx context.Context, params CreateParams) (*shared.LibraryInfo, error) {
	s.mu.RLock()
	provider, ok := s.providers[params.LibraryType]
	s.mu.RUnlock()

	if !ok {
		return nil, ErrInvalidLibraryType
	}

	return provider.CreateLibrary(ctx, shared.CreateLibraryRequest{
		Name:     params.Name,
		Paths:    params.Paths,
		Settings: params.Settings,
	})
}

// UpdateParams contains parameters for updating a library.
type UpdateParams struct {
	ID       uuid.UUID
	Name     *string
	Paths    []string
	Settings any
}

// Update updates a library.
func (s *Service) Update(ctx context.Context, params UpdateParams) (*shared.LibraryInfo, error) {
	// Find the library's module
	lib, err := s.GetByID(ctx, params.ID)
	if err != nil {
		return nil, err
	}

	s.mu.RLock()
	provider, ok := s.providers[lib.Module]
	s.mu.RUnlock()

	if !ok {
		return nil, ErrModuleNotFound
	}

	return provider.UpdateLibrary(ctx, params.ID, shared.UpdateLibraryRequest{
		Name:     params.Name,
		Paths:    params.Paths,
		Settings: params.Settings,
	})
}

// Delete removes a library.
func (s *Service) Delete(ctx context.Context, libraryID uuid.UUID) error {
	// Find the library's module
	lib, err := s.GetByID(ctx, libraryID)
	if err != nil {
		return err
	}

	s.mu.RLock()
	provider, ok := s.providers[lib.Module]
	s.mu.RUnlock()

	if !ok {
		return ErrModuleNotFound
	}

	return provider.DeleteLibrary(ctx, libraryID, false)
}

// Scan triggers a library scan.
func (s *Service) Scan(ctx context.Context, libraryID uuid.UUID, fullScan bool) error {
	// Find the library's module
	lib, err := s.GetByID(ctx, libraryID)
	if err != nil {
		return err
	}

	s.mu.RLock()
	provider, ok := s.providers[lib.Module]
	s.mu.RUnlock()

	if !ok {
		return ErrModuleNotFound
	}

	return provider.ScanLibrary(ctx, libraryID, fullScan)
}
