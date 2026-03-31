package config

// Config holds payments service configuration.
type Config struct {
	Stripe     *StripeConfig     `env:"init" envPrefix:"STRIPE_"     json:"stripe"`
	RevenueCat *RevenueCatConfig `env:"init" envPrefix:"REVENUECAT_" json:"revenueCat"`
}

// StripeConfig holds Stripe-specific configuration.
type StripeConfig struct {
	APIKey        string `env:"API_KEY"        json:"apiKey"`
	WebhookSecret string `env:"WEBHOOK_SECRET" json:"webhookSecret"`
}

// RevenueCatConfig holds RevenueCat-specific configuration.
type RevenueCatConfig struct {
	APIKey            string `env:"API_KEY"             json:"apiKey"`
	WebhookAuthHeader string `env:"WEBHOOK_AUTH_HEADER" json:"webhookAuthHeader"`
}
