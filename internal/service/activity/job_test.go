package activity

import (
	"testing"
	"time"

	infrajobs "github.com/lusoris/revenge/internal/infra/jobs"
	"github.com/stretchr/testify/assert"
)

func TestActivityLogArgs_Kind(t *testing.T) {
	t.Parallel()
	args := ActivityLogArgs{}
	assert.Equal(t, "activity_log", args.Kind())
}

func TestActivityLogArgs_InsertOpts(t *testing.T) {
	t.Parallel()
	args := ActivityLogArgs{}
	opts := args.InsertOpts()
	assert.Equal(t, infrajobs.QueueLow, opts.Queue)
	assert.Equal(t, 3, opts.MaxAttempts)
}

func TestIsSecurityAction(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		action   string
		expected bool
	}{
		{"login is security", ActionUserLogin, true},
		{"logout is security", ActionUserLogout, true},
		{"user create is security", ActionUserCreate, true},
		{"user delete is security", ActionUserDelete, true},
		{"password reset is security", ActionUserPasswordReset, true},
		{"session create is security", ActionSessionCreate, true},
		{"session revoke is security", ActionSessionRevoke, true},
		{"oidc login is security", ActionOIDCLogin, true},
		{"oidc link is security", ActionOIDCLink, true},
		{"admin role assign is security", ActionAdminRoleAssign, true},
		{"admin user ban is security", ActionAdminUserBan, true},
		{"failed suffix is security", "user.login.failed", true},
		{"library create is not security", ActionLibraryCreate, false},
		{"library scan is not security", ActionLibraryScan, false},
		{"settings update is not security", ActionSettingsUpdate, false},
		{"empty action is not security", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expected, isSecurityAction(tt.action))
		})
	}
}

func TestActivityLogWorker_Timeout(t *testing.T) {
	t.Parallel()
	worker := &ActivityLogWorker{}
	timeout := worker.Timeout(nil)
	assert.Equal(t, 10*time.Second, timeout)
}
