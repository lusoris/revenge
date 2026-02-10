//go:build integration
// +build integration

package cache_test

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/lusoris/revenge/internal/config"
	"github.com/lusoris/revenge/internal/infra/cache"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestCache(t *testing.T) *cache.Cache {
	t.Helper()
	cfg := &config.Config{
		Cache: config.CacheConfig{
			Enabled: true,
			URL:     "redis://localhost:6379/0",
		},
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelWarn}))
	client, err := cache.NewClient(cfg, logger)
	if err != nil {
		t.Skipf("Redis not available: %v", err)
	}
	c, err := cache.NewCache(client, 1000, 5*time.Minute)
	require.NoError(t, err)
	return c
}

func TestBasicOperations(t *testing.T) {
	c := newTestCache(t)
	defer c.Close()
	ctx := context.Background()

	err := c.Set(ctx, "test:key", []byte("value"), 1*time.Minute)
	require.NoError(t, err)

	val, err := c.Get(ctx, "test:key")
	require.NoError(t, err)
	assert.Equal(t, []byte("value"), val)

	exists, err := c.Exists(ctx, "test:key")
	require.NoError(t, err)
	assert.True(t, exists)

	err = c.Delete(ctx, "test:key")
	require.NoError(t, err)

	exists, err = c.Exists(ctx, "test:key")
	require.NoError(t, err)
	assert.False(t, exists)
}

func TestConcurrency(t *testing.T) {
	c := newTestCache(t)
	defer c.Close()
	ctx := context.Background()

	done := make(chan bool, 50)
	for i := 0; i < 50; i++ {
		go func(id int) {
			for j := 0; j < 100; j++ {
				key := fmt.Sprintf("test:%d:%d", id, j)
				c.Set(ctx, key, []byte("val"), 1*time.Minute)
				c.Get(ctx, key)
				c.Delete(ctx, key)
			}
			done <- true
		}(i)
	}

	for i := 0; i < 50; i++ {
		<-done
	}
}

func TestLargeValue(t *testing.T) {
	c := newTestCache(t)
	defer c.Close()
	ctx := context.Background()

	value := make([]byte, 5*1024*1024) // 5MB
	for i := range value {
		value[i] = byte(i % 256)
	}

	err := c.Set(ctx, "test:large", value, 1*time.Minute)
	require.NoError(t, err)

	retrieved, err := c.Get(ctx, "test:large")
	require.NoError(t, err)
	assert.Equal(t, value, retrieved)

	c.Delete(ctx, "test:large")
}
