package playback

import (
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelWarn}))
}

func TestSessionManager_CreateAndGet(t *testing.T) {
	sm, err := NewSessionManager(10, 30*time.Minute, testLogger())
	require.NoError(t, err)
	defer sm.Close()

	session := &Session{
		ID:        uuid.New(),
		UserID:    uuid.New(),
		MediaType: MediaTypeMovie,
		MediaID:   uuid.New(),
		FileID:    uuid.New(),
		FilePath:  "/media/movies/test.mkv",
	}

	err = sm.Create(session)
	require.NoError(t, err)

	got, ok := sm.Get(session.ID)
	require.True(t, ok)
	assert.Equal(t, session.ID, got.ID)
	assert.Equal(t, session.FilePath, got.FilePath)
	assert.False(t, got.CreatedAt.IsZero(), "CreatedAt should be set")
	assert.False(t, got.ExpiresAt.IsZero(), "ExpiresAt should be set")
}

func TestSessionManager_GetNotFound(t *testing.T) {
	sm, err := NewSessionManager(10, 30*time.Minute, testLogger())
	require.NoError(t, err)
	defer sm.Close()

	_, ok := sm.Get(uuid.New())
	assert.False(t, ok)
}

func TestSessionManager_MaxConcurrentSessions(t *testing.T) {
	sm, err := NewSessionManager(2, 30*time.Minute, testLogger())
	require.NoError(t, err)
	defer sm.Close()

	// Create 2 sessions — should succeed
	for range 2 {
		err = sm.Create(&Session{
			ID:        uuid.New(),
			UserID:    uuid.New(),
			MediaType: MediaTypeMovie,
			MediaID:   uuid.New(),
		})
		require.NoError(t, err)
	}

	// Third session should fail
	err = sm.Create(&Session{
		ID:        uuid.New(),
		UserID:    uuid.New(),
		MediaType: MediaTypeMovie,
		MediaID:   uuid.New(),
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "maximum concurrent sessions")
}

func TestSessionManager_Delete(t *testing.T) {
	sm, err := NewSessionManager(10, 30*time.Minute, testLogger())
	require.NoError(t, err)
	defer sm.Close()

	id := uuid.New()
	err = sm.Create(&Session{
		ID:        id,
		UserID:    uuid.New(),
		MediaType: MediaTypeMovie,
		MediaID:   uuid.New(),
	})
	require.NoError(t, err)
	assert.Equal(t, 1, sm.ActiveCount())

	removed := sm.Delete(id)
	assert.NotNil(t, removed)
	assert.Equal(t, id, removed.ID)
	assert.Equal(t, 0, sm.ActiveCount())

	// Get should return false
	_, ok := sm.Get(id)
	assert.False(t, ok)
}

func TestSessionManager_DeleteNotFound(t *testing.T) {
	sm, err := NewSessionManager(10, 30*time.Minute, testLogger())
	require.NoError(t, err)
	defer sm.Close()

	removed := sm.Delete(uuid.New())
	assert.Nil(t, removed)
}

func TestSessionManager_Touch(t *testing.T) {
	sm, err := NewSessionManager(10, 30*time.Minute, testLogger())
	require.NoError(t, err)
	defer sm.Close()

	id := uuid.New()
	err = sm.Create(&Session{
		ID:        id,
		UserID:    uuid.New(),
		MediaType: MediaTypeMovie,
		MediaID:   uuid.New(),
	})
	require.NoError(t, err)

	session, _ := sm.Get(id)
	originalExpiry := session.ExpiresAt

	time.Sleep(10 * time.Millisecond)

	ok := sm.Touch(id)
	assert.True(t, ok)

	session, _ = sm.Get(id)
	assert.True(t, session.ExpiresAt.After(originalExpiry), "ExpiresAt should be updated after touch")
}

func TestSessionManager_TouchNotFound(t *testing.T) {
	sm, err := NewSessionManager(10, 30*time.Minute, testLogger())
	require.NoError(t, err)
	defer sm.Close()

	ok := sm.Touch(uuid.New())
	assert.False(t, ok)
}

func TestSessionManager_DeleteFreesSlot(t *testing.T) {
	sm, err := NewSessionManager(1, 30*time.Minute, testLogger())
	require.NoError(t, err)
	defer sm.Close()

	id := uuid.New()
	err = sm.Create(&Session{
		ID:        id,
		UserID:    uuid.New(),
		MediaType: MediaTypeMovie,
		MediaID:   uuid.New(),
	})
	require.NoError(t, err)

	// Should fail — max reached
	err = sm.Create(&Session{
		ID:        uuid.New(),
		UserID:    uuid.New(),
		MediaType: MediaTypeMovie,
		MediaID:   uuid.New(),
	})
	assert.Error(t, err)

	// Delete frees the slot
	sm.Delete(id)

	// Should succeed now
	err = sm.Create(&Session{
		ID:        uuid.New(),
		UserID:    uuid.New(),
		MediaType: MediaTypeMovie,
		MediaID:   uuid.New(),
	})
	assert.NoError(t, err)
}
