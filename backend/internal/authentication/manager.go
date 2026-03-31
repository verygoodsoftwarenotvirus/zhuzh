package authentication

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/authentication/tokens"
	tokenscfg "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/authentication/tokens/config"
	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/audit"
	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/auth"
	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/identity"
	identitykeys "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/identity/keys"

	"github.com/verygoodsoftwarenotvirus/platform/v4/identifiers"
	"github.com/verygoodsoftwarenotvirus/platform/v4/messagequeue"
	msgconfig "github.com/verygoodsoftwarenotvirus/platform/v4/messagequeue/config"
	"github.com/verygoodsoftwarenotvirus/platform/v4/observability"
	"github.com/verygoodsoftwarenotvirus/platform/v4/observability/logging"
	"github.com/verygoodsoftwarenotvirus/platform/v4/observability/tracing"
)

const (
	name = "authentication_manager"
)

type (
	// LoginMetadata holds request metadata for session tracking.
	LoginMetadata struct {
		ClientIP  string
		UserAgent string
	}

	Manager interface {
		ProcessLogin(ctx context.Context, adminOnly bool, loginData *auth.UserLoginInput, meta *LoginMetadata) (*auth.TokenResponse, error)
		ProcessPasskeyLogin(ctx context.Context, userID, desiredAccountID string, meta *LoginMetadata) (*auth.TokenResponse, error)
		ExchangeTokenForUser(ctx context.Context, refreshToken, desiredAccountID string) (*auth.TokenResponse, error)
	}

	manager struct {
		tokenIssuer             tokens.Issuer
		authenticator           Authenticator
		tracer                  tracing.Tracer
		logger                  logging.Logger
		dataChangesPublisher    messagequeue.Publisher
		userAuthDataManager     identity.Repository
		sessionDataManager      auth.UserSessionDataManager
		maxAccessTokenLifetime  time.Duration
		maxRefreshTokenLifetime time.Duration
	}
)

func NewManager(
	ctx context.Context,
	queuesConfig *msgconfig.QueuesConfig,
	tokenIssuer tokens.Issuer,
	authenticator Authenticator,
	tracingProvider tracing.TracerProvider,
	logger logging.Logger,
	publisherProvider messagequeue.PublisherProvider,
	userAuthDataManager identity.Repository,
	sessionDataManager auth.UserSessionDataManager,
	cfg *tokenscfg.Config,
) (Manager, error) {
	dataChangesPublisher, err := publisherProvider.ProvidePublisher(ctx, queuesConfig.DataChangesTopicName)
	if err != nil {
		return nil, observability.PrepareError(err, nil, "creating data changes publisher")
	}

	m := &manager{
		maxRefreshTokenLifetime: cfg.MaxRefreshTokenLifetime,
		maxAccessTokenLifetime:  cfg.MaxAccessTokenLifetime,
		tracer:                  tracing.NewNamedTracer(tracingProvider, name),
		logger:                  logging.NewNamedLogger(logger, name),
		tokenIssuer:             tokenIssuer,
		authenticator:           authenticator,
		dataChangesPublisher:    dataChangesPublisher,
		userAuthDataManager:     userAuthDataManager,
		sessionDataManager:      sessionDataManager,
	}

	return m, nil
}

// validateLogin takes login information and returns whether the login is valid.
// In the event that there's an error, this function will return false and the error.
func (m *manager) validateLogin(ctx context.Context, user *identity.User, loginInput *auth.UserLoginInput) (bool, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	loginInput.TOTPToken = strings.TrimSpace(loginInput.TOTPToken)
	loginInput.Password = strings.TrimSpace(loginInput.Password)
	loginInput.Username = strings.TrimSpace(loginInput.Username)

	// alias the relevant data.
	logger := m.logger.WithValue(identitykeys.UsernameKey, user.Username)

	// check for login validity.
	loginValid, err := m.authenticator.CredentialsAreValid(
		ctx,
		user.HashedPassword,
		loginInput.Password,
		user.TwoFactorSecret,
		loginInput.TOTPToken,
	)

	if errors.Is(err, ErrInvalidTOTPToken) || errors.Is(err, ErrPasswordDoesNotMatch) {
		return false, err
	}

	if err != nil {
		return false, observability.PrepareError(err, span, "validating login")
	}

	logger.Debug("login validated")

	return loginValid, nil
}

func (m *manager) ProcessLogin(ctx context.Context, adminOnly bool, loginData *auth.UserLoginInput, meta *LoginMetadata) (*auth.TokenResponse, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.Clone()

	if err := loginData.ValidateWithContext(ctx); err != nil {
		return nil, observability.PrepareError(err, span, "validating input")
	}

	logger = logger.WithValue(identitykeys.UsernameKey, loginData.Username)

	userFunc := m.userAuthDataManager.GetUserByUsername
	if adminOnly {
		userFunc = m.userAuthDataManager.GetAdminUserByUsername
	}

	user, err := userFunc(ctx, loginData.Username)
	if err != nil || user == nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, observability.PrepareError(err, span, "user does not exist")
		}

		return nil, observability.PrepareError(err, span, "fetching user")
	}

	logger = logger.WithValue(identitykeys.UserIDKey, user.ID)
	tracing.AttachToSpan(span, identitykeys.UserIDKey, user.ID)

	if user.IsBanned() {
		return nil, observability.PrepareError(err, span, "user is banned")
	}

	loginValid, err := m.validateLogin(ctx, user, loginData)
	logger.WithValue("login_valid", loginValid)

	if err != nil {
		if errors.Is(err, ErrInvalidTOTPToken) {
			return nil, observability.PrepareError(err, span, "invalid TOTP AccessToken")
		}

		if errors.Is(err, ErrPasswordDoesNotMatch) {
			return nil, observability.PrepareError(err, span, "password did not match")
		}

		return nil, observability.PrepareError(err, span, "validating login")
	} else if !loginValid {
		return nil, observability.PrepareError(err, span, "login was invalid")
	}

	if user.TwoFactorSecretVerifiedAt != nil && loginData.TOTPToken == "" {
		return nil, observability.PrepareError(ErrTOTPRequired, span, "processing login")
	}

	var accountID string
	if loginData.DesiredAccountID != "" {
		var isMember bool
		isMember, err = m.userAuthDataManager.UserIsMemberOfAccount(ctx, user.ID, loginData.DesiredAccountID)
		if err != nil {
			return nil, observability.PrepareError(err, span, "validating account membership")
		}
		if !isMember {
			return nil, observability.PrepareError(errors.New("user does not have access to account"), span, "user does not have access to the desired account")
		}
		accountID = loginData.DesiredAccountID
	} else {
		var defaultAccountID string
		defaultAccountID, err = m.userAuthDataManager.GetDefaultAccountIDForUser(ctx, user.ID)
		if err != nil {
			return nil, observability.PrepareError(err, span, "validating input")
		}
		accountID = defaultAccountID
	}

	response, err := m.issueTokensWithSession(ctx, user, accountID, auth.LoginMethodPassword, meta)
	if err != nil {
		return nil, observability.PrepareError(err, span, "issuing tokens with session")
	}

	dcm := &audit.DataChangeMessage{
		EventType: identity.UserLoggedInServiceEventType,
		AccountID: accountID,
		UserID:    user.ID,
	}

	if err = m.dataChangesPublisher.Publish(ctx, dcm); err != nil {
		return nil, observability.PrepareError(err, span, "publishing data change")
	}

	return response, nil
}

// ProcessPasskeyLogin issues tokens for a user authenticated via passkey.
func (m *manager) ProcessPasskeyLogin(ctx context.Context, userID, desiredAccountID string, meta *LoginMetadata) (*auth.TokenResponse, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	tracing.AttachToSpan(span, identitykeys.UserIDKey, userID)

	user, err := m.userAuthDataManager.GetUser(ctx, userID)
	if err != nil || user == nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, observability.PrepareError(err, span, "user does not exist")
		}
		return nil, observability.PrepareError(err, span, "fetching user")
	}

	if user.IsBanned() {
		return nil, observability.PrepareError(errors.New("user is banned"), span, "user is banned")
	}

	var accountID string
	if desiredAccountID != "" {
		var isMember bool
		isMember, err = m.userAuthDataManager.UserIsMemberOfAccount(ctx, user.ID, desiredAccountID)
		if err != nil {
			return nil, observability.PrepareError(err, span, "validating account membership")
		}
		if !isMember {
			return nil, observability.PrepareError(errors.New("user does not have access to account"), span, "user does not have access to the desired account")
		}
		accountID = desiredAccountID
	} else {
		var defaultAccountID string
		defaultAccountID, err = m.userAuthDataManager.GetDefaultAccountIDForUser(ctx, user.ID)
		if err != nil {
			return nil, observability.PrepareError(err, span, "validating input")
		}
		accountID = defaultAccountID
	}

	response, err := m.issueTokensWithSession(ctx, user, accountID, auth.LoginMethodPasskey, meta)
	if err != nil {
		return nil, observability.PrepareError(err, span, "issuing tokens with session")
	}

	dcm := &audit.DataChangeMessage{
		EventType: identity.UserLoggedInServiceEventType,
		AccountID: accountID,
		UserID:    user.ID,
	}

	if err = m.dataChangesPublisher.Publish(ctx, dcm); err != nil {
		return nil, observability.PrepareError(err, span, "publishing data change")
	}

	return response, nil
}

func (m *manager) ExchangeTokenForUser(ctx context.Context, refreshToken, desiredAccountID string) (*auth.TokenResponse, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.Clone()

	userID, err := m.tokenIssuer.ParseUserIDFromToken(ctx, refreshToken)
	if err != nil {
		return nil, observability.PrepareError(err, span, "parsing userID from token")
	}

	user, err := m.userAuthDataManager.GetUser(ctx, userID)
	if err != nil || user == nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, observability.PrepareError(err, span, "user does not exist")
		}

		return nil, observability.PrepareError(err, span, "fetching user")
	}

	logger = logger.WithValue(identitykeys.UserIDKey, user.ID)
	tracing.AttachToSpan(span, identitykeys.UserIDKey, user.ID)

	if user.IsBanned() {
		return nil, observability.PrepareError(err, span, "user is banned")
	}

	// Validate the existing session via refresh token JTI.
	refreshJTI, jtiErr := m.tokenIssuer.ParseJTIFromToken(ctx, refreshToken)
	sessionID, sidErr := m.tokenIssuer.ParseSessionIDFromToken(ctx, refreshToken)

	if jtiErr != nil {
		refreshJTI = ""
	}
	if sidErr != nil {
		sessionID = ""
	}

	if refreshJTI != "" && sessionID != "" {
		session, sessErr := m.sessionDataManager.GetUserSessionByRefreshTokenID(ctx, refreshJTI)
		if sessErr != nil || session == nil {
			return nil, observability.PrepareError(errors.New("session has been revoked"), span, "validating session")
		}
	}

	var accountID string
	if desiredAccountID != "" {
		var isMember bool
		isMember, err = m.userAuthDataManager.UserIsMemberOfAccount(ctx, user.ID, desiredAccountID)
		if err != nil {
			return nil, observability.PrepareError(err, span, "validating account membership")
		}
		if !isMember {
			return nil, observability.PrepareError(errors.New("user does not have access to account"), span, "user does not have access to the desired account")
		}
		accountID = desiredAccountID
	} else {
		var defaultAccountID string
		defaultAccountID, err = m.userAuthDataManager.GetDefaultAccountIDForUser(ctx, user.ID)
		if err != nil {
			return nil, observability.PrepareError(err, span, "validating input")
		}
		accountID = defaultAccountID
	}

	// Issue new tokens with the same session ID.
	var accessJTI, refreshJTINew string
	response := &auth.TokenResponse{
		UserID:     user.ID,
		AccountID:  accountID,
		ExpiresUTC: time.Now().Add(m.maxAccessTokenLifetime).UTC(),
	}

	response.AccessToken, accessJTI, err = m.tokenIssuer.IssueToken(ctx, user, m.maxAccessTokenLifetime, accountID, sessionID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating access token")
	}

	response.RefreshToken, refreshJTINew, err = m.tokenIssuer.IssueToken(ctx, user, m.maxRefreshTokenLifetime, accountID, sessionID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating refresh token")
	}

	// Update the session with new token JTIs.
	if sessionID != "" {
		if updateErr := m.sessionDataManager.UpdateSessionTokenIDs(ctx, sessionID, accessJTI, refreshJTINew, time.Now().Add(m.sessionExpiryDuration()).UTC()); updateErr != nil {
			logger.Error("updating session token IDs", updateErr)
		}
	}

	dcm := &audit.DataChangeMessage{
		EventType: identity.UserLoggedInServiceEventType,
		AccountID: accountID,
		UserID:    user.ID,
	}

	if err = m.dataChangesPublisher.Publish(ctx, dcm); err != nil {
		return nil, observability.PrepareError(err, span, "publishing data change")
	}

	return response, nil
}

// issueTokensWithSession creates a session record and issues tokens with the session ID embedded.
func (m *manager) issueTokensWithSession(ctx context.Context, user *identity.User, accountID, loginMethod string, meta *LoginMetadata) (*auth.TokenResponse, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	sessionID := identifiers.New()

	var clientIP, userAgent string
	if meta != nil {
		clientIP = meta.ClientIP
		userAgent = meta.UserAgent
	}

	response := &auth.TokenResponse{
		UserID:     user.ID,
		AccountID:  accountID,
		ExpiresUTC: time.Now().Add(m.maxAccessTokenLifetime).UTC(),
	}

	var (
		err                   error
		accessJTI, refreshJTI string
	)
	response.AccessToken, accessJTI, err = m.tokenIssuer.IssueToken(ctx, user, m.maxAccessTokenLifetime, accountID, sessionID)
	if err != nil {
		return nil, observability.PrepareError(err, span, "creating access token")
	}

	response.RefreshToken, refreshJTI, err = m.tokenIssuer.IssueToken(ctx, user, m.maxRefreshTokenLifetime, accountID, sessionID)
	if err != nil {
		return nil, observability.PrepareError(err, span, "creating refresh token")
	}

	if _, err = m.sessionDataManager.CreateUserSession(ctx, &auth.UserSessionDatabaseCreationInput{
		ID:             sessionID,
		BelongsToUser:  user.ID,
		SessionTokenID: accessJTI,
		RefreshTokenID: refreshJTI,
		ClientIP:       clientIP,
		UserAgent:      userAgent,
		DeviceName:     deriveDeviceName(userAgent),
		LoginMethod:    loginMethod,
		ExpiresAt:      time.Now().Add(m.sessionExpiryDuration()).UTC(),
	}); err != nil {
		m.logger.Error("creating user session", err)
	}

	return response, nil
}

const defaultSessionExpiry = 72 * time.Hour

// sessionExpiryDuration returns the session expiry duration, defaulting to 72h if the
// configured refresh token lifetime is zero (e.g. in test configs).
func (m *manager) sessionExpiryDuration() time.Duration {
	if m.maxRefreshTokenLifetime > 0 {
		return m.maxRefreshTokenLifetime
	}
	return defaultSessionExpiry
}

// deriveDeviceName produces a simple friendly device name from a User-Agent string.
func deriveDeviceName(userAgent string) string {
	if userAgent == "" {
		return "Unknown Device"
	}

	ua := strings.ToLower(userAgent)
	switch {
	case strings.Contains(ua, "iphone"):
		return "iPhone"
	case strings.Contains(ua, "ipad"):
		return "iPad"
	case strings.Contains(ua, "android"):
		return "Android Device"
	case strings.Contains(ua, "macintosh") || strings.Contains(ua, "mac os"):
		return "Mac"
	case strings.Contains(ua, "windows"):
		return "Windows PC"
	case strings.Contains(ua, "linux"):
		return "Linux"
	default:
		return "Unknown Device"
	}
}
