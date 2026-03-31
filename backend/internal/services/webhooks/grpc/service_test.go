package grpc

import (
	"testing"

	webhookmgrmock "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/webhooks/manager/mock"
	webhookssvc "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/grpc/generated/services/webhooks"

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
		webhookManager := &webhookmgrmock.WebhookDataManager{}

		service := NewService(logger, tracerProvider, webhookManager)

		assert.NotNil(t, service)
		assert.Implements(t, (*webhookssvc.WebhooksServiceServer)(nil), service)
	})
}
