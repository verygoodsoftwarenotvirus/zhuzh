package datachangemessagehandler

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/audit"
	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/identity"
	identityfakes "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/identity/fakes"
	identitykeys "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/identity/keys"
	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/webhooks"
	webhooksfakes "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/webhooks/fakes"
	identityindexing "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/services/identity/indexing"

	msgqueuemock "github.com/verygoodsoftwarenotvirus/platform/v4/messagequeue/mock"
	"github.com/verygoodsoftwarenotvirus/platform/v4/reflection"
	textsearch "github.com/verygoodsoftwarenotvirus/platform/v4/search/text"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAsyncDataChangeMessageHandler_DataChangesEventHandler(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		handler, identityRepo, webhookRepo, _, _, analyticsReporter, _, _, _, decoder, _ := buildTestAsyncDataChangeMessageHandler(t)

		ctx := t.Context()

		// Create a test data change message
		dataChangeMessage := &audit.DataChangeMessage{
			EventType: identity.UserSignedUpServiceEventType,
			UserID:    "test-user-id",
			AccountID: "test-account-id",
			Context:   nil, // When marshaled as empty map {} and unmarshaled, becomes nil
		}

		// Keep webhook processing enabled for full test coverage

		rawMsg, err := json.Marshal(dataChangeMessage)
		assert.NoError(t, err)

		// Set up decoder mock: DataChangesEventHandler decodes rawMsg into a DataChangeMessage
		decoder.On(reflection.GetMethodName(decoder.DecodeBytes), mock.Anything, rawMsg, mock.Anything).Run(func(args mock.Arguments) {
			dest := args.Get(2).(*audit.DataChangeMessage)
			*dest = *dataChangeMessage
		}).Return(nil).Once()

		// Set up mock expectations
		analyticsReporter.On(reflection.GetMethodName(analyticsReporter.EventOccurred), mock.Anything, dataChangeMessage.EventType, dataChangeMessage.UserID, dataChangeMessage.Context).Return(nil).Once()
		analyticsReporter.On(reflection.GetMethodName(analyticsReporter.AddUser), mock.Anything, dataChangeMessage.UserID, dataChangeMessage.Context).Return(nil).Maybe()
		webhookRepo.On(reflection.GetMethodName(webhookRepo.GetWebhooksForAccountAndEvent), mock.Anything, dataChangeMessage.AccountID, dataChangeMessage.EventType).Return([]*webhooks.Webhook{}, nil).Once()
		identityRepo.On(reflection.GetMethodName(identityRepo.GetUser), mock.Anything, dataChangeMessage.UserID).Return(identityfakes.BuildFakeUser(), nil).Maybe()

		// Set up mock expectations for publishers
		if handler.searchDataIndexPublisher != nil {
			handler.searchDataIndexPublisher.(*msgqueuemock.Publisher).On("Publish", mock.Anything, mock.AnythingOfType("*textsearch.IndexRequest")).Return(nil).Maybe()
		}
		if handler.outboundEmailsPublisher != nil {
			handler.outboundEmailsPublisher.(*msgqueuemock.Publisher).On("Publish", mock.Anything, mock.Anything).Return(nil).Maybe()
		}
		if handler.webhookExecutionRequestPublisher != nil {
			handler.webhookExecutionRequestPublisher.(*msgqueuemock.Publisher).On("Publish", mock.Anything, mock.Anything).Return(nil).Maybe()
		}

		assert.NoError(t, handler.DataChangesEventHandler("data_changes")(ctx, rawMsg))

		mock.AssertExpectationsForObjects(t, analyticsReporter, webhookRepo, identityRepo, decoder)
	})

	t.Run("with invalid JSON", func(t *testing.T) {
		t.Parallel()

		handler, _, _, _, _, _, _, _, _, decoder, _ := buildTestAsyncDataChangeMessageHandler(t)

		ctx := t.Context()
		rawMsg := []byte("invalid json")

		decoder.On(reflection.GetMethodName(decoder.DecodeBytes), mock.Anything, rawMsg, mock.Anything).Return(errors.New("invalid character 'i' looking for beginning of value")).Once()

		err := handler.DataChangesEventHandler("data_changes")(ctx, rawMsg)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "decoding message body")

		mock.AssertExpectationsForObjects(t, decoder)
	})
}

func TestAsyncDataChangeMessageHandler_handleDataChangeMessage(t *testing.T) {
	t.Parallel()

	t.Run("with nil message", func(t *testing.T) {
		t.Parallel()

		handler, _, _, _, _, _, _, _, _, _, _ := buildTestAsyncDataChangeMessageHandler(t)

		ctx := t.Context()

		err := handler.handleDataChangeMessage(ctx, nil, "data_changes")
		assert.Error(t, err)
		assert.Equal(t, errRequiredDataIsNil, err)
	})

	t.Run("with analytics event reporting", func(t *testing.T) {
		t.Parallel()

		handler, identityRepo, _, _, _, analyticsEventReporter, _, _, _, _, _ := buildTestAsyncDataChangeMessageHandler(t)

		ctx := t.Context()

		dataChangeMessage := &audit.DataChangeMessage{
			EventType: identity.UserSignedUpServiceEventType,
			UserID:    "test-user-id",
			AccountID: "test-account-id",
			Context:   nil,
		}

		// Set non-webhook event types to exclude this event (focuses test on analytics only)
		handler.SetNonWebhookEventTypes([]string{dataChangeMessage.EventType})

		analyticsEventReporter.On(reflection.GetMethodName(analyticsEventReporter.EventOccurred), mock.Anything, dataChangeMessage.EventType, dataChangeMessage.UserID, dataChangeMessage.Context).Return(nil)
		analyticsEventReporter.On(reflection.GetMethodName(analyticsEventReporter.AddUser), mock.Anything, dataChangeMessage.UserID, dataChangeMessage.Context).Return(nil).Maybe()
		identityRepo.On(reflection.GetMethodName(identityRepo.GetUser), mock.Anything, dataChangeMessage.UserID).Return(identityfakes.BuildFakeUser(), nil).Maybe()

		// Mock the search index publisher (handleDataChangeMessage calls handleSearchIndexUpdates in a goroutine)
		mockSearchPublisher := &msgqueuemock.Publisher{}
		mockSearchPublisher.On(reflection.GetMethodName(mockSearchPublisher.Publish), mock.Anything, mock.AnythingOfType("*textsearch.IndexRequest")).Return(nil)
		handler.searchDataIndexPublisher = mockSearchPublisher

		err := handler.handleDataChangeMessage(ctx, dataChangeMessage, "data_changes")
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, analyticsEventReporter, identityRepo, mockSearchPublisher)
	})

	t.Run("with webhook execution", func(t *testing.T) {
		t.Parallel()

		handler, identityRepo, webhookRepo, _, _, analyticsEventReporter, _, _, _, _, _ := buildTestAsyncDataChangeMessageHandler(t)

		ctx := t.Context()

		webhook := webhooksfakes.BuildFakeWebhook()

		dataChangeMessage := &audit.DataChangeMessage{
			EventType: identity.UserSignedUpServiceEventType,
			UserID:    "test-user-id",
			AccountID: "test-account-id",
			Context:   nil,
		}

		// Set non-webhook event types to exclude this event
		handler.SetNonWebhookEventTypes([]string{})

		analyticsEventReporter.On(reflection.GetMethodName(analyticsEventReporter.EventOccurred), mock.Anything, dataChangeMessage.EventType, dataChangeMessage.UserID, dataChangeMessage.Context).Return(nil)
		analyticsEventReporter.On(reflection.GetMethodName(analyticsEventReporter.AddUser), mock.Anything, dataChangeMessage.UserID, dataChangeMessage.Context).Return(nil).Maybe()
		webhookRepo.On(reflection.GetMethodName(webhookRepo.GetWebhooksForAccountAndEvent), mock.Anything, dataChangeMessage.AccountID, dataChangeMessage.EventType).Return([]*webhooks.Webhook{webhook}, nil)
		identityRepo.On(reflection.GetMethodName(identityRepo.GetUser), mock.Anything, dataChangeMessage.UserID).Return(identityfakes.BuildFakeUser(), nil).Maybe()

		// Mock the webhook execution request publisher
		mockWebhookPublisher := &msgqueuemock.Publisher{}
		mockWebhookPublisher.On(reflection.GetMethodName(mockWebhookPublisher.Publish), mock.Anything, mock.MatchedBy(func(req *webhooks.WebhookExecutionRequest) bool {
			return req.WebhookID == webhook.ID && req.AccountID == dataChangeMessage.AccountID
		})).Return(nil)
		handler.webhookExecutionRequestPublisher = mockWebhookPublisher

		// Mock the search index publisher (handleDataChangeMessage calls handleSearchIndexUpdates in a goroutine)
		mockSearchPublisher := &msgqueuemock.Publisher{}
		mockSearchPublisher.On(reflection.GetMethodName(mockSearchPublisher.Publish), mock.Anything, mock.AnythingOfType("*textsearch.IndexRequest")).Return(nil)
		handler.searchDataIndexPublisher = mockSearchPublisher

		err := handler.handleDataChangeMessage(ctx, dataChangeMessage, "data_changes")
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, analyticsEventReporter, webhookRepo, identityRepo, mockWebhookPublisher, mockSearchPublisher)
	})
}

func TestAsyncDataChangeMessageHandler_handleSearchIndexUpdates(t *testing.T) {
	t.Parallel()

	t.Run("user signed up event", func(t *testing.T) {
		t.Parallel()

		handler, _, _, _, _, _, _, _, _, _, _ := buildTestAsyncDataChangeMessageHandler(t)

		ctx := t.Context()

		dataChangeMessage := &audit.DataChangeMessage{
			EventType: identity.UserSignedUpServiceEventType,
			UserID:    "test-user-id",
			AccountID: "test-account-id",
			Context:   nil,
		}

		// Mock the search data index publisher
		mockPublisher := &msgqueuemock.Publisher{}
		mockPublisher.On(reflection.GetMethodName(mockPublisher.Publish), mock.Anything, mock.MatchedBy(func(req *textsearch.IndexRequest) bool {
			return req.RowID == dataChangeMessage.UserID &&
				req.IndexType == identityindexing.IndexTypeUsers &&
				req.Delete == false
		})).Return(nil)
		handler.searchDataIndexPublisher = mockPublisher

		err := handler.handleSearchIndexUpdates(ctx, dataChangeMessage)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mockPublisher)
	})

	t.Run("user archived event", func(t *testing.T) {
		t.Parallel()

		handler, _, _, _, _, _, _, _, _, _, _ := buildTestAsyncDataChangeMessageHandler(t)

		ctx := t.Context()

		dataChangeMessage := &audit.DataChangeMessage{
			EventType: identity.UserArchivedServiceEventType,
			UserID:    "test-user-id",
			AccountID: "test-account-id",
			Context:   nil,
		}

		// Mock the search data index publisher
		mockPublisher := &msgqueuemock.Publisher{}
		mockPublisher.On(reflection.GetMethodName(mockPublisher.Publish), mock.Anything, mock.MatchedBy(func(req *textsearch.IndexRequest) bool {
			return req.RowID == dataChangeMessage.UserID &&
				req.IndexType == identityindexing.IndexTypeUsers &&
				req.Delete == true
		})).Return(nil)
		handler.searchDataIndexPublisher = mockPublisher

		err := handler.handleSearchIndexUpdates(ctx, dataChangeMessage)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mockPublisher)
	})

	t.Run("unhandled event type", func(t *testing.T) {
		t.Parallel()

		handler, _, _, _, _, _, _, _, _, _, _ := buildTestAsyncDataChangeMessageHandler(t)

		ctx := t.Context()

		dataChangeMessage := &audit.DataChangeMessage{
			EventType: "unhandled.event.type",
			UserID:    "test-user-id",
			AccountID: "test-account-id",
			Context:   nil,
		}

		err := handler.handleSearchIndexUpdates(ctx, dataChangeMessage)
		assert.NoError(t, err) // Should not error for unhandled event types
	})

	t.Run("with missing user ID", func(t *testing.T) {
		t.Parallel()

		handler, _, _, _, _, _, _, _, _, _, _ := buildTestAsyncDataChangeMessageHandler(t)

		ctx := t.Context()

		dataChangeMessage := &audit.DataChangeMessage{
			EventType: identity.UserSignedUpServiceEventType,
			UserID:    "", // Missing user ID
			AccountID: "test-account-id",
			Context:   nil,
		}

		// Mock the search index publisher (function still publishes even with missing user ID)
		mockPublisher := &msgqueuemock.Publisher{}
		mockPublisher.On(reflection.GetMethodName(mockPublisher.Publish), mock.Anything, mock.AnythingOfType("*textsearch.IndexRequest")).Return(nil)
		handler.searchDataIndexPublisher = mockPublisher

		err := handler.handleSearchIndexUpdates(ctx, dataChangeMessage)
		assert.NoError(t, err) // Should handle gracefully but log error

		mock.AssertExpectationsForObjects(t, mockPublisher)
	})
}

func TestAsyncDataChangeMessageHandler_handleOutboundNotifications(T *testing.T) {
	T.Run("with nil message", func(t *testing.T) {
		handler, _, _, _, _, _, _, _, _, _, _ := buildTestAsyncDataChangeMessageHandler(t)

		ctx := t.Context()

		err := handler.handleOutboundNotifications(ctx, nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "nil data change message")
	})

	T.Run("user signed up event", func(t *testing.T) {
		// Set environment variable needed for email configuration
		t.Setenv("ZHUZH_SERVICE_ENVIRONMENT", "testing")

		handler, identityRepo, _, _, _, analyticsEventReporter, _, _, _, _, _ := buildTestAsyncDataChangeMessageHandler(t)

		ctx := t.Context()

		user := identityfakes.BuildFakeUser()
		evf := "email-verification-token"

		dataChangeMessage := &audit.DataChangeMessage{
			EventType: identity.UserSignedUpServiceEventType,
			UserID:    user.ID,
			AccountID: "test-account-id",
			Context: map[string]any{
				identitykeys.UserEmailVerificationTokenKey: evf,
			},
		}

		identityRepo.On(reflection.GetMethodName(identityRepo.GetUser), mock.Anything, user.ID).Return(user, nil)
		analyticsEventReporter.On(reflection.GetMethodName(analyticsEventReporter.AddUser), mock.Anything, user.ID, dataChangeMessage.Context).Return(nil)

		mockPublisher := &msgqueuemock.Publisher{}
		mockPublisher.On(reflection.GetMethodName(mockPublisher.Publish), mock.Anything, mock.AnythingOfType("*email.OutboundEmailMessage")).Return(nil)
		handler.outboundEmailsPublisher = mockPublisher

		err := handler.handleOutboundNotifications(ctx, dataChangeMessage)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, identityRepo, analyticsEventReporter, mockPublisher)
	})

	T.Run("with user fetch error", func(t *testing.T) {
		// Set environment variable needed for email configuration
		t.Setenv("ZHUZH_SERVICE_ENVIRONMENT", "testing")

		handler, identityRepo, _, _, _, _, _, _, _, _, _ := buildTestAsyncDataChangeMessageHandler(t)

		ctx := t.Context()

		dataChangeMessage := &audit.DataChangeMessage{
			EventType: identity.UserSignedUpServiceEventType,
			UserID:    "test-user-id",
			AccountID: "test-account-id",
			Context:   nil,
		}

		expectedError := errors.New("user fetch error")
		identityRepo.On(reflection.GetMethodName(identityRepo.GetUser), mock.Anything, "test-user-id").Return((*identity.User)(nil), expectedError)

		err := handler.handleOutboundNotifications(ctx, dataChangeMessage)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "getting user")

		mock.AssertExpectationsForObjects(t, identityRepo)
	})

	T.Run("unhandled event type", func(t *testing.T) {
		// Set environment variable needed for email configuration
		t.Setenv("ZHUZH_SERVICE_ENVIRONMENT", "testing")

		handler, identityRepo, _, _, _, _, _, _, _, _, _ := buildTestAsyncDataChangeMessageHandler(t)

		ctx := t.Context()

		user := identityfakes.BuildFakeUser()

		dataChangeMessage := &audit.DataChangeMessage{
			EventType: "unhandled.event.type",
			UserID:    user.ID,
			AccountID: "test-account-id",
			Context:   nil,
		}

		identityRepo.On(reflection.GetMethodName(identityRepo.GetUser), mock.Anything, user.ID).Return(user, nil)

		err := handler.handleOutboundNotifications(ctx, dataChangeMessage)
		assert.NoError(t, err) // Should handle gracefully with no outbound emails

		mock.AssertExpectationsForObjects(t, identityRepo)
	})
}
