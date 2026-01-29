package scene

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

// Repository errors.
var (
	ErrSceneNotFound = errors.New("adult scene not found")
)

// ListParams contains pagination parameters.
type ListParams struct {
	Limit  int
	Offset int
}

// DefaultListParams returns default list parameters.
func DefaultListParams() ListParams {
	return ListParams{Limit: 20, Offset: 0}
}

// Repository defines data access for adult scenes.
type Repository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*Scene, error)
	List(ctx context.Context, params ListParams) ([]*Scene, error)
	ListByLibrary(ctx context.Context, libraryID uuid.UUID, params ListParams) ([]*Scene, error)
	Create(ctx context.Context, scene *Scene) error
	Update(ctx context.Context, scene *Scene) error
	Delete(ctx context.Context, id uuid.UUID) error
}
