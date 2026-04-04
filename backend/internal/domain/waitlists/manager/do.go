package manager

import (
	"context"

	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/waitlists"
	pgwaitlistsrepo "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/repositories/postgres/waitlists"
	sqlitewaitlistsrepo "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/repositories/sqlite/waitlists"

	databasecfg "github.com/verygoodsoftwarenotvirus/platform/v4/database/config"
	"github.com/verygoodsoftwarenotvirus/platform/v4/messagequeue"
	msgconfig "github.com/verygoodsoftwarenotvirus/platform/v4/messagequeue/config"
	"github.com/verygoodsoftwarenotvirus/platform/v4/observability/logging"
	"github.com/verygoodsoftwarenotvirus/platform/v4/observability/tracing"

	"github.com/samber/do/v2"
)

// RegisterWaitlistDataManager registers the waitlist data manager with the injector.
func RegisterWaitlistDataManager(i do.Injector) {
	// Bind the concrete repo (already registered in registry.go) to the waitlistRepository interface.
	// Pick the right concrete type based on the configured database provider.
	do.Provide[waitlistRepository](i, func(i do.Injector) (waitlistRepository, error) {
		cfg := do.MustInvoke[*databasecfg.Config](i)
		if cfg.Provider == databasecfg.ProviderSQLite {
			return do.MustInvoke[*sqlitewaitlistsrepo.Repository](i), nil
		}
		return do.MustInvoke[*pgwaitlistsrepo.Repository](i), nil
	})

	do.Provide[WaitlistsDataManager](i, func(i do.Injector) (WaitlistsDataManager, error) {
		return NewWaitlistDataManager(
			do.MustInvoke[context.Context](i),
			do.MustInvoke[tracing.TracerProvider](i),
			do.MustInvoke[logging.Logger](i),
			do.MustInvoke[waitlistRepository](i),
			do.MustInvoke[*msgconfig.QueuesConfig](i),
			do.MustInvoke[messagequeue.PublisherProvider](i),
		)
	})

	// Bind WaitlistsDataManager to waitlists.Repository
	do.Provide[waitlists.Repository](i, func(i do.Injector) (waitlists.Repository, error) {
		return do.MustInvoke[WaitlistsDataManager](i), nil
	})
}
