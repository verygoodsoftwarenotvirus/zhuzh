package sessions

import (
	"context"

	"github.com/samber/do/v2"
)

// ProvideContextDataFetcherFromContext provides a function that fetches context data from context.
func ProvideContextDataFetcherFromContext() func(context.Context) (*ContextData, error) {
	return FetchContextDataFromContext
}

// RegisterSessionProviders registers session providers with the injector.
func RegisterSessionProviders(i do.Injector) {
	do.Provide[func(context.Context) (*ContextData, error)](i, func(i do.Injector) (func(context.Context) (*ContextData, error), error) {
		return ProvideContextDataFetcherFromContext(), nil
	})
}
