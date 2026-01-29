// Package port provides adult studio domain models (QAR obfuscation: studios → ports).
package port

import (
	"context"

	"github.com/google/uuid"
)

// Repository defines the port (studio) data access interface.
type Repository interface {
	// GetByID retrieves a port by ID.
	GetByID(ctx context.Context, id uuid.UUID) (*Port, error)

	// List retrieves ports with pagination.
	List(ctx context.Context, limit, offset int) ([]Port, error)

	// ListRoot retrieves root ports (no parent).
	ListRoot(ctx context.Context, limit, offset int) ([]Port, error)

	// ListChildren retrieves child ports of a parent.
	ListChildren(ctx context.Context, parentID uuid.UUID) ([]Port, error)

	// Create creates a new port.
	Create(ctx context.Context, port *Port) error

	// Update updates an existing port.
	Update(ctx context.Context, port *Port) error

	// Delete removes a port by ID.
	Delete(ctx context.Context, id uuid.UUID) error

	// GetByStashDBID retrieves a port by StashDB ID.
	GetByStashDBID(ctx context.Context, stashdbID string) (*Port, error)

	// GetByTPDBID retrieves a port by TPDB ID.
	GetByTPDBID(ctx context.Context, tpdbID string) (*Port, error)

	// Search searches ports by name.
	Search(ctx context.Context, query string, limit, offset int) ([]Port, error)

	// CountExpeditions returns the number of expeditions for a port.
	CountExpeditions(ctx context.Context, id uuid.UUID) (int64, error)

	// CountVoyages returns the number of voyages for a port.
	CountVoyages(ctx context.Context, id uuid.UUID) (int64, error)
}
