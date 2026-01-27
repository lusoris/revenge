// Package genre provides the genre service implementation.
package genre

import (
	"context"
	"fmt"
	"log/slog"
	"regexp"
	"strings"

	"github.com/google/uuid"

	"github.com/lusoris/revenge/internal/domain"
)

// Service implements domain.GenreService.
type Service struct {
	repo   domain.GenreRepository
	logger *slog.Logger
}

// NewService creates a new genre service.
func NewService(repo domain.GenreRepository, logger *slog.Logger) *Service {
	return &Service{
		repo:   repo,
		logger: logger.With(slog.String("service", "genre")),
	}
}

// GetGenre retrieves a genre by ID.
func (s *Service) GetGenre(ctx context.Context, id uuid.UUID) (*domain.Genre, error) {
	genre, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get genre: %w", err)
	}
	return genre, nil
}

// GetGenreBySlug retrieves a genre by domain and slug.
func (s *Service) GetGenreBySlug(ctx context.Context, d domain.GenreDomain, slug string) (*domain.Genre, error) {
	genre, err := s.repo.GetBySlug(ctx, d, slug)
	if err != nil {
		return nil, fmt.Errorf("get genre by slug: %w", err)
	}
	return genre, nil
}

// ListGenres retrieves genres with filtering options.
func (s *Service) ListGenres(ctx context.Context, params domain.ListGenresParams) ([]*domain.Genre, error) {
	genres, err := s.repo.List(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("list genres: %w", err)
	}
	return genres, nil
}

// ListGenresByDomain retrieves all genres for a specific domain.
func (s *Service) ListGenresByDomain(ctx context.Context, d domain.GenreDomain) ([]*domain.Genre, error) {
	if !d.IsValid() {
		return nil, fmt.Errorf("invalid genre domain: %s", d)
	}

	genres, err := s.repo.ListByDomain(ctx, d)
	if err != nil {
		return nil, fmt.Errorf("list genres by domain: %w", err)
	}
	return genres, nil
}

// ListGenresForMediaItem retrieves all genres assigned to a media item.
func (s *Service) ListGenresForMediaItem(ctx context.Context, mediaItemID uuid.UUID) ([]*domain.Genre, error) {
	genres, err := s.repo.ListForMediaItem(ctx, mediaItemID)
	if err != nil {
		return nil, fmt.Errorf("list genres for media item: %w", err)
	}
	return genres, nil
}

// GetGenreHierarchy retrieves the genre hierarchy for a domain.
// Returns top-level genres with their children populated.
func (s *Service) GetGenreHierarchy(ctx context.Context, d domain.GenreDomain) ([]*domain.Genre, error) {
	// Get all genres for the domain
	allGenres, err := s.repo.ListByDomain(ctx, d)
	if err != nil {
		return nil, fmt.Errorf("list genres: %w", err)
	}

	// Build hierarchy
	genreMap := make(map[uuid.UUID]*domain.Genre)
	var topLevel []*domain.Genre

	// First pass: index all genres
	for _, g := range allGenres {
		g.Children = make([]*domain.Genre, 0)
		genreMap[g.ID] = g
	}

	// Second pass: build hierarchy
	for _, g := range allGenres {
		if g.ParentID == nil {
			topLevel = append(topLevel, g)
		} else if parent, ok := genreMap[*g.ParentID]; ok {
			parent.Children = append(parent.Children, g)
			g.Parent = parent
		}
	}

	return topLevel, nil
}

// SearchGenres searches for genres by name within a domain.
func (s *Service) SearchGenres(ctx context.Context, d domain.GenreDomain, query string, limit int) ([]*domain.Genre, error) {
	if !d.IsValid() {
		return nil, fmt.Errorf("invalid genre domain: %s", d)
	}
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	genres, err := s.repo.Search(ctx, d, query, limit)
	if err != nil {
		return nil, fmt.Errorf("search genres: %w", err)
	}
	return genres, nil
}

// CreateGenre creates a new genre.
func (s *Service) CreateGenre(ctx context.Context, params domain.CreateGenreParams) (*domain.Genre, error) {
	// Validate domain
	if !params.Domain.IsValid() {
		return nil, fmt.Errorf("invalid genre domain: %s", params.Domain)
	}

	// Validate name
	if strings.TrimSpace(params.Name) == "" {
		return nil, fmt.Errorf("genre name is required")
	}

	// Generate slug if not provided
	if params.Slug == "" {
		params.Slug = generateSlug(params.Name)
	}

	// Validate parent exists if provided
	if params.ParentID != nil {
		parent, err := s.repo.GetByID(ctx, *params.ParentID)
		if err != nil {
			return nil, fmt.Errorf("parent genre not found: %w", err)
		}
		// Parent must be in the same domain
		if parent.Domain != params.Domain {
			return nil, fmt.Errorf("parent genre must be in the same domain")
		}
	}

	genre, err := s.repo.Create(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("create genre: %w", err)
	}

	s.logger.Info("genre created",
		slog.String("id", genre.ID.String()),
		slog.String("domain", string(genre.Domain)),
		slog.String("name", genre.Name),
		slog.String("slug", genre.Slug))

	return genre, nil
}

// UpdateGenre updates an existing genre.
func (s *Service) UpdateGenre(ctx context.Context, params domain.UpdateGenreParams) (*domain.Genre, error) {
	// Verify genre exists
	existing, err := s.repo.GetByID(ctx, params.ID)
	if err != nil {
		return nil, fmt.Errorf("genre not found: %w", err)
	}

	// Validate parent if being changed
	if params.ParentID != nil {
		// Cannot be its own parent
		if *params.ParentID == params.ID {
			return nil, fmt.Errorf("genre cannot be its own parent")
		}

		parent, err := s.repo.GetByID(ctx, *params.ParentID)
		if err != nil {
			return nil, fmt.Errorf("parent genre not found: %w", err)
		}
		// Parent must be in the same domain
		if parent.Domain != existing.Domain {
			return nil, fmt.Errorf("parent genre must be in the same domain")
		}
	}

	genre, err := s.repo.Update(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("update genre: %w", err)
	}

	s.logger.Info("genre updated",
		slog.String("id", genre.ID.String()),
		slog.String("name", genre.Name))

	return genre, nil
}

// DeleteGenre deletes a genre.
func (s *Service) DeleteGenre(ctx context.Context, id uuid.UUID) error {
	// Verify genre exists
	genre, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("genre not found: %w", err)
	}

	// Check for children
	children, err := s.repo.ListChildren(ctx, id)
	if err != nil {
		return fmt.Errorf("check children: %w", err)
	}
	if len(children) > 0 {
		return fmt.Errorf("cannot delete genre with %d children", len(children))
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete genre: %w", err)
	}

	s.logger.Info("genre deleted",
		slog.String("id", id.String()),
		slog.String("name", genre.Name))

	return nil
}

// AssignGenreToMediaItem assigns a genre to a media item.
func (s *Service) AssignGenreToMediaItem(ctx context.Context, mediaItemID, genreID uuid.UUID, source string) error {
	// Validate source
	if source == "" {
		source = domain.GenreSourceManual
	}

	params := domain.AssignGenreParams{
		MediaItemID: mediaItemID,
		GenreID:     genreID,
		Source:      source,
		Confidence:  1.0,
	}

	if err := s.repo.AssignToMediaItem(ctx, params); err != nil {
		return fmt.Errorf("assign genre: %w", err)
	}

	s.logger.Debug("genre assigned to media item",
		slog.String("media_item_id", mediaItemID.String()),
		slog.String("genre_id", genreID.String()),
		slog.String("source", source))

	return nil
}

// RemoveGenreFromMediaItem removes a genre from a media item.
func (s *Service) RemoveGenreFromMediaItem(ctx context.Context, mediaItemID, genreID uuid.UUID) error {
	if err := s.repo.RemoveFromMediaItem(ctx, mediaItemID, genreID); err != nil {
		return fmt.Errorf("remove genre: %w", err)
	}

	s.logger.Debug("genre removed from media item",
		slog.String("media_item_id", mediaItemID.String()),
		slog.String("genre_id", genreID.String()))

	return nil
}

// SetMediaItemGenres replaces all genres on a media item.
func (s *Service) SetMediaItemGenres(ctx context.Context, mediaItemID uuid.UUID, genreIDs []uuid.UUID, source string) error {
	// First remove all existing genres
	if err := s.repo.RemoveAllFromMediaItem(ctx, mediaItemID); err != nil {
		return fmt.Errorf("clear existing genres: %w", err)
	}

	// Then assign new genres
	if len(genreIDs) > 0 {
		if err := s.repo.BulkAssignToMediaItem(ctx, mediaItemID, genreIDs, source); err != nil {
			return fmt.Errorf("assign genres: %w", err)
		}
	}

	s.logger.Debug("media item genres updated",
		slog.String("media_item_id", mediaItemID.String()),
		slog.Int("genre_count", len(genreIDs)))

	return nil
}

// GetDomainForMediaType maps a media type to a genre domain.
func (s *Service) GetDomainForMediaType(mediaType string) domain.GenreDomain {
	switch mediaType {
	case "movie":
		return domain.GenreDomainMovie
	case "series", "season", "episode":
		return domain.GenreDomainTV
	case "artist", "album", "audio":
		return domain.GenreDomainMusic
	case "book", "audiobook":
		return domain.GenreDomainBook
	case "podcast":
		return domain.GenreDomainPodcast
	default:
		return domain.GenreDomainMovie // Default fallback
	}
}

// generateSlug creates a URL-safe slug from a name.
func generateSlug(name string) string {
	// Convert to lowercase
	slug := strings.ToLower(name)

	// Replace spaces and special characters with hyphens
	slug = regexp.MustCompile(`[^a-z0-9]+`).ReplaceAllString(slug, "-")

	// Remove leading/trailing hyphens
	slug = strings.Trim(slug, "-")

	return slug
}
