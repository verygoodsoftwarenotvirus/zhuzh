package authorization

const (
	// CreateWebhooksPermission is an account admin permission.
	CreateWebhooksPermission Permission = "create.webhooks"
	// ReadWebhooksPermission is an account admin permission.
	ReadWebhooksPermission Permission = "read.webhooks"
	// UpdateWebhooksPermission is an account admin permission.
	UpdateWebhooksPermission Permission = "update.webhooks"
	// ArchiveWebhooksPermission is an account admin permission.
	ArchiveWebhooksPermission Permission = "archive.webhooks"
	// CreateWebhookTriggerConfigsPermission is an account permission for adding a trigger config to a webhook.
	CreateWebhookTriggerConfigsPermission Permission = "create.webhook_trigger_configs"
	// ArchiveWebhookTriggerConfigsPermission is an account permission for archiving a webhook trigger config.
	ArchiveWebhookTriggerConfigsPermission Permission = "archive.webhook_trigger_configs"
	// CreateWebhookTriggerEventsPermission is a permission for creating a catalog trigger event.
	CreateWebhookTriggerEventsPermission Permission = "create.webhook_trigger_events"
	// ReadWebhookTriggerEventsPermission is a permission for reading catalog trigger events.
	ReadWebhookTriggerEventsPermission Permission = "read.webhook_trigger_events"
	// UpdateWebhookTriggerEventsPermission is a permission for updating a catalog trigger event.
	UpdateWebhookTriggerEventsPermission Permission = "update.webhook_trigger_events"
	// ArchiveWebhookTriggerEventsPermission is a permission for archiving a catalog trigger event.
	ArchiveWebhookTriggerEventsPermission Permission = "archive.webhook_trigger_events"
)

var (
	// WebhooksPermissions contains all webhook-related permissions.
	WebhooksPermissions = []Permission{
		CreateWebhooksPermission,
		ReadWebhooksPermission,
		UpdateWebhooksPermission,
		ArchiveWebhooksPermission,
		CreateWebhookTriggerConfigsPermission,
		ArchiveWebhookTriggerConfigsPermission,
		CreateWebhookTriggerEventsPermission,
		ReadWebhookTriggerEventsPermission,
		UpdateWebhookTriggerEventsPermission,
		ArchiveWebhookTriggerEventsPermission,
	}
)
