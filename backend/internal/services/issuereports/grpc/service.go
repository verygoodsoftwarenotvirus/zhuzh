package grpc

import (
	"context"

	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/authentication/sessions"
	issuereportsmanager "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/issuereports/manager"
	issuereportssvc "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/grpc/generated/services/issue_reports"

	"github.com/verygoodsoftwarenotvirus/platform/v4/observability/logging"
	"github.com/verygoodsoftwarenotvirus/platform/v4/observability/tracing"
)

const (
	o11yName = "issue_reports_service"
)

var _ issuereportssvc.IssueReportsServiceServer = (*serviceImpl)(nil)

type (
	serviceImpl struct {
		issuereportssvc.UnimplementedIssueReportsServiceServer
		tracer                    tracing.Tracer
		logger                    logging.Logger
		sessionContextDataFetcher func(context.Context) (*sessions.ContextData, error)
		issueReportsManager       issuereportsmanager.IssueReportsDataManager
	}
)

func NewService(
	logger logging.Logger,
	tracerProvider tracing.TracerProvider,
	issueReportsManager issuereportsmanager.IssueReportsDataManager,
) issuereportssvc.IssueReportsServiceServer {
	return &serviceImpl{
		logger:                    logging.NewNamedLogger(logger, o11yName),
		tracer:                    tracing.NewNamedTracer(tracerProvider, o11yName),
		sessionContextDataFetcher: sessions.FetchContextDataFromContext,
		issueReportsManager:       issueReportsManager,
	}
}
