package movie

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/lusoris/revenge/internal/util"
	"github.com/shopspring/decimal"
)

// Matcher handles matching scanned files to movies
type Matcher struct {
	repo            Repository
	metadataService MetadataProvider
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
func NewMatcher(repo Repository, metadataService MetadataProvider) *Matcher {
	return &Matcher{
		repo:            repo,
		metadataService: metadataService,
	}
}

// MatchFiles attempts to match scan results to movies
func (m *Matcher) MatchFiles(ctx context.Context, results []ScanResult) ([]MatchResult, error) {
	var matchResults []MatchResult

	for _, result := range results {
		matchResult := m.MatchFile(ctx, result)
		matchResults = append(matchResults, matchResult)
	}

	return matchResults, nil
}

// MatchFile attempts to match a single file (public for use by jobs)
func (m *Matcher) MatchFile(ctx context.Context, result ScanResult) MatchResult {
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
	if result.ParsedTitle == "" {
		return nil, fmt.Errorf("no title to search")
	}

	// Search for movies with matching title
	movies, err := m.repo.SearchMoviesByTitle(ctx, result.ParsedTitle, 10, 0)
	if err != nil {
		return nil, fmt.Errorf("search failed: %w", err)
	}

	if len(movies) == 0 {
		return nil, fmt.Errorf("no existing movies found")
	}

	// Find best match considering year if available
	var bestMatch *Movie
	var bestScore float64

	for i := range movies {
		movie := &movies[i]
		score := m.scoreExistingMovie(result, movie)

		if score > bestScore {
			bestScore = score
			bestMatch = movie
		}
	}

	// Only return if we have a high-confidence match
	if bestMatch != nil && bestScore >= 0.8 {
		return bestMatch, nil
	}

	return nil, fmt.Errorf("no high-confidence match found")
}

// scoreExistingMovie calculates a match score for an existing movie
func (m *Matcher) scoreExistingMovie(result ScanResult, movie *Movie) float64 {
	score := 0.0

	// Title similarity using Levenshtein distance
	parsedTitle := strings.ToLower(result.ParsedTitle)
	movieTitle := strings.ToLower(movie.Title)

	if parsedTitle == movieTitle {
		score += 0.6
	} else {
		// Calculate normalized Levenshtein distance
		distance := levenshteinDistance(parsedTitle, movieTitle)
		maxLen := max(len(parsedTitle), len(movieTitle))
		if maxLen > 0 {
			similarity := 1.0 - float64(distance)/float64(maxLen)
			score += similarity * 0.6
		}
	}

	// Year match
	if result.ParsedYear != nil && movie.Year != nil {
		if *result.ParsedYear == int(*movie.Year) {
			score += 0.4
		} else if abs(*result.ParsedYear-int(*movie.Year)) <= 1 {
			score += 0.2
		}
	}

	return score
}

// levenshteinDistance calculates the edit distance between two strings
func levenshteinDistance(s1, s2 string) int {
	if len(s1) == 0 {
		return len(s2)
	}
	if len(s2) == 0 {
		return len(s1)
	}

	// Create distance matrix
	matrix := make([][]int, len(s1)+1)
	for i := range matrix {
		matrix[i] = make([]int, len(s2)+1)
		matrix[i][0] = i
	}
	for j := range matrix[0] {
		matrix[0][j] = j
	}

	// Fill in the matrix
	for i := 1; i <= len(s1); i++ {
		for j := 1; j <= len(s2); j++ {
			cost := 1
			if s1[i-1] == s2[j-1] {
				cost = 0
			}
			matrix[i][j] = min(
				matrix[i-1][j]+1,      // deletion
				matrix[i][j-1]+1,      // insertion
				matrix[i-1][j-1]+cost, // substitution
			)
		}
	}

	return matrix[len(s1)][len(s2)]
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
			for _, credit := range credits {
				creditParams := CreateMovieCreditParams{
					MovieID:      newMovie.ID,
					TMDbPersonID: credit.TMDbPersonID,
					Name:         credit.Name,
					CreditType:   credit.CreditType,
					Character:    credit.Character,
					Job:          credit.Job,
					Department:   credit.Department,
					CastOrder:    credit.CastOrder,
					ProfilePath:  credit.ProfilePath,
				}
				// Ignore errors for individual credits
				_, _ = m.repo.CreateMovieCredit(ctx, creditParams)
			}
		}
	}

	// Fetch and save genres if available
	if tmdbMovie.TMDbID != nil {
		genres, err := m.metadataService.GetMovieGenres(ctx, newMovie.ID, int(*tmdbMovie.TMDbID))
		if err == nil && len(genres) > 0 {
			for _, genre := range genres {
				// Ignore errors for individual genres
				_ = m.repo.AddMovieGenre(ctx, newMovie.ID, genre.TMDbGenreID, genre.Name)
			}
		}
	}

	return newMovie, nil
}

// calculateConfidence calculates match confidence score using Levenshtein distance
func (m *Matcher) calculateConfidence(result ScanResult, tmdbMovie *Movie) float64 {
	confidence := 0.0

	// Title similarity using Levenshtein distance
	parsedTitleLower := strings.ToLower(result.ParsedTitle)
	tmdbTitleLower := strings.ToLower(tmdbMovie.Title)

	if parsedTitleLower == tmdbTitleLower {
		confidence += 0.5
	} else {
		// Calculate normalized Levenshtein distance
		distance := levenshteinDistance(parsedTitleLower, tmdbTitleLower)
		maxLen := max(len(parsedTitleLower), len(tmdbTitleLower))
		if maxLen > 0 {
			similarity := 1.0 - float64(distance)/float64(maxLen)
			confidence += similarity * 0.5
		}
	}

	// Also check original title if available
	if tmdbMovie.OriginalTitle != nil && *tmdbMovie.OriginalTitle != tmdbMovie.Title {
		origTitleLower := strings.ToLower(*tmdbMovie.OriginalTitle)
		if parsedTitleLower == origTitleLower {
			confidence += 0.1 // Bonus for matching original title
		}
	}

	// Year match (high weight for exact match)
	if result.ParsedYear != nil && tmdbMovie.ReleaseDate != nil {
		tmdbYear := tmdbMovie.ReleaseDate.Year()
		if *result.ParsedYear == tmdbYear {
			confidence += 0.3
		} else if abs(*result.ParsedYear-tmdbYear) <= 1 {
			confidence += 0.15
		}
	} else if result.ParsedYear == nil && tmdbMovie.ReleaseDate != nil {
		// Penalize slightly if year not in filename
		confidence -= 0.05
	}

	// Popularity boost (higher popularity = more likely correct match)
	if tmdbMovie.Popularity != nil && !tmdbMovie.Popularity.IsZero() {
		pop, _ := tmdbMovie.Popularity.Float64()
		if pop > 100 {
			confidence += 0.1
		} else if pop > 50 {
			confidence += 0.05
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
	year := util.SafeIntToInt32(t.Year())
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
