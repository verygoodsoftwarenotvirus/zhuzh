package notifications

import (
	"context"
	"database/sql"

	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/audit"
	identitykeys "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/identity/keys"
	types "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/notifications"
	notificationkeys "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/notifications/keys"
	generated "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/repositories/sqlite/notifications/generated"

	"github.com/verygoodsoftwarenotvirus/platform/v4/database/filtering"
	platformerrors "github.com/verygoodsoftwarenotvirus/platform/v4/errors"
	"github.com/verygoodsoftwarenotvirus/platform/v4/identifiers"
	"github.com/verygoodsoftwarenotvirus/platform/v4/observability"
	"github.com/verygoodsoftwarenotvirus/platform/v4/observability/tracing"
)

const (
	resourceTypeUserDeviceTokens = "user_device_tokens"
)

var (
	_ types.UserDeviceTokenDataManager = (*Repository)(nil)
)

// UserDeviceTokenExists fetches whether a user device token exists from the database.
func (q *Repository) UserDeviceTokenExists(ctx context.Context, userID, tokenID string) (bool, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if userID == "" {
		return false, platformerrors.ErrInvalidIDProvided
	}
	logger = logger.WithValue(identitykeys.UserIDKey, userID)
	tracing.AttachToSpan(span, identitykeys.UserIDKey, userID)

	if tokenID == "" {
		return false, platformerrors.ErrInvalidIDProvided
	}
	logger = logger.WithValue(notificationkeys.UserDeviceTokenIDKey, tokenID)
	tracing.AttachToSpan(span, notificationkeys.UserDeviceTokenIDKey, tokenID)

	result, err := q.generatedQuerier.CheckUserDeviceTokenExistence(ctx, q.readDB, &generated.CheckUserDeviceTokenExistenceParams{
		ID:            tokenID,
		BelongsToUser: userID,
	})
	if err != nil {
		return false, observability.PrepareAndLogError(err, logger, span, "performing user device token existence check")
	}

	return result == "1", nil
}

// GetUserDeviceToken fetches a user device token from the database.
func (q *Repository) GetUserDeviceToken(ctx context.Context, userID, tokenID string) (*types.UserDeviceToken, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if userID == "" {
		return nil, platformerrors.ErrInvalidIDProvided
	}
	logger = logger.WithValue(identitykeys.UserIDKey, userID)
	tracing.AttachToSpan(span, identitykeys.UserIDKey, userID)

	if tokenID == "" {
		return nil, platformerrors.ErrInvalidIDProvided
	}
	logger = logger.WithValue(notificationkeys.UserDeviceTokenIDKey, tokenID)
	tracing.AttachToSpan(span, notificationkeys.UserDeviceTokenIDKey, tokenID)

	result, err := q.generatedQuerier.GetUserDeviceToken(ctx, q.readDB, &generated.GetUserDeviceTokenParams{
		BelongsToUser: userID,
		ID:            tokenID,
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching user device token")
	}

	decryptedDeviceToken, err := q.userDeviceTokenEncDec.Decrypt(ctx, result.DeviceToken)
	if err != nil {
		return nil, observability.PrepareError(err, span, "decrypting user device token")
	}

	return &types.UserDeviceToken{
		CreatedAt:     parseTime(result.CreatedAt),
		LastUpdatedAt: parseTimePtr(result.LastUpdatedAt),
		ArchivedAt:    parseTimePtr(result.ArchivedAt),
		ID:            result.ID,
		DeviceToken:   decryptedDeviceToken,
		Platform:      result.Platform,
		BelongsToUser: result.BelongsToUser,
	}, nil
}

// GetUserDeviceTokens fetches a list of user device tokens from the database that meet a particular filter.
func (q *Repository) GetUserDeviceTokens(ctx context.Context, userID string, filter *filtering.QueryFilter, platformFilter *string) (*filtering.QueryFilteredResult[types.UserDeviceToken], error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if userID == "" {
		return nil, platformerrors.ErrInvalidIDProvided
	}
	logger = logger.WithValue(identitykeys.UserIDKey, userID)
	tracing.AttachToSpan(span, identitykeys.UserIDKey, userID)

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	var platformFilterVal any
	if platformFilter != nil {
		platformFilterVal = *platformFilter
	}

	results, err := q.generatedQuerier.GetUserDeviceTokensForUser(ctx, q.readDB, &generated.GetUserDeviceTokensForUserParams{
		UserID:         userID,
		PlatformFilter: platformFilterVal,
		CreatedAfter:   timePtrToStringPtr(filter.CreatedAfter),
		CreatedBefore:  timePtrToStringPtr(filter.CreatedBefore),
		Cursor:         filter.Cursor,
		ResultLimit:    int64PtrFromUint8Ptr(filter.MaxResponseSize),
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing user device tokens list retrieval query")
	}

	var (
		data                      = []*types.UserDeviceToken{}
		filteredCount, totalCount uint64
	)
	for _, result := range results {
		var decryptedDeviceToken string
		decryptedDeviceToken, err = q.userDeviceTokenEncDec.Decrypt(ctx, result.DeviceToken)
		if err != nil {
			return nil, observability.PrepareError(err, span, "decrypting user device token")
		}
		token := &types.UserDeviceToken{
			CreatedAt:     parseTime(result.CreatedAt),
			LastUpdatedAt: parseTimePtr(result.LastUpdatedAt),
			ArchivedAt:    parseTimePtr(result.ArchivedAt),
			ID:            result.ID,
			DeviceToken:   decryptedDeviceToken,
			Platform:      result.Platform,
			BelongsToUser: result.BelongsToUser,
		}
		data = append(data, token)
		filteredCount = uint64(result.FilteredCount)
		totalCount = uint64(result.TotalCount)
	}

	x := filtering.NewQueryFilteredResult(
		data,
		filteredCount,
		totalCount,
		func(t *types.UserDeviceToken) string {
			return t.ID
		},
		filter,
	)

	return x, nil
}

// CreateUserDeviceToken creates or updates a user device token in the database (upsert).
// If the same user+token exists and is not archived, updates last_updated_at and returns the existing record.
func (q *Repository) CreateUserDeviceToken(ctx context.Context, input *types.UserDeviceTokenDatabaseCreationInput) (*types.UserDeviceToken, error) {
	return q.UpsertUserDeviceToken(ctx, input)
}

// UpsertUserDeviceToken creates or updates a user device token.
// If the same user+token exists and is not archived, updates last_updated_at and returns the existing record.
func (q *Repository) UpsertUserDeviceToken(ctx context.Context, input *types.UserDeviceTokenDatabaseCreationInput) (*types.UserDeviceToken, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, platformerrors.ErrNilInputProvided
	}
	tracing.AttachToSpan(span, notificationkeys.UserDeviceTokenIDKey, input.ID)
	logger := q.logger.WithValue(notificationkeys.UserDeviceTokenIDKey, input.ID)

	encryptedDeviceToken, err := q.userDeviceTokenEncDec.Encrypt(ctx, input.DeviceToken)
	if err != nil {
		return nil, observability.PrepareError(err, span, "encrypting user device token")
	}

	result, err := q.generatedQuerier.UpsertUserDeviceToken(ctx, q.writeDB, &generated.UpsertUserDeviceTokenParams{
		ID:            input.ID,
		BelongsToUser: input.BelongsToUser,
		DeviceToken:   encryptedDeviceToken,
		Platform:      input.Platform,
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "performing user device token upsert query")
	}

	decryptedDeviceToken, err := q.userDeviceTokenEncDec.Decrypt(ctx, result.DeviceToken)
	if err != nil {
		return nil, observability.PrepareError(err, span, "decrypting user device token")
	}

	x := &types.UserDeviceToken{
		CreatedAt:     parseTime(result.CreatedAt),
		LastUpdatedAt: parseTimePtr(result.LastUpdatedAt),
		ArchivedAt:    parseTimePtr(result.ArchivedAt),
		ID:            result.ID,
		DeviceToken:   decryptedDeviceToken,
		Platform:      result.Platform,
		BelongsToUser: result.BelongsToUser,
	}
	logger.Info("user device token upserted")
	return x, nil
}

// UpdateUserDeviceToken updates a particular user device token.
func (q *Repository) UpdateUserDeviceToken(ctx context.Context, updated *types.UserDeviceToken) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return platformerrors.ErrNilInputProvided
	}
	logger := q.logger.WithValue(notificationkeys.UserDeviceTokenIDKey, updated.ID)
	tracing.AttachToSpan(span, notificationkeys.UserDeviceTokenIDKey, updated.ID)

	tx, err := q.writeDB.BeginTx(ctx, nil)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "beginning transaction")
	}

	if _, err = q.generatedQuerier.UpdateUserDeviceToken(ctx, tx, &generated.UpdateUserDeviceTokenParams{
		Platform:      updated.Platform,
		ID:            updated.ID,
		BelongsToUser: updated.BelongsToUser,
	}); err != nil {
		q.RollbackTransaction(ctx, tx)
		return observability.PrepareAndLogError(err, logger, span, "updating user device token")
	}

	if _, err = q.auditLogEntryRepo.CreateAuditLogEntry(ctx, tx, &audit.AuditLogEntryDatabaseCreationInput{
		ID:            identifiers.New(),
		ResourceType:  resourceTypeUserDeviceTokens,
		RelevantID:    updated.ID,
		EventType:     audit.AuditLogEventTypeUpdated,
		BelongsToUser: updated.BelongsToUser,
	}); err != nil {
		q.RollbackTransaction(ctx, tx)
		return observability.PrepareError(err, span, "creating audit log entry")
	}

	if err = tx.Commit(); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "committing transaction")
	}

	logger.Info("user device token updated")
	return nil
}

// ArchiveUserDeviceToken archives a user device token.
func (q *Repository) ArchiveUserDeviceToken(ctx context.Context, userID, tokenID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if userID == "" {
		return platformerrors.ErrInvalidIDProvided
	}
	if tokenID == "" {
		return platformerrors.ErrInvalidIDProvided
	}
	logger := q.logger.WithValue(notificationkeys.UserDeviceTokenIDKey, tokenID)
	tracing.AttachToSpan(span, notificationkeys.UserDeviceTokenIDKey, tokenID)

	tx, err := q.writeDB.BeginTx(ctx, nil)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "beginning transaction")
	}

	rows, err := q.generatedQuerier.ArchiveUserDeviceToken(ctx, tx, &generated.ArchiveUserDeviceTokenParams{
		ID:            tokenID,
		BelongsToUser: userID,
	})
	if err != nil {
		q.RollbackTransaction(ctx, tx)
		return observability.PrepareAndLogError(err, logger, span, "archiving user device token")
	}
	if rows == 0 {
		q.RollbackTransaction(ctx, tx)
		return sql.ErrNoRows
	}

	if _, err = q.auditLogEntryRepo.CreateAuditLogEntry(ctx, tx, &audit.AuditLogEntryDatabaseCreationInput{
		ID:            identifiers.New(),
		ResourceType:  resourceTypeUserDeviceTokens,
		RelevantID:    tokenID,
		EventType:     audit.AuditLogEventTypeArchived,
		BelongsToUser: userID,
	}); err != nil {
		q.RollbackTransaction(ctx, tx)
		return observability.PrepareError(err, span, "creating audit log entry")
	}

	if err = tx.Commit(); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "committing transaction")
	}

	logger.Info("user device token archived")
	return nil
}
