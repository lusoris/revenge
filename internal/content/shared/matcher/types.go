// Package matcher provides a generic matching framework for content libraries.
// It supports different content types (movies, TV shows, music) through pluggable
// match strategies that implement the MatchStrategy interface.
package matcher

import (
	"context"

	"github.com/lusoris/revenge/internal/content/shared/scanner"
)

// MatchType indicates how content was matched
type MatchType string

const (
	// MatchTypeExact indicates an exact ID match (e.g., TMDb ID)
	MatchTypeExact MatchType = "exact"

	// MatchTypeTitle indicates a title and year match
	MatchTypeTitle MatchType = "title"

	// MatchTypeFuzzy indicates a fuzzy title match with lower confidence
	MatchTypeFuzzy MatchType = "fuzzy"

	// MatchTypeManual indicates manual user matching
	MatchTypeManual MatchType = "manual"

	// MatchTypeUnmatched indicates no match was found
	MatchTypeUnmatched MatchType = "unmatched"
)

// MatchResult represents the result of matching a scanned file to content.
// The generic Content field can hold any content type (Movie, TVShow, etc.).
type MatchResult[T any] struct {
	// ScanResult is the original scan result being matched
	ScanResult scanner.ScanResult

	// Content is the matched content item (nil if unmatched)
	Content *T

	// MatchType indicates how the match was made
	MatchType MatchType

	// Confidence is the match confidence score (0.0-1.0)
	Confidence float64

	// Error contains any error that occurred during matching
	Error error

	// CreatedNew indicates if a new content record was created
	CreatedNew bool
}

// IsMatched returns true if a match was found
func (r MatchResult[T]) IsMatched() bool {
	return r.Content != nil && r.Error == nil
}

// MatchStrategy defines the interface for content-specific matching logic.
// Each content type (movie, TV, music) implements this interface.
type MatchStrategy[T any] interface {
	// FindExisting searches for existing content in the database
	// Returns the content and confidence score, or nil if not found
	FindExisting(ctx context.Context, scanResult scanner.ScanResult) (*T, float64, error)

	// SearchExternal searches external metadata providers (TMDb, TVDB, etc.)
	// Returns a list of potential matches
	SearchExternal(ctx context.Context, scanResult scanner.ScanResult) ([]*T, error)

	// CalculateConfidence calculates a confidence score for a potential match
	CalculateConfidence(scanResult scanner.ScanResult, candidate *T) float64

	// CreateContent creates a new content record from a matched candidate
	CreateContent(ctx context.Context, candidate *T) (*T, error)
}

// Matcher is a generic content matcher that uses a pluggable strategy
type Matcher[T any] struct {
	strategy MatchStrategy[T]
}

// NewMatcher creates a new matcher with the given strategy
func NewMatcher[T any](strategy MatchStrategy[T]) *Matcher[T] {
	return &Matcher[T]{
		strategy: strategy,
	}
}

// MatchFiles matches multiple scan results to content
func (m *Matcher[T]) MatchFiles(ctx context.Context, results []scanner.ScanResult) ([]MatchResult[T], error) {
	matchResults := make([]MatchResult[T], 0, len(results))

	for _, result := range results {
		matchResult := m.MatchFile(ctx, result)
		matchResults = append(matchResults, matchResult)
	}

	return matchResults, nil
}

// MatchFile matches a single scan result to content
func (m *Matcher[T]) MatchFile(ctx context.Context, result scanner.ScanResult) MatchResult[T] {
	// Try to find existing content in DB first
	existing, confidence, err := m.strategy.FindExisting(ctx, result)
	if err == nil && existing != nil && confidence >= 0.8 {
		return MatchResult[T]{
			ScanResult: result,
			Content:    existing,
			MatchType:  MatchTypeTitle,
			Confidence: confidence,
		}
	}

	// Search external providers
	candidates, err := m.strategy.SearchExternal(ctx, result)
	if err != nil {
		return MatchResult[T]{
			ScanResult: result,
			MatchType:  MatchTypeUnmatched,
			Error:      err,
		}
	}

	if len(candidates) == 0 {
		return MatchResult[T]{
			ScanResult: result,
			MatchType:  MatchTypeUnmatched,
			Error:      ErrNoMatches,
		}
	}

	// Use the first candidate (highest relevance from provider)
	candidate := candidates[0]
	confidence = m.strategy.CalculateConfidence(result, candidate)

	// Create new content record
	created, err := m.strategy.CreateContent(ctx, candidate)
	if err != nil {
		return MatchResult[T]{
			ScanResult: result,
			MatchType:  MatchTypeUnmatched,
			Error:      err,
		}
	}

	matchType := MatchTypeTitle
	if confidence < 0.7 {
		matchType = MatchTypeFuzzy
	}

	return MatchResult[T]{
		ScanResult: result,
		Content:    created,
		MatchType:  matchType,
		Confidence: confidence,
		CreatedNew: true,
	}
}
