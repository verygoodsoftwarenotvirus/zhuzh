package datachangemessagehandler

import (
	"context"

	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/config"
	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/auth"
	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/dataprivacy"
	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/identity"
	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/internalops"
	notificationsmanager "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/notifications/manager"
	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/webhooks"
	identityindexing "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/services/identity/indexing"

	"github.com/verygoodsoftwarenotvirus/platform/v4/analytics"
	"github.com/verygoodsoftwarenotvirus/platform/v4/email"
	"github.com/verygoodsoftwarenotvirus/platform/v4/encoding"
	"github.com/verygoodsoftwarenotvirus/platform/v4/messagequeue"
	notifications "github.com/verygoodsoftwarenotvirus/platform/v4/notifications/mobile"
	"github.com/verygoodsoftwarenotvirus/platform/v4/observability/logging"
	"github.com/verygoodsoftwarenotvirus/platform/v4/observability/metrics"
	"github.com/verygoodsoftwarenotvirus/platform/v4/observability/tracing"
	"github.com/verygoodsoftwarenotvirus/platform/v4/uploads"

	"github.com/samber/do/v2"
)

// RegisterAsyncDataChangeMessageHandler registers the async data change message handler with the injector.
func RegisterAsyncDataChangeMessageHandler(i do.Injector) {
	do.Provide[*AsyncDataChangeMessageHandler](i, func(i do.Injector) (*AsyncDataChangeMessageHandler, error) {
		return NewAsyncDataChangeMessageHandler(
			do.MustInvoke[context.Context](i),
			do.MustInvoke[logging.Logger](i),
			do.MustInvoke[tracing.TracerProvider](i),
			do.MustInvoke[*config.AsyncMessageHandlerConfig](i),
			do.MustInvoke[identity.Repository](i),
			do.MustInvoke[dataprivacy.Repository](i),
			do.MustInvoke[webhooks.Repository](i),
			do.MustInvoke[internalops.InternalOpsDataManager](i),
			do.MustInvoke[messagequeue.ConsumerProvider](i),
			do.MustInvoke[messagequeue.PublisherProvider](i),
			do.MustInvoke[analytics.EventReporter](i),
			do.MustInvoke[email.Emailer](i),
			do.MustInvoke[uploads.UploadManager](i),
			do.MustInvoke[metrics.Provider](i),
			do.MustInvoke[encoding.ServerEncoderDecoder](i),
			do.MustInvoke[*identityindexing.UserDataIndexer](i),
			do.MustInvoke[auth.PasswordResetTokenDataManager](i),
			do.MustInvoke[notificationsmanager.NotificationsDataManager](i),
			do.MustInvoke[notifications.PushNotificationSender](i),
		)
	})
}
