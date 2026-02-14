package web

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnabled_NoFrontend(t *testing.T) {
	// Default build (no frontend tag) -- assets is nil
	// We are in the !frontend build, so Enabled should return false
	if assets == nil {
		assert.False(t, Enabled())
	}
}

func TestNewSPAHandler_NoFrontend(t *testing.T) {
	// Without frontend tag, handler should be nil
	if assets == nil {
		handler := NewSPAHandler()
		assert.Nil(t, handler)
	}
}

func TestIsImmutableAsset(t *testing.T) {
	t.Parallel()

	tests := []struct {
		path     string
		expected bool
	}{
		{"_app/immutable/chunks/0.abc123.js", true},
		{"_app/immutable/entry/start.js", true},
		{"assets/style.css", true},
		{"app.js", true},
		{"fonts/inter.woff2", true},
		{"fonts/inter.woff", true},
		{"favicon.png", false},
		{"robots.txt", false},
		{"images/logo.svg", false},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			assert.Equal(t, tt.expected, isImmutableAsset(tt.path))
		})
	}
}

func TestNewSPAHandler_Interface(t *testing.T) {
	// Verify that NewSPAHandler returns an http.Handler (or nil)
	handler := NewSPAHandler()
	if handler != nil {
		var _ http.Handler = handler // type check
	}
}
