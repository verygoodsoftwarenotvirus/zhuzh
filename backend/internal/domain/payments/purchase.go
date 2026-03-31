package payments

import (
	"time"
)

type (
	// Purchase represents a one-time purchase record.
	Purchase struct {
		_                     struct{}   `json:"-"`
		CreatedAt             time.Time  `json:"createdAt"`
		CompletedAt           *time.Time `json:"completedAt"`
		LastUpdatedAt         *time.Time `json:"lastUpdatedAt"`
		ArchivedAt            *time.Time `json:"archivedAt"`
		ID                    string     `json:"id"`
		BelongsToAccount      string     `json:"belongsToAccount"`
		ProductID             string     `json:"productId"`
		Currency              string     `json:"currency"`
		ExternalTransactionID string     `json:"externalTransactionId"`
		AmountCents           int32      `json:"amountCents"`
	}

	// PurchaseDatabaseCreationInput is used for creating a purchase in the database.
	PurchaseDatabaseCreationInput struct {
		_                     struct{}   `json:"-"`
		CompletedAt           *time.Time `json:"-"`
		ID                    string     `json:"-"`
		BelongsToAccount      string     `json:"-"`
		ProductID             string     `json:"-"`
		Currency              string     `json:"-"`
		ExternalTransactionID string     `json:"-"`
		AmountCents           int32      `json:"-"`
	}
)
