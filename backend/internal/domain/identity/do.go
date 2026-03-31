package identity

import "github.com/samber/do/v2"

// RegisterProviders registers identity domain providers with the injector.
func RegisterProviders(i do.Injector) {
	do.Provide[AccountDataManager](i, func(i do.Injector) (AccountDataManager, error) {
		return ProvideAccountDataManagerFromRepository(do.MustInvoke[Repository](i)), nil
	})
	do.Provide[AccountInvitationDataManager](i, func(i do.Injector) (AccountInvitationDataManager, error) {
		return ProvideAccountInvitationDataManagerFromRepository(do.MustInvoke[Repository](i)), nil
	})
	do.Provide[UserDataManager](i, func(i do.Injector) (UserDataManager, error) {
		return ProvideUserDataManagerFromRepository(do.MustInvoke[Repository](i)), nil
	})
	do.Provide[AccountUserMembershipDataManager](i, func(i do.Injector) (AccountUserMembershipDataManager, error) {
		return ProvideAccountUserMembershipDataManagerFromRepository(do.MustInvoke[Repository](i)), nil
	})
}

func ProvideAccountDataManagerFromRepository(r Repository) AccountDataManager {
	return r
}

func ProvideAccountInvitationDataManagerFromRepository(r Repository) AccountInvitationDataManager {
	return r
}

func ProvideUserDataManagerFromRepository(r Repository) UserDataManager {
	return r
}

func ProvideAccountUserMembershipDataManagerFromRepository(r Repository) AccountUserMembershipDataManager {
	return r
}
