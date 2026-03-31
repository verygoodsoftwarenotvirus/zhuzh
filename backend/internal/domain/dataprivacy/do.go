package dataprivacy

import "github.com/samber/do/v2"

// RegisterProviders registers data privacy domain providers with the injector.
func RegisterProviders(i do.Injector) {
	do.Provide[DataPrivacyDataManager](i, func(i do.Injector) (DataPrivacyDataManager, error) {
		return ProvideDataPrivacyDataManagerFromRepository(do.MustInvoke[Repository](i)), nil
	})
}

func ProvideDataPrivacyDataManagerFromRepository(r Repository) DataPrivacyDataManager {
	return r
}
