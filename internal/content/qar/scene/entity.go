// Package scene provides adult scene domain models.
package scene

import (
	"time"

	"github.com/google/uuid"

	"github.com/lusoris/revenge/internal/content/shared"
)

// Scene represents an adult scene.
type Scene struct {
	shared.ContentEntity
	ReleaseDate    *time.Time
	RuntimeMinutes int
	Overview       string
	StudioID       *uuid.UUID
}
