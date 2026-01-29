// Package crew provides adult performer domain models (QAR obfuscation: performers → crew).
package crew

import (
	"time"

	"github.com/google/uuid"
)

// Crew represents an adult performer (obfuscated as "crew").
type Crew struct {
	ID             uuid.UUID
	Name           string
	Disambiguation string
	Gender         string
	Christening    *time.Time // birth_date
	DeathDate      *time.Time
	BirthCity      string
	Origin         string // ethnicity
	Nationality    string
	Rigging        string // hair_color
	Compass        string // eye_color
	HeightCM       *int
	WeightKG       *int
	Measurements   string
	CupSize        string
	BreastType     string
	Markings       string // tattoos
	Anchors        string // piercings
	MaidenVoyage   *int   // career_start
	LastPort       *int   // career_end
	Bio            string
	StashID        string
	Charter        string // stashdb_id
	Registry       string // tpdb_id
	Manifest       string // freeones_id
	Twitter        string
	Instagram      string
	ImagePath      string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// CrewName represents an alias for a crew member.
type CrewName struct {
	CrewID uuid.UUID
	Name   string
}

// CrewPortrait represents a portrait image for a crew member.
type CrewPortrait struct {
	ID           uuid.UUID
	CrewID       uuid.UUID
	Path         string
	Type         string
	Source       string
	PrimaryImage bool
	CreatedAt    time.Time
}
