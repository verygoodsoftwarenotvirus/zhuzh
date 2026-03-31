package authentication

import (
	"context"
	"errors"
	"testing"
	"time"

	mocktokens "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/authentication/tokens/mock"
	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/auth"
	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/identity"
	identitymock "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/identity/mock"

	"github.com/verygoodsoftwarenotvirus/platform/v4/database/filtering"
	mockpublishers "github.com/verygoodsoftwarenotvirus/platform/v4/messagequeue/mock"
	"github.com/verygoodsoftwarenotvirus/platform/v4/observability/logging"
	"github.com/verygoodsoftwarenotvirus/platform/v4/observability/tracing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// mockAuthenticator is a local mock for the Authenticator interface (same package).
type mockAuthenticator struct {
	mock.Mock
}

func (m *mockAuthenticator) CredentialsAreValid(ctx context.Context, hash, password, totpSecret, totpCode string) (bool, error) {
	args := m.Called(ctx, hash, password, totpSecret, totpCode)
	return args.Bool(0), args.Error(1)
}

func (m *mockAuthenticator) HashPassword(ctx context.Context, password string) (string, error) {
	args := m.Called(ctx, password)
	return args.String(0), args.Error(1)
}

// mockSessionDataManager is a local mock for auth.UserSessionDataManager.
type mockSessionDataManager struct {
	mock.Mock
}

func (m *mockSessionDataManager) CreateUserSession(ctx context.Context, input *auth.UserSessionDatabaseCreationInput) (*auth.UserSession, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*auth.UserSession), args.Error(1)
}

func (m *mockSessionDataManager) GetUserSessionBySessionTokenID(ctx context.Context, sessionTokenID string) (*auth.UserSession, error) {
	args := m.Called(ctx, sessionTokenID)
	return args.Get(0).(*auth.UserSession), args.Error(1)
}

func (m *mockSessionDataManager) GetUserSessionByRefreshTokenID(ctx context.Context, refreshTokenID string) (*auth.UserSession, error) {
	args := m.Called(ctx, refreshTokenID)
	return args.Get(0).(*auth.UserSession), args.Error(1)
}

func (m *mockSessionDataManager) GetActiveSessionsForUser(ctx context.Context, userID string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[auth.UserSession], error) {
	args := m.Called(ctx, userID, filter)
	return args.Get(0).(*filtering.QueryFilteredResult[auth.UserSession]), args.Error(1)
}

func (m *mockSessionDataManager) RevokeUserSession(ctx context.Context, sessionID, userID string) error {
	return m.Called(ctx, sessionID, userID).Error(0)
}

func (m *mockSessionDataManager) RevokeAllSessionsForUser(ctx context.Context, userID string) error {
	return m.Called(ctx, userID).Error(0)
}

func (m *mockSessionDataManager) RevokeAllSessionsForUserExcept(ctx context.Context, userID, sessionID string) error {
	return m.Called(ctx, userID, sessionID).Error(0)
}

func (m *mockSessionDataManager) UpdateSessionTokenIDs(ctx context.Context, sessionID, newSessionTokenID, newRefreshTokenID string, newExpiresAt time.Time) error {
	return m.Called(ctx, sessionID, newSessionTokenID, newRefreshTokenID, newExpiresAt).Error(0)
}

func (m *mockSessionDataManager) TouchSessionLastActive(ctx context.Context, sessionTokenID string) error {
	return m.Called(ctx, sessionTokenID).Error(0)
}

type managerTestMocks struct {
	tokenIssuer         *mocktokens.Issuer
	authenticator       *mockAuthenticator
	userAuthDataManager *identitymock.RepositoryMock
	sessionDataManager  *mockSessionDataManager
	publisher           *mockpublishers.Publisher
}

// helper to build a minimal manager for testing.
func buildTestManager(t *testing.T) (*manager, *managerTestMocks) {
	t.Helper()

	mocks := &managerTestMocks{
		tokenIssuer:         &mocktokens.Issuer{},
		authenticator:       &mockAuthenticator{},
		userAuthDataManager: &identitymock.RepositoryMock{},
		sessionDataManager:  &mockSessionDataManager{},
		publisher:           &mockpublishers.Publisher{},
	}

	m := &manager{
		tokenIssuer:             mocks.tokenIssuer,
		authenticator:           mocks.authenticator,
		tracer:                  tracing.NewNamedTracer(tracing.NewNoopTracerProvider(), "test"),
		logger:                  logging.NewNoopLogger(),
		dataChangesPublisher:    mocks.publisher,
		userAuthDataManager:     mocks.userAuthDataManager,
		sessionDataManager:      mocks.sessionDataManager,
		maxAccessTokenLifetime:  15 * time.Minute,
		maxRefreshTokenLifetime: 24 * time.Hour,
	}

	return m, mocks
}

func buildExampleUser() *identity.User {
	return &identity.User{
		ID:             "user123",
		Username:       "testuser",
		HashedPassword: "hashedpassword",
		AccountStatus:  string(identity.GoodStandingUserAccountStatus),
		EmailAddress:   "test@example.com",
		FirstName:      "Test",
		LastName:       "User",
	}
}

func Test_deriveDeviceName(T *testing.T) {
	T.Parallel()

	tests := []struct {
		name      string
		userAgent string
		expected  string
	}{
		{
			name:      "empty user agent",
			userAgent: "",
			expected:  "Unknown Device",
		},
		{
			name:      "iPhone",
			userAgent: "Mozilla/5.0 (iPhone; CPU iPhone OS 16_0 like Mac OS X)",
			expected:  "iPhone",
		},
		{
			name:      "iPad",
			userAgent: "Mozilla/5.0 (iPad; CPU OS 16_0 like Mac OS X)",
			expected:  "iPad",
		},
		{
			name:      "Android",
			userAgent: "Mozilla/5.0 (Linux; Android 13; Pixel 7)",
			expected:  "Android Device",
		},
		{
			name:      "Mac via Macintosh",
			userAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7)",
			expected:  "Mac",
		},
		{
			name:      "Mac via Mac OS",
			userAgent: "Mozilla/5.0 (compatible; Mac OS X 12_0)",
			expected:  "Mac",
		},
		{
			name:      "Windows",
			userAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64)",
			expected:  "Windows PC",
		},
		{
			name:      "Linux",
			userAgent: "Mozilla/5.0 (X11; Linux x86_64)",
			expected:  "Linux",
		},
		{
			name:      "unknown user agent",
			userAgent: "SomeCustomBot/1.0",
			expected:  "Unknown Device",
		},
	}

	for _, tc := range tests {
		T.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			actual := deriveDeviceName(tc.userAgent)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestManager_ProcessLogin(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		m, mocks := buildTestManager(t)

		user := buildExampleUser()
		loginInput := &auth.UserLoginInput{
			Username: "testuser",
			Password: "validP@ssw0rd",
		}

		mocks.userAuthDataManager.On("GetUserByUsername", mock.Anything, loginInput.Username).Return(user, nil)
		mocks.authenticator.On("CredentialsAreValid", mock.Anything, user.HashedPassword, loginInput.Password, user.TwoFactorSecret, loginInput.TOTPToken).Return(true, nil)
		mocks.userAuthDataManager.On("GetDefaultAccountIDForUser", mock.Anything, user.ID).Return("account123", nil)

		// issueTokensWithSession: two IssueToken calls (access + refresh)
		mocks.tokenIssuer.On("IssueToken", mock.Anything, mock.AnythingOfType("*identity.User"), m.maxAccessTokenLifetime, "account123", mock.AnythingOfType("string")).Return("access-token", "access-jti", nil).Once()
		mocks.tokenIssuer.On("IssueToken", mock.Anything, mock.AnythingOfType("*identity.User"), m.maxRefreshTokenLifetime, "account123", mock.AnythingOfType("string")).Return("refresh-token", "refresh-jti", nil).Once()

		mocks.sessionDataManager.On("CreateUserSession", mock.Anything, mock.AnythingOfType("*auth.UserSessionDatabaseCreationInput")).Return(&auth.UserSession{}, nil)

		mocks.publisher.On("Publish", mock.Anything, mock.Anything).Return(nil)

		response, err := m.ProcessLogin(ctx, false, loginInput, &LoginMetadata{
			ClientIP:  "127.0.0.1",
			UserAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7)",
		})

		require.NoError(t, err)
		require.NotNil(t, response)
		assert.Equal(t, "access-token", response.AccessToken)
		assert.Equal(t, "refresh-token", response.RefreshToken)
		assert.Equal(t, user.ID, response.UserID)
		assert.Equal(t, "account123", response.AccountID)

		mock.AssertExpectationsForObjects(t, mocks.tokenIssuer, mocks.authenticator, mocks.userAuthDataManager, mocks.sessionDataManager, mocks.publisher)
	})

	T.Run("with desired account ID", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		m, mocks := buildTestManager(t)

		user := buildExampleUser()
		loginInput := &auth.UserLoginInput{
			Username:         "testuser",
			Password:         "validP@ssw0rd",
			DesiredAccountID: "specific-account",
		}

		mocks.userAuthDataManager.On("GetUserByUsername", mock.Anything, loginInput.Username).Return(user, nil)
		mocks.authenticator.On("CredentialsAreValid", mock.Anything, user.HashedPassword, loginInput.Password, user.TwoFactorSecret, loginInput.TOTPToken).Return(true, nil)
		mocks.userAuthDataManager.On("UserIsMemberOfAccount", mock.Anything, user.ID, "specific-account").Return(true, nil)

		mocks.tokenIssuer.On("IssueToken", mock.Anything, mock.AnythingOfType("*identity.User"), m.maxAccessTokenLifetime, "specific-account", mock.AnythingOfType("string")).Return("access-token", "access-jti", nil).Once()
		mocks.tokenIssuer.On("IssueToken", mock.Anything, mock.AnythingOfType("*identity.User"), m.maxRefreshTokenLifetime, "specific-account", mock.AnythingOfType("string")).Return("refresh-token", "refresh-jti", nil).Once()

		mocks.sessionDataManager.On("CreateUserSession", mock.Anything, mock.AnythingOfType("*auth.UserSessionDatabaseCreationInput")).Return(&auth.UserSession{}, nil)
		mocks.publisher.On("Publish", mock.Anything, mock.Anything).Return(nil)

		response, err := m.ProcessLogin(ctx, false, loginInput, nil)

		require.NoError(t, err)
		require.NotNil(t, response)
		assert.Equal(t, "specific-account", response.AccountID)

		mock.AssertExpectationsForObjects(t, mocks.tokenIssuer, mocks.authenticator, mocks.userAuthDataManager, mocks.sessionDataManager, mocks.publisher)
	})

	T.Run("with invalid credentials", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		m, mocks := buildTestManager(t)

		user := buildExampleUser()
		loginInput := &auth.UserLoginInput{
			Username: "testuser",
			Password: "wrongP@ssw0rd",
		}

		mocks.userAuthDataManager.On("GetUserByUsername", mock.Anything, loginInput.Username).Return(user, nil)
		mocks.authenticator.On("CredentialsAreValid", mock.Anything, user.HashedPassword, loginInput.Password, user.TwoFactorSecret, loginInput.TOTPToken).Return(false, ErrPasswordDoesNotMatch)

		response, err := m.ProcessLogin(ctx, false, loginInput, nil)

		assert.Error(t, err)
		assert.Nil(t, response)

		mock.AssertExpectationsForObjects(t, mocks.authenticator, mocks.userAuthDataManager)
	})

	T.Run("with banned user", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		m, mocks := buildTestManager(t)

		user := buildExampleUser()
		user.AccountStatus = string(identity.BannedUserAccountStatus)

		loginInput := &auth.UserLoginInput{
			Username: "testuser",
			Password: "validP@ssw0rd",
		}

		mocks.userAuthDataManager.On("GetUserByUsername", mock.Anything, loginInput.Username).Return(user, nil)

		response, err := m.ProcessLogin(ctx, false, loginInput, nil)

		// NOTE: the production code calls observability.PrepareError(err, ...) where err is nil
		// from the successful GetUserByUsername call, so PrepareError returns nil. This means
		// banned users currently get (nil, nil) back rather than an error. This is likely a bug.
		assert.NoError(t, err)
		assert.Nil(t, response)

		mock.AssertExpectationsForObjects(t, mocks.userAuthDataManager)
	})

	T.Run("with nonexistent user", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		m, mocks := buildTestManager(t)

		loginInput := &auth.UserLoginInput{
			Username: "nouser",
			Password: "validP@ssw0rd",
		}

		mocks.userAuthDataManager.On("GetUserByUsername", mock.Anything, loginInput.Username).Return((*identity.User)(nil), errors.New("not found"))

		response, err := m.ProcessLogin(ctx, false, loginInput, nil)

		assert.Error(t, err)
		assert.Nil(t, response)

		mock.AssertExpectationsForObjects(t, mocks.userAuthDataManager)
	})

	T.Run("with invalid login input", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		m, _ := buildTestManager(t)

		// empty username and password should fail validation
		loginInput := &auth.UserLoginInput{
			Username: "",
			Password: "",
		}

		response, err := m.ProcessLogin(ctx, false, loginInput, nil)

		assert.Error(t, err)
		assert.Nil(t, response)
	})

	T.Run("with TOTP required but not provided", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		m, mocks := buildTestManager(t)

		now := time.Now()
		user := buildExampleUser()
		user.TwoFactorSecretVerifiedAt = &now
		user.TwoFactorSecret = "ASECRET"

		loginInput := &auth.UserLoginInput{
			Username: "testuser",
			Password: "validP@ssw0rd",
			// TOTPToken intentionally left empty
		}

		mocks.userAuthDataManager.On("GetUserByUsername", mock.Anything, loginInput.Username).Return(user, nil)
		mocks.authenticator.On("CredentialsAreValid", mock.Anything, user.HashedPassword, loginInput.Password, user.TwoFactorSecret, "").Return(true, nil)

		response, err := m.ProcessLogin(ctx, false, loginInput, nil)

		assert.Error(t, err)
		assert.Nil(t, response)

		mock.AssertExpectationsForObjects(t, mocks.authenticator, mocks.userAuthDataManager)
	})

	T.Run("with user not member of desired account", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		m, mocks := buildTestManager(t)

		user := buildExampleUser()
		loginInput := &auth.UserLoginInput{
			Username:         "testuser",
			Password:         "validP@ssw0rd",
			DesiredAccountID: "other-account",
		}

		mocks.userAuthDataManager.On("GetUserByUsername", mock.Anything, loginInput.Username).Return(user, nil)
		mocks.authenticator.On("CredentialsAreValid", mock.Anything, user.HashedPassword, loginInput.Password, user.TwoFactorSecret, loginInput.TOTPToken).Return(true, nil)
		mocks.userAuthDataManager.On("UserIsMemberOfAccount", mock.Anything, user.ID, "other-account").Return(false, nil)

		response, err := m.ProcessLogin(ctx, false, loginInput, nil)

		assert.Error(t, err)
		assert.Nil(t, response)

		mock.AssertExpectationsForObjects(t, mocks.authenticator, mocks.userAuthDataManager)
	})

	T.Run("admin only", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		m, mocks := buildTestManager(t)

		user := buildExampleUser()
		loginInput := &auth.UserLoginInput{
			Username: "testuser",
			Password: "validP@ssw0rd",
		}

		mocks.userAuthDataManager.On("GetAdminUserByUsername", mock.Anything, loginInput.Username).Return(user, nil)
		mocks.authenticator.On("CredentialsAreValid", mock.Anything, user.HashedPassword, loginInput.Password, user.TwoFactorSecret, loginInput.TOTPToken).Return(true, nil)
		mocks.userAuthDataManager.On("GetDefaultAccountIDForUser", mock.Anything, user.ID).Return("account123", nil)

		mocks.tokenIssuer.On("IssueToken", mock.Anything, mock.AnythingOfType("*identity.User"), m.maxAccessTokenLifetime, "account123", mock.AnythingOfType("string")).Return("access-token", "access-jti", nil).Once()
		mocks.tokenIssuer.On("IssueToken", mock.Anything, mock.AnythingOfType("*identity.User"), m.maxRefreshTokenLifetime, "account123", mock.AnythingOfType("string")).Return("refresh-token", "refresh-jti", nil).Once()

		mocks.sessionDataManager.On("CreateUserSession", mock.Anything, mock.AnythingOfType("*auth.UserSessionDatabaseCreationInput")).Return(&auth.UserSession{}, nil)
		mocks.publisher.On("Publish", mock.Anything, mock.Anything).Return(nil)

		response, err := m.ProcessLogin(ctx, true, loginInput, nil)

		require.NoError(t, err)
		require.NotNil(t, response)

		mock.AssertExpectationsForObjects(t, mocks.tokenIssuer, mocks.authenticator, mocks.userAuthDataManager, mocks.sessionDataManager, mocks.publisher)
	})
}

func TestManager_ProcessPasskeyLogin(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		m, mocks := buildTestManager(t)

		user := buildExampleUser()

		mocks.userAuthDataManager.On("GetUser", mock.Anything, user.ID).Return(user, nil)
		mocks.userAuthDataManager.On("GetDefaultAccountIDForUser", mock.Anything, user.ID).Return("account123", nil)

		mocks.tokenIssuer.On("IssueToken", mock.Anything, mock.AnythingOfType("*identity.User"), m.maxAccessTokenLifetime, "account123", mock.AnythingOfType("string")).Return("access-token", "access-jti", nil).Once()
		mocks.tokenIssuer.On("IssueToken", mock.Anything, mock.AnythingOfType("*identity.User"), m.maxRefreshTokenLifetime, "account123", mock.AnythingOfType("string")).Return("refresh-token", "refresh-jti", nil).Once()

		mocks.sessionDataManager.On("CreateUserSession", mock.Anything, mock.AnythingOfType("*auth.UserSessionDatabaseCreationInput")).
			Run(func(args mock.Arguments) {
				input := args.Get(1).(*auth.UserSessionDatabaseCreationInput)
				assert.Equal(t, auth.LoginMethodPasskey, input.LoginMethod)
			}).
			Return(&auth.UserSession{}, nil)

		mocks.publisher.On("Publish", mock.Anything, mock.Anything).Return(nil)

		response, err := m.ProcessPasskeyLogin(ctx, user.ID, "", &LoginMetadata{
			ClientIP:  "10.0.0.1",
			UserAgent: "Mozilla/5.0 (iPhone; CPU iPhone OS 16_0 like Mac OS X)",
		})

		require.NoError(t, err)
		require.NotNil(t, response)
		assert.Equal(t, "access-token", response.AccessToken)
		assert.Equal(t, "refresh-token", response.RefreshToken)
		assert.Equal(t, user.ID, response.UserID)
		assert.Equal(t, "account123", response.AccountID)

		mock.AssertExpectationsForObjects(t, mocks.tokenIssuer, mocks.userAuthDataManager, mocks.sessionDataManager, mocks.publisher)
	})

	T.Run("with desired account ID", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		m, mocks := buildTestManager(t)

		user := buildExampleUser()

		mocks.userAuthDataManager.On("GetUser", mock.Anything, user.ID).Return(user, nil)
		mocks.userAuthDataManager.On("UserIsMemberOfAccount", mock.Anything, user.ID, "specific-account").Return(true, nil)

		mocks.tokenIssuer.On("IssueToken", mock.Anything, mock.AnythingOfType("*identity.User"), m.maxAccessTokenLifetime, "specific-account", mock.AnythingOfType("string")).Return("access-token", "access-jti", nil).Once()
		mocks.tokenIssuer.On("IssueToken", mock.Anything, mock.AnythingOfType("*identity.User"), m.maxRefreshTokenLifetime, "specific-account", mock.AnythingOfType("string")).Return("refresh-token", "refresh-jti", nil).Once()

		mocks.sessionDataManager.On("CreateUserSession", mock.Anything, mock.AnythingOfType("*auth.UserSessionDatabaseCreationInput")).Return(&auth.UserSession{}, nil)
		mocks.publisher.On("Publish", mock.Anything, mock.Anything).Return(nil)

		response, err := m.ProcessPasskeyLogin(ctx, user.ID, "specific-account", nil)

		require.NoError(t, err)
		require.NotNil(t, response)
		assert.Equal(t, "specific-account", response.AccountID)

		mock.AssertExpectationsForObjects(t, mocks.tokenIssuer, mocks.userAuthDataManager, mocks.sessionDataManager, mocks.publisher)
	})

	T.Run("with banned user", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		m, mocks := buildTestManager(t)

		user := buildExampleUser()
		user.AccountStatus = string(identity.BannedUserAccountStatus)

		mocks.userAuthDataManager.On("GetUser", mock.Anything, user.ID).Return(user, nil)

		response, err := m.ProcessPasskeyLogin(ctx, user.ID, "", nil)

		assert.Error(t, err)
		assert.Nil(t, response)

		mock.AssertExpectationsForObjects(t, mocks.userAuthDataManager)
	})

	T.Run("with nonexistent user", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		m, mocks := buildTestManager(t)

		mocks.userAuthDataManager.On("GetUser", mock.Anything, "nonexistent").Return((*identity.User)(nil), errors.New("not found"))

		response, err := m.ProcessPasskeyLogin(ctx, "nonexistent", "", nil)

		assert.Error(t, err)
		assert.Nil(t, response)

		mock.AssertExpectationsForObjects(t, mocks.userAuthDataManager)
	})

	T.Run("with user not member of desired account", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		m, mocks := buildTestManager(t)

		user := buildExampleUser()

		mocks.userAuthDataManager.On("GetUser", mock.Anything, user.ID).Return(user, nil)
		mocks.userAuthDataManager.On("UserIsMemberOfAccount", mock.Anything, user.ID, "other-account").Return(false, nil)

		response, err := m.ProcessPasskeyLogin(ctx, user.ID, "other-account", nil)

		assert.Error(t, err)
		assert.Nil(t, response)

		mock.AssertExpectationsForObjects(t, mocks.userAuthDataManager)
	})
}

func TestManager_ExchangeTokenForUser(T *testing.T) {
	T.Parallel()

	T.Run("standard with session validation", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		m, mocks := buildTestManager(t)

		user := buildExampleUser()
		refreshToken := "valid-refresh-token"

		// ParseUserIDFromToken delegates to ParseUserIDAndAccountIDFromToken on the mock
		mocks.tokenIssuer.On("ParseUserIDAndAccountIDFromToken", mock.Anything, refreshToken).Return(user.ID, "account123", nil)
		mocks.userAuthDataManager.On("GetUser", mock.Anything, user.ID).Return(user, nil)

		mocks.tokenIssuer.On("ParseJTIFromToken", mock.Anything, refreshToken).Return("refresh-jti-old", nil)
		mocks.tokenIssuer.On("ParseSessionIDFromToken", mock.Anything, refreshToken).Return("session-abc", nil)

		mocks.sessionDataManager.On("GetUserSessionByRefreshTokenID", mock.Anything, "refresh-jti-old").Return(&auth.UserSession{
			ID:             "session-abc",
			BelongsToUser:  user.ID,
			RefreshTokenID: "refresh-jti-old",
		}, nil)

		mocks.userAuthDataManager.On("GetDefaultAccountIDForUser", mock.Anything, user.ID).Return("account123", nil)

		mocks.tokenIssuer.On("IssueToken", mock.Anything, mock.AnythingOfType("*identity.User"), m.maxAccessTokenLifetime, "account123", "session-abc").Return("new-access-token", "new-access-jti", nil).Once()
		mocks.tokenIssuer.On("IssueToken", mock.Anything, mock.AnythingOfType("*identity.User"), m.maxRefreshTokenLifetime, "account123", "session-abc").Return("new-refresh-token", "new-refresh-jti", nil).Once()

		mocks.sessionDataManager.On("UpdateSessionTokenIDs", mock.Anything, "session-abc", "new-access-jti", "new-refresh-jti", mock.AnythingOfType("time.Time")).Return(nil)

		mocks.publisher.On("Publish", mock.Anything, mock.Anything).Return(nil)

		response, err := m.ExchangeTokenForUser(ctx, refreshToken, "")

		require.NoError(t, err)
		require.NotNil(t, response)
		assert.Equal(t, "new-access-token", response.AccessToken)
		assert.Equal(t, "new-refresh-token", response.RefreshToken)
		assert.Equal(t, user.ID, response.UserID)
		assert.Equal(t, "account123", response.AccountID)

		mock.AssertExpectationsForObjects(t, mocks.tokenIssuer, mocks.userAuthDataManager, mocks.sessionDataManager, mocks.publisher)
	})

	T.Run("with revoked session", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		m, mocks := buildTestManager(t)

		user := buildExampleUser()
		refreshToken := "revoked-refresh-token"

		mocks.tokenIssuer.On("ParseUserIDAndAccountIDFromToken", mock.Anything, refreshToken).Return(user.ID, "account123", nil)
		mocks.userAuthDataManager.On("GetUser", mock.Anything, user.ID).Return(user, nil)

		mocks.tokenIssuer.On("ParseJTIFromToken", mock.Anything, refreshToken).Return("old-jti", nil)
		mocks.tokenIssuer.On("ParseSessionIDFromToken", mock.Anything, refreshToken).Return("session-abc", nil)

		// Session not found means it was revoked
		mocks.sessionDataManager.On("GetUserSessionByRefreshTokenID", mock.Anything, "old-jti").Return((*auth.UserSession)(nil), errors.New("not found"))

		response, err := m.ExchangeTokenForUser(ctx, refreshToken, "")

		assert.Error(t, err)
		assert.Nil(t, response)

		mock.AssertExpectationsForObjects(t, mocks.tokenIssuer, mocks.userAuthDataManager, mocks.sessionDataManager)
	})

	T.Run("with desired account ID", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		m, mocks := buildTestManager(t)

		user := buildExampleUser()
		refreshToken := "valid-refresh-token"

		mocks.tokenIssuer.On("ParseUserIDAndAccountIDFromToken", mock.Anything, refreshToken).Return(user.ID, "account123", nil)
		mocks.userAuthDataManager.On("GetUser", mock.Anything, user.ID).Return(user, nil)

		mocks.tokenIssuer.On("ParseJTIFromToken", mock.Anything, refreshToken).Return("refresh-jti", nil)
		mocks.tokenIssuer.On("ParseSessionIDFromToken", mock.Anything, refreshToken).Return("session-abc", nil)

		mocks.sessionDataManager.On("GetUserSessionByRefreshTokenID", mock.Anything, "refresh-jti").Return(&auth.UserSession{
			ID:             "session-abc",
			RefreshTokenID: "refresh-jti",
		}, nil)

		mocks.userAuthDataManager.On("UserIsMemberOfAccount", mock.Anything, user.ID, "desired-account").Return(true, nil)

		mocks.tokenIssuer.On("IssueToken", mock.Anything, mock.AnythingOfType("*identity.User"), m.maxAccessTokenLifetime, "desired-account", "session-abc").Return("new-access-token", "new-access-jti", nil).Once()
		mocks.tokenIssuer.On("IssueToken", mock.Anything, mock.AnythingOfType("*identity.User"), m.maxRefreshTokenLifetime, "desired-account", "session-abc").Return("new-refresh-token", "new-refresh-jti", nil).Once()

		mocks.sessionDataManager.On("UpdateSessionTokenIDs", mock.Anything, "session-abc", "new-access-jti", "new-refresh-jti", mock.AnythingOfType("time.Time")).Return(nil)

		mocks.publisher.On("Publish", mock.Anything, mock.Anything).Return(nil)

		response, err := m.ExchangeTokenForUser(ctx, refreshToken, "desired-account")

		require.NoError(t, err)
		require.NotNil(t, response)
		assert.Equal(t, "desired-account", response.AccountID)

		mock.AssertExpectationsForObjects(t, mocks.tokenIssuer, mocks.userAuthDataManager, mocks.sessionDataManager, mocks.publisher)
	})

	T.Run("with banned user", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		m, mocks := buildTestManager(t)

		user := buildExampleUser()
		user.AccountStatus = string(identity.BannedUserAccountStatus)
		refreshToken := "valid-refresh-token"

		mocks.tokenIssuer.On("ParseUserIDAndAccountIDFromToken", mock.Anything, refreshToken).Return(user.ID, "account123", nil)
		mocks.userAuthDataManager.On("GetUser", mock.Anything, user.ID).Return(user, nil)

		response, err := m.ExchangeTokenForUser(ctx, refreshToken, "")

		// NOTE: same as ProcessLogin, the production code calls observability.PrepareError(err, ...)
		// where err is nil from the successful GetUser call, so PrepareError returns nil.
		// Banned users currently get (nil, nil) back rather than an error. This is likely a bug.
		assert.NoError(t, err)
		assert.Nil(t, response)

		mock.AssertExpectationsForObjects(t, mocks.tokenIssuer, mocks.userAuthDataManager)
	})

	T.Run("without JTI or session ID in token", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		m, mocks := buildTestManager(t)

		user := buildExampleUser()
		refreshToken := "legacy-refresh-token"

		mocks.tokenIssuer.On("ParseUserIDAndAccountIDFromToken", mock.Anything, refreshToken).Return(user.ID, "account123", nil)
		mocks.userAuthDataManager.On("GetUser", mock.Anything, user.ID).Return(user, nil)

		// No JTI or session ID in token (legacy token)
		mocks.tokenIssuer.On("ParseJTIFromToken", mock.Anything, refreshToken).Return("", errors.New("no jti"))
		mocks.tokenIssuer.On("ParseSessionIDFromToken", mock.Anything, refreshToken).Return("", errors.New("no session id"))

		mocks.userAuthDataManager.On("GetDefaultAccountIDForUser", mock.Anything, user.ID).Return("account123", nil)

		// Session ID is empty, so tokens are issued with empty session ID
		mocks.tokenIssuer.On("IssueToken", mock.Anything, mock.AnythingOfType("*identity.User"), m.maxAccessTokenLifetime, "account123", "").Return("new-access-token", "new-access-jti", nil).Once()
		mocks.tokenIssuer.On("IssueToken", mock.Anything, mock.AnythingOfType("*identity.User"), m.maxRefreshTokenLifetime, "account123", "").Return("new-refresh-token", "new-refresh-jti", nil).Once()

		// UpdateSessionTokenIDs should NOT be called because sessionID is empty

		mocks.publisher.On("Publish", mock.Anything, mock.Anything).Return(nil)

		response, err := m.ExchangeTokenForUser(ctx, refreshToken, "")

		require.NoError(t, err)
		require.NotNil(t, response)
		assert.Equal(t, "new-access-token", response.AccessToken)

		mock.AssertExpectationsForObjects(t, mocks.tokenIssuer, mocks.userAuthDataManager, mocks.publisher)
	})

	T.Run("with invalid refresh token", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		m, mocks := buildTestManager(t)

		refreshToken := "bad-token"

		mocks.tokenIssuer.On("ParseUserIDAndAccountIDFromToken", mock.Anything, refreshToken).Return("", "", errors.New("invalid token"))

		response, err := m.ExchangeTokenForUser(ctx, refreshToken, "")

		assert.Error(t, err)
		assert.Nil(t, response)

		mock.AssertExpectationsForObjects(t, mocks.tokenIssuer)
	})

	T.Run("with nonexistent user", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		m, mocks := buildTestManager(t)

		refreshToken := "valid-refresh-token"

		mocks.tokenIssuer.On("ParseUserIDAndAccountIDFromToken", mock.Anything, refreshToken).Return("nonexistent-user", "account123", nil)
		mocks.userAuthDataManager.On("GetUser", mock.Anything, "nonexistent-user").Return((*identity.User)(nil), errors.New("not found"))

		response, err := m.ExchangeTokenForUser(ctx, refreshToken, "")

		assert.Error(t, err)
		assert.Nil(t, response)

		mock.AssertExpectationsForObjects(t, mocks.tokenIssuer, mocks.userAuthDataManager)
	})

	T.Run("with user not member of desired account", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		m, mocks := buildTestManager(t)

		user := buildExampleUser()
		refreshToken := "valid-refresh-token"

		mocks.tokenIssuer.On("ParseUserIDAndAccountIDFromToken", mock.Anything, refreshToken).Return(user.ID, "account123", nil)
		mocks.userAuthDataManager.On("GetUser", mock.Anything, user.ID).Return(user, nil)

		mocks.tokenIssuer.On("ParseJTIFromToken", mock.Anything, refreshToken).Return("jti", nil)
		mocks.tokenIssuer.On("ParseSessionIDFromToken", mock.Anything, refreshToken).Return("session-abc", nil)

		mocks.sessionDataManager.On("GetUserSessionByRefreshTokenID", mock.Anything, "jti").Return(&auth.UserSession{
			ID:             "session-abc",
			RefreshTokenID: "jti",
		}, nil)

		mocks.userAuthDataManager.On("UserIsMemberOfAccount", mock.Anything, user.ID, "wrong-account").Return(false, nil)

		response, err := m.ExchangeTokenForUser(ctx, refreshToken, "wrong-account")

		assert.Error(t, err)
		assert.Nil(t, response)

		mock.AssertExpectationsForObjects(t, mocks.tokenIssuer, mocks.userAuthDataManager, mocks.sessionDataManager)
	})
}
