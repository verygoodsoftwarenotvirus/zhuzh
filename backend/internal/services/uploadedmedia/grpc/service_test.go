package grpc

import (
	"testing"

	uploadedmediamock "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/uploadedmedia/mock"
	uploadedmediasvc "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/grpc/generated/services/uploaded_media"

	"github.com/verygoodsoftwarenotvirus/platform/v4/observability/logging"
	"github.com/verygoodsoftwarenotvirus/platform/v4/observability/tracing"
	mockuploads "github.com/verygoodsoftwarenotvirus/platform/v4/uploads/mock"

	"github.com/stretchr/testify/assert"
)

func TestNewService(t *testing.T) {
	t.Parallel()

	t.Run("standard", func(t *testing.T) {
		t.Parallel()

		logger := logging.NewNoopLogger()
		tracerProvider := tracing.NewNoopTracerProvider()
		uploadedMediaManager := &uploadedmediamock.Repository{}

		service := NewService(logger, tracerProvider, uploadedMediaManager, &mockuploads.MockUploadManager{})

		assert.NotNil(t, service)
		assert.Implements(t, (*uploadedmediasvc.UploadedMediaServiceServer)(nil), service)
	})
}
