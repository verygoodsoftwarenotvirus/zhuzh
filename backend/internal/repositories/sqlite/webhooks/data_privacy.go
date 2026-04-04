package webhooks

import (
	"context"

	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/domain/webhooks"
	"github.com/verygoodsoftwarenotvirus/zhuzh/backend/internal/repositories/sqlite/webhooks/generated"

	"github.com/verygoodsoftwarenotvirus/platform/v4/database/filtering"
	"github.com/verygoodsoftwarenotvirus/platform/v4/observability"
)

func (r *repository) CollectUserData(ctx context.Context, accountIDs []string) (*webhooks.UserDataCollection, error) {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	logger := r.logger.WithSpan(span)
	filter := filtering.DefaultQueryFilter()

	x := &webhooks.UserDataCollection{
		Data: make(map[string][]webhooks.Webhook),
	}

	for _, accountID := range accountIDs {
		accountWebhooks, err := r.generatedQuerier.GetWebhooksForAccount(ctx, r.readDB, &generated.GetWebhooksForAccountParams{
			CreatedBefore:    timePtrToStringPtr(filter.CreatedBefore),
			CreatedAfter:     timePtrToStringPtr(filter.CreatedAfter),
			UpdatedBefore:    timePtrToStringPtr(filter.UpdatedBefore),
			UpdatedAfter:     timePtrToStringPtr(filter.UpdatedAfter),
			Cursor:           filter.Cursor,
			ResultLimit:      int64PtrFromUint8Ptr(filter.MaxResponseSize),
			BelongsToAccount: accountID,
		})
		if err != nil {
			return nil, observability.PrepareAndLogError(err, logger, span, "retrieving webhooks for account")
		}

		seen := make(map[string]struct{})
		for _, row := range accountWebhooks {
			if _, ok := seen[row.ID]; ok {
				continue
			}
			seen[row.ID] = struct{}{}
			x.Data[accountID] = append(x.Data[accountID], webhooks.Webhook{
				CreatedAt:        parseTime(row.CreatedAt_2),
				ArchivedAt:       parseTimePtr(row.ArchivedAt_2),
				LastUpdatedAt:    parseTimePtr(row.LastUpdatedAt),
				Name:             row.Name,
				URL:              row.URL,
				Method:           row.Method,
				ID:               row.ID,
				BelongsToAccount: row.BelongsToAccount,
				CreatedByUser:    row.CreatedByUser,
				ContentType:      row.ContentType,
				TriggerConfigs:   nil,
			})
		}
	}

	return x, nil
}
