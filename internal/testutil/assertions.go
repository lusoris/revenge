package testutil

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// AssertValidUUID asserts that the given string is a valid UUID.
func AssertValidUUID(t *testing.T, s string, msgAndArgs ...any) {
	t.Helper()
	_, err := uuid.Parse(s)
	assert.NoError(t, err, msgAndArgs...)
}

// RequireValidUUID requires that the given string is a valid UUID.
func RequireValidUUID(t *testing.T, s string, msgAndArgs ...any) {
	t.Helper()
	_, err := uuid.Parse(s)
	require.NoError(t, err, msgAndArgs...)
}

// AssertTimeEqual asserts that two times are equal within a tolerance.
// Default tolerance is 1 second.
func AssertTimeEqual(t *testing.T, expected, actual time.Time, msgAndArgs ...any) {
	t.Helper()
	AssertTimeEqualWithTolerance(t, expected, actual, time.Second, msgAndArgs...)
}

// AssertTimeEqualWithTolerance asserts that two times are equal within a given tolerance.
func AssertTimeEqualWithTolerance(t *testing.T, expected, actual time.Time, tolerance time.Duration, msgAndArgs ...any) {
	t.Helper()
	diff := expected.Sub(actual)
	if diff < 0 {
		diff = -diff
	}
	assert.LessOrEqual(t, diff, tolerance, msgAndArgs...)
}

// RequireTimeEqual requires that two times are equal within a tolerance.
// Default tolerance is 1 second.
func RequireTimeEqual(t *testing.T, expected, actual time.Time, msgAndArgs ...any) {
	t.Helper()
	RequireTimeEqualWithTolerance(t, expected, actual, time.Second, msgAndArgs...)
}

// RequireTimeEqualWithTolerance requires that two times are equal within a given tolerance.
func RequireTimeEqualWithTolerance(t *testing.T, expected, actual time.Time, tolerance time.Duration, msgAndArgs ...any) {
	t.Helper()
	diff := expected.Sub(actual)
	if diff < 0 {
		diff = -diff
	}
	require.LessOrEqual(t, diff, tolerance, msgAndArgs...)
}

// AssertRecentTime asserts that the given time is recent (within the last minute).
func AssertRecentTime(t *testing.T, tm time.Time, msgAndArgs ...any) {
	t.Helper()
	AssertRecentTimeWithTolerance(t, tm, time.Minute, msgAndArgs...)
}

// AssertRecentTimeWithTolerance asserts that the given time is recent within a tolerance.
func AssertRecentTimeWithTolerance(t *testing.T, tm time.Time, tolerance time.Duration, msgAndArgs ...any) {
	t.Helper()
	now := time.Now()
	diff := now.Sub(tm)
	assert.GreaterOrEqual(t, diff, time.Duration(0), append([]any{"time should not be in the future"}, msgAndArgs...)...)
	assert.LessOrEqual(t, diff, tolerance, append([]any{"time should be recent"}, msgAndArgs...)...)
}

// RequireRecentTime requires that the given time is recent (within the last minute).
func RequireRecentTime(t *testing.T, tm time.Time, msgAndArgs ...any) {
	t.Helper()
	RequireRecentTimeWithTolerance(t, tm, time.Minute, msgAndArgs...)
}

// RequireRecentTimeWithTolerance requires that the given time is recent within a tolerance.
func RequireRecentTimeWithTolerance(t *testing.T, tm time.Time, tolerance time.Duration, msgAndArgs ...any) {
	t.Helper()
	now := time.Now()
	diff := now.Sub(tm)
	require.GreaterOrEqual(t, diff, time.Duration(0), append([]any{"time should not be in the future"}, msgAndArgs...)...)
	require.LessOrEqual(t, diff, tolerance, append([]any{"time should be recent"}, msgAndArgs...)...)
}

// AssertNotZeroUUID asserts that the UUID is not zero/nil.
func AssertNotZeroUUID(t *testing.T, id uuid.UUID, msgAndArgs ...any) {
	t.Helper()
	assert.NotEqual(t, uuid.Nil, id, msgAndArgs...)
}

// RequireNotZeroUUID requires that the UUID is not zero/nil.
func RequireNotZeroUUID(t *testing.T, id uuid.UUID, msgAndArgs ...any) {
	t.Helper()
	require.NotEqual(t, uuid.Nil, id, msgAndArgs...)
}

// AssertSliceContains asserts that a slice contains the given element.
func AssertSliceContains[T comparable](t *testing.T, slice []T, element T, msgAndArgs ...any) {
	t.Helper()
	assert.Contains(t, slice, element, msgAndArgs...)
}

// RequireSliceContains requires that a slice contains the given element.
func RequireSliceContains[T comparable](t *testing.T, slice []T, element T, msgAndArgs ...any) {
	t.Helper()
	require.Contains(t, slice, element, msgAndArgs...)
}
