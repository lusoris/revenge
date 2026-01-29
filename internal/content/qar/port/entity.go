// Package port provides adult studio domain models (QAR obfuscation: studios → ports).
package port

import (
	"time"

	"github.com/google/uuid"
)

// Port represents an adult studio (obfuscated as "port").
type Port struct {
	ID        uuid.UUID
	Name      string
	ParentID  *uuid.UUID // Network/parent studio
	StashDBID string
	TPDBID    string
	URL       string
	LogoPath  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
