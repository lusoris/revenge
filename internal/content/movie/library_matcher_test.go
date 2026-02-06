package movie

import (
	"testing"
	"time"

	"github.com/govalues/decimal"
	"github.com/stretchr/testify/assert"
)

func TestCalculateConfidence(t *testing.T) {
	matcher := &Matcher{}

	t.Run("Exact title and year match", func(t *testing.T) {
		result := ScanResult{
			ParsedTitle: "The Matrix",
			ParsedYear:  intPtr(1999),
		}
		releaseDate := time.Date(1999, 3, 31, 0, 0, 0, 0, time.UTC)
		tmdbMovie := &Movie{
			Title:       "The Matrix",
			ReleaseDate: &releaseDate,
			Popularity:  decimalPtr("100.5"),
		}

		confidence := matcher.calculateConfidence(result, tmdbMovie)

		// 0.5 (exact title) + 0.3 (exact year) + 0.1 (popularity > 100) = 0.9
		assert.InDelta(t, 0.9, confidence, 0.001)
	})

	t.Run("Exact title, no year in filename", func(t *testing.T) {
		result := ScanResult{
			ParsedTitle: "The Matrix",
			ParsedYear:  nil,
		}
		releaseDate := time.Date(1999, 3, 31, 0, 0, 0, 0, time.UTC)
		tmdbMovie := &Movie{
			Title:       "The Matrix",
			ReleaseDate: &releaseDate,
			Popularity:  nil, // No popularity data
		}

		confidence := matcher.calculateConfidence(result, tmdbMovie)

		// 0.5 (exact title) - 0.05 (no year penalty) = 0.45
		assert.InDelta(t, 0.45, confidence, 0.001)
	})

	t.Run("Partial title match, exact year", func(t *testing.T) {
		result := ScanResult{
			ParsedTitle: "Matrix",
			ParsedYear:  intPtr(1999),
		}
		releaseDate := time.Date(1999, 3, 31, 0, 0, 0, 0, time.UTC)
		tmdbMovie := &Movie{
			Title:       "The Matrix",
			ReleaseDate: &releaseDate,
			Popularity:  nil,
		}

		confidence := matcher.calculateConfidence(result, tmdbMovie)

		// "matrix" vs "the matrix": after normalization (article removal), both become "matrix"
		// similarity = 1.0, title_confidence = 1.0 * 0.5 = 0.5
		// 0.5 (exact normalized title) + 0.3 (exact year) = 0.8
		assert.InDelta(t, 0.8, confidence, 0.001)
	})

	t.Run("Fuzzy title match - typo in filename", func(t *testing.T) {
		result := ScanResult{
			ParsedTitle: "Matrx", // Typo - missing 'i'
			ParsedYear:  intPtr(1999),
		}
		releaseDate := time.Date(1999, 3, 31, 0, 0, 0, 0, time.UTC)
		tmdbMovie := &Movie{
			Title:       "The Matrix",
			ReleaseDate: &releaseDate,
			Popularity:  nil,
		}

		confidence := matcher.calculateConfidence(result, tmdbMovie)

		// "matrx" vs "matrix": Levenshtein distance = 1, maxLen = 6
		// similarity = 1.0 - 1/6 ≈ 0.833, title_confidence ≈ 0.417
		// ≈0.417 (fuzzy title) + 0.3 (exact year) ≈ 0.717
		assert.InDelta(t, 0.717, confidence, 0.05) // Allow more variance for fuzzy
		assert.GreaterOrEqual(t, confidence, 0.6)
		assert.LessOrEqual(t, confidence, 0.8)
	})

	t.Run("Exact title, year off by one", func(t *testing.T) {
		result := ScanResult{
			ParsedTitle: "The Matrix",
			ParsedYear:  intPtr(1998),
		}
		releaseDate := time.Date(1999, 3, 31, 0, 0, 0, 0, time.UTC)
		tmdbMovie := &Movie{
			Title:       "The Matrix",
			ReleaseDate: &releaseDate,
			Popularity:  nil,
		}

		confidence := matcher.calculateConfidence(result, tmdbMovie)

		// 0.5 (exact title) + 0.15 (year off by 1) = 0.65
		assert.InDelta(t, 0.65, confidence, 0.001)
	})

	t.Run("Poor match - different title and year", func(t *testing.T) {
		result := ScanResult{
			ParsedTitle: "Inception",
			ParsedYear:  intPtr(2010),
		}
		releaseDate := time.Date(1999, 3, 31, 0, 0, 0, 0, time.UTC)
		tmdbMovie := &Movie{
			Title:       "The Matrix",
			ReleaseDate: &releaseDate,
			Popularity:  nil,
		}

		confidence := matcher.calculateConfidence(result, tmdbMovie)

		// "inception" vs "the matrix": very different, low Levenshtein similarity
		// Year diff = 11, no year bonus
		// Result should be low (< 0.2)
		assert.LessOrEqual(t, confidence, 0.2)
		assert.GreaterOrEqual(t, confidence, 0.0)
	})

	t.Run("Title match with high popularity boost", func(t *testing.T) {
		result := ScanResult{
			ParsedTitle: "The Matrix",
			ParsedYear:  intPtr(1999),
		}
		releaseDate := time.Date(1999, 3, 31, 0, 0, 0, 0, time.UTC)
		tmdbMovie := &Movie{
			Title:       "The Matrix",
			ReleaseDate: &releaseDate,
			Popularity:  decimalPtr("150.0"),
		}

		confidence := matcher.calculateConfidence(result, tmdbMovie)

		// 0.5 (exact title) + 0.3 (exact year) + 0.1 (popularity > 100) = 0.9
		assert.InDelta(t, 0.9, confidence, 0.001)
	})

	t.Run("Low popularity no boost", func(t *testing.T) {
		result := ScanResult{
			ParsedTitle: "Obscure Movie",
			ParsedYear:  intPtr(2020),
		}
		releaseDate := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
		tmdbMovie := &Movie{
			Title:       "Obscure Movie",
			ReleaseDate: &releaseDate,
			Popularity:  decimalPtr("5.0"),
		}

		confidence := matcher.calculateConfidence(result, tmdbMovie)

		// 0.5 (exact title) + 0.3 (exact year) = 0.8 (no popularity boost, pop < 50)
		assert.InDelta(t, 0.8, confidence, 0.001)
	})

	t.Run("Case insensitive title match", func(t *testing.T) {
		result := ScanResult{
			ParsedTitle: "the matrix",
			ParsedYear:  intPtr(1999),
		}
		releaseDate := time.Date(1999, 3, 31, 0, 0, 0, 0, time.UTC)
		tmdbMovie := &Movie{
			Title:       "The Matrix",
			ReleaseDate: &releaseDate,
			Popularity:  nil,
		}

		confidence := matcher.calculateConfidence(result, tmdbMovie)

		// 0.5 (exact title - case insensitive) + 0.3 (exact year) = 0.8
		assert.InDelta(t, 0.8, confidence, 0.001)
	})

	t.Run("Confidence minimum is 0", func(t *testing.T) {
		result := ScanResult{
			ParsedTitle: "Completely Different Movie",
			ParsedYear:  nil,
		}
		releaseDate := time.Date(1999, 3, 31, 0, 0, 0, 0, time.UTC)
		tmdbMovie := &Movie{
			Title:       "The Matrix",
			ReleaseDate: &releaseDate,
			Popularity:  nil,
		}

		confidence := matcher.calculateConfidence(result, tmdbMovie)

		// Very different titles, low Levenshtein similarity
		// title_confidence ≈ 0.05-0.1, minus 0.05 for no year
		// Result clamped to >= 0
		assert.GreaterOrEqual(t, confidence, 0.0)
		assert.LessOrEqual(t, confidence, 0.2)
	})

	t.Run("Confidence never exceeds 1", func(t *testing.T) {
		result := ScanResult{
			ParsedTitle: "The Matrix",
			ParsedYear:  intPtr(1999),
		}
		releaseDate := time.Date(1999, 3, 31, 0, 0, 0, 0, time.UTC)
		tmdbMovie := &Movie{
			Title:       "The Matrix",
			ReleaseDate: &releaseDate,
			Popularity:  decimalPtr("999999.0"),
		}

		confidence := matcher.calculateConfidence(result, tmdbMovie)

		// 0.5 (exact title) + 0.3 (exact year) + 0.1 (high popularity) = 0.9
		// Should not exceed 1.0
		assert.LessOrEqual(t, confidence, 1.0)
		assert.InDelta(t, 0.9, confidence, 0.001)
	})
}

func TestAbs(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		expected int
	}{
		{"Positive number", 5, 5},
		{"Negative number", -5, 5},
		{"Zero", 0, 0},
		{"Large positive", 1000, 1000},
		{"Large negative", -1000, 1000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := abs(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestExtractYear(t *testing.T) {
	t.Run("Valid date", func(t *testing.T) {
		date := time.Date(1999, 3, 31, 0, 0, 0, 0, time.UTC)
		year := extractYear(&date)
		assert.NotNil(t, year)
		assert.Equal(t, int32(1999), *year)
	})

	t.Run("Nil date", func(t *testing.T) {
		year := extractYear(nil)
		assert.Nil(t, year)
	})
}

func TestFormatDate(t *testing.T) {
	t.Run("Valid date", func(t *testing.T) {
		date := time.Date(1999, 3, 31, 12, 30, 45, 0, time.UTC)
		formatted := formatDate(&date)
		assert.NotNil(t, formatted)
		assert.Equal(t, "1999-03-31", *formatted)
	})

	t.Run("Nil date", func(t *testing.T) {
		formatted := formatDate(nil)
		assert.Nil(t, formatted)
	})
}

func TestFormatDecimal(t *testing.T) {
	t.Run("Valid decimal", func(t *testing.T) {
		d, _ := decimal.NewFromFloat64(123.456)
		formatted := formatDecimal(&d)
		assert.NotNil(t, formatted)
		assert.Equal(t, "123.456", *formatted)
	})

	t.Run("Zero decimal", func(t *testing.T) {
		d := decimal.Decimal{}
		formatted := formatDecimal(&d)
		assert.Nil(t, formatted)
	})

	t.Run("Nil decimal", func(t *testing.T) {
		formatted := formatDecimal(nil)
		assert.Nil(t, formatted)
	})
}

// Helper functions
func decimalPtr(s string) *decimal.Decimal {
	d, _ := decimal.Parse(s)
	return &d
}
