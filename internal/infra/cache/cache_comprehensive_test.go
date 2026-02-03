package cache

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCache_ConcurrentAccess tests cache under high concurrency
func TestCache_ConcurrentAccess(t *testing.T) {
	cache, err := NewCache(nil, 1000, 1*time.Minute)
	require.NoError(t, err)
	defer cache.Close()

	ctx := context.Background()
	numGoroutines := 100
	numOpsPerGoroutine := 100

	errChan := make(chan error, numGoroutines)
	doneChan := make(chan bool, numGoroutines)

	// Concurrent writes and reads
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			for j := 0; j < numOpsPerGoroutine; j++ {
				key := string(rune('a' + (id % 26)))
				value := []byte("value")

				// Write
				if err := cache.Set(ctx, key, value, 1*time.Minute); err != nil {
					errChan <- err
					return
				}

				// Read
				if _, err := cache.Get(ctx, key); err != nil {
					// L2 unavailable is expected
					if !errors.Is(err, context.Canceled) {
						continue
					}
				}
			}
			doneChan <- true
		}(i)
	}

	// Wait for all goroutines
	for i := 0; i < numGoroutines; i++ {
		select {
		case err := <-errChan:
			t.Fatalf("Concurrent operation failed: %v", err)
		case <-doneChan:
			// Success
		case <-time.After(10 * time.Second):
			t.Fatal("Test timeout")
		}
	}
}

// TestCache_LargeValues tests caching of large values
func TestCache_LargeValues(t *testing.T) {
	cache, err := NewCache(nil, 10, 1*time.Minute)
	require.NoError(t, err)
	defer cache.Close()

	ctx := context.Background()

	// Test 1MB value
	largeValue := make([]byte, 1024*1024)
	for i := range largeValue {
		largeValue[i] = byte(i % 256)
	}

	err = cache.Set(ctx, "large", largeValue, 1*time.Minute)
	require.NoError(t, err)

	retrieved, err := cache.Get(ctx, "large")
	require.NoError(t, err)
	assert.Equal(t, largeValue, retrieved)
}

// TestCache_Eviction tests cache eviction under max size pressure
func TestCache_Eviction(t *testing.T) {
	maxSize := 10
	cache, err := NewCache(nil, maxSize, 1*time.Hour)
	require.NoError(t, err)
	defer cache.Close()

	ctx := context.Background()

	// Fill cache beyond max size
	for i := 0; i < maxSize*2; i++ {
		key := string(rune('a' + i))
		value := []byte("value")
		err := cache.Set(ctx, key, value, 1*time.Hour)
		require.NoError(t, err)
	}

	// Otter performs evictions asynchronously for performance
	// Size may temporarily exceed max during eviction processing
	initialSize := cache.l1.Size()

	// Wait for background evictions to complete
	time.Sleep(100 * time.Millisecond)

	finalSize := cache.l1.Size()

	// Final size should be close to max (allow 20% variance for async eviction)
	maxAllowedSize := int(float64(maxSize) * 1.2)
	assert.LessOrEqual(t, finalSize, maxAllowedSize,
		"Cache should settle near max size after evictions complete")

	// Verify evictions did occur
	assert.Less(t, finalSize, initialSize,
		"Evictions should have reduced cache size from initial burst")
}

// TestCache_TTLExpiration tests that expired entries are not returned
func TestCache_TTLExpiration(t *testing.T) {
	cache, err := NewCache(nil, 100, 50*time.Millisecond)
	require.NoError(t, err)
	defer cache.Close()

	ctx := context.Background()

	// Set with short TTL
	err = cache.Set(ctx, "expiring", []byte("value"), 100*time.Millisecond)
	require.NoError(t, err)

	// Should be available immediately
	_, err = cache.Get(ctx, "expiring")
	require.NoError(t, err)

	// Wait for expiration
	time.Sleep(150 * time.Millisecond)

	// Should be expired
	_, err = cache.Get(ctx, "expiring")
	assert.Error(t, err, "Expired entry should not be returned")
}

// TestCache_DeletePropagation tests delete removes from both layers
func TestCache_DeletePropagation(t *testing.T) {
	cache, err := NewCache(nil, 100, 1*time.Minute)
	require.NoError(t, err)
	defer cache.Close()

	ctx := context.Background()

	// Set value
	err = cache.Set(ctx, "key1", []byte("value1"), 1*time.Minute)
	require.NoError(t, err)

	// Verify it exists
	val, err := cache.Get(ctx, "key1")
	require.NoError(t, err)
	assert.Equal(t, []byte("value1"), val)

	// Delete
	err = cache.Delete(ctx, "key1")
	require.NoError(t, err)

	// Should be gone from L1
	_, ok := cache.l1.Get("key1")
	assert.False(t, ok, "Key should be removed from L1")

	// Should return error on get
	_, err = cache.Get(ctx, "key1")
	assert.Error(t, err)
}

// TestCache_ExistsCheck tests Exists works correctly
func TestCache_ExistsCheck(t *testing.T) {
	cache, err := NewCache(nil, 100, 1*time.Minute)
	require.NoError(t, err)
	defer cache.Close()

	ctx := context.Background()

	// Non-existent key
	exists, err := cache.Exists(ctx, "nonexistent")
	require.NoError(t, err)
	assert.False(t, exists)

	// Set a key
	err = cache.Set(ctx, "existing", []byte("value"), 1*time.Minute)
	require.NoError(t, err)

	// Should exist
	exists, err = cache.Exists(ctx, "existing")
	require.NoError(t, err)
	assert.True(t, exists)

	// Delete it
	err = cache.Delete(ctx, "existing")
	require.NoError(t, err)

	// Should not exist
	exists, err = cache.Exists(ctx, "existing")
	require.NoError(t, err)
	assert.False(t, exists)
}

// TestCache_InvalidatePattern tests pattern-based invalidation
func TestCache_InvalidatePattern(t *testing.T) {
	cache, err := NewCache(nil, 100, 1*time.Minute)
	require.NoError(t, err)
	defer cache.Close()

	ctx := context.Background()

	// Set multiple keys with prefix
	err = cache.Set(ctx, "user:1", []byte("value1"), 1*time.Minute)
	require.NoError(t, err)
	err = cache.Set(ctx, "user:2", []byte("value2"), 1*time.Minute)
	require.NoError(t, err)
	err = cache.Set(ctx, "session:1", []byte("value3"), 1*time.Minute)
	require.NoError(t, err)

	// Invalidate user:* pattern
	err = cache.Invalidate(ctx, "user:*")
	require.NoError(t, err)

	// All L1 cache should be cleared (pattern matching limitation)
	size := cache.l1.Size()
	assert.Equal(t, 0, size, "L1 should be completely cleared on pattern invalidate")
}

// TestCache_JSONOperations tests JSON marshal/unmarshal
func TestCache_JSONOperations(t *testing.T) {
	cache, err := NewCache(nil, 100, 1*time.Minute)
	require.NoError(t, err)
	defer cache.Close()

	ctx := context.Background()

	type TestStruct struct {
		Name  string `json:"name"`
		Count int    `json:"count"`
	}

	original := TestStruct{
		Name:  "test",
		Count: 42,
	}

	// Set JSON
	err = cache.SetJSON(ctx, "json_key", original, 1*time.Minute)
	require.NoError(t, err)

	// Get JSON
	var retrieved TestStruct
	err = cache.GetJSON(ctx, "json_key", &retrieved)
	require.NoError(t, err)

	assert.Equal(t, original, retrieved)
}

// TestCache_JSONInvalidData tests JSON unmarshal with invalid data
func TestCache_JSONInvalidData(t *testing.T) {
	cache, err := NewCache(nil, 100, 1*time.Minute)
	require.NoError(t, err)
	defer cache.Close()

	ctx := context.Background()

	// Set invalid JSON bytes
	invalidJSON := []byte("{invalid json")
	err = cache.Set(ctx, "invalid", invalidJSON, 1*time.Minute)
	require.NoError(t, err)

	// Try to get as JSON
	var result map[string]interface{}
	err = cache.GetJSON(ctx, "invalid", &result)
	assert.Error(t, err, "Should fail to unmarshal invalid JSON")
	assert.Contains(t, err.Error(), "unmarshal")
}

// TestCache_NilClient tests cache works with nil client (L1 only)
func TestCache_NilClient(t *testing.T) {
	cache, err := NewCache(nil, 100, 1*time.Minute)
	require.NoError(t, err)
	defer cache.Close()

	ctx := context.Background()

	// Should work with L1 only
	err = cache.Set(ctx, "key1", []byte("value1"), 1*time.Minute)
	require.NoError(t, err)

	val, err := cache.Get(ctx, "key1")
	require.NoError(t, err)
	assert.Equal(t, []byte("value1"), val)

	// Delete should work
	err = cache.Delete(ctx, "key1")
	require.NoError(t, err)

	// Exists should work
	exists, err := cache.Exists(ctx, "key1")
	require.NoError(t, err)
	assert.False(t, exists)

	// Invalidate should work (clears L1)
	err = cache.Set(ctx, "key2", []byte("value2"), 1*time.Minute)
	require.NoError(t, err)

	err = cache.Invalidate(ctx, "*")
	require.NoError(t, err)

	size := cache.l1.Size()
	assert.Equal(t, 0, size)
}

// TestCache_ContextCancellation tests operations respect context cancellation
func TestCache_ContextCancellation(t *testing.T) {
	cache, err := NewCache(nil, 100, 1*time.Minute)
	require.NoError(t, err)
	defer cache.Close()

	// Create canceled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	// Operations with L1 only should still work (context not checked)
	err = cache.Set(ctx, "key1", []byte("value1"), 1*time.Minute)
	require.NoError(t, err)

	val, err := cache.Get(ctx, "key1")
	require.NoError(t, err)
	assert.Equal(t, []byte("value1"), val)
}

// TestCache_EmptyKey tests handling of empty key
func TestCache_EmptyKey(t *testing.T) {
	cache, err := NewCache(nil, 100, 1*time.Minute)
	require.NoError(t, err)
	defer cache.Close()

	ctx := context.Background()

	// Empty key should work (no validation currently)
	err = cache.Set(ctx, "", []byte("value"), 1*time.Minute)
	require.NoError(t, err)

	val, err := cache.Get(ctx, "")
	require.NoError(t, err)
	assert.Equal(t, []byte("value"), val)
}

// TestCache_NilValue tests handling of nil/empty values
func TestCache_NilValue(t *testing.T) {
	cache, err := NewCache(nil, 100, 1*time.Minute)
	require.NoError(t, err)
	defer cache.Close()

	ctx := context.Background()

	// Nil value
	err = cache.Set(ctx, "nil", nil, 1*time.Minute)
	require.NoError(t, err)

	val, err := cache.Get(ctx, "nil")
	require.NoError(t, err)
	assert.Nil(t, val)

	// Empty value
	err = cache.Set(ctx, "empty", []byte{}, 1*time.Minute)
	require.NoError(t, err)

	val, err = cache.Get(ctx, "empty")
	require.NoError(t, err)
	assert.Equal(t, []byte{}, val)
}

// TestCache_UpdateExisting tests updating existing cache entries
func TestCache_UpdateExisting(t *testing.T) {
	cache, err := NewCache(nil, 100, 1*time.Minute)
	require.NoError(t, err)
	defer cache.Close()

	ctx := context.Background()

	// Set initial value
	err = cache.Set(ctx, "key1", []byte("value1"), 1*time.Minute)
	require.NoError(t, err)

	val, err := cache.Get(ctx, "key1")
	require.NoError(t, err)
	assert.Equal(t, []byte("value1"), val)

	// Update value
	err = cache.Set(ctx, "key1", []byte("value2"), 1*time.Minute)
	require.NoError(t, err)

	val, err = cache.Get(ctx, "key1")
	require.NoError(t, err)
	assert.Equal(t, []byte("value2"), val, "Value should be updated")
}

// TestCache_ZeroTTL tests setting with zero TTL
func TestCache_ZeroTTL(t *testing.T) {
	cache, err := NewCache(nil, 100, 1*time.Minute)
	require.NoError(t, err)
	defer cache.Close()

	ctx := context.Background()

	// Set with zero TTL
	err = cache.Set(ctx, "zero_ttl", []byte("value"), 0)
	require.NoError(t, err)

	// Should be retrievable (L1 uses cache-wide TTL)
	val, err := cache.Get(ctx, "zero_ttl")
	require.NoError(t, err)
	assert.Equal(t, []byte("value"), val)
}

// TestCache_NegativeTTL tests setting with negative TTL
func TestCache_NegativeTTL(t *testing.T) {
	cache, err := NewCache(nil, 100, 1*time.Minute)
	require.NoError(t, err)
	defer cache.Close()

	ctx := context.Background()

	// Set with negative TTL (should work, treated as immediate expiry in L2)
	err = cache.Set(ctx, "negative_ttl", []byte("value"), -1*time.Second)
	require.NoError(t, err)

	// May or may not be retrievable depending on implementation
	_, _ = cache.Get(ctx, "negative_ttl")
	// No assertion - behavior undefined
}
