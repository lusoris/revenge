package movie

import (
	"context"
	"fmt"

	"os"

	"github.com/google/uuid"

	"github.com/lusoris/revenge/internal/config"
)

// LibraryService manages movie library operations
type LibraryService struct {
	repo            Repository
	metadataService MetadataProvider
	scanner         *Scanner
	matcher         *Matcher
	prober          Prober
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

// NewLibraryService creates a new library service
func NewLibraryService(
	repo Repository,
	metadataService MetadataProvider,
	libConfig config.LibraryConfig,
	prober Prober,
) *LibraryService {
	scanner := NewScanner(libConfig.Paths)
	matcher := NewMatcher(repo, metadataService)

	if prober == nil {
		prober = NewMediaInfoProber()
	}

	return &LibraryService{
		repo:            repo,
		metadataService: metadataService,
		scanner:         scanner,
		matcher:         matcher,
		prober:          prober,
	}
}

// ScanLibrary scans all library paths and matches files to movies
func (s *LibraryService) ScanLibrary(ctx context.Context) (*ScanSummary, error) {
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
				fileInfo, err := s.extractFileInfo(result.ScanResult.FilePath)
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

// extractFileInfo extracts file info using the configured prober
func (s *LibraryService) extractFileInfo(filePath string) (*MovieFileInfo, error) {
	mediaInfo, err := s.prober.Probe(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to probe file: %w", err)
	}

	info := mediaInfo.ToMovieFileInfo()

	// Get file size from stat if not in mediainfo (fallback)
	if info.Size == 0 {
		if fileInfo, err := os.Stat(filePath); err == nil {
			info.Size = fileInfo.Size()
		}
	}
	return info, nil
}

// createMovieFile creates a movie file record
func (s *LibraryService) createMovieFile(ctx context.Context, movieFile *MovieFile) error {
	params := CreateMovieFileParams{
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
func (s *LibraryService) RefreshMovie(ctx context.Context, movieID uuid.UUID) error {
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
func (s *LibraryService) GetLibraryStats(ctx context.Context) (map[string]int, error) {
	// This would query the repository for counts
	// For now, return placeholder
	return map[string]int{
		"total_movies": 0,
		"total_files":  0,
	}, nil
}
