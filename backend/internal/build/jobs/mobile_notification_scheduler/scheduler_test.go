package mobilenotificationscheduler

import (
	"testing"

	identitymock "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/identity/mock"

	msgqueuemock "github.com/verygoodsoftwarenotvirus/platform/v4/messagequeue/mock"
	"github.com/verygoodsoftwarenotvirus/platform/v4/observability/logging"
	"github.com/verygoodsoftwarenotvirus/platform/v4/observability/tracing"

	"github.com/stretchr/testify/require"
)

func TestScheduler_ScheduleNotifications_noHandlers(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	logger := logging.NewNoopLogger()
	tracerProvider := tracing.NewNoopTracerProvider()

	identityRepo := &identitymock.RepositoryMock{}
	publisher := &msgqueuemock.Publisher{}

	scheduler := NewScheduler(logger, tracerProvider, identityRepo, publisher)

	err := scheduler.ScheduleNotifications(ctx)

	require.NoError(t, err)
}
