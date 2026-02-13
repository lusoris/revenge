package cache

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"
)

func BenchmarkL1Cache_Set(b *testing.B) {
	c, err := NewL1Cache[string, []byte](10000, 5*time.Minute)
	if err != nil {
		b.Fatal(err)
	}
	defer c.Close()

	value := []byte("benchmark-value-data")
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		c.Set(fmt.Sprintf("key-%d", i), value)
	}
}

func BenchmarkL1Cache_Get_Hit(b *testing.B) {
	c, err := NewL1Cache[string, []byte](10000, 5*time.Minute)
	if err != nil {
		b.Fatal(err)
	}
	defer c.Close()

	value := []byte("benchmark-value-data")
	c.Set("hit-key", value)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		c.Get("hit-key")
	}
}

func BenchmarkL1Cache_Get_Miss(b *testing.B) {
	c, err := NewL1Cache[string, []byte](10000, 5*time.Minute)
	if err != nil {
		b.Fatal(err)
	}
	defer c.Close()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		c.Get("miss-key")
	}
}

func BenchmarkL1Cache_Concurrent(b *testing.B) {
	c, err := NewL1Cache[string, []byte](10000, 5*time.Minute)
	if err != nil {
		b.Fatal(err)
	}
	defer c.Close()

	value := []byte("benchmark-value-data")
	// Pre-populate
	for i := range 1000 {
		c.Set(fmt.Sprintf("key-%d", i), value)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := fmt.Sprintf("key-%d", i%1000)
			if i%3 == 0 {
				c.Set(key, value)
			} else {
				c.Get(key)
			}
			i++
		}
	})
}

func BenchmarkCache_SetGet_L1Only(b *testing.B) {
	// Cache with nil L2 client — pure L1 path
	cache, err := NewNamedCache(nil, 10000, 5*time.Minute, "bench")
	if err != nil {
		b.Fatal(err)
	}
	defer cache.Close()

	ctx := context.Background()
	value := []byte(`{"title":"Benchmark Movie","year":2024}`)

	b.Run("Set", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = cache.Set(ctx, fmt.Sprintf("movie:%d", i), value, 5*time.Minute)
		}
	})

	b.Run("Get_Hit", func(b *testing.B) {
		_ = cache.Set(ctx, "movie:hit", value, 5*time.Minute)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = cache.Get(ctx, "movie:hit")
		}
	})

	b.Run("Concurrent_Mixed", func(b *testing.B) {
		// Pre-populate
		for i := range 500 {
			_ = cache.Set(ctx, fmt.Sprintf("movie:%d", i), value, 5*time.Minute)
		}
		b.ResetTimer()

		var wg sync.WaitGroup
		b.RunParallel(func(pb *testing.PB) {
			wg.Add(1)
			defer wg.Done()
			i := 0
			for pb.Next() {
				key := fmt.Sprintf("movie:%d", i%500)
				if i%4 == 0 {
					_ = cache.Set(ctx, key, value, 5*time.Minute)
				} else {
					_, _ = cache.Get(ctx, key)
				}
				i++
			}
		})
		wg.Wait()
	})
}
