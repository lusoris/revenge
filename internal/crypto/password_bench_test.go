package crypto_test

import (
	"testing"

	"github.com/alexedwards/argon2id"
	"github.com/lusoris/revenge/internal/crypto"
)

// BenchmarkHashPassword measures the performance of password hashing
// This is critical as hashing is intentionally slow for security
func BenchmarkHashPassword(b *testing.B) {
	hasher := crypto.NewPasswordHasher()

	password := "TestPassword123!SecureAndLong"

	// Reset timer to exclude setup time
	b.ResetTimer()

	for b.Loop() {
		_, err := hasher.HashPassword(password)
		if err != nil {
			b.Fatalf("HashPassword failed: %v", err)
		}
	}
}

// BenchmarkHashPasswordParallel tests concurrent hashing performance
// Useful for understanding multi-user registration load
func BenchmarkHashPasswordParallel(b *testing.B) {
	hasher := crypto.NewPasswordHasher()

	password := "TestPassword123!SecureAndLong"

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := hasher.HashPassword(password)
			if err != nil {
				b.Fatalf("HashPassword failed: %v", err)
			}
		}
	})
}

// BenchmarkVerifyPassword measures password verification speed
// This is the hot path for login operations
func BenchmarkVerifyPassword(b *testing.B) {
	hasher := crypto.NewPasswordHasher()

	password := "TestPassword123!SecureAndLong"
	hash, err := hasher.HashPassword(password)
	if err != nil {
		b.Fatalf("Failed to hash password: %v", err)
	}

	b.ResetTimer()

	for b.Loop() {
		match, err := hasher.VerifyPassword(password, hash)
		if err != nil {
			b.Fatalf("VerifyPassword failed: %v", err)
		}
		if !match {
			b.Fatal("Password verification failed")
		}
	}
}

// BenchmarkVerifyPasswordParallel tests concurrent verification
// Critical for high-traffic login scenarios
func BenchmarkVerifyPasswordParallel(b *testing.B) {
	hasher := crypto.NewPasswordHasher()

	password := "TestPassword123!SecureAndLong"
	hash, err := hasher.HashPassword(password)
	if err != nil {
		b.Fatalf("Failed to hash password: %v", err)
	}

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			match, err := hasher.VerifyPassword(password, hash)
			if err != nil {
				b.Fatalf("VerifyPassword failed: %v", err)
			}
			if !match {
				b.Fatal("Password verification failed")
			}
		}
	})
}

// BenchmarkGenerateSecureToken measures token generation speed
// Used for session tokens, password reset tokens, etc.
func BenchmarkGenerateSecureToken(b *testing.B) {
	for b.Loop() {
		_, err := crypto.GenerateSecureToken(32)
		if err != nil {
			b.Fatalf("GenerateSecureToken failed: %v", err)
		}
	}
}

// BenchmarkGenerateSecureTokenParallel tests concurrent token generation
func BenchmarkGenerateSecureTokenParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := crypto.GenerateSecureToken(32)
			if err != nil {
				b.Fatalf("GenerateSecureToken failed: %v", err)
			}
		}
	})
}

// BenchmarkPasswordHasherWithCustomParams tests performance impact of different parameters
func BenchmarkPasswordHasherWithCustomParams(b *testing.B) {
	tests := []struct {
		name   string
		memory uint32
		time   uint32
	}{
		{"Low", 32 * 1024, 1},    // 32 MB, 1 iteration
		{"Medium", 64 * 1024, 3}, // 64 MB, 3 iterations (default)
		{"High", 128 * 1024, 5},  // 128 MB, 5 iterations
	}

	password := "TestPassword123!SecureAndLong"

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			hasher := crypto.NewPasswordHasherWithParams(&argon2id.Params{
				Memory:      tt.memory,
				Iterations:  tt.time,
				Parallelism: 2,
				SaltLength:  16,
				KeyLength:   32,
			})

			b.ResetTimer()

			for b.Loop() {
				_, err := hasher.HashPassword(password)
				if err != nil {
					b.Fatalf("HashPassword failed: %v", err)
				}
			}
		})
	}
}
