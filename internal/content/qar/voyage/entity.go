// Package voyage provides adult scene domain models (QAR obfuscation: scenes → voyages).
package voyage

import (
	"time"

	"github.com/google/uuid"

	"github.com/lusoris/revenge/internal/content/shared"
)

// Voyage represents an adult scene (obfuscated as "voyage").
type Voyage struct {
	shared.ContentEntity
	FleetID     uuid.UUID  // Library reference
	LaunchDate  *time.Time // release_date
	Distance    int        // runtime_minutes
	Overview    string
	PortID      *uuid.UUID // studio_id
	Coordinates string     // phash
	Oshash      string
	MD5         string
	CoverPath   string
	Charter     string // stashdb_id
	Registry    string // tpdb_id
	StashID     string
	WhisparrID  *int
}
