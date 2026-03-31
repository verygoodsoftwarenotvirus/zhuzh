package manager

import (
	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/waitlists"
)

// WaitlistsDataManager defines the interface for waitlists business logic.
// It embeds waitlists.Repository so the manager provides the full repository surface.
type WaitlistsDataManager interface {
	waitlists.Repository
}
