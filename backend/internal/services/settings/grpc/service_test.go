package grpc

import (
	"testing"

	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/settings/mock"

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
		settingsManager := &mock.Repository{}

		service := NewService(
			logger,
			tracerProvider,
			settingsManager,
		)

		assert.NotNil(t, service)

		// Type assertion to ensure we get the correct implementation
		impl, ok := service.(*serviceImpl)
		assert.True(t, ok)
		assert.NotNil(t, impl.logger)
		assert.NotNil(t, impl.tracer)
		assert.Equal(t, settingsManager, impl.settingsManager)
		assert.NotNil(t, impl.sessionContextDataFetcher)
	})
}
