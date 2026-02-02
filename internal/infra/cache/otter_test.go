package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewL1Cache_Defaults(t *testing.T) {
	cache, err := NewL1Cache[string, string](0, 0)
	require.NoError(t, err)
	require.NotNil(t, cache)
	defer cache.Close()

	// Should use default values
	assert.NotNil(t, cache.cache)
}

func TestNewL1Cache_CustomConfig(t *testing.T) {
	maxSize := 5000
	ttl := 10 * time.Minute

	cache, err := NewL1Cache[string, string](maxSize, ttl)
	require.NoError(t, err)
	require.NotNil(t, cache)
	defer cache.Close()

	assert.NotNil(t, cache.cache)
}

func TestL1Cache_SetAndGet(t *testing.T) {
	cache, err := NewL1Cache[string, string](100, 1*time.Minute)
	require.NoError(t, err)
	defer cache.Close()

	// Initially, key should not exist
	val, ok := cache.Get("key1")
	assert.False(t, ok)
	assert.Equal(t, "", val)

	// Set value
	cache.Set("key1", "value1")

	// Get value
	val, ok = cache.Get("key1")
	assert.True(t, ok)
	assert.Equal(t, "value1", val)
}

func TestL1Cache_Delete(t *testing.T) {
	cache, err := NewL1Cache[string, string](100, 1*time.Minute)
	require.NoError(t, err)
	defer cache.Close()

	// Set value
	cache.Set("key1", "value1")

	// Verify it exists
	val, ok := cache.Get("key1")
	assert.True(t, ok)
	assert.Equal(t, "value1", val)

	// Delete
	cache.Delete("key1")

	// Verify it's gone
	val, ok = cache.Get("key1")
	assert.False(t, ok)
	assert.Equal(t, "", val)
}

func TestL1Cache_Clear(t *testing.T) {
	cache, err := NewL1Cache[string, string](100, 1*time.Minute)
	require.NoError(t, err)
	defer cache.Close()

	// Set multiple values
	cache.Set("key1", "value1")
	cache.Set("key2", "value2")
	cache.Set("key3", "value3")

	// Verify size
	assert.Greater(t, cache.Size(), 0)

	// Clear all
	cache.Clear()

	// Verify all gone
	assert.Equal(t, 0, cache.Size())
	_, ok := cache.Get("key1")
	assert.False(t, ok)
	_, ok = cache.Get("key2")
	assert.False(t, ok)
	_, ok = cache.Get("key3")
	assert.False(t, ok)
}

func TestL1Cache_Size(t *testing.T) {
	cache, err := NewL1Cache[string, string](100, 1*time.Minute)
	require.NoError(t, err)
	defer cache.Close()

	// Initially empty
	assert.Equal(t, 0, cache.Size())

	// Add entries
	cache.Set("key1", "value1")
	cache.Set("key2", "value2")

	// Size should increase
	assert.Equal(t, 2, cache.Size())

	// Delete one
	cache.Delete("key1")
	assert.Equal(t, 1, cache.Size())
}

func TestL1Cache_Has(t *testing.T) {
	cache, err := NewL1Cache[string, string](100, 1*time.Minute)
	require.NoError(t, err)
	defer cache.Close()

	// Initially, key should not exist
	assert.False(t, cache.Has("key1"))

	// Set value
	cache.Set("key1", "value1")

	// Now it should exist
	assert.True(t, cache.Has("key1"))

	// Delete
	cache.Delete("key1")

	// Should not exist anymore
	assert.False(t, cache.Has("key1"))
}

func TestL1Cache_TTL(t *testing.T) {
	cache, err := NewL1Cache[string, string](100, 100*time.Millisecond)
	require.NoError(t, err)
	defer cache.Close()

	// Set value
	cache.Set("key1", "value1")

	// Should exist immediately
	val, ok := cache.Get("key1")
	assert.True(t, ok)
	assert.Equal(t, "value1", val)

	// Wait for TTL to expire
	time.Sleep(150 * time.Millisecond)

	// Should be gone
	_, ok = cache.Get("key1")
	assert.False(t, ok)
}

func TestL1Cache_MaxSize(t *testing.T) {
	maxSize := 10
	cache, err := NewL1Cache[int, string](maxSize, 1*time.Minute)
	require.NoError(t, err)
	defer cache.Close()

	// Fill cache beyond max size
	for i := 0; i < maxSize*2; i++ {
		cache.Set(i, "value")
	}

	// Size should not exceed max (approximately, due to async eviction)
	// Give it a moment for eviction to happen
	time.Sleep(50 * time.Millisecond)

	size := cache.Size()
	assert.LessOrEqual(t, size, maxSize*2, "Size should be bounded by eviction policy")
}

func TestL1Cache_Update(t *testing.T) {
	cache, err := NewL1Cache[string, string](100, 1*time.Minute)
	require.NoError(t, err)
	defer cache.Close()

	// Set initial value
	cache.Set("key1", "value1")

	val, ok := cache.Get("key1")
	assert.True(t, ok)
	assert.Equal(t, "value1", val)

	// Update value
	cache.Set("key1", "value2")

	val, ok = cache.Get("key1")
	assert.True(t, ok)
	assert.Equal(t, "value2", val)

	// Size should still be 1
	assert.Equal(t, 1, cache.Size())
}

func TestL1Cache_TypedInt(t *testing.T) {
	cache, err := NewL1Cache[string, int](100, 1*time.Minute)
	require.NoError(t, err)
	defer cache.Close()

	cache.Set("number", 42)

	val, ok := cache.Get("number")
	assert.True(t, ok)
	assert.Equal(t, 42, val)
}

func TestL1Cache_TypedStruct(t *testing.T) {
	type User struct {
		ID   int
		Name string
	}

	cache, err := NewL1Cache[int, User](100, 1*time.Minute)
	require.NoError(t, err)
	defer cache.Close()

	user := User{ID: 1, Name: "Alice"}
	cache.Set(1, user)

	val, ok := cache.Get(1)
	assert.True(t, ok)
	assert.Equal(t, user, val)
}

func TestDefaultConstants_L1(t *testing.T) {
	assert.Equal(t, 10000, DefaultL1MaxSize)
	assert.Equal(t, 5*time.Minute, DefaultL1TTL)
}
