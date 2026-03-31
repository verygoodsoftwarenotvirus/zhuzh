package payments

import (
	"time"
)

const (
	PaymentTransactionStatusSucceeded = "succeeded"
	PaymentTransactionStatusFailed    = "failed"
	PaymentTransactionStatusPending   = "pending"
	PaymentTransactionStatusRefunded  = "refunded"
)

type (
	// PaymentTransaction represents an audit record of a payment attempt.
	PaymentTransaction struct {
		_                     struct{}  `json:"-"`
		CreatedAt             time.Time `json:"createdAt"`
		SubscriptionID        *string   `json:"subscriptionId"`
		PurchaseID            *string   `json:"purchaseId"`
		ID                    string    `json:"id"`
		BelongsToAccount      string    `json:"belongsToAccount"`
		ExternalTransactionID string    `json:"externalTransactionId"`
		Currency              string    `json:"currency"`
		Status                string    `json:"status"`
		AmountCents           int32     `json:"amountCents"`
	}

	// PaymentTransactionDatabaseCreationInput is used for creating a payment transaction in the database.
	PaymentTransactionDatabaseCreationInput struct {
		_                     struct{} `json:"-"`
		SubscriptionID        *string  `json:"-"`
		PurchaseID            *string  `json:"-"`
		ID                    string   `json:"-"`
		BelongsToAccount      string   `json:"-"`
		ExternalTransactionID string   `json:"-"`
		Currency              string   `json:"-"`
		Status                string   `json:"-"`
		AmountCents           int32    `json:"-"`
	}
)
