package authorization

import (
	"encoding/gob"
	"slices"
)

func init() {
	gob.Register(serviceRoleCollection{})
}

const (
	serviceAdminRoleName     = "service_admin"
	serviceDataAdminRoleName = "service_data_admin"
	serviceUserRoleName      = "service_user"

	// ServiceUserRoleID is the database ID for the service_user role.
	ServiceUserRoleID = "role_service_user"
	// ServiceAdminRoleID is the database ID for the service_admin role.
	ServiceAdminRoleID = "role_service_admin"
	// ServiceDataAdminRoleID is the database ID for the service_data_admin role.
	ServiceDataAdminRoleID = "role_service_data_admin"

	invalidServiceRoleWarning = "INVALID_SERVICE_ROLE"

	// invalidServiceRole is a service role to apply for non-admin users to have one.
	invalidServiceRole ServiceRole = iota
	// ServiceUserRole is a service role to apply for non-admin users to have one.
	ServiceUserRole ServiceRole = iota
	// ServiceAdminRole is a role that allows a user to do basically anything.
	ServiceAdminRole ServiceRole = iota
)

type (
	// ServiceRole describes a role a user has for the Service context.
	ServiceRole role

	// ServiceRolePermissionChecker checks permissions for one or more service Roles.
	ServiceRolePermissionChecker interface {
		HasPermission(Permission) bool

		AsAccountRolePermissionChecker() AccountRolePermissionsChecker
		IsServiceAdmin() bool
		CanUpdateUserAccountStatuses() bool
		CanImpersonateUsers() bool
		CanManageUserSessions() bool
	}

	serviceRoleCollection struct {
		Permissions map[Permission]bool
		RoleNames   []string
	}
)

func (r ServiceRole) String() string {
	switch r {
	case invalidServiceRole:
		return invalidServiceRoleWarning
	case ServiceUserRole:
		return serviceUserRoleName
	case ServiceAdminRole:
		return serviceAdminRoleName
	default:
		return ""
	}
}

// NewServiceRolePermissionChecker returns a new checker from role names and a set of permissions.
func NewServiceRolePermissionChecker(roleNames []string, perms []Permission) ServiceRolePermissionChecker {
	m := make(map[Permission]bool, len(perms))
	for _, p := range perms {
		m[p] = true
	}
	return &serviceRoleCollection{
		Permissions: m,
		RoleNames:   roleNames,
	}
}

func (r serviceRoleCollection) AsAccountRolePermissionChecker() AccountRolePermissionsChecker {
	perms := make([]Permission, 0, len(r.Permissions))
	for p := range r.Permissions {
		perms = append(perms, p)
	}
	return NewAccountRolePermissionChecker(perms)
}

// HasPermission returns whether a user can do something or not.
func (r serviceRoleCollection) HasPermission(p Permission) bool {
	return r.Permissions[p]
}

// IsServiceAdmin returns if a role is an admin.
func (r serviceRoleCollection) IsServiceAdmin() bool {
	return slices.Contains(r.RoleNames, serviceAdminRoleName)
}

// CanUpdateUserAccountStatuses returns whether a user can update user account statuses.
func (r serviceRoleCollection) CanUpdateUserAccountStatuses() bool {
	return r.Permissions[UpdateUserStatusPermission]
}

// CanImpersonateUsers returns whether a user can impersonate others.
func (r serviceRoleCollection) CanImpersonateUsers() bool {
	return r.Permissions[ImpersonateUserPermission]
}

// CanManageUserSessions returns whether a user can manage other users' sessions.
func (r serviceRoleCollection) CanManageUserSessions() bool {
	return r.Permissions[ManageUserSessionsPermission]
}
