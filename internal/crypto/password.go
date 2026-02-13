// Package crypto provides cryptographic utilities for the application
package crypto

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"runtime"

	"github.com/alexedwards/argon2id"
	"golang.org/x/sync/semaphore"
)

// defaultMemoryBudget is the maximum total memory (in bytes) that concurrent
// Argon2id operations are allowed to consume. The semaphore size is computed
// as budget / (params.Memory * 1024).  Using a fixed memory budget instead
// of a CPU multiplier means the limit is safe regardless of core count or
// container memory limit. 512 MiB leaves ample headroom in a 2 GiB container
// for the Go runtime, HTTP connections, and other goroutines under heavy load.
const defaultMemoryBudget = 512 << 20 // 512 MiB

// ErrHasherBusy is returned when the password hasher cannot accept more
// concurrent operations. Callers should translate this into HTTP 503.
var ErrHasherBusy = errors.New("password hasher at capacity")

// maxConcurrentFromParams calculates a safe concurrency limit from the
// Argon2id memory parameter and the memory budget.
func maxConcurrentFromParams(params *argon2id.Params) int {
	memPerOp := uint64(params.Memory) * 1024 // KiB → bytes
	if memPerOp == 0 {
		return runtime.NumCPU()
	}
	n := max(int(defaultMemoryBudget/memPerOp), 1)
	// Also cap to 2×NumCPU — no point queueing more than the CPUs can handle.
	if cpuCap := 2 * runtime.NumCPU(); n > cpuCap {
		n = cpuCap
	}
	return n
}

// PasswordHasher provides password hashing and verification.
// It uses golang.org/x/sync/semaphore with TryAcquire for non-blocking
// load shedding: if all slots are busy, new requests are rejected immediately
// instead of queuing goroutines that consume memory and cause OOM.
type PasswordHasher struct {
	params        *argon2id.Params
	sem           *semaphore.Weighted
	maxConcurrent int64
}

// NewPasswordHasher creates a new password hasher with default Argon2id parameters
// and a memory-aware concurrency limit.
func NewPasswordHasher() *PasswordHasher {
	params := argon2id.DefaultParams
	n := int64(maxConcurrentFromParams(params))
	return &PasswordHasher{
		params:        params,
		sem:           semaphore.NewWeighted(n),
		maxConcurrent: n,
	}
}

// NewPasswordHasherWithParams creates a password hasher with custom parameters
// and a memory-aware concurrency limit.
func NewPasswordHasherWithParams(params *argon2id.Params) *PasswordHasher {
	n := int64(maxConcurrentFromParams(params))
	return &PasswordHasher{
		params:        params,
		sem:           semaphore.NewWeighted(n),
		maxConcurrent: n,
	}
}

// NewPasswordHasherWithConcurrency creates a password hasher with custom parameters
// and a custom concurrency limit.
func NewPasswordHasherWithConcurrency(params *argon2id.Params, maxConcurrent int) *PasswordHasher {
	if maxConcurrent <= 0 {
		maxConcurrent = maxConcurrentFromParams(params)
	}
	n := int64(maxConcurrent)
	return &PasswordHasher{
		params:        params,
		sem:           semaphore.NewWeighted(n),
		maxConcurrent: n,
	}
}

// MaxConcurrent returns the configured concurrency limit. Useful for tests.
func (h *PasswordHasher) MaxConcurrent() int64 {
	return h.maxConcurrent
}

// HashPassword hashes a password using Argon2id.
// Returns the hash in PHC string format: $argon2id$v=19$m=65536,t=3,p=2$...
// Returns ErrHasherBusy immediately if all concurrency slots are in use.
func (h *PasswordHasher) HashPassword(password string) (string, error) {
	return h.HashPasswordContext(context.Background(), password)
}

// HashPasswordContext is like HashPassword but accepts a context for cancellation.
// Uses TryAcquire for immediate rejection when at capacity.
func (h *PasswordHasher) HashPasswordContext(ctx context.Context, password string) (string, error) {
	if password == "" {
		return "", fmt.Errorf("password cannot be empty")
	}

	if !h.sem.TryAcquire(1) {
		return "", ErrHasherBusy
	}
	defer h.sem.Release(1)

	hash, err := argon2id.CreateHash(password, h.params)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	return hash, nil
}

// VerifyPassword verifies a password against an Argon2id hash.
// Returns ErrHasherBusy immediately if all concurrency slots are in use.
func (h *PasswordHasher) VerifyPassword(password, hash string) (bool, error) {
	return h.VerifyPasswordContext(context.Background(), password, hash)
}

// VerifyPasswordContext is like VerifyPassword but accepts a context for cancellation.
// Uses TryAcquire for immediate rejection when at capacity.
func (h *PasswordHasher) VerifyPasswordContext(ctx context.Context, password, hash string) (bool, error) {
	if password == "" {
		return false, fmt.Errorf("password cannot be empty")
	}
	if hash == "" {
		return false, fmt.Errorf("hash cannot be empty")
	}

	if !h.sem.TryAcquire(1) {
		return false, ErrHasherBusy
	}
	defer h.sem.Release(1)

	match, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		return false, fmt.Errorf("argon2id verification failed: %w", err)
	}

	return match, nil
}

// GenerateSecureToken generates a cryptographically secure random token
// Returns a hex-encoded string of the specified byte length
func GenerateSecureToken(byteLength int) (string, error) {
	if byteLength <= 0 {
		return "", fmt.Errorf("byte length must be positive")
	}

	bytes := make([]byte, byteLength)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	return hex.EncodeToString(bytes), nil
}
