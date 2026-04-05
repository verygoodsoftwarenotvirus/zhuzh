package notifications

import (
	"database/sql"

	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/audit"
	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/repositories/sqlite/notifications/generated"

	"github.com/verygoodsoftwarenotvirus/platform/v4/cryptography/encryption"
	"github.com/verygoodsoftwarenotvirus/platform/v4/cryptography/encryption/salsa20"
	"github.com/verygoodsoftwarenotvirus/platform/v4/database"
	databasecfg "github.com/verygoodsoftwarenotvirus/platform/v4/database/config"
	"github.com/verygoodsoftwarenotvirus/platform/v4/observability/logging"
	"github.com/verygoodsoftwarenotvirus/platform/v4/observability/tracing"
)

const (
	o11yName = "notifications_db_client"
)

// Repository is the notifications repository implementation.
// Exported so the manager can wrap it; the manager is the sole provider of notifications.Repository for services.
type Repository struct {
	database.Client
	tracer                tracing.Tracer
	logger                logging.Logger
	generatedQuerier      generated.Querier
	auditLogEntryRepo     audit.Repository
	userDeviceTokenEncDec encryption.EncryptorDecryptor
	readDB                *sql.DB
	writeDB               *sql.DB
}

// ProvideNotificationsRepository provides a new repository.
func ProvideNotificationsRepository(
	logger logging.Logger,
	tracerProvider tracing.TracerProvider,
	auditLogEntryRepo audit.Repository,
	cfg *databasecfg.Config,
	client database.Client,
) *Repository {
	encDec, err := salsa20.NewEncryptorDecryptor(tracerProvider, logger, []byte(cfg.UserDeviceTokenEncryptionKey))
	if err != nil {
		return nil
	}

	c := &Repository{
		Client:                client,
		readDB:                client.ReadDB(),
		writeDB:               client.WriteDB(),
		tracer:                tracing.NewNamedTracer(tracerProvider, o11yName),
		generatedQuerier:      generated.New(),
		auditLogEntryRepo:     auditLogEntryRepo,
		userDeviceTokenEncDec: encDec,
		logger:                logging.NewNamedLogger(logger, o11yName),
	}

	return c
}
