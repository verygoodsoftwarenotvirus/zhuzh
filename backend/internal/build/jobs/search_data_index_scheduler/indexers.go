package searchdataindexscheduler

import (
	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/identity"
	coreindexing "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/services/identity/indexing"

	"github.com/verygoodsoftwarenotvirus/platform/v4/search/text/indexing"
)

func ProvideIndexFunctions(identityRepo identity.Repository) map[string]indexing.Function {
	return coreindexing.BuildCoreDataIndexingFunctions(identityRepo)
}
