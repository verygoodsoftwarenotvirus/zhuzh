package repositories

import (
	domaindataprivacy "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/dataprivacy"
	pgaudit "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/repositories/postgres/auditlogentries"
	pgauth "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/repositories/postgres/auth"
	pgcomments "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/repositories/postgres/comments"
	pgdataprivacy "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/repositories/postgres/dataprivacy"
	pgidentity "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/repositories/postgres/identity"
	pginternalops "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/repositories/postgres/internalops"
	pgissuereports "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/repositories/postgres/issuereports"
	pgnotifications "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/repositories/postgres/notifications"
	pgoauth "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/repositories/postgres/oauth"
	pgpayments "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/repositories/postgres/payments"
	pgsettings "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/repositories/postgres/settings"
	pguploadedmedia "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/repositories/postgres/uploadedmedia"
	pgwaitlists "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/repositories/postgres/waitlists"
	pgwebhooks "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/repositories/postgres/webhooks"
	sqliteaudit "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/repositories/sqlite/auditlogentries"
	sqliteauth "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/repositories/sqlite/auth"
	sqlitecomments "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/repositories/sqlite/comments"
	sqlitedataprivacy "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/repositories/sqlite/dataprivacy"
	sqliteidentity "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/repositories/sqlite/identity"
	sqliteinternalops "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/repositories/sqlite/internalops"
	sqliteissuereports "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/repositories/sqlite/issuereports"
	sqlitenotifications "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/repositories/sqlite/notifications"
	sqliteoauth "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/repositories/sqlite/oauth"
	sqlitepayments "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/repositories/sqlite/payments"
	sqlitesettings "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/repositories/sqlite/settings"
	sqliteuploadedmedia "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/repositories/sqlite/uploadedmedia"
	sqlitewaitlists "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/repositories/sqlite/waitlists"
	sqlitewebhooks "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/repositories/sqlite/webhooks"

	databasecfg "github.com/verygoodsoftwarenotvirus/platform/v4/database/config"

	"github.com/samber/do/v2"
)

// RegisterRepositories registers the appropriate repository implementations
// based on the configured database provider.
func RegisterRepositories(i do.Injector, provider string) {
	do.Provide[[]domaindataprivacy.UserDataCollector](i, func(_ do.Injector) ([]domaindataprivacy.UserDataCollector, error) {
		return []domaindataprivacy.UserDataCollector{}, nil
	})

	switch provider {
	case databasecfg.ProviderSQLite:
		sqliteaudit.RegisterAuditLogRepository(i)
		sqliteauth.RegisterAuthRepository(i)
		sqlitecomments.RegisterCommentsRepository(i)
		sqliteidentity.RegisterIdentityRepository(i)
		sqliteissuereports.RegisterIssueReportsRepository(i)
		sqlitenotifications.RegisterNotificationsRepository(i)
		sqliteuploadedmedia.RegisterUploadedMediaRepository(i)
		sqlitewaitlists.RegisterWaitlistsRepository(i)
		sqlitewebhooks.RegisterWebhooksRepository(i)
		sqliteoauth.RegisterOAuthRepository(i)
		sqlitepayments.RegisterPaymentsRepository(i)
		sqlitesettings.RegisterSettingsRepository(i)
		sqlitedataprivacy.RegisterDataPrivacyRepository(i)
		sqliteinternalops.RegisterInternalOpsRepository(i)
	default:
		pgaudit.RegisterAuditLogRepository(i)
		pgauth.RegisterAuthRepository(i)
		pgcomments.RegisterCommentsRepository(i)
		pgidentity.RegisterIdentityRepository(i)
		pgissuereports.RegisterIssueReportsRepository(i)
		pgnotifications.RegisterNotificationsRepository(i)
		pguploadedmedia.RegisterUploadedMediaRepository(i)
		pgwaitlists.RegisterWaitlistsRepository(i)
		pgwebhooks.RegisterWebhooksRepository(i)
		pgoauth.RegisterOAuthRepository(i)
		pgpayments.RegisterPaymentsRepository(i)
		pgsettings.RegisterSettingsRepository(i)
		pgdataprivacy.RegisterDataPrivacyRepository(i)
		pginternalops.RegisterInternalOpsRepository(i)
	}
}
