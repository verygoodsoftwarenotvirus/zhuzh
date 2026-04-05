package auditlogentries

import (
	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/audit"

	"github.com/verygoodsoftwarenotvirus/platform/v4/database"
	"github.com/verygoodsoftwarenotvirus/platform/v4/observability/logging"
	"github.com/verygoodsoftwarenotvirus/platform/v4/observability/tracing"

	"github.com/samber/do/v2"
)

// RegisterAuditLogRepository registers the audit log repository with the injector.
func RegisterAuditLogRepository(i do.Injector) {
	do.Provide[audit.Repository](i, func(i do.Injector) (audit.Repository, error) {
		return ProvideAuditLogRepository(
			do.MustInvoke[logging.Logger](i),
			do.MustInvoke[tracing.TracerProvider](i),
			do.MustInvoke[database.Client](i),
		), nil
	})
}
