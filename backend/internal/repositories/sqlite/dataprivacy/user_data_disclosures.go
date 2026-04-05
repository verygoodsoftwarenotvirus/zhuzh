package dataprivacy

import (
	"context"
	"time"

	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/audit"
	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/dataprivacy"
	identitykeys "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/identity/keys"
	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/repositories/sqlite/dataprivacy/generated"

	"github.com/verygoodsoftwarenotvirus/platform/v4/database/filtering"
	platformerrors "github.com/verygoodsoftwarenotvirus/platform/v4/errors"
	"github.com/verygoodsoftwarenotvirus/platform/v4/identifiers"
	"github.com/verygoodsoftwarenotvirus/platform/v4/observability"
	"github.com/verygoodsoftwarenotvirus/platform/v4/observability/tracing"
)

const (
	disclosureIDKey                = "disclosure_id"
	resourceTypeUserDataDisclosure = "user_data_disclosure"
)

func parseTime(s string) time.Time {
	t, err := time.Parse(time.RFC3339Nano, s)
	if err != nil {
		return time.Time{}
	}
	return t
}

func parseTimePtr(s *string) *time.Time {
	if s == nil {
		return nil
	}
	t := parseTime(*s)
	return &t
}

func stringFromStringPointer(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// CreateUserDataDisclosure creates a new user data disclosure record.
func (r *repository) CreateUserDataDisclosure(ctx context.Context, input *dataprivacy.UserDataDisclosureCreationInput) (*dataprivacy.UserDataDisclosure, error) {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, platformerrors.ErrNilInputProvided
	}

	if input.ID == "" {
		return nil, platformerrors.ErrInvalidIDProvided
	}

	tracing.AttachToSpan(span, disclosureIDKey, input.ID)
	tracing.AttachToSpan(span, identitykeys.UserIDKey, input.BelongsToUser)

	logger := r.logger.WithValue(disclosureIDKey, input.ID).WithValue(identitykeys.UserIDKey, input.BelongsToUser)
	logger.Info("creating user data disclosure")

	if err := r.generatedQuerier.CreateUserDataDisclosure(ctx, r.writeDB, &generated.CreateUserDataDisclosureParams{
		ID:            input.ID,
		BelongsToUser: input.BelongsToUser,
		ExpiresAt:     input.ExpiresAt.Format(time.RFC3339Nano),
	}); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating user data disclosure")
	}

	if _, err := r.auditLogRepo.CreateAuditLogEntry(ctx, r.writeDB, &audit.AuditLogEntryDatabaseCreationInput{
		ID:            identifiers.New(),
		ResourceType:  resourceTypeUserDataDisclosure,
		RelevantID:    input.ID,
		EventType:     audit.AuditLogEventTypeCreated,
		BelongsToUser: input.BelongsToUser,
	}); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating audit log entry")
	}

	disclosure, err := r.GetUserDataDisclosure(ctx, input.ID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching created disclosure")
	}

	return disclosure, nil
}

// GetUserDataDisclosure fetches a user data disclosure by ID.
func (r *repository) GetUserDataDisclosure(ctx context.Context, disclosureID string) (*dataprivacy.UserDataDisclosure, error) {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	if disclosureID == "" {
		return nil, platformerrors.ErrInvalidIDProvided
	}

	tracing.AttachToSpan(span, disclosureIDKey, disclosureID)
	logger := r.logger.WithValue(disclosureIDKey, disclosureID)

	result, err := r.generatedQuerier.GetUserDataDisclosure(ctx, r.readDB, disclosureID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching user data disclosure")
	}

	disclosure := &dataprivacy.UserDataDisclosure{
		ID:            result.ID,
		BelongsToUser: result.BelongsToUser,
		Status:        dataprivacy.UserDataDisclosureStatus(result.Status),
		ExpiresAt:     parseTime(result.ExpiresAt),
		CreatedAt:     parseTime(result.CreatedAt),
		LastUpdatedAt: parseTimePtr(result.LastUpdatedAt),
		CompletedAt:   parseTimePtr(result.CompletedAt),
		ArchivedAt:    parseTimePtr(result.ArchivedAt),
		ReportID:      stringFromStringPointer(result.ReportID),
	}

	return disclosure, nil
}

// GetUserDataDisclosuresForUser fetches user data disclosures for a user.
func (r *repository) GetUserDataDisclosuresForUser(ctx context.Context, userID string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[dataprivacy.UserDataDisclosure], error) {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	if userID == "" {
		return nil, platformerrors.ErrInvalidIDProvided
	}

	tracing.AttachToSpan(span, identitykeys.UserIDKey, userID)
	logger := r.logger.WithValue(identitykeys.UserIDKey, userID)

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}

	args := &generated.GetUserDataDisclosuresForUserParams{
		UserID: userID,
	}

	if filter.CreatedAfter != nil {
		s := filter.CreatedAfter.Format(time.RFC3339Nano)
		args.CreatedAfter = &s
	}
	if filter.CreatedBefore != nil {
		s := filter.CreatedBefore.Format(time.RFC3339Nano)
		args.CreatedBefore = &s
	}
	if filter.Cursor != nil {
		args.Cursor = filter.Cursor
	}
	if filter.MaxResponseSize != nil {
		args.ResultLimit = int64(*filter.MaxResponseSize)
	}

	results, err := r.generatedQuerier.GetUserDataDisclosuresForUser(ctx, r.readDB, args)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching user data disclosures")
	}

	disclosures := make([]*dataprivacy.UserDataDisclosure, 0, len(results))
	var filteredCount, totalCount int64

	for _, result := range results {
		disclosure := &dataprivacy.UserDataDisclosure{
			ID:            result.ID,
			BelongsToUser: result.BelongsToUser,
			Status:        dataprivacy.UserDataDisclosureStatus(result.Status),
			ExpiresAt:     parseTime(result.ExpiresAt),
			CreatedAt:     parseTime(result.CreatedAt),
			LastUpdatedAt: parseTimePtr(result.LastUpdatedAt),
			CompletedAt:   parseTimePtr(result.CompletedAt),
			ArchivedAt:    parseTimePtr(result.ArchivedAt),
			ReportID:      stringFromStringPointer(result.ReportID),
		}

		disclosures = append(disclosures, disclosure)
		filteredCount = result.FilteredCount
		totalCount = result.TotalCount
	}

	return &filtering.QueryFilteredResult[dataprivacy.UserDataDisclosure]{
		Data: disclosures,
		Pagination: filtering.Pagination{
			FilteredCount: uint64(filteredCount),
			TotalCount:    uint64(totalCount),
		},
	}, nil
}

// MarkUserDataDisclosureCompleted marks a disclosure as completed.
func (r *repository) MarkUserDataDisclosureCompleted(ctx context.Context, disclosureID, reportID string) error {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	if disclosureID == "" || reportID == "" {
		return platformerrors.ErrInvalidIDProvided
	}

	tracing.AttachToSpan(span, disclosureIDKey, disclosureID)
	logger := r.logger.WithValue(disclosureIDKey, disclosureID)

	if err := r.generatedQuerier.MarkUserDataDisclosureCompleted(ctx, r.writeDB, &generated.MarkUserDataDisclosureCompletedParams{
		ID:       disclosureID,
		ReportID: &reportID,
	}); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "marking disclosure completed")
	}

	return nil
}

// MarkUserDataDisclosureFailed marks a disclosure as failed.
func (r *repository) MarkUserDataDisclosureFailed(ctx context.Context, disclosureID string) error {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	if disclosureID == "" {
		return platformerrors.ErrInvalidIDProvided
	}

	tracing.AttachToSpan(span, disclosureIDKey, disclosureID)
	logger := r.logger.WithValue(disclosureIDKey, disclosureID)

	if err := r.generatedQuerier.MarkUserDataDisclosureFailed(ctx, r.writeDB, disclosureID); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "marking disclosure failed")
	}

	return nil
}

// ArchiveUserDataDisclosure archives a disclosure.
func (r *repository) ArchiveUserDataDisclosure(ctx context.Context, disclosureID string) error {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	if disclosureID == "" {
		return platformerrors.ErrInvalidIDProvided
	}

	tracing.AttachToSpan(span, disclosureIDKey, disclosureID)
	logger := r.logger.WithValue(disclosureIDKey, disclosureID)

	disclosure, err := r.GetUserDataDisclosure(ctx, disclosureID)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "fetching disclosure for archive")
	}

	tx, err := r.writeDB.BeginTx(ctx, nil)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "beginning transaction")
	}

	if err = r.generatedQuerier.ArchiveUserDataDisclosure(ctx, tx, disclosureID); err != nil {
		r.RollbackTransaction(ctx, tx)
		return observability.PrepareAndLogError(err, logger, span, "archiving disclosure")
	}

	if _, err = r.auditLogRepo.CreateAuditLogEntry(ctx, tx, &audit.AuditLogEntryDatabaseCreationInput{
		ID:            identifiers.New(),
		ResourceType:  resourceTypeUserDataDisclosure,
		RelevantID:    disclosureID,
		EventType:     audit.AuditLogEventTypeArchived,
		BelongsToUser: disclosure.BelongsToUser,
	}); err != nil {
		r.RollbackTransaction(ctx, tx)
		return observability.PrepareAndLogError(err, logger, span, "creating audit log entry")
	}

	if err = tx.Commit(); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "committing transaction")
	}

	return nil
}
