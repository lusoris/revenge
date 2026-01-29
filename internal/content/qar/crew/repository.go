// Package crew provides adult performer domain models (QAR obfuscation: performers → crew).
package crew

import (
	"context"

	"github.com/google/uuid"
)

// Repository defines the crew (performer) data access interface.
type Repository interface {
	// GetByID retrieves a crew member by ID.
	GetByID(ctx context.Context, id uuid.UUID) (*Crew, error)

	// List retrieves crew members with pagination.
	List(ctx context.Context, limit, offset int) ([]Crew, error)

	// Create creates a new crew member.
	Create(ctx context.Context, crew *Crew) error

	// Update updates an existing crew member.
	Update(ctx context.Context, crew *Crew) error

	// Delete removes a crew member by ID.
	Delete(ctx context.Context, id uuid.UUID) error

	// GetByCharter retrieves a crew member by StashDB ID.
	GetByCharter(ctx context.Context, charter string) (*Crew, error)

	// GetByRegistry retrieves a crew member by TPDB ID.
	GetByRegistry(ctx context.Context, registry string) (*Crew, error)

	// Search searches crew members by name.
	Search(ctx context.Context, query string, limit, offset int) ([]Crew, error)

	// ListNames retrieves all names/aliases for a crew member.
	ListNames(ctx context.Context, crewID uuid.UUID) ([]CrewName, error)

	// AddName adds a name/alias to a crew member.
	AddName(ctx context.Context, crewID uuid.UUID, name string) error

	// RemoveName removes a name/alias from a crew member.
	RemoveName(ctx context.Context, crewID uuid.UUID, name string) error

	// ListPortraits retrieves all portraits for a crew member.
	ListPortraits(ctx context.Context, crewID uuid.UUID) ([]CrewPortrait, error)

	// AddPortrait adds a portrait to a crew member.
	AddPortrait(ctx context.Context, portrait *CrewPortrait) error

	// SetPrimaryPortrait sets the primary portrait for a crew member.
	SetPrimaryPortrait(ctx context.Context, crewID, portraitID uuid.UUID) error

	// ListExpeditionCrew retrieves crew for an expedition.
	ListExpeditionCrew(ctx context.Context, expeditionID uuid.UUID) ([]Crew, error)

	// ListVoyageCrew retrieves crew for a voyage.
	ListVoyageCrew(ctx context.Context, voyageID uuid.UUID) ([]Crew, error)

	// AddExpeditionCrew adds a crew member to an expedition.
	AddExpeditionCrew(ctx context.Context, expeditionID, crewID uuid.UUID, characterName string) error

	// AddVoyageCrew adds a crew member to a voyage.
	AddVoyageCrew(ctx context.Context, voyageID, crewID uuid.UUID, role string) error

	// RemoveExpeditionCrew removes a crew member from an expedition.
	RemoveExpeditionCrew(ctx context.Context, expeditionID, crewID uuid.UUID) error

	// RemoveVoyageCrew removes a crew member from a voyage.
	RemoveVoyageCrew(ctx context.Context, voyageID, crewID uuid.UUID) error
}
