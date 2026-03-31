package grpc

import (
	"context"

	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/authentication/webauthn"
	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/identity"
	identitymanager "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/identity/manager"
)

// passkeyUserStore adapts identityDataManager to webauthn.UserStore.
type passkeyUserStore struct {
	identityDataManager identitymanager.IdentityDataManager
}

func (s *passkeyUserStore) GetUserByID(ctx context.Context, userID string) (*identity.User, error) {
	return s.identityDataManager.GetUser(ctx, userID)
}

func (s *passkeyUserStore) GetUserByUsername(ctx context.Context, username string) (*identity.User, error) {
	return s.identityDataManager.GetUserByUsername(ctx, username)
}

// ProvidePasskeyService creates a WebAuthn passkey service.
func ProvidePasskeyService(
	cfg webauthn.Config,
	identityDataManager identitymanager.IdentityDataManager,
	identityRepo identity.Repository,
	sessionStore webauthn.SessionStore,
) (*webauthn.Service, error) {
	userStore := &passkeyUserStore{identityDataManager: identityDataManager}
	return webauthn.NewService(cfg, identityRepo, userStore, sessionStore)
}
