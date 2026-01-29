// Package flag provides adult tag domain models (QAR obfuscation: tags → flags).
package flag

import (
	"context"

	"github.com/google/uuid"
)

// Repository defines the flag (tag) data access interface.
type Repository interface {
	// GetByID retrieves a flag by ID.
	GetByID(ctx context.Context, id uuid.UUID) (*Flag, error)

	// GetByName retrieves a flag by name.
	GetByName(ctx context.Context, name string) (*Flag, error)

	// List retrieves flags with pagination.
	List(ctx context.Context, limit, offset int) ([]Flag, error)

	// ListRoot retrieves root flags (no parent).
	ListRoot(ctx context.Context, limit, offset int) ([]Flag, error)

	// ListChildren retrieves child flags of a parent.
	ListChildren(ctx context.Context, parentID uuid.UUID) ([]Flag, error)

	// ListByWaters retrieves flags in a specific category.
	ListByWaters(ctx context.Context, waters string) ([]Flag, error)

	// Create creates a new flag.
	Create(ctx context.Context, flag *Flag) error

	// Update updates an existing flag.
	Update(ctx context.Context, flag *Flag) error

	// Delete removes a flag by ID.
	Delete(ctx context.Context, id uuid.UUID) error

	// GetByStashDBID retrieves a flag by StashDB ID.
	GetByStashDBID(ctx context.Context, stashdbID string) (*Flag, error)

	// Search searches flags by name.
	Search(ctx context.Context, query string, limit, offset int) ([]Flag, error)

	// ListExpeditionFlags retrieves flags for an expedition.
	ListExpeditionFlags(ctx context.Context, expeditionID uuid.UUID) ([]Flag, error)

	// ListVoyageFlags retrieves flags for a voyage.
	ListVoyageFlags(ctx context.Context, voyageID uuid.UUID) ([]Flag, error)

	// AddExpeditionFlag adds a flag to an expedition.
	AddExpeditionFlag(ctx context.Context, expeditionID, flagID uuid.UUID) error

	// AddVoyageFlag adds a flag to a voyage.
	AddVoyageFlag(ctx context.Context, voyageID, flagID uuid.UUID) error

	// RemoveExpeditionFlag removes a flag from an expedition.
	RemoveExpeditionFlag(ctx context.Context, expeditionID, flagID uuid.UUID) error

	// RemoveVoyageFlag removes a flag from a voyage.
	RemoveVoyageFlag(ctx context.Context, voyageID, flagID uuid.UUID) error

	// ClearExpeditionFlags removes all flags from an expedition.
	ClearExpeditionFlags(ctx context.Context, expeditionID uuid.UUID) error

	// ClearVoyageFlags removes all flags from a voyage.
	ClearVoyageFlags(ctx context.Context, voyageID uuid.UUID) error
}
