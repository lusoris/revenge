package mfa

import (
	"context"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"

	db "github.com/lusoris/revenge/internal/infra/database/db"
	"github.com/lusoris/revenge/internal/testutil"
)

func setupBackupCodesService(t *testing.T) (*BackupCodesService, *db.Queries) {
	t.Helper()
	testDB := testutil.NewTestDB(t)
	queries := db.New(testDB.Pool())
	logger := zaptest.NewLogger(t)
	svc := NewBackupCodesService(queries, logger)
	return svc, queries
}

func createTestUser(t *testing.T, queries *db.Queries, ctx context.Context) uuid.UUID {
	t.Helper()
	user, err := queries.CreateUser(ctx, db.CreateUserParams{
		Username:     "testuser_" + uuid.New().String()[:8],
		Email:        "test_" + uuid.New().String()[:8] + "@example.com",
		PasswordHash: "$argon2id$v=19$m=65536,t=3,p=2$test",
	})
	require.NoError(t, err)
	return user.ID
}

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
}

func TestBackupCodesService_GenerateCodes(t *testing.T) {
	t.Parallel()
	svc, queries := setupBackupCodesService(t)

	ctx := context.Background()
	userID := createTestUser(t, queries, ctx)

	codes, err := svc.GenerateCodes(ctx, userID)
	require.NoError(t, err)

	// Should generate exactly 10 codes
	assert.Len(t, codes, 10)

	// All codes should be properly formatted
	for _, code := range codes {
		assert.Len(t, code, 19, "formatted code should be 19 chars (XXXX-XXXX-XXXX-XXXX)")
		assert.Regexp(t, `^[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}$`, code)
	}

	// All codes should be unique
	codeSet := make(map[string]bool)
	for _, code := range codes {
		assert.False(t, codeSet[code], "codes should be unique")
		codeSet[code] = true
	}

	// Verify codes are stored in database
	count, err := svc.GetRemainingCount(ctx, userID)
	require.NoError(t, err)
	assert.Equal(t, int64(10), count)
}

func TestBackupCodesService_VerifyCode(t *testing.T) {
	t.Parallel()
	svc, queries := setupBackupCodesService(t)

	ctx := context.Background()
	userID := createTestUser(t, queries, ctx)

	// Generate codes
	codes, err := svc.GenerateCodes(ctx, userID)
	require.NoError(t, err)
	require.Len(t, codes, 10)

	testCode := codes[0]
	clientIP := "192.168.1.1"

	t.Run("valid code", func(t *testing.T) {
		valid, err := svc.VerifyCode(ctx, userID, testCode, clientIP)
		require.NoError(t, err)
		assert.True(t, valid)

		// Count should decrement
		count, err := svc.GetRemainingCount(ctx, userID)
		require.NoError(t, err)
		assert.Equal(t, int64(9), count)
	})

	t.Run("already used code", func(t *testing.T) {
		// Try to use the same code again
		valid, err := svc.VerifyCode(ctx, userID, testCode, clientIP)
		require.NoError(t, err)
		assert.False(t, valid)
	})

	t.Run("normalized code formats", func(t *testing.T) {
		testCode2 := codes[1]
		// Remove dashes
		codeWithoutDashes := normalizeCode(testCode2)

		valid, err := svc.VerifyCode(ctx, userID, codeWithoutDashes, clientIP)
		require.NoError(t, err)
		assert.True(t, valid)

		// Try uppercase
		testCode3 := codes[2]
		upperCode := strings.ToUpper(testCode3)

		valid, err = svc.VerifyCode(ctx, userID, upperCode, clientIP)
		require.NoError(t, err)
		assert.True(t, valid)
	})

	t.Run("invalid code", func(t *testing.T) {
		valid, err := svc.VerifyCode(ctx, userID, "0000-0000-0000-0000", clientIP)
		require.NoError(t, err)
		assert.False(t, valid)
	})

	t.Run("no codes available", func(t *testing.T) {
		newUserID := uuid.New()
		valid, err := svc.VerifyCode(ctx, newUserID, testCode, clientIP)
		require.Error(t, err)
		assert.False(t, valid)
		assert.Equal(t, ErrNoBackupCodes, err)
	})
}

func TestBackupCodesService_RegenerateCodes(t *testing.T) {
	t.Parallel()
	svc, queries := setupBackupCodesService(t)

	ctx := context.Background()
	userID := createTestUser(t, queries, ctx)

	// Generate initial codes
	oldCodes, err := svc.GenerateCodes(ctx, userID)
	require.NoError(t, err)
	require.Len(t, oldCodes, 10)

	// Regenerate codes
	newCodes, err := svc.RegenerateCodes(ctx, userID)
	require.NoError(t, err)
	require.Len(t, newCodes, 10)

	// Old codes should not work
	clientIP := "192.168.1.1"
	valid, err := svc.VerifyCode(ctx, userID, oldCodes[0], clientIP)
	require.NoError(t, err)
	assert.False(t, valid)

	// New codes should work
	valid, err = svc.VerifyCode(ctx, userID, newCodes[0], clientIP)
	require.NoError(t, err)
	assert.True(t, valid)

	// Count should still be 9 (used 1 of 10 new codes)
	count, err := svc.GetRemainingCount(ctx, userID)
	require.NoError(t, err)
	assert.Equal(t, int64(9), count)
}

func TestBackupCodesService_GetRemainingCount(t *testing.T) {
	t.Parallel()
	svc, queries := setupBackupCodesService(t)

	ctx := context.Background()
	userID := createTestUser(t, queries, ctx)

	// No codes yet
	count, err := svc.GetRemainingCount(ctx, userID)
	require.NoError(t, err)
	assert.Equal(t, int64(0), count)

	// Generate codes
	codes, err := svc.GenerateCodes(ctx, userID)
	require.NoError(t, err)

	// Should be 10
	count, err = svc.GetRemainingCount(ctx, userID)
	require.NoError(t, err)
	assert.Equal(t, int64(10), count)

	// Use one code
	clientIP := "192.168.1.1"
	valid, err := svc.VerifyCode(ctx, userID, codes[0], clientIP)
	require.NoError(t, err)
	require.True(t, valid)

	// Should be 9
	count, err = svc.GetRemainingCount(ctx, userID)
	require.NoError(t, err)
	assert.Equal(t, int64(9), count)
}

func TestBackupCodesService_HasBackupCodes(t *testing.T) {
	t.Parallel()
	svc, queries := setupBackupCodesService(t)

	ctx := context.Background()
	userID := createTestUser(t, queries, ctx)

	// No codes yet
	has, err := svc.HasBackupCodes(ctx, userID)
	require.NoError(t, err)
	assert.False(t, has)

	// Generate codes
	_, err = svc.GenerateCodes(ctx, userID)
	require.NoError(t, err)

	// Should have codes
	has, err = svc.HasBackupCodes(ctx, userID)
	require.NoError(t, err)
	assert.True(t, has)
}

func TestBackupCodesService_DeleteAllCodes(t *testing.T) {
	t.Parallel()
	svc, queries := setupBackupCodesService(t)

	ctx := context.Background()
	userID := createTestUser(t, queries, ctx)

	// Generate codes
	_, err := svc.GenerateCodes(ctx, userID)
	require.NoError(t, err)

	// Verify they exist
	has, err := svc.HasBackupCodes(ctx, userID)
	require.NoError(t, err)
	require.True(t, has)

	// Delete all codes
	err = svc.DeleteAllCodes(ctx, userID)
	require.NoError(t, err)

	// Should have no codes
	has, err = svc.HasBackupCodes(ctx, userID)
	require.NoError(t, err)
	assert.False(t, has)

	count, err := svc.GetRemainingCount(ctx, userID)
	require.NoError(t, err)
	assert.Equal(t, int64(0), count)
}

func TestNewBackupCodesService(t *testing.T) {
	testDB := testutil.NewTestDB(t)
	queries := db.New(testDB.Pool())
	logger := zaptest.NewLogger(t)

	svc := NewBackupCodesService(queries, logger)

	assert.NotNil(t, svc)
	assert.NotNil(t, svc.queries)
	assert.NotNil(t, svc.hasher)
	assert.NotNil(t, svc.logger)
}
