package metadata

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseReleaseDate(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected *time.Time
	}{
		{
			name:     "valid date",
			input:    "1999-03-31",
			expected: new(time.Date(1999, 3, 31, 0, 0, 0, 0, time.UTC)),
		},
		{
			name:     "empty string",
			input:    "",
			expected: nil,
		},
		{
			name:     "invalid format",
			input:    "31-03-1999",
			expected: nil,
		},
		{
			name:     "partial date",
			input:    "1999",
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseReleaseDate(tt.input)
			if tt.expected == nil {
				assert.Nil(t, result)
			} else {
				assert.NotNil(t, result)
				assert.Equal(t, *tt.expected, *result)
			}
		})
	}
}

func TestExtractYearFromDate(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected *int
	}{
		{
			name:     "valid date",
			input:    "1999-03-31",
			expected: new(1999),
		},
		{
			name:     "year only",
			input:    "2020",
			expected: new(2020),
		},
		{
			name:     "empty string",
			input:    "",
			expected: nil,
		},
		{
			name:     "invalid year",
			input:    "abc-03-31",
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractYearFromDate(tt.input)
			if tt.expected == nil {
				assert.Nil(t, result)
			} else {
				assert.NotNil(t, result)
				assert.Equal(t, *tt.expected, *result)
			}
		})
	}
}

func TestParseOptionalString(t *testing.T) {
	t.Run("nil input", func(t *testing.T) {
		result := ParseOptionalString(nil)
		assert.Nil(t, result)
	})

	t.Run("empty string", func(t *testing.T) {
		empty := ""
		result := ParseOptionalString(&empty)
		assert.Nil(t, result)
	})

	t.Run("valid string", func(t *testing.T) {
		valid := "test"
		result := ParseOptionalString(&valid)
		assert.NotNil(t, result)
		assert.Equal(t, "test", *result)
	})
}

func TestParseOptionalStringValue(t *testing.T) {
	t.Run("empty string", func(t *testing.T) {
		result := ParseOptionalStringValue("")
		assert.Nil(t, result)
	})

	t.Run("valid string", func(t *testing.T) {
		result := ParseOptionalStringValue("test")
		assert.NotNil(t, result)
		assert.Equal(t, "test", *result)
	})
}

func TestSafeIntToInt32(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		expected int32
	}{
		{"normal value", 100, 100},
		{"zero", 0, 0},
		{"negative", -100, -100},
		{"max int32", 2147483647, 2147483647},
		{"min int32", -2147483648, -2147483648},
		{"overflow positive", 3000000000, 2147483647},
		{"overflow negative", -3000000000, -2147483648},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SafeIntToInt32(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestParseOptionalInt32(t *testing.T) {
	t.Run("zero", func(t *testing.T) {
		result := ParseOptionalInt32(0)
		assert.Nil(t, result)
	})

	t.Run("positive", func(t *testing.T) {
		result := ParseOptionalInt32(123)
		assert.NotNil(t, result)
		assert.Equal(t, int32(123), *result)
	})

	t.Run("negative", func(t *testing.T) {
		result := ParseOptionalInt32(-456)
		assert.NotNil(t, result)
		assert.Equal(t, int32(-456), *result)
	})
}

func TestParseOptionalInt32Ptr(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		result := ParseOptionalInt32Ptr(nil)
		assert.Nil(t, result)
	})

	t.Run("zero", func(t *testing.T) {
		zero := 0
		result := ParseOptionalInt32Ptr(&zero)
		assert.Nil(t, result)
	})

	t.Run("valid", func(t *testing.T) {
		val := 123
		result := ParseOptionalInt32Ptr(&val)
		assert.NotNil(t, result)
		assert.Equal(t, int32(123), *result)
	})
}

func TestParseOptionalInt64Ptr(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		result := ParseOptionalInt64Ptr(nil)
		assert.Nil(t, result)
	})

	t.Run("zero", func(t *testing.T) {
		zero := int64(0)
		result := ParseOptionalInt64Ptr(&zero)
		assert.Nil(t, result)
	})

	t.Run("valid", func(t *testing.T) {
		val := int64(123456789)
		result := ParseOptionalInt64Ptr(&val)
		assert.NotNil(t, result)
		assert.Equal(t, int64(123456789), *result)
	})
}

func TestGetAgeRatingSystem(t *testing.T) {
	tests := []struct {
		country  string
		expected AgeRatingSystem
	}{
		{"US", AgeRatingMPAA},
		{"DE", AgeRatingFSK},
		{"GB", AgeRatingBBFC},
		{"FR", AgeRatingCNC},
		{"JP", AgeRatingEirin},
		{"KR", AgeRatingKMRB},
		{"BR", AgeRatingDJCTQ},
		{"AU", AgeRatingACB},
		{"XX", AgeRatingSystem("XX")}, // Unknown country
	}

	for _, tt := range tests {
		t.Run(tt.country, func(t *testing.T) {
			result := GetAgeRatingSystem(tt.country)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestLanguageToISO(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"full format", "en-US", "en"},
		{"full format DE", "de-DE", "de"},
		{"already ISO", "en", "en"},
		{"short code", "ja", "ja"},
		{"no hyphen", "english", "english"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := LanguageToISO(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestISOToLanguage(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"English", "en", "en-US"},
		{"German", "de", "de-DE"},
		{"French", "fr", "fr-FR"},
		{"Spanish", "es", "es-ES"},
		{"Italian", "it", "it-IT"},
		{"Portuguese", "pt", "pt-BR"},
		{"Japanese", "ja", "ja-JP"},
		{"Korean", "ko", "ko-KR"},
		{"Chinese", "zh", "zh-CN"},
		{"Russian", "ru", "ru-RU"},
		{"Unknown", "xx", "xx"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ISOToLanguage(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Helper functions
