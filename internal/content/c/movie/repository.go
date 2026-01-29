package movie

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

// Repository errors.
var (
	ErrMovieNotFound = errors.New("adult movie not found")
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

// Repository defines data access for adult movies.
type Repository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*Movie, error)
	List(ctx context.Context, params ListParams) ([]*Movie, error)
	ListByLibrary(ctx context.Context, libraryID uuid.UUID, params ListParams) ([]*Movie, error)
	Create(ctx context.Context, movie *Movie) error
	Update(ctx context.Context, movie *Movie) error
	Delete(ctx context.Context, id uuid.UUID) error
}
