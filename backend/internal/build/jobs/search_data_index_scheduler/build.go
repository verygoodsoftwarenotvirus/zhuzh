package searchdataindexscheduler

import (
	"context"

	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/config"
	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/identity"
	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/repositories"

	databasecfg "github.com/verygoodsoftwarenotvirus/platform/v4/database/config"
	"github.com/verygoodsoftwarenotvirus/platform/v4/database/postgres"
	msgconfig "github.com/verygoodsoftwarenotvirus/platform/v4/messagequeue/config"
	"github.com/verygoodsoftwarenotvirus/platform/v4/observability"
	loggingcfg "github.com/verygoodsoftwarenotvirus/platform/v4/observability/logging/config"
	metricscfg "github.com/verygoodsoftwarenotvirus/platform/v4/observability/metrics/config"
	tracingcfg "github.com/verygoodsoftwarenotvirus/platform/v4/observability/tracing/config"
	"github.com/verygoodsoftwarenotvirus/platform/v4/search/text/indexing"

	"github.com/samber/do/v2"
)

// BuildInjector creates and configures the dependency injection container.
func BuildInjector(
	ctx context.Context,
	cfg *config.SearchDataIndexSchedulerConfig,
) *do.RootScope {
	i := do.New()

	do.ProvideValue(i, ctx)
	do.ProvideValue(i, cfg)

	RegisterConfigs(i)

	observability.RegisterO11yConfigs(i)
	tracingcfg.RegisterTracerProvider(i)
	loggingcfg.RegisterLogger(i)
	metricscfg.RegisterMetricsProvider(i)
	msgconfig.RegisterMessageQueue(i)
	databasecfg.RegisterClientConfig(i)
	postgres.RegisterDatabaseClient(i)
	repositories.RegisterRepositories(i, cfg.Database.Provider)

	do.Provide[map[string]indexing.Function](i, func(i do.Injector) (map[string]indexing.Function, error) {
		identityRepo := do.MustInvoke[identity.Repository](i)
		return ProvideIndexFunctions(identityRepo), nil
	})

	indexing.RegisterIndexScheduler(i)

	return i
}

// Build builds a server.
func Build(
	ctx context.Context,
	cfg *config.SearchDataIndexSchedulerConfig,
) (*indexing.IndexScheduler, error) {
	i := BuildInjector(ctx, cfg)
	return do.MustInvoke[*indexing.IndexScheduler](i), nil
}
