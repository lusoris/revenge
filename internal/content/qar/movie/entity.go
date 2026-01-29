// Package movie provides adult movie domain models.
package movie

import (
	"time"

	"github.com/google/uuid"

	"github.com/lusoris/revenge/internal/content/shared"
)

// Movie represents an adult movie.
type Movie struct {
	shared.ContentEntity
	ReleaseDate  *time.Time
	RuntimeTicks int64
	Overview     string
	StudioID     *uuid.UUID
	Director     string
	Series       string
}
