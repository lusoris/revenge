package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSessionKey(t *testing.T) {
	key := SessionKey("abc123")
	assert.Equal(t, "session:abc123", key)
}

func TestSessionByUserKey(t *testing.T) {
	key := SessionByUserKey("user-uuid")
	assert.Equal(t, "session:user:user-uuid", key)
}

func TestRBACEnforceKey(t *testing.T) {
	key := RBACEnforceKey("user1", "resource", "read")
	assert.Equal(t, "rbac:enforce:user1:resource:read", key)
}

func TestRBACUserRolesKey(t *testing.T) {
	key := RBACUserRolesKey("user-uuid")
	assert.Equal(t, "rbac:roles:user-uuid", key)
}

func TestRBACUserPermsKey(t *testing.T) {
	key := RBACUserPermsKey("user-uuid")
	assert.Equal(t, "rbac:perms:user-uuid", key)
}

func TestServerSettingKey(t *testing.T) {
	key := ServerSettingKey("theme.mode")
	assert.Equal(t, "settings:server:theme.mode", key)
}

func TestUserSettingKey(t *testing.T) {
	key := UserSettingKey("user-uuid", "notifications.enabled")
	assert.Equal(t, "settings:user:user-uuid:notifications.enabled", key)
}

func TestUserKey(t *testing.T) {
	key := UserKey("user-uuid")
	assert.Equal(t, "user:user-uuid", key)
}

func TestUserByNameKey(t *testing.T) {
	key := UserByNameKey("johndoe")
	assert.Equal(t, "user:name:johndoe", key)
}

func TestMovieKey(t *testing.T) {
	key := MovieKey("movie-uuid")
	assert.Equal(t, "movie:movie-uuid", key)
}

func TestMovieMetaKey(t *testing.T) {
	key := MovieMetaKey("tmdb", "12345")
	assert.Equal(t, "movie:meta:tmdb:12345", key)
}

func TestDefaultTTLs(t *testing.T) {
	// Verify TTLs are reasonable values
	assert.Equal(t, 30*time.Second, SessionTTL)
	assert.Equal(t, 5*time.Minute, RBACPolicyTTL)
	assert.Equal(t, 30*time.Second, RBACEnforceTTL)
	assert.Equal(t, 5*time.Minute, ServerSettingsTTL)
	assert.Equal(t, 2*time.Minute, UserSettingsTTL)
	assert.Equal(t, 1*time.Minute, UserTTL)
	assert.Equal(t, 10*time.Minute, MovieMetaTTL)
}

func TestKeyPrefixes(t *testing.T) {
	// Verify key prefixes are correct
	assert.Equal(t, "session:", KeyPrefixSession)
	assert.Equal(t, "session:user:", KeyPrefixSessionByUser)
	assert.Equal(t, "rbac:policy:", KeyPrefixRBACPolicy)
	assert.Equal(t, "rbac:enforce:", KeyPrefixRBACEnforce)
	assert.Equal(t, "rbac:roles:", KeyPrefixRBACUserRoles)
	assert.Equal(t, "rbac:perms:", KeyPrefixRBACUserPerms)
	assert.Equal(t, "settings:server:", KeyPrefixServerSetting)
	assert.Equal(t, "settings:user:", KeyPrefixUserSetting)
	assert.Equal(t, "user:", KeyPrefixUser)
	assert.Equal(t, "user:name:", KeyPrefixUserByName)
	assert.Equal(t, "user:email:", KeyPrefixUserEmail)
	assert.Equal(t, "movie:", KeyPrefixMovie)
	assert.Equal(t, "movie:meta:", KeyPrefixMovieMeta)
}
