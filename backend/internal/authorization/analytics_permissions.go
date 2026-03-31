package authorization

const (
	// ReportAnalyticsEventsPermission allows reporting analytics events via the proxy (TrackEvent).
	ReportAnalyticsEventsPermission Permission = "report.analytics_events"
)

var (
	// AnalyticsPermissions contains all analytics-related permissions.
	AnalyticsPermissions = []Permission{
		ReportAnalyticsEventsPermission,
	}
)
