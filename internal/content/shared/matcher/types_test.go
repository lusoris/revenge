package matcher

import (
	"context"
	"errors"
	"testing"

	"github.com/lusoris/revenge/internal/content/shared/scanner"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// testContent is a simple content type for testing
type testContent struct {
	ID    string
	Title string
	Year  int
}

// mockStrategy is a test implementation of MatchStrategy
type mockStrategy struct {
	findExistingFn    func(ctx context.Context, sr scanner.ScanResult) (*testContent, float64, error)
	searchExternalFn  func(ctx context.Context, sr scanner.ScanResult) ([]*testContent, error)
	calculateConfFn   func(sr scanner.ScanResult, c *testContent) float64
	createContentFn   func(ctx context.Context, c *testContent) (*testContent, error)
}

func (m *mockStrategy) FindExisting(ctx context.Context, sr scanner.ScanResult) (*testContent, float64, error) {
	if m.findExistingFn != nil {
		return m.findExistingFn(ctx, sr)
	}
	return nil, 0, errors.New("not found")
}

func (m *mockStrategy) SearchExternal(ctx context.Context, sr scanner.ScanResult) ([]*testContent, error) {
	if m.searchExternalFn != nil {
		return m.searchExternalFn(ctx, sr)
	}
	return nil, nil
}

func (m *mockStrategy) CalculateConfidence(sr scanner.ScanResult, c *testContent) float64 {
	if m.calculateConfFn != nil {
		return m.calculateConfFn(sr, c)
	}
	return 0.9
}

func (m *mockStrategy) CreateContent(ctx context.Context, c *testContent) (*testContent, error) {
	if m.createContentFn != nil {
		return m.createContentFn(ctx, c)
	}
	return c, nil
}

func TestMatchResult_IsMatched(t *testing.T) {
	t.Run("matched with content", func(t *testing.T) {
		result := MatchResult[testContent]{
			Content:   &testContent{ID: "1", Title: "Test"},
			MatchType: MatchTypeTitle,
		}
		assert.True(t, result.IsMatched())
	})

	t.Run("unmatched with nil content", func(t *testing.T) {
		result := MatchResult[testContent]{
			Content:   nil,
			MatchType: MatchTypeUnmatched,
		}
		assert.False(t, result.IsMatched())
	})

	t.Run("has error", func(t *testing.T) {
		result := MatchResult[testContent]{
			Content:   &testContent{ID: "1", Title: "Test"},
			MatchType: MatchTypeTitle,
			Error:     errors.New("some error"),
		}
		assert.False(t, result.IsMatched())
	})
}

func TestMatcher_MatchFile_FindsExisting(t *testing.T) {
	existingContent := &testContent{ID: "existing", Title: "The Matrix", Year: 1999}

	strategy := &mockStrategy{
		findExistingFn: func(ctx context.Context, sr scanner.ScanResult) (*testContent, float64, error) {
			return existingContent, 0.95, nil
		},
	}

	matcher := NewMatcher[testContent](strategy)

	scanResult := scanner.ScanResult{
		FilePath:    "/movies/The Matrix (1999).mkv",
		FileName:    "The Matrix (1999).mkv",
		ParsedTitle: "The Matrix",
		Metadata:    map[string]any{"year": 1999},
	}

	result := matcher.MatchFile(context.Background(), scanResult)

	assert.True(t, result.IsMatched())
	assert.Equal(t, MatchTypeTitle, result.MatchType)
	assert.Equal(t, existingContent, result.Content)
	assert.Equal(t, 0.95, result.Confidence)
	assert.False(t, result.CreatedNew)
}

func TestMatcher_MatchFile_SearchesExternal(t *testing.T) {
	externalContent := &testContent{ID: "tmdb-123", Title: "Inception", Year: 2010}

	strategy := &mockStrategy{
		findExistingFn: func(ctx context.Context, sr scanner.ScanResult) (*testContent, float64, error) {
			return nil, 0, errors.New("not found")
		},
		searchExternalFn: func(ctx context.Context, sr scanner.ScanResult) ([]*testContent, error) {
			return []*testContent{externalContent}, nil
		},
		calculateConfFn: func(sr scanner.ScanResult, c *testContent) float64 {
			return 0.85
		},
		createContentFn: func(ctx context.Context, c *testContent) (*testContent, error) {
			return &testContent{ID: "new-1", Title: c.Title, Year: c.Year}, nil
		},
	}

	matcher := NewMatcher[testContent](strategy)

	scanResult := scanner.ScanResult{
		FilePath:    "/movies/Inception.2010.mkv",
		FileName:    "Inception.2010.mkv",
		ParsedTitle: "Inception",
		Metadata:    map[string]any{"year": 2010},
	}

	result := matcher.MatchFile(context.Background(), scanResult)

	assert.True(t, result.IsMatched())
	assert.Equal(t, MatchTypeTitle, result.MatchType)
	require.NotNil(t, result.Content)
	assert.Equal(t, "new-1", result.Content.ID)
	assert.Equal(t, 0.85, result.Confidence)
	assert.True(t, result.CreatedNew)
}

func TestMatcher_MatchFile_NoMatches(t *testing.T) {
	strategy := &mockStrategy{
		findExistingFn: func(ctx context.Context, sr scanner.ScanResult) (*testContent, float64, error) {
			return nil, 0, errors.New("not found")
		},
		searchExternalFn: func(ctx context.Context, sr scanner.ScanResult) ([]*testContent, error) {
			return []*testContent{}, nil // Empty results
		},
	}

	matcher := NewMatcher[testContent](strategy)

	scanResult := scanner.ScanResult{
		FilePath:    "/movies/Unknown Movie.mkv",
		FileName:    "Unknown Movie.mkv",
		ParsedTitle: "Unknown Movie",
	}

	result := matcher.MatchFile(context.Background(), scanResult)

	assert.False(t, result.IsMatched())
	assert.Equal(t, MatchTypeUnmatched, result.MatchType)
	assert.ErrorIs(t, result.Error, ErrNoMatches)
}

func TestMatcher_MatchFile_ExternalSearchError(t *testing.T) {
	strategy := &mockStrategy{
		findExistingFn: func(ctx context.Context, sr scanner.ScanResult) (*testContent, float64, error) {
			return nil, 0, errors.New("not found")
		},
		searchExternalFn: func(ctx context.Context, sr scanner.ScanResult) ([]*testContent, error) {
			return nil, errors.New("API error")
		},
	}

	matcher := NewMatcher[testContent](strategy)

	scanResult := scanner.ScanResult{
		FilePath:    "/movies/Test.mkv",
		FileName:    "Test.mkv",
		ParsedTitle: "Test",
	}

	result := matcher.MatchFile(context.Background(), scanResult)

	assert.False(t, result.IsMatched())
	assert.Equal(t, MatchTypeUnmatched, result.MatchType)
	assert.Error(t, result.Error)
}

func TestMatcher_MatchFile_FuzzyMatchType(t *testing.T) {
	externalContent := &testContent{ID: "tmdb-456", Title: "Similar Movie", Year: 2020}

	strategy := &mockStrategy{
		findExistingFn: func(ctx context.Context, sr scanner.ScanResult) (*testContent, float64, error) {
			return nil, 0, errors.New("not found")
		},
		searchExternalFn: func(ctx context.Context, sr scanner.ScanResult) ([]*testContent, error) {
			return []*testContent{externalContent}, nil
		},
		calculateConfFn: func(sr scanner.ScanResult, c *testContent) float64 {
			return 0.6 // Below 0.7 threshold
		},
	}

	matcher := NewMatcher[testContent](strategy)

	scanResult := scanner.ScanResult{
		FilePath:    "/movies/Simliar Movie.mkv", // Typo
		FileName:    "Simliar Movie.mkv",
		ParsedTitle: "Simliar Movie",
	}

	result := matcher.MatchFile(context.Background(), scanResult)

	assert.True(t, result.IsMatched())
	assert.Equal(t, MatchTypeFuzzy, result.MatchType)
	assert.Equal(t, 0.6, result.Confidence)
}

func TestMatcher_MatchFiles(t *testing.T) {
	content1 := &testContent{ID: "1", Title: "Movie 1", Year: 2020}
	content2 := &testContent{ID: "2", Title: "Movie 2", Year: 2021}

	strategy := &mockStrategy{
		findExistingFn: func(ctx context.Context, sr scanner.ScanResult) (*testContent, float64, error) {
			if sr.ParsedTitle == "Movie 1" {
				return content1, 0.95, nil
			}
			return nil, 0, errors.New("not found")
		},
		searchExternalFn: func(ctx context.Context, sr scanner.ScanResult) ([]*testContent, error) {
			if sr.ParsedTitle == "Movie 2" {
				return []*testContent{content2}, nil
			}
			return nil, nil
		},
	}

	matcher := NewMatcher[testContent](strategy)

	scanResults := []scanner.ScanResult{
		{FilePath: "/movies/Movie 1.mkv", FileName: "Movie 1.mkv", ParsedTitle: "Movie 1"},
		{FilePath: "/movies/Movie 2.mkv", FileName: "Movie 2.mkv", ParsedTitle: "Movie 2"},
		{FilePath: "/movies/Movie 3.mkv", FileName: "Movie 3.mkv", ParsedTitle: "Movie 3"},
	}

	results, err := matcher.MatchFiles(context.Background(), scanResults)
	require.NoError(t, err)
	require.Len(t, results, 3)

	// Movie 1: found existing
	assert.True(t, results[0].IsMatched())
	assert.False(t, results[0].CreatedNew)

	// Movie 2: found via external search
	assert.True(t, results[1].IsMatched())
	assert.True(t, results[1].CreatedNew)

	// Movie 3: no match
	assert.False(t, results[2].IsMatched())
}
