package grpc

import (
	"testing"

	issuereportmock "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/issuereports/mock"
	issuereportssvc "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/grpc/generated/services/issue_reports"

	"github.com/verygoodsoftwarenotvirus/platform/v4/observability/logging"
	"github.com/verygoodsoftwarenotvirus/platform/v4/observability/tracing"

	"github.com/stretchr/testify/assert"
)

func TestNewService(t *testing.T) {
	t.Parallel()

	t.Run("standard", func(t *testing.T) {
		t.Parallel()

		logger := logging.NewNoopLogger()
		tracerProvider := tracing.NewNoopTracerProvider()
		issueReportsManager := &issuereportmock.Repository{}

		service := NewService(logger, tracerProvider, issueReportsManager)

		assert.NotNil(t, service)
		assert.Implements(t, (*issuereportssvc.IssueReportsServiceServer)(nil), service)
	})
}
