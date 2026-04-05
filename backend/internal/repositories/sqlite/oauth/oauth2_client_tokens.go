package oauth

import (
	"context"
	"time"

	types "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/oauth"
	oauthkeys "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/oauth/keys"
	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/repositories/sqlite/oauth/generated"

	platformerrors "github.com/verygoodsoftwarenotvirus/platform/v4/errors"
	"github.com/verygoodsoftwarenotvirus/platform/v4/observability"
	"github.com/verygoodsoftwarenotvirus/platform/v4/observability/tracing"
)

var _ types.OAuth2ClientTokenDataManager = (*repository)(nil)

func convertOAuth2ClientToken(result *generated.Oauth2ClientTokens) *types.OAuth2ClientToken {
	codeCreatedAt := parseTime(result.CodeCreatedAt)
	accessCreatedAt := parseTime(result.AccessCreatedAt)
	refreshCreatedAt := parseTime(result.RefreshCreatedAt)
	codeExpiresAt := parseTime(result.CodeExpiresAt)
	accessExpiresAt := parseTime(result.AccessExpiresAt)
	refreshExpiresAt := parseTime(result.RefreshExpiresAt)

	return &types.OAuth2ClientToken{
		RefreshCreatedAt:    refreshCreatedAt,
		AccessCreatedAt:     accessCreatedAt,
		CodeCreatedAt:       codeCreatedAt,
		RedirectURI:         result.RedirectUri,
		Code:                result.Code,
		CodeChallenge:       result.CodeChallenge,
		CodeChallengeMethod: result.CodeChallengeMethod,
		BelongsToUser:       result.BelongsToUser,
		Access:              result.Access,
		ClientID:            result.ClientID,
		Refresh:             result.Refresh,
		ID:                  result.ID,
		CodeExpiresAt:       codeExpiresAt.Sub(codeCreatedAt),
		AccessExpiresAt:     accessExpiresAt.Sub(accessCreatedAt),
		RefreshExpiresAt:    refreshExpiresAt.Sub(refreshCreatedAt),
	}
}

// GetOAuth2ClientTokenByCode fetches an OAuth2 client token from the database.
func (q *repository) GetOAuth2ClientTokenByCode(ctx context.Context, code string) (*types.OAuth2ClientToken, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if code == "" {
		return nil, platformerrors.ErrInvalidIDProvided
	}
	logger = logger.WithValue(oauthkeys.OAuth2ClientTokenCodeKey, code)
	tracing.AttachToSpan(span, oauthkeys.OAuth2ClientTokenCodeKey, code)

	encryptedCode, err := q.oauth2ClientTokenEncDec.Encrypt(ctx, code)
	if err != nil {
		return nil, observability.PrepareError(err, span, "encrypting oauth2 token code")
	}

	result, err := q.generatedQuerier.GetOAuth2ClientTokenByCode(ctx, q.readDB, encryptedCode)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "getting oauth2 client token by code")
	}

	oauth2ClientToken := convertOAuth2ClientToken(result)

	decryptedCode, err := q.oauth2ClientTokenEncDec.Decrypt(ctx, result.Code)
	if err != nil {
		return nil, observability.PrepareError(err, span, "decrypting oauth2 token code")
	}
	oauth2ClientToken.Code = decryptedCode

	decryptedAccess, err := q.oauth2ClientTokenEncDec.Decrypt(ctx, result.Access)
	if err != nil {
		return nil, observability.PrepareError(err, span, "decrypting oauth2 token access")
	}
	oauth2ClientToken.Access = decryptedAccess

	decryptedRefresh, err := q.oauth2ClientTokenEncDec.Decrypt(ctx, result.Refresh)
	if err != nil {
		return nil, observability.PrepareError(err, span, "decrypting oauth2 token refresh")
	}
	oauth2ClientToken.Refresh = decryptedRefresh

	return oauth2ClientToken, nil
}

// GetOAuth2ClientTokenByAccess fetches an OAuth2 client token from the database.
func (q *repository) GetOAuth2ClientTokenByAccess(ctx context.Context, access string) (*types.OAuth2ClientToken, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if access == "" {
		return nil, platformerrors.ErrInvalidIDProvided
	}
	logger = logger.WithValue(oauthkeys.OAuth2ClientTokenAccessKey, access)
	tracing.AttachToSpan(span, oauthkeys.OAuth2ClientTokenAccessKey, access)

	encryptedAccess, err := q.oauth2ClientTokenEncDec.Encrypt(ctx, access)
	if err != nil {
		return nil, observability.PrepareError(err, span, "encrypting oauth2 token access")
	}

	result, err := q.generatedQuerier.GetOAuth2ClientTokenByAccess(ctx, q.readDB, encryptedAccess)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "getting oauth2 client token by access")
	}

	oauth2ClientToken := convertOAuth2ClientToken(result)

	decryptedCode, err := q.oauth2ClientTokenEncDec.Decrypt(ctx, result.Code)
	if err != nil {
		return nil, observability.PrepareError(err, span, "decrypting oauth2 token code")
	}
	oauth2ClientToken.Code = decryptedCode

	decryptedAccess, err := q.oauth2ClientTokenEncDec.Decrypt(ctx, result.Access)
	if err != nil {
		return nil, observability.PrepareError(err, span, "decrypting oauth2 token access")
	}
	oauth2ClientToken.Access = decryptedAccess

	decryptedRefresh, err := q.oauth2ClientTokenEncDec.Decrypt(ctx, result.Refresh)
	if err != nil {
		return nil, observability.PrepareError(err, span, "decrypting oauth2 token refresh")
	}
	oauth2ClientToken.Refresh = decryptedRefresh

	return oauth2ClientToken, nil
}

// GetOAuth2ClientTokenByRefresh fetches an OAuth2 client token from the database.
func (q *repository) GetOAuth2ClientTokenByRefresh(ctx context.Context, refresh string) (*types.OAuth2ClientToken, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if refresh == "" {
		return nil, platformerrors.ErrInvalidIDProvided
	}
	logger = logger.WithValue(oauthkeys.OAuth2ClientTokenRefreshKey, refresh)
	tracing.AttachToSpan(span, oauthkeys.OAuth2ClientTokenRefreshKey, refresh)

	encryptedRefresh, err := q.oauth2ClientTokenEncDec.Encrypt(ctx, refresh)
	if err != nil {
		return nil, observability.PrepareError(err, span, "encrypting oauth2 token refresh")
	}

	result, err := q.generatedQuerier.GetOAuth2ClientTokenByRefresh(ctx, q.readDB, encryptedRefresh)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "getting oauth2 client token by refresh")
	}

	oauth2ClientToken := convertOAuth2ClientToken(result)

	decryptedCode, err := q.oauth2ClientTokenEncDec.Decrypt(ctx, result.Code)
	if err != nil {
		return nil, observability.PrepareError(err, span, "decrypting oauth2 token code")
	}
	oauth2ClientToken.Code = decryptedCode

	decryptedAccess, err := q.oauth2ClientTokenEncDec.Decrypt(ctx, result.Access)
	if err != nil {
		return nil, observability.PrepareError(err, span, "decrypting oauth2 token access")
	}
	oauth2ClientToken.Access = decryptedAccess

	decryptedRefresh, err := q.oauth2ClientTokenEncDec.Decrypt(ctx, result.Refresh)
	if err != nil {
		return nil, observability.PrepareError(err, span, "decrypting oauth2 token refresh")
	}
	oauth2ClientToken.Refresh = decryptedRefresh

	return oauth2ClientToken, nil
}

// CreateOAuth2ClientToken creates an OAuth2 client token in the database.
func (q *repository) CreateOAuth2ClientToken(ctx context.Context, input *types.OAuth2ClientTokenDatabaseCreationInput) (*types.OAuth2ClientToken, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, platformerrors.ErrNilInputProvided
	}

	logger := q.logger.WithValue(oauthkeys.OAuth2ClientTokenIDKey, input.ID)
	now := q.CurrentTime()

	encryptedCode, err := q.oauth2ClientTokenEncDec.Encrypt(ctx, input.Code)
	if err != nil {
		return nil, observability.PrepareError(err, span, "encrypting oauth2 token code")
	}

	encryptedAccess, err := q.oauth2ClientTokenEncDec.Encrypt(ctx, input.Access)
	if err != nil {
		return nil, observability.PrepareError(err, span, "encrypting oauth2 token access")
	}

	encryptedRefresh, err := q.oauth2ClientTokenEncDec.Encrypt(ctx, input.Refresh)
	if err != nil {
		return nil, observability.PrepareError(err, span, "encrypting oauth2 token refresh")
	}

	// create the oauth2 client token.
	if err = q.generatedQuerier.CreateOAuth2ClientToken(ctx, q.writeDB, &generated.CreateOAuth2ClientTokenParams{
		AccessExpiresAt:     now.Add(input.AccessExpiresIn).Format(time.RFC3339Nano),
		CodeExpiresAt:       now.Add(input.CodeExpiresIn).Format(time.RFC3339Nano),
		RefreshExpiresAt:    now.Add(input.RefreshExpiresIn).Format(time.RFC3339Nano),
		RefreshCreatedAt:    now.Format(time.RFC3339Nano),
		CodeCreatedAt:       now.Format(time.RFC3339Nano),
		AccessCreatedAt:     now.Format(time.RFC3339Nano),
		CodeChallenge:       input.CodeChallenge,
		CodeChallengeMethod: input.CodeChallengeMethod,
		ClientID:            input.ClientID,
		Access:              encryptedAccess,
		Code:                encryptedCode,
		ID:                  input.ID,
		Refresh:             encryptedRefresh,
		RedirectUri:         input.RedirectURI,
		BelongsToUser:       input.BelongsToUser,
	}); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "performing oauth2 client token creation query")
	}

	logger.Info("oauth2 client token created")

	oauth2ClientToken := &types.OAuth2ClientToken{
		RefreshCreatedAt:    input.RefreshCreatedAt,
		AccessCreatedAt:     input.AccessCreatedAt,
		CodeCreatedAt:       input.CodeCreatedAt,
		RedirectURI:         input.RedirectURI,
		Code:                input.Code,
		CodeChallenge:       input.CodeChallenge,
		CodeChallengeMethod: input.CodeChallengeMethod,
		BelongsToUser:       input.BelongsToUser,
		Access:              input.Access,
		ClientID:            input.ClientID,
		Refresh:             input.Refresh,
		ID:                  input.ID,
		CodeExpiresAt:       input.CodeExpiresIn,
		AccessExpiresAt:     input.AccessExpiresIn,
		RefreshExpiresAt:    input.RefreshExpiresIn,
	}

	return oauth2ClientToken, nil
}

// DeleteOAuth2ClientTokenByAccess deletes an OAuth2 client token from the database by its access token.
func (q *repository) DeleteOAuth2ClientTokenByAccess(ctx context.Context, access string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if access == "" {
		return platformerrors.ErrInvalidIDProvided
	}
	logger = logger.WithValue(oauthkeys.OAuth2ClientTokenAccessKey, access)
	tracing.AttachToSpan(span, oauthkeys.OAuth2ClientTokenAccessKey, access)

	encryptedAccess, err := q.oauth2ClientTokenEncDec.Encrypt(ctx, access)
	if err != nil {
		return observability.PrepareError(err, span, "encrypting oauth2 token access")
	}

	if _, err = q.generatedQuerier.DeleteOAuth2ClientTokenByAccess(ctx, q.writeDB, encryptedAccess); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "deleting oauth2 client token by access")
	}

	return nil
}

// DeleteOAuth2ClientTokenByCode deletes an OAuth2 client token from the database by its authorization code.
func (q *repository) DeleteOAuth2ClientTokenByCode(ctx context.Context, code string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if code == "" {
		return platformerrors.ErrInvalidIDProvided
	}
	logger = logger.WithValue(oauthkeys.OAuth2ClientTokenCodeKey, code)
	tracing.AttachToSpan(span, oauthkeys.OAuth2ClientTokenCodeKey, code)

	encryptedCode, err := q.oauth2ClientTokenEncDec.Encrypt(ctx, code)
	if err != nil {
		return observability.PrepareError(err, span, "encrypting oauth2 token code")
	}

	if _, err = q.generatedQuerier.DeleteOAuth2ClientTokenByCode(ctx, q.writeDB, encryptedCode); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "deleting oauth2 client token by code")
	}

	return nil
}

// DeleteOAuth2ClientTokenByRefresh deletes an OAuth2 client token from the database by its refresh token.
func (q *repository) DeleteOAuth2ClientTokenByRefresh(ctx context.Context, refresh string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if refresh == "" {
		return platformerrors.ErrInvalidIDProvided
	}
	logger = logger.WithValue(oauthkeys.OAuth2ClientTokenRefreshKey, refresh)
	tracing.AttachToSpan(span, oauthkeys.OAuth2ClientTokenRefreshKey, refresh)

	encryptedRefresh, err := q.oauth2ClientTokenEncDec.Encrypt(ctx, refresh)
	if err != nil {
		return observability.PrepareError(err, span, "encrypting oauth2 token refresh")
	}

	_, err = q.generatedQuerier.DeleteOAuth2ClientTokenByRefresh(ctx, q.writeDB, encryptedRefresh)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "deleting oauth2 client token by refresh")
	}

	return nil
}
