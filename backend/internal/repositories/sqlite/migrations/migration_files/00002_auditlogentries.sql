-- Audit Log Entries Domain Migration
-- Audit trail functionality

CREATE TABLE IF NOT EXISTS audit_log_entries (
    id TEXT NOT NULL PRIMARY KEY,
    resource_type TEXT NOT NULL,
    relevant_id TEXT NOT NULL DEFAULT '',
    event_type TEXT NOT NULL DEFAULT 'other' CHECK(event_type IN ('other', 'created', 'updated', 'archived')),
    changes TEXT NOT NULL,
    belongs_to_account TEXT REFERENCES accounts("id") ON DELETE CASCADE,
    belongs_to_user TEXT REFERENCES users("id") ON DELETE CASCADE,
    created_at TEXT DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')) NOT NULL
);

-- =============================================================================
-- INDEXES FOR AUDIT LOG TABLES
-- =============================================================================

-- Audit log entries indexes
CREATE INDEX idx_audit_log_user ON audit_log_entries (belongs_to_user);
CREATE INDEX idx_audit_log_account ON audit_log_entries (belongs_to_account);
CREATE INDEX idx_audit_log_resource_type ON audit_log_entries (resource_type);
CREATE INDEX idx_audit_log_event_type ON audit_log_entries (event_type);
CREATE INDEX idx_audit_log_relevant_id ON audit_log_entries (relevant_id);
CREATE INDEX idx_audit_log_user_created_at ON audit_log_entries (belongs_to_user, created_at);
CREATE INDEX idx_audit_log_account_created_at ON audit_log_entries (belongs_to_account, created_at);
CREATE INDEX idx_audit_log_resource_relevant ON audit_log_entries (resource_type, relevant_id);
