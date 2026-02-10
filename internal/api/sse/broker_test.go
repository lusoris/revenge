package sse

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatSSE(t *testing.T) {
	t.Run("standard message", func(t *testing.T) {
		result := formatSSE("1", "message", []byte(`{"hello":"world"}`))
		expected := "id: 1\nevent: message\ndata: {\"hello\":\"world\"}\n\n"
		assert.Equal(t, expected, string(result))
	})

	t.Run("empty data", func(t *testing.T) {
		result := formatSSE("42", "ping", []byte(""))
		expected := "id: 42\nevent: ping\ndata: \n\n"
		assert.Equal(t, expected, string(result))
	})

	t.Run("custom event type", func(t *testing.T) {
		result := formatSSE("abc", "library-update", []byte("test"))
		expected := "id: abc\nevent: library-update\ndata: test\n\n"
		assert.Equal(t, expected, string(result))
	})
}
