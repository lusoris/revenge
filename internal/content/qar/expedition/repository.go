// Package expedition provides adult movie domain models (QAR obfuscation: movies → expeditions).
package expedition

import (
	"context"

	"github.com/google/uuid"
)

// Repository defines the expedition (adult movie) data access interface.
type Repository interface {
	// GetByID retrieves an expedition by ID.
	GetByID(ctx context.Context, id uuid.UUID) (*Expedition, error)

	// List retrieves expeditions with pagination.
	List(ctx context.Context, limit, offset int) ([]Expedition, error)

	// ListByFleet retrieves expeditions for a specific fleet (library).
	ListByFleet(ctx context.Context, fleetID uuid.UUID, limit, offset int) ([]Expedition, error)

	// Create creates a new expedition.
	Create(ctx context.Context, expedition *Expedition) error

	// Update updates an existing expedition.
	Update(ctx context.Context, expedition *Expedition) error

	// Delete removes an expedition by ID.
	Delete(ctx context.Context, id uuid.UUID) error

	// GetByPath retrieves an expedition by file path.
	GetByPath(ctx context.Context, path string) (*Expedition, error)

	// GetByCoordinates retrieves an expedition by pHash (fingerprint).
	GetByCoordinates(ctx context.Context, coordinates string) (*Expedition, error)

	// GetByCharter retrieves an expedition by StashDB ID.
	GetByCharter(ctx context.Context, charter string) (*Expedition, error)

	// CountByFleet returns the number of expeditions in a fleet.
	CountByFleet(ctx context.Context, fleetID uuid.UUID) (int64, error)

	// Search searches expeditions by title.
	Search(ctx context.Context, query string, limit, offset int) ([]Expedition, error)

	// ListByCrewID retrieves expeditions featuring a specific crew member (performer).
	ListByCrewID(ctx context.Context, crewID uuid.UUID, limit, offset int) ([]Expedition, error)

	// CountByCrewID returns the number of expeditions featuring a crew member.
	CountByCrewID(ctx context.Context, crewID uuid.UUID) (int64, error)

	// ListByPortID retrieves expeditions from a specific port (studio).
	ListByPortID(ctx context.Context, portID uuid.UUID, limit, offset int) ([]Expedition, error)

	// CountByPortID returns the number of expeditions from a port.
	CountByPortID(ctx context.Context, portID uuid.UUID) (int64, error)

	// ListByFlagID retrieves expeditions with a specific flag (tag).
	ListByFlagID(ctx context.Context, flagID uuid.UUID, limit, offset int) ([]Expedition, error)

	// CountByFlagID returns the number of expeditions with a flag.
	CountByFlagID(ctx context.Context, flagID uuid.UUID) (int64, error)
}
