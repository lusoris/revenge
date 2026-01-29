// Package fleet provides adult library domain models (QAR obfuscation: libraries → fleets).
package fleet

import (
	"context"
	"errors"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	adultdb "github.com/lusoris/revenge/internal/content/qar/db"
	"github.com/lusoris/revenge/internal/content/shared"
)

// LibraryService errors.
var (
	ErrFleetNotFound   = errors.New("fleet not found")
	ErrFleetAccessDenied = errors.New("access denied to fleet")
)

// LibraryService provides fleet (adult library) management implementing shared.LibraryProvider.
// QAR obfuscation: LibraryService manages "fleets" which are adult content libraries.
type LibraryService struct {
	queries *adultdb.Queries
	logger  *slog.Logger
}

// NewLibraryService creates a new fleet library service.
func NewLibraryService(queries *adultdb.Queries, logger *slog.Logger) *LibraryService {
	return &LibraryService{
		queries: queries,
		logger:  logger.With("service", "fleet_library"),
	}
}

// Ensure LibraryService implements LibraryProvider.
var _ shared.LibraryProvider = (*LibraryService)(nil)

// ModuleName returns the module identifier.
func (s *LibraryService) ModuleName() string {
	return "qar"
}

// ListLibraries returns all fleets accessible by the user.
// For adult content, only fleets owned by the user are returned.
func (s *LibraryService) ListLibraries(ctx context.Context, userID uuid.UUID) ([]shared.LibraryInfo, error) {
	// Adult content: only show fleets owned by the user
	ownerUUID := pgtype.UUID{Bytes: userID, Valid: true}
	fleets, err := s.queries.ListFleetsByOwner(ctx, ownerUUID)
	if err != nil {
		return nil, err
	}

	result := make([]shared.LibraryInfo, len(fleets))
	for i, fleet := range fleets {
		// Get item count (expeditions + voyages)
		count := s.getFleetItemCount(ctx, fleet.ID)
		result[i] = s.toLibraryInfo(fleet, count)
	}

	return result, nil
}

// GetLibrary returns a specific fleet by ID.
func (s *LibraryService) GetLibrary(ctx context.Context, libraryID uuid.UUID) (*shared.LibraryInfo, error) {
	fleet, err := s.queries.GetFleetByID(ctx, libraryID)
	if err != nil {
		return nil, ErrFleetNotFound
	}

	count := s.getFleetItemCount(ctx, libraryID)
	info := s.toLibraryInfo(fleet, count)
	return &info, nil
}

// CreateLibrary creates a new fleet.
func (s *LibraryService) CreateLibrary(ctx context.Context, req shared.CreateLibraryRequest) (*shared.LibraryInfo, error) {
	params := adultdb.CreateFleetParams{
		Name:              req.Name,
		FleetType:         "expedition", // Default to expedition (adult movie)
		Paths:             req.Paths,
		TpdbEnabled:       true,
		WhisparrSync:      false,
		AutoTagCrew:       true,
		FingerprintOnScan: true,
	}

	// Apply module-specific settings if provided
	if settings, ok := req.Settings.(*FleetSettings); ok && settings != nil {
		params.FleetType = settings.FleetType
		if settings.StashDBEndpoint != "" {
			params.StashdbEndpoint = &settings.StashDBEndpoint
		}
		params.TpdbEnabled = settings.TPDBEnabled
		params.WhisparrSync = settings.WhisparrSync
		params.AutoTagCrew = settings.AutoTagCrew
		params.FingerprintOnScan = settings.FingerprintOnScan
		if settings.OwnerUserID != uuid.Nil {
			params.OwnerUserID = pgtype.UUID{Bytes: settings.OwnerUserID, Valid: true}
		}
	}

	fleet, err := s.queries.CreateFleet(ctx, params)
	if err != nil {
		return nil, err
	}

	s.logger.Info("fleet created",
		"id", fleet.ID,
		"name", fleet.Name,
		"type", fleet.FleetType,
	)

	info := s.toLibraryInfo(fleet, 0)
	return &info, nil
}

// UpdateLibrary updates fleet settings.
func (s *LibraryService) UpdateLibrary(ctx context.Context, libraryID uuid.UUID, req shared.UpdateLibraryRequest) (*shared.LibraryInfo, error) {
	// Get current fleet
	current, err := s.queries.GetFleetByID(ctx, libraryID)
	if err != nil {
		return nil, ErrFleetNotFound
	}

	// Build update params
	params := adultdb.UpdateFleetParams{
		ID:                libraryID,
		Name:              current.Name,
		FleetType:         current.FleetType,
		Paths:             current.Paths,
		StashdbEndpoint:   current.StashdbEndpoint,
		TpdbEnabled:       current.TpdbEnabled,
		WhisparrSync:      current.WhisparrSync,
		AutoTagCrew:       current.AutoTagCrew,
		FingerprintOnScan: current.FingerprintOnScan,
		OwnerUserID:       current.OwnerUserID,
	}

	// Apply updates
	if req.Name != nil {
		params.Name = *req.Name
	}
	if len(req.Paths) > 0 {
		params.Paths = req.Paths
	}

	// Apply module-specific settings if provided
	if settings, ok := req.Settings.(*FleetSettings); ok && settings != nil {
		if settings.FleetType != "" {
			params.FleetType = settings.FleetType
		}
		if settings.StashDBEndpoint != "" {
			params.StashdbEndpoint = &settings.StashDBEndpoint
		}
		params.TpdbEnabled = settings.TPDBEnabled
		params.WhisparrSync = settings.WhisparrSync
		params.AutoTagCrew = settings.AutoTagCrew
		params.FingerprintOnScan = settings.FingerprintOnScan
		if settings.OwnerUserID != uuid.Nil {
			params.OwnerUserID = pgtype.UUID{Bytes: settings.OwnerUserID, Valid: true}
		}
	}

	fleet, err := s.queries.UpdateFleet(ctx, params)
	if err != nil {
		return nil, err
	}

	s.logger.Info("fleet updated",
		"id", fleet.ID,
		"name", fleet.Name,
	)

	count := s.getFleetItemCount(ctx, libraryID)
	info := s.toLibraryInfo(fleet, count)
	return &info, nil
}

// DeleteLibrary removes a fleet and optionally its content.
func (s *LibraryService) DeleteLibrary(ctx context.Context, libraryID uuid.UUID, deleteContent bool) error {
	// Verify fleet exists
	fleet, err := s.queries.GetFleetByID(ctx, libraryID)
	if err != nil {
		return ErrFleetNotFound
	}

	if err := s.queries.DeleteFleet(ctx, libraryID); err != nil {
		return err
	}

	s.logger.Info("fleet deleted",
		"id", fleet.ID,
		"name", fleet.Name,
		"delete_content", deleteContent,
	)

	return nil
}

// ScanLibrary triggers a fleet scan.
func (s *LibraryService) ScanLibrary(ctx context.Context, libraryID uuid.UUID, fullScan bool) error {
	// Verify fleet exists
	fleet, err := s.queries.GetFleetByID(ctx, libraryID)
	if err != nil {
		return ErrFleetNotFound
	}

	s.logger.Info("fleet scan triggered",
		"fleet_id", libraryID,
		"fleet_name", fleet.Name,
		"full_scan", fullScan,
	)

	// TODO: Enqueue a River job for the actual scanning
	return nil
}

// getFleetItemCount returns the total item count (expeditions + voyages).
func (s *LibraryService) getFleetItemCount(ctx context.Context, fleetID uuid.UUID) int64 {
	expCount, err := s.queries.CountFleetExpeditions(ctx, fleetID)
	if err != nil {
		s.logger.Warn("failed to count expeditions", "fleet_id", fleetID, "error", err)
		expCount = 0
	}

	voyCount, err := s.queries.CountFleetVoyages(ctx, fleetID)
	if err != nil {
		s.logger.Warn("failed to count voyages", "fleet_id", fleetID, "error", err)
		voyCount = 0
	}

	return expCount + voyCount
}

// toLibraryInfo converts a database QarFleet to shared.LibraryInfo.
func (s *LibraryService) toLibraryInfo(fleet adultdb.QarFleet, itemCount int64) shared.LibraryInfo {
	var ownerID uuid.UUID
	if fleet.OwnerUserID.Valid {
		ownerID = fleet.OwnerUserID.Bytes
	}

	return shared.LibraryInfo{
		ID:        fleet.ID,
		Module:    "qar",
		Name:      fleet.Name,
		Paths:     fleet.Paths,
		IsAdult:   true, // Fleets are ALWAYS adult content
		ItemCount: itemCount,
		Settings: &FleetSettings{
			FleetType:         fleet.FleetType,
			StashDBEndpoint:   derefString(fleet.StashdbEndpoint),
			TPDBEnabled:       fleet.TpdbEnabled,
			WhisparrSync:      fleet.WhisparrSync,
			AutoTagCrew:       fleet.AutoTagCrew,
			FingerprintOnScan: fleet.FingerprintOnScan,
			OwnerUserID:       ownerID,
		},
	}
}

// FleetSettings contains fleet-specific settings.
// QAR obfuscation: these are adult library settings.
type FleetSettings struct {
	FleetType         string    `json:"fleet_type"`          // "expedition" or "voyage"
	StashDBEndpoint   string    `json:"stashdb_endpoint,omitempty"`
	TPDBEnabled       bool      `json:"tpdb_enabled"`
	WhisparrSync      bool      `json:"whisparr_sync"`
	AutoTagCrew       bool      `json:"auto_tag_crew"`
	FingerprintOnScan bool      `json:"fingerprint_on_scan"`
	OwnerUserID       uuid.UUID `json:"owner_user_id,omitempty"`
}

// derefString safely dereferences a string pointer.
func derefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
