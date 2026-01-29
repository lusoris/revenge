// Package flag provides adult tag domain models (QAR obfuscation: tags → flags).
package flag

import (
	"time"

	"github.com/google/uuid"
)

// Flag represents an adult content tag (obfuscated as "flag").
type Flag struct {
	ID          uuid.UUID
	Name        string
	Description string
	ParentID    *uuid.UUID
	StashDBID   string
	Waters      string // category
	CreatedAt   time.Time
}
