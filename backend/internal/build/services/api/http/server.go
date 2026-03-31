package api

import (
	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/auth"
	paymentswebhook "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/services/payments/http"

	"github.com/verygoodsoftwarenotvirus/platform/v4/healthcheck"
	"github.com/verygoodsoftwarenotvirus/platform/v4/observability/logging"
	"github.com/verygoodsoftwarenotvirus/platform/v4/observability/metrics"
	"github.com/verygoodsoftwarenotvirus/platform/v4/observability/tracing"
	"github.com/verygoodsoftwarenotvirus/platform/v4/routing"
	routingcfg "github.com/verygoodsoftwarenotvirus/platform/v4/routing/config"

	"github.com/samber/do/v2"
)

// RegisterAPIRouter registers the API router provider with the injector.
func RegisterAPIRouter(i do.Injector) {
	do.Provide[routing.Router](i, func(i do.Injector) (routing.Router, error) {
		return ProvideAPIRouter(
			*do.MustInvoke[*routingcfg.Config](i),
			do.MustInvoke[logging.Logger](i),
			do.MustInvoke[tracing.TracerProvider](i),
			do.MustInvoke[metrics.Provider](i),
			do.MustInvoke[auth.AuthDataService](i),
			do.MustInvoke[*paymentswebhook.WebhookHandler](i),
			do.MustInvoke[healthcheck.Registry](i),
		)
	})
}
