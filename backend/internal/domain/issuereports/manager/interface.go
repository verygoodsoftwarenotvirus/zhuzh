package manager

import (
	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/issuereports"
)

// IssueReportsDataManager defines the interface for issue reports business logic.
// It embeds issuereports.Repository so the manager provides the full repository surface.
type IssueReportsDataManager interface {
	issuereports.Repository
}
