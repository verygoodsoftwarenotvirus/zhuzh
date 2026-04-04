package dataprivacy

import (
	"database/sql"

	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/audit"
	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/dataprivacy"
	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/identity"
	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/issuereports"
	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/notifications"
	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/settings"
	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/uploadedmedia"
	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/waitlists"
	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/webhooks"
	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/repositories/sqlite/dataprivacy/generated"

	"github.com/verygoodsoftwarenotvirus/platform/v4/database"
	"github.com/verygoodsoftwarenotvirus/platform/v4/observability/logging"
	"github.com/verygoodsoftwarenotvirus/platform/v4/observability/tracing"
)

const (
	o11yName = "dataprivacy_db_client"
)

// repository is the data privacy repository client.
type repository struct {
	issueReportsRepo  issuereports.Repository
	uploadedMediaRepo uploadedmedia.Repository
	logger            logging.Logger
	generatedQuerier  generated.Querier
	auditLogRepo      audit.Repository
	identityRepo      identity.Repository
	tracer            tracing.Tracer
	webhooksRepo      webhooks.Repository
	database.Client
	settingsRepo      settings.Repository
	notificationsRepo notifications.Repository
	waitlistsRepo     waitlists.Repository
	readDB            *sql.DB
	writeDB           *sql.DB
	dataCollectors    []dataprivacy.UserDataCollector
}

// ProvideDataPrivacyRepository provides a new repository.
func ProvideDataPrivacyRepository(
	logger logging.Logger,
	tracerProvider tracing.TracerProvider,
	auditLogRepo audit.Repository,
	identityRepo identity.Repository,
	issueReportsRepo issuereports.Repository,
	notificationsRepo notifications.Repository,
	settingsRepo settings.Repository,
	uploadedMediaRepo uploadedmedia.Repository,
	waitlistsRepo waitlists.Repository,
	webhooksRepo webhooks.Repository,
	client database.Client,
	dataCollectors []dataprivacy.UserDataCollector,
) dataprivacy.Repository {
	c := &repository{
		Client:            client,
		readDB:            client.ReadDB(),
		writeDB:           client.WriteDB(),
		tracer:            tracing.NewNamedTracer(tracerProvider, o11yName),
		logger:            logging.NewNamedLogger(logger, o11yName),
		generatedQuerier:  generated.New(),
		auditLogRepo:      auditLogRepo,
		identityRepo:      identityRepo,
		issueReportsRepo:  issueReportsRepo,
		notificationsRepo: notificationsRepo,
		settingsRepo:      settingsRepo,
		uploadedMediaRepo: uploadedMediaRepo,
		waitlistsRepo:     waitlistsRepo,
		webhooksRepo:      webhooksRepo,
		dataCollectors:    dataCollectors,
	}

	return c
}
