package authorization

const (
	CreateProductsPermission  Permission = "create.products"
	ReadProductsPermission    Permission = "read.products"
	UpdateProductsPermission  Permission = "update.products"
	ArchiveProductsPermission Permission = "archive.products"

	CreateCheckoutSessionPermission Permission = "create.checkout_sessions"

	CreateSubscriptionsPermission  Permission = "create.subscriptions"
	ReadSubscriptionsPermission    Permission = "read.subscriptions"
	UpdateSubscriptionsPermission  Permission = "update.subscriptions"
	ArchiveSubscriptionsPermission Permission = "archive.subscriptions"
	CancelSubscriptionPermission   Permission = "cancel.subscriptions"

	ReadPurchasesPermission      Permission = "read.purchases"
	ReadPaymentHistoryPermission Permission = "read.payment_history"
)

var (
	PaymentsPermissions = []Permission{
		CreateProductsPermission,
		ReadProductsPermission,
		UpdateProductsPermission,
		ArchiveProductsPermission,
		CreateCheckoutSessionPermission,
		CreateSubscriptionsPermission,
		ReadSubscriptionsPermission,
		UpdateSubscriptionsPermission,
		ArchiveSubscriptionsPermission,
		CancelSubscriptionPermission,
		ReadPurchasesPermission,
		ReadPaymentHistoryPermission,
	}
)
