package middleware

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseSameSite(t *testing.T) {
	tests := []struct {
		input string
		want  http.SameSite
	}{
		{"strict", http.SameSiteStrictMode},
		{"Strict", http.SameSiteStrictMode},
		{"STRICT", http.SameSiteStrictMode},
		{"none", http.SameSiteNoneMode},
		{"None", http.SameSiteNoneMode},
		{"lax", http.SameSiteLaxMode},
		{"Lax", http.SameSiteLaxMode},
		{"", http.SameSiteLaxMode},
		{"unknown", http.SameSiteLaxMode},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.want, parseSameSite(tt.input))
		})
	}
}

func TestIsSafeMethod(t *testing.T) {
	safe := []string{http.MethodGet, http.MethodHead, http.MethodOptions, http.MethodTrace}
	for _, m := range safe {
		t.Run(m+"_safe", func(t *testing.T) {
			assert.True(t, isSafeMethod(m))
		})
	}

	unsafe := []string{http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch}
	for _, m := range unsafe {
		t.Run(m+"_unsafe", func(t *testing.T) {
			assert.False(t, isSafeMethod(m))
		})
	}
}

func TestGenerateCSRFToken(t *testing.T) {
	token, err := GenerateCSRFToken()
	require.NoError(t, err)
	assert.Len(t, token, 64) // 32 bytes = 64 hex chars

	// Tokens should be unique
	token2, err := GenerateCSRFToken()
	require.NoError(t, err)
	assert.NotEqual(t, token, token2)
}
