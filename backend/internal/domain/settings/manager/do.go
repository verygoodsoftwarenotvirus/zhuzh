package manager

import (
	"context"

	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/settings"

	"github.com/verygoodsoftwarenotvirus/platform/v4/messagequeue"
	msgconfig "github.com/verygoodsoftwarenotvirus/platform/v4/messagequeue/config"
	"github.com/verygoodsoftwarenotvirus/platform/v4/observability/logging"
	"github.com/verygoodsoftwarenotvirus/platform/v4/observability/tracing"

	"github.com/samber/do/v2"
)

// RegisterSettingsDataManager registers the settings data manager with the injector.
func RegisterSettingsDataManager(i do.Injector) {
	do.Provide[SettingsDataManager](i, func(i do.Injector) (SettingsDataManager, error) {
		return NewSettingsDataManager(
			do.MustInvoke[context.Context](i),
			do.MustInvoke[tracing.TracerProvider](i),
			do.MustInvoke[logging.Logger](i),
			do.MustInvoke[settings.Repository](i),
			do.MustInvoke[*msgconfig.QueuesConfig](i),
			do.MustInvoke[messagequeue.PublisherProvider](i),
		)
	})
}
