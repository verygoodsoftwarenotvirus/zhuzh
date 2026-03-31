package payments

import "context"

// ParsedWebhookEvent holds the result of parsing a provider webhook payload.
type ParsedWebhookEvent struct {
	EventType      string // e.g. "subscription.updated", "INITIAL_PURCHASE"
	AccountID      string // app_user_id, customer ID, etc.
	SubscriptionID string // external subscription or transaction ID
	ProductID      string // external product ID (e.g. StoreKit product_id for RevenueCat)
	Status         string // subscription status when known from payload (e.g. "active", "cancelled")
}

// PaymentProcessor defines the interface for payment provider webhook handling.
// Implementations (e.g., Stripe, RevenueCat) parse and verify webhooks; the manager writes to the database.
type PaymentProcessor interface {
	VerifyWebhookSignature(ctx context.Context, payload []byte, signature string, accountID string) bool
	ParseWebhookEvent(ctx context.Context, payload []byte) (*ParsedWebhookEvent, error)
}
