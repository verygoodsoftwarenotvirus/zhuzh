package authorization

const (
	// CreateUserNotificationsPermission is an admin user permission.
	CreateUserNotificationsPermission Permission = "create.user_notifications"
	// ReadUserNotificationsPermission is a permission.
	ReadUserNotificationsPermission Permission = "read.user_notifications"
	// UpdateUserNotificationsPermission is a permission.
	UpdateUserNotificationsPermission Permission = "update.user_notifications"
	// #nosec G101 CreateUserDeviceTokensPermission is a permission for registering device tokens.
	CreateUserDeviceTokensPermission Permission = "create.user_device_tokens"
	// #nosec G101 ReadUserDeviceTokensPermission is a permission for reading device tokens.
	ReadUserDeviceTokensPermission Permission = "read.user_device_tokens"
	// #nosec G101 ArchiveUserDeviceTokensPermission is a permission for archiving device tokens.
	ArchiveUserDeviceTokensPermission Permission = "archive.user_device_tokens"
)

var (
	// NotificationsPermissions contains all notification-related permissions.
	NotificationsPermissions = []Permission{
		CreateUserNotificationsPermission,
		ReadUserNotificationsPermission,
		UpdateUserNotificationsPermission,
		CreateUserDeviceTokensPermission,
		ReadUserDeviceTokensPermission,
		ArchiveUserDeviceTokensPermission,
	}
)
