package managers

import (
	"context"

	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/authentication"
	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/auth"
	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/identity"

	"github.com/verygoodsoftwarenotvirus/platform/v4/messagequeue"
	msgconfig "github.com/verygoodsoftwarenotvirus/platform/v4/messagequeue/config"
	"github.com/verygoodsoftwarenotvirus/platform/v4/observability/logging"
	"github.com/verygoodsoftwarenotvirus/platform/v4/observability/tracing"
	"github.com/verygoodsoftwarenotvirus/platform/v4/qrcodes"
	"github.com/verygoodsoftwarenotvirus/platform/v4/random"

	"github.com/samber/do/v2"
)

// RegisterAuthManager registers the auth manager with the injector.
func RegisterAuthManager(i do.Injector) {
	do.Provide[AuthManagerInterface](i, func(i do.Injector) (AuthManagerInterface, error) {
		return ProvideAuthManager(
			do.MustInvoke[context.Context](i),
			do.MustInvoke[logging.Logger](i),
			do.MustInvoke[tracing.TracerProvider](i),
			do.MustInvoke[auth.PasswordResetTokenDataManager](i),
			do.MustInvoke[auth.UserSessionDataManager](i),
			do.MustInvoke[identity.UserDataManager](i),
			do.MustInvoke[authentication.Authenticator](i),
			do.MustInvoke[messagequeue.PublisherProvider](i),
			do.MustInvoke[random.Generator](i),
			do.MustInvoke[qrcodes.Builder](i),
			do.MustInvoke[*msgconfig.QueuesConfig](i),
		)
	})
}
