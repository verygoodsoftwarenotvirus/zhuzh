package mock

import (
	"context"

	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/internalops"

	"github.com/stretchr/testify/mock"
)

var _ internalops.InternalOpsDataManager = (*InternalOpsDataManager)(nil)

type InternalOpsDataManager struct {
	mock.Mock
}

func (m *InternalOpsDataManager) DeleteExpiredOAuth2ClientTokens(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

func (m *InternalOpsDataManager) CreateQueueTestMessage(ctx context.Context, id, queueName string) error {
	return m.Called(ctx, id, queueName).Error(0)
}

func (m *InternalOpsDataManager) AcknowledgeQueueTestMessage(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}

func (m *InternalOpsDataManager) GetQueueTestMessage(ctx context.Context, id string) (*internalops.QueueTestMessage, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*internalops.QueueTestMessage), args.Error(1)
}

func (m *InternalOpsDataManager) PruneQueueTestMessages(ctx context.Context, queueName string) error {
	return m.Called(ctx, queueName).Error(0)
}
