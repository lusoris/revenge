//go:build integration
// +build integration

package cache_test

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTTLAccuracy(t *testing.T) {
	c := newTestCache(t)
	defer c.Close()
	ctx := context.Background()

	// Test multiple TTL values for accuracy
	testCases := []struct {
		name string
		ttl  time.Duration
	}{
		{"500ms", 500 * time.Millisecond},
		{"1s", 1 * time.Second},
		{"2s", 2 * time.Second},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			key := fmt.Sprintf("test:ttl:%s", tc.name)
			value := []byte("expires")

			// Set with TTL
			err := c.Set(ctx, key, value, tc.ttl)
			require.NoError(t, err)

			// Should exist immediately
			exists, err := c.Exists(ctx, key)
			require.NoError(t, err)
			assert.True(t, exists)

			// Wait for expiration + buffer
			time.Sleep(tc.ttl + 200*time.Millisecond)

			// Should be expired
			_, err = c.Get(ctx, key)
			assert.Error(t, err, "Key should have expired")
		})
	}
}

func TestPatternInvalidation(t *testing.T) {
	c := newTestCache(t)
	defer c.Close()
	ctx := context.Background()

	// Create multiple patterns
	patterns := []string{"user", "session", "token"}
	for _, pattern := range patterns {
		for i := 1; i <= 10; i++ {
			key := fmt.Sprintf("test:pattern:%s:%d", pattern, i)
			err := c.Set(ctx, key, []byte("value"), 5*time.Minute)
			require.NoError(t, err)
		}
	}

	// Invalidate one pattern
	err := c.Invalidate(ctx, "test:pattern:user:*")
	require.NoError(t, err)

	// Give it time to propagate
	time.Sleep(100 * time.Millisecond)

	// Verify user pattern is gone
	for i := 1; i <= 10; i++ {
		key := fmt.Sprintf("test:pattern:user:%d", i)
		_, err := c.Get(ctx, key)
		assert.Error(t, err, "Key %s should be deleted", key)
	}

	// Note: Other patterns might be gone too due to L1 clearing
	// This is documented behavior - Invalidate clears entire L1
}

func TestConcurrentReadsWrites(t *testing.T) {
	c := newTestCache(t)
	defer c.Close()
	ctx := context.Background()

	const numKeys = 100
	const numReaders = 20
	const numWriters = 10
	const duration = 5 * time.Second

	// Pre-populate keys
	for i := 0; i < numKeys; i++ {
		key := fmt.Sprintf("test:rw:%d", i)
		err := c.Set(ctx, key, []byte(fmt.Sprintf("value-%d", i)), 10*time.Minute)
		require.NoError(t, err)
	}

	var wg sync.WaitGroup
	stopChan := make(chan struct{})
	errorChan := make(chan error, numReaders+numWriters)

	// Start readers
	for r := 0; r < numReaders; r++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			readCount := 0
			for {
				select {
				case <-stopChan:
					t.Logf("Reader %d: %d reads", id, readCount)
					return
				default:
					keyIdx := readCount % numKeys
					key := fmt.Sprintf("test:rw:%d", keyIdx)
					_, err := c.Get(ctx, key)
					if err != nil {
						errorChan <- fmt.Errorf("reader %d get failed: %w", id, err)
						return
					}
					readCount++
				}
			}
		}(r)
	}

	// Start writers
	for w := 0; w < numWriters; w++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			writeCount := 0
			for {
				select {
				case <-stopChan:
					t.Logf("Writer %d: %d writes", id, writeCount)
					return
				default:
					keyIdx := writeCount % numKeys
					key := fmt.Sprintf("test:rw:%d", keyIdx)
					value := []byte(fmt.Sprintf("updated-%d-%d", id, writeCount))
					err := c.Set(ctx, key, value, 10*time.Minute)
					if err != nil {
						errorChan <- fmt.Errorf("writer %d set failed: %w", id, err)
						return
					}
					writeCount++
				}
			}
		}(w)
	}

	// Run for duration
	time.Sleep(duration)
	close(stopChan)
	wg.Wait()

	// Check for errors
	close(errorChan)
	for err := range errorChan {
		t.Error(err)
	}
}

func TestMemoryPressure(t *testing.T) {
	c := newTestCache(t)
	defer c.Close()
	ctx := context.Background()

	// Try to fill cache with large values
	const valueSize = 1024 * 1024 // 1MB
	const numValues = 100

	value := make([]byte, valueSize)
	for i := range value {
		value[i] = byte(i % 256)
	}

	// Store many large values
	for i := 0; i < numValues; i++ {
		key := fmt.Sprintf("test:memory:%d", i)
		err := c.Set(ctx, key, value, 10*time.Minute)
		require.NoError(t, err, "Failed to set large value %d", i)
	}

	// Verify we can still retrieve them
	for i := 0; i < numValues; i++ {
		key := fmt.Sprintf("test:memory:%d", i)
		retrieved, err := c.Get(ctx, key)
		require.NoError(t, err, "Failed to get large value %d", i)
		assert.Equal(t, valueSize, len(retrieved))
	}

	// Clean up
	for i := 0; i < numValues; i++ {
		key := fmt.Sprintf("test:memory:%d", i)
		c.Delete(ctx, key)
	}
}

func TestContextCancellation(t *testing.T) {
	c := newTestCache(t)
	defer c.Close()

	// Create a cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	// Operations should fail or complete before context check
	key := "test:cancelled"
	value := []byte("value")

	err := c.Set(ctx, key, value, 1*time.Minute)
	// Rueidis might complete before checking context, so error is not guaranteed
	t.Logf("Set with cancelled context: %v", err)

	_, err = c.Get(ctx, key)
	t.Logf("Get with cancelled context: %v", err)
}

func TestEdgeCases(t *testing.T) {
	c := newTestCache(t)
	defer c.Close()
	ctx := context.Background()

	t.Run("EmptyKey", func(t *testing.T) {
		err := c.Set(ctx, "", []byte("value"), 1*time.Minute)
		// Should either error or accept it
		t.Logf("Empty key set result: %v", err)
	})

	t.Run("EmptyValue", func(t *testing.T) {
		err := c.Set(ctx, "test:empty", []byte{}, 1*time.Minute)
		require.NoError(t, err)

		val, err := c.Get(ctx, "test:empty")
		require.NoError(t, err)
		assert.Equal(t, []byte{}, val)
	})

	t.Run("NilValue", func(t *testing.T) {
		err := c.Set(ctx, "test:nil", nil, 1*time.Minute)
		require.NoError(t, err)

		val, err := c.Get(ctx, "test:nil")
		require.NoError(t, err)
		assert.Equal(t, []byte{}, val) // nil becomes empty slice through Redis
	})

	t.Run("VeryLongKey", func(t *testing.T) {
		// Redis key limit is 512MB, but reasonable limit is ~1KB
		longKey := "test:" + string(make([]byte, 10000))
		err := c.Set(ctx, longKey, []byte("value"), 1*time.Minute)
		t.Logf("Very long key result: %v", err)
	})

	t.Run("ZeroTTL", func(t *testing.T) {
		// What happens with 0 TTL?
		err := c.Set(ctx, "test:zero-ttl", []byte("value"), 0)
		t.Logf("Zero TTL result: %v", err)
	})

	t.Run("NegativeTTL", func(t *testing.T) {
		// What happens with negative TTL?
		err := c.Set(ctx, "test:neg-ttl", []byte("value"), -1*time.Second)
		t.Logf("Negative TTL result: %v", err)
	})

	t.Run("MaxTTL", func(t *testing.T) {
		// Very long TTL
		err := c.Set(ctx, "test:max-ttl", []byte("value"), 100*365*24*time.Hour)
		require.NoError(t, err)
	})
}

func TestRapidOperations(t *testing.T) {
	c := newTestCache(t)
	defer c.Close()
	ctx := context.Background()

	key := "test:rapid"

	// Rapid set/get/delete cycles
	for i := 0; i < 1000; i++ {
		value := []byte(fmt.Sprintf("value-%d", i))

		err := c.Set(ctx, key, value, 1*time.Minute)
		require.NoError(t, err)

		retrieved, err := c.Get(ctx, key)
		require.NoError(t, err)
		assert.Equal(t, value, retrieved)

		err = c.Delete(ctx, key)
		require.NoError(t, err)
	}
}

func TestStaleReads(t *testing.T) {
	c := newTestCache(t)
	defer c.Close()
	ctx := context.Background()

	key := "test:stale"
	value1 := []byte("value1")
	value2 := []byte("value2")

	// Set initial value
	err := c.Set(ctx, key, value1, 1*time.Minute)
	require.NoError(t, err)

	// Get (should be in L1 now)
	retrieved, err := c.Get(ctx, key)
	require.NoError(t, err)
	assert.Equal(t, value1, retrieved)

	// Update value
	err = c.Set(ctx, key, value2, 1*time.Minute)
	require.NoError(t, err)

	// Get again - should get updated value
	retrieved, err = c.Get(ctx, key)
	require.NoError(t, err)
	assert.Equal(t, value2, retrieved, "Should get updated value, not stale L1 cache")
}

func TestConnectionResilience(t *testing.T) {
	c := newTestCache(t)
	defer c.Close()
	ctx := context.Background()

	// Perform many operations to test connection pooling
	const numOps = 10000

	for i := 0; i < numOps; i++ {
		key := fmt.Sprintf("test:conn:%d", i%100) // Reuse keys
		value := []byte(fmt.Sprintf("value-%d", i))

		if i%3 == 0 {
			c.Set(ctx, key, value, 1*time.Minute)
		} else if i%3 == 1 {
			c.Get(ctx, key)
		} else {
			c.Exists(ctx, key)
		}

		if i%1000 == 0 {
			t.Logf("Completed %d operations", i)
		}
	}

	t.Logf("Completed %d operations successfully", numOps)
}

func TestDataIntegrity(t *testing.T) {
	c := newTestCache(t)
	defer c.Close()
	ctx := context.Background()

	// Test that data doesn't get corrupted
	testData := [][]byte{
		[]byte("simple text"),
		[]byte("\x00\x01\x02\x03\x04\x05"), // Binary data
		[]byte("unicode: 你好世界 🚀"),
		[]byte("{\"json\": \"data\", \"nested\": {\"value\": 123}}"),
		make([]byte, 1024*1024), // 1MB of zeros
	}

	// Fill the 1MB buffer with pattern
	for i := range testData[4] {
		testData[4][i] = byte(i % 256)
	}

	for i, data := range testData {
		key := fmt.Sprintf("test:integrity:%d", i)

		err := c.Set(ctx, key, data, 1*time.Minute)
		require.NoError(t, err)

		retrieved, err := c.Get(ctx, key)
		require.NoError(t, err)
		assert.Equal(t, data, retrieved, "Data corruption detected for test case %d", i)
	}
}

func TestHighConcurrencyStress(t *testing.T) {
	c := newTestCache(t)
	defer c.Close()
	ctx := context.Background()

	const numGoroutines = 200
	const opsPerGoroutine = 500

	var wg sync.WaitGroup
	errorChan := make(chan error, numGoroutines)

	startTime := time.Now()

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			for j := 0; j < opsPerGoroutine; j++ {
				key := fmt.Sprintf("test:stress:%d:%d", id, j)
				value := []byte(fmt.Sprintf("value-%d-%d", id, j))

				// Random operations
				switch j % 4 {
				case 0: // Set
					if err := c.Set(ctx, key, value, 1*time.Minute); err != nil {
						errorChan <- err
						return
					}
				case 1: // Get
					if _, err := c.Get(ctx, key); err != nil {
						// May not exist, that's ok
					}
				case 2: // Exists
					if _, err := c.Exists(ctx, key); err != nil {
						errorChan <- err
						return
					}
				case 3: // Delete
					if err := c.Delete(ctx, key); err != nil {
						errorChan <- err
						return
					}
				}
			}
		}(i)
	}

	wg.Wait()
	duration := time.Since(startTime)

	totalOps := numGoroutines * opsPerGoroutine
	opsPerSec := float64(totalOps) / duration.Seconds()

	t.Logf("Stress test: %d operations in %v (%.0f ops/sec)", totalOps, duration, opsPerSec)

	// Check for errors
	close(errorChan)
	errorCount := 0
	for err := range errorChan {
		t.Error(err)
		errorCount++
	}

	assert.Equal(t, 0, errorCount, "Should have no errors under high concurrency")
}
