package dataprivacy

import (
	"context"
	"testing"

	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/audit"
	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/identity"
	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/identity/converters"
	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/identity/fakes"
	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/repositories/postgres/auditlogentries"
	identityrepo "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/repositories/postgres/identity"
	issue_reports "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/repositories/postgres/issuereports"
	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/repositories/postgres/migrations"
	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/repositories/postgres/notifications"
	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/repositories/postgres/settings"
	pgtesting "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/repositories/postgres/testing"
	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/repositories/postgres/uploadedmedia"
	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/repositories/postgres/waitlists"
	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/repositories/postgres/webhooks"

	"github.com/verygoodsoftwarenotvirus/platform/v4/database"
	"github.com/verygoodsoftwarenotvirus/platform/v4/database/postgres"
	"github.com/verygoodsoftwarenotvirus/platform/v4/observability/logging"
	"github.com/verygoodsoftwarenotvirus/platform/v4/observability/tracing"

	"github.com/stretchr/testify/require"
	pgcontainers "github.com/testcontainers/testcontainers-go/modules/postgres"
)

func buildDatabaseClientForTest(t *testing.T) (repo *repository, auditRepo audit.Repository, idRepo identity.Repository, pgContainer *pgcontainers.PostgresContainer) {
	t.Helper()

	ctx := t.Context()
	container, db, config := pgtesting.BuildDatabaseContainerForTest(t)
	require.NoError(t, migrations.NewMigrator(logging.NewNoopLogger()).Migrate(ctx, db))

	pgc, err := postgres.ProvideDatabaseClient(ctx, logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), config, nil)
	require.NotNil(t, pgc)
	require.NoError(t, err)

	auditLogEntryRepo := auditlogentries.ProvideAuditLogRepository(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), pgc)
	identityRepo := identityrepo.ProvideIdentityRepository(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), auditLogEntryRepo, pgc)
	issueReportsRepo := issue_reports.ProvideIssueReportsRepository(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), auditLogEntryRepo, pgc)
	notificationsRepo := notifications.ProvideNotificationsRepository(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), auditLogEntryRepo, config, pgc)
	settingsRepo := settings.ProvideSettingsRepository(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), auditLogEntryRepo, pgc)
	uploadedMediaRepo := uploadedmedia.ProvideUploadedMediaRepository(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), auditLogEntryRepo, pgc)
	waitlistsRepo := waitlists.ProvideWaitlistsRepository(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), auditLogEntryRepo, pgc)
	webhooksRepo := webhooks.ProvideWebhooksRepository(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), auditLogEntryRepo, pgc)

	c := ProvideDataPrivacyRepository(
		logging.NewNoopLogger(),
		tracing.NewNoopTracerProvider(),
		auditLogEntryRepo,
		identityRepo,
		issueReportsRepo,
		notificationsRepo,
		settingsRepo,
		uploadedMediaRepo,
		waitlistsRepo,
		webhooksRepo,
		pgc,
		nil,
	)

	return c.(*repository), auditLogEntryRepo, identityRepo, container
}

func buildInertClientForTest(t *testing.T) *repository {
	t.Helper()

	c := ProvideDataPrivacyRepository(
		logging.NewNoopLogger(),
		tracing.NewNoopTracerProvider(),
		nil, // auditLogRepo
		nil, // identityRepo
		nil, // issueReportsRepo
		nil, // notificationsRepo
		nil, // settingsRepo
		nil, // uploadedMediaRepo
		nil, // waitlistsRepo
		nil, // webhooksRepo
		&database.MockClient{},
		nil, // dataCollectors
	)

	return c.(*repository)
}

func createUserForTest(t *testing.T, ctx context.Context, exampleUser *identity.User, identityRepo identity.Repository) *identity.User {
	t.Helper()

	if exampleUser == nil {
		exampleUser = fakes.BuildFakeUser()
	}
	exampleUser.TwoFactorSecretVerifiedAt = nil
	dbInput := converters.ConvertUserToUserDatabaseCreationInput(exampleUser)

	created, err := identityRepo.CreateUser(ctx, dbInput)
	require.NoError(t, err)
	require.NotNil(t, created)

	return created
}
