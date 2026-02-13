package rbac

// Fine-grained permission constants for the RBAC system.
// Permissions follow the format: resource:action

// User permissions
const (
	PermUsersList   = "users:list"
	PermUsersGet    = "users:get"
	PermUsersCreate = "users:create"
	PermUsersUpdate = "users:update"
	PermUsersDelete = "users:delete"
)

// Profile permissions (own profile)
const (
	PermProfileRead   = "profile:read"
	PermProfileUpdate = "profile:update"
)

// Movie permissions
const (
	PermMoviesList   = "movies:list"
	PermMoviesGet    = "movies:get"
	PermMoviesCreate = "movies:create"
	PermMoviesUpdate = "movies:update"
	PermMoviesDelete = "movies:delete"
)

// Library permissions
const (
	PermLibrariesList   = "libraries:list"
	PermLibrariesGet    = "libraries:get"
	PermLibrariesCreate = "libraries:create"
	PermLibrariesUpdate = "libraries:update"
	PermLibrariesDelete = "libraries:delete"
	PermLibrariesScan   = "libraries:scan"
)

// Playback permissions
const (
	PermPlaybackStream   = "playback:stream"
	PermPlaybackProgress = "playback:progress"
)

// Request permissions (media requests)
const (
	PermRequestsList    = "requests:list"
	PermRequestsGet     = "requests:get"
	PermRequestsCreate  = "requests:create"
	PermRequestsApprove = "requests:approve"
	PermRequestsDelete  = "requests:delete"
)

// Settings permissions
const (
	PermSettingsRead      = "settings:read"
	PermSettingsWrite     = "settings:write"
	PermSettingsUserRead  = "settings:user_read"
	PermSettingsUserWrite = "settings:user_write"
)

// Audit permissions
const (
	PermAuditRead   = "audit:read"
	PermAuditExport = "audit:export"
)

// Integration permissions
const (
	PermIntegrationsList   = "integrations:list"
	PermIntegrationsGet    = "integrations:get"
	PermIntegrationsCreate = "integrations:create"
	PermIntegrationsUpdate = "integrations:update"
	PermIntegrationsDelete = "integrations:delete"
	PermIntegrationsSync   = "integrations:sync"
)

// Notification permissions
const (
	PermNotificationsList   = "notifications:list"
	PermNotificationsGet    = "notifications:get"
	PermNotificationsCreate = "notifications:create"
	PermNotificationsUpdate = "notifications:update"
	PermNotificationsDelete = "notifications:delete"
)

// Admin wildcard permission
const (
	PermAdminAll = "admin:*"
)

// FineGrainedResources defines the resource taxonomy.
var FineGrainedResources = []string{
	"users",
	"profile",
	"movies",
	"libraries",
	"playback",
	"requests",
	"settings",
	"audit",
	"integrations",
	"notifications",
	"admin",
}

// FineGrainedActions defines all fine-grained actions.
var FineGrainedActions = []string{
	"list",
	"get",
	"create",
	"update",
	"delete",
	"read",
	"write",
	"stream",
	"progress",
	"approve",
	"export",
	"sync",
	"scan",
	"user_read",
	"user_write",
	"*",
}

// DefaultRolePermissions defines the default permissions for each built-in role.
var DefaultRolePermissions = map[string][]string{
	"admin": {
		PermAdminAll, // Full access
	},
	"moderator": {
		// User management (limited)
		PermUsersList, PermUsersGet,
		// Own profile
		PermProfileRead, PermProfileUpdate,
		// Movies - full access
		PermMoviesList, PermMoviesGet, PermMoviesCreate, PermMoviesUpdate, PermMoviesDelete,
		// Libraries - full access
		PermLibrariesList, PermLibrariesGet, PermLibrariesCreate, PermLibrariesUpdate, PermLibrariesDelete, PermLibrariesScan,
		// Playback
		PermPlaybackStream, PermPlaybackProgress,
		// Requests - full access
		PermRequestsList, PermRequestsGet, PermRequestsCreate, PermRequestsApprove, PermRequestsDelete,
		// Settings - read only server, full user
		PermSettingsRead, PermSettingsUserRead, PermSettingsUserWrite,
		// Audit - read
		PermAuditRead,
		// Integrations - full access
		PermIntegrationsList, PermIntegrationsGet, PermIntegrationsCreate, PermIntegrationsUpdate, PermIntegrationsDelete, PermIntegrationsSync,
		// Notifications - full access
		PermNotificationsList, PermNotificationsGet, PermNotificationsCreate, PermNotificationsUpdate, PermNotificationsDelete,
	},
	"user": {
		// Own profile
		PermProfileRead, PermProfileUpdate,
		// Movies - list and view
		PermMoviesList, PermMoviesGet,
		// Libraries - list and view
		PermLibrariesList, PermLibrariesGet,
		// Playback
		PermPlaybackStream, PermPlaybackProgress,
		// Requests - can create and view own
		PermRequestsList, PermRequestsGet, PermRequestsCreate,
		// Settings - own settings only
		PermSettingsUserRead, PermSettingsUserWrite,
		// Notifications - own only
		PermNotificationsList, PermNotificationsGet,
	},
	"guest": {
		// Profile read only
		PermProfileRead,
		// Movies - list and view only
		PermMoviesList, PermMoviesGet,
		// Libraries - list and view only
		PermLibrariesList, PermLibrariesGet,
		// Playback - stream only (no progress tracking)
		PermPlaybackStream,
	},
}

// HasPermission checks if a permission string is in a list of permissions.
func HasPermission(permissions []string, perm string) bool {
	for _, p := range permissions {
		if p == perm || p == PermAdminAll {
			return true
		}
	}
	return false
}
