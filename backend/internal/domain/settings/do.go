package settings

import "github.com/samber/do/v2"

// RegisterProviders registers settings domain providers with the injector.
func RegisterProviders(i do.Injector) {
	do.Provide[ServiceSettingDataManager](i, func(i do.Injector) (ServiceSettingDataManager, error) {
		return ProvideServiceSettingDataManagerFromRepository(do.MustInvoke[Repository](i)), nil
	})
	do.Provide[ServiceSettingConfigurationDataManager](i, func(i do.Injector) (ServiceSettingConfigurationDataManager, error) {
		return ProvideServiceSettingConfigurationDataManagerFromRepository(do.MustInvoke[Repository](i)), nil
	})
}

func ProvideServiceSettingDataManagerFromRepository(r Repository) ServiceSettingDataManager {
	return r
}

func ProvideServiceSettingConfigurationDataManagerFromRepository(r Repository) ServiceSettingConfigurationDataManager {
	return r
}
