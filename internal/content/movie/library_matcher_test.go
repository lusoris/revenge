package movie

import (
	"testing"
	"time"

	"github.com/shopspring/decimal"
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

		// 0.6 (exact title) + 0.3 (exact year) + 0.1 (popularity > 50) = 1.0
		assert.InDelta(t, 1.0, confidence, 0.001) // Allow small floating point difference
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

		// 0.6 (exact title) - 0.05 (no year penalty) = 0.55
		assert.InDelta(t, 0.55, confidence, 0.001)
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

		// 0.4 (partial title) + 0.3 (exact year) = 0.7
		assert.InDelta(t, 0.7, confidence, 0.001)
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

		// 0.6 (exact title) + 0.1 (year off by 1) = 0.7
		assert.InDelta(t, 0.7, confidence, 0.001)
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

		// 0.2 (no title match) + 0 (year mismatch > 1) = 0.2
		assert.InDelta(t, 0.2, confidence, 0.001)
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

		// 0.6 (exact title) + 0.3 (exact year) + 0.1 (popularity > 50) = 1.0
		assert.InDelta(t, 1.0, confidence, 0.001)
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

		// 0.6 (exact title) + 0.3 (exact year) = 0.9 (no popularity boost)
		assert.InDelta(t, 0.9, confidence, 0.001)
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

		// 0.6 (exact title - case insensitive) + 0.3 (exact year) = 0.9
		assert.InDelta(t, 0.9, confidence, 0.001)
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

		// 0.2 (no title match) - 0.05 (no year penalty) = 0.15
		assert.InDelta(t, 0.15, confidence, 0.001)
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

		// Should be capped at 1.0
		assert.InDelta(t, 1.0, confidence, 0.001)
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
		d := decimal.NewFromFloat(123.456)
		formatted := formatDecimal(&d)
		assert.NotNil(t, formatted)
		assert.Equal(t, "123.456", *formatted)
	})

	t.Run("Zero decimal", func(t *testing.T) {
		d := decimal.Zero
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
	d, _ := decimal.NewFromString(s)
	return &d
}
