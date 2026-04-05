package internalops

import (
	"context"
	"time"

	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/internalops"
	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/repositories/sqlite/internalops/generated"

	platformerrors "github.com/verygoodsoftwarenotvirus/platform/v4/errors"
	"github.com/verygoodsoftwarenotvirus/platform/v4/observability"
)

func parseTime(s string) time.Time {
	t, err := time.Parse(time.RFC3339Nano, s)
	if err != nil {
		return time.Time{}
	}
	return t
}

func parseTimePtr(s *string) *time.Time {
	if s == nil {
		return nil
	}
	t := parseTime(*s)
	return &t
}

func (q *repository) CreateQueueTestMessage(ctx context.Context, id, queueName string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if id == "" || queueName == "" {
		return platformerrors.ErrInvalidIDProvided
	}

	if err := q.generatedQuerier.CreateQueueTestMessage(ctx, q.writeDB, &generated.CreateQueueTestMessageParams{
		ID:        id,
		QueueName: queueName,
	}); err != nil {
		return observability.PrepareError(err, span, "creating queue test message")
	}

	return nil
}

func (q *repository) AcknowledgeQueueTestMessage(ctx context.Context, id string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if id == "" {
		return platformerrors.ErrInvalidIDProvided
	}

	if err := q.generatedQuerier.AcknowledgeQueueTestMessage(ctx, q.writeDB, id); err != nil {
		return observability.PrepareError(err, span, "acknowledging queue test message")
	}

	return nil
}

func (q *repository) GetQueueTestMessage(ctx context.Context, id string) (*internalops.QueueTestMessage, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if id == "" {
		return nil, platformerrors.ErrInvalidIDProvided
	}

	row, err := q.generatedQuerier.GetQueueTestMessage(ctx, q.readDB, id)
	if err != nil {
		return nil, observability.PrepareError(err, span, "getting queue test message")
	}

	result := &internalops.QueueTestMessage{
		ID:             row.ID,
		QueueName:      row.QueueName,
		CreatedAt:      parseTime(row.CreatedAt),
		AcknowledgedAt: parseTimePtr(row.AcknowledgedAt),
	}

	return result, nil
}

func (q *repository) PruneQueueTestMessages(ctx context.Context, queueName string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if queueName == "" {
		return platformerrors.ErrInvalidIDProvided
	}

	if err := q.generatedQuerier.PruneQueueTestMessages(ctx, q.writeDB, queueName); err != nil {
		return observability.PrepareError(err, span, "pruning queue test messages")
	}

	return nil
}
