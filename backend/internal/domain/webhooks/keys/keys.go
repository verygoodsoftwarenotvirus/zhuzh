package keys

const (
	idSuffix = ".id"

	// WebhookIDKey is the standard key for referring to a webhook's ID.
	WebhookIDKey = "webhook" + idSuffix
	// WebhookTriggerConfigIDKey is the standard key for referring to a webhook trigger config's ID.
	WebhookTriggerConfigIDKey = "webhook_trigger_config" + idSuffix
	// WebhookTriggerEventIDKey is the standard key for referring to a webhook trigger event's ID.
	WebhookTriggerEventIDKey = "webhook_trigger_event" + idSuffix
)
