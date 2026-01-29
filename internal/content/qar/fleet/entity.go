// Package fleet provides adult library domain models (QAR obfuscation: libraries → fleets).
package fleet

import (
	"time"

	"github.com/google/uuid"
)

// FleetType represents the type of content in a fleet.
type FleetType string

const (
	FleetTypeExpedition FleetType = "expedition" // Full-length movies
	FleetTypeVoyage     FleetType = "voyage"     // Scenes
)

// Fleet represents an adult content library (obfuscated as "fleet").
type Fleet struct {
	ID                uuid.UUID
	Name              string
	FleetType         FleetType
	Paths             []string
	StashDBEndpoint   string
	TPDBEnabled       bool
	WhisparrSync      bool
	AutoTagCrew       bool
	FingerprintOnScan bool
	OwnerUserID       *uuid.UUID
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// FleetStats contains statistics about a fleet.
type FleetStats struct {
	ExpeditionCount int64
	VoyageCount     int64
	TotalSizeBytes  int64
}
