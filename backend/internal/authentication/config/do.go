package authcfg

import (
	tokenscfg "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/authentication/tokens/config"

	"github.com/samber/do/v2"
)

// RegisterConfigs registers auth config sub-fields with the injector.
func RegisterConfigs(i do.Injector) {
	do.Provide[*tokenscfg.Config](i, func(i do.Injector) (*tokenscfg.Config, error) {
		cfg := do.MustInvoke[*Config](i)
		return &cfg.Tokens, nil
	})
	do.Provide[SSOConfigs](i, func(i do.Injector) (SSOConfigs, error) {
		cfg := do.MustInvoke[*Config](i)
		return cfg.SSO, nil
	})
}
