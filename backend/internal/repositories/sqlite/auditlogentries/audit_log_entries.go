package auditlogentries

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/audit"
	auditkeys "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/audit/keys"
	identitykeys "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/identity/keys"
	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/repositories/sqlite/auditlogentries/generated"

	"github.com/verygoodsoftwarenotvirus/platform/v4/database"
	"github.com/verygoodsoftwarenotvirus/platform/v4/database/filtering"
	platformerrors "github.com/verygoodsoftwarenotvirus/platform/v4/errors"
	"github.com/verygoodsoftwarenotvirus/platform/v4/observability"
	"github.com/verygoodsoftwarenotvirus/platform/v4/observability/tracing"
)

var (
	_ audit.AuditLogEntryDataManager = (*repository)(nil)
)

func parseTime(s string) time.Time {
	t, err := time.Parse(time.RFC3339Nano, s)
	if err != nil {
		return time.Time{}
	}
	return t
}

func stringFromStringPointer(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func stringPointerFromStringPointer(s *string) *string {
	if s == nil {
		return nil
	}
	v := *s
	return &v
}

func stringPointerFromString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func timeStringFromTimePointer(t *time.Time) *string {
	if t == nil {
		return nil
	}
	s := t.Format(time.RFC3339Nano)
	return &s
}

func uint8PtrToInterfacePtr(v *uint8) any {
	if v == nil {
		return nil
	}
	return int64(*v)
}

// GetAuditLogEntry fetches an audit log entry from the database.
func (q *repository) GetAuditLogEntry(ctx context.Context, auditLogEntryID string) (*audit.AuditLogEntry, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if auditLogEntryID == "" {
		return nil, platformerrors.ErrInvalidIDProvided
	}
	logger = logger.WithValue(auditkeys.AuditLogEntryIDKey, auditLogEntryID)
	tracing.AttachToSpan(span, auditkeys.AuditLogEntryIDKey, auditLogEntryID)

	result, err := q.generatedQuerier.GetAuditLogEntry(ctx, q.readDB, auditLogEntryID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching audit log entry")
	}

	auditLogEntry := &audit.AuditLogEntry{
		CreatedAt:        parseTime(result.CreatedAt),
		BelongsToAccount: stringPointerFromStringPointer(result.BelongsToAccount),
		ID:               result.ID,
		ResourceType:     result.ResourceType,
		RelevantID:       result.RelevantID,
		EventType:        result.EventType,
		BelongsToUser:    stringFromStringPointer(result.BelongsToUser),
	}

	if err = json.Unmarshal([]byte(result.Changes), &auditLogEntry.Changes); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "parsing audit log entry JSON data")
	}

	return auditLogEntry, nil
}

// GetAuditLogEntriesForUser fetches a list of audit log entries from the database that meet a particular filter.
func (q *repository) GetAuditLogEntriesForUser(ctx context.Context, userID string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[audit.AuditLogEntry], error) {
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

	results, err := q.generatedQuerier.GetAuditLogEntriesForUser(ctx, q.readDB, &generated.GetAuditLogEntriesForUserParams{
		BelongsToUser: stringPointerFromString(userID),
		CreatedBefore: timeStringFromTimePointer(filter.CreatedBefore),
		CreatedAfter:  timeStringFromTimePointer(filter.CreatedAfter),
		Cursor:        filter.Cursor,
		ResultLimit:   uint8PtrToInterfacePtr(filter.MaxResponseSize),
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching audit log entries from database")
	}

	var (
		data                      []*audit.AuditLogEntry
		filteredCount, totalCount uint64
	)
	for _, result := range results {
		auditLogEntry := &audit.AuditLogEntry{
			CreatedAt:        parseTime(result.CreatedAt),
			BelongsToAccount: stringPointerFromStringPointer(result.BelongsToAccount),
			ID:               result.ID,
			ResourceType:     result.ResourceType,
			RelevantID:       result.RelevantID,
			EventType:        result.EventType,
			BelongsToUser:    stringFromStringPointer(result.BelongsToUser),
		}

		if err = json.Unmarshal([]byte(result.Changes), &auditLogEntry.Changes); err != nil {
			return nil, observability.PrepareAndLogError(err, logger, span, "parsing audit log entry JSON data")
		}

		data = append(data, auditLogEntry)
		filteredCount = uint64(result.FilteredCount)
		totalCount = uint64(result.TotalCount)
	}

	x := filtering.NewQueryFilteredResult(
		data,
		filteredCount,
		totalCount,
		func(t *audit.AuditLogEntry) string {
			return t.ID
		},
		filter,
	)

	return x, nil
}

// GetAuditLogEntriesForUserAndResourceTypes fetches a list of audit log entries from the database that meet a particular filter.
func (q *repository) GetAuditLogEntriesForUserAndResourceTypes(ctx context.Context, userID string, resourceTypes []string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[audit.AuditLogEntry], error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if userID == "" {
		return nil, platformerrors.ErrInvalidIDProvided
	}
	logger = logger.WithValue(identitykeys.UserIDKey, userID)
	tracing.AttachToSpan(span, identitykeys.UserIDKey, userID)

	if len(resourceTypes) == 0 {
		return nil, platformerrors.ErrEmptyInputProvided
	}
	logger = logger.WithValue(auditkeys.AuditLogEntryResourceTypesKey, resourceTypes)
	tracing.AttachToSpan(span, auditkeys.AuditLogEntryResourceTypesKey, resourceTypes)

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	results, err := q.generatedQuerier.GetAuditLogEntriesForUserAndResourceType(ctx, q.readDB, &generated.GetAuditLogEntriesForUserAndResourceTypeParams{
		BelongsToUser: stringPointerFromString(userID),
		Resources:     resourceTypes,
		CreatedBefore: timeStringFromTimePointer(filter.CreatedBefore),
		CreatedAfter:  timeStringFromTimePointer(filter.CreatedAfter),
		Cursor:        filter.Cursor,
		ResultLimit:   uint8PtrToInterfacePtr(filter.MaxResponseSize),
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching audit log entries from database")
	}

	var (
		data                      []*audit.AuditLogEntry
		filteredCount, totalCount uint64
	)
	for _, result := range results {
		auditLogEntry := &audit.AuditLogEntry{
			CreatedAt:        parseTime(result.CreatedAt),
			BelongsToAccount: stringPointerFromStringPointer(result.BelongsToAccount),
			ID:               result.ID,
			ResourceType:     result.ResourceType,
			RelevantID:       result.RelevantID,
			EventType:        result.EventType,
			BelongsToUser:    stringFromStringPointer(result.BelongsToUser),
		}

		if err = json.Unmarshal([]byte(result.Changes), &auditLogEntry.Changes); err != nil {
			return nil, observability.PrepareAndLogError(err, logger, span, "parsing audit log entry JSON data")
		}

		data = append(data, auditLogEntry)
		filteredCount = uint64(result.FilteredCount)
		totalCount = uint64(result.TotalCount)
	}

	x := filtering.NewQueryFilteredResult(
		data,
		filteredCount,
		totalCount,
		func(t *audit.AuditLogEntry) string {
			return t.ID
		},
		filter,
	)

	return x, nil
}

// GetAuditLogEntriesForAccount fetches a list of audit log entries from the database that meet a particular filter.
func (q *repository) GetAuditLogEntriesForAccount(ctx context.Context, accountID string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[audit.AuditLogEntry], error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if accountID == "" {
		return nil, platformerrors.ErrInvalidIDProvided
	}
	logger = logger.WithValue(identitykeys.AccountIDKey, accountID)
	tracing.AttachToSpan(span, identitykeys.AccountIDKey, accountID)

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	results, err := q.generatedQuerier.GetAuditLogEntriesForAccount(ctx, q.readDB, &generated.GetAuditLogEntriesForAccountParams{
		BelongsToAccount: stringPointerFromString(accountID),
		CreatedBefore:    timeStringFromTimePointer(filter.CreatedBefore),
		CreatedAfter:     timeStringFromTimePointer(filter.CreatedAfter),
		Cursor:           filter.Cursor,
		ResultLimit:      uint8PtrToInterfacePtr(filter.MaxResponseSize),
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching audit log entries from database")
	}

	var (
		data                      []*audit.AuditLogEntry
		filteredCount, totalCount uint64
	)
	for _, result := range results {
		auditLogEntry := &audit.AuditLogEntry{
			CreatedAt:        parseTime(result.CreatedAt),
			BelongsToAccount: stringPointerFromStringPointer(result.BelongsToAccount),
			ID:               result.ID,
			ResourceType:     result.ResourceType,
			RelevantID:       result.RelevantID,
			EventType:        result.EventType,
			BelongsToUser:    stringFromStringPointer(result.BelongsToUser),
		}

		if err = json.Unmarshal([]byte(result.Changes), &auditLogEntry.Changes); err != nil {
			return nil, observability.PrepareAndLogError(err, logger, span, "parsing audit log entry JSON data")
		}

		data = append(data, auditLogEntry)
		filteredCount = uint64(result.FilteredCount)
		totalCount = uint64(result.TotalCount)
	}

	x := filtering.NewQueryFilteredResult(
		data,
		filteredCount,
		totalCount,
		func(t *audit.AuditLogEntry) string {
			return t.ID
		},
		filter,
	)

	return x, nil
}

// GetAuditLogEntriesForAccountAndResourceTypes fetches a list of audit log entries from the database that meet a particular filter.
func (q *repository) GetAuditLogEntriesForAccountAndResourceTypes(ctx context.Context, accountID string, resourceTypes []string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[audit.AuditLogEntry], error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if accountID == "" {
		return nil, platformerrors.ErrInvalidIDProvided
	}
	logger = logger.WithValue(identitykeys.AccountIDKey, accountID)
	tracing.AttachToSpan(span, identitykeys.AccountIDKey, accountID)

	if len(resourceTypes) == 0 {
		return nil, platformerrors.ErrEmptyInputProvided
	}
	logger = logger.WithValue(auditkeys.AuditLogEntryResourceTypesKey, resourceTypes)
	tracing.AttachToSpan(span, auditkeys.AuditLogEntryResourceTypesKey, resourceTypes)

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	results, err := q.generatedQuerier.GetAuditLogEntriesForAccountAndResourceType(ctx, q.readDB, &generated.GetAuditLogEntriesForAccountAndResourceTypeParams{
		BelongsToAccount: stringPointerFromString(accountID),
		Resources:        resourceTypes,
		CreatedBefore:    timeStringFromTimePointer(filter.CreatedBefore),
		CreatedAfter:     timeStringFromTimePointer(filter.CreatedAfter),
		Cursor:           filter.Cursor,
		ResultLimit:      uint8PtrToInterfacePtr(filter.MaxResponseSize),
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching audit log entries from database")
	}

	var (
		data                      []*audit.AuditLogEntry
		filteredCount, totalCount uint64
	)
	for _, result := range results {
		auditLogEntry := &audit.AuditLogEntry{
			CreatedAt:        parseTime(result.CreatedAt),
			BelongsToAccount: stringPointerFromStringPointer(result.BelongsToAccount),
			ID:               result.ID,
			ResourceType:     result.ResourceType,
			RelevantID:       result.RelevantID,
			EventType:        result.EventType,
			BelongsToUser:    stringFromStringPointer(result.BelongsToUser),
		}

		if err = json.Unmarshal([]byte(result.Changes), &auditLogEntry.Changes); err != nil {
			return nil, observability.PrepareAndLogError(err, logger, span, "parsing audit log entry JSON data")
		}

		data = append(data, auditLogEntry)
		filteredCount = uint64(result.FilteredCount)
		totalCount = uint64(result.TotalCount)
	}

	x := filtering.NewQueryFilteredResult(
		data,
		filteredCount,
		totalCount,
		func(t *audit.AuditLogEntry) string {
			return t.ID
		},
		filter,
	)

	return x, nil
}

// CreateAuditLogEntry creates an audit log entry in a database.
func (q *repository) CreateAuditLogEntry(ctx context.Context, querier database.SQLQueryExecutor, input *audit.AuditLogEntryDatabaseCreationInput) (*audit.AuditLogEntry, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if input == nil {
		return nil, platformerrors.ErrNilInputProvided
	}

	tracing.AttachToSpan(span, identitykeys.AccountIDKey, input.BelongsToAccount)
	logger = logger.WithValue(identitykeys.AccountIDKey, input.BelongsToAccount)

	tracing.AttachToSpan(span, identitykeys.UserIDKey, input.BelongsToUser)
	logger = logger.WithValue(identitykeys.UserIDKey, input.BelongsToUser)

	marshaledChanges, err := json.Marshal(input.Changes)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "serializing audit log change list")
	}

	belongsToUser := strings.TrimSpace(input.BelongsToUser)
	var belongsToUserPtr *string
	if belongsToUser != "" {
		belongsToUserPtr = &belongsToUser
	}

	if err = q.generatedQuerier.CreateAuditLogEntry(ctx, querier, &generated.CreateAuditLogEntryParams{
		ID:               input.ID,
		ResourceType:     input.ResourceType,
		RelevantID:       input.RelevantID,
		EventType:        input.EventType,
		Changes:          string(marshaledChanges),
		BelongsToUser:    belongsToUserPtr,
		BelongsToAccount: input.BelongsToAccount,
	}); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "performing audit log creation query")
	}

	x := &audit.AuditLogEntry{
		ID:               input.ID,
		Changes:          input.Changes,
		BelongsToAccount: input.BelongsToAccount,
		CreatedAt:        q.CurrentTime(),
		ResourceType:     input.ResourceType,
		RelevantID:       input.RelevantID,
		EventType:        input.EventType,
		BelongsToUser:    input.BelongsToUser,
	}

	tracing.AttachToSpan(span, auditkeys.AuditLogEntryIDKey, x.ID)

	return x, nil
}
