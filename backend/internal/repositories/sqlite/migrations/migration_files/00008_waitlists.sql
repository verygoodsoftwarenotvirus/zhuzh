CREATE TABLE IF NOT EXISTS waitlists (
    id TEXT NOT NULL PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    valid_until TEXT NOT NULL,
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    last_updated_at TEXT,
    archived_at TEXT
);

CREATE TABLE IF NOT EXISTS waitlist_signups (
    id TEXT NOT NULL PRIMARY KEY,
    notes TEXT NOT NULL DEFAULT '',
    belongs_to_waitlist TEXT REFERENCES waitlists("id") ON DELETE CASCADE,
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    last_updated_at TEXT,
    archived_at TEXT,
    belongs_to_user TEXT REFERENCES users("id") ON DELETE CASCADE,
    belongs_to_account TEXT REFERENCES accounts("id") ON DELETE CASCADE
);
