package auth

type Repository interface {
	PasswordResetTokenDataManager
	UserSessionDataManager
}
