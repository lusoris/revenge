package mfa

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateRandomCode(t *testing.T) {
	// Generate multiple codes
	codes := make(map[string]bool)
	for i := 0; i < 100; i++ {
		code, err := generateRandomCode()
		require.NoError(t, err)
		
		// Should be hex string
		assert.Len(t, code, 16, "code should be 16 characters (8 bytes hex encoded)")
		assert.Regexp(t, `^[0-9a-f]{16}$`, code, "code should be lowercase hex")
		
		// Should be unique
		assert.False(t, codes[code], "code should be unique")
		codes[code] = true
	}
}

func TestFormatCode(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "16 character code",
			input:    "1234567890abcdef",
			expected: "1234-5678-90ab-cdef",
		},
		{
			name:     "short code (no formatting)",
			input:    "12345",
			expected: "12345",
		},
		{
			name:     "empty code",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatCode(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNormalizeCode(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "formatted code with dashes",
			input:    "1234-5678-90AB-CDEF",
			expected: "1234567890abcdef",
		},
		{
			name:     "code with spaces",
			input:    "1234 5678 90ab cdef",
			expected: "1234567890abcdef",
		},
		{
			name:     "mixed case",
			input:    "AbCdEfGh12345678",
			expected: "abcdefgh12345678",
		},
		{
			name:     "already normalized",
			input:    "1234567890abcdef",
			expected: "1234567890abcdef",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizeCode(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestConstantTimeCompare(t *testing.T) {
	tests := []struct {
		name     string
		a        string
		b        string
		expected bool
	}{
		{
			name:     "identical strings",
			a:        "1234567890abcdef",
			b:        "1234567890abcdef",
			expected: true,
		},
		{
			name:     "different strings",
			a:        "1234567890abcdef",
			b:        "1234567890abcdeg",
			expected: false,
		},
		{
			name:     "different lengths",
			a:        "12345",
			b:        "1234567890",
			expected: false,
		},
		{
			name:     "empty strings",
			a:        "",
			b:        "",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ConstantTimeCompare(tt.a, tt.b)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestBackupCodeLength(t *testing.T) {
	// Verify constants are correct
	assert.Equal(t, 10, BackupCodeCount, "should generate 10 backup codes")
	assert.Equal(t, 8, BackupCodeLength, "code should be 8 bytes (16 hex chars)")
	assert.Equal(t, 12, BackupCodeCost, "bcrypt cost should be 12")
}

// Integration tests would go here
// These would require actual database connection

func TestBackupCodesService_GenerateCodes(t *testing.T) {
	t.Skip("Integration test - requires database")

	// This test would cover:
	// 1. GenerateCodes creates exactly 10 codes
	// 2. All codes are unique
	// 3. All codes are properly formatted (XXXX-XXXX-XXXX-XXXX)
	// 4. All codes are stored as bcrypt hashes in database
	// 5. Plain text codes are returned
}

func TestBackupCodesService_VerifyCode(t *testing.T) {
	t.Skip("Integration test - requires database")

	// This test would cover:
	// 1. Valid code returns true
	// 2. Invalid code returns false
	// 3. Code is marked as used after successful verification
	// 4. Used code cannot be verified again
	// 5. Client IP is stored
	// 6. Normalized codes work (with/without dashes, different case)
}

func TestBackupCodesService_RegenerateCodes(t *testing.T) {
	t.Skip("Integration test - requires database")

	// This test would cover:
	// 1. Old codes are deleted
	// 2. New codes are generated
	// 3. Old codes cannot be used after regeneration
	// 4. New codes can be used
}

func TestBackupCodesService_GetRemainingCount(t *testing.T) {
	t.Skip("Integration test - requires database")

	// This test would cover:
	// 1. Returns 10 for newly generated codes
	// 2. Decrements after each use
	// 3. Returns 0 when all codes are used
}

func TestBackupCodesService_TimingAttack(t *testing.T) {
	t.Skip("Security test - requires careful timing analysis")

	// This test would verify:
	// 1. Verification time is constant regardless of:
	//    - Number of codes
	//    - Position of matching code
	//    - Whether code matches or not
	// 2. No information leak about code validity through timing
}
