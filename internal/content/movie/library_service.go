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
		FileName:    movieFile.FileName,
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

	tmdbID := int(*existingMovie.TMDbID)

	// Fetch fresh metadata from TMDb
	if err := s.metadataService.EnrichMovie(ctx, existingMovie); err != nil {
		return fmt.Errorf("failed to enrich movie: %w", err)
	}

	// Update movie in database
	params := UpdateMovieParams{
		ID:               existingMovie.ID,
		TMDbID:           existingMovie.TMDbID,
		IMDbID:           existingMovie.IMDbID,
		Title:            &existingMovie.Title,
		OriginalTitle:    existingMovie.OriginalTitle,
		Year:             existingMovie.Year,
		ReleaseDate:      formatDate(existingMovie.ReleaseDate),
		Runtime:          existingMovie.Runtime,
		Overview:         existingMovie.Overview,
		Tagline:          existingMovie.Tagline,
		Status:           existingMovie.Status,
		OriginalLanguage: existingMovie.OriginalLanguage,
		PosterPath:       existingMovie.PosterPath,
		BackdropPath:     existingMovie.BackdropPath,
		VoteAverage:      formatDecimal(existingMovie.VoteAverage),
		VoteCount:        existingMovie.VoteCount,
		Popularity:       formatDecimal(existingMovie.Popularity),
		Budget:           existingMovie.Budget,
		Revenue:          existingMovie.Revenue,
	}

	if _, err := s.repo.UpdateMovie(ctx, params); err != nil {
		return fmt.Errorf("failed to update movie: %w", err)
	}

	// Refresh credits
	credits, err := s.metadataService.GetMovieCredits(ctx, movieID, tmdbID)
	if err == nil && len(credits) > 0 {
		// Delete old credits
		if err := s.repo.DeleteMovieCredits(ctx, movieID); err != nil {
			return fmt.Errorf("failed to delete old credits: %w", err)
		}

		// Insert new credits
		for _, credit := range credits {
			creditParams := CreateMovieCreditParams{
				MovieID:      movieID,
				TMDbPersonID: credit.TMDbPersonID,
				Name:         credit.Name,
				CreditType:   credit.CreditType,
				Character:    credit.Character,
				Job:          credit.Job,
				Department:   credit.Department,
				CastOrder:    credit.CastOrder,
				ProfilePath:  credit.ProfilePath,
			}
			if _, err := s.repo.CreateMovieCredit(ctx, creditParams); err != nil {
				// Log but continue with other credits
				continue
			}
		}
	}

	// Refresh genres
	genres, err := s.metadataService.GetMovieGenres(ctx, movieID, tmdbID)
	if err == nil && len(genres) > 0 {
		// Delete old genres
		if err := s.repo.DeleteMovieGenres(ctx, movieID); err != nil {
			return fmt.Errorf("failed to delete old genres: %w", err)
		}

		// Add new genres
		for _, genre := range genres {
			if err := s.repo.AddMovieGenre(ctx, movieID, genre.TMDbGenreID, genre.Name); err != nil {
				// Log but continue with other genres
				continue
			}
		}
	}

	return nil
}

// GetLibraryStats returns statistics about the library
func (s *LibraryService) GetLibraryStats(ctx context.Context) (map[string]int, error) {
	stats := map[string]int{
		"total_movies": 0,
	}

	count, err := s.repo.CountMovies(ctx)
	if err != nil {
		return nil, fmt.Errorf("count movies: %w", err)
	}
	stats["total_movies"] = int(count)

	return stats, nil
}

// MatchFile attempts to match a single file to a movie
// This is used by the file match job to process individual files
func (s *LibraryService) MatchFile(ctx context.Context, filePath string, forceRematch bool) (*MatchResult, error) {
	// Check if file exists
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("file not found: %w", err)
	}

	// Check if it's a video file
	if !isVideoFile(filePath) {
		return nil, fmt.Errorf("not a video file: %s", filePath)
	}

	// Check if file is already matched (unless force rematch)
	if !forceRematch {
		existingFile, err := s.repo.GetMovieFileByPath(ctx, filePath)
		if err == nil && existingFile != nil {
			// File already matched
			movie, err := s.repo.GetMovie(ctx, existingFile.MovieID)
			if err == nil {
				return &MatchResult{
					ScanResult: ScanResult{
						FilePath: filePath,
						FileName: fileInfo.Name(),
						FileSize: fileInfo.Size(),
						IsVideo:  true,
					},
					Movie:      movie,
					MatchType:  MatchTypeExact,
					Confidence: 1.0,
				}, nil
			}
		}
	}

	// Parse filename to get title and year
	fileName := fileInfo.Name()
	title, year := parseMovieFilename(fileName)

	// Create scan result
	scanResult := ScanResult{
		FilePath:    filePath,
		FileName:    fileName,
		ParsedTitle: title,
		ParsedYear:  year,
		FileSize:    fileInfo.Size(),
		IsVideo:     true,
	}

	// Use matcher to match the file
	matchResult := s.matcher.MatchFile(ctx, scanResult)

	// If match was successful and created a new movie, create the file record
	if matchResult.Movie != nil && matchResult.Error == nil {
		fileInfoExtracted, err := s.extractFileInfo(filePath)
		if err == nil {
			movieFile := CreateMovieFile(matchResult.Movie.ID, fileInfoExtracted)
			if err := s.createMovieFile(ctx, movieFile); err != nil {
				// Log error but don't fail the match
				matchResult.Error = fmt.Errorf("matched but failed to create file record: %w", err)
			}
		}
	}

	return &matchResult, nil
}
