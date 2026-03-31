package identity

import (
	"testing"

	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/authentication/sessions"
	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/authorization"

	"github.com/verygoodsoftwarenotvirus/platform/v4/observability/logging"

	"github.com/stretchr/testify/assert"
)

func TestChangeActiveAccountInput_ValidateWithContext(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		x := &ChangeActiveAccountInput{
			AccountID: "123",
		}

		assert.NoError(t, x.ValidateWithContext(ctx))
	})
}

func TestSessionContextData_AttachToLogger(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &sessions.ContextData{
			Requester: sessions.RequesterInfo{ServicePermissions: authorization.NewServiceRolePermissionChecker([]string{t.Name()}, nil)},
		}

		assert.NotNil(t, x.AttachToLogger(logging.NewNoopLogger()))
	})
}

func TestSessionContextData_AccountRolePermissionsChecker(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &sessions.ContextData{
			ActiveAccountID: t.Name(),
			AccountPermissions: map[string]authorization.AccountRolePermissionsChecker{
				t.Name(): authorization.NewAccountRolePermissionChecker(nil),
			},
		}

		assert.NotNil(t, x.AccountRolePermissionsChecker())
	})
}

func TestSessionContextData_ServiceRolePermissionChecker(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &sessions.ContextData{
			ActiveAccountID: t.Name(),
			Requester: sessions.RequesterInfo{
				ServicePermissions: authorization.NewServiceRolePermissionChecker([]string{t.Name()}, nil),
			},
		}

		assert.NotNil(t, x.ServiceRolePermissionChecker())
	})
}
