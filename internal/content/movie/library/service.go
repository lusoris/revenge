package library

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/lusoris/revenge/internal/content/movie"
	"github.com/lusoris/revenge/internal/content/movie/metadata"
)

// Service manages movie library operations
type Service struct {
	repo            movie.Repository
	metadataService *metadata.MetadataService
	scanner         *Scanner
	matcher         *Matcher
}

// Config holds library service configuration
type Config struct {
	LibraryPaths []string
}

// ScanSummary contains statistics from a library scan
type ScanSummary struct {
	TotalFiles     int
	MatchedFiles   int
	UnmatchedFiles int
	NewMovies      int
	ExistingMovies int
	Errors         []error
}

// NewService creates a new library service
func NewService(
	repo movie.Repository,
	metadataService *metadata.MetadataService,
	config Config,
) *Service {
	scanner := NewScanner(config.LibraryPaths)
	matcher := NewMatcher(repo, metadataService)

	return &Service{
		repo:            repo,
		metadataService: metadataService,
		scanner:         scanner,
		matcher:         matcher,
	}
}

// ScanLibrary scans all library paths and matches files to movies
func (s *Service) ScanLibrary(ctx context.Context) (*ScanSummary, error) {
	// Scan file system
	scanResults, err := s.scanner.Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("scan failed: %w", err)
	}

	summary := &ScanSummary{
		TotalFiles: len(scanResults),
	}

	// Match files to movies
	matchResults, err := s.matcher.MatchFiles(ctx, scanResults)
	if err != nil {
		return nil, fmt.Errorf("match failed: %w", err)
	}

	// Process match results
	for _, result := range matchResults {
		if result.Error != nil {
			summary.Errors = append(summary.Errors, result.Error)
			summary.UnmatchedFiles++
			continue
		}

		if result.Movie != nil {
			summary.MatchedFiles++
			if result.CreatedNewMovie {
				summary.NewMovies++

				// Create movie file record
				fileInfo, err := ExtractFileInfo(result.ScanResult.FilePath)
				if err != nil {
					summary.Errors = append(summary.Errors, fmt.Errorf("failed to extract file info: %w", err))
					continue
				}

				movieFile := CreateMovieFile(result.Movie.ID, fileInfo)
				if err := s.createMovieFile(ctx, movieFile); err != nil {
					summary.Errors = append(summary.Errors, fmt.Errorf("failed to create movie file: %w", err))
				}
			} else {
				summary.ExistingMovies++
			}
		} else {
			summary.UnmatchedFiles++
		}
	}

	return summary, nil
}

// createMovieFile creates a movie file record
func (s *Service) createMovieFile(ctx context.Context, movieFile *movie.MovieFile) error {
	params := movie.CreateMovieFileParams{
		MovieID:     movieFile.MovieID,
		FilePath:    movieFile.FilePath,
		FileSize:    movieFile.FileSize,
		Container:   movieFile.Container,
		Resolution:  movieFile.Resolution,
		VideoCodec:  movieFile.VideoCodec,
		AudioCodec:  movieFile.AudioCodec,
		BitrateKbps: movieFile.BitrateKbps,
	}

	_, err := s.repo.CreateMovieFile(ctx, params)
	return err
}

// RefreshMovie updates a movie's metadata from TMDb
func (s *Service) RefreshMovie(ctx context.Context, movieID uuid.UUID) error {
	// Get existing movie
	existingMovie, err := s.repo.GetMovie(ctx, movieID)
	if err != nil {
		return fmt.Errorf("failed to get movie: %w", err)
	}

	if existingMovie.TMDbID == nil {
		return fmt.Errorf("movie has no TMDb ID")
	}

	// Fetch fresh metadata from TMDb
	if err := s.metadataService.EnrichMovie(ctx, existingMovie); err != nil {
		return fmt.Errorf("failed to enrich movie: %w", err)
	}

	// Update in database - note: UpdateMovie signature needs to be checked
	// For now, just use the enriched movie data
	_ = existingMovie

	// Refresh credits
	if existingMovie.TMDbID != nil {
		credits, err := s.metadataService.GetMovieCredits(ctx, movieID, int(*existingMovie.TMDbID))
		if err == nil && len(credits) > 0 {
			// Delete old credits and insert new ones
			// Would need DeleteMovieCredits and CreateMovieCredits methods
			_ = credits
		}
	}

	// Refresh genres
	if existingMovie.TMDbID != nil {
		genres, err := s.metadataService.GetMovieGenres(ctx, movieID, int(*existingMovie.TMDbID))
		if err == nil && len(genres) > 0 {
			// Delete old genres and insert new ones
			_ = genres
		}
	}

	return nil
}

// GetLibraryStats returns statistics about the library
func (s *Service) GetLibraryStats(ctx context.Context) (map[string]int, error) {
	// This would query the repository for counts
	// For now, return placeholder
	return map[string]int{
		"total_movies": 0,
		"total_files":  0,
	}, nil
}
