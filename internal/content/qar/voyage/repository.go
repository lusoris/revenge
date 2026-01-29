// Package voyage provides adult scene domain models (QAR obfuscation: scenes → voyages).
package voyage

import (
	"context"

	"github.com/google/uuid"
)

// Repository defines the voyage (adult scene) data access interface.
type Repository interface {
	// GetByID retrieves a voyage by ID.
	GetByID(ctx context.Context, id uuid.UUID) (*Voyage, error)

	// List retrieves voyages with pagination.
	List(ctx context.Context, limit, offset int) ([]Voyage, error)

	// ListByFleet retrieves voyages for a specific fleet (library).
	ListByFleet(ctx context.Context, fleetID uuid.UUID, limit, offset int) ([]Voyage, error)

	// ListByPort retrieves voyages for a specific port (studio).
	ListByPort(ctx context.Context, portID uuid.UUID, limit, offset int) ([]Voyage, error)

	// Create creates a new voyage.
	Create(ctx context.Context, voyage *Voyage) error

	// Update updates an existing voyage.
	Update(ctx context.Context, voyage *Voyage) error

	// Delete removes a voyage by ID.
	Delete(ctx context.Context, id uuid.UUID) error

	// GetByPath retrieves a voyage by file path.
	GetByPath(ctx context.Context, path string) (*Voyage, error)

	// GetByOshash retrieves a voyage by OpenSubtitles hash.
	GetByOshash(ctx context.Context, oshash string) (*Voyage, error)

	// GetByCoordinates retrieves a voyage by pHash (fingerprint).
	GetByCoordinates(ctx context.Context, coordinates string) (*Voyage, error)

	// GetByCharter retrieves a voyage by StashDB ID.
	GetByCharter(ctx context.Context, charter string) (*Voyage, error)

	// CountByFleet returns the number of voyages in a fleet.
	CountByFleet(ctx context.Context, fleetID uuid.UUID) (int64, error)

	// Search searches voyages by title.
	Search(ctx context.Context, query string, limit, offset int) ([]Voyage, error)
}
