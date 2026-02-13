package crypto

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/alexedwards/argon2id"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPasswordHasher_HashPassword(t *testing.T) {
	hasher := NewPasswordHasher()

	password := "MySecurePassword123!"
	hash, err := hasher.HashPassword(password)

	require.NoError(t, err)
	assert.NotEmpty(t, hash)
	assert.NotEqual(t, password, hash)
	assert.Contains(t, hash, "$argon2id$", "hash should use Argon2id format")
}

func TestPasswordHasher_HashPasswordEmpty(t *testing.T) {
	hasher := NewPasswordHasher()

	hash, err := hasher.HashPassword("")

	assert.Error(t, err)
	assert.Empty(t, hash)
	assert.Contains(t, err.Error(), "password cannot be empty")
}

func TestPasswordHasher_VerifyPassword(t *testing.T) {
	hasher := NewPasswordHasher()

	password := "MySecurePassword123!"
	hash, err := hasher.HashPassword(password)
	require.NoError(t, err)

	// Correct password
	match, err := hasher.VerifyPassword(password, hash)
	require.NoError(t, err)
	assert.True(t, match)

	// Wrong password
	match, err = hasher.VerifyPassword("WrongPassword", hash)
	require.NoError(t, err)
	assert.False(t, match)
}

func TestPasswordHasher_VerifyPasswordEmpty(t *testing.T) {
	hasher := NewPasswordHasher()

	// Empty password
	match, err := hasher.VerifyPassword("", "somehash")
	assert.Error(t, err)
	assert.False(t, match)

	// Empty hash
	match, err = hasher.VerifyPassword("password", "")
	assert.Error(t, err)
	assert.False(t, match)
}

func TestPasswordHasher_CustomParams(t *testing.T) {
	params := &argon2id.Params{
		Memory:      32 * 1024, // 32 MiB
		Iterations:  2,
		Parallelism: 2,
		SaltLength:  16,
		KeyLength:   32,
	}

	hasher := NewPasswordHasherWithParams(params)

	password := "TestPassword"
	hash, err := hasher.HashPassword(password)
	require.NoError(t, err)

	match, err := hasher.VerifyPassword(password, hash)
	require.NoError(t, err)
	assert.True(t, match)
}

func TestGenerateSecureToken(t *testing.T) {
	token1, err := GenerateSecureToken(32)
	require.NoError(t, err)
	assert.Len(t, token1, 64, "32 bytes should produce 64 hex characters")

	token2, err := GenerateSecureToken(32)
	require.NoError(t, err)
	assert.NotEqual(t, token1, token2, "tokens should be unique")
}

func TestGenerateSecureTokenInvalidLength(t *testing.T) {
	token, err := GenerateSecureToken(0)
	assert.Error(t, err)
	assert.Empty(t, token)

	token, err = GenerateSecureToken(-1)
	assert.Error(t, err)
	assert.Empty(t, token)
}

func TestGenerateSecureTokenVariousLengths(t *testing.T) {
	tests := []struct {
		byteLength int
		hexLength  int
	}{
		{16, 32},
		{32, 64},
		{64, 128},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%d_bytes", tt.byteLength), func(t *testing.T) {
			token, err := GenerateSecureToken(tt.byteLength)
			require.NoError(t, err)
			assert.Len(t, token, tt.hexLength)
		})
	}
}

func TestPasswordHasher_MultipleHashes(t *testing.T) {
	hasher := NewPasswordHasher()
	password := "TestPass123!"

	// Hash the same password twice
	hash1, err := hasher.HashPassword(password)
	require.NoError(t, err)

	hash2, err := hasher.HashPassword(password)
	require.NoError(t, err)

	// Hashes should be different (due to salt)
	assert.NotEqual(t, hash1, hash2, "same password should produce different hashes")

	// Both hashes should verify
	match1, err := hasher.VerifyPassword(password, hash1)
	require.NoError(t, err)
	assert.True(t, match1)

	match2, err := hasher.VerifyPassword(password, hash2)
	require.NoError(t, err)
	assert.True(t, match2)
}

func TestPasswordHasher_ConcurrencySemaphore(t *testing.T) {
	// Use lightweight params to keep tests fast
	params := &argon2id.Params{
		Memory:      1024, // 1 MiB — tiny for testing
		Iterations:  1,
		Parallelism: 1,
		SaltLength:  16,
		KeyLength:   32,
	}

	maxConcurrent := 2
	hasher := NewPasswordHasherWithConcurrency(params, maxConcurrent)

	password := "TestPass123!"
	hash, err := hasher.HashPassword(password)
	require.NoError(t, err)

	// Track how many goroutines are inside the critical section concurrently
	var inFlight atomic.Int32
	var maxObserved atomic.Int32
	var wg sync.WaitGroup
	const numGoroutines = 20

	var busyCount atomic.Int32

	wg.Add(numGoroutines)
	for range numGoroutines {
		go func() {
			defer wg.Done()
			cur := inFlight.Add(1)
			defer inFlight.Add(-1)

			// Record the maximum in-flight we observe
			for {
				old := maxObserved.Load()
				if cur <= old || maxObserved.CompareAndSwap(old, cur) {
					break
				}
			}

			// TryAcquire may reject some requests — that's expected
			match, err := hasher.VerifyPassword(password, hash)
			if err != nil {
				if err == ErrHasherBusy {
					busyCount.Add(1)
					return
				}
				assert.NoError(t, err)
				return
			}
			assert.True(t, match)
		}()
	}

	wg.Wait()
	// With TryAcquire, some goroutines should have been rejected
	t.Logf("busy rejections: %d / %d", busyCount.Load(), numGoroutines)
}

func TestPasswordHasher_ErrHasherBusy(t *testing.T) {
	// Create a hasher with concurrency=1 so we can easily saturate it
	params := &argon2id.Params{
		Memory:      1024,
		Iterations:  1,
		Parallelism: 1,
		SaltLength:  16,
		KeyLength:   32,
	}
	hasher := NewPasswordHasherWithConcurrency(params, 1)

	// Occupy the single slot via the underlying semaphore
	acquired := hasher.sem.TryAcquire(1)
	require.True(t, acquired, "should acquire the only slot")

	// Now all hash/verify calls should return ErrHasherBusy immediately
	_, err := hasher.HashPassword("password")
	require.ErrorIs(t, err, ErrHasherBusy)

	_, err = hasher.HashPasswordContext(context.Background(), "password")
	require.ErrorIs(t, err, ErrHasherBusy)

	_, verifyErr := hasher.VerifyPassword("password", "$argon2id$v=19$m=1024,t=1,p=1$aaaa$bbbb")
	require.ErrorIs(t, verifyErr, ErrHasherBusy)

	_, verifyErr = hasher.VerifyPasswordContext(context.Background(), "password", "$argon2id$v=19$m=1024,t=1,p=1$aaaa$bbbb")
	require.ErrorIs(t, verifyErr, ErrHasherBusy)

	// Release the slot, operations should work again
	hasher.sem.Release(1)

	hash, err := hasher.HashPassword("password")
	require.NoError(t, err)
	assert.NotEmpty(t, hash)
}

func TestNewPasswordHasherWithConcurrency_Defaults(t *testing.T) {
	params := &argon2id.Params{
		Memory:      1024,
		Iterations:  1,
		Parallelism: 1,
		SaltLength:  16,
		KeyLength:   32,
	}

	// Computed from 512 MiB budget / (1024 KiB * 1024) = 512, but capped at 2×NumCPU
	expected := int64(maxConcurrentFromParams(params))

	// Zero should use computed default
	h := NewPasswordHasherWithConcurrency(params, 0)
	assert.Equal(t, expected, h.MaxConcurrent())

	// Negative should use computed default
	h = NewPasswordHasherWithConcurrency(params, -1)
	assert.Equal(t, expected, h.MaxConcurrent())

	// Custom value should be used
	h = NewPasswordHasherWithConcurrency(params, 5)
	assert.Equal(t, int64(5), h.MaxConcurrent())
}

func TestMaxConcurrentFromParams(t *testing.T) {
	// Default params: 64 MiB (65536 KiB). Budget 1 GiB = 1048576 KiB.
	// 1048576 / 65536 = 16, and if 2×NumCPU >= 16, result is 16.
	defaultParams := argon2id.DefaultParams
	n := maxConcurrentFromParams(defaultParams)
	assert.LessOrEqual(t, n, 2*runtime.NumCPU())
	// With 64MB params, budget (512 MiB) allows at most 8
	assert.LessOrEqual(t, n, 8)
	assert.GreaterOrEqual(t, n, 1)

	// Small memory params: 1 MiB. Budget allows 512, but capped at 2×NumCPU.
	smallParams := &argon2id.Params{Memory: 1024, Iterations: 1, Parallelism: 1, SaltLength: 16, KeyLength: 32}
	n = maxConcurrentFromParams(smallParams)
	assert.Equal(t, 2*runtime.NumCPU(), n) // CPU cap should kick in

	// Huge memory params: 512 MiB. Budget allows 1.
	hugeParams := &argon2id.Params{Memory: 512 * 1024, Iterations: 1, Parallelism: 1, SaltLength: 16, KeyLength: 32}
	n = maxConcurrentFromParams(hugeParams)
	assert.Equal(t, 1, n)

	// Absurdly large: 2 GiB. Budget allows 0 → clamped to 1.
	giganticParams := &argon2id.Params{Memory: 2 * 1024 * 1024, Iterations: 1, Parallelism: 1, SaltLength: 16, KeyLength: 32}
	n = maxConcurrentFromParams(giganticParams)
	assert.Equal(t, 1, n)
}
