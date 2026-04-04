package comments

import (
	"context"
	"database/sql"
	"time"

	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/audit"
	types "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/comments"
	commentskeys "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/comments/keys"
	generated "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/repositories/sqlite/comments/generated"

	"github.com/verygoodsoftwarenotvirus/platform/v4/database/filtering"
	platformerrors "github.com/verygoodsoftwarenotvirus/platform/v4/errors"
	"github.com/verygoodsoftwarenotvirus/platform/v4/identifiers"
	"github.com/verygoodsoftwarenotvirus/platform/v4/observability"
	"github.com/verygoodsoftwarenotvirus/platform/v4/observability/tracing"
)

const (
	resourceTypeComments = "comments"
)

var (
	_ types.CommentDataManager = (*repository)(nil)
)

// parseTime parses an RFC3339Nano timestamp string into a time.Time.
func parseTime(s string) time.Time {
	t, err := time.Parse(time.RFC3339Nano, s)
	if err != nil {
		return time.Time{}
	}
	return t
}

// parseTimePtr parses an optional RFC3339Nano timestamp string pointer into a *time.Time.
func parseTimePtr(s *string) *time.Time {
	if s == nil {
		return nil
	}
	t, err := time.Parse(time.RFC3339Nano, *s)
	if err != nil {
		return nil
	}
	return &t
}

// stringPointerFromTimePointer converts a *time.Time to a *string in RFC3339Nano format.
func stringPointerFromTimePointer(t *time.Time) *string {
	if t == nil {
		return nil
	}
	s := t.Format(time.RFC3339Nano)
	return &s
}

// interfaceFromUint8Pointer converts a *uint8 to an interface{} for SQLite COALESCE parameters.
func interfaceFromUint8Pointer(i *uint8) any {
	if i == nil {
		return nil
	}
	return int32(*i)
}

func convertCommentFromGenerated(c *generated.Comments) *types.Comment {
	return &types.Comment{
		ID:              c.ID,
		Content:         c.Content,
		TargetType:      c.TargetType,
		ReferencedID:    c.ReferencedID,
		ParentCommentID: c.ParentCommentID,
		BelongsToUser:   c.BelongsToUser,
		CreatedAt:       parseTime(c.CreatedAt),
		LastUpdatedAt:   parseTimePtr(c.LastUpdatedAt),
		ArchivedAt:      parseTimePtr(c.ArchivedAt),
	}
}

func convertRowToComment(r *generated.GetCommentsForReferenceRow) *types.Comment {
	return &types.Comment{
		ID:              r.ID,
		Content:         r.Content,
		TargetType:      r.TargetType,
		ReferencedID:    r.ReferencedID,
		ParentCommentID: r.ParentCommentID,
		BelongsToUser:   r.BelongsToUser,
		CreatedAt:       parseTime(r.CreatedAt),
		LastUpdatedAt:   parseTimePtr(r.LastUpdatedAt),
		ArchivedAt:      parseTimePtr(r.ArchivedAt),
	}
}

// CreateComment creates a comment in the database.
func (q *repository) CreateComment(ctx context.Context, input *types.CommentDatabaseCreationInput) (*types.Comment, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, platformerrors.ErrNilInputProvided
	}
	tracing.AttachToSpan(span, "comment_id", input.ID)
	logger := q.logger.WithValue("comment_id", input.ID)

	tx, err := q.writeDB.BeginTx(ctx, nil)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "beginning transaction")
	}

	if err = q.generatedQuerier.CreateComment(ctx, tx, &generated.CreateCommentParams{
		ID:              input.ID,
		Content:         input.Content,
		TargetType:      input.TargetType,
		ReferencedID:    input.ReferencedID,
		ParentCommentID: input.ParentCommentID,
		BelongsToUser:   input.BelongsToUser,
	}); err != nil {
		q.RollbackTransaction(ctx, tx)
		return nil, observability.PrepareAndLogError(err, logger, span, "performing comment creation query")
	}

	x := &types.Comment{
		ID:              input.ID,
		Content:         input.Content,
		TargetType:      input.TargetType,
		ReferencedID:    input.ReferencedID,
		ParentCommentID: input.ParentCommentID,
		BelongsToUser:   input.BelongsToUser,
		CreatedAt:       q.CurrentTime(),
	}
	tracing.AttachToSpan(span, "comment_id", x.ID)
	logger.Info("comment created")

	if _, err = q.auditLogEntryRepo.CreateAuditLogEntry(ctx, tx, &audit.AuditLogEntryDatabaseCreationInput{
		ID:            identifiers.New(),
		ResourceType:  resourceTypeComments,
		RelevantID:    x.ID,
		EventType:     audit.AuditLogEventTypeCreated,
		BelongsToUser: x.BelongsToUser,
	}); err != nil {
		q.RollbackTransaction(ctx, tx)
		return nil, observability.PrepareError(err, span, "creating audit log entry")
	}

	if err = tx.Commit(); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "committing transaction")
	}

	return x, nil
}

// GetComment fetches a comment from the database.
func (q *repository) GetComment(ctx context.Context, id string) (*types.Comment, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if id == "" {
		return nil, platformerrors.ErrInvalidIDProvided
	}
	logger = logger.WithValue(commentskeys.CommentIDKey, id)
	tracing.AttachToSpan(span, commentskeys.CommentIDKey, id)

	result, err := q.generatedQuerier.GetComment(ctx, q.readDB, id)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching comment")
	}

	return convertCommentFromGenerated(result), nil
}

// GetCommentsForReference fetches comments for a reference (including replies).
func (q *repository) GetCommentsForReference(ctx context.Context, targetType, referencedID string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[types.Comment], error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if targetType == "" || referencedID == "" {
		return nil, platformerrors.ErrInvalidIDProvided
	}
	logger = logger.WithValue("target_type", targetType).WithValue("referenced_id", referencedID)
	tracing.AttachToSpan(span, "target_type", targetType)
	tracing.AttachToSpan(span, "referenced_id", referencedID)

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	results, err := q.generatedQuerier.GetCommentsForReference(ctx, q.readDB, &generated.GetCommentsForReferenceParams{
		CreatedAfter:  stringPointerFromTimePointer(filter.CreatedAfter),
		CreatedBefore: stringPointerFromTimePointer(filter.CreatedBefore),
		UpdatedBefore: stringPointerFromTimePointer(filter.UpdatedBefore),
		UpdatedAfter:  stringPointerFromTimePointer(filter.UpdatedAfter),
		TargetType:    targetType,
		ReferencedID:  referencedID,
		Cursor:        filter.Cursor,
		ResultLimit:   interfaceFromUint8Pointer(filter.MaxResponseSize),
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing comments list retrieval query")
	}

	var (
		data                      = []*types.Comment{}
		filteredCount, totalCount uint64
	)
	for _, result := range results {
		data = append(data, convertRowToComment(result))
		filteredCount = uint64(result.FilteredCount)
		totalCount = uint64(result.TotalCount)
	}

	return filtering.NewQueryFilteredResult(
		data,
		filteredCount,
		totalCount,
		func(t *types.Comment) string { return t.ID },
		filter,
	), nil
}

// UpdateComment updates a comment in the database.
func (q *repository) UpdateComment(ctx context.Context, id, belongsToUser, content string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if id == "" || belongsToUser == "" {
		return platformerrors.ErrInvalidIDProvided
	}
	logger := q.logger.WithValue(commentskeys.CommentIDKey, id)
	tracing.AttachToSpan(span, commentskeys.CommentIDKey, id)

	tx, err := q.writeDB.BeginTx(ctx, nil)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "beginning transaction")
	}

	rowsAffected, err := q.generatedQuerier.UpdateComment(ctx, tx, &generated.UpdateCommentParams{
		Content:       content,
		ID:            id,
		BelongsToUser: belongsToUser,
	})
	if err != nil {
		q.RollbackTransaction(ctx, tx)
		return observability.PrepareAndLogError(err, logger, span, "updating comment")
	}
	if rowsAffected == 0 {
		q.RollbackTransaction(ctx, tx)
		return sql.ErrNoRows
	}

	if _, err = q.auditLogEntryRepo.CreateAuditLogEntry(ctx, tx, &audit.AuditLogEntryDatabaseCreationInput{
		ID:            identifiers.New(),
		ResourceType:  resourceTypeComments,
		RelevantID:    id,
		EventType:     audit.AuditLogEventTypeUpdated,
		BelongsToUser: belongsToUser,
	}); err != nil {
		q.RollbackTransaction(ctx, tx)
		return observability.PrepareError(err, span, "creating audit log entry")
	}

	if err = tx.Commit(); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "committing transaction")
	}

	return nil
}

// ArchiveComment archives a comment.
func (q *repository) ArchiveComment(ctx context.Context, id string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if id == "" {
		return platformerrors.ErrInvalidIDProvided
	}
	logger := q.logger.WithValue(commentskeys.CommentIDKey, id)
	tracing.AttachToSpan(span, commentskeys.CommentIDKey, id)

	comment, err := q.GetComment(ctx, id)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "fetching comment for archive")
	}

	tx, err := q.writeDB.BeginTx(ctx, nil)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "beginning transaction")
	}

	rowsAffected, err := q.generatedQuerier.ArchiveComment(ctx, tx, id)
	if err != nil {
		q.RollbackTransaction(ctx, tx)
		return observability.PrepareAndLogError(err, logger, span, "archiving comment")
	}
	if rowsAffected == 0 {
		q.RollbackTransaction(ctx, tx)
		return sql.ErrNoRows
	}

	if _, err = q.auditLogEntryRepo.CreateAuditLogEntry(ctx, tx, &audit.AuditLogEntryDatabaseCreationInput{
		ID:            identifiers.New(),
		ResourceType:  resourceTypeComments,
		RelevantID:    id,
		EventType:     audit.AuditLogEventTypeArchived,
		BelongsToUser: comment.BelongsToUser,
	}); err != nil {
		q.RollbackTransaction(ctx, tx)
		return observability.PrepareError(err, span, "creating audit log entry")
	}

	if err = tx.Commit(); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "committing transaction")
	}

	return nil
}

// ArchiveCommentsForReference archives all comments for a reference (including replies).
func (q *repository) ArchiveCommentsForReference(ctx context.Context, targetType, referencedID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if targetType == "" || referencedID == "" {
		return platformerrors.ErrInvalidIDProvided
	}
	logger := q.logger.WithValue("target_type", targetType).WithValue("referenced_id", referencedID)
	tracing.AttachToSpan(span, "target_type", targetType)
	tracing.AttachToSpan(span, "referenced_id", referencedID)

	filter := filtering.DefaultQueryFilter()
	maxSize := uint8(filtering.MaxQueryFilterLimit)
	filter.MaxResponseSize = &maxSize
	commentsResult, err := q.GetCommentsForReference(ctx, targetType, referencedID, filter)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "fetching comments for archive")
	}

	tx, err := q.writeDB.BeginTx(ctx, nil)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "beginning transaction")
	}

	_, err = q.generatedQuerier.ArchiveCommentsForReference(ctx, tx, &generated.ArchiveCommentsForReferenceParams{
		TargetType:   targetType,
		ReferencedID: referencedID,
	})
	if err != nil {
		q.RollbackTransaction(ctx, tx)
		return observability.PrepareAndLogError(err, logger, span, "archiving comments for reference")
	}

	for _, c := range commentsResult.Data {
		if _, err = q.auditLogEntryRepo.CreateAuditLogEntry(ctx, tx, &audit.AuditLogEntryDatabaseCreationInput{
			ID:            identifiers.New(),
			ResourceType:  resourceTypeComments,
			RelevantID:    c.ID,
			EventType:     audit.AuditLogEventTypeArchived,
			BelongsToUser: c.BelongsToUser,
		}); err != nil {
			q.RollbackTransaction(ctx, tx)
			return observability.PrepareError(err, span, "creating audit log entry")
		}
	}

	if err = tx.Commit(); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "committing transaction")
	}

	return nil
}
