package authorization

const (
	// ReadAuditLogEntriesPermission is a service permission.
	ReadAuditLogEntriesPermission Permission = "read.audit_log_entries"
)

var (
	// AuditPermissions contains all audit-related permissions.
	AuditPermissions = []Permission{
		ReadAuditLogEntriesPermission,
	}
)
