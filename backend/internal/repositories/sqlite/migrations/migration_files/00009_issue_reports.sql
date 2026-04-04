CREATE TABLE IF NOT EXISTS issue_reports (
    id TEXT NOT NULL PRIMARY KEY,
    issue_type TEXT NOT NULL,
    details TEXT NOT NULL,
    relevant_table TEXT,
    relevant_record_id TEXT,
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    last_updated_at TEXT,
    archived_at TEXT,
    created_by_user TEXT NOT NULL REFERENCES users("id") ON DELETE CASCADE,
    belongs_to_account TEXT NOT NULL REFERENCES accounts("id") ON DELETE CASCADE
);
