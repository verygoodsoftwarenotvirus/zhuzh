package localdev

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/authentication"
	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/authorization"
	apiserver "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/build/services/api"
	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/config"
	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/audit"
	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/auth"
	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/identity"
	identityconverters "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/identity/converters"
	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/notifications"
	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/oauth"
	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/settings"
	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/webhooks"
	authsvc "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/grpc/generated/services/auth"
	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/repositories"
	pgauditlogentries "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/repositories/postgres/auditlogentries"
	pgauth "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/repositories/postgres/auth"
	pgidentity "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/repositories/postgres/identity"
	pgnotifications "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/repositories/postgres/notifications"
	pgoauth "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/repositories/postgres/oauth"
	pgsettings "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/repositories/postgres/settings"
	pgtesting "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/repositories/postgres/testing"
	pgwebhooks "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/repositories/postgres/webhooks"
	sqliteauditlogentries "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/repositories/sqlite/auditlogentries"
	sqliteauth "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/repositories/sqlite/auth"
	sqliteidentity "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/repositories/sqlite/identity"
	sqlitenotifications "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/repositories/sqlite/notifications"
	sqliteoauth "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/repositories/sqlite/oauth"
	sqlitesettings "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/repositories/sqlite/settings"
	sqlitewebhooks "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/repositories/sqlite/webhooks"
	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/pkg/client"

	"github.com/verygoodsoftwarenotvirus/platform/v4/database"
	databasecfg "github.com/verygoodsoftwarenotvirus/platform/v4/database/config"
	"github.com/verygoodsoftwarenotvirus/platform/v4/httpclient"
	"github.com/verygoodsoftwarenotvirus/platform/v4/identifiers"
	msgconfig "github.com/verygoodsoftwarenotvirus/platform/v4/messagequeue/config"
	"github.com/verygoodsoftwarenotvirus/platform/v4/messagequeue/redis"
	"github.com/verygoodsoftwarenotvirus/platform/v4/observability/logging"
	"github.com/verygoodsoftwarenotvirus/platform/v4/observability/tracing"
	"github.com/verygoodsoftwarenotvirus/platform/v4/random"

	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// ProvideAuditLogRepository returns the audit log repository for the given provider.
func ProvideAuditLogRepository(provider string, logger logging.Logger, tracerProvider tracing.TracerProvider, dbClient database.Client) audit.Repository {
	switch provider {
	case databasecfg.ProviderSQLite:
		return sqliteauditlogentries.ProvideAuditLogRepository(logger, tracerProvider, dbClient)
	default:
		return pgauditlogentries.ProvideAuditLogRepository(logger, tracerProvider, dbClient)
	}
}

// ProvideIdentityRepository returns the identity repository for the given provider.
func ProvideIdentityRepository(provider string, logger logging.Logger, tracerProvider tracing.TracerProvider, auditRepo audit.Repository, dbClient database.Client) identity.Repository {
	switch provider {
	case databasecfg.ProviderSQLite:
		return sqliteidentity.ProvideIdentityRepository(logger, tracerProvider, auditRepo, dbClient)
	default:
		return pgidentity.ProvideIdentityRepository(logger, tracerProvider, auditRepo, dbClient)
	}
}

// ProvideAuthRepository returns the auth repository for the given provider.
func ProvideAuthRepository(provider string, logger logging.Logger, tracerProvider tracing.TracerProvider, auditRepo audit.Repository, dbClient database.Client) auth.Repository {
	switch provider {
	case databasecfg.ProviderSQLite:
		return sqliteauth.ProvideAuthRepository(logger, tracerProvider, auditRepo, dbClient)
	default:
		return pgauth.ProvideAuthRepository(logger, tracerProvider, auditRepo, dbClient)
	}
}

// ProvideOAuthRepository returns the OAuth repository for the given provider.
func ProvideOAuthRepository(provider string, logger logging.Logger, tracerProvider tracing.TracerProvider, auditRepo audit.Repository, dbCfg *databasecfg.Config, dbClient database.Client) oauth.Repository {
	switch provider {
	case databasecfg.ProviderSQLite:
		return sqliteoauth.ProvideOAuthRepository(logger, tracerProvider, auditRepo, dbCfg, dbClient)
	default:
		return pgoauth.ProvideOAuthRepository(logger, tracerProvider, auditRepo, dbCfg, dbClient)
	}
}

func CreatePremadeAdminUser(
	ctx context.Context,
	logger logging.Logger,
	tracerProvider tracing.TracerProvider,
	identityRepo identity.Repository,
	dbClient database.Client,
	premadeAdminUser *identity.User,
	provider string,
) (*identity.User, error) {
	hasher := authentication.ProvideArgon2Authenticator(logger, tracerProvider)

	actuallyHashedPass, err := hasher.HashPassword(ctx, premadeAdminUser.HashedPassword)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}
	premadeAdminUser.HashedPassword = actuallyHashedPass

	var user *identity.User
	if user, err = identityRepo.GetUserByUsername(ctx, premadeAdminUser.Username); err == nil {
		return user, nil
	}

	user, err = identityRepo.CreateUser(ctx, identityconverters.ConvertUserToUserDatabaseCreationInput(premadeAdminUser))
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Promote user to service_admin by archiving old service role and assigning new one.
	var archiveSQL, insertSQL string
	switch provider {
	case databasecfg.ProviderSQLite:
		archiveSQL = "UPDATE user_role_assignments SET archived_at = strftime('%Y-%m-%dT%H:%M:%fZ', 'now') WHERE user_id = ? AND account_id IS NULL AND archived_at IS NULL"
		insertSQL = "INSERT INTO user_role_assignments (id, user_id, role_id) VALUES (?, ?, ?)"
	default:
		archiveSQL = "UPDATE user_role_assignments SET archived_at = NOW() WHERE user_id = $1 AND account_id IS NULL AND archived_at IS NULL"
		insertSQL = "INSERT INTO user_role_assignments (id, user_id, role_id) VALUES ($1, $2, $3)"
	}

	if _, err = dbClient.WriteDB().ExecContext(ctx, archiveSQL, user.ID); err != nil {
		return nil, fmt.Errorf("failed to archive old service role: %w", err)
	}
	if _, err = dbClient.WriteDB().ExecContext(ctx, insertSQL, identifiers.New(), user.ID, authorization.ServiceAdminRoleID); err != nil {
		return nil, fmt.Errorf("failed to assign service_admin role: %w", err)
	}

	if err = identityRepo.MarkUserTwoFactorSecretAsVerified(ctx, user.ID); err != nil {
		return nil, fmt.Errorf("failed to mark user as verified: %w", err)
	}

	adminUser, err := identityRepo.GetAdminUserByUsername(ctx, user.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to get admin user: %w", err)
	}

	return adminUser, nil
}

func CreateOAuth2ClientForService(ctx context.Context, pgc database.Client, dbCfg *databasecfg.Config, oauth2Input *oauth.OAuth2ClientDatabaseCreationInput, provider string) (*oauth.OAuth2Client, error) {
	auditRepo := ProvideAuditLogRepository(provider, nil, nil, pgc)
	oauth2ClientManager := ProvideOAuthRepository(provider, nil, nil, auditRepo, dbCfg, pgc)

	createdClient, err := oauth2ClientManager.CreateOAuth2Client(ctx, oauth2Input)
	if err != nil {
		return nil, fmt.Errorf("failed to create oauth2 client: %w", err)
	}

	return createdClient, nil
}

func BuildInProcessServer(ctx context.Context, cfg *config.APIServiceConfig) (server *apiserver.Server, databaseClient database.Client, dbCfg *databasecfg.Config, err error) {
	pillars, err := cfg.Observability.ProvidePillars(ctx)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("setting up observability pillars: %w", err)
	}
	logger := pillars.Logger

	switch cfg.Database.Provider {
	case databasecfg.ProviderSQLite:
		// Create a temp directory for the SQLite database file.
		tmpDir, mkdirErr := os.MkdirTemp("", "zhuzh-sqlite-integration-*")
		if mkdirErr != nil {
			return nil, nil, nil, fmt.Errorf("creating temp dir for sqlite: %w", mkdirErr)
		}

		// busy_timeout pragma waits up to 5 s before returning SQLITE_BUSY.
		// modernc.org/sqlite uses _pragma=<pragma> DSN params, not _busy_timeout.
		// foreign_keys pragma enforces referential integrity (off by default in SQLite).
		dbDSN := filepath.Join(tmpDir, "integration_test.db") + "?_pragma=busy_timeout(5000)&_pragma=foreign_keys(on)"
		cfg.Database.Provider = databasecfg.ProviderSQLite
		cfg.Database.ReadConnection = databasecfg.ConnectionDetails{Database: dbDSN}
		cfg.Database.WriteConnection = databasecfg.ConnectionDetails{Database: dbDSN}
		cfg.Database.OAuth2TokenEncryptionKey = "blahblahblahblahblahblahblahblah"
		cfg.Database.UserDeviceTokenEncryptionKey = "blahblahblahblahblahblahblahblah"

		// Use noop message queue (the platform library falls through to noop for unrecognized providers).
		cfg.Events.Publisher.Provider = "noop"
		cfg.Events.Consumer.Provider = "noop"

		dbCfg = &cfg.Database
	default:
		// PostgreSQL path: spin up a container via testcontainers.
		redisConfig, _, redisErr := redis.BuildContainerBackedRedisConfig(ctx)
		if redisErr != nil {
			return nil, nil, nil, fmt.Errorf("connecting to redis: %w", redisErr)
		}
		cfg.Events.Publisher.Provider = msgconfig.ProviderRedis
		cfg.Events.Publisher.Redis = *redisConfig
		cfg.Events.Consumer.Redis = *redisConfig

		_, _, pgDBCfg, pgErr := pgtesting.BuildDatabaseContainer(ctx, "integration_testing")
		if pgErr != nil {
			return nil, nil, nil, fmt.Errorf("connecting to postgres: %w", pgErr)
		}
		cfg.Database.WriteConnection = pgDBCfg.WriteConnection
		cfg.Database.ReadConnection = pgDBCfg.ReadConnection
		cfg.Database.OAuth2TokenEncryptionKey = pgDBCfg.OAuth2TokenEncryptionKey
		cfg.Database.UserDeviceTokenEncryptionKey = pgDBCfg.UserDeviceTokenEncryptionKey

		dbCfg = &cfg.Database
	}

	tracerProvider := tracing.NewNoopTracerProvider()
	migrator := repositories.ProvideMigrator(&cfg.Database, logger)
	databaseClient, err = databasecfg.ProvideDatabase(ctx, logger, tracerProvider, &cfg.Database, migrator, nil)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("initializing database client: %w", err)
	}

	server, err = apiserver.NewServer(ctx, pillars, cfg)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("building API server: %w", err)
	}

	return server, databaseClient, dbCfg, nil
}

// DatabaseInitFunc is a function that performs database initialization operations.
// It receives the database client, config, logger, and tracer to perform arbitrary operations.
type DatabaseInitFunc func(ctx context.Context, dbClient database.Client, dbCfg *databasecfg.Config, logger logging.Logger, tracerProvider tracing.TracerProvider) error

// WithIdentityRepository provides an identity repository for custom operations.
func WithIdentityRepository(provider string, fn func(ctx context.Context, repo identity.Repository, logger logging.Logger, tracerProvider tracing.TracerProvider, dbClient database.Client) error) DatabaseInitFunc {
	return func(ctx context.Context, dbClient database.Client, dbCfg *databasecfg.Config, logger logging.Logger, tracerProvider tracing.TracerProvider) error {
		auditLogRepo := ProvideAuditLogRepository(provider, logger, tracerProvider, dbClient)
		identityRepo := ProvideIdentityRepository(provider, logger, tracerProvider, auditLogRepo, dbClient)
		return fn(ctx, identityRepo, logger, tracerProvider, dbClient)
	}
}

// WithOAuth2Repository provides an OAuth2 repository for custom operations.
func WithOAuth2Repository(provider string, fn func(ctx context.Context, repo oauth.Repository, logger logging.Logger, tracerProvider tracing.TracerProvider) error) DatabaseInitFunc {
	return func(ctx context.Context, dbClient database.Client, dbCfg *databasecfg.Config, logger logging.Logger, tracerProvider tracing.TracerProvider) error {
		auditLogRepo := ProvideAuditLogRepository(provider, logger, tracerProvider, dbClient)
		oauthRepo := ProvideOAuthRepository(provider, logger, tracerProvider, auditLogRepo, dbCfg, dbClient)
		return fn(ctx, oauthRepo, logger, tracerProvider)
	}
}

// WithAuthRepository provides an auth repository for custom operations.
func WithAuthRepository(provider string, fn func(ctx context.Context, repo auth.Repository, logger logging.Logger, tracerProvider tracing.TracerProvider) error) DatabaseInitFunc {
	return func(ctx context.Context, dbClient database.Client, dbCfg *databasecfg.Config, logger logging.Logger, tracerProvider tracing.TracerProvider) error {
		auditLogRepo := ProvideAuditLogRepository(provider, logger, tracerProvider, dbClient)
		var authRepo auth.Repository
		switch provider {
		case databasecfg.ProviderSQLite:
			authRepo = sqliteauth.ProvideAuthRepository(logger, tracerProvider, auditLogRepo, dbClient)
		default:
			authRepo = pgauth.ProvideAuthRepository(logger, tracerProvider, auditLogRepo, dbClient)
		}
		return fn(ctx, authRepo, logger, tracerProvider)
	}
}

// WithSettingsRepository provides a settings repository for custom operations.
func WithSettingsRepository(provider string, fn func(ctx context.Context, repo settings.Repository, logger logging.Logger, tracerProvider tracing.TracerProvider) error) DatabaseInitFunc {
	return func(ctx context.Context, dbClient database.Client, dbCfg *databasecfg.Config, logger logging.Logger, tracerProvider tracing.TracerProvider) error {
		auditLogRepo := ProvideAuditLogRepository(provider, logger, tracerProvider, dbClient)
		var settingsRepo settings.Repository
		switch provider {
		case databasecfg.ProviderSQLite:
			settingsRepo = sqlitesettings.ProvideSettingsRepository(logger, tracerProvider, auditLogRepo, dbClient)
		default:
			settingsRepo = pgsettings.ProvideSettingsRepository(logger, tracerProvider, auditLogRepo, dbClient)
		}
		return fn(ctx, settingsRepo, logger, tracerProvider)
	}
}

// WithWebhooksRepository provides a webhooks repository for custom operations.
func WithWebhooksRepository(provider string, fn func(ctx context.Context, repo webhooks.Repository, logger logging.Logger, tracerProvider tracing.TracerProvider) error) DatabaseInitFunc {
	return func(ctx context.Context, dbClient database.Client, dbCfg *databasecfg.Config, logger logging.Logger, tracerProvider tracing.TracerProvider) error {
		auditLogRepo := ProvideAuditLogRepository(provider, logger, tracerProvider, dbClient)
		var webhooksRepo webhooks.Repository
		switch provider {
		case databasecfg.ProviderSQLite:
			webhooksRepo = sqlitewebhooks.ProvideWebhooksRepository(logger, tracerProvider, auditLogRepo, dbClient)
		default:
			webhooksRepo = pgwebhooks.ProvideWebhooksRepository(logger, tracerProvider, auditLogRepo, dbClient)
		}
		return fn(ctx, webhooksRepo, logger, tracerProvider)
	}
}

// WithNotificationsRepository provides a notifications repository for custom operations.
func WithNotificationsRepository(provider string, fn func(ctx context.Context, repo notifications.Repository, logger logging.Logger, tracerProvider tracing.TracerProvider) error) DatabaseInitFunc {
	return func(ctx context.Context, dbClient database.Client, dbCfg *databasecfg.Config, logger logging.Logger, tracerProvider tracing.TracerProvider) error {
		auditLogRepo := ProvideAuditLogRepository(provider, logger, tracerProvider, dbClient)
		var notifsRepo notifications.Repository
		switch provider {
		case databasecfg.ProviderSQLite:
			notifsRepo = sqlitenotifications.ProvideNotificationsRepository(logger, tracerProvider, auditLogRepo, dbCfg, dbClient)
		default:
			notifsRepo = pgnotifications.ProvideNotificationsRepository(logger, tracerProvider, auditLogRepo, dbCfg, dbClient)
		}
		return fn(ctx, notifsRepo, logger, tracerProvider)
	}
}

// AllInOne sets up a complete local development environment with a docker-backed server,
// database, and runs the provided database initialization functions.
func AllInOne(ctx context.Context, cfg *config.APIServiceConfig, initFuncs ...DatabaseInitFunc) (*apiserver.Server, error) {
	server, databaseClient, dbCfg, err := BuildInProcessServer(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("building in-process server: %w", err)
	}

	log.Printf("%sDATABASE CONNECTION URL: %s%s", strings.Repeat("\n", 10), dbCfg.ReadConnection.URI(), strings.Repeat("\n", 10))

	pillars, err := cfg.Observability.ProvidePillars(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting o11y pillars: %w", err)
	}

	// Run all database initialization functions
	for i, initFunc := range initFuncs {
		if err = initFunc(ctx, databaseClient, dbCfg, pillars.Logger, pillars.TracerProvider); err != nil {
			return nil, fmt.Errorf("running database init function %d: %w", i, err)
		}
	}

	return server, nil
}

func BuildInsecureOAuthedGRPCClient(
	ctx context.Context,
	createdClientID,
	createdClientSecret,
	httpTestServerAddress,
	grpcServerAddress,
	token string,
) (client.Client, error) {
	state, err := random.GenerateBase64EncodedString(ctx, 32)
	if err != nil {
		return nil, fmt.Errorf("generating state: %w", err)
	}

	oauth2Config := oauth2.Config{
		ClientID:     createdClientID,
		ClientSecret: createdClientSecret,
		Scopes:       []string{"anything"}, // TODO: This should be nil-able
		RedirectURL:  httpTestServerAddress,
		Endpoint: oauth2.Endpoint{
			AuthStyle: oauth2.AuthStyleInParams,
			AuthURL:   httpTestServerAddress + "/oauth2/authorize",
			TokenURL:  httpTestServerAddress + "/oauth2/token",
		},
	}

	authCodeURL := oauth2Config.AuthCodeURL(
		state,
		oauth2.SetAuthURLParam("code_challenge_method", "plain"),
	)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		authCodeURL,
		http.NoBody,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get oauth2 code: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Location", "localhost")

	httpClient := httpclient.ProvideHTTPClient(&httpclient.Config{EnableTracing: true})
	httpClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	res, err := httpClient.Do(req) //nolint:gosec // G704: authCodeURL from OAuth config (httpTestServerAddress), not user-controlled
	if err != nil {
		return nil, fmt.Errorf("failed to get oauth2 code: %w", err)
	}
	defer func() {
		if err = res.Body.Close(); err != nil {
			log.Println("failed to close oauth2 response body", err)
		}
	}()

	const (
		codeKey = "code"
	)

	rl, err := res.Location()
	if err != nil {
		return nil, fmt.Errorf("getting location value from response: %w", err)
	}

	code := rl.Query().Get(codeKey)
	if code == "" {
		return nil, fmt.Errorf("code not returned from oauth2 redirect")
	}

	oauth2Token, err := oauth2Config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("exchanging OAuth2 code: %w", err)
	}

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithPerRPCCredentials(&insecureOAuth{
			TokenSource: oauth2Config.TokenSource(ctx, oauth2Token),
		}),
	}

	c, err := client.BuildClient(grpcServerAddress, opts...)
	if err != nil {
		return nil, fmt.Errorf("building client: %w", err)
	}

	return c, nil
}

// Custom insecure OAuth2 credentials that skip security checks.
type insecureOAuth struct {
	TokenSource oauth2.TokenSource
}

func (i *insecureOAuth) GetRequestMetadata(_ context.Context, _ ...string) (map[string]string, error) {
	token, err := i.TokenSource.Token()
	if err != nil {
		return nil, err
	}

	return map[string]string{"authorization": token.Type() + " " + token.AccessToken}, nil
}

func (i *insecureOAuth) RequireTransportSecurity() bool {
	return false // Explicitly allow insecure transport
}

func FetchLoginTokenForUser(ctx context.Context, grpcServerAddr string, loginInput *authsvc.UserLoginInput) (string, error) {
	unauthedClient, err := client.BuildUnauthenticatedGRPCClient(grpcServerAddr)
	if err != nil {
		return "", fmt.Errorf("initializing client: %w", err)
	}

	return FetchLoginTokenForUserWithClient(ctx, unauthedClient, loginInput)
}

// FetchLoginTokenForUserWithClient calls LoginForToken using the given client.
// Use this when the client must use TLS (e.g. for api.zhuzh.dev).
func FetchLoginTokenForUserWithClient(ctx context.Context, c client.Client, loginInput *authsvc.UserLoginInput) (string, error) {
	tokenRes, err := c.LoginForToken(ctx, &authsvc.LoginForTokenRequest{
		Input: loginInput,
	})
	if err != nil {
		return "", fmt.Errorf("fetching login token: %w", err)
	}

	return tokenRes.Result.AccessToken, nil
}

// FetchOAuth2TokenForUser performs the OAuth2 authorization code flow using the given JWT
// and returns the OAuth2 access and refresh tokens. Used by integration tests for token revocation.
func FetchOAuth2TokenForUser(
	ctx context.Context,
	httpServerAddress, grpcServerAddress, clientID, clientSecret string,
	loginInput *authsvc.UserLoginInput,
) (*oauth2.Token, error) {
	jwt, err := FetchLoginTokenForUser(ctx, grpcServerAddress, loginInput)
	if err != nil {
		return nil, fmt.Errorf("fetching JWT for OAuth2 exchange: %w", err)
	}

	state, err := random.GenerateBase64EncodedString(ctx, 32)
	if err != nil {
		return nil, fmt.Errorf("generating state: %w", err)
	}

	oauth2Config := oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       []string{"anything"},
		RedirectURL:  httpServerAddress,
		Endpoint: oauth2.Endpoint{
			AuthStyle: oauth2.AuthStyleInParams,
			AuthURL:   httpServerAddress + "/oauth2/authorize",
			TokenURL:  httpServerAddress + "/oauth2/token",
		},
	}

	authCodeURL := oauth2Config.AuthCodeURL(
		state,
		oauth2.SetAuthURLParam("code_challenge_method", "plain"),
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, authCodeURL, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("creating auth request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+jwt)
	req.Header.Set("Location", "localhost")

	httpClient := httpclient.ProvideHTTPClient(&httpclient.Config{EnableTracing: true})
	httpClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	res, err := httpClient.Do(req) //nolint:gosec // G704: authCodeURL from OAuth config
	if err != nil {
		return nil, fmt.Errorf("fetching OAuth2 code: %w", err)
	}
	defer func() {
		if err = res.Body.Close(); err != nil {
			log.Println("failed to close oauth2 response body", err)
		}
	}()

	rl, err := res.Location()
	if err != nil {
		return nil, fmt.Errorf("getting location from response: %w", err)
	}

	code := rl.Query().Get("code")
	if code == "" {
		return nil, fmt.Errorf("code not returned from oauth2 redirect")
	}

	oauth2Token, err := oauth2Config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("exchanging OAuth2 code: %w", err)
	}

	return oauth2Token, nil
}
