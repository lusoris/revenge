package scanner

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScanResult_GetYear(t *testing.T) {
	tests := []struct {
		name     string
		metadata map[string]any
		expected *int
	}{
		{
			name:     "valid year",
			metadata: map[string]any{"year": 1999},
			expected: intPtr(1999),
		},
		{
			name:     "nil metadata",
			metadata: nil,
			expected: nil,
		},
		{
			name:     "no year key",
			metadata: map[string]any{"title": "test"},
			expected: nil,
		},
		{
			name:     "wrong type",
			metadata: map[string]any{"year": "1999"},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := ScanResult{Metadata: tt.metadata}
			result := r.GetYear()
			if tt.expected == nil {
				assert.Nil(t, result)
			} else {
				assert.NotNil(t, result)
				assert.Equal(t, *tt.expected, *result)
			}
		})
	}
}

func TestScanResult_GetSeason(t *testing.T) {
	tests := []struct {
		name     string
		metadata map[string]any
		expected *int
	}{
		{
			name:     "valid season",
			metadata: map[string]any{"season": 1},
			expected: intPtr(1),
		},
		{
			name:     "nil metadata",
			metadata: nil,
			expected: nil,
		},
		{
			name:     "no season key",
			metadata: map[string]any{"episode": 5},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := ScanResult{Metadata: tt.metadata}
			result := r.GetSeason()
			if tt.expected == nil {
				assert.Nil(t, result)
			} else {
				assert.NotNil(t, result)
				assert.Equal(t, *tt.expected, *result)
			}
		})
	}
}

func TestScanResult_GetEpisode(t *testing.T) {
	tests := []struct {
		name     string
		metadata map[string]any
		expected *int
	}{
		{
			name:     "valid episode",
			metadata: map[string]any{"episode": 5},
			expected: intPtr(5),
		},
		{
			name:     "nil metadata",
			metadata: nil,
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := ScanResult{Metadata: tt.metadata}
			result := r.GetEpisode()
			if tt.expected == nil {
				assert.Nil(t, result)
			} else {
				assert.NotNil(t, result)
				assert.Equal(t, *tt.expected, *result)
			}
		})
	}
}

func TestScanResult_GetString(t *testing.T) {
	tests := []struct {
		name     string
		metadata map[string]any
		key      string
		expected string
	}{
		{
			name:     "valid string",
			metadata: map[string]any{"artist": "Test Artist"},
			key:      "artist",
			expected: "Test Artist",
		},
		{
			name:     "nil metadata",
			metadata: nil,
			key:      "artist",
			expected: "",
		},
		{
			name:     "missing key",
			metadata: map[string]any{"album": "Test Album"},
			key:      "artist",
			expected: "",
		},
		{
			name:     "wrong type",
			metadata: map[string]any{"track": 5},
			key:      "track",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := ScanResult{Metadata: tt.metadata}
			assert.Equal(t, tt.expected, r.GetString(tt.key))
		})
	}
}

func TestDefaultScanOptions(t *testing.T) {
	opts := DefaultScanOptions()

	assert.False(t, opts.FollowSymlinks)
	assert.Equal(t, 0, opts.MaxDepth)
	assert.False(t, opts.IncludeHidden)
	assert.NotEmpty(t, opts.ExcludePatterns)
	assert.Contains(t, opts.ExcludePatterns, "@eaDir")
	assert.Contains(t, opts.ExcludePatterns, ".DS_Store")
}
