// Package repository provides PostgreSQL implementations of domain repositories.
package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/jellyfin/jellyfin-go/internal/domain"
	"github.com/jellyfin/jellyfin-go/internal/infra/database/db"
)

// GenreRepository implements domain.GenreRepository using PostgreSQL.
type GenreRepository struct {
	pool    *pgxpool.Pool
	queries *db.Queries
}

// NewGenreRepository creates a new PostgreSQL genre repository.
func NewGenreRepository(pool *pgxpool.Pool) *GenreRepository {
	return &GenreRepository{
		pool:    pool,
		queries: db.New(pool),
	}
}

// GetByID retrieves a genre by its unique ID.
func (r *GenreRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Genre, error) {
	genre, err := r.queries.GetGenreByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get genre by id: %w", err)
	}
	return mapDBGenreToDomain(&genre), nil
}

// GetBySlug retrieves a genre by domain and slug.
func (r *GenreRepository) GetBySlug(ctx context.Context, d domain.GenreDomain, slug string) (*domain.Genre, error) {
	genre, err := r.queries.GetGenreBySlug(ctx, db.GetGenreBySlugParams{
		Domain: db.GenreDomain(d),
		Slug:   slug,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get genre by slug: %w", err)
	}
	return mapDBGenreToDomain(&genre), nil
}

// List retrieves genres with filtering options.
func (r *GenreRepository) List(ctx context.Context, params domain.ListGenresParams) ([]*domain.Genre, error) {
	if params.Domain != nil {
		if params.IncludeAll {
			return r.ListByDomain(ctx, *params.Domain)
		}
		genres, err := r.queries.ListTopLevelGenresByDomain(ctx, db.GenreDomain(*params.Domain))
		if err != nil {
			return nil, fmt.Errorf("failed to list top-level genres: %w", err)
		}
		return mapDBGenresToDomain(genres), nil
	}

	// If no domain specified, return top-level genres from all domains
	// This is a simplified implementation - could be expanded
	return nil, fmt.Errorf("domain filter required")
}

// ListByDomain retrieves all genres for a specific domain.
func (r *GenreRepository) ListByDomain(ctx context.Context, d domain.GenreDomain) ([]*domain.Genre, error) {
	genres, err := r.queries.ListGenresByDomain(ctx, db.GenreDomain(d))
	if err != nil {
		return nil, fmt.Errorf("failed to list genres by domain: %w", err)
	}
	return mapDBGenresToDomain(genres), nil
}

// ListChildren retrieves child genres of a parent.
func (r *GenreRepository) ListChildren(ctx context.Context, parentID uuid.UUID) ([]*domain.Genre, error) {
	genres, err := r.queries.ListChildGenres(ctx, pgtype.UUID{Bytes: parentID, Valid: true})
	if err != nil {
		return nil, fmt.Errorf("failed to list child genres: %w", err)
	}
	return mapDBGenresToDomain(genres), nil
}

// ListForMediaItem retrieves all genres assigned to a media item.
func (r *GenreRepository) ListForMediaItem(ctx context.Context, mediaItemID uuid.UUID) ([]*domain.Genre, error) {
	genres, err := r.queries.ListGenresForMediaItem(ctx, mediaItemID)
	if err != nil {
		return nil, fmt.Errorf("failed to list genres for media item: %w", err)
	}
	return mapDBGenresToDomain(genres), nil
}

// Search searches for genres by name within a domain.
func (r *GenreRepository) Search(ctx context.Context, d domain.GenreDomain, query string, limit int) ([]*domain.Genre, error) {
	genres, err := r.queries.SearchGenres(ctx, db.SearchGenresParams{
		Domain:  db.GenreDomain(d),
		Column2: pgtype.Text{String: query, Valid: true},
		Limit:   int32(min(limit, 1000)), //nolint:gosec // limit is validated upstream
	})
	if err != nil {
		return nil, fmt.Errorf("failed to search genres: %w", err)
	}
	return mapDBGenresToDomain(genres), nil
}

// Create creates a new genre.
func (r *GenreRepository) Create(ctx context.Context, params domain.CreateGenreParams) (*domain.Genre, error) {
	externalIDsJSON, err := json.Marshal(params.ExternalIDs)
	if err != nil {
		externalIDsJSON = []byte("{}")
	}

	dbParams := db.CreateGenreParams{
		Domain:      db.GenreDomain(params.Domain),
		Name:        params.Name,
		Slug:        params.Slug,
		ExternalIds: externalIDsJSON,
	}

	if params.Description != nil {
		dbParams.Description = pgtype.Text{String: *params.Description, Valid: true}
	}
	if params.ParentID != nil {
		dbParams.ParentID = pgtype.UUID{Bytes: *params.ParentID, Valid: true}
	}

	genre, err := r.queries.CreateGenre(ctx, dbParams)
	if err != nil {
		if isUniqueViolation(err, "genres_domain_slug_key") {
			return nil, domain.ErrAlreadyExists
		}
		return nil, fmt.Errorf("failed to create genre: %w", err)
	}
	return mapDBGenreToDomain(&genre), nil
}

// Update updates an existing genre.
func (r *GenreRepository) Update(ctx context.Context, params domain.UpdateGenreParams) (*domain.Genre, error) {
	dbParams := db.UpdateGenreParams{
		ID: params.ID,
	}

	if params.Name != nil {
		dbParams.Name = pgtype.Text{String: *params.Name, Valid: true}
	}
	if params.Slug != nil {
		dbParams.Slug = pgtype.Text{String: *params.Slug, Valid: true}
	}
	if params.Description != nil {
		dbParams.Description = pgtype.Text{String: *params.Description, Valid: true}
	}
	if params.ParentID != nil {
		dbParams.ParentID = pgtype.UUID{Bytes: *params.ParentID, Valid: true}
	}
	if params.ExternalIDs != nil {
		externalIDsJSON, err := json.Marshal(params.ExternalIDs)
		if err == nil {
			dbParams.ExternalIds = externalIDsJSON
		}
	}

	genre, err := r.queries.UpdateGenre(ctx, dbParams)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		if isUniqueViolation(err, "genres_domain_slug_key") {
			return nil, domain.ErrAlreadyExists
		}
		return nil, fmt.Errorf("failed to update genre: %w", err)
	}
	return mapDBGenreToDomain(&genre), nil
}

// Delete deletes a genre by ID.
func (r *GenreRepository) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.queries.DeleteGenre(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete genre: %w", err)
	}
	return nil
}

// AssignToMediaItem assigns a genre to a media item.
func (r *GenreRepository) AssignToMediaItem(ctx context.Context, params domain.AssignGenreParams) error {
	err := r.queries.AssignGenreToMediaItem(ctx, db.AssignGenreToMediaItemParams{
		MediaItemID: params.MediaItemID,
		GenreID:     params.GenreID,
		Source:      params.Source,
		Confidence:  numericFromFloat64(params.Confidence),
	})
	if err != nil {
		return fmt.Errorf("failed to assign genre to media item: %w", err)
	}
	return nil
}

// RemoveFromMediaItem removes a genre from a media item.
func (r *GenreRepository) RemoveFromMediaItem(ctx context.Context, mediaItemID, genreID uuid.UUID) error {
	err := r.queries.RemoveGenreFromMediaItem(ctx, db.RemoveGenreFromMediaItemParams{
		MediaItemID: mediaItemID,
		GenreID:     genreID,
	})
	if err != nil {
		return fmt.Errorf("failed to remove genre from media item: %w", err)
	}
	return nil
}

// RemoveAllFromMediaItem removes all genres from a media item.
func (r *GenreRepository) RemoveAllFromMediaItem(ctx context.Context, mediaItemID uuid.UUID) error {
	err := r.queries.RemoveAllGenresFromMediaItem(ctx, mediaItemID)
	if err != nil {
		return fmt.Errorf("failed to remove all genres from media item: %w", err)
	}
	return nil
}

// BulkAssignToMediaItem assigns multiple genres to a media item.
func (r *GenreRepository) BulkAssignToMediaItem(ctx context.Context, mediaItemID uuid.UUID, genreIDs []uuid.UUID, source string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	qtx := r.queries.WithTx(tx)

	for _, genreID := range genreIDs {
		err := qtx.AssignGenreToMediaItem(ctx, db.AssignGenreToMediaItemParams{
			MediaItemID: mediaItemID,
			GenreID:     genreID,
			Source:      source,
			Confidence:  numericFromFloat64(1.0),
		})
		if err != nil {
			return fmt.Errorf("failed to assign genre %s: %w", genreID, err)
		}
	}

	return tx.Commit(ctx)
}

// mapDBGenreToDomain converts a database genre to a domain genre.
func mapDBGenreToDomain(g *db.Genre) *domain.Genre {
	if g == nil {
		return nil
	}

	genre := &domain.Genre{
		ID:        g.ID,
		Domain:    domain.GenreDomain(g.Domain),
		Name:      g.Name,
		Slug:      g.Slug,
		CreatedAt: g.CreatedAt,
		UpdatedAt: g.UpdatedAt,
	}

	if g.Description.Valid {
		genre.Description = &g.Description.String
	}
	if g.ParentID.Valid {
		parentID := g.ParentID.Bytes
		genre.ParentID = (*uuid.UUID)(&parentID)
	}

	// Parse external IDs
	if len(g.ExternalIds) > 0 {
		var externalIDs map[string]string
		if err := json.Unmarshal(g.ExternalIds, &externalIDs); err == nil {
			genre.ExternalIDs = externalIDs
		}
	}
	if genre.ExternalIDs == nil {
		genre.ExternalIDs = make(map[string]string)
	}

	return genre
}

// mapDBGenresToDomain converts a slice of database genres to domain genres.
func mapDBGenresToDomain[T db.Genre](genres []T) []*domain.Genre {
	result := make([]*domain.Genre, len(genres))
	for i := range genres {
		g := db.Genre(genres[i])
		result[i] = mapDBGenreToDomain(&g)
	}
	return result
}

// numericFromFloat64 creates a pgtype.Numeric from a float64.
func numericFromFloat64(f float64) pgtype.Numeric {
	var n pgtype.Numeric
	_ = n.Scan(fmt.Sprintf("%.2f", f))
	return n
}
