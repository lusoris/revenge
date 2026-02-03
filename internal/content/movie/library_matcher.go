package movie

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

// Matcher handles matching scanned files to movies
type Matcher struct {
	repo            Repository
	metadataService *MetadataService
}

// MatchResult represents the result of matching a file
type MatchResult struct {
	ScanResult      ScanResult
	Movie           *Movie
	MatchType       MatchType
	Confidence      float64
	Error           error
	CreatedNewMovie bool
}

// MatchType indicates how a file was matched
type MatchType string

const (
	MatchTypeExact     MatchType = "exact"     // Exact TMDb ID match
	MatchTypeTitle     MatchType = "title"     // Title and year match
	MatchTypeFuzzy     MatchType = "fuzzy"     // Fuzzy title match
	MatchTypeManual    MatchType = "manual"    // Manually matched
	MatchTypeUnmatched MatchType = "unmatched" // Could not match
)

// NewMatcher creates a new file matcher
func NewMatcher(repo Repository, metadataService *MetadataService) *Matcher {
	return &Matcher{
		repo:            repo,
		metadataService: metadataService,
	}
}

// MatchFiles attempts to match scan results to movies
func (m *Matcher) MatchFiles(ctx context.Context, results []ScanResult) ([]MatchResult, error) {
	var matchResults []MatchResult

	for _, result := range results {
		matchResult := m.matchFile(ctx, result)
		matchResults = append(matchResults, matchResult)
	}

	return matchResults, nil
}

// matchFile attempts to match a single file
func (m *Matcher) matchFile(ctx context.Context, result ScanResult) MatchResult {
	// Try to find existing movie in DB first
	existingMovie, err := m.findExistingMovie(ctx, result)
	if err == nil && existingMovie != nil {
		return MatchResult{
			ScanResult: result,
			Movie:      existingMovie,
			MatchType:  MatchTypeTitle,
			Confidence: 0.95,
		}
	}

	// Search TMDb for the movie
	tmdbMovies, err := m.metadataService.SearchMovies(ctx, result.ParsedTitle, result.ParsedYear)
	if err != nil {
		return MatchResult{
			ScanResult: result,
			MatchType:  MatchTypeUnmatched,
			Error:      fmt.Errorf("TMDb search failed: %w", err),
		}
	}

	if len(tmdbMovies) == 0 {
		return MatchResult{
			ScanResult: result,
			MatchType:  MatchTypeUnmatched,
			Error:      fmt.Errorf("no TMDb results for: %s", result.ParsedTitle),
		}
	}

	// Use the first result (highest TMDb relevance)
	tmdbMovie := tmdbMovies[0]

	// Calculate confidence
	confidence := m.calculateConfidence(result, tmdbMovie)

	// Create new movie in DB
	newMovie, err := m.createMovieFromTMDb(ctx, tmdbMovie)
	if err != nil {
		return MatchResult{
			ScanResult: result,
			MatchType:  MatchTypeUnmatched,
			Error:      fmt.Errorf("failed to create movie: %w", err),
		}
	}

	return MatchResult{
		ScanResult:      result,
		Movie:           newMovie,
		MatchType:       MatchTypeTitle,
		Confidence:      confidence,
		CreatedNewMovie: true,
	}
}

// findExistingMovie searches for an existing movie in the database
func (m *Matcher) findExistingMovie(ctx context.Context, result ScanResult) (*Movie, error) {
	// Note: We'd need to add SearchMoviesByTitle to the repository
	// For now, we'll return not found to trigger TMDb search
	return nil, fmt.Errorf("not implemented")
}

// createMovieFromTMDb creates a new movie record from TMDb data
func (m *Matcher) createMovieFromTMDb(ctx context.Context, tmdbMovie *Movie) (*Movie, error) {
	// Enrich with full metadata if we only have search result
	if tmdbMovie.TMDbID != nil {
		if err := m.metadataService.EnrichMovie(ctx, tmdbMovie); err != nil {
			// Log warning but continue with partial data
			_ = err
		}
	}

	// Save to database
	params := CreateMovieParams{
		TMDbID:        tmdbMovie.TMDbID,
		IMDbID:        tmdbMovie.IMDbID,
		Title:         tmdbMovie.Title,
		OriginalTitle: tmdbMovie.OriginalTitle,
		Year:          extractYear(tmdbMovie.ReleaseDate),
		ReleaseDate:   formatDate(tmdbMovie.ReleaseDate),
		Runtime:       tmdbMovie.Runtime,
		Overview:      tmdbMovie.Overview,
		Tagline:       tmdbMovie.Tagline,
		Status:        tmdbMovie.Status,
		PosterPath:    tmdbMovie.PosterPath,
		BackdropPath:  tmdbMovie.BackdropPath,
		VoteAverage:   formatDecimal(tmdbMovie.VoteAverage),
		VoteCount:     tmdbMovie.VoteCount,
		Popularity:    formatDecimal(tmdbMovie.Popularity),
	}

	newMovie, err := m.repo.CreateMovie(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create movie: %w", err)
	}

	// Fetch and save credits if available
	if tmdbMovie.TMDbID != nil {
		credits, err := m.metadataService.GetMovieCredits(ctx, newMovie.ID, int(*tmdbMovie.TMDbID))
		if err == nil && len(credits) > 0 {
			// Save credits to DB (would need CreateMovieCredits method)
			// For now, skip credits
			_ = credits
		}
	}

	// Fetch and save genres if available
	if tmdbMovie.TMDbID != nil {
		genres, err := m.metadataService.GetMovieGenres(ctx, newMovie.ID, int(*tmdbMovie.TMDbID))
		if err == nil && len(genres) > 0 {
			// Save genres to DB (would need CreateMovieGenres method)
			_ = genres
		}
	}

	return newMovie, nil
}

// calculateConfidence calculates match confidence score
func (m *Matcher) calculateConfidence(result ScanResult, tmdbMovie *Movie) float64 {
	confidence := 0.0

	// Title similarity
	parsedTitleLower := strings.ToLower(result.ParsedTitle)
	tmdbTitleLower := strings.ToLower(tmdbMovie.Title)

	if parsedTitleLower == tmdbTitleLower {
		confidence += 0.6
	} else if strings.Contains(tmdbTitleLower, parsedTitleLower) || strings.Contains(parsedTitleLower, tmdbTitleLower) {
		confidence += 0.4
	} else {
		confidence += 0.2
	}

	// Year match
	if result.ParsedYear != nil && tmdbMovie.ReleaseDate != nil {
		tmdbYear := tmdbMovie.ReleaseDate.Year()
		if *result.ParsedYear == tmdbYear {
			confidence += 0.3
		} else if abs(*result.ParsedYear-tmdbYear) <= 1 {
			confidence += 0.1
		}
	} else if result.ParsedYear == nil && tmdbMovie.ReleaseDate != nil {
		// Penalize slightly if year not in filename
		confidence -= 0.05
	}

	// Popularity boost (higher popularity = more likely correct)
	if tmdbMovie.Popularity != nil && !tmdbMovie.Popularity.IsZero() {
		pop, _ := tmdbMovie.Popularity.Float64()
		if pop > 50 {
			confidence += 0.1
		}
	}

	// Ensure confidence is between 0 and 1
	if confidence < 0 {
		confidence = 0
	}
	if confidence > 1 {
		confidence = 1
	}

	return confidence
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func extractYear(t *time.Time) *int32 {
	if t == nil {
		return nil
	}
	year := int32(t.Year())
	return &year
}

func formatDate(t *time.Time) *string {
	if t == nil {
		return nil
	}
	s := t.Format("2006-01-02")
	return &s
}

func formatDecimal(d *decimal.Decimal) *string {
	if d == nil || d.IsZero() {
		return nil
	}
	s := d.String()
	return &s
}
