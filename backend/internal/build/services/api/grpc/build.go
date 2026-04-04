package grpcapi

import (
	"context"

	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/authentication"
	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/authentication/sessions"
	tokenscfg "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/authentication/tokens/config"
	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/config"
	auditmanager "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/audit/manager"
	authmgr "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/auth/managers"
	commentsmanager "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/comments/manager"
	dataprivacymanager "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/dataprivacy/manager"
	identitymgr "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/identity/manager"
	issuereportsmanager "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/issuereports/manager"
	notificationsmanager "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/notifications/manager"
	oauthmgr "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/oauth/manager"
	paymentsmanager "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/payments/manager"
	settingsmanager "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/settings/manager"
	uploadedmediamanager "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/uploadedmedia/manager"
	waitlistsmanager "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/waitlists/manager"
	webhooksmanager "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/webhooks/manager"
	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/repositories"
	analyticssvc "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/services/analytics/grpc"
	auditsvc "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/services/audit/grpc"
	authsvc "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/services/auth/grpc"
	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/services/auth/grpc/interceptors"
	authhttpsvc "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/services/auth/handlers/authentication"
	commentssvc "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/services/comments/grpc"
	dataprivacysvc "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/services/dataprivacy/grpc"
	identitysvc "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/services/identity/grpc"
	internalopssvc "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/services/internalops/grpc"
	issuereportssvc "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/services/issuereports/grpc"
	notificationssvc "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/services/notifications/grpc"
	oauthsvc "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/services/oauth/grpc"
	paymentsadapters "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/services/payments/adapters"
	paymentssvc "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/services/payments/grpc"
	settingssvc "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/services/settings/grpc"
	uploadedmediacfg "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/services/uploadedmedia/config"
	uploadedmediasvc "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/services/uploadedmedia/grpc"
	waitlistssvc "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/services/waitlists/grpc"
	webhookssvc "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/services/webhooks/grpc"

	"github.com/verygoodsoftwarenotvirus/platform/v4/analytics/multisource"
	databasecfg "github.com/verygoodsoftwarenotvirus/platform/v4/database/config"
	featureflagscfg "github.com/verygoodsoftwarenotvirus/platform/v4/featureflags/config"
	"github.com/verygoodsoftwarenotvirus/platform/v4/httpclient"
	msgconfig "github.com/verygoodsoftwarenotvirus/platform/v4/messagequeue/config"
	"github.com/verygoodsoftwarenotvirus/platform/v4/observability"
	loggingcfg "github.com/verygoodsoftwarenotvirus/platform/v4/observability/logging/config"
	metricscfg "github.com/verygoodsoftwarenotvirus/platform/v4/observability/metrics/config"
	tracingcfg "github.com/verygoodsoftwarenotvirus/platform/v4/observability/tracing/config"
	"github.com/verygoodsoftwarenotvirus/platform/v4/qrcodes"
	"github.com/verygoodsoftwarenotvirus/platform/v4/random"
	"github.com/verygoodsoftwarenotvirus/platform/v4/server/grpc"
	uploadscfg "github.com/verygoodsoftwarenotvirus/platform/v4/uploads/config"
	"github.com/verygoodsoftwarenotvirus/platform/v4/uploads/objectstorage"

	"github.com/samber/do/v2"
)

// BuildInjector creates and configures the dependency injection container.
func BuildInjector(
	ctx context.Context,
	cfg *config.APIServiceConfig,
) *do.RootScope {
	i := do.New()

	do.ProvideValue(i, ctx)
	do.ProvideValue(i, cfg)

	// config field extraction
	RegisterConfigs(i)

	// platform providers
	observability.RegisterO11yConfigs(i)
	metricscfg.RegisterMetricsProvider(i)
	loggingcfg.RegisterLogger(i)
	tracingcfg.RegisterTracerProvider(i)
	httpclient.RegisterHTTPClient(i)
	msgconfig.RegisterMessageQueue(i)
	random.RegisterGenerator(i)
	repositories.RegisterMigrator(i)
	databasecfg.RegisterDatabase(i)
	grpc.RegisterGRPCServer(i)
	qrcodes.RegisterBuilder(i)
	uploadscfg.RegisterStorageConfig(i)
	objectstorage.RegisterUploadManager(i)
	featureflagscfg.RegisterFeatureFlagManager(i)
	multisource.RegisterMultiSourceEventReporter(i)

	// authentication
	authentication.RegisterAuth(i)
	sessions.RegisterSessionProviders(i)
	tokenscfg.RegisterTokenIssuer(i)
	interceptors.RegisterAuthInterceptor(i)

	// repositories (core) — provider-conditional
	repositories.RegisterRepositories(i, cfg.Database.Provider)

	// managers
	auditmanager.RegisterAuditDataManager(i)
	authmgr.RegisterAuthManager(i)
	commentsmanager.RegisterCommentsDataManager(i)
	identitymgr.RegisterIdentityDataManager(i)
	notificationsmanager.RegisterNotificationsDataManager(i)
	settingsmanager.RegisterSettingsDataManager(i)
	paymentsmanager.RegisterPaymentsDataManager(i)
	oauthmgr.RegisterOAuth2Manager(i)
	webhooksmanager.RegisterWebhookDataManager(i)
	waitlistsmanager.RegisterWaitlistDataManager(i)
	issuereportsmanager.RegisterIssueReportsDataManager(i)
	uploadedmediamanager.RegisterUploadedMediaManager(i)
	dataprivacymanager.RegisterDataPrivacyManager(i)
	paymentsadapters.RegisterPaymentProcessorRegistry(i)

	// services
	authsvc.RegisterAuthService(i)
	authhttpsvc.RegisterAuthHTTPService(i)
	analyticssvc.RegisterAnalyticsService(i)
	auditsvc.RegisterAuditService(i)
	commentssvc.RegisterCommentsService(i)
	dataprivacysvc.RegisterDataPrivacyService(i)
	identitysvc.RegisterIdentityService(i)
	internalopssvc.RegisterInternalOpsService(i)
	issuereportssvc.RegisterIssueReportsService(i)
	notificationssvc.RegisterNotificationsService(i)
	settingssvc.RegisterSettingsService(i)
	uploadedmediasvc.RegisterUploadedMediaService(i)
	webhookssvc.RegisterWebhooksService(i)
	oauthsvc.RegisterOAuthService(i)
	paymentssvc.RegisterPaymentsService(i)
	waitlistssvc.RegisterWaitlistsService(i)
	uploadedmediacfg.RegisterUploadedMediaConfig(i)

	// extras (functions from extras.go)
	RegisterExtras(i)

	return i
}

// Build builds a server.
func Build(
	ctx context.Context,
	cfg *config.APIServiceConfig,
) (*GRPCService, error) {
	i := BuildInjector(ctx, cfg)
	return do.MustInvoke[*GRPCService](i), nil
}
