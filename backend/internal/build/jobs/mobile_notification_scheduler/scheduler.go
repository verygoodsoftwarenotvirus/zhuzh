package mobilenotificationscheduler

import (
	"context"

	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/identity"

	"github.com/verygoodsoftwarenotvirus/platform/v4/messagequeue"
	"github.com/verygoodsoftwarenotvirus/platform/v4/observability/logging"
	"github.com/verygoodsoftwarenotvirus/platform/v4/observability/tracing"

	"github.com/hashicorp/go-multierror"
)

const schedulerTracerName = "mobile_notification_scheduler"

// Scheduler publishes notification requests to the mobile_notifications queue.
type Scheduler struct {
	tracer tracing.Tracer
}

// NewScheduler creates a new mobile notification scheduler.
func NewScheduler(
	_ logging.Logger,
	tracerProvider tracing.TracerProvider,
	_ identity.Repository,
	_ messagequeue.Publisher,
) *Scheduler {
	return &Scheduler{
		tracer: tracing.NewNamedTracer(tracerProvider, schedulerTracerName),
	}
}

// ScheduleNotifications runs all notification schedulers. Each scheduler queries for items
// that need notifications and publishes them to the queue.
func (s *Scheduler) ScheduleNotifications(ctx context.Context) error {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	errs := &multierror.Error{}
	// Future notification types: add handler calls here.

	return errs.ErrorOrNil()
}
