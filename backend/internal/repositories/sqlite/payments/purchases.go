package payments

import (
	"context"

	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/audit"
	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/payments"
	paymentskeys "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/payments/keys"
	generated "github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/repositories/sqlite/payments/generated"

	"github.com/verygoodsoftwarenotvirus/platform/v4/database/filtering"
	platformerrors "github.com/verygoodsoftwarenotvirus/platform/v4/errors"
	"github.com/verygoodsoftwarenotvirus/platform/v4/identifiers"
	"github.com/verygoodsoftwarenotvirus/platform/v4/observability"
	"github.com/verygoodsoftwarenotvirus/platform/v4/observability/tracing"
)

const (
	resourceTypePurchases = "purchases"
)

func (r *repository) CreatePurchase(ctx context.Context, input *payments.PurchaseDatabaseCreationInput) (*payments.Purchase, error) {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, platformerrors.ErrNilInputProvided
	}

	logger := r.logger.Clone()
	logger = logger.WithValue(paymentskeys.PurchaseIDKey, input.ID)
	tracing.AttachToSpan(span, paymentskeys.PurchaseIDKey, input.ID)

	arg := &generated.CreatePurchaseParams{
		ID:                    input.ID,
		BelongsToAccount:      input.BelongsToAccount,
		ProductID:             input.ProductID,
		AmountCents:           int64(input.AmountCents),
		Currency:              input.Currency,
		CompletedAt:           timePtrToStringPtr(input.CompletedAt),
		ExternalTransactionID: stringPtrFromString(input.ExternalTransactionID),
	}

	if err := r.generatedQuerier.CreatePurchase(ctx, r.writeDB, arg); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating purchase")
	}

	if _, err := r.auditLogEntryRepo.CreateAuditLogEntry(ctx, r.writeDB, &audit.AuditLogEntryDatabaseCreationInput{
		BelongsToAccount: &input.BelongsToAccount,
		ID:               identifiers.New(),
		ResourceType:     resourceTypePurchases,
		RelevantID:       input.ID,
		EventType:        audit.AuditLogEventTypeCreated,
	}); err != nil {
		return nil, observability.PrepareError(err, span, "creating audit log entry")
	}

	return r.GetPurchase(ctx, input.ID)
}

func (r *repository) GetPurchase(ctx context.Context, id string) (*payments.Purchase, error) {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	logger := r.logger.Clone()
	if id == "" {
		return nil, platformerrors.ErrInvalidIDProvided
	}
	logger = logger.WithValue(paymentskeys.PurchaseIDKey, id)
	tracing.AttachToSpan(span, paymentskeys.PurchaseIDKey, id)

	result, err := r.generatedQuerier.GetPurchase(ctx, r.readDB, id)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching purchase")
	}

	return convertPurchaseFromGenerated(result), nil
}

func (r *repository) GetPurchasesForAccount(ctx context.Context, accountID string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[payments.Purchase], error) {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	if accountID == "" {
		return nil, platformerrors.ErrInvalidIDProvided
	}

	logger := r.logger.Clone()
	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	params := &generated.GetPurchasesForAccountParams{
		BelongsToAccount: accountID,
		CreatedAfter:     timePtrToStringPtr(filter.CreatedAfter),
		CreatedBefore:    timePtrToStringPtr(filter.CreatedBefore),
		UpdatedBefore:    timePtrToStringPtr(filter.UpdatedBefore),
		UpdatedAfter:     timePtrToStringPtr(filter.UpdatedAfter),
		Cursor:           filter.Cursor,
		ResultLimit:      int64PtrFromUint8Ptr(filter.MaxResponseSize),
	}

	results, err := r.generatedQuerier.GetPurchasesForAccount(ctx, r.readDB, params)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching purchases for account")
	}

	return convertPurchasesResult(results, filter), nil
}

func convertPurchaseFromGenerated(m *generated.Purchases) *payments.Purchase {
	if m == nil {
		return nil
	}
	return &payments.Purchase{
		ID:                    m.ID,
		BelongsToAccount:      m.BelongsToAccount,
		ProductID:             m.ProductID,
		AmountCents:           int32(m.AmountCents),
		Currency:              m.Currency,
		CompletedAt:           parseTimePtr(m.CompletedAt),
		ExternalTransactionID: stringFromStringPtr(m.ExternalTransactionID),
		CreatedAt:             parseTime(m.CreatedAt),
		LastUpdatedAt:         parseTimePtr(m.LastUpdatedAt),
		ArchivedAt:            parseTimePtr(m.ArchivedAt),
	}
}

func convertPurchasesResult(rows []*generated.GetPurchasesForAccountRow, filter *filtering.QueryFilter) *filtering.QueryFilteredResult[payments.Purchase] {
	data := make([]*payments.Purchase, 0, len(rows))
	var filteredCount, totalCount uint64
	for _, row := range rows {
		data = append(data, &payments.Purchase{
			ID:                    row.ID,
			BelongsToAccount:      row.BelongsToAccount,
			ProductID:             row.ProductID,
			AmountCents:           int32(row.AmountCents),
			Currency:              row.Currency,
			CompletedAt:           parseTimePtr(row.CompletedAt),
			ExternalTransactionID: stringFromStringPtr(row.ExternalTransactionID),
			CreatedAt:             parseTime(row.CreatedAt),
			LastUpdatedAt:         parseTimePtr(row.LastUpdatedAt),
			ArchivedAt:            parseTimePtr(row.ArchivedAt),
		})
		filteredCount = uint64(row.FilteredCount)
		totalCount = uint64(row.TotalCount)
	}
	return filtering.NewQueryFilteredResult(
		data,
		filteredCount,
		totalCount,
		func(p *payments.Purchase) string { return p.ID },
		filter,
	)
}
