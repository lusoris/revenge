package middleware

import (
	"fmt"
	"log/slog"
	"testing"
)

func BenchmarkRateLimiter_GetLimiter_New(b *testing.B) {
	rl := NewRateLimiter(DefaultRateLimitConfig(), slog.Default())
	defer rl.Stop()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rl.getLimiter(fmt.Sprintf("192.168.1.%d", i%256))
	}
}

func BenchmarkRateLimiter_GetLimiter_Existing(b *testing.B) {
	rl := NewRateLimiter(DefaultRateLimitConfig(), slog.Default())
	defer rl.Stop()

	rl.getLimiter("192.168.1.1")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rl.getLimiter("192.168.1.1")
	}
}

func BenchmarkRateLimiter_Allow(b *testing.B) {
	rl := NewRateLimiter(RateLimitConfig{
		Enabled:           true,
		RequestsPerSecond: 1000000,
		Burst:             1000000,
	}, slog.Default())
	defer rl.Stop()

	limiter := rl.getLimiter("192.168.1.1")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		limiter.Allow()
	}
}

func BenchmarkRateLimiter_ShouldLimit_AllOps(b *testing.B) {
	rl := NewRateLimiter(DefaultRateLimitConfig(), slog.Default())
	defer rl.Stop()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rl.shouldLimit("SearchMoviesMetadata")
	}
}

func BenchmarkRateLimiter_ShouldLimit_SpecificOps(b *testing.B) {
	rl := NewRateLimiter(AuthRateLimitConfig(), slog.Default())
	defer rl.Stop()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rl.shouldLimit("Login")
	}
}

func BenchmarkRateLimiter_Concurrent(b *testing.B) {
	rl := NewRateLimiter(RateLimitConfig{
		Enabled:           true,
		RequestsPerSecond: 1000000,
		Burst:             1000000,
	}, slog.Default())
	defer rl.Stop()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			ip := fmt.Sprintf("10.0.%d.%d", (i/256)%256, i%256)
			limiter := rl.getLimiter(ip)
			limiter.Allow()
			i++
		}
	})
}
