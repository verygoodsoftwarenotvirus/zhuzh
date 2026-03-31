package authorization

import (
	"encoding/gob"
	"slices"
)

type (
	// AccountRole describes a role a user has for an account context.
	AccountRole role

	// AccountRolePermissionsChecker checks permissions for one or more account Roles.
	AccountRolePermissionsChecker interface {
		HasPermission(Permission) bool
	}
)

const (
	// AccountMemberRole is a role for a plain account participant.
	AccountMemberRole AccountRole = iota
	// AccountAdminRole is a role for someone who can manipulate the specifics of an account.
	AccountAdminRole AccountRole = iota

	AccountAdminRoleName  = "account_admin"
	AccountMemberRoleName = "account_member"

	// AccountAdminRoleID is the database ID for the account_admin role.
	AccountAdminRoleID = "role_account_admin"
	// AccountMemberRoleID is the database ID for the account_member role.
	AccountMemberRoleID = "role_account_member"
)

type accountRoleCollection struct {
	Permissions map[Permission]bool
	RoleNames   []string
}

func init() {
	gob.Register(accountRoleCollection{})
}

// NewAccountRolePermissionChecker returns a new checker from a set of permissions.
func NewAccountRolePermissionChecker(perms []Permission) AccountRolePermissionsChecker {
	m := make(map[Permission]bool, len(perms))
	for _, p := range perms {
		m[p] = true
	}
	return &accountRoleCollection{
		Permissions: m,
	}
}

func (r AccountRole) String() string {
	switch r {
	case AccountMemberRole:
		return AccountMemberRoleName
	case AccountAdminRole:
		return AccountAdminRoleName
	default:
		return ""
	}
}

// HasPermission returns whether a user can do something or not.
func (r accountRoleCollection) HasPermission(p Permission) bool {
	return r.Permissions[p]
}

// IsAccountAdmin returns whether a user is an account admin.
func (r accountRoleCollection) IsAccountAdmin() bool {
	return slices.Contains(r.RoleNames, AccountAdminRoleName)
}
