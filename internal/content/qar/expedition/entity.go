// Package expedition provides adult movie domain models (QAR obfuscation: movies → expeditions).
package expedition

import (
	"time"

	"github.com/google/uuid"

	"github.com/lusoris/revenge/internal/content/shared"
)

// Expedition represents an adult movie (obfuscated as "expedition").
type Expedition struct {
	shared.ContentEntity
	FleetID      uuid.UUID  // Library reference
	LaunchDate   *time.Time // release_date
	RuntimeTicks int64
	Overview     string
	PortID       *uuid.UUID // studio_id
	Director     string
	Series       string
	Coordinates  string // phash
	Charter      string // stashdb_id
	Registry     string // tpdb_id
	WhisparrID   *int
	HasFile      bool
	IsHDR        bool
	Is3D         bool
}
