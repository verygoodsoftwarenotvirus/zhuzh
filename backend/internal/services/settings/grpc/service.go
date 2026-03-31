package grpc

import (
	"context"

	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/authentication/sessions"
	settingsmanager "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/settings/manager"
	settingssvc "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/grpc/generated/services/settings"

	"github.com/verygoodsoftwarenotvirus/platform/v4/observability/logging"
	"github.com/verygoodsoftwarenotvirus/platform/v4/observability/tracing"
)

const (
	o11yName = "configuration_service"
)

var _ settingssvc.SettingsServiceServer = (*serviceImpl)(nil)

type (
	serviceImpl struct {
		settingssvc.UnimplementedSettingsServiceServer
		tracer                    tracing.Tracer
		logger                    logging.Logger
		sessionContextDataFetcher func(context.Context) (*sessions.ContextData, error)
		settingsManager           settingsmanager.SettingsDataManager
	}
)

func NewService(
	logger logging.Logger,
	tracerProvider tracing.TracerProvider,
	settingsManager settingsmanager.SettingsDataManager,
) settingssvc.SettingsServiceServer {
	return &serviceImpl{
		logger:                    logging.NewNamedLogger(logger, o11yName),
		tracer:                    tracing.NewNamedTracer(tracerProvider, o11yName),
		settingsManager:           settingsManager,
		sessionContextDataFetcher: sessions.FetchContextDataFromContext,
	}
}
