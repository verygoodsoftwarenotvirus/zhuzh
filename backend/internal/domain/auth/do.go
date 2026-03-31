package auth

import "github.com/samber/do/v2"

// ProvidePasswordResetTokenDataManagerFromRepository provides a PasswordResetTokenDataManager from a Repository.
func ProvidePasswordResetTokenDataManagerFromRepository(r Repository) PasswordResetTokenDataManager {
	return r
}

// RegisterProviders registers auth domain providers with the injector.
func RegisterProviders(i do.Injector) {
	do.Provide[PasswordResetTokenDataManager](i, func(i do.Injector) (PasswordResetTokenDataManager, error) {
		return ProvidePasswordResetTokenDataManagerFromRepository(do.MustInvoke[Repository](i)), nil
	})
}
