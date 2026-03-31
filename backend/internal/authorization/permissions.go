package authorization

type (
	role int

	// Permission is a simple string alias.
	Permission string
)

var (
	// ServiceAdminPermissions is every service admin permission.
	ServiceAdminPermissions = []Permission{
		ReadUserDataPermission,
		UpdateUserStatusPermission,
		ReadUserPermission,
		SearchUserPermission,
		ArchiveUserPermission,
		CreateOAuth2ClientsPermission,
		ArchiveOAuth2ClientsPermission,
		ArchiveServiceSettingsPermission,
		CreateUserNotificationsPermission,
		ImpersonateUserPermission,
		ManageUserSessionsPermission,
		PublishArbitraryQueueMessagePermission,
		// only admins can arbitrarily create these via the API, this is exclusively for integration test purposes.
		CreateServiceSettingsPermission,
		CreateWaitlistsPermission,
		UpdateWaitlistsPermission,
		ArchiveWaitlistsPermission,
		CreateProductsPermission,
		ReadProductsPermission,
		UpdateProductsPermission,
		ArchiveProductsPermission,
		CreateSubscriptionsPermission,
		ReadSubscriptionsPermission,
		UpdateSubscriptionsPermission,
		ArchiveSubscriptionsPermission,
	}

	// ServiceDataAdminPermissions is every service data admin permission.
	ServiceDataAdminPermissions = []Permission{}

	// AccountAdminPermissions is every account admin permission.
	AccountAdminPermissions = []Permission{
		UpdateAccountPermission,
		ArchiveAccountPermission,
		TransferAccountPermission,
		InviteUserToAccountPermission,
		ModifyMemberPermissionsForAccountPermission,
		RemoveMemberAccountPermission,
		CreateWebhooksPermission,
		UpdateWebhooksPermission,
		ArchiveWebhooksPermission,
		CreateIssueReportsPermission,
		UpdateIssueReportsPermission,
		ArchiveIssueReportsPermission,
		CreateWebhookTriggerConfigsPermission,
		ArchiveWebhookTriggerConfigsPermission,
		CreateWebhookTriggerEventsPermission,
		ReadWebhookTriggerEventsPermission,
		UpdateWebhookTriggerEventsPermission,
		ArchiveWebhookTriggerEventsPermission,
		CreateCheckoutSessionPermission,
		CancelSubscriptionPermission,
		ReadPurchasesPermission,
		ReadPaymentHistoryPermission,
		ReadSubscriptionsPermission,
	}

	// AccountMemberPermissions is every account member permission.
	AccountMemberPermissions = []Permission{
		ReportAnalyticsEventsPermission,
		ReadWebhooksPermission,
		ReadIssueReportsPermission,
		ReadAuditLogEntriesPermission,
		ReadOAuth2ClientsPermission,
		ReadServiceSettingsPermission,
		SearchServiceSettingsPermission,
		CreateUploadedMediaPermission,
		ReadUploadedMediaPermission,
		UpdateUploadedMediaPermission,
		ArchiveUploadedMediaPermission,
		CreateServiceSettingConfigurationsPermission,
		ReadServiceSettingConfigurationsPermission,
		UpdateServiceSettingConfigurationsPermission,
		ArchiveServiceSettingConfigurationsPermission,
		CreateCommentsPermission,
		ReadCommentsPermission,
		UpdateCommentsPermission,
		ArchiveCommentsPermission,
		ReadUserNotificationsPermission,
		UpdateUserNotificationsPermission,
		CreateUserDeviceTokensPermission,
		ReadUserDeviceTokensPermission,
		ArchiveUserDeviceTokensPermission,
		CreateWaitlistSignupsPermission,
		UpdateWaitlistSignupsPermission,
		ArchiveWaitlistSignupsPermission,
		ReadWaitlistsPermission,
		ReadWaitlistSignupsPermission,
		CreateCheckoutSessionPermission,
		CancelSubscriptionPermission,
		ReadPurchasesPermission,
		ReadPaymentHistoryPermission,
		ReadSubscriptionsPermission,
	}
)
