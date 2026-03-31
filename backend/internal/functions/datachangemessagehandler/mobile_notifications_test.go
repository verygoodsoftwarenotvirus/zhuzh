package datachangemessagehandler

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/identity"
	domainnotifications "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/notifications"
	notificationsmock "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/notifications/mock"

	"github.com/verygoodsoftwarenotvirus/platform/v4/database/filtering"
	notifications "github.com/verygoodsoftwarenotvirus/platform/v4/notifications/mobile"
	"github.com/verygoodsoftwarenotvirus/platform/v4/reflection"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMobileNotificationsEventHandler(t *testing.T) {
	t.Parallel()

	t.Run("invalid JSON", func(t *testing.T) {
		t.Parallel()

		handler, _, _, _, _, _, _, _, _, _, _ := buildTestAsyncDataChangeMessageHandler(t)

		err := handler.MobileNotificationsEventHandler("mobile_notifications")(t.Context(), []byte("not json"))

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "decoding")
	})

	t.Run("missing title", func(t *testing.T) {
		t.Parallel()

		handler, _, _, _, _, _, _, _, _, _, _ := buildTestAsyncDataChangeMessageHandler(t)

		req := notifications.MobileNotificationRequest{
			RequestType:      identity.MobileNotificationRequestTypeHouseholdInvitationAccepted,
			RecipientUserIDs: []string{"user-1"},
			Title:            "",
			Body:             "body",
		}
		raw, _ := json.Marshal(req)

		err := handler.MobileNotificationsEventHandler("mobile_notifications")(t.Context(), raw)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "title")
	})

	t.Run("missing body", func(t *testing.T) {
		t.Parallel()

		handler, _, _, _, _, _, _, _, _, _, _ := buildTestAsyncDataChangeMessageHandler(t)

		req := notifications.MobileNotificationRequest{
			RequestType:      identity.MobileNotificationRequestTypeHouseholdInvitationAccepted,
			RecipientUserIDs: []string{"user-1"},
			Title:            "title",
			Body:             "",
		}
		raw, _ := json.Marshal(req)

		err := handler.MobileNotificationsEventHandler("mobile_notifications")(t.Context(), raw)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "body")
	})

	t.Run("missing request type", func(t *testing.T) {
		t.Parallel()

		handler, _, _, _, _, _, _, _, _, _, _ := buildTestAsyncDataChangeMessageHandler(t)

		req := notifications.MobileNotificationRequest{
			RecipientUserIDs: []string{"user-1"},
			Title:            "title",
			Body:             "body",
		}
		raw, _ := json.Marshal(req)

		err := handler.MobileNotificationsEventHandler("mobile_notifications")(t.Context(), raw)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "request type")
	})

	t.Run("unknown request type", func(t *testing.T) {
		t.Parallel()

		handler, _, _, _, _, _, _, _, _, _, _ := buildTestAsyncDataChangeMessageHandler(t)

		req := notifications.MobileNotificationRequest{
			RequestType:      "unknown_type",
			RecipientUserIDs: []string{"user-1"},
			Title:            "title",
			Body:             "body",
		}
		raw, _ := json.Marshal(req)

		err := handler.MobileNotificationsEventHandler("mobile_notifications")(t.Context(), raw)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unknown request type")
	})

	t.Run("household invitation accepted with no recipients", func(t *testing.T) {
		t.Parallel()

		handler, _, _, _, _, _, _, _, _, _, _ := buildTestAsyncDataChangeMessageHandler(t)

		req := notifications.MobileNotificationRequest{
			RequestType:      identity.MobileNotificationRequestTypeHouseholdInvitationAccepted,
			RecipientUserIDs: []string{},
			Title:            "title",
			Body:             "body",
		}
		raw, _ := json.Marshal(req)

		err := handler.MobileNotificationsEventHandler("mobile_notifications")(t.Context(), raw)

		assert.NoError(t, err)
	})

	t.Run("no device tokens for recipients", func(t *testing.T) {
		t.Parallel()

		handler, _, _, _, _, _, _, _, _, _, _ := buildTestAsyncDataChangeMessageHandler(t)
		notificationsRepo := &notificationsmock.Repository{}
		handler.notificationsRepo = notificationsRepo

		req := notifications.MobileNotificationRequest{
			RequestType:      identity.MobileNotificationRequestTypeHouseholdInvitationAccepted,
			RecipientUserIDs: []string{"user-1"},
			Title:            "title",
			Body:             "body",
		}
		raw, _ := json.Marshal(req)

		notificationsRepo.On(reflection.GetMethodName(notificationsRepo.GetUserDeviceTokens), mock.Anything, "user-1", mock.Anything, (*string)(nil)).Return(&filtering.QueryFilteredResult[domainnotifications.UserDeviceToken]{Data: []*domainnotifications.UserDeviceToken{}}, nil).Once()

		err := handler.MobileNotificationsEventHandler("mobile_notifications")(t.Context(), raw)

		assert.NoError(t, err)
		mock.AssertExpectationsForObjects(t, notificationsRepo)
	})

	t.Run("success sends push notification", func(t *testing.T) {
		t.Parallel()

		handler, _, _, _, _, _, _, _, _, _, _ := buildTestAsyncDataChangeMessageHandler(t)
		notificationsRepo := &notificationsmock.Repository{}
		handler.notificationsRepo = notificationsRepo

		req := notifications.MobileNotificationRequest{
			RequestType:      identity.MobileNotificationRequestTypeHouseholdInvitationAccepted,
			RecipientUserIDs: []string{"user-1"},
			Title:            "Household invitation accepted",
			Body:             "A new member has joined your household",
		}
		raw, _ := json.Marshal(req)

		deviceToken := &domainnotifications.UserDeviceToken{
			ID:            "token-1",
			DeviceToken:   strings.Repeat("a", 64),
			Platform:      domainnotifications.UserDeviceTokenPlatformIOS,
			BelongsToUser: "user-1",
		}

		notificationsRepo.On(reflection.GetMethodName(notificationsRepo.GetUserDeviceTokens), mock.Anything, "user-1", mock.Anything, (*string)(nil)).Return(&filtering.QueryFilteredResult[domainnotifications.UserDeviceToken]{
			Data: []*domainnotifications.UserDeviceToken{deviceToken},
		}, nil).Once()

		err := handler.MobileNotificationsEventHandler("mobile_notifications")(t.Context(), raw)

		assert.NoError(t, err)
		mock.AssertExpectationsForObjects(t, notificationsRepo)
	})
}
