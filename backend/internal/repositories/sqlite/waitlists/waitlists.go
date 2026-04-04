package waitlists

import (
	"context"
	"database/sql"
	"time"

	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/audit"
	identitykeys "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/identity/keys"
	types "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/waitlists"
	waitlistkeys "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/waitlists/keys"
	generated "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/repositories/sqlite/waitlists/generated"

	"github.com/verygoodsoftwarenotvirus/platform/v4/database/filtering"
	platformerrors "github.com/verygoodsoftwarenotvirus/platform/v4/errors"
	"github.com/verygoodsoftwarenotvirus/platform/v4/identifiers"
	"github.com/verygoodsoftwarenotvirus/platform/v4/observability"
	"github.com/verygoodsoftwarenotvirus/platform/v4/observability/tracing"
)

const (
	resourceTypeWaitlists       = "waitlists"
	resourceTypeWaitlistSignups = "waitlist_signups"
)

var (
	_ types.Repository = (*Repository)(nil)
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

func timeToString(t time.Time) string {
	return t.Format(time.RFC3339Nano)
}

func timePtrToStringPtr(t *time.Time) *string {
	if t == nil {
		return nil
	}
	s := timeToString(*t)
	return &s
}

func int64PtrFromUint8Ptr(v *uint8) any {
	if v == nil {
		return nil
	}
	x := int64(*v)
	return &x
}

func stringFromStringPtr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// WaitlistIsNotExpired checks if a waitlist exists and is not expired.
func (r *Repository) WaitlistIsNotExpired(ctx context.Context, waitlistID string) (bool, error) {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	logger := r.logger.Clone()

	if waitlistID == "" {
		return false, platformerrors.ErrInvalidIDProvided
	}
	logger = logger.WithValue(waitlistkeys.WaitlistIDKey, waitlistID)
	tracing.AttachToSpan(span, waitlistkeys.WaitlistIDKey, waitlistID)

	existsStr, err := r.generatedQuerier.CheckWaitlistExistence(ctx, r.readDB, waitlistID)
	if err != nil {
		return false, observability.PrepareAndLogError(err, logger, span, "checking waitlist existence")
	}

	if existsStr != "1" {
		return false, sql.ErrNoRows
	}

	resultStr, err := r.generatedQuerier.WaitlistIsNotExpired(ctx, r.readDB, waitlistID)
	if err != nil {
		return false, observability.PrepareAndLogError(err, logger, span, "checking waitlist expiration status")
	}

	return resultStr == "1", nil
}

// GetWaitlist fetches a waitlist from the database.
func (r *Repository) GetWaitlist(ctx context.Context, waitlistID string) (*types.Waitlist, error) {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	logger := r.logger.Clone()

	if waitlistID == "" {
		return nil, platformerrors.ErrInvalidIDProvided
	}
	logger = logger.WithValue(waitlistkeys.WaitlistIDKey, waitlistID)
	tracing.AttachToSpan(span, waitlistkeys.WaitlistIDKey, waitlistID)

	result, err := r.generatedQuerier.GetWaitlist(ctx, r.readDB, waitlistID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching waitlist")
	}

	waitlist := &types.Waitlist{
		CreatedAt:     parseTime(result.CreatedAt),
		LastUpdatedAt: parseTimePtr(result.LastUpdatedAt),
		ArchivedAt:    parseTimePtr(result.ArchivedAt),
		ID:            result.ID,
		Name:          result.Name,
		Description:   result.Description,
		ValidUntil:    parseTime(result.ValidUntil),
	}

	return waitlist, nil
}

// GetWaitlists fetches waitlists with filtering.
func (r *Repository) GetWaitlists(ctx context.Context, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[types.Waitlist], error) {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	logger := r.logger.Clone()

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	results, err := r.generatedQuerier.GetWaitlists(ctx, r.readDB, &generated.GetWaitlistsParams{
		CreatedAfter:  timePtrToStringPtr(filter.CreatedAfter),
		CreatedBefore: timePtrToStringPtr(filter.CreatedBefore),
		UpdatedBefore: timePtrToStringPtr(filter.UpdatedBefore),
		UpdatedAfter:  timePtrToStringPtr(filter.UpdatedAfter),
		Cursor:        filter.Cursor,
		ResultLimit:   int64PtrFromUint8Ptr(filter.MaxResponseSize),
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching waitlists from database")
	}

	var (
		data                      []*types.Waitlist
		filteredCount, totalCount uint64
	)
	for _, result := range results {
		data = append(data, &types.Waitlist{
			CreatedAt:     parseTime(result.CreatedAt),
			LastUpdatedAt: parseTimePtr(result.LastUpdatedAt),
			ArchivedAt:    parseTimePtr(result.ArchivedAt),
			ID:            result.ID,
			Name:          result.Name,
			Description:   result.Description,
			ValidUntil:    parseTime(result.ValidUntil),
		})

		filteredCount = uint64(result.FilteredCount)
		totalCount = uint64(result.TotalCount)
	}

	x := filtering.NewQueryFilteredResult(
		data,
		filteredCount,
		totalCount,
		func(t *types.Waitlist) string {
			return t.ID
		},
		filter,
	)

	return x, nil
}

// GetActiveWaitlists fetches non-expired waitlists with filtering.
func (r *Repository) GetActiveWaitlists(ctx context.Context, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[types.Waitlist], error) {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	logger := r.logger.Clone()

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	results, err := r.generatedQuerier.GetActiveWaitlists(ctx, r.readDB, &generated.GetActiveWaitlistsParams{
		CreatedAfter:  timePtrToStringPtr(filter.CreatedAfter),
		CreatedBefore: timePtrToStringPtr(filter.CreatedBefore),
		UpdatedBefore: timePtrToStringPtr(filter.UpdatedBefore),
		UpdatedAfter:  timePtrToStringPtr(filter.UpdatedAfter),
		Cursor:        filter.Cursor,
		ResultLimit:   int64PtrFromUint8Ptr(filter.MaxResponseSize),
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching active waitlists from database")
	}

	var (
		data                      []*types.Waitlist
		filteredCount, totalCount uint64
	)
	for _, result := range results {
		data = append(data, &types.Waitlist{
			CreatedAt:     parseTime(result.CreatedAt),
			LastUpdatedAt: parseTimePtr(result.LastUpdatedAt),
			ArchivedAt:    parseTimePtr(result.ArchivedAt),
			ID:            result.ID,
			Name:          result.Name,
			Description:   result.Description,
			ValidUntil:    parseTime(result.ValidUntil),
		})

		filteredCount = uint64(result.FilteredCount)
		totalCount = uint64(result.TotalCount)
	}

	x := filtering.NewQueryFilteredResult(
		data,
		filteredCount,
		totalCount,
		func(t *types.Waitlist) string {
			return t.ID
		},
		filter,
	)

	return x, nil
}

// CreateWaitlist creates a waitlist in the database.
func (r *Repository) CreateWaitlist(ctx context.Context, input *types.WaitlistDatabaseCreationInput) (*types.Waitlist, error) {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, platformerrors.ErrNilInputProvided
	}

	logger := r.logger.WithValue(waitlistkeys.WaitlistIDKey, input.ID)
	tracing.AttachToSpan(span, waitlistkeys.WaitlistIDKey, input.ID)

	if err := r.generatedQuerier.CreateWaitlist(ctx, r.writeDB, &generated.CreateWaitlistParams{
		ID:          input.ID,
		Name:        input.Name,
		Description: input.Description,
		ValidUntil:  timeToString(input.ValidUntil),
	}); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "performing waitlist creation query")
	}

	if _, err := r.auditLogEntryRepo.CreateAuditLogEntry(ctx, r.writeDB, &audit.AuditLogEntryDatabaseCreationInput{
		ID:           identifiers.New(),
		ResourceType: resourceTypeWaitlists,
		RelevantID:   input.ID,
		EventType:    audit.AuditLogEventTypeCreated,
	}); err != nil {
		return nil, observability.PrepareError(err, span, "creating audit log entry")
	}

	x := &types.Waitlist{
		ID:          input.ID,
		CreatedAt:   r.CurrentTime(),
		Name:        input.Name,
		Description: input.Description,
		ValidUntil:  input.ValidUntil,
	}

	logger.Info("waitlist created")
	return x, nil
}

// UpdateWaitlist updates a waitlist.
func (r *Repository) UpdateWaitlist(ctx context.Context, updated *types.Waitlist) error {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return platformerrors.ErrNilInputProvided
	}
	logger := r.logger.WithValue(waitlistkeys.WaitlistIDKey, updated.ID)
	tracing.AttachToSpan(span, waitlistkeys.WaitlistIDKey, updated.ID)

	if _, err := r.generatedQuerier.UpdateWaitlist(ctx, r.writeDB, &generated.UpdateWaitlistParams{
		Name:        updated.Name,
		Description: updated.Description,
		ValidUntil:  timeToString(updated.ValidUntil),
		ID:          updated.ID,
	}); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating waitlist")
	}

	if _, err := r.auditLogEntryRepo.CreateAuditLogEntry(ctx, r.writeDB, &audit.AuditLogEntryDatabaseCreationInput{
		ID:           identifiers.New(),
		ResourceType: resourceTypeWaitlists,
		RelevantID:   updated.ID,
		EventType:    audit.AuditLogEventTypeUpdated,
	}); err != nil {
		return observability.PrepareError(err, span, "creating audit log entry")
	}

	return nil
}

// ArchiveWaitlist archives a waitlist.
func (r *Repository) ArchiveWaitlist(ctx context.Context, waitlistID string) error {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	if waitlistID == "" {
		return platformerrors.ErrInvalidIDProvided
	}
	logger := r.logger.WithValue(waitlistkeys.WaitlistIDKey, waitlistID)
	tracing.AttachToSpan(span, waitlistkeys.WaitlistIDKey, waitlistID)

	recordsChanged, err := r.generatedQuerier.ArchiveWaitlist(ctx, r.writeDB, waitlistID)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving waitlist")
	}

	if recordsChanged == 0 {
		return sql.ErrNoRows
	}

	if _, err = r.auditLogEntryRepo.CreateAuditLogEntry(ctx, r.writeDB, &audit.AuditLogEntryDatabaseCreationInput{
		ID:           identifiers.New(),
		ResourceType: resourceTypeWaitlists,
		RelevantID:   waitlistID,
		EventType:    audit.AuditLogEventTypeArchived,
	}); err != nil {
		return observability.PrepareError(err, span, "creating audit log entry")
	}

	logger.Info("waitlist archived")
	return nil
}

// GetWaitlistSignup fetches a waitlist signup from the database.
func (r *Repository) GetWaitlistSignup(ctx context.Context, waitlistSignupID, waitlistID string) (*types.WaitlistSignup, error) {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	logger := r.logger.Clone()

	if waitlistSignupID == "" {
		return nil, platformerrors.ErrInvalidIDProvided
	}
	logger = logger.WithValue(waitlistkeys.WaitlistSignupIDKey, waitlistSignupID)
	tracing.AttachToSpan(span, waitlistkeys.WaitlistSignupIDKey, waitlistSignupID)

	if waitlistID == "" {
		return nil, platformerrors.ErrInvalidIDProvided
	}
	logger = logger.WithValue(waitlistkeys.WaitlistIDKey, waitlistID)
	tracing.AttachToSpan(span, waitlistkeys.WaitlistIDKey, waitlistID)

	result, err := r.generatedQuerier.GetWaitlistSignup(ctx, r.readDB, &generated.GetWaitlistSignupParams{
		ID:                waitlistSignupID,
		BelongsToWaitlist: new(waitlistID),
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching waitlist signup")
	}

	waitlistSignup := &types.WaitlistSignup{
		CreatedAt:         parseTime(result.CreatedAt),
		LastUpdatedAt:     parseTimePtr(result.LastUpdatedAt),
		ArchivedAt:        parseTimePtr(result.ArchivedAt),
		ID:                result.ID,
		Notes:             result.Notes,
		BelongsToWaitlist: stringFromStringPtr(result.BelongsToWaitlist),
		BelongsToUser:     stringFromStringPtr(result.BelongsToUser),
		BelongsToAccount:  stringFromStringPtr(result.BelongsToAccount),
	}

	return waitlistSignup, nil
}

// GetWaitlistSignupsForWaitlist fetches waitlist signups for a waitlist with filtering.
func (r *Repository) GetWaitlistSignupsForWaitlist(ctx context.Context, waitlistID string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[types.WaitlistSignup], error) {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	logger := r.logger.Clone()

	if waitlistID == "" {
		return nil, platformerrors.ErrInvalidIDProvided
	}
	logger = logger.WithValue(waitlistkeys.WaitlistIDKey, waitlistID)
	tracing.AttachToSpan(span, waitlistkeys.WaitlistIDKey, waitlistID)

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	results, err := r.generatedQuerier.GetWaitlistSignupsForWaitlist(ctx, r.readDB, &generated.GetWaitlistSignupsForWaitlistParams{
		BelongsToWaitlist: new(waitlistID),
		CreatedAfter:      timePtrToStringPtr(filter.CreatedAfter),
		CreatedBefore:     timePtrToStringPtr(filter.CreatedBefore),
		UpdatedBefore:     timePtrToStringPtr(filter.UpdatedBefore),
		UpdatedAfter:      timePtrToStringPtr(filter.UpdatedAfter),
		Cursor:            filter.Cursor,
		ResultLimit:       int64PtrFromUint8Ptr(filter.MaxResponseSize),
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching waitlist signups from database")
	}

	var (
		data                      []*types.WaitlistSignup
		filteredCount, totalCount uint64
	)
	for _, result := range results {
		data = append(data, &types.WaitlistSignup{
			CreatedAt:         parseTime(result.CreatedAt),
			LastUpdatedAt:     parseTimePtr(result.LastUpdatedAt),
			ArchivedAt:        parseTimePtr(result.ArchivedAt),
			ID:                result.ID,
			Notes:             result.Notes,
			BelongsToWaitlist: stringFromStringPtr(result.BelongsToWaitlist),
			BelongsToUser:     stringFromStringPtr(result.BelongsToUser),
			BelongsToAccount:  stringFromStringPtr(result.BelongsToAccount),
		})

		filteredCount = uint64(result.FilteredCount)
		totalCount = uint64(result.TotalCount)
	}

	x := filtering.NewQueryFilteredResult(
		data,
		filteredCount,
		totalCount,
		func(t *types.WaitlistSignup) string {
			return t.ID
		},
		filter,
	)

	return x, nil
}

// GetWaitlistSignupsForUser fetches waitlist signups for a user with filtering.
func (r *Repository) GetWaitlistSignupsForUser(ctx context.Context, userID string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[types.WaitlistSignup], error) {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	logger := r.logger.Clone()

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

	results, err := r.generatedQuerier.GetWaitlistSignupsForUser(ctx, r.readDB, &generated.GetWaitlistSignupsForUserParams{
		BelongsToUser: new(userID),
		CreatedAfter:  timePtrToStringPtr(filter.CreatedAfter),
		CreatedBefore: timePtrToStringPtr(filter.CreatedBefore),
		UpdatedBefore: timePtrToStringPtr(filter.UpdatedBefore),
		UpdatedAfter:  timePtrToStringPtr(filter.UpdatedAfter),
		Cursor:        filter.Cursor,
		ResultLimit:   int64PtrFromUint8Ptr(filter.MaxResponseSize),
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching waitlist signups for user from database")
	}

	var (
		data                      []*types.WaitlistSignup
		filteredCount, totalCount uint64
	)
	for _, result := range results {
		data = append(data, &types.WaitlistSignup{
			CreatedAt:         parseTime(result.CreatedAt),
			LastUpdatedAt:     parseTimePtr(result.LastUpdatedAt),
			ArchivedAt:        parseTimePtr(result.ArchivedAt),
			ID:                result.ID,
			Notes:             result.Notes,
			BelongsToWaitlist: stringFromStringPtr(result.BelongsToWaitlist),
			BelongsToUser:     stringFromStringPtr(result.BelongsToUser),
			BelongsToAccount:  stringFromStringPtr(result.BelongsToAccount),
		})

		filteredCount = uint64(result.FilteredCount)
		totalCount = uint64(result.TotalCount)
	}

	return filtering.NewQueryFilteredResult(
		data,
		filteredCount,
		totalCount,
		func(t *types.WaitlistSignup) string {
			return t.ID
		},
		filter,
	), nil
}

// CreateWaitlistSignup creates a waitlist signup in the database.
func (r *Repository) CreateWaitlistSignup(ctx context.Context, input *types.WaitlistSignupDatabaseCreationInput) (*types.WaitlistSignup, error) {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, platformerrors.ErrNilInputProvided
	}

	logger := r.logger.WithValue(waitlistkeys.WaitlistSignupIDKey, input.ID)
	tracing.AttachToSpan(span, waitlistkeys.WaitlistSignupIDKey, input.ID)

	if err := r.generatedQuerier.CreateWaitlistSignup(ctx, r.writeDB, &generated.CreateWaitlistSignupParams{
		ID:                input.ID,
		Notes:             input.Notes,
		BelongsToWaitlist: new(input.BelongsToWaitlist),
		BelongsToUser:     new(input.BelongsToUser),
		BelongsToAccount:  new(input.BelongsToAccount),
	}); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "performing waitlist signup creation query")
	}

	if _, err := r.auditLogEntryRepo.CreateAuditLogEntry(ctx, r.writeDB, &audit.AuditLogEntryDatabaseCreationInput{
		BelongsToAccount: &input.BelongsToAccount,
		ID:               identifiers.New(),
		ResourceType:     resourceTypeWaitlistSignups,
		RelevantID:       input.ID,
		EventType:        audit.AuditLogEventTypeCreated,
	}); err != nil {
		return nil, observability.PrepareError(err, span, "creating audit log entry")
	}

	x := &types.WaitlistSignup{
		ID:                input.ID,
		CreatedAt:         r.CurrentTime(),
		Notes:             input.Notes,
		BelongsToWaitlist: input.BelongsToWaitlist,
		BelongsToUser:     input.BelongsToUser,
		BelongsToAccount:  input.BelongsToAccount,
	}

	logger.Info("waitlist signup created")
	return x, nil
}

// UpdateWaitlistSignup updates a waitlist signup.
func (r *Repository) UpdateWaitlistSignup(ctx context.Context, updated *types.WaitlistSignup) error {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return platformerrors.ErrNilInputProvided
	}
	logger := r.logger.WithValue(waitlistkeys.WaitlistSignupIDKey, updated.ID)
	tracing.AttachToSpan(span, waitlistkeys.WaitlistSignupIDKey, updated.ID)

	if _, err := r.generatedQuerier.UpdateWaitlistSignup(ctx, r.writeDB, &generated.UpdateWaitlistSignupParams{
		Notes: updated.Notes,
		ID:    updated.ID,
	}); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating waitlist signup")
	}

	if _, err := r.auditLogEntryRepo.CreateAuditLogEntry(ctx, r.writeDB, &audit.AuditLogEntryDatabaseCreationInput{
		BelongsToAccount: &updated.BelongsToAccount,
		ID:               identifiers.New(),
		ResourceType:     resourceTypeWaitlistSignups,
		RelevantID:       updated.ID,
		EventType:        audit.AuditLogEventTypeUpdated,
	}); err != nil {
		return observability.PrepareError(err, span, "creating audit log entry")
	}

	return nil
}

// ArchiveWaitlistSignup archives a waitlist signup.
func (r *Repository) ArchiveWaitlistSignup(ctx context.Context, waitlistSignupID string) error {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	if waitlistSignupID == "" {
		return platformerrors.ErrInvalidIDProvided
	}
	logger := r.logger.WithValue(waitlistkeys.WaitlistSignupIDKey, waitlistSignupID)
	tracing.AttachToSpan(span, waitlistkeys.WaitlistSignupIDKey, waitlistSignupID)

	recordsChanged, err := r.generatedQuerier.ArchiveWaitlistSignup(ctx, r.writeDB, waitlistSignupID)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving waitlist signup")
	}

	if recordsChanged == 0 {
		return sql.ErrNoRows
	}

	// ArchiveWaitlistSignup does not have account ID in signature; create audit entry without it
	if _, err = r.auditLogEntryRepo.CreateAuditLogEntry(ctx, r.writeDB, &audit.AuditLogEntryDatabaseCreationInput{
		ID:           identifiers.New(),
		ResourceType: resourceTypeWaitlistSignups,
		RelevantID:   waitlistSignupID,
		EventType:    audit.AuditLogEventTypeArchived,
	}); err != nil {
		return observability.PrepareError(err, span, "creating audit log entry")
	}

	logger.Info("waitlist signup archived")
	return nil
}
