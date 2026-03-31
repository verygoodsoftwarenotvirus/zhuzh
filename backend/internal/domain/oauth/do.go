package oauth

import "github.com/samber/do/v2"

// RegisterProviders registers OAuth domain providers with the injector.
func RegisterProviders(i do.Injector) {
	do.Provide[OAuth2ClientDataManager](i, func(i do.Injector) (OAuth2ClientDataManager, error) {
		return ProvideOAuth2ClientDataManagerFromRepository(do.MustInvoke[Repository](i)), nil
	})
	do.Provide[OAuth2ClientTokenDataManager](i, func(i do.Injector) (OAuth2ClientTokenDataManager, error) {
		return ProvideOAuth2ClientTokenDataManagerFromRepository(do.MustInvoke[Repository](i)), nil
	})
}

func ProvideOAuth2ClientDataManagerFromRepository(r Repository) OAuth2ClientDataManager {
	return r
}

func ProvideOAuth2ClientTokenDataManagerFromRepository(r Repository) OAuth2ClientTokenDataManager {
	return r
}
