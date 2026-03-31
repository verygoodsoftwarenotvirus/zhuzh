package authorization

const (
	// UpdateUserStatusPermission is a service admin permission.
	UpdateUserStatusPermission Permission = "update.user_status"
	// ImpersonateUserPermission is a service admin permission.
	ImpersonateUserPermission Permission = "imitate.user"
	// ReadUserPermission is a service admin permission.
	ReadUserPermission Permission = "read.user"
	// SearchUserPermission is a service admin permission.
	SearchUserPermission Permission = "search.user"
	// ArchiveUserPermission is a service admin permission.
	ArchiveUserPermission Permission = "archive.user"
	// ManageUserSessionsPermission is a service admin permission.
	ManageUserSessionsPermission Permission = "manage.user_sessions"
)

var (
	// AuthPermissions contains all authentication-related permissions.
	AuthPermissions = []Permission{
		UpdateUserStatusPermission,
		ImpersonateUserPermission,
		ReadUserPermission,
		SearchUserPermission,
		ArchiveUserPermission,
		ManageUserSessionsPermission,
	}
)
