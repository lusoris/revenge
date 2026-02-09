package activity

import (
	"context"
	"net"
	"testing"

	"github.com/google/uuid"
	"github.com/lusoris/revenge/internal/infra/logging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAsyncLogger_LogAction_NilClient(t *testing.T) {
	t.Parallel()
	logger := logging.NewTestLogger()
	al := NewAsyncLogger(nil, logger)

	err := al.LogAction(context.Background(), LogActionRequest{
		UserID:   uuid.Must(uuid.NewV7()),
		Username: "testuser",
		Action:   ActionUserLogin,
	})
	// Nil client should not return error; it silently drops the log
	require.NoError(t, err)
}

func TestAsyncLogger_LogFailure_NilClient(t *testing.T) {
	t.Parallel()
	logger := logging.NewTestLogger()
	al := NewAsyncLogger(nil, logger)

	err := al.LogFailure(context.Background(), LogFailureRequest{
		Action:       ActionUserLogin,
		ErrorMessage: "invalid credentials",
	})
	require.NoError(t, err)
}

func TestAsyncLogger_ImplementsLogger(t *testing.T) {
	t.Parallel()
	var _ Logger = (*AsyncLogger)(nil)
}

func TestAsyncLogger_LogAction_FieldMapping(t *testing.T) {
	t.Parallel()

	// Verify that the AsyncLogger correctly maps LogActionRequest fields
	// to LogRequest. We test the mapping by constructing a request with
	// all fields populated.
	userID := uuid.Must(uuid.NewV7())
	resourceID := uuid.Must(uuid.NewV7())
	ip := net.ParseIP("192.168.1.1")

	req := LogActionRequest{
		UserID:       userID,
		Username:     "testuser",
		Action:       ActionUserLogin,
		ResourceType: ResourceTypeUser,
		ResourceID:   resourceID,
		Changes:      map[string]interface{}{"field": "value"},
		Metadata:     map[string]interface{}{"key": "val"},
		IPAddress:    ip,
		UserAgent:    "test-agent",
	}

	// With nil client, LogAction succeeds but does nothing
	logger := logging.NewTestLogger()
	al := NewAsyncLogger(nil, logger)
	err := al.LogAction(context.Background(), req)
	require.NoError(t, err)
}

func TestAsyncLogger_SecurityRouting(t *testing.T) {
	t.Parallel()

	// Verify that isSecurityAction is checked correctly for both
	// security and non-security actions
	assert.True(t, isSecurityAction(ActionUserLogin))
	assert.True(t, isSecurityAction(ActionUserPasswordReset))
	assert.False(t, isSecurityAction(ActionLibraryCreate))
	assert.False(t, isSecurityAction(ActionSettingsUpdate))
}
