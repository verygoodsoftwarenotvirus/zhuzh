package datachangemessagehandler

import (
	"context"
	"testing"

	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/config"
	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/auth"
	dataprivacymock "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/dataprivacy/mock"
	identitymock "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/identity/mock"
	internalopsmock "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/internalops/mock"
	notificationsmock "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/notifications/mock"
	webhooksmock "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/webhooks/mock"
	identityindexing "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/services/identity/indexing"

	analyticsmock "github.com/verygoodsoftwarenotvirus/platform/v4/analytics/mock"
	emailmock "github.com/verygoodsoftwarenotvirus/platform/v4/email/mock"
	encodingmock "github.com/verygoodsoftwarenotvirus/platform/v4/encoding/mock"
	msgconfig "github.com/verygoodsoftwarenotvirus/platform/v4/messagequeue/config"
	msgqueuemock "github.com/verygoodsoftwarenotvirus/platform/v4/messagequeue/mock"
	noopnotifications "github.com/verygoodsoftwarenotvirus/platform/v4/notifications/mobile/noop"
	"github.com/verygoodsoftwarenotvirus/platform/v4/observability/logging"
	"github.com/verygoodsoftwarenotvirus/platform/v4/observability/metrics"
	mockmetrics "github.com/verygoodsoftwarenotvirus/platform/v4/observability/metrics/mock"
	"github.com/verygoodsoftwarenotvirus/platform/v4/observability/tracing"
	"github.com/verygoodsoftwarenotvirus/platform/v4/reflection"
	uploadsmock "github.com/verygoodsoftwarenotvirus/platform/v4/uploads/mock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// noopPasswordResetTokenDataManager implements auth.PasswordResetTokenDataManager for tests that do not exercise the password reset flow.
type noopPasswordResetTokenDataManager struct{}

func (noopPasswordResetTokenDataManager) GetPasswordResetTokenByID(context.Context, string) (*auth.PasswordResetToken, error) {
	return nil, nil
}
func (noopPasswordResetTokenDataManager) GetPasswordResetTokenByToken(context.Context, string) (*auth.PasswordResetToken, error) {
	return nil, nil
}
func (noopPasswordResetTokenDataManager) CreatePasswordResetToken(context.Context, *auth.PasswordResetTokenDatabaseCreationInput) (*auth.PasswordResetToken, error) {
	return nil, nil
}
func (noopPasswordResetTokenDataManager) RedeemPasswordResetToken(context.Context, string) error {
	return nil
}

//nolint:gocritic // I know this returns too many things
func buildTestAsyncDataChangeMessageHandler(t *testing.T) (*AsyncDataChangeMessageHandler, *identitymock.RepositoryMock, *webhooksmock.Repository, *msgqueuemock.ConsumerProvider, *msgqueuemock.PublisherProvider, *analyticsmock.EventReporter, *emailmock.Emailer, *uploadsmock.MockUploadManager, *mockmetrics.MetricsProvider, *encodingmock.EncoderDecoder, *dataprivacymock.Repository) {
	t.Helper()

	logger := logging.NewNoopLogger()
	tracer := tracing.NewTracerForTest(t.Name())

	identityRepo := &identitymock.RepositoryMock{}
	webhookRepo := &webhooksmock.Repository{}
	consumerProvider := &msgqueuemock.ConsumerProvider{}
	publisherProvider := &msgqueuemock.PublisherProvider{}
	analyticsEventReporter := &analyticsmock.EventReporter{}
	emailer := &emailmock.Emailer{}
	uploadManager := &uploadsmock.MockUploadManager{}
	metricsProvider := &mockmetrics.MetricsProvider{}
	decoder := &encodingmock.EncoderDecoder{}
	dataPrivacyRepo := &dataprivacymock.Repository{}

	// Create mock indexers with noop implementations for testing
	userDataIndexer := &identityindexing.UserDataIndexer{}

	// Set up mock publishers for the indexers to prevent nil pointer dereferences
	mockPublisher := &msgqueuemock.Publisher{}
	publisherProvider.On(reflection.GetMethodName(publisherProvider.ProvidePublisher), mock.AnythingOfType("string")).Return(mockPublisher, nil).Maybe()

	// Set up mock histograms and counters
	noopProvider := metrics.NewNoopMetricsProvider()
	noopHistogram, _ := noopProvider.NewFloat64Histogram("test")
	noopCounter, _ := noopProvider.NewInt64Counter("test")
	metricsProvider.On(reflection.GetMethodName(metricsProvider.NewFloat64Histogram), mock.AnythingOfType("string"), mock.Anything).Return(noopHistogram, nil).Maybe()
	metricsProvider.On(reflection.GetMethodName(metricsProvider.NewInt64Counter), mock.AnythingOfType("string"), mock.Anything).Return(noopCounter, nil).Maybe()

	internalOpsRepo := &internalopsmock.InternalOpsDataManager{}
	notificationsRepo := &notificationsmock.Repository{}
	pushNotificationSender := noopnotifications.NewPushNotificationSender()

	handler := &AsyncDataChangeMessageHandler{
		identityRepo:                         identityRepo,
		webhookRepo:                          webhookRepo,
		internalOpsRepo:                      internalOpsRepo,
		consumerProvider:                     consumerProvider,
		analyticsEventReporter:               analyticsEventReporter,
		emailer:                              emailer,
		uploadManager:                        uploadManager,
		decoder:                              decoder,
		userDataIndexer:                      userDataIndexer,
		logger:                               logger,
		tracer:                               tracer,
		nonWebhookEventTypes:                 []string{},
		dataChangesExecutionTimeHistogram:    noopHistogram,
		outboundEmailsExecutionTimeHistogram: noopHistogram,
		webhookExecutionTimestampHistogram:   noopHistogram,
		userDataAggregationExecutionTimeHistogram: noopHistogram,
		searchIndexRequestsExecutionTimeHistogram: noopHistogram,
		mobileNotificationsExecutionTimeHistogram: noopHistogram,
		messagesProcessedCounter:                  noopCounter,
		messageDecodeErrorsCounter:                noopCounter,
		handlerErrorsCounter:                      noopCounter,
		emailsSentCounter:                         noopCounter,
		emailsFailedCounter:                       noopCounter,
		pushNotificationsSentCounter:              noopCounter,
		badDeviceTokensArchivedCounter:            noopCounter,
		queuesConfig: msgconfig.QueuesConfig{
			SearchIndexRequestsTopicName: "search-index-requests",
		},
		searchDataIndexPublisher:         mockPublisher,
		outboundEmailsPublisher:          mockPublisher,
		webhookExecutionRequestPublisher: mockPublisher,
		mobileNotificationsPublisher:     mockPublisher,
		dataPrivacyRepo:                  dataPrivacyRepo,
		passwordResetTokenDataManager:    noopPasswordResetTokenDataManager{},
		notificationsRepo:                notificationsRepo,
		pushNotificationSender:           pushNotificationSender,
	}

	handler.searchIndexHandlers = []SearchIndexEventHandler{
		handler.handleIdentitySearchIndexUpdate,
	}
	handler.outboundNotificationHandlers = []OutboundNotificationHandler{
		handler.handleIdentityOutboundNotification,
	}

	return handler, identityRepo, webhookRepo, consumerProvider, publisherProvider, analyticsEventReporter, emailer, uploadManager, metricsProvider, decoder, dataPrivacyRepo
}

func TestNewAsyncDataChangeMessageHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		logger := logging.NewNoopLogger()
		tracerProvider := tracing.NewNoopTracerProvider()
		cfg := &config.AsyncMessageHandlerConfig{
			Queues: msgconfig.QueuesConfig{
				OutboundEmailsTopicName:           "outbound-emails",
				SearchIndexRequestsTopicName:      "search-index-requests",
				WebhookExecutionRequestsTopicName: "webhook-execution-requests",
				MobileNotificationsTopicName:      "mobile-notifications",
			},
		}
		identityRepo := &identitymock.RepositoryMock{}
		dataPrivacyRepo := &dataprivacymock.Repository{}
		webhookRepo := &webhooksmock.Repository{}
		consumerProvider := &msgqueuemock.ConsumerProvider{}
		publisherProvider := &msgqueuemock.PublisherProvider{}
		analyticsEventReporter := &analyticsmock.EventReporter{}
		emailer := &emailmock.Emailer{}
		uploadManager := &uploadsmock.MockUploadManager{}
		metricsProvider := &mockmetrics.MetricsProvider{}
		decoder := &encodingmock.EncoderDecoder{}
		coreDataIndexer := &identityindexing.UserDataIndexer{}

		// Set up metrics expectations
		noopProvider := metrics.NewNoopMetricsProvider()
		noopHistogram, _ := noopProvider.NewFloat64Histogram("test")
		noopCounter, _ := noopProvider.NewInt64Counter("test")
		metricsProvider.On(reflection.GetMethodName(metricsProvider.NewFloat64Histogram), "data_changes_execution_time", mock.Anything).Return(noopHistogram, nil)
		metricsProvider.On(reflection.GetMethodName(metricsProvider.NewFloat64Histogram), "outbound_emails_execution_time", mock.Anything).Return(noopHistogram, nil)
		metricsProvider.On(reflection.GetMethodName(metricsProvider.NewFloat64Histogram), "search_index_requests_execution_time", mock.Anything).Return(noopHistogram, nil)
		metricsProvider.On(reflection.GetMethodName(metricsProvider.NewFloat64Histogram), "user_data_aggregation_execution_time", mock.Anything).Return(noopHistogram, nil)
		metricsProvider.On(reflection.GetMethodName(metricsProvider.NewFloat64Histogram), "webhook_requests_execution_time", mock.Anything).Return(noopHistogram, nil)
		metricsProvider.On(reflection.GetMethodName(metricsProvider.NewFloat64Histogram), "mobile_notifications_execution_time", mock.Anything).Return(noopHistogram, nil)
		metricsProvider.On(reflection.GetMethodName(metricsProvider.NewInt64Counter), "messages_processed_total", mock.Anything).Return(noopCounter, nil)
		metricsProvider.On(reflection.GetMethodName(metricsProvider.NewInt64Counter), "message_decode_errors_total", mock.Anything).Return(noopCounter, nil)
		metricsProvider.On(reflection.GetMethodName(metricsProvider.NewInt64Counter), "handler_errors_total", mock.Anything).Return(noopCounter, nil)
		metricsProvider.On(reflection.GetMethodName(metricsProvider.NewInt64Counter), "emails_sent_total", mock.Anything).Return(noopCounter, nil)
		metricsProvider.On(reflection.GetMethodName(metricsProvider.NewInt64Counter), "emails_failed_total", mock.Anything).Return(noopCounter, nil)
		metricsProvider.On(reflection.GetMethodName(metricsProvider.NewInt64Counter), "push_notifications_sent_total", mock.Anything).Return(noopCounter, nil)
		metricsProvider.On(reflection.GetMethodName(metricsProvider.NewInt64Counter), "bad_device_tokens_archived_total", mock.Anything).Return(noopCounter, nil)

		// Set up publisher expectations
		mockPublisher := &msgqueuemock.Publisher{}
		publisherProvider.On(reflection.GetMethodName(publisherProvider.ProvidePublisher), "outbound-emails").Return(mockPublisher, nil)
		publisherProvider.On(reflection.GetMethodName(publisherProvider.ProvidePublisher), "search-index-requests").Return(mockPublisher, nil)
		publisherProvider.On(reflection.GetMethodName(publisherProvider.ProvidePublisher), "webhook-execution-requests").Return(mockPublisher, nil)
		publisherProvider.On(reflection.GetMethodName(publisherProvider.ProvidePublisher), "mobile-notifications").Return(mockPublisher, nil)

		internalOpsRepo := &internalopsmock.InternalOpsDataManager{}
		prtManager := noopPasswordResetTokenDataManager{}
		notificationsRepo := &notificationsmock.Repository{}
		pushNotificationSender := noopnotifications.NewPushNotificationSender()

		handler, err := NewAsyncDataChangeMessageHandler(
			ctx,
			logger,
			tracerProvider,
			cfg,
			identityRepo,
			dataPrivacyRepo,
			webhookRepo,
			internalOpsRepo,
			consumerProvider,
			publisherProvider,
			analyticsEventReporter,
			emailer,
			uploadManager,
			metricsProvider,
			decoder,
			coreDataIndexer,
			prtManager,
			notificationsRepo,
			pushNotificationSender,
		)

		assert.NoError(t, err)
		assert.NotNil(t, handler)
		assert.Equal(t, identityRepo, handler.identityRepo)
		assert.Equal(t, webhookRepo, handler.webhookRepo)
		assert.Equal(t, consumerProvider, handler.consumerProvider)
		assert.Equal(t, analyticsEventReporter, handler.analyticsEventReporter)
		assert.Equal(t, emailer, handler.emailer)
		assert.Equal(t, uploadManager, handler.uploadManager)
		assert.Equal(t, decoder, handler.decoder)
		assert.Equal(t, coreDataIndexer, handler.userDataIndexer)

		mock.AssertExpectationsForObjects(t, metricsProvider, publisherProvider)
	})
}

func TestAsyncDataChangeMessageHandler_SetNonWebhookEventTypes(t *testing.T) {
	t.Parallel()

	t.Run("standard", func(t *testing.T) {
		t.Parallel()

		handler, _, _, _, _, _, _, _, _, _, _ := buildTestAsyncDataChangeMessageHandler(t)

		eventTypes := []string{"event1", "event2", "event3"}
		handler.SetNonWebhookEventTypes(eventTypes)

		handler.nonWebhookEventTypesHat.RLock()
		assert.Equal(t, eventTypes, handler.nonWebhookEventTypes)
		handler.nonWebhookEventTypesHat.RUnlock()
	})
}
