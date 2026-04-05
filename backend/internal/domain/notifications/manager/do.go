package manager

import (
	"context"

	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/notifications"
	pgnotificationsrepo "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/repositories/postgres/notifications"
	sqlitenotificationsrepo "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/repositories/sqlite/notifications"

	databasecfg "github.com/verygoodsoftwarenotvirus/platform/v4/database/config"
	"github.com/verygoodsoftwarenotvirus/platform/v4/messagequeue"
	msgconfig "github.com/verygoodsoftwarenotvirus/platform/v4/messagequeue/config"
	"github.com/verygoodsoftwarenotvirus/platform/v4/observability/logging"
	"github.com/verygoodsoftwarenotvirus/platform/v4/observability/tracing"

	"github.com/samber/do/v2"
)

// RegisterNotificationsDataManager registers the notifications data manager with the injector.
func RegisterNotificationsDataManager(i do.Injector) {
	// Bind the concrete repo (already registered in registry.go) to the notificationsRepo interface.
	// Pick the right concrete type based on the configured database provider.
	do.Provide[notificationsRepo](i, func(i do.Injector) (notificationsRepo, error) {
		cfg := do.MustInvoke[*databasecfg.Config](i)
		if cfg.Provider == databasecfg.ProviderSQLite {
			return do.MustInvoke[*sqlitenotificationsrepo.Repository](i), nil
		}
		return do.MustInvoke[*pgnotificationsrepo.Repository](i), nil
	})

	do.Provide[NotificationsDataManager](i, func(i do.Injector) (NotificationsDataManager, error) {
		return NewNotificationsDataManager(
			do.MustInvoke[context.Context](i),
			do.MustInvoke[tracing.TracerProvider](i),
			do.MustInvoke[logging.Logger](i),
			do.MustInvoke[notificationsRepo](i),
			do.MustInvoke[*msgconfig.QueuesConfig](i),
			do.MustInvoke[messagequeue.PublisherProvider](i),
		)
	})

	// Bind NotificationsDataManager to notifications.Repository
	do.Provide[notifications.Repository](i, func(i do.Injector) (notifications.Repository, error) {
		return do.MustInvoke[NotificationsDataManager](i), nil
	})
}
