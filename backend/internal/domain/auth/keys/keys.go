package keys

const (
	idSuffix = ".id"

	// PasswordResetTokenKey is the standard key for referring to a password reset token's ID.
	PasswordResetTokenKey = "password_reset_token"
	// PasswordResetTokenIDKey is the standard key for referring to a password reset token's ID.
	PasswordResetTokenIDKey = PasswordResetTokenKey + idSuffix

	// UserSessionKey is the standard key for referring to a user session.
	UserSessionKey = "user_session"
	// UserSessionIDKey is the standard key for referring to a user session's ID.
	UserSessionIDKey = UserSessionKey + idSuffix
)
