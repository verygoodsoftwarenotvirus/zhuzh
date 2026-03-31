package manager

import (
	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/notifications"
)

// NotificationsDataManager defines the interface for notifications business logic.
// It embeds notifications.Repository so the manager provides the full repository surface.
type NotificationsDataManager interface {
	notifications.Repository
}
