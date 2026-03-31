package authentication

import (
	"context"

	authn "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/authentication"
	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/authentication/tokens"
	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/auth"
	identitymanager "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/identity/manager"
	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/oauth"

	"github.com/verygoodsoftwarenotvirus/platform/v4/analytics"
	"github.com/verygoodsoftwarenotvirus/platform/v4/encoding"
	"github.com/verygoodsoftwarenotvirus/platform/v4/messagequeue"
	msgconfig "github.com/verygoodsoftwarenotvirus/platform/v4/messagequeue/config"
	"github.com/verygoodsoftwarenotvirus/platform/v4/observability/logging"
	"github.com/verygoodsoftwarenotvirus/platform/v4/observability/tracing"
	"github.com/verygoodsoftwarenotvirus/platform/v4/routing"

	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/samber/do/v2"
)

// RegisterAuthHTTPService registers the auth HTTP service providers with the injector.
func RegisterAuthHTTPService(i do.Injector) {
	do.Provide[*manage.Manager](i, func(i do.Injector) (*manage.Manager, error) {
		return ProvideOAuth2ClientManager(
			do.MustInvoke[logging.Logger](i),
			do.MustInvoke[tracing.TracerProvider](i),
			do.MustInvoke[*OAuth2Config](i),
			do.MustInvoke[oauth.Repository](i),
		), nil
	})

	do.Provide[*server.Server](i, func(i do.Injector) (*server.Server, error) {
		return ProvideOAuth2ServerImplementation(
			do.MustInvoke[logging.Logger](i),
			do.MustInvoke[tracing.TracerProvider](i),
			do.MustInvoke[identitymanager.IdentityDataManager](i),
			do.MustInvoke[authn.Authenticator](i),
			do.MustInvoke[tokens.Issuer](i),
			do.MustInvoke[*manage.Manager](i),
		), nil
	})

	do.Provide[auth.AuthDataService](i, func(i do.Injector) (auth.AuthDataService, error) {
		return ProvideService(
			do.MustInvoke[context.Context](i),
			do.MustInvoke[logging.Logger](i),
			do.MustInvoke[*Config](i),
			do.MustInvoke[authn.Authenticator](i),
			do.MustInvoke[oauth.Repository](i),
			do.MustInvoke[identitymanager.IdentityDataManager](i),
			do.MustInvoke[encoding.ServerEncoderDecoder](i),
			do.MustInvoke[tracing.TracerProvider](i),
			do.MustInvoke[messagequeue.PublisherProvider](i),
			do.MustInvoke[analytics.EventReporter](i),
			do.MustInvoke[routing.RouteParamManager](i),
			do.MustInvoke[*msgconfig.QueuesConfig](i),
		)
	})
}
