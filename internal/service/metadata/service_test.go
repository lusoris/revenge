package metadata

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMergeExternalRatings(t *testing.T) {
	t.Run("merge disjoint sets", func(t *testing.T) {
		existing := []ExternalRating{
			{Source: "IMDb", Value: "9.0/10", Score: 90.0},
		}
		additional := []ExternalRating{
			{Source: "Rotten Tomatoes", Value: "94%", Score: 94.0},
		}
		result := mergeExternalRatings(existing, additional)
		assert.Len(t, result, 2)
	})

	t.Run("dedup by Source keeps first", func(t *testing.T) {
		existing := []ExternalRating{
			{Source: "IMDb", Value: "9.0/10", Score: 90.0},
		}
		additional := []ExternalRating{
			{Source: "IMDb", Value: "8.5/10", Score: 85.0},
		}
		result := mergeExternalRatings(existing, additional)
		assert.Len(t, result, 1)
		assert.Equal(t, "9.0/10", result[0].Value)
	})

	t.Run("nil existing", func(t *testing.T) {
		additional := []ExternalRating{
			{Source: "IMDb", Value: "9.0/10", Score: 90.0},
		}
		result := mergeExternalRatings(nil, additional)
		assert.Len(t, result, 1)
	})

	t.Run("nil additional", func(t *testing.T) {
		existing := []ExternalRating{
			{Source: "IMDb", Value: "9.0/10", Score: 90.0},
		}
		result := mergeExternalRatings(existing, nil)
		assert.Len(t, result, 1)
	})

	t.Run("both empty", func(t *testing.T) {
		result := mergeExternalRatings(nil, nil)
		assert.Empty(t, result)
	})
}

func TestPtrToString(t *testing.T) {
	assert.Equal(t, "", ptrToString(nil))

	s := "hello"
	assert.Equal(t, "hello", ptrToString(&s))

	empty := ""
	assert.Equal(t, "", ptrToString(&empty))
}

func TestDefaultServiceConfig(t *testing.T) {
	cfg := DefaultServiceConfig()
	assert.Equal(t, []string{"en"}, cfg.DefaultLanguages)
	assert.True(t, cfg.EnableProviderFallback)
	assert.False(t, cfg.EnableEnrichment)
}
