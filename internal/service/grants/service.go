// Package grants provides polymorphic resource-level access control.
package grants

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/lusoris/revenge/internal/infra/database/db"
)

var (
	// ErrGrantNotFound indicates the grant was not found.
	ErrGrantNotFound = errors.New("grant not found")
)

// GrantType represents the level of access granted.
type GrantType string

const (
	GrantView   GrantType = "view"   // Can view/browse the resource
	GrantEdit   GrantType = "edit"   // Can view + edit the resource
	GrantManage GrantType = "manage" // Can view + edit + manage items
	GrantOwner  GrantType = "owner"  // Full control including sharing
)

// ResourceType represents the type of resource being granted access to.
type ResourceType string

const (
	ResourceMovieLibrary ResourceType = "movie_library"
	ResourceTVLibrary    ResourceType = "tv_library"
	ResourceMusicLibrary ResourceType = "music_library"
	ResourceAdultLibrary ResourceType = "adult_library"
	ResourcePlaylist     ResourceType = "playlist"
	ResourceCollection   ResourceType = "collection"
)

// Grant represents a resource access grant.
type Grant struct {
	ID           uuid.UUID
	UserID       uuid.UUID
	ResourceType ResourceType
	ResourceID   uuid.UUID
	GrantType    GrantType
	GrantedBy    *uuid.UUID
	GrantedAt    time.Time
	ExpiresAt    *time.Time
}

// CreateParams contains parameters for creating a grant.
type CreateParams struct {
	UserID       uuid.UUID
	ResourceType ResourceType
	ResourceID   uuid.UUID
	GrantType    GrantType
	GrantedBy    uuid.UUID
	ExpiresAt    *time.Time
}

// Service provides resource grant operations.
type Service struct {
	queries *db.Queries
	logger  *slog.Logger
}

// NewService creates a new grants service.
func NewService(queries *db.Queries, logger *slog.Logger) *Service {
	return &Service{
		queries: queries,
		logger:  logger.With(slog.String("service", "grants")),
	}
}

// HasGrant checks if a user has access to a resource with at least one of the specified grant types.
func (s *Service) HasGrant(ctx context.Context, userID uuid.UUID, resourceType ResourceType, resourceID uuid.UUID, grantTypes ...GrantType) (bool, error) {
	if len(grantTypes) == 0 {
		grantTypes = []GrantType{GrantView}
	}

	// Convert to string slice for query
	types := make([]string, len(grantTypes))
	for i, gt := range grantTypes {
		types[i] = string(gt)
	}

	has, err := s.queries.HasResourceGrant(ctx, db.HasResourceGrantParams{
		UserID:       userID,
		ResourceType: string(resourceType),
		ResourceID:   resourceID,
		Column4:      types,
	})
	if err != nil {
		return false, fmt.Errorf("check grant: %w", err)
	}

	return has, nil
}

// HasViewGrant checks if a user can view a resource.
func (s *Service) HasViewGrant(ctx context.Context, userID uuid.UUID, resourceType ResourceType, resourceID uuid.UUID) (bool, error) {
	return s.HasGrant(ctx, userID, resourceType, resourceID, GrantView, GrantEdit, GrantManage, GrantOwner)
}

// HasEditGrant checks if a user can edit a resource.
func (s *Service) HasEditGrant(ctx context.Context, userID uuid.UUID, resourceType ResourceType, resourceID uuid.UUID) (bool, error) {
	return s.HasGrant(ctx, userID, resourceType, resourceID, GrantEdit, GrantManage, GrantOwner)
}

// HasManageGrant checks if a user can manage a resource.
func (s *Service) HasManageGrant(ctx context.Context, userID uuid.UUID, resourceType ResourceType, resourceID uuid.UUID) (bool, error) {
	return s.HasGrant(ctx, userID, resourceType, resourceID, GrantManage, GrantOwner)
}

// HasOwnerGrant checks if a user owns a resource.
func (s *Service) HasOwnerGrant(ctx context.Context, userID uuid.UUID, resourceType ResourceType, resourceID uuid.UUID) (bool, error) {
	return s.HasGrant(ctx, userID, resourceType, resourceID, GrantOwner)
}

// GetGrant retrieves a specific grant.
func (s *Service) GetGrant(ctx context.Context, userID uuid.UUID, resourceType ResourceType, resourceID uuid.UUID) (*Grant, error) {
	row, err := s.queries.GetResourceGrant(ctx, db.GetResourceGrantParams{
		UserID:       userID,
		ResourceType: string(resourceType),
		ResourceID:   resourceID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrGrantNotFound
		}
		return nil, fmt.Errorf("get grant: %w", err)
	}

	return s.rowToGrant(&row), nil
}

// CreateGrant creates or updates a resource grant.
func (s *Service) CreateGrant(ctx context.Context, params CreateParams) (*Grant, error) {
	var grantedBy pgtype.UUID
	if params.GrantedBy != uuid.Nil {
		grantedBy = pgtype.UUID{Bytes: params.GrantedBy, Valid: true}
	}

	var expiresAt pgtype.Timestamptz
	if params.ExpiresAt != nil {
		expiresAt = pgtype.Timestamptz{Time: *params.ExpiresAt, Valid: true}
	}

	row, err := s.queries.CreateResourceGrant(ctx, db.CreateResourceGrantParams{
		UserID:       params.UserID,
		ResourceType: string(params.ResourceType),
		ResourceID:   params.ResourceID,
		GrantType:    string(params.GrantType),
		GrantedBy:    grantedBy,
		ExpiresAt:    expiresAt,
	})
	if err != nil {
		return nil, fmt.Errorf("create grant: %w", err)
	}

	s.logger.Info("Resource grant created",
		slog.String("user_id", params.UserID.String()),
		slog.String("resource_type", string(params.ResourceType)),
		slog.String("resource_id", params.ResourceID.String()),
		slog.String("grant_type", string(params.GrantType)),
	)

	return s.rowToGrant(&row), nil
}

// DeleteGrant removes a specific grant.
func (s *Service) DeleteGrant(ctx context.Context, userID uuid.UUID, resourceType ResourceType, resourceID uuid.UUID) error {
	if err := s.queries.DeleteResourceGrant(ctx, db.DeleteResourceGrantParams{
		UserID:       userID,
		ResourceType: string(resourceType),
		ResourceID:   resourceID,
	}); err != nil {
		return fmt.Errorf("delete grant: %w", err)
	}

	s.logger.Info("Resource grant deleted",
		slog.String("user_id", userID.String()),
		slog.String("resource_type", string(resourceType)),
		slog.String("resource_id", resourceID.String()),
	)

	return nil
}

// DeleteByResource removes all grants for a resource (used when resource is deleted).
func (s *Service) DeleteByResource(ctx context.Context, resourceType ResourceType, resourceID uuid.UUID) error {
	if err := s.queries.DeleteResourceGrantsByResource(ctx, db.DeleteResourceGrantsByResourceParams{
		ResourceType: string(resourceType),
		ResourceID:   resourceID,
	}); err != nil {
		return fmt.Errorf("delete grants by resource: %w", err)
	}

	s.logger.Info("Resource grants deleted for resource",
		slog.String("resource_type", string(resourceType)),
		slog.String("resource_id", resourceID.String()),
	)

	return nil
}

// ListUserGrants returns all grants for a user.
func (s *Service) ListUserGrants(ctx context.Context, userID uuid.UUID, limit, offset int) ([]Grant, error) {
	rows, err := s.queries.ListUserResourceGrants(ctx, db.ListUserResourceGrantsParams{
		UserID: userID,
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("list user grants: %w", err)
	}

	return s.rowsToGrants(rows), nil
}

// ListUserGrantsByType returns grants of a specific type for a user.
func (s *Service) ListUserGrantsByType(ctx context.Context, userID uuid.UUID, resourceType ResourceType, limit, offset int) ([]Grant, error) {
	rows, err := s.queries.ListUserResourceGrantsByType(ctx, db.ListUserResourceGrantsByTypeParams{
		UserID:       userID,
		ResourceType: string(resourceType),
		Limit:        int32(limit),
		Offset:       int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("list user grants by type: %w", err)
	}

	return s.rowsToGrants(rows), nil
}

// ListResourceGrants returns all users with grants for a resource.
func (s *Service) ListResourceGrants(ctx context.Context, resourceType ResourceType, resourceID uuid.UUID) ([]Grant, error) {
	rows, err := s.queries.ListResourceGrants(ctx, db.ListResourceGrantsParams{
		ResourceType: string(resourceType),
		ResourceID:   resourceID,
	})
	if err != nil {
		return nil, fmt.Errorf("list resource grants: %w", err)
	}

	return s.rowsToGrants(rows), nil
}

// CleanupExpired removes all expired grants.
func (s *Service) CleanupExpired(ctx context.Context) error {
	if err := s.queries.DeleteExpiredResourceGrants(ctx); err != nil {
		return fmt.Errorf("cleanup expired grants: %w", err)
	}

	s.logger.Debug("Expired resource grants cleaned up")

	return nil
}

// rowToGrant converts a database row to a Grant.
func (s *Service) rowToGrant(row *db.ResourceGrant) *Grant {
	if row == nil {
		return nil
	}

	grant := &Grant{
		ID:           row.ID,
		UserID:       row.UserID,
		ResourceType: ResourceType(row.ResourceType),
		ResourceID:   row.ResourceID,
		GrantType:    GrantType(row.GrantType),
		GrantedAt:    row.GrantedAt,
	}

	if row.GrantedBy.Valid {
		id := uuid.UUID(row.GrantedBy.Bytes)
		grant.GrantedBy = &id
	}

	if row.ExpiresAt.Valid {
		grant.ExpiresAt = &row.ExpiresAt.Time
	}

	return grant
}

// rowsToGrants converts multiple database rows to Grants.
func (s *Service) rowsToGrants(rows []db.ResourceGrant) []Grant {
	grants := make([]Grant, len(rows))
	for i := range rows {
		if g := s.rowToGrant(&rows[i]); g != nil {
			grants[i] = *g
		}
	}
	return grants
}
