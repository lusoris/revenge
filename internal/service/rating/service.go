// Package rating provides the content rating service.
package rating

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/lusoris/revenge/internal/domain"
	"github.com/lusoris/revenge/internal/infra/database/db"
)

// Service handles content rating operations.
type Service struct {
	queries *db.Queries
	logger  *slog.Logger
}

// NewService creates a new rating service.
func NewService(queries *db.Queries, logger *slog.Logger) *Service {
	return &Service{
		queries: queries,
		logger:  logger.With(slog.String("service", "rating")),
	}
}

// GetRatingSystemByCode retrieves a rating system by its code.
func (s *Service) GetRatingSystemByCode(ctx context.Context, code string) (*domain.RatingSystem, error) {
	rs, err := s.queries.GetRatingSystemByCode(ctx, code)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get rating system: %w", err)
	}

	return s.ratingSystemToEntity(&rs), nil
}

// ListRatingSystems retrieves all active rating systems.
func (s *Service) ListRatingSystems(ctx context.Context) ([]*domain.RatingSystem, error) {
	rows, err := s.queries.ListRatingSystems(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list rating systems: %w", err)
	}

	systems := make([]*domain.RatingSystem, len(rows))
	for i := range rows {
		systems[i] = s.ratingSystemToEntity(&rows[i])
	}

	return systems, nil
}

// ListRatingSystemsByCountry retrieves rating systems for a specific country.
func (s *Service) ListRatingSystemsByCountry(ctx context.Context, countryCode string) ([]*domain.RatingSystem, error) {
	rows, err := s.queries.ListRatingSystemsByCountry(ctx, []string{countryCode})
	if err != nil {
		return nil, fmt.Errorf("failed to list rating systems by country: %w", err)
	}

	systems := make([]*domain.RatingSystem, len(rows))
	for i := range rows {
		systems[i] = s.ratingSystemToEntity(&rows[i])
	}

	return systems, nil
}

// GetRating retrieves a rating by system and code.
func (s *Service) GetRating(ctx context.Context, systemID uuid.UUID, code string) (*domain.Rating, error) {
	row, err := s.queries.GetRatingBySystemAndCode(ctx, db.GetRatingBySystemAndCodeParams{
		SystemID: systemID,
		Code:     code,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get rating: %w", err)
	}

	return s.ratingRowToRating(&row), nil
}

// ListRatingsBySystem retrieves all ratings for a rating system.
func (s *Service) ListRatingsBySystem(ctx context.Context, systemID uuid.UUID) ([]*domain.Rating, error) {
	rows, err := s.queries.ListRatingsBySystem(ctx, systemID)
	if err != nil {
		return nil, fmt.Errorf("failed to list ratings: %w", err)
	}

	ratings := make([]*domain.Rating, len(rows))
	for i := range rows {
		ratings[i] = s.listRatingRowToRating(&rows[i])
	}

	return ratings, nil
}

// GetEquivalentRatings retrieves equivalent ratings from other systems.
func (s *Service) GetEquivalentRatings(ctx context.Context, ratingID uuid.UUID) ([]*domain.Rating, error) {
	rows, err := s.queries.GetRatingEquivalents(ctx, ratingID)
	if err != nil {
		return nil, fmt.Errorf("failed to get equivalent ratings: %w", err)
	}

	ratings := make([]*domain.Rating, len(rows))
	for i := range rows {
		ratings[i] = s.equivalentRowToRating(&rows[i])
	}

	return ratings, nil
}

// GetContentRatings retrieves all ratings for a piece of content.
func (s *Service) GetContentRatings(ctx context.Context, contentID uuid.UUID, contentType string) ([]*domain.ContentRating, error) {
	rows, err := s.queries.GetContentRatings(ctx, db.GetContentRatingsParams{
		ContentID:   contentID,
		ContentType: contentType,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get content ratings: %w", err)
	}

	contentRatings := make([]*domain.ContentRating, len(rows))
	for i := range rows {
		contentRatings[i] = s.contentRatingRowToContentRating(&rows[i])
	}

	return contentRatings, nil
}

// GetContentMinLevel retrieves the minimum (most restrictive) rating level for content.
func (s *Service) GetContentMinLevel(ctx context.Context, contentID uuid.UUID, contentType string) (*domain.ContentMinRatingLevel, error) {
	row, err := s.queries.GetContentMinLevel(ctx, db.GetContentMinLevelParams{
		ContentID:   contentID,
		ContentType: contentType,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// Content has no ratings, return nil (unrestricted)
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get content min level: %w", err)
	}

	// MinLevel is interface{}, need to type assert
	minLevel := 0
	if row.MinLevel != nil {
		switch v := row.MinLevel.(type) {
		case int64:
			minLevel = int(v)
		case int32:
			minLevel = int(v)
		case float64:
			minLevel = int(v)
		}
	}

	return &domain.ContentMinRatingLevel{
		ContentID:   row.ContentID,
		ContentType: row.ContentType,
		MinLevel:    minLevel,
		IsAdult:     row.IsAdult,
	}, nil
}

// GetDisplayRating retrieves the rating to display for content in a preferred system.
func (s *Service) GetDisplayRating(ctx context.Context, contentID uuid.UUID, contentType string, preferredSystem string) (*domain.ContentRating, error) {
	row, err := s.queries.GetContentDisplayRating(ctx, db.GetContentDisplayRatingParams{
		ContentID:   contentID,
		ContentType: contentType,
		Code:        preferredSystem,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // No rating found
		}
		return nil, fmt.Errorf("failed to get display rating: %w", err)
	}

	return s.displayRatingRowToContentRating(&row), nil
}

// AddContentRating adds a rating to content.
func (s *Service) AddContentRating(ctx context.Context, contentID uuid.UUID, contentType string, ratingID uuid.UUID, source *string) (*domain.ContentRating, error) {
	var sourceText pgtype.Text
	if source != nil {
		sourceText = pgtype.Text{String: *source, Valid: true}
	}

	row, err := s.queries.CreateContentRating(ctx, db.CreateContentRatingParams{
		ContentID:   contentID,
		ContentType: contentType,
		RatingID:    ratingID,
		Source:      sourceText,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to add content rating: %w", err)
	}

	s.logger.Info("content rating added",
		slog.String("content_id", contentID.String()),
		slog.String("content_type", contentType),
		slog.String("rating_id", ratingID.String()),
	)

	var resultSource *string
	if row.Source.Valid {
		resultSource = &row.Source.String
	}

	return &domain.ContentRating{
		ID:          row.ID,
		ContentID:   row.ContentID,
		ContentType: row.ContentType,
		RatingID:    row.RatingID,
		Source:      resultSource,
		CreatedAt:   row.CreatedAt,
	}, nil
}

// RemoveContentRating removes a rating from content.
func (s *Service) RemoveContentRating(ctx context.Context, contentID uuid.UUID, ratingID uuid.UUID) error {
	err := s.queries.DeleteContentRating(ctx, db.DeleteContentRatingParams{
		ContentID: contentID,
		RatingID:  ratingID,
	})
	if err != nil {
		return fmt.Errorf("failed to remove content rating: %w", err)
	}

	s.logger.Info("content rating removed",
		slog.String("content_id", contentID.String()),
		slog.String("rating_id", ratingID.String()),
	)

	return nil
}

// IsContentAllowed checks if content is allowed for a user's rating level.
func (s *Service) IsContentAllowed(ctx context.Context, contentID uuid.UUID, contentType string, maxLevel int, includeAdult bool) (bool, error) {
	allowed, err := s.queries.IsContentAllowed(ctx, db.IsContentAllowedParams{
		ContentID:   contentID,
		ContentType: contentType,
		MinLevel:    int32(maxLevel),
		Column4:     includeAdult,
	})
	if err != nil {
		return false, fmt.Errorf("failed to check content allowed: %w", err)
	}

	return allowed, nil
}

// FilterAllowedContent filters a list of content IDs to only those allowed.
func (s *Service) FilterAllowedContent(ctx context.Context, contentIDs []uuid.UUID, contentType string, maxLevel int, includeAdult bool) ([]uuid.UUID, error) {
	allowed, err := s.queries.FilterAllowedContentIDs(ctx, db.FilterAllowedContentIDsParams{
		Column1:     contentIDs,
		ContentType: contentType,
		MinLevel:    int32(maxLevel),
		Column4:     includeAdult,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to filter allowed content: %w", err)
	}

	return allowed, nil
}

// Helper conversion functions

func (s *Service) ratingSystemToEntity(rs *db.RatingSystem) *domain.RatingSystem {
	return &domain.RatingSystem{
		ID:           rs.ID,
		Code:         rs.Code,
		Name:         rs.Name,
		CountryCodes: rs.CountryCodes,
		IsActive:     rs.IsActive,
		SortOrder:    int(rs.SortOrder),
		CreatedAt:    rs.CreatedAt,
	}
}

func (s *Service) ratingRowToRating(row *db.GetRatingBySystemAndCodeRow) *domain.Rating {
	var desc *string
	if row.Description.Valid {
		desc = &row.Description.String
	}
	var minAge *int
	if row.MinAge.Valid {
		age := int(row.MinAge.Int32)
		minAge = &age
	}
	var iconURL *string
	if row.IconUrl.Valid {
		iconURL = &row.IconUrl.String
	}

	return &domain.Rating{
		ID:              row.ID,
		SystemID:        row.SystemID,
		Code:            row.Code,
		Name:            row.Name,
		Description:     desc,
		MinAge:          minAge,
		NormalizedLevel: int(row.NormalizedLevel),
		SortOrder:       int(row.SortOrder),
		IsAdult:         row.IsAdult,
		IconURL:         iconURL,
		CreatedAt:       row.CreatedAt,
		System: &domain.RatingSystem{
			Code: row.SystemCode,
			Name: row.SystemName,
		},
	}
}

func (s *Service) listRatingRowToRating(row *db.ListRatingsBySystemRow) *domain.Rating {
	var desc *string
	if row.Description.Valid {
		desc = &row.Description.String
	}
	var minAge *int
	if row.MinAge.Valid {
		age := int(row.MinAge.Int32)
		minAge = &age
	}
	var iconURL *string
	if row.IconUrl.Valid {
		iconURL = &row.IconUrl.String
	}

	return &domain.Rating{
		ID:              row.ID,
		SystemID:        row.SystemID,
		Code:            row.Code,
		Name:            row.Name,
		Description:     desc,
		MinAge:          minAge,
		NormalizedLevel: int(row.NormalizedLevel),
		SortOrder:       int(row.SortOrder),
		IsAdult:         row.IsAdult,
		IconURL:         iconURL,
		CreatedAt:       row.CreatedAt,
		System: &domain.RatingSystem{
			Code: row.SystemCode,
			Name: row.SystemName,
		},
	}
}

func (s *Service) equivalentRowToRating(row *db.GetRatingEquivalentsRow) *domain.Rating {
	var desc *string
	if row.Description.Valid {
		desc = &row.Description.String
	}
	var minAge *int
	if row.MinAge.Valid {
		age := int(row.MinAge.Int32)
		minAge = &age
	}
	var iconURL *string
	if row.IconUrl.Valid {
		iconURL = &row.IconUrl.String
	}

	return &domain.Rating{
		ID:              row.ID,
		SystemID:        row.SystemID,
		Code:            row.Code,
		Name:            row.Name,
		Description:     desc,
		MinAge:          minAge,
		NormalizedLevel: int(row.NormalizedLevel),
		SortOrder:       int(row.SortOrder),
		IsAdult:         row.IsAdult,
		IconURL:         iconURL,
		CreatedAt:       row.CreatedAt,
		System: &domain.RatingSystem{
			Code: row.SystemCode,
			Name: row.SystemName,
		},
	}
}

func (s *Service) contentRatingRowToContentRating(row *db.GetContentRatingsRow) *domain.ContentRating {
	var source *string
	if row.Source.Valid {
		source = &row.Source.String
	}

	return &domain.ContentRating{
		ID:          row.ID,
		ContentID:   row.ContentID,
		ContentType: row.ContentType,
		RatingID:    row.RatingID,
		Source:      source,
		CreatedAt:   row.CreatedAt,
		Rating: &domain.Rating{
			Code:            row.RatingCode,
			Name:            row.RatingName,
			NormalizedLevel: int(row.NormalizedLevel),
			IsAdult:         row.IsAdult,
			System: &domain.RatingSystem{
				Code: row.SystemCode,
				Name: row.SystemName,
			},
		},
	}
}

func (s *Service) displayRatingRowToContentRating(row *db.GetContentDisplayRatingRow) *domain.ContentRating {
	var source *string
	if row.Source.Valid {
		source = &row.Source.String
	}

	return &domain.ContentRating{
		ID:          row.ID,
		ContentID:   row.ContentID,
		ContentType: row.ContentType,
		RatingID:    row.RatingID,
		Source:      source,
		CreatedAt:   row.CreatedAt,
		Rating: &domain.Rating{
			Code:            row.RatingCode,
			Name:            row.RatingName,
			NormalizedLevel: int(row.NormalizedLevel),
			IsAdult:         row.IsAdult,
			System: &domain.RatingSystem{
				Code: row.SystemCode,
				Name: row.SystemName,
			},
		},
	}
}
