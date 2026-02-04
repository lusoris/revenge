package storage

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/lusoris/revenge/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestLocalStorage_Store(t *testing.T) {
	// Create temp directory for test
	tmpDir, err := os.MkdirTemp("", "storage-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	cfg := config.AvatarConfig{
		StoragePath: tmpDir,
	}

	storage, err := NewLocalStorage(cfg, zap.NewNop())
	require.NoError(t, err)

	ctx := context.Background()

	t.Run("store and retrieve file", func(t *testing.T) {
		data := []byte("test file content")
		key := "test/file.txt"

		storedKey, err := storage.Store(ctx, key, bytes.NewReader(data), "text/plain")
		require.NoError(t, err)
		assert.Equal(t, key, storedKey)

		// Verify file exists
		exists, err := storage.Exists(ctx, key)
		require.NoError(t, err)
		assert.True(t, exists)

		// Get file
		reader, err := storage.Get(ctx, key)
		require.NoError(t, err)
		defer reader.Close()

		content, err := os.ReadFile(filepath.Join(tmpDir, key))
		require.NoError(t, err)
		assert.Equal(t, data, content)
	})

	t.Run("delete file", func(t *testing.T) {
		data := []byte("to be deleted")
		key := "delete/me.txt"

		_, err := storage.Store(ctx, key, bytes.NewReader(data), "text/plain")
		require.NoError(t, err)

		// Verify exists
		exists, err := storage.Exists(ctx, key)
		require.NoError(t, err)
		assert.True(t, exists)

		// Delete
		err = storage.Delete(ctx, key)
		require.NoError(t, err)

		// Verify deleted
		exists, err = storage.Exists(ctx, key)
		require.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("get URL", func(t *testing.T) {
		key := "avatars/user123/image.png"
		url := storage.GetURL(key)
		assert.Equal(t, "/api/v1/files/avatars/user123/image.png", url)
	})
}

func TestSanitizeKey(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"normal/path/file.txt", "normal/path/file.txt"},
		{"/leading/slash.txt", "leading/slash.txt"},
		{"../parent/traversal.txt", "parent/traversal.txt"},
		{"path/../sneaky/file.txt", "sneaky/file.txt"}, // filepath.Clean normalizes this
		{"./current/dir.txt", "current/dir.txt"},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := sanitizeKey(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGenerateAvatarKey(t *testing.T) {
	userID := uuid.New()

	t.Run("with extension", func(t *testing.T) {
		key := GenerateAvatarKey(userID, "profile.jpg")
		assert.Contains(t, key, "avatars/"+userID.String())
		assert.Contains(t, key, ".jpg")
	})

	t.Run("without extension", func(t *testing.T) {
		key := GenerateAvatarKey(userID, "noext")
		assert.Contains(t, key, "avatars/"+userID.String())
		assert.Contains(t, key, ".png") // Default extension
	})
}

func TestMockStorage(t *testing.T) {
	storage := NewMockStorage()
	ctx := context.Background()

	t.Run("store and retrieve", func(t *testing.T) {
		data := []byte("mock data")
		key := "mock/file.txt"

		storedKey, err := storage.Store(ctx, key, bytes.NewReader(data), "text/plain")
		require.NoError(t, err)
		assert.Equal(t, key, storedKey)

		exists, err := storage.Exists(ctx, key)
		require.NoError(t, err)
		assert.True(t, exists)

		reader, err := storage.Get(ctx, key)
		require.NoError(t, err)

		buf := new(bytes.Buffer)
		_, err = buf.ReadFrom(reader)
		require.NoError(t, err)
		assert.Equal(t, data, buf.Bytes())
	})

	t.Run("delete", func(t *testing.T) {
		key := "to-delete.txt"
		_, err := storage.Store(ctx, key, bytes.NewReader([]byte("x")), "text/plain")
		require.NoError(t, err)

		err = storage.Delete(ctx, key)
		require.NoError(t, err)

		exists, err := storage.Exists(ctx, key)
		require.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("get stored files", func(t *testing.T) {
		storage := NewMockStorage()
		_, err := storage.Store(ctx, "a.txt", bytes.NewReader([]byte("a")), "text/plain")
		require.NoError(t, err)
		_, err = storage.Store(ctx, "b.txt", bytes.NewReader([]byte("b")), "text/plain")
		require.NoError(t, err)

		files := storage.GetStoredFiles()
		assert.Len(t, files, 2)
		assert.Equal(t, []byte("a"), files["a.txt"])
		assert.Equal(t, []byte("b"), files["b.txt"])
	})
}
