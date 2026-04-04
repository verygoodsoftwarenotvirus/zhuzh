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
	resourceTypeProducts = "products"
)

func (r *repository) CreateProduct(ctx context.Context, input *payments.ProductDatabaseCreationInput) (*payments.Product, error) {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, platformerrors.ErrNilInputProvided
	}

	logger := r.logger.Clone()
	logger = logger.WithValue(paymentskeys.ProductIDKey, input.ID)
	tracing.AttachToSpan(span, paymentskeys.ProductIDKey, input.ID)

	arg := &generated.CreateProductParams{
		ID:                    input.ID,
		Name:                  input.Name,
		Description:           input.Description,
		Kind:                  input.Kind,
		AmountCents:           int64(input.AmountCents),
		Currency:              input.Currency,
		BillingIntervalMonths: int64PtrFromInt32Ptr(input.BillingIntervalMonths),
		ExternalProductID:     stringPtrFromString(input.ExternalProductID),
	}

	if err := r.generatedQuerier.CreateProduct(ctx, r.writeDB, arg); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating product")
	}

	if _, err := r.auditLogEntryRepo.CreateAuditLogEntry(ctx, r.writeDB, &audit.AuditLogEntryDatabaseCreationInput{
		ID:           identifiers.New(),
		ResourceType: resourceTypeProducts,
		RelevantID:   input.ID,
		EventType:    audit.AuditLogEventTypeCreated,
	}); err != nil {
		return nil, observability.PrepareError(err, span, "creating audit log entry")
	}

	return r.GetProduct(ctx, input.ID)
}

func (r *repository) GetProduct(ctx context.Context, id string) (*payments.Product, error) {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	logger := r.logger.Clone()
	if id == "" {
		return nil, platformerrors.ErrInvalidIDProvided
	}
	logger = logger.WithValue(paymentskeys.ProductIDKey, id)
	tracing.AttachToSpan(span, paymentskeys.ProductIDKey, id)

	result, err := r.generatedQuerier.GetProduct(ctx, r.readDB, id)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching product")
	}

	return convertProductFromGenerated(result), nil
}

func (r *repository) GetProductByExternalID(ctx context.Context, externalProductID string) (*payments.Product, error) {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	logger := r.logger.Clone()
	if externalProductID == "" {
		return nil, platformerrors.ErrInvalidIDProvided
	}
	logger = logger.WithValue("external_product_id", externalProductID)
	tracing.AttachToSpan(span, "external_product_id", externalProductID)

	result, err := r.generatedQuerier.GetProductByExternalID(ctx, r.readDB, stringPtrFromString(externalProductID))
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching product by external ID")
	}

	return convertProductFromGenerated(result), nil
}

func (r *repository) GetProducts(ctx context.Context, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[payments.Product], error) {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	logger := r.logger.Clone()
	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	params := &generated.GetProductsParams{
		CreatedAfter:  timePtrToStringPtr(filter.CreatedAfter),
		CreatedBefore: timePtrToStringPtr(filter.CreatedBefore),
		UpdatedBefore: timePtrToStringPtr(filter.UpdatedBefore),
		UpdatedAfter:  timePtrToStringPtr(filter.UpdatedAfter),
		Cursor:        filter.Cursor,
		ResultLimit:   int64PtrFromUint8Ptr(filter.MaxResponseSize),
	}

	results, err := r.generatedQuerier.GetProducts(ctx, r.readDB, params)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching products")
	}

	return convertProductsResult(results, filter), nil
}

func (r *repository) UpdateProduct(ctx context.Context, product *payments.Product) error {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	if product == nil {
		return platformerrors.ErrNilInputProvided
	}

	logger := r.logger.Clone()
	logger = logger.WithValue(paymentskeys.ProductIDKey, product.ID)
	tracing.AttachToSpan(span, paymentskeys.ProductIDKey, product.ID)

	arg := &generated.UpdateProductParams{
		ID:                    product.ID,
		Name:                  product.Name,
		Description:           product.Description,
		Kind:                  product.Kind,
		AmountCents:           int64(product.AmountCents),
		Currency:              product.Currency,
		BillingIntervalMonths: int64PtrFromInt32Ptr(product.BillingIntervalMonths),
		ExternalProductID:     stringPtrFromString(product.ExternalProductID),
	}

	_, err := r.generatedQuerier.UpdateProduct(ctx, r.writeDB, arg)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating product")
	}

	if _, err = r.auditLogEntryRepo.CreateAuditLogEntry(ctx, r.writeDB, &audit.AuditLogEntryDatabaseCreationInput{
		ID:           identifiers.New(),
		ResourceType: resourceTypeProducts,
		RelevantID:   product.ID,
		EventType:    audit.AuditLogEventTypeUpdated,
	}); err != nil {
		return observability.PrepareError(err, span, "creating audit log entry")
	}

	return nil
}

func (r *repository) ArchiveProduct(ctx context.Context, id string) error {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	logger := r.logger.Clone()
	if id == "" {
		return platformerrors.ErrInvalidIDProvided
	}
	logger = logger.WithValue(paymentskeys.ProductIDKey, id)
	tracing.AttachToSpan(span, paymentskeys.ProductIDKey, id)

	_, err := r.generatedQuerier.ArchiveProduct(ctx, r.writeDB, id)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving product")
	}

	if _, err = r.auditLogEntryRepo.CreateAuditLogEntry(ctx, r.writeDB, &audit.AuditLogEntryDatabaseCreationInput{
		ID:           identifiers.New(),
		ResourceType: resourceTypeProducts,
		RelevantID:   id,
		EventType:    audit.AuditLogEventTypeArchived,
	}); err != nil {
		return observability.PrepareError(err, span, "creating audit log entry")
	}

	return nil
}

func (r *repository) ProductExists(ctx context.Context, id string) (bool, error) {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	if id == "" {
		return false, platformerrors.ErrInvalidIDProvided
	}

	result, err := r.generatedQuerier.CheckProductExistence(ctx, r.readDB, id)
	if err != nil {
		return false, observability.PrepareAndLogError(err, r.logger.Clone(), span, "checking product existence")
	}
	return result != "", nil
}

func convertProductFromGenerated(m *generated.Products) *payments.Product {
	if m == nil {
		return nil
	}
	return &payments.Product{
		ID:                    m.ID,
		Name:                  m.Name,
		Description:           m.Description,
		Kind:                  m.Kind,
		AmountCents:           int32(m.AmountCents),
		Currency:              m.Currency,
		BillingIntervalMonths: int32PtrFromInt64Ptr(m.BillingIntervalMonths),
		ExternalProductID:     stringFromStringPtr(m.ExternalProductID),
		CreatedAt:             parseTime(m.CreatedAt),
		LastUpdatedAt:         parseTimePtr(m.LastUpdatedAt),
		ArchivedAt:            parseTimePtr(m.ArchivedAt),
	}
}

func convertProductFromRow(row *generated.GetProductsRow) *payments.Product {
	if row == nil {
		return nil
	}
	return &payments.Product{
		ID:                    row.ID,
		Name:                  row.Name,
		Description:           row.Description,
		Kind:                  row.Kind,
		AmountCents:           int32(row.AmountCents),
		Currency:              row.Currency,
		BillingIntervalMonths: int32PtrFromInt64Ptr(row.BillingIntervalMonths),
		ExternalProductID:     stringFromStringPtr(row.ExternalProductID),
		CreatedAt:             parseTime(row.CreatedAt),
		LastUpdatedAt:         parseTimePtr(row.LastUpdatedAt),
		ArchivedAt:            parseTimePtr(row.ArchivedAt),
	}
}

func convertProductsResult(rows []*generated.GetProductsRow, filter *filtering.QueryFilter) *filtering.QueryFilteredResult[payments.Product] {
	data := make([]*payments.Product, 0, len(rows))
	var filteredCount, totalCount uint64
	for _, row := range rows {
		data = append(data, convertProductFromRow(row))
		filteredCount = uint64(row.FilteredCount)
		totalCount = uint64(row.TotalCount)
	}
	return filtering.NewQueryFilteredResult(
		data,
		filteredCount,
		totalCount,
		func(p *payments.Product) string { return p.ID },
		filter,
	)
}
