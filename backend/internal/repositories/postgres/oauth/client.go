package oauth

import (
	"database/sql"

	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/audit"
	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/oauth"
	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/repositories/postgres/oauth/generated"

	"github.com/verygoodsoftwarenotvirus/platform/v4/cryptography/encryption"
	encryptioncfg "github.com/verygoodsoftwarenotvirus/platform/v4/cryptography/encryption/config"
	"github.com/verygoodsoftwarenotvirus/platform/v4/database"
	databasecfg "github.com/verygoodsoftwarenotvirus/platform/v4/database/config"
	"github.com/verygoodsoftwarenotvirus/platform/v4/observability/logging"
	"github.com/verygoodsoftwarenotvirus/platform/v4/observability/tracing"
)

const (
	o11yName = "oauth_db_client"
)

// repository is the oauth2 client and token repo implemenation.
type repository struct {
	database.Client
	tracer                  tracing.Tracer
	logger                  logging.Logger
	generatedQuerier        generated.Querier
	auditLogEntryRepo       audit.Repository
	oauth2ClientTokenEncDec encryption.EncryptorDecryptor
	readDB                  *sql.DB
	writeDB                 *sql.DB
}

// ProvideOAuthRepository provides a new repository.
func ProvideOAuthRepository(
	logger logging.Logger,
	tracerProvider tracing.TracerProvider,
	auditLogEntryRepo audit.Repository,
	cfg *databasecfg.Config,
	client database.Client,
) oauth.Repository {
	encDec, err := encryptioncfg.ProvideEncryptorDecryptor(&cfg.Encryption, tracerProvider, logger, []byte(cfg.OAuth2TokenEncryptionKey))
	if err != nil {
		return nil
	}

	c := &repository{
		Client:                  client,
		readDB:                  client.ReadDB(),
		writeDB:                 client.WriteDB(),
		tracer:                  tracing.NewNamedTracer(tracerProvider, o11yName),
		generatedQuerier:        generated.New(),
		auditLogEntryRepo:       auditLogEntryRepo,
		oauth2ClientTokenEncDec: encDec,
		logger:                  logging.NewNamedLogger(logger, o11yName),
	}

	return c
}
