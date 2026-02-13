package matcher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLevenshteinDistance(t *testing.T) {
	tests := []struct {
		name     string
		s1       string
		s2       string
		expected int
	}{
		{
			name:     "identical strings",
			s1:       "hello",
			s2:       "hello",
			expected: 0,
		},
		{
			name:     "empty strings",
			s1:       "",
			s2:       "",
			expected: 0,
		},
		{
			name:     "one empty string",
			s1:       "hello",
			s2:       "",
			expected: 5,
		},
		{
			name:     "single substitution",
			s1:       "hello",
			s2:       "hallo",
			expected: 1,
		},
		{
			name:     "single insertion",
			s1:       "hello",
			s2:       "hellos",
			expected: 1,
		},
		{
			name:     "single deletion",
			s1:       "hello",
			s2:       "hell",
			expected: 1,
		},
		{
			name:     "multiple edits",
			s1:       "kitten",
			s2:       "sitting",
			expected: 3,
		},
		{
			name:     "case sensitive",
			s1:       "Hello",
			s2:       "hello",
			expected: 1,
		},
		{
			name:     "unicode characters",
			s1:       "café",
			s2:       "cafe",
			expected: 1,
		},
		{
			name:     "movie titles",
			s1:       "the matrix",
			s2:       "matrix",
			expected: 4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := LevenshteinDistance(tt.s1, tt.s2)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNormalizedSimilarity(t *testing.T) {
	tests := []struct {
		name     string
		s1       string
		s2       string
		minScore float64
		maxScore float64
	}{
		{
			name:     "identical strings",
			s1:       "hello",
			s2:       "hello",
			minScore: 1.0,
			maxScore: 1.0,
		},
		{
			name:     "empty strings",
			s1:       "",
			s2:       "",
			minScore: 1.0,
			maxScore: 1.0,
		},
		{
			name:     "completely different",
			s1:       "abc",
			s2:       "xyz",
			minScore: 0.0,
			maxScore: 0.1,
		},
		{
			name:     "similar strings",
			s1:       "hello",
			s2:       "hallo",
			minScore: 0.7,
			maxScore: 0.9,
		},
		{
			name:     "movie title variations",
			s1:       "the matrix",
			s2:       "matrix",
			minScore: 0.5,
			maxScore: 0.7,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NormalizedSimilarity(tt.s1, tt.s2)
			assert.GreaterOrEqual(t, result, tt.minScore)
			assert.LessOrEqual(t, result, tt.maxScore)
		})
	}
}

func TestTitleSimilarity(t *testing.T) {
	tests := []struct {
		name     string
		title1   string
		title2   string
		minScore float64
	}{
		{
			name:     "identical titles",
			title1:   "The Matrix",
			title2:   "The Matrix",
			minScore: 1.0,
		},
		{
			name:     "case difference",
			title1:   "THE MATRIX",
			title2:   "the matrix",
			minScore: 1.0,
		},
		{
			name:     "leading article stripped",
			title1:   "The Matrix",
			title2:   "Matrix",
			minScore: 0.9,
		},
		{
			name:     "punctuation removed",
			title1:   "Spider-Man: No Way Home",
			title2:   "Spider Man No Way Home",
			minScore: 0.9,
		},
		{
			name:     "minor typo",
			title1:   "Inception",
			title2:   "Inceptoin",
			minScore: 0.7,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := TitleSimilarity(tt.title1, tt.title2)
			assert.GreaterOrEqual(t, result, tt.minScore)
		})
	}
}

func TestYearMatch(t *testing.T) {
	tests := []struct {
		name     string
		year1    *int
		year2    *int
		expected float64
	}{
		{
			name:     "exact match",
			year1:    new(2020),
			year2:    new(2020),
			expected: 1.0,
		},
		{
			name:     "one year off",
			year1:    new(2020),
			year2:    new(2021),
			expected: 0.5,
		},
		{
			name:     "two years off",
			year1:    new(2020),
			year2:    new(2022),
			expected: 0.0,
		},
		{
			name:     "first year nil",
			year1:    nil,
			year2:    new(2020),
			expected: 0.0,
		},
		{
			name:     "both nil",
			year1:    nil,
			year2:    nil,
			expected: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := YearMatch(tt.year1, tt.year2)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestConfidenceScore(t *testing.T) {
	t.Run("simple weighted average", func(t *testing.T) {
		score := NewConfidenceScore().
			Add(1.0, 0.5).
			Add(0.0, 0.5).
			Calculate()

		assert.InDelta(t, 0.5, score, 0.01)
	})

	t.Run("weighted towards higher score", func(t *testing.T) {
		score := NewConfidenceScore().
			Add(1.0, 0.8).
			Add(0.0, 0.2).
			Calculate()

		assert.InDelta(t, 0.8, score, 0.01)
	})

	t.Run("with bonus", func(t *testing.T) {
		score := NewConfidenceScore().
			Add(0.5, 1.0).
			AddBonus(0.1).
			Calculate()

		assert.InDelta(t, 0.6, score, 0.01)
	})

	t.Run("with penalty", func(t *testing.T) {
		score := NewConfidenceScore().
			Add(0.5, 1.0).
			AddBonus(-0.1).
			Calculate()

		assert.InDelta(t, 0.4, score, 0.01)
	})

	t.Run("clamped to 1", func(t *testing.T) {
		score := NewConfidenceScore().
			Add(1.0, 1.0).
			AddBonus(0.5).
			Calculate()

		assert.Equal(t, 1.0, score)
	})

	t.Run("clamped to 0", func(t *testing.T) {
		score := NewConfidenceScore().
			Add(0.0, 1.0).
			AddBonus(-0.5).
			Calculate()

		assert.Equal(t, 0.0, score)
	})

	t.Run("empty returns 0", func(t *testing.T) {
		score := NewConfidenceScore().Calculate()
		assert.Equal(t, 0.0, score)
	})
}
