// Package fleet provides adult library domain models (QAR obfuscation: libraries → fleets).
package fleet

import (
	"context"

	"github.com/google/uuid"
)

// Repository defines the fleet (adult library) data access interface.
type Repository interface {
	// GetByID retrieves a fleet by ID.
	GetByID(ctx context.Context, id uuid.UUID) (*Fleet, error)

	// List retrieves fleets with pagination.
	List(ctx context.Context, limit, offset int) ([]Fleet, error)

	// ListByOwner retrieves fleets owned by a specific user.
	ListByOwner(ctx context.Context, ownerID uuid.UUID) ([]Fleet, error)

	// ListByType retrieves fleets of a specific type.
	ListByType(ctx context.Context, fleetType FleetType, limit, offset int) ([]Fleet, error)

	// Create creates a new fleet.
	Create(ctx context.Context, fleet *Fleet) error

	// Update updates an existing fleet.
	Update(ctx context.Context, fleet *Fleet) error

	// Delete removes a fleet by ID.
	Delete(ctx context.Context, id uuid.UUID) error

	// GetStats returns statistics for a fleet.
	GetStats(ctx context.Context, id uuid.UUID) (*FleetStats, error)

	// CountExpeditions returns the number of expeditions in a fleet.
	CountExpeditions(ctx context.Context, id uuid.UUID) (int64, error)

	// CountVoyages returns the number of voyages in a fleet.
	CountVoyages(ctx context.Context, id uuid.UUID) (int64, error)
}
