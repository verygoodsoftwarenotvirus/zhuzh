package manager

import (
	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/uploadedmedia"
)

// UploadedMediaManager defines the interface for uploaded media business logic.
// It embeds uploadedmedia.Repository so the manager provides the full repository surface.
type UploadedMediaManager interface {
	uploadedmedia.Repository
}
